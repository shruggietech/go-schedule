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

Give the maintainer the exact linked MSI/hash and ask for an ordinary usability
walkthrough, not a second formal release ceremony. Collect display, DPI,
artifact, and other machine-readable facts automatically. Record what the
maintainer actually observes in `checklists/attended-demo.md`, including any
visible failure, unavailable hardware, or accepted follow-up.

The maintainer decides whether the local demo is a valid release target. A
reproducible in-scope blocker is recorded before code changes, gets a failing
test, and requires a new MSI/hash plus a repeat of affected checks. Accepted
post-v1 polish is filed as a GitHub issue and does not expand S047. The full
43-scenario matrix is executed later against the formal post-merge candidate.

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
