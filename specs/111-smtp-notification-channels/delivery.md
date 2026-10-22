# S111 local delivery

## Outcome

The full SMTP feature #176 is implemented with shared policies, durable task/test/health delivery, protected write-only credentials, recipient-aware retries, CLI/desktop/local/remote administration, typed portable configuration, documentation and unattended evaluation. This completes the functional notification coordinator #19 when merged.

| Issue | Functional delivery |
| --- | --- |
| #176 | SMTP settings, supported encryption/authentication, MIME, multiple recipients, durable lifecycle, all supported interfaces and provider guidance |
| #19 | Webhook and SMTP share daemon-owned conditions/history; desktop-local popups remain distinct |
| #341 | Complete authenticated remote notification route family, OpenAPI and selected-daemon desktop setup |
| #342 | SMTP task evaluation requires the existing acknowledgement and retains exact task/run/delivery/receiver/webhook correlation |
| #343 | Quoted mailboxes preserve identity through CLI/desktop authoring, envelopes and MIME |
| #344 | Bounded run/alert observations preserve empty text, pagination, byte limits and truncation |

These issues remain open pending merge. Their individual functional acceptance criteria are delivered; CI, review and evidence activity are recorded separately and are not issue-closure gates. #146 remains open for its separate distributed-scheduling decision.

## Workflow and decisions

Installed Spec-Kit specify, clarify, checklist, plan, tasks, analyze and implement instructions were followed. Five owner-scope clarifications were resolved autonomously, both requirement-quality checklists passed, and analysis found complete functional coverage with no critical conflict. The installed planning skill specifically required three read-only research agents; implementation was not delegated.

The maintained MIT go-smtp dependency supplies explicit RCPT/final DATA control. Existing protected authorization storage holds write-only SMTP authentication. A forward v23 migration preserves old state. Recipient checkpoints and atomic completion retain acceptance and disclose interrupted submission uncertainty. Portable definitions reuse logical names and frozen target-local IDs/revisions rather than creating a second secret/identity system. Existing Enroll authority covers stored-credential retargeting.

## Verification

All eight canonical gates pass, along with three browser/accessibility tests and the fresh finalized 1,766-pair evaluation (1,763 pass, three declared platform exclusions, 1,109 contextual reviews, zero findings). See [verification.md](verification.md) for exact source, artifact identities and distinctions between local, fixture and future hosted evidence.

## Publication authorization

The owner explicitly authorized automatic push and official PR publication, CI/security review handling and at most two Codex rounds without another pre-publication halt. Stop before merge. No merge, auto-merge, tag or release is authorized. The review branch is codex/111-smtp-notification-channels.

The initial commit records local implementation. Publication identity, exact hosted checks and bounded review evidence belong in the PR and [review.md](review.md).
