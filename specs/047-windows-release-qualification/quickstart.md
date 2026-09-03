# Quickstart: Windows Release Qualification

## 1. Validate the specification context

```powershell
./.specify/scripts/powershell/check-prerequisites.ps1 -Json -RequireTasks -IncludeTasks
```

Confirm the feature resolves to `specs/047-windows-release-qualification` and
both requirement-quality checklists have no incomplete items.

## 2. Run focused automated checks

```powershell
go test ./internal/releasegate ./scripts/windows-release-gate -count=1
go test -race ./internal/releasegate ./scripts/windows-release-gate -count=1
go test ./test/integration -run 'Test(ReleaseWorkflowStagesEveryUploadAsDraft|PromotionOrdersExactGateChecksumsAndPublication|AttendedCollectorUsesCanonicalScenariosAndHiddenChildren|MSIInspectorCanWriteExactCandidateManifest)$' -count=1
pwsh -NoProfile -Command "[void][scriptblock]::Create((Get-Content -Raw test/windows/Invoke-ReleaseCandidateAttended.ps1))"
pwsh -NoProfile -File build/windows/verify_wxs.ps1
```

Expected: the passing fixture contains 43 scenarios; deleting or corrupting a
new desktop scenario fails with its identity/metric; production validation
rejects the fixture class; the collector parses; installer authoring passes.

## 3. Run canonical CI parity

```powershell
& 'C:\Program Files\Git\bin\bash.exe' scripts/verify.sh all
```

All eight gates must pass in order: format, vet, lint, race, GUI, coverage,
docs, automation. Do not substitute hosted CI for a missing local prerequisite.

## 4. Commit before creating the demo identity

After source, specs, tests, and verification are stable, create the local S047
implementation commit. Record its full commit and derive:

```text
Embedded version: 1.0.0-s047-demo.<short-commit>
MSI ProductVersion: 1.0.0
Artifact class: local-demo
```

No push, PR, tag, workflow dispatch, or release is permitted at this step.

## 5. Build and inspect the S047 local demo

Mirror the pinned Windows build section in `.github/workflows/release.yml`:

1. Build `gosched.exe`, `goschedd.exe`, `gosched-gui.exe`, and
   `gosched-cleanup.exe` for Windows amd64 with the embedded demo version.
2. Compile `build/windows/goschedule.wxs` with WiX 6.0.2 and numeric version
   `1.0.0`.
3. Name the output
   `dist/go-schedule_s047-demo_<short-commit>_windows_amd64.msi`.
4. Run `test/windows/inspect-installer.ps1` with `ArtifactClass local-demo`.
5. Record source commit, embedded version, ProductVersion, ProductCode, byte
   size, SHA-256, build time, and inspection path in `verification.md`.

The demo is ready only when compiled inspection and all prior checks pass.

## 6. Perform the pre-push attended walkthrough

Test only the exact linked MSI/hash and record the result in
`checklists/attended-demo.md`:

- seamless install, shortcuts, Finish actions, startup, and no error spam;
- Dark and Light at 100 percent and the available scaled QHD configuration;
- System default/reset plus System, Geist, Inter, Ubuntu, and Monospace;
- control states, rail boundary/Exit, compact storage/selector behavior;
- conventional wheel at 1x/2x/4x and touchpad if available;
- at least 100 Tasks, Schedule, and Activity rows at 1280x800 and 800x600;
- guided preserve, cancel, and wipe uninstall journeys where safely practical.

A failure is recorded before code changes. Any correction gets a failing test,
a new commit, a rebuilt MSI/hash, and affected checks repeat.

## 7. Stop before publication

At the S047 halt, report:

- completed Spec Kit and implementation artifacts;
- focused and canonical gate results;
- exact demo MSI identity and link;
- attended checklist status;
- formal candidate, issue disposition, and publication work still outstanding;
- the exact branch push and PR commands, but do not execute them.

## Later formal qualification (outside the pre-push run)

After review and merge, an independently authorized tag stages the formal MSI.
Initialize `Invoke-ReleaseCandidateAttended.ps1` from that exact draft asset,
complete all 43 scenarios, finalize the archive, reconcile each issue, and only
then authorize promotion. Local S047 demo results must be repeated rather than
copied into the formal bundle.
