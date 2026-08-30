# Contract: v0.9.0 Publication

## Preparation Contract

- The release-preparation pull request references issue #79 without closing it.
- `CHANGELOG.md` contains the complete v0.9.0 section and a fresh Unreleased
  section.
- `.github/release-notes/v0.9.0.md` contains only four to six highlights and a
  full-changelog link.
- The workflow chooses `.github/release-notes/${tag}.md` and disables generated
  notes.
- All local and hosted review gates pass before merge.

## Tag Contract

- Tag creation requires explicit authorization distinct from pull-request
  publication.
- `main` must still equal the reviewed S034 merge commit immediately before tag
  creation.
- Neither tag `v0.9.0` nor its GitHub release may already exist.
- The tag name is exactly `v0.9.0` and points to the reviewed merge commit.

## Publication Contract

- The release workflow publishes exactly these eight expected assets:
  - `go-schedule_v0.9.0_linux_amd64.tar.gz`
  - `go-schedule_v0.9.0_linux_arm64.tar.gz`
  - `go-schedule_v0.9.0_darwin_amd64.tar.gz`
  - `go-schedule_v0.9.0_darwin_arm64.tar.gz`
  - `go-schedule-desktop_v0.9.0_linux_amd64.tar.gz`
  - `go-schedule-desktop_v0.9.0_darwin_arm64.tar.gz`
  - `go-schedule_v0.9.0_windows_amd64.msi`
  - `SHA256SUMS.txt`
- Every publishing job concludes successfully.
- The GitHub release is not a draft or prerelease unless the maintainer later
  gives an explicit conflicting instruction.
- The public body matches the reviewed tag-specific highlights file.

## Completion Contract

- Every expected asset is present and non-empty.
- Every non-checksum asset appears exactly once in `SHA256SUMS.txt` and verifies
  against its downloaded bytes.
- README release badge and sample health version identify v0.9.0.
- Publication evidence is posted to issue #79 before it is closed.
- Any missing evidence leaves issue #79 open and the release incomplete.
