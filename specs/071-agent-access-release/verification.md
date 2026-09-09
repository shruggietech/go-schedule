# Verification: Agent Access Controls and MCP Release Gates

**Branch**: `codex/071-agent-access-controls`

**Date**: 2026-09-09

**Issue**: [#164](https://github.com/shruggietech/go-schedule/issues/164)

## Spec-kit analysis

The specification, plan, research, data model, contract, checklists, and tasks were audited before implementation. Every functional requirement maps to a concrete task, all three user stories retain independent tests, and no constitution conflict or unresolved clarification remains. The analysis preserved one intentional boundary: S071 labels and revokes the single ephemeral S070 credential, while durable or simultaneous client grants remain assigned to issue #181. First-round Codex review identified that clipboard rollback was reported as successful without checking the disable result and that an initial frontend load failure remained on an indefinite loading panel. The service now projects confirmed disabled status, clearly reports unconfirmed revocation with an immediate CLI recovery action, and the page renders an actionable retry state; regression tests cover all three paths.

## Focused evidence

`go test ./...` passed every root-module package, including the complete Observe conformance matrix and the real package-shaped `gosched mcp serve` official-SDK subprocess smoke.

Focused daemon and API tests cover client-name compatibility and validation, successful-request evidence, saturation, rejected-request exclusion, concurrent lifecycle behavior, rotation reset, disable clearing, and non-secret additive serialization. Desktop Go tests cover workspace projection, native credential copying, confirmed and unconfirmed clipboard-failure rollback, validation, fixed guide navigation, and lifecycle results. All 80 frontend component tests and the production TypeScript/Vite build passed. Four Playwright Agent Access cases passed WCAG 2.2 AA serious/critical and horizontal-overflow checks at 80, 100, 150, and 200 percent zoom.

## Canonical verification

`scripts/verify.sh all` passed in the foreground on the completed review-branch tree:

- `format`: repository-authored Markdown and GitHub text contain no Unicode em dash or hard-wrapped prose defects.
- `vet`: passed.
- `lint`: passed with zero issues.
- `race`: passed across daemon, CLI, core, scripts, and integration packages, including the built-command subprocess and Agent Access service.
- `gui`: desktop Go suites, native Wails Windows build, all 80 frontend tests, and the production frontend bundle passed.
- `coverage`: engine 82.9%, schedule 89.1%, timezone 91.3%, store 80.7%, catchup 88.9%, and logbus 91.1%.
- `docs`: documentation policy, fixtures, links, front matter, fences, theme, and product policy passed across 17 pages.
- `automation`: workflow, CodeQL, Dependabot, release, brand, lifecycle, eight-gate, and fixture audits passed.

## Security and publication audit

Changed content passed `git diff --check`, the repository GitHub formatter, UTF-8 without BOM and mojibake inspection, secret-boundary review, spec lifecycle validation, and issue #164 traceability review. Plaintext MCP credentials exist only in one-time backend lifecycle results long enough for native clipboard transfer and have no frontend model field. Status and evidence exclude authorization values, request content, peer metadata, and failed-attempt history. No dependency, persisted schema, remote listener, mutation tool, pinned workflow, or release artifact changed.
