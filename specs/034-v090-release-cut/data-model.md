# Data Model: v0.9.0 Release Cut

S034 adds no application database entity. Its durable records are repository and GitHub release objects.

## Release Boundary

- **Version**: `v0.9.0`
- **Previous version**: `v0.8.0`
- **Release date**: `2026-08-30`, refreshed if review crosses a date boundary
- **Reviewed commit**: the S034 pull-request merge commit
- **State**: prepared, reviewed, tagged, publishing, verified, or failed

Transitions are monotonic except that a failed publication may be retried for the same immutable tag after the cause is corrected without changing content.

## Release Notes

- **Identity**: exact version tag
- **Title**: Highlights
- **Highlights**: four to six user-meaningful bullets
- **Detailed record**: one link to the v0.9.0 changelog at the same tag
- **Excluded content**: generated change inventories, installation guidance, and copied detailed changelog entries

## Artifact Set

- **Headless archives**: Linux amd64/arm64 and macOS amd64/arm64
- **Desktop bundles**: Linux amd64 and macOS arm64
- **Windows installer**: Windows amd64 MSI
- **Integrity record**: `SHA256SUMS.txt`

Each non-checksum asset has a unique filename and exactly one checksum entry.

## Publication Evidence

- Tag and target commit
- Release URL and published state
- Release workflow run and job conclusions
- Asset names and byte sizes
- Checksum entries and local verification results
- README badge and health-version synchronization commit
- Final issue #79 comment and closed state
