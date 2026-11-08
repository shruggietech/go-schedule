# Research: Maintainer command gates

## 2026-11-08 accepted-source assessment

Accepted main e41005b4dfbc16a54f8a3e8df77bd54583bdded6 includes the compiled policy corpus but omits lifecycle/public adapter regressions from canonical automation. The native test job already has Ubuntu, Windows and macOS runners. Reuse that boundary rather than creating another workflow or gate list.

The old failure-propagation case executes GOFMT=false against the actual aggregate. If stopping regresses after integration it can reach automation again. An isolated exact driver copy plus inert GO, GOFMT, SH and npm children makes both normal and deliberately broken stopping behavior finite. A status-37 sentinel distinguishes native failure preservation from a generic failure.

Reuse the existing automation fixture writer as source-only functions. Target fixtures intentionally lack checker source so the real wrappers must compile in the source checkout and inspect the caller-selected root. A fixed wrapper/root matrix adds adapter evidence without changing the one-build policy corpus. The skill-directed research agent performed read-only assessment; the primary agent owns implementation and analyze.
