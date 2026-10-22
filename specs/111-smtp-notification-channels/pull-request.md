# S111: SMTP notification channels and complete administration parity

SMTP email now uses the existing daemon-owned notification queue, task/group conditions and delivery history. CLI, desktop, local IPC and authenticated remote clients can administer the same channels, with write-only credentials and recipient-level results.

## What this changes

- Add maintained SMTP transport with verified STARTTLS/implicit TLS, explicit plain relay, bounded attempts, Unicode MIME text and optional escaped HTML.
- Persist frozen settings and recipient checkpoints through v23 migration, partial acceptance, retry and restart; confirmed recipients are excluded from later attempts.
- Add complete authoring, credential rotation/clear, test and history support, including quoted mailbox handling across CLI, desktop and MIME.
- Restore the missing remote notification routes using existing capability boundaries. Changing the destination authority of stored credentials requires Enroll.
- Add typed portable SMTP definitions bound to provisioned target-local channels by name and frozen revision, without exporting credentials.
- Extend owned Windows/WSL evaluation and correct the SMTP task-case rendering contract and pre-existing empty-text bounded history defect.

Closes #176. Closes #19. Closes #341. Closes #342. Closes #343. Closes #344. Refs #146.

Spec-Kit slice: S111, specs/111-smtp-notification-channels/. Existing issue milestones remain authoritative. #146 stays open for its separate distributed-scheduling decision.

## Verification

Foreground `sh scripts/verify.sh all` passed all eight gates: format, vet, lint, race, GUI, coverage, docs and automation. All 161 frontend tests and three notification browser/accessibility tests passed. Core coverage is 80.2%-91.3%.

Fresh receipt-bound Windows/WSL run c29c5da5-15e9-4adb-924d-90b204aed0fa finalized with 1,763 passing executions, three declared platform exclusions, 1,109 contextual message reviews and zero findings. Source-content SHA256: b294faa7d031106214aeab7c14aa5b3cab932b6ea5804295026d91a2aa712050. All owned runtimes/process trees were cleaned. See [verification](specs/111-smtp-notification-channels/verification.md) for exact evidence and the separately retained first diagnostic capture.

- [x] Canonical local verification passed.

## Checklist

- [x] No safety-critical test surface was weakened or skipped.
- [x] Dependency and architecture decisions are recorded in CHANGELOG.md.
- [x] Documentation updated for configuration, access, lifecycle, portability and evaluation.

## Delivery boundaries

This PR provides unreleased source. Owned fixtures exercise SMTP, TLS and authentication without a production mail account. A later real-provider exercise uses the owner's chosen server and credentials. The owner authorized publication, CI and at most two Codex rounds; this PR remains unmerged and auto-merge stays disabled.
