## Summary

S099 extends portable bundles to the automation sources and notification policy references deferred by S098. Bundle v2 carries standalone external triggers, trigger sets, filesystem watcher selection rules, and notification conditions. Version 1 bundles remain readable.

## What changed

- Export produces deterministic, secret-free v2 documents with stable portable identities for source records.
- Preview accepts target-local watcher paths outside the canonical document, checks channel-name bindings, and reports item-level conflicts before apply.
- Apply creates triggers, trigger sets, and watchers disabled with fresh local trigger keys. Existing source updates preserve keys, watcher paths, and enabled state.
- CLI and desktop expose watcher path bindings and reviewed plan outcomes through the existing target-bound workflow.
- Bundle and API tests cover canonical ordering, v1 compatibility, missing bindings, disabled imports, secret exclusion, and a second preview with no drift.

## Verification

- `scripts/verify.sh all`: format, vet, lint, Go race, and desktop gates passed. The first coverage run found store at 79.9%; focused identity tests raised it to 80.1%, and coverage, docs, and automation gates passed on rerun.
- Native Windows Wails build passed. Frontend tests passed (149 tests), and production frontend build passed.
- `scripts/spec-lifecycle-check.sh .` and `go run ./scripts/github-format` passed after the final spec updates.

## Review focus

- Secret and machine-local field exclusion from exported documents.
- Binding of watcher paths and notification channels to the selected target during preview and apply.
- Per-item outcomes after partial apply and preservation of target-only state.

## Scope boundary

Existing trigger-set rename and member-count changes are reported as conflicts; the current store supports atomic retargeting but has no atomic structural edit API. No deletion or background synchronization is inferred from bundle omission.

## Complexity / Deviation

No constitution deviation. The v2 schema and preview-only path bindings are needed to preserve v1 semantics and keep local paths out of portable data.

Refs #184.
