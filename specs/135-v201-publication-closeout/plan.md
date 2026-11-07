# Implementation Plan: S135 v2.0.1 publication closeout

**Branch**: `codex/135-v201-publication-closeout` | **Date**: 2026-11-06 | **Spec**: [spec.md](spec.md)

## Summary

Complete #419's actual public maintenance-package outcome using accepted S134 main, then reconcile all current release-discovery and first-use availability notices in one normal review PR. Source push/official PR, CI monitoring and in-scope fixes are authorized. The owner approved exact-source tagging and draft staging on 2026-11-06. The owner separately approved the inspected asset selection; promotion run 42335391283 published the exact release. Stop before owner merge.

## Technical Context

**Language/Version**: Existing Go 1.26.0 and Node 26 toolchains; Markdown, JSON and one existing TypeScript browser assertion.

**Primary Dependencies**: Existing Release, promote-release-waived and Pages workflows; existing release-promotion, windows-release-gate, release-source-check and comparison-check commands. No new dependency.

**Storage**: Public release assets and source-owned discovery records; private immutable evidence under .test-temp/s135.

**Testing**: Canonical foreground format, vet, lint, race, GUI, coverage, docs and automation gates; production site browser suite; exact package/readback checks. Existing security and tenancy tests remain in the canonical suite.

**Target Platform**: Existing Linux/macOS daemon/CLI and desktop distributions, Windows MSI, static website.

**Project Type**: Release-delivery and documentation increment for an existing scheduler.

**Performance Goals**: No runtime or hot-path change; existing budgets and checks retained.

**Constraints**: UTF-8 without BOM, GitHub formatting, hidden noninteractive Windows children, immutable accepted tag/assets/history, standing native-testing waiver, two Codex requests including automatic review.

**Scale/Scope**: One open release issue, eight staged content assets and a ninth checksum asset; bounded current discovery/guide edits. No new dependency campaign, framework, historical migration or cluster reactivation.

## Constitution Check

Pre-research and post-design: PASS. Principle I retains existing code and error-handling guards; II retains all race, security, tenancy and coverage gates and updates the existing changed-content browser contract; III aligns user-facing release identity; IV changes no runtime path or performance budget; V runs the installed full Spec-Kit sequence and review-branch workflow, honoring explicit source publication authority and separate release decisions. Pinned Windows guide edits require a dated changelog decision. Issue closure describes delivered package/discovery behavior, not verification tasks. No constitutional deviation is requested.

## Project Structure

- specs/135-v201-publication-closeout/: spec, checklists, research, plan, data model, publication contract, quickstart, tasks, analysis and delivery/evidence records.
- README.md, site/data/published.json and current docs/test guide availability notices: actual public-release reconciliation.
- site/e2e/site.spec.ts: existing first-task availability assertion.
- CHANGELOG.md: Unreleased discovery line and dated pinned-guide decision; preserve tagged v2.0.1 section and notes/review bytes.
- .specify/feature.json, CLAUDE.md and specs/README.md: current slice selection and truthful lifecycle inventory.
- .test-temp/s135/: fresh hidden foreground runners and private source/asset/check evidence. Prior inventories and artifacts remain intact.

## Chronological execution

1. Complete specify, clarify and requirements checklists; run installed planning setup and read-only research; complete design artifacts and agent-context hook.
2. Generate story-ordered tasks, run the blocking read-only analysis gate and resolve any findings before implementation.
3. Record baseline identities and explicit staging decision. On approval, tag only accepted main and monitor the existing Release run to terminal completion.
4. Inspect the actual draft asset inventory and source/run/attempt binding. Present that concrete set for promotion approval. On approval, publish through the existing waived workflow and independently read back downloads/checksums/latest/Pages.
5. Reconcile current availability guides, fallback and existing browser assertion after publication. Preserve historical records; explicitly record the live fallback transition from its frozen baseline.
6. Run canonical and focused validation in the foreground. Finalize factual source/publication evidence and lifecycle inventory, commit, push the authorized branch and open the official PR.
7. Monitor exact-candidate checks and at most two Codex requests including automatic review; resolve findings under existing authority. Keep PR review/merge bookkeeping outside required feature tasks. Stop ready for final owner review/merge.

## Decisions

- Reuse accepted S134 source for release instead of tagging the future closeout branch. No release identity is inferred from successful source checks.
- Include every current availability notice revealed by research, preserving historically accurate baseline statements.
- Keep public package delivery, automatic comparison deployment and later merged guide deployment separately attributable.
- Preserve historical inventories; the live fallback's authorized post-publication change is a declared transition, not qualification-artifact drift.
- Use fresh S135 output paths and actual package metadata instead of another unchanged broad CLI QA replay. This does not substitute old QA for new package readback.

## Complexity Tracking

No deviation, new abstraction or new release framework.
