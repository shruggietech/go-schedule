# Quickstart: Desktop Interaction and Appearance Polish

## Prerequisites

- Go 1.25 toolchain and the repository dependencies
- Git Bash at `C:\Program Files\Git\bin\bash.exe` for canonical verification
- Windows 11 desktop access for native release-candidate observations

## Focused headless validation

```powershell
go test ./gui -run "Appearance|Theme|Interaction|Options|Storage|Navigation|Scroll|Info|Shutdown" -count=1
```

Expected: all preference, font, contrast, selector, compact-row, navigation,
scroll, Info, and shutdown regressions pass.

```powershell
go test -race ./gui -run "AppearancePreferences|ScrollPreference|ScrollDelta|Shutdown" -count=1
```

Expected: pure preference, scroll calculation, and shutdown coordination pass
under the race detector.

## Complete headless GUI validation

```powershell
go test ./gui -count=1
```

Expected: the entire GUI package passes without opening a native window.

## Canonical repository validation

```powershell
& 'C:\Program Files\Git\bin\bash.exe' scripts/verify.sh all
```

Expected: format prints no files; vet, lint, race, GUI, coverage, docs, and
automation all pass.

## Native Windows release-candidate checks

Use the exact candidate produced by the release-readiness flow:

1. Start from clean/reset preferences. Confirm Dark and System are selected.
2. In dark and light modes, hover and keyboard-focus the active navigation
   destination, Save/Okay, ordinary buttons, selectors, links, and Exit. Labels
   and glyphs remain readable and selection/focus remains visible.
3. Switch through Geist (brand), Inter, Ubuntu, Monospace, and System. Confirm
   existing content and an open task dialog update together, then restart and
   verify persistence.
4. Compare Info version and attribution with nearby title/link text at 100% and
   125% or 150% scaling. Resize, minimize/restore, and reopen Info.
5. Inspect Options at 1280 by 800 and the supported small clamp. Storage rows
   remain aligned with no horizontal scrollbar; unavailable rows are wholly
   muted; available Copy actions copy exact paths.
6. Confirm selector menus omit their current value and show it again only after
   a different selection makes it an alternative.
7. Use a conventional mouse wheel on Options at 1x, default 2x, and 4x. Each
   detent moves immediately and higher settings remain controllable.
8. If a precision touchpad is available, confirm small gestures remain precise.
   Otherwise record the hardware as unavailable.
9. Confirm the rail boundary is visible, labels retain breathing room, Exit is
   bottom-right with a semantic glyph, and repeated close requests exit once.

Record exact installer identity, Windows build, DPI settings, devices, and
observations. Do not convert unavailable hardware or headless results into a
native pass.
