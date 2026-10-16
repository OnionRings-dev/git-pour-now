# git-pour-now Technical Overview

The default git-pour-now workflow uses `main`, `feature/`, and temporary
`prerelease/<version>` branches. Stable tags are created on main; prerelease tags
on the matching prerelease branch. No develop branch is created. Examples below
using develop/release/hotfix belong to the explicit `--preset=classic` workflow.
The test suite requests that preset explicitly for inherited GitFlow scenarios.


git-pour-now is a modern, flexible implementation of Git workflow management that builds upon the original git-pour concepts while providing extensive customization capabilities for modern development practices.

> **See also:** [CODE_REFERENCE.md](CODE_REFERENCE.md) for a quick codebase navigation guide with file paths, struct definitions, and implementation details.

## Project Structure

The project follows Go best practices with clear separation of concerns:

```
git-pour-now/
├── cmd/                    # Command implementations
│   ├── root.go            # Root CLI command setup with Cobra
│   ├── init.go            # Repository initialization command
│   ├── config.go          # config subcommands (add/edit/rename/delete/list)
│   ├── start.go           # Branch starting logic
│   ├── finish.go          # Branch finishing logic (most complex)
│   ├── update.go          # Branch updating from parent
│   ├── integrate.go       # Integrate a base branch into its parent
│   ├── topicbranch.go     # Dynamic command registration for branch types
│   ├── shorthand.go       # Shorthand commands (auto-detect branch type)
│   ├── continue_abort.go  # Continue/abort in-progress operations; foreign-op guard
│   ├── preflight.go       # Shared pre-merge fetch/sync preflight guard
│   ├── publish.go         # Publish a topic branch to the remote
│   ├── track.go           # Track a remote topic branch
│   ├── list.go            # Branch listing commands
│   ├── checkout.go        # Branch checkout functionality
│   ├── delete.go          # Branch deletion
│   ├── rename.go          # Branch renaming
│   ├── worktree.go        # Worktree management by branch name
│   └── overview.go        # Repository overview
├── internal/              # Internal packages (not exported)
│   ├── config/           # Git configuration management
│   │   └── config.go     # Branch type definitions, config loading
│   ├── git/              # Git command wrapper
│   │   ├── repo.go       # Git operations with error handling
│   │   └── worktree.go   # Worktree add/remove/list/prune/detach
│   ├── worktree/         # Worktree path templates and provenance markers
│   │   ├── path.go       # gitpour.worktreePath expansion
│   │   └── provenance.go # Markers recording which worktrees git-pour created
│   ├── navigate/         # GIT_POUR_CD_FILE destination channel
│   │   └── cdfile.go     # Hand a destination directory to the calling shell
│   ├── mergestate/       # Merge conflict state persistence
│   │   └── mergestate.go # State management for multi-step operations
│   ├── hooks/            # Client-side hook execution
│   │   ├── hooks.go      # Hook runner and context construction
│   │   ├── filters.go    # Hook matching/filtering
│   │   └── types.go      # Hook type definitions
│   ├── errors/           # Custom error types and exit codes
│   │   └── errors.go     # Structured error handling
│   ├── util/             # Validation and utility functions
│   │   └── validation.go # Input validation helpers
│   └── update/           # Branch updating logic
│       └── update.go     # Shared update functionality
├── test/                  # Test files (mirrors source structure)
│   ├── cmd/              # Command-level integration tests
│   ├── internal/         # Internal package unit tests
│   └── testutil/         # Test utilities and Git repo helpers
├── scripts/              # Build and deployment scripts
│   └── build.sh          # Multi-platform build script
├── main.go               # Application entry point
└── [documentation files] # README.md, ARCHITECTURE.md, etc.
```

### Key Organizational Principles

- **cmd/**: Contains all CLI command implementations using the Cobra framework
- **internal/**: Private packages that handle core functionality (config, git operations, state management)
- **test/**: Mirrors the source structure with comprehensive test coverage
- **testutil/**: Shared testing utilities, especially Git repository helpers

## Core Architecture

### Branch Dependency Model

The foundation of git-pour-now is a **branch dependency model** that formalizes the parent-child relationship between branches. This model enables:

- Automatic tracking of dependencies between branches
- Intelligent synchronization of changes across branch hierarchies
- Consistent propagation of changes through the dependency tree
- Visualization of branch relationships

Every branch (except root branches) has a parent, and branches can have multiple children. This simple paradigm enables powerful workflow customization.

### Branch Types

git-pour-now defines two fundamental branch types:

#### Base Branches (Long-living)
- Exist throughout the project lifecycle
- Serve as integration points for features and releases
- Examples: `main`, `develop`, `staging`, `production`
- Configured with parent-child relationships for change propagation

#### Topic Branches (Short-living)
- Created for specific purposes (features, hotfixes, releases)
- Always have a defined parent base branch
- Automatically cleaned up after completion
- Examples: `feature/login`, `hotfix/security-fix`, `release/v1.0`

## Single Topic Branch Implementation

### Unified Command Structure

Instead of separate commands for `feature`, `hotfix`, and `release` branches, git-pour-now implements a **single topic branch mechanism**:

```bash
# Traditional git-pour
git pour feature start my-feature
git pour hotfix start critical-fix
git pour release start v1.0

# git-pour-now unified approach
git pour feature start my-feature
git pour hotfix start critical-fix
git pour release start v1.0
```

All topic branches use the same `start` and `finish` commands, with behavior determined by configuration rather than branch type.

### Configurable Behavior

Topic branch behavior is defined through Git configuration at two levels. **Branch type configuration** (Layer 1) defines the branch type's identity and process characteristics:

- **Parent branch**: Which base branch to branch from (structural)
- **Start point**: Where to create the branch (structural)
- **Merge strategies**: How changes flow upstream and downstream (process)
- **Tag creation**: Whether the branch type produces tags on finish (process)
- **Child branch updates**: Automatic updating of child base branches after finish (process)

**Command-specific configuration** (Layer 2) then controls operational details like fetch behavior, tag signing, and branch retention. CLI flags (Layer 3) override everything for one-off situations.

## Configuration System

### Default Configuration Overview

git-pour-now provides sensible defaults that work for most teams while remaining fully customizable.

#### Branch Structure
```
main/master     ← Production releases
    ↓
develop         ← Integration branch (auto-updated from main)
    ↓
feature/        ← New features
release/        ← Release preparation
hotfix/         ← Emergency fixes
```

#### Base Branches

| Branch | Type | Parent | Config Key | Auto-Update from Parent |
|--------|------|--------|------------|------------------------|
| `main` | base | (root) | `gitpour.branch.main` | None |
| `develop` | base | `main` | `gitpour.branch.develop` | ✅ Yes |

#### Topic Branches

| Branch Type | Prefix | Parent | Start Point | Config Key | Created by Default |
|-------------|--------|--------|-------------|------------|-------------------|
| Feature | `feature/` | `develop` | `develop` | `gitpour.branch.feature` | ✅ Yes |
| Release | `release/` | `main` | `develop` | `gitpour.branch.release` | ✅ Yes |
| Hotfix | `hotfix/` | `main` | `main` | `gitpour.branch.hotfix` | ✅ Yes |

#### Merge Strategies

**Upstream Strategy (Finish Operations)** - How topic branches merge into their parent:

| Branch Type | Default | Options | Target Branch |
|-------------|---------|---------|---------------|
| Feature | `merge` | `merge`, `rebase`, `squash` | → `develop` |
| Release | `merge` | `merge`, `rebase`, `squash` | → `main` |
| Hotfix | `merge` | `merge`, `rebase`, `squash` | → `main` |

**Downstream Strategy (Update Operations)** - How parent updates are pulled into topic branches:

| Branch Type | Default | Options | Source Branch |
|-------------|---------|---------|---------------|
| Feature | `rebase` | `merge`, `rebase` | ← `develop` |
| Release | `merge` | `merge`, `rebase` | ← `main` |
| Hotfix | `rebase` | `merge`, `rebase` | ← `main` |

**Note**: The `--rebase` flag can be used with the `update` command to override the configured strategy and force rebase behavior.

#### Tag Configuration

| Branch Type | Default Tagging | Tag Prefix | When Tagged |
|-------------|-----------------|------------|-------------|
| Feature | ❌ Disabled | (configurable) | Never by default |
| Bugfix | ❌ Disabled | (configurable) | Never by default |
| Release | ✅ Enabled | v (during init) | On finish |
| Hotfix | ✅ Enabled | v (during init) | On finish |
| Support | ❌ Disabled | (configurable) | Never by default |

#### Branch Retention (After Finish)

| Setting | Default | Description |
|---------|---------|-------------|
| Delete Local | ✅ Yes | Remove local branch after successful merge |
| Delete Remote | ✅ Yes | Remove remote branch after successful merge |
| Force Delete | ❌ No | Use safe delete (checks for unmerged commits) |

#### Core Configuration Commands

```bash
# Base branch names
git config gitpour.branch.main main
git config gitpour.branch.develop develop

# Base branch relationships
git config gitpour.branch.develop.parent main
git config gitpour.branch.develop.upstreamStrategy merge
git config gitpour.branch.develop.downstreamStrategy merge
git config gitpour.branch.develop.autoUpdate true

# Topic branch prefixes
git config gitpour.branch.feature.prefix feature/
git config gitpour.branch.release.prefix release/
git config gitpour.branch.hotfix.prefix hotfix/

# Branch relationships
git config gitpour.branch.feature.parent develop
git config gitpour.branch.release.parent main
git config gitpour.branch.hotfix.parent main

# Merge strategies (upstream - finish operations)
git config gitpour.branch.feature.upstreamStrategy merge
git config gitpour.branch.release.upstreamStrategy merge
git config gitpour.branch.hotfix.upstreamStrategy merge

# Merge strategies (downstream - update operations)
git config gitpour.branch.feature.downstreamStrategy rebase
git config gitpour.branch.release.downstreamStrategy merge
git config gitpour.branch.hotfix.downstreamStrategy rebase

# Tag settings
git config gitpour.feature.finish.notag true
git config gitpour.release.finish.notag false
git config gitpour.hotfix.finish.notag false

# Usage examples
git pour update feature/my-feature              # Uses configured strategy
git pour update feature/my-feature --rebase     # Forces rebase strategy
git pour rebase                                 # Shorthand for update --rebase
```

**Note**: Release and hotfix branches merge only into `main`, then `develop` is automatically updated from `main` to stay synchronized.

### Branch Configuration Structure

Base branches are configured with dependency relationships:

```ini
[gitpour "branch.main"]
    type = base
    parent =
    upstreamStrategy = none
    downstreamStrategy = none

[gitpour "branch.develop"]
    type = base
    parent = main
    upstreamStrategy = merge
    downstreamStrategy = merge
```

Topic branch types are configured with the same key format:

```ini
[gitpour "branch.feature"]
    type = topic
    parent = develop
    startPoint = develop
    upstreamStrategy = rebase
    downstreamStrategy = merge
```

### Configurable Properties (Layer 1 — Branch Type Definition)

These properties define the branch type's identity and process characteristics. They describe *what the branch type is*, not how individual commands behave.

#### For Base Branches:
- **parent**: The parent base branch for dependency tracking (structural)
- **upstreamStrategy**: How changes flow to parent (process)
- **downstreamStrategy**: How updates flow from parent (process)
- **autoUpdate**: Whether the branch receives updates automatically on finish (process)

#### For Topic Branch Types (using gitpour.branch.* keys):
- **parent**: Default parent base branch (structural)
- **startPoint**: Branch to create from — can differ from parent (structural)
- **prefix**: Branch name prefix (structural)
- **upstreamStrategy**: How to merge back to parent on finish (process)
- **downstreamStrategy**: How to receive updates from parent (process)
- **tag**: Whether the branch type produces tags on finish (process)
- **tagPrefix**: Prefix for created tags (process)
- **deleteRemote**: Whether finishing also deletes the remote branch (process)
- **worktree**: Whether `start` gives every new branch its own worktree (process)

### Merge Strategies

git-pour-now supports multiple merge strategies:

- **merge**: Standard Git merge with merge commit
- **rebase**: Rebase changes onto target branch
- **squash**: Squash all commits into single commit
- **none**: No automatic merging

## Example Workflow Configurations

### 1. Simple GitHub Flow

Perfect for continuous deployment with hotfix capability:

```ini
[gitpour "branch.main"]
    type = base
    parent =
    upstreamStrategy = none
    downstreamStrategy = none

[gitpour "branch.feature"]
    type = topic
    parent = main
    startPoint = main
    upstreamStrategy = rebase
    downstreamStrategy = rebase

[gitpour "branch.hotfix"]
    type = topic
    parent = main
    startPoint = main
    upstreamStrategy = merge
    downstreamStrategy = none
```

**Branch Structure:**
```
main
├── feature/user-interface
├── feature/api-integration
└── hotfix/security-patch
```

### 2. Traditional Git-Flow

Classic git-pour with develop branch and release management:

```ini
[gitpour "branch.main"]
    type = base
    parent =
    upstreamStrategy = none
    downstreamStrategy = none

[gitpour "branch.develop"]
    type = base
    parent = main
    upstreamStrategy = merge
    downstreamStrategy = merge

[gitpour "branch.feature"]
    type = topic
    parent = develop
    startPoint = develop
    upstreamStrategy = rebase
    downstreamStrategy = merge

[gitpour "branch.release"]
    type = topic
    parent = main
    startPoint = develop
    upstreamStrategy = merge
    downstreamStrategy = none
    tag = true

[gitpour "branch.hotfix"]
    type = topic
    parent = main
    startPoint = main
    upstreamStrategy = merge
    downstreamStrategy = none
    tag = true
```

**Branch Structure:**
```
main
├── hotfix/critical-fix
├── release/v1.0
└── develop
     ├── feature/payment-gateway
     │    └── feature/card-processing
     ├── feature/user-authentication
     │    └── feature/two-factor-auth
```

### 3. Web Application Flow

Multi-environment deployment with staging and production:

```ini
[gitpour "branch.production"]
    type = base
    parent =
    upstreamStrategy = none
    downstreamStrategy = none

[gitpour "branch.staging"]
    type = base
    parent = production
    upstreamStrategy = merge
    downstreamStrategy = merge

[gitpour "branch.main"]
    type = base
    parent = staging
    upstreamStrategy = merge
    downstreamStrategy = merge

[gitpour "branch.feature"]
    type = topic
    parent = main
    startPoint = main
    upstreamStrategy = rebase
    downstreamStrategy = rebase

[gitpour "branch.hotfix"]
    type = topic
    parent = production
    startPoint = production
    upstreamStrategy = merge
    downstreamStrategy = none
```

**Branch Structure:**
```
production
├── hotfix/urgent-fix
└── staging
    └── main
        ├── feature/new-feature
        └── feature/ui-improvement
```

## Advanced Features

### Automatic Branch Updates

Configure branches to automatically receive updates from their parent:

```ini
[gitpour "branch.develop"]
    parent = main
    autoUpdate = true
    downstreamStrategy = merge
```

When `autoUpdate` is enabled, finishing a topic branch into `main` automatically propagates changes to `develop`.

### Cascade Updates

Changes can cascade through multiple levels of the dependency tree:

1. Finish `hotfix/security-patch` into `production`
2. Changes automatically flow to `staging` (if configured)
3. Changes then flow to `main` (if configured)
4. Finally cascade to `develop` (if configured)

### Tag Creation

Automatic tag creation with configurable naming:

```ini
[gitpour "branch.release"]
    tag = true
    tagPrefix = "v"
```

### Child Branch Updates

When finishing a topic branch, git-pour-now automatically updates child base branches that depend on the target parent branch. This ensures consistency across the branch hierarchy:

```bash
# Finishing a hotfix into main automatically updates develop
git pour hotfix finish security-patch

# The system will:
# 1. Merge hotfix/security-patch into main
# 2. Automatically update develop from main (if configured)
# 3. Update any other child branches of main
```

Configure automatic updates in base branch settings:

```ini
[gitpour "branch.develop"]
    parent = main
    autoUpdate = true
    downstreamStrategy = merge
```

## Command Structure

### Core Commands

```bash
# Initialize git-pour configuration
git pour init

# Topic branch operations
git pour <type> start <name>
git pour <type> finish <name>
git pour <type> list

# Base branch operations
git pour integrate <branch>        # Integrate a base branch into its parent
git pour update <branch>
git pour update <branch> --rebase  # Force rebase strategy

# Shorthand commands (auto-detect branch type)
git pour rebase                    # Shorthand for: git pour <type> update --rebase
git pour update                    # Shorthand for: git pour <type> update
git pour finish                    # Shorthand for: git pour <type> finish

# Worktree operations (addressed by full branch name)
git pour worktree add <branch>
git pour worktree remove <branch>
git pour worktree list
git pour worktree prune
git pour worktree path <branch>

# Overview
git pour overview
```

### Dynamically Generated Commands

Topic-branch commands are generated per branch type from configuration. Each
configured type (feature, release, hotfix, support, and any custom type) gets
the same set of subcommands:

```bash
git pour feature start <name>    # start a feature branch
git pour hotfix finish <name>    # finish a hotfix branch
git pour release list            # list release branches
```

### Shorthand Commands

git-pour-now provides convenient shorthand commands that automatically detect your current topic branch:

```bash
git pour rebase                   # → git pour <type> update --rebase
git pour update                   # → git pour <type> update
git pour finish                   # → git pour <type> finish
git pour delete                   # → git pour <type> delete
git pour rename <name>            # → git pour <type> rename <name>
```

## Command Implementation

### Command Structure Overview

Commands in git-pour-now are implemented using the Cobra CLI framework with a clear architectural pattern:

1. **Root Command** (`cmd/root.go`): Sets up the main CLI structure and global flags
2. **Dynamic Registration** (`cmd/topicbranch.go`): Automatically registers branch type commands based on configuration
3. **Individual Commands** (`cmd/*.go`): Each major operation has its own file with specific logic

All commands follow a consistent pattern: validate inputs, load configuration, execute Git operations, and handle errors gracefully.

### The Finish Command: A Deep Dive

The finish command (`cmd/finish.go`) is the most complex command in the system, demonstrating the sophisticated architecture used throughout git-pour-now.

#### Step-Based State Machine

The finish command uses a **step-based state machine** approach to handle complex multi-step operations that can be interrupted by merge conflicts:

```
Steps: merge → create_tag → update_children → delete_branch
```

This architecture allows the command to:
- Resume operations after conflict resolution
- Provide clear progress feedback
- Handle complex branching scenarios
- Maintain consistency across interruptions

State is persisted to disk (`mergestate.MergeState`) so the operation can resume after conflict resolution via `--continue` or be cancelled with `--abort`. Each step has a dedicated handler, and child branch updates respect individual downstream strategies (see [Advanced Features > Child Branch Updates](#child-branch-updates) above).

For implementation details—struct definitions, handler functions, and code examples—see [CODE_REFERENCE.md](CODE_REFERENCE.md#state-machine-finish-command).

## Integration Points

### Tower Integration

git-pour-now integrates seamlessly with [Tower](https://www.git-tower.com), providing graphical workflow management while using the same configuration system.

### CI/CD Integration

The flexible configuration system enables easy integration with modern CI/CD pipelines by supporting:

- Webhook-triggered deployments based on branch patterns
- Environment-specific deployment strategies
- Automatic tag-based releases

## Migration from git-flow-avh

git-pour-now maintains compatibility with existing git-flow-avh configurations while providing migration tools for enhanced features:

```bash
# Confirm import of existing gitflow configuration, then optionally remove local legacy keys
# (when AVH config is present and no explicit options are given)
git pour init
```

## Extensibility

The unified topic branch implementation and configuration-driven approach make git-pour-now highly extensible:

- Add custom branch types through configuration
- Define organization-specific workflow templates
- Create custom merge strategies through hooks
- Extend functionality through plugin architecture

This technical foundation enables teams to implement any branching strategy while maintaining the automation and convenience that made git-pour popular.
