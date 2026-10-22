# S111 research

## Workflow

Executed installed /speckit-specify and /speckit-clarify; all five clarifications were answered from owner scope, constitution and existing contracts. Requirements checklist passes 16/16. /speckit-checklist covers SMTP security, delivery and parity. /speckit-plan setup used the installed PowerShell script. Its optional agent-context hooks are handled by the prescribed CLAUDE.md plan-reference update. No mandatory extension hook is registered.

The installed planning command explicitly dispatches research agents. Three read-only research tasks inspected transport/library choices, persistence/bundles and administration/evaluation surfaces; no implementation was delegated.

## SMTP client

**Decision**: github.com/emersion/go-smtp v0.25.0 and its go-sasl dependency, MIT.

**Rationale**: Explicit Mail/Rcpt/Data APIs expose per-recipient outcomes and final submission confirmation. The release was published 2026-08-13. Library commands reset socket deadlines, so context cancellation must close the owned connection across greeting, STARTTLS, authentication and final DATA.

**Alternatives considered**: Standard-library net/smtp is frozen and does not meet the maintained-library requirement; wneessen/go-mail v0.8.1 adds useful MIME conveniences but gives less direct control over durable transaction progress. Do not add both.

Primary sources: [release](https://github.com/emersion/go-smtp/releases/tag/v0.25.0), [MIT license](https://github.com/emersion/go-smtp/blob/v0.25.0/LICENSE), [client implementation](https://github.com/emersion/go-smtp/blob/v0.25.0/client.go), [dependency manifest](https://github.com/emersion/go-smtp/blob/v0.25.0/go.mod), [SMTP completion semantics](https://datatracker.ietf.org/doc/html/rfc5321#section-4.2.5).

## Delivery semantics

**Decision**: Reuse the queue, persist bounded recipient outcomes and checkpoint submission before DATA; atomically persist final outcomes with retry/terminal transitions.

**Rationale**: RCPT acceptance alone does not confirm a message. Final DATA success confirms only those RCPT recipients. Temporary failures retry within three attempts; permanent rejection does not. Interrupted submission remains uncertain and may be retried with duplicate-risk disclosure. Stable message IDs aid correlation but are not deduplication guarantees.

**Alternatives considered**: Resend the full recipient list (duplicates confirmed recipients); fail an entire recipient group at the first RCPT rejection (suppresses valid recipients); build a separate mail queue (duplicates existing lifecycle).

## Stored state and secrets

**Decision**: Add v23, rebuilding only the constrained channel table while preserving IDs and foreign keys, and add typed transport snapshots and recipient evidence to deliveries.

**Rationale**: v22 CHECK permits webhook only. Existing Complete/Recover erasure paths must retain credential scrubbing. Preserve Windows DPAPI and Unix daemon-only DB behavior rather than claiming uniform encryption at rest.

## Remote and desktop parity

**Decision**: Add the complete notifications family to remote operations/OpenAPI and remove the desktop local-only setup gate after route support is established.

**Rationale**: Existing local authorization catalog already defines Observe/Manage/Operate/Enroll boundaries, but remote routes are absent and the UI explicitly gates remote notifications. API client request additions alone would be insufficient.

## Portable configuration

**Decision**: Optional typed SMTP channel definitions in supported portable v2, with explicit target-local provisioning/credential binding and frozen binding identities.

**Rationale**: Existing bundles carry policy assignments bound by channel name but no destination definitions. Export must include nonsecret SMTP settings and a credential requirement without exporting credentials; missing or stale bindings must block apply.

## QA and messages

**Decision**: Shared owned SMTP/TLS fixtures, redacted authored diagnostics and existing evaluator report schemas.

**Rationale**: No production credentials or user interaction is needed; numeric SMTP codes and phase categories are useful without publishing arbitrary receiver text or authentication transcripts.
