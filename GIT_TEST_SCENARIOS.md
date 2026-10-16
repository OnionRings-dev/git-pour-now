# Git Test Scenario Setup Guide

This document provides instructions for setting up various Git scenarios in tests, including merge conflicts, remote operations, and state verification.

> **Note**: This document is essential reading when implementing tests. It contains patterns for creating merge conflicts, setting up remotes, and verifying Git states.

For general testing conventions and patterns, see [TESTING_GUIDELINES.md](TESTING_GUIDELINES.md).

## Table of Contents

- [Basic Test Repository Setup](#basic-test-repository-setup)
- [Creating Merge and Rebase Conflicts](#creating-merge-and-rebase-conflicts)
- [Setting Up Local Remotes](#setting-up-local-remotes)
- [Verifying Git Conflict States](#verifying-git-conflict-states)
- [Common Mistakes to Avoid](#common-mistakes-to-avoid)

## Basic Test Repository Setup

All tests use temporary Git repositories created through test utilities.

```go
func TestExample(t *testing.T) {
    // Setup temporary repository
    dir := testutil.SetupTestRepo(t)
    defer testutil.CleanupTestRepo(t, dir)

    // Initialize git-pour with defaults
    output, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
    if err != nil {
        t.Fatalf("Failed to initialize git-pour: %v\nOutput: %s", err, output)
    }

    // Test implementation...
}
```

### Available Helper Functions

Available in `test/testutil/git.go`:

- `SetupTestRepo(t *testing.T) string` - Creates temporary Git repository
- `SetupTestRepoWithRemote(t *testing.T) (string, string)` - Creates repo with local remote
- `CleanupTestRepo(t *testing.T, dir string)` - Removes temporary repository
- `RunGit(t *testing.T, dir string, args ...string)` - Executes Git commands
- `RunGitFlow(t *testing.T, dir string, args ...string)` - Executes git-pour commands
- `WriteFile(t *testing.T, dir string, name string, content string)` - Creates files
- `BranchExists(t *testing.T, dir string, branch string) bool` - Checks branch existence
- `GetCurrentBranch(t *testing.T, dir string) string` - Gets current branch name
- `AddRemote(t *testing.T, dir, name string, bare bool) (string, error)` - Adds a remote repository

## Creating Merge and Rebase Conflicts

To create merge/rebase conflicts for testing, **sequence matters**.

### Correct Sequence for Conflict Generation

1. **Create branches BEFORE adding conflicting content**:
   - Start with clean base branch
   - Create topic branch from clean base
   - Add content to topic branch first
   - Switch back to base branch and add different content to same file
   - Attempting to merge/rebase will create conflict

```go
// CORRECT: Create branch first, then add conflicting content
output, err := testutil.RunGitFlow(t, dir, "feature", "start", "conflict-test")
if err != nil {
    t.Fatalf("Failed to create feature branch: %v", err)
}

// Add content to feature branch
testutil.WriteFile(t, dir, "test.txt", "feature content")
_, err = testutil.RunGit(t, dir, "add", "test.txt")
_, err = testutil.RunGit(t, dir, "commit", "-m", "Add test.txt in feature")

// Switch to base branch and add conflicting content
_, err = testutil.RunGit(t, dir, "checkout", "develop")
testutil.WriteFile(t, dir, "test.txt", "develop content")
_, err = testutil.RunGit(t, dir, "add", "test.txt")
_, err = testutil.RunGit(t, dir, "commit", "-m", "Add test.txt in develop")

// Now finish will create conflict
output, err = testutil.RunGitFlow(t, dir, "feature", "finish", "conflict-test")
```

### Incorrect Sequence (Won't Create Conflicts)

```go
// WRONG: Adding file to base first, then branching
testutil.WriteFile(t, dir, "test.txt", "develop content")
_, err = testutil.RunGit(t, dir, "add", "test.txt")
_, err = testutil.RunGit(t, dir, "commit", "-m", "Add test.txt in develop")

// Branch inherits the file - no conflict will occur
output, err := testutil.RunGitFlow(t, dir, "feature", "start", "conflict-test")
testutil.WriteFile(t, dir, "test.txt", "feature content") // Modifying existing file
```

### Multiple Rebase Conflicts

For testing complex rebase scenarios:
- Repeat the conflict generation pattern with additional commits and files
- Each conflicting commit will create a separate rebase step that requires resolution

### Child Branch Update Conflicts

To test conflicts during the child branch update phase of finish operations:

```go
// 1. Initialize and setup
dir := testutil.SetupTestRepo(t)
defer testutil.CleanupTestRepo(t, dir)
testutil.RunGitFlow(t, dir, "init", "--defaults")

// 2. Create a release branch
testutil.RunGitFlow(t, dir, "release", "start", "1.0.0")
testutil.WriteFile(t, dir, "release.txt", "release content")
testutil.RunGit(t, dir, "add", "release.txt")
testutil.RunGit(t, dir, "commit", "-m", "Release changes")

// 3. Add conflicting content to develop (child branch)
testutil.RunGit(t, dir, "checkout", "develop")
testutil.WriteFile(t, dir, "release.txt", "develop content")
testutil.RunGit(t, dir, "add", "release.txt")
testutil.RunGit(t, dir, "commit", "-m", "Develop changes")

// 4. Finish release - merge to main succeeds, but develop update conflicts
testutil.RunGit(t, dir, "checkout", "release/1.0.0")
output, err := testutil.RunGitFlow(t, dir, "release", "finish", "1.0.0")
// Conflict occurs during develop auto-update
```

## Setting Up Local Remotes

To test remote functionality, create a local bare repository as a remote.

### Basic Remote Setup

```go
// Add a remote repository
remoteDir, err := testutil.AddRemote(t, dir, "origin", true)
if err != nil {
    t.Fatalf("Failed to add remote: %v", err)
}
defer testutil.CleanupTestRepo(t, remoteDir)

// Create feature branch and make changes
output, err = testutil.RunGitFlow(t, dir, "feature", "start", "fetch-test")
testutil.WriteFile(t, dir, "feature.txt", "feature content")
_, err = testutil.RunGit(t, dir, "add", "feature.txt")
_, err = testutil.RunGit(t, dir, "commit", "-m", "Add feature file")

// Test operations with remote (fetch, push, etc.)
output, err = testutil.RunGitFlow(t, dir, "feature", "finish", "fetch-test", "--fetch")
```

### Using SetupTestRepoWithRemote Helper

For simpler setup when you need a repo with remote pre-configured:

```go
func TestWithRemote(t *testing.T) {
    dir, remoteDir := testutil.SetupTestRepoWithRemote(t)
    defer testutil.CleanupTestRepo(t, dir)
    defer testutil.CleanupTestRepo(t, remoteDir)

    // Initialize git-pour
    output, err := testutil.RunGitFlow(t, dir, "init", "--defaults")
    if err != nil {
        t.Fatalf("Failed to initialize git-pour: %v", err)
    }

    // Remote 'origin' is already configured and tracking branches pushed
    // Test remote operations...
}
```

### Simulating Remote Changes

To test fetch scenarios where remote has changes:

```go
// 1. Setup repo with remote
dir, remoteDir := testutil.SetupTestRepoWithRemote(t)

// 2. Clone the remote to a second working copy
secondDir := t.TempDir()
testutil.RunGit(t, secondDir, "clone", remoteDir, ".")

// 3. Make changes in second copy and push
testutil.WriteFile(t, secondDir, "remote-change.txt", "from remote")
testutil.RunGit(t, secondDir, "add", ".")
testutil.RunGit(t, secondDir, "commit", "-m", "Remote change")
testutil.RunGit(t, secondDir, "push", "origin", "main")

// 4. Original repo now has outdated refs - fetch will get new changes
output, err := testutil.RunGitFlow(t, dir, "feature", "finish", "test", "--fetch")
```

## Verifying Git Conflict States

Different Git operations create different internal states that require specific verification.

### Merge Conflict Verification

```go
// Check for merge conflict state
if _, err := os.Stat(filepath.Join(dir, ".git", "MERGE_HEAD")); os.IsNotExist(err) {
    t.Error("Expected to be in merge conflict state")
}

// Check conflicted files
output, _ := testutil.RunGit(t, dir, "status", "--porcelain")
if !strings.Contains(output, "UU") { // UU = both modified (conflict)
    t.Error("Expected conflicted files")
}
```

### Rebase Conflict Verification

During rebase conflicts, standard Git commands show misleading information:

```go
// WRONG: Current branch shows "HEAD" during rebase conflict
currentBranch := testutil.GetCurrentBranch(t, dir) // Returns "HEAD"
if currentBranch != "feature/branch-name" { // This will always fail
    t.Error("Not on expected branch")
}

// CORRECT: Check rebase state via Git internal files
rebaseMergeDir := filepath.Join(dir, ".git", "rebase-merge")
if _, err := os.Stat(rebaseMergeDir); os.IsNotExist(err) {
    t.Error("Expected to be in rebase conflict state")
}

// Verify which branch is being rebased
headNameFile := filepath.Join(rebaseMergeDir, "head-name")
if headNameBytes, err := os.ReadFile(headNameFile); err != nil {
    t.Errorf("Failed to read head-name file: %v", err)
} else {
    headName := strings.TrimSpace(string(headNameBytes))
    expectedBranch := "refs/heads/feature/branch-name" // Full ref path
    if headName != expectedBranch {
        t.Errorf("Expected rebasing %s, got %s", expectedBranch, headName)
    }
}
```

### Squash Merge State

```go
// After a squash merge with conflicts, check SQUASH_MSG exists
squashMsgPath := filepath.Join(dir, ".git", "SQUASH_MSG")
if _, err := os.Stat(squashMsgPath); os.IsNotExist(err) {
    t.Error("Expected squash merge state")
}
```

### Application Merge State

git-pour-now maintains its own merge state for multi-step operations:

```go
// Check application state
state, err := testutil.LoadMergeState(t, dir)
if err != nil {
    t.Fatalf("Failed to load merge state: %v", err)
}

if state.Action != "finish" {
    t.Errorf("Expected action to be 'finish', got '%s'", state.Action)
}

if state.CurrentStep != "merge" {
    t.Errorf("Expected step to be 'merge', got '%s'", state.CurrentStep)
}
```

## Common Mistakes to Avoid

### Wrong Conflict Generation Sequence

```go
// BAD: Creates file on base, then branches (inherits file - no conflict)
testutil.WriteFile(t, dir, "test.txt", "base content")
testutil.RunGit(t, dir, "add", ".")
testutil.RunGit(t, dir, "commit", "-m", "Add file to base")
testutil.RunGitFlow(t, dir, "feature", "start", "test")
testutil.WriteFile(t, dir, "test.txt", "feature content") // Modifying, not conflicting
```

### Incorrect Rebase State Verification

```go
// BAD: Expects branch name during rebase conflict (always returns "HEAD")
currentBranch := testutil.GetCurrentBranch(t, dir)
if currentBranch != "feature/branch-name" {
    t.Error("Not on expected branch") // Will always fail
}
```

### Testing Multiple Behaviors with Fragile Dependencies

```go
// BAD: Test fails if fetch errors, even though it's testing continue behavior
// TestFinishContinueDoesNotFetch - but tries to test fetch too!

// This will fail fatally if no remote exists
output, _ := RunGitFlow(t, dir, "feature", "finish", "test", "--fetch")

// Never reaches here because fetch failed and aborted before merge
if !strings.Contains(output, "conflict") {
    t.Error("Expected conflict") // Test fails for wrong reason
}

// Continue never runs because no merge state was created
RunGitFlow(t, dir, "feature", "finish", "--continue") // Error: no merge in progress
```

**Problem**: The test name says "ContinueDoesNotFetch" but setup requires fetch to succeed. If fetch fails (no remote), the merge never starts, and continue can't be tested.

**Solution**: Match test setup to test goal:

```go
// GOOD: Test only what the test name promises
// TestFinishContinueDoesNotFetch - tests ONLY that continue doesn't fetch

// Finish without fetch to create conflict (simple setup)
output, _ := RunGitFlow(t, dir, "feature", "finish", "test")

// Verify conflict occurred
if !strings.Contains(output, "conflict") {
    t.Error("Expected conflict")
}

// Resolve conflict
WriteFile(t, dir, "conflict.txt", "resolved")
RunGit(t, dir, "add", "conflict.txt")

// Continue and verify NO fetch happens
continueOutput, _ := RunGitFlow(t, dir, "feature", "finish", "--continue")
if strings.Contains(continueOutput, "Fetching") {
    t.Error("Continue should not fetch") // Tests exactly one behavior
}
```

### Placeholder Helper Functions

```go
// BAD: Helper function that doesn't actually test anything
func captureConfigAddBase(name string) error {
    defer func() {
        if r := recover(); r != nil {
            // Convert panic to error for testing
        }
    }()
    return nil // Placeholder - doesn't actually call git-pour
}

// GOOD: Actually calls the command being tested
func captureConfigAddBase(t *testing.T, dir, name, parent string) error {
    args := []string{"config", "add", "base", name, parent}
    _, err := testutil.RunGitFlow(t, dir, args...)
    return err
}
```

## Quick Reference

| Scenario | Key Setup Step |
|----------|----------------|
| Merge conflict | Create branch first, then add conflicting content to both branches |
| Rebase conflict | Same as merge, but use rebase strategy |
| Remote operations | Use `AddRemote()` or `SetupTestRepoWithRemote()` |
| Verify merge conflict | Check `.git/MERGE_HEAD` exists |
| Verify rebase conflict | Check `.git/rebase-merge/` directory exists |
| Child branch conflict | Add conflicting content to child branch before finishing parent |
