# Quickstart: v1.1.1 Release Recovery

## 1. Verify recovery preparation

```powershell
go run ./scripts/github-format
git diff --check
```

Confirm the v1.1.1 changelog boundary, README version, four release highlights, tagged changelog footer, and unchanged v1.1.0 tag.

## 2. Run repository gates

```powershell
bash -lc 'cd /mnt/a/Code/go-schedule && GO="/mnt/c/Program Files/Go/bin/go.exe" GOFMT="/mnt/c/Program Files/Go/bin/gofmt.exe" sh scripts/verify.sh all'
```

## 3. Review and merge

Push the preparation branch, open one pull request referencing #140, satisfy hosted CI and every review finding, and merge only the reviewed green commit.

## 4. Retire the old draft and stage v1.1.1

Fetch and prune origin, fast-forward main, verify that v1.1.0 remains unpublished and immutable, delete only its GitHub draft, create the annotated v1.1.1 tag at the reviewed S058 merge commit, and require successful main CI for that exact commit before staging artifacts.

## 5. Qualify, promote, and audit

Download all draft assets, verify the candidate manifest and MSI, complete the smallest valid attended Windows qualification against that exact MSI, promote without rebuilding, freshly download every public asset, verify all checksums and version surfaces, then close #140 and milestone v1.1.1.
