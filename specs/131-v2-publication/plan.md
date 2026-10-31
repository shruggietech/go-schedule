# Implementation Plan: S131 v2.0.0 Publication

**Branch**: `codex/131-v2-publication` | **Date**: 2026-10-31 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from specs/131-v2-publication/spec.md.

## Summary

Execute the already prepared #412 release continuation from accepted main 75eae1ae11ad40db38341dc7d63f32876cba6483. Reuse release.yml and promote-release-waived.yml, retain exact candidate bytes, then reconcile public fallback/current wording and execution records in this review branch. Source publication is authorized; exact tag/draft and concrete inspected-asset promotion approvals remain separate owner decisions. S131 must deliver actual publication, not a new preparation-only design slice.

## Technical Context

**Language/Version**: Go 1.26.0, two /v2 Go modules; existing Node/frontend and PowerShell packaging tools.

**Primary Dependencies**: Existing GitHub Actions, gh/Git APIs, Wails 2.17.0, SQLite 1.61.1/libc 1.79.0. No dependency changes planned.

**Storage**: Private immutable candidate downloads, source/hash inventories and receipts under .test-temp/s131; committed factual delivery summary. No application schema changes.

**Testing**: Fresh exact-main CI/CodeQL/website and source/comparison guards before staging; actual staged archive/metadata/MSI/manifest inspection; final foreground sh scripts/verify.sh all before source commit; existing browser discovery cases and negative promotion/selection tests remain mandatory in canonical gates.

**Target Platform**: Four headless Linux/macOS archives, two desktop Linux/macOS archives and one Windows amd64 MSI plus candidate manifest.

**Project Type**: Existing CLI/daemon/desktop release and static public discovery reconciliation.

**Performance Goals**: No runtime or hot-path change; preserve existing budgets. No new benchmark campaign required for factual release reconciliation.

**Constraints**: UTF-8 without BOM, no em dash or wrapped prose, actual github-format before publication; no focus-stealing console/GUI or active-workstation installation; immutable accepted bytes and all 221 retained inputs.

**Scale/Scope**: One release outcome #412/milestone 19, eight original assets/nine final public files, one source reconciliation PR, no more than two Codex requests including automatic review.

## Constitution Check

Pre-research and post-design checks pass: reuse existing product and release architecture (I); retain all security/current-state/race/coverage gates and source-bound claims (II); truthful version selection and consistent install wording (III); no performance-sensitive runtime changes (IV); full installed Spec-Kit sequence, explicit source push/PR authorization, separate release approvals and final owner merge (V). Current-only policy applies; no migrations or historical upgrade qualification. Engineering evidence is separate from functional #412 closure.

## Project Structure

`specs/131-v2-publication/` contains spec, plan, research, data-model, contracts/publication.md, quickstart, tasks, analysis, verification and delivery records. Existing workflow/helper paths remain authoritative: .github/workflows/release.yml, .github/workflows/promote-release-waived.yml, scripts/release-source-check, scripts/comparison-check, scripts/release-promotion and scripts/windows-release-gate. Current-source reconciliation may change site/data/published.json, README.md, current docs/install.md, docs/INSTALL-{windows,linux,macos}.md, docs/website.md, docs/build-autopilot.md and test/windows/README.md only after actual publication, along with factual specification inventory/context records. Do not modify immutable tagged notes/comparison or add a release framework.

## Chronological execution and decisions

1. Specify, clarify, checklist, plan, tasks and blocking analyze precede release execution. Optional agent-context hooks are satisfied by updating CLAUDE.md's active-plan reference; no other hooks are registered.
2. Read-only release research inspects existing staging/promotion authorities. Alternatives: another preparation PR (rejected: missing public outcome persists), new framework (rejected: existing exact-identity workflow already exists), direct/manual publication (rejected: bypasses accepted identity guards), existing same-source staged process (selected).
3. Recheck accepted source, all exact-main workflows/checks, local/remote tag absence and authenticated release state. Run source/comparison guards for the exact accepted SHA. Record exact approval request only after that concrete candidate is reviewable.
4. After exact staging approval, create/push v2.0.0 at the named source and observe the actual Release run to terminal completion. Download into a fresh owned path, retain every identity/digest, inspect six archives, CLI/daemon embedded build metadata, macOS bundle version and Windows MSI/manifest product/source agreement. Do not install packages on this workstation.
5. Ask separate promotion approval naming actual release ID/source/run/attempt/eight-asset inventory. After approval, dispatch the existing standing-waiver workflow, observe terminal result and independently verify final public/latest identity, unchanged originals and checksum inventory.
6. Verify actual Pages deployment and HTTPS discovery. Reconcile fallback and current source wording after publication; record the dated release-document decision in CHANGELOG.md if pinned docs/INSTALL-windows.md changes. Preserve tagged release notes/comparison identity and prior changelog sections.
7. Run full canonical foreground verification on the final source tree with pre/post source and frozen-input hashes. Finalize factual records without reattributing prior QA, then commit, automatically push and open the official PR. Automatic review is request one; optionally request one final manual review. Address every genuine comment/CI failure, accurately bind reviewed heads, then report ready and leave merge to the owner.

## Validation and recovery

A pending or failed run is never acceptance. Preserve first failures and actual attempt identities. Before retry, inspect existing tag, draft and assets; never move tag or overwrite accepted bytes. Promotion is guarded by complete allowlist, source/run/attempt, manifest/MSI agreement and exact remote asset IDs/bytes. Existing tests exercise altered identities, starter/unexpected assets, conflicting checksum, changed tag and unsuccessful staging. Any actual defect becomes a separate bug issue; it does not turn CI/attended checks into issue closure criteria. Attended native observations remain untested under the standing waiver.

## Complexity Tracking

No architecture deviation or new dependency. Owner's explicit push/PR authorization satisfies the usual pre-push approval, so no redundant pre-push halt is required. Tag/draft and concrete-asset approvals are retained as explicitly accepted release boundaries.
