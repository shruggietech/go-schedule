# Quickstart: v1.0.0 Release Execution and Audit

## 1. Confirm the immutable tag boundary

Require a clean synchronized repository, one annotated `v1.0.0` tag, and the same peeled commit locally and remotely:

```powershell
git status --short --branch
git rev-parse origin/main
git rev-parse 'v1.0.0^{}'
git ls-remote --tags origin refs/tags/v1.0.0 refs/tags/v1.0.0^{}
```

Expected peeled commit: `ff47b4410d1aecbfadb8165d1ebf025ca1dadde7`.

## 2. Accept the tag-triggered draft

Accept only Release run `33838072246` if it completes successfully with event `push`, head branch `v1.0.0`, head SHA equal to the peeled tag commit, and one successful Windows staging job. The release must remain draft and contain exactly seven packages plus `windows-candidate-manifest.json`.

## 3. Download and verify the candidate

Use an absent fixed-volume workspace outside the repository. Download all draft assets, then run:

```powershell
go run ./scripts/windows-release-gate verify-candidate `
  --candidate-manifest '<workspace>\windows-candidate-manifest.json' `
  --artifact '<workspace>\go-schedule_v1.0.0_windows_amd64.msi' `
  --repository 'shruggietech/go-schedule' `
  --tag 'v1.0.0' `
  --commit 'ff47b4410d1aecbfadb8165d1ebf025ca1dadde7'
```

## 4. Collect formal Windows evidence

Initialize the attended workspace from the exact downloaded MSI and accepted run identity. Complete every generated fragment from genuine native observation, including required screenshots and task-run evidence. Do not copy S043 or S047 results.

Finalize once with `Invoke-ReleaseCandidateAttended.ps1 -Action Finalize`, then independently run `windows-release-gate verify-bundle` against the separately downloaded MSI and manifest. The only passing result is 47 of 47 observations.

## 5. Upload evidence and reconcile issues

Upload exactly `go-schedule_v1.0.0_windows-attended-evidence.zip` to the draft, producing nine pre-checksum assets. Render the S049 disposition packet in an absent output directory, inspect `packet.json`, and review each Markdown record against the live issue acceptance criteria.

Apply and close eligible leaf issues first. Reconcile #96 from actual child and prerequisite states. Keep #122 open.

## 6. Promote through the reviewed workflow

After every readiness issue is satisfied and the tag is unchanged:

```powershell
gh workflow run "Promote Release" --ref main -f tag=v1.0.0
```

Accept only a successful workflow that validates the nine staged assets, candidate, evidence, tag, and workflow provenance; creates and verifies one checksum per payload; re-downloads the final set; and publishes the same draft.

## 7. Run the final audit

Download all public assets into an absent directory and require:

- exactly nine payloads plus `SHA256SUMS.txt`;
- every checksum passes and every file is non-empty;
- public/latest release, annotated tag, manifest, packages, notes, README, and changelog identify v1.0.0 consistently;
- all ten readiness issues are accurately closed;
- #122 contains the complete audit and closes only now;
- the v1.0.0 milestone has no open issues before closure; and
- S050 changes no candidate/runtime files.

## Stop conditions

Stop before the next mutation for any failed, partial, unavailable, ambiguous, missing, extra, mismatched, moved, replaced, or publicly pre-existing state. Record the failure without moving the tag or substituting assets.
