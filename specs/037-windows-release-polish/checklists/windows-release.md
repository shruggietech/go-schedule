# Windows Release Requirements Checklist: Windows Release Polish

**Purpose**: Review the clarity and completeness of the Windows GUI and MSI requirements before implementation **Created**: 2026-08-31 **Feature**: [spec.md](../spec.md)

## Requirement Completeness

- [x] CHK001 Are preferred, capped, fallback, and hard-bound sizing rules all defined? [Completeness, Spec FR-001 through FR-004]
- [x] CHK002 Is launch-monitor selection specified for constrained secondary-monitor behavior? [Completeness, Spec FR-005]
- [x] CHK003 Are user-controlled window states explicitly preserved? [Completeness, Spec FR-007]
- [x] CHK004 Are all installer identity fields that must remain stable named? [Completeness, Spec FR-011]

## Requirement Clarity

- [x] CHK005 Is the 90 percent calculation ordered after physical-to-logical conversion? [Clarity, Spec FR-002]
- [x] CHK006 Is the exact approved Subject value stated without stylistic alternatives? [Clarity, Spec FR-009]
- [x] CHK007 Is the former specification 003 behavior explicitly superseded? [Clarity, Spec FR-008]

## Acceptance Criteria Quality

- [x] CHK008 Are required display cases and exact dimensions measurable? [Measurability, Spec SC-001]
- [x] CHK009 Is effective 800x600 reachability tied to a measurable startup viewport? [Measurability, Spec SC-002]
- [x] CHK010 Does artifact evidence require both hash and observed Subject? [Measurability, Spec SC-004]

## Evidence Boundaries

- [x] CHK011 Are source, compiled-MSI, and native Explorer evidence kept distinct? [Consistency, Spec Edge Cases]
- [x] CHK012 Is unavailable native evidence prohibited from being reported as passed? [Consistency, Spec Assumptions]
