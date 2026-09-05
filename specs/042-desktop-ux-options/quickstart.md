# Quickstart: Validate Desktop UX and Options

## Safety boundary

Headless tests are safe in this repository session. Do not launch the windowed GUI or attended installer through an unverified console path. Native visual validation belongs on the clean Windows 11 exact-candidate environment in #94.

## 1. Run focused headless tests

```powershell
go test ./gui -run "Appearance|Theme|Options|Storage|Navigation|TaskRow|TaskEditor|Shutdown|Info" -count=1
```

Expected:

- all nine appearance combinations and invalid-value fallbacks pass;
- every storage row has accurate scope, existence, and lifecycle copy;
- copying returns the displayed path exactly;
- navigation order, rail width, selection, Activity badge, and Exit placement pass;
- single click opens no editor and double activation retains stable identity;
- competing close requests cancel and close once;
- Info version and attribution labels remain centered and unwrapped.

## 2. Run GUI and race verification

```powershell
go test ./gui -count=1
go test -race ./gui -run "AppearancePreferences|ResolveStorage|StorageExistence|EditorOwnership|ShutdownCoordinator" -count=1
```

Expected: the complete GUI suite passes under the headless driver, and pure concurrent state contracts pass under the race detector.

## 3. Run canonical verification

```powershell
C:\Program Files\Git\bin\bash.exe scripts/verify.sh all
```

Expected: format, vet, lint, race, GUI, coverage, documentation, and automation gates all pass.

## 4. Perform exact-candidate native validation

After CI produces the immutable staged candidate, use the #94 clean Windows 11 runbook at 100 percent and at least one scaled-DPI setting:

1. Confirm the default window remains conservatively bounded.
2. Inspect the Info version and attribution beside adjacent body text in Dark and Light modes.
3. Confirm all destination labels have balanced spacing at 1280 by 800 and at the supported 800 by 600 clamp.
4. Confirm Options immediately precedes Info and Exit remains at the rail's bottom-right at both sizes.
5. Exercise keyboard navigation, focus visibility, all appearance selections, storage-path selection/copy, task-row single and double activation, and both Exit and title-bar close.
6. Record screenshots and the exact candidate identity in the #94 evidence bundle.

This native step is mandatory before the DPI-dependent acceptance criteria in #101, #104, and #105 close. Headless screenshots do not substitute for it.
