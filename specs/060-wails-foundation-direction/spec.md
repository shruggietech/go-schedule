# Feature Specification: Wails Foundation and Experience Direction

**Feature Branch**: `codex/060-wails-foundation-direction`

**Created**: 2026-09-07

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Local canonical eight-gate verification and hosted Windows, macOS, Linux, and Chromium proof jobs passed in [CI run 34091460576](https://github.com/shruggietech/go-schedule/actions/runs/34091460576); first-round review fixes add direct state and action regressions on PR [#193](https://github.com/shruggietech/go-schedule/pull/193).

**Input**: Bundle GitHub issues [#149](https://github.com/shruggietech/go-schedule/issues/149) and [#150](https://github.com/shruggietech/go-schedule/issues/150) into one end-to-end slice that selects and proves the supported Wails foundation, establishes one fresh accessible desktop direction, and completes under the operator-authorized autopilot publication and review workflow.

## Clarifications

### Session 2026-09-07

- Q: Should the production foundation select the current stable Wails line or adopt the Wails v3 prerelease? -> A: Select the latest stable v2 release and record a bounded upgrade policy; prerelease v3 is not a production dependency.
- Q: Should the experience direction preserve the existing Fyne structure or begin from the approved control-center outcome and brand system? -> A: Begin from the control-center outcome and approved brand system, with no widget-for-widget Fyne translation.
- Q: How does the operator's up-front publication direction interact with the usual autopilot pre-publication halt? -> A: It supplies publication authorization for this S060 review branch, pull request, and verified in-scope review fixes only; merge, tag, and release remain unauthorized.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Choose a Supportable Desktop Foundation (Priority: P1)

As a maintainer, I can rely on one recorded desktop technology decision backed by a small proof on every supported platform, so the replacement program starts from evidence instead of framework preference.

**Why this priority**: Every later Wails shell, connection, workflow, packaging, and accessibility slice depends on a stable, maintainable foundation.

**Independent Test**: Review the decision record and run the proof checks on Windows, macOS, and Linux to show local health, a list result, one live event, one native action boundary, offline assets, and clean application shutdown.

**Acceptance Scenarios**:

1. **Given** the supported desktop platforms and current upstream release state, **When** a maintainer reads the decision record, **Then** one supported Wails release, frontend language, UI approach, build tool, and test stack are named with versions, rationale, ownership cost, and rejected alternatives.
2. **Given** a clean Windows, macOS, or Linux runner, **When** the disposable proof is built and tested, **Then** it exercises daemon health, one list operation, one event, and one desktop-native action boundary without network-hosted application assets.
3. **Given** a future dependency update, **When** the upgrade policy is applied, **Then** stable releases are evaluated deliberately and prerelease adoption requires a separate recorded decision.

---

### User Story 2 - See One Fresh Control-Center Direction (Priority: P1)

As a scheduler operator, I can understand the proposed desktop control center at a glance across representative screens, states, appearances, and window sizes before production screens are migrated.

**Why this priority**: Technology selection without a concrete experience target would allow the replacement to reproduce the current widget tree and postpone the hardest interaction decisions.

**Independent Test**: Open the committed prototype at ordinary laptop and compact widths, switch appearance and connection states, navigate by keyboard, and compare Tasks, task editing, Schedule, Activity, and target switching against the measurable experience contract.

**Acceptance Scenarios**:

1. **Given** a local operator at an ordinary laptop width, **When** they inspect each representative screen, **Then** the active target, page purpose, primary action, critical status, and next useful step are distinguishable without opening secondary help.
2. **Given** a compact window, **When** the same screens are inspected, **Then** the primary workflow remains usable without horizontal page scrolling or loss of target identity.
3. **Given** empty, loading, disconnected, degraded, destructive, and success states, **When** the prototype presents each state, **Then** it uses consistent hierarchy, plain language, semantic cues beyond color, and a clear recovery or continuation action.

---

### User Story 3 - Inherit a Measurable Accessible Design Contract (Priority: P2)

As an implementer or reviewer of later migration slices, I can use a compact set of tokens, component rules, interaction behaviors, and acceptance checks that prevent page-specific styling and accessibility drift.

**Why this priority**: A visual concept only becomes a dependable foundation when later work can implement and verify it consistently.

**Independent Test**: Trace every token and interaction rule to the prototype, inspect automated accessibility and responsive checks, and confirm that the contract covers the shared patterns required by issues #151 through #157.

**Acceptance Scenarios**:

1. **Given** the approved brand kit, **When** the design contract is inspected, **Then** typography, spacing, color, elevation, icon, focus, interaction-state, and motion rules are defined without inventing a competing identity.
2. **Given** keyboard-only, zoomed, reduced-motion, high-contrast, or assistive-technology use, **When** the representative prototype is evaluated, **Then** its requirements and automated checks have objective pass conditions.
3. **Given** a later production screen, **When** its author needs a table, form, status, empty state, error, confirmation, or disclosure pattern, **Then** the direction defines the intended reusable pattern or explicitly assigns it to the shell slice.

### Edge Cases

- An upstream stable release can appear after this decision; the recorded pin remains authoritative until a deliberate upgrade passes the same proof.
- A platform runner can compile the proof without a usable interactive desktop session; automated checks must distinguish build and contract evidence from attended native observation.
- Windows can lack a preinstalled WebView2 runtime, Linux distributions can expose different WebKitGTK packages, and macOS can require platform SDK updates; prerequisites and user-facing diagnosis must be recorded.
- Native clipboard or dialog APIs can be unavailable in headless tests; the proof must keep the native boundary replaceable and prove both the adapter contract and a real platform build.
- A visual state can meet contrast ratios while still relying on color alone; text, icon, shape, or position must carry the same meaning.
- A compact layout can hide secondary controls, but it must never hide active target identity, destructive consequences, or the route back to primary navigation.
- Offline capability excludes remote fonts, CDN scripts, analytics, telemetry, and network-hosted application assets.
- Prototype code is evidence for the decision, not the production shell; later slices must not silently treat disposable shortcuts as application architecture.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The slice MUST record one supported stable Wails release, frontend language, UI approach, build tool, and test stack with exact baseline versions.
- **FR-002**: The decision MUST compare a bounded shortlist and document maintainability, release status, dependency footprint, accessibility support, packaging consequences, upgrade policy, and reasons for rejection.
- **FR-003**: The selected foundation MUST use platform-native web rendering and MUST NOT bundle Electron or another standalone browser runtime.
- **FR-004**: The selected foundation MUST keep application assets local and support complete offline operation.
- **FR-005**: A disposable proof MUST expose representative local daemon health, list, event, and desktop-native action boundaries through typed contracts.
- **FR-006**: The proof MUST build and pass its automated contract checks on Windows, macOS, and Linux.
- **FR-007**: The proof MUST define deterministic substitutes for daemon and native boundaries so frontend behavior can be tested without a running daemon or interactive desktop.
- **FR-008**: The proof MUST document Windows WebView2, macOS SDK, Linux WebKitGTK, Node, and Go prerequisites plus their packaging consequences.
- **FR-009**: The proof MUST remain isolated from the shipped Fyne application and MUST NOT alter current production behavior, package contents, or release identity.
- **FR-010**: The selected path MUST include a license inventory and MUST justify every new dependency against the repository's ownership constraints.
- **FR-011**: The experience direction MUST be visibly distinct from the current Fyne application and MUST organize the product as a target-aware desktop control center rather than a widget-for-widget port.
- **FR-012**: Representative views MUST cover Tasks, task editing, Schedule, Activity, target switching, and a useful compact-window layout.
- **FR-013**: Representative states MUST cover new, empty, loading, connected, disconnected, degraded, destructive, validation-error, and success conditions.
- **FR-014**: The direction MUST define reusable typography, spacing, color, surface, border, elevation, icon, interaction-state, and motion tokens derived from the approved brand system.
- **FR-015**: The direction MUST preserve unmistakable active-target identity throughout every mutating workflow, with `This computer` as the default local target.
- **FR-016**: The direction MUST define keyboard order, visible focus, skip behavior, semantic landmarks, accessible names, status announcements, contrast, zoom, compact reflow, target sizes, and reduced-motion behavior with measurable thresholds.
- **FR-017**: Text and essential meaning MUST meet WCAG 2.2 AA contrast, MUST NOT rely on color alone, and MUST remain usable at 200 percent zoom.
- **FR-018**: All primary workflows in the prototype MUST be operable by keyboard without a trap, and visible focus MUST remain distinguishable in both light and dark appearances.
- **FR-019**: Motion MUST be nonessential, bounded by the approved token range, and removed or reduced when the operating system requests reduced motion.
- **FR-020**: The prototype and design contract MUST use sentence case, plain operational language, restrained density, and no decorative dashboard metric that lacks a user decision.
- **FR-021**: Automated checks MUST cover type safety, unit behavior, representative component behavior, accessibility rules, responsive overflow, offline asset references, and deterministic Go boundary behavior.
- **FR-022**: The slice MUST identify which proof elements graduate into issue #151 and which are intentionally discarded, preventing disposable code from becoming an implicit production contract.
- **FR-023**: The slice MUST satisfy and close issues #149 and #150 while leaving #151, #152, and later migration issues open.
- **FR-024**: The eventual pull request MUST complete hosted cross-platform proof checks and no more than two Codex review rounds before maintainer merge review.

### Key Entities

- **Foundation decision**: The supported release, frontend stack, test stack, prerequisites, upgrade policy, and rejected alternatives that govern later desktop work.
- **Disposable proof**: An isolated, non-shipping Wails application that demonstrates typed Go/frontend boundaries, local assets, events, native adapters, representative interface states, and cross-platform buildability.
- **Experience contract**: The measurable hierarchy, layout, state, token, language, and accessibility rules for the future control center.
- **Representative view**: A prototype state for Tasks, task editing, Schedule, Activity, target switching, or compact-window behavior.
- **Target identity**: The persistent user-visible identity and condition of the daemon against which an operation will run.
- **Evidence matrix**: Platform and interaction results that distinguish automated build, automated behavior, and attended native observations.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: One decision record names exact baseline versions and resolves 100 percent of the technology categories required by issue #149, with no unresolved production dependency choice.
- **SC-002**: The disposable proof builds on Windows, macOS, and Linux and its automated suite demonstrates all four required boundaries: health, list, event, and native action.
- **SC-003**: The proof contains zero network-hosted fonts, scripts, styles, images, analytics, or runtime application assets.
- **SC-004**: Six representative views and nine required state classes are inspectable at both 1440 by 900 and 900 by 650 CSS-pixel viewports without horizontal page scrolling.
- **SC-005**: Automated accessibility checks report zero serious or critical violations across representative light, dark, compact, disconnected, degraded, and destructive states.
- **SC-006**: All interactive prototype controls are reachable and operable by keyboard, focus never becomes trapped, and focus indication reaches at least 3:1 contrast against adjacent colors.
- **SC-007**: Text and meaningful UI graphics meet WCAG 2.2 AA contrast, and the primary workflows remain usable at 200 percent zoom with target identity still visible.
- **SC-008**: A reviewer can identify the active target, page purpose, primary action, critical status, and next useful step on every representative view in no more than 10 seconds per view.
- **SC-009**: The complete repository verification suite remains green and the shipped Fyne binary, installer, and release workflow are unchanged by the disposable proof.
- **SC-010**: Issues #149 and #150 have direct closing evidence, while no completion claim is made for the production shell or connection layer.

## Assumptions

- The current stable Wails v2 line is the appropriate production baseline; Wails v3 remains prerelease software on 2026-09-07 and is evaluated only as a rejected production alternative.
- The disposable proof is committed for review and reproducibility but isolated under an experimental directory with its own dependency manifests.
- The prototype uses representative local data and replaceable boundaries rather than requiring a running daemon or changing the existing IPC contract.
- The approved go-schedule brand kit remains authoritative for identity and tokens, while layout and interaction may change materially.
- Maintainer selection of the direction is represented by approval and merge of the S060 pull request; production shell implementation does not begin in this slice.
- Native attended validation belongs to the later packaging and cutover issue; S060 proves platform builds and adapter contracts without claiming installed-user observation.
- The current Fyne application remains the shipped desktop until issue #157 completes the cutover.
