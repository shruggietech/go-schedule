# Quickstart: Verify Activity Diagnostics

1. Start the daemon with its default configuration and open Activity.
2. Confirm Activity identifies itself as a limited recent view and displays the configured full-log path.
3. Restart with a `log_file_path` override containing spaces and confirm Activity preserves it exactly.
4. Confirm Clear View, severity filtering, detail selection, alerts, and live records still behave normally.
5. Inspect daemon output and confirm one `daemon startup complete` record carries `endpoint`, `db`, and `log_path`.
6. Run `goschedule logs` and `goschedule logs --json` and confirm their existing output shapes remain unchanged.
7. Run `sh scripts/verify.sh all` and the encoding/whitespace audits before publication.
