# Implementation Plan: Dedicated IPC Administrative Group

**Branch**: `codex/030-ipc-admin-group` | **Date**: 2026-08-30 | **Spec**: [spec.md](spec.md)

## Summary

Turn the existing inert `admin_group` setting into the authorization boundary for the daemon's local endpoint. Resolve and apply a fail-closed group policy before serving, retain the former broad policy only when the operator explicitly configures an empty value, provision the default Windows group and installer user, and document equivalent Linux and macOS setup. Add platform-specific contract tests without requiring a manual VM session.

## Technical Context

**Language/Version**: Go 1.25, WiX Toolset 6.0.2 XML, POSIX shell, PowerShell, Markdown **Primary Dependencies**: Go standard library, `github.com/Microsoft/go-winio`, `golang.org/x/sys/windows`, `golang.org/x/sys/unix`, WiX Util extension **Storage**: Existing configuration file and data directory; no database schema change **Testing**: Platform-specific Go unit tests, daemon structured-log tests, Windows installer XML contracts, release workflow/build validation, and eight canonical verification gates **Target Platform**: Linux, macOS, and Windows service and foreground modes **Performance Goals**: Group lookup and permission setup occur once at daemon startup; no per-request authorization overhead or dispatch-latency impact **Constraints**: Fail closed for every non-empty unresolved group, preserve endpoint names and client protocol, keep daemon/CLI cgo-free, UTF-8 without BOM, no manual Windows VM requirement **Scale/Scope**: One daemon endpoint, one configured group, one Windows installer membership, three operator guides

## Constitution Check

| Principle | Gate | Design evidence |
| --- | --- | --- |
| I. Code Quality | PASS | One small access-policy value crosses the platform boundary; errors retain setting, group, endpoint, and failed operation context. |
| II. Testing Standards | PASS | Policy construction, failure cleanup, installer contracts, and logging receive regression tests before implementation; race and coverage remain mandatory. |
| III. UX Consistency | PASS | `admin_group` becomes a single documented schema field with a secure default, explicit compatibility value, structured state logging, and actionable startup errors. |
| IV. Performance | PASS | Authorization is delegated to operating-system endpoint permissions and adds startup-only work. |
| V. Autonomous Execution | PASS | S030 follows the full Spec-Kit sequence, blocking analysis, local CI parity, commit, and single pre-publication halt. |

### Post-design re-check

The design preserves the API protocol and scheduling hot path. The security boundary is enforced before readiness, every platform implementation has an injected test seam, and the required WiX pin change is recorded as a dated decision. No constitution violation remains.

## Architecture and Decision Log

### Fail closed when a configured group is unavailable

An empty `admin_group` is the only compatibility opt-in. Any non-empty lookup or permission failure aborts startup. This intentionally rejects issue #13's suggested automatic permissive fallback because a typo or missing provisioning step must not silently grant arbitrary-command execution to a broader audience.

### Return access-policy evidence from the listener

Change the internal listener contract to accept the configured group and return an `AccessInfo` value with the selected mode and group. Platform code owns resolution and permission application; daemon code owns structured operator logging. This keeps logging policy testable without coupling the transport package to a global logger.

### Secure the Unix endpoint without taking over unrelated directories

Restricted Unix mode resolves the numeric group, creates or tightens the managed default parent to group-owned `0770`, creates the socket, then sets and verifies group ownership and `0660`. A missing custom parent is safe to create with the same restricted policy. An existing custom parent is verified but never chowned or chmodded; unsafe ownership or mode fails startup with an actionable error. Compatibility mode explicitly manages the default parent as `0755` and socket as `0666` while leaving an existing custom parent unchanged. Any post-listen failure closes the listener and removes the socket.

This enforces traversal at the parent boundary without allowing a custom endpoint such as `/tmp/goschedd.sock` to retake ownership or permissions of a shared system directory.

### Build Windows descriptors only from resolved SIDs

Restricted Windows mode resolves the configured account through the operating system, accepts only group-like SID types, and inserts the canonical SID string into a protected descriptor granting full control to SYSTEM and built-in Administrators plus read/write to the configured group. Compatibility mode retains the Authenticated Users descriptor. Using a resolved SID avoids descriptor injection through configuration text.

### Provision the Windows group with WiX 6.0.2

WiX 5 cannot create local groups. Upgrade the pinned tool and both extensions to WiX 6.0.2, the first stable line supporting nested `util:Group` creation. The installer creates or reuses local `goschedadmin`, preserves it on uninstall, and uses `util:User` with `[LogonUser]` plus `GroupRef` to enroll the interactive installing user. WiX's group-before-membership sequencing handles initial install and repair idempotently.

### Keep archive platforms operator-managed

Linux and macOS have no package installer in this repository. Their guides will create `goschedadmin`, add the intended account, explain login-session refresh, and only then install/start the service. Creating new distribution packages would expand the slice without improving the runtime boundary.

## Project Structure

```text
internal/config/config.go
internal/config/config_test.go
internal/ipc/ipc.go
internal/ipc/ipc_unix.go
internal/ipc/ipc_unix_test.go
internal/ipc/ipc_windows.go
internal/ipc/ipc_windows_test.go
cmd/goschedd/main.go
cmd/goschedd/main_test.go
build/windows/goschedule.wxs
build/windows/verify_wxs.ps1
.github/workflows/release.yml
test/integration/windows_installer_contract_test.go
docs/INSTALL-windows.md
docs/INSTALL-linux.md
docs/INSTALL-macos.md
README.md
SECURITY.md
CHANGELOG.md
specs/030-ipc-admin-group/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── verification.md
├── contracts/ipc-access.md
├── checklists/
└── tasks.md
```

**Structure Decision**: Extend the existing config, IPC, daemon, installer, and documentation surfaces. No new package, background service, database entity, or network protocol is introduced.

## Complexity Tracking

No constitution violations require justification. The WiX major-version update is necessary because group creation does not exist in the pinned v5 extension; using a custom action would add more security-sensitive code and weaker rollback semantics.
