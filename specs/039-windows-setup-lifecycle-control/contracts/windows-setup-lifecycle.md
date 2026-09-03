# Contract: Windows Setup Lifecycle

## Shortcut feature contract

| Feature | Fresh default | Managed install | Maintenance | Upgrade |
| --- | --- | --- | --- | --- |
| Core `Main` | installed, cannot be absent | always local | not removable independently | retained |
| `StartMenuShortcut` | installed | standard feature selection | add/remove | S039-and-later state migrates |
| `DesktopShortcut` | absent | standard feature selection | add/remove | S039-and-later state migrates |

The first upgrade from a pre-S039 MSI uses the new defaults because no earlier matching optional-feature identities exist.

## Completion action contract

| Flow | Launch application default | Open docs default | May execute actions? |
| --- | --- | --- | --- |
| Fresh full-UI success | selected | unselected | Yes, independently on Finish |
| Silent/basic install | property ignored | property ignored | No |
| Repair/modify/upgrade | not offered | not offered | No |
| Uninstall/rollback/failure/cancel | not offered | not offered | No |

The application target is the installed `gosched-gui.exe`. The documentation target is `https://shruggietech.github.io/go-schedule/`. Both are delegated to the interactive user's shell without installer elevation.

## Removal contract

| Invocation | Property | Outcome |
| --- | --- | --- |
| Full UI remove, default | `0` | Remove software, preserve application data |
| Full UI remove, confirmed wipe | `1` set only by confirmation | Remove software and wipe declared local application roots after successful commit |
| Silent remove | absent or `0` | Preserve |
| Silent remove | exact `1` | Wipe |
| Any operation | other explicit value | Fail validation before execute sequence |
| Upgrade, repair, reinstall | any | Never wipe |

All destructive scheduling uses the equivalent of:

```text
Installed AND REMOVE~="ALL" AND NOT UPGRADINGPRODUCTCODE
AND NOT REINSTALL AND GOSCHEDULE_REMOVE_DATA="1"
```

## Cleanup helper contract

The embedded helper accepts exactly one public operation, `wipe`. It accepts no target path. All owned roots and its protected result location are derived internally.

1. Discover the machine root and every safely registered accessible local-profile preference root.
2. Canonicalize and validate every existing target before deleting any target.
3. On any unsafe target, record `refused` and delete nothing.
4. After complete preflight, remove each root without following reparse entries and atomically record progress.
5. Record `complete` only when no declared root remains; otherwise record `partial` or `internal-error`.

The MSI ignores the commit helper's process return for transaction rollback purposes. Hosted lifecycle automation consumes the protected result, and attended completion guidance identifies its stable location when cleanup is incomplete; an MSI success code proves software removal, not unconditional cleanup completion.

## Removal inventory contract

Always removed: installed binaries, service registration, machine PATH integration, product registration, installer registry markers, and selected installer-created shortcuts.

Preserved by default: `C:\ProgramData\goschedule` and each safely registered local profile's `AppData\Roaming\fyne\tech.shruggie.goschedule` leaf.

Removed by confirmed wipe: the same declared local application roots after safety validation.

Always preserved separately: `goschedadmin`, its memberships, exports, configured paths outside the default product-owned roots, unrelated shortcuts/files, disconnected or detached profile storage, and any refused unsafe root.

## Evidence contract

Source and compiled-MSI evidence records feature/component relationships, dialog controls/events, custom actions and conditions, secure property behavior, service/GUI close behavior, and absence of completion actions from the execute sequence.

Silent native evidence records MSI hash, Windows identity, command mode, feature selection, shortcut paths, seeded root hashes, cleanup-result state, security-state preservation, and out-of-scope sentinel hashes.

Attended candidate evidence remains #94 and must additionally record visible defaults, four completion combinations, shell handler, process integrity, cancellation, multiple interactive profiles, window bounds, and recurring-error observations.
