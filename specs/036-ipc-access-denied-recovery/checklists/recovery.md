# Recovery Requirements Checklist: IPC Access-Denied Recovery

**Purpose**: Review recovery, classification, security, and evidence requirement quality **Created**: 2026-08-31 **Feature**: [spec.md](../spec.md)

## Requirement Completeness

- [x] CHK001 Are all required transport and API failure categories explicitly defined? [Completeness, Spec FR-001]
- [x] CHK002 Are the incident creation, update, clear, and cancellation states specified? [Completeness, Spec FR-003, FR-008, FR-009]
- [x] CHK003 Are Retry and Exit requirements defined for every connection category? [Coverage, Spec FR-005]
- [x] CHK004 Are native Windows evidence fields and unavailable-evidence handling documented? [Completeness, Spec FR-013, SC-007]

## Requirement Clarity

- [x] CHK005 Is access denial explicitly separated from daemon absence? [Clarity, Spec FR-001, FR-002]
- [x] CHK006 Is the maximum visible incident count objectively measurable? [Clarity, Spec SC-001]
- [x] CHK007 Is stale-token guidance conditioned on verified account and token evidence? [Clarity, Spec FR-006, FR-007]
- [x] CHK008 Is bounded retry behavior quantified with initial and maximum delays? [Clarity, Clarifications]

## Requirement Consistency

- [x] CHK009 Do recovery requirements preserve the IPC authorization boundary consistently? [Consistency, Spec FR-014, Out of Scope]
- [x] CHK010 Do background and user-triggered retries share one incident lifecycle without conflicting behavior? [Consistency, Spec FR-008, FR-009]

## Scenario and Edge-Case Coverage

- [x] CHK011 Are simultaneous model, calendar, and stream failures covered as one incident? [Coverage, Spec FR-011]
- [x] CHK012 Are unknown diagnostics, in-flight exit, repeated Retry, and connectivity return covered? [Coverage, Edge Cases]
- [x] CHK013 Are unrelated API operation errors kept distinct from transport incidents? [Coverage, Spec FR-010, FR-012]

## Acceptance Criteria Quality

- [x] CHK014 Can dialog suppression, deduplication, backoff bounds, and recovery be measured independently? [Measurability, Spec SC-001 through SC-006]
- [x] CHK015 Does the done gate require both canonical verification and native Windows evidence? [Traceability, Spec FR-013, FR-015]
