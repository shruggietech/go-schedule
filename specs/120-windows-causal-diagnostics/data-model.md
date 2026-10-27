# S120 diagnostic model

| Entity | Required identity and fields | Invariants |
| --- | --- | --- |
| Original evidence | Original path/reference, bytes/hash, source/runtime attribution, native outcome | Immutable; missing identity remains unknown |
| Experiment | Hypothesis, source/input hashes, environment, owned target, supplied bounds, controlled difference | Diagnostic-only changes never inherit package qualification |
| Process observation | Resolved executable/hash, harness hash, PID/start state, duration, native exit when available, collection and cleanup outcomes | Collection failure cannot replace or erase native exit |
| Pool-wait observation | Request/phase category, queue-entry barrier, client cancellation, server completion, pool counters/allowlisted stack category | Completion is checked before releasing held connection |
| Disposition | Issue, demonstrated defect/correction, upstream match or unconfirmed cause, source references | Passing repeats do not imply historical fix |

Ownership lifecycle: prepared -> owned target started -> controlled observation -> primary result retained -> bounded collection/cleanup -> immutable receipt -> separate disposition. Unknown or unavailable collection is explicit. Original failed capture state is never rewritten.
