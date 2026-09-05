# Contract: Windows Installer Identity

## Canonical asset

The Windows package and GUI artifact use `cmd/gosched-gui/icon.ico`. No second installer-specific copy is permitted.

## MSI relationships

The compiled MSI must contain:

| Table | Relationship |
|---|---|
| Icon | One row whose name is the canonical icon identifier |
| Property | `ARPPRODUCTICON` equals the canonical icon identifier |
| Shortcut | `GuiShortcut.Icon_` equals the canonical icon identifier |
| Environment | The existing machine PATH row remains bound to `Gosched` |

The source definition must express those relationships explicitly. Falling back to the shortcut target executable does not satisfy this contract.

## GUI executable relationship

Before the Windows GUI build, the release pipeline generates a 64-bit Windows resource from the canonical `.ico` into the GUI command package. The GUI build must consume that package in the same job. This contract does not assert that a specific Windows shell cache has refreshed; that belongs to native observation.

## Preserved package behavior

- Per-machine amd64 MSI.
- Existing install directory and three installed executables.
- Existing service registration and lifecycle.
- Existing machine PATH component behavior.
- Existing data-directory preservation on uninstall.
