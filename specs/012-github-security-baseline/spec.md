# Feature Specification: GitHub Security Baseline

**Feature Branch**: `012-github-security-baseline` (trunk-based: committed
directly to local `main`)

**Created**: 2026-08-26

**Status**: Local Complete - Hosted Activation Pending

**Input**: User description: "Combine GitHub issues #38 and #39 into the next
focused Spec-Kit work slice and run it end-to-end under the build-phase
autopilot protocol: restore the advertised private vulnerability-reporting
route, enable supported secret protections, and add repository-controlled
CodeQL analysis without changing product behavior."

**Traceability**: Closes GitHub issues
[#38](https://github.com/shruggietech/go-schedule/issues/38) and
[#39](https://github.com/shruggietech/go-schedule/issues/39).

## Overview

The repository currently advertises a private vulnerability-reporting route
that GitHub reports as disabled. A security researcher following the documented
link therefore reaches a dead end and may fall back to an unsafe public issue.
The policy also promises no explicit secondary channel even though the issue
that discovered the defect assumed one existed.

The same repository has no code-scanning analysis and reports every available
secret-scanning control as disabled. Release automation, system-service
packaging, and a daemon that runs configured commands make early detection of
credential leaks and common code vulnerabilities materially valuable.

This feature establishes one repository security baseline: researchers receive
a working private route with a truthful fallback and accountable triage, while
maintainers receive repeatable secret and code analysis through supported
GitHub-native controls. Unsupported account-plan capabilities remain visible as
limitations rather than being silently counted as passing.

### Scope in

- GitHub private vulnerability reporting for this repository.
- The vulnerability-reporting instructions, fallback channel, ownership, and
  acknowledgement expectation in `SECURITY.md`.
- Repository-level secret scanning, non-provider pattern scanning, validity
  checks, and push protection where GitHub makes them available.
- Go CodeQL analysis on pull requests, the default branch, and a maintenance
  schedule.
- Offline regression coverage for the new workflow action references and
  security-workflow contract.
- Dated changelog decisions for changed pinned automation.
- Hosted evidence after explicit authorization to publish and activate the
  repository settings.

### Scope out

- Product, scheduler, daemon, CLI, GUI, store, IPC, or packaging behavior.
- Branch protection or changing the integration constitution (#23).
- Dependabot or dependency-version updates (#40).
- Spec-Kit lifecycle repair (#42).
- Resolving or dismissing any security alert discovered by the new controls.
- Creating a fake vulnerability report or publishing vulnerability details as
  a validation technique.
- Tagging, releasing, or changing organization-wide security policy.

## Clarifications

### Session 2026-08-26

- Q: What fallback should be documented when private reporting is unavailable?
  → A: The organization’s published contact, `info@shruggie.tech`.
- Q: Who owns initial private-report triage? → A: Repository administrators,
  with the existing one-week acknowledgement commitment.
- Q: What happens when a requested GitHub protection is unavailable to the
  repository plan? → A: Enable every available control and record the exact
  unavailable status plus a compensating control; never report it as enabled.
- Q: Must validation submit a sample advisory? → A: No. Validate the private
  form and maintainer access without creating fabricated vulnerability data.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Researchers can report privately (Priority: P1)

As a security researcher, I need the repository’s advertised vulnerability
route to open a usable private report form, with a truthful fallback if GitHub
is unavailable, so I do not have to disclose a vulnerability publicly.

**Why this priority**: The current P0 defect breaks the repository’s primary
security-reporting promise at the moment a reporter needs it.

**Independent Test**: Inspect the repository setting and open the advertised
route in authenticated and unauthenticated contexts. The story is complete when
GitHub presents a private-report path, repository administrators can access the
private triage surface, and `SECURITY.md` provides an independently verified
fallback and response expectation.

**Acceptance Scenarios**:

1. **Given** a researcher follows the repository’s security contact link,
   **when** GitHub loads the advisory route, **then** the researcher is offered
   a private reporting flow rather than a public issue form.
2. **Given** an unauthenticated visitor follows the same link, **when** GitHub
   requires authentication, **then** the visitor can authenticate and continue
   to the private report flow without losing the destination.
3. **Given** GitHub private reporting is temporarily unavailable, **when** a
   reporter reads `SECURITY.md`, **then** a published organization contact is
   available and the policy still says not to open a public issue.
4. **Given** a private report arrives, **when** repository administrators inspect
   the security triage surface, **then** ownership and the acknowledgement
   expectation are explicit.

---

### User Story 2 - Maintainers receive automated security findings (Priority: P2)

As a maintainer, I need supported secret protections and repeatable Go code
analysis to run as part of normal repository activity, so credential leaks and
common code vulnerabilities become visible before release.

**Why this priority**: These controls close a meaningful prevention and
detection gap, but the private reporting route is already publicly broken and
therefore comes first.

**Independent Test**: Inspect the repository security settings, validate the
security workflow offline, and observe hosted analysis on an ordinary change
and on the default branch. The story is complete when every available requested
secret control is enabled, CodeQL reports a successful Go analysis, and any
unavailable control has a precise limitation record.

**Acceptance Scenarios**:

1. **Given** the repository security settings, **when** their status is queried,
   **then** secret scanning and push protection are enabled when supported and
   each additional requested control has an explicit enabled or unavailable
   result.
2. **Given** a pull request targeting the default branch, **when** security
   analysis runs, **then** Go analysis completes and publishes its result as a
   named check without weakening the existing validation workflow.
3. **Given** a push to the default branch, **when** the security workflow runs,
   **then** the repository’s Go source is analyzed and results appear in the
   security tooling surface.
4. **Given** no code changes occur for a week, **when** the maintenance schedule
   fires, **then** analysis reruns against current security queries.
5. **Given** an obsolete or unknown security action reference or a weakened
   workflow contract, **when** offline automation validation runs, **then** the
   repository check fails with an actionable diagnostic.

### Edge Cases

- The organization plan may support provider-pattern secret scanning but not
  non-provider patterns or validity checks. Each control needs its own observed
  result; a partial enablement cannot be summarized as fully green.
- The authenticated operator token may administer the repository but lack a
  scope needed to read alerts. Enablement and verification failures must be
  distinguished from missing token scope.
- Enabling secret scanning may immediately reveal existing alerts. This feature
  reports their count and leaves triage to follow-up work; it does not dismiss,
  expose, or silently remediate them.
- A CodeQL workflow introduced by a pull request can analyze that pull request,
  but its scheduled trigger cannot be observed until the workflow exists on the
  default branch. Static validation and hosted evidence must be reported
  separately.
- Code analysis for a compiled language can fail when an automatic build reaches
  desktop dependencies. The workflow must define a repeatable analysis build
  compatible with the repository’s existing headless CI environment.
- Pull requests from forks receive reduced token permissions. The workflow must
  retain least privilege and must not expose secrets to untrusted code.
- The existing offline action policy rejects newly introduced action families
  until they have been researched, approved, and covered by negative fixtures.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Private vulnerability reporting MUST be enabled for the repository
  before the feature is declared active.
- **FR-002**: The advertised advisory link MUST lead an authenticated reporter to
  a private reporting form and MUST NOT route vulnerability reports into public
  issues.
- **FR-003**: `SECURITY.md` MUST name `info@shruggie.tech` as the fallback channel
  and retain the instruction not to disclose vulnerabilities publicly.
- **FR-004**: `SECURITY.md` MUST name repository administrators as initial triage
  owners and retain a measurable acknowledgement expectation of one week.
- **FR-005**: Secret scanning and push protection MUST be enabled when supported
  by the repository’s GitHub plan.
- **FR-006**: Non-provider pattern scanning and secret validity checks MUST each
  be enabled when supported and otherwise recorded individually with the exact
  unavailable status and compensating control.
- **FR-007**: Go code analysis MUST run for pull requests targeting the default
  branch, pushes to that branch, and at least weekly.
- **FR-008**: The security analysis workflow MUST use least-privilege permissions
  sufficient to read source and publish security results, without access to
  repository or organization secrets.
- **FR-009**: Every external action used by the security workflow MUST target a
  researched, supported hosted runtime and MUST pass the repository’s offline
  action policy.
- **FR-010**: The security workflow MUST use the module-selected Go version and a
  repeatable headless build that does not depend on desktop graphics libraries.
- **FR-011**: Offline regression coverage MUST reject obsolete CodeQL action
  majors, unknown security actions, missing required triggers, insufficient
  permissions, and omission of required analysis steps.
- **FR-012**: Existing CI triggers, jobs, permissions, verification gates, and
  product-build behavior MUST remain unchanged.
- **FR-013**: Every changed pinned workflow or automation policy MUST receive a
  dated unreleased changelog decision explaining the security control and its
  compatibility constraints.
- **FR-014**: Hosted activation evidence MUST distinguish enabled, unavailable,
  unverified, and failed controls and MUST report discovered alert counts without
  exposing alert contents.
- **FR-015**: No local validation step MUST create an advisory, push a branch,
  publish a release, dismiss an alert, or mutate repository security settings
  before the autopilot pre-push authorization.
- **FR-016**: The feature MUST complete without changes to product behavior,
  dependency versions, branch protection, or the integration constitution.

### Key Entities

- **Reporting Route**: The private advisory destination, authentication behavior,
  fallback channel, triage owner, and acknowledgement expectation.
- **Security Control**: A repository protection with a requested state, observed
  state, availability, evidence source, and compensating control when unavailable.
- **Security Analysis Contract**: The events, cadence, permissions, language,
  build behavior, and result publication required of hosted code analysis.
- **Security Evidence Record**: A non-sensitive result that identifies a control
  as enabled, unavailable, unverified, or failed without exposing alert details.
- **Pinned-Artifact Decision**: A dated record naming the automation artifact,
  the security change, the rationale, and compatibility constraints.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The private-vulnerability-reporting API reports enabled, and 100%
  of repository security links direct reporters to the private advisory route.
- **SC-002**: Both authenticated and unauthenticated route checks reach a usable
  authentication/private-report sequence without exposing a public issue form.
- **SC-003**: The security policy identifies one primary private route, one
  published fallback, one triage owner group, and a one-week acknowledgement
  target.
- **SC-004**: 100% of requested secret controls have an evidence-backed status;
  every supported control is enabled and no unavailable control is reported as
  passing.
- **SC-005**: Code analysis successfully analyzes Go on a pull request and on the
  default branch, with a weekly trigger present for future query updates.
- **SC-006**: Offline negative fixtures detect all five required drift classes:
  obsolete action major, unknown action, missing trigger, insufficient
  permission, and missing analysis step.
- **SC-007**: The existing eight-gate aggregate verification contract remains
  green and unchanged after the security workflow is added.
- **SC-008**: Every changed pinned automation artifact has a dated changelog
  decision before commit.
- **SC-009**: The local feature commit contains no product-code, dependency,
  branch-protection, release, advisory, alert-triage, or remote-setting mutation.

## Assumptions

- The repository remains public and GitHub Actions remains enabled.
- Repository administrators are the GitHub-defined recipients of private
  vulnerability reports and can access the advisories triage surface.
- `info@shruggie.tech` is the organization’s published fallback contact; this
  feature verifies publication but does not send a synthetic vulnerability by
  email.
- Advanced repository-owned CodeQL configuration is preferred over opaque
  default setup because this issue explicitly requires reviewable workflow
  configuration, stable triggers, and integration with the offline action
  policy.
- Weekly scheduled analysis is justified because query updates can detect newly
  understood vulnerabilities even when source code is unchanged.
- Remote settings are activated only after the local implementation passes the
  autopilot halt and the operator explicitly authorizes publication.
