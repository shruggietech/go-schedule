# Data Model: MCP Manage Authority

## Manage session policy

- Capability: `manage`
- Require confirmation: runtime boolean selected when the connection starts
- Lifetime and revocation: inherited from the runtime MCP session

## Manage request envelope

- Daemon ID: exact installation UUID
- Request ID: caller-generated UUID
- Action: create, update, or delete where supported
- Object ID: required for update and delete, absent for create
- Confirmed: required only when connection policy demands it
- Definition: one family-specific bounded input

## Manage result

- Schema version
- Permission
- Operation
- Daemon ID
- Request ID
- Object kind
- Object ID
- Outcome: accepted, rejected, denied, or uncertain
- Safe message

The result never includes definition contents or secrets.

## Deduplication record

- Request ID
- Exact mutation fingerprint
- First result
- Expiry

The cache holds at most 256 records for 10 minutes per runtime session.
