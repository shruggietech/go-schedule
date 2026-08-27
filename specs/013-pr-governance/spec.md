# Feature Specification: Lightweight Pull-Request Workflow

**Feature Branch**: `codex/013-pr-governance`

**Created**: 2026-08-26

**Status**: Local Complete - Publication Pending

**Traceability**: Closes GitHub issue
[#23](https://github.com/shruggietech/go-schedule/issues/23).

## Overview

Maintainer work has increasingly used pull requests so third-party AI reviewers
can inspect it, while the written project policy still says to push directly to
`main`. This feature makes the useful practice official without imposing
enterprise controls on a one-developer project with no users.

The workflow is deliberately small: use a review branch, verify locally, halt
before publication, open a pull request, consider reviewer feedback, and leave
the final merge to the maintainer.

### Scope in

- One pull-request workflow for maintainer, agent, and outside work.
- Synchronized constitution, contributor, autopilot, agent, and PR-template
  guidance.
- Correct issue-closing keywords when a pull request completes an issue.
- A lightweight third-party review record.

### Scope out

- Branch protection, rulesets, required approvals, mandatory conversations,
  administrator enforcement, or a fixed list of required checks.
- Product code, dependencies, CI behavior, releases, tags, or merge methods.
- Requiring every reviewer suggestion to be accepted.

## User Story - Obtain useful third-party review (Priority: P1)

As the sole maintainer, I want changes proposed through pull requests so AI
reviewers can inspect them before I decide whether to merge.

**Independent Test**: Follow the maintained guidance from synchronized `main`.
It creates a review branch, runs local verification, halts before publication,
opens a PR with `Closes #N` when appropriate, considers review feedback, and
does not require repository protection settings.

### Acceptance Scenarios

1. **Given** locally complete work, **when** the publication halt is authorized,
   **then** the review branch is pushed and a pull request targets `main`.
2. **Given** AI review comments, **when** they are assessed, **then** warranted
   changes are made and unwarranted suggestions receive a reasoned response.
3. **Given** the pull request fully completes an issue, **when** its description
   is written, **then** it uses `Closes`, `Fixes`, or `Resolves` rather than
   `Refs`.
4. **Given** the maintainer is satisfied with the review and available CI
   evidence, **when** they choose to merge, **then** no artificial approval or
   branch-protection prerequisite is imposed by this feature.

## Requirements

- **FR-001**: Maintainer and agent work MUST use a review branch and pull request
  targeting `main` so third-party review has a durable place to occur.
- **FR-002**: Autopilot MUST halt once before pushing the review branch and
  opening the pull request.
- **FR-003**: The canonical local aggregate MUST pass before the halt.
- **FR-004**: Maintained guidance MUST describe the same lightweight workflow
  and MUST NOT instruct maintainers to push directly to `main`.
- **FR-005**: Review comments MUST be considered individually; warranted changes
  are implemented and other suggestions receive a concise rationale.
- **FR-006**: Issue-closing keywords MUST be used only when the PR completes the
  issue; partial work uses `Refs`.
- **FR-007**: The feature MUST NOT add branch protection, required approvals,
  required conversations, fixed hosted-check requirements, or settings changes.
- **FR-008**: The maintainer retains final review and merge control.

## Success Criteria

- **SC-001**: Constitution, contributor, autopilot, agent, and PR-template
  guidance all describe the same branch-to-PR review path.
- **SC-002**: No active maintained guidance instructs direct maintainer pushes
  to `main`.
- **SC-003**: The pull-request template uses the canonical aggregate and includes
  `Closes #`.
- **SC-004**: Full local verification passes with no product or CI workflow diff.
- **SC-005**: No repository setting is changed by this feature.

## Assumptions

- Pull requests are valuable here primarily as a venue for AI review and a
  durable discussion record.
- The sole maintainer can judge reviewer feedback and available CI results
  without technical merge enforcement.
- If the project gains users or maintainers later, stronger controls can be
  proposed as a separate, evidence-driven change.
