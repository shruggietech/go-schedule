# Contract: Authorization and Management Audit

## Capability semantics

`observe < operate < manage < enroll`. Every operation has one minimum capability. Unknown operations, actors, kinds, capabilities, and states are denied.

## Local actor resolution

Requests accepted through the protected local transport use the persisted built-in local actor. No header or credential is accepted in S076. Future transports must resolve an authenticated client to an actor before using the same authorizer.

## Actor API

- `GET /v1/access/actors` lists actors and requires Enroll.
- `POST /v1/access/actors` creates a non-built-in actor and requires Enroll.
- `PATCH /v1/access/actors/{id}` changes display name, capability, state, or expiration and requires Enroll.
- `POST /v1/access/actors/{id}/revoke` irrevocably revokes a non-built-in actor and requires Enroll.

Actor responses never contain credentials. Validation failures use the established JSON error envelope. Missing actors return not found, and protected built-in mutations return conflict.

## Audit API

- `GET /v1/audit` lists retained events and requires Enroll.
- `GET /v1/audit/export` exports matching events as `application/x-ndjson` and requires Enroll.

Supported filters are `actor_id`, `operation`, `result`, `since`, `until`, and `limit`. The default limit is 100 and the maximum is 1,000. Results sort by `occurred_at`, then `id`, both ascending. Export uses the same filter and ordering contract.

## Audit lifecycle

Mutations and privileged reads insert an uncertain intent before the handler runs. Successful 2xx and 3xx responses become succeeded; 4xx and 5xx responses become failed. Authorization denials are inserted directly as denied. Audit persistence failure prevents audited handler execution.

## Redaction boundary

The durable contract contains only the fields defined in the data model. Implementations must not add bodies, headers, raw errors, commands, environment values, standard input, paths, keys, credentials, or secrets.
