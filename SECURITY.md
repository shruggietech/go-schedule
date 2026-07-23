# Security policy

**Audience:** anyone reporting a vulnerability, or deciding how to deploy
go-schedule\
**Applies to:** go-schedule 0.6.0 and later

## Contents

- [Reporting a vulnerability](#reporting-a-vulnerability)
- [Supported versions](#supported-versions)
- [What go-schedule is, security-wise](#what-go-schedule-is-security-wise)
- [Known limitations](#known-limitations)

## Reporting a vulnerability

Please report privately, using GitHub's private vulnerability reporting:

**<https://github.com/shruggietech/go-schedule/security/advisories/new>**

Do not open a public issue for a vulnerability. This is a small project and a
public report is a disclosure with no lead time.

Useful things to include, in rough order of how much they help: what an attacker
gains, the exact steps to reproduce, the version and install method, and whether
the attacker needs local access or an account on the machine. A proof of concept
is welcome but not required — a clear description of the mechanism is often
enough.

You can expect an acknowledgement within a week. This is a one-to-two developer
project, so please read that as a realistic commitment rather than a service
level.

## Supported versions

| Version | Supported |
| --- | --- |
| 0.6.x | Yes |
| 0.5.x and earlier | No — please upgrade |

Fixes land on `main` and ship in the next release. There is no backport branch.

## What go-schedule is, security-wise

Being explicit about this matters more than a policy statement, because the
answer is not obvious from the outside: **go-schedule runs arbitrary commands,
and when installed as a system service it runs them with system privileges.**
That is the entire purpose of the software. It follows that anyone who can
create or edit a task can execute code at that privilege level.

So the security boundary that matters is **who can reach the daemon**, and the
answer today is deliberately permissive, because the design target is a
single-user desktop or a machine whose local users are all trusted:

- **Windows.** The IPC named pipe carries an explicit ACL granting SYSTEM and
  Administrators full control, and Authenticated Users read and write. The last
  of those is what lets a non-elevated GUI reach a `LocalSystem` daemon at all —
  a non-elevated administrator token carries its Administrators SID as
  deny-only, so Authenticated Users is required for the ordinary case to work.
  The consequence is that any authenticated local user can manage the scheduler.
- **Linux and macOS.** The daemon listens on a Unix domain socket in its data
  directory, created with default permissions. Any local user who can traverse
  that directory can connect.

If you are deploying onto a multi-user or locked-down machine, treat access to
the daemon as equivalent to the privilege the daemon runs with, and restrict it
at the operating-system level. Narrowing the IPC ACL to a dedicated
administrative group is tracked as open work; it is not implemented today, and
this document says so rather than implying otherwise.

## Known limitations

These are current and acknowledged, not undiscovered:

- **Release artifacts are unsigned.** The Windows `.msi` is not
  Authenticode-signed and the macOS builds are neither signed nor notarized.
  Verify downloads against `SHA256SUMS.txt`; that checksum file is the integrity
  guarantee on offer.
- **IPC access control is coarse**, as described above.
- **Task credentials.** A task's environment is stored in the local database in
  plain text. Do not put secrets in `--env`; point the task at a secret store
  your platform already provides.
- **No authentication on the local API.** Reaching the socket or pipe *is* the
  authorization.

Reports that a documented limitation exists are welcome as issues rather than
advisories. Reports that one can be exploited beyond its stated scope — privilege
escalation across a boundary this document claims holds — are exactly what the
private route is for.
