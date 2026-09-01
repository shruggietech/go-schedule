# Data Model: IPC Access-Denied Recovery

## ConnectionFailure

- `Kind`: unavailable, access denied, timeout, or other transport
- `Operation`: request method/path or event-stream operation
- `Cause`: wrapped original error
- Validation: API response errors are never represented by this entity

## ConnectionIncident

- `Failure`: latest classified connection failure
- `Title`: stable category label
- `Guidance`: actionable, evidence-conditioned user copy
- `Retrying`: whether an immediate recovery attempt is active
- `Revision`: monotonically increasing update number for deterministic tests
- Lifecycle: absent -> active -> updated/retrying -> cleared on success
- Identity: one process-wide incident; category changes update the same entity

## AccessDiagnosis

- `Service`: running, stopped, missing, or unknown
- `GroupExists`: yes, no, or unknown
- `AccountMember`: yes, no, or unknown
- `TokenMember`: yes, no, or unknown
- Derived stale-token state: group exists AND account member AND token not member

## ReconnectPolicy

- `InitialDelay`: 2 seconds
- `MaximumDelay`: 30 seconds
- `Attempt`: reset after successful connectivity
- Transition: 2s -> 4s -> 8s -> 16s -> 30s -> 30s
- Waits are context-cancelable and can be interrupted by Retry
