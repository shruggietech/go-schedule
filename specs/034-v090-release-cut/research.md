# Research: v0.9.0 Release Cut

## Tag-Specific Public Notes

**Decision**: Use `.github/release-notes/${tag}.md` as the release action's
dynamic body path.

**Rationale**: The workflow fails honestly if the current tag has no reviewed
copy and cannot silently reuse the prior release's highlights.

**Alternatives considered**:

- Keep a single `RELEASE_NOTES.md`: rejected because future releases can reuse
  stale copy without a failure.
- Embed v0.9.0 highlights directly in the workflow: rejected because the next
  tag would publish v0.9.0 text unless another workflow edit happened first.

## Generated Notes

**Decision**: Disable generated release notes.

**Rationale**: GitHub-generated commit and pull-request inventories violate the
explicit highlights-only request and duplicate the detailed changelog.

**Alternatives considered**:

- Prefix generated notes with highlights: rejected because the page remains
  exhaustive.
- Edit the release body after publication: rejected because the first published
  state is wrong and the workflow is no longer deterministic.

## Changelog Boundary

**Decision**: Cut the complete current Unreleased section as v0.9.0 dated
2026-08-30 and leave a new empty Unreleased section above it.

**Rationale**: The repository has no release after v0.8.0, and the existing
v0.9.0 milestone names the intended minor-version boundary.

**Alternatives considered**:

- Publish v0.8.1: rejected because the accumulated capabilities are substantial
  minor-version additions rather than a patch-only correction.
- Publish v1.0.0: rejected because the maintainer describes the product as
  pre-beta and has not selected a stable major release.

## Release Verification

**Decision**: Use local canonical gates, hosted pull-request CI, and post-tag
asset and checksum readback. Do not require an interactive installation session.

**Rationale**: The release workflow already builds the formal packages on their
native runners, and the user explicitly deferred manual Windows observation.
Artifact completeness and checksums provide objective publication evidence.

**Alternatives considered**:

- Require a VM install and uninstall walkthrough: rejected as outside the
  maintainer-approved release gate.
- Trust workflow success without an asset audit: rejected because partial
  release states can still leave a release page present.
