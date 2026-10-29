# S126 Data Model

## Dependency selection

Proposal number, upstream version, selected version, source/lockfile paths and selected/deferred rationale. Selection is immutable within a frozen acceptance run; a changed graph requires fresh source attribution. #400/#404 are deferred, not delivered.

## Candidate receipt

Source commit plus dirty/content identity, six binary hashes, QA run ID, fixed-schema report/corpus/review identities, verification timings and first failures. Historical receipts remain byte-identical. Source verification, PR checks, merged-main checks and staged-asset acceptance are different stages.

## Delivery continuation

Source proposal -> owner merge -> exact-main acceptance -> separately approved tag/draft -> observed release/run/asset inventory -> separately approved concrete-asset promotion -> public discovery reconciliation. Unobserved identities remain unresolved. Only the source proposal stage is authorized in this run.
