# Implementation Plan: Windows incident replay and causal triage

**Branch**: `codex/122-windows-incident-replay` | **Date**: 2026-10-27 | **Spec**: [spec.md](spec.md)

## Summary

Reuse the delivered script observers, private incident channel, catalog definitions and real native IPC. Add small opt-in replay/observation fixtures and a fixed maintainer recipe, run one bounded campaign against retained/current runtimes and all original IPC shapes, and correct only separately demonstrated defects. Investigate existing-file launch classification and canceled MCP lifecycle operations with failing regressions before claiming either as a defect. Neither is assumed to explain the historical original.

## Technical Context

**Language/Version**: Existing Go 1.26.0, Windows PowerShell/Bash, existing Ubuntu-24.04.

**Primary Dependencies**: Existing mcphttp, API, CLI, clifixture, incident and SQLite packages. No new dependency.

**Storage**: Existing private filesystem evidence, in-memory SQLite and owned file WAL. Packaged daemon remains file-backed.

**Testing**: Go race regressions; exact unchanged catalog step execution; canonical sh scripts/verify.sh all; one source-bound full Windows/WSL evaluation after application changes.

**Target Platform**: Windows 11 primary investigation, native Linux/WSL parity; hosted macOS regression coverage.

**Project Type**: Maintainer investigation plus narrowly demonstrated runtime/IPC fixes.

**Performance Goals**: Existing production dispatch and QA resource budgets unchanged. Lifecycle cancellation must return while another operation remains held, without applying canceled mutations. Campaign condition groups run sequentially.

**Constraints**: At most ten repetitions/condition, one thirty-minute shared wall-clock campaign execution deadline; development regressions, implementation, reporting and independent verification remain separately identified. No installed-runtime/service, session, DNS, release or public profiler changes.

**Scale/Scope**: Three original symptom tracks; retained/current PowerShell and Bash contrasts; five original scenario IDs; memory/WAL component fixtures plus packaged-daemon WAL.

## Constitution Check

Pre-research: PASS on all five principles. Correctness and cancellation are explicit; regressions precede fixes, all race/access/identity/coverage contracts remain; diagnostics retain original results and use plain-language collection guidance; existing performance/resource controls remain. Source publication is explicitly owner-authorized, so no extra pre-push halt is required. Merge/tag/promotion remain excluded.

Post-design: PASS. Test-only replay is opt-in and reports unavailable prerequisites honestly. Memory mode is explicitly a component fixture; no unsupported packaged-daemon memory capability is advertised. Use the real MCP manager and catalog protocol assertions; preserve rotation old-token rejection and expected bad-port exit 2. No complexity deviation or new framework/dependency is needed.

## Research and decisions

Two read-only research agents required by /speckit-plan inspected the existing runtime and IPC seams. Findings are recorded in research.md. Both existing collectors already supply the needed evidence; missing work is condition-to-incident correlation and bounded replay. Source-inspected launch classification and lifecycle cancellation candidates require executed failing regressions before filing bugs or fixing source.

## Project Structure

- test/integration/testscripts_binding_test.go and testscripts_contrasts_test.go: opt-in observation correlation.
- test/integration/testscripts_incident_test.go and testscripts_observer_test.go: existing-but-unlaunchable interpreter regression and demonstrated correction.
- internal/clievaluate/incident_replay_test.go: opt-in unchanged catalog replay with owned memory/WAL native IPC and packaged-daemon WAL sessions.
- internal/mcphttp/manager_cancellation_test.go and manager.go: discriminating cancellation regression and demonstrated lifecycle correction if confirmed.
- docs/maintainers/windows-incident-replay.md: repeatable recipe, fixed hypothesis/outcome report and limits.
- specs/122-windows-incident-replay/: design, analysis, investigation, verification and delivery.
- specs/121-windows-incident-capture/ and specs/README.md: reconcile actual owner merge and exact-main acceptance.

## Design and execution

1. Preserve all original evidence hashes and source baseline. Prepare private exact-source binary inputs before the shared campaign starts.
2. Add opt-in script observation correlation, one existing-invalid-executable regression and controlled lifecycle cancellation tests. Record first failures without assigning original causality.
3. Create the smallest test-only replay fixture from the existing native IPC snapshot fixture. Use store.Open for memory/WAL, server.NewWithRuntimeInfo, actual mcphttp.NewWithSessions, private incident.Serve and unchanged session.execute. Seed group.create/group.child before removal and preserve enable-before-rotate ordering. Run seven catalog steps per fresh session, within eight potential snapshots. Explicitly report actual selected storage separately from logical owned runtime paths.
4. Own cleanup independently: MCP shutdown while API is live, incident cancel/join, native API shutdown/join, store close and runtime removal. Never stop the Go test host through a synthetic process handle. Packaged sessions use their existing close path.
5. Run the declared runtime and IPC matrix with ten repetitions maximum and one shared thirty-minute deadline. Deliberate prerequisite fault cases run once. Stop to analyze a real spontaneous recurrence; keep controlled negative tests separate.
6. Correct independently demonstrated defects, rerun their discriminating checks, then perform independent full source-bound evaluation if application source changed. Original reports remain immutable.
7. Finalize factual source delivery records, run all eight canonical foreground gates with before/after input manifests, format, commit, push one PR and handle at most two Codex rounds.

## Alternatives considered

A new public replay command or report framework duplicates delivered QA and is rejected. Adding a production in-memory daemon mode only for investigation expands configuration and is rejected; actual memory/WAL component fixtures plus packaged file-backed replay preserve the intended cross-product honestly. Serializing writers, retrying failures into success and increasing collection limits would change the observed contract and are excluded. Repeating unchanged full catalogs adds little causal evidence; target original shapes and perform one fresh integrated evaluation after actual application changes.

## Complexity Tracking

No constitution violation. Any confirmed MCP lifecycle gate correction uses a small context-aware operation gate with the existing manager, preserving serialized mutation and committed cleanup rather than introducing a new dependency.
