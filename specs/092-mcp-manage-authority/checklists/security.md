# Security Checklist: MCP Manage Authority

- [x] Observe and Operate discovery boundaries are explicit.
- [x] Manage cannot modify actors, grants, sessions, or enrollment state.
- [x] Exact daemon, request, and object identities are required.
- [x] Secret-bearing task fields and credential reveal are excluded.
- [x] Trigger creation output is redacted.
- [x] Optional confirmation policy fails closed before dispatch.
- [x] Retries are deduplicated and differing request reuse is rejected.
- [x] Revocation is evaluated on subsequent requests.
- [x] Audit attribution uses the runtime MCP actor.
- [x] Adversarial and conformance tests are required.
