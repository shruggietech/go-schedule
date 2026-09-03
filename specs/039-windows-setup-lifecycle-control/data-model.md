# Data Model: Windows Setup Lifecycle Control

S039 changes installer and cleanup state only. It does not change the scheduler database or application configuration schema.

## Shortcut Feature

| Field | Values | Rules |
| --- | --- | --- |
| Identity | `StartMenuShortcut`, `DesktopShortcut` | Stable across S039-and-later upgrades |
| Desired state | local, absent | Independently selectable |
| Fresh default | Start Menu local; desktop absent | Determined by feature install level |
| Installed component | one shortcut component per feature | Each owns one shortcut and one machine registry key path |

**Transitions**: fresh default or explicit selection -> installed state -> maintenance add/remove -> upgrade migration -> removal.

## Completion Selection

| Field | Values | Rules |
| --- | --- | --- |
| Launch application | selected, unselected | Fresh attended default selected |
| Open documentation | selected, unselected | Fresh attended default unselected |
| Eligibility | fresh full-UI success only | Never persisted and never evaluated in execute sequence |

The two fields are independent, producing four valid combinations.

## Removal Mode

| Field | Values | Rules |
| --- | --- | --- |
| `GOSCHEDULE_REMOVE_DATA` | `0`, `1` | Secure public property; any other explicit value is invalid |
| Default | `0` | Preserve in attended and unattended removal |
| Attended transition to wipe | `0` -> confirmation -> `1` | Back/cancel returns to `0` |
| Unattended transition to wipe | exact command property `1` | Never persisted |

## Owned Root Candidate

| Field | Description | Validation |
| --- | --- | --- |
| Kind | machine data or user preference | Closed internal enumeration |
| Source | trusted Windows known location or registered local profile | No caller-supplied paths |
| Product leaf | `goschedule` or `AppData\Roaming\fyne\tech.shruggie.goschedule` | Exact case-insensitive component match |
| Path | canonical absolute owned root | Fixed local volume; within expected parent; no root/ancestor reparse |
| State | absent, validated, deleting, deleted, refused, residual | Monotonic within one commit invocation |

## Cleanup Result

| State | Meaning | MSI behavior |
| --- | --- | --- |
| not requested | preserve mode | No data discovery or mutation |
| running | explicit commit wipe began | Atomic ledger identifies interrupted work |
| complete | every declared existing root was removed | Remove stale result evidence |
| refused | preflight found unsafe state | Delete nothing; retain protected report |
| partial | deletion began but one or more owned roots remain | Complete software removal; retain protected report and summary |
| internal-error | helper could not maintain trustworthy evidence | Complete software removal; retain safest available diagnostic state |

## Cleanup Ledger

| Field | Purpose |
| --- | --- |
| Schema version | Enables compatible result parsing |
| Transaction identifier | Distinguishes cleanup attempts without carrying a path input |
| Started/completed UTC | Orders and diagnoses interrupted work |
| Overall state | `running`, `complete`, `refused`, `partial`, or `internal-error` |
| Candidate kind and profile SID | Identifies declared scope without publishing account names |
| Canonical owned path and outcome | Supports actionable residual cleanup |
| Operating-system error | Explains an owned-path failure |

## Invariants

- A wipe action has no path-valued public input.
- No deletion begins until every existing candidate passes preflight.
- Recursive deletion never traverses a reparse point.
- A commit result cannot cause Windows Installer to pretend deleted data was rolled back.
- Security state is outside the owned-root model.
