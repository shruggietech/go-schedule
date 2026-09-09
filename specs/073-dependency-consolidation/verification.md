# Verification: Dependency Consolidation

**Branch**: `codex/073-dependency-consolidation`

**Date**: 2026-09-09

**Issue**: [#215](https://github.com/shruggietech/go-schedule/issues/215)

**Source pull requests**: [#201](https://github.com/shruggietech/go-schedule/pull/201), [#202](https://github.com/shruggietech/go-schedule/pull/202), [#203](https://github.com/shruggietech/go-schedule/pull/203), [#204](https://github.com/shruggietech/go-schedule/pull/204), [#205](https://github.com/shruggietech/go-schedule/pull/205), [#206](https://github.com/shruggietech/go-schedule/pull/206), [#207](https://github.com/shruggietech/go-schedule/pull/207), and [#208](https://github.com/shruggietech/go-schedule/pull/208)

## Spec Kit Analysis

The specification, plan, research, data model, dependency contract, quickstart, two checklists, and 20 tasks were audited before implementation. All 14 functional requirements and seven measurable outcomes map to explicit tasks. No ambiguity, duplication, constitutional conflict, or uncovered product behavior was found.

Two coverage findings were remediated before implementation. A high-severity task coverage gap named Wails 2.15.0 in the desktop module but did not explicitly include the current build pins, CLI recovery guidance, desktop README, or packaging contract tests. T007 and T011 now cover those operational references. A medium-severity validation gap required browser and installer evidence in T013 without runnable quickstart commands. The quickstart now includes frontend end-to-end, native Wails, and focused packaging contract commands. Analysis after remediation reports complete requirement and task coverage with zero critical or high findings.

## Selected Baseline

| Source | Selected baseline | Disposition |
| --- | --- | --- |
| #201 | Wails 2.15.0 | Requested version selected; current build pins, guidance, and contract tests advance with the module. |
| #202 | modernc.org/sqlite 1.58.0 | Requested version selected in both resolved Go graphs; libc 1.75.6 and memory 1.12.1 follow module resolution. |
| #203 | fsnotify 1.10.1 | Requested version selected in both resolved Go graphs. |
| #204 | React and React DOM 19.2.8, React types 19.2.18 and React DOM types 19.2.7 | Current compatible patch releases selected within the requested update line. |
| #205 | Testing Library jest-dom 7.0.1 | Requested version selected. |
| #206 | jsdom 30.0.1 | Requested version selected. |
| #207 | Vite React plugin 6.1.1 | Requested version selected with Vite 8.2.2 because plugin 6 requires the Vite 8 peer line. |
| #208 | Node types 26.5.0 | Current compatible patch selected within the requested Node 26 line; the engine and workflow runtime advance to Node 26 to keep runtime and types aligned. |

## Clean Restoration Evidence

Root and desktop `go mod tidy` plus `go mod verify` completed successfully. A Node 26 and npm 11.6.1 `npm ci` and high-severity audit restored 138 packages with zero vulnerabilities. SHA-256 comparisons before and after restoration confirmed that all six manifests and integrity files were unchanged.

## Focused Compatibility Evidence

- Root storage, watcher, and integration race suites passed.
- All eight desktop Go packages passed with the race detector.
- All 23 frontend unit-test files and 82 tests passed on Node 26.
- TypeScript checking and the Vite 8 production bundle passed.
- All 21 Playwright accessibility, responsive, scale, and workflow tests passed on Node 26.
- A native Windows Wails 2.15.0 clean build produced `gosched-gui.exe`.
- The positive automation contract passed, and dedicated negative fixtures rejected both CI and release workflow regressions to Node 24.

## Canonical Verification

The foreground `scripts/verify.sh all` contract passed on 2026-09-09 with the Node 26 runtime selected for frontend commands:

| Gate | Result |
| --- | --- |
| format | Passed, no em dashes or hard-wrapped Markdown prose |
| vet | Passed |
| lint | Passed, zero issues |
| race | Passed across root and integration packages |
| gui | Passed desktop Go tests, native Wails 2.15.0 build, 82 frontend tests, and Vite 8 bundle |
| coverage | Passed all six core thresholds: engine 82.9%, schedule 89.1%, timezone 91.3%, store 80.7%, catchup 88.9%, logbus 91.1% |
| docs | Passed policy fixtures and 17-page documentation validation |
| automation | Passed the positive policy gate and all mutation fixtures |

## Review Evidence

Pending publication and external review.

## Publication Audit

Pull requests #201 through #208 and issue #215 were confirmed open after implementation. The committed scope contains dependency manifests, integrity files, required runtime and build pins, compatibility assertions, changelog decisions, and S073 planning evidence only. No source pull request, issue, release, tag, or public artifact was closed or published.
