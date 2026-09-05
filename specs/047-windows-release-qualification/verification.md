# Verification: Windows Release Qualification

**Date**: 2026-09-03

**Branch**: `codex/047-windows-release-qualification`

## Current result

Implemented and accepted locally. The production gate requires 47 scenarios, the collector emits all eleven desktop templates, the exact local demo passed compiled inspection, and the maintainer accepted it as a valid release target. All three accepted usability findings are tracked as Post-v1 issues. Push and pull-request publication remain intentionally unperformed.

## Spec Kit analysis

PASS after one task-path correction.

- Functional requirements and success criteria: 30/30 mapped to tasks (100%).
- Tasks: 38 dependency-ordered items across setup, four user stories, qualification, demo packaging, and the attended boundary.
- User stories: 4/4 independently testable.
- Requirements checklists: 16/16 and 30/30 complete.
- Unresolved clarification markers or template placeholders: 0.
- Constitution conflicts: 0.
- Critical findings: 0.
- High findings: 0.
- One medium finding was resolved before implementation: T006 named the nonexistent `test/integration/windows_installer_test.go`; it now names the established `test/integration/windows_release_gate_contract_test.go`.
- Unmapped tasks: 0.

At analysis time, the four Phase 8 tasks described the physical attended work required before S047 could advance to Implemented. They were completed from the maintainer's exact-demo walkthrough plus automatically collected artifact and environment facts, not from headless inference.

## Red evidence

Tests were added before validator or collector implementation.

| Command | Expected failure |
| --- | --- |
| `go test ./internal/releasegate -run TestRequiredScenarioIDsIncludeDesktopQualification -count=1` | `RequiredScenarioIDs() count = 36, want 43` |
| `go test ./test/integration -run TestAttendedCollectorUsesCanonicalScenariosAndHiddenChildren -count=1` | S047 collector missing `New-ScenarioTemplate` |

Both failures isolate the missing S047 behavior. No product/runtime test was weakened or skipped.

Second-round review corrections also began with focused red tests:

| Command | Expected failure |
| --- | --- |
| `go test ./internal/releasegate -run 'TestRequiredScenarioIDsIncludeDesktopQualification|TestValidateRejectsDeclaredImageWithNonRasterBytes' -count=1` | `RequiredScenarioIDs() count = 43, want 47`; mislabeled text produced no validation failure |

The corrections add distinct scaled observations for interaction, navigation, Tasks, and Schedule/Activity and inspect raster signatures from attachment bytes instead of trusting declared media types or extensions.

## Evidence and issue boundaries

- Formal candidate execution is impossible before review, merge, and separately authorized tag staging. The local S047 MSI will remain `local-demo` evidence.
- Issue #94 is closed in GitHub while coordinator #96 still lists it unchecked. S047 does not infer acceptance, reopen it, or rewrite its historical state.
- Issues #98, #101, #104, #105, #106, #109, #111, #112, and #113 remain open until their individual formal evidence is reviewed.
- Post-v1 diagnostics issue #102 is excluded.

## Automated verification

PASS.

- `go test ./internal/releasegate ./scripts/windows-release-gate -count=1`: PASS.
- `go test -race ./internal/releasegate ./scripts/windows-release-gate -count=1`: PASS.
- The four release workflow, promotion, collector, and inspector integration contract tests: PASS. The draft quickstart filter `WindowsRelease` matched no test names, so it was replaced with the four exact established names before this result was accepted.
- Strict PowerShell parse of `test/windows/Invoke-ReleaseCandidateAttended.ps1`: PASS.
- `pwsh -NoProfile -File build/windows/verify_wxs.ps1`: PASS.
- Foreground `scripts/verify.sh all`: PASS, all eight ordered gates: format, vet, lint (0 issues), race, GUI, coverage, docs, and automation.
- Final post-acceptance `scripts/verify.sh all`: PASS. Its first execution correctly rejected noncanonical Implemented delivery wording; the Delivery field was changed to cite the review branch and committed evidence, then the complete eight-gate aggregate passed without exclusions.
- Coverage thresholds: engine 86.4%, schedule 89.2%, timezone 91.3%, store 84.1%, catchup 88.9%, and logbus 91.1%.
- Specification lifecycle: PASS for 47 lifecycle-consistent specifications.
- Task format: PASS for 44 unique S047 task identities.
- Strict UTF-8/no-BOM and mojibake audit: PASS for all 21 changed or added paths.
- Unresolved template placeholder audit: PASS.
- `git diff --check`: PASS (the existing Git line-ending advisory for the PowerShell collector is non-failing and does not indicate whitespace damage).

After the second-round review corrections, the focused release-gate and CLI tests, collector integration contracts, PowerShell parser, WiX source check, and all eight canonical gates passed again. The task inventory now contains 44 unique completed S047 implementation tasks; GitHub thread resolution and hosted CI remain publication operations outside the Implemented lifecycle contract.

No validation threshold, existing scenario, test, or release-origin check was removed or weakened.

## Local demo identity

PASS for local, non-release handoff.

| Field | Value |
| --- | --- |
| Artifact | `dist/go-schedule_s047-demo_5151c069_windows_amd64.msi` |
| Artifact class | `local-demo` |
| Built at | `2026-09-03T21:46:54.0060000Z` |
| Source commit | `5151c0698aaa65021b32ce697f84a6c7469bf3f6` |
| Embedded version | `1.0.0-s047-demo.5151c069` |
| ProductVersion | `1.0.0` |
| ProductCode | `{98534477-793A-473C-9E3A-8C7E245C2220}` |
| UpgradeCode | `{B6F3C2E1-7A4D-4C9E-9B2A-1F6D8E5A0C34}` |
| Byte size | `24162304` |
| SHA-256 | `31feabdb3f61334b7b46a81350297cce07f8b59882bcb80336f26248c260d2f7` |
| Compiled inspection | PASS; `dist/s047-demo-5151c069-installer-inspection.md` |
| Authenticode | Not signed (local demo) |

The MSI was compiled through WiX 6.0.2 from the committed source above. `go version -m` reports the exact embedded version in all four staged executables, `gosched --version` reports the same value, and independent MSI-table inspection matches ProductVersion and ProductCode. The compiled inspector also proves the canonical product identity, Summary Subject, icon, PATH, administrative group, shortcut features, lifecycle dialogs, completion actions, wipe property and cleanup action, and stale-process handling.

This artifact is suitable only for the pre-push attended walkthrough. It is not signed, staged by the Release workflow, or eligible for promotion.

Post-build documentation policy, documentation integrity, specification lifecycle, strict UTF-8/no-BOM, mojibake, task-state, whitespace, and recorded artifact-identity checks all pass. The task audit reports 38 tasks with exactly four attended tasks intentionally open.

## Attended pre-push result

PASS for the local-demo purpose. The maintainer installed the linked S047 MSI, performed an ordinary workflow, completed uninstall with wipe, and explicitly called it a valid release target.

Automatically collected environment:

- primary display: 3440x1440, 3440x1392 working area, 120 Hz, 100 percent scaling, 96 effective DPI;
- secondary display: 2560x1440, 2560x1392 working area; and
- Windows interactive intended-user context against the installed application.

Observed passes:

- desktop and Start Menu shortcuts;
- default window size;
- tab and button hover readability;
- command-line editor size and usability;
- all offered font selections;
- Dark and Light appearance (Light was the maintainer's preference);
- scrolling, dropdown behavior, and compact application-storage rows;
- successful real task execution and completion;
- effective group suppression of an individually active task; and
- guided uninstall with wipe.

No release-blocking defect was reported and no S047 product correction was authorized or required. Three explicitly issue-only findings were created:

| Issue | Disposition |
| --- | --- |
| #118 | Post-v1: new tasks should remain inactive unless an unchecked creation-time activation option is selected. |
| #119 | Post-v1: Schedule and Activity should have wider defaults and user-adjustable persisted column widths. |
| #120 | Post-v1: Tasks should expose effective suppression by a disabled group or ancestor. |

Process deviation: the initial pre-push checklist duplicated much of the formal 47-scenario release ceremony and imposed excessive manual work. The maintainer rejected that burden. S047 now uses ordinary exploratory acceptance for the local demo, collects machine-readable facts automatically, and preserves the exhaustive matrix for the only artifact it can formally qualify: the later post-merge Release-workflow candidate. This changes no validator threshold or formal promotion requirement.

Issues #96, #98, #101, #104, #105, #106, #109, #111, #112, and #113 remain open pending their individual formal candidate evidence. The local demo does not close them by association.
