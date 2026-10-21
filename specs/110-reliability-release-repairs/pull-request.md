# S110 reliability and release-promotion repairs

Standing-waiver promotion can lose access to a draft through tag lookup, observer recovery can discard an atomic replacement, and SQLite CLI binding can collapse distinct numeric-looking heartbeat IDs into the same value. S110 corrects all three demonstrated mechanisms.

## Changes

- Resolve one authenticated release ID across all pages and use it for downloads, immutable content checks, checksum recovery, publication and latest-pointer repair. Build the helper from the exact workflow source before checking out the release tag.
- Preserve unchanged watcher revisions, pending work and observed file identities through full or root observer recovery; suppress recovery duplicates while retaining later writes.
- Bind payload values as exact UTF-8 text in both Bash and PowerShell, retain final SQLite errors, and keep existing bounded contention attempts and concurrent record-count assertions.

## Issue traceability

Closes #262. Closes #271. Closes #335. Refs #277. Existing issue milestones remain authoritative.

## Validation

All eight local canonical gates passed with aggregate exit 0: format, vet, lint, race, GUI, coverage, docs and automation. The GUI suite passed 160 tests. Focused watcher, release-helper and payload checks also passed under the race detector on native Windows and Linux. See [verification](specs/110-reliability-release-repairs/verification.md) for failing-before and passing-after evidence and exact aggregate log identity.

## Boundaries

The earlier hosted heartbeat failure discarded the diagnostics needed to establish its exact cause. The binding defect is independently reproduced and corrected; passing repetitions are not presented as a diagnosis of that historical run. All 22 hosted checks, including native macOS, passed on the first reviewed head 8badade. Final-head results and the two bounded Codex review rounds are tracked in the PR handoff and [publication record](specs/110-reliability-release-repairs/review.md). No live release promotion, tag, release or DNS operation was performed.
