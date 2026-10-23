# S114 research

## Prepared release notes

Decision: correct prepared v1.6.0 release notes, leaving the existing validator intact.

Rationale: scripts/release-promotion/promotion.go requires the release body to contain untested in both resolution and promotion. release.yml stages the exact tag-specific notes, which currently omit this disclosure. A regression can read actual notes into the existing authenticated loopback fixture, call production resolve and require zero mutations.

Alternatives: weaken the validator or invent a new release path. Rejected because they remove existing owner-approved disclosure and identity controls.

## Original Windows payload failures

Decision: preserve the CLR and prerequisite observations independently, add missing launch/failure attribution and investigate without changing the owner's runtime.

Evidence: private S112 verify-all-corrected.log lines 69-113 records writer 2 access violation 0xc0000005 with CLR 0x80131506 and a PowerShell type-resolution stack; Bash writer 0 separately exits 2 with no usable SQLite. Candidate details were suppressed by that historical helper. S113 diagnostic and passing-run evidence is separate.

Current code evidence: Bash sqlite.sh captures timeout-wrapper stderr in probe_dir/launcher and deletes it without reading it. The six-writer fixture currently prints only failed writer results, losing successful stages and row attribution during a failing run. PowerShell's post-exit redirected reads have a separately investigated inherited-pipe deadline hypothesis; no historical causation is claimed.

Authoritative references: [collect .NET crash dumps](https://learn.microsoft.com/en-us/dotnet/core/diagnostics/collect-dumps-crash), [Windows dump analysis](https://learn.microsoft.com/en-us/dotnet/core/diagnostics/debug-windows-dumps). No matching upstream resolution for the original stack was established. No CET change or runtime replacement is justified.

## Owned IPC

Decision: retain the shared case context and limits, exercise actual cancellation and mixed typed requests over owned endpoints, and separate command, follow-up and cleanup phases.

Evidence: original retained failures include successful command results followed by inspection/cleanup deadlines. Existing Windows DialContext delegates to go-winio with the caller context. Some protocol cleanup has an independent bounded reserve, explaining elapsed-time arithmetic but not the original stall. Existing store methods lack request context in several paths, and idle transports lack explicit release; neither is demonstrated as the original cause. A broad storage/pooling rewrite is rejected without causal proof.

## Release and upgrade guidance

Decision: use an S114 delivery runbook, link active operator guidance to it and reconcile stopped-daemon backup/rollback instructions on Linux/macOS with Windows schema-v24 guidance.

Rationale: the macOS GUI-started daemon outlives the window; replacing a bundle alone is not a stopped-state backup. Older binaries must use matching older full-state backup rather than a migrated database. Preserve protected credentials and SQLite sidecars.

Alternatives: rewrite historical S112 v23 documentation or change published v1.5.1 discovery. Rejected because those records are historical and promotion is not authorized.
