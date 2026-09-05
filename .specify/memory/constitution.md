<!--
SYNC IMPACT REPORT
==================
Version change: 2.0.1 → 3.0.0 (2026-08-26, reinstating pull requests as the
  integration path for third-party review).
Bump rationale: Requiring a review branch and pull request changes the
  integration contract for every author and is therefore a MAJOR amendment.
  The rule is intentionally lightweight for a one-developer project with no
  users: it creates a durable AI-review venue without branch protection,
  approval requirements, or other hosted enforcement.

Modified principles:
  - V. Autonomous Build-Phase Execution (2026-08-26, v3.0.0: the mandatory halt
    precedes review-branch publication and PR creation)
Added principles:
  - I. Code Quality
  - II. Testing Standards (NON-NEGOTIABLE)
  - III. User Experience Consistency
  - IV. Performance Requirements
  - V. Autonomous Build-Phase Execution (added 2026-07-22, v1.1.0)
Added sections:
  - Engineering Constraints
  - Development Workflow & Quality Gates
  - Governance
Modified sections (2026-08-26, v3.0.0):
  - Development Workflow & Quality Gates (lightweight PR-first review)
  - Governance (the maintainer retains final review and merge judgment)

Templates requiring updates:
  ✅ .specify/templates/plan-template.md (Constitution Check gates are dynamic;
     no edit required for either amendment)
  ✅ .specify/templates/spec-template.md (no mandatory-section changes required)
  ✅ .specify/templates/tasks-template.md (principle-driven task categories covered)
  ✅ README.md (contributor entry point synchronized for PR-first review)
  ✅ CLAUDE.md (Integration workflow synchronized for PR-first, v3.0.0)
  ✅ CONTRIBUTING.md (one integration workflow and closing-keyword guidance)
  ✅ docs/build-autopilot.md (review branch, halt, PR, review, and cleanup flow)
  ✅ .github/PULL_REQUEST_TEMPLATE.md (canonical verification and `Closes`)
  ✅ .github/workflows/ci.yml (no change needed; already runs on pull requests)

Deferred TODOs: None
-->

# go-schedule Constitution

## Core Principles

### I. Code Quality

Code MUST be correct, readable, and idiomatic before it is considered complete.

- All Go code MUST pass `gofmt`, `go vet`, and the project linter (`golangci-lint`) with zero warnings; CI MUST reject any change that does not.
- Public packages, exported types, and exported functions MUST carry doc comments that explain intent and contract, not restate the signature.
- Errors MUST be handled explicitly: wrap with context using `fmt.Errorf("...: %w", err)`, never silently discard with `_`, and never `panic` in library code for recoverable conditions.
- Functions MUST do one thing; cyclomatic complexity SHOULD stay low and any function the linter flags MUST be refactored or justified in review.
- Concurrency primitives (goroutines, channels, locks) MUST have a documented ownership and lifecycle; every goroutine MUST have a defined termination path and the race detector (`go test -race`) MUST pass.

**Rationale**: A scheduler is long-running, concurrent infrastructure. Defects in error handling or goroutine lifecycle leak resources and corrupt timing guarantees, so quality is enforced mechanically rather than left to discretion.

### II. Testing Standards (NON-NEGOTIABLE)

Tests are written alongside or before the code they verify, and the suite is the source of truth for correctness.

- Every behavioral change MUST ship with tests; bug fixes MUST include a regression test that fails before the fix and passes after.
- Unit tests MUST cover scheduling logic, time/clock handling, and error paths. Time MUST be injected through an abstraction (e.g., a `Clock` interface), tests MUST NOT depend on real wall-clock `time.Sleep` for deterministic assertions.
- Integration tests MUST cover job persistence, recovery after restart, and concurrent job execution.
- All tests MUST run under `go test -race`; flaky tests MUST be fixed or quarantined with a tracking issue, never ignored.
- Coverage on core scheduling packages MUST be ≥ 80%; new code MUST NOT lower package coverage. CI MUST enforce these gates.

**Rationale**: Correct timing and reliable recovery are the product's core promise. Without deterministic, race-checked tests these guarantees cannot be verified, so testing discipline is non-negotiable.

### III. User Experience Consistency

Every interface the project exposes, CLI, configuration, logs, and API, MUST behave predictably and uniformly.

- The CLI MUST follow a consistent verb-noun command structure, support both human-readable and `--json` output, write results to stdout and errors/diagnostics to stderr, and return conventional exit codes (0 success, non-zero failure).
- Configuration MUST have a single documented schema with sensible defaults; invalid configuration MUST fail fast at startup with a clear, actionable message naming the field.
- Time inputs and outputs MUST use a consistent format (RFC 3339) and timezone handling MUST be explicit; durations MUST use Go duration syntax consistently.
- Error messages MUST be actionable: state what failed, why, and what the user can do.
- Logging MUST be structured and consistent across components (consistent field names, levels, and correlation/job identifiers).

**Rationale**: Operators interact with a scheduler under time pressure during incidents. Consistent, self-explanatory interfaces reduce operator error and the cost of recovery.

### IV. Performance Requirements

The scheduler MUST be efficient and meet stated timing and resource budgets.

- Scheduling decisions MUST be measured: job dispatch latency (scheduled time → execution start) MUST stay within a documented budget (default target: p99 < 100ms under nominal load) and the budget MUST live next to the code it governs.
- Performance-sensitive changes MUST include benchmarks (`go test -bench`); changes MUST NOT regress an existing benchmark by more than 10% without explicit, recorded justification.
- The system MUST NOT leak goroutines or memory under sustained load; resource usage MUST be bounded and verified for the supported job-count target.
- Hot paths MUST avoid unnecessary allocations and unbounded data structures; algorithmic complexity of scheduling operations MUST be documented.
- Premature optimization is rejected: optimizations MUST be justified by a benchmark or profile, not by intuition.

**Rationale**: A scheduler that drifts, stalls, or leaks under load fails silently and erodes trust. Performance is therefore a measured, budgeted, and continuously verified property rather than an afterthought.

### V. Autonomous Build-Phase Execution

Build-phase work runs under the Build-Phase Autopilot Protocol (`docs/build-autopilot.md`): a single verbal kickoff authorizes the full spec-kit sequence (specify, clarify, checklist, plan, tasks, analyze, implement, verify, commit) to run end to end without per-step authorization.

- The standing authorization covers features traceable to the master specification (`specs/001-task-scheduler/spec.md`) or to an open issue on the GitHub tracker, which is where the roadmap of remaining work lives, plus any feature or task the operator explicitly places under autopilot by request. An explicit request authorizes autopilot for the named work and is itself the renewal; absent one, work outside that scope stays in normal interactive mode.
- The agent decides routine questions itself: it enumerates the alternatives, evaluates them against this constitution, the master specification, and the feature's scope and acceptance criteria, picks the best, proceeds, and records the rationale in the feature's `plan.md` or `spec.md` and in `CHANGELOG.md` when the choice is architecture-affecting. It halts to the user only when no option is clearly best on an irreversible or architecture-defining choice, when the feature's intent is genuinely ambiguous, or when a CRITICAL conflict needs a human decision.
- All development MUST use a review branch and a pull request targeting `main`.
- Exactly one halt per feature is mandatory: before the review branch is pushed or its pull request is opened. The halt includes a breakdown of notable decisions and what was built. Publishing, tagging, and cutting a release always require explicit authorization. The authorization that publishes a pull request also covers verified, in-scope commits pushed to that same branch to address CI or review feedback until the PR is merged or closed; it does not authorize material scope expansion, a different PR, a tag, or a release.
- Every feature MUST still be spec'd through spec-kit before implementation. The `/speckit-analyze` gate MUST NOT be skipped or weakened.

**Rationale**: In practice the routine decisions raised between spec-kit steps were approved as recommended, so the per-step pause cost time without adding review value. Concentrating review at a single pre-push halt keeps the human gate where it actually matters. This principle grants autonomy of execution only: it does not relax principles I through IV or the CI quality gates below.

## Engineering Constraints

- Language and tooling: Go (latest stable minor release), managed with Go modules. Avoid third-party dependencies where the standard library suffices; every new dependency MUST be justified in review and pass a license check.
- Supported platforms: the project MUST build and pass tests on Linux and Windows.
- Backward compatibility: persisted job state and configuration schemas MUST migrate forward; breaking changes to either require a MAJOR version bump and a documented migration path.
- Security: secrets MUST NOT be logged; inputs from configuration and APIs MUST be validated at the boundary.

## Development Workflow & Quality Gates

- Features run under the Build-Phase Autopilot Protocol (`docs/build-autopilot.md`, see principle V): one kickoff runs the spec-kit sequence end to end, the agent decides routine questions itself and records the rationale, and it halts once before the review branch is published or its pull request is opened.
- All development MUST use a review branch and a pull request targeting `main`. This applies equally to maintainers, automation agents, and outside contributors. Direct pushes to `main` are prohibited.
- CI MUST pass and MUST enforce: `gofmt`/`go vet`, linter, `go test -race`, coverage thresholds, and benchmark regression checks for performance-sensitive packages. The agent MUST run CI-parity locally before the halt, in the foreground and watched to completion. Existing hosted CI provides additional review evidence after the PR opens. A red run MUST be investigated rather than misrepresented as green.
- Third-party review comments MUST be considered individually. Warranted changes are made; suggestions that do not fit receive a concise rationale. AI feedback is advisory and the maintainer retains final review and merge judgment.
- The single pre-publication halt MUST verify compliance with all five core principles, and MUST surface any change that weakens one without recorded justification.
- Any deviation from a principle MUST be recorded in the pull-request description under a "Complexity / Deviation" note explaining why a simpler compliant approach was rejected.

## Governance

This constitution supersedes ad-hoc practices and conventions. When a technical decision conflicts with these principles, the principles win unless an explicit, recorded amendment changes them.

- **Authority**: Design documents, the pre-publication halt, and pull-request review MUST verify compliance with the five core principles and the constraints above. The pull request is the durable review record; the maintainer is the final integration authority.
- **Guiding decisions**: Technical and implementation choices (architecture, dependencies, data structures, interface design) MUST be evaluated against the principles. The default bias is the simplest design that satisfies all four; added complexity MUST be justified in writing against a named principle (typically Performance or Testing) and recorded in the pull-request description and, where it is architecture-affecting, in `CHANGELOG.md`.
- **Amendment procedure**: Amendments require (1) a written proposal describing the change and rationale, (2) operator approval, and (3) a synchronized update of dependent templates and guidance docs. The Sync Impact Report at the top of this file MUST be updated on every amendment.
- **Versioning policy**: This constitution is versioned with semantic versioning. MAJOR = backward-incompatible governance/principle removals or redefinitions; MINOR = a new principle/section or materially expanded guidance; PATCH = clarifications and non-semantic refinements.
- **Compliance review**: Compliance is checked at every pre-publication halt and pull request. Periodically (at minimum each release), maintainers MUST review whether the principles still reflect reality and amend rather than let practice silently drift.
- **Runtime guidance**: Use `CLAUDE.md`, `docs/build-autopilot.md`, and `.specify/` templates for day-to-day development guidance; those documents MUST stay consistent with this constitution. Where they appear to conflict with it, the constitution wins.

**Amendments**:

- 2026-07-22, v1.1.0: added principle V (Autonomous Build-Phase Execution), adopting the Build-Phase Autopilot Protocol at `docs/build-autopilot.md`. A single kickoff authorizes the full spec-kit sequence with one mandatory halt before anything is pushed. Autonomy of execution only: principles I through IV, the CI quality gates, and the pull-request integration requirement are unchanged. Mirrored in `CLAUDE.md`.
- 2026-07-22, v2.0.0: **removed the pull-request integration requirement.** Development is trunk-based: work is committed directly onto `main`, with no feature branches and no pull requests. The requirement never described this project's practice, it is a one-to-two developer project that has never used pull requests for review, and a PR with no reviewer is ceremony that adds latency without adding scrutiny. Removing a mandated gate is a backward-incompatible governance change, hence MAJOR. Nothing else is relaxed: the single pre-push halt remains mandatory and becomes the sole human review point, CI still runs on every push to `main`, and the local CI-parity requirement is *strengthened*, with no PR to block a bad merge, a red local run must halt rather than be pushed and sorted out afterwards. Mirrored in `CLAUDE.md` and `docs/build-autopilot.md`.
- 2026-07-23, v2.0.1: **retired `TODO.md`; the roadmap is now the GitHub issue tracker.** Principle V's standing authorization previously named `TODO.md` as the second source of traceable scope alongside the master specification. That file duplicated, in prose a reader had to be told to look at, work that belongs in the tracker where it can be labelled, discussed, and closed; its eight remaining open items were filed as issues
  #13 through #20 and the file was removed. This is a PATCH clarification: what
  autopilot may run without further authorization is unchanged, only where that scope is recorded. Mirrored in `CLAUDE.md` and `docs/build-autopilot.md`.
- 2026-08-26, v3.0.0: **reinstated pull requests as the mandatory integration path.** Recent maintainer work established that a review branch and PR are useful venues for third-party AI review. Autopilot now halts before review-branch publication and PR creation. That publication authorization continues through verified, in-scope review-fix pushes to the same open PR, avoiding contradictory extra halts during review. The amendment deliberately adds no branch protection, approval requirement, fixed check list, or mandatory conversation rule; those controls are disproportionate for the current one-developer project with no users. The maintainer considers feedback and retains final merge judgment. Mirrored in `README.md`, `CLAUDE.md`, `CONTRIBUTING.md`, `docs/build-autopilot.md`, and the PR template.

**Version**: 3.0.0 | **Ratified**: 2026-06-19 | **Last Amended**: 2026-08-26
