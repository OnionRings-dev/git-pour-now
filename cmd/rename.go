package cmd

import (
	"fmt"
	"os"

	"github.com/OnionRings-dev/git-pour-now/internal/config"
	"github.com/OnionRings-dev/git-pour-now/internal/errors"
	"github.com/OnionRings-dev/git-pour-now/internal/git"
	"github.com/OnionRings-dev/git-pour-now/internal/worktree"
)

// RenameCommand handles renaming a topic branch
func RenameCommand(branchType string, oldName string, newName string) {
	repo := mustOpenRepo()
	if err := executeRename(repo, branchType, oldName, newName); err != nil {
		var exitCode errors.ExitCode
		if pourErr, ok := err.(errors.Error); ok {
			exitCode = pourErr.ExitCode()
		} else {
			exitCode = errors.ExitCodeGitError
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(int(exitCode))
	}
}

// executeRename performs the actual branch renaming logic and returns any errors
func executeRename(repo *git.Repo, branchType string, oldName string, newName string) error {
	// Validate that git-pour is initialized before resolving branch types.
	// LoadConfig falls back to DefaultConfig when uninitialized, so this gate
	// must run first or the default branch types mask the uninitialized state.
	initialized, err := config.IsInitialized(repo)
	if err != nil {
		return &errors.GitError{Operation: "check if git-pour is initialized", Err: err}
	}
	if !initialized {
		return &errors.NotInitializedError{}
	}

	// Load configuration
	cfg, err := config.Load(repo)
	if err != nil {
		return &errors.GitError{Operation: "load configuration", Err: err}
	}

	// Get branch configuration
	branchConfig, ok := cfg.Branches[branchType]
	if !ok {
		return &errors.InvalidBranchTypeError{BranchType: branchType}
	}

	// Construct full branch names
	oldFullBranchName := oldName
	newFullBranchName := newName
	if branchConfig.Prefix != "" {
		oldFullBranchName = branchConfig.Prefix + oldName
		newFullBranchName = branchConfig.Prefix + newName
	}

	// Check if old branch exists
	err = repo.BranchExists(oldFullBranchName)
	if err != nil {
		return &errors.BranchNotFoundError{BranchName: oldFullBranchName}
	}

	// Check if new branch name already exists
	err = repo.BranchExists(newFullBranchName)
	if err == nil {
		return &errors.GitError{Operation: "rename branch", Err: fmt.Errorf("branch '%s' already exists", newFullBranchName)}
	}

	// git branch -m handles every checkout state itself, including a branch
	// checked out in the current worktree or in a linked one, which follows the
	// rename.
	if err := repo.RenameBranch(oldFullBranchName, newFullBranchName); err != nil {
		return &errors.GitError{Operation: fmt.Sprintf("rename branch '%s' to '%s'", oldFullBranchName, newFullBranchName), Err: err}
	}

	// The ref has moved, so git-pour's per-branch state has to follow it. This
	// runs only after RenameBranch succeeds: a refused or failed rename must
	// leave config untouched.
	migrateBranchState(repo, oldFullBranchName, newFullBranchName)

	fmt.Printf("Renamed branch '%s' to '%s'\n", oldFullBranchName, newFullBranchName)
	return nil
}

// migrateBranchState moves the config keys that git-pour keys by BRANCH NAME onto
// the branch's new name: the worktree provenance marker and the recorded start
// point. Both are repository-local runtime state, never user settings, and both
// are already excluded from the shared-config set.
//
// A failure is a warning naming the key, not an error: the branch has already
// been renamed and there is nothing to roll back to. delete and finish treat a
// failed .base cleanup the same way. One key failing must not stop the others
// from migrating.
//
// worktree.MarkManaged cannot perform the marker move — it always writes "true",
// which would turn a hand-written "false" into a claim that git-pour created the
// worktree. The value is copied raw.
func migrateBranchState(repo *git.Repo, oldBranch string, newBranch string) {
	moves := []struct{ from, to string }{
		{worktree.MarkerKey(oldBranch), worktree.MarkerKey(newBranch)},
		{git.BaseBranchKey(oldBranch), git.BaseBranchKey(newBranch)},
	}
	for _, move := range moves {
		if err := repo.MoveConfigLocal(move.from, move.to); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to move config %s to %s: %v\n", move.from, move.to, err)
		}
	}
}
