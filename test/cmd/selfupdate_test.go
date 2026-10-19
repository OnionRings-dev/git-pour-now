package cmd_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/OnionRings-dev/git-pour-now/test/testutil"
)

func selfUpdateFixture(t *testing.T) (string, []string, string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("self-update installation requires Unix")
	}
	dir := t.TempDir()
	binary := filepath.Join(dir, "git-pour")
	data, err := os.ReadFile(testutil.GitFlowPath())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(binary, data, 0755); err != nil {
		t.Fatal(err)
	}
	tools := filepath.Join(dir, "tools")
	if err := os.Mkdir(tools, 0755); err != nil {
		t.Fatal(err)
	}
	log := filepath.Join(dir, "go-commands")
	script := `#!/bin/sh
set -eu
printf '%s\n' "$*" >> "$SELFUPDATE_LOG"
case "$*" in
  'list -m -json github.com/OnionRings-dev/git-pour-now@latest')
    echo '{"Path":"github.com/OnionRings-dev/git-pour-now","Version":"v1.2.0"}'
    ;;
  'install github.com/OnionRings-dev/git-pour-now@v1.2.0')
    [ "${SELFUPDATE_FAIL_BUILD:-}" != 1 ] || exit 1
    printf '#!/bin/sh\necho updated-fixture\n' > "$GOBIN/git-pour-now"
    ;;
  *) echo "unexpected Go command: $*" >&2; exit 2 ;;
esac
`
	if err := os.WriteFile(filepath.Join(tools, "go"), []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
	env := append(os.Environ(), "PATH="+tools+string(os.PathListSeparator)+os.Getenv("PATH"), "SELFUPDATE_LOG="+log, "SELFUPDATE_FAIL_BUILD=")
	return binary, env, log
}

// TestSelfUpdateCheckOutsideRepository tests a read-only check outside Git.
// Steps:
// 1. Copies the CLI to an isolated installation and mocks module lookup.
// 2. Runs self-update --check in a directory without a Git repository.
// 3. Verifies latest is reported, no build runs, and the executable is unchanged.
func TestSelfUpdateCheckOutsideRepository(t *testing.T) {
	binary, env, log := selfUpdateFixture(t)
	before, err := os.ReadFile(binary)
	if err != nil {
		t.Fatal(err)
	}
	command := exec.Command(binary, "self-update", "--check")
	command.Dir = t.TempDir()
	command.Env = env
	output, err := command.CombinedOutput()
	if err != nil || !strings.Contains(string(output), "Latest: v1.2.0") || !strings.Contains(string(output), "Installed:") {
		t.Fatalf("check failed: %v\n%s", err, output)
	}
	after, err := os.ReadFile(binary)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatalf("check changed executable: %v", err)
	}
	commands, err := os.ReadFile(log)
	if err != nil || strings.Contains(string(commands), "install ") {
		t.Fatalf("check unexpectedly built an update: %v\n%s", err, commands)
	}
}

// TestSelfUpdateReplacesRunningExecutable tests installation from the CLI itself.
// Steps:
// 1. Copies the CLI to isolate the replacement from the test runner's binary.
// 2. Runs self-update with a mocked Go build of the pinned latest version.
// 3. Verifies success and that subsequent invocations execute the replacement.
func TestSelfUpdateReplacesRunningExecutable(t *testing.T) {
	binary, env, _ := selfUpdateFixture(t)
	command := exec.Command(binary, "self-update")
	command.Dir = t.TempDir()
	command.Env = env
	output, err := command.CombinedOutput()
	if err != nil || !strings.Contains(string(output), "Updated git-pour to v1.2.0.") {
		t.Fatalf("update failed: %v\n%s", err, output)
	}
	output, err = exec.Command(binary, "version").CombinedOutput()
	if err != nil || string(output) != "updated-fixture\n" {
		t.Fatalf("replacement did not run: %v\n%s", err, output)
	}
}

// TestSelfUpdateBuildFailure tests command failure without corrupting the installation.
// Steps:
// 1. Copies the CLI and configures the mocked build to fail.
// 2. Runs self-update and verifies a nonzero exit status.
// 3. Verifies the binary remains byte-for-byte identical to the original.
func TestSelfUpdateBuildFailure(t *testing.T) {
	binary, env, _ := selfUpdateFixture(t)
	before, err := os.ReadFile(binary)
	if err != nil {
		t.Fatal(err)
	}
	command := exec.Command(binary, "self-update")
	command.Dir = t.TempDir()
	command.Env = append(env, "SELFUPDATE_FAIL_BUILD=1")
	output, err := command.CombinedOutput()
	if err == nil || !strings.Contains(string(output), "existing executable was not changed") {
		t.Fatalf("missing build error: %v\n%s", err, output)
	}
	after, err := os.ReadFile(binary)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatalf("failed update changed executable: %v", err)
	}
}

// TestSelfUpdateDoesNotActivateSharedConfig tests independence from repository setup.
// Steps:
// 1. Creates an uninitialized Git repository with automatic shared config activation.
// 2. Runs self-update -c and verifies the lookup succeeds.
// 3. Verifies the repository's shared configuration was not activated.
func TestSelfUpdateDoesNotActivateSharedConfig(t *testing.T) {
	binary, env, _ := selfUpdateFixture(t)
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)
	if err := os.WriteFile(filepath.Join(dir, ".gitpour"), []byte("[gitpour]\nversion = 1.0\n[gitpour.shared]\nautoInit = true\n"), 0644); err != nil {
		t.Fatal(err)
	}
	command := exec.Command(binary, "self-update", "-c")
	command.Dir = dir
	command.Env = env
	if output, err := command.CombinedOutput(); err != nil || strings.Contains(string(output), "auto-initialized") {
		t.Fatalf("check triggered repository setup: %v\n%s", err, output)
	}
	if output, err := testutil.RunGit(t, dir, "config", "--local", "--get", "gitpour.version"); err == nil || output != "" {
		t.Fatalf("shared configuration was activated: %v, %s", err, output)
	}
}
