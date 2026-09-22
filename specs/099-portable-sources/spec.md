# Feature Specification: Portable Automation Sources

**Feature Branch**: `codex/099-portable-sources`

**Created**: 2026-09-22

**Status**: Implemented

**Delivery**: S099 review branch implementation with all eight local verification gates passed on 2026-09-22; pull request pending.

**Input**: S099, completing the remaining functional scope of [#184](https://github.com/shruggietech/go-schedule/issues/184) after S098.

## Clarifications

### Session 2026-09-22

- Bundle v2 adds external triggers, trigger sets, filesystem watcher definitions, and notification policy references. Bundle v1 remains readable and retains its original meaning.
- Export omits trigger keys, webhook endpoints, authorization, watcher paths, credentials, machine identities, and execution inputs. These values never appear in the bundle digest.
- An imported external trigger or trigger set receives fresh local keys and remains disabled until the operator configures the target. An imported watcher requires an explicit target-local path binding and starts disabled.
- Notification assignments refer to channel names only. A target must already contain exactly one channel of that name; the bundle does not create a channel or transport secrets.
- Omission never deletes or disables target records. Unsupported bindings are named conflicts in preview, not silent approximations.
- Existing trigger sets can be retargeted. A changed set name or member count is a named conflict because the current set API has no atomic edit operation for those fields.

## User Scenarios & Testing

### User Story 1 - Review source intent without secrets (Priority: P1)

An operator exports a deterministic bundle containing source names, target task relationships, watcher selection rules, and notification conditions while secret and machine-local fields are absent.

**Independent Test**: Export equivalent stores in different insertion orders, verify identical canonical JSON, and assert that known keys, paths, and endpoints are absent.

### User Story 2 - Preview on a selected target (Priority: P1)

An operator previews source changes with target-local watcher paths and sees conflicts for missing tasks or notification channels before apply.

**Independent Test**: Preview source definitions on a target with missing paths and channels, then provide the missing local bindings and confirm that only the resolved items become applicable.

### User Story 3 - Apply source definitions safely (Priority: P1)

An operator applies a reviewed target-bound plan. New triggers have fresh target keys, new watchers have bound local paths, and existing target-only records remain untouched.

**Independent Test**: Apply once, export again, and confirm matching portable intent, no transferred keys, disabled new sources, and item-level outcomes.

## Requirements

- **FR-001**: Bundle v2 MUST support portable external triggers, trigger sets, watcher definitions, and notification policy references while preserving v1 decoding.
- **FR-002**: Export MUST never include trigger keys, notification destinations or credentials, watcher paths, daemon identity, task commands, or run history.
- **FR-003**: Every collection and reference MUST have deterministic ordering and stable portable identities.
- **FR-004**: Preview MUST validate target task, source, local path, and channel bindings, report conflicts per item, and remain read-only.
- **FR-005**: Apply MUST use only the reviewed target-bound plan, preserve target-local values on updates, create sources disabled with fresh local credentials, and report unsupported trigger-set edits as conflicts.
- **FR-006**: Notification policy application MUST bind existing target channels by unique name and preserve unspecified target policies.
- **FR-007**: CLI and desktop MUST expose the same target context, conflicts, outcomes, and local watcher binding workflow.
- **FR-008**: Apply MUST never infer removal or perform background synchronization.

## Success Criteria

- Canonical export is stable for equivalent state and contains none of the known secret or machine-local fixture values.
- Preview identifies every missing watcher path, task, or channel by portable identity without mutation.
- A reviewed apply reports a terminal outcome for every planned item and a second preview reports unchanged portable intent.
