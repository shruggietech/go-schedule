# Quickstart: v1.0.0 Release Cut

## 1. Verify the preparation context

```powershell
./.specify/scripts/powershell/check-prerequisites.ps1 -Json -RequireTasks -IncludeTasks
git status --short --branch
git tag --list v1.0.0
gh release view v1.0.0
```

Before merge, the feature path resolves to `specs/048-v100-release-cut`, the
tree is clean after committed work, and both tag/release lookups confirm
v1.0.0 does not exist.

## 2. Run focused preparation checks

```powershell
go test ./internal/releasegate ./scripts/windows-release-gate -count=1
go test ./test/integration -run 'Test(ReleaseWorkflowStagesEveryUploadAsDraft|PromotionOrdersExactGateChecksumsAndPublication|AttendedCollectorUsesCanonicalScenariosAndHiddenChildren|MSIInspectorCanWriteExactCandidateManifest)$' -count=1
& 'C:\Program Files\Git\bin\bash.exe' test/scripts/automation-check_test.sh
& 'C:\Program Files\Git\bin\bash.exe' scripts/automation-check.sh
```

Confirm the release-note policy accepts v1.0.0, the workflows retain staged
candidate/promotion ordering, and the evidence validator requires 47 scenarios.

## 3. Audit release copy

Verify mechanically:

- one README badge derived from GitHub's latest published release;
- one `daemon ok (version 1.0.0)` example;
- one empty Unreleased section;
- one dated v1.0.0 section containing the 33 original top-level entries;
- correct Unreleased/v1.0.0 comparison links; and
- five release-note bullets plus one tagged changelog link.

## 4. Run canonical CI parity

```powershell
& 'C:\Program Files\Git\bin\bash.exe' scripts/verify.sh all
```

All eight ordered gates must pass: format, vet, lint, race, GUI, coverage, docs,
and automation.

## 5. Publish the preparation PR

Push `codex/048-v100-release-cut`, open a PR targeting `main` with `Refs #122`,
wait for hosted CI and Codex review, address every comment, and request at most
one explicit second Codex round. Stop for the maintainer's merge ritual when all
checks and threads are green.

## 6. Post-merge release ritual (separate authorization required)

Follow [contracts/publication.md](contracts/publication.md) exactly. In summary:

1. synchronize and verify the reviewed merge boundary;
2. obtain explicit authorization and push annotated tag `v1.0.0`;
3. wait for successful draft staging and download the exact MSI/manifest;
4. collect and finalize all 47 attended Windows observations;
5. upload and independently verify the evidence archive;
6. reconcile all ten v1 readiness issues independently;
7. dispatch Promote Release for `v1.0.0`; and
8. audit assets, checksums, public/latest identity, #122, and the milestone.

Any failed prerequisite leaves the release draft and affected issues open.
