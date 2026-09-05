# Data Model: v1.0.0 Release Operations

## DispositionPacket

| Field | Type | Rules |
| --- | --- | --- |
| `schema_version` | integer | Exactly `1` |
| `candidate` | CandidateIdentity | Copied from validated formal evidence |
| `issues` | ordered IssueFile list | Exactly ten entries in ascending issue order |

The packet is created only after complete production validation. It contains no generation timestamp, local path, credential, or operator name so identical inputs remain byte-for-byte deterministic.

## CandidateIdentity

| Field | Rules |
| --- | --- |
| repository | Exact expected `owner/repository` |
| tag | Exact expected semantic release tag |
| commit | Exact expected 40-character commit |
| workflow | Canonical `Release` workflow |
| run ID / attempt | Positive staged-run identity |
| filename / bytes / SHA-256 | Exact MSI file identity |
| product version / ProductCode | Exact compiled installer identity |

The evidence, independent manifest, physical MSI, and command-line expectations must all agree before this entity can be rendered.

## IssueFile

| Field | Type | Rules |
| --- | --- | --- |
| `issue` | integer | One of 96, 98, 101, 104, 105, 106, 109, 111, 112, 113 |
| `file` | string | Canonical `issue-NNN.md` filename |
| `observations` | ordered string list | Exact compiled mapping for that issue |

## IssueDisposition

| Component | Rules |
| --- | --- |
| Header | Identifies formal v1.0.0 evidence and issue number |
| Candidate table | Complete immutable identity and exact workflow link |
| Validator statement | Explicitly states complete production validation passed |
| Evidence archive | Tag-specific immutable release-asset URL and archive member semantics |
| Environments | Only environments referenced by mapped observations, sorted by ID |
| Observations | Exact canonical mapping order; no missing or extra scenario |
| Metrics | Deterministic indented JSON with sorted object keys |
| Attachments | Exact validated archive-relative paths for each observation |
| Disposition boundary | States that the record supports review and does not itself close the issue |

## CanonicalIssueMapping

| Issue | Observations |
| --- | --- |
| 98 | Nine `setup.*` and seven `remove.*` lifecycle observations |
| 101 | Two standard/scaled appearance observations |
| 104 | Two navigation/options and two interaction-state observations |
| 105 | Two navigation/options and two interaction-state observations |
| 106 | Two appearance, two navigation/options, and scroll-input observations |
| 109 | Two interaction-state and four structured-table observations |
| 111 | Scroll-input observation |
| 112 | Two Tasks-table and two interaction-state observations |
| 113 | Two Schedule/Activity-table and two interaction-state observations |
| 96 | All 36 pre-desktop lifecycle observations and coordinator references #97, #98, #94, #89, and #90 |

The exact identifiers are normative in `spec.md` and the packet contract.

## Packet State Transitions

```text
inputs selected
  -> decoded
  -> candidate validated
  -> evidence validated
  -> archive contents validated
  -> manifest identity validated
  -> all files rendered in memory
  -> private sibling directory written
  -> destination atomically committed
```

Any failure before the final transition leaves the destination absent. A destination that already exists is a terminal conflict, not a resumable state.

## Release State Transitions

```text
reviewed S049 merge
  -> separately authorized immutable tag
  -> successful draft staging with eight payloads
  -> exact candidate verification
  -> 47 attended passing observations
  -> evidence archive upload (nine pre-checksum payloads)
  -> ten independent issue dispositions
  -> authorized promotion and checksum generation
  -> public latest release with ten payloads
  -> #122 and milestone closure after final audit
```

No packet state transition performs a release state transition.
