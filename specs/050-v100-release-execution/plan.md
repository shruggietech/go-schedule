# Implementation Plan: v1.0.0 Release Execution and Audit

**Branch**: `codex/050-v100-release-execution` | **Date**: 2026-09-04 | **Spec**: [spec.md](spec.md)

## Execution disposition

The maintainer later directed immediate publication after installing and
approving the exact candidate. Execution therefore preserved the immutable
candidate and public checksum audit but waived the formal archive, disposition,
and promotion-workflow stages. The deviation, remaining open issue state, and
absence of formal evidence are recorded in `verification.md`; no missing result
is relabeled as passing.

**Input**: Feature specification from `specs/050-v100-release-execution/spec.md`

## Summary

Execute issue #122 against the immutable `v1.0.0` tag at the reviewed S049
commit. Accept exactly one tag-triggered draft, verify its candidate manifest
and MSI, collect and independently validate all 47 attended Windows
observations, reconcile ten readiness issues from individual evidence, promote
the existing draft, and commit only a post-release specification and immutable
audit record to the S050 review branch.

## Technical Context

**Language/Version**: Go 1.25.7 release validators; PowerShell 7 attended collector; GitHub Actions and GitHub CLI for remote release operations

**Primary Dependencies**: Existing standard-library `internal/releasegate`; WiX 6.0.2 candidate metadata; GitHub Releases and Actions APIs; no new dependencies

**Storage**: Draft/public GitHub release assets; fixed-volume local evidence workspace outside the repository; Markdown audit record in `specs/050-v100-release-execution/`

**Testing**: Production candidate and bundle validators; disposition renderer; exact asset/checksum audits; canonical `scripts/verify.sh all`; hosted CI and Codex review

**Target Platform**: Windows 11 desktop for formal attended evidence; GitHub-hosted Linux for promotion; repository builds remain Linux, macOS, and Windows compatible

**Project Type**: Release-operations and audit slice for a Go desktop/daemon/CLI project

**Performance Goals**: No product hot-path change; every remote wait is bounded and every validator completes before the next mutation

**Constraints**: Immutable tag and candidate bytes; exactly 47 passing observations; exact 8/9/10 asset cardinalities; no local-demo substitution; no release rebuild; UTF-8 without BOM; no visible console child windows from project automation

**Scale/Scope**: One tag, one accepted staging run, one MSI, one evidence archive, ten issue records, ten final public assets, eleven coordinated v1 issues including #122

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Gate | Result |
| --- | --- | --- |
| I. Code Quality | Use reviewed production tooling, explicit diagnostics, immutable identities, and no recoverable panic path | PASS. S050 changes no product code and records every failure boundary. |
| II. Testing Standards | Run exact candidate/bundle validation, asset/checksum audit, race-aware canonical gates, and no weakened checks | PASS. Formal evidence is stricter than fixtures and cannot be bypassed. |
| III. User Experience Consistency | Publish only the MSI already accepted through the complete native Windows matrix | PASS. All user-visible v1 surfaces require formal evidence. |
| IV. Performance Requirements | Avoid product performance change and preserve existing benchmark gate | PASS. Canonical verification retains benchmark evidence. |
| V. Autonomous Execution | Use Spec Kit, review branch, PR, hosted CI/review, and explicit release authorization | PASS. The tag was authorized before the audit branch; conditional publication is explicitly authorized. |

Post-design re-check: PASS. The design adds no dependency, runtime code, schema,
workflow, or candidate mutation. The only irreversible operations are already
specified as downstream results of complete fail-closed gates.

## Phase 0 Research Decisions

Research is captured in [research.md](research.md). The governing decisions are:

1. Treat the authorized S049 tag staging as an immutable prerequisite and keep
   S050's review branch outside candidate identity.
2. Use only the existing production validators and attended collector; never
   create a shortcut evidence format or promote fixture data.
3. Keep the evidence workspace outside the repository and commit only hashes,
   immutable URLs, identifiers, results, and redacted environment facts.
4. Reconcile leaf issues before coordinators, and keep #122 open through the
   public-release audit.
5. Dispatch the existing Promote Release workflow rather than directly editing
   the draft, so asset cardinality, provenance, and checksum enforcement remain
   centralized.

## Phase 1 Design

- [data-model.md](data-model.md) defines candidate, evidence, disposition,
  promotion, and audit state.
- [contracts/release-audit.md](contracts/release-audit.md) defines the exact
  S050 evidence that may be committed.
- [quickstart.md](quickstart.md) defines the chronological operator path and
  stop conditions.
- [checklists/release-gate.md](checklists/release-gate.md) tests whether release
  requirements are complete and unambiguous.

## Project Structure

### Documentation (this feature)

```text
specs/050-v100-release-execution/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── verification.md
├── contracts/
│   └── release-audit.md
├── checklists/
│   ├── requirements.md
│   └── release-gate.md
└── tasks.md
```

### Existing operational surfaces

```text
.github/workflows/
├── release.yml
└── promote-release.yml

internal/releasegate/
├── validate.go
├── archive.go
└── disposition.go

scripts/windows-release-gate/
└── main.go

test/windows/
├── Invoke-ReleaseCandidateAttended.ps1
└── README.md
```

**Structure Decision**: S050 reuses the reviewed release implementation without
modification. The repository receives only the Spec Kit tree, inventory pointer,
changelog note, and final audit record. Candidate MSI, manifest, formal archive,
and issue packet remain external artifacts identified by immutable hashes and
URLs.

## Complexity Tracking

No constitution violation or added implementation complexity is required.
