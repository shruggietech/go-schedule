# Data Model: General Five-Field Cron Breadth

## Cron field set

An in-memory parsed field with:

- ordered, deduplicated integer values;
- fixed field bounds;
- wildcard identity;
- optional source step;
- the existing focused calendar or ordinal selector metadata where applicable.

Validation rejects empty sets, descending or out-of-range values, invalid steps, and unsupported modifier mixtures. Sunday values `0` and `7` normalize to `0`.

## Composite cron specification

Five field sets in minute, hour, day-of-month, month, and day-of-week order. The specification is executable only when:

- it contains no unsupported extension;
- day-of-month and day-of-week are not both restricted;
- focused `L`, `W`, or `#` selectors satisfy their existing narrow contracts.

## Durable schedule mapping

Newly broad standard expressions map to the existing schedule entity:

- `Kind`: recurring;
- `RRULE`: one constrained daily recurrence with explicit selected hours, minutes, and zero seconds, plus optional selected months, month dates, or weekdays;
- `Anchor`: creation instant in UTC;
- `HumanSummary`: stable readable description;
- `Expression`: normalized submitted cron source for editing only;
- `CalendarAdjustment`: empty for standard composite fields.

No column or schema version changes.

## Lifecycle

1. Parse and validate five cron fields.
2. Produce a readable description.
3. Compile bounded field sets into the durable recurrence.
4. Preview or persist through the shared task boundary.
5. Evaluate strictly after the anchor using task timezone and policy.
6. Export by recognizing and rendering recurrence values, not source text.

## Missing dates

For a day-of-month set under `skip`, absent dates produce no run. Under `last_valid` or `next_valid`, every absent target resolves independently. Candidate instants are ordered and deduplicated so two targets resolving to the same wall time produce one run.

## Compatibility invariants

- Existing schedule rows remain unchanged and readable.
- Existing simple cron expressions retain concise descriptions.
- Source text remains inert.
- Existing focused calendar adjustments remain separate from the general standard-field mapping.
