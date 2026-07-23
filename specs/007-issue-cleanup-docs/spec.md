# Feature Specification: Issue cleanup, README refresh, and documentation completion

**Feature Branch**: `007-issue-cleanup-docs`

**Created**: 2026-07-23

**Status**: Draft

**Input**: User description: "Close issues #5 #6 #7 (MSI PATH, least-privilege
Windows service status, issue templates), refresh README to house style, and
complete the repo documentation set"

## User Scenarios & Testing *(mandatory)*

Three open issues share one shape: none is about the scheduling engine. Each
sits on the seam between the shipped product and the people trying to use it or
report against it. The documentation work in this feature is the same seam seen
from the other side — what a reader is told the product does, versus what it
does.

### User Story 1 - The documented commands work after an MSI install (Priority: P1)

Someone installs go-schedule on Windows from the `.msi`, opens a fresh
PowerShell, and follows the README. Today every command fails at the first
word, because the installer puts the binaries in *Program Files* and never puts
that directory on `PATH`. The reader concludes the install is broken. It is
not; the documentation is describing a shell the reader does not have.

**Why this priority**: This breaks the first five minutes of every Windows
install. Someone who cannot run `gosched health` has no reason to believe
anything else in the documentation either, and the failure mode — an
unrecognized command — reads as a broken package rather than a missing `PATH`
entry.

**Independent Test**: Install the MSI on a machine that has never had the repo
checked out, open a new shell, and run `gosched --version`. It resolves. Then
uninstall and confirm neither `PATH` scope retains a `go-schedule` fragment.

**Acceptance Scenarios**:

1. **Given** a machine with no prior go-schedule install, **When** the `.msi` is
   installed and a **new** shell is opened, **Then** `gosched --version`
   resolves without a path prefix.
2. **Given** an installed go-schedule, **When** it is uninstalled via *Apps &
   features*, **Then** no `go-schedule` entry remains in either the machine or
   the user `PATH`.
3. **Given** an installed go-schedule, **When** a newer `.msi` performs a major
   upgrade, **Then** exactly one `go-schedule` entry is present in `PATH`, not
   two.
4. **Given** the fixed installer, **When** a reader follows `README.md` or
   `docs/INSTALL-windows.md` verbatim, **Then** every command in those files
   runs as written.

---

### User Story 2 - An ordinary user can ask whether the service is running (Priority: P1)

A standard, non-elevated user wants to know whether the scheduler is running.
The installed service's ACL explicitly grants Interactive Users
`SERVICE_QUERY_STATUS`, so the answer is permitted by policy. Today the command
returns `Access is denied`, which tells the reader the ACL forbids it — the
opposite of the truth, and a message that sends them looking in entirely the
wrong place.

**Why this priority**: Status is the one service subcommand an unprivileged
user has a legitimate reason to run, and it is the first thing anyone runs when
the GUI reports the daemon unreachable. A wrong diagnosis at that moment costs
far more than the command itself.

**Independent Test**: From a non-elevated shell, as a user not in
`Administrators`, run `gosched service status` against the installed service
and observe a real answer. Separately confirm `service start` from the same
shell still fails, because the ACL genuinely withholds that right.

**Acceptance Scenarios**:

1. **Given** an installed service and a non-elevated standard user, **When**
   `gosched service status` runs, **Then** it reports `running` or `stopped`
   rather than an access error.
2. **Given** no installed service, **When** `gosched service status` runs,
   **Then** it reports the service is not installed, using wording unchanged
   from the current behavior.
3. **Given** the same non-elevated shell, **When** `gosched service start`
   runs, **Then** it still fails on privileges — that restriction is correct
   and must not be relaxed as a side effect.
4. **Given** Linux or macOS, **When** any service subcommand runs, **Then**
   behavior and output are identical to before this feature.

---

### User Story 3 - A bug report arrives with the facts needed to act on it (Priority: P2)

Someone hits a problem and opens an issue. Today they get an empty box, so what
arrives is prose without a version, without an install method, and without
whether they were elevated. Each of those three has already, on this repo,
been the fact that decided the diagnosis.

**Why this priority**: It costs a round trip per report rather than breaking
anything outright — but the round trip is unbounded, because the reporter may
never come back.

**Independent Test**: Open the *New issue* page on the repository and confirm
the blank-issue route is gone and that a form cannot be submitted without the
deciding fields.

**Acceptance Scenarios**:

1. **Given** the repository issue page, **When** a visitor clicks *New issue*,
   **Then** they are offered forms rather than an empty box.
2. **Given** the bug form, **When** a visitor omits version, component, install
   method, OS, or elevation state, **Then** the form cannot be submitted.
3. **Given** a submitted bug report, **When** a maintainer reads it, **Then**
   it carries pasted command output rather than a description of output.

---

### User Story 4 - A reader can install, run, and contribute without leaving the docs (Priority: P2)

Someone arriving at the repository should be able to get from the front page to
a first running task, find what every CLI command does, and — if they want to
contribute — learn how this project actually works, all without opening a spec
artifact. Today `README.md` routes command-level questions into `specs/`,
`TODO.md` describes a project mid-build that has in fact shipped and still
advertises a feature that was removed, only Windows has a real install guide,
and there is no `CONTRIBUTING.md`, `SECURITY.md`, or `CODE_OF_CONDUCT.md`.

**Why this priority**: Nothing here is broken, but the documentation set reads
as unfinished, which is a claim about the software that is not true.

**Independent Test**: Read `README.md` top to bottom as a newcomer and reach a
running task without opening `specs/`. Then follow the docs index to a CLI
reference that covers every command the binary actually exposes.

**Acceptance Scenarios**:

1. **Given** `README.md`, **When** a newcomer follows it, **Then** they install
   and create a first task without opening a `specs/` artifact.
2. **Given** the CLI reference, **When** it is compared against the command
   tree the binary exposes, **Then** every command and subcommand is covered.
3. **Given** `TODO.md`, **When** it is read, **Then** it describes only work
   that is genuinely open and mentions no removed feature.
4. **Given** each platform, **When** a reader looks for install instructions,
   **Then** a guide exists for that platform.

### Edge Cases

- **An already-open shell after install.** The `PATH` change is broadcast by
  the installer but does not reach shells that were already running. The
  documentation must say so, and must keep the full-path invocation available
  as the fallback for exactly that case rather than deleting it outright.
- **Upgrade rather than fresh install.** A major upgrade must not leave two
  `PATH` entries; component reference-counting handles this, and the
  requirement is stated so it is checked rather than assumed.
- **Service absent versus access denied.** These are different answers and must
  not collapse into one message; today the second is reported when the first is
  false and the third — "permitted, here is the answer" — is unreachable.
- **A partial upgrade where CLI and daemon versions differ.** The bug form must
  ask for both, because `gosched --version` and the version reported by
  `gosched health` can legitimately disagree.
- **A reader on a platform whose install guide is new.** The Linux and macOS
  guides describe paths and service managers that the Windows guide does not;
  each must stand alone rather than referring across.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The Windows installer MUST add its install directory to the
  machine `PATH` so that documented commands resolve by name.
- **FR-002**: The installer MUST remove that `PATH` entry on uninstall, and
  MUST NOT accumulate duplicate entries across upgrades.
- **FR-003**: A pre-build sanity check MUST assert the `PATH` entry is declared,
  so the regression cannot silently return.
- **FR-004**: `service status` MUST succeed for an unprivileged user whenever
  the service ACL grants status-query rights.
- **FR-005**: `service status` MUST request no more access than the query
  itself requires, and MUST NOT request start or stop rights.
- **FR-006**: `service start`, `stop`, `restart`, `install`, and `uninstall`
  MUST retain their existing privilege requirements.
- **FR-007**: Status output wording MUST be unchanged, including the
  not-installed case, and MUST be identical across platforms.
- **FR-008**: Non-Windows platforms MUST be behaviorally unaffected.
- **FR-009**: The repository MUST offer structured issue forms and MUST NOT
  offer a blank issue.
- **FR-010**: The bug form MUST require version, component, install method, OS
  and version, and elevation state, and MUST solicit pasted command output.
- **FR-011**: The feature form MUST ask for the problem before the proposed
  solution and MUST ask whether the request traces to the master specification.
- **FR-012**: The repository MUST carry a pull-request template that reflects
  the trunk-based workflow honestly rather than implying a review gate that
  does not exist for internal work.
- **FR-013**: `README.md` MUST be authored to the house Markdown style and MUST
  let a newcomer reach a first running task without opening a spec artifact.
- **FR-014**: A user-facing CLI reference MUST document every command the
  binary exposes, derived from the implementation rather than from the spec
  contract, which remains a spec contract.
- **FR-015**: Each supported platform MUST have an install guide.
- **FR-016**: `docs/` MUST carry an index distinguishing user-facing guides
  from maintainer material.
- **FR-017**: `TODO.md` MUST reflect delivered state and MUST NOT list a
  feature that has been removed from the product.
- **FR-018**: The repository MUST carry contribution, security, and conduct
  documents.
- **FR-019**: Every changed pinned artifact MUST be accompanied by a dated
  decision entry in `CHANGELOG.md`.

### Key Entities

- **Install-time `PATH` entry**: a machine-scoped, non-permanent, appended
  entry whose lifetime is bound to the presence of the CLI binary.
- **Service access mask**: the set of rights requested when opening the service
  control manager and the service handle; the defect is a mask larger than the
  operation needs.
- **Issue form**: a structured template whose required fields are the facts
  that decide triage — version, component, install method, OS, elevation.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Every command written in `README.md` and `docs/**` runs verbatim
  on a clean MSI install in a new shell, with no path prefix.
- **SC-002**: Uninstalling leaves zero `go-schedule` fragments in either `PATH`
  scope; a major upgrade leaves exactly one.
- **SC-003**: `gosched service status` returns a real status for a non-elevated
  standard user against a service whose ACL grants query rights.
- **SC-004**: `gosched service start` from that same shell still fails, and the
  Linux and macOS service paths are unchanged.
- **SC-005**: A new issue cannot be opened without the deciding fields, and the
  blank-issue route is unavailable.
- **SC-006**: The CLI reference covers every command and subcommand the binary
  exposes.
- **SC-007**: Every relative link in `README.md` and `docs/**` resolves against
  the working tree.
- **SC-008**: All six CI-parity gates are green, run in the foreground.

## Assumptions

- The install directory belongs on the **machine** `PATH` rather than the user
  `PATH`, because the package is per-machine and installs a system service; a
  per-user entry would leave the CLI invisible to the administrator who
  installed it for everyone.
- Fixing the status access mask is done in this repository rather than by
  forking or patching the upstream service library. The upstream helper that
  causes the defect is shared by paths that genuinely need broader access, so a
  local, status-only path is both smaller and safer than an upstream change.
- The end-to-end verification of the `PATH` fix cannot happen before release:
  it needs an installer built by the release workflow and installed on a clean
  machine. The pre-release checks are the sanity check and a read of the
  declaration; the issue stays open until the post-release check passes.
- `.github/ISSUE_TEMPLATE/**` and `.github/PULL_REQUEST_TEMPLATE.md` are not
  pinned artifacts — only `.github/workflows/**` is — so they require no dated
  decision. `build/**` and `docs/INSTALL-windows.md` are pinned and do.
- The documentation is written for the current release surface. Nothing here
  changes the scheduling engine, the store, the IPC layer, or the GUI.
