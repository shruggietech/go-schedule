# Data Model: Post-release dependency refresh

## Dependency Selection

| Field | Meaning |
| --- | --- |
| Ecosystem | Root Go, desktop Go, or frontend npm graph |
| Source | Dependabot pull request proposing the update |
| Direct version | Explicit manifest selection |
| Resolved companions | Transitive versions selected by native tooling |
| Runtime floor | Minimum Go or Node version required by the graph |
| Evidence | Restore, test, build, security, and graph-cleanliness results |

## Source Disposition

| Field | Meaning |
| --- | --- |
| Pull request | #240, #241, or #242 |
| Representation | Manifest and integrity-file locations covering the proposal |
| Replacement | Official S090 pull request URL |
| State | Open before replacement publication; superseded afterward |
| Comment | Explicit explanation and replacement link |

## Review Finding

| Field | Meaning |
| --- | --- |
| Round | Automatic first round or authorized second round |
| Author | Codex, security bot, CI, or human reviewer |
| Finding | Actionable request or approval signal |
| Disposition | Fixed, explained, or not applicable with rationale |
| Evidence | Commit and final-head check result |

## State Transitions

1. Proposed dependency selections become resolved only after native graph regeneration.
2. Resolved selections become verified only after clean restoration, focused evidence, and canonical gates pass.
3. Source pull requests become superseded only after the replacement PR is published and linked.
4. S090 becomes merge-ready only after every finding is dispositioned and latest-head CI is green.

