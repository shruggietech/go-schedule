# S142 v2.0.3 comparison research

Reviewed on 2026-11-09 (owner-local date) against the v2.0.2 comparison snapshot and merged S141 source. This patch corrects Windows build identity and release metadata. The ten capability decisions are no-change: scheduling, events, dependencies, interfaces, remote, observability, access, mcp, highavailability and license. No new scheduler, CLI, API or desktop interaction capability is introduced.

The proposal uses the source-bound product context: v2.0.3 source, independent local/remote daemons; Apache 2.0. Publication availability belongs to the page current baseline rather than this immutable comparison snapshot. The top-level review date and Go-schedule source review dates advance to 2026-11-09; source links advance to v2.0.3. All vendor notes, statuses, source URLs and assessment dates remain identical to v2.0.2, including explicit older uncertainty dates. Reopening an overview does not imply reassessing every product capability.

Official sources reopened successfully during this research:

- JAMS scheduling: https://docs.jamsscheduler.com/docs/concepts/scheduling (natural-language dates, intervals, recurring windows and event-trigger distinction still documented).
- ActiveBatch interfaces: https://www.advsyscon.com/activebatch/core-capabilities/activebatch-architecture/interfaces/ (Console and Web Console remain documented; selected interface coverage unchanged).
- ScriptRunner platform: https://www.scriptrunner.com/platform (scheduled PowerShell actions, central credentials, delegation and logging remain documented).
- Control-M 9.0.22 architecture: https://documents.bmc.com/supportu/9.0.22/en-US/Documentation/Architecture.htm (central servers, clients, agents and agentless hosts still documented).
- Azure Logic Apps overview: https://learn.microsoft.com/en-us/azure/logic-apps/logic-apps-overview (Consumption/Standard selection and cloud/on-premises/hybrid workflow context remain documented).
- Resolve platform brief: https://go.resolve.io/hubfs/Documents/Platform-Overview-Solution-Brief.pdf (SaaS/on-premises/hybrid, orchestration, visual/API interfaces and RBAC context still documented).
- OpCon 26.0 on-premises introduction: https://help.smatechnologies.com/opcon/core/v26.0/server-programs/introduction (SAM, agent communication and supporting-service failover context remain documented).

No material claim change was established by these source checks. They do not resolve the existing ActiveBatch, ScriptRunner and OpCon MCP unknowns, Control-M MCP availability ambiguity or selected Control-M licensing unknown. Preserve those statuses and their original dates. No vendor runtime was installed or tested. No desktop focus was acquired. No tracked files were edited by this research agent. Final dataset/changelog/release-note hashes must be generated from the parent final content using the existing review-release tooling.

## Execution decision

Owner authorization covers the entire release, including source merges, tag, staging, inspection and promotion. Preserve technical identity checks and remove human confirmation stops. Existing implemented Windows handling/diagnostics can close independently of historical forensic uncertainty; unimplemented deferred capabilities remain open.
