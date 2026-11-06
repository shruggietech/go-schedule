# Feature Specification: Maintenance Convergence and First-use Release Handoff

**Feature Branch**: `codex/134-maintenance-convergence-publication`

**Created**: 2026-11-06

**Status**: Implemented

**Delivery**: Review branch codex/134-maintenance-convergence-publication; all 16 source tasks, eight gates, 35 desktop/11 site browser cases and 1763 current QA passes complete; see [verification](verification.md) and [delivery](delivery.md). Owner merge and actual #419 public delivery remain separate.

**Input**: Owner kickoff authorizes S134 source autopilot, automatic push/official PR, CI monitoring and finding resolution, with two total Codex review requests including automatic review. Stop ready for final owner review and merge. Scope traces to #419, #429, #430 and dependency proposals #426-#428.

## Clarifications

### Session 2026-11-06

Resolved under the owner's autopilot decision policy; no user questions were required.

- Q: Does source autopilot authorize release mutation? A: No. The explicit stop before owner review/merge governs this run; #419 and the later exact-source/asset approvals remain in the continuation.
- Q: Should Pages retain the unsupported one-hour input? A: Use the selected action's actual ten-minute ceiling under #373's existing exception classification, preserving the outer job budget and all other verification limits.
- Q: Which dependency scope is authorized? A: Consolidate compatible changes from #426-#428 and warranted same-package patch corrections; name any superseded pin or materially blocked proposal explicitly.
- Q: What happens when a failure precedes diagnostic creation? A: Retain generated diagnostics and the actual unsuccessful outcome; record absence rather than claim a nonexistent artifact.

## User Scenarios & Testing

### User Story 1 - Finish each first-task navigation (Priority: P1)

A maintainer exercises the first-task entry points from all three platform guides. Each transition reaches its intended route and content before another guide is opened, so a previous action cannot interrupt the next navigation.

**Why this priority**: The current first-task browser flow produced an attributable overlapping-navigation failure, despite the proposal changing only desktop dependency files.

**Independent Test**: Exercise the existing production first-task flow and inspect route/content completion between each platform entry, retaining the original failed observation.

**Acceptance Scenarios**:

1. **Given** a platform guide, **When** its first-task link is selected, **Then** the expected first-task route and document finish before the next platform guide is requested.
2. **Given** a failed website exercise, **When** its runner exits, **Then** its unsuccessful result and retrievable bounded context/trace diagnostics remain available.
3. **Given** the first-task guide, **When** its content and mobile layout are inspected, **Then** the existing exact-run, failure/recovery, source/public availability and overflow expectations retain their meaning.

### User Story 2 - Know the effective Pages deadline (Priority: P1)

A maintainer can see the deployment deadline actually enforced by the selected service action and recover an incomplete deployment using the accepted source and build artifact.

**Why this priority**: The current one-hour setting was clamped to ten minutes and the first merged-source deployment timed out while syncing files.

**Independent Test**: Check the declared policy against the selected action's enforced ceiling and follow the existing source/artifact-bound recovery procedure.

**Acceptance Scenarios**:

1. **Given** the selected action's ten-minute ceiling, **When** the repository declares its effective completion policy, **Then** configuration and documentation agree with that ceiling while the outer job allowance remains intact.
2. **Given** an unfinished deployment after a successful build, **When** the outcome is reported, **Then** build success, incomplete deployment and a later successful recovery remain separately attributable.
3. **Given** a later retry, **When** it uses the existing accepted artifact, **Then** the original failure and source/artifact identity remain retained.

### User Story 3 - Review one coherent maintenance candidate (Priority: P2)

A maintainer receives one bounded candidate containing the compatible changes proposed in #426-#428 together with the navigation and deployment-policy corrections, accurate release records and a concrete continuation for public delivery.

**Why this priority**: These proposals arrived after S133 integration and share one maintenance-release verification and review boundary.

**Independent Test**: Inspect the consolidated candidate's declared inputs, source/version records and delivery handoff against each included proposal and functional issue.

**Acceptance Scenarios**:

1. **Given** the three accepted proposals, **When** their changes are consolidated, **Then** each proposal is traceable and every affected dependency manifest and lock/graph remains coherent.
2. **Given** prepared v2.0.1 source, **When** source readiness is described, **Then** public downloads continue to identify v2.0.0 until actual approved publication.
3. **Given** the owner's pre-merge stop, **When** this source phase finishes, **Then** #429/#430 have implemented source outcomes, #419's public-delivery outcome remains explicitly incomplete, and the existing tag/draft and concrete-asset approval steps remain in the handoff.

### Edge Cases

- A prior first-task heading is visible before the current route transition finishes.
- Website execution fails before diagnostic files exist; no absent artifact is advertised as collected evidence.
- An input requests more than the selected deployment action supports, or omits/shortens the declared effective allowance.
- The Pages backend remains queued or syncing beyond its ceiling; no source cause is inferred from that delay alone.
- An individual dependency proposal exposes a material regression; defer only that update with explicit traceability rather than advertise incomplete integration.
- Historical tagged release bytes and prior failed observations remain immutable.
- A source PR, QA report or website deployment does not itself publish maintenance packages.

## Requirements

### Functional Requirements

- **FR-001**: Every platform-guide first-task transition MUST finish at its intended route and document before the flow advances.
- **FR-002**: Existing guide-content, exact-run, recovery and responsive expectations MUST remain intact.
- **FR-003**: A failed website browser exercise MUST retain its generated bounded context/trace diagnostics when produced and its unsuccessful exit; diagnostic absence remains explicit.
- **FR-004**: The effective Pages completion allowance MUST match the selected action's supported ceiling, explicitly classified under #373's existing exception rule.
- **FR-005**: The outer job allowance, other owner-directed verification limits and deliberate expiry scenarios MUST remain intact.
- **FR-006**: Build success, incomplete deployment and successful recovery MUST remain distinct and bound to their actual source/artifact identities.
- **FR-007**: Compatible changes from proposals #426-#428 MUST form one coherent, individually traceable maintenance candidate; warranted same-package patch corrections MUST be recorded, and any deferred proposal MUST be named with its incomplete outcome.
- **FR-008**: Candidate notes, changelog and comparison attribution MUST describe the consolidated source accurately while preserving historical records and current public selection.
- **FR-009**: The source handoff MUST preserve the two-request review cap, owner merge boundary and later exact-source tag/draft plus concrete-asset promotion controls.
- **FR-010**: #419 MUST remain open until actual public package/discovery delivery; verification activity MUST NOT become a functional issue closure gate.

### Key Entities

- **Maintenance candidate**: Source identity, selected version, included proposals and implemented issue outcomes.
- **Navigation transition**: Origin platform guide, requested first-task destination, completed route/document and diagnostic outcome.
- **Deployment receipt**: Source identity, artifact identity, attempt, effective deadline and terminal or incomplete outcome.
- **Release continuation**: Accepted source, later approvals, actual staged asset inventory and public delivery/discovery state.

## Success Criteria

### Measurable Outcomes

- **SC-001**: All three platform entry paths finish at the intended first-task route before the next navigation starts.
- **SC-002**: The declared deployment allowance is the effective ten-minute ceiling, with no advertised one-hour allowance that the selected action shortens.
- **SC-003**: Every failed website exercise retains an unsuccessful result and retrievable diagnostics when produced.
- **SC-004**: All three selected dependency proposals are accounted for by delivered changes or an explicit bounded deferral.
- **SC-005**: No historical release record/asset is changed and no unperformed merge, staging or publication is represented as completed.

## Assumptions

- Accepted baseline is e09231bf0074480e5765a776f7fbf61c8b5610b1, tree-equivalent to the reviewed S133 candidate.
- The planned Next patch is 16.4.5, superseding proposal #428's 16.4.4 pin after primary-source confirmation of the published security release. No affected request-time surface or application exploitability is demonstrated for this static export.
- Prepared v2.0.1 remains an appropriate maintenance target; public latest remains v2.0.0 until later approved publication.
- #429 and #430 are source outcomes; #419 retains the distinct public-release outcome.
- The existing canonical verification, browser contracts, current-candidate QA and publication workflows are reused.
- Current-state initialization/persistence/refusal obligations remain; historical migrations and old-version qualification are outside the owner-approved prelaunch baseline.
- Attended native installation/GUI remains untested under the standing waiver. Original Windows causes #355/#356 and owner-deferred clustered work remain independent.
