# Remote Security Architecture Checklist

**Purpose**: Validate the bounded security, deployment, compatibility, and ownership contract before implementation
**Created**: 2026-09-09
**Feature**: [spec.md](../spec.md)

## Boundary Completeness

- [x] CHK001 Is one primary remote transport named and are alternative wire protocols rejected? [Completeness, Spec FR-001]
- [x] CHK002 Are local IPC and remote HTTPS separate adapters over shared daemon operations? [Clarity, Spec FR-002 and FR-003]
- [x] CHK003 Does the ordered request boundary place authentication and authorization before domain behavior? [Coverage, Contract Boundary order]
- [x] CHK004 Are local-only and secret-disclosure routes denied by omission from the remote allowlist? [Security, Contract Remote allowlist record]

## Identity and Authority

- [x] CHK005 Are Observe, Operate, Manage, and Enroll sufficient and explicitly bounded? [Clarity, Spec FR-004]
- [x] CHK006 Are actor authority and credential bytes stored separately? [Security, Spec FR-005 and FR-006]
- [x] CHK007 Are pairing phrases isolated from durable bearer use, replay, and recovery? [Coverage, Spec FR-006]
- [x] CHK008 Do revocation and expiry affect new requests and ongoing streams? [Coverage, Spec Edge Cases]

## Transport and Deployment

- [x] CHK009 Do all network modes retain HTTPS without a plaintext exception? [Consistency, Spec FR-007 and FR-008]
- [x] CHK010 Does every deployment mode identify product and operator ownership? [Completeness, Spec SC-002]
- [x] CHK011 Are invalid, wildcard, public, proxy, and certificate configurations fail-closed? [Security, Spec Edge Cases]
- [x] CHK012 Are direct public exposure and certificate automation claims explicitly bounded? [Scope, Spec FR-015]

## Abuse and Failure

- [x] CHK013 Are pre-authentication source and post-authentication actor limits both required? [Coverage, Spec FR-009]
- [x] CHK014 Do authentication failures avoid credential or enrollment enumeration? [Security, Contract Failure contract]
- [x] CHK015 Are headers, bodies, collections, streams, timeouts, and limiter state bounded? [Completeness, Spec FR-005 and FR-013]
- [x] CHK016 Are uncertain mutations never automatically replayed? [Safety, Spec FR-012]

## Compatibility and Maintenance

- [x] CHK017 Is OpenAPI source ownership distinct from generated server and client code? [Clarity, Spec FR-010]
- [x] CHK018 Are compatible and breaking changes separated by explicit URI-major rules? [Coverage, Spec FR-011]
- [x] CHK019 Does every selected dependency have an introduction owner and replacement rule? [Maintenance, Spec FR-014]
- [x] CHK020 Are then-current version, license, advisory, and platform checks deferred to the slice that introduces each dependency? [Scope, Research Dependency ownership policy]

## Traceability

- [x] CHK021 Does the threat table map every named risk to a control and an acceptance-test class? [Traceability, Spec FR-013]
- [x] CHK022 Does the executable guard cover transport, modes, owners, threats, non-goals, and downstream order? [Measurability, Spec FR-016]
- [x] CHK023 Is #166 excluded from S074 so architecture review precedes implementation? [Dependency, Spec Clarifications]
- [x] CHK024 Do #165 and all downstream issues remain open until their own completion evidence exists? [Lifecycle, Spec FR-018]

## Notes

- All security-architecture questions are answered by the specification, research decisions, data model, or remote-boundary contract.
- This checklist validates requirements quality. Runtime penetration, certificate, proxy, revocation, and platform tests remain mandatory in their owning implementation and release slices.
