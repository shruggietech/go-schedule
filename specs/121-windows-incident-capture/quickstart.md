# S121 validation guide

Use the existing hidden noninteractive launcher on Windows. Build source-bound binaries with the existing gosched-qa build command. Run capture with --incidents and the existing explicit nonroot WSL distribution/user selection, then review/finalize through the existing schema.

For script experiments set GOSCHEDULE_INCIDENT_DIR to a new private owned absolute directory and GOSCHEDULE_PROBE_DIAGNOSTICS=1. Run the incident tests, six-writer fresh/existing-schema tests and metadata/SQLite contrasts. The controlled FailFast smoke must either yield an actual native dump or an explicit unavailable outcome, and must retain its nonzero native result.

Targeted Go checks cover internal/incident, internal/clievaluate, internal/store, internal/api/server and test/integration. Exercise live held-pool memory/WAL capture, handler-only stall, expired/canceled collection, wrong identity, symlink/duplicate paths, oversize data and disabled mode. No installed-service or account-session changes are involved.

After all source changes, run the complete canonical foreground aggregate: sh scripts/verify.sh all. Preserve the original failed artifacts and bind new results to source/binary identities. See contracts/incident.md for exact bounds and independent outcomes.
