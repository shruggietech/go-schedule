# S115 data model

## Source preflight input

Repository root, canonical owner/repository and strict selected vMAJOR.MINOR.PATCH tag identify local README, CHANGELOG and tag-specific notes. Validate version syntax before using it as a path. Read-only results are success or actionable failure naming the affected file and correction. No Git, credentials, network or ref mutation.

## Candidate source and observations

Existing build receipt binds commit, dirty state, sorted source digest and six executable hashes. Capture binds existing receipt, environment and resource budget. A changed source or binary refuses reuse; partial observations retain their own outcome and never become complete implicitly. Finalization binds actual semantic decisions to the current message corpus. Metadata recorded later must name the captured checkpoint explicitly.

## Release handoff and issue state

Pending source version remains v1.6.0, public latest remains v1.5.1, staged package identity is absent until separately authorized. #354 remains open for actual publication. #355/#356 remain independent unknown-cause records unless their own functional correction is established. New demonstrated defects receive separate native issue identities and source references.

## IPC workload

Each named fixture owns a store, endpoint, seeded task and bounded goroutines. Two snapshot workers overlap exactly sixteen create/delete pairs. Successful snapshots contain the stable task and all expected state keys, final groups are empty, all workers terminate through context cancellation and join. File-backed fixture owns and removes its temporary WAL database through normal test cleanup.
