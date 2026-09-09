# Remote Security Requirements Checklist

**Purpose**: Review the completeness and precision of S077 security requirements before implementation
**Created**: 2026-09-09
**Feature**: [spec.md](../spec.md)

## Boundary and Exposure

- [x] CHK001 Are default-disabled, pre-bind validation, wildcard acknowledgement, TLS minimum, timeout, and shutdown requirements explicit? [Completeness, Spec FR-001 to FR-004]
- [x] CHK002 Are remote allowlist and local-only exclusions stated independently of local route registration? [Consistency, Spec FR-005 to FR-007]
- [x] CHK003 Are hostile browser-origin and forwarded-metadata boundaries addressed? [Coverage, Spec FR-022 and Edge Cases]

## Authentication and Authorization

- [x] CHK004 Is the only accepted credential transport unambiguous and are alternate transports rejected? [Clarity, Spec FR-008]
- [x] CHK005 Are entropy, digest-only storage, constant-time comparison, actor reload, and generic failure requirements complete? [Completeness, Spec FR-009 to FR-012]
- [x] CHK006 Are source and actor rate limits required with bounded state and non-enumerating errors? [Coverage, Spec FR-013 and SC-008]

## Enrollment and Lifecycle

- [x] CHK007 Are phrase entropy, lifetime, attempt budget, verifier, and single-use semantics measurable? [Measurability, Spec FR-014 to FR-016]
- [x] CHK008 Are concurrency, replay, wrong-daemon, cancellation, expiry, and exhaustion outcomes specified? [Edge Cases, Spec FR-017 and FR-018]
- [x] CHK009 Are rotation, revocation, immediate-request denial, and stream termination requirements consistent? [Consistency, Spec FR-019 and FR-020]

## Secret Safety and Contract Integrity

- [x] CHK010 Is the protected-value inventory defined across every output and evidence surface? [Completeness, Spec FR-023]
- [x] CHK011 Is complete machine-contract traceability required for every reachable remote route? [Traceability, Spec FR-021 and SC-003]
- [x] CHK012 Are unsupported security and deployment mechanisms explicitly excluded to prevent scope drift? [Clarity, Spec FR-025]

## Notes

- Reviewer-depth checklist covering the S074 threat classes owned by issues #168 and #169.
