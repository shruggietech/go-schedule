# S112 delivery and later publication

## Authorized boundary

Owner kickoff authorizes specification, implementation, verification, push, official PR publication and review fixes with at most two Codex rounds. Stop before merge. No tag, release workflow or public promotion is authorized by this source-preparation request.

## Delivered source preparation

- #352: terminal run selection, bounded hidden probe commands and retained primary/diagnostic/cleanup evidence; real LocalSystem assertions remain in disposable hosted CI.
- #353: separate bounded packaged CLI setup and runtime contexts with setup cancellation diagnostics; native Windows and Ubuntu WSL race regressions retain peer/image/parent checks.
- Dependency proposals #351/#273, #336, #337, #295, #297/#299/#301 converge in this branch. #273 overlaps #351. Existing PRs remain open until replacement delivery reaches main; do not merge or close them early.
- v1.6.0 changelog, concise notes, README health example, v23 upgrade guidance and release-bound comparison/review are prepared. site/data/published.json and the v1.5.1 snapshot remain unchanged.
- S111 merge metadata is reconciled chronologically without erasing its historical two-round review evidence.

## Independent discovery

#355 records a Windows payload-runtime/SQLite-prerequisite observation from local verification. Five unchanged focused concurrency repetitions passed. The interpreter crash and SQLite discovery refusal are retained as independent evidence; neither is claimed corrected by #352/#353 or by a successful repetition.

Independent evaluation findings #356/#357 retain the interrupted default-budget capture and its unresolved cause/reporting behavior. The fresh full capture passed all 1,763 executions but finalized with the P3 SMTP diagnostic-label finding #358 (two communication failures). All four fresh runtimes were removed. These observations remain open; the accepted canonical repository run passed all eight gates. See verification.md for exact source, corpus and native outcome-code evidence.

## After owner review and merge

1. Reconcile the eight superseded dependency proposals against actual main delivery and close redundant PRs with this source PR reference, preserving their traceability.
2. Obtain explicit owner authorization for the exact version/tag and release publication. Verify the selected main commit, exact-head successful CI, clean refs and prepared content digests.
3. Create only the authorized v1.6.0 tag on that reviewed source, triggering the existing tag-only Release workflow. Its exact-source metadata preflight, comparison guard and artifact controls remain required.
4. Wait for staging to finish and retain its exact release ID, run/attempt, source/tag identity, eight candidate assets and SHA256SUMS. Pre-tag local output is not an immutable staged-candidate manifest and must never be presented as one.
5. Use the existing standing-waiver promotion route with the authenticated exact release ID and candidate inputs. Preserve staging-attempt, immutable-tag, candidate bytes, complete asset allowlist, checksum and comparison checks. Absent attended observations remain untested under the standing waiver.
6. Verify actual public asset discovery, latest pointer, Pages deployment and v1.6.0 comparison selection. Update the local published.json default and pending-status prose only after public promotion, with a new reviewed source change if required.
7. Close #354 only when its public-release functional outcome is delivered. Conditional clustered scheduling remains outside S112.

## Evidence

Verification logs and original failed attempts are retained outside the source checkout under the S112 evidence archive. The final PR handoff will identify exact source/CI/review evidence. No owner-installed service, production SMTP account, account session, public DNS or deployment setting is mutated by this source slice.
