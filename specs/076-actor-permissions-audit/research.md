# Research: Actor Permissions and Management Audit

## Decisions

### Closed capability hierarchy

Use the architecture-approved Observe, Operate, Manage, and Enroll ordering. An actor satisfies a requirement when its level is equal or stronger. Parsing rejects any unknown spelling so future transports cannot invent authority outside the shared model.

### Built-in local actor

Persist exactly one protected local operating-system actor with Enroll authority. The existing IPC ownership boundary authenticates local use, and all current local requests resolve to this actor without credentials. Actor records remain credential-independent for issues #168 and #169.

### Shared operation catalog

Keep route registration and authorization metadata in one server catalog. Each operation declares method, path pattern, identifier, required capability, target kind, and audit class. Tests compare registered routes against catalog entries, while the authorizer independently denies unknown operation identifiers.

### Intent-first audit

Insert an `uncertain` audit event before an audited handler runs. If insertion fails, return service unavailable and do not execute the handler. Afterward update that event to succeeded or failed. A process interruption intentionally leaves uncertain evidence.

### Redaction and retention

Store only identifiers, classifications, results, correlation identifiers, and timestamps. Do not store payloads, headers, raw errors, commands, environment values, standard input, paths, keys, credentials, or secrets. On each intent insert, transactionally delete records older than 90 days and retain only the newest 10,000.

### Actor administration and export

Expose actor list, create, update, and revoke plus audit list and newline-delimited JSON export through the protected local API, typed client, and CLI. Stable ordering uses occurrence time followed by event identifier. Filters accept actor, operation, result, time range, and bounded limit.

## Alternatives rejected

- Per-handler authorization was rejected because missing checks would fail open and future transports could drift.
- Audit-after-execution was rejected because crashes could erase evidence of attempted mutations.
- Persisting client credentials with actors was rejected because credential lifecycle belongs to #169.
- Unbounded audit history was rejected because it creates unmanaged disk growth.
- Replacing IPC protection with a local login was rejected because it adds friction without improving the established local trust boundary.
