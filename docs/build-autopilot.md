---
title: Build-phase autopilot
nav_order: 7
---

# Build-Phase Autopilot Protocol

Version: 2.0.0
Adopted: 2026-07-22
Status: operating procedure for the coding agent
Project: go-schedule

This document is the operating procedure for running a spec-kit feature under
autopilot. It is governed by the project constitution
(`.specify/memory/constitution.md`) and reaffirms, never weakens, its
principles. The constitution is project law; this document is the how. Where
they appear to conflict, the constitution wins.

## Purpose

Every feature runs the same spec-kit sequence. Feature scope traces to the
master specification (`specs/001-task-scheduler/spec.md`, whose authoritative
task list is `specs/001-task-scheduler/tasks.md`) or to an open issue on the
GitHub tracker, which is where the roadmap of remaining work lives; the
per-feature `specs/NNN-name/` tree is where a feature is actually specified. The default agent behavior pauses for
authorization between each step and raises routine decisions to the user that,
in practice, are approved as recommended.

Autopilot removes that friction: one verbal kickoff runs a full feature end to
end, the agent makes the routine decisions itself and records them, and the
agent halts once, right before the review branch and pull request leave local
control, with a breakdown for review.

## Trigger

The user starts an autopilot feature run with a verbal kickoff naming the
feature or the next feature, for example:

- "Kick off the catch-up feature"
- "Run the next feature"
- "Autopilot 005"

The operator may also place any other feature or task under autopilot with an
explicit request, for example:

- "Run the MSI signing work under autopilot"
- "Autopilot this"

On trigger, the agent runs the entire feature sequence below without pausing for
inter-step authorization.

## Preconditions

Confirm setup before running the sequence. Do not assume; read the project.

- Spec-kit is initialized (`.specify/` exists) and the project exposes the
  `/speckit-*` command skills under `.claude/skills/`. This project gitignores
  `.claude/`, so a fresh clone has no command skills until they are installed;
  restore them with `specify integration upgrade claude`. If they are absent,
  halt and say so. Never fake or reimplement a spec-kit command.
- This project uses the hyphenated command form (`/speckit-specify`,
  `/speckit-plan`). Use it consistently for the whole run.
- The constitution `.specify/memory/constitution.md` governs every decision
  below.

## Per-feature sequence

The agent runs these steps in order, with no halt between them:

1. `/speckit-specify` creates `specs/NNN-*/`, `spec.md`, and
   `checklists/requirements.md`, drawing scope from the relevant sections of the
   master specification and the issue it traces to.
2. `/speckit-clarify` runs under the decision policy below. The agent answers
   clarification questions itself from the feature spec, the constitution, the
   master specification, and the feature's stated scope and acceptance criteria.
   Only genuinely unanswerable questions are escalated.
3. `/speckit-checklist` adds domain checklists where the feature warrants them.
4. `/speckit-plan` produces `research.md`, `data-model.md`, `contracts/`, and
   `quickstart.md`.
5. `/speckit-tasks` produces `tasks.md`.
6. `/speckit-analyze` is the blocking gate. The agent resolves findings. A
   genuine CRITICAL conflict that needs a human decision triggers an early halt.
7. `/speckit-implement` executes the tasks under test-driven discipline
   (constitution principle II). Tests covering the safety-critical behaviors
   below are required, not optional.
8. Verify with CI parity (see the next section). A red result that cannot be
   fixed within the feature triggers a halt with the failure.
9. Commit locally as `feat(NNN): <title>` (NNN is the spec-kit feature number)
   with the agent's `Co-Authored-By:` attribution trailer, and update the
   `CHANGELOG.md` `[Unreleased]` section (an Added or Changed line, plus a dated
   Decisions entry for any architecture-affecting choice).
10. Halt before the work leaves the machine. Present the breakdown below and
    wait for explicit authorization.

## Verification: CI parity

The canonical aggregate mirrors `.github/workflows/ci.yml`. Run it in the
foreground and watch it finish:

```bash
sh scripts/verify.sh all
```

It runs all eight required gates, in order: `format`, `vet`, `lint`, `race`,
`gui`, `coverage`, `docs`, and `automation`. The format gate must print no files.
The race gate excludes the cgo-only GUI entry point and the Fyne widget package
(whose races are inside Fyne's own font cache, not this project's code); the
pure-Go `gui/viewmodel` package stays race-tested, and the GUI is covered by the
headless gate. Coverage on the core packages must stay at or above 80 percent
(constitution principle II). The automation gate independently guards approved
hosted-action majors and the exact verification manifest.

A missing POSIX shell, Go command, formatter, or C toolchain is a failed or
unrun prerequisite, never a green gate. Record the named gate and halt until it
can be run; do not omit it from the aggregate or substitute a later CI result.

Run all of these in the foreground and watch them to completion. Never launch
the test suite in the background and poll for its output. `go test` buffers a
package's output until that package finishes, so a background run cannot be
distinguished from a dead one, and treating one as the other has caused
misdiagnosed hangs elsewhere.

## Decision policy

This is the core behavioral change. For any decision point that the default
behavior would raise to the user, the agent instead:

- Enumerates the viable alternatives.
- Evaluates them against the constitution, the master specification, the
  feature's stated scope and acceptance criteria, and existing code patterns.
- Picks the best-supported option, proceeds, and records the decision and its
  rationale in the feature's `plan.md` or `spec.md`, and in `CHANGELOG.md`
  Decisions when the choice is architecture-affecting.

The agent halts to the user only when one of these holds:

- No option is clearly best and the choice is materially irreversible or
  architecture-defining.
- The feature's intent or scope is genuinely ambiguous in the master
  specification and in the issue it traces to.
- A constitution CRITICAL conflict cannot be resolved without a human decision.

## Branching and integration

Autopilot works on a review branch. Every change targets `main` through a pull
request, using the same path for maintainers, automation agents, and outside
contributors. Direct pushes to `main` are prohibited.

The halt precedes the review-branch push and pull-request creation. It also
marks the boundary before anything is published. Nothing is pushed or opened
without explicit authorization.

The authorization that publishes the PR remains valid for verified, in-scope
commits pushed to the same review branch to address CI failures or review
feedback. It expires when the PR is merged or closed and does not cover material
scope expansion, another PR, a tag, or a release. These review-fix pushes do not
introduce another autopilot halt.

Local CI parity must be green before the halt. After authorization, hosted CI
runs independently on the pull request and third-party AI reviewers can comment.
Consider each comment, implement warranted changes, and explain why other
suggestions do not fit. AI review is advisory and the maintainer decides when
the PR is ready to merge. This protocol adds no branch protection, approval
count, fixed hosted-check list, or mandatory conversation rule.

When the pull request fully completes an issue, its description uses a supported
closing keyword such as `Closes #N`. Partial or related work uses `Refs #N` and
leaves the issue open. After the maintainer merges, synchronize local `main` and
remove the merged local and remote review branch plus any stale worktree that
belongs to it.

## The pre-publication halt breakdown

At the single halt, the agent presents:

- The feature number and title, and what was built: the spec, plan, and tasks
  artifacts, the code packages, and the tests.
- The notable decisions made and why (the decision log).
- The verification results for format, vet, lint, race, GUI, coverage, docs, and
  automation, with evidence of pass, fail, or an honestly reported unavailable
  prerequisite.
- Any deviations or open risks against the feature's acceptance criteria.
- The exact review-branch push and pull-request creation commands awaiting
  authorization.

## Always-halt guardrails

These hold regardless of the decision policy:

- Never push a branch, open a pull request, tag a release, or run the release
  workflow without explicit authorization. The open PR's publication
  authorization is sufficient only for verified, in-scope review-fix pushes to
  its existing branch.
- Never weaken or skip the `/speckit-analyze` gate.
- Never weaken or skip the safety-critical test surfaces of this project:
  - Clock discipline: engine code takes time through the injected `Clock`
    interface, never `time.Now()` directly, and scheduling tests stay
    deterministic rather than leaning on real `time.Sleep`.
  - Timezone and DST correctness: per-task IANA timezone handling, including
    next-valid and first-occurrence resolution across DST transitions.
  - Store migrations: forward-only, non-destructive, and covered by a test that
    runs the migration against a database from the prior schema version.
  - Restart and catch-up recovery: persisted state survives a daemon restart and
    downtime catch-up runs a task once before resuming.
  - Concurrency: the race detector passes and every goroutine has a defined
    termination path (constitution principle I).
  - Local IPC access control: the Unix socket and Windows named pipe stay
    restricted to authorized callers.

Pinned process artifacts (`.github/workflows/**`, `build/**`, `Makefile`,
`.golangci.yml`, the `go`/`toolchain` lines of `go.mod`, `.gitattributes`,
`.gitignore`, `LICENSE`, `docs/INSTALL-windows.md`) may be modified when a
feature's scope requires it, provided the change is recorded as a dated decision
in the changelog. Autopilot does not halt separately for this. The changes
surface at the once-per-feature halt and must pass the CI merge gate before
merge. Cutting a release (a `vX.Y.Z` tag) still requires explicit
authorization.

## Scope and expiry

Autopilot is valid for features traceable to `specs/001-task-scheduler/spec.md`
or to an open issue on the GitHub tracker, which is where the roadmap of
remaining work lives. It also applies to any other feature or task when
the operator explicitly requests an autopilot run (for example "run the MSI
signing work under autopilot" or "autopilot this"). Such an explicit request
authorizes autopilot for the named work and is itself the renewal.

Absent an explicit request, work not traceable to the master specification or an
open issue falls back to normal interactive mode. When the master specification
is superseded by a new version, the standing authorization lapses and requires
renewal against the new document; per-request autopilot remains available
regardless.
