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
