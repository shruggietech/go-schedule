# Verification: Stable Daemon Identity and Capability Manifest

## Test-first evidence

- Store tests initially failed to compile because `DaemonIdentity`, `RenameDaemon`, `ResetDaemonIdentity`, the domain contract, and schema v16 did not exist.
- API and client tests initially failed to compile because the manifest response and typed lifecycle methods did not exist.
- CLI tests initially failed to compile because the daemon command family did not exist.
- Desktop integration initially failed until connection health and target publication carried daemon-owned identity, name, platform, version, and capabilities.

## Focused verification

| Command | Result |
|---|---|
| `go test ./internal/store ./internal/api/server ./internal/api/client ./internal/cli` | PASS |
| `go test ./...` from `desktop/` | PASS |
| `go test ./...` from the repository root | PASS |
| `bash scripts/spec-lifecycle-check.sh .` | PASS, 75 specifications lifecycle-consistent |
| `bash scripts/docs-check.sh` | PASS, 19 pages plus policy and remote-architecture fixtures |
| `go run ./scripts/github-format` | PASS |
| `git diff --check` | PASS |

## Acceptance coverage

- One hundred independent stores produce unique UUIDs and reopen with byte-for-byte equivalent identity records; existing identities reopen without requiring fresh randomness.
- Schema v16 upgrades seeded v15 group, schedule, and task records and initializes exactly one identity.
- Concurrent identity reads observe one singleton; corrupt persistent identity fails closed.
- Rename trims valid names and rejects blank, oversized, and control-bearing values without mutation.
- Reset rejects mismatched and stale confirmations, replaces the ID atomically, and preserves name and scheduler data.
- The manifest has exactly the allowlisted top-level fields, deterministic protocol and capability collections, explicit local-only mode, and no host or secret fields.
- Shared-client and CLI tests cover typed reads, writes, JSON output, validation, and conflict envelopes.
- Backup-copy tests prove restore and clone retain logical identity until an explicit reset, after which the original remains unchanged.
- Desktop tests prove connected target context, including separate OS and architecture facts, comes from the daemon while command examples still receive the expected OS key and established local permissions remain intact.

## Canonical verification

`bash scripts/verify.sh all` passed on 2026-09-09 after local remediations for goimports grouping and preserving separate desktop OS and architecture fields. All eight gates passed: format, vet, lint, race, GUI, coverage, documentation, and automation.

## Review remediation

- First-round Codex review identified that the clone instructions needed to account for SQLite WAL state. The guide now requires stopping the daemon before a raw database-file copy or using a WAL-aware SQLite backup operation while it runs.
- Second-round Codex review identified that trimming occurred before control-character validation. Validation now inspects the original UTF-8 input first, rejecting leading and trailing control characters without mutation while retaining normalization for non-control Unicode space.
