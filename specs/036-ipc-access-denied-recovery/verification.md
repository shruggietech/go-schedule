# Verification Evidence: IPC Access-Denied Recovery

## Baseline

- Issue: #90
- Released tree: v0.9.1 (`52202601b73ba63654f339564909ca989846e4f3`)
- Observed behavior: model, calendar, and stream paths independently create
  recurring modal errors; shared client appends daemon-absence copy to denial.
- Native Windows host: installed v0.9.1 service and stale-token condition
  verified on `NOVX999`.

## Test-first record

- Red: `go test ./internal/api/client ./gui/...` failed on the missing typed
  connection categories and incident coordinator before production code existed.
- Green: the same focused command passes classification, simultaneous
  tasks/calendar/events deduplication, evidence-conditioned guidance, reachable
  Retry/Exit controls, API-status separation, bounded backoff, and recovery
  clearing without reinstall.
- Windows cross-build: `GOOS=windows CGO_ENABLED=0 go test
  ./internal/api/client` and `go test ./gui -run '^$'` pass.
- Live pipe: `go run ./cmd/gosched health` returns one accurately classified
  `daemon IPC access was denied for the current session` error and contains no
  daemon-absence suffix.

## Canonical gates

- 2026-08-31 first canonical run through Git for Windows `sh.exe`:
  format PASS, vet PASS, lint PASS (0 issues), race PASS, GUI PASS, coverage PASS
  (engine 86.4%, schedule 89.2%, timezone 91.3%, store 84.1%, catchup 88.9%,
  logbus 91.1%), docs PASS. Automation correctly rejected a Draft/In Progress
  inventory mismatch; state was corrected before rerun.
- 2026-09-01 final canonical run: `scripts/verify.sh all` PASS through
  `C:\Program Files\Git\bin\sh.exe`; format, vet, lint (0 issues), race, GUI,
  coverage, docs, and automation all passed in order. Coverage remained engine
  86.4%, schedule 89.2%, timezone 91.3%, store 84.1%, catchup 88.9%, and logbus
  91.1%.

## Native Windows evidence

- Host: native Windows desktop session on `NOVX999`, current v0.9.1 service.
- Service: `goschedd` Running, Automatic.
- Local group: `goschedadmin` exists with SID
  `S-1-5-21-106076546-252188030-3528287134-1010`.
- Installing account: `NOVX999\h8rt3rmin8r` is a direct member with SID
  `S-1-5-21-106076546-252188030-3528287134-1000`.
- Current GUI/command-runner token: `whoami /groups` does not contain the
  `goschedadmin` SID, verifying the stale-token diagnosis.
- Live endpoint: the current token receives Windows `Access is denied` opening
  `\\.\pipe\goschedd`.
- In-frame incident and recovery controls: headless native GUI tests PASS and
  assert one visible state with Retry, Exit, and normal tabs still present.
- Opt-in native stale probe: `GOSCHED_NATIVE_EXPECT=stale go test ./gui -run
  '^TestNativeWindowsConnectionRecovery$' -count=1 -v` PASS (0.07 seconds).
- Authorized recovery transition: `TestSuccessfulRetryClearsIncidentWithoutReinstall`
  changes the same GUI/backend from classified denial to success, invokes the
  actual Retry control, and requires the incident to clear with the normal tabs
  retained. PASS.

## Verification-boundary decision

- 2026-09-01: issue #90 was re-authored under direct operator authority to
  separate native diagnosis from deterministic recovery. Native Windows proves
  the installed service, group, account, token, restricted-pipe denial, and
  selected guidance. The automated authorization transition proves Retry and
  complete GUI recovery. A forced sign-out is no longer a completion gate
  because it terminates the active development session without adding coverage
  beyond those two independently verified contracts.

## Completion

- 2026-09-01 native installed-host stale diagnosis: PASS.
- 2026-09-01 deterministic authorized Retry recovery: PASS.
- 2026-09-01 post-decision canonical eight-gate suite: PASS.
- Issue #90 acceptance coverage: complete under the operator-approved
  non-disruptive verification boundary.
