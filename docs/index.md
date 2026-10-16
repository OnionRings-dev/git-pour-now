# git-pour-now Documentation Index

Examples using `develop`, release branches, or hotfix branches refer to the
explicit `--preset=classic` compatibility workflow. The default workflow uses
`main`, `feature/`, and temporary `prerelease/<version>` branches, with releases
identified by SemVer tags.

Comprehensive manpage-style documentation for git-pour-now commands and configuration.

## Quick Reference

| Command | Purpose | Documentation |
|---------|---------|---------------|
| **git-pour** | Main command overview | [git-pour(1)](git-pour.1.md) |
| **git pour release** | Create SemVer tags | [git-pour-release(1)](git-pour-release.1.md) |
| **git pour prerelease** | Prepare a version on a temporary branch | [git-pour-prerelease(1)](git-pour-prerelease.1.md) |
| **git-pour init** | Initialize git-pour | [git-pour-init(1)](git-pour-init.1.md) |
| **git-pour config** | Manage configuration | [git-pour-config(1)](git-pour-config.1.md) |
| **git-pour overview** | Repository status | [git-pour-overview(1)](git-pour-overview.1.md) |
| **git-pour integrate** | Integrate a base branch into its parent | [git-pour-integrate(1)](git-pour-integrate.1.md) |
| **git-pour worktree** | Manage worktrees for branches | [git-pour-worktree(1)](git-pour-worktree.1.md) |
| **git-pour completion** | Generate shell completion scripts | [git-pour-completion(1)](git-pour-completion.1.md) |
| **git-pour shell-init** | Print the shell integration script for directory switching | [git-pour-shell-init(1)](git-pour-shell-init.1.md) |

## Topic Branch Commands

| Command | Purpose | Documentation |
|---------|---------|---------------|
| **git-pour \<topic\> start** | Create topic branches | [git-pour-start(1)](git-pour-start.1.md) |
| **git-pour \<topic\> finish** | Complete topic branches | [git-pour-finish(1)](git-pour-finish.1.md) |
| **git-pour \<topic\> list** | List topic branches | [git-pour-list(1)](git-pour-list.1.md) |
| **git-pour \<topic\> update** | Update topic branches | [git-pour-update(1)](git-pour-update.1.md) |
| **git-pour \<topic\> delete** | Delete topic branches | [git-pour-delete(1)](git-pour-delete.1.md) |
| **git-pour \<topic\> rename** | Rename topic branches | [git-pour-rename(1)](git-pour-rename.1.md) |
| **git-pour \<topic\> checkout** | Switch to topic branches | [git-pour-checkout(1)](git-pour-checkout.1.md) |

## Configuration Reference

| Topic | Documentation |
|-------|---------------|
| **Configuration System** | [gitpour-config(5)](gitpour-config.5.md) |
| **git-flow-avh Compatibility** | [gitpour-config(5)#git-flow-avh-compatibility](gitpour-config.5.md#git-flow-avh-compatibility) |

## Getting Started

1. **[Initialize git-pour](git-pour-init.1.md)** - Set up your workflow
2. **[Configuration guide](git-pour-config.1.md)** - Customize your setup
3. **[Topic branch workflow](git-pour-start.1.md)** - Basic development workflow
4. **[Repository overview](git-pour-overview.1.md)** - Monitor your setup

## Workflow Guides

- **Classic GitFlow** - Traditional git-pour workflow
- **GitHub Flow** - Simplified continuous deployment workflow
- **GitLab Flow** - Multi-environment staged workflow
- **Custom Workflows** - Build your own with config commands

## Advanced Topics

- **[Three-layer configuration](gitpour-config.5.md#configuration-hierarchy)** - Understanding precedence
- **[Merge strategies](gitpour-config.5.md#merge-strategy-reference)** - Choose the right strategy
- **[Command overrides](gitpour-config.5.md#command-overrides)** - Fine-tune behavior
- **[Dynamic commands](git-pour.1.md#topic-branch-commands)** - Auto-generated commands

## Integration

- **CI/CD Integration** - Use JSON output for automation
- **IDE Integration** - Consume workflow data
- **Shell Integration** - Add status to prompts
- **Git Hooks** - Customize workflow behavior

All documentation follows Unix manpage conventions and includes practical examples.
