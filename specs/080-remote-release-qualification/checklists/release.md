# Release Qualification Requirements Checklist

**Purpose**: Test whether the S080 requirements are complete enough to govern a secure v1.4 source qualification

**Created**: 2026-09-10

**Feature**: [spec.md](../spec.md)

## Lifecycle Coverage

- [x] Are clean installation, explicit enablement, pairing, use, revocation, disablement, restart, and upgrade requirements all stated?
- [x] Is default-off behavior required independently from an unreachable endpoint?
- [x] Is registered-service configuration persistence defined for every platform?
- [x] Are local IPC continuity and retained state required throughout the journey?

## Security and Deployment Coverage

- [x] Are recommended and advanced deployment modes visibly distinguished?
- [x] Are application TLS and operator-owned certificate duties required in every network mode?
- [x] Are permission denial, attribution, rate limits, identity changes, version skew, network loss, and certificate changes represented by the qualification scope or its dependencies?
- [x] Are secret exclusions explicit for evidence, logs, diagnostics, and documentation?

## Evidence and Publication Boundary

- [x] Are all three supported operating systems named?
- [x] Is every issue #173 criterion required to map to objective evidence?
- [x] Is source qualification distinguished from an immutable public release artifact?
- [x] Are tags, releases, and public artifacts explicitly excluded?
