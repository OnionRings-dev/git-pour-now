package git_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OnionRings-dev/git-pour-now/test/testutil"
)

func TestRemoveLegacyConfigLocalPreservesIncludedSettings(t *testing.T) {
	t.Parallel()
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)
	included := filepath.Join(t.TempDir(), "included.config")
	const contents = "[gitflow]\n\torigin = inherited\n"
	if err := os.WriteFile(included, []byte(contents), 0600); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{
		{"config", "include.path", included},
		{"config", "gitflow.origin", "local"},
		{"config", "--add", "gitflow.feature.publish.pushoption", "first"},
		{"config", "--add", "gitflow.feature.publish.pushoption", "second"},
		{"config", "gitpour.version", "1.0"},
	} {
		if out, err := testutil.RunGit(t, dir, args...); err != nil {
			t.Fatalf("configure: %s: %v", out, err)
		}
	}
	repo := openRepo(t, dir)
	if err := repo.RemoveLegacyConfigLocal(); err != nil {
		t.Fatal(err)
	}
	if out, err := testutil.RunGit(t, dir, "config", "--local", "--no-includes", "--get-regexp", "^gitflow\\."); err == nil || strings.TrimSpace(out) != "" {
		t.Fatalf("local legacy settings remain: %s: %v", out, err)
	}
	if got, err := repo.GetConfig("gitflow.origin"); err != nil || got != "inherited" {
		t.Fatalf("included setting changed: %q: %v", got, err)
	}
	if got, err := os.ReadFile(included); err != nil || string(got) != contents {
		t.Fatalf("included file changed: %q: %v", got, err)
	}
	if got, err := repo.GetConfig("gitpour.version"); err != nil || got != "1.0" {
		t.Fatalf("git-pour setting changed: %q: %v", got, err)
	}
	if err := repo.RemoveLegacyConfigLocal(); err != nil {
		t.Fatalf("cleanup of already removed settings: %v", err)
	}
}
