# Git-pour Hook Examples

This directory contains example hook and filter scripts for git-pour-now. Copy the scripts you want to use to your repository's `.git/hooks/` directory and make them executable.

Release/hotfix branch examples require the explicit `--preset=classic`
compatibility workflow. The default workflow prepares prereleases with the
`prerelease` topic type and creates release tags through `git pour release`.

## Installation

```bash
# Copy a hook to your repository
cp contrib/hooks/pre-pour-release-start /path/to/repo/.git/hooks/

# Make it executable
chmod +x /path/to/repo/.git/hooks/pre-pour-release-start
```

## Available Examples

### Filters

Filters transform values during git-pour operations:

| Script | Description |
|--------|-------------|
| `filter-pour-release-start-version` | Validates and normalizes version numbers (adds 'v' prefix) |
| `filter-pour-release-finish-tag-message` | Appends changelog to tag messages |

### Hooks

Hooks execute before (pre) or after (post) operations:

| Script | Description |
|--------|-------------|
| `pre-pour-feature-start` | Validates feature branch naming conventions |
| `pre-pour-release-start` | Checks for uncommitted changes and remote sync |
| `post-pour-release-finish` | Shows next steps and can trigger notifications |

## Customization

These are starting points - modify them to fit your workflow:

- **CI Integration**: Add checks for CI status using `gh` or API calls
- **Notifications**: Send Slack/Discord messages on release completion
- **Versioning**: Enforce your team's version number format
- **Naming**: Set your own branch naming conventions

## Creating Your Own

1. Create a file with the correct name (see patterns below)
2. Make it executable: `chmod +x filename`
3. Test it before relying on it

### Naming Patterns

**Filters:**
```
filter-pour-{type}-{action}-{target}
  type: release, hotfix
  action: start, finish
  target: version, tag-message
```

**Hooks:**
```
{pre,post}-pour-{type}-{action}
  pre/post: when to run
  type: feature, release, hotfix, support
  action: start, finish, publish, track, delete
```

## Environment Variables

All hooks and filters receive context via environment variables:

| Variable | Description |
|----------|-------------|
| `BRANCH` | Full branch name (e.g., `feature/my-feature`) |
| `BRANCH_NAME` | Short name (e.g., `my-feature`) |
| `BRANCH_TYPE` | Type (e.g., `feature`, `release`) |
| `BASE_BRANCH` | Parent branch (e.g., `develop`, `main`) |
| `ORIGIN` | Remote name |
| `VERSION` | Version (for release/hotfix) |
| `EXIT_CODE` | Post-hooks only: operation result |

## Sharing Hooks

To share hooks across your team:

```bash
# Option 1: Use Git's core.hooksPath
mkdir .githooks
cp contrib/hooks/* .githooks/
git config core.hooksPath .githooks

# Option 2: Create symbolic links
ln -s ../../.githooks/pre-pour-release-start .git/hooks/
```

## Documentation

For more details, see `docs/gitpour-hooks.7.md` or run:
```bash
git pour hooks --help
```
