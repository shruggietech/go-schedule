# Implementation Plan: S115 release candidate handoff

**Branch**: `codex/115-release-candidate-handoff` | **Date**: 2026-10-24 | **Spec**: [spec.md](spec.md)

## Summary

Carry 18e6fd8 after blocking analysis, extract the production release-source preflight into a read-only Go command, strengthen its dated-heading contract, and exercise current native writer and file-backed IPC behavior. Keep #354 open for publication and #355/#356 open while their original causes remain unknown.

## Technical Context

**Language/Version**: Go 1.26.0, PowerShell 7, existing POSIX shell and Node tooling.

**Dependencies**: Standard library for the new command; unchanged root/desktop and locked frontend/site dependencies.

**Storage**: Existing SQLite, owned temporary file-backed WAL fixtures and versioned QA evidence. No production schema change.

**Testing**: Eight canonical gates, executable release metadata regressions, unchanged six-writer native paths, concurrent memory/file-backed IPC and a receipt-bound full owned Windows/Linux evaluator capture.

**Platforms and budgets**: Windows 11 and native Linux through existing WSL2; hosted Linux/macOS/Windows acceptance. Preserve 15-second case, existing resource limits, six writers and exact counts. Use the previously authorized two-hour total evaluator budget with original per-case/resource fields unchanged.

**Constraints**: No installed service, account/session, dependency/runtime replacement, public endpoint, DNS, tag, staging or promotion changes. All process launches hidden and noninteractive.

## Decisions

See [research.md](research.md). All research questions are resolved. A reusable source preflight provides useful pre-tag preparation without requiring package publication. No current source defect is demonstrated for the historical native failures; extend the missing concurrent file-backed workload rather than speculate about transports or runtime replacements.

## Structure and contracts

- `scripts/release-source-check/main.go` and `main_test.go`: strict local source validation and regression fixtures.
- `.github/workflows/release.yml`: shared command after setup-go, before comparison and staging dependencies.
- `scripts/automation-check.sh` and `test/scripts/automation-check_test.sh`: enforce invocation and failure propagation instead of obsolete inline strings.
- `internal/clievaluate/ipc_diagnostics_test.go`: bounded simultaneous snapshot/write coverage with memory and file-backed cases.
- `specs/115-release-candidate-handoff/`: design, source/candidate observations and handoff. Private full output stays ignored under `.test-temp/s115`.
- Existing release preparation docs receive the carried reconciliation and concise preflight/handoff guidance; refresh comparison digests only for governed content changes.

## Execution order

Complete tasks and read-only analyze before implementation. Carry 18e6fd8, demonstrate the dated-heading defect with a fixture, implement the shared validator and policy changes, then extend IPC coverage. Run focused regressions and native investigations once per declared experiment. Complete all canonical gates before source checkpoint; freeze that checkpoint for the six-binary full owned capture and actual semantic review. Record any later metadata changes as distinct from the captured subject. Push one PR, monitor exact-head CI and up to two review rounds outside implementation-task checkboxes.

## Constitution Check

Pre-research and post-design gates pass. Standard library suffices; no production schema or hot-path change. Every changed behavior has an executable regression; original assertion strength, deadlines and issue dispositions remain intact. Foreground canonical parity, source-bound observations and immutable publication contracts are required. Explicit owner push authorization satisfies the publication halt; merge/tag/staging/promotion remain excluded.

## Context extension disposition

The optional agent-context hook was inspected. Its PowerShell fallback launches unwrapped Python probes, so it is not executed on this desktop. The managed CLAUDE.md plan pointer was updated directly to the same S115 path; no mandatory hook is registered.

## Verification discovery adjustment (2026-10-24)

Initial canonical run failed the unchanged ten-minute integration-package timeout while TestV13PackageDefaultsRemainOptIn was 52 seconds into its daemon build. Preserve that failure as its own observation, not a demonstrated CLI hang. The existing packaged fixtures repeatedly compile identical targets. Add process-local once-per-target build preparation using the existing hidden owned process runner, a bounded preparation deadline and immutable in-memory candidate bytes copied to each original package path. Keep every original runtime/assertion and global/QA deadline unchanged. This is a source-preparation efficiency correction under FR-005, not a correction of #355/#356.

## Timeout fixture chronology discovery (2026-10-24)

Second canonical run completed integration in 473.051 seconds but failed the existing injected 150 ms QA-wrapper capture fixture because expected pre-timeout child output was absent. The fixture currently assumes startup and output precede its tiny timeout. Synchronize a private test copy at the existing wait boundary with an explicit child-ready marker emitted after both outputs, then enforce the same 150 ms policy and all original output/input assertions. Readiness preparation has its own bounded failure path. Production wrapper and real QA timing policy remain unchanged. Explicitly disable Git/GCM prompts for the new owned package builds through inherited environment overrides.

## Frozen message-review discovery adjustment (2026-10-24)

After finalize/validate/render preserved native outcome 2 on fc39261, correct two independently tracked communication defects under FR-005. Override only the default PowerShell completion Long guidance in the existing CLI command-contract initialization, telling readers to add the invocation itself to their profile. Preserve generated script bytes, current-session pipeline, flags, other shells and all profile state. Use permission-neutral whole-session MCP title/description/instructions so callers may add authorized tools without inheriting an Observe-only claim. Keep resource-specific Observe descriptions, content-trust warning, tool discovery and authority controls unchanged. Verify actual CLI help/generator output and SDK Observe/Operate/Manage discovery; retain the full failed report under its original source and distinguish these later corrections.
