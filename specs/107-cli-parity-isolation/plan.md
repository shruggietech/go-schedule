# Implementation Plan: CLI parity and isolated operation

**Branch**: `codex/107-cli-parity-isolation` | **Date**: 2026-10-19 | **Spec**: [spec.md](spec.md)

## Summary

Deliver all functional outcomes in #283 through #288 under #277. Reuse current API/scheduler behavior, replace CLI globals with invocation-scoped dependencies, extract headless preference persistence and add owned fixtures. The comprehensive coordinator/report product remains #289/#290.

## Technical Context

**Language/Version**: Current root and desktop Go modules/toolchain, unchanged frontend shapes.

**Primary Dependencies**: Existing Cobra, API clients/types, clientprofile/clientsecret, clock, platform and x/sys. No new external dependency.

**Storage**: Current SQLite and preference/profile schemas; atomic persistence and cross-process scoped updates.

**Testing**: Request/transport, invocation, persistence, identity/ownership, adapter and packaged-fixture regressions, then canonical `sh scripts/verify.sh all` through the hidden foreground launcher.

**Target Platform**: Native Windows and supported Unix targets.

**Project Type**: Thin CLI/API clients, shared headless persistence and focused qualification helpers.

**Performance Goals**: At most four observation workers, per-target deadlines, bounded results and calendar intervals of at most 30 days. Scheduler hot paths unchanged.

**Constraints**: No default-target fallback, secret output or host mutation by fixtures. Preserve existing output/persistence and local-only runtime metadata. Hidden noninteractive children and exact owned cleanup.

**Scale/Scope**: Six issues, five user stories, one review branch. No full evaluator/report or release.

## Constitution Check

Pre-research and post-design: PASS. I: focused shared components and cancelable workers. II: behavioral/regression/identity tests, full race/coverage gates. III: compatible text/JSON, output channels and actionable validation. IV: bounded observation, no scheduler hot-path change. V: installed Spec-Kit sequence, analyze before implementation and local commit before publication halt.

No exception or weakened gate. Shared preference extraction and invocation architecture receive dated changelog decisions. Any necessary pinned-script edit is recorded explicitly.

## Project Structure

Feature artifacts: spec.md, plan.md, research.md, data-model.md, contracts/cli.md, contracts/isolation.md, checklists/, quickstart.md, tasks.md and verification.md in this directory.

Source: internal/cli (invocations/commands), internal/api/server and cmd/goschedd (clear intent/runtime identity), internal/userprefs (shared schema/store), internal/clientobserve (bounded observations), internal/clifixture (owned process/state), desktop/settings and desktop/taskgroup (shared behavior), test/integration and docs/cli.md.

**Structure Decision**: Root CLI never imports Wails. Attach invocation state to Cobra context, preserve helper constructors and pass its switchable client to captured-client constructors. Replace globals and direct process I/O with scoped dependencies.

## Decisions and order

1. Use --ipc-path for explicit local selection, incompatible with remote flags. Avoid a second --config meaning that shadows service install. Validate paths and never fall back.
2. Invocation owns target/JSON/client/store/secret/service/launcher/clock/I/O dependencies. Production native secrets remain default; explicit harness injection never becomes an environment-variable downgrade.
3. Local preferences/profile commands skip daemon target resolution. Saved desktop startup choice remains distinct from invocation --profile.
4. Add ClearEnv to shared task-update intent, consistent with existing ClearWorkingDir/ClearRunAs. Desktop and CLI use it; omitted and nonempty replacement retain meaning.
5. Add task set/clear, disabled create, syntax intent, group edit and schedule preview using current request/evaluator semantics.
6. Strict bounded UTF-8 notification JSON input supports heterogeneous replacement and canonical export; reject shorthand conflicts and preserve inheritance.
7. Keep history/log defaults; new active/get/details/metadata modes are opt-in. Calendar validates timestamps and 30-day bounds; runtime is local-only.
8. Share transport-neutral observation plumbing with expected/observed identity, access verification, four workers, bounded API results and independent errors. Preserve successful partial results with nonzero aggregate status.
9. Extract preference v1 schema/defaults/migration/validation/atomic writes. Cross-process locking prevents CLI/GUI lost updates. Keep desktop aliases/frontend shapes; use clientprofile.SetActive for startup choice.
10. Owned fixtures assert installation identity, PID, IPC and effective paths before mutation. Add missing authoritative runtime IPC/PID. Windows job containment must precede child execution; Unix groups are owned, cleanup exact and artifacts retained.
11. A disposable regression demonstrated that the inherited group API persisted a rename before rejecting a cyclic parent. Fix #304 within the authoring slice using one shared-store transaction; preserve rename-only, parent-only and notification reset semantics.

12. Round-1 review found that context checks around synchronous credential access did not bound it (#306). Require LoadContext for observation adapters and run native reads in a hidden, cancellable CLI helper. Keep tokens exclusively in private redirected pipes; never detach blocking native calls into unjoinable goroutines. Selected-target preflight receives a request deadline too.
13. Share bounded Unix group termination/drain between fixtures and native-read helpers (#307). Wait for group absence after SIGKILL; preserve runtime and report cleanup failure at the deadline. Windows native reads use the exact helper process handle and do not spawn descendants.
14. Hosted macOS CI demonstrated oversized default socket paths (#308). Use a private UUID root beneath canonical /private/tmp on macOS and validate explicit Unix placement before creating state.

15. Round 2 found that an internal argument did not authenticate the credential helper (#309). Replace standard-stream requests/responses with a private operation-owned local channel, mutually checked kernel peer PIDs and the actual parent executable image. Reject ordinary direct invocation before native access. Use bounded frames, cancellation and an independent helper lifetime. Windows uses native pipe peer/image APIs; Linux uses SO_PEERCRED and /proc executable identity; macOS uses LOCAL_PEERPID and the kernel executable-vnode path (not mutable argv). No third Codex round is requested.

## Complexity Tracking

No constitution exception. These small shared components replace global/duplicated behavior; no second scheduler, general test framework or model-service integration.
