# Data Model: Dedicated IPC Administrative Group

S030 changes configuration and ephemeral startup state only. It adds no stored database entity or migration.

## Administrative Group Setting

| Field | Type | Rules |
| --- | --- | --- |
| `admin_group` | string | Defaults to `goschedadmin`; empty selects compatibility mode; a non-empty value must have no leading/trailing whitespace and must resolve to a group on the host. |

## IPC Access Info

| Field | Type | Rules |
| --- | --- | --- |
| Mode | enum | `restricted` or `compatibility`; selected once during listener creation. |
| Admin group | string | Required and non-empty in restricted mode; empty in compatibility mode. |

### State transitions

```text
configuration loaded
  ├─ admin_group empty ────────────────> compatibility policy
  └─ admin_group non-empty
       ├─ group resolved + applied ────> restricted policy
       └─ lookup/apply/verify failure ─> startup failed

policy ready ─> structured policy log ─> daemon readiness ─> serving
```

There is no transition from restricted setup failure to compatibility mode.

## Unix Endpoint Policy

| Object | Restricted owner group | Restricted mode | Compatibility mode |
| --- | --- | --- | --- |
| Managed default or newly created custom parent | Resolved configured group | `0770` | `0755` for the managed default |
| Existing custom parent | Must already match configured group | Must already be `0770` | Preserved |
| Unix socket | Resolved configured group | `0660` | `0666` |

Restricted ownership and modes are read back before the listener is returned. An unsafe existing custom parent fails without mutation. A post-listen failure transitions through listener close and socket removal before returning an error.

## Windows Endpoint Policy

| Trustee | Restricted rights | Compatibility rights |
| --- | --- | --- |
| LocalSystem | Full control | Full control |
| Built-in Administrators | Full control | Full control |
| Configured group SID | Generic read/write | Not applicable |
| Authenticated Users | None | Generic read/write |

The configured group name is resolved to a canonical SID before descriptor construction. Accepted account types are local aliases, domain groups, and well-known groups.

## Installer Provisioning Record

The installer owns no application database record. Its declarative contract is:

- group name: `goschedadmin`
- group scope: local machine
- group creation: create or reuse
- initial member: Windows Installer `LogonUser`
- uninstall: preserve group and membership
- service start dependency: provisioning completes first
