# Data Model: Authenticated Remote Access and Pairing

## Pairing Session

| Field | Rule |
| --- | --- |
| ID | Opaque UUID, unique |
| Display name | Normalized with the actor display-name rule |
| Actor kind | Desktop, CLI, JSON, or MCP; never local OS |
| Capability | Observe, Operate, Manage, or Enroll |
| Phrase salt | 16 random bytes, never returned after creation |
| Phrase verifier | 32-byte Argon2id output, never returned |
| Attempts remaining | Starts at five; decremented atomically on failed exchange |
| State | Active, consumed, cancelled, expired, or exhausted |
| Created / expires / completed | UTC timestamps; expiry is creation plus ten minutes |

Transitions are `active -> consumed`, `active -> cancelled`, `active -> expired`, or `active -> exhausted`. All terminal states are irreversible. Cleanup may remove terminal sessions after their evidence is no longer operationally useful because management audit remains durable.

## Client Credential

| Field | Rule |
| --- | --- |
| ID | Opaque UUID, unique |
| Actor ID | Unique relationship to one non-local actor |
| Digest | SHA-256 of the raw 32-byte value, unique and indexed |
| Fingerprint | First 12 base64url digest characters for safe identification |
| State | Active or revoked |
| Created / updated / last used / revoked / expires | UTC lifecycle timestamps |

The raw bearer value exists only in issuance memory and the single response. Rotation replaces digest and fingerprint in one transaction. Revocation changes both credential and actor state so either boundary independently denies use.

## Remote Operation

| Field | Rule |
| --- | --- |
| Method and remote pattern | Unique `/api/v1` route |
| Local pattern | Existing `/v1` route, or the isolated enrollment exception |
| Operation ID | Existing S076 identifier or `enrollment.exchange` |
| Capability | Closed S076 capability |
| Audit and target | Existing operation classification |
| Retry class | Safe read, resumable stream, or no replay |
| Body limit | Explicit bytes, zero for requests without bodies |

Runtime records and OpenAPI operation IDs must form the same set.

## Remote Listener Configuration

The configuration contains `enabled`, `bind_address`, `certificate_file`, `private_key_file`, and `acknowledge_public_exposure`. Defaults are disabled with empty values. Enabled configuration validates all fields and certificate material before binding.

## Relationships

- One pairing session can create zero or one actor and zero or one credential.
- One client actor owns at most one active credential in S077.
- One credential always resolves through its actor before authorization.
- One remote operation maps to exactly one application operation.
- Audit events contain actor and operation identity but never credential or pairing identifiers or secret material.
