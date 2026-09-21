# S094 Verification: Agent Grant Controls

## Outcome

S094 completes issue #181 and the remaining acceptance sweep for epic #148. Agent Access now presents a secret-free MCP grant inventory, independent transport state, duration-bounded remote enrollment, monotonic authority and expiry controls, permanent revocation, and the newest 25 actor-attributed shared audit records.

## Acceptance evidence

- Grant inventory projects only MCP actors and includes client, daemon, authority, inferred transport, creation, last use when known, expiry, state, and safe fingerprint metadata.
- Observe, Operate, and Manage use distinct names, descriptions, and non-color-only card treatment. Manage explicitly identifies definition-changing authority.
- Pairing schema v19 carries a fixed optional grant deadline into the actor created by atomic exchange. Supported presets are 1 hour, 24 hours, 7 days, 30 days, and deliberate non-expiring access.
- Enrollment phrases cross only the native clipboard boundary. The React bridge receives no phrase or bearer value, and clipboard failure cancels the pending pairing.
- Existing grants can only be narrowed, assigned an earlier expiry, or revoked. Widening, extending, clearing a finite expiry, built-in actor modification, and revoked-actor reactivation are rejected.
- Shared authorization reloads current actor state on the next local or remote request. Focused tests prove capability narrowing and revocation invalidate existing request paths without restart.
- Remote listener configuration is reported independently from remote grant existence. A stored remote credential does not make the interface claim that Remote HTTPS is active.
- Recent actions use bounded shared audit fields and newest-first ordering with a 25-event limit.
- The keyboard Playwright workflow creates a bounded grant, inspects recent actions, verifies dialog focus return, checks zoom from 80 through 200 percent, and reports no serious or critical WCAG 2.2 AA violations.

## Commands and results

- `go test ./...` at the repository root: passed after formatting reconciliation.
- `go test ./...` in `desktop`: passed.
- `npm test -- --run` in `desktop/frontend`: 24 files and 133 tests passed.
- Review regressions passed for complete enrollment bundles, atomic stale-update rejection, stdio last-use audit projection, and stale recent-action response suppression.
- `npm run build` in `desktop/frontend`: TypeScript and Vite production build passed.
- `npm run test:e2e -- --grep "Agent Access"` in `desktop/frontend`: 5 Playwright tests passed.
- `go test ./internal/api/server ./internal/remotemcp`: focused next-request authority tests passed.
- `go run ./scripts/github-format .`: passed with no em dashes or hard-wrapped Markdown prose.
- `scripts/verify.sh all` through Git Bash: formatting, vet, lint, repository and desktop race suites, integration, native Windows Wails production build, frontend unit and production build, coverage thresholds, documentation, specification lifecycle, and automation policy all passed.

## Security and compatibility

No durable credential, pairing phrase, verifier, digest, protected configuration, request body, task input, environment value, or raw internal error is added to the webview or audit projection. Local IPC remains authoritative for grant administration, all optional network listeners remain disabled unless separately configured, and stdio continues to open no listening port. Existing Observe, Operate, Manage, OAuth resource binding, hostile-content, redaction, protocol compatibility, and client-attributed audit guarantees remain covered by the canonical race and integration suites.

## Traceability

- Closes #181.
- Completes the final child and acceptance matrix for #148.
