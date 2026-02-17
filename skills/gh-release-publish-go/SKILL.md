---
name: gh-release-publish-go
description: Use GitHub CLI to prepare and publish a new GitHub release for a Go module by inspecting the latest existing release, diffing against master, choosing the next semantic version, and creating release notes from the new changes. Use when asked to cut a new release from master/main and publish the newest Go package version.
---

# GitHub Release Publisher For Go

Use `gh` to inspect the current release state, compute what changed on `master`, select the next version tag, and publish a new release.

## Preconditions

- Ensure `gh` is installed and authenticated.
- Ensure the repository has a remote named `origin`.
- Prefer `master` as the release branch; if missing, use the repository default branch.

## Workflow

1. Identify repository and release branch.
- Run `gh repo view --json nameWithOwner,defaultBranchRef`.
- Set `OWNER_REPO` from `nameWithOwner`.
- Set `RELEASE_BRANCH` to `master` if it exists; otherwise use `defaultBranchRef.name`.

2. Inspect the latest release.
- Run `gh release list --limit 1`.
- If a release exists, read its tag and date with:
  `gh release view <latest-tag> --json tagName,publishedAt,name,url`.
- If no release exists, treat this as first release and use commit history from repo start.

3. Collect new changes on release branch.
- Fetch branch state: `git fetch origin`.
- If latest tag exists, compare tag to branch head:
  `gh api repos/$OWNER_REPO/compare/<latest-tag>...$RELEASE_BRANCH`.
- If no latest tag exists, gather recent commits from the release branch:
  `gh api repos/$OWNER_REPO/commits?sha=$RELEASE_BRANCH&per_page=100`.
- Stop and report when there are no commits ahead.

4. Choose next semantic version.
- Parse latest tag as SemVer, allowing optional `v` prefix.
- Default bump policy:
  - Bump major for breaking changes (`BREAKING CHANGE` or `!:`).
  - Else bump minor for `feat:` commits.
  - Else bump patch.
- If no prior release, start at `v0.1.0` unless project policy says otherwise.

5. Draft release notes from changes.
- Build sections:
  - Features
  - Fixes
  - Maintenance
- Include commit subjects and short SHAs.
- Add compare link context when previous tag exists.

6. Create and publish release.
- Create the release directly:
  `gh release create <new-tag> --target $RELEASE_BRANCH --title <new-tag> --notes-file <notes-file>`.
- If requested, create as prerelease using `--prerelease`.

7. Confirm result.
- Run `gh release view <new-tag> --json tagName,publishedAt,url`.
- Report the tag, publish time, and URL.

## Safety Rules

- Never overwrite or delete an existing release/tag.
- Never infer a major bump without explicit breaking-change evidence.
- When tag format in repository is non-SemVer, stop and ask for versioning policy.
- Before publish, summarize:
  - latest release
  - commits included
  - proposed new tag
  - release title
  - release notes outline

## Command Snippets

```bash
# repository metadata
gh repo view --json nameWithOwner,defaultBranchRef

# latest release
gh release list --limit 1
gh release view v1.2.3 --json tagName,publishedAt,name,url

# compare previous release to master
gh api repos/OWNER/REPO/compare/v1.2.3...master

# create release
gh release create v1.2.4 --target master --title v1.2.4 --notes-file /tmp/release-notes.md
```
