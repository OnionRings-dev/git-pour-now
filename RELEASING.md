# Releasing git-pour-now

Release versions come exclusively from immutable Git tags and follow
[Semantic Versioning 2.0.0](https://semver.org/). The public API is the documented
CLI, configuration keys, and machine-readable output. Use PATCH for compatible
fixes, MINOR for compatible additions/deprecations, and MAJOR for incompatible
changes. Start initial development at `0.1.0`.

1. Merge completed features into `main` and run `go test ./...` and `go vet ./...`.
2. Fetch tags with `git fetch --tags` and review the changes since the last release.
3. Run `git pour release <version>` on a clean `main`, for example
   `git pour release 0.1.0` or `git pour release 1.0.0-rc.1`.
4. Push `main`, then the individual release tag:
   `git push origin main` and `git push origin tag v0.1.0`.
5. The release workflow validates SemVer and the expected branch ancestry, builds the archives,
   and creates a GitHub release. A pre-release identifier marks a preview release.

For prereleases, create `git pour prerelease start 1.2.0`, commit preparation
changes, then tag with `git pour release 1.2.0-rc.1`. Push the prerelease branch
before the tag. Keep it on the remote until the prerelease workflow finishes.
Use `git pour prerelease finish 1.2.0` to merge into main, then tag the stable
release explicitly with `git pour release 1.2.0`.

Do not edit version constants, move published tags, or
replace published artifacts. Corrections require a new version.
Untagged local builds report `dev`. `./scripts/build.sh` derives the version from
an exact tag at HEAD; an optional explicit version supports local build testing.

The fork currently distributes unsigned archives through GitHub Releases.
Upstream's signing and package registry integrations are not configured for this
project. The BSD license and upstream attribution must accompany distribution.
