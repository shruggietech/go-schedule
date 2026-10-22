# S111 SMTP contracts

## Authoring and administration

Channel create accepts optional kind (omitted means webhook), optional smtp settings and typed write-only smtp_auth. SMTP requires typed settings and rejects webhook endpoint/authorization combinations; webhook rejects SMTP fields. Update accepts an entire typed SMTP settings replacement while retaining immutable kind and omitted credentials. Rotation accepts replacement SMTP credentials or explicit clearing and preserves the existing Enroll capability; channel management remains Manage, test Operate, reads Observe. Retargeting existing stored credentials also requires Enroll.

All existing notification channel, task/group policy, effective-policy and delivery endpoints are available over local IPC and the authenticated remote API. Remote request bodies are bounded and mutations are never automatically replayed. Public reads expose typed nonsecret settings and credential presence, never stored credentials.

CLI uses existing notification channel verbs with --kind smtp and smtp-prefixed settings; credential provisioning supports bounded UTF-8 file/stdin input to avoid requiring command-line secrets. Human delivery output identifies the transport and recipient outcomes; JSON stdout remains one response.

Desktop supports a labelled immutable-on-edit transport selector, typed mail settings, password/username replacement or explicit clearing and recipient-level history. Stored credentials are never prefilled. Pending/error states and remote target selection use existing UI conventions.

## Transport

- Required STARTTLS and implicit TLS verify trust and hostname with minimum TLS 1.2; explicit plain relay is supported without silent downgrade.
- One shared attempt context bounds dial, greeting, TLS, AUTH, RCPT and DATA. Cancellation closes the owned raw socket and transport cleanup terminates.
- Use advertised PLAIN/LOGIN for configured credentials; missing/unsupported AUTH is a clear failure, never anonymous fallback.
- RCPT 4xx is temporary, 5xx permanent; valid recipients may proceed. Checkpoint the DATA subset as submitting before transmission. Confirm recipients only after final positive DATA completion.
- Persist retry/terminal outcomes atomically. Retry pending/temporary/uncertain recipients up to the existing three-attempt budget; do not retry accepted/permanent recipients. Interrupted submission carries duplicate-risk wording.
- QUIT failure after DATA confirmation does not undo success. Message-ID remains stable for correlation, not exactly-once guarantees.
- Use authored diagnostics based on phase/category/numeric code. Never return receiver prose, protocol transcripts, wrapped credential errors or authorization values.

## Portable configuration and QA

Optional typed SMTP definitions carry no credentials. Provision exactly one target-local SMTP channel by logical name and required credential presence; existing credentials require the same host/port/encryption mode. ID/revision bindings are resolved and frozen during preview; missing/changed bindings invalidate apply. New definitions integrate with validation/export/compare/drift and existing policy resolution.

Owned SMTP/TLS fixtures and scenarios cover authenticated/relay delivery, Unicode text/escaped HTML, mixed recipient outcomes, temporary/permanent rejection, TLS/authentication failures, cancellation, restart and credential changes. Existing versioned event/report/message-review contracts remain unchanged.
