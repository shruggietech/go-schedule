# Implementation Plan: SMTP notification channels

**Branch**: `codex/111-smtp-notification-channels` | **Date**: 2026-10-21 | **Spec**: [spec.md](spec.md)

**Input**: Owner-authorized S111, #176 in full and notification coordinator #19.

## Summary

Add typed SMTP destinations to the existing durable notification pipeline. Preserve webhook wire behavior, policy inheritance and the four bounded workers. Integrate one additive schema migration, recipient checkpoints, supported administration surfaces, portable configuration and owned SMTP/TLS evaluation. The owner explicitly authorized push/official PR and at most two Codex rounds, overriding the default pre-publication pause; stop before merge.

## Technical Context

**Language/Version**: Go version and toolchain pinned by existing go.mod; Wails Go module and React/TypeScript frontend retain their current versions.

**Primary Dependencies**: Add github.com/emersion/go-smtp v0.25.0 (MIT) and its go-sasl dependency for explicit SMTP transaction control. Use standard-library MIME, mail parsing, TLS and HTML escaping.

**Storage**: Existing SQLite v22 advances to v23 without destructive migration. Retain existing secretstore protection semantics (Windows DPAPI; daemon-only database permissions on Unix).

**Testing**: Existing Go race tests, storage migration tests, frontend tests/native build, OpenAPI generation and versioned unattended CLI evaluation.

**Target Platform**: Windows, Linux and macOS daemon/CLI/desktop.

**Project Type**: Existing daemon, CLI, local/remote JSON API and desktop application.

**Performance Goals**: Scheduling hot paths retain the existing p99 dispatch budget; outbound mail stays in four bounded notification workers. One attempt owns a shared timeout of 1-60 seconds (default 10 seconds), at most 100 unique recipients and a bounded message.

**Constraints**: No production email account, installed-service mutation, public release or merge. No credential values in normal output, export, logging, fingerprints or diagnostic messages. No TLS verification bypass.

**Scale/Scope**: One end-to-end notification feature, existing durable history bound and three-attempt budget. No separate queue or provider-specific integration.

## Constitution Check

| Principle | Before research | After design |
| --- | --- | --- |
| I, idiomatic code and defined ownership | Pass: existing package boundaries | Pass: explicit owned sockets, cancellation closes transport; no protocol transcript logging |
| II, regression/security/migration/race coverage | Pass: tests precede behavior | Pass: prior-v22 preservation, secret/TLS/authorization/recipient/recovery cases and canonical eight gates required |
| III, interface consistency | Pass: full desktop/headless scope | Pass: shared typed validation, human/JSON diagnostics, remote capability reuse and maintained inventory |
| IV, bounded performance | Pass: no scheduler hot-path expansion | Pass: bounded recipient/message/time limits and existing worker pool |
| V, Spec-Kit and PR integration | Pass: source issue and owner kickoff | Pass: specify/clarify/checklist/plan/tasks/analyze before implementation; authorized push, two-round maximum, no merge |

## Research and decisions

See [research.md](research.md), [data-model.md](data-model.md) and [contracts/smtp.md](contracts/smtp.md).

1. Use go-smtp rather than implementing SMTP or using a second high-level queue. Its explicit transaction API exposes RCPT and final DATA acceptance.
2. Use typed nonsecret settings plus the existing protected authorization slot for serialized SMTP username/password. Public requests accept a typed write-only SMTP authentication object; old webhook authorization input remains compatible.
3. Keep kind immutable after creation. Default an omitted kind to webhook for existing callers.
4. Snapshot settings and credentials when work is created. Edits/rotation affect new work; disable stops new work; existing deletion semantics remove unfinished work.
5. Persist recipient outcomes before retry/terminal transitions. Checkpoint recipients entering DATA as submitting; recovery marks interrupted submission uncertain. Retry pending/temporary/uncertain recipients within the existing three-attempt budget, with explicit duplicate-risk wording for uncertain submission. Confirmed accepted and permanent failures are excluded.
6. Use one SMTP transaction per attempt for eligible recipients; DATA confirmation applies only to its accepted RCPT subset. A failed QUIT after successful DATA does not turn accepted mail into failure.
7. Close the owned socket from context.AfterFunc to enforce the shared attempt deadline even when the SMTP library resets deadlines.
8. Complete the existing remote notification route family, which research found absent from the remote allowlist and desktop gating. Reuse existing operation identities/capabilities, including Enroll for secret rotation.
9. Extend portable v2 with optional typed SMTP definitions and unique-name target-local provisioning/bindings. Reuse reviewed channel IDs/UpdatedAt and target fingerprints; never export or hash secret values. Existing webhook-only documents remain valid.
10. Preserve the existing unattended report schema; extend cases, observations, owned receivers and initialized inventory.

## Project Structure

### Documentation

spec.md, plan.md, research.md, data-model.md, contracts/smtp.md, quickstart.md, checklists/requirements.md, checklists/smtp.md, tasks.md, analysis.md, verification.md and delivery.md under specs/111-smtp-notification-channels/.

### Source Code

- internal/domain/notification.go: typed SMTP settings, credentials input and recipient/result contract.
- internal/notification/: validation, MIME rendering, SMTP transport, routing and dispatcher.
- internal/store/: v23 migration, settings snapshots and atomic recipient completion.
- internal/api/server/notifications.go, internal/api/client/notifications.go: shared authoring and rotation contracts.
- internal/remote/, api/openapi/: notification route parity, bounded bodies, generated contract and authorization tests.
- internal/cli/notification.go: SMTP flags, credential provisioning and useful transport-aware output.
- desktop/notifications/ and desktop/frontend/src/notifications/: shared drafts/view models and accessible forms/history.
- internal/bundle/ and internal/api/server/bundles*: typed portable SMTP definitions and explicit local provisioning/bindings.
- internal/clievaluate/, internal/testsmtp/: owned mail fixtures, lifecycle/failure observations and inventory.
- docs/: self-contained configuration and behavior guides.

## Execution and verification

Execute the dependency-ordered tasks after the blocking analysis gate. Run focused failing regressions before implementation, then affected packages and UI checks, fresh isolated evaluation inputs, and foreground `sh scripts/verify.sh all`. Save exact outputs and source identity in verification.md. Run github-format before repository/GitHub publication. Keep publication/review stages out of unchecked implementation tasks.

## Complexity Tracking

No constitution exception. The new dependency is justified by standards-compliant explicit SMTP transaction control and reviewed MIT licensing. The v23 table rebuild is necessary because v22 restricts channel kind to webhook. Recipient checkpoints are required by partial acceptance/recovery semantics, not a second queue.

## Implementation decisions, 2026-10-22

Portable definitions use the existing logical-name policy convention instead of adding opaque portable IDs or a second secret-provisioning interface. Exactly one provisioned same-name local SMTP channel, matching credential presence and credential authority, is required. Preview freezes ID/UpdatedAt and apply preflights every binding before writes. Stored-secret destination changes reuse Enroll. This delivers typed portability without credential export, copies or a parallel identity system.

## Diagnostic findings and corrections, 2026-10-22

A complete first capture observed 1,760 successful executions, three failures and three platform exclusions. #342 records the new SMTP scenario mismatch (text acknowledgement incorrectly required to be JSON); exact run/delivery/receiver/webhook assertions remain. #343 records quoted mailbox serialization and authoring separators; canonical envelopes, MIME parsing and quoted-list CLI/desktop inputs preserve identity. #344 records a legacy bounded history query returning NULL for empty text; COALESCE retains empty strings without changing byte limits, pagination or truncation. The diagnostic run completed owned cleanup before source changed. A fresh source-bound capture and canonical aggregate follow the corrections.
