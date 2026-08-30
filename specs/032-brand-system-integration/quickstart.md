# Quickstart: Validate the Repository Brand System

## Prerequisites

- Go toolchain selected by `go.mod`
- POSIX shell used by the repository verification driver
- No graphics, font, browser, Python, or PDF package is required for routine checks

## 1. Inspect the canonical entry point

Open `brand/README.md` for identity rules and `brand/REPOSITORY.md` for repository integration. Confirm `brand/VERIFY.md` reports PASS and `brand/brand-guide.pdf` opens as an eleven-page tagged guide.

## 2. Run focused integrity checks

```text
go test ./scripts/brand-check
go run ./scripts/brand-check
```

Expected result: all tests pass, followed by one `brand-check: OK` summary.

## 3. Validate documentation

```text
sh scripts/verify.sh docs
```

Expected result: every page has valid front matter, all local brand downloads resolve inside `docs/`, and existing documentation policy checks pass.

## 4. Validate packaging contracts

```text
pwsh -File build/windows/verify_wxs.ps1
go test ./gui/...
```

Expected result: the installer and embedded desktop identities reference approved assets and GUI tests remain green.

## 5. Exercise negative drift detection

Copy one mapped asset to a temporary location, alter its repository consumer, run `go run ./scripts/brand-check`, observe the exact source/target failure, then restore the consumer from the temporary copy. Do not commit this deliberate mutation.

## 6. Run complete CI parity

```text
sh scripts/verify.sh all
```

Expected result: format, vet, lint, race, GUI, coverage, docs, and automation all pass in the foreground.

## 7. Review issue disposition

The eventual pull request description uses `Closes #10` and `Closes #34`. It must not close #33.
