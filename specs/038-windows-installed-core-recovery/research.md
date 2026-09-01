# Research: Windows Installed Core Recovery

## Verified local evidence

On 2026-09-01, the current workstation reported:

- current user SID `S-1-5-21-106076546-252188030-3528287134-1000`;
- that user as a direct local `goschedadmin` member;
- no `goschedadmin` SID in the medium-integrity process token;
- Built-in Administrators present only as `Group used for deny only`;
- no currently installed `goschedd` service or `C:\\Program Files\\go-schedule` tree, so a fresh real-service run could not be claimed from this checkout.

The first four facts reproduce the authorization mismatch recorded by S036. The missing installation means task-execution root cause cannot be inferred from this workstation's current state.

## Decision 1: Enumerate local-group members by SID

**Decision**: Call `NetLocalGroupGetMembers` at level 1 through `Netapi32.dll`, retain direct entries whose `SID_NAME_USE` is `SidTypeUser`, and add their validated SIDs to the restricted pipe descriptor.

**Rationale**: Microsoft documents that named-pipe access checks compare the client's token SIDs with the pipe DACL. A stable user SID remains in a standard token even when a newly enrolled alias SID is absent. The membership API returns both SID and account type, allowing the daemon to avoid recursively expanding group-valued members.

**Alternatives rejected**:

- Authenticated Users or Everyone ACE: violates least privilege and issue acceptance.
- Built-in Administrators as normal path: requires elevation and fails UAC-filtered standard tokens.
- Sign-out only: does not address filtered-token observations and allowed the installed first launch to remain broken.
- Installer-only per-user config: duplicates authorization state and becomes stale after membership changes.
- Per-connection impersonation plus membership query: increases complexity and performs authorization work on every connection.

## Decision 2: Preserve the configured group ACE

**Decision**: Keep the group SID ACE in addition to direct user member ACEs.

**Rationale**: Fresh tokens and supported nested/global group membership already work through Windows access checks. Removing the group ACE would regress those users and force recursive enumeration.

## Decision 3: Treat enumeration failures as startup failures

**Decision**: Fail daemon startup if the configured group cannot be enumerated. Skip only individual deleted or unknown entries that supply no valid SID.

**Rationale**: Falling back to a broader descriptor would violate the authorization policy; silently retaining only the group ACE would recreate the reported broken state for direct members.

## Decision 4: Extend native lifecycle evidence, not the fake scheduler runner

**Decision**: Add a bounded installed execution probe to the existing Windows lifecycle script.

**Rationale**: The executor already passes in-process Windows tests, while the reported failure occurs across the LocalSystem service boundary. The existing scheduling integration's recording runner intentionally cannot prove process creation, service PATH, working directory, output capture, or marker effects.

## Decision 5: Use an absolute inbox executable and marker side effects

**Decision**: Use absolute inbox Windows PowerShell with `-NoLogo -NoProfile -NonInteractive -Command`, direct .NET marker-file I/O, explicit marker paths, and no interactive input.

**Rationale**: An absolute inbox path removes service PATH ambiguity. Output plus marker files proves both record persistence and external process effects. Direct .NET file I/O avoids `cmd.exe`'s documented nonstandard argument parsing around redirection and quoting.

## Decision 6: Improve only verified diagnostic weakness

**Decision**: Prefix no-exit-code failures with a stable process-start diagnostic naming only the executable.

**Rationale**: Current history can contain a raw OS error but does not define the failed boundary. Including the executable is actionable; including arguments or environment values can expose secrets. No evidence supports changing the service identity, executor environment, or command parsing.

## Sources

- Microsoft Learn, `NetLocalGroupGetMembers`: member SID and account-type enumeration, level semantics, buffer ownership, and failure codes.
- Microsoft Learn, `LOCALGROUP_MEMBERS_INFO_1`: SID, `SID_NAME_USE`, and account name structure.
- Microsoft Learn, Named Pipe Security and Access Rights: client token-to-DACL access checks.
- Microsoft Learn, SID Attributes in an Access Token: enabled and deny-only SID behavior.
- GitHub issues #90 and #93 and S036 verification evidence.
