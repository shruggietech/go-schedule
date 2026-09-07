# Verification: Wails Release Cutover

**Branch**: `codex/066-wails-cutover`

**Date**: 2026-09-07

## Chronological local evidence

1. The new integration contracts were run before implementation and failed against the shipping Fyne path, missing Wails release payload, retained legacy directories and dependencies, and stale current-product guidance.
2. `go test ./test/integration -run "DesktopCutover|WindowsInstallerGUIResourceContract" -count=1` passed after the release, CI, dependency, removal, naming, and documentation changes.
3. `sh test/scripts/automation-check_test.sh automation` passed after its Wails production, WebKitGTK, four-boundary Dependabot, and retired-path fixtures were updated.
4. `sh scripts/verify.sh gui` passed. Desktop Go packages passed under the race detector, Wails v2.14.0 built `desktop/build/bin/gosched-gui.exe`, 18 frontend files containing 60 tests passed, and the production frontend bundle compiled.
5. `npm audit --audit-level=high` reported zero vulnerabilities.
6. `npm run test:e2e` passed 16 Chromium tests covering accessibility, keyboard use, reduced motion, local assets, responsive behavior, 100-row workspaces, and 80, 100, 150, and 200 percent zoom.
7. `go test ./...` passed for the complete root module, including integration and release-contract packages.
8. `go run ./scripts/brand-check` passed with 108 artifacts, 28 SVGs, and 58 current consumers after retired desktop mappings were removed.
9. `pwsh -NoProfile -File build/windows/verify_wxs.ps1 -StageDir desktop/build/bin` passed against the Wails executable and freshly built daemon, CLI, and cleanup companions.
10. The final artifact audit found that the existing WiX source omitted the README, license, and changelog copied into the Windows staging directory. S066 added each as an installed component plus a failure-closed contract. WiX 6.0.2 then compiled `go-schedule_S066_windows_amd64.msi` from the corrected Wails payload. `test/windows/inspect-installer.ps1` accepted it as local-demo evidence with SHA-256 `4c2a0842ea0272d585f6295b3d5a59abbf60c7cb845e5d3a8678e8f58ae8e990` and proved the package subject, canonical identity, icon, PATH, and local-group rows.
11. The first canonical run passed format, vet, lint, race, GUI, coverage, and documentation, then stopped because the lifecycle inventory said In Progress while completed task checkboxes had not yet been reconciled. The task record was corrected and `sh scripts/automation-check.sh .` passed independently.
12. `sh scripts/verify.sh all` was rerun from the corrected state and passed all eight gates in order: format, vet, lint, race, gui, coverage, docs, and automation. Coverage was engine 82.8 percent, schedule 89.1 percent, timezone 91.3 percent, store 80.2 percent, catchup 88.9 percent, and logbus 91.1 percent.
13. Final read-only spec-kit analysis found 20 requirements, 8 success criteria, 35 completed tasks, no unresolved markers, no lifecycle conflict across 66 specifications, and no constitution conflict or unexplained requirement gap.
14. After the Windows documentation payload and hosted staging contract were added, `sh scripts/verify.sh all` was run once more from the complete state and passed all eight gates, including the full automation mutation matrix.

## Hosted evidence boundary

Pull-request CI must still provide same-revision Windows, macOS, and Linux Wails builds, Chromium accessibility and responsive coverage, root race and coverage gates, CodeQL, Windows compiled MSI and silent lifecycle checks, and release-policy checks. No public tag, draft release, release publication, or milestone closure is authorized by this verification record.
