# Quickstart: Validate the Wails Release Cutover

## Preconditions

- Work from the S066 review branch with Go 1.25, Node.js 24, npm, Wails platform prerequisites, and a C toolchain available.
- Use the exact branch revision for every local and hosted result.
- Do not create a tag or invoke the release workflow while validating this slice.

## 1. Validate specification and current-product residue

```text
pwsh -File .specify/scripts/powershell/check-prerequisites.ps1 -Json -RequireTasks -IncludeTasks
sh scripts/spec-lifecycle-check.sh
go run ./scripts/github-format
```

Inspect remaining Fyne references and require each to be historical, release-history, or preference-migration context.

## 2. Validate the production desktop

```text
cd desktop/frontend
npm ci
npm audit --audit-level=high
npm test
npm run build
npm run test:e2e

cd ..
go test -race ./...
go run github.com/wailsapp/wails/v2/cmd/wails@v2.14.0 build -clean
```

The native output must use the stable `gosched-gui` name.

## 3. Validate root and release contracts

```text
go test ./test/integration/... ./scripts/...
sh scripts/verify.sh all
```

The canonical verifier must complete in the foreground. No gate may be omitted or reported green when it did not run.

## 4. Validate hosted evidence

After the pull request is published, require green native Wails builds on Windows, macOS, and Linux, the browser accessibility contract, root race and coverage gates, Windows installer contracts, package-payload checks, and CodeQL for the same head revision.

## 5. Audit parity and publication boundary

Review [release-candidate.md](contracts/release-candidate.md), the final parity inventory, and intentional deviations. There must be no `blocked` parity entry and no release tag, draft release, public release, or milestone closure created by the slice.
