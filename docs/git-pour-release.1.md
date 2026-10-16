# GIT-POUR-RELEASE(1)

## NAME

git-pour-release - Create an annotated SemVer release tag

## SYNOPSIS

**git pour release** *version* [**-m**|**--message** *message*] [**-s**|**--sign**]

## DESCRIPTION

Create an annotated **v<version>** tag at HEAD. Accept a SemVer 2.0.0 version,
optionally prefixed by **v**. Stable versions require the configured trunk
(normally **main**). Prerelease versions require the corresponding temporary
**prerelease/<MAJOR.MINOR.PATCH>** branch. Build metadata does not select a branch.

The repository must be initialized, with no pending changes or integration.
The version must have greater SemVer precedence than every existing local
SemVer tag. Unrelated tag names are ignored. Existing tags are never moved.
Fetch tags from the remote before releasing in a shared repository.

The command creates no branch, changes no version file, and does not push.

## OPTIONS

**-m**, **--message** *message*
: Annotation text. Defaults to **Release v<version>**.

**-s**, **--sign**
: Sign the annotated tag using Git's signing configuration.

## EXAMPLES

```
git pour init --defaults
git pour prerelease start 1.2.0
# Commit release preparation.
git pour release 1.2.0-rc.1
git push -u origin prerelease/1.2.0
git push origin tag v1.2.0-rc.1
# Keep the remote branch until the prerelease build finishes.
git pour prerelease finish 1.2.0
git pour release 1.2.0
git push origin main
git push origin tag v1.2.0
```

## COMPATIBILITY

When a compatibility preset configures a **release** topic type, that name
continues to expose legacy **start**/**finish** branch commands. Direct tag
releases are available when no release topic is configured.

## EXIT STATUS

**0** on success; nonzero for invalid input, branch, repository state, or Git error.

## SEE ALSO

**git-pour-init**(1), **git-pour-prerelease**(1), **gitpour-config**(5),
[Semantic Versioning 2.0.0](https://semver.org/)
