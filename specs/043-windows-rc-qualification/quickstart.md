# Quickstart: Test the S043 Demo MSI

Use only the MSI linked in the S043 handoff and confirm its SHA-256 before
installing. This is a local demo, not an official release candidate.

## 1. Record the environment

- Windows edition and build
- display resolution and scaling percentage
- whether prior go-schedule data/preferences are present
- whether the current account is standard or elevated

Do not include account names or other personal data in screenshots.

## 2. Install and launch

- Confirm Start Menu is selected by default and can be deselected.
- Confirm Desktop shortcut is optional and both selected shortcuts work.
- Exercise the Finish-page app and documentation choices independently and
  together; confirm silent operations do not launch either.
- Confirm the GUI opens restored with comfortable margins, reachable title bar
  and taskbar, and no repeating error dialogs.

## 3. Exercise the desktop UI

- Confirm the navigation rail has comfortable spacing.
- Confirm Options is directly above Info and Exit is isolated at bottom-right.
- Switch Dark, Light, and System modes and each font; reopen to confirm
  persistence, then restore defaults.
- Compare Info version and attribution sharpness with nearby text.
- Select and copy storage paths; confirm scope and uninstall wording are clear.
- Double-click a task before and after a list refresh; exactly one correct editor
  should open.
- Confirm Exit and the title-bar close action both shut down cleanly.

## 4. Exercise real work

- Create and run one simple task manually.
- Create a second short-interval task and observe a scheduled success.
- Record run history, output, exit codes, and an external marker where practical.
- If a task fails, retain the task definition and Activity/run details rather
  than guessing the cause.

## 5. Exercise removal

- From Windows Settings, choose Uninstall and confirm maintenance opens.
- Cancel once and confirm nothing is removed.
- Choose preserve removal with populated data, reinstall, and confirm tasks
  remain available.
- Choose wipe removal only after preserving any wanted test data; confirm its
  warning and verify machine data plus per-user preferences are removed.
- Confirm unrelated files outside product-owned locations remain untouched.

Report each item as pass, fail, partial, unavailable, or not run. Include the
demo SHA-256 in the response so results cannot be attached to the wrong rebuild.

