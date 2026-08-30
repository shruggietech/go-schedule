# Quickstart: Verify S030

## Repository verification

1. Run the focused configuration, IPC, daemon, and installer contract tests.
2. Cross-compile daemon and CLI for Linux, macOS, and Windows with cgo disabled.
3. Build the WiX source with the pinned 6.0.2 tool and extensions.
4. Run `sh scripts/verify.sh all` in the foreground through all eight gates.

## Restricted-mode scenarios

1. Supply an existing group in `admin_group` and create a listener in the platform-specific test harness.
2. Assert the returned access evidence is `restricted` and names the group.
3. Assert Windows SDDL or Unix owner/mode matches [the contract](contracts/ipc-access.md).
4. Assert structured logs record the restricted mode before readiness.

## Fail-closed scenarios

1. Supply an unknown non-empty group.
2. Assert listener creation fails and the error names `admin_group` and the supplied value.
3. Inject a post-listen permission failure on Unix.
4. Assert the listener closes and the socket is removed.

## Compatibility scenario

1. Supply an explicitly empty `admin_group`.
2. Assert the historical broad policy from [the contract](contracts/ipc-access.md) is selected.
3. Assert one warning records `access_mode=compatibility`.

## Installer and documentation scenarios

1. Run the Windows installer XML contract against valid and deliberately broken group/membership relationships.
2. Run the pinned WiX sanity script and build validation.
3. Inspect each platform installation guide and confirm group provisioning precedes service startup, membership/session refresh is explained, and compatibility mode is documented as weaker.

## Expected result

Every supported platform has a testable restricted policy, compatibility requires an explicit setting, invalid security configuration prevents readiness, installer provisioning is reviewable offline, and canonical verification remains green.
