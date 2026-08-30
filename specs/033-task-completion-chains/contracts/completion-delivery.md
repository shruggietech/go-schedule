# Completion Delivery Contract

## Matching

- Source `success` matches `success` and `any`.
- Source `failure` matches `failure` and `any`.
- Other outcomes create no deliveries.
- Each chain/source-run pair is inserted at most once.

## Processing

- Startup returns every claimed delivery to pending and records replay evidence.
- A claim atomically changes pending to claimed and increments attempts.
- Current target eligibility is checked after claim.
- Eligible targets enter existing overlap handling with completion correlation.
- Missing/inactive/disabled targets or removed relationships resolve without execution.
- Queue-one preserves correlation on the actual pending execution; its informational queued record does not close the delivery.
- Skip overlap resolves the delivery with correlated skipped history and creates no downstream deliveries.
- Actual success/failure closes the incoming delivery and creates matching downstream deliveries atomically.

## Recovery

Pending work is durable. Completed/resolved work never returns. Claimed work at unclean shutdown returns to pending and may replay because external process state is unknowable. Attempts and diagnostics expose recovery.

## Bounds

Claims use finite indexed batches. Actual execution remains bounded by the engine worker pool. There is no polling loop.
