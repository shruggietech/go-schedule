# S144 research and decisions

Three research agents were dispatched by the installed speckit-plan command for persistence, engine/API/portable integration and interfaces. Their source investigations are consolidated below; the settled contract is plan.md and contracts/joins.md. No production file was edited during research.

## Identity and delivery

Select an explicit reusable parallel block that launches its own named invocations after entry success. Ambient latest-task history cannot disambiguate overlapping entries or repeated task references. Add role/instance/invocation metadata to the existing completion queue, with separate partial unique indexes. Preserve the entry task/run as the stable source for branch and target work. Each invocation has its own delivery identity; no fake chain ID or task-ID deduplication.

## Atomic decisions and recovery

Select the existing run-recording transaction for branch contribution, new instance creation and unique target admission. A repeated identical committed observation is idempotent; conflicting identity reuse is rejected. Reuse current claims/resolution/recovery with exact immutable bindings. Claimed is not running: sample existing active runs before stored instances and merge by DeliveryID. An extra attempt table and early synthetic persisted Run rows are unnecessary. Preserve the external-command at-least-once crash boundary.

## Mutable definitions and current storage

Select schema25 current initialization, immutable runtime snapshots and non-cascading historical identifiers. Explicit task/definition removal resolves outstanding decisions; a native command that already began retains its result without reviving the resolved instance. Edits affect future bindings only. Pending snapshot destinations remain activation sources, avoiding indirect auto-disablement after an edit. Typed task.Sources replaces positional optional booleans. No historical database conversion or automatic deletion.

## Graph and portable contracts

Select one combined ordinary-chain/join graph validator and one replacement transaction for the projected portable graph. Independent per-kind validators can reject valid swaps or admit a combined cycle. Portable v3 introduces joins and excludes all runtime state; maintained v1/v2 readers retain their previous semantics, refusing join fields under those older declarations. This is portable-file decoding, not database migration.

## Confirmed proportional corrections

- #471: a deterministic private race-enabled source overlay against accepted28cd3ba reproduced queued target execution after task disablement and after chain deletion. Both assertions failed, without race warnings. Reuse one fresh automatic admission check for initial dispatch, queue release and worker acquisition. Preserve first-pending queue-one collapse and manual semantics.
- #472: a private desktop source overlay supplied MutationUncertainError from chain creation/deletion and trigger creation. All three incorrectly returned unavailable/Try again. Fix the shared Automation classifier and same-target frontend/facade reconciliation for every mutation. Retain original red receipts privately.

## Interfaces

Select explicit join CRUD and bounded definition/instance observation with nested branch paging. Add remote chain CRUD because the full existing Automation workspace loads chains alongside joins/triggers/watchers; a joins-only degraded remote page is not selected. Reuse current JSON authority/audit/retry classification, selected-daemon identity, CLI family/input conventions, controlled MCP cache/revocation and keyboard-accessible desktop forms. No new orchestration runtime, approval policy or visual-canvas framework.

## Competitive and execution boundaries

The owner prioritizes the dependent-workflow gap in the maintained competitor matrix. This source increment improves parallel orchestration without rewriting published v2.0.3 comparisons/reviews or claiming full visual-language parity. Cluster work remains explicitly deferred. All execution remains hidden, unattended and no-focus; actual native Wails UI is not used. Source publication and in-scope CI/review fixes are already authorized, limited to two total Codex reviews and owner final merge.
