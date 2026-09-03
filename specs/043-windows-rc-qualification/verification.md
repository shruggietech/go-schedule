# Verification: Windows Demo Qualification

**Date**: 2026-09-03

**Branch**: `codex/043-windows-rc-qualification`

## Current result

Automated qualification and the operator's attended local-demo walkthrough are
complete. The demo found usability defects, all of which remain outside S043
implementation and are tracked in GitHub. No release-blocking regression in the
setup, shortcut, Finish-action, double-click, Exit, guided-uninstall, or wipe
paths was reported. The operator authorized branch publication through the pull-
request review boundary; formal release-candidate evidence remains separate.

## Evidence classes

- `local-demo`: useful for pre-PR exploratory testing, never formal #94 proof.
- `automated`: repository, source, compiled-table, and headless behavior checks.
- `attended-demo`: operator observations against the exact linked local hash.
- `attended-windows-release-candidate`: later draft-release evidence required for formal release readiness.

## Analyze gate

PASS.

- Functional requirements: 14/14 mapped to implementation or handoff tasks.
- Success criteria: 6/6 mapped to automated or attended evidence.
- Unresolved clarification markers or template placeholders: 0.
- Duplicate or contradictory requirements: 0.
- Constitution conflicts: 0.
- Critical findings: 0.

The ordering tension between pre-PR testing and #94's workflow-staged exact
candidate is resolved through separate evidence classes. The demo cannot be
promoted or relabeled as formal evidence.

## Automated verification

| Check | Result |
| --- | --- |
| `go test -race ./internal/releasegate ./scripts/windows-release-gate ./internal/winuninstall ./internal/executor ./internal/api/server ./internal/api/client ./cmd/goschedd -count=1` | PASS |
| `go test ./gui -count=1` | PASS in 26.721 seconds |
| `go test ./test/integration -run 'WindowsInstaller|WindowsRelease|Cleanup' -count=1` | PASS |
| PowerShell parser check for six Windows build/evidence scripts | PASS |
| `pwsh -NoProfile -File build/windows/verify_wxs.ps1` | PASS |
| `scripts/verify.sh all` | PASS, all eight gates |

Canonical coverage remained above the required threshold: engine 86.4 percent,
schedule 89.2 percent, timezone 91.3 percent, store 84.1 percent, catchup 88.9
percent, and logbus 91.1 percent. Automation validated CI, CodeQL, Dependabot,
release staging/promotion, release notes, brand, lifecycle, and all eight gates.

The hosted destructive installer contract is unavailable on this development
host because it intentionally requires an elevated GitHub-hosted disposable
runner. Its safety check was not bypassed. Compiled inspection and operator-
driven attended testing cover the local demo boundary.

After the tracked documentation changes, the documentation and automation gates
were rerun independently and passed. This confirms the new Spec Kit inventory,
links, lifecycle state, changelog text, and release-policy descriptions are
accepted by the repository checks.

## Local review and integrity

The specification, issue traceability, evidence classes, build ordering, safety
boundaries, and changed documentation were reviewed against the constitution,
#94, #98, and the S040 through S042 verification history. No finding at or above
the 80-percent confidence threshold remains.

All changed or added text files decode as strict UTF-8, none has a UTF-8 BOM,
and no mojibake signature was found. The template and clarification scan found
no unresolved placeholder.

## Demo identity

The first local compilation succeeded, but its inspection report used the only
available non-published class, `candidate`. Although the origin explicitly said
it was not staged, the generated heading and status still called it a candidate
artifact. That report and MSI are superseded and will not be handed off.

A regression assertion was changed first and failed because the inspector did
not expose `[ValidateSet('candidate', 'published', 'local-demo')]`. After adding
only that provenance value, the focused integration contract and PowerShell
parser check passed. Candidate-manifest generation remains restricted to
`candidate`, and the published class retains its repository HTTPS-origin check.

The complete eight-gate suite then passed again with the correction in place.
The race gate included the 54.026-second integration suite; lint reported zero
issues, all coverage packages remained above 80 percent, and documentation plus
automation completed cleanly.

The corrected demo was rebuilt from commit
`cd877817fdfe3eac9ceeb7f8500f7ef7a55ceffa`. That commit remains reachable in
the review branch and is an ancestor of its final head. The later commits are
enumerated rather than described collectively as evidence-only:

- `3add766` changes only the S043 task and verification records.
- `576fbad` records the attended result in CHANGELOG, the S043 checklist, spec,
  tasks, verification, and specification inventory.
- `9baed96` removes extra EOF whitespace from four S043 Markdown files.
- The final review-response commit changes only this verification record to make
  the source and payload boundary explicit.

The runtime product code, compiled executables, WiX authoring, icons, README,
license, and every MSI input except CHANGELOG are unchanged after `cd877817`.
The tested MSI therefore contains the earlier CHANGELOG revision and is not
claimed to be byte-equivalent to a build from the final pull-request head.
Rebuilding it after the walkthrough would invalidate the operator-bound hash for
a documentation-only payload difference. The later formal candidate will be
staged from the merged reviewed source and must repeat #94's exact-artifact gate.

| Field | Value |
| --- | --- |
| Artifact | `dist/go-schedule_s043-demo_cd877817_windows_amd64.msi` |
| Artifact class | `local-demo` |
| Built at | `2026-09-03T09:17:01.782Z` |
| Source commit | `cd877817fdfe3eac9ceeb7f8500f7ef7a55ceffa` |
| Embedded version | `1.0.0-s043-demo.cd877817` |
| ProductVersion | `1.0.0` |
| ProductCode | `{593C2C48-6C61-4359-BB30-BFDBA2E59C7A}` |
| UpgradeCode | `{B6F3C2E1-7A4D-4C9E-9B2A-1F6D8E5A0C34}` |
| Byte size | `23511040` |
| SHA-256 | `ef4a869af0e6971d445c53e8f7a237df655245287c7ae64d8e719941bda8ad59` |
| Compiled inspection | PASS; `dist/s043-demo-cd877817-installer-inspection.md` |
| Authenticode | Not signed (local demo) |

Independent inspection queried the MSI Property table and matched the values
above. `gosched --version` returned the embedded demo version, and a byte scan
found the same marker in the CLI, daemon, GUI, and cleanup helper. The compiled
inspector proved summary metadata, Explorer Subject, canonical product identity,
icon and shortcut rows, PATH integration, local-group authoring, maintenance
routing, completion actions, preserve default, explicit wipe condition, invalid-
input guard, and close-application rows.

## Native evidence boundary

No unattended or headless result will be represented as proof of visual quality,
Windows Settings interaction, native window geometry, or user judgment. The demo
remains unofficial and cannot close #94, #98, or #96.

## Attended demo result

The operator installed the one linked S043 demo artifact and explicitly ended
testing on 2026-09-03. An independent operator-side hash computation and exact
Windows build/scale record were not supplied, so the observations are bound to
the handoff artifact but are not promoted to formal exact-candidate evidence.

Passed observations:

- setup was seamless;
- Desktop and Start Menu shortcuts worked;
- the selected Finish actions opened the app and documentation together;
- task-row double-click opened Edit;
- the bottom Exit action worked;
- Windows Settings reached the guided uninstall and presented wipe choices.

The operator intentionally did not execute manual, scheduled, nonzero-exit, or
process-start-failure task cases. The headless suite still exercised scheduler,
store, API/client-server, executor, Windows process creation, and integration
behavior, including successful and failed execution classes. That is meaningful
source-level proof, but it does not replace #94's installed LocalSystem and
ordinary-user exact-candidate matrix.

### Post-wipe system audit

A read-only audit captured at `2026-09-03T06:22:41-04:00` found all declared
owned and installer state absent:

| State | Result |
| --- | --- |
| `C:\Program Files\go-schedule` | Absent |
| `goschedd` service | Absent |
| MSI product registration | Absent |
| Public Desktop/current-user/all-users Start Menu shortcuts | Absent |
| `C:\ProgramData\goschedule` | Absent |
| Fyne app root in all five registered Windows profiles | Absent |
| cleanup-result ledger and cleanup-status registry key | Absent |

The registered-profile audit covered SYSTEM, LocalService, NetworkService, the
maintainer profile, and CodexSandboxOffline. The implementation retains a ledger
and registry status for incomplete cleanup, so absence of both owned roots and
failure evidence is consistent with a complete wipe. The unrelated-sentinel and
preserve/cancel matrices were not performed in this attended walkthrough and
remain open under #98/#94.

### Issue-only disposition of walkthrough feedback

No UI product code was added to S043. New bounded issues were created for:

- #109 systemic hover/active/focus contrast;
- #110 natural command entry and a larger arguments editor;
- #111 Windows wheel sensitivity and its Options setting;
- #112 a labeled, frozen-header Tasks table; and
- #113 tabulated, color-coded Schedule and Activity rows.

Existing issues were updated instead of duplicated: #101 and #106 now carry the
System-font A/B result and font-family/default work; #106 also owns compact
storage rows and selector behavior; #104 owns the sidebar boundary; and #105
owns the Exit-glyph refinement. Confirmation evidence was added to #103, #97,
#98, and #94.

Accessibility acceptance in #109, #112, and #113 uses measurable shared-state
criteria: at least 4.5:1 contrast for normal text, 3:1 for large text and
essential non-text indicators, visible hover and keyboard focus, and no color-
only meaning. This prevents another one-control color patch from masking the
systemic defect.

## Final disposition

S043 achieved its pre-publication purpose: one identified demo was built,
automated checks passed, attended findings were recorded before any correction,
the wipe result was independently inspected, and unfinished exact-candidate
proof remains open. The branch may proceed through pull-request CI and review,
but no tag, release, merge, or issue closure is authorized by this record.

After recording the attended results and issue disposition, the complete eight-
gate suite passed again: format, vet, lint (zero issues), race, GUI, coverage,
documentation, and automation. The coverage results remained engine 86.4%,
schedule 89.2%, timezone 91.3%, store 84.1%, catchup 88.9%, and logbus 91.1%.
After the first external review identified four pre-existing EOF blank lines,
commit `9baed96` removed them and `git diff --check origin/main...HEAD` passed
against the complete committed pull-request range.
