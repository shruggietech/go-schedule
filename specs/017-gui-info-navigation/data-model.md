# Data Model: GUI Navigation and Information

No persistent or daemon-owned data changes. This feature models immutable GUI
presentation values only.

## Navigation Item

| Field | Meaning | Constraint |
|---|---|---|
| label | Visible sidebar name | Exactly one of Tasks, Groups, Schedule, Activity, Info |
| position | Zero-based location | Matches the required five-item sequence |
| content | Existing or new local view | Non-nil and independently selectable |
| dynamic label | Optional bounded Activity count | May change text, never collection order |

## Application Identity

| Field | Meaning | Source |
|---|---|---|
| mark | Full official product artwork | Existing embedded application icon resource |
| version | Exact running build identifier | Existing compile-time build-version value |
| maintainer | Builder and maintainer attribution | ShruggieTech |

## Official Resource Link

| Label | Canonical destination |
|---|---|
| ShruggieTech | `https://shruggie.tech` |
| Source repository | `https://github.com/shruggietech/go-schedule` |
| Documentation | `https://shruggietech.github.io/go-schedule/` |

## State and Relationships

- Info has no loading, persistence, refresh, or error state because all content
  is compiled into the process.
- Activity owns the only dynamic navigation label. Updating its badge changes
  the existing item text but not its identity, position, or neighboring Info item.
- External destination availability is outside application state; the standard
  desktop hyperlink control delegates opening to the operating system.
