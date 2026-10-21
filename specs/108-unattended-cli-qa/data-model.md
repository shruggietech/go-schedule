# S108 Data Model

- Manifest: schema_version, command metadata/digest, stable scenarios with command_path/class/kind/flag coverage/expectations and applicability. Change requires explicit inventory update; newly generated probes cannot silently count as behavioral coverage.
- Subject: repository, commit, dirty state, version, binaries with role/OS/architecture/path/SHA-256.
- Run: unique ID, UTC start/finish, phase/completion state, subject, runner version, manifest hash, budgets, complete environment/case ledger and cleanup.
- Environment: stable ID, platform/architecture/tool versions, selected distro/user/kernel when applicable, applied limits and ownership identity.
- Case: stable ID/environment/command/class/kind, expected exit/checks, observed exit/duration/checks, execution_status plus effective status/reason, redacted argv/stdin/stdout/stderr/log refs, side effects and review/finding IDs. Nullable observations remain explicit.
- Artifact: portable relative path/media type/bytes/SHA-256/redaction policy. Root-confined regular-file content is the authority.
- Event: monotonically sequenced typed state transition with run/environment/case identity and bounded sanitized payload.
- Corpus: immutable normalized redacted messages with context-derived IDs and all case/channel occurrences, plus explicit no-message occurrences.
- Review: run/schema/corpus hash binding, exactly one record per required message, exact case/channel/excerpt, named rubric outcomes, verdict and finding refs.
- Finding: stable semantic dedup key, category/severity/product-or-harness, exact affected identity/cases, reproducible fixture invocation, expected/actual/impact, hashed evidence spans and functional acceptance suggestions.
- Cleanup: exact owned resources, attempted actions/results, retained evidence and remaining owned resources requiring action.

State flow: planned -> capturing -> pending_review -> finalized. Interruption produces complete blocked/uncompleted case records and a resumable checkpoint. Functional fail remains fail throughout review. Resume validates every immutable identity and hash before completed cases are reused; changed subjects start a new run.
