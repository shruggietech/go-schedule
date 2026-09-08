# Research: Desktop Notification Management

## Decision 1: Dedicated desktop feature service

**Decision**: Add `desktop/notifications` with its own narrow backend interface and Wails-safe models.

**Rationale**: Notification administration composes channels, task/group summaries, policies, and delivery evidence. Keeping this composition outside `taskgroup` and `operations` preserves their bounded responsibilities and makes secret exclusion reviewable at the type boundary.

**Alternatives considered**: Extend `taskgroup` and `operations` independently (fragments one workflow and duplicates refresh behavior); bind the API client directly to Wails (exposes server/domain shapes and weakens the secret boundary).

## Decision 2: Complete initial snapshot, selected-scope policy read

**Decision**: Load channels, tasks, groups, and 200 deliveries as one complete workspace. Load direct/effective policy only after a scope is selected.

**Rationale**: The first four reads are bounded aggregates needed by every section. Loading effective policy for every task would create avoidable request fan-out and inconsistent partial state. Selection-time loading keeps the UI truthful and fast.

**Alternatives considered**: One policy request per task at startup (unbounded fan-out); calculate inheritance in the desktop (duplicates authoritative server precedence).

## Decision 3: Secret-free response types with explicit replacement drafts

**Decision**: Desktop response structs omit endpoint, authorization, and payload fields entirely. Existing channel drafts carry blank write-only values plus `replaceEndpoint` and `replaceAuthorization` booleans.

**Rationale**: Type-level omission is safer than relying on UI masking. Explicit replacement distinguishes retain from clear for authorization and prevents a blank endpoint from accidentally erasing a destination.

**Alternatives considered**: Return masked secrets (still creates disclosure and state retention risk); overload blank values (cannot distinguish retain from clear).

## Decision 4: Preserve complete state and sequence asynchronous work

**Decision**: Preserve the last successful complete workspace or policy, discard stale responses by request sequence, suppress duplicate mutations, and debounce relevant notification/task/group events.

**Rationale**: This matches existing desktop operational patterns while avoiding flicker, stale overwrite, and mutation duplication during daemon events.

**Alternatives considered**: Clear state on error (destroys useful evidence); accept completion order (allows older state to replace newer state); poll (adds unnecessary local load).

## Decision 5: Explicit delivery semantics

**Decision**: Map pending with zero attempts to `queued`, pending with attempts to `retrying`, claimed to `sending`, succeeded to `successful`, and failed to `failed`; map event kind `test` to Test and `run.completed` to Task outcome.

**Rationale**: Server lifecycle values are operationally accurate but require interpretation. The mapping preserves truth while meeting the issue requirement that test and production records cannot be confused.

**Alternatives considered**: Display raw enums only (insufficient meaning); merge queued and retrying (hides recovery progress); label tests as successful runs (false).

## Decision 6: No delivery payload rendering or replay

**Decision**: Display correlation and safe failure evidence, not retained event payload bodies, and do not add replay.

**Rationale**: Payload rendering adds no requirement-critical value and can expose command/task context. Replay changes delivery semantics and belongs in separately authorized work.

**Alternatives considered**: JSON payload viewer (larger disclosure surface); replay button (new mutation contract and duplicate-delivery risk).
