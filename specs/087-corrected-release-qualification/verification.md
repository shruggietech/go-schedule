# S087 verification

## 2026-09-15: Candidate refresh

The obsolete unpublished draft was backed up before mutation. All eight downloaded assets matched their GitHub lengths and SHA-256 digests. Release metadata, annotated tag metadata, and the assets are retained locally in `dist/s087/obsolete-candidate/`. The original annotated tag object `db8559876d4ebf8020736968eb32b765adea0959` is preserved at local reference `refs/s087-backup/v1.4.0`.

The `v1.4.0` tag was replaced with reviewed S086 main commit `951d864d9683ec3bdcb1e37535ecddd67738bf13` using an exact old-object force-with-lease. This triggered [Release staging run 34968064596](https://github.com/shruggietech/go-schedule/actions/runs/34968064596), attempt 1. All three preflight jobs and the daemon/CLI staging job passed. GUI staging remains in progress at this checkpoint. No public promotion was dispatched.

Issue #226 was reopened because its exact-candidate qualification criteria are unfinished. Issue #228 remains open. No native qualification pass or release completion is claimed.

## 2026-09-15: Offline prerequisites

The official PowerShell 7.6.6 x64 portable ZIP is 106328873 bytes and has SHA-256 `02fe458be20493fbdf43f61ea20610b811ee6c738ab1676c61b9cfcd1a33c860`, matching the upstream release's published hash listing.

The Microsoft WebView2 x64 standalone installer is 212745424 bytes and has SHA-256 `ebebc5ec130378ff1ab513f3917be791a9cf84f849e970b1695ff01801a9d348`. Windows Authenticode verification returned `Valid`, signed by Microsoft Corporation. The resolved source is `https://msedge.sf.dl.delivery.mp.microsoft.com/filestreamingservice/files/5bdafa0b-8f06-4cab-b664-252bfeca6c75/MicrosoftEdgeWebView2RuntimeInstallerX64.exe`. The small online bootstrapper returned by a different Microsoft link was rejected for the offline package.

Native Windows window discovery through the installed computer-use runtime succeeded. No UI qualification pass is claimed.

## 2026-09-15: Verified staging and prepared sessions

Release run 34968064596 completed successfully with all seven jobs passing at commit `951d864d9683ec3bdcb1e37535ecddd67738bf13`. The remote annotated tag object is `524eeaa5c0a8a4a26f18c20207fefd82ce73309f`, peeled to that same commit. The draft retains exactly eight assets. All eight downloaded assets match their GitHub API byte lengths and SHA-256 digests.

`windows-release-gate verify-candidate` passed with the explicit repository, tag and reviewed commit. Candidate ProductVersion is `1.4.0`, ProductCode is `{7E38AA57-B909-43A6-8F08-0B624C77040F}`, byte length is 18907136, and SHA-256 is `e40baf062079a0d023c429c9c1677c40e43c6eb3806f93acea018f1923e06bee`. The manifest identifies staging run 34968064596, attempt 1, and has SHA-256 `67ae6c59dc8e92e21a8944592bebfd365682d0623ae92b9d31e021c6546a3f5b`.

The independently retained public v1.1.1 baseline is 24338432 bytes with expected SHA-256 `f7ac8f56f28330b016eb6e505e424b19e9bfbe435591cfbc54a723c91ac8e567`. `windows-qualification-session` successfully validated all four role inputs and created `dist/s087/sessions/` with separate fresh and upgrade export folders, offline networking, verified packaged helper copies, and the unchanged collector. The fresh session was launched through the command runner with a hidden parent launch. Installation observations remain pending.

The fresh guest completed PowerShell verification, collector initialization, and offline WebView2 installation. WebView2 returned exit code 0 after 54.2249029 seconds. Phase records and streaming logs are retained in `dist/s087/sessions/fresh-exports/`. The supported Windows control runtime selected and activated the returned Sandbox window and observed the corrected-candidate installation instruction dialog. Execution paused before the Windows UI installation action for the computer-use skill's mandatory action-time confirmation. No candidate installation pass, native matrix completion, PR publication, or merge readiness is claimed.

The maintainer approved the installation step. On resumption, the MSI was already running. Native control advanced the welcome, project license, default feature choices, and Install pages. Start Menu shortcut was selected and desktop shortcut unselected by default. The initial client-side package security-policy check took approximately 121 seconds; the service-side check introduced a second pause after Install. `dist/s087/screenshots/001-candidate-installing.png` retains the native Sandbox raster. These are focused observed facts, not complete setup or medium-integrity release-gate evidence.

## 2026-09-15: Fresh installation and qualification blocker

The candidate installation completed with exit code 0 and no reboot. The mechanical phase ran from `2026-09-15T12:28:40.2313668Z` to `2026-09-15T12:38:16.6349377Z` (576.4056414 seconds, including operator wizard time). Offline WebView2 had already completed with exit code 0. The Finish page reported successful setup with Launch go-schedule selected and documentation unselected. The production Wails application launched, displayed v1.4.0, and connected to the installed daemon.

The native task editor opened as a bounded modal, kept advanced settings collapsed, and presented Insert example as a secondary button. Cancel closed it and restored focus to Create task. Timing mode still appeared taller than neighboring selectors. Retained local rasters are `dist/s087/screenshots/002-candidate-completed.png`, `003-system-tasks.png`, and `004-task-modal.png`.

A new real-browser regression in `desktop/frontend/e2e/tasks.spec.ts` confirmed a 10.296875 px difference between Group, Timing mode, and Schedule syntax at 1440 by 900 before the CSS correction. Shared `.field` intrinsic grid alignment corrected the difference; all three task Playwright tests passed afterward. The fix and CSS skills guided test-first root-cause repair rather than an isolated fixed-height override. Issue #231 is reopened for this existing acceptance-criteria failure, not replaced by a duplicate.

The entire canonical eight-gate local verification passed before the repair (including 24 frontend files and 128 tests, native Wails build, core coverage, docs, and automation). Those results do not verify the subsequent repair; its separate aggregate run is required before publication.

Full qualification remains blocked. Upgrade was not executed; no native task run or deletion was performed; required normal-user, multiple-profile, display measurements, and other native matrix observations remain unavailable. No collector Finalize or complete native gate pass is claimed. T007 and T008 remain incomplete. The draft still contains the reviewed pre-repair bytes, not this unreviewed branch's CSS. A repair/evidence PR must be reviewed and merged before a new reviewed candidate can be staged and qualified. This publication-order deviation avoids claiming an incomplete release passed and does not relax any native or local test gate. The PR uses references, not issue-closing keywords, for #226, #228, and #231.

## 2026-09-15: Repair verification

The complete Playwright suite passed with 28 tests. A separate foreground `sh scripts/verify.sh all` run after the repair exited 0: format, vet, lint (zero issues), root race tests, desktop race tests, production Wails build, 24 frontend files and 128 tests, frontend production bundle, core coverage, documentation (20 pages plus policy/architecture fixtures), and automation including mutation fixtures all passed. Final output was `automation-check-test: OK (automation)`.

Core coverage remained engine 82.9%, schedule 89.1%, timezone 91.3%, store 80.1%, catchup 88.9%, and logbus 91.1%. `git diff --check` passed. These establish repair source verification, not a complete native qualification pass or public release. Issue #231's finding is retained in [its progress comment](https://github.com/shruggietech/go-schedule/issues/231#issuecomment-5680378658); #228's child index now correctly shows #231 incomplete.
