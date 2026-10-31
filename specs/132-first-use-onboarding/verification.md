# S132 Verification

## Evaluated scope and provenance

Accepted starting source: owner-merged S131 main 1c36dae70af3c0315827b2bb063f8ba4f6dc743f. Package exercises use actual published v2.0.0 release 430039348 at immutable source 75eae1ae11ad40db38341dc7d63f32876cba6483; the new own Run ID presentation is separate source evidence and is not in those public binaries. No runtime, API, schema, authority or dependency change.

## First failures retained

1. Installed checklist prerequisite required plan.md before the documented pre-plan checklist stage. Paths-only resolution succeeded; unchanged prerequisite passed after setup-plan.
2. Selected-own-versus-source Run ID component regression failed because Run details lacked Run ID. After exposing the existing field, all six ActivityPage tests passed.
3. The first Linux owned-runtime setup used 0700 on the socket directory; the existing 0770 restricted-access gate refused startup before any task. Correcting only the new owned directory permissions allowed the fresh run. Original receipt/log retained; no product/access check was weakened.
4. The first new site heading assertion used exact names without Fumadocs' Copy Anchor Link suffix. Ten existing cases passed; the new locator failed. Its error context/trace were retained separately; level/name matching corrected the assertion and all eleven cases then passed.

## Published-package workload

Windows MSI cabinet was read read-only and extracted to a fresh owned path without installation, UAC or service registration. Linux binaries came from the inspected public archive. Windows and Linux each used one task with a space-containing working path, explicit native shell and UTC schedule. Automatic scheduling was disabled during the deliberate failure. Exact terminal records were success (exit 0), missing-input failure (exit 2 and explanatory output) and success after restoration (exit 0). Snapshot content matched input, prior snapshot survived failure, and history retained all three independent IDs. Explicit named-pipe/Unix-socket endpoints and non-empty restricted OS groups remained enabled. All owned daemon processes exited; the installed services were unchanged.

The guide PowerShell setup and payload also passed independently, with only the workspace redirected to an owned private path. Windows execution used the current owner process, not LocalSystem. Linux used an owned WSL process, not an installed systemd service. No attended installation, native GUI or local macOS execution is claimed. Existing standing native-testing waiver remains.

## Targeted source verification

All 168 frontend cases across 28 files passed, including the failing-first own/source run distinction and selection changes. All 35 existing desktop browser cases passed in 27.0 seconds, including selected Run ID/output, accessibility and reflow. Seven site unit tests and the complete 25-document static export/link/metadata/product-policy checks passed. All eleven site browser cases passed in 8.0 seconds after the locator correction; the same first-use journey adds a focused 390px overflow assertion and retained mobile capture.

## Canonical verification

Full foreground sh scripts/verify.sh all passed all eight gates in 1108.337 seconds, from 2026-10-31T21:58:22.446Z to 2026-10-31T22:16:50.784Z. The frozen evaluation covered 2570 source paths with zero source drift and all 255 protected inputs unchanged. Evaluated source-manifest SHA256: 45531b0f36cc040fc39e3223e1cb5ffd6dde495bbd3b30f1bab145dcefe5235b. Existing Go cache behavior is preserved; no forced uncached native-runtime campaign is claimed. Only five factual spec/status records were finalized afterward; evaluated product, guide, browser-test, dependency, workflow and tagged release bytes remain unchanged. Final format/lifecycle/immutable-input checks apply to those factual updates.
