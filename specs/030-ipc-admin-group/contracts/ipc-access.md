# Contract: Local IPC Access Control

## Configuration contract

| `admin_group` value | Required result |
| --- | --- |
| omitted | Use `goschedadmin` and enter restricted mode. |
| non-empty valid group | Enter restricted mode for the resolved group. |
| non-empty unknown/non-group value | Fail startup; name `admin_group` and the requested value. |
| empty string | Enter explicit compatibility mode and emit a warning. |
| leading/trailing or whitespace-only value | Reject configuration as invalid; do not reinterpret as compatibility mode. |

## Listener contract

The listener is not returned until the selected platform policy has been applied and verified. Success returns immutable evidence:

```text
restricted:    mode=restricted, admin_group=<configured group>
compatibility: mode=compatibility, admin_group=""
```

Any listener created before a later policy error is closed. A Unix socket created during that attempt is removed.

## Logging contract

Restricted startup emits an informational record:

```text
message="IPC access configured"
access_mode="restricted"
admin_group="<configured group>"
```

Compatibility startup emits a warning record:

```text
message="IPC compatibility mode enabled"
access_mode="compatibility"
admin_group=""
```

The access record precedes `daemon startup complete`.

## Windows descriptor contract

Restricted:

```text
D:P(A;;GA;;;SY)(A;;GA;;;BA)(A;;GRGW;;;<resolved group SID>)
```

Compatibility:

```text
D:P(A;;GA;;;SY)(A;;GA;;;BA)(A;;GRGW;;;AU)
```

No configured account name is interpolated into SDDL.

## Unix filesystem contract

Restricted mode:

- the managed default parent or newly created custom parent group is the configured group and mode is exactly `0770`;
- an existing custom parent is left unchanged and must already have that group and exact `0770` mode;
- socket group is the configured group and mode is exactly `0660`.

Compatibility mode:

- the managed default parent mode is exactly `0755`, while an existing custom parent is preserved;
- socket mode is exactly `0666`.

## Windows installer contract

- The MSI creates or reuses local `goschedadmin` before starting `goschedd`.
- The logged-on installing user is added through declarative WiX group membership.
- Repair and upgrade tolerate existing group and membership.
- Uninstall leaves the group and membership intact.
- The release build pins the WiX tool and UI/Util extensions to the same 6.0.2 version.
