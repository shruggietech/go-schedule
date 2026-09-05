# Implementation Plan: v0.9 Closure and Maintenance Automation

**Branch**: `codex/029-v09-closure` | **Date**: 2026-08-30 | **Spec**: [spec.md](spec.md)

## Summary

Complete the locally actionable v0.9 maintenance work in one reviewable slice. Establish and enforce trustworthy specification lifecycle metadata, add low-noise Dependabot proposals, activate the useful secret-scanning controls named by #39, and move the unverified Windows icon report to the future backlog without implementing it. Keep #40 open only for its inherently post-merge first-proposal observation.

## Technical Context

**Language/Version**: POSIX shell, YAML, Markdown, GitHub repository APIs **Primary Dependencies**: Existing shell utilities, GitHub Actions, Dependabot, and Spec-Kit; no application dependency additions **Storage**: Repository files and hosted GitHub issue/security metadata **Testing**: Shell fixture tests, offline automation contract, hosted-state readback, and eight canonical gates **Target Platform**: Cross-platform repository checkout; hosted settings on `shruggietech/go-schedule` **Performance Goals**: Lifecycle and automation checks complete within the existing automation gate without network access **Constraints**: No app behavior or dependency-version changes, no branch protection, no Windows observation, UTF-8 without BOM **Scale/Scope**: 29 specifications, two dependency ecosystems, three requested security controls, and issue #33 triage

## Constitution Check

| Principle | Gate | Design evidence |
| --- | --- | --- |
| I. Code Quality | PASS | One documented lifecycle model replaces ad hoc status strings; validators stay small and explicit. |
| II. Testing Standards | PASS | Negative fixtures precede validator and Dependabot contract implementation; the full gate remains mandatory. |
| III. UX Consistency | PASS | Maintainer-facing inventory and failures use consistent state names and actionable paths. |
| IV. Performance | PASS | Checks are bounded local scans and add no runtime application work. |
| V. Autonomous Execution | PASS | S029 follows full Spec-Kit, blocking analysis, hosted-state verification, all gates, local commit, and one pre-push halt. |

### Post-design re-check

The design changes process automation but no application or release behavior. The automation gate remains offline, the hosted mutations are narrow and reversible, and every process-artifact decision is recorded in the changelog.

## Architecture and Decision Log

### Separate lifecycle maturity from publication

Use `Draft`, `Ready`, `In Progress`, `Implemented`, `Deferred`, `Superseded`, and `Abandoned`. `Implemented` means required implementation and verification work is resolved in the repository. A separate delivery field cites the commit, release, or pull request, avoiding compound strings such as "Local Complete - Publication Pending" that become stale immediately after merge.

### Reconcile historical tasks without rewriting decisions

For checkbox task lists whose complete delivery is proven by a merge or release, reconcile stale markers and add an audit note. Explicitly record a waiver when a manual observation was never performed. Preserve the legacy 007 task notation and identify it as resolved historical format rather than mechanically rewriting the whole file.

### Put lifecycle enforcement inside the existing automation gate

Add one offline lifecycle checker and fixture tests, then call it from `automation-check.sh`. This reuses the canonical eight-gate contract instead of creating a ninth gate or hosted workflow.

### Keep dependency proposals low noise

Run Go module checks weekly and GitHub Actions checks monthly. Group compatible minor and patch routine updates, cap open proposals, and leave major and security updates separate. This follows the sibling project's proven approach while matching this Go repository's smaller dependency surface.

### Activate only requested hosted controls

Enable push protection, non-provider patterns, and validity checks individually and verify their final API state. Keep AI detection and delegated dismissal unchanged. No workflow substitutes or governance expansion are justified.

### Defer evidence-blocked Windows work proportionally

Move #33 to `Post-v1`, retain `needs: verification`, and lower it from P1 to P3. The report remains valid backlog, but an unverified issue should not block v0.9 or interrupt the maintainer for a manual install session.

## Project Structure

```text
.github/dependabot.yml
.specify/templates/spec-template.md
specs/README.md
specs/001-*/spec.md ... specs/029-*/spec.md
specs/003-*/tasks.md ... selected historical tasks.md
specs/029-v09-closure/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── verification.md
├── contracts/maintenance-contract.md
├── checklists/
└── tasks.md
scripts/spec-lifecycle-check.sh
scripts/automation-check.sh
test/scripts/spec-lifecycle-check_test.sh
test/scripts/automation-check_test.sh
docs/build-autopilot.md
CHANGELOG.md
CLAUDE.md
```

**Structure Decision**: Extend existing documentation and shell-policy surfaces. No new service, language runtime, workflow, or application package is introduced.

## Complexity Tracking

No constitution violations require justification.
