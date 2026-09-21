# Security Requirements Checklist: Agent Grant Controls

**Purpose**: Validate that S094's authorization, secret, transport, and lifecycle requirements are explicit before planning
**Created**: 2026-09-20
**Feature**: [spec.md](../spec.md)

## Grant Boundaries

- [x] Observe, Operate, and Manage have explicit and testable authority boundaries
- [x] Existing grants can only be narrowed, expired sooner, or revoked from the desktop
- [x] Existing connections must re-evaluate actor state on the next request
- [x] Wrong-daemon, expired, and revoked behavior is fail closed

## Secret Handling

- [x] Durable credentials and protected secrets are excluded from all interface projections
- [x] The one-time pairing phrase is copied only through the native boundary
- [x] Clipboard failure cancels the unexchanged pairing
- [x] Audit presentation excludes bodies, task inputs, secrets, and raw errors

## Transport Isolation

- [x] Stdio, localhost HTTP, and remote HTTPS states are independently represented
- [x] Listener lifecycle is explicitly separate from permission grants
- [x] Enabling one transport cannot imply or enable another transport

## Evidence

- [x] Recent shared audit evidence is bounded and actor-associated
- [x] Epic-level conformance, compatibility, redaction, and hostile-content gates remain required

## Notes

- All security requirements are specific, measurable, and covered by acceptance scenarios.
