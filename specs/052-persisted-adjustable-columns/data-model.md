# Data Model: Persisted Adjustable Columns

## Column layout profile

| Field | Meaning | Validation |
|-------|---------|------------|
| View identity | Schedule or Activity preference namespace | Non-empty and application-defined |
| Schema version | Serialization contract version | Must equal supported version |
| Column identities | Ordered stable names | Exact match to current view |
| Proportions | Relative width intent | Same count, finite, positive, normalized to 1 |
| Defaults | Built-in relative intent | Valid for current schema |

## Boundary adjustment

| Field | Meaning | Validation |
|-------|---------|------------|
| Boundary index | Left column index | Zero through column count minus two |
| Delta | Logical horizontal movement | Finite and clamped by adjacent minimums |
| Available width | Current table content width | Non-negative after gaps |

Transition: resolve widths, transfer the clamped delta between neighbors,
normalize, refresh consumers, and persist when the interaction completes.

## Stored preference

One atomic current-user JSON value per view. It is accepted only when version,
identities, length, and every numeric value validate. Rejection returns the
complete default rather than a partially recovered mixture.
