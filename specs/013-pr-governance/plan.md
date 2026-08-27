# Implementation Plan: Lightweight Pull-Request Workflow

**Branch**: `codex/013-pr-governance` | **Date**: 2026-08-26 | **Spec**:
`specs/013-pr-governance/spec.md`

## Summary

Replace stale direct-to-`main` guidance with a small branch-to-PR workflow whose
purpose is third-party AI review. Amend the constitution to v3.0.0 and
synchronize the README, contributor guide, autopilot protocol, agent context,
and pull request template. Do not add branch protection, repository settings, policy
scripts, required approvals, or new CI behavior.

## Technical Context

**Language/Version**: Markdown and existing repository conventions

**Dependencies**: Existing GitHub pull requests and CI only

**Testing**: `sh scripts/verify.sh all`, documentation link check, diff and
encoding audits

**Constraints**: No product, dependency, CI workflow, settings, CodeQL, release,
tag, or merge-method changes

## Constitution Check

| Gate | Result | Evidence |
| --- | --- | --- |
| Amendment authority | PASS | Issue #23 and the operator's explicit direction authorize replacing the stale integration rule. |
| Code and testing | PASS | No product code changes; the existing full aggregate remains the definition of green. |
| Autonomous execution | PASS | The single halt moves before review-branch push and PR creation. |
| Proportionality | PASS | The implementation is documentation-only and explicitly rejects branch protection as unjustified overhead for the current project. |
| Integration | PASS by amendment | Constitution v3.0.0 reintroduces PRs solely as a review venue. |

## Files

- `.specify/memory/constitution.md`
- `README.md`
- `CLAUDE.md`
- `CONTRIBUTING.md`
- `docs/build-autopilot.md`
- `.github/PULL_REQUEST_TEMPLATE.md`
- `CHANGELOG.md`
- `specs/013-pr-governance/`

## Explicit deviation from the earlier draft

The first Feature 013 draft proposed branch protection, eight hard-coded check
contexts, conversation resolution, administrator enforcement, and an offline
policy checker. The operator rejected that as disproportionate red tape for a
one-developer application with no users. Those artifacts are removed. The
remaining design keeps only the stated value: a PR where AI reviewers can
comment before the maintainer merges.
