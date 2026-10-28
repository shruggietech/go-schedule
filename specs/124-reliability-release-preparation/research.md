# Research: S124

**Reviewed**: 2026-10-27

## Release identity

**Decision**: Prepare v1.7.0.

**Rationale**: The project declares [Semantic Versioning](https://semver.org/spec/v2.0.0.html); sections 6 and 7 distinguish fixes from backward-compatible public functionality. `cmd/gosched-qa/main.go`, `scripts/cli-qa.ps1` and `docs/cli-qa.md` expose additive incident/replay selections since v1.6.0. No major/schema change was found.

**Alternatives considered**: Proposed v1.6.1 is inappropriate for these additions. A major increment is unsupported by the source delta.

**Decision**: Leave buildinfo at its development fallback and keep package identity derived from an actual future tag.

**Rationale**: Existing workflow injection supplies shipped CLI, daemon, desktop and MSI versions. Source preparation cannot create an actual Windows candidate manifest or future asset hashes. Public v1.6.0 selection and immutable tagged content stay unchanged.

## Comparison review

**Decision**: Reassess all 80 cells and create v1.7.0 records. Preserve all status values except qualified Resolve MCP and HA, now supported for Actions Pro 8.1.

**Rationale**: [Pro 8.1 release notes](https://docs.resolve.io/pro/Release%20Notes/release-notes-8-1-0-0/) and [Dynamic MCP configuration](https://docs.resolve.io/pro/Development%20Guides/Dynamic%20MCP%20Server/create-and-configure/) establish shipped selected-runbook endpoints; [architecture](https://docs.resolve.io/pro/Getting%20Started/architecture-overview/) documents clustered components, state replication and primary/backup messaging. Neither establishes arbitrary in-flight workload recovery.

**Alternatives considered**: Retaining stale unknowns would understate current evidence. Expanding OpCon to hosted failover would change its selected on-premises v26.0 scope. ScriptRunner conflicting MCP announcements do not establish release availability; retain unknown. Control-M and Logic Apps preview limitations remain explicit.

**Evidence refresh**: JAMS [platform](https://jamsscheduler.com/), [security](https://jamstechdocs.z19.web.core.windows.net/docs/concepts/configuration-and-security), [MCP](https://jamsscheduler.com/product/mcp-server/), [HA](https://jamstechdocs.z19.web.core.windows.net/docs/highavailability); ActiveBatch [platform](https://www.advsyscon.com/activebatch/), [architecture](https://www.advsyscon.com/activebatch/core-capabilities/activebatch-architecture/), [security](https://www.advsyscon.com/activebatch/security-driven-automation/); ScriptRunner [platform](https://www.scriptrunner.com/platform), [workflows](https://www.scriptrunner.com/product/workflows), [triggers](https://www.scriptrunner.com/blog-admin-architect/reinventing-enterprise-it-automation-with-scriptrunner-part-2-automation-trigger-types), [launch announcement](https://www.scriptrunner.com/news/scriptrunner-unveils-the-first-independent-agentic-workflow-platform-for-microsoft-ecosystems); Control-M [9.0.22 overview](https://documents.bmc.com/supportu/9.0.22/en-US/Documentation/Introduction_to.htm), [MCP preview](https://documents.bmc.com/supportu/9.0.22/en-US/Documentation/MCP_Server.htm), [HA](https://documents.bmc.com/supportu/9.0.22/en-US/Documentation/High_Availability.htm); Logic Apps [overview](https://learn.microsoft.com/en-us/azure/logic-apps/logic-apps-overview), [security](https://learn.microsoft.com/en-us/azure/logic-apps/set-up-security-permissions), [MCP](https://learn.microsoft.com/en-us/azure/logic-apps/create-mcp-server-api-center), [recovery](https://learn.microsoft.com/en-us/azure/logic-apps/multi-region-disaster-recovery); OpCon [v26 history](https://help.smatechnologies.com/opcon/core/v26.0/Files/UI/Enterprise-Manager/Using-History-Management), [audit](https://help.smatechnologies.com/opcon/core/v26.0/Files/UI/Enterprise-Manager/Using-Audit-Management), [vendor CLI](https://github.com/smatechnologies/opcon-cli-java/blob/master/documentation/opcon-cli.md). Each selected cell retains its directly relevant source and edition context.

## Candidate QA

**Decision**: Build and evaluate the new assembled source through existing Windows/WSL coordinator contracts with incidents enabled; review all current messages.

**Rationale**: S122 provides an intact fixed-schema precedent, not a current result. The builder reports 0.0.1-dev and has no release-version override. Actual run/source hashes and observed outcomes are authoritative. All six inputs stay frozen through capture and later processing.

**Alternatives considered**: Relabeling old reports or recompiling after capture would break provenance. Repeating the unchanged S123 historical investigation would add no new evidence and is outside this release-preparation goal.

## Historical uncertainty and delivery

**Decision**: Keep #355/#356 open with their original uncertainty; retain the historical full-tree mismatch and controller/collector limitations. Keep #394 open after source preparation.

**Rationale**: Passing current runs does not establish historical causes or actual release delivery. Source PR, owner merge, tag/staging and concrete-asset promotion are separate states.

## Demonstrated release-validation defect (#395)

**Decision**: Derive the repository smoke case's selected version from the first declared changelog version and retain fixed-version negative fixtures.

**Rationale**: The first S124 aggregate failed because TestCurrentCandidate reused v1.6.0's fixture constant after correct v1.7.0 preparation. Advancing the constant to v1.7.0 would repeat the defect at the next release. The existing strict validator remains the authority for health/version, date, badge and notes integrity.

**Alternatives considered**: A new production selector/API or weakened metadata assertions is unnecessary. A test-only selection helper plus future-version and malformed-selection coverage resolves the demonstrated boundary.
