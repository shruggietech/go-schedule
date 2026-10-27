# Data Model: S124

No runtime or persisted schema changes are introduced.

| Entity | Identity and fields | Invariants |
|---|---|---|
| Prepared release | v1.7.0, accepted base, source proposal, notes, upgrade guidance, known issues | Preparation is not publication; #394 remains open. |
| Comparison snapshot | version, eight product contexts, ten capability definitions, 80 evidence cells | All dates are actual review dates; historical v1.6.0 bytes stay unchanged. |
| Comparison review | version, source_ref, review date, settled content digests, ten changed/no-change decisions | Existing validator schema and digest calculation apply; future tag existence is not inferred. |
| QA input identity | source tree/index/patch identity, six binary hashes, actual versions | Frozen through capture and finalization; no release identity substitution. |
| QA report and semantic reviews | actual run UUID, corpus SHA256, cases, occurrences, seven rubric criteria, findings | Use existing schema and generated report; retain first failure and cleanup receipts. |
| Publication ledger | source commit, PR, exact-head checks, two-round receipts, final audit | No merge/tag/staging/promotion transitions in S124. |

State progression: planned source, assembled source, locally qualified source, published source PR, owner review. Actual release states follow only after separate authorization.
