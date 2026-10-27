# Implementation Plan: Reliability release preparation

**Branch**: `codex/124-reliability-release-preparation` | **Date**: 2026-10-27 | **Spec**: [spec.md](spec.md)

**Input**: `specs/124-reliability-release-preparation/spec.md`, owner kickoff and issue #394.

## Summary

Prepare v1.7.0 source metadata, notes, upgrade/rollback guidance and a fresh qualified comparison. Preserve the public v1.6.0 selection and all historical evidence. Exercise the assembled source through resource-bounded Windows/WSL QA, semantic message review and the canonical aggregate, then publish one source PR and handle no more than two Codex rounds. #394 remains open for actual release delivery; #355/#356 remain unconfirmed incidents.

## Technical Context

**Language/Version**: Existing Go module/toolchain (Go 1.26.0), Markdown, JSON, Next.js/TypeScript website, PowerShell 7 and Bash maintenance interfaces.

**Primary Dependencies**: Existing Wails desktop, SQLite storage, QA coordinator, comparison/release guards, GitHub Actions; no new dependencies.

**Storage**: Versioned release/comparison/review files; ignored `.test-temp/s124` qualification records. Existing SQLite schema remains v24.

**Testing**: Source-only release guard, comparison structural/digest guard, generated Windows/WSL CLI QA, semantic evaluation and `sh scripts/verify.sh all` in the foreground.

**Target Platform**: Windows host and existing Ubuntu-24.04 WSL distribution as h8rt3rmin8r; Linux/Windows/macOS source CI and existing packaging remain unchanged.

**Project Type**: CLI/daemon, Wails desktop and static project/documentation website. This slice edits release preparation and documentation only unless a demonstrated candidate defect requires a separate tracked fix.

**Performance Goals**: Preserve product budgets. QA retains one daemon, sequential environments, worker pool 1, GOMAXPROCS 2, 120-second case allowance, original two-hour coordinator deadline, channel/artifact caps and existing Windows/WSL resource controls.

**Constraints**: No visible consoles, user-session transitions, installed daemon interference, new WSL provisioning or rewritten historical reports. No merge, tag, staging, promotion or public version change.

**Scale/Scope**: One release issue, eight products, ten comparison capabilities, six source-built QA binaries and the actual current manifest/corpus counts.

## Constitution Check

Pre-research and post-design assessment: PASS for all five principles. Existing validators and full race/coverage aggregate remain mandatory; tenancy/security assertions are unchanged. No behavioral or dependency change is planned, so no artificial implementation-mirroring tests or performance expansion is introduced. A demonstrated defect requires regression-first correction and separate issue traceability. All GitHub issue completion remains outcome based. Explicit owner authorization already covers push/PR and in-scope review fixes; it does not cover merge or release operations.

## Project Structure

```text
specs/124-reliability-release-preparation/
  spec.md, plan.md, research.md, data-model.md, quickstart.md
  contracts/source-handoff.md, checklists/, tasks.md, verification.md
.github/release-notes/v1.7.0.md
site/data/comparisons/v1.7.0.json
site/data/reviews/v1.7.0.json
docs/INSTALL-{windows,linux,macos}.md
README.md, CHANGELOG.md, docs/build-autopilot.md, docs/website.md
test/windows/README.md
scripts/release-source-check/main_test.go
specs/123-historical-windows-diagnosis/
.test-temp/s124/
```

**Structure Decision**: Extend the existing release metadata and maintainers guides, retaining the existing public selector, buildinfo fallback and package identity derivation. Resolve the actual S123 directory from the repository before editing its factual handoff. No new release orchestration or QA schema is needed.

## Phase 0: Research decisions

Three read-only research agents assessed version/package identity, all comparison cells and existing Windows/WSL QA contracts as required by the installed planning command. Consolidation is in [research.md](research.md). Select v1.7.0 because additive documented maintainer CLI selections are public interfaces under this project's declared SemVer policy. A patch increment would misclassify these additions. Keep package identities derived from the later actual tag; leave the development fallback 0.0.1-dev.

## Phase 1: Design and contracts

[Data model](data-model.md) and [source handoff contract](contracts/source-handoff.md) preserve separate source, runtime, CI/review and future asset identities. [Quickstart](quickstart.md) gives exact current source guard/QA/aggregate commands. All future tagged release identities remain pending. Post-design constitution check passes without exceptions.

## Phase 2: Implementation and qualification sequence

1. Finish specify, autonomous clarify, domain checklist, plan, tasks and blocking analyze before implementation. Update agent context through the installed hook.
2. Freeze historical evidence identities; prepare release notes, dated changelog, current version examples, recovery guides and 80-cell comparison plus settled-content digests. Reconcile S123 owner merge records.
3. Run source guards and default/public-selection checks, format all prose, then freeze assembled source and build six QA inputs. Capture Windows/WSL sequentially using the existing coordinator and opt-in incidents. Do not edit source while capture is active.
4. Read every actual message, evaluate all seven rubric criteria, condition any prior judgment reuse on exact message and complete occurrence/context/outcome equivalence, and bind reviews to this run/corpus. Finalize, validate and render through the original coordinator. Preserve first failures and separate newly demonstrated bugs.
5. Run all eight canonical gates in one watched foreground aggregate against frozen implementation; preserve terminal exit and timing. Add only factual finalization afterwards, run final formatting/encoding/lifecycle and evidence audits, commit conventionally with the Codex trailer.
6. Push and open the official PR automatically. Update #394/roadmap with selected-version and source-only status. Monitor exact-head CI and all review comments; count automatic Codex round one and at most one requested second round. Audit final head, threads, metadata, clean refs and preserved historical inputs before owner handoff.

## Complexity Tracking

No constitution violations or additional architectural layers. Publication and release operations remain independent owner decisions.

## Canonical discovery follow-up

The first aggregate demonstrated #395. T016 includes its bounded test-only correction and regression acceptance; the complete aggregate is rerun against separately frozen corrected source. The original application QA and its artifacts remain unchanged because no application or QA runner implementation changed. All original analysis requirements and five constitution principles remain satisfied.

## Source lifecycle and publication bookkeeping

Local source implementation completes after the final eight-gate acceptance. The 19 source tasks are complete; the original three post-implementation publishing actions remain mandatory in the publication ledger and official PR receipts. This follows the existing S123 distinction between implemented source and later publication evidence. Every owner-authorized push, review, CI and handoff action is preserved, with no merge/tag/staging/promotion authority added. Final factual/lifecycle edits receive separate audits and hosted source acceptance.
