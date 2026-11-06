# S133 Data Model

No application schema, API or preference changes.

| Entity | Fields and rules | Transition |
| --- | --- | --- |
| Source candidate | selected tag, source revision, content hashes, review/check evidence | prepared, reviewed, owner-integrated |
| Published release | immutable tag/source, exact asset IDs/hashes, checksum inventory | separately approved staging, inspection, promotion |
| Review ledger | request ordinal, trigger receipt, reviewed revision, findings disposition | at most two total including automatic review |

Public v2.0.0 stays selected during source preparation. Current valid state remains durable; unsupported state refusal and developer-controlled reset are unchanged.
