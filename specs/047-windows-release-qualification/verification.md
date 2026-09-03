# Verification: Windows Release Qualification

**Date**: 2026-09-03

**Branch**: `codex/047-windows-release-qualification`

## Current result

Implementation and unattended qualification are complete. The production gate
now requires 43 scenarios, the collector emits all seven desktop templates, and
the complete canonical repository aggregate passes. Local demo compilation and
attended pre-push qualification follow.

## Spec Kit analysis

PASS after one task-path correction.

- Functional requirements and success criteria: 30/30 mapped to tasks (100%).
- Tasks: 38 dependency-ordered items across setup, four user stories,
  qualification, demo packaging, and the attended boundary.
- User stories: 4/4 independently testable.
- Requirements checklists: 16/16 and 30/30 complete.
- Unresolved clarification markers or template placeholders: 0.
- Constitution conflicts: 0.
- Critical findings: 0.
- High findings: 0.
- One medium finding was resolved before implementation: T006 named the
  nonexistent `test/integration/windows_installer_test.go`; it now names the
  established `test/integration/windows_release_gate_contract_test.go`.
- Unmapped tasks: 0.

The deliberate open Phase 8 tasks describe the physical attended work required
before S047 can advance from In Progress to Implemented. They are not
publication bookkeeping and cannot be marked complete from headless evidence.

## Red evidence

Tests were added before validator or collector implementation.

| Command | Expected failure |
| --- | --- |
| `go test ./internal/releasegate -run TestRequiredScenarioIDsIncludeDesktopQualification -count=1` | `RequiredScenarioIDs() count = 36, want 43` |
| `go test ./test/integration -run TestAttendedCollectorUsesCanonicalScenariosAndHiddenChildren -count=1` | S047 collector missing `New-ScenarioTemplate` |

Both failures isolate the missing S047 behavior. No product/runtime test was
weakened or skipped.

## Evidence and issue boundaries

- Formal candidate execution is impossible before review, merge, and separately
  authorized tag staging. The local S047 MSI will remain `local-demo` evidence.
- Issue #94 is closed in GitHub while coordinator #96 still lists it unchecked.
  S047 does not infer acceptance, reopen it, or rewrite its historical state.
- Issues #98, #101, #104, #105, #106, #109, #111, #112, and #113 remain open
  until their individual formal evidence is reviewed.
- Post-v1 diagnostics issue #102 is excluded.

## Automated verification

PASS.

- `go test ./internal/releasegate ./scripts/windows-release-gate -count=1`:
  PASS.
- `go test -race ./internal/releasegate ./scripts/windows-release-gate
  -count=1`: PASS.
- The four release workflow, promotion, collector, and inspector integration
  contract tests: PASS. The draft quickstart filter `WindowsRelease` matched no
  test names, so it was replaced with the four exact established names before
  this result was accepted.
- Strict PowerShell parse of
  `test/windows/Invoke-ReleaseCandidateAttended.ps1`: PASS.
- `pwsh -NoProfile -File build/windows/verify_wxs.ps1`: PASS.
- Foreground `scripts/verify.sh all`: PASS, all eight ordered gates:
  format, vet, lint (0 issues), race, GUI, coverage, docs, and automation.
- Coverage thresholds: engine 86.4%, schedule 89.2%, timezone 91.3%, store
  84.1%, catchup 88.9%, and logbus 91.1%.
- Specification lifecycle: PASS for 47 lifecycle-consistent specifications.
- Task format: PASS for 38 unique S047 task identities.
- Strict UTF-8/no-BOM and mojibake audit: PASS for all 21 changed or added
  paths.
- Unresolved template placeholder audit: PASS.
- `git diff --check`: PASS (the existing Git line-ending advisory for the
  PowerShell collector is non-failing and does not indicate whitespace damage).

No validation threshold, existing scenario, test, or release-origin check was
removed or weakened.

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

The MSI was compiled through WiX 6.0.2 from the committed source above. `go
version -m` reports the exact embedded version in all four staged executables,
`gosched --version` reports the same value, and independent MSI-table inspection
matches ProductVersion and ProductCode. The compiled inspector also proves the
canonical product identity, Summary Subject, icon, PATH, administrative group,
shortcut features, lifecycle dialogs, completion actions, wipe property and
cleanup action, and stale-process handling.

This artifact is suitable only for the pre-push attended walkthrough. It is not
signed, staged by the Release workflow, or eligible for promotion.

Post-build documentation policy, documentation integrity, specification
lifecycle, strict UTF-8/no-BOM, mojibake, task-state, whitespace, and recorded
artifact-identity checks all pass. The task audit reports 38 tasks with exactly
four attended tasks intentionally open.

## Attended pre-push result

Pending maintainer execution of
`checklists/attended-demo.md` against the exact identity above.
