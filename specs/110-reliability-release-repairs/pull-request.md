# S110 reliability and release-promotion repairs

Standing-waiver promotion can lose access to a draft through tag lookup, observer recovery can discard an atomic replacement, and SQLite CLI binding can collapse distinct numeric-looking heartbeat IDs into the same value. S110 corrects all three demonstrated mechanisms.

## Changes

- Resolve one authenticated release ID across all pages and use it for downloads, immutable content checks, checksum recovery, publication and latest-pointer repair. Build the helper from the exact workflow source before checking out the release tag.
- Preserve unchanged watcher revisions, pending work and observed file identities through full or root observer recovery; suppress recovery duplicates while retaining later writes.
- Bind payload values as exact UTF-8 text in both Bash and PowerShell, retain final SQLite errors, and keep existing bounded contention attempts and concurrent record-count assertions.

- Stabilize the Windows native credential forged-parent fixture by establishing real endpoints before authentication can reject and close the peer. Kernel peer, parent and image checks remain unchanged (#339).

## Issue traceability

Closes #262. Closes #271. Closes #335. Closes #339. Refs #277. Existing issue milestones remain authoritative.

## Validation

The initial S110 implementation passed all eight local canonical gates with aggregate exit 0: format, vet, lint, race, GUI, coverage, docs and automation. The GUI suite passed 160 tests. Focused watcher, release-helper and payload checks also passed under the race detector on native Windows and Linux. See [verification](specs/110-reliability-release-repairs/verification.md) for failing-before and passing-after evidence and exact aggregate log identity.

## Boundaries

The earlier hosted heartbeat failure discarded the diagnostics needed to establish its exact cause. The binding defect is independently reproduced and corrected; passing repetitions are not presented as a diagnosis of that historical run. All 22 hosted checks, including native macOS, passed on the first reviewed head 8badade. Both bounded Codex rounds completed without findings before a subsequent Windows fixture failure (#339). The correction passes 100 native Windows race repetitions and the full credential-helper package on Windows and Linux. Final-head results and precise review subjects are tracked in the PR handoff and [publication record](specs/110-reliability-release-repairs/review.md); no third review is requested. No live release promotion, tag, release or DNS operation was performed.
