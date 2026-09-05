# Research: External Trigger Lifecycle

## Decision 1: Reuse the existing protected local API

The existing daemon-owned IPC transport remains the only trigger ingress. The CLI sends a normal authenticated local API request, so S054 adds no TCP listener, socket protocol, resident helper, or extra goroutine.

## Decision 2: Use opaque 256-bit trigger keys

Each key is generated from 32 cryptographically random bytes and encoded as unpadded base64url with a `gst_` prefix. This is easy to pass as one CLI argument, contains no shell metacharacters that require quoting in ordinary use, and provides ample resistance to guessing.

## Decision 3: Persist recoverable keys in the protected local database

The desktop requirement includes revealing and copying an existing key after restart. Hash-only storage cannot satisfy that workflow, so this local-only design persists the recoverable value inside the existing user-scoped SQLite database. Raw keys are returned only by create, rotate, and explicit reveal operations and are excluded from normal list, detail, history, event, and error payloads.

## Decision 4: Dispatch through the scheduler engine

An accepted fire request enters the same overlap, worker-limit, execution, cancellation, event, and history path as manual and scheduled runs. The new source is `external_trigger`, with the trigger identifier retained as provenance and the key omitted.

## Decision 5: Treat enabled triggers as automatic sources

An enabled trigger makes a command-ready task eligible for activation even if the task has no schedule or incoming completion chain. Creating or enabling a trigger never silently activates a task. Removing, disabling, or retargeting the final enabled automatic source atomically deactivates the old target so persisted state cannot claim an impossible automatic configuration.

## Decision 6: Validate eligibility before dispatch

The fire path distinguishes unknown key, disabled trigger, missing target, incomplete command, inactive task, disabled task, blocked ancestor group, and unavailable dispatch capacity. Stable machine-readable codes support CLI messaging and future integrations without exposing the key.

## Decision 7: Make rotation immediately exclusive

Rotation replaces the stored key in one transaction. The old key becomes unknown as soon as the operation succeeds, and the new value is shown once in the response. Rotation and deletion require confirmation in the desktop interface.

## Decision 8: Keep S054 bounded

This slice delivers single-trigger lifecycle and desktop management from issues #132 and #133. Trigger sets from #134 and filesystem watchers from #135 remain separate because they add bulk identity, fan-out, filesystem semantics, and materially different validation risks.
