# S110 Research and decisions

## Release identity (#262)

Decision: build a small standard-library Go helper from explicit github.sha workflow checkout into RUNNER_TEMP before checking out the requested tag. Four commands and one temporary state document bind later operations. Keep candidate/staging/comparison/immutable-tag and shell checksum validation in the workflow.

Rationale: GitHub documents tag lookup as a published-release lookup; authenticated listing includes accessible drafts. A helper added on current main cannot be expected in an old tagged tree. Download/upload by tag would repeat the original ambiguity even after fixing initial lookup. Follow all pages, select one positive ID, re-fetch and validate it, then use ID-bound asset and mutation endpoints. Existing published releases allow latest-only recovery without content edits.

Alternatives: replacing only initial lookup leaves downloads/uploads defective; inline Bash complicates meaningful cross-platform failure injection; a broad release framework exceeds scope. Build the narrow helper once from the trusted workflow revision and test it with httptest.

Sources: [GitHub releases API](https://docs.github.com/en/rest/releases/releases), [release assets API](https://docs.github.com/en/rest/releases/assets), [pagination](https://docs.github.com/en/rest/using-the-rest-api/using-pagination-in-the-rest-api), issue #262. Binary downloads can redirect; never forward authentication to an unrelated origin. API errors retain useful status/context without credentials. A lost upload response is not grounds to repeat POST blindly; resume by listing and comparing bytes.

## Observer recovery (#271)

Decision: recovery retains state only for enabled definitions with identical FilesystemWatcherRevision. Explicit Reload retains generation-reset behavior. Keep bounded last-dispatch identity for recovery suppression and suspend inactive candidate timers rather than dropping work or spinning.

Demonstration: unchanged production code fails a direct runtimeState/fake-clock Go overlay test: active -> degraded -> active, replacement dispatches=0, expected1. degradeAll/root degradation clear fileStates/pending/dispatched; rebuild/recovery unconditionally baselines the replacement. This proves a production loss mechanism independently of the macOS observer's original failure cause. Reproducer: .test-temp/s110-watcher-research/repro_test.go and result.log.

Alternatives: preserving only fileStates loses pending debounce/stability work; preserving candidates without timer guards deletes them or spins while inactive; deduplicating all native hints permanently changes ordinary write behavior. Preserve revision-scoped state and apply suppression specifically to recovery/reconciliation provenance.

## Heartbeat record loss (#335)

Decision: encode non-null values as UTF-8 hexadecimal inside a CLI-quoted text-cast expression in both script twins; retain explicit null handling, numeric casts, success/row-count assertions, diagnostic preservation and bounded retry behavior.

Demonstration: distinct valid 32-character IDs (1e followed by 30 nines, and 2e followed by 30 nines) both persist as text Inf through current .param quoting. The second heartbeat exits1 with UNIQUE constraint failed: beat.session_id, beat.sequence (19); only one row remains. Evidence: .test-temp/s110-heartbeat-research/fault-sessions.py and fault-sessions.log. This is a demonstrated data-loss mechanism, not proof of the earlier hosted run's exact cause. Two hundred additional unchanged local concurrency rounds passed and do not prove a correction.

Alternatives: merely increasing waits is unsupported for this non-contention failure; more SQL quote escaping does not prevent SQLite CLI expression evaluation; generating different session IDs hides corrupt general text binding. Correct the shared boundary and exercise arbitrary literal values plus distinct identities concurrently.

## Operational decisions

Research agents were dispatched by installed speckit-plan. Process launches use hidden native children with redirected non-interactive I/O. The owner reported focus loss during research, then directed that it need not be investigated further; all diagnostic sessions had ended. Continue the original S110 scope without claiming that transient focus cause was established.
