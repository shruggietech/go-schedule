# Implementation Plan: Maintenance Convergence and First-use Release Handoff

**Branch**: `codex/134-maintenance-convergence-publication` | **Date**: 2026-11-06 | **Spec**: [spec.md](spec.md)

**Input**: S134 source specification, #419/#429/#430 and proposals #426-#428.

## Summary

Finish each first-task browser transition before the next platform navigation, retain failure diagnostics in the existing website CI job, and align the Pages policy/checker with the selected action's enforced 600000 ms ceiling. Consolidate the three bounded dependency proposals, selecting Next 16.4.5 as the warranted same-package patch correction to #428. Prepare coherent v2.0.1 source records and the existing post-merge publication continuation. Stop before owner merge; actual public delivery remains in #419.

## Technical Context

**Language/Version**: Go 1.26.0; Node 26; desktop TypeScript 5.6.3 and site TypeScript 7.0.2 remain selected.

**Primary Dependencies**: Wails runtime/CLI 2.17.1 unchanged; pflag 1.1.0; desktop plugin-react 6.2.1 and Vitest 5.1.0; site Fumadocs core/UI 16.17.10, Mermaid 12.3.0, Next 16.4.5 and Tailwind/PostCSS 4.3.4. Existing Playwright 1.64.1 and approved hosted-action majors remain.

**Storage**: Existing current SQLite/application/helper formats and valid-state persistence remain unchanged. Only source manifests, lockfiles and deployment/test records change.

**Testing**: Existing Go timeout-policy regression suite; production site browser case; existing frontend/race/canonical eight gates; desktop/site browser contracts; fresh receipt-bound current-candidate CLI QA with complete message review.

**Target Platform**: Existing Windows/Linux/macOS build matrix and static GitHub Pages website. Attended native installation/GUI remains untested under the standing waiver.

**Project Type**: Daemon/CLI/Wails desktop plus static documentation/comparison site.

**Performance Goals**: Preserve the existing dispatcher and verification budgets. No scheduler hot path, new benchmark target or runtime polling framework is introduced.

**Constraints**: UTF-8 without BOM, no published em dash/hard-wrapped prose, hidden noninteractive Windows child execution, unchanged source during qualification, original failed observations retained, two total Codex review requests including automatic review.

**Scale/Scope**: Two functional source issues, one remaining public-release outcome, three bounded proposals across root Go and two npm projects. Root/desktop graphs are verified together; no unused desktop pflag requirement is forced.

## Constitution Check

Pre-research and post-design checks pass against constitution 4.0.0:

- Code quality: use the existing navigation assertions, workflow, checker and publication machinery; replace the false one-hour expectation rather than preserve it.
- Testing: retain #429's original failed production observation, add failing-first cases to the existing timeout-policy test, preserve every existing assertion, and run the full foreground aggregate before commit.
- UX: preserve CLI human/JSON/error/exit behavior, all platform entry paths, current API and source/public availability distinctions.
- Performance: no production hot-path changes; retain existing performance and all resource/ownership limits.
- Autopilot: installed specify/clarify/checklist/plan/tasks/analyze/implement sequence, review branch/official PR, publication already authorized by kickoff, owner merge/release boundaries retained.
- Engineering: license audit records actual package bytes/metadata; current-state initialization, refusal, IPC security and all safety-critical race/DST/storage surfaces remain.
- #373 already permits an explicit supported-ceiling classification. The selected action's ten-minute cap is such a classification; no constitution amendment or reduction of other verification allowances is needed.

## Research and Decisions

See [research.md](research.md). All planning unknowns are resolved.

1. Coordinate the platform loop with the expected URL and existing heading before its next navigation. A new router abstraction, sleep, blind retry or weakened assertion would obscure the actual sequencing contract.
2. Upload generated site/test-results diagnostics only after the existing website browser step fails, using the approved upload action, bounded retention and a run/attempt-specific name. Keep the exercise unsuccessful. Move the existing review captures from fixed S106 filenames into testInfo.outputPath so current outputs are collected and prior captures remain untouched.
3. Select 600000 ms for the current Pages action and require that effective value in the existing checker. Explicitly reject the old clamped 3600000 ms input and preserve 360-minute outer jobs plus every other nested allowance.
4. Use existing source/artifact-bound rerun controls after a Pages timeout; preserve distinct failed and recovered receipts. Do not add a deployment framework.
5. Include #426/#427 and #428's other updates, superseding its Next 16.4.4 pin with published 16.4.5. This same-line patch has compatible published engines/peers. Static-export applicability is an inference; no application exploitability or all-seven-advisory exposure is claimed.
6. Record elkjs's changed dual-license metadata and retained EPL-2.0 option. Do not label all license bytes unchanged, change the root Apache license, or create an unrelated inventory framework.
7. Keep v2.0.1 as the maintenance target, refresh its notes/changelog/review hashes and retain the same comparison capability conclusions. Competitor facts were reviewed in S133 on this date; preserve dates and existing unknown/edition qualifiers.
8. #429/#430 close for independently integrated source behavior; #419 remains open until actual public packages/discovery are delivered. Native observations and original Windows causal incidents remain separate.

## Project Structure

### Feature documentation

```text
specs/134-maintenance-convergence-publication/
  spec.md
  plan.md
  research.md
  data-model.md
  contracts/maintenance.md
  quickstart.md
  checklists/requirements.md
  checklists/navigation-release.md
  tasks.md
  verification.md
  delivery.md
  handoff.md
```

### Existing source surfaces

```text
site/e2e/site.spec.ts
.github/workflows/ci.yml
.github/workflows/pages.yml
scripts/ci-timeout-check/main.go
scripts/ci-timeout-check/main_test.go
docs/ci-timeouts.md
go.mod and go.sum
desktop/go.mod and desktop/go.sum (coherence validation; change only if required)
desktop/frontend/package.json and package-lock.json
site/package.json and package-lock.json
CHANGELOG.md
.github/release-notes/v2.0.1.md
site/data/reviews/v2.0.1.json
```

**Structure Decision**: Reuse these existing components. Private transcripts, inventories, original failures and current QA receipts stay under .test-temp/s134; no new product subsystem is required.

## Delivery and Verification Sequence

Complete the independent navigation and Pages-policy source outcomes, consolidate dependencies, then freeze one candidate. Run the canonical foreground aggregate, existing browser contracts and fresh current-candidate QA. Review every current message, finalize/validate/render through the established QA wrapper, and audit source/protected-input identity. Update only bounded factual delivery/status records afterward. Commit and publish the official PR under existing kickoff authorization, monitor terminal applicable hosted checks and resolve findings within the two-request cap. Owner merge remains the stop.

## Complexity Tracking

No constitution violations or additional architectural components. The existing attended-native waiver is disclosed rather than treated as runtime evidence.

## Responsive assertion settlement

The first complete S134 site run completed every navigation, then sampled 38 pixels of horizontal overflow immediately after its live mobile resize. A bounded current-bundle diagnostic measured 47, 38, then zero pixels across successive animation frames and retained zero through later samples. Use Playwright assertion polling for the unchanged one-pixel limit, preserving the live-resize scenario, all route/content assertions, existing expectation timeout and zero test retries. This is condition-based layout synchronization, not a test/navigation rerun or a fixed sleep. Persistent overflow must still fail. The original failure context/trace and diagnostic samples remain retained.
