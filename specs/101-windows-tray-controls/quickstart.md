# S101 validation guide

## Prerequisites

- Windows build environment with WiX and Wails dependencies for the MSI and desktop GUI.
- A disposable Windows installation where service start/stop and UAC prompts can be exercised.
- Go and frontend dependencies installed as described in the repository.

## Automated checks

Run `scripts/verify.sh all` from the repository root, then build the Windows MSI using the documented release workflow. Check the installer contract tests and staged artifacts include the companion.

## Functional walkthrough

1. Install MSI as an administrator, sign in, and check exactly one go-schedule icon appears. Inspect tooltip and menu.
2. Close GUI. Verify icon remains. Restart Explorer and verify exactly one icon recovers.
3. Stop and start the local service via tray, approving UAC. Check pending and observed results. Repeat with UAC cancelled.
4. Use Open multiple times when GUI is visible, minimized, and hidden. Check exactly one GUI appears and receives focus.
5. Stop the installed service. Open GUI and confirm it remains stopped. Select a remote connection and confirm local state/actions remain independently visible.
6. Upgrade and remove MSI. Verify companion termination, icon cleanup, and no stale logon registration.

Any defects found in verification are tracked as defects against the affected revision, not as extra closure criteria on #255.
