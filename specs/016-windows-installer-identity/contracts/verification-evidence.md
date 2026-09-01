# Contract: Windows Installer Verification Evidence

## Evidence classes

1. **Source definition**: portable structural tests over tracked WXS and release
   workflow sources.
2. **Candidate artifact**: MSI table inspection of a locally built package.
3. **Published artifact**: the same inspection and lifecycle steps against a
   downloaded release asset.
4. **Native observation**: Start Menu, installed-apps, window title area, and
   taskbar observed in a clean Windows desktop session.

## Status vocabulary

- `proven`: the named evidence class ran and met its expected result.
- `failed`: the named evidence class ran and did not meet its expected result.
- `unavailable`: the evidence class did not run because a named prerequisite
  was absent.

No other word counts as a completion status. An `unavailable` status cannot be
used to close an issue whose acceptance criteria require that evidence class.

## Clean lifecycle sequence

This identity-era sequence is historical. S035 supersedes its installer
security and lifecycle boundary with compiled `Wix4Group` inspection, separate
fresh and v0.9.0-upgrade hosts, explicit repair/reinstall, group SID/member
invariance, service-order diagnostics, and a non-elevated post-refresh access
probe. See `specs/035-windows-msi-local-group/contracts/lifecycle-evidence.md`.

1. Prove no CLI resolution, install directory, or machine PATH entry.
2. Record Windows build and artifact version/origin/hash.
3. Install per-machine.
4. Open a new shell and run version plus task-list commands.
5. Inspect installed-apps and Start Menu identity.
6. Launch and inspect native running-GUI identity.
7. Reinstall or upgrade and prove one PATH entry remains.
8. Uninstall and prove the entry is removed without unrelated PATH drift.
