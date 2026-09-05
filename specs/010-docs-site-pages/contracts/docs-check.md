# Contract: `scripts/docs-check.sh`

The documentation-integrity gate. Run locally before pushing and as the `docs` CI job. Pure POSIX `sh` + coreutils; no network, no Ruby, no build.

## Invocation

```sh
sh scripts/docs-check.sh
```

No arguments. Run from the repository root (like `scripts/coverage-gate.sh`). It discovers its own inputs.

## Inputs

- Every `*.md` file directly under `docs/`.
- `docs/_sass/custom/custom.scss` for the pinned dark-theme contract.
- The pointer README(s): `test/scripts/README.md`.
- The explicit current-product inventory in `scripts/docs-policy-check.sh` and its generated good/stale/missing fixtures.

## Rules (all must hold)

1. **Front matter present.** Each `docs/*.md` MUST begin with a YAML front-matter block (`---` … `---`) whose keys include `title:` and `nav_order:`.
2. **On-disk link integrity.** For each Markdown link `[text](target)` in a `docs/*.md`:
   - Skip if `target` starts with `http://` or `https://`.
   - Skip if `target` is a pure `#fragment`.
   - Otherwise strip a trailing `#fragment`, resolve `target` relative to the containing file's directory, and require the resolved path to exist.
3. **No escape from `docs/`.** A resolved link target from rule 2 MUST NOT lie outside the `docs/` directory (no surviving `../` that climbs above `docs/`). Content outside `docs/` must be referenced by absolute `https://github.com/…` URL, which rule 2 skips.
4. **Pointer validity.** Each pointer README's `docs/*.md` link target MUST exist.
5. **Fence vocabulary.** Every published fenced block MUST use `sh`, `bash`, `powershell`, or `text`, close with plain triple backticks, and be balanced.
6. **Dark-theme contract.** The custom stylesheet MUST retain safe token fallbacks, named-role colors, highlighted/selected code treatment, and responsive endorsement spacing.
7. **Current product policy.** Current README, guide, help, and authoritative S001 surfaces MUST retain human-first dual-syntax positioning, equivalent examples, the exact five-field/fidelity contract, and API source identity. Targeted categorical human-only claims MUST fail. Historical specs and the changelog are not scanned globally.
8. **Policy fixtures.** The aligned fixture MUST pass. Fixtures containing an obsolete categorical claim or omitting required posture MUST fail.

Anchors are intentionally not validated (fragments are stripped, not resolved).

## Output & exit codes

| Exit | Meaning | stdout/stderr |
|------|---------|----------------|
| `0`  | All rules pass. | A one-line success summary (page count checked). |
| `1`  | One or more rules failed. | Each failure on its own line as `docs/<file>: <reason>: <link-or-field>`, then a total. |

Errors go to stderr; the summary to stdout, consistent with principle III's stream conventions. The script uses `set -eu` and quotes all expansions.

## CI binding

`.github/workflows/ci.yml` gains a `docs` job:

```yaml
docs:
  name: Docs link check
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - name: docs-check
      run: sh scripts/docs-check.sh
```

Triggers: the workflow's existing `push`/`pull_request` to `main`. No `needs` (runs in parallel). `permissions: contents: read` is inherited from the workflow default. This is a change to a pinned artifact and is recorded as a dated decision in `CHANGELOG.md`.

## Test scenarios (drive implementation)

| # | Setup | Expect |
|---|-------|--------|
| T1 | Repo as shipped (all front matter present, links valid) | exit 0 |
| T2 | A `docs/*.md` with no front-matter block | exit 1, names the file: missing front matter |
| T3 | A `docs/*.md` front matter missing `nav_order` | exit 1, names the file + missing key |
| T4 | A relative link to a non-existent `docs/` file | exit 1, names file + broken link |
| T5 | A surviving `../outside.md` escape link | exit 1, names file + escape link |
| T6 | A `#fragment`-only or `https://` link | ignored (no failure) |
| T7 | `test/scripts/README.md` pointing at a missing `docs/` file | exit 1, names the pointer + target |
| T8 | Current dual-syntax product copy and contracts | exit 0 |
| T9 | Fixture adds a targeted human-only cron prohibition | exit 1, names obsolete posture |
| T10 | Fixture removes the required README posture | exit 1, names missing policy |

T1-T7 are exercised by the integrity gate; T8-T10 run through `test/scripts/docs-policy-check_test.sh` with temporary fixture directories. The shipped tree must satisfy every scenario.
