# Feature Specification: Windows Release Qualification

**Feature Branch**: `codex/047-windows-release-qualification`

**Created**: 2026-09-03

**Status**: Implemented

**Delivery**: Review branch `codex/047-windows-release-qualification`; committed
implementation, local-demo acceptance, and canonical verification completed;
publication intentionally withheld

**Input**: Complete the Windows release-readiness slice for issues #96, #98,
#101, #104, #105, #106, #109, #111, #112, and #113 as far as possible under
autopilot, then halt before the review branch is pushed.

## Scope and Evidence Boundary

S047 consolidates the remaining v1 Windows acceptance work into one traceable
qualification path. It extends the existing release-promotion evidence gate so
native desktop appearance, interaction, navigation, Options, scrolling, and
structured-table outcomes cannot be omitted from a formal candidate. It also
produces a clearly non-release S047 demo MSI for pre-push attended testing.

A branch-local build cannot be the formal release candidate. The repository
creates that immutable artifact only after reviewed source is merged, an
authorized semantic-version tag is staged by the Release workflow, and the
attended evidence is collected from the exact staged MSI. Local demo evidence
can find defects before publication but cannot be relabeled, copied, or counted
as formal release evidence.

### Scope in

- Add fixed, issue-traceable desktop qualification observations to the existing
  fail-closed Windows release evidence contract.
- Require native image evidence, both supported palettes, 100 percent and
  scaled-DPI coverage, supported window sizes, and explicit interaction/input
  outcomes where applicable.
- Preserve the existing installer, uninstall, startup, error, access, and real
  task-execution matrix, including #98 preserve/wipe verification.
- Update the resumable attended collector, passing synthetic fixture, validator
  tests, Windows runbook, and promotion documentation together.
- Build and inspect one hash-bound, clearly labeled S047 local demo MSI from a
  committed source identity after all unattended gates pass.
- Record pre-push attended results without closing issues whose formal staged-
  candidate requirements remain unsatisfied.
- Fix only reproducible defects discovered during S047 qualification, with a
  failing regression test before implementation.

### Scope out

- Pushing the review branch, opening a pull request, merging, tagging, staging,
  promoting, or publishing a release without the separately required authority.
- Claiming that headless tests, synthetic fixtures, or the local demo prove
  native rendering, physical input feel, Windows Settings behavior, or the
  identity of a later workflow-staged candidate.
- Implementing Post-v1 diagnostics issue #102.
- Reworking S044, S045, or S046 product behavior without a concrete S047
  reproduction.
- Reopening or rewriting the historical state of closed issue #94. Its S040
  release-gate contract remains the mechanism extended by this slice.

## Clarifications

### Session 2026-09-03

- Q: Can the pre-push S047 MSI satisfy the formal exact-candidate requirement?
  -> A: No. It is a `local-demo` artifact used only for exploratory attended
  qualification; the exact staged candidate requires reviewed, merged, tagged
  source and must repeat the formal matrix.
- Q: Should all remaining v1 visual issues share one formal evidence bundle?
  -> A: Yes. They depend on the same Windows environment, themes, DPI settings,
  window sizes, populated data, and interaction states, so one issue-traceable
  bundle is the smallest coherent and auditable unit.
- Q: How is optional precision-touchpad hardware represented? -> A: The scroll
  observation must always cover a conventional wheel. It records either passing
  touchpad evidence or an explicit hardware-unavailable reason without
  misrepresenting that optional check as executed.
- Q: When may an open issue close? -> A: Only after its individual criteria are
  mapped to passing exact-candidate observations and reviewed; shared bundle
  success does not implicitly close every referenced issue.
- Q: What remains before the branch may be pushed? -> A: The local demo must be
  built and inspected, all unattended gates must pass, and the maintainer must
  complete or explicitly disposition the pre-push attended checklist. Any found
  defect must be resolved and requalified first.

## User Scenarios & Testing

### User Story 1 - Enforce Complete Desktop Evidence (Priority: P1)

As the release maintainer, I cannot promote a Windows candidate whose evidence
omits the remaining v1 desktop acceptance surfaces.

**Why this priority**: The current promotion gate enforces installer and core
runtime evidence, but the later S044 and S045 native requirements exist only as
runbook prose and can therefore be skipped without a machine-detectable failure.

**Independent Test**: Validate a complete synthetic bundle successfully, then
remove or corrupt each new desktop observation and confirm validation fails with
the exact scenario and metric named.

**Acceptance Scenarios**:

1. **Given** a candidate bundle without one required desktop observation,
   **when** promotion validation runs, **then** validation fails and identifies
   the missing observation.
2. **Given** a desktop observation without the required palette, DPI, size,
   state, population, or input evidence, **when** validation runs, **then** it
   fails with an actionable metric diagnostic.
3. **Given** every old and new observation passes against one exact candidate,
   **when** validation runs, **then** the bundle remains eligible for the later
   remote-origin and promotion checks.

---

### User Story 2 - Record Native Results Safely (Priority: P1)

As the Windows operator, I receive a resumable, bounded workspace with explicit
templates for every new visual and interaction observation so I can record
results without editing the canonical evidence file by hand.

**Why this priority**: Native judgment cannot be automated credibly, but its
inputs and completeness can be constrained, hashed, and reviewed.

**Independent Test**: Initialize the collector against an inert candidate,
inspect the generated placeholders/templates, import reviewed fragments, and
confirm finalization rejects missing, duplicate, malformed, or attachment-free
desktop observations.

**Acceptance Scenarios**:

1. **Given** a new evidence workspace, **when** it is initialized, **then** every
   existing and S047 desktop observation exists once as an unavailable
   placeholder with a scenario-specific metrics template.
2. **Given** native screenshots and reviewed metrics, **when** an observation is
   recorded, **then** the collector binds it to one environment and refuses to
   overwrite or duplicate it.
3. **Given** a physical touchpad is unavailable, **when** scroll evidence is
   recorded, **then** the conventional-wheel result remains mandatory and the
   absent optional hardware has a non-empty reason.

---

### User Story 3 - Test One Pre-Push Demo (Priority: P1)

As the maintainer, I receive one clearly identified local installer only after
all safe unattended checks pass, so the last pre-push work is a focused native
walkthrough rather than repository diagnosis.

**Why this priority**: A local demo catches Windows rendering and interaction
defects before review while preserving the formal candidate identity boundary.

**Independent Test**: Compare the delivered file with its recorded source
commit, embedded demo version, ProductVersion, ProductCode, byte size, SHA-256,
and compiled inspection report, then complete the condensed attended checklist.

**Acceptance Scenarios**:

1. **Given** the committed S047 source and passing unattended gates, **when** the
   local build completes, **then** its filename, embedded version, inspection,
   and verification record consistently label it as an S047 local demo.
2. **Given** a failing automated check or compiled-MSI inspection, **when** the
   handoff is prepared, **then** the MSI is not reported ready for testing.
3. **Given** attended feedback, **when** any observation fails, **then** the
   failure is recorded before a correction and the corrected build receives a
   new identity and repeats affected checks.

---

### User Story 4 - Close Work by Evidence, Not Association (Priority: P2)

As the project maintainer, I can determine exactly which open v1 issues are
complete from the final evidence, while incomplete criteria remain visible.

**Why this priority**: A bundled slice is efficient only if issue-level
traceability prevents accidental collective closure.

**Independent Test**: Review the traceability matrix and verify every issue has
named observations, acceptance criteria, evidence class, and an explicit close
or remain-open disposition.

**Acceptance Scenarios**:

1. **Given** a passing local demo observation but no formal candidate result,
   **when** issue disposition is reviewed, **then** exact-candidate issues remain
   open.
2. **Given** the later exact candidate passes all criteria for one issue,
   **when** records are reconciled, **then** that issue can close independently
   even if another bundled issue remains open.
3. **Given** coordinator #96, **when** child states are reconciled, **then** its
   index reflects actual GitHub state and does not treat a reference as proof.

### Edge Cases

- A rebuilt local MSI has a different byte identity even when the source tree
  appears unchanged; prior attended observations remain bound to the old hash.
- A synthetic passing fixture must never satisfy the production evidence class.
- Screenshots can be supplied without matching metrics, or metrics without a
  screenshot; both forms are incomplete and block finalization.
- A text or vector attachment can declare an image media type; the gate must
  inspect bytes and require a supported raster format.
- Dark and Light observations can accidentally come from one unchanged palette;
  required palette sets must identify both values explicitly.
- A scaled-DPI observation can repeat the 96-DPI environment; the gate must
  require one effective DPI greater than 96.
- A table can be populated below the required 100 rows or tested at only one
  window size; both are incomplete.
- Precision touchpad hardware may be absent, but mouse-wheel evidence may not be
  skipped and optional absence must be explained.
- Windows elevation or an undisposable host may make destructive lifecycle
  automation unavailable; it remains a blocking formal-candidate result rather
  than being silently omitted.
- Issue #94 is closed while coordinator #96 still lists it unchecked. S047
  reports the inconsistency without inventing completion evidence or rewriting
  history.

## Requirements

### Functional Requirements

- **FR-001**: The promotion evidence contract MUST require every pre-existing
  release scenario without weakening its validation.
- **FR-002**: The contract MUST add fixed standard- and scaled-DPI desktop
  observations for appearance, interaction states, navigation and Options, the
  Tasks table, and the Schedule/Activity tables, plus standard-DPI scroll input.
- **FR-003**: Every desktop observation MUST use an intended user at medium
  integrity against the installed service and MUST reference at least one
  integrity-protected native raster screenshot validated from file bytes.
- **FR-004**: Appearance evidence MUST cover Dark and Light modes, System-font
  default/reset/persistence, sharp Info and ordinary body text, unclipped
  centered labels, resize/minimize/restore/reopen behavior, 100 percent scaling,
  and one effective DPI greater than 96.
- **FR-005**: Interaction evidence MUST cover navigation, selector, ordinary,
  primary, danger, and dialog controls across rest, hover, keyboard focus,
  pressed, selected, and disabled states; it MUST record a 4.5:1 normal-text
  floor and 3:1 essential non-text floor without color-only meaning.
- **FR-006**: Navigation and Options evidence MUST cover destination ordering,
  balanced unclipped rail spacing, the full-height boundary, bottom-right
  non-selected Exit treatment, compact storage rows, unavailable-row muting,
  exact Copy behavior, selector alternatives, no horizontal scrollbar, and both
  1280-by-800 and 800-by-600 content sizes.
- **FR-007**: Scroll evidence MUST cover conventional-wheel behavior at 1x, 2x,
  and 4x across every application-owned vertical scroll surface, immediate
  application, persistence, absence of nested multiplication, and preserved
  keyboard behavior.
- **FR-008**: Scroll evidence MUST represent precision-touchpad behavior as
  either executed and passing or unavailable with a non-empty hardware reason.
- **FR-009**: Tasks-table evidence MUST use at least 100 rows and cover fixed
  Task, Enabled, Lifecycle, Time zone, and Group headers; distinct status
  semantics; no unexplained brackets; full-value disclosure; odd/even hover,
  focus, selection, live-refresh identity, toolbar actions, and double-click at
  both supported sizes and palettes.
- **FR-010**: Schedule/Activity evidence MUST use at least 100 applicable rows
  per view and cover fixed headers, normalized INFO/WARNING/ERROR casing,
  restrained matching glyph/text semantics, full detail and value disclosure,
  all documented event states, overlapping row states, live-refresh identity,
  filters, clearing, acknowledgement, range/calendar switching, both palettes,
  and both supported sizes.
- **FR-011**: The collector MUST initialize exactly one unavailable placeholder
  and one metrics template for every required observation.
- **FR-012**: The collector MUST remain resumable and MUST reject overwriting an
  existing workspace, observation, attachment, or final archive.
- **FR-013**: Finalization MUST reject missing, duplicate, skipped, partial,
  timed-out, unavailable, failed, malformed, or attachment-free mandatory
  evidence with actionable diagnostics.
- **FR-014**: The synthetic fixture MUST exercise every new semantic rule and
  MUST remain barred from the attended production path.
- **FR-015**: The Windows runbook MUST map each new desktop observation to issues
  #101, #104, #105, #106, #109, #111, #112, and #113, and retain #98/#96
  lifecycle traceability.
- **FR-016**: S047 MUST produce one committed-source local demo MSI whose
  filename and embedded version unmistakably say `s047-demo`, while its numeric
  ProductVersion remains suitable for MSI compilation and inspection.
- **FR-017**: The demo handoff MUST record the full source commit, ProductCode,
  ProductVersion, embedded version, byte size, SHA-256, build timestamp, artifact
  class, and compiled inspection result.
- **FR-018**: The local demo MUST pass all eight canonical repository gates,
  focused release-gate tests, PowerShell parser checks, installer-source checks,
  and compiled-MSI inspection before it is reported ready.
- **FR-019**: Attended local-demo results MUST be recorded separately from the
  later formal candidate and MUST NOT close an issue whose criteria require the
  formal evidence class.
- **FR-020**: Any S047 product correction MUST begin with a failing regression
  test and preserve the proof-before-commit rule for recurring error behavior.
- **FR-021**: S047 MUST stop before push or pull-request creation and report each
  unfinished attended, formal-candidate, issue-disposition, or publication step.
- **FR-022**: S047 MUST NOT implement or imply completion of Post-v1 issue #102.

### Key Entities

- **Required Observation**: A fixed scenario identity, issue mapping, environment
  class, mandatory metrics, and attachment requirements enforced by promotion.
- **Desktop Qualification Metrics**: The native palette, DPI, window-size,
  control-state, input-device, row-population, and interaction results for one
  desktop observation.
- **Evidence Bundle**: The existing versioned candidate identity, operator
  attestation, environments, observations, and hashed attachments, expanded by
  S047 desktop scenarios.
- **Local Demo**: A non-release MSI bound to one committed S047 source revision,
  inspection record, and optional pre-push attended observations.
- **Issue Disposition**: One issue's acceptance-to-observation map and explicit
  status based on the evidence class actually completed.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Removing any one of the eleven new desktop observations from a
  complete bundle produces a non-zero validation result naming that identity.
- **SC-002**: Mutation tests cover 100 percent of the new required metric groups,
  and each invalid group produces at least one actionable diagnostic.
- **SC-003**: A newly initialized workspace contains exactly one placeholder and
  one template for every canonical observation, with no duplicate identifiers.
- **SC-004**: A complete synthetic bundle containing all existing and new
  scenarios passes fixture validation, while the same bundle is rejected by the
  attended production entry point.
- **SC-005**: The local demo MSI identity in the handoff exactly matches its
  computed SHA-256, byte size, ProductVersion, ProductCode, source commit, and
  embedded S047 demo version.
- **SC-006**: Format, vet, lint, race, GUI, coverage, documentation, and
  automation gates all pass without exclusions added by S047.
- **SC-007**: Every included open issue has a documented evidence mapping and no
  issue is closed solely from local-demo, headless, or synthetic evidence.
- **SC-008**: Zero branch pushes, pull requests, tags, staged drafts, promotions,
  or public releases occur before the required user authorization.

## Assumptions

- The target release line remains v1.0.0, but S047 does not authorize a tag or
  release version.
- The development host can compile and inspect a local MSI but is not treated as
  a clean disposable attended environment.
- The operator will perform only the native interactions that cannot be proved
  honestly by headless automation.
- Existing S039 through S046 implementations are the baseline and are changed
  only when S047 produces a reproducible failure.
- Issue #98's full preserve, wipe, cancel, locked-file, multi-profile, repair,
  upgrade, and reinstall matrix remains mandatory in the later exact candidate
  even though the earlier S043 demo supplied useful wipe-path evidence.
