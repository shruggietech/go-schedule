# Data Model: Windows Installer Identity and Verification

No product database or runtime entity changes. This feature models build and
verification evidence only.

## Installer identity contract

| Field | Meaning | Validation |
|---|---|---|
| canonical icon source | Tracked multi-resolution Windows brand asset | Exactly `cmd/gosched-gui/icon.ico` |
| icon identifier | Key used by installer consumers | One stable identifier with `.ico` suffix |
| installed-apps reference | Package primary-icon relationship | Equals icon identifier |
| shortcut reference | Start Menu shortcut icon relationship | Equals icon identifier |
| executable resource source | Release input for native GUI icon | Equals canonical icon source |
| executable resource output | Build-time Windows resource | `resource_windows_amd64.syso` in GUI command package |

## Search-path contract

| Field | Required state |
|---|---|
| scope | Machine-wide |
| value | Install directory |
| position | Append without replacing unrelated entries |
| persistence | Removed with owning CLI component |
| cardinality after install/upgrade | Exactly one |
| cardinality after uninstall | Zero |

## Evidence record

| Field | Description |
|---|---|
| date | Date of observation |
| environment | Windows build and clean/disposable status |
| artifact | Version, origin, path, and checksum when available |
| evidence class | Source, candidate artifact, published artifact, or native observation |
| status | `proven`, `failed`, or `unavailable` |
| observations | Commands, table rows, or visual surfaces inspected |
| blocker | Exact missing prerequisite when unavailable |

### State transitions

```text
unavailable -> proven  (prerequisite supplied and scenario passes)
unavailable -> failed  (prerequisite supplied and scenario fails)
failed      -> proven  (fix applied and same evidence class rerun)
```

`unavailable` never transitions implicitly and never counts as `proven`.
