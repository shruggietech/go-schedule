# S113 data model

- Evaluation timing policy: validated existing Budget plus explicit bounded shutdown/report allowance, returned by the coordinator for the wrapper. Case and total execution bounds retain their meanings.
- Case diagnostics: bounded phase/request timings, independent command/inspection/protocol/cleanup categories, snapshot completeness and failed request inside existing state evidence. The strict v1 report envelope remains unchanged; new evidence is hash-bound through its existing artifact references.
- Observation: completed partial map, complete flag and failing request/category. Partial equality never proves unchanged complete state.
- Capture lifecycle: capturing/running advances to pending_review/completed or pending_review/interrupted; finalization failure yields explicit recovery information with original identities. No fabricated pass or cleanup claim.
- Delivery: Diagnostic is optional neutral result evidence; LastError is actual unresolved failure evidence. Migration23->24 preserves rows, recipient state, credentials and timestamps, copies SMTP old error evidence into diagnostic, clears only succeeded SMTP error values and is idempotent on reopen.
- Process diagnostic: executable identity, bounded version/probe outcome and original exit/output category; no environment/credential dump and no retry-to-success.
