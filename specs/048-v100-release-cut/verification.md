# Verification: v1.0.0 Release Cut

## Baseline (2026-09-03)

- Repository: `shruggietech/go-schedule`.
- Starting commit: `d4ad27540331eae943e43a3830d4f5c9a6c38afa`, the merged
  S047 pull request #121, on clean synchronized `main`.
- Working branch: `codex/048-v100-release-cut`.
- Latest local tag and public GitHub release: `v0.9.1`.
- Neither a local v1.0.0 tag nor a GitHub v1.0.0 release existed.
- No pull request was open when S048 began.
- Milestone #2, `v1.0.0 - Release readiness`, was open with ten open
  verification-dependent issues.
- Release coordinator issue #122 was created for S048 under milestone #2.
- The pre-cut Unreleased section contained 33 top-level entries and 16,665
  characters. Its LF-normalized SHA-256 was
  `bf2b8f2fcedd2f099c28531ee29f0d6ddb3f55a04e0d2878cd7a349a8ab5af2d`.
- Twelve commits were present in `v0.9.1..d4ad275`.
- The existing workflows already stage drafts and require exact-candidate
  attended evidence before promotion; no workflow modification was planned.

## Expected v1.0.0 staged payloads

1. `go-schedule_v1.0.0_darwin_amd64.tar.gz`
2. `go-schedule_v1.0.0_darwin_arm64.tar.gz`
3. `go-schedule_v1.0.0_linux_amd64.tar.gz`
4. `go-schedule_v1.0.0_linux_arm64.tar.gz`
5. `go-schedule-desktop_v1.0.0_darwin_arm64.tar.gz`
6. `go-schedule-desktop_v1.0.0_linux_amd64.tar.gz`
7. `go-schedule_v1.0.0_windows_amd64.msi`
8. `windows-candidate-manifest.json`
9. `go-schedule_v1.0.0_windows-attended-evidence.zip`

`SHA256SUMS.txt` is generated during promotion and must contain exactly nine
passing payload entries.

## Authorization and lifecycle boundary

- The current user authorization covers Spec Kit, local implementation,
  verification, commit, push, one official S048 PR, in-scope review-fix pushes,
  hosted CI, and at most one explicit second Codex review round.
- It does not explicitly authorize the separate v1.0.0 tag or release mutation
  required by constitution Principle V and the project autopilot contract.
- Therefore S048 implementation completes at the reviewed preparation PR. The
  tag-triggered draft, formal evidence, issue closure, promotion, and milestone
  closure remain chronological post-merge release work under issue #122.

## Spec Kit analysis

PASS after one asset-cardinality wording correction.

- Functional requirements: 18/18 mapped to implementation or verification
  tasks (100 percent).
- Success criteria: 6/6 objectively testable and mapped (100 percent).
- Tasks: 23 unique, chronological identities with no duplicates.
- User stories: 3/3 independently testable.
- Requirement checklist: 20/20 complete.
- Release-contract checklist: 22/22 complete.
- Unresolved clarification markers or template placeholders: 0.
- Constitution conflicts: 0.
- Critical findings: 0.
- High findings: 0.
- One medium finding was corrected before implementation: the draft initially
  described eight build outputs plus the manifest, while the workflow actually
  stages seven build packages plus the manifest before attended evidence becomes
  the ninth pre-checksum payload.
- Unmapped requirements or tasks: 0.

## Preparation verification

Focused checks: PASS.

- Changelog retention: 33/33 original top-level entries retained. Replacing the
  new v1.0.0 heading with the prior Unreleased heading reproduces baseline
  SHA-256 `bf2b8f2fcedd2f099c28531ee29f0d6ddb3f55a04e0d2878cd7a349a8ab5af2d`.
- Future Unreleased content: zero non-heading, nonblank lines.
- README identity: exactly one `release-v1.0.0-58A6FF` badge and exactly one
  `daemon ok (version 1.0.0)` example.
- Release notes: exactly five bullets, one Highlights heading, and one tagged
  v1.0.0 changelog link.
- `test/scripts/automation-check_test.sh`: PASS for all positive and negative
  release-policy fixtures.
- `scripts/automation-check.sh`: PASS with staging/promotion, release notes,
  lifecycle, and eight-gate contracts intact.
- Release-gate package and CLI tests: PASS.
- Four exact workflow/collector/inspector integration contracts: PASS.
- No release workflow, promotion workflow, validator, product source, package
  format, supported platform, or dependency was changed.

Canonical `scripts/verify.sh all`: PASS on 2026-09-03.

- Format: PASS.
- Vet: PASS.
- Lint: PASS with zero issues.
- Race: PASS for all selected packages, including integration tests.
- GUI: PASS.
- Coverage: PASS (`engine` 86.4%, `schedule` 89.2%, `timezone` 91.3%,
  `store` 84.1%, `catchup` 88.9%, and `logbus` 91.1%).
- Documentation: PASS for 15 pages, links, front matter, fences, theme, and
  current product policy, including positive and negative policy fixtures.
- Automation: PASS for Actions, CodeQL, Dependabot, release staging and
  promotion, release notes, brand, lifecycle, and the eight-gate contract.

Final integrity audits: PASS.

- All 16 changed or added files decode as strict UTF-8 and have no BOM.
- Mojibake and unresolved-template-placeholder scans found zero matches.
- All 23 task identifiers are unique, correctly formatted, and complete.
- The release-copy cardinality and exact changelog reconstruction checks pass.
- `git diff --check` reports no whitespace errors.
- `scripts/spec-lifecycle-check.sh .` reports all 48 specifications
  lifecycle-consistent.
- Spec Kit prerequisite discovery resolves S048 and all required design
  artifacts successfully.
- No local v1.0.0 tag or GitHub v1.0.0 release was created.
