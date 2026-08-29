# Research: General Five-Field Cron Breadth

## Decision: Compile cron directly into the durable recurrence model

**Rationale**: The current cron parser already understands standard lists, ranges, steps, and names, but `Phrase` rejects them because the human grammar cannot describe every combination. Routing cron through a phrase therefore makes language breadth an accidental execution limit. The durable schedule already separates authoritative recurrence from display summary and retained source, so cron can compile into that same recurrence without making source text executable.

**Alternatives considered**:

- Expand the natural-language grammar until every cron combination has a phrase. Rejected because it creates awkward machine-shaped prose and makes an unrelated language parser a correctness dependency.
- Add a second cron evaluator. Rejected because it would duplicate timezone, restart, catch-up, and next-run behavior.
- Split one expression into several tasks. Rejected because task identity, history, enablement, editing, and overlap policy would change.

## Decision: Use one constrained daily recurrence with explicit field sets

**Rationale**: A daily recurrence with explicit hour and minute sets plus optional month, month-day, and weekday filters represents cron conjunctions directly. Wildcard and stepped fields are expanded into bounded ordered sets, preserving field-boundary restart behavior such as `*/7`. The evaluator remains the existing recurrence engine.

**Alternatives considered**:

- Translate steps into elapsed intervals. Rejected because uneven steps drift at hour and day boundaries.
- Select a different recurrence frequency for every shape. Rejected because the resulting export recognizer and missing-date behavior become more complex without adding fidelity.
- Persist the cron AST. Rejected because normalized recurrence values already carry execution semantics and retained source already carries editing identity.

## Decision: Keep concise phrases for existing shapes and add exact descriptions for broad shapes

**Rationale**: Existing outputs such as `weekdays at 09:00` are stable user contracts. Broader shapes use a deterministic description that names selected minutes, hours, dates, months, or weekdays. The description is display metadata, consistent with `HumanSummary`, and is never reparsed for execution.

**Alternatives considered**:

- Pretend every description is an authorable human phrase. Rejected because that recreates the coupling S026 removes.
- Echo only the raw cron expression. Rejected because explain and import preview must help an operator understand it.

## Decision: Export from recurrence values, never retained source

**Rationale**: Export is a fidelity claim about what will run. It recognizes the constrained daily recurrence, reconstructs all five fields, and renders deterministic numeric wildcards, ranges, wildcard steps, or lists. Existing source is used only for editing and cannot override execution.

**Alternatives considered**:

- Return the retained expression verbatim. Rejected because legacy or externally constructed schedules may have no source, and inert source must not become authoritative.
- Always emit expanded lists. Rejected for established wildcard-step outputs such as `*/15`; deterministic compression preserves compatibility and readability.

## Decision: Preserve missing-date policy for date sets

**Rationale**: A task policy must not become silently inert because its cron source selects several dates. Non-skip evaluation resolves each absent intended date independently, applies all selected times, suppresses duplicate instants, and retains the strict anchor boundary. Export refuses a non-skip date set whenever cron would skip an occurrence the task would run.

**Alternatives considered**:

- Force composite cron tasks permanently to `skip`. Rejected because task policy is editable and would require surprising cross-field validation.
- Ignore non-default policy. Rejected as silent semantic loss.

## Decision: Keep semantic boundaries explicit

**Rationale**: Restricted day-of-month plus restricted day-of-week remains refused because cron unions those fields while one recurrence intersects filters. Modifier composites, Quartz fields, boot events, crontab environment, and shell behavior remain separate work. Issue #22 stays open.

## Performance and compatibility findings

- Every field set is bounded by cron's fixed ranges, so compilation and storage are bounded.
- The broadest time set contains 1,440 wall times per date; next-run benchmarks will cover `* * * * *`, a business-hours step, and a sparse month/date rule.
- No schema migration or dependency is required.
- Existing simple cron, human schedules, calendar adjustments, and stored rows remain readable.
