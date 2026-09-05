# Research: Activity Diagnostics Clarity

## Decision 1: Carry `log_path` in `GET /v1/logs`

**Decision**: Add an explicit `log_path` field to the existing logs response and return the full typed response from the API client.

**Rationale**: Activity already refreshes this endpoint. Returning records and their full-log location together costs no additional request and follows the client's existing typed-response conventions.

**Alternatives considered**:

- Add the path to health: rejected because Activity would need an unrelated extra request.
- Recompute the path in the GUI: rejected because overrides and platform semantics belong to daemon configuration.
- Add a new endpoint: rejected as needless surface area for one diagnostic value.

## Decision 2: Keep `log_path` present when empty

**Decision**: Do not omit the field. Empty means unavailable, including when an older daemon response decodes without the new metadata.

**Rationale**: A stable response shape is simpler to document and test. The GUI can render a truthful unavailable state without inferring a default.

## Decision 3: Preserve CLI output explicitly

**Decision**: The CLI consumes only `response.Logs` for both human and JSON output.

**Rationale**: Serializing the richer response would silently change JSON from a bare array to an object. The API may grow while the CLI's existing public shape remains stable.

## Decision 4: Use passive, semantic Activity guidance

**Decision**: Add a word-wrapped label below the existing toolbar and retain the separate Clear View explanation. Use `limited set` rather than the numeric ring size.

**Rationale**: The wording remains accurate if the internal bound changes, while the exact configured path gives operators the actionable next location. A label does not imply file interaction.

## Decision 5: Extract a startup logging helper

**Decision**: Emit `daemon startup complete` through a small helper with `endpoint`, `db`, and `log_path` attributes at the existing pre-serve site.

**Rationale**: A helper makes the single structured record directly testable. The message is a discrete event, and preserving `db` avoids unnecessary field churn.

## File and test inventory

- Server response and nil-ring behavior: `internal/api/server/logs.go`, `logs_test.go`.
- Server constructor and call sites: `internal/api/server/server.go`, daemon and four server test files.
- Client/CLI compatibility: `internal/api/client/logs.go`, `internal/cli/logs.go`.
- State propagation: `gui/viewmodel/viewmodel.go`, `viewmodel_test.go`.
- Presentation and integration: `gui/logs.go`, `gui/app_test.go`.
- Startup record: `cmd/goschedd/main.go`, new `main_test.go`.
- Contract/user copy: S004 logs API contract and `docs/cli.md`.

## Risks and controls

- Nil rings can skip metadata: assert the path in that branch.
- Richer response can leak into CLI JSON: unwrap records explicitly.
- GUI can normalize a path: use the response string directly and test unusual paths.
- Live events can erase metadata: mutate only the log slice during event folding.
- Long paths can clip: use wrapped text and a pure text function for deterministic tests.
