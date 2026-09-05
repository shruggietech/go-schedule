# Implementation Plan: GUI Navigation and Information

**Branch**: `codex/017-gui-info-navigation` | **Date**: 2026-08-28 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/017-gui-info-navigation/spec.md`

## Summary

Reorder the existing leading desktop navigation to group Tasks and Groups, then add a final Info view that reuses the embedded application mark and build-version source while presenting three canonical external links. The implementation stays inside the GUI package, uses existing Fyne controls and resources, adds no backend call or dependency, and protects both issues with headless regression tests.

## Technical Context

**Language/Version**: Go 1.25.0

**Primary Dependencies**: Fyne 2.7.4; Go standard library `net/url`; existing `internal/buildinfo` package

**Storage**: N/A; the view is derived entirely from compiled resources

**Testing**: Go `testing` with Fyne's shared headless test application; canonical repository verification through `sh scripts/verify.sh all`

**Target Platform**: Supported Linux, macOS, and Windows desktop builds

**Project Type**: Go desktop application backed by a local daemon

**Performance Goals**: Info builds synchronously from local constants and adds no daemon, disk, or network latency

**Constraints**: No new backend surface, dependency, brand asset, state, update check, telemetry, or external-link customization; preserve current Activity terminology and badge behavior

**Scale/Scope**: One reordered five-item navigation collection, one new local view, three hyperlinks, and focused GUI tests; closes #29 and #32

## Constitution Check

*GATE: Passed before research and rechecked after design.*

- **I. Code Quality**: The change is small, idiomatic, and confined to GUI composition. Canonical URLs are constructed without fallible runtime parsing.
- **II. Testing Standards**: Tests are written first and must fail on the old four-item order and absent Info view before implementation. Headless GUI and race verification remain mandatory.
- **III. User Experience Consistency**: Current Activity terminology is preserved; visible text identifies the application and every destination.
- **IV. Performance Requirements**: No scheduler hot path changes. Info uses bounded local resources and performs no I/O while rendering.
- **V. Autonomous Build-Phase Execution**: The full Spec-Kit sequence runs on a review branch and halts once before publication.
- **Engineering constraints**: No schema, configuration, IPC, dependency, or platform-support change.

Post-design recheck: passed. The design uses existing dependencies and resources, keeps behavior local, and introduces no constitutional exception.

## Project Structure

### Documentation (this feature)

```text
specs/017-gui-info-navigation/
├── checklists/
│   ├── requirements.md
│   └── ux.md
├── contracts/
│   └── gui-navigation-info.md
├── data-model.md
├── plan.md
├── quickstart.md
├── research.md
├── spec.md
├── tasks.md
└── verification.md
```

### Source Code (repository root)

```text
gui/
├── app.go          # tab construction and Activity badge ownership
├── app_test.go     # exact navigation-order and badge-position regressions
├── icon.go         # existing embedded full mark, unchanged
├── info.go         # new local Info view
└── info_test.go    # identity, version, mark, and canonical-link regressions

internal/buildinfo/
└── buildinfo.go    # existing version source, unchanged
```

**Structure Decision**: Add one `gui/info.go` file following the package's existing per-view convention. Keep navigation ownership in `gui/app.go`, reuse `appIcon` from `gui/icon.go`, and read `buildinfo.Version` directly. No API or viewmodel layer is justified for immutable process-local identity.

## Delivery Sequence

1. Capture the current tab-order failure and missing Info contract in headless tests.
2. Reorder existing tab construction while retaining the Activity item reference.
3. Build Info from local image, version, attribution, and canonical links.
4. Run focused GUI tests, then all eight repository gates.
5. Record issue closure in the changelog and commit locally for PR review.

## Decision Log

### Preserve Activity instead of restoring Logs

- **Options**: follow issue #29 literally and restore Logs; preserve Activity; rename again.
- **Decision**: Preserve Activity.
- **Rationale**: PR #46 intentionally established Activity as the accurate name for the combined alert and log surface. Reverting it would contradict current behavior and reopen solved UX ambiguity.

### Build Info entirely inside the GUI process

- **Options**: add a daemon health/API field; load metadata from the network; use existing compiled resources.
- **Decision**: Use existing `buildinfo.Version`, `appIcon`, and fixed official destinations locally.
- **Rationale**: All required facts already exist in the binary. An API would add coupling and make a basic information view fail when the daemon is down.

### Use the standard hyperlink control

- **Options**: custom buttons with browser-launch callbacks; rich-text markup; standard hyperlink widgets.
- **Decision**: Use standard Fyne hyperlinks backed by explicit URL values.
- **Rationale**: The control already supplies focus, keyboard activation, hyperlink styling, and operating-system URL opening. Custom behavior would duplicate toolkit responsibilities without user value.

### Reuse the full application mark

- **Options**: create a banner; use the reduced window tile; reuse `appIcon`.
- **Decision**: Reuse `appIcon` with aspect-preserving containment.
- **Rationale**: It is the existing large-surface identity resource. A new banner belongs to the broader brand-package backlog and would expand this slice.

## Complexity Tracking

No constitution violations. Table intentionally empty.
