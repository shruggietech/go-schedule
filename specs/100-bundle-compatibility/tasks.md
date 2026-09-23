# Tasks: Target-Aware Portable Bundle Compatibility

## Tests first

- [x] T001 Add bundle compatibility matrix tests for required capabilities, supported and unsupported watcher platforms, v1/v2 documents, deterministic unique findings, and unchanged canonical digests.
- [x] T002 Add shared client transport tests for missing manifest, missing bundle and source-family capabilities, normal validation findings, no-forwarding, and selection changes between discovery and request.
- [x] T003 Extend manifest and bundle API regressions for advertised support, target-local watcher path interpretation, read-only drift, and preserved no-removal behavior.

## Implementation

- [x] T004 Derive target requirements and findings in `internal/bundle` without changing bundle wire format.
- [x] T005 Advertise the bundle-operation capability in the daemon manifest and update its contract tests.
- [x] T006 Freeze the concrete selected target and preflight every shared-client bundle operation, with validation findings and actionable errors.

## Delivery

- [x] T007 Update relevant API and CLI guidance, changelog, spec lifecycle inventory, and issue traceability without claiming publication.
- [x] T008 Run the spec-kit analyze gate, full local CI parity, and final functional acceptance audit; record delivery evidence and commit S100.
