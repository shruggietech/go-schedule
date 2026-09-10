# Release Requirements Checklist: Cumulative v1.4.0 Release Preparation

**Purpose**: Test whether the release requirements are complete, unambiguous, measurable, and safe before implementation

**Created**: 2026-09-10

**Feature**: [spec.md](../spec.md)

**Audience**: Pull-request reviewers and the maintainer conducting the later release ritual

## Release identity

- [x] CHK001 Is the direct cumulative boundary from public v1.1.1 to v1.4.0 explicit? [Clarity, Spec §FR-001, §FR-002]
- [x] CHK002 Is retrospective publication of v1.2.0 and v1.3.0 explicitly excluded? [Coverage, Spec §FR-002]
- [x] CHK003 Are changelog, release-note, README, and installation-document identities defined consistently? [Consistency, Spec §FR-001 through §FR-005]
- [x] CHK004 Is the release-note shape objectively countable and tag-specific? [Measurability, Spec §FR-003]

## Staging and provenance

- [x] CHK005 Are all metadata conditions that must fail before artifact mutation enumerated? [Completeness, Spec §FR-006]
- [x] CHK006 Is successful main CI bound to the exact immutable tagged commit rather than a branch or pull-request result? [Clarity, Spec §FR-006]
- [x] CHK007 Are the required draft artifact classes and draft-only state specified? [Completeness, Spec §FR-007, §SC-004]
- [x] CHK008 Are conflicting pre-existing tag and release states addressed? [Edge Case, Spec §Edge Cases]

## Installation and upgrade

- [x] CHK009 Are fresh installation and upgrade from the latest public v1.1.1 MSI independently required? [Coverage, Spec §FR-008]
- [x] CHK010 Are all state and service behaviors that must survive upgrade enumerated? [Completeness, Spec §FR-009]
- [x] CHK011 Are newly delivered notification, MCP, and remote-access capabilities required to remain disabled by default? [Security, Spec §FR-009, §SC-006]
- [x] CHK012 Is all attended evidence bound to the exact staged candidate and established observation set? [Traceability, Spec §SC-005]

## Promotion and administration

- [x] CHK013 Is no-rebuild promotion defined with candidate, evidence, tag, asset, and checksum validation? [Completeness, Spec §FR-010]
- [x] CHK014 Are issue, project, and milestone closure conditions separated from preparation completion? [Consistency, Spec §FR-011]
- [x] CHK015 Is the current authorization boundary explicit for tag, draft, installation, evidence, and publication operations? [Clarity, Spec §FR-012]
- [x] CHK016 Are failure states required to preserve a safe reviewed-source or draft-release state? [Recovery, Spec §Edge Cases]

## Notes

- All sixteen requirements-quality checks pass. The checklist uses release-review depth and covers identity, provenance, upgrade safety, promotion, and administrative closure.
