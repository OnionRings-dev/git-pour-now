# Configuration Reference

The default git-pour-now workflow uses `main`, `feature/`, and temporary
`prerelease/<version>` branches. Stable tags are created on main; prerelease tags
on the matching prerelease branch. No develop branch is created. Examples below
using develop/release/hotfix belong to the explicit `--preset=classic` workflow.
The test suite requests that preset explicitly for inherited GitFlow scenarios.


This document provides a comprehensive reference for all git-pour-now configuration options. All configuration is stored in Git configuration using the `gitpour.*` key namespace.

## Configuration Overview

git-pour-now uses a hierarchical configuration system with **three levels of precedence**:

1. **Branch Type Definition** (`gitpour.branch.*`) - The identity and process characteristics of a branch type
2. **Command-Specific Configuration** (`gitpour.<branchtype>.*`) - How specific commands behave for a branch type
3. **Command-Line Flags** - **Highest priority** - one-off overrides that always win

### Understanding the Layers

**Layer 1 defines what a branch type *is*** — its role in the workflow, not tunable defaults. These properties describe the branch type's identity (where it lives in the hierarchy, its naming convention) and its process characteristics (how it participates in the workflow). For example, `tag=true` on a release branch means "releases are the kind of branch that produces tags on finish" — it describes the release process, not merely a default that you toggle.

**Layer 2 controls how commands execute** — operational settings like whether to fetch before an operation, sign tags, or keep branches after finishing. These are genuine behavioral knobs that adjust command execution without changing what the branch type fundamentally is.

**Layer 3 provides one-off overrides** via CLI flags for ad-hoc situations.

Many command options intentionally have no Layer 1 equivalent — they exist only in Layer 2 and Layer 3. For example, publish `push-options` are purely operational and don't define any branch type characteristic.

> **For developers adding new options**: Most new options need only Layer 2 + Layer 3. Add Layer 1 support only when the option defines a branch type characteristic. See [CODING_GUIDELINES.md](CODING_GUIDELINES.md) → "Adding New Configuration Options" for the full decision guide and checklist.

## Core System Configuration

### System Settings

| Key | Description | Default | Example |
|-----|-------------|---------|---------|
| `gitpour.version` | Internal version marker for compatibility | `1.0` | `1.0` |
| `gitpour.initialized` | Marks repository as git-pour initialized | `false` | `true` |
| `gitpour.origin` | Remote name to use for operations | `origin` | `upstream` |
| `gitpour.remote` | Alias for `gitpour.origin` | `origin` | `upstream` |
| `gitpour.worktreePath` | Path template for a branch's worktree. Supports `{{ repo }}`, `{{ branch }}`, `{{ branchName }}`, `{{ topicType }}` (spaced or unspaced), expands a leading `~`, resolves a relative template against the main worktree root, and always yields an absolute path | `../{{ repo }}-worktrees/{{ branch }}` | `~/worktrees/{{ topicType }}/{{ branchName }}` |

`gitpour.worktree.<branch>.managed` is **not** a setting: it is repository-local state git-pour writes when it creates a worktree, so later commands can tell its own worktrees from the user's. It is carried to the new name when the branch is renamed, and cleared when git-pour removes the worktree. It is excluded from the shared-config set and never reaches a committed `.gitpour`. See [git-pour-worktree(1)](docs/git-pour-worktree.1.md).

`gitpour.branch.<branch>.base` is likewise state rather than a setting: it records the start point a topic branch was created from. It is carried through a rename alongside the marker, and cleared when the branch is finished or deleted. See [git-pour-rename(1)](docs/git-pour-rename.1.md).

## Branch Type Configuration (Layer 1)

Branch type configuration defines the **identity and process characteristics** of each branch type using the pattern:
`gitpour.branch.<branchname>.<property>`

These properties describe *what the branch type is* — its structural role and how it participates in the workflow. They should only contain essential branch-type-relevant configuration.

Branch names (the `<branchname>` subsection) are matched **case-insensitively**, mirroring Git's `core.ignorecase`. The case used when a branch is first added is preserved as its canonical form and written back verbatim; git-pour always operates on that canonical name (including the Git ref, so it is correct on case-sensitive filesystems). Any later reference — a parent, a `--starting-point`, or the target of an `edit`/`rename`/`delete` — may be given in any case and resolves to the single canonical entry. Adding or renaming to a name that differs only in case from a *different* existing branch is rejected as "already exists"; a case-only rename of a branch to itself is allowed and re-canonicalizes it. Two branch types differing only in case cannot coexist. Branch names may contain dots (e.g., `custom.main`, `V10.5`): the property is always the final dot-separated segment, so the full dotted name round-trips through the config reader. (Names Git itself rejects as refs — a leading dot, consecutive dots, a trailing `.lock` — remain invalid. `config add`/`rename` additionally reject names beginning with `-`, which Git would accept as a ref but git-pour passes as a bare operand to git commands.)

### Structural Properties

Define where the branch type fits in the hierarchy:

| Property | Description | Values | Default |
|----------|-------------|--------|---------|
| `type` | Branch type classification | `base`, `topic` | Required |
| `parent` | Parent branch for this branch type | Branch name | Varies by type |
| `startPoint` | Branch to start new branches from | Branch name | Same as parent |
| `prefix` | Branch name prefix (topic only) | String | `<type>/` |

The `prefix` value is a literal string prepended to the short branch name, used exactly as configured. The trailing slash in the defaults is a naming convention, not a requirement: with `prefix = feature_`, `git pour feature start login` creates the branch `feature_login`.

### Process Characteristics

Define how the branch type participates in the workflow:

| Property | Description | Values | Default |
|----------|-------------|--------|---------|
| `upstreamStrategy` | How changes flow TO parent on finish | `merge`, `rebase`, `squash` | `merge` |
| `downstreamStrategy` | How updates flow FROM parent | `merge`, `rebase` | `merge` |
| `tag` | Branch type produces tags on finish (topic only) | `true`, `false` | `false` |
| `tagprefix` | Prefix for created tags (topic only) | String | `""` |
| `autoUpdate` | Auto-update from parent on finish (base only) | `true`, `false` | `false` |
| `deleteRemote` | Delete remote branch on finish (topic only) | `true`, `false` | `false` |
| `worktree` | `start` creates a worktree for new branches (topic only) | `true`, `false` | `false` |

`worktree` says the branch type's process gives every new branch its own working directory: `git pour <type> start` then creates the branch in a worktree at the path [`gitpour.worktreePath`](#system-settings) computes, instead of checking it out. `--worktree` and `--no-worktree` override it per invocation. git-pour writes the key for topic branch types only — a worktree default on a base branch has no effect.

For example, setting `tag=true` on release branches means "releases produce tags" — it characterizes the release process. This can still be overridden per-command (Layer 2: `gitpour.release.finish.notag`) or per-invocation (Layer 3: `--notag`), but the branch config establishes the branch type's intended role.

### Default Branch Types

#### Base Branches

```bash
# Main branch (production releases)
gitpour.branch.main.type=base
gitpour.branch.main.parent=             # No parent (root)
gitpour.branch.main.upstreamStrategy=merge
gitpour.branch.main.downstreamStrategy=merge
gitpour.branch.main.autoUpdate=false

# Develop branch (integration)
gitpour.branch.develop.type=base
gitpour.branch.develop.parent=main
gitpour.branch.develop.startPoint=main
gitpour.branch.develop.upstreamStrategy=merge
gitpour.branch.develop.downstreamStrategy=merge
gitpour.branch.develop.autoUpdate=true    # Auto-updates from main
```

Base branches carry no `worktree` key: `start` never creates a worktree for them, so a default there would have no effect and git-pour does not write one.

#### Topic Branches

```bash
# Feature branches
gitpour.branch.feature.type=topic
gitpour.branch.feature.parent=develop
gitpour.branch.feature.startPoint=develop
gitpour.branch.feature.prefix=feature/
gitpour.branch.feature.upstreamStrategy=merge
gitpour.branch.feature.downstreamStrategy=rebase
gitpour.branch.feature.tag=false
gitpour.branch.feature.autoUpdate=false
gitpour.branch.feature.worktree=false

# Release branches
gitpour.branch.release.type=topic
gitpour.branch.release.parent=main
gitpour.branch.release.startPoint=develop
gitpour.branch.release.prefix=release/
gitpour.branch.release.upstreamStrategy=merge
gitpour.branch.release.downstreamStrategy=merge
gitpour.branch.release.tag=true
gitpour.branch.release.autoUpdate=false
gitpour.branch.release.worktree=false

# Hotfix branches
gitpour.branch.hotfix.type=topic
gitpour.branch.hotfix.parent=main
gitpour.branch.hotfix.startPoint=main
gitpour.branch.hotfix.prefix=hotfix/
gitpour.branch.hotfix.upstreamStrategy=merge
gitpour.branch.hotfix.downstreamStrategy=merge
gitpour.branch.hotfix.tag=true
gitpour.branch.hotfix.autoUpdate=false
gitpour.branch.hotfix.worktree=false
```

#### Custom Branch Types

You can define custom branch types by setting the appropriate configuration:

```bash
# Example: Support branches
gitpour.branch.support.type=topic
gitpour.branch.support.parent=main
gitpour.branch.support.startPoint=main
gitpour.branch.support.prefix=support/
gitpour.branch.support.tag=true
gitpour.branch.support.tagprefix=support-
```

## Command-Specific Configuration (Layer 2)

Command-specific configuration controls **how commands execute** for a branch type, using the pattern:
`gitpour.<branchtype>.<command>.<option>`

These are operational settings — they adjust command behavior without changing the branch type's identity. Some of these can override Layer 1 process characteristics (e.g., `notag` overrides the branch type's `tag` setting), while others exist only at this layer (e.g., `sign`, `keep`, `fetch`).

### Finish Command Options

The finish command supports per-branch-type configuration:

| Option | Description | Values | Default |
|--------|-------------|--------|---------|
| `merge` | Override merge strategy | `merge`, `rebase`, `squash` | From branch config |
| `no-ff` | Always create a merge commit | `true`, `false` | `false` |
| `ff` | Allow a fast-forward merge when possible | `true`, `false` | `true` |
| `ff-only` | Require the merge into the parent to be a fast-forward, aborting otherwise (wins over `no-ff` and `ff` within this layer) | `true`, `false` | `false` |
| `notag` | Disable tag creation | `true`, `false` | Opposite of branch `tag` |
| `sign` | Sign created tags | `true`, `false` | `false` |
| `signingkey` | GPG key for signing | Key ID | Git default |
| `messagefile` | File containing tag message | File path | None |
| `keep` | Keep branch after finish | `true`, `false` | `false` |
| `keepremote` | Keep remote branch | `true`, `false` | `false` |
| `keeplocal` | Keep local branch | `true`, `false` | `false` |
| `force-delete` | Force delete branch | `true`, `false` | `false` |
| `fetch` | Fetch before operation (`gitpour.<type>.<cmd>.fetch`; default → config → `--fetch`/`--no-fetch`) | `true`, `false` | `true` for both `finish` and `start` |
| `push` | Push finished branches to remote | `true`, `false` | `false` |
| `pushtag` | Push created tag to remote | `true`, `false` | Follows `push` |

#### Examples

```bash
# Override feature finish to use rebase
gitpour.feature.finish.rebase=true

# Disable tags for feature branches
gitpour.feature.finish.notag=true

# Sign all release tags
gitpour.release.finish.sign=true
gitpour.release.finish.signingkey=ABC123DEF

# Keep hotfix branches locally after finish
gitpour.hotfix.finish.keeplocal=true

# Push finished release branches and the created tag to the remote
gitpour.release.finish.push=true

# Push the created tag even when branches are not pushed
gitpour.release.finish.pushtag=true
```

When `push` is enabled, finish pushes the target/parent branch and each
auto-updated child base branch to the remote (`gitpour.remote`, default
`origin`) as a final stage, then the created tag. `pushtag` follows the
resolved `push` value unless set explicitly, letting you decouple the tag push
from the branch push. A missing remote skips the push with a note (finish still
succeeds); a rejected non-fast-forward push is a hard error, leaving the
already-completed local finish intact.

### Delete Command Options

| Option | Description | Values | Default |
|--------|-------------|--------|---------|
| `fetch` | Fetch (and fast-forward the parent when it is checked out) before delete | `true`, `false` | `false` |

#### Examples

```bash
# Always fetch before deleting feature branches
gitpour.feature.delete.fetch=true
```

When `fetch` is enabled, delete fetches from the remote so Git can detect a branch that was merged remotely (e.g., via a GitHub PR merge). The parent is fast-forwarded from its remote only when it is the branch currently checked out (`git merge --ff-only` acts on HEAD); delete auto-checks-out the parent when you delete the branch you are on, which is the common case. A fetch failure against an unreachable remote is a non-fatal note — deletion is not blocked by the fetch itself, but the topic sync check still runs against existing local tracking data and can still abort (behind/diverged) unless `--force` is given. `--no-fetch` skips only the fetch, not the sync check.

### Update Command Options

| Option | Description | Values | Default |
|--------|-------------|--------|---------|
| `downstreamStrategy` | Override update strategy | `merge`, `rebase` | From branch config |

#### Examples

```bash
# Force rebase for feature updates
gitpour.branch.feature.downstreamStrategy=rebase

# Use merge for release updates
gitpour.branch.release.downstreamStrategy=merge
```

## Configuration Precedence

git-pour-now follows a strict three-layer precedence hierarchy:

### Layer 1: Branch Type Definition (Lowest Priority)
Defines the branch type's identity and process characteristics via `gitpour.branch.<branchtype>.*`
```bash
gitpour.branch.release.tag=true    # "Releases produce tags" — a process characteristic
```

### Layer 2: Command-Specific Configuration (Medium Priority)
Operational settings that control how commands execute via `gitpour.<branchtype>.<command>.*`
```bash
gitpour.release.finish.notag=true  # Override: skip tagging for this command
```

### Layer 3: Command-Line Flags (Highest Priority)
One-off overrides that **always win**
```bash
git pour release finish v1.0 --tag  # Force tagging for this invocation
```

Not every option spans all three layers. Layer 1 is reserved for essential branch-type properties — many command options exist only at Layer 2 and Layer 3. For example, `sign`, `keep`, and publish `push-options` are purely operational and have no Layer 1 equivalent.

### Boolean Values

Options listed with `Values: true, false` accept every git-config boolean spelling, at both Layer 1 and Layer 2: `true`, `yes`, `on` and any non-zero integer are true; `false`, `no`, `off`, `0` and an empty value are false; matching is case-insensitive. An unrecognized value is treated as `false`. Note that inverse options such as `no-rebase` and `notag` are separate keys rather than falsy spellings of the positive option. See the **Boolean Values** section of `docs/gitpour-config.5.md` for the full rule.

## Shared Configuration (.gitpour)

By default git-pour configuration lives in `.git/config`, which is never committed, so every clone must run `git pour init`. A committable `.gitpour` file at the repository root lets a team share one configuration.

git-pour does not read `.gitpour` directly during operations — it **copies** the `gitpour.*` keys into the clone's local `.git/config`, so every read site sees the shared values and native git precedence still applies.

### Creating and using .gitpour

```bash
# Author .gitpour and copy it into local config, then commit it
git pour init --defaults --shared
git add .gitpour && git commit -m "Add shared git-pour configuration"

# On a fresh clone, activate the committed configuration
git pour config sync     # or set gitpour.shared.autoInit=true to activate automatically

# Keep local config in step with .gitpour
git pour config status   # exit 0 in sync, exit 6 on drift
git pour config sync     # re-copy .gitpour into local config

# Edit the shared file (and re-sync local) instead of local-only config
git pour config edit topic feature --shared --prefix feat/
git pour config add topic qa develop --shared
```

### Shared-managed keys and control keys

The copy moves the whole `gitpour.*` namespace **except** the control keys `gitpour.shared.*` and per-branch runtime metadata `gitpour.branch.<branch>.base`. Keys outside `gitpour.*` (e.g. `core.*`, `alias.*`) are never copied, and multi-value keys keep their order.

```bash
# Activate a committed .gitpour automatically on your clones (read from global/local, never copied)
git config --global gitpour.shared.autoInit true

# Allow copying the executable hook path key from a trusted .gitpour
git config gitpour.shared.trustHooks true
```

When `gitpour.shared.trustHooks` is unset, the `gitpour.path.hooks` key is not copied from `.gitpour` (and is removed from local on the next sync if it was previously copied); non-interactive auto-init refuses when an untrusted shared hook path is present. See `docs/gitpour-config.5.md` for the full first-run activation and hook-trust semantics.

## Legacy Compatibility

### Git-Flow-AVH Compatibility

git-pour-now offers to import git-flow-avh configuration during `git pour init`, after explicit confirmation:

| AVH Key | git-pour-now Equivalent | Description |
|---------|--------------------------|-------------|
| `gitflow.branch.master` | `gitpour.branch.main` | Main branch name |
| `gitflow.branch.develop` | `gitpour.branch.develop` | Develop branch name |
| `gitflow.prefix.feature` | `gitpour.branch.feature.prefix` | Feature prefix |
| `gitflow.prefix.release` | `gitpour.branch.release.prefix` | Release prefix |
| `gitflow.prefix.hotfix` | `gitpour.branch.hotfix.prefix` | Hotfix prefix |
| `gitflow.prefix.support` | `gitpour.branch.support.prefix` | Support prefix |
| `gitflow.prefix.versiontag` | `gitpour.branch.*.tagprefix` | Tag prefix |

### Migration from git-flow-avh

1. **Confirmed Import**: Run `git pour init` and answer yes to import existing configuration, including operational preferences. After a successful import, a separate prompt offers to remove the previous repository-local `gitflow.*` keys. Both prompts default to no; global/system and included settings are preserved.
2. **Manual Migration**: Use the mapping table above to convert AVH config to git-pour-now format
3. **Verification**: Run `git pour overview` to verify imported configuration. Ordinary commands do not import legacy settings; explicit initialization options such as `--defaults` skip migration.

## Configuration Examples

### Simple GitHub Flow

```bash
# Minimal configuration - feature branches only
git config gitpour.branch.main.type base
git config gitpour.branch.feature.type topic
git config gitpour.branch.feature.parent main
git config gitpour.branch.feature.prefix feature/
```

### Traditional Git-Flow

```bash
# Full git-pour with main, develop, and all topic types
git config gitpour.branch.main.type base
git config gitpour.branch.develop.type base
git config gitpour.branch.develop.parent main
git config gitpour.branch.develop.autoUpdate true

git config gitpour.branch.feature.type topic
git config gitpour.branch.feature.parent develop
git config gitpour.branch.feature.prefix feature/

git config gitpour.branch.release.type topic
git config gitpour.branch.release.parent main
git config gitpour.branch.release.startPoint develop
git config gitpour.branch.release.tag true

git config gitpour.branch.hotfix.type topic
git config gitpour.branch.hotfix.parent main
git config gitpour.branch.hotfix.tag true
```

### Custom Workflow

```bash
# Custom branch types with specific merge strategies
git config gitpour.branch.bugfix.type topic
git config gitpour.branch.bugfix.parent develop
git config gitpour.branch.bugfix.prefix bugfix/
git config gitpour.branch.bugfix.upstreamStrategy squash

# Custom finish behavior
git config gitpour.bugfix.finish.notag true
git config gitpour.bugfix.finish.keep true
```

## Configuration Management

### Viewing Configuration

```bash
# View all git-pour configuration
git config --get-regexp '^gitpour\.'

# View branch-specific configuration
git config --get-regexp '^gitpour\.branch\.feature\.'

# View command-specific configuration
git config --get-regexp '^gitpour\.feature\.finish\.'
```

### Setting Configuration

```bash
# Set branch type configuration
git config gitpour.branch.feature.upstreamStrategy rebase

# Set command-specific configuration
git config gitpour.feature.finish.squash true

# Set global configuration (all repositories)
git config --global gitpour.feature.finish.notag true
```

### Removing Configuration

```bash
# Remove specific setting
git config --unset gitpour.feature.finish.rebase

# Remove all git-pour configuration (reset)
git config --remove-section gitpour
```

## Best Practices

1. **Start with defaults** - Use `git pour init --defaults` for standard configuration
2. **Customize gradually** - Override only what you need to change
3. **Document team conventions** - Use consistent configuration across team repositories
4. **Use command precedence** - Rely on command-line flags for one-off operations
5. **Test configuration** - Use `git pour overview` to verify your setup
6. **Backup configuration** - Export configuration before major changes:
   ```bash
   git config --get-regexp '^gitpour\.' > gitpour-backup.txt
   ```
