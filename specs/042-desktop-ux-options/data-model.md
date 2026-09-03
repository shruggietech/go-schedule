# Data Model: Desktop UX and Options

S042 adds no daemon database entity or migration. It introduces bounded
per-user appearance preferences and in-memory desktop interaction state.

## Appearance preferences

| Field | Allowed values | Default | Persistence |
| --- | --- | --- | --- |
| Mode | `dark`, `light`, `system` | `dark` | Fyne per-user preferences |
| Font | `brand`, `system`, `monospace` | `brand` | Fyne per-user preferences |

Unknown, empty, or malformed identifiers normalize independently to their
defaults. Applying a preference creates a new immutable theme value and replaces
the application setting. Reset persists and applies both defaults through one
initiating UI action; persistence does not claim a multi-key transaction.

## Storage location

| Field | Meaning |
| --- | --- |
| Category | Stable user-facing purpose such as Machine data or Desktop preferences |
| Path | Clean absolute resolved value, or empty when unavailable |
| Availability | Available or unavailable with honest explanatory text |
| Scope | Machine, Current user, or Running application when installer ownership cannot be inferred safely |
| Exists | Present, absent, or unknown when inspection is not possible |
| Software-only removal | Plain-language retention behavior |
| Explicit data wipe | Plain-language removal or retention behavior |

Resolution is a pure ordered transformation over daemon-reported runtime paths,
the platform-owned root, executable, application-storage, and stat inputs. It
never creates, deletes, opens, or recursively searches a path.

## Navigation destination

| Field | Meaning |
| --- | --- |
| ID | Stable internal identifier |
| Label | Current visible text, including Activity badge when present |
| Content | Already-built destination canvas object |
| Selected | Exactly one ordinary destination is selected |
| Button | Focusable activation control |

Required destination order is Tasks, Groups, Chains, Schedule, Activity,
Options, Info. Exit is a navigation command, never a destination, and therefore
has no selected state or content.

## Task-row identity

| Field | Meaning |
| --- | --- |
| Task ID | Stable daemon identity bound during row update |
| Display text | Current rendered task summary |
| Activate callback | Requests edit by stable identity |

State transition:

```text
Unbound row
  -> list update with current task
Bound task ID
  -> single click
Selected only
  -> double activation and current ID exists
Editor requested for current task
  -> double activation and current ID is stale
Ignored
```

## Editor ownership

```text
Idle
  -> New, toolbar Edit, or valid row double activation
Open
  -> any additional edit activation
Open (request ignored)
  -> Save, Cancel, or dismissal
Idle
```

The state is guarded so concurrent event callbacks cannot own two task dialogs.

## Shutdown state

```text
Running
  -> title-bar close, connection Exit, or navigation Exit
Closing (cancel once, clear intercept once, close once)
  -> any later close request
Closing (no-op)
```

No transition returns to Running within the same `App` instance.
