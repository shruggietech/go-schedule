# S119 release identity model

| Entity | Fields | Validation and transition |
| --- | --- | --- |
| Approved source | repository, tag, commit, owner approval, exact-main checks | v1.6.0 only at 0f8834ba3b88b69c8f55aaa5d72add70e5b6252b; later branch commits cannot replace it |
| Stage receipt | run, attempt, head, event, workflow path, individual jobs, outcome | Actual terminal successful Release run and manifest-recorded successful Windows attempt |
| Draft | release ID, tag, target, name/body, draft flag, eight assets | Authenticated exact matching release, complete immutable candidate inventory |
| Asset | ID, name, size, state, downloaded SHA-256 | Original bytes retained; MSI matches manifest identity |
| Windows candidate | repository,tag,commit,workflow,run_id,run_attempt,filename,bytes,sha256,product_version,product_code | Existing validator binds exact approved source and actual MSI |
| Promotion decision | owner approval, concrete draft/run/asset identities | Received only after real staging; permits existing standing-waiver route |
| Public release | same release ID and assets, SHA256SUMS, latest pointer, publication time | Eight original staged assets plus checksum inventory; original content retained |
| Discovery record | actual selected public tag, canonical prose, prepared comparison | Current source follows actual publication; generated content is rebuilt rather than hand-edited |

State order: approved source -> approved tag/draft staging -> concrete candidate review -> separately approved promotion -> published asset/latest verification -> reviewed source discovery PR -> owner merge. Failed attempts retain their own outcomes and cannot be relabelled by later success.
