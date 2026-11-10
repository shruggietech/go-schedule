# S143 data model

Active execution remains keyed by Run ID, with one cancellation cause/handle and a finishing result retained until persistence. Per-task running state becomes an in-flight count including executions waiting for a worker.

Run adds transient cancellation_status (requested or failed) and cancellation_error observation fields. Existing unrestricted outcome storage adds canceled and cancellation_failed terminal values. Native exit, output caps and correlation fields remain unchanged; no task or database shape migration is required.

RunCancellation acknowledgement contains run_id, task_id, status (requested or already_finished), optional observed outcome and transient cancellation status/error. It does not include execution inputs or raw output. Unknown ID is missing; queued/skipped history is not an active executable handle.

State transitions: active -> cancellation requested -> observed canceled, or active -> requested -> failed/unconfirmed active -> eventual cancellation_failed after actual reap. Natural result already observed -> finishing/persisted -> immutable already_finished. A five-second expiry does not manufacture native exit or end time.
