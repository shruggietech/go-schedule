# Research: Scheduler Startup Trigger

## Decision 1: Reuse the retained event schedule row

**Decision**: Persist startup schedules as `kind=event` with `trigger_id=scheduler_startup`, `human_summary=At scheduler startup`, and the retained source expression.

**Rationale**: The domain model, `schedules.trigger_id` column, CRUD paths, and no-next-run evaluator survived removal of the old trigger feature. They are the minimal discriminated schedule shape this feature needs and require no migration.

**Alternatives considered**: Encoding startup as RRULE or one-off time invents wrong semantics. A new table duplicates the existing discriminator. Restoring triggers and dedup tables adds unrelated source-task behavior.

## Decision 2: Make engine `Start` the only firing boundary

**Decision**: After initial active-task recompute, take one startup-task snapshot and dispatch it once before catch-up and the normal timer loop. Reload only recomputes schedules.

**Rationale**: One `Start` invocation is the explicit lifecycle boundary. Keeping dispatch outside `recompute` mechanically prevents mutation and reload from firing startup.

**Alternatives considered**: Firing from `recompute` repeats on reload. Firing from daemon main through `RunNow` loses origin and duplicates eligibility. A persisted boot token conflicts with daemon-start semantics.

## Decision 3: Add a dedicated startup run origin

**Decision**: Add `startup` to run origins and preserve the legacy generic `event` value.

**Rationale**: Startup must be distinguishable from scheduled, catch-up, manual, and future generic events. SQLite stores the enum as text, so the change is additive and migration-free.

**Alternatives considered**: Reusing `event` cannot identify startup. A second run-detail column duplicates causality already owned by the trigger enum.

## Decision 4: Treat startup as syntax with no occurrence list

**Decision**: Both parsers compile to the same event schedule; explain returns `At scheduler startup`; preview and detail return no next runs; conversion handles the canonical pair before recurrence logic.

**Rationale**: The schedule evaluator already returns no next run for events. Direct compilation avoids fake RRULE values and keeps clients on the shared input boundary.

## Decision 5: Keep the existing desktop mode

**Decision**: Startup remains in the Schedule field under Recurring mode. Event schedules prefill there; recurrence policy controls stay wire-compatible but inert.

**Rationale**: A third mode adds broad UI state for one phrase and no capability. One-off remains reserved for a timestamp.

## Decision 6: No schema migration

**Decision**: Test a prior version-8 database plus persistence/reopen behavior, but do not create schema version 9.

**Rationale**: Both required schedule values already have columns and run origins are text. An empty migration adds risk without storage change.

