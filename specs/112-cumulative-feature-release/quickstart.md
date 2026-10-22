# S112 validation

1. Run the queued-before-completed PowerShell regression and packaged credential-helper regression on the unchanged source, retaining original outcomes.
2. Apply the bounded corrections and rerun the same regressions. Use only disposable hosted Windows CI for LocalSystem service operations.
3. Review each dependency patch and upstream migration, then run `sh scripts/verify.sh all` in the foreground with the configured Go, C, Node and WebView tools.
4. Run the owned CLI QA entry point with fresh binary inputs and frozen Windows/WSL captures; finalize against actual communication reviews.
5. Validate pending v1.6.0 comparison and candidate instructions while preserving current-public v1.5.1 selection.
6. Publish the authorized PR, resolve every comment and thread, request no more than a second Codex round and report exact-head terminal CI. Stop before merge or tag.
