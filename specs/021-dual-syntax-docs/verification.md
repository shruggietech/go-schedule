# Verification Evidence: S021

## Baseline

- `sh scripts/docs-check.sh`: PASS (11 pages, links, front matter, fences, theme)
- `go test ./internal/cli ./gui -count=1`: PASS

## Reproduced correctness defect

Before implementation, `0 9 */2 * *`, `0 9 * */2 *`, and
`0 9 * * */2` were each accepted as `every day at 09:00`, losing the step.

## Final evidence

- Parser/shared-boundary/API/CLI focused suites: PASS.
- Documentation policy fixtures: PASS (aligned accepted, stale and missing rejected).
- Documentation integrity: PASS (11 pages plus theme and product policy).
- Canonical `scripts/verify.sh all`: PASS after adding the already installed
  Scoop GCC directory to this process's `PATH`.
- Eight gates: format, vet, lint, race, GUI, coverage, docs, automation.
- Core coverage: engine 88.1%, schedule 91.5%, timezone 88.9%, store 87.0%,
  catchup 87.5%, logbus 91.1%.
