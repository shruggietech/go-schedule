## Summary

S100 completes target-aware compatibility for portable bundles. The daemon now advertises bundle support, and the shared client checks the selected daemon's manifest against the bundle's actual object families before forwarding export, validation, comparison, preview, or apply. This finishes the remaining functional acceptance criterion in #184 without changing the portable v1 or v2 format.

## What changed

- Derive deterministic, unique missing-capability findings from the actual document or reviewed plan. A task requires both `schedule` and `tasks`; source families require their respective daemon capabilities.
- Report unsupported or unknown target operating systems for watcher bundles while leaving target-local watcher path syntax to the daemon.
- Freeze the selected concrete target across manifest discovery and the operation. Reject preview and comparison plans whose returned daemon identity differs, and require the reviewed target identity at apply.
- Return ordinary machine-readable validation findings for incompatible targets; stop other incompatible bundle operations before forwarding their request.
- Document the CLI workflow and keep target-only drift, explicit no-removal behavior, secret exclusion, and canonical digests unchanged.

## Verification

- Wrote focused bundle, client transport, manifest, and watcher path regressions before implementing the change.
- `scripts/verify.sh all`: format, vet, lint, race, GUI, coverage, and docs passed. The initial automation gate found unchecked completed spec tasks; after correcting the lifecycle record, `scripts/verify.sh automation` and its fixture suite passed. All eight gates passed on final product code.
- `go run ./scripts/github-format`, `scripts/spec-lifecycle-check.sh .`, and `git diff --check` passed after the spec and documentation updates.

## Review focus

- Capability derivation from untrusted bundle content and reviewed apply items, with no bundle-authored downgrade path.
- One selected target for discovery and forwarding when desktop selection changes in parallel.
- Incompatible validation response shape and fail-closed behavior on missing manifest or changed target identity.
- Target-side watcher path interpretation across Windows, Linux, and macOS.

## Scope and traceability

Closes #184. Parent #174 remains open for its independent fleet outcome. No bundle schema migration, inferred deletion, background synchronization, new dependency, or public release is included.
