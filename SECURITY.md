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

Do not open a public issue for a vulnerability. This is a small project and a public report is a disclosure with no lead time.

If GitHub's private reporting form is unavailable, email **info@shruggie.tech** instead. Use the same private-reporting precautions and do not include vulnerability details in a public issue.

Useful things to include, in rough order of how much they help: what an attacker gains, the exact steps to reproduce, the version and install method, and whether the attacker needs local access or an account on the machine. A proof of concept is welcome but not required, a clear description of the mechanism is often enough.

Repository administrators own initial triage. You can expect an acknowledgement within a week. This is a one-to-two developer project, so please read that as a realistic commitment rather than a service level.

## Supported versions

| Version | Supported |
| --- | --- |
| 0.6.x | Yes |
| 0.5.x and earlier | No N/A please upgrade |

Fixes land on `main` and ship in the next release. There is no backport branch.

## What go-schedule is, security-wise

Being explicit about this matters more than a policy statement, because the answer is not obvious from the outside: **go-schedule runs arbitrary commands, and when installed as a system service it runs them with system privileges.** That is the entire purpose of the software. It follows that anyone who can create or edit a task can execute code at that privilege level.

So the security boundary that matters is **who can reach the daemon**. By default, that boundary is the dedicated operating-system group `goschedadmin`:

- **Windows.** The IPC named pipe carries an explicit ACL granting SYSTEM and built-in Administrators full control, and the resolved `goschedadmin` SID read and write. The MSI creates or reuses that group and adds the interactive installing account.
- **Linux and macOS.** The daemon listens on a Unix domain socket in its data directory. In restricted mode the directory is owned by `goschedadmin` with mode `0770`, and the socket is group-owned with mode `0660`. A custom existing socket directory must already match that group and mode; the daemon verifies it rather than mutating unrelated operator-owned paths.

Any non-empty group lookup, SID type, ownership, permission application, or readback failure stops startup before the API is served. An explicitly empty `admin_group` selects compatibility mode instead: Authenticated Users regain named-pipe access on Windows, or the Unix socket uses mode `0666`, and the daemon emits a warning. Treat membership in the configured group, or selecting compatibility mode, as granting the daemon's execution privilege.

## Known limitations

These are current and acknowledged, not undiscovered:

- **Release artifacts are unsigned.** The Windows `.msi` is not Authenticode-signed and the macOS builds are neither signed nor notarized. Verify downloads against `SHA256SUMS.txt`; that checksum file is the integrity guarantee on offer.
- **IPC access control is group-based**, as described above; it does not provide per-command roles within the API.
- **Task credentials.** A task's environment is stored in the local database in plain text. Do not put secrets in `--env`; point the task at a secret store your platform already provides.
- **No authentication on the local API.** Reaching the socket or pipe *is* the authorization.

Reports that a documented limitation exists are welcome as issues rather than advisories. Reports that one can be exploited beyond its stated scope, privilege escalation across a boundary this document claims holds, are exactly what the private route is for.
