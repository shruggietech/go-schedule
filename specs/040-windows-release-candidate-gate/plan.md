# Implementation Plan: Windows Release Candidate Gate

**Branch**: `codex/040-windows-release-candidate-gate` | **Date**: 2026-09-02 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/040-windows-release-candidate-gate/spec.md`

## Summary

Replace immediate tag-triggered publication with a build-once, two-phase release path. The tag workflow stages all assets in a draft release and emits exact Windows candidate identity. A separate manual promotion workflow retrieves that same MSI and a hash-bound attended Windows 11 evidence bundle, runs one cross-platform validator, creates final all-asset checksums, and only then publishes. A resumable PowerShell collector supplies native HWND, DPI, token, process, and operator-observation fragments without allowing hosted CI or fixture data to masquerade as attended proof.

## Technical Context

**Language/Version**: Go 1.25.0; PowerShell 7; POSIX shell for repository automation; GitHub Actions YAML

**Primary Dependencies**: Go standard library; Windows Installer COM and Win32 APIs for local capture; existing WiX 6.0.2/Fyne 2.8.1/product dependencies; GitHub CLI on hosted runners

**Storage**: Versioned JSON evidence bundle and hashed attachments; draft GitHub release assets; no product database changes

**Testing**: Go table and mutation tests, PowerShell parser/compliance tests, shell automation fixtures, existing compiled-MSI and hosted Windows lifecycle jobs, canonical eight-gate verification

**Target Platform**: Cross-platform validator and Linux promotion runner; Windows 11 client for attended collection; current hosted Linux/macOS/Windows release builders

**Project Type**: Desktop application release tooling and CI/CD control

**Performance Goals**: Validate a normal evidence bundle and candidate in under 10 seconds excluding network retrieval; collect two-minute observations in one process without process-per-sample polling

**Constraints**: Exact MSI bytes only; fail closed on every non-pass state; no speculative GUI/error fix; no secret or unnecessary personal data; hidden console children; UTF-8 without BOM; no release, tag, or merge during S040 implementation

**Scale/Scope**: One candidate manifest, 36 fixed attended scenarios, two or more display environments, three account roles, one evidence archive, two release workflows, and one cross-platform semantic validator

## Constitution Check

| Principle | Gate | Design evidence |
| --- | --- | --- |
| I. Code Quality | PASS | Evidence decoding, candidate hashing, scenario validation, attachment verification, and diagnostics have separate testable boundaries; unknown fields and ambiguous duplicates are rejected. |
| II. Testing Standards | PASS | Validator and automation mutation tests are written failing first; the existing MSI, LocalSystem, race, GUI, coverage, docs, and automation gates remain intact. |
| III. UX Consistency | PASS | The operator flow is phase-based, resumable, explicit about non-pass states, and uses the same fixed scenario names in schema, runbook, diagnostics, and promotion. |
| IV. Performance | PASS | No scheduler runtime changes; hashing is streaming, validation is linear in observations/attachments, and native sampling reuses one process. |
| V. Autonomous Execution | PASS WITH RECORDED OPERATOR OVERRIDE | Full Spec Kit and analyze gates remain mandatory. The operator explicitly authorized automatic S040 push and PR plus at most one additional Codex review round; merge, tag, and release remain outside authority. |

### Post-design re-check

All principles remain satisfied. Changes to pinned release workflows, automation checks, and Windows test documentation are directly required by #94 and will receive dated changelog decisions. The two-phase draft/promotion path is additional release machinery, but it is the smallest design that can prove and later publish the exact tested MSI. Product runtime behavior is intentionally unchanged without native reproduction.

## Architecture and Decision Log

### Stage every release asset without publication

The existing tag workflow remains the build authority but every release upload explicitly sets `draft: true`. Its checksum publication job is removed because the evidence archive does not yet exist. The Windows build writes a deterministic candidate manifest beside the MSI after compiled inspection. This preserves current platform packaging while ensuring the tag cannot become a public release before attended validation.

### Promote only an already-staged exact candidate

The new manual promotion workflow accepts one strict semantic-version tag. It inspects the existing draft through the GitHub API, verifies the target commit, downloads all assets, and requires the canonical MSI, candidate manifest, and attended evidence archive. It invokes the shared validator with the expected repository, tag, and tag commit. Only after validation does it generate `SHA256SUMS.txt` across the complete final asset set, upload that checksum, and make the draft public. Re-running against a public release fails safely.

A new workflow does not create a second release path. It is the sole transition from draft to public. Direct administrator actions in GitHub remain outside repository-enforceable guarantees and are documented honestly.

### Validate one canonical manifest

`internal/releasegate` owns strict JSON decoding, schema validation, identity comparison, scenario coverage, scenario-specific metrics, time bounds, safe attachment paths, and streaming hashes. `scripts/windows-release-gate` exposes local-directory and ZIP-bundle validation with a stable 0/1/2 exit contract. Validation collects all independent defects instead of stopping at the first, while true I/O or invocation failures remain distinct.

The checked-in positive fixture uses an inert fake artifact and representative text attachments. Mutation tests alter every critical identity, scenario, outcome, duration, dimension, environment, path, and digest rule. Fixture acceptance proves the validator, never the product.

### Collect native fragments without pretending to automate observation

`test/windows/Invoke-ReleaseCandidateAttended.ps1` initializes a bounded evidence workspace, captures native window/process measurements, imports explicit operator-reviewed observation fragments, and finalizes only through the shared validator. Its Win32 adapter enumerates the exact process HWND and records outer/client rectangles, monitor/work areas, DPI, state, session, and process identity. A narrowly scoped opt-in GUI evidence file records the exact Fyne canvas size and scale from the installed process; it remains distinct from the collector's independently derived native client size. Screenshots plus attended counts cover canvas overlays.

Existing lifecycle probes remain authoritative for installed access and production task execution. The collector references their structured attachments and adds fresh-process PATH and attended lifecycle observations. Machine mutation remains explicit and `ShouldProcess`-guarded; console child processes are hidden and redirected.

### Stop on native regression evidence

S040 contains no GUI sizing or connection presentation change by default. If future native execution finds a failure, the candidate evidence remains non-pass. A follow-up correction must retain the exact baseline reproduction and, for recurring errors, prove the uncommitted fix natively before the first correction commit. This is a deliberate deviation from speculative “fix while here” behavior because #94 makes evidence ordering part of acceptance.

## Project Structure

```text
.github/workflows/release.yml
.github/workflows/promote-release.yml
internal/releasegate/
├── model.go
├── validate.go
├── bundle.go
└── *_test.go
scripts/windows-release-gate/
├── main.go
└── main_test.go
test/fixtures/windows-release-gate/
├── passing/
└── mutations/
test/windows/Invoke-ReleaseCandidateAttended.ps1
gui/window_evidence.go
gui/window_evidence_test.go
test/windows/README.md
scripts/automation-check.sh
test/scripts/automation-check_test.sh
docs/INSTALL-windows.md
specs/040-windows-release-candidate-gate/
CHANGELOG.md
CLAUDE.md
specs/README.md
```

**Structure Decision**: Extend the established internal-package/repository-command pattern and the existing Windows evidence surface. The release gate is not part of the shipped application, and the attended collector remains under `test/windows`. No new external dependency or public product command is introduced.

## Complexity Tracking

| Added complexity | Why needed | Simpler alternative rejected because |
| --- | --- | --- |
| Draft staging plus manual promotion | Native testing must occur after exact MSI creation and before publication. | Rebuilding after evidence breaks byte identity; immediate tag publication cannot be tested in time. |
| Scenario-specific evidence validator | Missing, duplicated, stale, elevated-only, or underspecified observations must fail deterministically. | A signed checkbox or generic JSON schema cannot enforce cross-observation and measurement rules. |
| Operator-assisted native collector | HWND/DPI state and Fyne overlay visibility require the real interactive Windows desktop. | Hosted Windows Server, headless Fyne tests, and screenshots alone cannot satisfy #94. |
