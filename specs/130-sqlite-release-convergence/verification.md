# S130 Verification

## Candidate and evidence boundary

The evaluated source is the dirty S130 implementation on accepted baseline a58651724bfe4a9582d52464273e9825943771c7, with SQLite1.61.1/libc1.79.0 and complete root/desktop metadata. Original receipts and first failures remain private under .test-temp/s130; all 221 retained inputs are frozen. This is source evidence, not acceptance of a tag, staged package, public release or attended native behavior.

## Test-first failures and selected-pair results

The first commit-hook fixture was invalid because SQLite invokes the hook after obtaining the blocking lock. Its diagnostic SQLITE_BUSY result and exact fixture are retained, not represented as proof of cancellation. The corrected final fixture finishes both read/write transaction setup, retains the blocking reader through synchronous Commit and arms a 100ms cancellation timer with explicit callback completion. On the old pair it returned SQLITE_BUSY after 5.0242866 seconds, losing the context error; the first failure is preserved. On the selected pair all focused cancellation and memory/WAL prepared-binding regressions passed five times under the race detector in 3.639 seconds.

The test establishes a real held-lock deadline contrast, not instrumentation of the busy-handler entry. It also verifies rollback leaves no partial transaction and the same connection independently commits its next write. Production one-connection/WAL policy and pragmas are unchanged. Existing current initialization/refusal/reopen/audit/security suites remain required by the aggregate below.

After root selection, normal desktop -mod=readonly dependency resolution failed with updates to go.mod needed, disabled by -mod=readonly. After regenerating desktop metadata, root and desktop read-only graphs, both checksum verifications and desktop connection race tests passed (3.373 seconds). All four manifest/checksum hashes remained unchanged during read-only acceptance. This resolves the demonstrated combined-graph defect rather than hiding metadata writes during builds.

## Performance

Initial unpaired samples varied substantially and appeared to exceed 10 percent in prepared workloads. Those raw receipts are retained. Seven paired old/new runs, alternating order on the same host/toolchain with identical workloads, did not reproduce that apparent timing regression. Median results:

| Workload | Old ns/op | Candidate ns/op | Change | Old/candidate B/op | Old/candidate allocations/op |
| --- | --- | --- | --- | --- | --- |
| BenchmarkFreshOpen-2 | 1402247 | 1403866 | 0.12% | 9062/9410 | 268/270 |
| BenchmarkSQLitePreparedBindings/memory-2 | 17729 | 16811 | -5.18% | 2528/2544 | 143/143 |
| BenchmarkSQLitePreparedBindings/wal-2 | 20045 | 18858 | -5.92% | 2624/2640 | 145/145 |
| BenchmarkTriggerSetMaximumLifecycle-2 | 4764870 | 4802571 | 0.79% | 461196/479192 | 9570/9578 |

All four median time changes satisfy the existing 10 percent regression policy. Prepared paths have unchanged allocation counts and a 16-byte increase; fresh-open and maximum-trigger lifecycle allocation/byte changes remain reported above. These are source benchmark measurements on this host, not end-user latency or tagged-package acceptance.

## Canonical verification

The complete foreground sh scripts/verify.sh all ran from 2026-10-31T11:15:57.028Z to 2026-10-31T11:40:21.717Z, exited 0 and passed all eight gates in 1464.688 seconds. All 2545 source files/paths and all 221 retained inputs remained unchanged during the run. Evaluated source manifest SHA256: 5fd4b6aa998ffbe864265cce125270ec5c6826fc070183391490089b6e0bc80a. Original manifests, logs and terminal receipt remain private.

| Gate | Terminal result |
| --- | --- |
| format | PASS, no Go or GitHub prose defects |
| vet | PASS |
| lint | PASS, zero issues |
| race | PASS, root daemon/CLI/security/current-state/integration suites; store 41.750s and integration 268.671s |
| gui | PASS, all desktop race packages, native Wails2.17.0 build, 167 frontend tests and production bundle |
| coverage | PASS, engine83.6%, schedule89.1%, timezone91.3%, store80.9%, catchup88.9%, logbus91.1%; no decrease |
| docs | PASS, 24-document export and existing link/metadata/policy checks |
| automation | PASS, all positive/negative automation fixtures and eight-gate contract |

After this quiescent acceptance, only factual specification/lifecycle records and the stale current desktop dependency guide are finalized. The guide now uses the same Wails2.17.0 command that passed the native gate and points to authoritative module/frontend manifests instead of copying version lists. Final documentation, source/comparison, formatting/lifecycle and integrity checks follow these edits; no evaluated Go, dependency, workflow, notes or comparison bytes change.

Top-level license texts for all six changed declared dependency versions are byte-identical to the prior versions. This license check does not equate module-only tooling with redistributed code or certify upstream provenance.

## History and release continuation

Candidate notes/changelog and content-bound review facts are updated; the v2.0.0 comparison dataset and historical releases stay immutable. Existing S129 source acceptance retains its original source identity. Continue actual #412 public delivery through the existing S129 handoff after owner integration and distinct exact tag/draft and concrete-asset promotion approvals. Public latest stays v1.7.0 now. #355/#356 causes remain unconfirmed; clustered execution is owner-deferred. No migrations or historical upgrade matrix.

## Final documentation and factual finalization

After correcting desktop/README.md, the full documentation policy/architecture/link checks, seven site tests and 24-document production export passed. Source/comparison guards, all 128 specification lifecycle checks, UTF-8/no-BOM/punctuation and history/input integrity checks passed. All 20 historical changelog sections, 17 immutable release/comparison/public-selection records and 221 retained inputs have zero drift. Only factual specification/lifecycle records and the current dependency guide changed after the aggregate; a dedicated finalization audit rejects any other source change.

This is the local pre-publication acceptance record. The official PR carries its actual committed head, hosted CI results, every bot review and the final owner handoff. No tag, draft, promotion or recurring monitor was created.
