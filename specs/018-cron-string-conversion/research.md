# Research: Pure Schedule String Conversion

**Date**: 2026-08-28

## D1 - Keep the service in the existing cron boundary

**Decision**: Add pure classification and conversion to `internal/cron`.

**Rationale**: That package already owns the five-field parser, cron-to-human phrasing, and human-schedule-to-cron export. It imports the schedule domain only in the cron-facing direction, so using the human parser does not create a cycle. This is the smallest reusable foundation for #50 without prematurely defining task source-syntax storage.

**Alternatives considered**:

- A new `internal/conversion` package would add a boundary whose only purpose is to call the two existing packages.
- Moving detection into `internal/schedule` would reverse the deliberate dependency direction and make cron part of normal authoring before #50.
- Implementing all logic in the CLI would make future reuse and fidelity tests depend on command plumbing.

## D2 - Use structural classification with a destination override

**Decision**: In auto mode, trimmed `@` input and five fields beginning with a cron-shaped numeric/wildcard minute field are cron; everything else is human. `--to cron` forces human input and `--to human` forces cron input.

**Rationale**: The rule is deterministic, explainable, and ensures malformed cron never falls through into a permissive human parser. Current supported existing five-field human phrases begin with words or ordinals, not a cron minute field. The explicit destination remains an escape hatch for future grammar growth.

**Alternatives considered**:

- Try cron, then human would violate the no-fallback requirement.
- Character-level scoring is hard to explain and can reclassify malformed expressions.
- Requiring `--to` every time harms the primary one-command workflow.

## D3 - Refuse implicit-anchor human schedules

**Decision**: Parse human input with a fixed UTC anchor for determinism, but refuse any output whose cron timing would be derived from that implicit anchor. Sub-daily intervals require `starting at`/`from`; calendar recurrences must encode the complete time and calendar position needed by five-field cron. For the already-supported cron-to-human minute-step, every-minute, and hourly forms, emit `starting at 00:00` so the phrase itself retains cron's phase.

**Rationale**: `every 15 minutes` begins relative to task creation, while `*/15` restarts at the top of each hour. Choosing conversion time as an anchor would silently manufacture semantics. A fixed instant makes parsing repeatable; the explicit-anchor check ensures it never changes accepted output.

**Alternatives considered**:

- Use current local time, which makes output machine- and moment-dependent.
- Assume midnight, which produces convenient but invented timing.
- Add a conversion-time or timezone flag, which still cannot recover the missing creation instant and expands beyond issue #51.
- Preserve the old phase-less phrases, which would knowingly retain a timing drift in the new round-trip surface.

## D4 - Extract schedule rendering from task export

**Decision**: Add a schedule-only renderer and keep `Export(task, schedule)` as the task-policy wrapper.

**Rationale**: Pure string conversion has no task enabled state or command, but the RRULE-to-cron logic is already correct and tested. Extraction preserves one mapping while retaining task-specific refusals for normal export.

**Alternatives considered**:

- Construct a fake enabled task, which embeds irrelevant state assumptions and obscures the true contract.
- Duplicate the recurrence switch, which invites fidelity drift.
- Change task export semantics, which is unrelated and risks regressions.

## D5 - Put structured failures on stderr

**Decision**: Structured success is one object on stdout. Structured malformed or unfaithful conversion is one object on stderr with exit 2 and no stdout. A reported-error marker prevents the root executor from appending plain text.

**Rationale**: This satisfies the issue's stable JSON requirement and its explicit no-stdout failure contract simultaneously. It also follows the constitution's result-versus-diagnostic stream rule.

**Alternatives considered**:

- JSON refusal on stdout would violate the failure contract and invite scripts to consume it as a successful conversion.
- Plain stderr in JSON mode would omit the required stable refusal shape.
- Printing both JSON and plain text would cease to be machine-readable.

## D6 - Preserve the current product boundary

**Decision**: Update only converter-specific CLI/cron documentation. Continue to state that task schedule fields reject cron until #50 lands.

**Rationale**: #51 is independently useful and local. Accepting or retaining cron in tasks requires cross-surface source identity, migrations, API behavior, and documentation handled by #50 and #52.

**Alternatives considered**:

- Combine #50, #51, and #52, which would turn a medium slice into the XL epic.
- Rewrite the product posture in advance, which would make documentation claim behavior that does not exist.

## Primary implementation sources inspected

- The repository-pinned Cobra command implementation and output-writer APIs.
- The local `internal/cron` parser, phrase, export, and refusal contracts.
- The local `internal/schedule.Parse` grammar and anchor behavior.
- Existing CLI tests and `docs/cli.md` / `docs/cron.md` contracts.

No external dependency or unstable remote behavior is required.
