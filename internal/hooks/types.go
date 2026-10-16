// Package hooks provides filter and hook execution for git-pour operations.
//
// Filters transform input/output values (e.g., version numbers, tag messages).
// Hooks execute scripts before/after operations (e.g., pre-pour-feature-start).
//
// All scripts are located in .git/hooks/ following these patterns:
//   - Filters: filter-pour-{type}-{action}-{target}
//   - Hooks: {pre,post}-pour-{type}-{action}
package hooks

import "fmt"

// FilterTarget represents what is being filtered.
type FilterTarget string

const (
	// FilterTargetVersion filters version/name for start operations.
	FilterTargetVersion FilterTarget = "version"

	// FilterTargetTagMessage filters tag message for finish operations.
	FilterTargetTagMessage FilterTarget = "tag-message"
)

// GetFilterName returns the filter script name for a given branch type, action, and target.
// Example: GetFilterName("release", "start", FilterTargetVersion) returns "filter-pour-release-start-version"
func GetFilterName(branchType string, action string, target FilterTarget) string {
	return fmt.Sprintf("filter-pour-%s-%s-%s", branchType, action, target)
}

// FilterContext contains data passed to filters.
type FilterContext struct {
	BranchType string // The type of branch (e.g., "release", "hotfix", "feature")
	BranchName string // The short name of the branch
	Version    string // The version/name being filtered
	TagMessage string // The tag message to filter
	BaseBranch string // The base/parent branch
	FullBranch string // The full branch name with prefix
	Origin     string // The remote name
}

// HookPhase represents pre or post execution.
type HookPhase string

const (
	// HookPre runs before the operation.
	HookPre HookPhase = "pre"

	// HookPost runs after the operation.
	HookPost HookPhase = "post"
)

// HookAction represents the git-pour action being performed.
type HookAction string

const (
	// HookActionStart is the start action.
	HookActionStart HookAction = "start"

	// HookActionFinish is the finish action.
	HookActionFinish HookAction = "finish"

	// HookActionPublish is the publish action.
	HookActionPublish HookAction = "publish"

	// HookActionTrack is the track action.
	HookActionTrack HookAction = "track"

	// HookActionDelete is the delete action.
	HookActionDelete HookAction = "delete"

	// HookActionUpdate is the update action.
	HookActionUpdate HookAction = "update"
)

// HookContext contains data passed to hooks via environment variables.
type HookContext struct {
	BranchType string // The type of branch (e.g., "feature", "release")
	BranchName string // The short name of the branch
	FullBranch string // The full branch name with prefix
	BaseBranch string // The base/parent branch
	Origin     string // The remote name
	Version    string // The version (for branches with tagging)
	ExitCode   int    // For post-hooks: exit code of the operation
}

// HookResult contains the result of hook execution.
type HookResult struct {
	Executed bool   // Whether the hook was found and executed
	ExitCode int    // Exit code of the hook script
	Output   string // Combined stdout/stderr output
	Error    error  // Any error that occurred during execution
}
