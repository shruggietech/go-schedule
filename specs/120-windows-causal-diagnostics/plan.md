# Implementation Plan: S120 Windows causal diagnostics

**Branch**: `codex/120-windows-causal-diagnostics` | **Date**: 2026-10-27 | **Spec**: [spec.md](spec.md)

## Summary

Investigate every original #355/#356 failure category through narrow causal experiments, preserve original receipts and implement only demonstrated component corrections. Add missing outer process provenance and bounded original-result retention, then exercise the real owned API/store cancellation boundary. Record unresolved historical causes independently from any newly demonstrated defect.

## Technical Context

**Language/Version**: Existing Go 1.26.0 modules, PowerShell/Bash script fixtures and Windows/WSL environments.

**Primary Dependencies**: Existing os/exec, context, database/sql, modernc SQLite, typed API/IPC and internal/clifixture containment. No dependency or toolchain update.

**Storage**: Existing single-connection memory/file-WAL fixtures; no schema or migration change.

**Testing**: Behavioral before/after regressions for inherited output handles and canceled real pool wait, native six-writer/probe acceptance, API authorization/audit/state contracts, Windows/Linux targeted checks and all eight canonical foreground gates.

**Target Platform**: Native Windows and existing nonroot WSL Linux; hosted disposable runners provide independent acceptance after PR publication.

**Performance Goals**: Normal successful payload/output behavior unchanged; diagnostics use existing owned resource limits and bounded per-experiment collection. Preserve current supported case/total budgets and explicit caller limits.

**Constraints**: No installed-system/session operation, public diagnostic listener, raw credential export, original-evidence overwrite, speculative rollback, source merge or release operation.

**Scale/Scope**: Two historical investigations, two targeted current-component hypotheses, proportionate corrections and source-bound maintainer guidance. New bugs receive separate issues only after demonstrated.

## Constitution Check

Before research and after design: all five principles pass. Standard-library/context additions preserve idiomatic errors and explicit goroutine ownership. Real regressions precede fixes; writer, migration, state-completeness, authorization/audit, race and coverage contracts remain required. No unmeasured hot-path optimization or new production dependency. All eight local gates precede commits. No constitution deviation.

## Project Structure

- test/integration/testscripts_{test,binding_test,probe_test}.go and a focused process-observation test helper: outer provenance, bounded collection and original result retention.
- internal/clifixture/: reuse existing process/tree/pipes ownership rather than duplicate it.
- internal/store/{access,identity,crud,remote_access}.go, context-focused tests and test-only pool acquisition/stats seam: explicit request cancellation where demonstrated.
- internal/api/server/{server,access,manifest,tasks,groups,pairing}.go: affected request contexts and existing authorization/audit semantics.
- internal/clievaluate/ipc_diagnostics_test.go and source-bound investigation artifacts: reuse existing phase/partial-data logic, add only missing causal discriminants.
- docs/test-scripts.md and docs/cli-qa.md: concise maintainer operation/interpretation guidance; CHANGELOG.md: actual increment and any architecture decision.
- specs/119-release-execution/{spec,delivery,verification}.md and specs/README.md: factual S119 merge reconciliation bundled with S120.

## Execution design

1. Freeze originals, exact accepted source and runtime identities. Identify each experiment and condition separately.
2. Demonstrate inherited-handle and canceled-pool-wait behavior on unchanged source through owned fixtures; preserve primary exits and server completion, not just client cancellation. File any demonstrated bug with parent #355 or #356 under existing milestone 16.
3. Correct the demonstrated boundary. For process observation, use existing contained capabilities for new bounded runners and preserve launch/native/collection/cleanup separately. For SQL, use explicit context-aware methods and legacy delegates on affected paths, retaining intent/terminal audit semantics.
4. Run targeted native experiments for metadata-only versus writer and implicit versus explicit prerequisite selection; compare Windows/Linux where it discriminates a hypothesis. No unchanged full-catalog replay.
5. Produce individual issue dispositions and maintainer guidance, reconcile S119 history, then freeze all source inputs and run scripts/verify.sh all. Finalize only factual lifecycle/task/verification metadata afterward.

## Complexity Tracking

No violation. Conditional source corrections depend on observed regressions, not assumed historical causes. Added types/helpers remain scoped to the shared process-result and request-cancellation contracts.

## Concrete collection and audit bounds

Owned integration script execution retains the existing 120-second operation allowance and permits a two-second output drain after the root exits, then stops the owned tree and closes only owned pipes. Stdout and stderr are captured separately with a one-MiB limit each; overflow remains a collection error and cannot replace the native exit. Opt-in provenance records the resolved executable path and SHA-256, harness path and SHA-256, root PID, elapsed time, actual native status, collection status and completed existing phase markers. No environment values or harness contents are logged. Runtime version experiments are separate named invocations.

Request cancellation applies to authorization, identity reads, audit intent and affected snapshot/cleanup SQL. Terminal audit persistence after an operation has executed uses a separate 30-second completion allowance, because disconnecting the caller must not discard an already committed operation outcome. A completion timeout retains the uncertain intent and returns the existing attributable audit-persistence error. This bound is independent of CI budgets.
