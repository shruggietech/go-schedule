# S111 publication and bounded review record

## 2026-10-22: Local completion and authorization

The owner authorized automatic S111 branch push and official PR creation, monitoring CI and third-party reviews, resolving findings autonomously, and no more than two Codex rounds. Stop before merge. No merge, auto-merge, tag or release is authorized.

S111 completes #176, #19 and independently tracked defects #341, #342, #343 and #344. The local canonical aggregate and corrected source-bound Windows/WSL evaluation passed; see [verification.md](verification.md). The initial diagnostic capture and its failures remain separate evidence.

Initial local completion preceded publication. Zero Codex rounds had been requested or completed at that point. Record subsequent publication and review identities chronologically; an automatic ready-for-review round counts toward the same two-round maximum. Do not request a duplicate or third review.

## 2026-10-22: Authorized publication

Pushed b34c58e66411cb6d4c6fd725c29bfc598bb0d4d6 and created [PR #345](https://github.com/shruggietech/go-schedule/pull/345) as a draft to reconcile publication metadata before automatic review. Authentication worked without owner involvement. The PR is open and unmerged, with auto-merge disabled, and closes #176, #19 and #341-#344 upon merge. Existing milestone and label conventions are retained.

The next commit changes delivery/publication metadata only. Marking the PR ready starts automatic Codex round 1; do not send a duplicate manual request. Subsequent actual review subjects, findings, thread responses and final-head CI conclusions are recorded in PR comments. At most one second-round request is authorized after round 1. No third review or merge is authorized.

## 2026-10-22: Round 1 and hosted corrections

Automatic Codex round 1 reviewed 2fe74f089eadc6b98fed2a172c84365ea646c4f3 and completed with one finding: selected remote notification actions lacked capability gating. #346 records its functional correction. #347 records hosted Linux/macOS cancellation/timeout findings in CI run 39747657783. The sender recognizes typed socket timeouts before the context timer callback becomes observable; the greeting fixture synchronizes after EHLO rather than server-side TCP accept. Timeout ceilings, retry policies and CI thresholds are unchanged.

CodeQL alert 7 traced a generated SMTP fixture password through the generic substitution map into evidence redaction and SHA-256 artifact identity. The sink hashes sanitized evidence rather than a stored password verifier. SMTP credentials now stay in typed receiver credential state outside generic command substitutions; fixture secrets remain registered for evidence redaction. No analyzer suppression or password-hashing substitution is introduced. Exact hosted security conclusions will be recorded on the corrected head.

## 2026-10-22: Round 2 and final corrections

The second and final authorized Codex round reviewed 8fa2a7ae1b1ce073b164c7d29e33c07f18a1dfab and reported two findings. #348 records bodyless enable/disable/test calls rejected by the remote JSON Content-Type guard. #349 records raw SMTP settings incorrectly compared with credential authority before normalization. Both reproduced through an authenticated HTTPS client with a verified owned certificate and daemon identity, then passed after their targeted corrections. The capability fixture now sends actual nil request bodies rather than JSON null. A genuine changed server still receives forbidden without partial update. OpenAPI already declares these actions without request bodies and needs no schema change.

Two Codex rounds have completed. Do not request a third review. Final corrections, review responses, thread resolution and exact-head CI evidence are recorded in the PR.

## 2026-10-22: Independent final-head heartbeat CI defect

CI run 39752749979 on 1fadcedc80c695dba208317506896136a27bbcf3 exposed #350: concurrent heartbeat records failed the existing finished_ms >= started_ms constraint. This is distinct from S110's proven literal-byte/session-ID fix for #335. Injecting reversed timestamp samples into disposable script copies reproduces the record loss in both twins, including replacement of an intentional exit code 3 by runtime failure 1. No host clock is changed. The source of the hosted time inversion remains unconfirmed.

Both twins retain the beat with finish equal to start and zero duration, while an unsuppressed warning preserves both observed timestamps and identifies the timing anomaly. Schema constraints, simultaneous-writer assertions, ordinary values, intended task exits and retry budgets stay intact. #350 remains in the current S111 follow-up milestone; completed QA milestone 15 and its historical issues are not reopened. The final source correction receives fresh CI without a third Codex round.
