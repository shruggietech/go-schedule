# S141 verification

## Scope and evidence boundaries

Issue #458 is a selected-build metadata correction. Published v2.0.2 and its nine assets remain unchanged. No metadata candidate was installed or launched. Canonical verification ran its owned headless CLI/daemon fixtures. No installed runtime, Explorer UI, upgrade or native desktop outcome is inferred. No operator desktop focus or non-admin interaction occurred. Console children use hidden, noninteractive execution.

## Spec-Kit and regressions

Installed specify, clarify, checklist, plan, tasks, analyze and implement workflows completed. Six functional requirements and four success criteria map to eight completed tasks; all eleven checklist items passed. Analysis found zero critical conflicts, ambiguities, duplicate requirements or unmapped tasks. Planning research used exact cached Wails 2.17.1 source.

Regression-first unit compilation failed before the builder existed. The first actual selected Windows PE check failed: fixed file version was 2.0.3.0 but product metadata was empty. Corrected fixed product and textual file/product fields with the application language passed read-only inspection. Selected-version, development, invalid input, exact configuration restoration, child failure, exclusive ownership and malformed-config tests passed. A discriminating regression showed forwarded --ldflags and --ldflags= could bypass the selected version; both failed before the fix and passed after it. Original observations remain retained.

Private synthetic version 2.0.3 PE/MSI fixtures passed matching metadata readback and contained selected Go linker identity; a separately compiled MSI ProductVersion 2.0.4 with that 2.0.3 GUI was refused. These implementation fixtures precede the final commit and are not release candidates. The corrected canonical source additionally built the development desktop: FileVersion and ProductVersion 0.0.1, Comments Development build (0.0.1-dev), embedded Go version 0.0.1-dev. The program was not launched. Exact committed-head hosted packaging evidence is recorded on the official PR.

## Source ownership

Earlier attempt 05 passed all eight gates in 710.147 seconds, enumerating 2723 files. Staging inspection then found the Windows resource template was inside an ignored generated directory. It is now explicitly tracked through a narrow ignore exception. Moving the generated icon and manifest aside, then building with only the source template, regenerated the missing inputs and passed selected PE inspection. The final aggregate includes all 2724 enumerated inputs, including the template.

## Complete foreground verification

Canonical attempt 06 ran from 2026-11-09T22:01:39.491Z to 2026-11-09T22:13:36.108Z, 716.617 seconds, exit zero. All 2724 frozen source files and 8401 protected inputs were unchanged. Earlier lint failures, the two old inline-build integration expectations and the intentionally stopped owned attempt 04 remain separately retained.

| Gate | Result |
| --- | --- |
| format | PASS, Go and GitHub publication formatting |
| vet | PASS, root module |
| lint | PASS, zero issues |
| race | PASS, root suites including integration in 254.258 seconds |
| gui | PASS, desktop race suites, managed native Windows development build, 168 frontend tests and production bundle |
| coverage | PASS, engine 83.6%, schedule 89.1%, timezone 91.3%, store 80.9%, catchup 88.9%, logbus 91.1% |
| docs | PASS, policy/architecture mutation corpus, historical comparisons, seven site tests, 25 documents and 30 routes |
| automation | PASS, timeouts, shared-builder pin resolution, 97 assertion cases, lifecycle corpus, 30 adapters and bounded broken-stop behavior |

Shell fixture compliance passed with shellcheck -s sh. The source-only library requires an explicit shell selection; its initial no-shell invocation remains retained. After this aggregate only the CI immediate nonzero-builder-exit check, its changelog decision and factual delivery records change. An actual invalid-version helper returned 1; a separate PowerShell status-7 contrast exited 7 without reaching the following staging marker. The original combined scratch probe was inconclusive and remains retained. Fresh format, lifecycle, automation and affected workflow contracts cover the final changes. No evaluated Go, dependency, resource-template, documentation guide or packaging-check bytes change after the aggregate.

## Publication boundary

Owner authorization covers the source review branch, official PR, CI/finding resolution and at most two total Codex reviews including automatic review. Exact candidate, terminal hosted checks, reviewer findings and clean branch state are maintained on that PR. Auto-merge stays disabled; the owner performs final review and merge. No tag or public release mutation occurs in S141.
