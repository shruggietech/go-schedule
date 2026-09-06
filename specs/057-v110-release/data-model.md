# Data Model: v1.1.0 Release

## Release identity

| Field | Rule |
| --- | --- |
| Version | `1.1.0` |
| Tag | `v1.1.0` |
| Source commit | Reviewed S057 merge commit |
| Release state | Absent, draft, then public |
| Latest | True only after successful promotion |

## Release copy

| Field | Rule |
| --- | --- |
| Heading | `## Highlights` |
| Highlights | Exactly four one-line bullets |
| Footer | Exactly one tagged full-changelog link |

## Artifact set

The draft owns four daemon and CLI archives, two desktop archives, one Windows MSI, and one Windows candidate manifest. Promotion adds one attended evidence archive and one checksum inventory. Every payload must be non-empty and bound to the same tag and source commit.

## State transitions

```text
prepared source -> reviewed merge -> immutable tag -> draft artifacts -> qualified candidate -> public release -> audited release
```

No transition may skip forward. Failure preserves the last safe state, normally reviewed source or a draft release.
