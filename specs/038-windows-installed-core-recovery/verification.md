# Verification: Windows Installed Core Recovery

**Date**: 2026-09-01

**Branch**: `codex/038-windows-installed-core-recovery`

**Status**: implemented; local and hosted LocalSystem verification passed

## Pre-change baseline

- `go test ./internal/ipc ./internal/executor`: PASS before production changes.
- Existing descriptor granted SYSTEM, Built-in Administrators, and only the
  configured group SID. The direct enrolled user SID was absent.
- Existing executor process-start failure exposed only the raw `os/exec` error
  and did not name the failed boundary.

## Test-first evidence

- New IPC tests initially failed to compile because `listGroupMembers` and the
  member record did not exist.
- The new process-start regression failed against the raw error
  `exec: "definitely-not-a-real-go-schedule-executable-xyz": executable file not found in %PATH%`.
- After implementation, `go test -race ./internal/ipc ./internal/executor`:
  PASS.
- `go test ./test/integration -run WindowsInstaller -count=1`: PASS after the
  lifecycle and CI probe contracts were updated.
- Both PowerShell files parse without errors through
  `System.Management.Automation.Language.Parser`.

## Native ordinary-token named-pipe evidence

The current Windows workstation provided the exact authorization state that
previously failed:

- current identity: `NOVX999\\h8rt3rmin8r`;
- user SID: `S-1-5-21-106076546-252188030-3528287134-1000`;
- direct `goschedadmin` membership: present;
- current medium-integrity token: `goschedadmin` SID absent;
- Built-in Administrators SID: deny-only in the current token.

`TestListenWindowsCurrentDirectMemberWithoutAliasToken` used production
Netapi32 enumeration and a production `go-winio` listener. It confirmed the
current user was a direct member, confirmed the alias SID was absent from the
token, opened a unique restricted named pipe, and connected through the
standard token. Result: PASS.

Unit descriptor tests additionally prove:

- SYSTEM and Built-in Administrators retain full access;
- the configured group retains its read/write ACE;
- direct user SIDs are validated, sorted, and de-duplicated;
- group-valued direct members are not recursively expanded;
- Authenticated Users and Everyone are absent from restricted mode;
- lookup, enumeration, and invalid-SID failures occur before listener creation;
- a non-alias configured group preserves its prior group-ACE behavior.

An unrelated account was not created on this workstation. Its exclusion is
verified by exact descriptor assertions rather than claimed as a second-user
native login.

## Service-hosted execution evidence

The repository now contains two non-fake proof paths:

1. `test/windows/Invoke-ServiceCoreCI.ps1` installs staged binaries as the real
   `goschedd` LocalSystem service on a disposable Windows CI runner. It proves
   manual and scheduled exit-code-0 runs, output, marker effects, a controlled
   exit-code-7 child, and a missing-executable process-start failure. CI uploads
   the resulting JSON evidence artifact.
2. `Invoke-InstallerLifecycle.ps1 -Scenario installed-core-probe` applies the
   same assertions to a release-equivalent installed candidate from a normal,
   non-elevated session and records MSI, authorization, service, run, marker,
   and daemon-log evidence.

The hosted `Windows LocalSystem execution` job passed in workflow run
[`33508591264`](https://github.com/shruggietech/go-schedule/actions/runs/33508591264),
job [`99858512438`](https://github.com/shruggietech/go-schedule/actions/runs/33508591264/job/99858512438),
at commit `5c644bf`. Its `windows-service-core-evidence` artifact
(`9800693172`, SHA-256 digest
`d90ed84fde5aaf6b25bf189160c16d61b2c8472e13d225bea4be523ebb47bb77`)
recorded:

- `LocalSystem` service account and `Running` service state;
- direct local-group membership with the group SID absent from the ordinary
  user's token;
- manual and scheduled success with exit code 0 and `S038-output`;
- a controlled child-process exit code of 7;
- a missing-executable process-start failure without argument, standard-input,
  or environment disclosure; and
- three service-authored marker lines with artifact SHA-256
  `2f5573d4b6bd4fc6d5c835b0a8b233c82d409ea5a0dd204b59e99d3ee0bdbc39`.

The first two hosted attempts correctly failed rather than producing false
proof. Run `33507430229` exposed that the service could not write into the
runner account's temporary directory. Run `33508009559` exposed a quoting flaw
in the probe's `cmd.exe` command that allowed a final `echo` to hide the failed
marker write. The probe now uses the absolute inbox Windows PowerShell path,
non-interactive flags, and direct .NET file output to a service-writable
ProgramData location. Run `33508591264` then proved the intended path.

The candidate-MSI installed-core probe remains unavailable on this workstation
because `goschedd` and `C:\\Program Files\\go-schedule` are absent and the
current session is not elevated. No fake runner, foreground daemon, or injected
backend substitutes for that unavailable proof. The broader candidate release
gate remains tracked by #94, as required by FR-017.

## Canonical verification

Command (foreground and watched):

```text
C:\\Program Files\\Git\\bin\\bash.exe scripts/verify.sh all
```

First run stopped at vet with `possible misuse of unsafe.Pointer`. The Netapi32
buffer was changed from a `uintptr` round-trip to a typed out-pointer and the
complete suite was restarted from gate one.

Final result: PASS.

| Gate | Result | Evidence |
| --- | --- | --- |
| format | PASS | no unformatted Go files |
| vet | PASS | all packages |
| lint | PASS | `0 issues.` |
| race | PASS | all selected packages; integration completed in 55.023s |
| gui | PASS | `gui` and `gui/viewmodel` |
| coverage | PASS | engine 86.4%, schedule 89.2%, timezone 91.3%, store 84.1%, catchup 88.9%, logbus 91.1% |
| docs | PASS | 14 pages plus policy fixtures |
| automation | PASS | actions, CodeQL, Dependabot, release notes, brand, lifecycle, and eight-gate contract |

## Delivery disposition

S038 implementation and required verification are complete. Pull-request CI,
external review, maintainer approval, and merge remain publication workflow
evidence rather than open feature tasks under the repository lifecycle policy.
