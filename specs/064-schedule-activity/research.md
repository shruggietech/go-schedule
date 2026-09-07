# Research: Operational Schedule and Activity

## Decision 1: Adapt existing daemon APIs through one desktop service

**Decision**: Add `desktop/operations` with a narrow backend interface for calendar, runs, logs, alerts, and acknowledgement.

**Rationale**: The daemon APIs already own scheduling, history, retention, and alert semantics. A Wails adapter can shape stable presentation models without changing those contracts.

**Alternatives considered**: A new aggregate daemon endpoint was rejected because it expands the public API without delivery value. Calling generated Wails bindings throughout React was rejected because it scatters error handling and mapping.

## Decision 2: Keep Schedule and Activity as distinct routes

**Decision**: Schedule answers what should or did happen in a selected window; Activity answers what recently happened and why.

**Rationale**: Combining both into one page would recreate the density and hierarchy problem in issue #155. Shared refresh mechanics do not require shared navigation.

**Alternatives considered**: One operations dashboard was rejected because calendar navigation and diagnostic filtering have different primary tasks.

## Decision 3: Preserve source record types

**Decision**: Model calendar predictions, calendar run occurrences, persisted runs, daemon logs, and alerts as explicit types with stable identities.

**Rationale**: A future calculation is not an execution record, and a log line is not an alert. Explicit types prevent accidental claims and make filters predictable.

**Alternatives considered**: Flattening everything to a generic message was rejected because it hides provenance and drops diagnostics.

## Decision 4: Publish complete snapshots only

**Decision**: Schedule publishes only a complete calendar response. Activity loads runs, logs, and alerts concurrently but publishes the workspace only when every read succeeds.

**Rationale**: A timestamped complete snapshot prevents a failed section from appearing current beside fresh data.

**Alternatives considered**: Independent section loading was rejected because mixed freshness is hard to explain and test.

## Decision 5: Preserve selection in React by stable identity

**Decision**: Store selected identifiers separately from loaded records, retain them across refresh, and show a missing-record explanation if the identity leaves the current window.

**Rationale**: Live updates should not steal focus or silently redirect investigation.

**Alternatives considered**: Selecting the first refreshed row was rejected because it changes context without user intent.

## Decision 6: Keep Clear View local and acknowledgement explicit

**Decision**: Clear View records a local timestamp and acknowledges only alerts visible in the current filtered result. Individual acknowledgement remains a separate daemon mutation.

**Rationale**: This preserves the established non-destructive behavior and makes its limited mutation scope testable.

**Alternatives considered**: Deleting or purging records was rejected because no such backend behavior exists and it would violate the issue scope.

## Decision 7: Debounce relevant domain events and sequence requests

**Decision**: Run, log, alert, and task events schedule a short debounced reload. Monotonic request identifiers discard stale responses.

**Rationale**: This handles event bursts without polling or allowing slow older responses to overwrite newer state.

**Alternatives considered**: Refreshing immediately for every event was rejected because one execution can emit several related events. Periodic polling was rejected because the event stream already exists.
