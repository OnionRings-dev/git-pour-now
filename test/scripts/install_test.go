package scripts_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type installerEnv struct {
	home   string
	script string
	env    []string
}

func setupInstaller(t *testing.T, shell string) installerEnv {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("installer supports macOS and Linux")
	}
	root := t.TempDir()
	home := filepath.Join(root, "home with spaces")
	bin := filepath.Join(root, "tools")
	for _, dir := range []string{home, bin} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}
	// Replace only the network build; installation and shell setup run for real.
	mockGo := `#!/bin/sh
set -eu
[ "$*" = 'install github.com/OnionRings-dev/git-pour-now@latest' ] || exit 2
[ "${INSTALLER_TEST_FAIL:-}" != 1 ] || exit 3
cat > "$GOBIN/git-pour-now" <<'BINARY'
#!/bin/sh
[ "$1" = version ] || exit 1
echo installer-test-version
BINARY
`
	if err := os.WriteFile(filepath.Join(bin, "go"), []byte(mockGo), 0755); err != nil {
		t.Fatal(err)
	}
	script, err := filepath.Abs("../../scripts/install.sh")
	if err != nil {
		t.Fatal(err)
	}
	env := []string{}
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		switch key {
		case "HOME", "SHELL", "ZDOTDIR", "PATH", "GOBIN", "INSTALLER_TEST_FAIL":
			continue
		}
		env = append(env, entry)
	}
	return installerEnv{home, script, append(env,
		"HOME="+home, "SHELL="+shell,
		"PATH="+bin+string(os.PathListSeparator)+os.Getenv("PATH"),
		"GOBIN="+filepath.Join(root, "unused-gobin"),
		"INSTALLER_TEST_FAIL=",
	)}
}

func (e installerEnv) run(t *testing.T) {
	t.Helper()
	cmd := exec.Command("sh", e.script)
	cmd.Env = e.env
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("installer failed: %v\n%s", err, output)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

// TestInstallerZsh tests installation and command discovery in a new zsh.
// Steps:
// 1. Installs into a temporary home whose path contains spaces.
// 2. Starts zsh with the generated configuration when zsh is available.
// 3. Verifies both git-pour version and git pour version resolve successfully.
func TestInstallerZsh(t *testing.T) {
	e := setupInstaller(t, "/bin/zsh")
	e.run(t)
	if !strings.Contains(readFile(t, filepath.Join(e.home, ".zshrc")), `export PATH="$HOME/.local/bin:$PATH"`) {
		t.Fatal("missing PATH configuration")
	}
	if _, err := exec.LookPath("zsh"); err != nil {
		t.Log("zsh unavailable; verified generated profile")
		return
	}
	cmd := exec.Command("zsh", "-ic", "git-pour version && git pour version")
	cmd.Env = e.env
	output, err := cmd.CombinedOutput()
	if err != nil || strings.Count(string(output), "installer-test-version") != 2 {
		t.Fatalf("commands unavailable: %v\n%s", err, output)
	}
}

// TestInstallerRepeatPreservesProfile tests repeat installation without duplicate setup.
// Steps:
// 1. Creates an existing bash profile without a trailing newline.
// 2. Runs the installer twice.
// 3. Verifies existing configuration survives and PATH setup appears only once.
func TestInstallerRepeatPreservesProfile(t *testing.T) {
	e := setupInstaller(t, "/bin/bash")
	profile := filepath.Join(e.home, ".bash_profile")
	if err := os.WriteFile(profile, []byte("# existing configuration"), 0644); err != nil {
		t.Fatal(err)
	}
	e.run(t)
	e.run(t)
	content := readFile(t, profile)
	if !strings.HasPrefix(content, "# existing configuration\n") || strings.Count(content, "export PATH=") != 1 {
		t.Fatalf("profile damaged or duplicated: %s", content)
	}
	if strings.Count(readFile(t, filepath.Join(e.home, ".bashrc")), "export PATH=") != 1 {
		t.Fatal("duplicate bashrc PATH setup")
	}
	cmd := exec.Command("bash", "-c", `. "$HOME/.bash_profile"; git pour version`)
	cmd.Env = e.env
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("bash PATH setup failed: %v\n%s", err, output)
	}
}

// TestInstallerZdotdir tests zsh configuration outside the home directory.
// Steps:
// 1. Sets ZDOTDIR to a nonexistent directory containing spaces.
// 2. Runs the installer and verifies it creates the zsh configuration there.
// 3. Verifies it does not create a home-directory zshrc.
func TestInstallerZdotdir(t *testing.T) {
	e := setupInstaller(t, "/bin/zsh")
	zdotdir := filepath.Join(e.home, "custom zsh")
	e.env = append(e.env, "ZDOTDIR="+zdotdir)
	e.run(t)
	readFile(t, filepath.Join(zdotdir, ".zshrc"))
	if _, err := os.Stat(filepath.Join(e.home, ".zshrc")); !os.IsNotExist(err) {
		t.Fatalf("unexpected home zshrc: %v", err)
	}
}

// TestInstallerBuildFailurePreservesInstall tests failure before replacing the binary.
// Steps:
// 1. Creates an existing executable and simulates a failed Go installation.
// 2. Verifies the installer returns an error.
// 3. Verifies the executable survives and shell configuration is not created.
func TestInstallerBuildFailurePreservesInstall(t *testing.T) {
	e := setupInstaller(t, "/bin/zsh")
	binary := filepath.Join(e.home, ".local", "bin", "git-pour")
	if err := os.MkdirAll(filepath.Dir(binary), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(binary, []byte("existing binary"), 0755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("sh", e.script)
	cmd.Env = append(e.env, "INSTALLER_TEST_FAIL=1")
	if output, err := cmd.CombinedOutput(); err == nil || !strings.Contains(string(output), "Go installation failed") {
		t.Fatalf("missing build failure: %v\n%s", err, output)
	}
	if readFile(t, binary) != "existing binary" {
		t.Fatal("existing executable was changed")
	}
	if _, err := os.Stat(filepath.Join(e.home, ".zshrc")); !os.IsNotExist(err) {
		t.Fatalf("unexpected profile after failure: %v", err)
	}
}

// TestInstallerWithoutShell tests the default profile when SHELL is empty.
// Steps:
// 1. Runs the installer without a configured login shell.
// 2. Verifies it creates the POSIX login profile successfully.
func TestInstallerWithoutShell(t *testing.T) {
	e := setupInstaller(t, "")
	e.run(t)
	readFile(t, filepath.Join(e.home, ".profile"))
}

func copiedInstallCommand(t *testing.T) string {
	t.Helper()
	for _, line := range strings.Split(readFile(t, "../../README.md"), "\n") {
		if strings.HasPrefix(line, "curl -fsSL ") {
			return line
		}
	}
	t.Fatal("README installation command missing")
	return ""
}

func mockCurl(t *testing.T, e installerEnv, body string) {
	t.Helper()
	for _, entry := range e.env {
		if strings.HasPrefix(entry, "PATH=") {
			bin := strings.SplitN(strings.TrimPrefix(entry, "PATH="), string(os.PathListSeparator), 2)[0]
			if err := os.WriteFile(filepath.Join(bin, "curl"), []byte("#!/bin/sh\n"+body), 0755); err != nil {
				t.Fatal(err)
			}
			return
		}
	}
	t.Fatal("test PATH missing")
}

// TestInstallerCopiedCommand tests command discovery after installing with the README command.
// Steps:
// 1. Serves the local installer through a mock curl command.
// 2. Runs the exact README command in bash.
// 3. Verifies git pour works in a new shell using the generated configuration.
func TestInstallerCopiedCommand(t *testing.T) {
	e := setupInstaller(t, "/bin/bash")
	mockCurl(t, e, `cat "$INSTALLER_SCRIPT"`)
	cmd := exec.Command("bash", "-c", copiedInstallCommand(t))
	cmd.Env = append(e.env, "INSTALLER_SCRIPT="+e.script)
	output, err := cmd.CombinedOutput()
	if err != nil || !strings.Contains(string(output), "installer-test-version") {
		t.Fatalf("copied command failed: %v\n%s", err, output)
	}
	cmd = exec.Command("bash", "--noprofile", "--rcfile", filepath.Join(e.home, ".bashrc"), "-ic", "git pour version")
	cmd.Env = e.env
	if output, err := cmd.CombinedOutput(); err != nil || !strings.Contains(string(output), "installer-test-version") {
		t.Fatalf("command unavailable in new shell: %v\n%s", err, output)
	}
}
