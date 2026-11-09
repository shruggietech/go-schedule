# S140 research and comparison decisions

Assessed on 2026-11-08 (America/New_York). Three read-only vendor research tasks covered 70 cells; parent source inspection covered the ten go-schedule cells. All 80 cells were individually assessed. Retrieval is documentation evidence, not installed vendor-product execution. The comparison keeps the existing product identities and capability definitions.

## Decisions

Use the existing JSON contracts and review helper. Preserve every historical snapshot/review byte, public v2.0.1 selection and original S139 observations. A new framework, runtime feature or native walkthrough adds no value to this source checkpoint. Candidate-tag links describe the prepared release namespace and become resolvable after its later accepted-source tag.

## Selected product contexts

- **go-schedule**: Prepared v2.0.2 source; independent local/remote daemons, Apache 2.0. Published packages remain v2.0.1. [Official source](https://github.com/shruggietech/go-schedule/tree/d5bb2632d24bbeb0615b041533acae4866761d28).
- **JAMS**: Current Core/Advanced offering and vendor documentation; installed Windows-server scheduler with platform agents; entitlements vary. [Official source](https://jamsscheduler.com/).
- **ActiveBatch**: Current vendor client-server platform; licensed cross-platform agents; Non-Cluster Failover is an add-in and interface entitlements vary by version. [Official source](https://www.advsyscon.com/activebatch/).
- **ScriptRunner**: ScriptRunner Enterprise 8 and current Microsoft/PowerShell platform material; licensed workflows/connectors and customer deployment context apply. [Official source](https://www.scriptrunner.com/platform).
- **Control-M**: Control-M 9.0.22 self-hosted documentation and Automation API; SaaS and API-version prerequisites differ. [Official source](https://documents.bmc.com/supportu/9.0.22/en-US/Documentation/Introduction_to.htm).
- **Azure Logic Apps**: Consumption and Standard; Standard cloud/ASE/hybrid options and preview capabilities differ. [Official source](https://learn.microsoft.com/en-us/azure/logic-apps/logic-apps-overview).
- **Resolve Actions**: Current platform brief; Actions cloud/hybrid Remote Executor documentation and Actions Pro 8.1 documentation, with separate edition, deployment and licensing prerequisites. [Official source](https://go.resolve.io/hubfs/Documents/Platform-Overview-Solution-Brief.pdf).
- **OpCon**: OpCon 26.0 On-Prem core plus current vendor-owned CLI; VisualCron/RPA and MFT integration documentation is separately scoped, with license and deployment prerequisites. [Official source](https://help.smatechnologies.com/opcon/core/v26.0/server-programs/introduction).

## Rating changes

Ten cells change status. These correct the evidence assessment; they do not imply ten newly shipped vendor features. All ten go-schedule ratings remain unchanged.

| Product / capability | Before | Prepared assessment | Reason and official evidence |
| --- | --- | --- | --- |
| jams / access | partial | supported | The features page now explicitly establishes audit logging alongside RBAC; the prior partial status was evidence-limited, not a documented missing capability. Core/Advanced distinctions remain. [Source](https://jamsscheduler.com/product/features). |
| jams / license | unknown | partial | The deployment model and license requirement are documented; partial denotes bounded commercial deployment evidence, not open-source parity. [Source](https://jamsscheduler.com/). |
| activebatch / license | unknown | partial | A vendor-specific subscription statement and current quote/architecture pages establish a bounded commercial deployment model; this does not establish source-code access. [Source](https://www.advsyscon.com/blog/perpetual-license-vs-subscription/). |
| scriptrunner / events | partial | supported | Version-specific workflow release material explicitly establishes webhook triggers, stronger than the earlier integration-only evidence. Individual external systems still require appropriate connectors. [Source](https://www.scriptrunner.com/blog-admin-architect/scriptrunner-enterprise-version-8-workflow-engine). |
| scriptrunner / highavailability | unknown | partial | Current operational guidance gives concrete multi-server executability and target failover mechanisms, justifying partial support while preserving the unverified scheduler/in-flight boundary. [Source](https://www.scriptrunner.com/blog-admin-architect/scriptrunner-server-groups-control-action-execution). |
| scriptrunner / license | unknown | partial | The 2026 agreement explicitly documents commercial executable distribution, customer installation and source restriction; pricing clarifies connector scope. [Source](https://lp.scriptrunner.com/hubfs/AGB/2026/ScriptRunner-License-Agreement%2BTerms-and-Conditions-2026.pdf). |
| opcon / events | partial | supported | Core event support is broader than the earlier RPA-only evidence; integrations are not prerequisites for every event. [Source](https://help.smatechnologies.com/opcon/core/v26.0/job-components/events). |
| opcon / interfaces | partial | supported | Operator and programmatic surfaces are documented; the license prerequisite belongs in the note rather than implying a missing interface. [Source](https://github.com/smatechnologies/opcon-cli-java/blob/master/documentation/opcon-cli.md). |
| opcon / highavailability | unknown | supported | Replace unknown with directly documented controller failover, preserving its configuration and failback limits. [Source](https://help.smatechnologies.com/opcon/core/v26.0/Files/Database-Information/Failover-Scenarios). |
| opcon / license | unknown | partial | The versioned licensing page supports a bounded deployment/license description; an open-source entitlement is not established. [Source](https://help.smatechnologies.com/opcon/core/v26.0/server-programs/licensing). |

## Evidence limits and dates

The snapshot date is the assessment date, not a statement that every source was freshly verified. Five cells retain their earlier evidence dates: ActiveBatch MCP, ScriptRunner MCP, Control-M MCP and OpCon MCP retain 2026-11-06; Control-M licensing retains 2026-10-27. Missing or conflicting documentation never establishes feature absence.

- ActiveBatch MCP availability remains unestablished for the selected product. A sibling product is not evidence for it.
- ScriptRunner material conflicts on MCP readiness; retain unknown rather than choosing a marketing date.
- The same Control-M MCP endpoint returned conflicting preview/general-availability descriptions and API prerequisites. Preserve the earlier preview assessment; the selected self-hosted product remains 9.0.22. Current selected-version licensing terms were not verified.
- Cloud-only Otto REST tooling does not establish an on-premises OpCon MCP endpoint. Retain the prior unknown and its date.
- ScriptRunner redundant execution server groups do not establish scheduler ownership or in-flight recovery. ActiveBatch failover is an add-in. Resolve Actions executor routing and Actions Pro 8.1 cluster behavior are separate architectures. OpCon failover requires configured primary/secondary topology and does not prove universal in-flight recovery. Logic Apps zone resilience and regional disaster-recovery design are separate.

## go-schedule source and bounded documentation finding

Accepted S139 source is d5bb2632d24bbeb0615b041533acae4866761d28. Structured acknowledgement and guidance corrections are included; source-tree CLI QA and opt-in pre-entry incident capture are not new installed CLI modes. Actual host/WSL connection tests were blocked by the test environment, and native UI/keyring behavior was not observed under the required isolation. Source and headless UI evidence do not replace those observations. Original #355/#356 remain open; independent daemons do not provide the owner-deferred clustered lane. See [S139 audit](../139-full-product-audit-v202/audit.md).

The existing cron guide incorrectly calls conversion and direct CLI/desktop cron input unreleased after 0.8.0. The v2.0.1 tag contains internal/cli/cron.go; the tagged cron implementation, desktop task-group implementation and guide are byte-identical to accepted main. Deduplicated bug [#453](https://github.com/shruggietech/go-schedule/issues/453) removes only that obsolete release-status paragraph, preserving dialect and fidelity boundaries.

## Review decisions

Each capability receives changed because at least one explanation, context, source qualification or rating changes materially; a refreshed date or tag link alone would not justify changed. Scheduling, dependencies, remote and observability clarify actual deployment/runtime scope; events, interfaces, access, high availability and licensing also contain rating changes; MCP explicitly records unresolved/conflicting evidence. The ten decisions bind the final dataset, changelog section and tag-specific notes through the existing helper.

Private per-cell retrieval records and original failures remain in ignored local evidence. The public dataset carries every cell status, concise explanation, direct official source and evidence date. No raw workstation identity, credentials or vendor page dumps are published.

## Preparation correction #454

The first full foreground aggregate passed the first seven gates, then correctly rejected the newly prepared seven-bullet notes at automation. Consolidate the cron correction into the existing diagnostic/documentation highlight to retain six bullets. This is a pending S140 preparation defect, not a published-package or checker defect. The old content-bound review then independently rejected changed notes/changelog; preserve that negative result, regenerate the review through the existing helper, reassess all ten decisions and rerun the affected format/docs/automation and release guards. No product code, test, dependency, workflow, frontend or earlier comparison bytes change.
