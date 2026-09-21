# Security Checklist: Remote MCP Authorization

- [ ] Metadata identifies one exact HTTPS resource and authorization server.
- [ ] Client authentication requires matching credential ID and secret.
- [ ] Access tokens are random, opaque, short-lived, resource-bound, and memory-only.
- [ ] Source credential rotation, revocation, actor expiry, and daemon restart revoke access.
- [ ] Pairing phrases, non-MCP credentials, query tokens, cookies, duplicated authorization, and excessive scopes fail closed.
- [ ] Host, origin, TLS, request-size, concurrency, source, and actor limits precede MCP dispatch.
- [ ] Tool discovery and invocation do not exceed issued scopes.
- [ ] Persistent MCP actor identity appears in audit records.
- [ ] Responses and logs exclude secrets, token digests, commands, protected inputs, and raw failures.
- [ ] Official SDK interoperability and supported-version tests pass.
