# Contract: Canonical Verification Command

## Entry points

Canonical direct invocation:

```sh
sh scripts/verify.sh all
```

Make convenience wrapper:

```sh
make verify
```

The Make wrapper MUST delegate to the canonical direct invocation and MUST NOT
restate child commands.

## Modes

```text
sh scripts/verify.sh list
sh scripts/verify.sh all
sh scripts/verify.sh format
sh scripts/verify.sh vet
sh scripts/verify.sh lint
sh scripts/verify.sh race
sh scripts/verify.sh gui
sh scripts/verify.sh coverage
sh scripts/verify.sh docs
sh scripts/verify.sh automation
```

- `list` prints the required gate names, one per line, in aggregate order.
- `all` runs every listed gate once, sequentially, in the foreground.
- A named gate runs only that gate.
- Missing or unknown modes print usage to stderr and exit non-zero.

## Gate contract

| Gate | Required behavior |
| --- | --- |
| `format` | List unformatted Go files under `internal`, `cmd`, and `test`; fail without rewriting |
| `vet` | Vet the complete Go module |
| `lint` | Run the project-pinned golangci-lint version using the module Go toolchain |
| `race` | Run race tests with cgo, excluding only `cmd/gosched-gui` and `gui`, retaining `gui/viewmodel` |
| `gui` | Run the headless GUI test packages |
| `coverage` | Invoke `scripts/coverage-gate.sh` with its enforced thresholds |
| `docs` | Invoke `scripts/docs-check.sh` |
| `automation` | Invoke `scripts/automation-check.sh` against the repository root |

## Output and exit behavior

- Before a gate starts, stdout identifies its stable gate name.
- Child stdout/stderr remains visible in the foreground.
- A successful named mode exits zero.
- A failed or unavailable gate exits non-zero and the most recent heading names
  the failing gate.
- Aggregate mode stops at the first non-zero child result.
- No verification mode runs a formatter, modifies tracked source, commits,
  pushes, tags, or starts release automation.

## Environment

- The caller may set `GOTOOLCHAIN`; otherwise lint mode derives an appropriate
  toolchain selector from the `go` directive in `go.mod`.
- A POSIX shell, Go, Git, and a C toolchain for race mode are prerequisites.
- Windows uses Git Bash or an equivalent environment already supported by the
  coverage and documentation scripts.
