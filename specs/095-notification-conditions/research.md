# Research: Actionable Notification Conditions

## Durable evaluation location

**Decision**: Evaluate task conditions inside `RecordRunAndCreateDeliveries` after the source run insert and before commit.

**Rationale**: The run, condition state, and zero or more resulting deliveries commit together. A crash cannot record a run without advancing its streak or create a delivery without its source run.

**Alternatives rejected**: Dispatcher-time evaluation risks duplicate state transitions after restart. In-memory evaluation loses streak and suppression state. A separate event broker adds failure modes and dependencies.

## Condition representation

**Decision**: Extend notification assignments with explicit scalar fields and use one deterministic primary problem per channel and run.

**Rationale**: The fixed vocabulary is small, validates clearly, migrates safely, and avoids a general event-query language. One primary problem prevents several messages for one run.

**Alternatives rejected**: JSON rule documents weaken database validation and complicate compatibility. A generic expression language exceeds issue scope. Independent simultaneous condition rows create duplicate noise.

## State reset

**Decision**: Hash the effective assignment identity and settings into a policy fingerprint stored with condition state.

**Rationale**: Direct and inherited policy changes reset streak and active state lazily on the next run without rewriting every descendant task when a group policy changes.

**Alternatives rejected**: Eager descendant updates are complex and expensive. Keeping state across policy changes can immediately fire using history collected under different rules.

## Daemon health

**Decision**: Add opt-in healthy-presence heartbeats on notification channels, scheduled through durable state and the existing dispatcher.

**Rationale**: A stopped or disconnected daemon cannot emit an outage notification. Receivers can reliably detect absence, while resumed heartbeats provide recovery evidence. This is an honest monitoring contract.

**Alternatives rejected**: A second resident watchdog expands packaging and lifecycle scope. Shutdown notifications are unreliable. Treating notification transport failures as daemon health is circular.

## Time configuration

**Decision**: Persist whole seconds, expose integer seconds in JSON, accept Go duration syntax in the CLI, and present compact duration controls in the desktop.

**Rationale**: Seconds are stable across languages and SQLite, while Go duration strings remain consistent with existing CLI practice. Sub-second alert timing has no operational value here.

## Webhook compatibility

**Decision**: Keep `go-schedule.webhook.v1`, add optional condition information to run events, and add the distinct `daemon.health` event kind.

**Rationale**: Optional JSON fields are backward compatible for ordinary decoders. Event headers already carry event kind. Health has no fabricated task or run.

## Security and privacy

**Decision**: Persist and serialize only enumerated condition kinds, bounded generated explanations, thresholds, streaks, and safe timing metadata.

**Rationale**: Commands, arguments, output, raw process errors, endpoint secrets, and authorization remain excluded from all new surfaces.
