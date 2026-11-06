# Implementation Plan: S133 First-use Maintenance Release

**Branch**: `codex/133-first-use-maintenance-release` | **Date**: 2026-11-06 | **Spec**: [spec.md](spec.md)

**Input**: Owner-authorized #419/#424 and bounded PRs #420-#423.

## Summary

Integrate the four dependency proposals, align every active Wails CLI pin with runtime 2.17.1 and make the existing automation contract detect future drift. Prepare v2.0.1 metadata, notes and content-bound comparison while public selection remains v2.0.0. Deliver an official, verified source PR and stop before owner merge; record later same-slice immutable publication.

## Technical Context

**Language/Version**: Existing Go 1.26.0 root/desktop modules, POSIX shell automation, React/TypeScript and Next.js documentation.

**Primary Dependencies**: Wails 2.17.1 and its proposal-selected indirect updates; YAML 3.0.6; Node types 26.6.6 in both npm projects; Fumadocs core/UI 16.17.9.

**Storage**: No state format or schema changes. Private evidence under .test-temp/s133; protect all retained historical inputs.

**Testing**: Failing-first positive/mutation automation regressions, read-only module/checksum verification, existing eight-gate foreground aggregate, desktop/site browser contracts and existing current-candidate QA.

**Target Platform**: Existing Linux/macOS/Windows build matrices; local native Windows plus owned WSL fixtures. No active-workstation GUI or installed service mutation.

**Project Type**: Bounded maintenance integration and release-source preparation.

**Performance Goals**: No scheduler/storage hot-path change. Existing race, coverage and benchmark policies remain.

**Constraints**: Hidden noninteractive child processes; UTF-8/no BOM, no em dashes or prose wrapping; at most two review requests including automatic review; no merge/tag/draft/promotion.

**Scale/Scope**: Two issues, four dependency proposals, one official PR.

## Constitution Check

Pre-research: I use existing module/automation/release structures; II add a failing-first contract regression and require all canonical gates; III preserve version/availability clarity and existing UI behavior; IV no hot-path change; V run installed Spec-Kit steps and honor explicit push/PR authorization plus owner merge boundary.

Post-design: all five principles pass. Deriving the existing automation expectation from desktop/go.mod prevents repeating its brittle independent hard-coded version; all actual build pins remain explicit. No dependency framework, new product authority, schema or verification bypass.

## Project Structure

S133 artifacts live here with checklists/, research.md, data-model.md, contracts/release.md, quickstart.md, tasks.md, verification.md and handoff.md. Source changes use desktop/go.mod/go.sum, root go.mod/go.sum, both npm manifests/lockfiles, Makefile, scripts/verify.sh, scripts/automation-check.sh, test/scripts/automation-check_test.sh, CI/release workflows and desktop/README.md. Release source uses README.md, CHANGELOG.md, .github/release-notes/v2.0.1.md and site/data/comparisons/reviews/v2.0.1.json. Existing published.json and historical tag records remain unchanged.

## Chronological execution and decisions

1. Specify, clarify and requirement checklist resolve source/public and review boundaries. The checklist script's pre-plan dependency is recorded and resolved through setup-plan.
2. Dispatch read-only research under the installed plan command. Reuse current metadata/comparison/QA routines and preserve all historical bytes.
3. Choose v2.0.1 under SemVer: presentation/documentation fixes and dependency patches add no public API or format. Alternatives are an unnecessary minor or postponing the merged field.
4. Prefer runtime-derived automation expectations plus targeted local/hosted pin checks over another independent literal. Retain explicit commands and use existing mutation fixtures.
5. Generate tasks and perform read-only analyze. Address any finding before implementation.
6. Implement the automation regression first; retain its pre-fix failure. Integrate proposals without merging their branches and align module graphs/lockfiles and all build guidance.
7. Prepare candidate notes/changelog/comparison with accurate current-public selection. Preserve the original v2.0.0 section and all tagged records. Explain first-use source versus candidate/package behavior.
8. Freeze source and historical inputs; run the existing foreground canonical driver, browser contracts and current-candidate QA. Finalize only factual evidence/status after acceptance; any behavioral correction requires relevant renewed acceptance.
9. Commit, automatically push and open the official PR. Count automatic review as one, use at most one additional request, resolve warranted findings and wait for exact-head terminal checks. Stop ready for owner review/merge.

## Validation and recovery

Use verified hidden launchers, redirected I/O and existing owned runtime namespaces. Record source identity, original failures and each evidence surface separately. Never overwrite retained receipts, old release content or published selection. #424 receives the source closing reference; #419 remains a reference because actual public delivery follows later approvals. No recurring monitor.

## Complexity Tracking

No constitution deviation. The small existing automation extension removes duplicate version truth and detects the exact proposal integration gap without a new subsystem.

## Implementation finding (2026-11-06)

The first aggregate failed two integration assertions that independently hard-coded CLI 2.17.0. Update test/integration/desktop_cutover_contract_test.go and windows_installer_contract_test.go to share the declared-runtime expectation, preserving their production payload/resource assertions. This is the same #424 alignment correction. Retain the failed 635.458-second aggregate and frozen manifests, then rerun the full aggregate against corrected quiescent source.
