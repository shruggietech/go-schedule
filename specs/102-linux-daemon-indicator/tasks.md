# Tasks: Linux desktop daemon presence and controls

**Input**: [spec.md](spec.md), [plan.md](plan.md), [research.md](research.md), [data-model.md](data-model.md), [desktop contract](contracts/desktop-service.md)

## Phase 1: Setup

- [x] T001 Record the Linux-only StatusNotifier dependency and license rationale in `go.mod`, `go.sum`, and `specs/102-linux-daemon-indicator/research.md`.
- [x] T002 [P] Add a package and build target for `cmd/gosched-indicator/` without changing the Windows tray binary.

## Phase 2: Foundation

- [x] T003 Add Linux installed-service state classification and tests in `internal/service/state_linux.go` and `internal/service/state_linux_test.go`.
- [x] T004 Add Linux service mutation/error-mapping tests in `internal/desktopcontrol/` while preserving the S101 Windows contract.

## Phase 3: User Story 1 - Truthful presence

**Goal**: A supported Linux session shows one accurate daemon item while the GUI may be closed.
**Independent test**: Stop/start the installed unit and remove local health; indicator and GUI never falsely show Running.

- [x] T005 [US1] Add status projection regression tests in `internal/desktopcontrol/control_test.go` for Linux missing, transition, and unhealthy service states.
- [x] T006 [US1] Implement indicator status, icon, and menu refresh in `cmd/gosched-indicator/` using `gogpu/systray` and the shared monitor.
- [x] T007 [US1] Expose the shared monitor in `desktop/main.go` and update platform-neutral frontend service wording in `desktop/frontend/src/`.

## Phase 4: User Story 2 - Safe local controls

**Goal**: Indicator and GUI can control the local installed service with confirmation and truthful outcomes.
**Independent test**: Start, Stop, and Restart from either surface, then deny or cancel graphical authorization.

- [x] T008 [US2] Extend action tests for fixed command arguments, authorization denial, cancellation, timeout, and final observations in `internal/desktopcontrol/`.
- [x] T009 [US2] Implement the narrow Linux graphical authorization path in `internal/desktopcontrol/`.
- [x] T010 [US2] Connect indicator Open, Start, Stop, Restart, and Quit handlers in `cmd/gosched-indicator/`.
- [x] T011 [US2] Extend Connections-page tests and UI messages in `desktop/frontend/src/settings/` and `desktop/frontend/src/App.tsx`.

## Phase 5: User Story 3 - Clean session lifecycle

**Goal**: One indicator and GUI per session, no accidental restart of a deliberately stopped service.
**Independent test**: Repeat Open and process launch, close GUI, restart host, and log out.

- [x] T012 [US3] Add and test Linux session singleton and GUI activation in `internal/desktopcontrol/activation_linux.go`.
- [x] T013 [US3] Correct GUI autospawn selection and tests in `desktop/main.go` and `desktop/main_service_test.go`.
- [x] T014 [US3] Add the XDG autostart contract, Linux archive build, and automation checks in `brand/platform/linux/`, `.github/workflows/release.yml`, and `scripts/automation-check.sh`.

## Phase 6: Polish and delivery

- [x] T015 [P] Document install, use, limitations, and removal in `docs/INSTALL-linux.md` and compact mark use in `docs/brand.md`.
- [x] T016 Run `sh scripts/verify.sh all` and relevant Linux build tests, then resolve failures in the affected code.
- [x] T017 Update `CHANGELOG.md`, mark this specification Implemented with delivery evidence, and commit S102 on `codex/102-linux-daemon-indicator`.
- [x] T018 Resolve the effective installed systemd `--config` path for both Linux desktop processes and test custom IPC path handling after review.
- [x] T019 Accept both daemon-supported `--config` forms and propagate session cancellation to a pending Linux authorization helper after second-round review.

## Dependencies and strategy

T001-T004 establish shared contracts. US1 can then establish truthful read-only status, US2 adds guarded mutation, and US3 closes lifecycle and packaging gaps. Documentation can proceed after behavior is stable. Tests precede or accompany each behavioral change. Publication and review are tracked on the pull request, not as incomplete implementation tasks.
