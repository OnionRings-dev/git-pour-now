package selfupdate_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/OnionRings-dev/git-pour-now/internal/selfupdate"
)

func mockGo(t *testing.T) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("shell fixtures require Unix")
	}
	dir := t.TempDir()
	script := `#!/bin/sh
set -eu
case "$*" in
  'list -m -json github.com/OnionRings-dev/git-pour-now@latest')
    [ "${SELFUPDATE_FAIL_LOOKUP:-}" != 1 ] || { echo 'network unavailable' >&2; exit 1; }
    printf '%s\n' "$SELFUPDATE_MODULE_JSON"
    ;;
  'install github.com/OnionRings-dev/git-pour-now@v1.2.0')
    [ "${SELFUPDATE_FAIL_BUILD:-}" != 1 ] || exit 1
    [ "$GOOS" = "$SELFUPDATE_NATIVE_OS" ]
    [ "$GOARCH" = "$SELFUPDATE_NATIVE_ARCH" ]
    printf '#!/bin/sh\necho updated\n' > "$GOBIN/git-pour-now"
    ;;
  *) echo "unexpected Go command: $*" >&2; exit 2 ;;
esac
`
	if err := os.WriteFile(filepath.Join(dir, "go"), []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("SELFUPDATE_MODULE_JSON", fmt.Sprintf(`{"Path":%q,"Version":"v1.2.0"}`, selfupdate.Module))
	t.Setenv("SELFUPDATE_NATIVE_OS", runtime.GOOS)
	t.Setenv("SELFUPDATE_NATIVE_ARCH", runtime.GOARCH)
	t.Setenv("SELFUPDATE_FAIL_LOOKUP", "")
	t.Setenv("SELFUPDATE_FAIL_BUILD", "")
	return dir
}

func check(t *testing.T, current string) selfupdate.Status {
	t.Helper()
	status, err := selfupdate.Check(context.Background(), current)
	if err != nil {
		t.Fatal(err)
	}
	return status
}

func installedBinary(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "git-pour")
	if err := os.WriteFile(path, []byte("original binary"), 0755); err != nil {
		t.Fatal(err)
	}
	return path
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// TestCheckNewVersion tests detection of an available update.
// Steps:
// 1. Mocks the module resolver returning v1.2.0.
// 2. Checks a v1.1.0 installation and verifies an update is available.
func TestCheckNewVersion(t *testing.T) {
	mockGo(t)
	status := check(t, "v1.1.0")
	if !status.Known || !status.Available || status.Latest != "v1.2.0" {
		t.Fatalf("unexpected status: %+v", status)
	}
}

// TestCheckSameVersion tests that build metadata does not trigger an update.
// Steps:
// 1. Mocks latest as v1.2.0 and checks an installation with build metadata.
// 2. Verifies SemVer precedence treats it as already current.
func TestCheckSameVersion(t *testing.T) {
	mockGo(t)
	status := check(t, "1.2.0+local")
	if !status.Known || status.Available {
		t.Fatalf("unexpected status: %+v", status)
	}
}

// TestCheckNewerInstalledVersion tests downgrade prevention.
// Steps:
// 1. Checks v2.0.0 against latest v1.2.0.
// 2. Verifies no update is offered.
func TestCheckNewerInstalledVersion(t *testing.T) {
	mockGo(t)
	status := check(t, "2.0.0")
	if !status.Known || status.Available {
		t.Fatalf("unexpected status: %+v", status)
	}
}

// TestCheckDevelopmentVersion tests honest reporting for unversioned builds.
// Steps:
// 1. Checks an installation reporting dev.
// 2. Verifies latest is resolved but no comparison or update claim is made.
func TestCheckDevelopmentVersion(t *testing.T) {
	mockGo(t)
	status := check(t, "dev")
	if status.Known || status.Available || status.Latest != "v1.2.0" {
		t.Fatalf("unexpected status: %+v", status)
	}
}

// TestCheckInvalidModule tests rejection of an unexpected resolver response.
// Steps:
// 1. Makes the resolver return a different module path.
// 2. Verifies the lookup fails instead of offering an update.
func TestCheckInvalidModule(t *testing.T) {
	mockGo(t)
	t.Setenv("SELFUPDATE_MODULE_JSON", `{"Path":"example.com/other","Version":"v1.2.0"}`)
	if _, err := selfupdate.Check(context.Background(), "1.0.0"); err == nil {
		t.Fatal("unexpected module was accepted")
	}
}

// TestCheckNetworkFailure tests lookup error reporting.
// Steps:
// 1. Simulates an unavailable module resolver.
// 2. Verifies the command reports the underlying network failure.
func TestCheckNetworkFailure(t *testing.T) {
	mockGo(t)
	t.Setenv("SELFUPDATE_FAIL_LOOKUP", "1")
	if _, err := selfupdate.Check(context.Background(), "1.0.0"); err == nil || !strings.Contains(err.Error(), "network unavailable") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestInstallReplacesSymlinkTarget tests replacement through a symbolic link.
// Steps:
// 1. Creates an existing binary and a symlink to it.
// 2. Builds the pinned update despite cross-compilation environment settings.
// 3. Verifies the target changes, permissions survive, and the symlink remains.
func TestInstallReplacesSymlinkTarget(t *testing.T) {
	mockGo(t)
	t.Setenv("GOOS", "windows")
	t.Setenv("GOARCH", "386")
	target := installedBinary(t)
	link := filepath.Join(t.TempDir(), "git-pour-link")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	if err := selfupdate.Install(context.Background(), "v1.2.0", link, io.Discard); err != nil {
		t.Fatal(err)
	}
	if readFile(t, target) != "#!/bin/sh\necho updated\n" {
		t.Fatal("target was not updated")
	}
	if resolved, err := os.Readlink(link); err != nil || resolved != target {
		t.Fatalf("symlink was changed: %q, %v", resolved, err)
	}
	info, err := os.Stat(target)
	if err != nil || info.Mode().Perm() != 0755 {
		t.Fatalf("executable permissions changed: %v", err)
	}
	entries, err := os.ReadDir(filepath.Dir(target))
	if err != nil || len(entries) != 1 {
		t.Fatalf("update left staged files behind: %v", err)
	}
}

// TestInstallBuildFailurePreservesBinary tests recovery from a failed build.
// Steps:
// 1. Simulates a failing go install with an existing executable.
// 2. Verifies the executable is unchanged and staged files are cleaned up.
func TestInstallBuildFailurePreservesBinary(t *testing.T) {
	mockGo(t)
	t.Setenv("SELFUPDATE_FAIL_BUILD", "1")
	target := installedBinary(t)
	if err := selfupdate.Install(context.Background(), "v1.2.0", target, io.Discard); err == nil {
		t.Fatal("expected build failure")
	}
	if readFile(t, target) != "original binary" {
		t.Fatal("failed build replaced the executable")
	}
	entries, err := os.ReadDir(filepath.Dir(target))
	if err != nil || len(entries) != 1 {
		t.Fatalf("failed build left staged files behind: %v", err)
	}
}

// TestInstallCanceledPreservesBinary tests cancellation before replacement.
// Steps:
// 1. Cancels the update context before starting the build.
// 2. Verifies the update fails and leaves the existing executable unchanged.
func TestInstallCanceledPreservesBinary(t *testing.T) {
	mockGo(t)
	target := installedBinary(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := selfupdate.Install(ctx, "v1.2.0", target, io.Discard); err == nil {
		t.Fatal("canceled update succeeded")
	}
	if readFile(t, target) != "original binary" {
		t.Fatal("canceled update replaced the executable")
	}
}
