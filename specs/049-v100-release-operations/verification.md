# Verification: v1.0.0 Release Operations

## Baseline (2026-09-03)

- Repository: `shruggietech/go-schedule`.
- Starting commit: `33d53298bc8103ea3b70985b4a08697b442b0431`, the merged
  S048 pull request #123, on clean synchronized `main`.
- Working branch: `codex/049-v100-release-operations`.
- Latest local tag and public GitHub release: `v0.9.1`.
- Neither a local/remote `v1.0.0` tag nor a GitHub v1.0.0 release existed.
- No pull request was open when S049 began.
- Milestone #2, `v1.0.0 - Release readiness`, had eleven open issues:
  #96, #98, #101, #104, #105, #106, #109, #111, #112, #113, and #122.
- S047 pull request #121 and S048 pull request #123 were merged.
- Existing production validation required 47 observations and the Release and
  Promote Release workflows already preserved draft staging and exact-candidate
  promotion.

## Clarification Decisions

- The reviewed S049 merge commit supersedes the S048 merge commit as the final
  `v1.0.0` tag target. A required S049 pull request advances `main`, making the
  old `HEAD == S048 merge` preflight impossible. S049 is restricted to release
  tooling and documentation, so packaged runtime behavior remains unchanged.
- Generated issue records are offline review input. They do not comment, close,
  tag, upload, dispatch, promote, or change milestone state.

## Spec Kit Analysis

PASS after two corrections.

- Functional requirements: 18/18 mapped to implementation or verification
  tasks (100 percent).
- Success criteria: 8/8 mapped to implementation or verification tasks
  (100 percent).
- Tasks: 31 unique chronological identities.
- User stories: 3/3 independently testable.
- Requirements checklist: 16/16 complete.
- Release-operations checklist: 27/27 complete.
- Unresolved clarification markers or template placeholders: 0.
- Constitution conflicts: 0.
- Critical findings: 0.
- High findings: 2, both corrected before implementation.
- Unmapped requirements or implementation tasks: 0.

Corrections:

1. The first mapping draft gave #98 only seven Removal observations and treated
   #96 as coordinator for all desktop issues. S047's authoritative model gives
   #98 all 16 Setup/Removal observations and #96 aggregate lifecycle
   reconciliation. S049 now maps #96 to the 36 pre-desktop observations and its
   actual child/prerequisite references #97, #98, #94, #89, and #90.
2. The generated task graph initially placed `/speckit-analyze` after
   implementation. It was moved before the In Progress transition and red-test
   phase, restoring the constitution's blocking-gate order.

## Red-to-Green Evidence

Red phase: PASS on 2026-09-03.

- `internal/releasegate` failed to compile only because
  `IssueDispositionMappings`, `RenderDispositionPacket`,
  `WriteDispositionPacket`, and its injectable write seam did not exist.
- `scripts/windows-release-gate` tests reached the existing usage path and
  rejected `render-dispositions` as an unknown command.
- The integration contract failed because command routing and renderer source
  did not exist.
- No production implementation file was present when these failures were
  captured.

Green phase: PASS on 2026-09-03.

- The renderer emits the fixed ten-issue mapping plus `packet.json`, preserves
  candidate/workflow/environment/observation/metric/attachment traceability,
  suppresses Markdown mentions, and produces identical bytes in two distinct
  output directories.
- The writer renders before mutation, refuses existing targets and linked
  parent paths, uses a private sibling staging directory, cleans injected write
  failures, and commits the complete packet through one rename.
- The CLI requires every identity and input option and reuses `Validate`,
  `ValidateBundleContents`, and `ValidateCandidateManifest` before writing.
- Formal evidence-class, candidate-field mutation, output-conflict, help,
  stream, exit-code, no-network, and production-validator-reuse contracts pass.

## Focused Verification

PASS on 2026-09-03.

- `go test ./internal/releasegate ./scripts/windows-release-gate -count=1`
- `go test -race ./internal/releasegate ./scripts/windows-release-gate -count=1`
- `go test ./test/integration -run 'TestWindowsReleaseGate' -count=1`
- `sh test/scripts/automation-check_test.sh automation`
- `sh scripts/spec-lifecycle-check.sh` (49 lifecycle-consistent
  specifications)

The automation fixture includes positive S049 boundary and command assertions
plus negative stale-boundary, missing-command, and missing-non-authority cases.

## Canonical Verification

PASS on 2026-09-03.

`sh scripts/verify.sh all` completed as one uninterrupted green run:

1. format: PASS
2. vet: PASS
3. lint: PASS (0 issues)
4. race: PASS, including `internal/releasegate`, the CLI, and integration
5. GUI: PASS
6. coverage: PASS (engine 86.4%, schedule 89.2%, timezone 91.3%, store
   84.1%, catchup 88.9%, logbus 91.1%)
7. docs: PASS (15 pages, links, front matter, fences, theme, and policy)
8. automation: PASS (release operations and the established automation
   baseline)

## Integrity Audit

PASS on 2026-09-03.

- Strict UTF-8 decoding and UTF-8-without-BOM inspection passed for all 23
  changed or added files.
- Mojibake scan passed.
- No unresolved clarification, placeholder, TODO, or TBD markers remain on
  authored S049 requirement, design, task, or verification surfaces. The one
  checklist sentence that names the forbidden marker is itself checked and was
  excluded from the semantic marker scan.
- The lifecycle/task-format audit passed for all 49 specifications.
- Packet tests verified the exact eleven-file inventory and byte-identical
  output in separate absent destinations.
- The canonical documentation gate verified all links and authored surfaces.
- `git diff --check` passed.
