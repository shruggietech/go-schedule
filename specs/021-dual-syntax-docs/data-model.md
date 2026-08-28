# Data Model: Documentation Concepts

S021 adds no persisted entities. It standardizes these concepts:

- **Schedule input**: recurring text classified as `human` or `cron`.
- **Human phrase**: readable product syntax compiled by the human parser.
- **Cron expression**: five timing fields, or a supported macro, without a
  command or crontab-file context.
- **Supported cron subset**: expressions faithfully representable by
  go-schedule.
- **Fidelity refusal**: a named rejection when conversion would change timing.
- **Task timezone**: the timezone governing recurrence and DST resolution.
- **Source identity**: retained expression plus `human` or `cron` classification.
- **Current surface**: operator-facing documentation, help, or authoritative
  contract subject to the policy check.
- **Historical artifact**: chronological specification evidence retained with
  a supersession notice when its boundary is no longer current.
