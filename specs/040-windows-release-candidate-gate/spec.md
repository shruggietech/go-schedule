# Feature Specification: Windows Release Candidate Gate

**Feature Branch**: `codex/040-windows-release-candidate-gate`

**Created**: 2026-09-02

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Review branch `codex/040-windows-release-candidate-gate`; repository implementation and canonical eight-gate verification completed on 2026-09-02; exact-candidate attended acceptance remains intentionally unclaimed

**Input**: Implement GitHub issue #94 as the fail-closed Windows release-candidate gate, consume the attended installer and uninstall acceptance remaining in #98 and coordinator #96, and ensure the exact tested MSI is the artifact eligible for publication.

## Context and Scope

The v0.9.1 release passed the repository's automated checks while ordinary installed use still exposed restricted IPC access, failed task execution, an oversized first window, recurring error presentation, and incomplete installer lifecycle controls. Subsequent slices corrected and automated important parts of those behaviors, but no release controller currently requires an attended Windows 11 observation of the complete installed product.

S040 creates that missing control. It stages one immutable Windows candidate in a draft release, records its identity, guides an operator through the required native scenarios, validates a machine-readable evidence bundle, and prevents publication unless the exact candidate has complete passing evidence. It does not publish a release, merge its own pull request, or manufacture successful evidence for a desktop environment that was not observed.

## Clarifications

### Session 2026-09-02

- Q: May release rebuild a Windows MSI after a candidate with the same source revision passes? -> A: No. Windows publication must use the byte-identical MSI whose SHA-256, ProductCode, version, source commit, workflow run, and repository were validated.
- Q: What happens when attended evidence is skipped, unavailable, incomplete, stale, timed out, or manually marked inconclusive? -> A: Release readiness fails closed, and the result names every missing or non-passing requirement.
- Q: Can one desktop configuration satisfy both display requirements? -> A: One clean Windows 11 candidate environment may be reconfigured, but evidence must contain separate standard-DPI and high-DPI or mixed-DPI observations with the active monitor and effective DPI recorded for each.
- Q: Should S040 change GUI or installer behavior when a native scenario fails? -> A: Only after the exact baseline failure is reproduced and recorded. Any new recurring-error correction must also pass the issue's proof-before-commit sequence against the uncommitted working tree. Otherwise S040 reports the failed gate without inventing corrective code.
- Q: When are #94, #98, and #96 complete? -> A: Tooling implementation alone does not close them. They close only after a real, exact release candidate produces complete passing attended evidence and all issue-level acceptance criteria are audited.
- Q: Does the operator's automatic-publication direction authorize S040 push and pull-request creation? -> A: Yes, for this review branch and pull request, including one additional `@Codex review` round if needed. Merge, tag, release, and more than two review rounds remain unauthorized.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Prepare One Traceable Candidate (Priority: P1)

As a release maintainer, I create one Windows candidate whose source and package identity remain unambiguous from preparation through publication.

**Why this priority**: Native results are meaningful only when they apply to the exact bytes users will receive.

**Independent Test**: Stage a candidate from an immutable source commit, then prove its manifest and artifact agree on repository, commit, version, ProductCode, workflow origin, filename, and SHA-256. Mutating any field or byte must invalidate it.

**Acceptance Scenarios**:

1. **Given** an approved version tag, **When** candidate staging succeeds, **Then** the release remains a draft with one MSI and a machine-readable identity manifest retained together.
2. **Given** a retained candidate, **When** its identity is checked, **Then** the recorded repository, source commit, workflow run, version, ProductCode, filename, size, and SHA-256 match authoritative sources.
3. **Given** a rebuilt, renamed, altered, expired, cross-repository, or different-commit artifact, **When** it is submitted as the tested candidate, **Then** the gate rejects it with a specific mismatch.

---

### User Story 2 - Prove Installed Core Behavior (Priority: P1)

As a release operator using a clean Windows 11 desktop, I execute one guided candidate walkthrough that proves normal-user access, conservative native window behavior, non-spamming recovery, and real scheduled work.

**Why this priority**: These are the application-breaking v0.9.1 regressions that prior headless and source-level checks failed to detect.

**Independent Test**: Install the exact candidate on the documented clean snapshot and complete every required access, display, failure, recovery, and task scenario. The evidence validator must accept only observations with the required native measurements and passing outcomes.

**Acceptance Scenarios**:

1. **Given** a fresh per-machine installation, **When** the intended non-elevated user launches the product, **Then** that user can reach the LocalSystem service while an unrelated local user remains denied.
2. **Given** no prior profile state and then retained v0.9.1-era profile state, **When** the GUI first opens at standard DPI and high or mixed DPI, **Then** it is restored, has visible margins, requests 1280 by 800 logical content on sufficiently large work areas, and never exceeds 90 percent of a smaller logical work area in either dimension.
3. **Given** the restored GUI, **When** the operator maximizes, restores, resizes, minimizes, closes, and relaunches it, **Then** title bar, borders, taskbar, and supported window-state transitions remain reachable and functional.
4. **Given** induced unavailable, access-denied, timeout, stream-disconnect, repeated-refresh, and retry conditions, **When** the product reconnects for at least two minutes per required observation, **Then** each incident occupies one persistent in-frame surface, no recurring modal or top-level error dialogs appear, access denial is accurately described, and successful recovery restores the interface.
5. **Given** the installed public interfaces, **When** the operator creates deterministic manual and scheduled tasks plus deliberate nonzero-exit and process-start failures, **Then** production execution produces the expected output, marker, history, exit status, run identity, and distinct failure diagnostics without a fake runner.

---

### User Story 3 - Prove the Attended Installer Lifecycle (Priority: P1)

As a release operator, I prove that the exact candidate's visible setup and teardown choices behave as documented for a normal desktop user and multiple supported profiles.

**Why this priority**: S039 supplied compiled and silent coverage, but #98 and coordinator #96 require attended observations before their lifecycle claims are releasable.

**Independent Test**: Exercise shortcut, completion, cancel, preserve, wipe, repair, upgrade, locked cleanup, and multiple-profile scenarios on the exact candidate, with populated product data and out-of-scope controls.

**Acceptance Scenarios**:

1. **Given** fresh attended setup, **When** the shortcut and completion controls are exercised, **Then** all supported selections have their documented defaults and effects, and finish-page GUI launch runs once as the unelevated interactive user.
2. **Given** populated machine data and supported-user preferences, **When** preserve uninstall is selected, **Then** software integration is removed and all declared data remains byte-for-byte available to reinstall.
3. **Given** populated owned data plus out-of-scope controls, **When** wipe is explicitly selected and confirmed, **Then** every verified owned root is removed, controls remain unchanged, shared security state follows the documented policy, and reinstall starts clean.
4. **Given** cancel, repair, upgrade, invalid input, locked content, or partial cleanup, **When** the lifecycle operation ends, **Then** destructive cleanup never runs in an unauthorized path and incomplete cleanup is reported rather than presented as success.

---

### User Story 4 - Block or Release the Exact Candidate (Priority: P1)

As a release maintainer, I receive a deterministic pass or block decision from reviewable evidence, and a passing decision makes only the exact tested Windows package eligible for publication.

**Why this priority**: Evidence without enforcement can be skipped accidentally and would repeat the v0.9.1 release-process failure.

**Independent Test**: Run the evidence validator against passing and systematically mutated bundles, then exercise the release workflow contract to prove missing or invalid evidence blocks before publication and valid evidence supplies the byte-identical MSI.

**Acceptance Scenarios**:

1. **Given** complete evidence, **When** the gate validates it, **Then** every requirement is accounted for, all required outcomes pass, referenced attachments are intact, and the decision identifies the exact candidate.
2. **Given** a failed, unavailable, skipped, timed-out, partial, stale, malformed, duplicated, or contradictory observation, **When** validation runs, **Then** the gate fails with all discovered defects and never converts absence into a pass.
3. **Given** a draft release without matching passing evidence, **When** promotion starts, **Then** publication is blocked before the release becomes public.
4. **Given** a draft release with matching passing evidence, **When** promotion prepares Windows assets, **Then** it retrieves, revalidates, and publishes the exact tested MSI rather than rebuilding it.

### Edge Cases

- A candidate artifact or draft release may be deleted, be renamed, or belong to a fork; each condition blocks rather than triggering a rebuild.
- A tag may identify a different commit or version than the evidence even when the MSI filename looks correct.
- Two evidence records may claim the same scenario or one required scenario may be omitted; duplicates and omissions both block.
- A screenshot may exist while native geometry, work-area, DPI, window-state, token, service, or visible-surface counts are absent; screenshots alone never satisfy structured requirements.
- Display scaling may yield fractional logical dimensions and a window frame larger than its content area; the gate evaluates logical content against the logical work area and separately retains native pixel rectangles.
- A monitor may disconnect or the primary monitor may change between launches; each observation records the monitor actually containing the window.
- Automatic reconnect may recover before two minutes; failure observations still need a controlled trigger long enough to demonstrate the no-spam interval, followed by a separate recovery observation.
- Destructive uninstall evidence may be partial because a file is locked; the result blocks wipe acceptance and retains actionable residual paths without exposing unrelated user data.
- A workflow rerun may have the same run number but a different attempt; origin identity includes the immutable run identifier and attempt.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Candidate staging MUST build a Windows MSI from an immutable tagged revision and MUST leave the release in draft state.
- **FR-002**: Candidate identity MUST include repository, source commit, workflow name and run identifier, run attempt, version, ProductCode, filename, byte size, and SHA-256.
- **FR-003**: The candidate manifest MUST be retained with the MSI and MUST be verifiable independently of attended evidence.
- **FR-004**: The evidence model MUST identify its schema version, candidate identity, clean-snapshot provenance, operator, timestamps, operating-system build, account and token context, service identity, display configuration, and cleanup disposition.
- **FR-005**: Evidence attachments MUST use relative paths within one bounded evidence bundle and MUST include byte size and SHA-256; absolute paths, parent traversal, missing files, and digest mismatches MUST be rejected.
- **FR-006**: Every required observation MUST have exactly one result from pass, fail, unavailable, skipped, timed-out, or partial, and only pass may satisfy release readiness.
- **FR-007**: Validation MUST report every discovered identity, schema, coverage, outcome, attachment, and cross-observation defect in one run and MUST return a nonzero status when any defect exists.
- **FR-008**: The attended workflow MUST prove first-use access from the intended non-elevated account to the installed LocalSystem service and denial for an unrelated local user.
- **FR-009**: Window evidence MUST record native top-level window rectangle, logical content size, native monitor work area, logical monitor work area, effective DPI, monitor identity, and restored, maximized, minimized, and fullscreen flags.
- **FR-010**: Fresh and retained-profile first-launch evidence MUST show a restored, non-fullscreen window with visible desktop margins and reachable title bar, resize borders, and taskbar.
- **FR-011**: On a sufficiently large logical work area, requested content MUST be 1280 by 800; on smaller work areas neither requested dimension may exceed 90 percent of the corresponding available logical dimension.
- **FR-012**: Window scenarios MUST include separate standard-DPI and high-DPI or mixed-DPI configurations plus maximize, restore, resize, minimize, close, and subsequent-launch observations.
- **FR-013**: Error evidence MUST cover daemon unavailable, access denied, timeout, stream disconnect, automatic reconnect or repeated refresh, manual retry, and successful recovery.
- **FR-014**: Each induced failure MUST be observed for at least two minutes where repetition could occur, record timestamps and visible in-frame, modal, and top-level surface counts, and show no recurring modal or top-level dialogs.
- **FR-015**: Error messages MUST distinguish access denied from daemon absence, MUST NOT recommend permanent elevation, and successful recovery MUST clear the incident without reinstalling.
- **FR-016**: Any new recurring-error correction MUST record the issue's complete baseline reproduction and uncommitted-fix native proof before its first corrective commit.
- **FR-017**: Task evidence MUST use public installed product interfaces and production process creation to prove deterministic manual and scheduled success, expected output, marker side effect, run identity, history, and exit code zero.
- **FR-018**: Task evidence MUST also prove a deliberate nonzero child exit and a process-start failure produce distinct accurate diagnostics.
- **FR-019**: Installer evidence MUST cover shortcut defaults and selections, independent completion actions, unelevated finish-page GUI launch, cancel behavior, repair, upgrade, preserve, wipe, locked cleanup, and multiple supported user profiles.
- **FR-020**: Preserve and wipe evidence MUST inventory populated owned state before removal, outcome state after removal, reinstall behavior, unaffected controls, and separately governed shared security state.
- **FR-021**: Operator tooling MUST use noninteractive hidden child-process execution for console programs and MUST avoid foreground or focus-stealing helper consoles.
- **FR-022**: Operator tooling MUST support safe resumption without silently overwriting prior evidence and MUST record cancellations or unavailable steps as non-passing outcomes.
- **FR-023**: The repository MUST include deterministic positive and negative evidence fixtures that exercise all validation rules without pretending fixture data is native release evidence.
- **FR-024**: The promotion controller MUST reject missing, malformed, stale, non-passing, or identity-mismatched attended evidence.
- **FR-025**: The promotion controller MUST verify the draft release tag, target commit, version, and exact candidate identity before any publication transition.
- **FR-026**: Windows release packaging MUST consume the revalidated exact candidate MSI; it MUST NOT rebuild or substitute a nominally equivalent MSI after native testing.
- **FR-027**: No release may become public before the candidate gate passes, and a gate failure MUST leave the release in draft state.
- **FR-028**: Documentation MUST define candidate staging, clean-snapshot setup, every attended scenario, evidence completion, local validation, evidence upload, promotion, artifact retention, and failure recovery.
- **FR-029**: Evidence records MUST avoid secrets and unnecessary personal data while retaining auditable account role, machine configuration, timestamps, and artifact origin.
- **FR-030**: Existing canonical unit, integration, race, GUI, coverage, documentation, automation, installer, and LocalSystem gates MUST remain required and pass.

### Key Entities

- **Candidate manifest**: Immutable identity and origin of one staged Windows MSI, including source, workflow, package metadata, size, and digest.
- **Evidence bundle**: Versioned structured results plus integrity-protected attachments for one candidate's attended Windows validation.
- **Observation**: One required native scenario with environment, start and end times, enumerated outcome, measurements, artifacts, and operator notes.
- **Environment profile**: Clean-snapshot, OS, account/token, service, monitor, work-area, and DPI context shared by related observations.
- **Release decision**: Deterministic pass or block result that lists candidate identity, satisfied requirements, and every defect.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Changing any candidate identity field or one byte of the MSI causes candidate validation to fail before publication.
- **SC-002**: The evidence gate accounts for 100 percent of required access, window, error, task, installer, and uninstall scenarios exactly once; every non-pass outcome blocks readiness.
- **SC-003**: Accepted window evidence contains both required DPI classes and 100 percent of required native and logical measurements, with every first-launch dimension satisfying the bounded-size rule.
- **SC-004**: Accepted error evidence records at least 120 seconds for every repetition-sensitive failure and zero recurring modal or top-level dialogs while preserving one reachable in-frame incident.
- **SC-005**: Accepted execution evidence contains two successful production runs (one manual and one scheduled), one nonzero exit, and one process-start failure with distinct recorded outcomes.
- **SC-006**: Accepted lifecycle evidence covers 100 percent of required attended setup, maintenance, preserve, wipe, cancellation, locked-item, multiple-profile, and reinstall outcomes without modifying any out-of-scope control.
- **SC-007**: Automated mutation tests reject every required missing field, duplicate scenario, invalid state, stale identity, unsafe attachment path, altered attachment, and altered candidate artifact.
- **SC-008**: A promotion attempt with no matching reviewed evidence or exact candidate fails before publication; a valid attempt supplies a byte-identical MSI and complete final checksums.
- **SC-009**: All eight canonical repository gates pass with the release-candidate gate added and no existing safety-critical coverage weakened.

## Assumptions

- Attended native evidence is collected by a trusted maintainer on a documented clean Windows 11 snapshot with permission to install and remove the candidate.
- The same clean snapshot may be reset and reconfigured between standard-DPI and high or mixed-DPI observations if each configuration is identified separately.
- A draft release is the candidate-staging mechanism. It is not public and its retained assets cover attended testing, review, and promotion.
- The exact evidence bundle is uploaded to the draft release, retained as a public release asset on promotion, and included in final checksums.
- Screenshots and logs support the structured record but do not replace required native measurements or explicit outcomes.
- S040 may complete its tooling and automated contracts without claiming the real candidate has passed. GitHub issues #94, #98, and #96 remain open until exact-candidate attended evidence actually satisfies their acceptance criteria.

## Dependencies and Traceability

- Primary implementation issue: #94.
- Parent coordinator: #96.
- Consumes the setup lifecycle implementation from #97 and the remaining attended acceptance from #98.
- Revalidates regression outcomes from #89, #90, #93, and #95 without reopening completed implementation work unless exact-candidate proof finds a reproducible failure.
- Blocks the next Windows release-readiness claim after v0.9.1.

## Out of Scope

- Publishing or tagging a release, merging the S040 pull request, or fabricating a passing candidate evidence record.
- Signing, notarization, certificate procurement, release promotion across external channels, or indefinite artifact archival.
- Replacing the scheduler, IPC authorization, installer, GUI toolkit, task persistence, or public product interfaces.
- Corrective product changes that lack required native baseline reproduction and proof.
- Treating automated fixtures, source inspection, hosted Windows Server CI, headless GUI tests, or screenshots alone as attended Windows 11 evidence.
