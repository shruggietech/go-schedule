# S144 data model

- ParallelJoin: ID/name, entry task, target task, ordered named branches and creation/update timestamps. Each ParallelJoinBranch has its own stable ID/name and TaskID; TaskID need not be unique. Optional expected update revision is a mutation precondition, not stored intent.
- ParallelJoinInstance: unique (DefinitionID, EntryRunID), immutable definition/name/entry/target/name snapshots, ordered invocation snapshots, state/reason, target delivery/run/outcome, created/updated/completed timestamps. Runtime identifiers remain attributable after live-object removal.
- ParallelJoinInvocation: unique immutable invocation ID plus definition BranchID, name/task snapshot, exact DeliveryID, delivery state, final RunID/outcome/reason. Same-task slots are independent.
- CompletionDelivery: existing durable queue plus Role=chain|join_branch|join_target and JoinInstanceID/JoinBranchID (invocation ID). Separate role-specific unique indexes prevent duplicate native admission decisions without fake chain IDs.
- Run: add optional DeliveryID, JoinInstanceID, JoinBranchID and JoinRole=branch|target, persisted and carried through active/queued/skipped/final metadata. Ordinary source task/run remain the entry for join-owned work.

Instance transitions: successful entry -> waiting; unsuccessful required contribution/resolution -> blocked; all exact successes -> target_pending plus one target delivery; target terminal result -> completed with actual outcome; definition/task invalidation -> resolved. Blocked instances retain remaining branch work but cannot create a target. Resolved instances never reopen; late native evidence is still retained. Claimed is not running; authoritative active metadata supplies the actual running phase.

Current schema 25 initializes directly. No historical database conversion/deletion. Portable schema v3 contains reusable definitions/portable task references only, never instances, deliveries or run/output data. Mixed chain/join replacement validates and writes the projected final graph atomically.
