# Quickstart: Validate S038

## Repository verification

From the repository root:

```powershell
go test ./internal/ipc ./internal/executor
go test -race ./internal/ipc ./internal/executor
& 'C:\Program Files\Git\bin\bash.exe' scripts/verify.sh all
```

## Native installed verification

Use a disposable Windows 11 client that satisfies the lifecycle script's clean-host requirements. Build or obtain the exact candidate MSI, then follow `test/windows/README.md` for installation and the S038 access/execution probe.

The proof must retain:

- candidate SHA-256 and product version;
- service state/account and group/token SID evidence;
- generated pipe descriptor and normal-client health result;
- created task JSON with environment values redacted;
- manual, scheduled, nonzero-exit, and process-start-failure run JSON;
- marker files and correlated daemon logs.

If a disposable host, elevated setup, or release-equivalent MSI is unavailable, record the native result as unavailable. Do not substitute an in-process runner or fake backend.
