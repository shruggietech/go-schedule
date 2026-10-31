# S129 Release Identity and State

No application database, preference, helper or HTTP schema changes. Existing formats remain application 24, preferences 1, metrics 4 and full qualification evidence 2.

| Entity | Identity and required fields | Relationship and validation |
| --- | --- | --- |
| Prepared source | v2.0.0, reviewed commit, dated changelog section, tag-specific notes | Root module /v2, local desktop /v2/desktop and linker symbol agree; preparation does not create a tag/package. |
| Comparison snapshot/review | tag, refs/tags/v2.0.0, review date, dataset/section/notes SHA256, ten decisions | Review binds exact content; final tagged source is resolved by workflow, not self-embedded in its own commit. |
| Public selection | actual release tag and release ID | Remains v1.7.0/428427523 until later approved publication. |
| Staged candidate | tag/source, Release run/attempt, successful Windows job, release ID, eight asset IDs/names/sizes/SHA256 | Manifest and MSI identities agree; staging is a draft, not availability. |
| Published candidate | same original eight assets plus SHA256SUMS.txt, public/latest, Pages selection | Only approved retained bytes; eight checksum lines cover original assets. |
| Receipt | actual source/artifact identity, result, time, untested scope | No reassignment of historical observations or implication of attended native acceptance. |
| Authorization | source PR, owner merge, exact tag/draft, concrete promotion | Current kickoff supplies source push/PR only; later boundaries remain pending. |

State sequence: prepared source -> reviewed PR -> owner integration -> exact tag/draft approval -> staged actual assets -> concrete promotion approval -> public release and discovery. Every interrupted/failed state remains distinct. #412 closes only when the public functional outcome exists.
