# Tasks: Actionable Notification Conditions

## Phase 1: Specification and contracts

- [x] T001 Define S095 user stories, condition semantics, daemon-health boundary, requirements, and traceability in specs/095-notification-conditions/spec.md
- [x] T002 Record technical decisions, persistence, contracts, research, data model, quickstart, and specification checklist in specs/095-notification-conditions/

## Phase 2: Durable condition foundation

- [x] T003 Add failing domain, executor, migration v20, and store tests for structured start failure, compatibility migration, condition state, transition precedence, reminders, recovery, policy reset, and heartbeat scheduling
- [x] T004 Extend domain models, executor evidence, SQLite schema, scans, and assignment validation for advanced conditions and safe delivery metadata
- [x] T005 Implement transactional condition evaluation, policy fingerprints, durable state transitions, and heartbeat creation in the existing store and dispatcher lifecycle

## Phase 3: User Story 1, actionable task problems

**Goal**: Alert on consecutive failure, failure to start, and excessive duration through direct and inherited policies.

**Independent Test**: Configure each problem condition, record controlled task runs, and inspect one correctly attributed delivery per transition.

- [x] T006 [US1] Extend API server, client, authorization-compatible MCP Manage input, and contract tests for complete advanced assignments
- [x] T007 [US1] Extend CLI policy flags, parsing, human output, JSON compatibility, and tests for thresholds and duration conditions
- [x] T008 [US1] Add integration coverage proving task and group condition evaluation remains atomic and isolated from task outcome

## Phase 4: User Story 2, suppression and recovery

**Goal**: Persist active problems, suppress duplicates, send bounded reminders, and report recovery.

**Independent Test**: Advance an injected clock through problem, early repeat, reminder, restart, and recovery transitions.

- [x] T009 [US2] Add deterministic clock injection and store tests for reminder, quiet-period, restart, and policy-fingerprint reset behavior
- [x] T010 [US2] Extend webhook payload and delivery evidence with safe condition kind, explanation, streak, threshold, and reminder metadata

## Phase 5: User Story 3, daemon presence

**Goal**: Provide opt-in durable healthy-presence heartbeats without unsupported outage claims.

**Independent Test**: Advance the dispatcher clock across due boundaries and restart while proving one non-terminal heartbeat maximum per channel.

- [x] T011 [US3] Implement channel heartbeat configuration, persistence, due-work creation, dispatcher wake scheduling, and tests
- [x] T012 [US3] Add heartbeat API and CLI controls plus documented receiver-side missing-heartbeat behavior

## Phase 6: User Story 4, configuration and explanation

**Goal**: Make advanced conditions configurable and understandable in the desktop.

**Independent Test**: Configure and inspect direct and inherited rules by keyboard at supported compact and zoomed layouts.

- [x] T013 [US4] Extend desktop Go projections, drafts, validation, coverage summaries, and service tests with complete condition configuration and evidence
- [x] T014 [US4] Extend frontend models, store, assignment editor, delivery summaries, heartbeat controls, and component tests
- [x] T015 [US4] Add accessible responsive styling and Playwright coverage for advanced condition configuration and explanations

## Phase 7: Documentation and completion

- [x] T016 Update docs/notifications.md, docs/cli.md, CHANGELOG.md, and specs/README.md with condition and heartbeat behavior
- [x] T017 Run focused Go, frontend, accessibility, Playwright, integration, race, formatting, coverage, and full verification suites; record results in specs/095-notification-conditions/verification.md
- [x] T018 Mark all tasks complete, transition the specification to Implemented, update issue #175 and parent #19, and prepare the S095 pull request

## Dependencies

- Phase 2 blocks all user stories.
- User Stories 1 and 2 share transactional evaluation and execute sequentially.
- User Story 3 reuses the durable delivery pipeline after Phase 2.
- User Story 4 follows the stable API and projection contracts from User Stories 1 through 3.
- Phase 7 begins after every functional story passes its focused tests.

## Implementation Strategy

Build the data-compatible state machine first, expose it through API and CLI, add heartbeats through the existing dispatcher, then complete the desktop explanation and documentation. Preserve current threshold-one failure behavior at every intermediate step.
