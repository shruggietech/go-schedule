# Implementation Plan: Authenticated Remote Access and Pairing

**Branch**: `codex/077-remote-https-pairing` | **Date**: 2026-09-09 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/077-remote-https-pairing/spec.md`

## Summary

Deliver issues #168 and #169 as one usable trust-establishment boundary. The daemon keeps its protected local API and adds a separately constructed, explicitly configured TLS 1.3 listener. A remote allowlist maps `/api/v1` routes to approved local operations, bearer authentication resolves current credential and actor state on every request, and isolated phrase enrollment atomically creates one actor and opaque credential. Local Enroll-authorized API and CLI commands administer pairing sessions and credential lifecycle. OpenAPI, contract checks, documentation, and adversarial tests keep the boundary fail-closed.

## Technical Context

**Language/Version**: Go 1.25.0

**Primary Dependencies**: Go standard library HTTP/TLS/crypto, `golang.org/x/time/rate` v0.15.0, `golang.org/x/crypto/argon2` v0.55.0, `github.com/oapi-codegen/oapi-codegen/v2` v2.8.0 as a pinned generation command, and `github.com/zalando/go-keyring` v0.2.8 behind a small client credential-store adapter

**Storage**: Existing SQLite store with forward-only schema v18 for pairing sessions and credential verifiers; raw phrases and bearer values are never persisted

**Testing**: Standard Go tests, race detector, HTTPS integration tests with ephemeral certificates, migration tests from schema v17, deterministic limiter tests, OpenAPI/allowlist drift checks, and the canonical eight-gate verification suite

**Target Platform**: Linux, macOS, and Windows daemon and client tooling

**Project Type**: Single daemon, CLI, typed client, and Wails desktop repository

**Performance Goals**: Routine authenticated requests add less than 25 ms p95 locally; phrase verification remains between 100 ms and 750 ms on the supported development baseline; event-stream authority is revalidated at least every 30 seconds

**Constraints**: No listener by default, TLS 1.3 only, no raw-secret persistence, no browser CORS, no custom cryptographic protocol, no automatic mutation replay, no local IPC regression, and no unbounded rate-limiter or connection state

**Scale/Scope**: Individually administered daemons with at most 1,000 durable credentials, 64 concurrent remote connections, 256 KiB request bodies, and the reviewed initial remote operation subset

## Constitution Check

*GATE: Passed before research and passed again after design.*

- **Performance**: The remote server has explicit connection, body, concurrency, and limiter bounds. Authentication uses one indexed digest lookup; Argon2id is isolated to short-lived enrollment.
- **Testing**: Tests precede implementation for migration, exposure defaults, TLS, authentication, capability denial, phrase lifecycle, concurrency, rate limits, secret exclusion, and contract completeness. Full race and coverage gates remain mandatory.
- **UX consistency**: The local API remains unchanged. Remote errors retain stable JSON envelopes, and administrative CLI commands follow existing output and validation conventions.
- **Simplicity**: The implementation reuses the local operation handlers through an explicit adapter, standard HTTP/TLS/crypto, four existing capabilities, opaque server-owned credentials, and two narrow dependencies approved by S074. No accounts, JWTs, PKI ownership, policy language, or browser client are introduced.
- **Autopilot and integration**: The full spec-kit chain, analysis, implementation, canonical verification, changelog, review branch, official PR, and review handling remain required. The user supplied publication authorization in the S077 kickoff.

## Project Structure

### Documentation

```text
specs/077-remote-https-pairing/
├── checklists/
│   ├── remote-security.md
│   └── requirements.md
├── contracts/
│   └── remote-v1.md
├── data-model.md
├── plan.md
├── quickstart.md
├── research.md
├── spec.md
├── tasks.md
└── verification.md
```

### Source Code

```text
api/openapi/
└── remote-v1.yaml
cmd/goschedd/
└── main.go
internal/
├── api/server/
├── authorization/
├── cli/
├── clientsecret/
├── config/
├── domain/
├── remote/
└── store/
```

**Structure Decision**: Keep all product code in the existing root Go module and Wails frontend. `internal/remote` owns only the network adapter, enrollment endpoint, authentication, limits, and TLS lifecycle. Existing `internal/api/server` handlers remain the single application behavior. `internal/clientsecret` isolates native keyring use, and the existing Connections workspace receives only the pairing form and handoff needed by #169 without creating the profile system owned by #170 and #171.

## Complexity Tracking

No constitution violation is accepted. Three direct dependencies are justified by the already reviewed S074 architecture: Argon2id avoids inventing password hashing, `x/time/rate` avoids inventing concurrent token buckets, and `go-keyring` avoids three platform credential-store implementations. OpenAPI generation is a pinned development command, not a runtime framework.
