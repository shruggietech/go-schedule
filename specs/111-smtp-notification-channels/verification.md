# S111 verification

## 2026-10-22: Local canonical verification

Foreground `sh scripts/verify.sh all` completed with exit 0 on the corrected S111 implementation. All eight gates passed without omissions or weakened assertions.

| Gate | Result |
| --- | --- |
| Format | gofmt clean; github-format clean |
| Vet | Passed |
| Lint | Zero issues |
| Race | All root packages and integration tests passed |
| GUI | All desktop Go race tests, native Windows Wails build, 161 frontend tests and production frontend build passed |
| Coverage | Engine 83.6%, schedule 89.1%, timezone 91.3%, store 80.2%, catchup 88.9%, logbus 91.1%; all above the required 80% |
| Docs | 22 canonical documents, links, fences, policy checks, seven website tests and static export passed |
| Automation | Policy, lifecycle and automation mutation fixtures passed |

The retained aggregate log is .test-temp/s111/verify-corrected.log, 24,420 bytes, SHA256 a32b3a5d1bb8b1716aca58c2ba0d31be8a21c2cc71259912d07c4f9913328044. Earlier runs identified import/staticcheck issues, outdated schema/capability expectations, the former remote-only UI expectation, an unsupported documentation fence and store coverage below threshold. Those were corrected before the passing aggregate.

## Focused functional and security regressions

- SMTP: authenticated/anonymous submission, required STARTTLS and implicit TLS, untrusted/expired/wrong-host certificates, no downgrade, shared cancellation/deadline through greeting and DATA, 450/550 recipient and final DATA rejection, partial acceptance, accepted-recipient retry exclusion, lost final acknowledgement, and failed QUIT after confirmed DATA.
- MIME: Unicode task text and subject, escaped HTML, plain text, task/test/daemon-health context, quoted sender/recipient identity and bounded settings/messages.
- Persistence: real prior-v22 migration preserves protected webhook strings, IDs, assignments, heartbeat state and queued work; recipient snapshots/checkpoints/recovery, stale attempt refusal, retry budget and terminal credential erasure pass.
- Interfaces: local/remote API capability and redaction checks, Enroll for stored-credential retargeting, CLI credential file/stdin and explicit clearing, quoted recipient lists, desktop binding/form reset, and portable provisioned-name/revision checks pass.
- Dispatcher: SMTP test and daemon-health messages succeed, a rejected mail destination fails independently, a sibling webhook succeeds and the source task outcome remains unchanged.
- Shared history: empty run output and alert messages remain readable with positive byte limits, truncation and pagination.
- Browser: all three notification Playwright/accessibility tests passed, including keyboard controls, quoted sender and secret reset. Log .test-temp/s111/playwright-quoted.log SHA256 1bf2ed7b0bc93ebbd334911f751434386117de9e8cc63afb32b08cbe7b640398.

## Diagnostic evaluation and independently tracked defects

The first receipt-bound capture, 84daa1cb-1093-41f8-bd85-367a8380fbe0, used source digest d516ffe8cae016c4bb67bdee522748a5c98615e00cc425541aca5c785dd7811a. It completed 1,760 successful executions, three failures and three platform exclusions, then drained owned process trees and removed runtime state before source changes.

Two failures were the new SMTP task case incorrectly expecting JSON from the existing text acknowledgement (#342). Its independent task, recipient, receiver and webhook checks passed. The Windows MCP resource failure exposed the legacy empty-text SQLite projection defect (#344), reproduced directly through Store.ListRunsPage and Store.ListAlertsPage. Separate deterministic probes reproduced quoted mailbox corruption in envelopes/MIME and authoring separators (#343). Missing remote notification routes were reproduced before their correction (#341). Each defect is tracked independently; none is hidden by changing expected product behavior.

## Corrected Windows and WSL evaluation

Final run c29c5da5-15e9-4adb-924d-90b204aed0fa used .test-temp/s111/qa-inputs-02 and .test-temp/s111/qa-run-02. The source subject is baseline e5aa3240676b1bbc0d215bdb6a7b120ce1240d86 with dirty=true and exact content SHA256 b294faa7d031106214aeab7c14aa5b3cab932b6ea5804295026d91a2aa712050. The six executable hashes and exact environments are in its build receipt and canonical report. Later delivery metadata does not change the evaluated implementation.

| Outcome | Count |
| --- | --- |
| Planned environment/case pairs | 1,766 |
| Passed execution and communication review | 1,763 |
| Failed / blocked | 0 / 0 |
| Platform exclusions | 3 |
| Distinct contextual message reviews | 1,109 |
| Final findings | 0 |

All 881 applicable Windows cases and 882 applicable Ubuntu-24.04 WSL cases passed. Windows Unix-only socket/run-as cases and the WSL Windows-only run-as case are the three declared exclusions. Default Windows Job and Linux user-cgroup limits were independently verified. Both JSON and human fixture runtimes and process trees were cleaned; no remaining resources are reported.

The semantic review compared every frozen command/scenario/channel context with retained prior language and current evidence: 901 equivalent normalized prior contexts, 133 changed existing contexts, 74 fresh notification contexts and one new owned SMTP policy observation. Fixture identities, timestamps, paths and repeated records aided comparison; exact original excerpts, occurrence cases, IDs and channels remain in all 1,109 seven-part review records. Finalize, validate and render each exited 0 with outcome pass. The first diagnostic capture remains pending_review and is not represented as a passing finalized report.

| Artifact | Bytes | SHA256 |
| --- | --- | --- |
| qa-run-02/report.json | 10578249 | 84497eab540c03c183a781704e2fe185358d6b9a428602ce566c7637fbe1b586 |
| qa-run-02/report.md | 1895935 | 80065cc5d30a387db5b942aff853d3ed69b6532494bd716ac5296a48052db2e4 |
| qa-run-02/corpus.json | 1807375 | 52f6638e870b424aef625eb587fb82bc875d815be9b07090c1220d1aa45b2ad4 |
| reviews-02.json | 4868487 | 91661eb59fae4ce61642a6a924f6ebc1b78c5ed9d7c13279d24a84ed860ee87c |

Paths in the table are under .test-temp/s111/. The complete evidence directory, six original input binaries and receipt are retained together locally. report.json is authoritative and report.md is generated from it.

## Observation boundaries

The local root race and native build results are Windows evidence. WSL capture executes the Linux packaged CLI/daemon and declared harness/shared-logic scenarios under owned limits. Browser tests exercise the frontend harness. These are separate from future exact-head hosted CI, native macOS checks, production SMTP account delivery and release publication. No live provider credentials, installed-service mutation, public DNS change, merge, tag or release is part of S111.
