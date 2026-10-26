# S118 record model

| Entity | Fields | Relationship and validation |
| --- | --- | --- |
| Implementation record | Slice, status, required task disposition, Delivery reference | S117 becomes Implemented independently of later release publication; every historical task ID retains an explicit disposition |
| Source identity | Reviewed head, merged commit, tree relationship, exact-head checks | A squash merge retains both SHAs; CI belongs only to the tested SHA |
| Candidate evidence | Frozen source, dirty flag, input digests, original outcome | Immutable S117/S115 failed reports remain failed after later source corrections |
| Issue disposition | Issue, delivered outcome, remaining cause/scope, publication state | #354/#355/#356 remain separate and open; no verification-only closure gate |
| Release state | Pending version, public version, selected content hashes, actual tag/run/release/assets | Current public v1.5.1 remains distinct from pending v1.6.0; nonexistent future IDs stay absent |

Transitions: source preparation -> owner-reviewed merge -> later exact-main tag selection -> authorized draft staging with actual asset/run identities -> separately approved promotion -> public discovery reconciliation. S118 reaches only the reviewed source PR boundary.
