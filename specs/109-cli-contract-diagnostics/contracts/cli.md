# S109 CLI and response contracts

## Invocation

Success is 0, usage/validation refusal is 2, operational failure is 1. Extra operands fail before target resolution or host adapters. Existing explicit variadic/custom operand validators keep their accepted forms. Built-in help/completion remains usable without remote authentication.

## Output

Structured results occupy stdout as one document. Alerts migration notices occupy stderr. Existing human list columns gain explicit headings. Help describes effects, identities and actual service configuration options. Lookup/validation/document errors include relevant safe context and a corrective action without losing underlying error category or parser detail.

## Targets and plans

Use --daemon-endpoint for an explicit daemon when a command also owns a destination --endpoint. Legacy global --endpoint remains accepted where unambiguous. Both remote aliases together are refused. Explicit IPC conflicts only with actual global remote selectors. Bundle bindings reject unknown/duplicate identities before a usable plan exists. Invalid cron lines retain their complete report with exit 2; operational import failures remain exit 1.

## Credential and watcher responses

Credential rotation keeps its response type, associated identities and one-time token delivery. Obtain fallible association context before credential mutation. Watcher replies retain the actual last observation and separately identify pending reconciliation against saved configuration. Observation metadata advances even when health state/reason remains unchanged. No waiting or new operational permission is required.

## Compatibility warning

Explain compatibility-mode access and identify the daemon configuration's top-level admin_group field. An explicitly empty group enables broad local access; configuring an OS administrative group and restarting applies restricted access. Retain structured access_mode/admin_group log fields and current authorization behavior. access_mode is diagnostic metadata, not a configuration field.
