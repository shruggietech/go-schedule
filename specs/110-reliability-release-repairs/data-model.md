# S110 Data model

## Release selection

Versioned temporary JSON binds repository, requested tag, positive release ID, original content metadata fingerprint (tag, name, body, target and prerelease), observed publication state, original asset IDs/names/sizes/upload state and verified SHA-256 digests, plus accepted checksum identity. Exclude credentials, signed download URLs and response headers. Write atomically. Updated timestamps do not enter content fingerprints because legitimate upload can change them.

Transitions: unresolved -> one authenticated match -> verified initial download -> accepted matching checksum -> verified final download -> published/latest. Missing/ambiguous/changed identity refuses without mutation. Published recovery requires complete existing checksum and contents and permits only latest-pointer repair. Interrupted upload resumes by observing existing bytes.

## Watcher recovery state

Definition revision identifies the state eligible for retention. Last observed file identity, pending debounce/stability candidate and bounded last-dispatched identity survive transient observer loss for that revision. Inactive candidates are suspended until registration recovers. Changed/disabled/deleted revisions and explicit Reload invalidate obsolete generation work. Suppression provenance distinguishes recovery/reconciled dispatch from ordinary later native hints.

## SQLite payload values

Explicit null remains null. Every other parameter is exact UTF-8 text expressed through a hex-literal text cast; numeric conversion stays at existing SQL CAST sites. Session IDs remain exact strings with one record per session/sequence. No schema migration or persisted product state change.
