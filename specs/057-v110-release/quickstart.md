# Quickstart: v1.1.0 Release

## 1. Verify preparation

```powershell
go run ./scripts/github-format
git diff --check
```

Confirm the changelog boundary, README version, four release highlights, and tagged changelog footer.

## 2. Run repository gates

```powershell
bash -lc 'cd /mnt/a/Code/go-schedule && GO="/mnt/c/Program Files/Go/bin/go.exe" GOFMT="/mnt/c/Program Files/Go/bin/gofmt.exe" sh scripts/verify.sh all'
```

## 3. Review and merge

Push the preparation branch, open one pull request referencing #140, satisfy hosted CI and review findings, and merge only the reviewed green commit.

## 4. Stage and qualify

Create the annotated v1.1.0 tag at synchronized main, wait for draft staging, download all assets, verify the candidate manifest and MSI, and complete the smallest valid attended Windows qualification against that exact MSI.

## 5. Promote and audit

Upload the evidence archive, run the existing promotion workflow, wait for success, freshly download the public assets, verify all checksums and version surfaces, then close #140 and the v1.1.0 milestone.
