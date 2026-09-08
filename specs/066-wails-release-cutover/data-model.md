# Data Model: Wails Release Cutover

This slice changes release and verification state rather than persisted scheduler data. The entities below define the auditable candidate contract.

## DesktopCandidate

- **revision**: Full Git commit identifier used by every qualification result.
- **version_context**: Tag-shaped version supplied only when the release workflow is later invoked.
- **source_modules**: Root module for daemon and CLI, desktop module for Wails.
- **status**: `unqualified`, `qualifying`, `qualified`, or `rejected`.

### Validation

- Every required result refers to the same revision.
- A skipped, cancelled, timed-out, or missing required result prevents `qualified`.
- Qualification does not imply publication.

## DesktopDistribution

- **platform**: `windows`, `darwin`, or `linux`.
- **architecture**: `amd64` for Windows/Linux or `arm64` for macOS.
- **format**: MSI or compressed desktop archive containing a native application bundle or executable.
- **desktop_entry**: Stable `gosched-gui` executable or bundle entry.
- **companions**: `goschedd`, `gosched`, README, LICENSE, and CHANGELOG.
- **platform_assets**: Windows installer identity, macOS icon and bundle metadata, or Linux desktop and icon tree.

### Validation

- Missing required payload rejects the distribution.
- The desktop entry is built from `desktop/`, while daemon and CLI are built from the root module.
- A current distribution contains no Fyne code or runtime dependency.

## ParityEntry

- **legacy_capability**: One materially supported Fyne workflow or behavior.
- **disposition**: `wails`, `retired`, or `blocked`.
- **destination**: Wails route/control or retirement rationale.
- **evidence**: Test, specification, or documentation reference.
- **user_impact**: Observable difference, if any.

### Validation

- Each capability has exactly one disposition.
- `blocked` prevents completion of #157.
- `retired` requires rationale and user-facing transition guidance when discoverable behavior changes.

## QualificationResult

- **candidate_revision**: Parent DesktopCandidate revision.
- **gate**: Stable gate name.
- **platform**: Optional platform dimension.
- **outcome**: `passed`, `failed`, `skipped`, `cancelled`, or `missing`.
- **evidence**: Local output, hosted check URL, or release-contract record.

### State transitions

```text
unqualified -> qualifying -> qualified
                         -> rejected
rejected -> qualifying
```

Only complete passing required results permit `qualified`.
