# Data Model: Windows Demo Qualification

## DemoCandidate

- `slice`: fixed value `S043`
- `class`: fixed value `local-demo`
- `source_commit`: full Git commit ID
- `embedded_version`: non-release binary version
- `product_version`: numeric MSI version
- `product_code`: compiled MSI ProductCode
- `filename`: basename containing `s043-demo`
- `byte_size`: exact integer size
- `sha256`: lowercase SHA-256
- `built_at`: RFC 3339 timestamp
- `origin`: local build command and host boundary

Invariant: Any byte change creates a new DemoCandidate identity. DemoCandidate
never transitions into FormalCandidate.

## QualificationCheck

- `name`: stable check identifier
- `command`: exact invocation or inspection action
- `status`: `pass`, `fail`, `unavailable`, or `pending-operator`
- `candidate_sha256`: optional artifact binding
- `evidence`: concise result or attachment path
- `boundary`: what the check cannot prove

Invariant: Only an observed success may be `pass`; absent prerequisites use
`unavailable`, never `pass` or `skipped`.

## AttendedDemoObservation

- `id`: stable checklist identifier
- `candidate_sha256`: exact demo binding
- `environment`: Windows edition/build, display geometry/scaling, and account
  integrity role without personal identifiers
- `expected`: observable result
- `actual`: operator description
- `status`: `pass`, `fail`, `unavailable`, `partial`, or `not-run`
- `attachment`: optional screenshot or log reference

Invariant: A failed observation is recorded before related product code changes.

## FormalCandidate

- Release-workflow repository, tag, commit, run ID, run attempt, artifact name,
  ProductCode, byte size, and SHA-256
- Evidence class `attended-windows-release-candidate`

Invariant: FormalCandidate is created only by authorized tag staging after review.
DemoCandidate evidence cannot be relabeled or copied into it as a pass.

