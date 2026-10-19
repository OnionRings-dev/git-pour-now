// Package selfupdate checks and installs the same Go module as scripts/install.sh.
package selfupdate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/OnionRings-dev/git-pour-now/internal/semver"
)

const Module = "github.com/OnionRings-dev/git-pour-now"

type Status struct {
	Current   string
	Latest    string
	Known     bool
	Available bool
}

// Check queries Go's latest module version without changing the executable.
func Check(ctx context.Context, current string) (Status, error) {
	status := Status{Current: current}
	dir, err := os.MkdirTemp("", "git-pour-check-*")
	if err != nil {
		return status, err
	}
	defer os.RemoveAll(dir)
	command := exec.CommandContext(ctx, "go", "list", "-m", "-json", Module+"@latest")
	command.Dir = dir
	command.Env = append(os.Environ(), "GOWORK=off")
	var stderr strings.Builder
	command.Stderr = &stderr
	output, err := command.Output()
	if err != nil {
		return status, fmt.Errorf("check latest Go module (Go 1.23.6 or newer is required): %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	var module struct{ Path, Version string }
	if err := json.Unmarshal(output, &module); err != nil {
		return status, fmt.Errorf("decode latest Go module: %w", err)
	}
	latest := strings.TrimPrefix(module.Version, "v")
	if module.Path != Module || !strings.HasPrefix(module.Version, "v") || !semver.Valid(latest) {
		return status, fmt.Errorf("invalid latest Go module: path %q, version %q", module.Path, module.Version)
	}
	status.Latest = module.Version
	normalized := strings.TrimPrefix(current, "v")
	status.Known = semver.Valid(normalized)
	status.Available = status.Known && semver.Compare(normalized, latest) < 0
	return status, nil
}

// Install builds a pinned module version, then atomically replaces executable.
// A symlink is preserved by replacing its target. Windows cannot rename a
// running executable this way, so it must use an external installer instead.
func Install(ctx context.Context, version, executable string, output io.Writer) error {
	if runtime.GOOS != "darwin" && runtime.GOOS != "linux" {
		return fmt.Errorf("self-update installation supports macOS and Linux; reinstall git-pour manually on %s", runtime.GOOS)
	}
	if !strings.HasPrefix(version, "v") || !semver.Valid(strings.TrimPrefix(version, "v")) {
		return fmt.Errorf("invalid update version %q", version)
	}
	target, err := filepath.EvalSymlinks(executable)
	if err != nil {
		return fmt.Errorf("resolve executable: %w", err)
	}
	info, err := os.Stat(target)
	if err != nil {
		return err
	}
	// Stage beside the executable so the final rename stays on one filesystem.
	staged, err := os.CreateTemp(filepath.Dir(target), ".git-pour-update-*")
	if err != nil {
		return fmt.Errorf("cannot write to installation directory; reinstall in a writable directory: %w", err)
	}
	defer os.Remove(staged.Name())
	defer staged.Close()
	dir, err := os.MkdirTemp("", "git-pour-build-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	command := exec.CommandContext(ctx, "go", "install", Module+"@"+version)
	command.Dir = dir
	command.Env = append(os.Environ(), "GOBIN="+dir, "GOWORK=off", "GOOS="+runtime.GOOS, "GOARCH="+runtime.GOARCH)
	command.Stdout, command.Stderr = output, output
	if err := command.Run(); err != nil {
		return fmt.Errorf("build update; existing executable was not changed: %w", err)
	}
	source, err := os.Open(filepath.Join(dir, "git-pour-now"))
	if err != nil {
		return fmt.Errorf("open built update: %w", err)
	}
	defer source.Close()
	if _, err := io.Copy(staged, source); err != nil {
		return fmt.Errorf("stage update: %w", err)
	}
	if err := staged.Chmod(info.Mode().Perm()); err != nil {
		return err
	}
	if err := staged.Sync(); err != nil {
		return err
	}
	if err := staged.Close(); err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := os.Rename(staged.Name(), target); err != nil {
		return fmt.Errorf("replace executable; existing executable was not changed: %w", err)
	}
	return nil
}
