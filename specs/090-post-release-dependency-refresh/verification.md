# Verification: Post-release dependency refresh

**Branch**: `codex/090-post-release-dependency-refresh`  
**Date**: 2026-09-17  
**Issue**: [#243](https://github.com/shruggietech/go-schedule/issues/243)  
**Source pull requests**: [#240](https://github.com/shruggietech/go-schedule/pull/240), [#241](https://github.com/shruggietech/go-schedule/pull/241), and [#242](https://github.com/shruggietech/go-schedule/pull/242)

## Spec-kit Analysis

The specification, plan, research, data model, dependency contract, quickstart, two checklists, and 17 tasks map every requirement to implementation and verification work. The analysis identified one material cross-artifact concern and resolved it before implementation: the Go proposals are also a required language-baseline change because x/crypto 0.57.0, x/sys 0.48.0, and x/time 0.16.0 require Go 1.26. Both modules and operational guidance therefore advance together while workflows continue reading module files.

## Selected Baseline

| Ecosystem | Selected baseline |
| --- | --- |
| Root Go | go-sdk 1.8.0; x/crypto 0.57.0; x/sys 0.48.0; x/time 0.16.0; x/mod 0.41.0; x/sync 0.23.0; x/text 0.42.0; x/tools 0.49.0 |
| Desktop Go | Go 1.26.0; x/crypto 0.57.0; x/net 0.58.0; x/sys 0.48.0; x/text 0.42.0 |
| Frontend | React and React DOM 19.3.0; React and React DOM types 19.3.0; Node types 26.5.1; Vite 8.3.0 |

## Clean Restoration Evidence

Root and desktop `go mod tidy` plus `go mod verify` completed successfully under the automatically provisioned Go 1.26 toolchain. Hash comparison before and after restoration showed zero manifest or checksum drift. Node 26.5.0 with npm 11.17.0 completed `npm install`, `npm ci`, and the high-severity audit with zero vulnerabilities and no engine, peer, force, or legacy-peer bypass.

## Focused Compatibility Evidence

- MCP HTTP, MCP Observe, remote transport, enrollment, API client, and API server race suites passed.
- All ten desktop Go packages passed with the race detector.
- All 24 frontend unit-test files and 128 tests passed.
- TypeScript checking and the Vite 8.3 production bundle passed.
- All 28 Playwright accessibility, responsive, zoom, workflow, and local-asset tests passed.
- A clean native Windows Wails 2.15.0 production build produced `desktop/build/bin/gosched-gui.exe` under Go 1.26.
- The production desktop identity, release workflow payload, and Windows installer GUI resource integration contracts passed.

## Canonical Verification

The foreground canonical aggregate passed all eight gates on 2026-09-17 after correcting one S090 lifecycle-status metadata defect exposed by the first automation pass.

| Gate | Result |
| --- | --- |
| format | Passed, including GitHub publication formatting |
| vet | Passed |
| lint | Passed with zero issues |
| race | Passed across root and integration packages |
| gui | Passed desktop race tests, native Wails build, 128 frontend tests, and Vite production build |
| coverage | Passed: engine 82.9%, schedule 89.1%, timezone 91.3%, store 80.1%, catchup 88.9%, logbus 91.1% |
| docs | Passed policy fixtures, remote architecture checks, and 20-page documentation validation |
| automation | Passed workflow, CodeQL, Dependabot, release, brand, lifecycle, and mutation-fixture contracts |

## Review Evidence

- The automatic first Codex round completed on commit `ecc761e` with no findings and an approval reaction.
- The authorized second and final Codex round found one valid operational-guidance omission: `.claude/skills/go-schedule-verify/SKILL.md` still showed a Go 1.25 linter recovery command. The command now uses Go 1.26 and the final verification contract was rerun before resolution.
- No third Codex review round was requested.
