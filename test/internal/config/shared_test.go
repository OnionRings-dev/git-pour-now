package config_test

import (
	"strings"
	"testing"

	"github.com/OnionRings-dev/git-pour-now/internal/config"
	"github.com/OnionRings-dev/git-pour-now/test/testutil"
)

// setSharedConfig writes a key into the repository's .gitpour file via
// `git config --file`.
func setSharedConfig(t *testing.T, dir, key, value string) {
	t.Helper()
	if _, err := testutil.RunGit(t, dir, "config", "--file", config.SharedConfigFileName, key, value); err != nil {
		t.Fatalf("Failed to set shared config %s: %v", key, err)
	}
}

// localConfigValue returns a single value from local git config, and whether the
// key is present.
func localConfigValue(t *testing.T, dir, key string) (string, bool) {
	t.Helper()
	out, err := testutil.RunGit(t, dir, "config", "--local", "--get", key)
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(out), true
}

// TestSharedManagedKeyPredicate covers scenario 38: the shared-managed key
// predicate classifies keys correctly (pure function, table-driven per the
// TESTING_GUIDELINES exception).
// Steps:
// 1. Enumerates representative keys with their expected managed classification
// 2. Calls config.IsSharedManagedKey on each
// 3. Verifies control keys and runtime .base keys are excluded, real config included
func TestSharedManagedKeyPredicate(t *testing.T) {
	t.Parallel()
	cases := []struct {
		key     string
		managed bool
	}{
		{"gitpour.branch.feature.prefix", true},
		{"gitpour.feature.publish.push-option", true},
		{"gitpour.version", true},
		{"gitpour.initialized", true},
		{"gitpour.shared.autoInit", false},
		{"gitpour.shared.trustHooks", false},
		{"gitpour.branch.feature/foo.base", false},
		{"core.sshCommand", false},
	}
	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			if got := config.IsSharedManagedKey(tc.key); got != tc.managed {
				t.Errorf("IsSharedManagedKey(%q) = %v, want %v", tc.key, got, tc.managed)
			}
		})
	}
}

// TestIsSharedManagedKeyExcludesWorktreeMarker verifies the worktree provenance
// marker is machine-local state and never participates in the shared-managed set,
// while gitpour.worktreePath stays shared-managed.
// Steps:
// 1. Calls config.IsSharedManagedKey on gitpour.worktree.feature/x.managed
// 2. Repeats with a mixed-case spelling of the same key
// 3. Verifies both are excluded (false)
// 4. Verifies gitpour.worktreePath and gitpour.branch.feature.prefix remain included (true)
func TestIsSharedManagedKeyExcludesWorktreeMarker(t *testing.T) {
	t.Parallel()

	if config.IsSharedManagedKey("gitpour.worktree.feature/x.managed") {
		t.Error("expected gitpour.worktree.feature/x.managed to be excluded from the shared-managed set")
	}
	if config.IsSharedManagedKey("gitpour.Worktree.feature/x.Managed") {
		t.Error("expected the mixed-case worktree marker to be excluded from the shared-managed set")
	}
	if !config.IsSharedManagedKey("gitpour.worktreePath") {
		t.Error("expected gitpour.worktreePath to stay shared-managed")
	}
	if !config.IsSharedManagedKey("gitpour.branch.feature.prefix") {
		t.Error("expected gitpour.branch.feature.prefix to stay shared-managed")
	}
}

// TestCopySharedToLocalStaleRemovalPreservesLocalOnly covers scenario 39:
// CopySharedToLocal removes stale managed keys but preserves local-only keys.
// Steps:
// 1. Seeds local with a stale managed key, a shared-control key, and a runtime base key
// 2. Writes a .gitpour with feature defaults (no "old" branch)
// 3. Calls config.CopySharedToLocal
// 4. Verifies old.* is gone, feature.* copied, and the local-only keys preserved
func TestCopySharedToLocalStaleRemovalPreservesLocalOnly(t *testing.T) {
	t.Parallel()
	dir := setupTestRepo(t)
	defer cleanupTestRepo(t, dir)

	setConfig(t, dir, "gitpour.branch.old.type", "topic")
	setConfig(t, dir, "gitpour.shared.autoInit", "true")
	setConfig(t, dir, "gitpour.branch.feature/foo.base", "develop")

	setSharedConfig(t, dir, "gitpour.version", "1.0")
	setSharedConfig(t, dir, "gitpour.branch.feature.type", "topic")
	setSharedConfig(t, dir, "gitpour.branch.feature.prefix", "feature/")

	if _, err := config.CopySharedToLocal(openRepo(t, dir)); err != nil {
		t.Fatalf("CopySharedToLocal failed: %v", err)
	}

	if _, ok := localConfigValue(t, dir, "gitpour.branch.old.type"); ok {
		t.Error("expected stale gitpour.branch.old.type to be removed from local")
	}
	if v, _ := localConfigValue(t, dir, "gitpour.branch.feature.prefix"); v != "feature/" {
		t.Errorf("expected local feature.prefix=feature/, got %q", v)
	}
	if v, _ := localConfigValue(t, dir, "gitpour.shared.autoInit"); v != "true" {
		t.Errorf("expected local gitpour.shared.autoInit preserved as true, got %q", v)
	}
	if v, _ := localConfigValue(t, dir, "gitpour.branch.feature/foo.base"); v != "develop" {
		t.Errorf("expected runtime gitpour.branch.feature/foo.base preserved as develop, got %q", v)
	}
}

// TestCopySharedToLocalIgnoresExcludedKeysInFile covers scenario 39b:
// CopySharedToLocal never copies excluded keys that appear inside .gitpour and
// never overwrites their local counterparts.
// Steps:
// 1. Seeds local counterparts for each excluded key with differing values
// 2. Writes a .gitpour that (wrongly) contains excluded keys plus feature defaults
// 3. Calls config.CopySharedToLocal
// 4. Verifies excluded keys are inert (not copied / not overwritten) and defaults copied
func TestCopySharedToLocalIgnoresExcludedKeysInFile(t *testing.T) {
	t.Parallel()
	dir := setupTestRepo(t)
	defer cleanupTestRepo(t, dir)

	setConfig(t, dir, "gitpour.shared.autoInit", "false")
	setConfig(t, dir, "gitpour.branch.feature/foo.base", "develop")

	setSharedConfig(t, dir, "gitpour.version", "1.0")
	setSharedConfig(t, dir, "gitpour.shared.autoInit", "true")
	setSharedConfig(t, dir, "gitpour.shared.trustHooks", "true")
	setSharedConfig(t, dir, "gitpour.branch.feature/foo.base", "main")
	setSharedConfig(t, dir, "gitpour.branch.feature.type", "topic")
	setSharedConfig(t, dir, "gitpour.branch.feature.prefix", "feature/")

	if _, err := config.CopySharedToLocal(openRepo(t, dir)); err != nil {
		t.Fatalf("CopySharedToLocal failed: %v", err)
	}

	if v, _ := localConfigValue(t, dir, "gitpour.shared.autoInit"); v != "false" {
		t.Errorf("expected local gitpour.shared.autoInit untouched (false), got %q", v)
	}
	if v, _ := localConfigValue(t, dir, "gitpour.branch.feature/foo.base"); v != "develop" {
		t.Errorf("expected local runtime base untouched (develop), got %q", v)
	}
	if _, ok := localConfigValue(t, dir, "gitpour.shared.trustHooks"); ok {
		t.Error("expected local gitpour.shared.trustHooks NOT created from .gitpour")
	}
	if v, _ := localConfigValue(t, dir, "gitpour.branch.feature.prefix"); v != "feature/" {
		t.Errorf("expected local feature.prefix copied as feature/, got %q", v)
	}
}

// TestIsUnconfiguredFreshRepo covers scenario 40a: a fresh repo with no gitpour.*
// reads as unconfigured.
// Steps:
// 1. Creates a fresh repository with no git-pour config
// 2. Calls config.IsUnconfiguredIgnoringShared
// 3. Verifies it reports unconfigured (true)
func TestIsUnconfiguredFreshRepo(t *testing.T) {
	t.Parallel()
	dir := setupTestRepo(t)
	defer cleanupTestRepo(t, dir)

	unconfigured, err := config.IsUnconfiguredIgnoringShared(openRepo(t, dir))
	if err != nil {
		t.Fatalf("IsUnconfiguredIgnoringShared failed: %v", err)
	}
	if !unconfigured {
		t.Error("expected fresh repo to read as unconfigured")
	}
}

// TestIsUnconfiguredSharedControlOnly covers scenario 40b: a repo carrying only
// gitpour.shared.* control keys still reads as unconfigured.
// Steps:
// 1. Sets only gitpour.shared.autoInit and gitpour.shared.trustHooks
// 2. Calls config.IsUnconfiguredIgnoringShared
// 3. Verifies it reports unconfigured (shared-control keys are ignored)
func TestIsUnconfiguredSharedControlOnly(t *testing.T) {
	t.Parallel()
	dir := setupTestRepo(t)
	defer cleanupTestRepo(t, dir)

	setConfig(t, dir, "gitpour.shared.autoInit", "true")
	setConfig(t, dir, "gitpour.shared.trustHooks", "true")

	unconfigured, err := config.IsUnconfiguredIgnoringShared(openRepo(t, dir))
	if err != nil {
		t.Fatalf("IsUnconfiguredIgnoringShared failed: %v", err)
	}
	if !unconfigured {
		t.Error("expected shared-control-only repo to read as unconfigured")
	}
}

// TestIsUnconfiguredCommandOnlyKey covers scenario 40c: a repo with only a
// command-scoped gitpour.* key (no version, no branch definitions) still reads as
// unconfigured.
// Steps:
// 1. Sets only gitpour.feature.start.fetch=true
// 2. Calls config.IsUnconfiguredIgnoringShared
// 3. Verifies it reports unconfigured (a stray command key is not branch config)
func TestIsUnconfiguredCommandOnlyKey(t *testing.T) {
	t.Parallel()
	dir := setupTestRepo(t)
	defer cleanupTestRepo(t, dir)

	setConfig(t, dir, "gitpour.feature.start.fetch", "true")

	unconfigured, err := config.IsUnconfiguredIgnoringShared(openRepo(t, dir))
	if err != nil {
		t.Fatalf("IsUnconfiguredIgnoringShared failed: %v", err)
	}
	if !unconfigured {
		t.Error("expected command-only-key repo to read as unconfigured")
	}
}

// TestIsUnconfiguredAfterGitFlowInit covers scenario 40d: a repo initialized by
// git-pour-now reads as configured.
// Steps:
// 1. Initializes git-pour with defaults via the binary
// 2. Calls config.IsUnconfiguredIgnoringShared
// 3. Verifies it reports configured (false)
func TestIsUnconfiguredAfterGitFlowInit(t *testing.T) {
	t.Parallel()
	dir := setupTestRepo(t)
	defer cleanupTestRepo(t, dir)

	if out, err := testutil.RunGitFlow(t, dir, "init", "--defaults", "--preset=classic"); err != nil {
		t.Fatalf("git pour init failed: %v\n%s", err, out)
	}

	unconfigured, err := config.IsUnconfiguredIgnoringShared(openRepo(t, dir))
	if err != nil {
		t.Fatalf("IsUnconfiguredIgnoringShared failed: %v", err)
	}
	if unconfigured {
		t.Error("expected git-pour-initialized repo to read as configured")
	}
}

// TestIsUnconfiguredAvhOnlyNoVersion covers scenario 40e: AVH-only config with no
// gitpour.version reads as configured.
// Steps:
// 1. Sets git-flow-avh style keys (gitpour.branch.master, gitpour.prefix.feature) and no version
// 2. Calls config.IsUnconfiguredIgnoringShared
// 3. Verifies it reports configured (AVH branch config counts as configuration)
func TestIsUnconfiguredAvhOnlyNoVersion(t *testing.T) {
	t.Parallel()
	dir := setupTestRepo(t)
	defer cleanupTestRepo(t, dir)

	setConfig(t, dir, "gitflow.branch.master", "master")
	setConfig(t, dir, "gitflow.branch.develop", "develop")
	setConfig(t, dir, "gitflow.prefix.feature", "feature/")

	unconfigured, err := config.IsUnconfiguredIgnoringShared(openRepo(t, dir))
	if err != nil {
		t.Fatalf("IsUnconfiguredIgnoringShared failed: %v", err)
	}
	if unconfigured {
		t.Error("expected AVH-only repo to read as configured")
	}
}

// TestIsUnconfiguredBranchDefsNoVersion covers scenario 40f: valid
// git-pour-now branch definitions without gitpour.version read as configured.
// Steps:
// 1. Sets gitpour.branch.feature.type/parent/prefix but no gitpour.version
// 2. Calls config.IsUnconfiguredIgnoringShared
// 3. Verifies it reports configured (branch definitions alone are enough)
func TestIsUnconfiguredBranchDefsNoVersion(t *testing.T) {
	t.Parallel()
	dir := setupTestRepo(t)
	defer cleanupTestRepo(t, dir)

	setConfig(t, dir, "gitpour.branch.feature.type", "topic")
	setConfig(t, dir, "gitpour.branch.feature.parent", "develop")
	setConfig(t, dir, "gitpour.branch.feature.prefix", "feature/")

	unconfigured, err := config.IsUnconfiguredIgnoringShared(openRepo(t, dir))
	if err != nil {
		t.Fatalf("IsUnconfiguredIgnoringShared failed: %v", err)
	}
	if unconfigured {
		t.Error("expected repo with branch definitions to read as configured")
	}
}
