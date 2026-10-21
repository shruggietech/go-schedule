# S109 verification record

## 2026-10-21: Scope and red-green evidence

Baseline is merged main c6d50a17c4b16a9025f9f7137e3cc05b0a6391bd. The full #311-#322 bundle is implemented locally on codex/109-cli-contract-diagnostics. Meaningful regressions demonstrated the parser/operand, JSON notice, help/list, target collision, duplicate/unknown binding, cron status, diagnostic context, credential identity/atomicity and watcher observation defects before their fixes.

All console commands use the verified hidden launcher and redirected noninteractive I/O. Product validation uses owned targets, disposable credentials/profiles and recording host adapters. No installed service, host session, public site or DNS mutation is performed.

Focused native Windows checks passed for CLI/daemon, bundle server race, enrollment/store/credential API race, watcher/domain/API race, desktop automation race, frontend watcher cards and TypeScript. The first canonical aggregate compiled a late-added missing-profile diagnostic regression before its root-owned context correction and failed that assertion. The correction preserves the original usage exit 2 and underlying sentinel, with identity/action; focused checks pass. Read-only integration review also found generated help/completion classification and explicit-empty remote selector fallback gaps within #311/#312/#314. New regressions reproduced them; generated commands now use the shared contracts, help stays local and propagates writer errors, and changed empty remote selectors refuse before local fallback. Fresh scenario entries cover these boundaries. The final stable-source aggregate and native Linux checks follow.

## Issue traceability

| Issue | Delivered correction |
| --- | --- |
| #311 | Typed parser/operand/document usage classification, operational failures preserved |
| #312 | Default/explicit operand validation before target resolution and host operations |
| #313 | Alerts notice on diagnostic channel, one JSON stdout document |
| #314 | Owned global selectors, --daemon-endpoint, legacy compatibility and alias conflicts |
| #315 | CLI duplicate refusal and strict server duplicate/unknown binding validation before planning |
| #316 | Complete mixed-input report, invalid-input exit 2 and report-write runtime failure |
| #317 | Meaningful credential/pairing/service help |
| #318 | Object/field/file/certificate/profile context, useful correction and retained causes |
| #319 | Accurate rotation actor/daemon context, transactionally assembled response and unchanged eligibility |
| #320 | Labeled human actor/audit/pairing/profile/channel/watcher tables with unchanged JSON shapes |
| #321 | Exact observed-definition revision, pending/time context in API/CLI/desktop without waits |
| #322 | Compatibility warning explains broad local access, admin_group configuration and restart |

## 2026-10-21: Fresh evaluation correction

Diagnostic capture s109-run-01 found completion.callback and completion.no-desc returning 2 because Cobra creates its hidden completion protocol commands only inside ExecuteC, after the public preflight. New failing regressions cover both valid raw protocol forms and the missing-word refusal. The shared execution entry point now bypasses ordinary lookup only for those exact protocol command names, preserving Cobra completion semantics without admitting unknown ordinary commands. New receipt-bound inputs and a complete fresh capture will verify the correction; the diagnostic run is allowed to finish owned cleanup and is not claimed as the final evaluation.

The same diagnostic run exposed an existing QA interpretation defect once human channel creation succeeded: the snapshot identity mapping omitted channels, so the evaluator marked the successful operation failed and blocked nine dependent human cases. A new regression demonstrated the failure with an unrelated new task in the snapshot; adding the channels mapping preserves unchanged/ambiguous-state refusal and actual stdout identity checks. The independently actionable [publication draft](findings/human-channel-identity.md) keeps this discovery separate from the 12 product issues. The fresh runner includes this correction.

An intervening fresh-input rebuild was correctly refused because the source changed during preparation, while final source corrections and automation fixture creation/cleanup were in progress. It produced no receipt. Preparation was repeated after the source and fixtures stabilized.

## 2026-10-21: Final engineering verification

The complete canonical foreground sh scripts/verify.sh all aggregate finished with exit 0 after the final completion and evaluator corrections. Its retained log is .test-temp/verify-s109-delivery.log. No gate was waived or replaced by hosted CI.

| Gate | Result and concrete evidence |
| --- | --- |
| format | PASS, no unformatted Go files; github-format reports no em dashes or hard-wrapped Markdown prose |
| vet | PASS, complete root module |
| lint | PASS, pinned golangci-lint v2.12.0, zero issues |
| race | PASS, complete root module including CLI, API, enrollment, store, watcher, evaluator and integration packages |
| gui | PASS, desktop Go race tests, Wails v2.15.0 native Windows application build, 28 frontend files and 160 tests, TypeScript and Vite production bundle |
| coverage | PASS, engine 83.6%, schedule 89.1%, timezone 91.3%, store 80.1%, catchup 88.9%, logbus 91.1% |
| docs | PASS, policy/architecture/comparison checks and fixtures, seven site tests, 22 canonical docs and static export/link audit |
| automation | PASS, automation-check and its complete automation fixtures |

Current Linux-target lint also completed with zero issues. Focused native Linux race checks passed for the final generated-command, empty-selector, certificate/profile, webhook targeting, contextual-error and human-channel-identity regressions: CLI 1.903s and evaluator 1.204s. Earlier affected-package native checks passed CLI, domain, watcher, API server and enrollment; that first six-package unit reached its deliberate 600-second resource deadline before store completed, so it is not represented as a fully successful suite. A separate focused CLI/store native race run passed (2.595s and 11.703s). The final canonical Windows root race gate covers all affected packages after all source corrections. Native Linux checks used the existing Go 1.26.4 installation; prepared packaged subjects use Go 1.26.0.

## 2026-10-21: Final receipt-bound application evaluation

Final run ea995dd2-8639-448f-b48e-c3aa0cfe9c17 completed from .test-temp/s109-inputs-03 into .test-temp/s109-run-03. The subject is baseline commit c6d50a17c4b16a9025f9f7137e3cc05b0a6391bd with dirty=true and exact source SHA256 15ff4f1f4dba0b1b5eff2b098b405af2321ba5ff9846dfa91647c04e03dbf855. Source remained frozen throughout preparation/capture; subsequent local changes only complete delivery metadata. This receipt-bound subject is distinct from the later local delivery commit and from future hosted CI.

| Subject binary | SHA256 |
| --- | --- |
| windows/cli (amd64) | 1cf128b23b219083f32ee0bbd5879e4896cf534c969c02c2656906d08f513fb9 |
| windows/daemon (amd64) | 6073ae8eb731ec8d7b5e73c97e2a56607dd51e94d9705d9981cfb7e2df6538d4 |
| windows/runner (amd64) | 6d9a71e0c60c9d317d3c1e2a4a4e96902438ac90abec14cc34a0a547dc9f478c |
| linux/cli (amd64) | c6d8dfb967d0a5d9a6b92db254aca813b6b32713388dd18d55096ff2dbd179db |
| linux/daemon (amd64) | 044867a12c133acecb32e3b36479cbc70e84b5f569be0171b19a49ab854db37a |
| linux/runner (amd64) | f79bb0a661721b5abc357c9430c9d20f11045876d325c7f9ed772ee4d8c0e5e7 |

Windows and selected Ubuntu-24.04 WSL2 executed sequentially with verified independent controls: 25% CPU, 1 GiB memory, Windows Job Object process cap 32 and Linux cgroup task cap 128; Linux swap is zero with nice 10. The existing selected environment is Ubuntu 24.04.3 LTS, kernel 6.18.33.2-microsoft-standard-WSL2, amd64. Both environment cleanup records contain two successful drain/runtime-removal actions and empty remaining-resource lists.

| Final canonical result | Count |
| --- | --- |
| Planned combinations | 1,688 |
| Passed execution and communication review | 1,685 |
| Failed | 0 |
| Blocked | 0 |
| Manifest-defined platform exclusions | 3 |
| Distinct complete contextual message reviews | 1,049 |
| Final findings | 0 |

The platform exclusions are Windows isolation.unix-socket-bound and task.run-as-unix, plus WSL task.run-as-windows. All 842 applicable Windows cases and all 843 applicable WSL cases passed. Existing case IDs and assertions were retained, the cron invalid-line expected status was corrected to the documented usage contract, and five targeted scenarios extend the existing inventory. The separate evaluator mapping correction retains actual stdout identity and persisted-state checks.

The executing agent reviewed every frozen command/scenario/channel context using the seven-part rubric, including parser diagnostics, changed help, object/field/file corrections, labeled human output, structured rotation/watcher responses, daemon logs and newly reached human notification cases. Previously reviewed equivalent language was compared only after normalizing fixture identities, timestamps and owned paths for review assistance; every final record retains its exact original excerpt, cases, message ID and channel. Finalization validated all 1,049 records against corpus SHA256 ec5542330c2018055bac1fee651dd0c28a6d98399faab43517c1fb3be1ccbbb0. Finalize, validate and two render operations all exited 0. The two Markdown renders have identical bytes.

| Retained artifact | Bytes | SHA256 |
| --- | --- | --- |
| .test-temp/s109-run-03/report.json | 9653099 | 28642460bd145836a32a61455f7bdf144158cdfbdcfd76cfdf58507aa89f5732 |
| .test-temp/s109-run-03/report.md | 1812332 | 1cbc65a732a87aaf73d1602aeae986d3bdb5a93598f653a4e6eda3a1c32a5a2d |
| .test-temp/s109-run-03/corpus.json | 1743889 | ec5542330c2018055bac1fee651dd0c28a6d98399faab43517c1fb3be1ccbbb0 |
| .test-temp/s109-run-03/manifest.json | 1175008 | 11fb63671f756ed26578051ba8cd98a071930e47d1d3b76a57067df7bd029e42 |
| .test-temp/s109-run-03/evidence-manifest.json | 1535693 | 5c3d43a11a974f1f49937c4b2d3d5d93fcb5a3f114f67cd185123c4dbf38195d |
| .test-temp/s109-inputs-03/build-receipt.json | 1513 | e218cc51657f50262d848a246217fe74ced4c3397eb67afcb216aafe2580f814 |
| .test-temp/s109-reviews-03.json | 4226415 | 9df286c8a4abc8209b389d9a535d3f491f9c0238dfea9624bb8e62aeab25a0f5 |
| .test-temp/verify-s109-delivery.log | 24044 | 348146b33be4bb6c16f8ba12d389d51277d219c3ec45cbf4d6de16e839cc52a7 |

## Disposition and delivery boundaries

F001-F012 map respectively to #311-#322 as recorded in the S108 findings index. Every functional correction in the traceability table above is implemented and covered by the final application evaluation and targeted regressions. No issue is closed merely because it was bundled or evaluated; closing references in the local PR draft apply after the corrections merge. The separate evaluator defect is implemented locally with a [publication draft](findings/human-channel-identity.md), not an invented GitHub issue identity.

These results establish scoped Windows/WSL execution, handler/shared-logic host adapters, deterministic evidence/report contracts and contextual communication review. Installed-service execution, interactive desktop operation and native secret-store mutation are distinct observation boundaries; this feature uses the existing owned fixtures and recording adapters for those surfaces. No public service, DNS, repository publication or release was changed. Remote main was rechecked at c6d50a17c4b16a9025f9f7137e3cc05b0a6391bd before the local delivery commit. See [delivery.md](delivery.md) for the pre-publication handoff and exact commands.

Final delivery metadata was checked separately after the runtime capture. The lifecycle check initially required its canonical review-branch evidence wording in the Delivery field; that wording was corrected without a product change. Final publication-format, lifecycle (107 specifications), whitespace and UTF-8 checks passed. All 71 changed files decoded as UTF-8 without BOM or detected mojibake.
