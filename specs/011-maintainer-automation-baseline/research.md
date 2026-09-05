# Research: Maintainer Automation Baseline

**Date**: 2026-08-26

## D1 - Upgrade all four Node 20 action families to verified Node 24 majors

**Decision**: Use `actions/checkout@v7`, `actions/setup-go@v7`, `actions/upload-artifact@v7`, and `softprops/action-gh-release@v3`.

**Rationale**: Inspection of each currently selected major's `action.yml` confirms that checkout v4, setup-go v5, upload-artifact v4, and action-gh-release v2 declare Node 20. Their current supported major lines declare Node 24. Official releases available on 2026-08-26 are checkout v7.0.1, setup-go v7.0.0, upload-artifact v7.0.1, and action-gh-release v3.0.2. The repository already follows floating-major action references, so retaining that convention minimizes unrelated policy change.

The repository uses only contracts retained by the new majors: checkout's default repository/ref behavior and explicit `ref: main`; setup-go's `go-version-file` and cache inputs; upload-artifact's archived named upload; and action-gh-release's name, files, draft, prerelease, and release-note inputs. Upload-artifact v7 adds optional direct upload but keeps archived upload as the default. Action-gh-release v3's stated major change is the Node 24 runtime.

**Alternatives considered**:

- Patch-level SHA pinning: stronger supply-chain immutability, but a material policy change outside #21 and inconsistent with current repository practice.
- Minimum intermediate Node 24 majors: no benefit over the current supported lines and a shorter maintenance horizon.
- Upgrade only the three actions named in the CI warning: rejected because the release workflow's `softprops/action-gh-release@v2` also declares Node 20.

**Primary sources**:

- <https://github.com/actions/checkout/releases/tag/v7.0.0>
- <https://github.com/actions/setup-go/releases/tag/v7.0.0>
- <https://github.com/actions/upload-artifact/releases/tag/v7.0.0>
- <https://github.com/softprops/action-gh-release/releases/tag/v3.0.0>

## D2 - Make a POSIX shell driver the command source of truth

**Decision**: Add `scripts/verify.sh` with named gate modes and aggregate `all` mode. CI invokes named modes; `Makefile` targets delegate to them; the documented canonical local command is `sh scripts/verify.sh all` (or the same script via Git Bash on Windows).

**Rationale**: The repository already relies on POSIX shell for coverage and documentation gates, including Git Bash on Windows. A shell driver can compose those scripts and Go commands without adding a dependency. Making the Makefile the source would exclude Windows environments without Make and would not be directly consumable by all hosted runner shells. Making workflow YAML the source would leave local execution duplicated.

**Alternatives considered**:

- `make verify` as the canonical source: useful as a wrapper, but Make is absent on some supported maintainer hosts, including the current Windows host.
- Duplicate the canonical list in Make and CI: this is the existing failure mode and does not satisfy #41.
- Add a Go verification binary: cross-platform, but disproportionate for command orchestration and still unable to remove the existing shell prerequisites without porting two mature gates.

## D3 - Use eight explicit gates and fail closed

**Decision**: The driver exposes `format`, `vet`, `lint`, `race`, `gui`, `coverage`, `docs`, `automation`, `list`, and `all`. `all` runs the eight gates sequentially in that order, prints a stable heading before each, stops on first failure, and never rewrites files.

**Rationale**: Explicit modes let CI retain job and platform parallelism while sharing exact commands with local verification. Sequential aggregate execution matches the autopilot protocol's foreground-output rule. A `list` mode provides a machine-readable manifest for the independent automation check.

**Alternatives considered**:

- Parallel local gates: faster, but mixes buffered output and conflicts with the documented requirement to watch foreground completion.
- Continue-on-error summary: more diagnostics per run, but makes the driver more complex and weakens the simple fail-closed contract.

## D4 - Validate automation offline with an explicit allowlist and manifest

**Decision**: Add `scripts/automation-check.sh [repo-root]`. It scans every workflow `uses:` reference, requires it to match the four audited action-major allowlist entries, rejects Node 20-era or unknown references, and compares `verify.sh list` against the eight required gates. Fixture tests operate against temporary repository roots.

**Rationale**: Resolving arbitrary marketplace action runtime metadata would require network access and would make ordinary CI dependent on remote source inspection. An allowlist turns every new action or major into an explicit review event. Keeping the required manifest outside the driver prevents `verify.sh` from silently omitting a gate and certifying its own reduced list.

**Alternatives considered**:

- Network-resolve every `action.yml` on each run: accurate but slow, brittle, rate-limited, and unnecessary after an audited version choice.
- Search only for the four known old version strings: misses a newly added Node 20 action family.
- Put the expected manifest inside `verify.sh`: cannot detect simultaneous drift in the implementation and its self-description.

## D5 - Preserve current workflow behavior and action permissions

**Decision**: Change action major selectors and gate command routing only. Keep workflow triggers, permissions, matrices, environment, inputs, artifact names, release globs, and job boundaries unchanged.

**Rationale**: Issue #21 is runtime maintenance, not workflow redesign. Narrow diffs make compatibility review reliable and honor the pinned-artifact rule.

**Alternatives considered**:

- Consolidate CI jobs while touching the workflow: would reduce repetition but alter platform coverage, failure isolation, and runtime behavior.
- Add unrelated action hardening or permissions changes: valuable candidates for separate issues, but outside the acceptance criteria.

## D6 - Treat unavailable prerequisites as failures, not skips

**Decision**: Let shell/command discovery and child commands fail with a named gate heading and actionable prerequisite guidance. For the linter, derive the module's Go language version for `GOTOOLCHAIN` when the caller has not already set it, preserving the documented old-base-toolchain workaround without a second hardcoded version.

**Rationale**: The constitution forbids reporting an unrun safety gate as green. Deriving the toolchain value from `go.mod` prevents the Makefile, docs, and module directive from drifting.

**Alternatives considered**:

- Skip unavailable race or shell gates: constitution violation.
- Hardcode `go1.25.0` in the driver: works now but creates another version source.
- Require a preinstalled golangci-lint binary: recreates the mismatch #41 asks to remove.
