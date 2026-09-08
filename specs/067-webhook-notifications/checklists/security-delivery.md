# Security and Delivery Requirements Checklist: Dependable Webhook Notifications

**Purpose**: Review high-risk secret, transport, durability, retry, and scheduler-isolation requirements before implementation
**Created**: 2026-09-07
**Feature**: [spec.md](../spec.md)

## Secret Boundary

- [x] CHK001 Are every accepted secret field and every prohibited output surface explicitly named? [Completeness, Spec §FR-004 and FR-005]
- [x] CHK002 Is the storage protection claim limited to the daemon-owned database, restricted IPC, and operating-system file permissions? [Clarity, Spec §FR-005 and Assumptions]
- [x] CHK003 Are terminal completion and channel-removal secret-erasure requirements defined? [Recovery, Spec §FR-023]
- [x] CHK004 Are credential-bearing URL components excluded from redacted summaries and diagnostics? [Security, Spec §Edge Cases and FR-004]

## Outbound Transport

- [x] CHK005 Are allowed schemes, loopback development behavior, URL rejection cases, and redirect behavior unambiguous? [Clarity, Spec §FR-006]
- [x] CHK006 Is authorization forwarding constrained to the originally validated destination? [Security, Spec §FR-015]
- [x] CHK007 Are success status, timeout, attempt count, and backoff values measurable? [Measurability, Spec §FR-017]
- [x] CHK008 Is receiver deduplication tied to a stable documented identifier without an exactly-once promise? [Consistency, Spec §FR-018]

## Durability and Isolation

- [x] CHK009 Is the atomic boundary between run persistence and delivery creation specified? [Consistency, Spec §FR-011]
- [x] CHK010 Is notification execution explicitly outside scheduler task-worker capacity? [Performance, Spec §FR-012 and SC-002]
- [x] CHK011 Are claimed-delivery recovery and attempt preservation stated for daemon restart? [Recovery, Spec §FR-019]
- [x] CHK012 Are pending-work preservation and terminal-history pruning rules compatible and bounded? [Consistency, Spec §FR-022]
- [x] CHK013 Is source-run immutability stated for every notification failure and lifecycle path? [Safety, Spec §FR-020]

## Lifecycle and Evidence

- [x] CHK014 Are configured and effective policies independently inspectable with their source scopes? [Coverage, Spec §FR-010]
- [x] CHK015 Are test deliveries required to avoid fabricating task runs? [Data Integrity, Spec §FR-016]
- [x] CHK016 Is channel removal defined for assignments, unfinished deliveries, secrets, and terminal evidence? [Completeness, Spec §FR-023]
- [x] CHK017 Are delivery filters and response limits complete enough for bounded incident inspection? [Coverage, Spec §FR-021]
- [x] CHK018 Are validation requirements traceable to unit, contract, restart, race, installed-daemon, and canonical gates? [Traceability, Spec §FR-025 and SC-006]

## Notes

- Formal reviewer-depth checklist passes all 18 requirements-quality checks.
- The checklist evaluates the written security and delivery contract, not implementation behavior.
