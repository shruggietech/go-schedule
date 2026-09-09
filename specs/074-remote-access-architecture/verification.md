# Verification: Remote Access Architecture

**Branch**: `codex/074-remote-access-architecture`

**Date**: 2026-09-09

**Issue**: [#165](https://github.com/shruggietech/go-schedule/issues/165)

## Spec Kit Analysis

The specification, plan, research, conceptual data model, remote-boundary contract, quickstart, two checklists, and 24 chronological implementation tasks were audited before implementation. All 18 functional requirements and eight measurable outcomes map to explicit artifacts and tasks. The design preserves issue-level sequencing and intentionally excludes #166 because #165 requires review before network implementation. Publication and review remain a separate runbook because they are workflow evidence rather than implementation state.

Two medium-severity completeness findings were resolved before implementation. The original research separated enrollment phrases from durable credentials but did not name a slow verifier for lower-entropy human phrases; the final design assigns salted Argon2id through `golang.org/x/crypto/argon2`, with parameters benchmarked by #169. The original boundary also did not state browser-origin behavior; the final design rejects `Origin` by default and requires a separately reviewed allowlist before any browser client can call the API directly. Reanalysis found no critical, high, or unresolved medium findings.

## Test-First Evidence

`test/scripts/remote-architecture-check_test.sh` was written before the checker and maintained page. Its first execution failed with `missing checker: /a/Code/go-schedule/scripts/remote-architecture-check.sh`, establishing red evidence. After implementation, the positive contract and all adversarial mutations passed.

The fixture suite removes or corrupts the primary HTTPS transport, separation between the local mux and remote allowlist, each of the five deployment modes, operator ownership, one threat-to-test row, the offline-mutation non-goal, and the #167-before-#168 dependency order. Every mutation is rejected with a focused diagnostic.

## Focused Verification

| Check | Result |
| --- | --- |
| `sh scripts/remote-architecture-check.sh .` | Passed; transport, trust, deployment, threats, dependencies, and issue order preserved |
| `sh test/scripts/remote-architecture-check_test.sh` | Passed; all required boundary mutations rejected |
| `sh scripts/docs-check.sh` | Passed; 18 pages, links, front matter, fences, theme, product policy, and remote architecture clean |
| `sh scripts/spec-lifecycle-check.sh .` | Passed; 74 specifications lifecycle-consistent while S074 is In Progress |
| `go run ./scripts/github-format` | Passed; no Unicode em dashes or hard-wrapped Markdown prose |

## Scope and Planning Audit

`git diff` over `internal`, `cmd`, `desktop`, `go.mod`, `go.sum`, `desktop/go.mod`, and `desktop/go.sum` is empty. S074 adds no daemon listener, configuration field, schema, endpoint, runtime dependency, generated client, public artifact, or release claim.

GitHub issues #165 through #173 were queried after focused implementation and all remain open. Issue #165 is assigned Slice `S074` and moved from Backlog through Specced to In progress. Downstream issue titles, dependency order, labels, milestones, and states are unchanged.

## Canonical Verification

The foreground `scripts/verify.sh all` aggregate passed on 2026-09-09 with Node 26.8.1 and npm 11.19.0 selected for frontend commands.

| Gate | Result |
| --- | --- |
| format | Passed; publication formatter reported no Unicode em dashes or hard-wrapped Markdown prose |
| vet | Passed |
| lint | Passed with zero issues |
| race | Passed across root, daemon, CLI, integration, and repository-tool packages |
| gui | Passed all desktop Go race tests, native Wails 2.15.0 build, 23 frontend files and 82 tests, TypeScript checking, and Vite 8.2.2 production build |
| coverage | Passed all six thresholds: engine 82.9%, schedule 89.1%, timezone 91.3%, store 80.7%, catchup 88.9%, logbus 91.1% |
| docs | Passed product policy, 18-page integrity, remote architecture, and all adversarial fixtures |
| automation | Passed action, CodeQL, Dependabot, release, brand, lifecycle, gate-manifest, and mutation contracts |

## Review Evidence

The first automatic Codex review on commit `b2280e0` raised two P2 findings. Both were accepted: the architecture guard now requires exactly one row for every threat identifier `T01` through `T12`, with a fixture that replaces `T12` with duplicate `T11`; and reverse-proxy deployment now assigns public and backend certificate issuance, renewal, hostname trust, and key permissions to the operator while go-schedule owns application TLS configuration validation. Focused architecture, fixture, documentation, lifecycle, and publication-format checks passed after the fix, and all hosted checks passed on commit `81ea4fa`.

The authorized second and final Codex review on commit `81ea4fa` raised two P2 findings. Both were accepted: the architecture guard now requires exactly one row per deployment mode and verifies that each row retains a nonempty operator-ownership cell; and the T12 guard now preserves the complete cross-origin default-deny control and its acceptance-test class rather than matching only an identifier or keyword. Adversarial fixtures remove the reverse-proxy operator contract and replace the T12 control with a permissive-origin policy. Focused and hosted final-head checks are required after these fixes.
