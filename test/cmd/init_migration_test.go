package cmd_test

import (
	"strings"
	"testing"

	"github.com/OnionRings-dev/git-pour-now/test/testutil"
)

func TestInitLegacyConsent(t *testing.T) {
	for _, tc := range []struct {
		name, input       string
		imported, removed bool
	}{
		{"decline", "no\n\n\n", false, false},
		{"empty answer", "\n\n\n", false, false},
		{"EOF", "", false, false},
		{"import and retain", "yes\nno\n", true, false},
		{"import and remove", "yes\nyes\n", true, true},
		{"import and EOF", "yes\n", true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := testutil.SetupTestRepo(t)
			defer testutil.CleanupTestRepo(t, dir)
			setupGitFlowAVH(t, dir)
			for key, value := range map[string]string{
				"gitflow.feature.finish.keep": "true",
				"gitflow.origin":              "upstream",
				"review.preserve":             "untouched",
			} {
				if out, err := testutil.RunGit(t, dir, "config", key, value); err != nil {
					t.Fatalf("configure: %s: %v", out, err)
				}
			}
			for _, value := range []string{"ci.skip", "merge_request.create"} {
				if out, err := testutil.RunGit(t, dir, "config", "--add", "gitflow.feature.publish.pushoption", value); err != nil {
					t.Fatalf("configure push option: %s: %v", out, err)
				}
			}
			out, err := testutil.RunGitFlowWithInput(t, dir, tc.input, "init")
			if err != nil {
				t.Fatalf("init: %s: %v", out, err)
			}
			wantPrefix := "feature/"
			if tc.imported {
				wantPrefix = "feat/"
				if getGitConfig(t, dir, "gitpour.feature.finish.keep") != "true" || getGitConfig(t, dir, "gitpour.origin") != "upstream" {
					t.Fatal("imported operational settings were not persisted")
				}
			} else if strings.Contains(out, "Delete the previous") {
				t.Fatal("cleanup offered after declining import")
			}
			if got := getGitConfig(t, dir, "gitpour.branch.feature.prefix"); got != wantPrefix {
				t.Fatalf("prefix: got %q, want %q", got, wantPrefix)
			}
			legacy, err := testutil.RunGit(t, dir, "config", "--get-regexp", "^gitflow\\.")
			if tc.removed {
				if err == nil || strings.TrimSpace(legacy) != "" {
					t.Fatalf("legacy config remains: %s", legacy)
				}
			} else if err != nil {
				t.Fatal("legacy config was removed without consent")
			}
			if getGitConfig(t, dir, "review.preserve") != "untouched" {
				t.Fatal("unrelated config changed")
			}
			if tc.imported {
				if out, err := testutil.RunGit(t, dir, "config", "--get-all", "gitpour.feature.publish.pushoption"); err != nil || strings.TrimSpace(out) != "ci.skip\nmerge_request.create" {
					t.Fatalf("multi-valued preferences lost: %s: %v", out, err)
				}
				// A fresh process must honor the persisted retention preference.
				if out, err := testutil.RunGitFlow(t, dir, "feature", "start", "retained"); err != nil {
					t.Fatalf("start: %s: %v", out, err)
				}
				if out, err := testutil.RunGitFlow(t, dir, "feature", "finish", "retained"); err != nil {
					t.Fatalf("finish: %s: %v", out, err)
				}
				if !testutil.BranchExists(t, dir, "feat/retained") {
					t.Fatal("finish ignored imported keep preference")
				}
			}
		})
	}
}

func TestInitLegacyNextSharedImport(t *testing.T) {
	t.Parallel()
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)
	for key, value := range map[string]string{
		"gitflow.version":               "1.0",
		"gitflow.branch.main.type":      "base",
		"gitflow.branch.feature.type":   "topic",
		"gitflow.branch.feature.parent": "main",
		"gitflow.branch.feature.prefix": "custom/",
		"gitflow.feature.finish.keep":   "true",
	} {
		if out, err := testutil.RunGit(t, dir, "config", key, value); err != nil {
			t.Fatalf("configure: %s: %v", out, err)
		}
	}
	if out, err := testutil.RunGitFlowWithInput(t, dir, "y\ny\n", "init", "--shared"); err != nil {
		t.Fatalf("init: %s: %v", out, err)
	}
	if getGitConfig(t, dir, "gitpour.branch.feature.prefix") != "custom/" {
		t.Fatal("git-flow-next branch definition was not imported")
	}
	if out, err := testutil.RunGit(t, dir, "config", "--file", ".gitpour", "--get", "gitpour.feature.finish.keep"); err != nil || strings.TrimSpace(out) != "true" {
		t.Fatalf("shared imported setting: %s: %v", out, err)
	}
}

func TestInitLegacyDefaultsSkipMigration(t *testing.T) {
	t.Parallel()
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)
	setupGitFlowAVH(t, dir)
	if out, err := testutil.RunGitFlow(t, dir, "feature", "start", "before-import"); err == nil {
		t.Fatalf("ordinary command accepted legacy config without import: %s", out)
	}
	out, err := testutil.RunGitFlowWithInput(t, dir, "yes\nyes\n", "init", "--defaults")
	if err != nil {
		t.Fatalf("init: %s: %v", out, err)
	}
	if strings.Contains(out, "Import it?") || strings.Contains(out, "Delete the previous") {
		t.Fatalf("explicit defaults prompted for migration: %s", out)
	}
	if getGitConfig(t, dir, "gitpour.branch.feature.prefix") != "feature/" || getGitConfig(t, dir, "gitflow.prefix.feature") != "feat/" {
		t.Fatal("explicit defaults imported or removed legacy settings")
	}
}

func TestInitLegacyFailurePreservesOldConfig(t *testing.T) {
	t.Parallel()
	dir := testutil.SetupTestRepo(t)
	defer testutil.CleanupTestRepo(t, dir)
	setupGitFlowAVH(t, dir)
	if out, err := testutil.RunGit(t, dir, "config", "gitflow.branch.master", "bad..name"); err != nil {
		t.Fatalf("configure: %s: %v", out, err)
	}
	out, err := testutil.RunGitFlowWithInput(t, dir, "yes\nyes\n", "init")
	if err == nil {
		t.Fatalf("initialization with invalid branch succeeded: %s", out)
	}
	if strings.Contains(out, "Delete the previous") {
		t.Fatalf("failed initialization offered cleanup: %s", out)
	}
	if getGitConfig(t, dir, "gitflow.prefix.feature") != "feat/" {
		t.Fatal("failed initialization removed legacy config")
	}
}
