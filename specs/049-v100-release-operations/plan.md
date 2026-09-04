# Implementation Plan: v1.0.0 Release Operations

**Branch**: `codex/049-v100-release-operations` | **Date**: 2026-09-03 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from
`specs/049-v100-release-operations/spec.md`

## Summary

Add a fail-closed `render-dispositions` operation to the existing Windows
release-gate command. It will reuse the production candidate and evidence
validators, then atomically write ten deterministic GitHub-Markdown records and
one machine-readable packet index for the v1.0.0 readiness issues. The slice
also revises the release contract so the reviewed S049 merge commit is the final
tag boundary and documents the exact post-merge staging, packet review,
promotion, and audit sequence. The pull request performs no remote release
mutation.

## Technical Context

**Language/Version**: Go 1.25.0; Markdown documentation; GitHub Actions YAML

**Primary Dependencies**: Go standard library; existing `internal/releasegate`
validation and bounded ZIP extraction; existing `windows-release-gate` command

**Storage**: Read-only MSI, JSON manifest, and ZIP inputs; one local atomic
directory containing ten Markdown files plus `packet.json`

**Testing**: Go unit and CLI tests, integration contract tests, release-policy
shell tests, `go test -race`, and `scripts/verify.sh all`

**Target Platform**: Generator supports the repository's Windows and Linux
development environments; evidence remains Windows 11-only

**Project Type**: Maintainer CLI and release-operations documentation inside a
cross-platform desktop application repository

**Performance Goals**: Generate the packet within 60 seconds; render time is
linear in the already bounded 47 observations and attachment inventory

**Constraints**: No network access; no GitHub mutation; no new dependency; no
overwrite; no partial output; strict UTF-8 without BOM; deterministic output;
reuse rather than restate the production validator

**Scale/Scope**: One v1.0.0 candidate, 47 observations, ten issue records, nine
child issue mappings, one coordinator mapping, and one packet index

## Constitution Check

*GATE: Must pass before research and again after design.*

| Principle | Gate | Result |
| --- | --- | --- |
| I. Code Quality | Small focused package surface, contextual errors, documented exported API, no concurrency | PASS |
| II. Testing Standards | Red tests precede renderer and CLI implementation; mutation, race, full-suite, and canonical gates remain mandatory | PASS |
| III. User Experience Consistency | One verb-noun operation, stdout success, stderr diagnostics, conventional exit codes, actionable failures | PASS |
| IV. Performance Requirements | Work is bounded by existing archive limits and fixed 47-observation cardinality; no hot scheduler path changes | PASS |
| V. Autonomous Execution | Full Spec Kit sequence, review branch, PR, analyze gate, and separate tag/release authority | PASS |

Post-design re-check: PASS. The design introduces no dependency, runtime code,
background process, secret handling, persistence migration, or hosted mutation.
It routes all trust decisions through the existing production validator and
stages output atomically only after validation succeeds.

## Technical Design

### Validation and rendering boundary

The command loads the independent manifest and extracts the bounded evidence
archive through existing helpers. It evaluates `Validate`,
`ValidateBundleContents`, and `ValidateCandidateManifest` with mandatory
repository, tag, and commit expectations. Rendering is unreachable until the
combined failure set is empty.

The renderer accepts only an already validated `Evidence` value. It exposes the
canonical issue mapping for direct unit testing, assembles all files in memory,
writes them into a private sibling staging directory, and renames that directory
onto an absent destination. Any validation, rendering, write, or rename error
removes the staging directory and leaves the requested destination absent.

### Deterministic packet

The output inventory is fixed:

```text
packet.json
issue-096.md
issue-098.md
issue-101.md
issue-104.md
issue-105.md
issue-106.md
issue-109.md
issue-111.md
issue-112.md
issue-113.md
```

`packet.json` has a schema version, candidate identity, and ordered issue/file
index. Markdown records share one candidate table, production-validator
statement, immutable workflow and evidence-archive links, relevant environment
table, and ordered observation sections. User-controlled strings are escaped in
tables and encoded JSON metrics are indented, preventing Markdown structure or
mention injection.

### Release boundary

S049 deliberately supersedes only S048's commit-equality condition. The
annotated `v1.0.0` tag will target the reviewed S049 merge commit because the
required PR advances `main`. The packaged application code and behavior are
unchanged; the delta is release tooling, tests, specifications, and
documentation. All other S048 staging, evidence, promotion, asset, checksum,
and audit requirements remain authoritative.

## Project Structure

### Documentation (this feature)

```text
specs/049-v100-release-operations/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── verification.md
├── contracts/
│   └── disposition-packet.md
├── checklists/
│   ├── requirements.md
│   └── release-operations.md
└── tasks.md
```

### Source Code (repository root)

```text
internal/releasegate/
├── disposition.go
└── disposition_test.go

scripts/windows-release-gate/
├── main.go
└── main_test.go

test/
├── integration/windows_release_gate_contract_test.go
└── windows/README.md

specs/048-v100-release-cut/contracts/publication.md
specs/README.md
CHANGELOG.md
CLAUDE.md
```

**Structure Decision**: Extend the existing release-gate package and command.
A second executable or network-aware release orchestrator would duplicate
validation, increase credential risk, and provide no independent value.

## Implementation Phases

1. Establish S049 lifecycle, exact GitHub baseline, boundary decision, and
   specification artifacts.
2. Write failing issue-mapping, deterministic-rendering, escaping, atomicity,
   invalid-input, and CLI contract tests.
3. Implement canonical disposition models, Markdown rendering, JSON index, and
   atomic packet writing in `internal/releasegate`.
4. Add the `render-dispositions` CLI path that requires all identity inputs and
   reuses the complete production validation chain.
5. Update the Windows runbook and S048 publication contract for packet use and
   the reviewed S049 merge boundary.
6. Run Spec Kit analysis, focused tests, race tests, full canonical verification,
   encoding and corruption audits, then complete lifecycle evidence and commit.

## Decision Log

- Use the existing release-gate executable, not a new script or workflow. This
  preserves one trust boundary and avoids Windows child-process launch risk.
- Generate local files but perform no GitHub API calls. Credentialed mutation
  remains deliberate and separately reviewable.
- Use Markdown records plus a JSON index. Markdown is directly consumable by
  issue comments; JSON provides exact inventory and candidate binding.
- Use a fixed compiled issue mapping. These issue numbers and scenarios define
  the v1.0.0 gate and should fail review visibly when changed.
- Stage into a sibling temporary directory and rename only after every file is
  complete. Per-file direct writes would leave ambiguous partial evidence.
- Move the tag boundary from S048 to S049 because merging the required S049 PR
  makes S048's `main == S048 merge` precondition impossible. Restricting S049 to
  release operations keeps the product candidate unchanged.
- Do not add the generic development skill's public blog post. This is internal
  release tooling, the repository has no product-blog surface, and adding one
  would expand scope beyond issue #122; the runbook and changelog are the
  project-standard documentation surfaces.
- Do not run the development skill's multi-agent phases. Session-level
  instructions prohibit delegation unless explicitly requested; the same
  discovery, design, review, and validation phases are performed locally.

## Complexity Tracking

No constitution violation requires justification.
