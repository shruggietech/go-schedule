# Research: Task and Group Authoring

## Decision 1: Add detailed task-list negotiation to the existing endpoint

**Decision**: `GET /v1/tasks?details=true` returns a `TaskResponse` array under the existing `tasks` envelope. Requests without the flag retain the current `domain.Task` array exactly.

**Rationale**: The overview requires readiness, schedule summary, policy summary, and next runs for one hundred tasks. Fetching each detail separately creates avoidable IPC serialization and unstable partial loading. A query-gated response is backward compatible and reuses the authoritative detail builder.

**Alternatives considered**: One request per row was rejected as an N+1 design. A new `/v2` route was rejected because the representation is additive and opt-in. Duplicating readiness or occurrence calculation in the desktop was rejected because it would drift from daemon behavior.

## Decision 2: Use one desktop application service for safe task/group contracts

**Decision**: Add a `desktop/taskgroup` package with a narrow backend interface, local client adapter, transport-neutral models, and service methods used by the Wails facade.

**Rationale**: The package keeps daemon types and errors out of the frontend, centralizes full-path/effective-state derivation and stale checks, gives tests a deterministic seam, and matches S061's connection separation.

**Alternatives considered**: Binding the root API client directly to Wails was rejected because it exposes backend types and errors. Reimplementing domain logic in TypeScript was rejected for correctness. Extending the connection manager with CRUD was rejected because connection lifecycle and feature operations have different ownership.

## Decision 3: Keep writes pessimistic and add stale draft timestamps

**Decision**: Entity records change only after accepted daemon responses and refresh. Edit drafts carry the source `updated_at`; the service compares it to current authority and returns a field-neutral stale result unless the operator explicitly overwrites.

**Rationale**: The daemon remains authoritative, duplicate and late responses cannot fabricate state, and live updates can preserve dirty drafts without silent loss.

**Alternatives considered**: Optimistic mutation was rejected because rollback across task schedules and group cascades is error-prone. Silent last-write-wins was rejected because it can discard another client's change. Full ETag protocol expansion was rejected as disproportionate for local v1.2 operation.

## Decision 4: Centralize the safe example in the root module

**Decision**: `internal/commandexample` owns exact allowlisted mappings for Windows (`cmd.exe /d /c ver`), macOS (`/usr/bin/sw_vers`), and Linux (`uname -a`), including parsed program/arguments, explanation, and output recognition.

**Rationale**: Both Wails and Fyne can import the root package, documentation can be contract-checked against it, and native CI can execute the exact current-platform mapping under a five-second deadline. The commands are installed platform utilities, read-only, noninteractive, and observable in captured output.

**Alternatives considered**: Separate UI constants were rejected as drift-prone. Python was rejected as optional and the prior example opened a listener. Shell scripts and file-writing heartbeat examples were rejected as platform-specific or mutating.

## Decision 5: Parse command and environment input in Go

**Decision**: The desktop service uses the existing direct command-line parser and accepts environment input as ordered `KEY=value` rows that it validates into the existing map contract.

**Rationale**: Exact parsing already handles platform-neutral quoting and preserves the no-shell guarantee. Central validation produces safe field-specific failures shared by tests and the frontend.

**Alternatives considered**: Browser-side shell parsing was rejected as duplicate security-sensitive grammar. Separate program and argument controls were rejected because current users edit a single natural command line and #192 requires a one-keystroke example.

## Decision 6: Preserve full task values through differential updates

**Decision**: For edits, the service retrieves current detail, compares each submitted value, sends only intended changes, and uses new explicit clear flags for working directory and run identity where the existing partial request cannot currently represent clearing.

**Rationale**: Resubmitting schedule text on every edit would create replacement schedules and could shift timing. Differential updates preserve untouched data and make clearing every exposed field honest.

**Alternatives considered**: Sending the full form on every patch was rejected because omitted-versus-clear semantics differ by field. Leaving fields uncleareable was rejected because it violates complete editor fidelity. Replacing PATCH with PUT was rejected as unnecessary protocol expansion.

## Decision 7: Use accessible native controls and bounded rendering

**Decision**: Build the route from semantic tables, buttons, forms, details, dialogs, and a nested group list, with client-side filtering over one bounded workspace snapshot and no new component dependency.

**Rationale**: Native controls provide keyboard and assistive semantics, the S061 primitives already cover focus and announcements, and one hundred rows do not warrant virtualization complexity.

**Alternatives considered**: A tree-grid dependency was rejected because its accessibility and bundle cost exceed the scale. Card-only layouts were rejected because comparisons become slow at one hundred tasks. Virtualization was deferred until measured data exceeds the specified scale.

## Decision 8: Prove guided execution without migrating Activity

**Decision**: The Wails journey saves and runs the suggested task, then links to the honest Activity placeholder. A root cross-platform integration creates the platform task through the product boundary, invokes Run now, and reads the successful captured Activity record; browser tests prove the guided UI handoff.

**Rationale**: #155 owns Activity screen migration. The integration still proves the complete execution and captured-record contract now without presenting a fake Activity screen or expanding S062.

**Alternatives considered**: Migrating Activity early was rejected as scope coupling. Skipping output evidence was rejected because #192 requires a trustworthy known-good path.

## Decision 9: Create disabled draft groups atomically

**Decision**: Add an optional enabled value to group creation with the same omitted-versus-explicit-false convention already used by task creation.

**Rationale**: Creating an enabled group and immediately disabling it exposes a transient runnable subtree, emits two mutations, and can fail halfway. The additive optional field preserves older clients while making the specified draft-group path one authoritative mutation.

**Alternatives considered**: Create-then-disable was rejected as non-atomic. Making all new groups disabled was rejected as a behavior change for existing clients.
