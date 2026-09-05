# Quickstart: Verify GUI Dual-Syntax Scheduling

## Focused automated check

```powershell
go test ./gui ./internal/scheduleinput ./internal/api/server -count=1
```

## Manual review path

1. Open a new recurring task and enter `0 9 * * 1-5`.
2. Confirm Save becomes available and preview describes weekdays at 09:00.
3. Save, reopen, and confirm the field still contains `0 9 * * 1-5`.
4. Replace it with `weekdays at 09:00` and confirm preview/save still work.
5. Enter an invalid or refused cron expression and confirm Save disables with a specific reason.
6. Confirm a one-off task and a human interval with Start at behave unchanged.

## Canonical repository check

```powershell
$env:Path = "$env:USERPROFILE\scoop\apps\gcc\current\bin;$env:Path"
& 'C:\Program Files\Git\bin\bash.exe' scripts/verify.sh all
```

The canonical command must report all eight gates green before the local commit.
