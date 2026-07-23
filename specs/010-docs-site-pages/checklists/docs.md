# Checklist: Documentation-site requirements quality

**Purpose**: Unit-test the *requirements* for feature 010 (docs site + README
consolidation) for completeness, clarity, consistency, and measurability before
planning. Tests what the spec says, not whether the site works.
**Created**: 2026-07-23
**Feature**: [spec.md](../spec.md)

## Requirement Completeness

- [x] CHK001 - Are requirements defined for every capability the site must expose — search, spanning navigation, and next/previous movement? [Completeness, Spec §FR-002]
- [x] CHK002 - Is the requirement that the site build from the docs folder (source == served content, no committed build output) stated unambiguously? [Completeness, Spec §FR-001]
- [x] CHK003 - Are front-matter requirements specified for *every* documentation page, not just a representative subset? [Completeness, Spec §FR-004]
- [x] CHK004 - Are requirements defined for the automated docs-check's three failure conditions (broken on-disk link, stale pointer, missing front matter)? [Completeness, Spec §FR-007]
- [x] CHK005 - Is the requirement to run the docs-check in CI on the same triggers as existing checks stated? [Completeness, Spec §FR-008]
- [x] CHK006 - Are requirements documented for every inbound reference that must be repointed to the site (README link, issue-form contact links)? [Completeness, Spec §FR-010]
- [x] CHK007 - Is the pinned-artifact / dated-CHANGELOG-decision obligation captured as a requirement rather than left implicit? [Completeness, Spec §FR-011]

## Requirement Clarity

- [x] CHK008 - Is "single source of truth" defined concretely enough to distinguish a compliant pointer README from a non-compliant duplicate? [Clarity, Spec §FR-009]
- [x] CHK009 - Is the boundary "a relative link that escapes the documentation folder" specified precisely enough to be checkable? [Clarity, Spec §FR-005]
- [x] CHK010 - Is "required front matter" enumerated (which fields) so a page can be objectively judged complete or not? [Clarity, Spec §FR-004]
- [x] CHK011 - Is the docs-check's explicit non-goal (ignore http/https links, no network dependency) stated clearly? [Clarity, Spec §FR-007]
- [x] CHK012 - Is the theme/generator compatibility constraint expressed as a requirement (theme pinned to a version the host's builder supports) rather than buried as prose? [Clarity, Spec §Assumptions, Edge Cases]

## Requirement Consistency

- [x] CHK013 - Do the "no committed build artifacts" requirement (§FR-001) and the "served from docs folder" requirement agree — i.e., no requirement implies checking in generated HTML? [Consistency, Spec §FR-001]
- [x] CHK014 - Are the within-docs-links-resolve requirement (§FR-006) and the escape-links-are-absolute requirement (§FR-005) mutually consistent and jointly exhaustive of link cases? [Consistency, Spec §FR-005/FR-006]
- [x] CHK015 - Is the home-page requirement (§FR-003) consistent with the "docs index renders in-repo too" requirement (no duplicate index file)? [Consistency, Spec §FR-003]

## Acceptance Criteria Quality

- [x] CHK016 - Is "zero 404s across the published site" objectively measurable as written? [Measurability, Spec §SC-002]
- [x] CHK017 - Is "reach any topic in at most two clicks or one search" verifiable without implementation detail? [Measurability, Spec §SC-001]
- [x] CHK018 - Is the "exactly one operator settings change, no ongoing deploy workflow" outcome measurable? [Measurability, Spec §SC-006]
- [x] CHK019 - Does each functional requirement map to at least one acceptance scenario or success criterion? [Traceability, Spec §FR-001..FR-012]

## Scenario & Edge-Case Coverage

- [x] CHK020 - Are requirements defined for the "link points outside the docs set" case (absolute repo URL) rather than left to implementer discretion? [Coverage, Edge Case, Spec §FR-005]
- [x] CHK021 - Is the "page added without navigation metadata" failure path covered by a requirement (the check catches it)? [Coverage, Edge Case, Spec §FR-007]
- [x] CHK022 - Is the "front matter must not degrade in-repo reading" case addressed as a requirement/assumption? [Coverage, Edge Case, Spec §FR-003]
- [x] CHK023 - Are the deferred scope boundaries (versioned docs, custom domain, branding assets) explicitly excluded so they are not silently expected? [Coverage, Spec §Assumptions]

## Dependencies & Assumptions

- [x] CHK024 - Is the assumption that the host serves and builds a site from a docs folder on the default branch stated and marked as an external dependency? [Assumption, Spec §Assumptions]
- [x] CHK025 - Is the operator-performed settings change documented as an out-of-band dependency the repo content cannot satisfy itself? [Assumption, Spec §SC-006, Assumptions]
- [x] CHK026 - Is the "no Go code touched; existing quality/test/race/coverage/bench gates stay green" constraint recorded as a requirement? [Assumption, Spec §FR-012]

## Notes

- All items test the **requirements**, not the eventual site. Items that fail
  here are spec gaps to close before or during `/speckit-plan`, not
  implementation bugs.
