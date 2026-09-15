# Implementation Plan: Corrected v1.4.0 qualification

**Branch**: `codex/087-corrected-release-qualification` | **Date**: 2026-09-15 | **Spec**: [spec.md](spec.md)

## Summary

Refresh the unpublished candidate to reviewed main commit `951d864d9683ec3bdcb1e37535ecddd67738bf13`, preserve obsolete assets and provenance locally, stage through the existing Release workflow, and reuse S086 sessions and the existing attended gate. This is release execution and evidence work, not a new application feature.

## Technical Context

**Language/Version**: Existing Go 1.25 and Windows PowerShell 5.1/7.

**Primary Dependencies**: Existing GitHub Actions, Windows Sandbox, offline WebView2 and portable PowerShell; no new production dependency.

**Storage**: Dedicated local backups, assets, session package, and exports; no secrets in tracked evidence.

**Testing**: Existing canonical eight gates, candidate verification, unchanged native matrix and evidence finalization.

**Target Platform**: Hosted builds and disposable Windows 11.

**Project Type**: Release execution and durable evidence documentation.

**Performance Goals**: Preserve S086 diagnostic deadlines and progress; make no latency-fix claim without observation.

**Constraints**: No active-host install, visible console children, invented passes, or public promotion.

**Scale/Scope**: #226 publication qualification and #228 corrected desktop walkthrough.

## Constitution Check

Pre-research and post-design checks require spec-kit before execution, preserved safety-critical gates, exact reviewed source, scoped authorization, and full foreground local verification before the evidence PR commit. No scheduling behavior, migration, custom cryptography, or new dependency is proposed. Native unavailable evidence remains nonpassing. No exception is introduced.

## Project Structure

```text
specs/087-corrected-release-qualification/
  spec.md, plan.md, research.md, data-model.md, quickstart.md
  contracts/qualification.md, checklists/, tasks.md, verification.md
scripts/windows-qualification-session/
  existing preparation command and tests
test/windows/
  existing bootstrap, collector, native matrix and probes
.github/workflows/
  existing release.yml and promote-release.yml
```

**Structure Decision**: Reuse delivered preparation and validation; keep large binary assets and native exports outside tracked source.

## Execution and decisions

1. Reopen #226 with unfinished criteria intact and record scoped authorization.
2. Back up old annotated tag, draft metadata, and complete asset set before refreshing assets. Replace the tag with an exact old-reference lease. Reuse the existing workflow's same-name asset replacement rather than deleting the draft.
3. Stage reviewed commit through tag-push Release automation, then verify successful workflow, eight assets, manifest, candidate bytes, and independent baseline.
4. Download official portable PowerShell and Microsoft offline WebView2; prepare independent S086 fresh and upgrade exports.
5. Execute disposable qualification through supported launch/control paths and complete genuine observations, using suitable separate environments where Sandbox is insufficient.
6. Finalize only complete passing evidence. Retain actual blockers and unfinished tasks if native requirements cannot be met.
7. Verify and publish the evidence PR when its stated deliverable is complete. Handle at most two review rounds; public promotion is excluded.

The maintainer's continuation after the explicit operation list authorizes that scope. Updating only a tag without assets was rejected as contradictory provenance. Locally rebuilt bytes were rejected because qualification requires exact hosted bytes. Automatically passing unavailable observations was rejected as false evidence.

## Qualification finding and bounded repair

The fresh walkthrough found that the task Timing mode selector stretches beside Schedule helper text. A browser test reproduced a 10.296875 px height difference. Correct shared `.field` grid alignment in desktop/frontend/src/styles.css and retain a real-browser regression in desktop/frontend/e2e/tasks.spec.ts. This explicitly deviates from the initial evidence-only plan because a reproduced UI release blocker requires a source correction, not another passing evidence claim. The repair must pass the canonical verification driver and PR review before it can enter a newly staged candidate. Keep T007 and T008 incomplete and #226/#228/#231 open until their individual native criteria are satisfied. Publishing a repair/evidence PR is not a declaration that this release-execution slice is complete.
