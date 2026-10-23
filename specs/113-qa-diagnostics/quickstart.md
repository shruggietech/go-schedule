# S113 validation guide

Use the existing supported Go/C toolchain, PowerShell7, Bash/SQLite and unelevated Windows/WSL user. All console children must use the existing hidden redirected launcher. Tests use disposable files, owned daemons and loopback receivers.

1. Run focused regression-first checks for internal/clievaluate, cmd/gosched-qa, internal/api/client and test/integration payload fixtures. Exercise a failed intermediate request, command success with failed inspection, selected nondefault budgets and forced interruption/recovery. Inspect partial data, failed request, phase timings and original identities.
2. Run migration23->24 and SMTP store/client/remote/desktop/CLI regressions, including saved successful history, partial failure, retry, uncertain acknowledgment, credential redaction and terminal erasure.
3. Exercise all six concurrent writers for both script twins and both initialized cases; retain original native crash/prerequisite outcomes without retries or weaker assertions.
4. Run production site and desktop checks, then foreground sh scripts/verify.sh all and wait for all eight gates. Preserve exact output and meaningful failure evidence.
5. Build fresh six-binary owned QA inputs with a source receipt. Run the complete Windows/WSL catalog through the PowerShell entry point using a supported selected total budget, keeping the default 15-second case and every resource/assertion limit. Review all exact-bound messages, finalize, validate and render the fixed report. Distinguish code2 defect outcome from integrity errors and code3 interruption.
6. Record each issue's actual functional disposition, unconfirmed causes, exact source/binary identities and cleanup. Publication workflow follows the owner authorization, with no third Codex request and no merge/tag/release.
