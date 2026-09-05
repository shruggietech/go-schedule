# Security Requirements Checklist: GitHub Security Baseline

**Purpose**: Test the completeness, clarity, consistency, and measurability of the security baseline requirements before implementation
**Created**: 2026-08-26
**Feature**: [spec.md](../spec.md)

**Note**: This checklist evaluates the written requirements, not the eventual GitHub settings or workflow implementation.

## Requirement Completeness

- [x] CHK001 Are both the primary private-report route and its fallback channel specified? [Completeness, Spec §FR-001–FR-004]
- [x] CHK002 Are ownership and response-time requirements defined for incoming private reports? [Completeness, Spec §FR-004]
- [x] CHK003 Are all requested secret-protection capabilities named individually rather than grouped into an ambiguous security toggle? [Completeness, Spec §FR-005–FR-006]
- [x] CHK004 Are analysis events, cadence, language, permissions, build constraints, and result publication all covered? [Completeness, Spec §FR-007–FR-010]
- [x] CHK005 Are regression requirements defined for both action-runtime drift and security-workflow contract drift? [Completeness, Spec §FR-009–FR-011]

## Requirement Clarity

- [x] CHK006 Is the meaning of a usable private-report route distinct from a public issue route? [Clarity, Spec §FR-002]
- [x] CHK007 Is the fallback address concrete and independently attributable to the repository owner? [Clarity, Spec §FR-003]
- [x] CHK008 Is “when supported” bounded by an evidence requirement for every requested control? [Clarity, Spec §FR-005–FR-006, FR-014]
- [x] CHK009 Is least privilege stated in terms of the capabilities the workflow needs and the sensitive capabilities it must not receive? [Clarity, Spec §FR-008]
- [x] CHK010 Is the headless analysis-build constraint specific enough to reject accidental desktop-library dependence? [Clarity, Spec §FR-010]

## Requirement Consistency

- [x] CHK011 Do private-route validation requirements remain consistent with the prohibition on fabricated advisory creation? [Consistency, Spec §FR-002, FR-015]
- [x] CHK012 Do hosted activation requirements remain consistent with the pre-push remote-mutation boundary? [Consistency, Spec §FR-001, FR-014–FR-015]
- [x] CHK013 Does adding a security workflow preserve the existing CI contract and product scope boundaries? [Consistency, Spec §FR-012, FR-016]

## Acceptance Criteria Quality

- [x] CHK014 Can every reporting-route outcome be measured without submitting a vulnerability? [Measurability, Spec §SC-001–SC-003]
- [x] CHK015 Does the secret-control outcome require a status for 100% of the requested controls? [Measurability, Spec §SC-004]
- [x] CHK016 Are the required hosted analysis contexts and maintenance cadence objectively observable? [Measurability, Spec §SC-005]
- [x] CHK017 Are all required negative drift classes enumerated and countable? [Measurability, Spec §SC-006]

## Scenario and Edge-Case Coverage

- [x] CHK018 Are authenticated, unauthenticated, fallback, and maintainer-triage reporting scenarios all defined? [Coverage, Spec §User Story 1]
- [x] CHK019 Are pull-request, default-branch, scheduled, and offline-failure analysis scenarios all defined? [Coverage, Spec §User Story 2]
- [x] CHK020 Are account-plan limitations, token-scope limitations, newly discovered alerts, fork permissions, and compiled-build failures addressed? [Edge Cases, Spec §Edge Cases]
- [x] CHK021 Is the required behavior for partial or unavailable controls defined without weakening the baseline? [Recovery, Spec §FR-006, FR-014]

## Dependencies and Assumptions

- [x] CHK022 Are repository visibility, GitHub Actions availability, triage ownership, published contact, setup mode, and post-halt activation documented as assumptions? [Assumption, Spec §Assumptions]
- [x] CHK023 Are branch protection, Dependabot, dependency upgrades, product behavior, alert resolution, and organization-wide policy explicitly excluded? [Scope, Spec §Scope out]

## Notes

- All 23 requirements-quality checks pass against the clarified specification.
- Focus: private disclosure safety, least privilege, hosted-control evidence, failure honesty, and remote-mutation boundaries.
- Audience/timing: formal maintainer and third-party review before activation.
