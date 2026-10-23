# S112 verification

## Baseline and regression evidence, 2026-10-22

Unchanged main 2b827e0 reproduced queued-first/completed-second manual-run masking in an isolated AST-loaded probe function, without installing a service. The unchanged packaged helper passed locally in 2.50 seconds; this does not refute its prior hosted 45-second build timeout.

Executable regressions now cover queued/terminal filtering, skipped outcomes, failure checkpoint retention, missing diagnostic commands, bounded log tail, primary/cleanup distinction, cleanup-only failed status and UTF-8 without BOM. Packaged helper tests retain real executable authentication and setup cancellation. Existing contained process trees own compilation, including descendants, and all recorder adapters retain their existing semantics at caller composition roots.

The original credential correction passed three native Ubuntu-24.04 WSL race repetitions (28.69, 4.73 and 4.70 seconds for packaged rejection). The later contained-build correction passed the focused Windows race run, final canonical checks and a fresh native Linux packaged-helper race test in 12.64 seconds; earlier Linux timings are not represented as containing that later correction.

## Retained failed attempts

- Initial canonical run passed format, vet, lint, root race, desktop/native build, 166 frontend tests and coverage, then found the stale SMTP anchor in docs/cli.md. The link is corrected.
- The second run failed an unchanged concurrent payload-writer fixture with a PowerShell internal CLR crash and a separate SQLite prerequisite refusal. #355 records the unknown cause. Five unchanged focused race repetitions passed; no source fix is claimed.
- A later run passed the earlier gates and website export, then identified an inconsistent S112 lifecycle inventory label. The inventory now matches the active spec state.
- One preliminary automation driver exceeded the tool response lifetime. Its identified orphan process tree was stopped and its untracked fixture directory removed. It is not counted as a passing check.

## Browser and tooling

All 35 desktop Playwright checks and all 10 static website Playwright checks passed using the matching managed Chromium build. Both npm audits found zero vulnerabilities. ShellCheck passes for all test/scripts/*.sh and lib scripts; precise fixture comments identify intentionally literal workflow variables and Markdown. Publication formatting passes with no em dashes or hard-wrapped prose. The supplementary PowerShell structure checker reports the same two pre-existing section-divider/section-heading deviations on unchanged HEAD and the corrected probe; all byte/encoding checks pass. The probe keeps the existing project layout rather than claiming that supplementary structure check passed.

## Final canonical result

The final foreground sh scripts/verify.sh all completed successfully with all eight gates: format, vet, lint, root race, Wails Go/frontend/native build, canonical coverage, docs and automation. The accepted log is verify-all-accepted.log (exit 0). All 166 frontend tests passed. Canonical coverage was engine 83.6%, schedule 89.1%, timezone 91.3%, store 80.2%, catchup 88.9% and logbus 91.1%. Local execution uses GOMAXPROCS=2 and GOFLAGS=-p=2 to bound build/test-package load; concurrent writer counts and all assertions are unchanged. Coverage retains the canonical six-package cross-package measurement, not per-package approximations.

## Owned Windows/WSL QA

The fresh complete capture dfe7cbf4-759e-4ec0-8643-3a58a57b05bd evaluated base 2b827e0a128cff5698b84ff84a83628e31063423 plus the dirty S112 tree, source_sha256 ff465894a3ab9e7a90720f5b7388a5f44d7a40b25bf3bddaf2c8637f61f991af. The six binaries and build receipt stayed fixed. Windows 11 and existing Ubuntu-24.04 WSL user h8rt3rmin8r used the existing contained Windows Jobs and delegated user cgroup. All 1,763 applicable executions passed; three platform exclusions remained explicit. Both JSON and human runtimes in each environment were removed, with no remaining runtime cleanup entries.

All 1,109 corpus messages received exact run/hash/message/occurrence bindings and seven-part rubrics. Review retained 947 unchanged prior contexts after explicit fixture UUID/time/port/digest/PID normalization, with 160 freshly reviewed contexts and two failed delivery-output reviews. Source-specific task, profile, notification, bundle, MCP, system and log variants were inspected, including native Linux run-as, watcher health and log metadata. The exact excerpts and review-basis index remain in the evidence archive.

The canonical final outcome is fail: 1,761 effective passes, two communication failures, zero blocked cases, three exclusions and one P3 finding. #358 records successful SMTP acceptance summaries labelled last_error on both platforms. Functional delivery succeeded; no source fix or all-green QA verdict is claimed. Finalize, validate, render and validation of the copied archive each returned native outcome code 2, confirming a valid failed evaluation rather than processing/integrity failure.

The original default-budget attempt cb3c46fc-4134-4d90-9971-029284d04e6c remains interrupted, with five Windows command/inspection deadline observations in #356 and wrapper/terminal-reporting behavior in #357. Its original checkpoints, case progress and error output are retained. The fresh attempt selected a bounded two-hour total catalog budget through the direct coordinator; the 15-second case deadline, CPU/memory/process/output limits and assertions were unchanged. All five earlier cases passed in the fresh attempt; their cause remains unconfirmed and #356 stays open. The first run's worker/daemon PIDs were absent after termination; unresolved runtime evidence is retained separately and is not represented as a completed cleanup action.

The evidence archive is C:/Users/h8rt3rmin8r/.codex/artifacts/go-schedule/s112/, including qa-run-02/report.json, generated report.md, the complete evidence directories, reviews-02.json, review-basis-02.json, six input binaries and qualified-source.patch plus the exact changed/new source files. Later lifecycle/delivery edits change only documentation metadata; this QA source digest is not represented as a final commit identity.

## Observation boundaries

Native LocalSystem and installer operations remain reserved for disposable hosted CI. Local browser harness, packaged CLI evaluation, source tests and exact-head hosted checks are distinct evidence. No installed service, production SMTP account, account session, DNS, merge, tag or release is changed by local verification.
