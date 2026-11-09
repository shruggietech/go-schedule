# S139 full-product audit

## Source and conclusion

S139 assesses accepted base 8b10419b916e04891609b3cef9727da16f486a23 and delivers six compatible source corrections under [#444](https://github.com/shruggietech/go-schedule/issues/444). Public packages remain v2.0.1. This assessment is complete with explicit platform limits; it does not qualify an installed v2.0.2 package or resolve historical #355/#356.

The full discovery capture used six immutable binaries with source-content SHA-256 d85e06909b31f9acef84a1704fa11de3d201dd2336b51eea4ad22295eb327d94, dirty source at the accepted base. All 1766 cases completed in 298.031 seconds: 1763 execution passes, zero execution failures and three not-applicable cases. The complete initialized tree is listed in [cli-surfaces.md](cli-surfaces.md). Evidence kinds are 1012 packaged-CLI cases, 730 command-handler/adapter cases and 24 deterministic-engine cases, including platform-not-applicable rows. Help is not functional execution coverage.

Every one of 1116 distinct contextual messages was read with all associated case contexts, yielding 7812 rubric judgments. Finalize, validate and render used the original coordinator and produced a valid failed report: 1747 case passes, 16 cases linked to four findings, zero blocked cases and three not applicable. Five distinct messages failed communication review. The original result remains failed. It is not reclassified by later fixes. The seven rubrics cover language, outcome, correction, severity, terminology, identifiers/units and credential absence.

Fresh six-binary inputs use source-content SHA-256 0f35e2894faa0e271d629a2bca69667550b680be71646219bc13e4c369f96108. Four compiled component-test binaries and ten fresh Windows/WSL observations passed after the corrections. These establish affected current-source behavior; they do not replace the original full-catalogue identity or claim a second complete catalogue. [Verification](verification.md) separates both generations and the final canonical gates.

## Product-area register

| Area and source boundary | Observed scope | Limit or finding |
| --- | --- | --- |
| Scheduling, timezone and catch-up: internal/schedule, timezone, engine | Full catalogue deterministic DST gap/fold, month-boundary, catch-up, startup, one-off and overlap scenarios; actual owned task/history and chain dispatch | Existing policy retained; no physical clock or performance claim |
| SQLite/current state: internal/store, config, userprefs | Current initialization, reopen/refusal, transactions, durable delivery, private scoped profiles/preferences and canonical race/coverage suites | Historical migrations remain outside the owner-selected prelaunch policy |
| Process execution: internal/executor, clifixture, processgroup, incident | Benign native fixture tasks, stdin/env/args, output bounds, failures, owned tree drain and retained diagnostics | Original #355/#356 causes remain unconfirmed; no new incident claimed |
| Local API/IPC/access: internal/api, ipc, authorization | Real owned Windows pipe and Linux socket; ownership/path guards, authorizations, readbacks and partial-result handling | No owner installation or live state changed |
| Remote enrollment/identity: internal/remote, remoteenroll, clientsecret | Same-host native catalogue, TLS/certificate, actor/credential and permission cases; fresh actual TLS identity-comparison component checks | #450 corrects the old wrong-identity coverage; real cross-OS matrix below is blocked |
| MCP: internal/mcpobserve, mcphttp, remotemcp | Actual bounded protocol dialogues, observe/operate/manage/enroll controls, confirmation/replay/refusal and listener lifecycle | No owner's MCP credential or client configuration used |
| Notifications: internal/notification, testsmtp, store | Actual owned SMTP/webhook receivers with independent task/history/delivery/receiver evidence; retry and credential boundaries | #447 fixes selected acknowledgement checks without weakening correlation |
| Bundles, cron, watchers and trigger sets | Full catalogue import/export/preview/apply, native watcher recovery, hierarchy, source readiness and dispatch cases | Saved definitions, readiness, dispatch acknowledgement and observed completion remain distinct |
| Desktop backend/frontend: desktop and internal/desktopcontrol | Canonical module/frontend/build checks, 35 headless browser workflows and eight inspected screenshots | Injected frontend fixtures, not native Wails or installed backend evidence |
| Installation/service/removal: packaging, internal/service, winuninstall | Existing source/build/installer/repair/removal regression contracts; supplemental actual usage/refusals and Windows read-only service status | No install/start/stop/restart/uninstall or destructive cleanup against the owner's hardcoded goschedd service identity |
| Source/dependencies/docs | Both Go module caches verified; both npm audits returned zero advisories; source/manifest alignment, maintained workflow utilities and canonical static/race/coverage/docs gates | govulncheck is not installed, so no Go advisory scan is claimed; cache verification is integrity evidence only |

## Supplemental entrypoint register

The private supplemental batch contains 55 expected native outcomes. Successful help/refusal checks prove only those paths. The original private driver syntax failure is retained and is not a product finding. Ten maintainer utilities were also observed safely, and both Go module-cache checks passed.

| Entrypoint | Actual observation | Deferred operation and reason |
| --- | --- | --- |
| goschedd | Windows and Linux help, unknown flags, extra operands, missing/malformed configuration; full owned startup/shutdown/recovery catalogue; new operand integration regression | No installed-service daemon replacement |
| gosched-qa | Actual build, inventory/catalogue, timing policy, capture, finalize, validate and render; subcommand help, invalid input/refusal and owned fixture-child success | Resume missing-input refusal observed; no completed immutable capture was mutated merely to create an interruption |
| gosched-cleanup | Fresh Windows binary rejects empty/help/unknown/extra operands with exit 64 before wipe | Valid wipe not run: no dedicated installed product/ledger boundary; existing contained cleanup suites exercise removal/refusal logic |
| gosched-tray | Fresh Windows binary rejects help/unknown arguments and invalid elevated action with exit 2 before singleton/UI/service handling | Native tray and valid service action blocked by shared-session/fixed-service boundary |
| gosched-indicator | Source preflight establishes that arguments, including --help, do not bypass native initialization | Not launched: no isolated graphical display; default WSLg is prohibited |
| gosched-gui and gosched gui | Actual Windows/WSL missing-sibling refusal, then fresh corrected guidance; launcher intent/error adapters in the catalogue | Native desktop launch blocked by isolation constraints |
| Test-Heartbeat, Test-GetSystemInfo, Test-ReadTestDB | PowerShell Windows and POSIX Linux help; actual Linux owned SQLite heartbeat, system snapshot and readback | Real snapshot data stays private; installed scheduling recipes and owner metrics databases not used |
| cli-qa.ps1 and canonical POSIX/Make adapters | Existing hidden wrapper built/captured real six-binary inputs; eight-gate driver and maintained adapter corpus exercised | Convenience clean/tidy/install/dev and release mutation modes were inspected or tested with inert fixtures, not run against owner state |
| github-format, release-source-check, comparison-check, brand-check, ci-timeout-check, maintainer-check | Read-only successful current-source policy, release metadata, historical comparison and brand/lifecycle checks | A source metadata pass establishes no tag, staged package or promotion |
| release-promotion, windows-release-gate, windows-qualification-session, brand-import | Safe help/usage refusal paths and existing canonical package tests | No authenticated publish, candidate install, native qualification session or brand import performed |

## Actual Windows/WSL transport

| Direction | Native local control | Actual opposite-OS attempt | Disposition |
| --- | --- | --- | --- |
| WSL2 Linux client to Windows 11 daemon | Owned Windows daemon started; local HTTPS health and pairing control passed | Actual Linux typed Go client attempted pinned enrollment, with a three-second connect bound | Blocked at enrollment connection; no cross-OS task, remote grant, rotation, restart or history claim |
| Windows 11 client to WSL2 Linux daemon | Owned Linux daemon started under an existing bounded systemd unit; local HTTPS health/pairing passed | Actual Windows typed Go client attempted pinned enrollment, with the same bound | Blocked at enrollment connection; downstream matrix unavailable |

The existing WSL2 environment reports networking mode none, no default route and no nonloopback address. Both real connection attempts were made; same-host successes are not substituted. Global networking, firewall and WSL configuration were not changed. Both valid native server attempts exited successfully and reported runtime removal. An initial private Linux fixture supplied Windows-style limits and was rejected before daemon launch; its correction inherited the owned systemd cgroup. The initial toolchain selection failure and original helper/source/binary identities are retained separately. These are fixture setup failures, not product bugs. Actual cross-OS probes use the typed API/enrollment client, not the packaged profile CLI or native credential storage.

## UI and unattended boundary

All 35 existing Playwright frontend workflows passed headlessly. Eight selected screenshots were individually inspected and copied unchanged into [evidence](evidence/README.md). Their PNG chunks contain only IHDR, IDAT and IEND, with no text/path metadata. Synthetic labels such as Production and fixture version strings describe injected test data. Keyboard accessibility tests target headless Chromium only.

Windows GUI singleton/activation objects are fixed in the Local logon-session namespace. A second desktop in the same session can still signal the owner's existing GUI. No separate existing logon session was available, so no native Windows GUI or tray was launched. Existing WSL has GTK/WebKit but no isolated display server, so neither native Linux GUI nor indicator was launched and WSLg was not used. Linux Secret Service can request unlock prompts; no owner keyring was probed. Native keyring persistence, native windows, tray behavior and attended installed MSI/repair remain unverified. There was no operator-desktop focus acquisition, global input injection, manual walkthrough or administrator prompt during this audit.

## Findings and functional disposition

| Issue | Delivered source behavior | Original evidence and remaining boundary |
| --- | --- | --- |
| [#445](https://github.com/shruggietech/go-schedule/issues/445) | 17 acknowledgement paths honor explicit JSON, retain human text and return write failures without replay | Failing-first observations retained; 68 mode/output regressions and full initial catalogue execution passed |
| [#446](https://github.com/shruggietech/go-schedule/issues/446) | Unexpected daemon operands exit 2 before startup/configuration/state creation | Original owned startup failure retained; packaged integration regression and both native supplemental refusals passed |
| [#447](https://github.com/shruggietech/go-schedule/issues/447) | Evaluator enforces selected JSON acknowledgement and exact SMTP task acknowledgement | All independent SMTP task/history/delivery/receiver checks retained |
| [#448](https://github.com/shruggietech/go-schedule/issues/448) | Missing-desktop guidance uses maintained make build-gui target and sibling location | Both original native obsolete-pin refusals and both fresh corrections retained; no UI launched |
| [#449](https://github.com/shruggietech/go-schedule/issues/449) | Pairing help explains stdin/EOF; timestamp errors give example; known cycle refusals give actual CLI corrections | Three semantic findings in original failed report; fresh native CLI/component checks pass; arbitrary server content still excluded |
| [#450](https://github.com/shruggietech/go-schedule/issues/450) | Wrong-identity fixture retrieves its owned credential and requires actual manifest identity refusal | Original two credential-lookup cases do not establish identity coverage; fresh Windows/Linux TLS component checks reach and reject the identity |

All six source issues remain open pending owner integration and will be closing references on the official PR. Parent #444 remains open for actual v2.0.2 publication. No other actionable product defect was substantiated by this audit. Missing environment prerequisites are explicit coverage limits, not newly invented product defects. Closed historical issues are retained as delivered; original #355/#356 and owner-deferred clustered execution remain open and independent.

## Evidence and privacy

Raw reports, excerpts, daemon logs, system snapshots, environment selections, local paths and pairing materials remain under ignored private S139 evidence. Public derivatives contain only sanitized descriptions, relative repository paths, synthetic screenshots and content digests. Original failures were preserved. Pre-finalization checks found zero drift across 2682 frozen source files and all 808 protected historical inputs. Final source and protected-input checks are recorded in verification. No raw private capture, credential, crash dump or machine identity is included in the PR.
