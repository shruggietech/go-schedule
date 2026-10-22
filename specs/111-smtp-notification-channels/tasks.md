# Tasks: SMTP notification channels

**Feature**: S111, #176 in full; completes #19 on functional delivery.

## Phase 1: Setup

- [x] T001 Record resolved dependency/license and workflow decisions in specs/111-smtp-notification-channels/research.md and CHANGELOG.md.
- [x] T002 Add pinned maintained SMTP dependency in go.mod/go.sum and retain desktop module compatibility.

## Phase 2: Foundations

- [x] T003 Write settings/credential/injection regressions in internal/notification/smtp_test.go (FR-001 through FR-004).
- [x] T004 Add typed settings, recipient/result and validation contracts in internal/domain/notification.go and internal/notification/smtp_config.go.
- [x] T005 Write prior-v22 preservation and recipient/restart tests in internal/store/notifications_smtp_test.go (including its prior-v22 migration fixture) (FR-007/FR-010).
- [x] T006 Add v23 migration and transport snapshots/recipient checkpoints in internal/store/store.go, internal/store/notifications.go and internal/store/notifications_smtp.go.

## Phase 3: User Story 1, reliable unattended email

- [x] T007 [US1] Add owned deterministic SMTP/TLS fixture and transport regression tests in internal/testsmtp/ and internal/notification/smtp_test.go (FR-002/FR-003/FR-006/FR-007/FR-013).
- [x] T008 [US1] Implement bounded MIME rendering and SMTP transactions in internal/notification/smtp.go and internal/notification/mail.go.
- [x] T009 [US1] Integrate recipient-aware dispatch and atomic completion in internal/notification/dispatcher.go while retaining webhook behavior (FR-005/FR-007).
- [x] T010 [US1] Cover task/test/health events, partial rejection, retries, cancellation and restart in internal/notification/ and internal/store/ tests.

## Phase 4: User Story 2, administration parity

- [x] T011 [US2] Add shared API settings/authentication and wrong-kind validation tests in internal/api/server/notifications_test.go (FR-003/FR-004/FR-008).
- [x] T012 [US2] Extend create/update/rotation contracts and API client methods in internal/api/server/notifications.go and internal/api/client/notifications.go.
- [x] T013 [US2] Add CLI SMTP request/secret-input/output regressions in internal/cli/notification_smtp_test.go.
- [x] T014 [US2] Implement complete transport-aware CLI flags, bounded secret file/stdin provisioning and recipient history in internal/cli/notification.go (FR-008/FR-011).
- [x] T015 [US2] Extend remote registry/OpenAPI/generated contracts and capability/target regressions in internal/remote/, api/openapi/remote-v1.yaml (FR-008).
- [x] T016 [US2] Extend desktop draft/view models and binding tests in desktop/notifications/ (FR-004/FR-008/FR-014).
- [x] T017 [US2] Implement accessible SMTP form, credential clear/replace, transport/recipient history and remote availability in desktop/frontend/src/notifications/ and desktop/frontend/src/App.tsx.
- [x] T018 [US2] Extend frontend unit/browser contracts for settings/error/keyboard/secret behavior in desktop/frontend/src/notifications/ and desktop/frontend/e2e/notifications.spec.ts.

## Phase 5: User Story 3, portable configuration and unattended evaluation

- [x] T019 [US3] Add typed SMTP portable-definition/binding/drift regressions in internal/bundle/ and internal/api/server/bundles_test.go (FR-009/FR-015).
- [x] T020 [US3] Extend portable validation/canonicalization/export/preview/apply/compare and stale binding checks in internal/bundle/ and internal/api/server/bundles*.go.
- [x] T021 [US3] Add owned SMTP fixture lifecycle, redaction and receiver/delivery observations in internal/clievaluate/ (FR-012/FR-013).
- [x] T022 [US3] Extend scenario catalog, intentional failure coverage and independent postconditions; regenerate internal/clievaluate/inventory-v1.json.
- [x] T023 [US3] Run fresh receipt-bound isolated evaluation and review contextual messages using existing canonical report contracts.
- [x] T024 [US3] Update docs/notifications.md, docs/cli.md, docs/api.md, docs/remote-access.md, docs/cli-qa.md and docs/gui-fields.md with complete behavior and provisioning guidance (FR-011).

## Phase 6: Cross-cutting completion

- [x] T025 Reconcile issue-level traceability and spec lifecycle in specs/README.md and specs/111-smtp-notification-channels/delivery.md.
- [x] T026 Run focused security, migration, remote, frontend and transport regressions and record actual results in specs/111-smtp-notification-channels/verification.md.
- [x] T027 Run foreground sh scripts/verify.sh all and resolve every required gate without weakening tests.
- [x] T028 Run github-format, audit encoding/diff/issue scope, mark delivered tasks and record local completion in specs/111-smtp-notification-channels/.

## Dependencies and execution

T001-T006 establish shared contracts/storage. US1 T007-T010 depends on foundations. US2 API precedes CLI/desktop/remote integration; each surface uses the same settings and credentials contract. US3 portable configuration and evaluator depend on settled authoring/transport semantics. All three stories and final verification ship together; the owner authorized the full slice, not an SMTP-only partial MVP.

The conceptual MVP is US1 plus supported authoring, but it is not this slice's delivery boundary. Independent tests are described in spec.md. Different UI/docs files can be edited independently after contracts settle; shared API/domain/store files are serialized. No agent implementation delegation is planned.

Publication, review and merge are not unchecked implementation tasks. Owner authorization already covers push/official PR and in-scope fixes, with no more than two Codex rounds and no merge.
