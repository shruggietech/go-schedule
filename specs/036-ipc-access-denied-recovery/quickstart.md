# Quickstart: Validate IPC Access-Denied Recovery

## Automated regression

```sh
go test ./internal/api/client ./gui/...
sh scripts/verify.sh all
```

Expected: access denied, unavailable, timeout, other transport, and API status failures classify separately; three simultaneous startup failures yield one incident; all eight gates pass.

## Native Windows walkthrough

1. Install the candidate MSI on a supported disposable Windows 11 host.
2. Before refreshing the login session, record `Get-Service goschedd`, `Get-LocalGroup goschedadmin`, `Get-LocalGroupMember goschedadmin`, and `whoami /groups` output.
3. Launch the GUI and record the one in-frame state, copy, Retry, Exit, and lack of recurring dialogs.
4. Sign out and back in, confirm `whoami /groups` includes the group SID, and record that a new GUI process reaches usable tabs without reinstalling.
5. Store redacted observations in `verification.md`; unavailable prerequisites remain explicitly unresolved.

The opt-in native probe records the installed stale-token diagnosis without entering ordinary CI:

```powershell
$env:GOSCHED_NATIVE_EXPECT = 'stale'
go test ./gui -run '^TestNativeWindowsConnectionRecovery$' -count=1 -v
```

The deterministic `TestSuccessfulRetryClearsIncidentWithoutReinstall` test supplies the authorized transition: it activates the same classified incident, changes the backend to success, invokes Retry, and requires the incident to clear while normal GUI tabs remain available. `GOSCHED_NATIVE_EXPECT=connected` remains an optional field diagnostic after a natural future login refresh, not a slice completion prerequisite.
