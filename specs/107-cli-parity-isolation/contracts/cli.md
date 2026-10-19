# S107 CLI contract

All new commands use scoped context and I/O. Results go to stdout, diagnostics to stderr; invalid input fails before mutation. Root --json remains supported. Existing history/log defaults remain unchanged unless an explicit new mode is requested.

| Surface | Contract |
| --- | --- |
| Local target | --ipc-path selects absolute Unix socket or Windows named pipe; incompatible remote flags fail, no fallback. |
| task add/edit | stdin/string/file (- for stdin), run identity, explicit edit clears for stdin/run-as/args/cwd/env, --disabled on create, --syntax auto/human/cron; existing repeatable flags remain. |
| group add/edit | Disabled create, rename/reparent/root using existing hierarchy rules. |
| schedule preview | Expression plus syntax/timezone/missing-date/time-basis/DST policies; predicted text or existing structured response, no mutation. |
| notification task/group set | --file PATH or - reads strict assignments request; conflicts with changed shorthand flags. |
| notification task/group show | Opt-in canonical input export supports lossless read/modify/reapply. |
| calendar | --from/--to RFC3339, default seven days, maximum 30; predicted and recorded source identifiers. |
| runs | Compatible history, active/get and optional readable detail with timestamps, output/truncation/provenance. |
| runtime | Local RuntimeInfo only, actionable unsupported-remote guidance. |
| logs | Compatible array/table; --metadata includes log path. |
| systems summary/search | Selected target by default; --all queries registrations with bounded work and independent identified outcomes. |
| preferences | show/set/reset every durable appearance/popup choice, omit retains state, explicit reset scope; identify effective user/path. |
| profile default | show/set PROFILE/reset saved desktop startup choice, no daemon/secret lookup. |

Structured input is UTF-8, at most 1 MiB, rejects unknown fields/trailing values and invalid collection shapes. Existing authorization and local/remote data availability remain unchanged.
