# Verification: v1.0.0 Release Execution and Audit

## Immutable boundary

- Repository: `shruggietech/go-schedule`.
- Reviewed S049 merge and peeled `v1.0.0` commit:
  `ff47b4410d1aecbfadb8165d1ebf025ca1dadde7`.
- Annotated remote tag object:
  `35a232c03f528d14335d042de8f46f62baa4e30e`.
- S050 audit branch: `codex/050-v100-release-execution`.
- The audit branch changes no candidate runtime, packaging, or workflow file.

## Spec Kit analysis

The pre-execution analysis covered 17 functional requirements, six success
criteria, and 34 tasks with 100 percent mapping. It found zero critical or high
findings and no constitution conflict. The original plan remained deliberately
fail-closed around formal attended evidence.

After staging, the maintainer explicitly directed publication following an
exact-candidate install. That decision is recorded as an execution waiver, not
as retroactive satisfaction of the original formal-evidence requirements.

Post-deviation analysis found no ambiguous task identity, unmapped requirement,
or internal candidate-identity conflict. It found one intentional release-
contract nonconformance: FR-005 through FR-013 and SC-002 through SC-005 could
not all pass after publication was ordered without formal evidence. The
specification, plan, tasks, verification, changelog, release notes, and #122
audit now state that exception consistently. All 34 task identities remain
unique at their definition sites; 17 functional requirements and six success
criteria remain preserved rather than rewritten to make the exception appear
compliant.

## Tag staging

Release workflow
[33838072246](https://github.com/shruggietech/go-schedule/actions/runs/33838072246)
was triggered by the `v1.0.0` tag push with head commit
`ff47b4410d1aecbfadb8165d1ebf025ca1dadde7`. All six jobs passed:

1. README release metadata preflight;
2. absent-or-draft release guard;
3. daemon and CLI package staging;
4. Windows GUI and MSI staging;
5. macOS GUI staging;
6. Linux GUI staging.

Draft release ID `382483152` contained exactly the seven platform packages and
`windows-candidate-manifest.json`. All were downloaded independently into the
new fixed-volume workspace `A:/_tmp/go-schedule-v1.0.0-s050`.

## Candidate identity

Production `windows-release-gate verify-candidate` passed with:

- repository: `shruggietech/go-schedule`;
- tag: `v1.0.0`;
- commit: `ff47b4410d1aecbfadb8165d1ebf025ca1dadde7`;
- workflow: `Release`;
- run/attempt: `33838072246` / `1`;
- filename: `go-schedule_v1.0.0_windows_amd64.msi`;
- bytes: `24154112`;
- SHA-256:
  `433e9a311228d2e1dc7d5d7cae0ec6cee1831de03f9a52f2505e5a6cda552c8b`;
- ProductVersion: `1.0.0`;
- ProductCode: `{C46BE695-74CD-4323-8E7C-B2F4DC315B15}`.

The attended collector initialized a fresh candidate-bound workspace and
generated all 47 unavailable placeholders plus 27 lifecycle/desktop templates.
No placeholder was imported as a pass and no incomplete archive was finalized.

## Native installed checks and maintainer decision

The maintainer ran through installation of the exact staged MSI and reported
that everything was fine, then explicitly directed publication. Independent
post-install inspection confirmed:

- installed product display version `1.0.0`;
- `gosched version v1.0.0`;
- `goschedd` present, running, and configured as `LocalSystem`;
- ordinary-user `gosched health` returned
  `daemon ok (version v1.0.0)`;
- ordinary-user `gosched task list --json` completed successfully with an empty
  task set;
- the installed GUI ran in the interactive user session from
  `C:/Program Files/go-schedule/gosched-gui.exe`.

These checks and the maintainer's approval are real exact-candidate evidence,
but they are not represented as the unperformed 47-observation archive.

## Publication exception

The planned `Promote Release` workflow could not pass without a completed
formal archive. Following the maintainer's direct publication instruction, S050
published the existing green draft through the GitHub release operation instead
of manufacturing evidence or weakening the tagged workflow.

Before publication, `SHA256SUMS` was generated from and verified against all
eight unchanged staged payloads. Release notes were corrected to describe the
evidence actually available and no longer claim formal candidate-bound evidence.
No package was rebuilt, replaced, renamed, or resigned.

The exception and its remaining work are recorded on
[#122](https://github.com/shruggietech/go-schedule/issues/122#issuecomment-5536182541).
No formal disposition packet was created, no evidence-dependent leaf was
closed, #96 remains open, #122 remains open, and milestone #2 remains open.

## Public release audit

Release [v1.0.0](https://github.com/shruggietech/go-schedule/releases/tag/v1.0.0)
was published at `2026-09-04T05:32:59Z`. Both the tag-specific and latest-release
APIs resolve to that release, report draft `false` and prerelease `false`, and
identify tag `v1.0.0`.

A fresh download into
`A:/_tmp/go-schedule-v1.0.0-s050/public-audit` produced these nine non-empty
assets:

| Asset | Bytes | SHA-256 |
| --- | ---: | --- |
| `go-schedule_v1.0.0_darwin_amd64.tar.gz` | 9,394,029 | `771d81def93f77b0752abeccfad8977c532709ff6a57ed7024f986c2c934d7e1` |
| `go-schedule_v1.0.0_darwin_arm64.tar.gz` | 8,960,598 | `b3d0fad04be1b0509a9983d72dcfebedce583532267fd5479a7068d1fd6b8bd8` |
| `go-schedule_v1.0.0_linux_amd64.tar.gz` | 9,355,632 | `78cb3361c53e88cbdf4867e027bd16947706da3ebee4e0d3d469d751ca82d707` |
| `go-schedule_v1.0.0_linux_arm64.tar.gz` | 8,741,915 | `d5a5418d1bda17250b5d4323cbb92cc776bb0f4b58ea51825053a01d6115676f` |
| `go-schedule_v1.0.0_windows_amd64.msi` | 24,154,112 | `433e9a311228d2e1dc7d5d7cae0ec6cee1831de03f9a52f2505e5a6cda552c8b` |
| `go-schedule-desktop_v1.0.0_darwin_arm64.tar.gz` | 22,925,485 | `47894f91c6ce9f842d715270dcbebe0a2a8d9d0d967ff6830e2b644d249e8322` |
| `go-schedule-desktop_v1.0.0_linux_amd64.tar.gz` | 23,923,734 | `3c0a71ed1b866fd9e1505aeecebed157d40723c3c0da8936f415ac2509dfca64` |
| `windows-candidate-manifest.json` | 449 | `334bf4f9ed0e2a98f3982faf1925ff7c379ca2f4cb1c8e0ab10a6308fb489d99` |
| `SHA256SUMS` | 844 | `bbed47fba351c9c9b60de7ba69f7e2f744f580ee1133b4d70166b814e351251e` |

All eight payload entries in `SHA256SUMS` passed after the fresh public
download. The README uses the latest-public-release badge, the tagged changelog
contains the v1.0.0 boundary, and the release notes link to that boundary.

## GitHub planning state

At publication audit time, milestone #2 had 20 closed and 11 open issues. The
open set was #96, #98, #101, #104, #105, #106, #109, #111, #112, #113, and
#122. This is intentional: publication occurred by explicit exception, while
formal-evidence closure criteria remain unresolved and visible.

## Repository verification

Focused release validation passed on 2026-09-04:

- `go test ./internal/releasegate ./scripts/windows-release-gate -count=1`;
- `go test ./test/integration -run TestWindowsReleaseGate -count=1`;
- `scripts/spec-lifecycle-check.sh .` for all 50 specifications;
- `git diff --check`.

The complete foreground `scripts/verify.sh all` run passed all eight canonical
gates on 2026-09-04:

1. format: PASS;
2. vet: PASS;
3. lint: PASS with zero issues;
4. race: PASS, including the 81.934-second integration suite;
5. GUI: PASS;
6. coverage: PASS (`engine` 86.4%, `schedule` 89.2%, `timezone` 91.3%,
   `store` 84.1%, `catchup` 88.9%, and `logbus` 91.1%);
7. documentation: PASS for 15 pages, links, front matter, fences, theme,
   current product policy, and positive/negative fixtures;
8. automation: PASS for Actions, CodeQL, Dependabot, release operations,
   release notes, brand, lifecycle, and the canonical-gate contract.

Final integrity audit: PASS.

- All 14 changed or added text files decode as strict UTF-8 without BOM.
- Mojibake and unresolved TODO, TBD, FIXME, and template-placeholder scans
  found zero matches.
- All 34 task definitions are correctly formatted and uniquely identified.
- The lifecycle checker accepted all 50 specifications.
- Documentation verification covered Markdown links and current product policy.
- The branch changes only Spec Kit metadata, S050 audit documentation, and the
  Unreleased changelog; candidate runtime, packaging, and workflow paths have
  zero changes.
- `git diff --check` reports no whitespace errors.
