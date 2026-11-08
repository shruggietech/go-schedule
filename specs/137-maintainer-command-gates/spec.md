# Feature Specification: Maintainer command gates

**Feature Branch**: `codex/137-maintainer-command-gates`

**Created**: 2026-11-08

**Status**: Implemented

**Delivery**: Review branch codex/137-maintainer-command-gates; integrated #439 outcome and all eight local gates in [verification](verification.md); authorized official PR and owner merge remain delivery steps.

**Input**: Owner-approved S137 autopilot, [#439](https://github.com/shruggietech/go-schedule/issues/439), following accepted S136 at e41005b4dfbc16a54f8a3e8df77bd54583bdded6.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Normal verification includes maintainer contracts (Priority: P1)

A contributor uses the normal automation command and receives failures from lifecycle, public command and driver-stop regressions without remembering separate manual checks.

**Why this priority**: Existing regressions must participate in the delivery path that contributors already use.

**Independent Test**: Exercise the automation entry point with an intentionally failing contract and observe its failure result.

**Acceptance Scenarios**:

1. **Given** the canonical eight-gate interface, **When** automation runs, **Then** lifecycle and public-command regressions participate alongside existing automation policy checks.
2. **Given** a controlled early child failure, **When** the driver runs, **Then** the original exit status and diagnostic remain observable and later gates do not run.
3. **Given** broken stopping behavior under assessment, **When** its guard runs, **Then** it remains bounded and cannot execute the real suite or recurse.

### User Story 2 - Supported native hosts enforce the same command contract (Priority: P1)

Contributors receive the same public-command behavior on Windows, Linux and macOS through the existing native CI matrix.

**Why this priority**: Shell adapters cross a platform boundary that direct native-checker fixtures do not cover.

**Independent Test**: Invoke both public entry points against valid isolated roots through the focused native command corpus.

**Acceptance Scenarios**:

1. **Given** normal, caller-relative, spaced, balanced-bracket and unmatched-bracket roots, **When** either public command selects a root explicitly or by its caller directory, **Then** it inspects that literal root and reports its successful outcome.
2. **Given** the existing native matrix, **When** a supported platform runs, **Then** the same focused corpus enforces these contracts without another gate-set definition.

### Edge Cases

- An early failure must retain a non-one native status and useful stderr.
- A removed, bypassed or misplaced CI corpus invocation must not satisfy integration policy.
- Caller roots contain no checker implementation, so accidental compilation from the target root must fail visibly.
- A driver-stop regression must never execute real coverage, GUI, package or nested automation work in its guard.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Canonical automation MUST include the existing lifecycle corpus, both public-command contracts and bounded driver-stop guards alongside existing policy and coverage-profile checks.
- **FR-002**: The eight named gates and their order MUST remain format, vet, lint, race, gui, coverage, docs, automation.
- **FR-003**: Both public commands MUST honor explicit and caller-default literal roots across the specified path classes on Windows, Linux and macOS, with caller-relative selection covered explicitly.
- **FR-004**: Early-child failures MUST preserve original status and diagnostic, and the driver guard MUST remain finite even when stopping regresses.
- **FR-005**: Existing native CI MUST invoke the shared focused command corpus on all three platforms, with policy detecting removal, bypass or missing platform coverage.
- **FR-006**: Existing policy decisions, compiled-fixture reuse, hidden noninteractive Windows child processes, dependency baseline and timeouts MUST remain intact.
- **FR-007**: Contributor guidance MUST describe the integrated commands and the boundary between native contract checks and product/release qualification.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: One normal automation invocation includes all three maintainer-contract categories without manual supplemental commands.
- **SC-002**: Both public entry points cover every required root class on all three supported native hosts.
- **SC-003**: Controlled status 37 and its diagnostic survive driver failure, and no later child is invoked in that scenario.
- **SC-004**: Broken stopping behavior is rejected by a finite guard with zero real-suite child executions; the public list remains exactly eight gates.

## Assumptions

- Accepted S136 supplies the native checker and fixture decisions; this slice adds integrated enforcement rather than a new product capability.
- #439 is the only currently ready issue sharing this implementation boundary. Historical replay and clustered scheduling work remain deferred under existing authority.
- No new dependencies, timeout changes, release publication or installed-service changes are authorized. The owner authorizes push and official PR creation, two total Codex reviews including automatic review, and fixes through green CI, stopping before owner merge.
- Functional delivery governs issue completion. Tests, CI, evidence archives and attended walkthroughs support review and do not become issue-closure gates.
