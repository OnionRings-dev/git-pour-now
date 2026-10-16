# Git-Flow-AVH Configuration Reference

Examples using `develop`, release branches, or hotfix branches refer to the
explicit `--preset=classic` compatibility workflow. The default workflow uses
`main`, `feature/`, and temporary `prerelease/<version>` branches, with releases
identified by SemVer tags.

This document provides a complete reference of all git-flow-avh configuration options for compatibility analysis with git-pour-now.

## Global Configuration Settings

### Remote Repository
- **`gitpour.origin`** - Set different remote repository name
  - Default: `origin`
  - Example: `git config gitpour.origin myorigin`

### Environment Variables
- Can be set in `.gitpour_export` file in home directory
- **`GIT_MERGE_AUTOEDIT`** - Control merge commit editing
  - Example: `export GIT_MERGE_AUTOEDIT=no`

### Working Tree State
- **`gitpour.allowdirty`** - Allow starting branches with dirty working tree
  - Type: Boolean
  - Default: `false`
  - Example: `git config --bool gitpour.allowdirty true`

## Branch Prefixes (Legacy Format)

### Core Branch Names
- **`gitpour.branch.master`** - Production branch name
  - Default: `master`
- **`gitpour.branch.develop`** - Development branch name
  - Default: `develop`

### Topic Branch Prefixes
- **`gitpour.prefix.feature`** - Feature branch prefix
  - Default: `feature/`
- **`gitpour.prefix.bugfix`** - Bugfix branch prefix
  - Default: `bugfix/`
- **`gitpour.prefix.release`** - Release branch prefix
  - Default: `release/`
- **`gitpour.prefix.hotfix`** - Hotfix branch prefix
  - Default: `hotfix/`
- **`gitpour.prefix.support`** - Support branch prefix
  - Default: `support/`
- **`gitpour.prefix.versiontag`** - Version tag prefix
  - Default: empty string

## Command-Specific Configuration

### Configuration Variable Convention
- **Format**: `gitpour.<branchtype>.<command>.<option>`
- **Boolean Values**:
  - Affirmative: `Y`, `YES`, `T`, `TRUE`
  - Negative: `N`, `NO`, `F`, `FALSE`

### Feature Branch Options

#### Start Command
- **`gitpour.feature.start.fetch`** - Fetch from remote before starting
  - Type: Boolean
  - Default: `false`

#### Finish Command
- **`gitpour.feature.finish.fetch`** - Fetch from remote before finishing
  - Type: Boolean
  - Default: `false`
- **`gitpour.feature.finish.rebase`** - Use rebase instead of merge
  - Type: Boolean
  - Default: `false`
- **`gitpour.feature.finish.preserve-merges`** - Preserve merges during rebase
  - Type: Boolean
  - Default: `false`
- **`gitpour.feature.finish.push`** - Push to remote after finishing
  - Type: Boolean
  - Default: `false`
- **`gitpour.feature.finish.keep`** - Keep branch after finishing
  - Type: Boolean
  - Default: `false`
- **`gitpour.feature.finish.force-delete`** - Force delete branch after finishing
  - Type: Boolean
  - Default: `false`
- **`gitpour.feature.finish.squash`** - Squash commits when finishing
  - Type: Boolean
  - Default: `false`

### Bugfix Branch Options

#### Start Command
- **`gitpour.bugfix.start.fetch`** - Fetch from remote before starting
  - Type: Boolean
  - Default: `false`

#### Finish Command
- **`gitpour.bugfix.finish.fetch`** - Fetch from remote before finishing
  - Type: Boolean
  - Default: `false`
- **`gitpour.bugfix.finish.rebase`** - Use rebase instead of merge
  - Type: Boolean
  - Default: `false`
- **`gitpour.bugfix.finish.preserve-merges`** - Preserve merges during rebase
  - Type: Boolean
  - Default: `false`
- **`gitpour.bugfix.finish.push`** - Push to remote after finishing
  - Type: Boolean
  - Default: `false`
- **`gitpour.bugfix.finish.keep`** - Keep branch after finishing
  - Type: Boolean
  - Default: `false`
- **`gitpour.bugfix.finish.force-delete`** - Force delete branch after finishing
  - Type: Boolean
  - Default: `false`
- **`gitpour.bugfix.finish.squash`** - Squash commits when finishing
  - Type: Boolean
  - Default: `false`

### Release Branch Options

#### Start Command
- **`gitpour.release.start.fetch`** - Fetch from remote before starting
  - Type: Boolean
  - Default: `false`

#### Finish Command
- **`gitpour.release.finish.fetch`** - Fetch from remote before finishing
  - Type: Boolean
  - Default: `false`
- **`gitpour.release.finish.sign`** - Sign the release tag
  - Type: Boolean
  - Default: `false`
- **`gitpour.release.finish.signingkey`** - GPG key for signing
  - Type: String
  - Default: Git configuration default
- **`gitpour.release.finish.push`** - Push to remote after finishing
  - Type: Boolean
  - Default: `false`
- **`gitpour.release.finish.keep`** - Keep branch after finishing
  - Type: Boolean
  - Default: `false`
- **`gitpour.release.finish.force-delete`** - Force delete branch after finishing
  - Type: Boolean
  - Default: `false`
- **`gitpour.release.finish.squash`** - Squash commits when finishing
  - Type: Boolean
  - Default: `false`
- **`gitpour.release.finish.notag`** - Skip tag creation
  - Type: Boolean
  - Default: `false`

### Hotfix Branch Options

#### Start Command
- **`gitpour.hotfix.start.fetch`** - Fetch from remote before starting
  - Type: Boolean
  - Default: `false`

#### Finish Command
- **`gitpour.hotfix.finish.fetch`** - Fetch from remote before finishing
  - Type: Boolean
  - Default: `false`
- **`gitpour.hotfix.finish.sign`** - Sign the hotfix tag
  - Type: Boolean
  - Default: `false`
- **`gitpour.hotfix.finish.signingkey`** - GPG key for signing
  - Type: String
  - Default: Git configuration default
- **`gitpour.hotfix.finish.push`** - Push to remote after finishing
  - Type: Boolean
  - Default: `false`
- **`gitpour.hotfix.finish.keep`** - Keep branch after finishing
  - Type: Boolean
  - Default: `false`
- **`gitpour.hotfix.finish.force-delete`** - Force delete branch after finishing
  - Type: Boolean
  - Default: `false`
- **`gitpour.hotfix.finish.squash`** - Squash commits when finishing
  - Type: Boolean
  - Default: `false`
- **`gitpour.hotfix.finish.notag`** - Skip tag creation
  - Type: Boolean
  - Default: `false`

### Support Branch Options

#### Start Command
- **`gitpour.support.start.fetch`** - Fetch from remote before starting
  - Type: Boolean
  - Default: `false`

#### Rebase Command
- **`gitpour.support.rebase.interactive`** - Use interactive rebase
  - Type: Boolean
  - Default: `false`

## Configuration Examples

### Setting Global Configuration
```bash
# Enable fetching before operations
git config --global gitpour.feature.start.fetch true
git config --global gitpour.feature.finish.fetch true

# Use rebase for feature finishing
git config gitpour.feature.finish.rebase true

# Sign release and hotfix tags
git config gitpour.release.finish.sign true
git config gitpour.hotfix.finish.sign true

# Set custom remote origin
git config gitpour.origin upstream
```

### Repository-Specific Configuration
```bash
# Allow dirty working tree for this repository
git config --bool gitpour.allowdirty true

# Custom branch prefixes
git config gitpour.prefix.feature "feat/"
git config gitpour.prefix.release "rel/"
```

## Compatibility Notes

### Branch Types
- git-flow-avh includes `bugfix` branches not in original git-pour
- Support branches are available but less commonly used
- Master branch is default (not main as in newer systems)

### Key Differences from git-pour-now
1. **Configuration Format**: AVH uses `gitpour.prefix.*` vs git-pour-now's `gitpour.branch.*.prefix`
2. **Branch Names**: AVH defaults to `master` vs git-pour-now's `main`
3. **Merge Strategies**: git-pour-now reads the same `finish.rebase`/`finish.squash` flags AVH uses, and adds branch-level `upstreamStrategy`/`downstreamStrategy` configuration on top
4. **Tag Management**: AVH has `notag` flags vs git-pour-now's `tag` boolean configuration

### Migration Path
When migrating from git-flow-avh to git-pour-now, the following mappings apply:

| AVH Configuration | git-pour-now Equivalent |
|------------------|--------------------------|
| `gitpour.branch.master` | `gitpour.branch.main` |
| `gitpour.prefix.feature` | `gitpour.branch.feature.prefix` |
| `gitpour.prefix.release` | `gitpour.branch.release.prefix` |
| `gitpour.prefix.hotfix` | `gitpour.branch.hotfix.prefix` |
| `gitpour.prefix.versiontag` | `gitpour.branch.*.tagprefix` |
| `gitpour.feature.finish.rebase` | `gitpour.feature.finish.rebase` (unchanged — read directly) |
| `gitpour.release.finish.squash` | `gitpour.release.finish.squash` (unchanged — read directly) |
| `gitpour.release.finish.notag` | `gitpour.release.finish.notag` (unchanged — read directly; or `gitpour.branch.release.tag=false`) |

---

*This reference is based on git-flow-avh wiki documentation and is maintained for compatibility analysis with git-pour-now.*
