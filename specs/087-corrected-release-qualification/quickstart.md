# S087 validation workflow

1. Verify the refreshed tag, completed Release run and eight draft assets identify the reviewed source. Download manifest and MSI and run existing `windows-release-gate verify-candidate` with their recorded identity.
2. Record candidate, independent baseline, official portable PowerShell ZIP and offline WebView2 installer in the manifest described by `test/windows/README.md`. Run `go run ./scripts/windows-qualification-session --manifest <absolute-manifest> --output <new-absolute-session-directory>`.
3. Run fresh.wsb and upgrade.wsb independently in disposable Windows 11. S086 handles setup and exports; complete actual native observations and suitable separate normal-user, multiple-profile, and high/mixed-DPI environments.
4. Finalize genuine complete evidence through the original host collector and shared gate. Record unavailable observations as unfinished work. Keep the release a draft until promotion approval.
