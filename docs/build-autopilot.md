# Build-Phase Autopilot Protocol

Version: 1.0.0
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
master specification (`specs/001-task-scheduler/spec.md`) and the ordered
roadmap (`TODO.md`, whose authoritative task list is
`specs/001-task-scheduler/tasks.md`); the per-feature `specs/NNN-name/` tree is
where a feature is actually specified. The default agent behavior pauses for
authorization between each step and raises routine decisions to the user that,
in practice, are approved as recommended.

Autopilot removes that friction: one verbal kickoff runs a full feature end to
end, the agent makes the routine decisions itself and records them, and the
agent halts once, right before the work leaves the machine, with a breakdown for
review.

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
   master specification and the roadmap.
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

These mirror `.github/workflows/ci.yml`. Run every one of them and watch it
finish:

```bash
gofmt -l internal cmd test
```

```bash
go vet ./...
```

```bash
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6 run ./...
```

```bash
CGO_ENABLED=1 go test -race $(go list ./... | grep -vE '/cmd/gosched-gui|/gui$')
```

```bash
go test ./gui/...
```

`gofmt` must print nothing. The race run excludes the cgo-only GUI entry point
and the Fyne widget package (whose races are inside Fyne's own font cache, not
this project's code); the pure-Go `gui/viewmodel` package stays race-tested, and
the GUI is covered by the headless run above it. Coverage on the core packages
must stay at or above 80 percent (constitution principle II); the CI coverage
job is the authority and a feature that lowers a core package below the gate is
a halt.

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
  specification and the roadmap.
- A constitution CRITICAL conflict cannot be resolved without a human decision.

## Branching and integration

This project integrates through pull requests: `main` is protected by
convention and the constitution forbids direct pushes to it. Autopilot commits
the feature onto a working branch named for the feature (for example
`005-catchup-hardening`), never onto `main`.

The single halt therefore precedes both the branch push and the pull request.
Nothing is pushed and no pull request is opened without explicit authorization.
Once authorized, the branch is pushed and the pull request opened against
`main`, and CI is the merge gate.

If the operator explicitly directs the work onto `main` instead, that
instruction governs for that run, and the halt still precedes the push.

## The pre-push halt breakdown

At the single halt, the agent presents:

- The feature number and title, and what was built: the spec, plan, and tasks
  artifacts, the code packages, and the tests.
- The notable decisions made and why (the decision log).
- The verification results for gofmt, vet, lint, race tests, and GUI tests, with
  evidence of pass or fail.
- Any deviations or open risks against the feature's acceptance criteria.
- The exact `git push` command (and pull-request intent) awaiting authorization.

## Always-halt guardrails

These hold regardless of the decision policy:

- Never `git push`, open a pull request, tag a release, or run the release
  workflow without explicit authorization.
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
and the roadmap in `TODO.md`. It also applies to any other feature or task when
the operator explicitly requests an autopilot run (for example "run the MSI
signing work under autopilot" or "autopilot this"). Such an explicit request
authorizes autopilot for the named work and is itself the renewal.

Absent an explicit request, work not traceable to the master specification or
the roadmap falls back to normal interactive mode. When the master specification
is superseded by a new version, the standing authorization lapses and requires
renewal against the new document; per-request autopilot remains available
regardless.
