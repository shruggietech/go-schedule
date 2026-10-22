# S111 validation guide

Use an isolated daemon and owned mail fixture; no production account or installed-service operation is required.

1. Run focused SMTP settings, MIME, transport, recipient/recovery and v22 migration tests under the race detector.
2. Exercise the shared channel lifecycle through API, CLI, remote permissions and desktop bindings; verify normal reads and exports omit credentials.
3. Preview/export/apply/compare SMTP portable configuration with valid, missing and stale local bindings.
4. Regenerate the initialized CLI inventory after flag/help changes, build fresh receipt-bound inputs and run the existing Windows/selected-WSL unattended evaluator with SMTP observations.
5. Run foreground `sh scripts/verify.sh all` through all eight canonical gates. Record actual commands/results in verification.md.
6. Publish the authorized review branch and official PR; monitor exact-head CI and all review comments. Complete at most two Codex rounds, then stop before merge.

Expected results: recipient-aware history, no confirmed-recipient replay during retries, explicit uncertainty after interrupted submission, useful redacted diagnostics and unchanged webhook/task outcomes. Live provider delivery and release publication are separate operations.
