# Verification Evidence: MCP Operate Authority

**Status**: Implementation verified; pull-request publication pending

## Required evidence

| Area | Required result |
| --- | --- |
| Spec-kit | Requirements, design, tasks, contracts, and checklists complete |
| Session security | Creation, digest storage, resolution, expiry, and revocation tests pass |
| Authorization | Observe denial and Operate-only task action boundaries pass |
| Audit | Accepted, rejected, denied, and uncertain attempts retain correct MCP actor and target evidence |
| Conformance | Observe has zero tools; Operate has exactly three schema-valid tools |
| Retry safety | Identical retries mutate once; conflicting reuse is rejected |
| Transport | Stdio and localhost HTTP explicit activation, rotation, disable, and shutdown pass |
| Regression | Existing Observe resources and local API behavior remain intact |
| Repository | Format, vet, lint, race, GUI, coverage, docs, and automation gates pass |

## Final evidence

Verification completed on Windows amd64 on 2026-09-17.

| Gate | Evidence |
| --- | --- |
| Focused Go | `go test ./internal/mcpsession ./internal/mcpoperate ./internal/api/server ./internal/api/client ./internal/mcphttp ./internal/authorization ./internal/cli ./cmd/goschedd` passed |
| Desktop Go | `go test ./...` from `desktop/` passed |
| Frontend | `npm test -- --run` passed 24 files and 128 tests |
| Canonical aggregate | `scripts/verify.sh all` passed format, vet, lint, race, GUI, coverage, docs, and automation |
| Native desktop | Wails 2.15.0 production Windows amd64 build passed |
| Coverage | engine 82.9%, schedule 89.1%, timezone 91.3%, store 80.2%, catchup 88.9%, logbus 91.1% |
| MCP conformance | Observe retained zero tools; explicit Operate exposed exactly three typed tools over in-memory and authenticated HTTP SDK sessions |
| Session security | Random digest-only session secrets, expiry, malformed credentials, revocation, listener teardown, actor attribution, and Manage denial passed |
| Retry and targets | One hundred identical calls produced one mutation; conflicting request reuse and wrong-daemon requests were rejected before mutation |

The first canonical run stopped at lifecycle validation because the newly generated specification used the unsupported status `Approved`. The status was corrected to the repository's `In Progress` lifecycle state, and the complete aggregate then passed. The Spec Kit agent-context updater also hard-wrapped its managed sentence, so that generated line was normalized to the repository's required one-paragraph format before the passing run.
