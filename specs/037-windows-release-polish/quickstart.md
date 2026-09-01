# Quickstart: Verify Windows Release Polish

## Startup sizing

Run the focused GUI tests:

```sh
go test ./gui/...
```

Expected: the required sizing matrix, 800x600 reachability contract, and existing GUI behavior pass.

On Windows, launch from a constrained secondary monitor and confirm the restored frame, title bar, resize borders, and taskbar remain reachable. Explicit maximize and restore must still work.

## MSI source contract

```powershell
pwsh build/windows/verify_wxs.ps1
```

Expected: the exact approved Summary Subject passes alongside all existing installer checks.

## Compiled MSI contract

Build the candidate with the pinned WiX process, then run:

```powershell
pwsh test/windows/inspect-installer.ps1 `
  -MsiPath C:\verify\candidate.msi `
  -EvidencePath C:\verify\artifact.md `
  -ArtifactClass candidate `
  -ArtifactOrigin 'local build from commit <full-commit-id>'
```

Expected: evidence records the candidate SHA-256 and Subject `go-schedule: cross-platform task scheduler`.

## Canonical verification

```sh
sh scripts/verify.sh all
```

All eight gates must pass in the foreground: format, vet, lint, race, gui, coverage, docs, and automation.
