# git-pour-now 🍻

A personal fork of [git-flow-next](https://github.com/gittower/git-flow-next),
maintained by OnionRings-dev. Keep development on `main`, work in short-lived
`feature/` branches, and identify releases exclusively with Git tags.

## Installation

Requires Go 1.23.6 or newer and Git.

```sh
go install github.com/OnionRings-dev/git-pour-now@latest
```

`go install` names the executable after the module (`git-pour-now`). Rename it to
`git-pour` in your PATH to enable `git pour`, or build it with the correct name:

```sh
go build -o git-pour main.go
# Copy git-pour to a directory in your PATH.
```

Release archives contain the `git-pour` executable (`git-pour.exe` on Windows).

## Workflow

```sh
git pour init --defaults
git pour feature start my-feature
# Edit files and commit your work.
git pour feature finish my-feature
# Feature changes are now on main.
git pour release 0.1.0
git push origin main
git push origin tag v0.1.0
```

The default configuration has one base branch, `main`, and two topic types: `feature/` and `prerelease/`. Features start from `main`, finish by merging into `main`, and update
by rebasing onto `main`. Finishing a feature does not create a release.
No `develop`, release, or hotfix branch is created by the default workflow.
Prerelease branches are created explicitly when needed.

`git pour release <version>` creates an annotated `v<version>` tag. Stable versions require the
configured trunk (normally `main`); prereleases require the matching
`prerelease/<MAJOR.MINOR.PATCH>` branch. It requires a clean working tree and no
unfinished merge or rebase. Versions must be valid [SemVer 2.0.0](https://semver.org/)
and greater than all existing local SemVer tags; build metadata does not affect
precedence. Both `1.2.3` and `v1.2.3` input create `v1.2.3`.
Prereleases such as `1.0.0-rc.1` and metadata such as `1.0.0+build.42` are supported.
Fetch remote tags before releasing in a shared repository (`git fetch --tags`).
Tags are never overwritten or automatically pushed.

Use `--message/-m` for a custom annotation and `--sign/-s` for a signed tag.

## Prereleases

Create a temporary branch for the target release core version. Prerelease tags
belong to that branch; finishing the branch merges its work into `main` and
removes the branch locally. The stable tag is created explicitly afterwards.

```sh
git pour prerelease start 1.2.0
# Prepare and commit the release.
git pour release 1.2.0-rc.1
git push -u origin prerelease/1.2.0
git push origin tag v1.2.0-rc.1
# Make further fixes and tag rc.2 if needed.
git pour prerelease finish 1.2.0
git pour release 1.2.0
git push origin main
git push origin tag v1.2.0
```

Keep the remote prerelease branch until its tag-triggered builds finish. You can
then delete it with `git push origin --delete prerelease/1.2.0`; tags retain the
prerelease history. The release workflow checks prerelease ancestry against that
branch and stable release ancestry against `main`.

## Versioning

The public API consists of the documented CLI commands, options, configuration
keys, and machine-readable output. Releases follow Semantic Versioning:

- PATCH: compatible bug fixes.
- MINOR: compatible additions or deprecations.
- MAJOR: incompatible public API changes.

The fork is in initial development until its first `1.0.0` release. Git tags are
the release version source; there are no hard-coded release version numbers.
`./scripts/build.sh` uses an exact tag at HEAD, or `dev` for an untagged checkout.
The release workflow validates pushed `v*` tags and their branch ancestry,
then builds archives and creates the GitHub release. Untagged source builds report
`dev`; `./scripts/build.sh v0.1.0` can prepare a versioned build locally.

## Configuration and compatibility

Settings live under `gitpour.*`; shared configuration lives in `.gitpour`.
Shell integration uses `GIT_POUR_*` environment variables. Hooks use names such
as `pre-pour-feature-start` and `filter-pour-feature-start-version`; rename any
existing hook scripts when opting into the fork.
This fork does not automatically rewrite an existing git-flow-next configuration;
initialize it with `git pour init --defaults` to opt into the new workflow.
Legacy git-flow-avh configuration can still be imported and its historical
`gitflow.*` keys are read only by the compatibility importer.

The explicit `--preset=classic`, `--preset=github`, and `--preset=gitlab`
compatibility presets and custom topic types remain available. With a configured
`release` topic type, `git pour release start/finish` manages that legacy branch;
the direct tag release command is available when no `release` topic is configured.
Advanced inherited capabilities are described in [docs](docs/README.md).

## Development

```sh
go test ./...
go vet ./...
go build -o git-pour main.go
```

The shell integration tests require bash, zsh, and fish.

## License and attribution

BSD 2-Clause; see [LICENSE](LICENSE). This fork retains the original copyright
and attribution. The historical changelog describes the upstream project.
