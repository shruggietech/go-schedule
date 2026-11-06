# S133 Research

## Selected baseline

Decision: integrate PRs #420-#423 into one candidate and select Wails runtime/CLI 2.17.1 together. Rationale: one common build/release boundary and a concrete pin-alignment gap. Alternatives: merge module-only #420 and retain older tooling, or launch a broader dependency campaign. Both are less coherent.

## Maintenance classification

Decision: v2.0.1. Presentation, first-use guidance and dependency corrections preserve public contracts and current formats. No change to /v2 module identity or development fallback is needed; releases inject their full tag.

## Build contract

Decision: derive automation's expected CLI version from desktop/go.mod, require every hosted invocation and explicit local build pin to match, and use existing negative fixtures. This removes independent hard-coded checker truth. No shared version service or new command wrapper.

## Release boundary

Decision: prepare new tag-specific records and source-bound review; preserve historical v2.0.0 and public selection. README may show the candidate health example while explicitly linking actual v2.0.0 downloads. No tag is created. Existing same-slice staging/promotion procedure follows later owner decisions.

## Evidence

The installed plan command dispatches read-only research. Existing S129-S132 artifacts and current source determine exact candidate QA/comparison commands. Local attended UI/installation and native macOS remain untested under the standing waiver. Passing source evidence is not staged-byte acceptance.

## Current primary-source comparison review (2026-11-06)

The dispatched research inspected official sources for all seven competitors. All ten capability classifications remain unchanged; S133 changes presentation and dependencies, not scheduler capabilities. Preserve vendor edition, preview, deployment and unknown qualifications.

| Product | Sources and bounded conclusion |
| --- | --- |
| JAMS | jamsscheduler.com overview/MCP and jamstechdocs.z19.web.core.windows.net security/HA retain current execution, ACL and failover scope; AI job creation remains roadmap. |
| ActiveBatch | advsyscon.com/activebatch overview, architecture and security retain active-standby add-in and unknown MCP qualifiers. |
| ScriptRunner | scriptrunner.com platform/workflows and February/April trigger/launch articles retain contradictory present/planned MCP descriptions. |
| Control-M | documents.bmc.com 9.0.22 overview, jobs, settings, MCP and HA plus controlm.github.io retain preview and topology limits. |
| Azure Logic Apps | learn.microsoft.com overview, security, API Center MCP and regional recovery retain preview, connector/tool and state-location limits. |
| Resolve Actions | official platform brief plus docs.resolve.io Pro 8.1 architecture and Dynamic MCP retain edition/deployment and in-flight recovery limits. |
| OpCon | help.smatechnologies.com v26.0 history/audit, VisualCron integration, official CLI and copyright-2020 PDF retain scoped and unknown qualifiers. |

The Control-M legacy xwiki license-context URL returned 403; preserve that cell's prior 2026-10-31 review date and unknown conclusion rather than claiming a fresh read. Other inspected cells and the candidate snapshot use 2026-11-06. New candidate source links select v2.0.1; historical v2.0.0 records remain unchanged.
