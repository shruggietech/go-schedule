# Research: Trigger-Ready Task Authoring

## R1 - Unscheduled Persistence

**Decision**: Make `tasks.schedule_id` nullable and represent no automatic time source as no schedule reference.

**Rationale**: This expresses the user's intent directly and lets scheduler, API, CLI, and GUI distinguish unscheduled work without special IDs or fabricated occurrences.

**Alternatives considered**: A shared placeholder schedule and a per-task synthetic event schedule were rejected because both claim an activation source that does not exist and complicate later trigger readiness.

## R2 - Readiness Ownership

**Decision**: Derive command readiness, automatic activation readiness, and effective eligibility from authoritative records rather than storing duplicated booleans.

**Rationale**: Readiness changes whenever a command, schedule, lifecycle, enabled flag, ancestor group, completion chain, or future trigger changes. Persisting a summary would require fragile synchronization across every mutation path.

**Alternatives considered**: Persisted `runnable` and `activation_ready` columns were rejected as denormalized state with no independent business identity.

## R3 - Activation Sources

**Decision**: Count a valid schedule, startup schedule, or incoming completion chain as an automatic source during S053; reserve extension for external triggers in #132.

**Rationale**: Completion chains already dispatch targets without changing their schedules, so a timeless target can be activation-ready today. The same query boundary can include triggers later.

**Alternatives considered**: Treat only schedules as sources was rejected because it would mislabel chain targets. Treat manual execution as automatic activation was rejected because it would conflate an operator action with enabled background behavior.

## R4 - Enable and Source-Removal Atomicity

**Decision**: Enforce enable eligibility and source-removal auto-disable within store transactions.

**Rationale**: The store is the single writer and can evaluate the task and its incoming sources at the same committed boundary as the mutation. API-level read then write logic can race.

**Alternatives considered**: GUI-only prevention and server-side prechecks were rejected because CLI or concurrent API callers could bypass or invalidate them.

## R5 - Draft Input Boundary

**Decision**: Persist omitted task name, command, and schedule values, while continuing to reject malformed supplied command lines, schedules, dates, timezones, and policies.

**Rationale**: Blank values clearly express unfinished work. Persisting malformed executable or timing syntax would require a second raw-draft schema and postpone validation into dangerous dispatch paths.

**Alternatives considered**: Persist every raw editor field was deferred because it materially expands task configuration ownership and is not necessary for the omissions explicitly requested in #129.

## R6 - Update Intent

**Decision**: Add explicit clear intent for name, command, and schedule updates while preserving omission as no change.

**Rationale**: PATCH must distinguish leaving a field untouched from intentionally clearing it. Existing clients remain compatible because their nonempty fields retain the same JSON representation.

**Alternatives considered**: Treat every empty string as clear was rejected because existing clients omit optional string fields through empty values. Replacing all fields with a full PUT was rejected as unnecessary API churn.

## R7 - Schedule Detail Shape

**Decision**: Return the existing schedule object for configured tasks and JSON `null` for an unscheduled task.

**Rationale**: Clients receive an unambiguous absence signal while configured-task responses remain structurally identical.

**Alternatives considered**: A zero-valued schedule object and an extra `has_schedule` boolean were rejected because they create ambiguous or duplicated state.

## R8 - Navigation Grouping

**Decision**: Assign each destination stable section metadata and render two separated vertical groups above the independently anchored Exit control.

**Rationale**: Metadata preserves current destination identity and order while giving #133 a direct place to insert Triggers.

**Alternatives considered**: A separator hard-coded after the third button was rejected because it encodes position rather than meaning. A disabled Triggers button was rejected because it advertises unfinished behavior.

## R9 - Group Keyboard Submission

**Decision**: Use one guarded submission function for Enter and Create, with validation before asynchronous creation and duplicate suppression after acceptance.

**Rationale**: A shared path prevents interaction drift and double creation. The single-line entry's submitted callback occurs after input-method composition is resolved by the toolkit.

**Alternatives considered**: Simulating a Create-button tap and keeping the stock dialog callback were rejected because they make dialog lifetime and duplicate suppression difficult to test or control.
