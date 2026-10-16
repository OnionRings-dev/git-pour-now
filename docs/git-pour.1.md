# GIT-POUR(1)

## NAME

git-pour - A modern implementation of the git-pour branching model

## SYNOPSIS

**git-pour** [**--verbose**|**-v**] *command* [*args*]

## DEFAULT WORKFLOW

Initialize with **git pour init --defaults**. Features start from and finish into
**main**. Create an annotated release tag with **git pour release 1.2.3**.
The default workflow has no develop, release, or hotfix branch; prereleases use
**prerelease/<version>** and stable releases are tagged on main.
Examples below involving those branches require **--preset=classic**.

## DESCRIPTION

git-pour-now is a modern, maintainable implementation of the git-pour branching model written in Go. It provides a set of Git extensions to provide high-level repository operations for Vincent Driessen's branching model.

This implementation maintains compatibility with existing git-pour repositories and supports both preset workflows (Classic GitFlow, GitHub Flow, GitLab Flow) and fully customizable branch configurations.

## GLOBAL OPTIONS

**--verbose**, **-v**
: Enable verbose output showing detailed operation information

**--help**, **-h**
: Show help information for any command

## COMMANDS

### Core Commands

**init**
: Initialize git-pour in the current repository. See **git-pour-init**(1).

**config**
: Manage git-pour configuration for base branches and topic branch types. See **git-pour-config**(1).

**overview**
: Display repository workflow overview. See **git-pour-overview**(1).

**integrate** [*branch*]
: Integrate a base branch into its parent (e.g. develop into main). See **git-pour-integrate**(1).

**worktree**
: Manage worktrees for branches: add, remove, list, prune, and print the computed path. See **git-pour-worktree**(1).

**shell-init** *shell*
: Print a shell script that lets git-pour change your shell's working directory when a command navigates. See **git-pour-shell-init**(1).

**version**
: Show version information. Prints `<version> (git-pour-now)`, so tooling that parses the first whitespace-separated token reads a bare version number.

### Topic Branch Commands

Topic branch commands are dynamically generated based on your configuration. Default types include **feature**, **release**, **hotfix**, **support**, plus any custom types you define.

Each topic branch type supports these subcommands:

**start** *name* [*base*]
: Create new topic branch. See **git-pour-start**(1).

**finish** [*name*]
: Complete and merge topic branch. See **git-pour-finish**(1).

**list** [**--worktrees**]
: List existing topic branches. See **git-pour-list**(1).

**update** [*name*]
: Update topic branch from parent. See **git-pour-update**(1).

**delete** *name*
: Delete topic branch. See **git-pour-delete**(1).

**rename** *old-name* *new-name*
: Rename topic branch. See **git-pour-rename**(1).

**checkout** *name*|*prefix*
: Switch to topic branch. See **git-pour-checkout**(1).

**track** *name*
: Create local branch tracking a remote topic branch. See **git-pour-track**(1).

### Shorthand Commands

**delete** [*name*]
: Delete current or specified topic branch. See **git-pour-delete**(1).

**update** [*name*]
: Update current or specified topic branch from parent. See **git-pour-update**(1).

**rebase** [*name*]
: Rebase current or specified topic branch (alias for update --rebase).

**rename** [*new-name*]
: Rename current topic branch. See **git-pour-rename**(1).

**finish** [*name*]
: Finish current or specified topic branch. See **git-pour-finish**(1).

**publish**
: Publish current topic branch to remote (planned feature).

## WORKFLOW PRESETS

git-pour-now supports three workflow presets:

**Classic GitFlow**
: Traditional git-pour with main, develop, feature/, release/, and hotfix/ branches.

**GitHub Flow**
: Simplified workflow with main and feature/ branches only.

**GitLab Flow**
: Multi-environment workflow with production, staging, main, feature/, and hotfix/ branches.

## CONFIGURATION

git-pour-now uses Git's configuration system to store all settings under the **gitpour.*** namespace. Configuration follows a three-layer hierarchy, but some options intentionally skip branch defaults:

1. **Branch Type Defaults** (gitpour.branch.*) - Default behavior for branch types (when applicable)
2. **Command Overrides** (gitpour.*type*.*command*.*) - Override defaults for specific operations
3. **Command-line Flags** - Always take highest precedence

See **git-pour-config**(1) and **gitpour-config**(5) for detailed configuration reference.

## GIT-POUR-AVH COMPATIBILITY

git-pour-now automatically detects and translates git-flow-avh configuration at runtime without modifying existing settings. Legacy configuration is mapped to the new format transparently.

## EXAMPLES

Initialize a repository with Classic GitFlow:
```bash
git pour init --preset=classic
```

Start a new feature:
```bash
git pour feature start my-feature
```

Finish current branch (shorthand):
```bash
git pour finish
```

Configure custom branch type:
```bash
git pour config add topic bugfix develop --prefix=bug/
```

## FILES

**.git/config**
: Repository-specific git-pour configuration stored under gitpour.* keys

## SEE ALSO

**git-pour-init**(1), **git-pour-config**(1), **git-pour-completion**(1), **git-pour-start**(1), **git-pour-finish**(1), **git-pour-integrate**(1), **git-pour-update**(1), **git-pour-delete**(1), **git-pour-track**(1), **git-pour-worktree**(1), **git-pour-shell-init**(1), **gitpour-config**(5), **git**(1)

## AUTHORS

git-pour-now is maintained by OnionRings-dev and derives from git-flow-next.

Original git-pour by Vincent Driessen. git-pour AVH Edition by Peter van der Does.

## BUGS

Report bugs at https://github.com/OnionRings-dev/git-pour-now/issues

## COPYRIGHT

This is free software; see the source for copying conditions.
