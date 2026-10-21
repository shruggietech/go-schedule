# S110 Validation guide

Run from the review branch with Go, native C compiler, Node/npm and the existing POSIX shell available. Process launches must be hidden and non-interactive on Windows. Use disposable fixtures only.

1. Run go test -race ./scripts/release-promotion for ID selection, paginated drafts, redirects, checksum recovery, content-change refusal and latest-only publication fixtures.
2. Run focused watcher recovery tests under go test -race ./internal/watcher, including fake-clock outage/recovery, suspended pending work, definition changes and late hints.
3. Run the payload binding/concurrency regressions in test/integration with native Bash/sqlite3 and PowerShell/sqlite3. Verify distinct numeric-looking IDs and exact literal text before and after the binding correction. Retain actual per-writer errors.
4. Run sh scripts/verify.sh all in the foreground through the hidden command path and watch all eight named gates finish. The aggregate is the local completion gate.

Record source identity, failing-before/passing-after evidence, native platform actually exercised, gate outcomes and remaining deviations in verification.md. Hosted macOS and live release mutation remain distinct from local fixtures. No promotion or actual release is required by this guide.
