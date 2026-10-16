# GIT-POUR-PRERELEASE(1)

## NAME

git-pour-prerelease - Prepare a release on a temporary version branch

## SYNOPSIS

**git pour prerelease start** *MAJOR.MINOR.PATCH*

**git pour prerelease finish** *MAJOR.MINOR.PATCH*

## DESCRIPTION

Start creates **prerelease/<version>** from **main** (or the configured trunk).
The name must be a SemVer core version without a prefix, prerelease identifier,
or build metadata. Commit preparation and fixes on this branch, then create
preview tags using **git pour release <version>-rc.1**.

Finish merges the preparation branch into its parent using the normal topic
finish machinery and removes it locally. It creates no tag. Create the stable
release tag explicitly on main after integration. Update uses merge by default.
The usual topic branch options and conflict continuation commands apply.

## EXAMPLES

```
git pour prerelease start 1.2.0
# Commit changes.
git pour release 1.2.0-beta.1
git pour prerelease finish 1.2.0
git pour release 1.2.0
```

## EXIT STATUS

**0** on success; nonzero on invalid input, integration conflict, or Git error.

## SEE ALSO

**git-pour-release**(1), **git-pour-start**(1), **git-pour-finish**(1)
