# Contract: Agent Access Workspace

## Secret-free workspace

The desktop workspace contains daemon identity, independent transport states, MCP grant projections, and authority descriptions. It may contain safe credential fingerprints but never credentials, pairing phrases, verifiers, digests, certificates, private keys, task inputs, or raw errors.

## Create grant

Input fields are client name, MCP kind, Observe, Operate, or Manage capability, and one approved duration. The service creates a pairing whose enrollment expiry remains ten minutes and whose persistent grant deadline is fixed at creation. Success copies this bundle through the native clipboard boundary:

```text
GO_SCHEDULE_MCP_ENROLLMENT_V1
daemon_id=<uuid>
pairing_id=<uuid>
phrase=<one-time phrase>
capability=<observe|operate|manage>
grant_expires_at=<RFC3339 timestamp or non-expiring>
```

The bundle is never returned to React. If clipboard copy fails, the service cancels the pairing before returning.

## Modify grant

Input identifies one MCP actor and may contain a lower capability or an earlier future expiry. The service rejects built-in, non-MCP, inactive, widening, expiry-removal, and expiry-extension requests before calling the API.

## Revoke grant

Input identifies one active non-built-in MCP actor. Revocation uses the existing actor lifecycle. Authorization reloads current actor state on every applicable request.

## Recent actions

Input identifies one projected MCP actor. Output contains at most 25 shared audit records and only daemon, operation, target, result, and occurrence fields.
