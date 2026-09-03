# Verification: Windows Release Candidate Gate

**Date**: 2026-09-02

**Branch**: `codex/040-windows-release-candidate-gate`

## Result

S040's implementation and repository-level verification are complete. The gate, collector, fixtures, workflow controls, and operator documentation are implemented. This result does not claim that a real Windows 11 release candidate has passed attended acceptance.

GitHub issues #94, #98, and coordinator #96 were confirmed open on 2026-09-02. They remain open until a byte-identical staged MSI produces complete reviewed attended evidence and each issue's acceptance criteria are audited.

## Focused verification

| Check | Result |
| --- | --- |
| `go test -race ./internal/releasegate ./scripts/windows-release-gate ./test/integration` | PASS |
| `go test ./gui/...` | PASS |
| PowerShell parser check for `Invoke-ReleaseCandidateAttended.ps1` and `inspect-installer.ps1` under PowerShell 7.6.5 | PASS |
| Compilation of the collector's embedded Win32 C# adapter with `Add-Type` | PASS |
| `test/scripts/automation-check_test.sh` positive and negative workflow fixtures | PASS (`automation-check-test: OK (all)`) |
| Candidate-only CLI validation against the inert fixture manifest/MSI | PASS |
| Production CLI validation against `evidence_class: automated-fixture` | Expected block, exit 1 with only the evidence-class defect |
| Strict UTF-8 without BOM scan across all changed text files | PASS |
| Mojibake scan across all changed text files | PASS |
| `git diff --check` | PASS |

The integration suite includes the WiX source, installer-inspector, release-staging, promotion-order, canonical-scenario, and hidden-child-process contracts. No non-fixture compiled MSI is present in this checkout and WiX is not installed on this host, so no new local compiled-MSI inspection was available. S039's compiled-MSI checks remain unchanged; S040's real compiled-MSI proof is intentionally deferred to exact-candidate staging and attended validation.

The repository's Fyne test driver is not part of the canonical race package set because its existing asynchronous test renderer reports framework cache races. The new GUI evidence tests pass in the canonical GUI gate; all new non-GUI release-gate and integration code passed race detection.

## Canonical verification

Command:

```text
C:\Program Files\Git\bin\bash.exe scripts/verify.sh all
```

All eight gates passed:

1. `format`
2. `vet`
3. `lint` (`0 issues.`)
4. `race`
5. `gui`
6. `coverage` (engine 86.4%, schedule 89.2%, timezone 91.3%, store 84.1%, catchup 88.9%, logbus 91.1%)
7. `docs`
8. `automation` (`release staging/promotion` included)

## Independent review disposition

Independent reviews covered project directives, architecture and correctness, Windows process and filesystem safety, Git and issue history, and adversarial fail-closure. Confirmed findings were remediated before canonical verification. Notable corrections included:

- a production-only attended evidence class distinct from automated fixtures;
- semantic parsing and cross-checking of native window and task-run attachments;
- process SID, integrity RID, session, HWND, DPI, rectangle, Fyne, and LocalSystem binding;
- duplicate-key, invalid UTF-8, fractional-integer, linked-parent, overwrite, traversal, and archive safety;
- hidden noninteractive child processes with concurrent output draining;
- explicit invalid-input, rollback, retry-duration, Exit reachability, preserve/wipe, and reinstall assertions;
- exact workflow identity, draft retry behavior, asset allowlisting, final checksum ordering, and last-moment tag verification.

No product correction for the error-spam issue was committed because S040 did not obtain the issue-mandated real baseline reproduction and uncommitted-fix native proof. The implemented gate makes that absence blocking rather than manufacturing a pass.
