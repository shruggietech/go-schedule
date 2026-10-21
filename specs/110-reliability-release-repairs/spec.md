# Feature Specification: Reliability and release-promotion repairs

**Feature Branch**: `codex/110-reliability-release-repairs`

**Created**: 2026-10-21

**Status**: Implemented

**Delivery**: [PR #338](https://github.com/shruggietech/go-schedule/pull/338), open and unmerged; [delivery](delivery.md), [local verification](verification.md) and [publication/review record](review.md).

**Input**: Owner-authorized S110 autopilot, bundling #262, #271 and #335 in full, with one local pre-publication halt.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Resume the intended release promotion (Priority: P1)

A maintainer promotes the exact authenticated staged release even when tag-specific lookup cannot expose its draft, and resumes an interrupted attempt without changing the release identity or staged contents. Priority follows #262's existing release blocker.

**Independent Test**: Model authenticated draft lookup and interrupted or already-published promotion without mutating a live release.

**Acceptance Scenarios**:

1. **Given** one authenticated matching draft unavailable through tag-specific lookup, **When** promotion starts, **Then** that exact release remains the identity used throughout the operation.
2. **Given** multiple matches, no match, or an identity mismatch, **When** promotion starts or resumes, **Then** it stops with a useful diagnostic before publication.
3. **Given** an interrupted checksum upload, **When** promotion resumes, **Then** matching checksums are reused and differing checksums are refused without overwriting staged bytes.
4. **Given** an already-published matching release with a stale latest pointer, **When** promotion resumes, **Then** the pointer is repaired without modifying its tag, artifacts, notes or checksums.

### User Story 2 - Keep exact-file replacements through observer recovery (Priority: P2)

An enabled exact-file watcher dispatches a completed atomic replacement that overlaps a transient observer failure and recovery without replaying a replacement already dispatched.

**Independent Test**: Control observer failure, replacement identity and recovery ordering and inspect dispatches, including late native hints.

**Acceptance Scenarios**:

1. **Given** a watched file with a known last-observed identity, **When** it is replaced during observer degradation and recovered, **Then** the replacement is dispatched once.
2. **Given** a replacement already pending or dispatched, **When** observer recovery and a late duplicate hint occur, **Then** recovery retains the event without duplicate dispatch.
3. **Given** a changed definition or disabled watcher, **When** recovery occurs, **Then** obsolete events are not dispatched and the current definition gets an appropriate baseline.

### User Story 3 - Persist every concurrent heartbeat invocation (Priority: P2)

Concurrent heartbeat payloads persist one durable record per successful invocation. Failed writers retain their actual diagnostics, and contention uses existing bounded waits and retries.

**Independent Test**: Reproduce a concrete failure in disposable state, inspect each writer's output and records, and exercise equivalent script contracts on supported hosts.

**Acceptance Scenarios**:

1. **Given** concurrent writers against one disposable record store, **When** they succeed, **Then** each has its independently persisted record.
2. **Given** temporary contention, **When** a writer retries, **Then** waits and retries remain bounded and the status accurately reports its outcome.
3. **Given** a non-contention failure, **When** a writer fails, **Then** its actual output remains available and distinguishable from contention.

### Edge Cases

- Drafts across multiple pages, duplicate exact tags, invalid IDs, changed tag identity, missing/unexpected assets, partial upload and differing checksums.
- Already-published recovery retains the same reviewed release and validates existing staged identity and bytes.
- Observer restart versus individual registration recovery, repeated recovery, missing files, changed definitions, pending debounce, late hints and dispatch rejection.
- Concurrent first-use setup and writes, an initialized store, malformed helper output, unavailable tools, exhausted contention and unrelated failures.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Promotion MUST select exactly one authenticated release matching the exact requested tag and bind its identity through every later operation (#262).
- **FR-002**: Promotion MUST refuse missing, ambiguous or changed identity before mutating publication state (#262).
- **FR-003**: Promotion MUST retain successful staging-attempt, exact candidate, eight-asset allowlist, complete checksum, immutable-tag and disclosure requirements (#262).
- **FR-004**: Interrupted promotion MUST resume without overwriting differing checksums or publishing incomplete assets; already-published recovery MUST repair the latest pointer without changing contents (#262).
- **FR-005**: Operator guidance MUST name the working path and refusal conditions (#262).
- **FR-006**: An enabled exact-file watcher MUST dispatch a completed replacement overlapping observer failure and recovery exactly once, retaining its prior meaningful identity until reconciliation (#271).
- **FR-007**: Recovery MUST retain valid pending events and suppress already-dispatched replacements and late duplicate hints (#271).
- **FR-008**: Definition changes, disablement and shutdown MUST retain lifecycle, health, bounded-memory and generation-isolation contracts (#271).
- **FR-009**: Successful concurrent heartbeat invocations MUST each persist an independently identifiable record, including concurrent first use (#335).
- **FR-010**: Contention MUST retain existing bounded wait/retry behavior, distinguish unrelated failures and retain failed-writer diagnostics (#335).
- **FR-011**: The correction MUST preserve equivalent Bash/PowerShell behavior, documented exit codes and existing concurrent success and record-count assertions (#335).
- **FR-012**: Each correction MUST trace to demonstrated defective behavior; repeated passes alone MUST NOT be represented as a corrected intermittent cause (#271, #335).
- **FR-013**: Bound payload parameter values MUST retain exact text, including numeric-looking identifiers, quotes, newlines, empty values and international text, while explicit nulls remain null (#335).

### Key Entities

- **Release identity**: Exact tag, authenticated stable identifier, publication state, staged candidate and asset identities.
- **Observed file identity**: Last meaningful version, definition revision, pending replacement and previously dispatched replacement.
- **Writer outcome**: Invocation exit status, actual diagnostic output and persisted record identity.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Every promotion/resumption path selects one identity and preserves reviewed contents or refuses without publication.
- **SC-002**: Each completed replacement crossing recovery produces exactly one dispatch, with zero duplicate replays from recovery or late hints.
- **SC-003**: Six simultaneous successful heartbeat invocations produce six distinct records; failures retain real diagnostics and bounded retries.
- **SC-004**: All three issues' functional outcomes are independently accounted for, without skipped assertions or deferred acceptance criteria presented as completion.

## Assumptions

- Baseline is merged S109 at 95c05692e41bdb3630aa0b3c1bf946e5d443c2af. Existing authentication, owned fixtures and hosted native macOS checks remain available.
- No persisted application schema, notification capability, authorization policy or distributed-scheduling behavior changes.
- Fixtures model release mutations; actual promotion, tagging, publication and hosted review remain outside local implementation authorization.
- Existing issue milestones and dependencies remain authoritative; #335 remains a child of #277.

## Clarifications

### Session 2026-10-21

- Q: Does release repair authorize live promotion? A: No; fixtures cover mutation paths, and any real promotion needs separate release authorization.
- Q: Can passing retry runs establish intermittent defect completion? A: No; establish the concrete failure and retain the original concurrent and dispatch assertions.
- Q: Should unrelated roadmap features join this slice? A: No; the authorized full bundle is #262, #271 and #335.
