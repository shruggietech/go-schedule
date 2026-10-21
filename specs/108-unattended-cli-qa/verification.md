# S108 Implementation and Verification Record

Status: Implemented and locally verified. All feature tasks are complete; local delivery is on review branch codex/108-unattended-cli-qa. The initial implementation halted locally; subsequent authorized publication and review are recorded in [review.md](review.md).

## Workflow and baseline

2026-10-21: Clean main was fetched/pruned and fast-forwarded to 57461fcf350967f9b50cc9d999a0d6c84ada8411 before source research. Review branch is codex/108-unattended-cli-qa. Installed speckit-specify, clarify, checklist, plan, tasks and read-only analyze ran in order. Both requirements checklists passed. Eighteen functional requirements and six success criteria map to 27 tasks. The plan required three read-only research agents, all completed; implementation remained in the primary agent.

## Delivered behavior

The full #289/#290 bundle implements a quiet PowerShell entry point and Go coordinator, six fixed Windows/Linux executable inputs, initialized 152-node Cobra inventory and 839 metadata/behavior/negative cases per environment. Fresh JSON and human fixtures run sequentially against owned daemons. Applied Windows Jobs and Linux cgroups independently verify CPU, memory, process/task and priority controls. Native service, GUI and credential operations use scoped adapters; private-helper entry and parent-channel regressions execute the actual boundary without requesting a native credential.

The evaluator covers task/group mutation and clears, real child execution, schedules and deterministic timing policies, chains, reusable triggers/sets, native watcher recovery, webhook delivery/retry/timeout/redirect rejection, independent notification policies, TLS identity and authority, pairing, profile/preferences persistence, MCP stdio/HTTP operations, history, audit, logs, summaries, search, bundles and deliberate malformed inputs. Every outcome retains its declared packaged_cli, command_harness or deterministic_engine boundary.

Six strict Draft 2020-12 schemas, duplicate-key rejection, confined bounded redaction, channel/corpus identity, complete event transitions, case/environment completeness, truthful summaries, deterministic seven-section Markdown and exact contextual agent reviews govern capture-review-finalize. Frozen capture identity binds source, budgets and environment claims. Incremental case checkpoints and append-only event writes avoid rewriting the entire ledger per case. Full worker snapshots occur every 32 cases and at lifecycle boundaries.

## Recovery evidence

The stable-ID input set .test-temp/s108-inputs-16 was used for a controlled cancellation after the first incremental records. The standalone .test-temp/s108-interrupted-16 package contains all 1,678 planned results, five execution passes, zero execution failures, 1,675 effective blocked results and three manifest-defined platform exclusions. Exact owned Windows process-tree drainage and runtime removal succeeded, with no remaining owned resource. Attempt 2 reused the same source, binaries, manifest, budgets and selected existing WSL environment in .test-temp/s108-run-16, replayed every case under a fresh run identity, carried prior cleanup obligations forward, and completed all applicable cases.

The final fresh capture is .test-temp/s108-run-17 with frozen .test-temp/s108-inputs-17. It ran from 2026-10-21T09:19:06.1968282Z to 2026-10-21T09:24:50.203687Z on Windows 11 build 26300 and Ubuntu 24.04.3 LTS in the existing Ubuntu-24.04 WSL2 distribution. No runtime provisioning, elevation or account change was needed. Both output modes completed, with four successful exact cleanup actions and no remaining owned resource.

A cancellation regression now preserves an active interrupted case as blocked/run_interrupted with captured channels, rather than inventing an application defect. Focused tests reject live launcher leases, changed source/binary/budget/manifest/environment inputs, unknown prior WSL state, corrupt incremental records, truncated/reordered event journals, forged pass claims, missing or stale reviews, omitted channels, escaped artifact paths, hash mismatches and rewritten frozen source/resource claims. Real child tests cover capture limits, deadlines, EOF cleanup and Job controls.

## Final report and retained evidence

Run ID: 0f1e9217-fd59-4f0e-bee9-f9c1d8d441a3. Subject commit: 57461fcf350967f9b50cc9d999a0d6c84ada8411, dirty=true, version gosched version 0.0.1-dev. Manifest SHA-256: 5c07180d2d7a59d88d587900c8174f4da66e12cebf1d3808c855e80ded6d42b5. The report records exact executable identities rather than claiming the dirty source is a published revision. The subsequent lint-only tagged-switch rewrite changes no capture or rendering behavior.

| Result | Count |
| --- | --- |
| Planned combinations | 1,678 |
| Applicable cases executed | 1,675 |
| Execution pass / fail | 1,085 / 590 |
| Final pass / fail / blocked | 1,009 / 666 / 0 |
| Manifest-defined platform exclusions | 3 |
| Distinct messages reviewed / failing reviews | 1,040 / 75 |
| Independently bounded product findings | 12 |

The 76 additional failing cases have communication findings. All seven rubric fields and exact excerpts are retained for every reviewed message. Finalization and independent validation exited zero while truthfully printing fail (finalized); that is the evaluated application outcome, not a malformed or incomplete report. Re-rendering produced identical Markdown bytes.

Canonical JSON: .test-temp/s108-run-17/report.json, SHA-256 ddb3be057c98c75540b244d5c223c085445d474b5d7c06e2d0c780bd8d945733. Human report: .test-temp/s108-run-17/report.md, SHA-256 2fc001f0c498fcabf349e9a38b4238d049a2eb83c8965aa764a34394abfaf37c. The same owned root retains corpus.json, events.jsonl, evidence manifest, frozen capture identity, case channels and source/state/log evidence. The external submitted review set is .test-temp/s108-reviews-17.json. Exact inputs are:

| Role | OS / architecture | SHA-256 |
| --- | --- | --- |
| cli | windows / amd64 | dfcd3c6d96c5c8c2570c3cf9eeaa3e44fb6ca74dfe2dd57d100ed340d81f3c9b |
| daemon | windows / amd64 | 0073cc71ab59a8a0a0627e68299e980edb324784b02865ac079787caaed69805 |
| runner | windows / amd64 | 47a937687e16ab66ce80a5e848dd44bd53468ae32b75c1e3f116a81df7863f1b |
| cli | linux / amd64 | c27b50b435b86db1c5b24b984f3b46da918ec726115e792d239cb6647d62164c |
| daemon | linux / amd64 | 6b41a907d30a3c042c32d5a7eaf5f5f988fff00f83ef8323b8aeb6297e2e3819 |
| runner | linux / amd64 | d8c8f24b7c5e8d14cb2322ad51633b852dc4497152ef456be83f1272a912c999 |

The final report verified independent 25% CPU denominators (aggregate Windows CPU and one Linux logical CPU), 1 GiB memory, Windows 32 processes, Linux 128 tasks/threads and zero swap, lower priority, 15-second case deadlines, a shared 30-minute coordinator deadline, 256 KiB channels and 64 MiB retained artifacts. No resource-limit value was inferred from the Windows host for Linux.

## Final checks

The final canonical foreground sh scripts/verify.sh all completed with exit zero. All eight required gates passed: format, vet, lint, race, GUI, coverage, docs and automation. Format printed no Go filenames; lint reported zero issues. The automation policy check and its adversarial fixture regression both passed. The aggregate log is .test-temp/verify-s108-delivery.log.

Current-source native Linux race checks passed with count=1 for internal/clievaluate, internal/clifixture, internal/cli and cmd/gosched-qa under an owned CPU/memory/task-bounded transient user service using its existing Go 1.26.4 toolchain (exit zero). The evaluated Windows/Linux packaged inputs were built by the repository-selected Go 1.26.0 toolchain. Current-source darwin/amd64 and darwin/arm64 coordinator compilation passed. Native Windows behavior and Job controls ran through the canonical root race gate and actual capture. The GUI gate passed the separate desktop Go race module, 158 frontend tests in 28 files, TypeScript/Vite production output and a native Wails application build. Core coverage remained engine 83.6%, schedule 89.1%, timezone 91.3%, store 80.0%, catchup 88.9% and logbus 91.1%. Static docs exported all 22 canonical pages with links and metadata audited.

Adversarial regressions reject missing or duplicate cases, stale or incomplete reviews, unsupported passes, missing message channels, invalid finding spans, escaped paths, hash mismatches, duplicate/unknown JSON fields, reordered/truncated events and rewritten frozen source/resource claims. Stable path-derived generated case IDs retain references when unrelated commands are added. The renderer preserves canonical detail while presenting compact summaries and collapsed complete ledgers. The final encoding audit found valid UTF-8 without BOM, no mojibake and no prohibited em dash across all 124 changed text files. Git diff whitespace and GitHub-format checks passed. The final lifecycle check passed for all 106 specifications after S108 was marked Implemented with all 27 tasks checked. No pinned process artifact or dependency changed.

## Functional findings and publication boundary

All #289/#290 functional acceptance criteria are delivered; [spec.md](spec.md) and [tasks.md](tasks.md) retain issue-level traceability. Their evaluated product defects remain explicit failures, independently bounded in the [12 issue drafts and publication manifest](findings/README.md). No new defect was silently folded into S108 product implementation, and no issue is closed merely by inclusion in the bundle.

Native service, GUI launch and credential-store mutation cases use their intended scoped adapters; private-helper entry/foreign-parent/channel behavior exercises the actual boundary without requesting a native credential. These results do not claim installed-service, interactive GUI or native credential-storage qualification. macOS was compiled, not executed. No material functional scope was deferred.

The initial implementation stopped at the requested pre-publication boundary. The owner subsequently authorized publication and bounded review on 2026-10-21. PR #310 and product defects #311 through #322 are published; CI/review corrections #323 through #329 and their validation are documented in [review.md](review.md). No merge, release or public-site change is performed.
