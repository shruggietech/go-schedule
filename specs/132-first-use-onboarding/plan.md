# Implementation Plan: S132 First-use Onboarding

**Branch**: `codex/132-first-use-onboarding` | **Date**: 2026-10-31 | **Spec**: [spec.md](spec.md)

**Input**: specs/132-first-use-onboarding/spec.md, owner-authorized S132 and preceding first-use recommendation.

## Summary

Deliver a self-contained local text snapshot task with native Windows and POSIX scripts, accurate asynchronous status/output and same-task recovery guidance. Replace uncreated backup.sh examples and obsolete desktop behavior/status descriptions. Show the existing own Run ID in the desktop Run details pane so the selected record can be correlated with gosched runs get. No execution/API/state change is planned unless a new defect is demonstrated.

## Technical Context

**Language/Version**: Existing Go 1.26.0 modules, React/TypeScript Wails frontend, Markdown docs and Next.js static export. No new dependency.

**Primary Dependencies**: Actual published v2.0.0 MSI/headless Linux archive; existing local named-pipe/Unix-socket API and frontend test/Playwright tooling.

**Storage**: An explicitly owned example workspace and isolated daemon data/config. Existing current store schema remains unchanged. Private receipts in .test-temp/s132, with old protected artifacts frozen before/after.

**Testing**: Failing-first component regression for own Run ID versus source run ID, actual packaged CLI/daemon task success/failure/recovery, existing desktop/browser checks, 25-document export/link policy, full foreground sh scripts/verify.sh all. New rendered guide route uses existing browser suite. Preserve first failures.

**Target Platform**: Windows package binaries in an owned hidden runtime and Linux package binaries under WSL; no service installation or active-workstation GUI. POSIX script remains /bin/sh portable for macOS, with no local native macOS acceptance claim.

**Project Type**: Focused end-to-end documentation and existing run-detail presentation correction.

**Performance Goals**: No scheduler/hot-path change; no new benchmark campaign. Existing timing/coverage/race budgets remain enforced.

**Constraints**: No visible consoles, UAC or focus stealing; UTF-8 without BOM, no em dashes/width-wrapped prose; exact GitHub formatter before publication; no migration, cluster, release mutation or unchanged historical replay.

**Scale/Scope**: One coherent first-task outcome across README, platform guides, CLI/desktop references and the new first-task guide; one existing frontend detail field; one official PR and at most two total review requests.

## Constitution Check

Pre-research: I reuse the existing product/UI/CLI patterns; II require a failing-first behavioral regression and unchanged eight canonical gates; III correct the actual labels and accepted-versus-terminal semantics; IV no hot-path change; V installed Spec-Kit sequence and explicit owner push/PR authorization. Current-only data and IPC access remain unchanged. Engineering verification is not an issue closure criterion.

Post-design: all five principles pass. The UI uses the already supplied run.id, no new query or authority. Example scripts touch only named local files; explicit cleanup does not remove daemon data. No new external listener, elevated host installation, framework or dependency.

## Project Structure

Feature artifacts live in specs/132-first-use-onboarding/. Product changes: desktop/frontend/src/operations/ActivityPage.tsx and its test. Canonical guidance: docs/first-task.md, README.md, docs/install.md, docs/INSTALL-{windows,linux,macos}.md, docs/cli.md and docs/gui-fields.md. The site derives routes automatically from canonical docs; add a first-use route assertion to site/e2e/site.spec.ts. No workflow/helper subsystem is added. Current context selector/inventory and CHANGELOG.md record the slice.

## Chronological execution and decisions

1. Specify/clarify/checklist run against installed templates before planning. Optional agent-context hook is fulfilled by updating CLAUDE.md plan reference. The installed checklist prerequisite expects plan.md before its documented pre-plan stage; retain that first failure and rerun successfully after setup-plan, without weakening scripts.
2. Phase 0 read-only research agent checks current CLI/desktop contracts and first-use guidance. Findings: Linux/macOS backup.sh does not exist; GUI guide claims obsolete live preview/two-pane behavior, outdated timing labels and unreleased cron/SMTP; the Run details pane lacks its own stable ID.
3. Select a local text snapshot workload using supplied native shell scripts. Alternatives: another hello message (does not teach useful file work/failure), new onboarding wizard (unneeded architecture), existing product plus one complete guide (selected). Preserve installed daemon identity/permissions and explicitly state paths run on the daemon host.
4. Resolve the Run ID visibility gap with the existing model field and a failing-first component test that distinguishes own run from source run and follows selection. Do not reinterpret it as an original #355/#356 causal defect.
5. Create/update the native GitHub first-use outcome under a bounded milestone after duplicate search. Record each demonstrated correction; no verification closure gate or invented date/priority.
6. Retain/freeze v2.0.0 assets and earlier protected inputs. Extract MSI files read-only using an owned private path; run only CLI/daemon binaries headlessly, never install a system service or launch GUI. Reuse archive binaries and existing runtime configuration/probes. Execute success, missing-input failure and repaired success with one stable task; inspect terminal records and output.
7. Write the canonical guide and accurate CLI/desktop instructions based on actual observations. Exercise Windows and Linux package workflows separately; macOS portability and hosted evidence remain distinct. Add screenshot/route checks only through existing headless browser tooling, redirecting generated captures to new owned paths.
8. Run final formatting/lifecycle/source/immutable-input checks and the full canonical eight-gate foreground driver. Freeze evaluated source before execution. Later changes may finalize only factual spec/status records, with explicit provenance and final audits.
9. Commit on this branch, automatically push and open official PR. Automatic code review is request one; at most one final manual request, preferably separate security review. Handle every genuine finding/CI failure with in-scope verified pushes; publish final-head CI/review receipts, then ping owner and leave unmerged. No scheduled monitor.

## Validation and recovery

Security/current-store/IPC gates remain unchanged. No ambient default endpoint/data paths in private probes: configure fresh explicit local paths and an existing permitted group, record process ownership and terminate only owned processes. Refuse an occupied or unexpected test path. Preserve original release hashes and diagnostic inputs. A functional regression gets its own issue; independent original #355/#356 remain open. Attended native installation/UI observations remain untested under the standing waiver.

## Complexity Tracking

One existing detail field plus guide corrections is the smallest supported solution. No new framework or authority. The owner explicitly authorizes publication, satisfying the normal pre-push approval; final source merge and any release remain outside scope. Pinned Windows installation-document edits are limited to the first-task link and receive a dated CHANGELOG decision.
