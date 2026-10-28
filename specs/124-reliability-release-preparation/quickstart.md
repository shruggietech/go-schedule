# Validation Guide: S124

Run all non-Git tools using the verified hidden redirected launcher. Use existing tools and caches; do not provision a distribution or alter installed services.

## Source guards

```text
go run ./scripts/release-source-check -root . -repository shruggietech/go-schedule -tag v1.7.0
go run ./scripts/comparison-check
go run ./scripts/comparison-check -tag v1.7.0 -commit a2e4653021b6007135550ed8f84f1421f9072db0
```

The selected-tag invocation above uses accepted base only as a syntactically valid source-only exercise. It is not a resolved v1.7.0 tag identity. It validates source structure/digests, not remote tag existence. Preserve `site/data/published.json` at v1.6.0.

## New source QA

Freeze the assembled source and historical input hashes, then use new ignored child directories.

```powershell
./scripts/cli-qa.ps1 -Action build -BinaryDirectory A:/Code/go-schedule/.test-temp/s124/qa-inputs-01
./scripts/cli-qa.ps1 -Action capture -BinaryDirectory A:/Code/go-schedule/.test-temp/s124/qa-inputs-01 -OutputDirectory A:/Code/go-schedule/.test-temp/s124/qa-run-01 -WSLDistribution Ubuntu-24.04 -WSLUser h8rt3rmin8r -Incidents
./scripts/cli-qa.ps1 -Action finalize -BinaryDirectory A:/Code/go-schedule/.test-temp/s124/qa-inputs-01 -OutputDirectory A:/Code/go-schedule/.test-temp/s124/qa-run-01 -Reviews A:/Code/go-schedule/.test-temp/s124/reviews-01.json
./scripts/cli-qa.ps1 -Action validate -BinaryDirectory A:/Code/go-schedule/.test-temp/s124/qa-inputs-01 -OutputDirectory A:/Code/go-schedule/.test-temp/s124/qa-run-01
./scripts/cli-qa.ps1 -Action render -BinaryDirectory A:/Code/go-schedule/.test-temp/s124/qa-inputs-01 -OutputDirectory A:/Code/go-schedule/.test-temp/s124/qa-run-01
```

Capture success means processing completed and awaits message review. Read all actual corpus messages and contexts before creating schema-valid reviews with this run/corpus binding. Finalize/validate/render distinguish pass (0), valid failed (2), valid blocked (3) and processing/integrity error (1). Preserve originals and actual counts.

## Canonical aggregate

```sh
sh scripts/verify.sh all
```

Watch the foreground command to a terminal receipt; it may take around 18 minutes. Require all eight existing gates, preserve first failures, full timing and source/evidence hashes. Final factual records and formatting receive separate final audits.

## Source publication

Push the verified S124 branch and open one official PR referencing #394. Audit exact-head CI, all comments/threads and no more than two Codex rounds. Hand off unmerged with future tag, staging and promotion identities pending.
