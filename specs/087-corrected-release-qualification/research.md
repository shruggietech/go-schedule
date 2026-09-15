# S087 research

**Decision**: Stage reviewed S086 main commit `951d864d9683ec3bdcb1e37535ecddd67738bf13`, not S087 documentation changes.

**Rationale**: This revision completed PR review and exact-commit main CI; the obsolete candidate predates the corrected desktop.

**Alternatives**: Stale bytes, an unreviewed branch, or an invented later version are invalid or outside scope.

**Decision**: Use the existing annotated tag-push Release workflow after backing up the old draft and assets, replacing the tag with an exact old-reference lease.

**Rationale**: Release has exact-main-CI and absent-or-draft preflight; no workflow_dispatch trigger is provided. Reuse existing staging instead of bypassing it.

**Decision**: Reuse S086 guest-local offline setup and retain the complete existing attended matrix.

**Rationale**: Setup and server-runner tests cannot prove medium-integrity Windows 11 UI, real profiles, display scaling, or physical input. Missing runtime/control capability must remain nonpassing, not be worked around with inferred observations.

**Finding**: The fresh candidate stalled during the client-side and service-side `SOFTWARE RESTRICTION POLICY: Verifying package` stages, before application custom actions. The client-side log advanced after approximately 121 seconds. [Microsoft Windows-Sandbox issue #68](https://github.com/microsoft/Windows-Sandbox/issues/68) reports the same stage as a Sandbox installation bottleneck. This is consistent with the observed pause but does not independently prove its cause on this host.

**Decision**: Retain the unchanged security posture and wait within the existing bounded diagnostic session. Do not adopt tracker suggestions that disable Smart App Control. Those are security-setting changes, not an application installer repair or credible proof of release behavior under normal controls.
