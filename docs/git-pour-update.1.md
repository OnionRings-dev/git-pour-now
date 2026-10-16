# GIT-POUR-UPDATE(1)

Examples using `develop`, release branches, or hotfix branches refer to the
explicit `--preset=classic` compatibility workflow. The default workflow uses
`main`, `feature/`, and temporary `prerelease/<version>` branches, with releases
identified by SemVer tags.

## NAME

git-pour-update - Update topic branches with parent changes

## SYNOPSIS

**git-pour** *topic* **update** [*name*] [*options*]

**git-pour update** [*base-branch*] [*options*]

**git-pour rebase** [*base-branch*] [**--continue**|**--abort**]

## DESCRIPTION

Update a branch with the latest changes from its parent branch using the configured downstream merge strategy. On the topic surface (**git-pour** *topic* **update**) it works with any topic branch type (feature, release, hotfix, support, or custom types); on the top-level surface (**git-pour update**) it also updates the current base branch from its parent (for example, **develop** from **main**).

The top-level **update** and **rebase** shorthands use the base-branch fallback only when the current branch is not a topic branch. Other branch-detection errors, including a cancelled ambiguous-branch prompt, stop the command without starting an update.

The update operation merges or rebases changes from the parent branch into the target branch, keeping it current with the latest development.

If a conflict occurs, the update saves a persistent state file and can be resumed with **--continue** or rolled back with **--abort** after resolving the conflict — the same resume/abort model as **git-pour finish** and **git-pour integrate**. The **--continue**/**--abort** flags act only on an in-progress update; an in-progress finish or integrate is never affected.

## ARGUMENTS

*topic*
: The topic branch type (feature, release, hotfix, support, or any configured custom type)

*name*
: Name of the topic branch to update on the **git-pour** *topic* **update** surface. If omitted, the current topic branch is used.

*base-branch*
: Base branch to update on the top-level **git-pour update** or **git-pour rebase** surface. If omitted, the current base branch is used.

## OPTIONS

**--rebase**
: Force rebase strategy instead of the configured downstream strategy

**--continue**, **-c**
: Continue the update operation after resolving merge conflicts. This acts only on an in-progress update. If a **finish** or **integrate** operation is in progress instead, update refuses non-destructively, names the owning operation, prints its resume/abort commands, and exits 3 without touching it. With nothing in progress, **--continue** reports "no merge in progress" and exits 3.

**--abort**, **-a**
: Abort the update operation and return to the original state. When no update operation is in progress, **--abort** is a no-op and exits successfully. Like **--continue**, it acts only on an update: a foreign in-progress finish or integrate is refused (exit 3) rather than aborted.

## MERGE STRATEGIES

The merge strategy used when updating is determined by configuration:

**merge** (default for most branches)
: Standard Git merge preserving both branch histories

**rebase**
: Rebase topic branch commits on top of parent branch for linear history

The strategy is configured via:
1. **Branch defaults**: `gitpour.branch.<type>.downstreamStrategy`
2. **Command overrides**: `gitpour.<type>.downstreamStrategy`
3. **Flag override**: `--rebase` always forces rebase strategy

## EXAMPLES

### Basic Usage

Update current topic branch (shorthand):
```bash
git pour update
```

Update specific feature branch:
```bash
git pour feature update user-authentication
```

Update release branch with latest hotfixes:
```bash
git pour release update 1.2.0
```

### Force Rebase Strategy

Update using rebase regardless of configuration:
```bash
git pour feature update my-feature --rebase
```

Update current branch with rebase:
```bash
git pour update --rebase
```

### Resume or Abort After a Conflict

Continue an update after resolving conflicts:
```bash
git pour feature update --continue my-feature
```

Abort an in-progress update and restore the pre-update state:
```bash
git pour feature update --abort my-feature
```

Resume or abort a base-branch update (top-level surface):
```bash
git pour update --continue
git pour update --abort
```

### Typical Workflows

Before finishing a long-running feature:
```bash
# Make sure feature is current before finishing
git pour feature update long-feature
git pour feature finish long-feature
```

Update release with production hotfixes:
```bash
# Hotfix was applied to main, update release
git pour release update 1.2.0
```

Keep hotfix current during development:
```bash
# Update hotfix with any other main changes
git pour hotfix update critical-fix
```

## CONFLICT RESOLUTION

When a conflict occurs, the update saves its state and stops so you can resolve it. Resume or roll back with git-pour's own **--continue**/**--abort** rather than raw Git commands, so the merge state is cleared correctly:

```bash
# Start update
git pour feature update my-feature

# Conflicts occur - Git shows conflict markers.
# Edit files to resolve, then stage them:
vim conflicted-file.js
git add conflicted-file.js

# Complete the update (merge or rebase, whichever the strategy uses):
git pour feature update --continue my-feature
```

To discard the in-progress update and return to the pre-update state:

```bash
git pour feature update --abort my-feature
```

For a base branch on the top-level surface, use the same flags without a topic type:

```bash
git pour update --continue
git pour update --abort
```

A rebase-strategy update that conflicts again on a later commit stays resumable: resolve, stage, and run **--continue** again. **--abort** with nothing in progress is a forgiving no-op (exit 0).

**Known gap:** resuming a **squash**-strategy update is not supported (there is no merge commit to complete and **--abort** cannot roll it back). Complete or discard a squash-strategy conflict with raw Git, then re-run the update.

## CONFIGURATION

Update behavior is controlled by these configuration keys:

### Branch-Level Defaults
```bash
git config gitpour.branch.feature.downstreamStrategy rebase
git config gitpour.branch.release.downstreamStrategy merge
git config gitpour.branch.hotfix.downstreamStrategy merge
```

### Command-Level Overrides
```bash
# Always use rebase for feature updates
git config gitpour.branch.feature.downstreamStrategy rebase

# Always use merge for release updates
git config gitpour.branch.release.downstreamStrategy merge
```

## STRATEGY RECOMMENDATIONS

### Feature Branches
**Rebase** (recommended) - Creates clean, linear history without unnecessary merge commits

### Release Branches
**Merge** (recommended) - Preserves the integration history of hotfixes and changes

### Hotfix Branches
**Merge** (recommended) - Preserves context of what changes were integrated

## PARENT RELATIONSHIPS

Each topic branch type updates from its configured parent:

- **feature** → updates from **develop**
- **release** → updates from **main** (to get hotfixes)
- **hotfix** → updates from **main** (rare, but possible)
- **support** → updates from **main**

## EXIT STATUS

**0**
: Successful update, or **--abort** with nothing in progress (forgiving no-op)

**1**
: git-pour is not initialized, or shorthand branch detection failed or was cancelled

**2**
: Invalid input (e.g. unknown branch type)

**3**
: A Git operation failed, a merge or rebase is already in progress, unresolved conflicts remain, or there is no in-progress operation to continue

**5**
: Target or parent branch not found

**6**
: Validation error

## SEE ALSO

**git-pour**(1), **git-pour-start**(1), **git-pour-finish**(1), **git-pour-config**(1), **git-rebase**(1), **git-merge**(1), **gitpour-config**(5)

## NOTES

- The **--rebase** flag always overrides configured strategy
- Update operations preserve your local commits while integrating parent changes
- Use **git pour update** shorthand when on a topic branch
- On conflict, resume with **git pour ... update --continue** or roll back with **--abort** rather than raw Git, so the saved state is cleared correctly
- Regular updates keep topic branches easier to merge when finishing
- Update from parent before finishing long-running branches
