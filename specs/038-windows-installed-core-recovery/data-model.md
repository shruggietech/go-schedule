# Data Model: Windows Installed Core Recovery

## Authorized pipe principal

An in-memory principal used only while constructing the Windows pipe descriptor.

| Field | Type | Rules |
| --- | --- | --- |
| SID | string | Required, valid Windows SID string, case-insensitive identity |
| Source | enum | `system`, `administrators`, `configured_group`, or `direct_user_member` |
| Access | enum | `full` for service owners; `read_write` for operators |

Principals are de-duplicated by SID and direct user members are sorted before rendering. Principal identities are not logged.

## Direct group member record

Transient projection of `LOCALGROUP_MEMBERS_INFO_1`.

| Field | Type | Rules |
| --- | --- | --- |
| SID | string | Validated before use; empty only for unusable deleted/unknown records |
| SID usage | integer enum | Only `SidTypeUser` is expanded into a user ACE |

Group, alias, and well-known group records remain authorized through the configured group ACE. They are not recursively expanded.

## Execution probe task

Existing `domain.Task` persisted through the public product boundary.

| Field | Probe contract |
| --- | --- |
| Name | Unique evidence-scoped name |
| Command | Absolute inbox Windows PowerShell executable |
| Args | Explicit noninteractive switches plus output and ProgramData .NET marker operation |
| WorkingDir | Empty or explicit evidence directory, recorded in evidence |
| Env | Keys recorded; values redacted |
| Schedule | Near-term recurring schedule for scheduled proof |
| RunAs | Empty, preserving LocalSystem execution |

## Execution evidence record

Markdown evidence assembled from existing task/run JSON, filesystem effects, service metadata, and daemon logs.

| Field | Required evidence |
| --- | --- |
| Candidate | MSI path/origin, SHA-256, product version/code |
| Host | Windows caption/version, timestamp, elevation state |
| Authorization | installing identity, user/group SIDs, direct membership, token SIDs, pipe descriptor, connection result |
| Service | name, state, start mode, account, executable path |
| Task | ID, command, args, working directory, environment key names |
| Manual run | trigger, outcome, exit code, output, service-writable marker content/hash |
| Scheduled run | trigger, outcome, exit code, output, retained marker evidence |
| Failure controls | nonzero exit record and process-start failure record |
| Logs | bounded correlated daemon-log excerpt with secrets redacted |

No new persistent application entities or schema migration are introduced.
