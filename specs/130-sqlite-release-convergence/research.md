# S130 Research

## Published correction and licensing

Decision: select modernc.org/sqlite v1.61.1 and the exact libc v1.79.0 declared by its published module. SQLite info records 2026-10-31T10:08:39Z at 014351f110f66c8f485cff948fea608d7dcfafb1; libc info records 2026-10-15T21:00:57Z at 02bd2316433516d22b0cc2946cbd2ebc08fbd44f. Both info/module endpoints return HTTP 200. Go 1.26.0 remains compatible. Sources: [SQLite info](https://proxy.golang.org/modernc.org/sqlite/@v/v1.61.1.info), [module](https://proxy.golang.org/modernc.org/sqlite/@v/v1.61.1.mod), [libc info](https://proxy.golang.org/modernc.org/libc/@v/v1.79.0.info).

Rationale: the [tagged changelog](https://gitlab.com/cznic/sqlite/-/blob/v1.61.1/CHANGELOG.md) reports the prepared-query reprepare correction, commit/rows lock-wait cancellation and VFS error corrections. These are upstream claims, not product acceptance. No VFS adoption, OFD locking, new authorizer or production pooling change is required.

Licenses: SQLite driver and libc retain BSD-3-Clause with unchanged top-level texts relative to the held pair. Upstream [third-party inventory](https://gitlab.com/cznic/sqlite/-/blob/v1.61.1/LICENSE-3RD-PARTY.md) distinguishes linked components from module-only tooling; preserve existing redistribution obligations. SQLite is public domain and inherited musl material is MIT. Do not label every go.sum entry as shipped code or claim verified provenance from the changelog.

Alternatives: v1.61.0 remains the known-regression proposal; unpublished forks are unnecessary; indefinite hold no longer has the availability rationale.

## Current storage and meaningful regressions

Decision: retain production Open's one connection, foreign keys, transactional current schema and WAL. Existing context_test.go and request_cancellation_test.go cover pool waits and actual IPC, not SQLite busy-handler waits. Existing bootstrap/refusal/reopen/audit tests cover current durable behavior and stay intact.

Add an owned driver-boundary lock test (rollback-journal commit and WAL execution as warranted), with reader/writer ownership and cancellation armed only after SQL setup. The first fixture established that the commit hook runs after lock acquisition, so it is not a usable wait-entry signal. Use the upstream real-I/O deadline pattern, retain the old-driver five-second SQLITE_BUSY failure, and join the cancellation timer callback. Deadline/elapsed watchdogs must be far shorter than configured busy timeout and generous for CI; no sleep-driven ordering. Explicitly release connections/transactions/hooks and observe result before releasing the blocking lock. Require context error and no durable partial mutation.

Prepared correctness/benchmark should include a parameterized join/ordered query with alternating integer bindings 0/1 on memory and disk/WAL. Before/after runs use the identical workload, toolchain and host. Existing current application state tests remain the product acceptance boundary. Alternatives rejected: pool-only cancellation duplicates existing coverage; VFS/cluster frameworks and migration campaigns are unrelated.

## Module and release records

Decision: regenerate all four root/desktop metadata files. Preserve /v2, desktop /v2/desktop, local replace and required root v2.0.0. Existing desktop/remote read-only builds expose the original graph defect, so no new generalized dependency guard is needed.

Update current v2.0.0 notes, the prepared 2.0.0 changelog decision and its review changelog/notes digests. Keep site/data/comparisons/v2.0.0.json unchanged, preserve historical notes/reviews, and keep published.json selecting v1.7.0. S129 verification remains historical source acceptance; its active handoff receives a factual dependency supersession link to S130 only.

Existing #412 and milestone 19 track actual delivery. #404 is linked as the current dependency child and closes only on functional source integration. Proposal #400 is superseded only after replacement integration. Source PR authorization is explicit; no release mutation or recurring monitor is inferred.
