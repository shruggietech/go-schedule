# Verification: v1.1.1 Release Recovery

## Specification analysis

The specification defines 12 functional requirements, 5 measurable outcomes, 3 independently testable user stories, 14 completed requirements-quality checks, and 13 chronological repository tasks. Every requirement maps to historical preservation, source preparation, automation regression closure, staging, qualification, promotion, or final-audit evidence. No unresolved ambiguity, duplicate planning record, scope leak, constitution conflict, or missing requirement coverage remains.

## Historical and hosted evidence

GitHub API inspection confirmed that the annotated v1.1.0 tag object `24e116c7436f737aafd136c6f7190d0fce301d8b` still points to commit `d3b47e44c18faab7474ed383021f483488064a06`. The v1.1.0 GitHub release remains an unpublished draft, and no v1.1.1 tag or release exists. Exact-commit main CI run [34005457330](https://github.com/shruggietech/go-schedule/actions/runs/34005457330) completed successfully for corrective merge `d207daf24b53bd351d671d1b4a90221b4dbd151d`.

Issue [#140](https://github.com/shruggietech/go-schedule/issues/140) and milestone #4 now identify v1.1.1 as the authoritative corrected release target. Their acceptance criteria preserve the historical v1.1.0 tag, prohibit publishing its draft, require exact-commit main CI, and retain the existing qualification and no-rebuild promotion controls.

## Repository evidence

The v1.1.1 release note passed the established highlights-only contract with exactly four bullets and one final tagged changelog link. The changelog contains an empty Unreleased boundary, a dated v1.1.1 correction section, the preserved v1.1.0 historical section, and correct comparison references. The README contains exactly one synchronized 1.1.1 health example.

The initial canonical `scripts/verify.sh all` run passed all eight gates through the installed WSL shell with Windows Go binaries explicitly selected: format, vet, lint with zero findings, race, GUI, coverage, documentation, and direct automation validation. Coverage results were engine 81.9 percent, schedule 89.2 percent, timezone 91.3 percent, store 80.1 percent, catchup 88.9 percent, and logbus 91.1 percent. The separate specification lifecycle audit reported all 58 specifications consistent, and `git diff --check` reported no whitespace errors.

Second-round Codex review identified that the approved release-workflow fixture was stale even though direct automation validation passed. The finding reproduced exactly. S058 updated the fixture with the exact-commit CI job and dependency contract, corrected the associated negative cases, and added `test/scripts/automation-check_test.sh automation` to the canonical automation gate. The focused fixture suite then passed. The first complete post-fix run passed format, vet, lint, race, GUI, coverage, and documentation before the automation gate correctly rejected T013 while it remained open. The next run proved that a Windows Go executable cannot consume a fixture rooted in WSL `/tmp`, so the harness now creates its verified disposable tree beneath the repository where both toolchains resolve it consistently.

The final canonical `scripts/verify.sh all` rerun passed all eight gates, including direct automation validation and `automation-check-test: OK (automation)` under the Windows-Go-from-WSL configuration. Coverage remained engine 81.9 percent, schedule 89.2 percent, timezone 91.3 percent, store 80.1 percent, catchup 88.9 percent, and logbus 91.1 percent.

## Publication evidence

Public release operations remain tracked by issue [#140](https://github.com/shruggietech/go-schedule/issues/140) and begin only after the reviewed S058 preparation merge. Their first destructive step must revalidate the unpublished v1.1.0 draft before deleting that draft only; the immutable v1.1.0 Git tag remains untouched.
