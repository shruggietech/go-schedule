# Research: v1.0.0 Release Operations

## R1. Final release commit after an additional reviewed pull request

**Decision**: Make the reviewed S049 merge commit the exact `v1.0.0` tag target
and update the S048 publication contract accordingly.

**Rationale**: S048 requires `HEAD`, `origin/main`, and its own merge commit to
match before tagging. Once the operator-required S049 PR merges, that condition
cannot be true. Tagging an older commit after `main` advances would contradict
the synchronized-main safety rule. S049 changes only release tooling, tests, and
documentation, so advancing the reviewed boundary does not alter packaged
application behavior.

**Alternatives considered**:

- Tag S048 before S049 merges. Rejected because no tag authorization exists and
  the new disposition generator would not yet be on synchronized `main`.
- Tag the historical S048 commit after S049 merges. Rejected because it violates
  the established main-equality preflight and obscures why current main is not
  the release source.
- Skip S049's PR. Rejected because the operator explicitly required an official
  pull request and third-party review.

## R2. Existing command versus a new orchestrator

**Decision**: Add `render-dispositions` to `windows-release-gate`.

**Rationale**: The command already owns strict decoding, candidate validation,
archive extraction, production evidence validation, and conventional exit
codes. Extending it guarantees the report uses the same trust boundary as
promotion and avoids a second implementation.

**Alternatives considered**:

- A PowerShell release script. Rejected because it would duplicate validation
  and create an avoidable Windows child-process-launch surface.
- A GitHub workflow that comments and closes issues. Rejected because it couples
  evidence rendering to credentials and could turn a mapping defect into bulk
  remote mutation.
- Manual templates only. Rejected because they do not mechanically prevent
  missing, mismatched, or copied evidence.

## R3. Validation order

**Decision**: Decode and validate every input, combine all production failures,
and render only when the complete failure set is empty.

**Rationale**: Rendering first could leave credible-looking partial output. The
existing validator intentionally reports all independently discoverable defects
in one run, which is more actionable than failing on only the first scenario.

**Alternatives considered**:

- Validate the evidence JSON but not archive inventory or manifest. Rejected
  because extra/missing files and independently staged identity are release
  requirements.
- Render files incrementally as observations pass. Rejected because partial
  success is ambiguous and violates the fail-closed contract.

## R4. Packet format

**Decision**: Write ten deterministic Markdown records and one deterministic
JSON index.

**Rationale**: Markdown can be reviewed and applied directly to GitHub issues.
The JSON index makes candidate identity, issue/file inventory, and future audit
mechanically checkable without parsing prose.

**Alternatives considered**:

- One combined Markdown document. Rejected because manual extraction reopens
  the copy/omission risk.
- JSON only. Rejected because it creates another rendering step before issue
  comments.
- Store copied attachments beside reports. Rejected because the canonical ZIP
  already integrity-protects them and copies could drift.

## R5. Atomic and overwrite-safe output

**Decision**: Require an absent target, build the complete packet in a private
sibling directory, and rename it into place only after all writes succeed.

**Rationale**: The user must be able to distinguish a complete packet from a
failed attempt. Refusing all existing targets also prevents silent replacement
of evidence already reviewed or applied.

**Alternatives considered**:

- Write directly to the destination. Rejected because interruption can leave a
  partial directory.
- Delete and recreate the target. Rejected because evidence replacement must
  never be implicit.
- Merge missing files into an existing target. Rejected because mixed candidate
  identities would be difficult to detect.

## R6. Canonical issue mapping

**Decision**: Compile the S047 mapping into ordered code and expose a copied view
for tests.

**Rationale**: The mapping is part of the v1.0.0 release contract, not arbitrary
runtime configuration. A fixed map provides deterministic ordering and makes a
reviewed source change necessary before any issue gains or loses evidence.

**Alternatives considered**:

- Read mappings from the evidence archive. Rejected because the evidence
  producer could then decide which issue its observations satisfy.
- Read mappings from a mutable external file. Rejected because promotion would
  not know which version was authoritative.

## R7. GitHub-safe rendering

**Decision**: Escape table-delimiter, HTML, code-delimiter, newline, and mention
characters in user-controlled text; render structured metrics as indented JSON;
link only to tag- and run-specific GitHub URLs derived from validated identity.

**Rationale**: Observation summaries and environment labels are evidence input.
They must not create headings, tables, mentions, or links that masquerade as
generator-authored disposition content.

**Alternatives considered**:

- Trust GitHub's sanitizer. Rejected because sanitized HTML does not prevent
  Markdown restructuring or unintended mentions.
- Remove all punctuation. Rejected because it would damage useful diagnostic
  evidence.

## R8. Remote mutation boundary

**Decision**: The generator is offline and never comments, closes, tags,
uploads, dispatches, or promotes.

**Rationale**: A human or authorized agent must still compare each issue's full
acceptance criteria with its record. Keeping generation pure makes it safe to
rerun and test with fixtures without credentials.

**Alternatives considered**:

- Optional `--apply` mode. Rejected because one typo would expand a read-only
  evidence operation into ten issue mutations.
- Auto-close on all-pass. Rejected because mapped observations do not replace
  acceptance-criterion review.

## R9. Compatibility and dependency posture

**Decision**: Add no dependency and leave the evidence schema, existing
subcommands, workflows, and application binaries unchanged.

**Rationale**: Standard-library Markdown/JSON rendering and filesystem staging
are sufficient. The release tooling belongs outside product runtime and must
not perturb the exact behavior already accepted in S047.

**Alternatives considered**:

- A Markdown templating dependency. Rejected because fixed output is simpler and
  safer to render directly.
- Evidence schema issue fields. Rejected because that would invalidate existing
  collectors and make evidence self-authorize its disposition.

## R10. Documentation surface

**Decision**: Extend the Windows evidence runbook, the S048 publication
contract, S049 quickstart, and changelog. Do not create a blog post.

**Rationale**: This is maintainer-only release tooling, and the repository has
no product-blog tree. The selected documents are already the authoritative
operator surfaces for issue #122.

**Alternatives considered**:

- Add a new public documentation page. Rejected as duplicate release-procedure
  content.
- Add a blog scaffold solely to satisfy generic feature guidance. Rejected as
  unrelated product scope and inconsistent with project practice.
