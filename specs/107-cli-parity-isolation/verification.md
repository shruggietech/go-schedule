# S107 verification

## Planning and scope

2026-10-19: Installed specify, clarify, checklist, plan, tasks and read-only analyze completed. Fourteen functional requirements map to 30 tasks, 100% coverage. No critical/high finding, constitutional conflict or unmapped task. Requirements checklist 15/15 and contract/isolation checklist 13/13 pass. Optional agent-context update applied through the CLAUDE.md plan reference.

Baseline: bae303071d1dd1af050cffb201267052ae58b2c7. Delivery branch: codex/107-cli-parity-isolation. Initial local implementation was published after owner approval as PR #305 at 44535cd336afc63cc539090e14c01295f40a9057. Publication remains limited to this review branch and PR. No merge, release, public DNS change or installed service operation is authorized by this handoff.

## Functional coverage

| Issue | Delivered outcome | Direct evidence |
| --- | --- | --- |
| #283 | Set/clear/omit task execution inputs, disabled task/group creation, schedule syntax intent, group rename/reparent/root and complete non-mutating timing preview | internal/cli/authoring_contract_test.go, group_atomic_test.go |
| #284 | Heterogeneous direct notification policies from strict UTF-8 file/stdin JSON, canonical export/reapply, validation before replacement and inheritance restoration | internal/cli/notification_input_test.go |
| #285 | Bounded calendar, active/exact/detailed runs, local runtime/log metadata and selected/all-system summary/search with identities and partial failures | internal/cli/operational_test.go, logs_test.go, runs_test.go, internal/clientobserve tests |
| #286 | Shared durable appearance/popup preference migration, scoped atomic updates/resets and saved desktop startup connection choice without launching the GUI | internal/userprefs tests, internal/cli/preferences_test.go, desktop/settings tests |
| #287 | Explicit shared environment-clear transport used by CLI, desktop and MCP; omission/replacement preserved and conflicting intent rejected before mutation | internal/api/server/environment_clear_test.go, internal/api/client/environment_clear_test.go, internal/remote/environment_clear_test.go, desktop/taskgroup/environment_clear_test.go, internal/mcpmanage/environment_clear_test.go |
| #288 | Invocation-scoped adapters, explicit IPC selection without fallback, authoritative runtime identity, owned workspace guards, hidden process containment and focused real-handler/package fixtures | internal/cli/invocation_test.go, internal/clifixture tests, test/integration/cli_isolation_test.go |
| #304 | Rejected combined group rename/reparent leaves persisted state unchanged | internal/cli/group_atomic_test.go |
| #306 | Cancellation-aware credential setup and joined native-read helpers, preserving successful observations | internal/clientobserve/observe_test.go, internal/clientsecret/context_test.go |
| #307 | Complete bounded Unix process-group drainage shared by fixtures and native reads | internal/processgroup/terminate_unix_test.go, internal/clifixture/process_test.go |
| #308 | Usable macOS default IPC placement and early corrective errors for overlong explicit parents | internal/clifixture/workspace_unix_test.go, workspace_test.go, test/integration/cli_isolation_test.go |
| #309 | Parent-authenticated private credential access, rejecting direct standard-stream invocation and forged peers | internal/clientsecret/context_test.go |

All six original issues retain their own functional scope. Issue #304 records the inherited group mutation defect discovered during authoring work. The comprehensive unattended scenario coordinator, Windows/WSL orchestration, fixed report schema and maintainer prompts remain #289/#290; this slice supplies their required command coverage and isolation foundation.

## Regression findings and corrections

Before #287 was fixed, disposable API regressions failed because clear_env was ignored and conflicting clear/replacement mutated the task. The desktop wire regression produced {"args":null}, demonstrating the empty environment map was omitted. The MCP regression likewise demonstrated that ClearSensitive sent an omitted empty map. Explicit ClearEnv now carries that intent across all three clients and the remote schema; Observe/Operate authority remains denied and Manage succeeds.

Before #304 was fixed, a disposable CLI/API regression returned the expected cyclic-parent error but persisted the proposed new group name. One store transaction now validates and changes the complete edit, preserving rename-only, parent-only and notification-state behavior.

Initial preference locking treated an abandoned lock file as a live writer. A failing regression led to shared operating-system file locks that release on process death, retain a stable lock inode and permit lock-free reads of an existing atomically replaced document. Exact-child death, concurrent scoped writers and unrelated preference preservation are covered.

The first aggregate attempt stopped at lint. Unused wrappers were removed, the authoring helper's error return moved last, imports normalized and handle-close error handling corrected. A second aggregate attempt passed code, race, GUI and coverage gates, then caught a disallowed JSON fence tag and an out-of-tree changelog link in the new CLI documentation. Both were corrected to existing documentation conventions. No lint rule or verification gate was weakened.

Windows process containment is established before resuming the child, inherited handles are explicitly limited, and teardown retains the job handle until the owned tree has exited or its deadline fails. Cleanup errors preserve runtime files for diagnosis. Native Windows regressions cover descendants, unrelated inheritable handles, cancellation, identity/path mismatch and retained evidence. Microsoft references: [job objects](https://learn.microsoft.com/en-us/windows/win32/procthread/job-objects), [TerminateJobObject](https://learn.microsoft.com/en-us/windows/win32/api/jobapi2/nf-jobapi2-terminatejobobject).

## Verification status

2026-10-19: The complete canonical sh scripts/verify.sh all command finished successfully with exit 0. All eight gates ran in order; no gate was skipped or substituted.

| Gate | Result |
| --- | --- |
| format | PASS: Go format and GitHub publishing format |
| vet | PASS: all root packages |
| lint | PASS: zero issues |
| race | PASS: all root packages and integration tests, including packaged CLI isolation |
| gui | PASS: desktop race tests, native Wails build, 28 frontend files and 158 tests, TypeScript and production frontend build |
| coverage | PASS: engine 82.9%, schedule 89.1%, timezone 91.3%, store 80.0%, catchup 88.9%, logbus 91.1% |
| docs | PASS: product/remote policy fixtures, comparison evidence, seven website tests and static export of 21 canonical documentation pages |
| automation | PASS: production automation contract and all configuration-mutation fixtures |

Final review covers all six functional issue outcomes, complete issue-level traceability, UTF-8 without BOM, mojibake checks and git diff whitespace checks. Source changes are complete; owner-authorized publication and review follow-up are recorded below, with a stop before merge.

The canonical command runs in a monitored foreground chain with hidden redirected Windows launchers. GOMAXPROCS=4 and GOFLAGS=-p=4 bound build/test parallelism without changing the eight gates. Existing CI, verification scripts, toolchain pins and frontend contracts are unchanged. The remote OpenAPI client was regenerated with the repository's existing pinned generator.

## Evidence interpretation

The packaged fixture builds real gosched/goschedd binaries, selects a unique owned IPC endpoint, binds installation ID/PID/effective paths, and checks ownership before each mutation. It creates, executes, inspects and removes a benign task, retains command stdout/stderr and daemon logs under an evidence directory, and removes only disposable runtime state after process cleanup.

Evidence kinds remain explicit: packaged_cli means an actual CLI executable and daemon; command_harness means real handlers with explicitly injected adapters; deterministic_engine means controlled engine fixtures. Recording service/launcher/secret adapters demonstrate requested handler behavior without touching installed services, native credentials, desktop sessions or personal preference files. These boundaries are not interchangeable evidence.

Daemon and CLI cross-compilation passed for linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64 and windows/arm64 using CGO_ENABLED=0, matching the existing CI build matrix. This is compile evidence, not runtime execution on those targets.

Native execution evidence includes the complete Windows aggregate and the focused Linux/amd64 WSL checks recorded below. macOS runtime and the complete hosted platform matrix are reported separately in PR #305; neither Windows results nor cross-compilation imply those outcomes.

## Publication and review round 1

PR #305 is open in milestone 15. Codex completed round 1 against 44535cd336afc63cc539090e14c01295f40a9057 with two findings, independently tracked as #306 and #307. Hosted macOS race failures exposed the additional socket-placement defect #308. All three records link parent #277 and their original #285/#288 scope. No third review round is authorized.

A new regression reproduced the credential hang in summary, search and cancellation. After context-aware adapter/native-helper changes, Windows race checks pass for clientsecret, clientobserve and CLI; private-pipe success, redacted failure, oversized output rejection and joined cancellation are covered without native keyring access. Unix group-drain contract tests and real owned process-tree tests pass under native Linux/amd64 Go 1.26.4 in WSL Ubuntu 24.04, alongside the affected CLI/observation/credential packages and actual packaged CLI/daemon execution. A WSL resolver cache miss was handled by reading the existing checksum-verified Windows Go module download cache as a file proxy; host DNS and checksum policies were unchanged.

The corrected source completed sh scripts/verify.sh all with exit 0 after all three fixes: format, vet, lint, root race tests, GUI tests/build, coverage, documentation/static export and automation fixtures all pass. Core coverage remains engine 82.9%, schedule 89.1%, timezone 91.3%, store 80.0%, catchup 88.9% and logbus 91.1%. The final placement and response-cap regressions also pass in the renewed Linux/WSL run, including actual packaged CLI/daemon execution. Daemon/CLI cross-compilation passes all six supported targets again. Final metadata-only format and lifecycle checks precede the review-fix commit. No corrected-head hosted result is claimed from these local checks.

## Codex review round 2

The second and final requested Codex round completed against 73f00d1fddd04570d670b2ca1f134cd6d7c0524c with one finding (#309): the private flag alone did not authorize native access and exposed a standard-stream token protocol. The correction removes standard-stream credential endpoints, authenticates kernel peer PIDs and the actual parent executable image, and retains bounded cancellation/cleanup. Windows and native Linux/WSL race regressions pass for authorized private reads, direct flag rejection, a packaged CLI rejecting an unrelated parent image before request consumption, forged client/parent endpoints, response bounds, cancellation and helper self-expiry. Native user credentials are never used by these regressions.

The final complete foreground canonical aggregate finished with exit 0 after the authenticated-channel correction: format, vet, lint, root race tests, GUI tests/build, coverage, documentation/static export and automation fixtures all pass. The final six-target daemon/CLI build matrix also passes. After the subsequent Unix-only probe correction, native Linux/WSL race checks pass for processgroup, clifixture, clientsecret, clientobserve and CLI, and the packaged CLI/daemon test passes. Final format and delivery-metadata checks precede the commit. No third review round will be requested. Exact-head hosted and review-resolution evidence remains in PR #305.

The 73f00d1 macOS run confirms usable fixture sockets and packaged operation, but its process-tree case encounters EPERM during a post-termination group probe. A new native Linux regression fails on the original immediate abort, then passes after bounded retries of ambiguous post-signal probes. Initial signal permission failures still fail immediately, persistent probe ambiguity still errors at the same deadline, and successful cleanup still requires ESRCH. The complete affected Linux/WSL race and packaged checks are repeated after that Unix-only correction.
