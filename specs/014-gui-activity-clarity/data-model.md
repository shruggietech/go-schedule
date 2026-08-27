# Data Model: GUI Activity Clarity

This feature introduces no persisted entities or API schema changes.

## Display state

### Activity tab label

- **Input**: unacknowledged alert count (integer)
- **Output**:
  - count <= 0: `Activity`
  - count 1..99: `Activity (count)`
  - count >= 100: `Activity (99+)`

### Activity row

The existing internal `logEntry` remains the normalized display record for a
daemon log record or scheduler alert. No fields or semantics change.

### Clear-view cutoff

The existing in-memory timestamp remains local to the built Activity view.
Entries at or before the cutoff are hidden. Later entries are visible. Alert
IDs in the currently displayed rows are acknowledged through the existing
Backend method. No persisted record is deleted.
