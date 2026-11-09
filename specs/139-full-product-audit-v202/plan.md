# Implementation Plan: Full-product audit and v2.0.2 maintenance delivery

**Branch**: `codex/139-full-product-audit-v202` | **Date**: 2026-11-08 | **Spec**: [spec.md](spec.md)

## Summary

Deliver owner-approved S139 under [#444](https://github.com/shruggietech/go-schedule/issues/444), starting at accepted main 8b10419b916e04891609b3cef9727da16f486a23. Reuse the complete source-bound Windows/WSL CLI evaluator and canonical gates, supplement every entrypoint outside its main Cobra tree, exercise bounded real cross-platform reachability and isolated UI feasibility, file all demonstrated defects, correct compatible functional outcomes and prepare the v2.0.2 source handoff. The operator desktop must never acquire focus from any audit child.

## Technical Context

**Language/Version**: Existing Go modules/toolchain, PowerShell 7 hidden coordinator, existing TypeScript/Playwright frontend.

**Primary Dependencies**: Existing application dependencies and audit tools only. No new QA framework, global host configuration, user/session provisioning or dependency is planned.

**Storage**: Existing UUID-owned runtime/evidence fixtures and strict version-1 CLI report contracts. Private scratch and raw environment inputs remain under ignored .test-temp/s139; only sanitized reports/selected images are published.

**Testing**: Failing-first focused product regressions; complete actual Windows/WSL catalogue with every emitted message semantically assessed; supplemental entrypoint matrix; real cross-platform reachability; headless frontend screenshots; all eight canonical foreground gates; exact-head hosted CI and bounded reviews.

**Target Platform**: Windows 11 host and existing nonroot WSL2 Linux environment. Existing native CI supplies separate macOS/build evidence, not a local native walkthrough claim.

**Project Type**: Whole-product source/behavior assessment and maintenance correction/release preparation.

**Performance Goals**: Retain existing supported budgets, one catalogue worker and bounded owned processes. No timeout-policy change or synthetic performance claim.

**Constraints**: Zero operator-focus acquisition, zero non-admin operator requests, no WSLg/shared-session native GUI launch, no native Linux keyring probe that could prompt, no installed-service mutation through hardcoded product service identity, immutable build/capture inputs and independent cleanup reporting.

**Scale/Scope**: The checked-in main catalogue currently has 152 command nodes, 125 runnable nodes, two aliases and 1590 effective flag records; rebuild inventory instead of hardcoding these historical counts. Supplement daemon, QA coordinator, cleanup/tray/indicator/desktop and documented script/utility entrypoints. No clustered execution or historical migration work.

## Constitution Check

Pre-research and post-design: PASS. I: existing exact process ownership, bounded streams and cleanup; no global UI input. II: new behavior corrections receive discriminating regressions, existing race/coverage and canonical gates remain. III: every surface and execution boundary remains explicit; missing observations cannot become passes. IV: reuse bounded fixtures instead of a new framework. V: approved Spec-Kit sequence, authorized official PR and two-review cap; owner merge/tag/promotion remain separate. No waiver.

## Research and design decisions

- Use the existing six-binary build receipt, source-content digest and fixed CLI report/corpus/reviews schema. Capture must use frozen source and binaries. A source correction needs new build identities and fresh affected observations; original failures remain retained.
- The main catalogue covers only gosched. Maintain a supplemental per-entrypoint disposition matrix for goschedd, gosched-qa, cleanup, tray/indicator, desktop and documented wrappers/utilities. Help is not functional coverage. Adapter service/GUI/native-secret results retain their actual boundary.
- Source research identified hypotheses in daemon extra-operand handling and JSON acknowledgements currently excluded from catalogue JSON checks. Establish failing observations before filing/correcting them. Reuse existing functional command handlers and report schema; avoid broad refactoring.
- The Windows GUI singleton uses fixed Local session objects and can signal the owner's existing GUI even from another desktop in the same logon session. Never launch it in that session. Existing WSL has GTK/WebKit but no isolated display server. Native UI is blocked unless a genuinely isolated existing environment is established first; no default WSLg launch. Headless Chromium screenshots use the actual frontend with clearly labeled injected backend fixtures, not native application evidence.
- Existing WSL reports networking mode none, no default route and no nonloopback address. Perform bounded actual endpoint attempts before reporting cross-OS reachability blocked. Do not reconfigure networking or substitute same-host transport as cross-platform proof. Where reachable, use actual supported HTTPS, exact certificate/daemon identity and owned synthetic tasks.
- Native Linux keyring operations may invoke Secret Service unlock prompts. Do not touch the owner's session keyring. Standard catalogue memory-secret coverage remains adapter evidence; native credential coverage needs proven isolated unlocked infrastructure or a literal blocked outcome.
- Service identity is hardcoded to goschedd; actual service mutation could affect the owner's installation. Catalogue handler-adapter observations remain explicit. Supplemental safe status/usage/refusal paths and installer-owned fixtures can run, but no operator-installed service or data is replaced.
- Use synthetic data, explicit aliases for host/user/path/network identities and a public-output inspection pass across text, images, metadata and filenames. Hash sanitized derivatives separately; keep original dumps/stacks/private identities local.
- Every actionable finding is deduplicated against open/closed GitHub records under milestone 25 and parent #444. Functional outcomes control closure; verification is evidence rather than a closure criterion. No original #355/#356 cause is inferred.
- The checklist helper initially required plan.md before planning despite the documented checklist-before-plan sequence. The failed prerequisite is retained, setup-plan copied the real template, and the helper then passed. This is a workflow ordering compatibility observation, not a product defect.
- Optional after-specify/after-plan context hooks are omitted; the core plan reference update supplies active context. Three skill-directed read-only research agents supplied CLI, transport and isolated-UI findings without test/UI launches or formal review requests.

## Project Structure

- specs/139-full-product-audit-v202/: requirements, design, complete audit disposition, verification, release continuation and source delivery.
- internal/cli and cmd/goschedd: narrowly scoped corrections only after demonstrated findings.
- internal/clievaluate: catalogue contract corrections where it masks actual behavior; existing schema stays fixed.
- test/integration and existing package tests: failing-first regressions and fresh functional observations.
- desktop/frontend/e2e and temporary Playwright configuration: existing headless frontend coverage and screenshot artifacts.
- .test-temp/s139/: frozen build/capture inputs, private environment observations, original failures and sanitized publication derivatives.

## Verification and delivery

Complete focused discovery/corrections before freezing final capture source. Run the full Windows/WSL catalogue in the foreground, semantically review every distinct corpus message, finalize/validate/render using original binaries and preserve failure or blocked dispositions. Exercise supplemental surfaces, actual transport and safe UI boundaries without changing frozen capture inputs. Finalize factual source records, run all eight canonical foreground gates and re-audit protected evidence. Publish the official PR automatically, monitor CI and at most two total Codex reviews including automatic review, resolve in-scope findings, then stop for owner final review/merge. Release coordinator #444 stays open until separately approved immutable-source staging and actual public package delivery.
