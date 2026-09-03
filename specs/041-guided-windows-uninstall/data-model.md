# Data Model: Guided Windows Uninstall Entry

S041 adds no application database entity or persistence migration. It defines an operating-system registration contract.

## Application-management registration

| Field | Required state | Meaning |
| --- | --- | --- |
| Visible product identity | Exactly one current go-schedule entry | The Windows application list must not show duplicate native or shadow records. |
| `NoRemove` | `1` | The system must not offer the reduced-interface direct attended removal entry. |
| `NoModify` | Absent | The system must retain the full maintenance entry. |
| `ModifyPath` | Present, native Windows Installer maintenance command for the current ProductCode | Opens the package-owned maintenance wizard. |
| `UninstallString` | Absent | Prevents Windows Settings from starting the bypassing direct removal command. |
| ProductCode | Current compiled package identity | Binds maintenance to the installed candidate. |
| Display name/version/publisher/icon | Existing canonical values | Keeps the visible entry recognizable and current. |

## State transitions

```text
Absent
  -> fresh install
Registered with maintenance only
  -> repair
Registered with maintenance only
  -> major upgrade
Registered with new ProductCode and maintenance only
  -> maintenance wizard -> Remove -> preserve or confirmed wipe
Absent
```

Direct unattended `/x` transitions from either registered state to absent without relying on the visible application-management action.

## Validation invariants

- A current install has exactly one visible go-schedule entry.
- `NoRemove` and `ModifyPath` are present together.
- `NoModify` and `UninstallString` are absent.
- Registration changes do not alter `GOSCHEDULE_REMOVE_DATA` or cleanup sequencing.
- Upgrade removes the old visible identity and registers only the current ProductCode.
