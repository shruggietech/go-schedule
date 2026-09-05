# Feature Specification: v1.0.0 Release Execution and Audit

**Feature Branch**: `codex/050-v100-release-execution`

**Created**: 2026-09-04

**Status**: In Progress

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Public release [v1.0.0](https://github.com/shruggietech/go-schedule/releases/tag/v1.0.0), staging run [33838072246](https://github.com/shruggietech/go-schedule/actions/runs/33838072246), exception audit [#122 comment](https://github.com/shruggietech/go-schedule/issues/122#issuecomment-5536182541), verified public checksums, and local canonical verification completed 2026-09-04; formal qualification and issue reconciliation remain actionable under #122

**Input**: Complete issue
[#122](https://github.com/shruggietech/go-schedule/issues/122) by qualifying
and publishing the exact `v1.0.0` candidate staged from the reviewed S049 merge
commit, reconciling every remaining v1 readiness issue from its own formal
evidence, and preserving an auditable post-release record. No runtime or
candidate bytes may change in this slice.

## Problem Statement

The reviewed S048 and S049 slices prepared the v1.0.0 source boundary, staging
workflow, 47-observation Windows qualification gate, deterministic issue
disposition packet, and fail-closed promotion workflow. The remaining work is
operational and irreversible: accept one exact tag-triggered draft, qualify its
MSI in a normal Windows desktop session, reconcile ten readiness issues, publish
the already-staged artifacts without rebuilding them, and prove the final public
state is internally consistent.

The annotated `v1.0.0` tag was explicitly authorized and created at
`ff47b4410d1aecbfadb8165d1ebf025ca1dadde7` under the approved S049 post-merge
contract before this audit branch was created. S050 records and verifies the
result. It must never move that tag, substitute candidate bytes, or allow its
review branch to affect the release identity.

## Execution Deviation

After the exact staged MSI was installed and accepted, the maintainer explicitly
directed immediate publication. That instruction superseded the plan's
conditional-publication sequence, but it did not turn missing observations into
evidence. S050 therefore deviates from the planned promotion path as follows:

- no 47-observation archive was manufactured or copied from prior local-demo
  evidence;
- no evidence disposition packet was generated and no evidence-dependent issue
  or milestone was closed;
- the green draft was published directly after tag, manifest, MSI, installed
  version, service identity, ordinary-user health, and payload hashes passed;
- `SHA256SUMS.txt` was generated for the eight staged payloads and verified again
  from a fresh public download;
- issue #122 remains open with the exact exception and unfinished criteria.

This is an explicit release-governance exception, not a claim that FR-005
through FR-013 or SC-002 through SC-005 passed as originally written. The audit
branch records the decision and preserves the unresolved work visibly. S050
therefore remains In Progress after publication until that work is completed or
separately superseded by an approved specification.

## Scope

### In scope

- Validate the one tag-triggered Release workflow and draft release against the
  reviewed tag, commit, event, run, manifest, and exact staged asset inventory.
- Download and independently verify the exact Windows MSI and candidate
  manifest before installation.
- Complete all 47 canonical attended Windows observations against that MSI,
  including the required attachments and environment facts.
- Finalize and independently verify one immutable formal evidence archive.
- Generate, review, and apply the ten deterministic issue dispositions.
- Close only readiness issues whose individual acceptance criteria and mapped
  formal observations pass, then reconcile coordinators from actual child state.
- Promote the existing draft without rebuilding or replacing artifacts.
- Audit public assets, checksums, release metadata, documentation identity,
  issue state, milestone state, and latest-release identity.
- Record exact immutable evidence and any stop/failure/recovery in the S050
  verification artifact reviewed through the normal pull-request workflow.

### Out of scope

- Runtime, GUI, installer, daemon, CLI, schema, or workflow behavior changes.
- Moving, recreating, or force-updating `v1.0.0`.
- Treating S043 or S047 local-demo evidence as formal candidate evidence.
- Closing Post-v1 issues or representing deferred work as part of v1.0.0.
- Rebuilding, repackaging, resigning, or otherwise substituting staged assets
  during qualification or promotion.

## Clarifications

### Session 2026-09-04

- Q: Can S050 use the ordinary implementation-PR merge as the release tag
  boundary? -> A: No. S049 intentionally fixed the reviewed release boundary;
  S050 is a post-tag operational audit branch and cannot affect candidate
  identity.
- Q: Does the authorization to "do as you wish" and "see it done" cover tag
  staging and publication? -> A: Yes, in the direct context of the requested
  S050 release operation, but publication remains conditional on every
  candidate, evidence, issue, and checksum gate passing. This initial answer
  was superseded by the maintainer's later direct order to publish after exact-
  candidate installation, which is recorded in Execution Deviation.
- Q: Can an issue be closed from a passing mapped observation alone? -> A: No.
  The complete issue acceptance criteria must also be reviewed and satisfied.
- Q: What happens when the maintainer directs publication after exact-candidate
  install acceptance but before formal evidence exists? -> A: Publish only as
  an explicit, auditable exception; do not fabricate evidence, run disposition
  closure, or close the evidence-dependent issues and milestone.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Qualify One Immutable Candidate (Priority: P1)

As a release steward, I can prove that every formal Windows observation was
performed against the exact MSI staged by the one accepted tag workflow, so
local demonstrations or mismatched artifacts cannot be mistaken for release
evidence.

**Why this priority**: Candidate identity and complete attended evidence are the
primary release safety boundary. Nothing may be reconciled or published without
them.

**Independent Test**: Starting from the draft release, independently verify the
candidate manifest and MSI, complete the generated 47-observation workspace,
finalize it once, and validate the resulting archive against the same manifest
and MSI.

**Acceptance Scenarios**:

1. **Given** the authorized tag at the reviewed S049 commit, **when** staging
   completes, **then** exactly one successful push run produces one draft with
   the contractually exact assets and matching candidate identity.
2. **Given** the staged MSI and manifest, **when** formal qualification is
   completed, **then** all 47 canonical observations pass in their required
   environment classes with every required supported attachment.
3. **Given** the finalized archive, **when** an independent validation pass is
   run, **then** the archive, manifest, MSI, repository, tag, commit, workflow
   run, and product identity all match without warnings or exceptions.

---

### User Story 2 - Reconcile Readiness Issues Individually (Priority: P1)

As a maintainer, I can judge and close each v1 readiness issue from its own
acceptance criteria and mapped formal evidence, so release association never
stands in for proof.

**Why this priority**: Promotion requires every readiness issue to be genuinely
satisfied while preserving independent issue and coordinator boundaries.

**Independent Test**: Generate the disposition packet from the independently
verified evidence and compare every record with its issue acceptance criteria.
Only passing issues receive their record and closure; coordinator state is then
derived from the actual child states.

**Acceptance Scenarios**:

1. **Given** the valid formal archive, manifest, and MSI, **when** dispositions
   are rendered, **then** exactly ten deterministic issue records plus one
   packet index are produced.
2. **Given** a child record with passing mapped observations, **when** an
   unmapped acceptance criterion remains unproven, **then** that issue and all
   dependent coordinators remain open and promotion stops.
3. **Given** every child issue is individually satisfied, **when** coordinator
   issues are reconciled, **then** their checklists and state reflect actual
   child and prerequisite state rather than assumed completion.

---

### User Story 3 - Publish and Audit v1.0.0 (Priority: P2)

As a user evaluating go-schedule, I can obtain one public v1.0.0 release whose
tag, packages, checksums, notes, documentation, and project records all identify
the same qualified software.

**Why this priority**: Publication is the user-visible outcome, but it is valid
only after immutable candidate qualification and issue reconciliation.

**Independent Test**: Dispatch promotion for the qualified tag, then independently
inventory and hash every public payload and compare the public release, latest
pointer, tag, documentation, issues, and milestone with the formal candidate.

**Acceptance Scenarios**:

1. **Given** nine exact pre-checksum payloads and satisfied readiness issues,
   **when** promotion runs, **then** it adds one checksum manifest and publishes
   the existing draft without rebuilding any payload.
2. **Given** the public release, **when** the final audit runs, **then** all ten
   payloads are non-empty, every recorded digest passes, and all release identity
   surfaces agree on v1.0.0 and the tagged commit.
3. **Given** a successful final audit, **when** issue #122 and the milestone are
   reconciled, **then** they close with durable links to the exact release,
   workflow, evidence archive, disposition records, and audit result.

### Edge Cases

- The tag-triggered workflow fails, is canceled, is rerun, or produces more than
  one plausible staging run.
- A draft already exists but contains missing, extra, zero-byte, or replaced
  assets.
- The tag object resolves correctly but the peeled commit, manifest commit, or
  workflow head differs.
- The formal workspace is interrupted, partially populated, or finalized twice.
- A screenshot is missing, corrupt, not a supported raster format, or does not
  show the required state at the required display class.
- One observation passes while its issue has another unproven acceptance
  criterion.
- GitHub becomes unavailable after evidence finalization or during issue
  reconciliation; local evidence remains valid but no remote completion is
  claimed.
- Promotion fails before publication or after checksum upload; the release stays
  draft until the same immutable assets can pass the documented recovery path.
- The latest-release pointer or documentation cache lags after publication; the
  audit distinguishes temporary propagation from contradictory durable state.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: S050 MUST preserve `v1.0.0` as one immutable annotated tag whose
  peeled commit is `ff47b4410d1aecbfadb8165d1ebf025ca1dadde7`.
- **FR-002**: The accepted staging run MUST be triggered by the tag push, target
  `v1.0.0`, use the tagged commit, complete successfully, and leave the release
  in draft state.
- **FR-003**: The draft MUST contain exactly the seven packages and one candidate
  manifest defined by the reviewed release contract before formal evidence is
  uploaded.
- **FR-004**: Candidate validation MUST bind the manifest and exact MSI bytes to
  the repository, tag, commit, workflow run, workflow attempt, package name,
  byte length, SHA-256 digest, ProductVersion, and ProductCode.
- **FR-005**: Formal qualification MUST contain exactly 47 passing canonical
  observations collected from the exact candidate in the required Windows
  environment classes.
- **FR-006**: Every required attachment MUST be present, supported, non-empty,
  archive-contained, and associated with the declared observation.
- **FR-007**: The evidence archive MUST be finalized once and independently
  validated against the independently downloaded candidate manifest and MSI.
- **FR-008**: The immutable evidence archive MUST be uploaded under the exact
  contract filename without overwriting unexplained state, producing exactly
  nine pre-checksum release payloads.
- **FR-009**: The disposition operation MUST produce the exact ten issue records
  and packet index from the validated formal inputs without remote mutation.
- **FR-010**: Each readiness issue MUST be reviewed against both its mapped
  formal observations and its complete individual acceptance criteria before a
  comment or closure is applied.
- **FR-011**: Issue #96 MUST be reconciled only after actual child and prerequisite
  states are known, while preserving #98's independent closure boundary.
- **FR-012**: Issue #122 MUST remain open until promotion and the final public
  audit both pass.
- **FR-013**: Promotion MUST accept only the exact qualified tag and nine
  pre-checksum assets, generate one checksum entry for every payload, verify the
  downloaded final set, and publish the existing draft without rebuilding.
- **FR-014**: The public release MUST contain exactly ten non-empty assets: the
  nine qualified payloads plus `SHA256SUMS.txt`.
- **FR-015**: The final audit MUST verify release state, latest-release identity,
  tag and commit, all asset names and sizes, every checksum, release notes,
  README and changelog identity, readiness issue state, #122, and milestone state.
- **FR-016**: Any failed or ambiguous gate MUST stop downstream mutation and be
  recorded accurately without moving the tag or substituting candidate bytes.
- **FR-017**: The S050 pull request MUST contain only specifications and the
  immutable operational audit record; it MUST NOT alter shipped v1.0.0 bytes or
  claim release success before the external state is verified.

### Key Entities

- **Release identity**: Repository, annotated tag, peeled commit, workflow run
  and attempt, and release database identity that define one candidate.
- **Candidate manifest**: The staged record binding workflow provenance to each
  expected package and the Windows MSI product identity.
- **Formal evidence archive**: The immutable 47-observation Windows record with
  environment facts, metrics, and attachments.
- **Disposition packet**: Ten issue-specific evidence records and one index
  generated from validated candidate inputs.
- **Published asset set**: Nine qualified payloads plus the final checksum file.
- **Release audit record**: Durable S050 evidence linking every operational
  mutation and verification result without containing mutable secrets.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Exactly one immutable v1.0.0 tag, one accepted staging run, and one
  public release resolve to the same reviewed 40-character commit.
- **SC-002**: Formal evidence contains 47 of 47 canonical observations in passing
  state and independently validates with zero missing or extra observations,
  attachments, or candidate-identity mismatches.
- **SC-003**: All ten readiness dispositions are generated, individually
  reviewed, and applied with zero premature closures.
- **SC-004**: The public release exposes exactly ten non-empty assets and every
  one of the nine payload checksums verifies byte-for-byte.
- **SC-005**: The latest-release pointer, release page, tag, candidate manifest,
  packages, README, changelog, issue #122, and v1.0.0 milestone contain no
  contradictory version or completion state at final audit.
- **SC-006**: S050 changes zero candidate runtime files and the post-release PR
  passes all canonical repository and hosted review gates.

## Assumptions

- The current Windows desktop session can run the exact staged MSI and capture
  the required native evidence without delegating observations to local-demo
  fixtures.
- GitHub Actions and Releases remain available; temporary outages pause the
  operation without changing candidate identity.
- The reviewed S047 evidence schema, S048 publication contract, and S049
  disposition mapping remain authoritative for v1.0.0.
- The user's instruction to see the release completed is explicit authorization
  for conditional publication after all documented gates pass. The later
  instruction to publish immediately superseded that condition and is handled
  as the explicit exception above.
