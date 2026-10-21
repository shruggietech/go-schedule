# CLI contract requirements quality checklist

Purpose: Validate the completeness and clarity of the CLI and response requirements before implementation. Created: 2026-10-21. Feature: [specification](../spec.md). Audience: implementation author and PR reviewer, standard depth.

## Completeness and clarity

- [x] CHK001 Are all 12 independently closeable outcomes explicitly represented? [Completeness, Spec FR-001 through FR-012]
- [x] CHK002 Are usage/validation and runtime categories defined separately for parser, local file and server errors? [Clarity, Spec FR-001, Edge Cases]
- [x] CHK003 Is refusal before connection or host mutation required for invalid operands? [Completeness, Spec FR-002]
- [x] CHK004 Are destination and daemon selection independently expressible for every target kind? [Clarity, Spec FR-004]
- [x] CHK005 Are duplicate and unknown binding identities both covered before usable plan creation? [Coverage, Spec FR-005]
- [x] CHK006 Is complete cron reporting required even when the resulting status is unsuccessful? [Consistency, Spec FR-006]

## Diagnostic and state meaning

- [x] CHK007 Are stdout and stderr roles explicit for deprecated structured commands? [Clarity, Spec FR-003]
- [x] CHK008 Are help, error and human-list requirements specific enough to assess without guessing? [Measurability, Spec FR-007, FR-008, FR-010]
- [x] CHK009 Are unavailable identities and one-time token delivery distinguished? [Coverage, Spec FR-009]
- [x] CHK010 Does the specification distinguish saved and observed watcher state without claiming immediate reconciliation? [Consistency, Spec FR-011, Clarifications]
- [x] CHK011 Does warning guidance preserve actual access behavior and metadata? [Consistency, Spec FR-012]

## Compatibility and dependencies

- [x] CHK012 Are credential privacy and relevant context differentiated? [Clarity, Spec FR-013, Clarifications]
- [x] CHK013 Are existing consumers, valid commands and persisted schemas retained? [Coverage, Spec FR-013]
- [x] CHK014 Are S108 prerequisites and the pre-publication boundary explicit? [Dependencies, Spec Assumptions]

## Notes

All 14 items pass against the clarified specification. No checklist extension hooks are registered. The checklist prerequisite helper expects plan.md although the repository orders checklist before plan; the paths-only prerequisite result and feature.json identify the active specification. The normal full prerequisite check is repeated after planning.
