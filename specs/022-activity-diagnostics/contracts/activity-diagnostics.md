# Contract: Activity Diagnostics

## `GET /v1/logs`

Successful responses retain the bounded `logs` collection and add the exact
configured on-disk path:

```json
{
  "logs": [],
  "log_path": "C:\\Users\\operator\\Schedule Data\\logs\\goschedule.log"
}
```

- `log_path` is present even when `logs` is empty.
- An empty `log_path` means metadata is unavailable; clients must not infer a default.
- Clients must treat the string as display metadata and must not normalize or probe it.
- Existing filters and ordering are unchanged.

## CLI compatibility

`goschedule logs --json` continues returning the record array rather than the
response envelope. Human output continues iterating the same records.

## Activity presentation

Activity states that it shows a limited set of recent daemon records plus
scheduler alerts. It identifies older daemon records as residing in the full
log and displays either the exact path or the unavailable-until-response state.

## Startup record

After initialization and immediately before serving, the daemon emits one
`daemon startup complete` info record with `endpoint`, `db`, and `log_path` attributes.
