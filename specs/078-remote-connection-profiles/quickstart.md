# Quickstart: Validate Remote Connection Profiles and Target-Safe Clients

## Prerequisites

- A clean checkout with supported Go, Node, POSIX shell, and Wails prerequisites.
- Two disposable remote daemon instances with trusted test certificates and distinct installation IDs.
- A disposable native credential-store namespace and user configuration directory.

## Scenarios

1. Run `sh scripts/verify.sh all` and require all eight gates.
2. Launch the desktop with no profile file and confirm This computer uses local IPC with unchanged workflows.
3. Pair two same-named test daemons, confirm both profiles persist without bearer values, restart, and select each exact identity.
4. Confirm the global target bar, Connections view, mutation confirmation, feature availability, and request destination all match the selected profile.
5. Rename one profile, repair it with a new phrase for the same daemon, verify the old native credential is removed, then reject a wrong-daemon repair without changes.
6. Remove inactive and active profiles, verify native deletion, and confirm active removal returns to This computer before deletion.
7. Run existing CLI commands without target flags and compare stdout plus transport with the pre-slice local behavior.
8. Pair a CLI profile through standard input, run supported reads and mutations with `--profile`, then exercise complete explicit endpoint selection.
9. Exercise ambiguous labels, incomplete flags, missing keyring values, certificate failure, daemon mismatch, revocation, unsupported remote routes, and concurrent profile writes.
10. Follow documented JSON examples through private HTTPS and an SSH tunnel, then prove a lost mutation response is never automatically replayed.
11. Scan profile files, Wails payloads, CLI stdout and stderr, logs, errors, and test artifacts for phrase and bearer canaries.

## Expected Outcome

Desktop and CLI users can deliberately select an exact remote daemon through a durable secret-free profile, every request is pinned to its trusted identity and native credential, This computer remains the default, and unsupported or ambiguous operations fail visibly without local fallback.
