# Data Model: v0.9 Closure and Maintenance Automation

## Specification lifecycle record

- Feature identifier and title
- Lifecycle state from the canonical vocabulary
- Delivery evidence (pull request, release, or commit)
- Task reconciliation note when historical metadata needed correction

### State transitions

```text
Draft -> Ready -> In Progress -> Implemented
  |        |          |
  +--------+----------+-> Deferred
  +--------+----------+-> Abandoned
Any non-final state ------> Superseded
```

`Implemented`, `Superseded`, and `Abandoned` are final for a feature record. A materially new body of work receives a new specification rather than reopening delivered history.

## Dependency proposal policy

- Ecosystem (`gomod` or `github-actions`)
- Root directory
- Cadence
- Open proposal limit
- Group eligibility
- Applied labels

## Hosted security control record

- Control name
- Requested state
- Observed state
- Provider constraint, when activation is unavailable
- Observation date

## Deferred issue record

- Issue number
- Open state
- Milestone
- Priority
- Evidence requirement
- Deferral rationale
