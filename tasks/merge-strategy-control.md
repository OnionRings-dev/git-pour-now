# Add Merge Strategy Control to Finish Command

Implement comprehensive merge strategy control for the finish command, allowing users to customize how topic branches are integrated back into target branches.

## Overview

Add command-line flags and configuration options to control merge behavior during finish operations. This includes rebase options, fast-forward control, and squash behavior while maintaining the existing three-layer configuration hierarchy.

## Current State Analysis

**✅ Existing Infrastructure:**
- Robust configuration resolver system in `internal/config/resolver.go`
- Three-layer configuration hierarchy (branch defaults → command config → CLI flags)
- Basic merge strategy execution in `cmd/finish.go` supporting `merge`, `rebase`, and `squash`
- Support for existing `gitpour.branch.<topic>.upstreamstrategy` configuration

**🔧 Current Limitations:**
- No granular control over rebase options (preserve-merges)
- No control over fast-forward behavior (--no-ff/--ff)
- No command-line flags to override merge strategies
- Limited merge strategy configuration options

## Implementation Plan

### Phase 1: Extend Configuration Structures

#### 1.1 Add Merge Strategy Options to Configuration

**File**: `internal/config/resolver.go`

```go
// Add new struct for merge strategy command-line options
type MergeStrategyOptions struct {
    Strategy        *string  // Override for entire strategy
    Rebase          *bool    // --rebase/--no-rebase override
    PreserveMerges  *bool    // --preserve-merges/--no-preserve-merges
    NoFF            *bool    // --no-ff/--ff
    Squash          *bool    // --squash/--no-squash override
}

// Extend ResolvedFinishOptions
type ResolvedFinishOptions struct {
    // Existing fields...

    // New merge strategy fields
    MergeStrategy     string  // Final resolved strategy (merge/rebase/squash)
    UseRebase         bool    // Whether to use rebase
    PreserveMerges    bool    // Whether to preserve merges during rebase
    NoFastForward     bool    // Whether to create merge commit for fast-forward
    UseSquash         bool    // Whether to squash commits
}
```

#### 1.2 Update Command Flag System

**File**: `cmd/topicbranch.go` - `addFinishFlags` function

```go
func addFinishFlags(cmd *cobra.Command) {
    // Existing flags (tag, retention, operation control)...

    // New merge strategy flags
    cmd.Flags().Bool("rebase", false, "Rebase topic branch before merging")
    cmd.Flags().Bool("no-rebase", false, "Don't rebase topic branch (use configured strategy)")
    cmd.Flags().Bool("preserve-merges", false, "Preserve merges during rebase")
    cmd.Flags().Bool("no-preserve-merges", false, "Flatten merges during rebase")
    cmd.Flags().Bool("no-ff", false, "Create merge commit even for fast-forward")
    cmd.Flags().Bool("ff", false, "Allow fast-forward merge when possible")
    cmd.Flags().Bool("squash", false, "Squash all commits into single commit")
    cmd.Flags().Bool("no-squash", false, "Keep individual commits (don't squash)")
}
```

#### 1.3 Update Command Execution

**File**: `cmd/finish.go` - Update command execution to parse new flags

```go
// Add merge strategy options parsing
func FinishCommand(branchType string, name string, continueOp bool, abortOp bool, force bool, tagOptions *config.TagOptions, retentionOptions *config.BranchRetentionOptions, mergeOptions *config.MergeStrategyOptions) {
    // Implementation with new merge options parameter
}
```

### Phase 2: Configuration Resolution Logic

#### 2.1 New Configuration Keys

Support these configuration keys in three-layer hierarchy:

```bash
# Layer 1: Branch defaults (existing)
gitpour.branch.<topic>.upstreamstrategy = merge|rebase|squash

# Layer 2: Command-specific overrides (new)
gitpour.<topic>.finish.rebase = true|false
gitpour.<topic>.finish.preserve-merges = true|false
gitpour.<topic>.finish.ff = true|false
gitpour.<topic>.finish.squash = true|false

# Layer 3: Command-line flags (highest priority)
```

#### 2.2 Resolution Functions

**File**: `internal/config/resolver.go`

Add these resolution functions:

```go
func resolveMergeStrategy(cfg *Config, branchConfig BranchConfig, branchType string, mergeOpts *MergeStrategyOptions) (string, bool, bool, bool, bool)

func resolveFinishRebase(cfg *Config, branchType string, baseStrategy string, mergeOpts *MergeStrategyOptions) bool

func resolveFinishPreserveMerges(cfg *Config, branchType string, mergeOpts *MergeStrategyOptions) bool

func resolveFinishNoFF(cfg *Config, branchType string, mergeOpts *MergeStrategyOptions) bool

func resolveFinishSquash(cfg *Config, branchType string, baseStrategy string, mergeOpts *MergeStrategyOptions) bool
```

#### 2.3 Strategy Precedence Logic

1. **Squash flag/config** - Highest precedence when explicitly set
2. **Rebase flag/config** - Second precedence when explicitly set
3. **Base strategy** - From `gitpour.branch.<topic>.upstreamstrategy`
4. **Default** - Fall back to "merge"

### Phase 3: Enhanced Git Operations

#### 3.1 Add Git Operation Functions

**File**: `internal/git/repo.go`

```go
// Add new Git operation functions
func RebaseWithOptions(targetBranch string, preserveMerges bool) error

func MergeWithOptions(branchName string, noFF bool) error

func MergeSquashWithMessage(branchName string, message string) error
```

#### 3.2 Update Merge Logic

**File**: `cmd/finish.go` - `handleMergeStep` function

Replace current merge logic with enhanced strategy handling:

```go
func handleMergeStep(cfg *config.Config, state *mergestate.MergeState, branchConfig config.BranchConfig, resolvedOptions *config.ResolvedFinishOptions) error {
    // Enhanced merge logic using resolvedOptions.MergeStrategy, UseRebase, PreserveMerges, NoFastForward, UseSquash

    switch resolvedOptions.MergeStrategy {
    case "rebase":
        // Rebase with preserve-merges option
        // Then fast-forward merge with no-ff option
    case "squash":
        // Squash merge
    case "merge":
        // Standard merge with no-ff option
    }
}
```

### Phase 4: Comprehensive Test Suite

#### 4.1 Configuration Resolution Tests

**File**: `test/internal/config/resolver_merge_test.go`

Individual test functions for each scenario:

- `TestResolveMergeStrategyBranchDefault`
- `TestResolveMergeStrategyRebaseFlagOverride`
- `TestResolveMergeStrategySquashFlagOverride`
- `TestResolveMergeStrategyConfigOverride`
- `TestResolveMergeStrategyFlagOverridesConfig`
- `TestResolveMergeStrategyPreserveMergesDefault`
- `TestResolveMergeStrategyNoFFDefault`

#### 4.2 Integration Tests for Each Flag

**File**: `test/cmd/finish_merge_strategy_test.go`

One test function per flag and scenario:

- `TestFinishWithRebaseFlag`
- `TestFinishWithNoRebaseFlag`
- `TestFinishWithSquashFlag`
- `TestFinishWithNoSquashFlag`
- `TestFinishWithNoFFFlag`
- `TestFinishWithFFFlag`
- `TestFinishWithPreserveMergesFlag`
- `TestFinishWithNoPreserveMergesFlag`

#### 4.3 Cross-Branch Type Tests

**File**: `test/cmd/finish_merge_strategy_branch_types_test.go`

Separate test functions for each branch type and flag combination:

- `TestFinishFeatureWithRebaseFlag`
- `TestFinishReleaseWithRebaseFlag`
- `TestFinishHotfixWithRebaseFlag`
- `TestFinishSupportWithRebaseFlag`
- `TestFinishFeatureWithSquashFlag`
- `TestFinishReleaseWithSquashFlag`
- etc.

#### 4.4 Configuration Hierarchy Tests

**File**: `test/cmd/finish_merge_config_hierarchy_test.go`

Individual tests for configuration precedence:

- `TestMergeStrategyBranchDefaultUsed`
- `TestMergeStrategyCommandConfigOverridesBranch`
- `TestMergeStrategyFlagOverridesCommandConfig`
- `TestMergeStrategyFlagOverridesBranchDefault`

#### 4.5 Error Handling Tests

**File**: `test/cmd/finish_merge_strategy_errors_test.go`

- `TestMergeStrategyConflictHandlingWithRebase`
- `TestMergeStrategyConflictHandlingWithSquash`
- `TestMergeStrategyConflictHandlingWithMerge`
- `TestMergeStrategyContinueAfterConflict`
- `TestMergeStrategyAbortAfterConflict`

### Phase 5: Documentation Updates

#### 5.1 Update Command Documentation

**File**: `docs/git-pour-finish.1.md`

Add comprehensive merge strategy section with:
- All new command-line flags
- Configuration hierarchy explanation
- Examples for each merge strategy
- Conflict resolution workflows

#### 5.2 Update Configuration Reference

**File**: `docs/gitpour-config.5.md`

Add new configuration keys:
- `gitpour.<branchtype>.finish.rebase`
- `gitpour.<branchtype>.finish.preserve-merges`
- `gitpour.<branchtype>.finish.ff`
- `gitpour.<branchtype>.finish.squash`

## Default Behaviors

- **Rebase**: Determined by `gitpour.branch.<topic>.upstreamstrategy`
- **Preserve Merges**: Default `false` (flatten merges during rebase)
- **Fast Forward**: Default `true` (allow fast-forward when possible)
- **Squash**: Determined by `gitpour.branch.<topic>.upstreamstrategy`

## Validation Rules

1. **`--preserve-merges`** only valid with rebase operations
2. **`--squash`** and **`--rebase`** are mutually exclusive when both explicitly set
3. **Flag precedence**: `--squash` > `--rebase` > base strategy
4. **Positive flags override negative flags** when both present
5. **Explicit config validation** - warn about incompatible configurations

## Testing Strategy

**CRITICAL**: Follow the "One Test Case Per Function" rule established in TESTING_GUIDELINES.md:

- Each test function tests exactly ONE scenario
- No table-driven tests with multiple test cases in integration tests
- Use descriptive function names that clearly indicate what is being tested
- Include proper test comments with numbered steps
- Test functions should be focused and easy to debug when they fail

**Exception**: Table-driven tests are acceptable only for pure validation functions with simple input/output validation.

## Implementation Order

1. **Phase 1**: Extend configuration structures and flag system
2. **Phase 2**: Implement configuration resolution logic
3. **Phase 3**: Update Git operations and merge logic
4. **Phase 4**: Add comprehensive test suite (following one-test-per-function rule)
5. **Phase 5**: Update documentation

## Backward Compatibility

- All existing `gitpour.branch.<topic>.upstreamstrategy` configurations continue to work
- New flags default to `false` (no behavior change unless explicitly used)
- Existing merge strategies (`merge`, `rebase`, `squash`) work as before
- New configuration keys are optional and have sensible defaults

## Success Criteria

- [ ] All new command-line flags implemented and working
- [ ] Three-layer configuration hierarchy respected for all new options
- [ ] Comprehensive test suite with one test case per function
- [ ] Documentation updated with examples and configuration reference
- [ ] Backward compatibility maintained
- [ ] All existing tests continue to pass
- [ ] Integration tests cover all branch types (feature, release, hotfix, support)
- [ ] Error handling and conflict resolution work with all merge strategies
