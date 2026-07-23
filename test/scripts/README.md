# Maintainer test scripts

Test payloads for verifying an installed `goschedd` actually fires on time,
survives restarts, catches up, and honors overlap policies.

**Full documentation: [`docs/test-scripts.md`](../../docs/test-scripts.md).**
Read that, not this file — this is only a signpost.

```bash
pwsh -File test/scripts/Test-Heartbeat.ps1 -Label smoke
```

```bash
pwsh -File test/scripts/Test-ReadTestDB.ps1 -Database heartbeat -Query summary
```

Each script is a matched pair: a PowerShell `.ps1` and a POSIX `.sh` with
identical behavior. A PowerShell `-FooBar` is `--foo-bar` in the shell twin.

Requires `sqlite3` 3.33.0+. If it is missing, the scripts exit 2 and name both
`-InstallSqlite` and your platform's package-manager command.
