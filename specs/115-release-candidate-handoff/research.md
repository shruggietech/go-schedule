# S115 research decisions

## Shared source preflight

**Decision**: Extract the Release workflow source checks into `go run ./scripts/release-source-check -root . -repository shruggietech/go-schedule -tag v1.6.0` and call the same command in Release. Preserve badge, health, unique version heading, notes and exact whole-line changelog-link semantics; require an exact calendar-valid YYYY-MM-DD date.

**Rationale**: Existing executable validation lives only after tag push. Automation policy checks its text rather than executing candidate semantics. Current v1.6.0 metadata is valid, but a heading ending in not-a-date satisfies the existing substring guard when notes follow its derived anchor. Record that demonstrated validator fault independently after regression reproduction.

**Alternatives**: Documentation-only preparation leaves behavior unexercised before tags; duplicating shell validation creates drift; staging would cross the owner boundary. Comparison validation remains a separate existing command and no supplied SHA is labelled a verified tag before a tag exists.

## Native writer investigation

**Decision**: One opt-in diagnostic six-writer run on native Windows and one native Linux run, retaining fresh/existing schema cases, six identities, exact counts, runtime identity and original outcomes. An unavailable PowerShell on Linux is explicitly unavailable.

**Rationale**: No current helper fault was demonstrated by inspection. S114 corrected separate launcher defects under #364; those are not the original CLR crash or independent Bash refusal under #355.

**Alternatives**: Repetition until green, runtime replacement, serialization, extra retries or raised case limits would not establish the original cause and would change the experiment.

## Concurrent file-backed IPC

**Decision**: Extend the existing owned IPC fixture with separately named memory/file-backed cases, two simultaneous complete-snapshot workers and exactly sixteen group create/delete pairs. Preserve cancellation, partial-data and fresh-snapshot checks; every successful snapshot retains the seeded task and complete expected collections.

**Rationale**: Existing fixture completes its snapshots before mutation and uses memory storage, whereas production store uses a file and WAL. This is a workload coverage gap, not evidence of a production bug.

**Alternatives**: Changing pipes, pools, dependencies, SQLite connection policy or deadlines without reproduction would be speculative.

## Source and publication identities

**Decision**: Carry 18e6fd8 after analyze; preserve original 574d071 archive. Perform a fresh full owned capture on a frozen current source checkpoint and review its actual messages, retaining original command/observation/cleanup outcomes. Keep full private logs ignored and share attributable summaries.

**Rationale**: Historical S113/S114 captures have distinct sources. Passing fresh observations do not close #355/#356, and no tag or staged package exists for pending v1.6.0.

**Alternatives**: Relabelling old captures or staging under source-only authorization violates identity and publication contracts.
