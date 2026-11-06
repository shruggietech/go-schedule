# S134 Operational Data Model

No application schema, API shape, stored identity or state conversion is introduced.

| Entity | Fields | Constraints and transitions |
| --- | --- | --- |
| Maintenance candidate | base/head/tree, target version, proposal versions, issue outcomes | One coherent source identity; superseded pins and any incomplete outcome are explicit |
| Navigation transition | origin guide, destination URL, document readiness | Guide -> requested first-task route -> expected document -> next guide; a visible old heading is insufficient |
| Browser failure artifact | run/attempt, producing step, existing diagnostics, retention | Only generated owned diagnostics are retained; unsuccessful exercise remains unsuccessful; missing files are not evidence |
| Pages receipt | source, artifact ID, attempt, effective deadline, outcome | Build success and deployment success differ; timeout remains failed; later source/artifact-bound recovery is a distinct receipt |
| Release review record | target tag, dataset/changelog/notes hashes, capability decisions | Hash actual prepared bytes; preserve historical tagged records and unknown/edition qualifiers |
| Review ledger | request ordinal, automatic/manual kind, head, status/findings | Maximum two requests including automatic; completion/security scope follows actual receipts |
| Publication continuation | accepted main, approvals, staged asset identities, discovery | Source-ready -> owner integration -> exact-source staging approval -> actual asset inspection/approval -> public delivery |

All owned runtime resources are drained through the existing QA workflow. Explicit/resumed budgets and CPU/memory/process/output/identity limits retain their existing contracts.
