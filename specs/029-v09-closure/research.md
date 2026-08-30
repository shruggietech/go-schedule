# Research: v0.9 Closure and Maintenance Automation

## R1. Lifecycle vocabulary

**Decision**: Use seven plain states: Draft, Ready, In Progress, Implemented, Deferred, Superseded, and Abandoned.

**Rationale**: They cover the real planning transitions without mixing local implementation with hosted publication. Compound historical states are rejected because they require later cleanup that the process never performs reliably.

## R2. Historical evidence and task reconciliation

**Decision**: Cite merge PRs where they exist, otherwise cite the delivery commit and nearest release. Reconcile stale checkboxes only when repository history demonstrates the feature landed. Preserve an explicit waiver for manual evidence that was never collected.

**Rationale**: Git history and merged PRs are objective. Marking a manual walkthrough as performed without evidence would be false, while leaving an unexplained open task would keep the contradiction.

## R3. Lifecycle enforcement location

**Decision**: Add a portable offline checker invoked by the current automation gate and tested with temporary fixtures.

**Alternatives rejected**:

- A new GitHub workflow adds ceremony and another policy surface.
- A network-dependent issue-state check makes local verification flaky and unavailable offline.
- Documentation alone does not catch the drift described by #42.

## R4. Dependabot cadence and grouping

**Decision**: Weekly Go modules, monthly GitHub Actions, at most five routine PRs per ecosystem, with compatible minor and patch updates grouped and major/security updates separate.

**Rationale**: This adapts the proven `shruggie-graph` configuration to a small Go repository. It minimizes routine review noise without hiding high-risk changes.

## R5. Secret-scanning scope

**Decision**: Activate exactly push protection, non-provider patterns, and validity checks, then read each status back.

**Rationale**: Those are the explicit remaining items in #39. AI detection and delegated dismissal have distinct semantics and no stated need.

## R6. Windows icon issue disposition

**Decision**: Move #33 to Post-v1 and P3 while retaining `needs: verification` and leaving it open.

**Rationale**: The report predates the brand rewrite, has no current reproduction, and requires a manual environment session the maintainer explicitly deferred. Closing it would discard a plausible issue; retaining it as a v0.9 P1 would misstate urgency.
