package cmd_test

import (
	"github.com/OnionRings-dev/git-pour-now/test/testutil"
	"strings"
	"testing"
)

func initPour(t *testing.T) string {
	t.Helper()
	dir := testutil.SetupTestRepo(t)
	t.Cleanup(func() { testutil.CleanupTestRepo(t, dir) })
	if out, err := testutil.RunGitFlow(t, dir, "init", "--defaults"); err != nil {
		t.Fatalf("init: %s: %v", out, err)
	}
	return dir
}

// TestPourDefaults verifies that fresh repositories use main and feature only.
func TestPourDefaults(t *testing.T) {
	dir := initPour(t)
	if testutil.BranchExists(t, dir, "develop") {
		t.Fatal("develop must not be created")
	}
	out, err := testutil.RunGit(t, dir, "config", "gitpour.branch.feature.parent")
	if err != nil || strings.TrimSpace(out) != "main" {
		t.Fatalf("feature parent: %s: %v", out, err)
	}
}

// TestPourFeatureFinish verifies feature changes land on main without a release tag.
func TestPourFeatureFinish(t *testing.T) {
	dir := initPour(t)
	if out, err := testutil.RunGitFlow(t, dir, "feature", "start", "example"); err != nil {
		t.Fatalf("start: %s: %v", out, err)
	}
	testutil.WriteFile(t, dir, "feature.txt", "new feature")
	if out, err := testutil.RunGit(t, dir, "add", "feature.txt"); err != nil {
		t.Fatalf("add: %s: %v", out, err)
	}
	if out, err := testutil.RunGit(t, dir, "commit", "-m", "Add feature"); err != nil {
		t.Fatalf("commit: %s: %v", out, err)
	}
	if out, err := testutil.RunGitFlow(t, dir, "feature", "finish", "example"); err != nil {
		t.Fatalf("finish: %s: %v", out, err)
	}
	if out, err := testutil.RunGit(t, dir, "show", "main:feature.txt"); err != nil || out != "new feature" {
		t.Fatalf("main contents: %s: %v", out, err)
	}
	if out, err := testutil.RunGit(t, dir, "tag", "--list"); err != nil || strings.TrimSpace(out) != "" {
		t.Fatalf("unexpected tags: %s: %v", out, err)
	}
}

// TestPourRelease verifies an annotated SemVer tag is created on main.
func TestPourRelease(t *testing.T) {
	dir := initPour(t)
	if out, err := testutil.RunGitFlow(t, dir, "prerelease", "start", "1.0.0"); err != nil {
		t.Fatalf("prerelease start: %s: %v", out, err)
	}
	if out, err := testutil.RunGitFlow(t, dir, "release", "1.0.0-rc.1+build.001"); err != nil {
		t.Fatalf("release: %s: %v", out, err)
	}
	if out, err := testutil.RunGit(t, dir, "cat-file", "-t", "refs/tags/v1.0.0-rc.1+build.001"); err != nil || strings.TrimSpace(out) != "tag" {
		t.Fatalf("annotated tag: %s: %v", out, err)
	}
}

func TestPourReleaseRejectsInvalidVersion(t *testing.T) {
	dir := initPour(t)
	if out, err := testutil.RunGitFlow(t, dir, "release", "1.02.0"); err == nil {
		t.Fatalf("accepted invalid version: %s", out)
	}
}

func TestPourReleaseRejectsFeatureBranch(t *testing.T) {
	dir := initPour(t)
	if out, err := testutil.RunGitFlow(t, dir, "feature", "start", "example"); err != nil {
		t.Fatalf("start: %s: %v", out, err)
	}
	if out, err := testutil.RunGitFlow(t, dir, "release", "1.0.0"); err == nil {
		t.Fatalf("released feature: %s", out)
	}
}

func TestPourReleaseRejectsDirtyTree(t *testing.T) {
	dir := initPour(t)
	testutil.WriteFile(t, dir, "untracked.txt", "pending")
	if out, err := testutil.RunGitFlow(t, dir, "release", "1.0.0"); err == nil {
		t.Fatalf("released dirty tree: %s", out)
	}
}

func TestPourReleaseRejectsExistingVersion(t *testing.T) {
	dir := initPour(t)
	if out, err := testutil.RunGitFlow(t, dir, "release", "1.0.0"); err != nil {
		t.Fatalf("first release: %s: %v", out, err)
	}
	if out, err := testutil.RunGitFlow(t, dir, "release", "1.0.0"); err == nil {
		t.Fatalf("reused release: %s", out)
	}
}

func TestPourReleaseRejectsLowerVersion(t *testing.T) {
	dir := initPour(t)
	if out, err := testutil.RunGitFlow(t, dir, "release", "2.0.0"); err != nil {
		t.Fatalf("first release: %s: %v", out, err)
	}
	if out, err := testutil.RunGitFlow(t, dir, "release", "1.0.0"); err == nil {
		t.Fatalf("accepted older release: %s", out)
	}
}

// TestPourStableReleaseRejectsPrereleaseBranch keeps stable releases on main.
func TestPourStableReleaseRejectsPrereleaseBranch(t *testing.T) {
	dir := initPour(t)
	if out, err := testutil.RunGitFlow(t, dir, "prerelease", "start", "1.0.0"); err != nil {
		t.Fatalf("start: %s: %v", out, err)
	}
	if out, err := testutil.RunGitFlow(t, dir, "release", "1.0.0"); err == nil {
		t.Fatalf("stable release from prerelease branch: %s", out)
	}
}

func TestPourPrereleaseRejectsMain(t *testing.T) {
	dir := initPour(t)
	if out, err := testutil.RunGitFlow(t, dir, "release", "1.0.0-rc.1"); err == nil {
		t.Fatalf("prerelease from main: %s", out)
	}
}

func TestPourPrereleaseRejectsWrongVersionBranch(t *testing.T) {
	dir := initPour(t)
	if out, err := testutil.RunGitFlow(t, dir, "prerelease", "start", "1.0.0"); err != nil {
		t.Fatalf("start: %s: %v", out, err)
	}
	if out, err := testutil.RunGitFlow(t, dir, "release", "2.0.0-rc.1"); err == nil {
		t.Fatalf("wrong version branch: %s", out)
	}
}

// TestPourPrereleasePromotion verifies preparation commits reach the stable tag.
func TestPourPrereleasePromotion(t *testing.T) {
	dir := initPour(t)
	if out, err := testutil.RunGitFlow(t, dir, "prerelease", "start", "1.0.0"); err != nil {
		t.Fatalf("start: %s: %v", out, err)
	}
	testutil.WriteFile(t, dir, "prepared.txt", "ready")
	if out, err := testutil.RunGit(t, dir, "add", "prepared.txt"); err != nil {
		t.Fatalf("add: %s: %v", out, err)
	}
	if out, err := testutil.RunGit(t, dir, "commit", "-m", "Prepare release"); err != nil {
		t.Fatalf("commit: %s: %v", out, err)
	}
	if out, err := testutil.RunGitFlow(t, dir, "release", "1.0.0-rc.1"); err != nil {
		t.Fatalf("prerelease: %s: %v", out, err)
	}
	if out, err := testutil.RunGitFlow(t, dir, "prerelease", "finish", "1.0.0"); err != nil {
		t.Fatalf("finish: %s: %v", out, err)
	}
	if out, err := testutil.RunGitFlow(t, dir, "release", "1.0.0"); err != nil {
		t.Fatalf("stable release: %s: %v", out, err)
	}
	if out, err := testutil.RunGit(t, dir, "show", "v1.0.0:prepared.txt"); err != nil || out != "ready" {
		t.Fatalf("stable contents: %s: %v", out, err)
	}
}

func TestPourPrereleaseRejectsNonCoreVersion(t *testing.T) {
	dir := initPour(t)
	if out, err := testutil.RunGitFlow(t, dir, "prerelease", "start", "1.0.0-rc.1"); err == nil {
		t.Fatalf("accepted prerelease branch with non-core version: %s", out)
	}
	if testutil.BranchExists(t, dir, "prerelease/1.0.0-rc.1") {
		t.Fatal("invalid branch was created")
	}
}

func TestPourReleaseUsesCustomTrunk(t *testing.T) {
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)
	if out, err := testutil.RunGitFlow(t, dir, "init", "--defaults", "--main", "trunk"); err != nil {
		t.Fatalf("init: %s: %v", out, err)
	}
	if out, err := testutil.RunGit(t, dir, "checkout", "trunk"); err != nil {
		t.Fatalf("checkout: %s: %v", out, err)
	}
	if out, err := testutil.RunGitFlow(t, dir, "release", "1.0.0"); err != nil {
		t.Fatalf("release: %s: %v", out, err)
	}
	if testutil.BranchExists(t, dir, "develop") {
		t.Fatal("custom trunk created develop")
	}
}

func TestPourReleaseRejectsHiddenUntrackedFiles(t *testing.T) {
	dir := initPour(t)
	if out, err := testutil.RunGit(t, dir, "config", "status.showUntrackedFiles", "no"); err != nil {
		t.Fatalf("configure status: %s: %v", out, err)
	}
	testutil.WriteFile(t, dir, "untracked.txt", "pending")
	if out, err := testutil.RunGitFlow(t, dir, "release", "1.0.0"); err == nil {
		t.Fatalf("released dirty tree: %s", out)
	}
	if out, err := testutil.RunGit(t, dir, "tag", "--list", "v1.0.0"); err != nil || strings.TrimSpace(out) != "" {
		t.Fatalf("unexpected release tag: %s: %v", out, err)
	}
}
