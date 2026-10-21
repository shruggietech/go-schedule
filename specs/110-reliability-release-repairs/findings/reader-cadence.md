# S110 CI correction: deterministic reader cadence

Published issue: [#340](https://github.com/shruggietech/go-schedule/issues/340), parent #277, milestone 15.

## Outcome

Make canned heartbeat reader-query coverage independent of incidental host-clock timestamps, while retaining automatic interval inference, insufficient-data refusal and explicit interval override behavior.

## Evidence and authority

Parent: #277. Discovered during owner-authorized S110 CI monitoring in PR #338. Pending source 5be817f033b9125e70d6db6842c9631c49514887 failed [the native Linux race job](https://github.com/shruggietech/go-schedule/actions/runs/39671902546/job/119471665372). TestScriptsReaderQueries/posix reports query jitter exiting 2 with interval not supplied and could not be inferred (too few records), followed by jitter needs an interval. Its recorder seeded two beats using a real one-second loop, but the fixture did not retain their actual timestamps in that failure.

Twenty unchanged native Linux repetitions passed with row diagnostics (33.792s). They do not diagnose the hosted failure. A controlled valid two-row seed with equal started_ms values reproduces the same inference refusal and old fixture failure (1.681s). This demonstrates the fixture's unstated positive observed-cadence assumption; it does not establish that a particular host clock adjustment caused the historical run.

## Bounded scope

Retain the real recorder seed and its two-record requirement, then give the reader fixture a fixed, consistent one-second observation timeline. Keep positive reader calls on automatic inference rather than passing an interval to hide that path. Assert that drift, jitter and gaps disclose their inferred one-second cadence. Include independent equal-timestamp and single-beat scenarios that retain jitter's usage failure and actionable guidance, plus successful explicit interval override. Failed query diagnostics must retain the seed rows. No reader implementation or production timestamp behavior changes belong in this correction.

## Functional acceptance criteria

- The canned-reader fixture owns deterministic, internally consistent timestamps instead of deriving its valid analysis input from an unasserted wall-clock gap.
- Automatic interval inference and its disclosure remain exercised for drift, jitter and gaps, and existing canned-query success behavior remains represented.
- Equal-time and single-beat inputs retain the insufficient-cadence usage error and interval guidance; an explicitly supplied interval remains usable.
- A failed query includes the actual fixture row identities, timestamps and declared intervals needed to distinguish a future data or reader defect.

## Dependencies

Parent #277. Uses the existing script twins and disposable SQLite fixtures; follows #335's literal binding correction but is an independent reader-fixture outcome. The earlier #339 credential fixture correction is unrelated and remains included in PR #338. Both Codex rounds are complete, so no third is requested. Keep failed hosted results and their source identities distinct from later passing checks.
