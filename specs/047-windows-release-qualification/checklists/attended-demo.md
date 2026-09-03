# S047 Attended Demo Checklist

This checklist is intentionally incomplete until a maintainer tests the exact
local demo linked in `verification.md`. Record failures before requesting any
correction. A local-demo pass reduces pre-push risk but does not replace the
formal post-merge exact-candidate evidence bundle.

## Artifact identity and environment

- [ ] Confirm the MSI filename, byte size, SHA-256, source commit, embedded
  version, ProductVersion, and ProductCode match `verification.md`.
- [ ] Record Windows edition/build, user role, integrity level, display scaling,
  effective DPI, screen resolution, and whether a precision touchpad is
  physically available.
- [ ] Use the installed LocalSystem service and ordinary intended-user GUI for
  all routine desktop observations.

## Installer and core smoke

- [ ] Install seamlessly and verify selected desktop/Start Menu shortcuts plus
  independent Finish actions.
- [ ] Confirm startup is conservatively sized and produces no error spam.
- [ ] Complete a real task success, nonzero exit, and start-failure smoke check.

## Desktop scenarios

- [ ] `desktop.appearance-standard`: at 96 DPI, complete both palettes and all
  five fonts, verify System default/reset and persistence, and capture sharp,
  centered, unclipped text after resize, minimize/restore, and reopen.
- [ ] `desktop.appearance-scaled`: repeat the appearance checks above 96 DPI and
  record the exact effective DPI.
- [ ] `desktop.interaction-states`: exercise every listed control family at
  rest, hover, focus, pressed, selected, and disabled in both palettes; record
  contrast results and native screenshots.
- [ ] `desktop.navigation-options`: at 1280x800 and 800x600, verify rail order,
  spacing/boundary, Exit, storage rows, Copy, selectors, and no horizontal
  scrollbar.
- [ ] `desktop.scroll-input`: exercise every listed surface at 1x, 2x, and 4x;
  verify immediate persistence, keyboard behavior, and touchpad precision or
  record the specific hardware-unavailability reason.
- [ ] `desktop.tasks-table`: populate at least 100 rows and verify both palettes,
  both sizes, headers, row states, full-value disclosure, refresh identity,
  removal, toolbar/double-click actions, and no horizontal scrollbar.
- [ ] `desktop.schedule-activity-tables`: populate at least 100 rows in each view
  and verify both palettes, both sizes, exact state/severity sets, row states,
  full-value disclosure, refresh identity, details, calendar/range/filter/
  clear/acknowledge actions, and no horizontal scrollbar.

## Removal and disposition

- [ ] Verify guided preserve, cancel, and explicitly confirmed wipe flows where
  safely practical, including owned-data results and protected state.
- [ ] Record every pass, failure, and unavailable hardware check in
  `verification.md`; do not infer results from headless tests.
- [ ] If any failure is reproducible, add a failing regression first, correct
  only that defect, build a new hash-bound demo, and repeat affected checks.
- [ ] Reconcile each linked issue independently without claiming formal
  exact-candidate closure, then authorize the final pre-push verification run.
