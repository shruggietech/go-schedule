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

Pending committed source, build, and compiled inspection.

## Attended pre-push result

Pending exact local-demo handoff and maintainer walkthrough.
