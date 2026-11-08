# Validation guide

Use resolved project Go, Git/POSIX shell and native build prerequisites. Run `sh scripts/spec-lifecycle-check.sh .`, `sh scripts/automation-check.sh .`, `sh test/scripts/spec-lifecycle-check_test.sh`, `sh test/scripts/automation-check_test.sh automation` and `sh test/scripts/automation-check_test.sh failure-propagation`. Positive fixtures succeed; invalid policies fail with their existing reason categories.

Run `go test -race ./internal/maintainercheck ./scripts/maintainer-check` and the maintainer benchmarks, then `sh scripts/verify.sh all` in the foreground. Keep before/after checker timings separate from the canonical aggregate and compilation cost. Hosted CI supplies independent Windows/Linux/macOS evidence. No tag, promotion or installed service is needed.
