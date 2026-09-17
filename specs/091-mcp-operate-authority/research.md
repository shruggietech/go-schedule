# Research: MCP Operate Authority

## Decision 1: Preserve Observe as a separate default contract

`mcpobserve.NewServer` continues to advertise resources only. Operate tools are added by a separate package only after a validated Operate session exists.

**Rationale**: Tool registration is the discovery boundary. An Observe client cannot invoke a tool that was never registered, even if hostile resource content resembles a call.

## Decision 2: Use runtime MCP actors and local session secrets

The protected local API creates a short-lived MCP actor and returns a random session secret once. The server stores only its digest and resolves the actor on every request. Revocation removes the digest and revokes the actor.

**Rationale**: Calling the local API through the ordinary client without a session would attribute every mutation to the built-in operating-system actor. A memory-only identity bridge preserves the existing API authorization and intent-first audit implementation without adding durable grant administration.

## Decision 3: Require exact daemon and task identity

Every tool requires `daemon_id`, `task_id`, and `request_id`. The adapter reads the current manifest before mutation and refuses a mismatch.

**Rationale**: MCP hosts can retain stale context or connect to several servers. Explicit target identity prevents an apparently valid call from reaching the wrong scheduler.

## Decision 4: Deduplicate at the MCP operation boundary

Each Operate server owns a bounded, expiring cache keyed by request identifier. The cache stores the exact operation fingerprint and first result. An identical retry returns that result; conflicting reuse is rejected.

**Rationale**: Run-now is not naturally idempotent. Automatic retry after a lost response could otherwise execute a task twice.

## Decision 5: Return domain outcomes as structured tool results

Tool responses use accepted, rejected, denied, or uncertain outcomes. API validation and not-found responses are rejected; authorization failures are denied; connection or deadline ambiguity is uncertain.

**Rationale**: Returning every failure as a generic protocol error loses the operational distinction required for safe automation.

## Decision 6: Keep Manage and remote authorization out of S091

Task definition changes, deletion, durable grants, expiry controls, and remote OAuth-style MCP remain in #179, #181, and #180.

**Rationale**: Those capabilities introduce different security and interaction boundaries and would make this slice materially harder to review.
