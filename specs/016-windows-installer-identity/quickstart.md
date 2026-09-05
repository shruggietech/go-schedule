# Quickstart: Validate Windows Installer Identity

## Portable source contract

From the repository root:

```powershell
go test ./test/integration -run WindowsInstaller
pwsh build/windows/verify_wxs.ps1
```

Expected: both commands pass. The Go test also proves independent failures for missing or misdirected icon relationships.

## Candidate MSI table contract

Build the MSI with WiX 5.0.2 using the same inputs as the release workflow, then run:

```powershell
pwsh test/windows/inspect-installer.ps1 `
  -MsiPath <candidate.msi> `
  -ArtifactClass candidate `
  -ArtifactOrigin 'local build from commit <full-commit-id>'
```

Expected: the Icon, Property, Shortcut, and Environment relationships are reported as `proven`. This does not install the package.

## Clean Windows lifecycle

Copy the MSI and repository `test/windows/` directory to a disposable Windows machine, then follow [test/windows/README.md](../../test/windows/README.md).

The lifecycle script refuses a host where go-schedule is already installed or on PATH. Pass the candidate or published evidence class and its build reference or release-asset URL. During a paused run, answer the native observation prompts; the script writes those responses into the generated evidence record.

## Repository verification

```sh
sh scripts/verify.sh all
```

All eight gates must pass before the local feature commit. Unavailable native or published evidence remains visible in `verification.md`; it is not converted to a pass.
