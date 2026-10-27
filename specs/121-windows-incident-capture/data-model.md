# S121 incident data model

Incident v1 binds kind, unique ID, category, phase, target PID/executable hash, original exit/inspection fields, collection status/detail/duration and bounded artifact roles/lengths/hashes. Missing, truncated, oversized, denied and timed-out collection are explicit values. Collection cannot set the original command exit or snapshot completeness.

Daemon request binds nonce, incident ID and expected PID. Response echoes those identities and records active authored route/stage categories, omitted count, pool connection/wait counters and a bounded local goroutine artifact. Pending -> collected/unavailable/failed/timed_out is independent of command execution. No request body, credential or arbitrary route identity is accepted into structured fields.

CaptureOptions, captureIdentity and WorkerConfig carry an explicit incidents boolean. Resume retains the original choice. Existing strict report v1 remains unchanged; case State carries incident references/status. Raw native/stack files remain private local evidence.
