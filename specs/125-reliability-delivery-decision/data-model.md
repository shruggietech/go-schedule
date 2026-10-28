# S125 Data Model

No production persistence, API, credential or coordination schema changes.

| Entity | Fields | Rules |
|---|---|---|
| Release source | reviewed head, owner merge, exact main acceptance, immutable tag, staging run/attempt, release/assets | Future identities stay unselected until observed; public v1.6.0 is distinct |
| Demand finding | source, actor/job, host/task scale, unmet alternative, availability and failure tradeoff, evidence gaps | Unknown is explicit; owner intent is not measured use |
| Decision gate | predecessor, proceed/hold/stop, authority, reason, next action | #186 positive recommendation precedes #187; approved #187 semantics precede #188 |
| Browser receipt | source, test, exit status, context/trace artifact | First failure is immutable; corrected-source evidence is separate |

State sequence: source proposal -> owner merge -> exact-main acceptance -> separately authorized immutable tag/draft -> observed assets -> separately approved promotion -> public discovery. This slice ends at source proposal. Distributed assessment may hold reversibly without changing product architecture or issue completion.
