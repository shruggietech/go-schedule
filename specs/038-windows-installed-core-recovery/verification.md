# Verification: Windows Installed Core Recovery

**Date**: 2026-09-01

**Branch**: `codex/038-windows-installed-core-recovery`

**Status**: local implementation verified; hosted LocalSystem CI evidence pending

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

The hosted CI job has not run before pull-request publication and remains a
required delivery check. The candidate-MSI installed-core probe is unavailable
on this workstation because `goschedd` and `C:\\Program Files\\go-schedule`
are absent and the current session is not elevated. No fake runner, foreground
daemon, or injected backend is substituted for that missing installed proof.
The broader candidate release gate remains tracked by #94.

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

## Remaining delivery gate

S038 remains `In Progress` until the pull request's `Windows LocalSystem
execution` job passes and its artifact is available. After that evidence and
the required external reviews are satisfied, the specification can transition
to `Implemented` before the final maintainer merge ritual.
