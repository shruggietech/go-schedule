# Research: GUI Dual-Syntax Scheduling

## Decision 1: Reuse the central input boundary directly

**Decision**: The editor calls `scheduleinput.Parse` with no explicit hint and
uses its returned expression, syntax, and compiled schedule.

**Rationale**: S019 already owns trimming, deterministic structural detection,
no fallback, cron fidelity refusal, and timezone-aware compilation. Reuse keeps
GUI behavior identical to the API and CLI.

**Alternatives considered**: A GUI cron detector duplicates policy. Trying the
human parser then cron changes refusal semantics. A syntax selector adds state
the operator must keep synchronized with the field.

## Decision 2: Derive intent from current text

**Decision**: Prefill `Schedule.Expression`, but classify the live text again
for each validation, preview, and form build.

**Rationale**: Exact prefill preserves operator-authored cron. Reclassification
allows syntax switching and makes mismatched or stale `SourceSyntax` harmless.

**Alternatives considered**: Caching response syntax breaks after edits.
Persisting editor syntax would add redundant state and a migration without an
execution benefit.

## Decision 3: Carry syntax in the private form model

**Decision**: Add a recurring syntax value to `taskForm`; populate it only after
successful shared parsing and send it beside the normalized expression.

**Rationale**: Preview and save then consume the same contract. One-offs and
blank legacy preservation naturally omit both fields.

**Alternatives considered**: Reclassifying inside `submitTask` loses the
editor/timezone context and creates a second path. Hard-coding cron based on
field shape recreates the detector.

## Decision 4: Preserve Start-at as a human affordance

**Decision**: Keep existing `schedule.IsSubDailyInterval` checks for visibility
and anchor composition only.

**Rationale**: Cron already carries its timing fields. Appending a human anchor
would corrupt it. Existing human interval behavior remains stable.

## Decision 5: Keep documentation local to the GUI

**Decision**: Update editor hint/help and `docs/gui-fields.md`; link to
`docs/cron.md` for the supported five-field and fidelity contract.

**Rationale**: These statements become false when S020 ships. The broader
cross-product documentation work remains tracked by issue #52.

## Existing compatibility evidence

- `splitAnchorClause` passes cron expressions through unchanged.
- Expressionless recurring details already prefill an empty field and preserve
  the stored schedule when timing remains blank.
- One-off requests already omit recurring fields.
- The recording GUI backend exposes preview/create/update requests for precise
  request-contract assertions.
- The daemon remains responsible for human summaries and upcoming-run timing.
