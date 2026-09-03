# Attended Demo Checklist: Windows Demo Qualification

**Purpose**: Condense the pre-PR native checks that require the operator's Windows desktop
**Created**: 2026-09-03
**Feature**: [spec.md](../spec.md)

## Candidate binding

- [ ] DEMO001 Confirm the downloaded/local file SHA-256 matches the handoff.
- [ ] DEMO002 Record Windows build, display resolution, scaling, and account integrity role.

## Setup and first launch

- [ ] DEMO003 Start Menu defaults on and can be turned off.
- [ ] DEMO004 Desktop shortcut defaults off and can be turned on.
- [ ] DEMO005 Created shortcuts launch the GUI.
- [ ] DEMO006 Finish actions work independently and together.
- [ ] DEMO007 First launch is restored, conservatively sized, and fully reachable.
- [ ] DEMO008 No recurring error dialogs appear during two minutes of normal use.

## Desktop UX

- [ ] DEMO009 Navigation spacing, order, selected/focus states, and bottom Exit placement are correct.
- [ ] DEMO010 Dark, Light, System, font selection, persistence, and reset work.
- [ ] DEMO011 Info version and attribution text are sharp, centered, and unclipped.
- [ ] DEMO012 Storage paths are selectable, accurate, and use honest lifecycle wording.
- [ ] DEMO013 Task double-click opens exactly one correct editor before and after refresh.
- [ ] DEMO014 Exit and title-bar close each perform orderly one-shot shutdown.

## Task execution

- [ ] DEMO015 Manual task execution succeeds with expected history/output.
- [ ] DEMO016 Scheduled task execution succeeds with expected history/output.
- [ ] DEMO017 Nonzero exit and process-start failure remain distinguishable in run history.

## Maintenance and removal

- [ ] DEMO018 Windows Settings Uninstall opens the guided maintenance path.
- [ ] DEMO019 Cancel removes nothing.
- [ ] DEMO020 Preserve removal keeps populated machine and per-user data across reinstall.
- [ ] DEMO021 Wipe removal warns clearly and removes owned machine and per-user data.
- [ ] DEMO022 Unrelated sentinels outside owned roots remain untouched.

## Result

- [ ] DEMO023 Every failure is described before any related correction is committed.
- [ ] DEMO024 The operator explicitly states when demo testing is complete.

