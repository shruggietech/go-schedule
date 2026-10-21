# Implementation Plan: CLI contract and diagnostic repairs

**Branch**: codex/109-cli-contract-diagnostics | **Date**: 2026-10-21 | **Spec**: [spec.md](spec.md)

**Input**: Full #311-#322 bundle, parent #277, milestone 15. Baseline c6d50a17c4b16a9025f9f7137e3cc05b0a6391bd.

## Summary

Repair the complete CLI invocation, validation and communication contract using existing Cobra handlers and shared application boundaries. Reuse S108 for comprehensive Windows/WSL evaluation. Preserve operational error classification, valid command forms, secret handling and consumers.

## Technical Context

**Language/Version**: Go 1.26.0, existing PowerShell entry point and TypeScript desktop consumers.

**Primary Dependencies**: Existing Cobra 1.10.2, pflag and shared client/server/domain packages. No dependency additions.

**Storage**: Existing SQLite and disposable QA stores, no schema migration.

**Testing**: Go regression tests before fixes, race detector, existing frontend tests, full eight-gate scripts/verify.sh all and fresh receipt-bound S108 capture/review/finalize.

**Target Platform**: Windows and Linux, supported macOS compilation and existing desktop compatibility gates.

**Project Type**: CLI with shared daemon API and desktop clients.

**Performance Goals**: No scheduling hot-path change, no blocking reconciliation waits or new polling; command metadata traversal is bounded by the initialized tree.

**Constraints**: Hidden noninteractive child processes; owned fixtures only; no production secrets or installed-service mutation; UTF-8 without BOM; GitHub format gate; stop before push/PR.

**Scale/Scope**: All 12 defects in one increment, 152 initialized command nodes plus bounded bundle/enrollment/watcher response paths.

## Constitution Check

Pre-research and post-design: PASS. I: typed/wrapped errors and proportionate helpers. II: meaningful failing regressions, deterministic observation tests and full race/coverage gates. III: exact result/diagnostic and help/state contracts. IV: no waits or new unbounded work. V: installed command workflow, local review branch and one pre-publication halt. Persisted schemas, authorization and credentials remain compatible. No waiver or principle deviation is needed.

## Phase 0: Research

Installed /speckit-plan research delegates inspect target/bundle and enrollment/watcher consumers read-only. [research.md](research.md) records the decisions and rejected alternatives. No external mutable platform standard is needed for these source-defined contracts.

## Phase 1: Design and contracts

- Shared invocation setup wraps known parser/operand validation as usage; command lookup rejects unknown commands without matching runtime error prose. Existing explicit operand rules are preserved, and omitted product-command rules become NoArgs before side effects.
- Add an unambiguous global --daemon-endpoint selector, retain legacy --endpoint and channel-local destination --endpoint, inspect global flag ownership rather than a same-name child flag, and refuse simultaneous alias selection.
- Validate duplicate bindings before map conversion and unknown bindings at the shared bundle service before storing a plan. Preserve parser detail and identify the input file.
- Complete cron reports before returning invalid-input exit 2; preserve operational exit 1 and successful unsupported-line reporting.
- Contextual CLI diagnostics retain typed causes and API code/field while naming the object and useful next step. Do not echo secrets or arbitrary payloads.
- Populate rotation association context without a post-mutation lookup that could lose token delivery or impose new authentication requirements.
- Correlate watcher observations with the observed definition revision, including unchanged health; expose additive observation/reconciliation context and describe it in the CLI without waits or fabricated state.
- Add meaningful help and headings, route alerts notices to stderr, and explain compatibility configuration in the existing structured warning.

[Data model](data-model.md), [interface contract](contracts/cli.md) and [quickstart](quickstart.md) define concrete boundaries. Update CLAUDE.md's existing Spec-Kit plan reference; the optional agent-context hook is satisfied by this focused refresh without regenerating unrelated project instructions.

## Project Structure

Feature records: specs/109-cli-contract-diagnostics/{spec,plan,research,data-model,quickstart,tasks,analysis,verification,delivery}.md and contracts/cli.md, checklists/*.md.

Implementation: internal/cli/{cli,invocation,target,bundle,cron,access,service,profile,notification,watcher,schedule}.go, internal/bundle and bundle service, internal/enrollment, internal/domain, internal/watcher, internal/api/server and cmd/goschedd. Regressions live beside these packages; existing desktop/API consumers are checked if response fields change. Regenerate internal/clievaluate/inventory-v1.json after metadata changes, preserve existing case IDs and extend refusal scenarios as necessary.

## Dependencies and execution order

S108 and S107 are merged prerequisites. Foundation (#311/#312) precedes integrated refusal paths (#315/#316) and contextual diagnostics (#318). Target ownership (#314) precedes affected help/docs/inventory. Response semantics (#319/#321) precede state wording. Human output (#313/#317/#320/#322) can be developed independently in separate files once contracts are established. The final aggregate integrates all three user stories; no issue is intentionally deferred.

## Complexity Tracking

No constitution exception. Reuse initialized metadata, existing adapters and the existing evaluator; do not create a second generic validation or report framework.
