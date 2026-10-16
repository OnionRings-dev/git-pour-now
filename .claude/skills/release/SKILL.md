---
name: release
description: Prepare and verify tag-based git-pour-now releases
allowed-tools: Bash, Read, Edit, Grep, Glob
---

# Release git-pour-now

Read RELEASING.md and follow its tag-based SemVer workflow.

1. Review changes since the latest SemVer tag. Propose the appropriate major,
   minor, or patch version based on the documented public API. Initial
   development versions start at 0.1.0; never edit version constants in Go files.
2. Run `go test ./...`, `go vet ./...`, and review the diff.
3. Stable releases are tagged on main using `git pour release <version>`.
   Prereleases use a temporary `prerelease/<MAJOR.MINOR.PATCH>` branch created
   through `git pour prerelease start <version>`, and tags such as
   `git pour release 1.2.0-rc.1`.
4. Show the exact commit, tag, and publication plan. Obtain explicit user
   confirmation before pushing a branch or tag unless already authorized.
5. Push the branch before its individual tag. For prereleases, keep the remote
   version branch until its release workflow has completed. Verify CI and the
   generated GitHub release on OnionRings-dev/git-pour-now.
6. To promote a prerelease, finish its branch into main and create a new stable
   tag explicitly. Never move a published tag or replace published artifacts.

No upstream signing credentials, Homebrew tap, WinGet submission, or external
announcement is part of the fork's release process.
