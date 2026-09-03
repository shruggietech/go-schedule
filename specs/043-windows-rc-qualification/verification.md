# Verification: Windows Demo Qualification

**Date**: 2026-09-03

**Branch**: `codex/043-windows-rc-qualification`

## Current result

In progress. The local demo candidate is not ready for handoff until the Spec Kit
analyze gate, automated checks, compiled inspection, identity audit, integrity
scan, local commit, and final rebuild disposition are complete.

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

All 13 changed or added text files decode as strict UTF-8, none has a UTF-8 BOM,
and no mojibake signature was found. `git diff --check` passed. The template and
clarification scan found no unresolved placeholder.

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

Final demo identity is pending a clean rebuild from the correction commit.

## Native evidence boundary

No unattended or headless result will be represented as proof of visual quality,
Windows Settings interaction, native window geometry, or user judgment. The demo
remains unofficial and the pull request remains withheld until the operator
finishes attended testing.
