# Attended Demo Checklist: Windows Demo Qualification

**Purpose**: Condense the pre-PR native checks that require the operator's Windows desktop **Created**: 2026-09-03 **Feature**: [spec.md](../spec.md)

## Candidate binding

- [x] DEMO001 **BOUND BY HANDOFF** The operator installed the one linked local artifact. The recorded SHA-256 is `ef4a869af0e6971d445c53e8f7a237df655245287c7ae64d8e719941bda8ad59`; an independent operator-side hash was not reported.
- [x] DEMO002 **PARTIAL** Windows on the previously reported QHD ultrawide was used. Exact Windows build, scale percentage, and integrity role were not recorded, so formal candidate proof remains pending.

## Setup and first launch

- [x] DEMO003 **PARTIAL PASS** The Start Menu entry was created and worked; the default/on-off matrix was not repeated.
- [x] DEMO004 **PARTIAL PASS** The selected Desktop shortcut was created and worked; the default/on-off matrix was not repeated.
- [x] DEMO005 **PASS** Both created shortcuts launched the GUI.
- [x] DEMO006 **PASS** The selected application and documentation Finish actions opened together.
- [ ] DEMO007 First launch is restored, conservatively sized, and fully reachable.
- [ ] DEMO008 No recurring error dialogs appear during two minutes of normal use.

## Desktop UX

- [x] DEMO009 **PARTIAL FAIL** Spacing and bottom Exit placement were accepted, but the rail lacks a boundary separator and active-tab hover text becomes unreadable in dark and light modes. Tracked by #104 and #109.
- [x] DEMO010 **PARTIAL FAIL** Dark, Light, and System font switching were used. System text was sharp, while the Brand default remained fuzzy; selectors also present the current value as an apparent action. Tracked by #101 and #106.
- [x] DEMO011 **FAIL** Ordinary Brand-font body text, including the affected Info treatment, remained fuzzy. Switching to System corrected it. Tracked by #101 and #106.
- [x] DEMO012 **PARTIAL FAIL** Storage information was present, but the cards are excessively wide, unavailable state is applied only to Copy, and compact tabular presentation is required. Tracked by #106.
- [x] DEMO013 **PASS** Double-click opened the correct task editor. Recorded on closed issue #103.
- [x] DEMO014 **PARTIAL PASS** Exit performed the expected shutdown. Title-bar close was not separately reported; the Exit glyph refinement is tracked by
  #105.

## Task execution

- [ ] DEMO015 Manual task execution succeeds with expected history/output.
- [ ] DEMO016 Scheduled task execution succeeds with expected history/output.
- [ ] DEMO017 Nonzero exit and process-start failure remain distinguishable in run history.

The operator intentionally ended attended testing without running these task cases. Headless engine, executor, API, store, schedule, and integration suites passed; the installed LocalSystem and exact-candidate cases remain formal #94 work rather than being inferred from the headless result.

## Maintenance and removal

- [x] DEMO018 **PASS** Windows Settings opened the guided uninstall and presented wipe choices.
- [ ] DEMO019 Cancel removes nothing.
- [ ] DEMO020 Preserve removal keeps populated machine and per-user data across reinstall.
- [x] DEMO021 **PASS WITH POST-AUDIT** The wipe completed, and a read-only audit found no machine root or app-specific Fyne root in any registered profile.
- [ ] DEMO022 Unrelated sentinels outside owned roots remain untouched.

## Result

- [x] DEMO023 **PASS** No product correction was committed. Every out-of-scope finding was filed or attached to GitHub issues #101, #104, #105, #106, #109,
  #110, #111, #112, and #113.
- [x] DEMO024 **PASS** The operator explicitly ended the walkthrough and authorized work through the final-review boundary.
