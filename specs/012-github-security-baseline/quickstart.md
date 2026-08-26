# Quickstart: GitHub Security Baseline

## Local validation

Run from the repository root before commit:

```sh
sh test/scripts/automation-check_test.sh automation
sh scripts/automation-check.sh .
shellcheck scripts/automation-check.sh test/scripts/automation-check_test.sh
sh scripts/verify.sh all
git status --short
```

Expected local evidence:

- all positive automation fixtures pass;
- obsolete CodeQL, unknown action, missing trigger, insufficient permission,
  and missing analysis fixtures fail for the expected reason;
- the real workflow satisfies the advanced CodeQL contract;
- all eight canonical verification gates pass in order;
- verification introduces no additional working-tree changes.

Test-first evidence: before the CodeQL actions and workflow contract were added
to `scripts/automation-check.sh`, the extended positive fixture failed by
naming both `github/codeql-action/init@v4` and
`github/codeql-action/analyze@v4` as unaudited references.

Post-implementation focused evidence:

- `sh test/scripts/automation-check_test.sh automation`: PASS.
- `sh scripts/automation-check.sh .`: PASS (`approved actions, CodeQL contract,
  and 8-gate manifest`).
- `shellcheck scripts/automation-check.sh
  test/scripts/automation-check_test.sh`: PASS on the Windows host. The WSL
  environment does not have ShellCheck installed, so the same files were
  checked with the available host binary rather than reported as skipped.

Full local evidence:

- `sh scripts/verify.sh all`: PASS on 2026-08-26 through `format`, `vet`,
  `lint`, `race`, `gui`, `coverage`, `docs`, and `automation` in order.
- Core coverage: engine 88.1%, schedule 91.5%, timezone 88.9%, store 87.0%,
  catchup 87.5%, and logbus 91.1%.
- `.github/workflows/ci.yml`, `.github/workflows/release.yml`, `Makefile`,
  product code, `go.mod`, and `go.sum`: no diff.
- Diff whitespace, UTF-8 BOM, and mojibake audits: PASS.
- Spec/plan/task analysis: 16/16 functional requirements mapped to executable
  tasks, no unresolved clarification or template markers, and both requirement
  checklists complete.
- Verification introduced no additional working-tree changes.

## Pre-activation baseline (read-only, observed 2026-08-26)

- Repository visibility: public.
- Private vulnerability reporting: disabled.
- Secret scanning: disabled.
- Secret scanning push protection: disabled.
- Non-provider patterns: disabled.
- Validity checks: disabled.
- Code scanning: no analysis found.
- Organization fallback publication: `info@shruggie.tech`.

## Hosted activation (after explicit authorization only)

1. Push the verified local `main` commit.
2. Enable private vulnerability reporting.
3. Request each secret protection individually.
4. Read every setting back and classify it using
   `contracts/security-settings.md`.
5. Inspect the private advisory route without submitting a report.
6. Observe the CodeQL run on `main`; verify the workflow name, Go analysis,
   successful conclusion, and security result publication.
7. Query aggregate alert counts when token scope permits. Never print contents.
8. Record enabled, unavailable, unverified, and failed results separately.

The pull-request trigger and weekly schedule are proven locally by contract;
their future hosted executions cannot both be forced without out-of-scope
remote artifacts or waiting for time to pass.
