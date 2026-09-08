# Feature Specification: Dependable Webhook Notifications

**Feature Branch**: `codex/067-webhook-notifications`

**Created**: 2026-09-07

**Status**: Implemented

**Delivery**: Implemented and locally verified on 2026-09-07. Pull-request CI and external review remain the hosted delivery gates.

**Input**: User description: "Define reusable notification policies and deliver dependable webhook notifications for task and group run outcomes, satisfying GitHub issues #158 and #159."

## User Scenarios & Testing

### User Story 1 - Send a safe run-outcome webhook (Priority: P1)

As an operator, I can create a reusable webhook channel and assign success or failure notifications to a task so an ordinary JSON receiver learns about selected terminal outcomes without gaining access to task secrets.

**Why this priority**: A useful notification feature begins with one complete, secure path from a recorded run to an external receiver.

**Independent Test**: Create a channel against a local receiver, assign it to a task for one outcome, run the task, and validate the documented payload and delivery record while proving that command, arguments, environment, standard input, captured output, authorization values, and URL credentials are absent.

**Acceptance Scenarios**:

1. **Given** an enabled webhook channel assigned to a task for success, **When** that task records a successful terminal run, **Then** the receiver obtains one versioned JSON payload correlated to that task and run.
2. **Given** a channel assigned only for failure, **When** the task succeeds, **Then** no delivery is created or attempted.
3. **Given** a webhook URL or authorization value containing credentials, **When** any channel or delivery is returned, logged, streamed, exported, or captured as evidence, **Then** the secret value and credential-bearing URL are absent.

---

### User Story 2 - Apply task and group policy predictably (Priority: P1)

As an operator managing related tasks, I can assign notification policies to a group and selectively replace them at a task so the effective channels and conditions are deterministic and explainable.

**Why this priority**: Group policy is essential for operating more than a few tasks, but an unclear merge rule would create duplicate or surprising notifications.

**Independent Test**: Assign policies at an ancestor group, descendant group, and task, query the effective policy, and show that the nearest non-empty scope fully replaces inherited assignments while an empty scope inherits.

**Acceptance Scenarios**:

1. **Given** a task with no direct notification assignments, **When** its group chain has assignments, **Then** the nearest ancestor with assignments supplies the complete effective policy.
2. **Given** a task with one or more direct assignments, **When** its effective policy is queried or a run completes, **Then** only the task assignments apply and inherited assignments do not also fire.
3. **Given** a nested group with its own assignments, **When** a descendant task has no direct assignments, **Then** the nested group's assignments replace those of more distant ancestors.

---

### User Story 3 - Survive delivery failure and restart (Priority: P1)

As an operator, I can rely on bounded asynchronous retries and durable redacted history so slow, unreachable, or rejecting webhooks remain observable without blocking task workers or changing task results.

**Why this priority**: Notification transport is less trustworthy than local scheduling and must never become part of task correctness.

**Independent Test**: Exercise timeouts, rejected responses, daemon cancellation, and restart with controlled receivers; validate finite attempts, recovery of interrupted work, terminal delivery evidence, and unchanged task outcomes.

**Acceptance Scenarios**:

1. **Given** an unreachable, slow, or rejecting receiver, **When** a matching run completes, **Then** delivery attempts use a 5-second per-request timeout, at most three total attempts, and backoff delays of 1 then 2 seconds.
2. **Given** notification failure or delay, **When** the task result is queried, **Then** its recorded outcome, timing, and output metadata remain exactly those produced by task execution.
3. **Given** a daemon stop during an in-progress delivery, **When** the daemon restarts, **Then** the interrupted delivery returns to pending and continues only within its existing three-attempt budget.
4. **Given** a receiver accepted a request but its response was lost, **When** retry occurs, **Then** the repeated request carries the same stable delivery identifier and is documented as possible at-least-once delivery.

---

### User Story 4 - Manage and inspect webhook channels (Priority: P2)

As an operator, I can create, test, update, disable, rotate, list, and remove webhook channels and inspect bounded delivery history through the local API and CLI.

**Why this priority**: Lifecycle and evidence operations make the channel maintainable after its first successful delivery.

**Independent Test**: Exercise the documented local API and CLI lifecycle against an isolated daemon, including secret replacement, disablement, channel deletion, and history filtering.

**Acceptance Scenarios**:

1. **Given** a new or updated channel, **When** the operator supplies an HTTPS endpoint and optional authorization value, **Then** ordinary responses return only a redacted endpoint summary and whether authorization is configured.
2. **Given** a test request, **When** it completes, **Then** a delivery record distinguishes it from task-run deliveries and reports its bounded result.
3. **Given** a disabled channel, **When** a matching run completes, **Then** no new delivery is created for that channel.
4. **Given** a channel is removed, **When** deletion succeeds, **Then** its assignments and pending work are removed while completed delivery history remains with a deleted-channel snapshot and no secret.

### Edge Cases

- Endpoint URLs containing user information or sensitive query parameters are accepted only for local submission and are represented elsewhere by scheme plus host with query and user information removed.
- Plain HTTP endpoints are rejected except loopback hosts, which remain available for local development and deterministic installed-daemon tests.
- Redirects are not followed, preventing a validated destination from silently forwarding authorization to another host.
- Success means an HTTP status from 200 through 299; every other response is a failed attempt and only a bounded diagnostic excerpt may be retained.
- A task or group is removed while deliveries exist; pending work cascades away, while completed history retains immutable task, group, channel, and outcome snapshots.
- A channel is disabled or changed after delivery creation; existing deliveries use their immutable redacted destination and payload snapshots, but resolve without sending if their protected credential reference no longer exists.
- Multiple channels in one effective policy create independent deliveries, each with its own retry budget and stable identifier.
- Daemon shutdown waits only for the normal bounded shutdown window; claimed deliveries not completed are recovered on the next start.
- Delivery history is pruned to the newest 1,000 terminal records while pending or claimed records are never removed by retention.

## Requirements

### Functional Requirements

- **FR-001**: The system MUST model reusable notification channels independently from tasks, groups, task runs, and notification deliveries.
- **FR-002**: The first supported channel kind MUST be webhook and MUST carry a display name, enabled state, endpoint, optional authorization value, and timestamps.
- **FR-003**: The local API and CLI MUST support create, list, get, update, enable, disable, rotate authorization, test, and remove operations for webhook channels.
- **FR-004**: Ordinary API responses, CLI output, logs, live events, exports, diagnostics, and verification evidence MUST never expose endpoint user information, endpoint query values, authorization values, or protected credential references.
- **FR-005**: Channel secrets MUST be stored outside task environment fields as write-only values in the daemon-owned SQLite store, protected by Windows DPAPI under the daemon service identity on Windows and daemon-only directory and database permissions on Unix; no custom cryptography or broader cross-platform encryption-at-rest claim is permitted.
- **FR-006**: Channel create and update MUST accept HTTPS destinations and loopback HTTP destinations only, MUST reject fragments, non-HTTP schemes, missing hosts, and non-loopback HTTP hosts, and MUST never follow redirects.
- **FR-007**: A notification assignment MUST bind one channel to exactly one task or group and select at least one terminal condition from success and failure.
- **FR-008**: A scope MAY have multiple channel assignments but MUST NOT have duplicate assignments for the same channel.
- **FR-009**: Effective-policy precedence MUST use the nearest non-empty scope as a full replacement boundary: task assignments replace all group assignments; otherwise the nearest ancestor group with assignments replaces more distant ancestors; no assignments means no notifications.
- **FR-010**: The API MUST expose both configured assignments and the effective task policy, including the source scope and selected conditions, without exposing secrets.
- **FR-011**: Recording a terminal success or failure run and creating its matching notification deliveries MUST occur in one database transaction before any outbound request begins.
- **FR-012**: Notification work MUST run outside the task worker semaphore so network latency, retries, and receiver failure cannot occupy or block scheduler execution capacity.
- **FR-013**: A delivery MUST be a durable record distinct from its source run and MUST include a stable delivery ID, source run and task IDs, immutable safe payload and destination snapshots, state, attempt count, timestamps, last status, and bounded diagnostic text.
- **FR-014**: A webhook payload MUST use schema identifier `go-schedule.webhook.v1` and include delivery, daemon, task, run outcome, trigger, schedule, and timing data while excluding command, arguments, working directory, environment, standard input, run output, channel authorization, and credential-bearing URLs.
- **FR-015**: Each webhook request MUST use `Content-Type: application/json`, `User-Agent: go-schedule/<version>`, `X-Go-Schedule-Delivery`, and `X-Go-Schedule-Event: run.completed`; optional authorization MUST be sent only to the validated original destination.
- **FR-016**: A test delivery MUST use the same transport and evidence path as a run delivery, MUST be distinguishable by event kind, and MUST not fabricate or mutate a task run.
- **FR-017**: Each request attempt MUST have a 5-second timeout, success MUST be any 2xx status, and failure MUST retry at most twice after the initial attempt with delays of 1 and 2 seconds.
- **FR-018**: Retry behavior MUST provide at-least-once rather than exactly-once delivery; every retry of one delivery MUST reuse its stable ID so receivers can deduplicate.
- **FR-019**: On daemon startup, claimed but unfinished deliveries MUST return to pending without resetting their attempt counts; deliveries already at three attempts MUST resolve as exhausted without another request.
- **FR-020**: Notification failure, delay, restart recovery, history pruning, and channel lifecycle operations MUST NOT modify the source run's outcome or prevent the scheduler from accepting other execution work.
- **FR-021**: Delivery history queries MUST support channel, task, run, and state filters and a bounded limit no greater than 1,000, ordered newest first.
- **FR-022**: The daemon MUST retain at most the newest 1,000 terminal delivery records after each terminal transition, excluding pending and claimed work from pruning.
- **FR-023**: Removing a channel MUST atomically remove its assignments and non-terminal deliveries, erase its secret, and preserve terminal history using immutable safe channel and destination snapshots.
- **FR-024**: The implementation MUST use Go's maintained standard `net/http`, `net/url`, and `encoding/json` packages for outbound transport and payload handling, SQLite transactions for durability, the already-pinned `golang.org/x/sys/windows` DPAPI binding on Windows, and operating-system file and IPC permissions; it MUST add no message broker, custom cryptography, credential-vault dependency, or vendor-specific SDK.
- **FR-025**: Unit, API contract, integration, restart-recovery, race, and installed-daemon validation MUST cover successful delivery, filtering, policy precedence, timeout or cancellation, rejection, secret redaction, retry exhaustion, duplicate identifiers, channel removal, and task-outcome isolation.
- **FR-026**: User documentation MUST publish the versioned payload, receiver example, lifecycle commands, precedence rule, security boundary, retry schedule, duplicate risk, retention rule, and common generic receiver guidance.
- **FR-027**: The completed slice MUST preserve issue-level traceability for #158 and #159 and report any unmet acceptance criterion instead of closing it implicitly.

### Key Entities

- **Notification channel**: A reusable named outbound destination with kind, enabled state, safe endpoint summary, protected endpoint and authorization values, and lifecycle timestamps.
- **Notification assignment**: A task-scoped or group-scoped binding between a channel and one or more terminal run conditions.
- **Effective notification policy**: The deterministic assignment set selected for a task from its own scope or nearest configured ancestor group.
- **Notification delivery**: Durable asynchronous work and redacted evidence for one channel, one event, and optionally one immutable source run.
- **Webhook payload**: The stable receiver-facing JSON document stored as an immutable safe snapshot for retry.

## Success Criteria

### Measurable Outcomes

- **SC-001**: A generic local receiver can validate every documented webhook field and correlate the payload to exactly one delivery and source run without importing go-schedule code.
- **SC-002**: In an integration run with all task workers occupied and webhook receivers delayed for five seconds, an additional eligible task begins as soon as task capacity is available with no notification worker occupying that capacity.
- **SC-003**: Every unreachable or rejecting delivery reaches a terminal failed state after exactly three or fewer attempts, and no attempt begins more than eight seconds after its delivery becomes eligible, excluding operating-system scheduling delay.
- **SC-004**: Restart tests preserve the delivery ID and attempt count and never exceed three outbound requests for one delivery across daemon lifecycles.
- **SC-005**: Automated redaction tests find zero submitted secret values across channel, assignment, effective-policy, delivery, log, event, export, and evidence representations.
- **SC-006**: Existing task-run tests and the canonical eight-gate verifier pass unchanged in outcome semantics, including the race detector and supported desktop gate.
- **SC-007**: The local API and CLI complete every channel and assignment lifecycle operation, including removal with terminal-history preservation, in installed-daemon integration tests.

## Assumptions

- Windows service identity remains stable across daemon restarts so its DPAPI-protected values remain decryptable. Unix service data remains daemon-only at directory mode 0700 and database mode 0600. Full-disk protection beyond these boundaries remains platform administration and is not claimed by this feature.
- Success and failure are the only notification conditions in S067. Consecutive failures, duration thresholds, and failure-to-start can extend the assignment condition vocabulary later without changing channel or delivery identity.
- Webhook receivers implement idempotency when duplicate side effects matter by keying on `X-Go-Schedule-Delivery` or the payload delivery ID.
- Desktop channel management belongs to #160 and is outside S067; the local API and CLI provide the complete management surface in this slice.
- The existing authenticated local IPC boundary remains the authority for channel and assignment management. No remote management listener is introduced.

## Dependencies and Traceability

- Parent: #19.
- Completes: #158 and #159 when all acceptance criteria and verification gates pass.
- Depends on the shipped daemon, SQLite store, task/group hierarchy, run history, local IPC access controls, and service packaging baseline.
- Blocks: #160 notification-management desktop UI.
