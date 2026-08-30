# Verification Record: S029

## Baseline (2026-08-30)

- Local `main` and `origin/main` both pointed to `931b484`; the worktree was clean before `codex/029-v09-closure` was created.
- Specifications 001-010 and 016-028 reported `Draft` despite delivery. Specs 012-015 used publication-pending compound states. Specs 003 and 006 had entirely unchecked generated task lists, 005 retained one unperformed manual walkthrough, 007 used legacy non-checkbox tasks, and 012-015 retained hosted or publication bookkeeping.
- The v0.9 milestone contained issues #33, #39, #40, and #42.
- Secret scanning and Dependabot security updates were enabled. Push protection, non-provider patterns, validity checks, AI detection, and delegated alert dismissal were disabled.
- No Dependabot version-update configuration existed.

## Test-First Evidence

- `test/scripts/spec-lifecycle-check_test.sh` first failed because the lifecycle checker did not exist, then passed after implementation.
- The new missing-Go-ecosystem automation fixture first failed because the old automation checker accepted it, then the focused automation fixture suite passed after the Dependabot contract was implemented.

## Hosted State Readback

GitHub reported the following after the narrow S029 update:

| Control | Before | After |
| --- | --- | --- |
| Dependabot security updates | enabled | enabled |
| Secret scanning | enabled | enabled |
| Push protection | disabled | enabled |
| Non-provider patterns | disabled | enabled |
| Validity checks | disabled | enabled |
| AI detection | disabled | disabled |
| Delegated alert dismissal | disabled | disabled |
| Delegated bypass | not returned in baseline | disabled |

Issue #33 remains open, retains `needs: verification`, moved from `v0.9.0 - Polish and hardening` to `Post-v1`, and changed from P1 to P3. Its S029 comment explicitly states that no VM or install observation was performed. The v0.9 milestone now contains only #39, #40, and #42. The S029 pull request will close #39 and #42; #40 must remain open until its first generated dependency proposal is reviewed.

## Analysis Gate

The cross-artifact analysis found no constitution conflict or uncovered functional requirement. It did find one workflow inconsistency: the initial S029 task list included commit/pre-push bookkeeping and placed the Implemented transition before final verification. The task list was corrected so publication stays outside implementation tasks and lifecycle completion follows all gates.

## Final Verification

- `shellcheck scripts/spec-lifecycle-check.sh scripts/automation-check.sh test/scripts/spec-lifecycle-check_test.sh test/scripts/automation-check_test.sh`: PASS with no findings.
- `sh test/scripts/spec-lifecycle-check_test.sh`: PASS.
- `sh scripts/spec-lifecycle-check.sh .`: PASS for 29 specifications.
- `sh test/scripts/automation-check_test.sh automation`: PASS.
- `sh scripts/verify.sh all`: PASS in the foreground through format, vet, lint, race, GUI, coverage, docs, and automation.
- Coverage: engine 88.1%, schedule 89.0%, timezone 91.3%, store 86.9%, catchup 88.9%, and logbus 91.1%.
- `git diff --check`: PASS.
- UTF-8 without BOM, mojibake, task-state, issue-closing-language, and final hosted-state audits: PASS.
