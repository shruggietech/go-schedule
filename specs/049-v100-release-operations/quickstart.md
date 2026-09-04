# Quickstart: v1.0.0 Release Operations

## 1. Verify the implementation context

```powershell
./.specify/scripts/powershell/check-prerequisites.ps1 -Json -RequireTasks -IncludeTasks
git status --short --branch
git tag --list v1.0.0
gh release view v1.0.0
```

During the S049 pull request, the branch is
`codex/049-v100-release-operations`, the tree becomes clean after committed
work, and both v1.0.0 lookups must remain absent.

## 2. Run focused packet tests

```powershell
go test ./internal/releasegate ./scripts/windows-release-gate -count=1
go test -race ./internal/releasegate ./scripts/windows-release-gate -count=1
go test ./test/integration -run 'TestWindowsReleaseGate' -count=1
```

The tests cover exact mappings, Markdown escaping, deterministic files, atomic
commit, destination conflicts, local-demo rejection, candidate mismatch, and
the command's stream/exit contract.

## 3. Run canonical CI parity

```powershell
& 'C:\Program Files\Git\bin\bash.exe' scripts/verify.sh all
```

All eight ordered gates must pass: format, vet, lint, race, GUI, coverage, docs,
and automation.

## 4. Pull-request boundary

Push the review branch, open one PR targeting `main` with `Refs #122`, resolve
hosted CI and every Codex comment, and request no more than one explicit second
Codex round. The PR must leave v1.0.0 absent and every readiness issue open.

## 5. Post-merge release execution (separate authorization required)

After S049 merges, follow the amended
[`publication.md`](../048-v100-release-cut/contracts/publication.md). In order:

1. synchronize clean `main` and capture the reviewed S049 merge commit;
2. verify local/remote tag and GitHub release absence;
3. obtain explicit tag authorization, create one annotated `v1.0.0` tag at the
   S049 merge commit, and push only that tag;
4. wait for the tag-triggered Release workflow and require a draft containing
   seven packages plus `windows-candidate-manifest.json`;
5. download and verify the exact MSI and candidate manifest;
6. collect, finalize, and independently verify all 47 attended observations;
7. upload the immutable evidence ZIP, yielding exactly nine pre-checksum
   payloads;
8. generate the disposition packet with the command below;
9. review and apply each child record individually, update #96 only from actual
   child states, and keep #122 open;
10. obtain publication authorization and dispatch `Promote Release` only after
    every issue and candidate-bound gate passes; and
11. audit ten public payloads, checksums, latest identity, issue states, #122,
    and the milestone.

## 6. Generate the formal disposition packet

```powershell
go run ./scripts/windows-release-gate render-dispositions `
  --bundle '<downloads>\go-schedule_v1.0.0_windows-attended-evidence.zip' `
  --candidate-manifest '<downloads>\windows-candidate-manifest.json' `
  --artifact '<downloads>\go-schedule_v1.0.0_windows_amd64.msi' `
  --repository 'shruggietech/go-schedule' `
  --tag 'v1.0.0' `
  --commit '<reviewed-s049-merge-commit>' `
  --output-dir '<absent-workspace>\issue-dispositions'
```

Inventory `packet.json` and all ten Markdown files. Re-run `verify-bundle`
independently before any issue comment. Apply no record whose issue acceptance
criteria have not also been reviewed.

## 7. Stop conditions

Stop without publication if any of these occurs:

- tag or release state exists unexpectedly;
- `main`, the remote default branch, reviewed S049 merge commit, and intended
  tag target differ;
- staging is not one successful tag-push run or the release is already public;
- the eight staged or nine pre-checksum payloads differ from the contract;
- the MSI, manifest, evidence, attachments, archive contents, or packet fail;
- any readiness issue remains unsatisfied or its generated record is incomplete;
- the tag moves, candidate bytes change, checksums fail, or promotion rebuilds;
  or
- the final public release does not contain exactly nine payloads plus
  `SHA256SUMS.txt` and identify v1.0.0 consistently.
