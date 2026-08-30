# Completion-Chain Management Contract

## Local API

Representation includes `id`, source and target task IDs and names, `on_outcome`, `created_at`, and `updated_at`.

- `GET /v1/chains` returns `{"chains": [...]}`.
- `POST /v1/chains` requires source, target, and outcome; returns 201.
- `GET /v1/chains/{id}` returns one chain.
- `PATCH /v1/chains/{id}` accepts a non-empty subset of mutable fields.
- `DELETE /v1/chains/{id}` returns 204.

Missing resources return `not_found`. Invalid outcome, self-link, duplicate, or cycle return `validation_failed` with a field and actionable message. Validation never partially mutates state.

## CLI

```text
gosched chain create --source <task-id> --target <task-id> --on <success|failure|any>
gosched chain list [--json]
gosched chain show <chain-id> [--json]
gosched chain update <chain-id> [--source <task-id>] [--target <task-id>] [--on <success|failure|any>]
gosched chain rm <chain-id>
```

Text identifies chains and names tasks. JSON emits the API representation. Usage/validation exits 2; runtime/transport errors exit 1.

## Live Event

Event kind `chain` carries `created`, `updated`, or `deleted`, stable chain ID, and the complete chain for create/update. The view model folds it into chain state.

## Desktop

The Chains view shows source name, target name, and outcome. Toolbar actions create, edit, and delete. Forms select name-plus-ID labels and plain-language conditions. Empty state explains how to begin; backend errors leave state intact.

## Run History

Text completion history includes source task/run. JSON/API add optional `source_task_id` and `source_run_id`. Other origins remain compatible with fields absent.
