# GIT-POUR-SELF-UPDATE(1)

## NAME

git-pour-self-update - Check for updates or update the git-pour executable

## SYNOPSIS

**git pour self-update** [**--check**|**-c**]

## DESCRIPTION

Query the latest version of the **github.com/OnionRings-dev/git-pour-now** Go
module. Without **--check**, build that exact version and replace the running
executable. This uses the same **@latest** channel as **scripts/install.sh**,
including Go's configured module proxy and checksum verification settings.
It does not require a published GitHub Release archive.

Go 1.23.6 or newer and network access are required. Git may also be needed when
Go fetches directly from the repository. Installation supports macOS and Linux;
**--check** also works on Windows with Go installed.

The command works outside Git repositories and does not activate or change
repository configuration. It does not modify shell profiles or PATH.

## OPTIONS

**--check**, **-c**
: Print installed and latest versions without building or replacing the executable.
  If an update is available, print the command to install it.

**--help**, **-h**
: Show command help.

## VERSION COMPARISON

Release builds use the version injected from their Git tag. Go-installed builds
can use the tag or pseudo-version recorded in their Go build information even
when **git pour version** reports **dev**. The latest version follows Go's
**@latest** selection rules, rather than always selecting the newest commit on
**main**. No downgrade is performed when the installed version is equal to or
newer than that version, using SemVer precedence.

A local development build may have no comparable version. In that case,
**--check** reports that the installed version cannot be compared; it does not
claim that an update is available. Running without **--check** replaces that
development build with the latest module version.

## INSTALLATION

The update is built in a temporary directory for the current operating system
and architecture. A staged copy beside the installed executable preserves its
permissions and is renamed over the executable after the build and copy
succeed. A failed build or staging operation leaves the existing executable
untouched. If invoked through a symbolic link, its target is updated and the
link is preserved. No terminal restart is needed.

The installation directory must be writable. For installations managed by a
package manager or a system administrator, use that installation method instead
if the command cannot write to the directory. On Windows, reinstall externally
because the running executable cannot be replaced by this command.

Version checks time out after one minute; builds time out after ten minutes.

## EXAMPLES

Check for updates:

```sh
git pour self-update --check
```

Install the latest module version:

```sh
git pour self-update
```

## EXIT STATUS

**0**
: The check completed, including when an update is available or the installed
  development version cannot be compared; or installation succeeded; or the
  installed version is already equal to or newer than latest.

**1**
: Invalid arguments, missing Go, lookup or build failure, timeout, unsupported
  installation platform, or inability to replace the executable.

## SEE ALSO

**git-pour**(1), **go-install**(1)
