# Feature Specification: Public website and project discovery

**Feature Branch**: `codex/106-public-web-presence`

**Created**: 2026-10-19

**Status**: Implemented

**Delivery**: [PR #291](https://github.com/shruggietech/go-schedule/pull/291) delivered the S106 implementation. The owner's subsequent 2026-10-19 cutover authorization establishes https://go-schedule.io/ through organization-verified GitHub Pages Actions and Cloudflare DNS, with enforced HTTPS and verified IPv4/IPv6 service. Follow-up [PR #302](https://github.com/shruggietech/go-schedule/pull/302) retires the temporary legacy publisher and updates discovery links; [delivery record](delivery.md) records the deployment and verification evidence.

**Input**: Owner kickoff of S106 under autopilot, full six-issue bundle #278, #279, #280, #281, #9 and #282, coordinated by #276. Stop before branch publication, pull request creation or live cutover.

## Clarifications

### Session 2026-10-19

- Autopilot resolution: preserve canonical docs source, static search and legacy routes; publish no site/domain changes before owner review. Authentication is preparation, not cutover authority.
- Coverage review: functional scope, entities, interaction, non-functional needs, external dependencies, failure/recovery, constraints, terminology and completion signals are clear; no human clarification questions required.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Understand and start using the project (Priority: P1)

A visitor understands the scheduling and workload-automation capabilities, sees practical examples and follows a clear download or documentation path from the branded home or concise README.

**Why this priority**: The current documentation-index home and long README hide useful capability and first use.

**Independent Test**: Inspect the home and README for purpose, platforms, interfaces, one complete first-use example and all required project links.

**Acceptance Scenarios**:

1. **Given** a new visitor, **When** they open the home, **Then** they can find installation, documentation, releases, changelog, comparison, repository and company links.
2. **Given** a reader choosing appearance, **When** they select light or dark mode, **Then** their choice persists and keyboard/reduced-motion preferences remain usable.
3. **Given** a README reader, **When** they follow the first-use example, **Then** its supported commands are sufficient and detailed workflow material is linked separately.
4. **Given** the README footer, **When** the company name is selected, **Then** it opens https://shruggie.tech and remains the final subordinate text line.

### User Story 2 - Find complete documentation (Priority: P1)

Users and maintainers browse and search the existing documentation through a consistent public site while old document links remain useful.

**Why this priority**: The migration must preserve the reference material people already use.

**Independent Test**: Compare the published-document inventory with exported routes, search results, assets and source-edit links.

**Acceptance Scenarios**:

1. **Given** a current document or asset, **When** it is migrated, **Then** its content and usable destination remain available.
2. **Given** a legacy public document link, **When** it is opened, **Then** the reader reaches its corresponding content.
3. **Given** a search term, **When** a reader searches, **Then** document results are available without a runtime server.

### User Story 3 - Compare capabilities and maintain release truth (Priority: P2)

Readers compare released go-schedule capability with seven named alternatives, and maintainers connect relevant opportunities to planning records and refresh the comparison with every release.

**Why this priority**: Competitive ambition needs concrete evidence and a repeatable maintenance path.

**Independent Test**: Inspect the comparison data and run the release-review contract against matching and mismatched release identities.

**Acceptance Scenarios**:

1. **Given** a comparison cell, **When** it is read, **Then** its status, bounded explanation, product edition/deployment context and official evidence/date are available.
2. **Given** planned or main-only capability, **When** the page is rendered, **Then** it is distinct from released capability.
3. **Given** release preparation or promotion, **When** comparison review is absent, stale or identifies another release, **Then** publication cannot proceed through either supported promotion path.
4. **Given** a reviewed release with no capability change, **When** its explicit no-change review is supplied, **Then** the review contract accepts it without fabricating changed competitor facts.

### User Story 4 - Reach the canonical project domain (Priority: P2)

Visitors reach the project through go-schedule.io with HTTPS and compatible www/old-site links; maintainers can repeat and recover deployment.

**Why this priority**: Domain ownership, source deployment and path changes must agree.

**Independent Test**: Check the prepared record/settings plan locally; after separately authorized cutover inspect DNS, TLS and representative public paths.

**Acceptance Scenarios**:

1. **Given** owner-authenticated domain access, **When** the authorized cutover runs, **Then** the apex serves the repository site over HTTPS and www leads to the canonical apex.
2. **Given** interrupted or conflicting setup, **When** a maintainer follows recovery instructions, **Then** only intended records/settings are changed and unrelated DNS is preserved.

### Edge Cases

- A document has literal JSX-like text, relative images, tables, code fences, Mermaid or legacy heading anchors.
- Search has no matches; JavaScript or local preference storage is unavailable.
- A competitor overview does not establish absence of a capability; the cell remains unknown rather than asserting absence.
- A candidate tag is not public yet; the public page continues to identify the released baseline.
- Authentication or organization permission is unavailable; local implementation proceeds and live domain work stays explicitly pending.
- Existing DNS or Pages settings conflict with the intended domain; report the conflict and retain a record-level rollback plan.
- A release entry point omits the comparison check; the shared verification contract identifies that omission.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The public home MUST explain concrete released capability and link repository, company, docs, changelog, releases, installation and comparison.
- **FR-002**: The site MUST use the approved brand, responsive layout, keyboard-accessible controls, persistent light/dark choice and reduced-motion behavior.
- **FR-003**: All current documentation and downloadable resources MUST have corresponding destinations, working relative links, navigation and source-edit links.
- **FR-004**: Search MUST work from the static publication without a runtime service, including an understandable empty state.
- **FR-005**: The README MUST prioritize purpose, supported platforms, installation and one first-use example, retain help/contribution/security/license links and end with the linked company attribution.
- **FR-006**: The comparison MUST cover JAMS, ActiveBatch, ScriptRunner, Control-M, Azure Logic Apps, Resolve Actions and OpCon through one evidence-backed dataset with per-cell status and product context.
- **FR-007**: Supported, partial, planned, unavailable and unknown capability MUST remain distinct; go-schedule claims MUST identify the released baseline.
- **FR-008**: Every release preparation/staging/promotion path MUST require a review record bound to the exact intended tag/source and an explicit capability change or no-change decision.
- **FR-009**: Comparison opportunities MUST link relevant existing planning records and MUST NOT imply automatic approval for every competitor feature.
- **FR-010**: The domain plan MUST bind go-schedule.io to this repository, preserve ownership verification, unrelated DNS and a bounded rollback path, and use HTTPS with canonical apex and www handling.
- **FR-011**: Deployment MUST retain a single canonical documentation authoring source and repeatable site/brand consumer updates.
- **FR-012**: The build phase MUST finish with local source, tests, documentation and review evidence; live settings, DNS, deployment, branch push and PR creation MUST remain behind the owner's pre-publication review boundary.

### Key Entities

- **Document**: Canonical path, title, navigation order, headings, content, legacy URLs and assets.
- **Comparison capability**: Stable identifier, product/cell status, bounded explanation, official source, review date and linked opportunities.
- **Release comparison review**: Intended tag, source revision, review date, capability decisions and released-baseline relationship.
- **Domain plan**: Zone identity, intended records, repository settings, ownership challenge and rollback inventory.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All six child outcomes have explicit implementation and delivery status, with none silently dropped.
- **SC-002**: Every current published document/asset in the migration inventory has a corresponding destination; all required home/README links are present.
- **SC-003**: All seven named products have sourced comparison entries and every release entry point enforces the same review contract.
- **SC-004**: Appearance selection, keyboard use, static search and representative documentation journeys are usable on desktop and mobile layouts.
- **SC-005**: The owner can review the complete local change and exact pending publication/domain actions before any public cutover.

## Assumptions

- The approved framework/hosting selection is an implementation constraint recorded in the plan: static Next.js, Fumadocs and GitHub Pages, with Cloudflare DNS.
- No application scheduler, GUI or CLI feature change is included; CLI gaps remain tracked in #277.
- Cloudflare login/MFA/token provisioning may require the owner. Authentication enables read-only preparation before the publication boundary; it does not authorize early public cutover.
- Existing release metadata contracts are retained or changed coherently with the README rewrite.
