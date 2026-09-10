# Data Model: Cumulative v1.4.0 Release Preparation

## Release boundary

| Field | Rule |
| --- | --- |
| Public predecessor | `v1.1.1` |
| Target | `v1.4.0` |
| History shape | One cumulative dated section after an empty Unreleased section |
| Source commit | Reviewed S081 squash merge commit on `main` |
| Intermediate milestones | Preserved as delivery history, not synthesized releases |

## Release metadata

| Surface | Required identity |
| --- | --- |
| README health example | `1.4.0` |
| Changelog heading | One dated `[1.4.0]` section |
| Changelog comparison | `v1.1.1...v1.4.0` |
| Release note | Four one-line highlights plus one tagged changelog link |
| Git tag | `v1.4.0`, absent until separately authorized |

## Draft artifact set

The draft contains four daemon and CLI archives, two Wails desktop archives, one Windows MSI, and one Windows candidate manifest. Qualification adds one attended evidence archive. Promotion adds one checksum inventory that covers every other public asset exactly once.

## Windows upgrade baseline

| Field | Rule |
| --- | --- |
| Baseline artifact | Public `go-schedule_v1.1.1_windows_amd64.msi` from its exact GitHub release URL, SHA-256 `f7ac8f56f28330b016eb6e505e424b19e9bfbe435591cfbc54a723c91ac8e567` |
| Candidate artifact | Draft `go-schedule_v1.4.0_windows_amd64.msi` |
| Preserved state | Tasks, run history, supported appearance intent, daemon identity, actor and credential records, notification configuration, remote configuration |
| Default-off proof | No notification channel, localhost MCP listener, or remote HTTPS listener enabled by upgrade alone |
| Native observations | All 47 established attended scenario identities pass against the exact candidate |
| Current window contract | Wails 1440 by 900 logical content on a sufficiently large work area, otherwise clamped within 90 percent of the logical work area |

## Publication lifecycle

```mermaid
flowchart TB
    A[Reviewed S081 pull request] --> B[Squash merge to main]
    B --> C[Successful exact-commit main CI]
    C --> D[Explicit tag authorization]
    D --> E[Immutable v1.4.0 tag]
    E --> F[Draft artifact staging]
    F --> G[Fresh and v1.1.1 upgrade qualification]
    G --> H[Evidence archive upload]
    H --> I[Explicit promotion authorization]
    I --> J[No-rebuild public promotion]
    J --> K[Final identity and checksum audit]
    K --> L[Issue and milestone closure]
```

Failure cannot advance the lifecycle. The safe retained states are reviewed source before tagging and an unpublished draft after staging.
