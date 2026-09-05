# Research: Dual-Syntax Task Input Foundation

**Date**: 2026-08-28

## D1 - Add a dedicated schedule-input orchestration package

**Decision**: Add `internal/scheduleinput`, depending on the existing cron and human schedule packages.

**Rationale**: `internal/schedule` cannot import `internal/cron` because cron already depends on schedule. The new package has a distinct responsibility: select one source grammar, compile it into the existing schedule model, retain source text, and report source identity. It keeps the engine's human parser and cron converter independently testable.

**Alternatives considered**:

- Put general input parsing in `internal/cron`. This is cycle-safe but makes a syntax-specific package own all human authoring.
- Add cron support directly to `schedule.Parse`. This creates an import cycle or forces cron code into the engine-facing grammar.
- Parse in API handlers. This duplicates classification and blocks reuse by import and future GUI work.

**Explicit deviation from S018**: S018 rejected a generic conversion wrapper because it would only forward calls. S019's package owns a new cross-surface authoring contract, so the additional boundary now has a concrete reason.

## D2 - Export and reuse S018 structural detection

**Decision**: Export `cron.DetectSyntax` and have both `cron.Convert` and `scheduleinput` use it.

**Rationale**: Automatic classification must not drift between the pure converter and task authoring. The current detector already preserves five-word human phrases while treating malformed cron-shaped values as cron.

**Alternatives considered**:

- Copy `looksLikeCron`, which would create two policies.
- Try cron then human, which violates no-fallback and obscures invalid cron.
- Require a hint everywhere, which harms normal CLI use.

## D3 - Use explicit optional syntax hints at the API boundary

**Decision**: Add optional `schedule_syntax` with normalized values `human` and `cron` to preview/create/update.

**Rationale**: Auto mode remains convenient, while integrations and import can select a parser deterministically. Invalid hints can name their own field.

**Alternatives considered**:

- No hint leaves future grammar ambiguity without an escape hatch.
- Separate cron fields duplicate schedule request shapes.
- A boolean cannot represent automatic mode clearly.

## D4 - Compile cron through the existing human parser

**Decision**: For cron input, call `cron.Explain`, parse the resulting phrase in the actual task timezone/current anchor, then replace only the inert expression with the submitted cron source.

**Rationale**: This retains one authoritative RRULE/anchor execution path and all existing DST/month/missing-date behavior. `cron.Explain` already provides specific refusal reasons for unsupported/lossy constructs.

**Alternatives considered**:

- Store and execute raw cron creates a second scheduler.
- Use the human-to-cron conversion path, which addresses the opposite direction and its deterministic conversion anchor rather than task creation context.
- Approximate unsupported input, which violates the fidelity contract.

## D5 - Derive response syntax instead of persisting it

**Decision**: Add transient response identity to `domain.Schedule`, derive it from retained expression, and leave store SQL unchanged.

**Rationale**: `schedules.expression` already persists the source needed for editing. A syntax column would be redundant state that can drift. Existing rows need no migration and expressionless rows remain valid.

**Alternatives considered**:

- Add a syntax column and forward-only migration. This adds no information that cannot be derived and creates a consistency invariant.
- Make every client detect syntax. This duplicates policy.
- Omit identity. Future GUI/API users would have to guess how to render errors and help.

## D6 - Preserve import reporting while changing retained source

**Decision**: Keep printing the converter's explanatory phrase, but preview and create from the scanner's `Line.Expr` with explicit cron syntax.

**Rationale**: Users retain familiar cron while the dry-run remains readable. Both preview evidence and persisted creation exercise the same central input. The scanner's existing single-space timing normalization is accepted as the import source contract; command tokenization is untouched.

**Alternatives considered**:

- Continue creating from `Line.Phrase`, which loses cron identity.
- Rewrite scanner slicing to preserve arbitrary timing whitespace, which adds parser risk without changing meaning.
- Remove the phrase from reports, which loses useful audit explanation.

## D7 - Explicitly contain GUI behavior

**Decision**: Existing GUI recurring requests send `schedule_syntax: human`.

**Rationale**: API auto-detection should not accidentally claim GUI cron support before live validation, prefill, error, help, and save flows are designed. This is compatibility wiring, not GUI adoption.

**Alternatives considered**:

- Let GUI requests auto-detect, which creates an undocumented bypass of its human-only validator.
- Implement complete GUI adoption here, which expands this focused slice.

## D8 - Correct live falsehoods without completing issue #52

**Decision**: Change current help/comments/docs that categorically prohibit cron on the newly supported surfaces, and record the architecture decision in the Unreleased changelog. Leave historical specs and the broad documentation posture for #52.

**Rationale**: Shipping behavior that current docs deny is incorrect. A narrow truth correction is separable from a full dialect/positioning rewrite.

**Alternatives considered**:

- Defer all prose, leaving contradictions in current documentation.
- Rewrite all user documentation, which conflates S019 with #52.

## Primary implementation sources inspected

- `internal/cron`, `internal/schedule`, and their S018 tests/history.
- API preview/create/update/read contracts and generic client methods.
- `domain.Schedule`, store v4 expression persistence, and CRUD SQL.
- CLI task and cron import paths, test doubles, and current reporting behavior.
- GUI local validation and preview/create/update request construction.
- Current README/docs plus historical Spec 008 and Spec 018 decisions.

No external dependency or unstable remote behavior is required.
