# Feature Specification: Windows Demo Qualification

**Feature Branch**: `codex/043-windows-rc-qualification`

**Created**: 2026-09-03

**Status**: Implemented

**Delivery**: review branch `codex/043-windows-rc-qualification`; S043 demo SHA-256 `ef4a869af0e6971d445c53e8f7a237df655245287c7ae64d8e719941bda8ad59`; canonical verification, attended walkthrough, and post-wipe audit recorded in `verification.md`

**Input**: Qualify the Windows release work as far as possible without operator help, then provide a locally built demo installer for attended testing before publishing a pull request.

## Scope and evidence boundary

S043 packages the merged S039 through S042 Windows work into one locally built, identity-pinned demo MSI, performs every safe unattended and headless check available on the development host, and hands that exact file to the operator for attended desktop testing. It is traceable to #94, #98, #96, #101, #104, #105, and #106.

The demo MSI is deliberately not represented as the formal release candidate. The formal #94 gate requires an MSI staged by the tag-triggered Release workflow from a reviewed and merged commit. Because the operator requires testing before the S043 pull request exists, the byte-identical staged candidate can only be produced after that review path. Passing demo evidence reduces risk but does not close #94, #98, or #96 by itself.

### Scope in

- Create a reproducible local Windows demo build from the exact S043 commit.
- Embed an unmistakable non-release build version while retaining the intended numeric MSI ProductVersion.
- Inspect the compiled MSI and record its source commit, ProductCode, size, and SHA-256.
- Run canonical, installer-source, compiled-MSI, GUI, API, executor, and release- gate verification that does not require destructive elevation or human visual judgment.
- Prepare a compact attended checklist covering the remaining #94, #98, and S042 native observations.
- Record failures before making product corrections. For error-spam work, retain #94's proof-before-commit ordering.
- Keep the branch local and withhold the pull request until the operator reports that demo testing is complete.

### Scope out

- Pushing a branch or opening a pull request before the operator completes demo testing; tagging, staging or publishing a GitHub release; or claiming formal release-candidate qualification.
- Closing issues whose native or exact-candidate evidence remains incomplete.
- Implementing #102 unless testing proves its current behavior blocks the v1 release gate rather than remaining a documented Post-v1 diagnostic improvement.
- Treating headless widget assertions, compiled tables, or a local demo build as substitutes for attended Windows 11 evidence.

## Clarifications

### Session 2026-09-03

- Q: Can the pre-PR demo MSI satisfy the exact-candidate gate in #94? -> A: No. The Release workflow can stage the formal candidate only from a pushed tag after reviewed source is merged. Demo findings are exploratory qualification evidence and must be repeated against the staged candidate.
- Q: Which release identity should the demo display? -> A: Use `1.0.0-s043-demo.<commit>` in binaries and `1.0.0` as the numeric MSI ProductVersion. The filename and evidence must say `s043-demo`; it must not be named like a published release artifact.
- Q: May local automation install or remove the MSI? -> A: Only a repository harness that positively proves it is running elevated on a disposable host may do so unattended. This development host will use non-destructive compilation, table inspection, tests, and operator-driven installation.
- Q: When may the S043 branch leave the machine? -> A: Only after the operator completes demo testing and explicitly authorizes publication.
- Q: What is the disposition of UI complaints found in the completed demo? -> A: File or update GitHub issues only. Do not expand S043 product scope. The operator authorized publication through CI and pull-request review, with merge still reserved for the final review ritual.

## User Scenarios & Testing

### User Story 1 - Receive an identifiable demo installer (Priority: P1)

As the maintainer, I can open one clearly labeled local MSI whose identity and hash are recorded, so I know exactly what I am testing.

**Independent Test**: Compare the linked MSI with the recorded SHA-256, inspect its ProductVersion and ProductCode, and confirm the embedded application version contains the S043 demo marker.

**Acceptance Scenarios**:

1. **Given** the locally committed S043 source, **when** the build completes, **then** one demo MSI and one inspection report identify the same source and artifact hash.
2. **Given** the demo MSI, **when** it is distinguished from a release artifact, **then** its filename, origin, and embedded binary version all identify it as non-release testing material.

---

### User Story 2 - Consume automated qualification first (Priority: P1)

As the maintainer, I receive a demo only after all non-attended checks available on the host have passed, so my time is reserved for native interactions and visual judgment.

**Independent Test**: Run the recorded commands and confirm canonical gates, focused Windows contracts, compiled-MSI inspection, and release-gate fixtures all pass against the same source tree and artifact.

**Acceptance Scenarios**:

1. **Given** a demo build, **when** any automated qualification check fails, **then** the installer is not handed off as ready and the failure is retained.
2. **Given** all automated checks pass, **when** the installer is handed off, **then** the verification record lists exact commands and outcomes without claiming attended proof.

---

### User Story 3 - Complete a bounded attended walkthrough (Priority: P1)

As the maintainer, I can follow a concise checklist for the interactions only a real Windows desktop user can credibly verify before the pull request is opened.

**Independent Test**: Install the exact linked MSI and report each required setup, GUI, task, and uninstall observation against its recorded hash.

**Acceptance Scenarios**:

1. **Given** a Windows 11 desktop, **when** the demo is installed interactively, **then** shortcut options, completion actions, first launch, navigation, appearance, storage disclosure, and task interactions can be verified.
2. **Given** populated application data, **when** maintenance removal is entered from Windows Settings, **then** preserve and wipe choices, cancellation, and the actual cleanup result can be reported separately.
3. **Given** a failed observation, **when** the result is returned, **then** it is recorded before any corrective implementation or commit.

## Edge Cases

- A rebuilt MSI can have a different hash even when source is unchanged; every rebuild replaces the demo identity and invalidates observations against the prior file.
- The development host may not be elevated or disposable. Destructive lifecycle automation must remain unavailable rather than being counted as a pass.
- Existing application data or Fyne preferences can contaminate fresh-launch observations; the checklist separates fresh and retained-profile cases.
- Windows scaling and monitor geometry vary. Visual checks record scale and work area rather than imposing a screen-percentage rule on every user.
- Installer cancellation, preserve, and wipe are different outcomes and cannot share one generic pass statement.
- A passing local demo does not prove the later staged MSI is byte-identical.

## Requirements

### Functional Requirements

- **FR-001**: S043 MUST produce one local MSI explicitly labeled as an S043 demo.
- **FR-002**: The demo MUST be built from a recorded full Git commit plus any separately enumerated evidence-only files, with product binaries built from committed source.
- **FR-003**: The GUI, daemon, CLI, and cleanup helper MUST embed a non-release S043 demo version derived from the source commit.
- **FR-004**: The MSI MUST use numeric ProductVersion `1.0.0` while its filename and evidence origin prevent confusion with the formal v1.0.0 MSI.
- **FR-005**: Compiled inspection MUST verify installer identity, summary data, shortcut features, completion actions, maintenance routing, preserve/wipe properties, cleanup conditions, and candidate hash.
- **FR-006**: All eight canonical repository gates MUST pass before handoff.
- **FR-007**: Focused release-gate, installer, cleanup, GUI, API, daemon, and real Windows executor tests MUST pass before handoff.
- **FR-008**: Verification MUST distinguish passed, failed, unavailable, and operator-pending evidence without converting missing native evidence to pass.
- **FR-009**: The handoff MUST link the exact MSI and report SHA-256, byte size, ProductVersion, ProductCode, source commit, and embedded demo version.
- **FR-010**: The attended checklist MUST cover setup choices, completion actions, native sizing and DPI, recurring-error behavior, real manual and scheduled tasks, S042 GUI acceptance, and uninstall preserve/wipe/cancel outcomes.
- **FR-011**: No S043 branch push or pull request may occur before the operator declares demo testing complete and separately authorizes publication.
- **FR-012**: No issue requiring exact-candidate or attended proof may close from local demo evidence alone.
- **FR-013**: Any product correction discovered by testing MUST include a failing regression test and, for recurring error presentation, comply with #94's baseline-reproduction and uncommitted-proof sequence before commit.
- **FR-014**: The demo handoff MUST state that the formal staged-candidate matrix remains a later release gate even after exploratory demo success.

### Key Entities

- **Demo Candidate**: The local MSI plus source identity, product identity, embedded version, hash, size, and build timestamp.
- **Automated Qualification Record**: Commands, outcomes, environment boundaries, and attachments produced without attended interaction.
- **Attended Observation**: One operator-reported behavior tied to the demo hash, display/account context, expected result, and pass/fail state.
- **Formal Candidate**: A later Release-workflow MSI staged from a reviewed tag; only this entity can satisfy #94's byte-identity requirement.

## Success Criteria

### Measurable Outcomes

- **SC-001**: The linked demo MSI's computed SHA-256 exactly matches the handoff.
- **SC-002**: All eight canonical gates and every selected focused check pass.
- **SC-003**: Compiled inspection reports no failed installer contract.
- **SC-004**: The operator can complete the condensed checklist without needing repository implementation knowledge.
- **SC-005**: Zero pushes, pull requests, tags, or releases occur before attended demo testing is complete.
- **SC-006**: Every unresolved native or formal-candidate requirement remains explicitly open rather than being implied complete.

## Assumptions

- The next intended release line is v1.0.0, consistent with the active milestone and installer test versions; this slice does not authorize cutting it.
- The operator will perform attended installation and visual checks on the current Windows desktop, while the agent performs all headless work.
- A local build is suitable for exploratory approval before PR publication but cannot replace the later draft-release artifact.
