# Implementation Plan: S127 Clustered-Execution Deferral

**Branch**: `codex/127-cluster-decision-deferral` | **Date**: 2026-10-31 | **Spec**: [spec.md](spec.md)

**Input**: Owner-selected explicit deferral; source base 6483fc81bbaa75ec0bddeb0fdf17fc4cc66be426.

## Summary

Deliver the owner's deferral as current documentation and native issue dispositions. Preserve all six unresolved cluster/roadmap outcomes and their dependencies. Reuse the existing intake and decision sequence; do not build a cluster, prototype, collector or release framework.

## Technical Context

**Language/Version**: Markdown/JSON planning records; existing Go and site toolchains only for required verification.
**Primary Dependencies**: Existing Spec-Kit 0.11.3 commands, GitHub Issues and repository verification scripts; no added dependency.
**Storage**: Source documentation and native issue bodies/comments. Private execution receipts stay under .test-temp/s127.
**Testing**: Requirements review, native readback, source/evidence integrity, canonical sh scripts/verify.sh all, specification lifecycle and applicable site browser checks.
**Target Platform**: Existing Windows/Linux/macOS product support; this slice changes documentation only.
**Project Type**: Maintainer decision and documentation.
**Performance Goals**: No runtime or performance change.
**Constraints**: UTF-8 without BOM, no em dash/hard-wrapped paragraphs; no visible console, scheduled monitor, issue closure, tag or release mutation.
**Scale/Scope**: Six existing issues, one current documentation page, one S127 Spec-Kit package and factual inventory/context/changelog updates.

## Constitution Check

- I Code quality: no runtime/code architecture added; reuse existing documents and issue identities.
- II Testing: run all eight unchanged canonical gates in the foreground; no new tests mirroring prose. Existing security/race/storage contracts stay intact.
- III UX: maintain one current decision page; distinguish deferred work from shipped independent daemons.
- IV Performance: runtime unchanged, so no new benchmark claim or benchmark-specific test.
- V Autopilot: installed specify, clarify, checklist, plan, tasks, analyze, implement sequence; feature branch and one authorized PR. Explicit owner publication authorization satisfies the default pre-push approval boundary.
- Functional issue closure remains independent of CI/review and remains unearned for all six issues.
- Pre-research and post-design review: PASS; no amendment or deviation required.

## Project Structure

### Documentation (this feature)

Use spec.md, plan.md, research.md, data-model.md, contracts/decision.md, quickstart.md, tasks.md, analysis.md, decision.md, delivery.md, verification.md and checklists/requirements.md plus checklists/decision.md.

### Source Code (repository root)

- docs/distributed-scheduling.md: current owner direction and reactivation path; preserve S125 historical facts.
- CHANGELOG.md: dated Unreleased decision line.
- specs/README.md: S127 lifecycle row.
- CLAUDE.md: managed current plan reference.
- .specify/feature.json: S127 feature directory.
- GitHub #186-#189, #20 and #146: native current disposition and preserved outcomes.

**Structure Decision**: Keep the existing page and native hierarchy. No new product policy, durable collector, schema, API, runtime package or duplicate planning issue.

## Decisions and sequence

1. The owner explicitly selected deferral, rather than a permanent stop or immediate architecture. Record that authority; do not manufacture workload evidence.
2. Preserve all six issues open. #186's actual demand outcome and #189's final go/no-go remain incomplete; #187/#188 stay gated; #20/#146 retain their children.
3. Owner revisit is required to reactivate; actual reports are useful inputs, not an automatic resume trigger. Preserve the semantics approval boundary.
4. Keep earlier S125/S126 records historical; update the current page's stale prepared-v1.7.0 sentence using the already published source identity.
5. Store before/after issue bodies and evidence hashes privately; format all public prose with the actual github-format command before publication and read it back.
6. No native issue deadline, priority, assignee, new label, milestone or closure is introduced.
7. The checklist command's default prerequisite requires plan.md before the prescribed checklist phase; use its supported PathsOnly identity receipt and the actual spec to perform the checklist before plan creation.
8. No unresolved research or technology unknown exists; no research-agent dispatch is needed. The optional context hook runs after this plan exists.

## Complexity Tracking

None. Required engineering verification is retained without creating a test or review as an issue closure criterion.
