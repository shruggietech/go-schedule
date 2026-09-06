# Data Model: v1.1.1 Release Recovery

## Historical release attempt

| Field | Rule |
| --- | --- |
| Version | `1.1.0` |
| Tag | Preserved at original commit |
| GitHub release | Unpublished draft, then absent |
| Public status | Never public |

## Corrected release identity

| Field | Rule |
| --- | --- |
| Version | `1.1.1` |
| Tag | `v1.1.1` |
| Source commit | Reviewed S058 merge commit |
| Required CI | Successful `main` CI for the exact source commit |
| Release state | Absent, draft, then public |
| Latest | True only after successful promotion |

## Release copy

| Field | Rule |
| --- | --- |
| Heading | `## Highlights` |
| Highlights | Exactly four one-line bullets |
| Footer | Exactly one tagged full-changelog link |

## Artifact set

The draft owns four daemon and CLI archives, two desktop archives, one Windows MSI, and one Windows candidate manifest. Promotion adds one attended evidence archive and one checksum inventory. Every payload must be non-empty and bound to the same v1.1.1 tag and source commit.

## State transitions

```text
preserved v1.1.0 tag + unpublished draft
    -> reviewed v1.1.1 preparation
    -> retired v1.1.0 draft
    -> immutable v1.1.1 tag
    -> exact-commit CI confirmation
    -> draft v1.1.1 artifacts
    -> qualified candidate
    -> public v1.1.1 release
    -> audited release
```

No transition may skip forward. Failure preserves the last safe state, normally reviewed source or a draft release.
