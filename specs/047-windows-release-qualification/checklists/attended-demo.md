# S047 Local Demo Record

**Result**: Accepted as a valid release target by the maintainer on 2026-09-03.

This is the concise pre-push usability record for the exact local demo. The assistant collected artifact and environment metadata and translated the maintainer's ordinary walkthrough into evidence. The exhaustive 47-scenario checklist remains mandatory only for the formal post-merge release candidate.

## Identity collected by automation

- [x] MSI: `dist/go-schedule_s047-demo_5151c069_windows_amd64.msi`
- [x] SHA-256: `31feabdb3f61334b7b46a81350297cce07f8b59882bcb80336f26248c260d2f7`
- [x] Source: `5151c0698aaa65021b32ce697f84a6c7469bf3f6`
- [x] Primary display: 3440x1440, 3440x1392 working area, 120 Hz, 100 percent scaling, 96 effective DPI
- [x] Secondary display: 2560x1440, 2560x1392 working area
- [x] Intended-user interactive install and application walkthrough

## Maintainer observations

- [x] Desktop icon and Start Menu entry work correctly.
- [x] Default window size is comfortable on the QHD ultrawide display.
- [x] Navigation-tab and ordinary-button hover treatments remain readable.
- [x] The single command-line editor is roomy and understandable.
- [x] Font selections work; System is sharp, and Light and Dark both look good.
- [x] Scrolling, dropdown choices, and compact application-storage rows work well.
- [x] Tasks run and complete successfully through the installed application.
- [x] Disabling a group prevents an individually active member task from running.
- [x] Guided uninstall with wipe completes cleanly.
- [x] The maintainer explicitly accepts this build as a valid release target.

## Accepted Post-v1 follow-ups

- [x] #118: new task creation should offer an unchecked **Activate task after saving** option instead of enabling every task automatically.
- [x] #119: Schedule and Activity need wider defaults plus user-adjustable, persisted column widths, especially for **When**.
- [x] #120: Tasks should indicate when an individually active task is effectively blocked by a disabled group or ancestor.

These findings were explicitly classified as issue-only follow-ups and do not expand S047 implementation scope. Formal candidate qualification will repeat the exact native scenario matrix rather than inferring formal evidence here.
