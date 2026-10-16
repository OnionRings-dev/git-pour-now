# Git-Flow-AVH Compatibility Analysis

Examples using `develop`, release branches, or hotfix branches refer to the
explicit `--preset=classic` compatibility workflow. The default workflow uses
`main`, `feature/`, and temporary `prerelease/<version>` branches, with releases
identified by SemVer tags.

This document provides a comprehensive compatibility analysis between git-flow-avh and git-pour-now, covering both configuration options and command implementations.

## Configuration Format Comparison

### git-flow-avh Format
```bash
# Branch names (legacy)
gitpour.branch.master
gitpour.branch.develop

# Branch prefixes (legacy)
gitpour.prefix.feature
gitpour.prefix.release
gitpour.prefix.hotfix
gitpour.prefix.bugfix
gitpour.prefix.support
gitpour.prefix.versiontag

# Command options (matches our Layer 2)
gitpour.<branchtype>.<command>.<option>
```

### git-pour-now Format
```bash
# System settings
gitpour.origin
gitpour.version
gitpour.initialized

# Branch type configuration (our Layer 1)
gitpour.branch.<branchtype>.type
gitpour.branch.<branchtype>.parent
gitpour.branch.<branchtype>.prefix
gitpour.branch.<branchtype>.tag
gitpour.branch.<branchtype>.tagprefix
gitpour.branch.<branchtype>.upstreamStrategy
gitpour.branch.<branchtype>.downstreamStrategy

# Command-specific overrides (our Layer 2)
gitpour.<branchtype>.<command>.<option>
```

## Supported Options Analysis

### ✅ Fully Compatible Options
These git-flow-avh options are fully supported in git-pour-now:

| AVH Option | git-pour-now Equivalent | Status |
|------------|--------------------------|---------|
| `gitpour.origin` | `gitpour.origin` | ✅ Direct match |
| `gitpour.feature.start.fetch` | `gitpour.feature.start.fetch` | ✅ Direct match |
| `gitpour.feature.finish.fetch` | `gitpour.feature.finish.fetch` | ✅ Direct match |
| `gitpour.release.finish.sign` | `gitpour.release.finish.sign` | ✅ Direct match |
| `gitpour.release.finish.signingkey` | `gitpour.release.finish.signingkey` | ✅ Direct match |
| `gitpour.hotfix.finish.sign` | `gitpour.hotfix.finish.sign` | ✅ Direct match |
| `gitpour.release.finish.keep` | `gitpour.release.finish.keep` | ✅ Direct match |
| `gitpour.feature.finish.keep` | `gitpour.feature.finish.keep` | ✅ Direct match |
| `gitpour.*.finish.no-ff` | `gitpour.*.finish.no-ff` | ✅ Direct match |
| `gitpour.*.finish.preserve-merges` | `gitpour.*.finish.preserve-merges` | ✅ Direct match |

### 🔄 Requires Translation (Auto-Import)
These AVH options need translation but are automatically handled during migration:

| AVH Option | git-pour-now Translation | Notes |
|------------|--------------------------|-------|
| `gitpour.branch.master` | `gitpour.branch.main` | Branch name modernization |
| `gitpour.prefix.feature` | `gitpour.branch.feature.prefix` | New hierarchical format |
| `gitpour.prefix.release` | `gitpour.branch.release.prefix` | New hierarchical format |
| `gitpour.prefix.hotfix` | `gitpour.branch.hotfix.prefix` | New hierarchical format |
| `gitpour.prefix.support` | `gitpour.branch.support.prefix` | New hierarchical format |
| `gitpour.prefix.versiontag` | `gitpour.branch.*.tagprefix` | Applied to all tag-creating types |

### ✅ Merge Strategy & Tag Flags (Direct Match)
These AVH boolean flags are read directly by git-pour-now's resolver at runtime (Layer 2) — no translation or conversion is needed. There is no `gitpour.<type>.finish.merge` key; the finish strategy is resolved from `gitpour.branch.<name>.upstreamStrategy` plus these `finish.rebase`/`finish.squash` booleans:

| AVH Option | git-pour-now | Notes |
|------------|---------------|-------|
| `gitpour.feature.finish.rebase=true` | `gitpour.feature.finish.rebase=true` | ✅ Direct match — resolver reads it as-is |
| `gitpour.feature.finish.squash=true` | `gitpour.feature.finish.squash=true` | ✅ Direct match — resolver reads it as-is |
| `gitpour.release.finish.squash=true` | `gitpour.release.finish.squash=true` | ✅ Direct match — resolver reads it as-is |
| `gitpour.release.finish.notag=true` | `gitpour.release.finish.notag=true` | ✅ Direct match — read at runtime (equivalently `gitpour.branch.release.tag=false`) |
| `gitpour.hotfix.finish.notag=true` | `gitpour.hotfix.finish.notag=true` | ✅ Direct match — read at runtime (equivalently `gitpour.branch.hotfix.tag=false`) |

### ⚠️ Missing Options in git-pour-now
These git-flow-avh options are **not currently supported** in git-pour-now:

| Missing Option | Description | Impact |
|----------------|-------------|--------|
| `gitpour.allowdirty` | Allow operations with dirty working tree | 🔴 High |
| `gitpour.*.finish.pushproduction` | Push to production branch after finish | 🟢 Covered by `gitpour.*.finish.push` |
| `gitpour.*.finish.pushdevelop` | Push to develop branch after finish | 🟢 Covered by `gitpour.*.finish.push` |
| `gitpour.*.finish.nobackmerge` | Skip back-merge to develop | 🟡 Medium |
| `gitpour.*.finish.ff-master` | Fast-forward merge to master | 🟡 Medium |

**Note on push-on-finish**: git-pour-now implements push-on-finish generically via `gitpour.*.finish.push` (and `gitpour.*.finish.pushtag`) plus the `--push`/`--no-push` and `--pushtag`/`--no-pushtag` flags. A single `--push` pushes the target branch **and** every auto-updated child base branch (plus the created tag), so avh's separate `--pushproduction`/`--pushdevelop` knobs are unnecessary — the child set is derived from the branch topology rather than hard-coded to production/develop.

### 🆕 git-pour-now Exclusive Features
These features are **only available** in git-pour-now:

| Feature | Description | Benefit |
|---------|-------------|---------|
| Branch type system | `gitpour.branch.*.type=base\|topic` | Flexible branch taxonomy |
| Parent relationships | `gitpour.branch.develop.parent=main` | Automatic change propagation |
| Auto-update | `gitpour.branch.develop.autoUpdate=true` | Keeps branches synchronized |
| Downstream strategies | `gitpour.*.downstreamStrategy` | Control update merge behavior |
| Universal branch properties | Common properties across all types | Consistent configuration |
| Custom branch types | User-defined branch types | Workflow customization |
| Config management | `git pour config` subcommands | Add, edit, rename, delete branch types |

## Migration Requirements

### Automatic Import (Implemented)
✅ Our `git pour init` automatically handles these conversions:
- `gitpour.branch.master` → `gitpour.branch.main`
- `gitpour.prefix.*` → `gitpour.branch.*.prefix`
- `gitpour.prefix.versiontag` → `gitpour.branch.*.tagprefix`
- Command options are preserved as-is

### Runtime Compatibility (No Migration Needed)
🔄 These AVH options work at runtime without conversion:
- `gitpour.*.finish.notag` — handled by the config resolver at Layer 2

### Manual Translation Required
🔄 These require user intervention or enhanced import logic:
- Boolean strategy flags (`rebase=true`, `squash=true`) → strategy configuration

### Missing Feature Implementation
❌ These require new feature development:

#### 1. Allow Dirty Working Tree
```bash
# AVH
gitpour.allowdirty=true

# Needed in git-pour-now
gitpour.allowdirty=true  # Add system-level support
```

#### 2. Auto-Push on Finish
```bash
# AVH
gitpour.release.finish.push=true

# git-pour-now (supported)
gitpour.release.finish.push=true     # Push target + auto-updated children + tag
gitpour.release.finish.pushtag=false # Optional: decouple the tag from the branch push
```

---

## Command Compatibility Analysis

### Command Status Overview

| Command | AVH | git-pour-now | Status |
|---------|-----|---------------|--------|
| `init` | ✅ | ✅ | ✅ Full parity (`-d/--defaults`, `-f/--force`) |
| `start` | ✅ | ✅ | ✅ Full parity (`-F/--fetch`, `<base>` parameter) |
| `finish` | ✅ | ✅ | 🟡 Most flags supported (see details) |
| `publish` | ✅ | ✅ | ✅ Fully implemented |
| `track` | ✅ | ✅ | ✅ Fully implemented |
| `list` | ✅ | ✅ | 🟡 Missing `-v/--verbose` |
| `checkout` | ✅ | ✅ | ✅ Complete |
| `delete` | ✅ | ✅ | ✅ Complete parity |
| `rename` | ✅ | ✅ | ✅ Complete |
| `diff` | ✅ | ❌ | ❌ Not implemented |
| `pull` | ✅ | ❌ | ❌ Not implemented |
| `config` | ✅ | ✅ | 🟡 Different approach (add/edit/rename/delete vs set/base) |

### General Shorthand Commands
git-pour-now implements shorthand commands that work on the current branch:
- ✅ **`git pour finish`** — Finish current branch
- ✅ **`git pour delete`** — Delete current branch
- ✅ **`git pour rebase`** — Rebase current branch (via `update --rebase`)
- ✅ **`git pour update`** — Update current branch
- ✅ **`git pour rename`** — Rename current branch
- ✅ **`git pour publish`** — Publish current branch

### Command-by-Command Option Analysis

#### ✅ `init` Command
| Flag | AVH | git-pour-now |
|------|-----|---------------|
| `-d/--defaults` | ✅ | ✅ |
| `-f/--force` | ✅ | ✅ |
| `--no-create-branches` | ❌ | ✅ (git-pour-now exclusive) |
| `--local/--global/--system/--file` | ✅ | ❌ |

#### ✅ `start` Command
| Flag | AVH | git-pour-now |
|------|-----|---------------|
| `-F/--fetch` | ✅ | ✅ |
| `<base>` parameter | ✅ | ✅ |
| `--showcommands` | ✅ | ❌ |

#### 🟡 `finish` Command
| Flag | AVH | git-pour-now |
|------|-----|---------------|
| `-F/--fetch` | ✅ | ✅ |
| `-r/--rebase` | ✅ | ✅ (via merge strategy) |
| `-S/--squash` | ✅ | ✅ (via merge strategy) |
| `-k/--keep` | ✅ | ✅ |
| `-D/--force-delete` | ✅ | ✅ |
| `--no-ff` / `--ff` | ✅ | ✅ |
| `--ff-only` | ❌ | ✅ (git-pour-now exclusive; avh has `--[no]ff` but no fast-forward precondition) |
| `-p/--preserve-merges` | ✅ | ✅ |
| `-s/--sign` | ✅ | ✅ |
| `--signingkey` | ✅ | ✅ |
| `-T/--tagname` | ✅ | ✅ |
| `-m/--message` | ✅ | ✅ |
| `--notag` | ✅ | ✅ |
| `--continue/--abort` | ❌ | ✅ (git-pour-now exclusive) |
| `--push` / `--no-push` | ✅ | ✅ (pushes target + auto-updated children + tag) |
| `--pushtag` / `--no-pushtag` | ❌ | ✅ (git-pour-now exclusive) |
| `--pushproduction` | ✅ | 🟢 Covered by `--push` |
| `--pushdevelop` | ✅ | 🟢 Covered by `--push` |
| `-b/--nobackmerge` | ✅ | ❌ |
| `--ff-master` | ✅ | ❌ |
| `--showcommands` | ✅ | ❌ |

**Landing branch parity**: like AVH, a completed `release` or `hotfix` finish leaves you on `develop` rather than on `main`. git-pour-now derives that branch from the configured topology — the last auto-update child of the finished branch's parent, or the parent when it has none — instead of hardcoding a branch name, so the same rule holds for custom topologies where no branch is called `develop`.

#### Other Commands
| Command | Missing from git-pour-now |
|---------|---------------------------|
| `list` | `-v/--verbose` |
| `checkout` | `--showcommands` (partial — present but not on all commands) |
| `delete` | ✅ Complete parity |
| `rename` | ✅ Complete parity |

---

## Compatibility Scores

### Configuration Compatibility: 85%

- **Core Features**: 100% compatible (branch management, merge strategies)
- **Configuration Translation**: 90% automatic (prefix mapping, branch names)
- **Command Options**: 75% compatible (most finish flags supported, including push-on-finish)
- **Advanced Features**: 60% compatible (preserve-merges and no-ff now supported)

### Command Compatibility: 80%

- **Core Commands**: 95% compatible (start, finish, list, delete, checkout, rename)
- **Shorthand Commands**: 100% compatible (all shorthands implemented)
- **Remote Operations**: 70% compatible (publish and track implemented; pull missing)
- **Advanced Operations**: 40% compatible (diff and full config parity still missing)

### Combined Compatibility Score: 82%

- **Core Workflow**: 95% compatible (start → work → finish → delete)
- **Team Collaboration**: 75% compatible (publish, track, and push-on-finish work)
- **Configuration**: 85% compatible
- **Commands**: 80% compatible

---

## Remaining Implementation Roadmap

### Phase 1: Critical Gaps
**Configuration:**
1. **`gitpour.allowdirty`** — Allow dirty working tree operations

**Commands:**
1. **`diff` command** — Show changes compared to parent branch

### Phase 2: Important Workflow Features
**Configuration:**
1. **Enhanced strategy translation** — Auto-convert boolean flags to strategies during import

**Commands:**
1. **`--verbose` on list** — Detailed branch listing
2. **Universal `--showcommands` flag** — Debug/learning aid across all commands

### Phase 3: Advanced Features
**Commands:**
1. **`pull` command** — Pull from remote
2. **`--nobackmerge` / `--ff-master`** — Specialized finish merge options
3. **Config file location flags** (`--local/--global/--system/--file`) in init

---

*Last updated: 2026-02-01*
