# Data Model: Windows MSI Local-Group Recovery

## Administrative Group State

| Field | Rule |
| --- | --- |
| Name | Exactly `goschedadmin`. |
| Scope | Machine-local; effective installer group domain is empty. |
| SID | Stable through repair, reinstall, upgrade, and uninstall. |
| Members | Contains the intended account before service start. |
| Removal | Preserved during uninstall. |

## Installer Operation Evidence

| Field | Rule |
| --- | --- |
| Scenario | `fresh` or `upgrade`. |
| Operation | Install, repair, reinstall, upgrade, uninstall, or cleanup. |
| Exit code | Numeric; only 0 and 3010 are successful. |
| Verbose log | Absolute path and SHA-256. |
| Diagnostics | Access denied, 26421, rollback, local-group action, service order. |

## Phase Observation

| Field | Rule |
| --- | --- |
| Product | Registration present or absent. |
| Group | Existence, name, SID, and normalized members. |
| Service | Existence, status, startup type, and account. |
| Files/PATH | Install-directory state and exact product PATH cardinality. |

## State Transitions

```text
fresh -> candidate -> repair -> reinstall -> uninstall (group preserved)
upgrade -> preprovision -> v0.9.0 -> candidate -> uninstall (group preserved)
```
