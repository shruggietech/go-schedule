# Data Model: Structured Desktop Data Tables

S045 introduces presentation models and exposes the existing stored run ID as
an optional Calendar occurrence response field. No database, event, request, or
persisted preference schema changes are required.

## Structured Column

Defines one aligned field across a header and all rows in a view.

| Field | Meaning | Rules |
| --- | --- | --- |
| `header` | Visible column label | Non-empty and unique within a view |
| `minimumShare` | Protected fraction of narrow available width | Greater than zero; all shares normalized by the allocator |
| `weight` | Share of width remaining after protected allocation | Non-negative |
| `alignment` | Leading, center, or trailing text alignment | Same for header and cells |
| `primary` | Whether the column carries record identity or primary meaning | At least one per view |

## Structured Cell

Represents one display field.

| Field | Meaning | Rules |
| --- | --- | --- |
| `text` | Normalized single-line display value | Never empty after fallback normalization |
| `fullText` | Unabridged value used in disclosure | Defaults to `text`; preserves Unicode exactly |
| `importance` | Theme semantic role | Medium, low/disabled, high/info, warning, danger/error, or success |
| `style` | Optional emphasis such as monospace or bold | Does not replace semantic meaning |

## Structured Row Model

| Field | Meaning | Rules |
| --- | --- | --- |
| `identity` | Stable source-record identity | Non-empty; never derived solely from current visual index |
| `cells` | Ordered values matching the view's columns | Count exactly equals column count |
| `summary` | Selectable full-value disclosure | Contains every cell's complete value and label |
| `alternate` | Quiet alternating surface choice | Derived from current row index only; never semantic |

## Task Row

Source: existing `domain.Task` plus the existing group snapshot.

| Column | Display normalization | Priority |
| --- | --- | --- |
| Task | Exact task name; `Unnamed task` fallback | Primary |
| Enabled | `Enabled` or `Disabled` from the independent boolean flag | High |
| Lifecycle | `Active`, `Completed`, `Disabled`, or title-cased neutral unknown | High |
| Time zone | Exact IANA time zone; `Unknown` fallback | Secondary |
| Group | Resolved group label; `Not assigned` fallback | Secondary |

Identity is the existing task ID. Selection, toolbar actions, and double activation resolve this identity against the current snapshot before acting.

## Schedule Row

Source: `server.Occurrence`, extended with optional `run_id` for past records.

| Column | Display normalization | Semantic role |
| --- | --- | --- |
| When | Existing localized formatted time | Neutral |
| Task | Exact task name; `Unnamed task` fallback | Neutral |
| Event | `SCHEDULED` for future or `COMPLETED` for past | Informational or neutral |
| Outcome | Glyph plus normalized outcome or `— Not available` | Success, error, warning, informational, disabled, or neutral |

Known outcomes: `SUCCESS`, `FAILURE`, `SKIPPED`, `CAUGHT UP`, `QUEUED`. Unknown values are uppercased with underscores changed to spaces and retain the neutral glyph `?`. Past occurrence identity uses the stored run ID so equal-time queued and completed records remain distinct across query-order changes. Computed future occurrences fall back to task identity plus timestamp because no run exists yet; a duplicate ordinal only disambiguates repeated representations of the same source identity.

## Activity Row

Source: existing merged `logEntry`, itself derived from `domain.LogRecord` or `domain.Alert`.

| Column | Display normalization | Semantic role |
| --- | --- | --- |
| When | Existing localized formatted time | Neutral |
| Severity | Glyph plus `INFO`, `WARNING`, `ERROR`, or neutral normalized unknown | Informational, warning, error, or neutral |
| Source | Existing source fallback (`daemon`) | Neutral |
| Summary | Exact message; `No message` fallback | Neutral |

Identity uses the existing log or alert ID when available, otherwise a deterministic composite of timestamp, source, message, and duplicate ordinal. Activation always reads the current visible row model before opening detail.

## State Transitions

```text
source snapshot
    ↓ normalize
ordered row models
    ↓ reconcile stable identity
retained selection OR cleared selection
    ↓ bind virtualized rows
visible structured table
```

- Refresh with retained identity: update/reorder models, find identity, select new index, refresh disclosure.
- Refresh without retained identity: clear visual selection and disclosure.
- Appearance/font change: theme refreshes row surfaces and labels; shared width allocation runs on subsequent layout.
- Empty source: zero rows, fixed header remains, disclosure shows the view's empty guidance.
