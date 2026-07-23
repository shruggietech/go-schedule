# Feature Specification: Maintainer Test Scripts and Vendored Skills

**Feature Branch**: `006-maintainer-test-scripts`

**Created**: 2026-07-23

**Status**: Draft

**Input**: User description: "Maintainer test scripts and vendored skills. Ship three
cross-platform script pairs under `test/scripts/` so a maintainer can verify an installed
`goschedd` actually fires tasks on time, survives restarts, catches up, and honors overlap
policies. Consolidate documentation into `docs/test-scripts.md`. Track `.claude/skills/`
in git and vendor the ShruggieTech house-standard skills."

## Clarifications

### Session 2026-07-23

Answered under the Build-Phase Autopilot decision policy (constitution principle V) from
the approved plan, the constitution, and the existing code, rather than escalated.

- Q: The scheduled moment a beat is compared against — where does it come from? → A:
  Three-tier precedence, recorded per beat so every drift figure is self-describing:
  scheduler-supplied environment value if present; otherwise snap to the nearest boundary
  of the caller-supplied interval; otherwise none.
- Q: Where do the databases live by default? → A: A user-writable per-user test directory,
  not the daemon's system-wide data directory, overridable by environment variable or an
  explicit path.
- Q: What does "record the run" mean for overlap detection? → A: Each beat records both its
  start and its finish, so overlap is the intersection of two intervals rather than an
  inference from start times.
- Q: Which platforms does the opt-in installer cover? → A: Only those the upstream project
  publishes prebuilt tools for; anything else exits with the prerequisite-failure code and
  names the package manager. No building from source.
- Q: How is database growth bounded? → A: It is not. The databases are disposable test
  artifacts; deleting one is the documented reset.

**Rationale for the drift answer, recorded because it is the feature's load-bearing
decision.** Inspection of the executor (`internal/executor/executor.go`) confirmed it
passes only the task's configured environment plus the inherited process environment — it
injects no scheduler-supplied variables at all, so no scheduled moment reaches a spawned
script today. The three alternatives were: (a) infer drift as deviation from the observed
cadence, (b) change the executor to inject the scheduled moment, (c) snap to the interval
boundary. Option (a) was rejected because it measures jitter, not latency: a scheduler
uniformly five seconds late looks perfect against its own history. Option (b) was rejected
because it changes a safety-critical product surface for a maintainer-tooling feature, and
the spec's no-product-change assumption is what keeps this release's binaries provably
identical. Option (c) was chosen: for a wall-clock-aligned interval schedule, the expected
moment is the beat's own timestamp snapped to the nearest interval boundary, which yields
*absolute* dispatch latency with no product change, and is unambiguous as long as drift
stays well below half the interval — a condition whose violation the missed-firing query
already reports. Tier one is kept ahead of it so that if a future release does inject the
scheduled moment, these scripts consume it with no change.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Prove the scheduler fires on time (Priority: P1)

A maintainer has just installed a release on a machine. They want evidence — not a
hopeful glance at a log — that the scheduler dispatches tasks when it said it would.
They register a recurring task pointing at the supplied heartbeat script, walk away,
come back later, and run the supplied reader against the resulting database. The reader
tells them how many times the task fired, how far each firing drifted from its scheduled
moment, and whether any expected firing was missed entirely.

**Why this priority**: This is the product's central promise. Every other verification
in this feature is a variation on it, and none of the existing surfaces (`gosched runs`,
the GUI logs view) quantify dispatch drift or detect a silently missed firing.

**Independent Test**: Register one recurring task against the heartbeat script, let it
run for several intervals, and run the reader's cadence, drift, and gaps queries. Delivers
a complete on-time-firing verdict with no other part of this feature present.

**Acceptance Scenarios**:

1. **Given** a running daemon and a task scheduled every minute against the heartbeat
   script, **When** ten minutes elapse, **Then** the heartbeat database contains ten
   beat records and the reader's gaps query reports none.
2. **Given** the same ten-minute window, **When** the maintainer runs the drift query,
   **Then** it reports the distribution of the difference between each beat's actual and
   scheduled moment, including a 99th-percentile figure comparable against the project's
   documented dispatch budget.
3. **Given** a heartbeat script invoked directly from a terminal with no daemon involved,
   **When** it completes, **Then** exactly one beat record is written and the script exits
   reporting success.
4. **Given** a task whose script invocation is instructed to fail, **When** it runs,
   **Then** the beat record is still written, the recorded outcome marks the failure, and
   the script reports failure to the scheduler so the daemon's own alerting reacts.

---

### User Story 2 - Prove behavior across restarts, downtime, and overlap (Priority: P2)

A maintainer wants to confirm the harder guarantees: that scheduled work survives a daemon
restart, that a task missed during downtime is made up exactly once rather than zero or
many times, and that a long-running task under each overlap policy behaves as documented.
They use the same heartbeat script — instructed to run slowly when overlap is the subject —
and the reader's restart and overlap queries to read the verdict out of recorded history.

**Why this priority**: These are the constitution's named safety-critical surfaces and the
behaviors most likely to regress unnoticed, but they depend on the P1 recording and reading
machinery already existing.

**Independent Test**: With the heartbeat machinery in place, stop the daemon across a
scheduled firing, restart it, and confirm the reader reports exactly one make-up beat and a
visible session boundary.

**Acceptance Scenarios**:

1. **Given** beats recorded before a daemon restart, **When** the daemon restarts and more
   beats are recorded, **Then** the reader's restart query shows the boundary and confirms
   recording continued across it.
2. **Given** the daemon stopped past one scheduled firing under a make-up-once policy,
   **When** it restarts, **Then** exactly one catch-up beat is recorded before the normal
   cadence resumes.
3. **Given** a task whose run takes longer than its interval, **When** it runs under each
   of the three overlap policies in turn, **Then** the reader's overlap query distinguishes
   the three outcomes: queued-and-serialized, skipped, and genuinely concurrent.
4. **Given** two invocations writing to the same database at the same moment, **When** both
   complete, **Then** both records are present and neither invocation fails on contention.

---

### User Story 3 - Capture machine state as a realistic scheduled workload (Priority: P2)

A maintainer wants a scheduled task that does something substantive rather than writing a
single row — a workload that touches the operating system, takes real time, and produces
inspectable output. The system-information script records a timestamped snapshot of the
host: identity, logged-in user, process count, uptime, network addresses, and listening
ports. Run on a schedule, it becomes both a scheduler test and a rudimentary machine-state
history the maintainer can diff.

**Why this priority**: It exercises a materially different execution profile from the
heartbeat — subprocess spawning, platform-specific tooling, multi-row writes — which is
where cross-platform execution bugs actually surface. It is independent of the timing
verification, hence equal priority rather than blocking.

**Independent Test**: Invoke the script once by hand on each supported platform and confirm
a complete snapshot lands with its address and port detail attached.

**Acceptance Scenarios**:

1. **Given** any supported platform, **When** the script runs once, **Then** exactly one
   snapshot is recorded carrying timestamp, host identity, user, process count, and uptime.
2. **Given** the same run, **When** the maintainer inspects the snapshot, **Then** the
   machine's IPv4 and IPv6 addresses and its listening ports are attached to it as
   individually queryable detail rather than as opaque text.
3. **Given** a platform where one probe is unavailable, **When** the script runs, **Then**
   the snapshot is still recorded with that datum marked absent and a warning emitted,
   rather than the whole snapshot failing.
4. **Given** two snapshots taken at different times, **When** the maintainer runs the
   listeners query, **Then** the change in listening ports between them is reported.

---

### User Story 4 - Get the house standards on a fresh clone (Priority: P3)

A maintainer clones the repository onto a new machine. Today they get no spec-kit commands
and no house-standard authoring guidance, because the agent-configuration directory is
excluded from version control wholesale; the build-autopilot document already names this as
a setup failure. After this feature, the skills directory is tracked, so a fresh clone
arrives with the spec-kit commands and the house standards already present, while
credential-bearing agent configuration stays excluded.

**Why this priority**: Real friction, but it blocks nobody who already has a working
checkout, and it carries no runtime behavior.

**Independent Test**: Clone the repository into a clean directory and confirm the skills
are present and no credential-bearing configuration file came with them.

**Acceptance Scenarios**:

1. **Given** a fresh clone, **When** the maintainer inspects the agent skills directory,
   **Then** the spec-kit command skills and the vendored house-standard skills are present.
2. **Given** a working checkout, **When** agent tooling writes a settings or credential file
   into the agent configuration directory, **Then** that file is not staged for commit.
3. **Given** a maintainer preparing to push, **When** they invoke the project's verification
   skill, **Then** they receive the exact verification command sequence and the documented
   local-environment traps without having to re-derive them from prose.

---

### Edge Cases

- **The database tool is missing.** The scripts depend on an external command-line SQLite
  tool. When it cannot be located, the script must exit with the reserved
  prerequisite-failure code and name both the opt-in installer and the platform package
  manager command — never fail with a generic error, and never silently produce no data.
- **The installer's download is tampered with or corrupted.** The opt-in installer verifies
  a pinned checksum. A mismatch is a hard failure; an unverified binary is never installed.
- **Concurrent writers.** Overlap-policy testing produces simultaneous writers by
  construction. Contention must be waited out and retried, not surfaced as a failure that
  looks like a scheduler defect.
- **A runaway loop.** The heartbeat's opt-in continuous mode must be bounded by a maximum
  beat count or duration; an unbounded loop under a scheduler is a resource incident.
- **Cross-platform command absence.** Host-inspection commands differ per platform and some
  are absent even on their own platform. Each probe degrades to a recorded absence.
- **A PowerShell script run on a non-Windows host.** The PowerShell scripts must not assume
  Windows-only host-inspection commands are present, because the shell itself is
  cross-platform and a maintainer may reasonably run either twin anywhere.
- **First run against a nonexistent database.** The database and its structure are created
  on demand; a first run is indistinguishable from a subsequent one to the maintainer.
- **A run with no expected moment available.** When neither a scheduler-supplied value nor
  a caller-supplied interval is present — a bare invocation from a terminal, say — the beat
  records no expected moment and no difference, and the drift query excludes it and says how
  many records it excluded. Reporting a drift of zero for an unmeasurable run would be
  worse than reporting nothing.
- **Drift approaching half the interval.** Snapping to the nearest interval boundary is
  unambiguous only while drift stays well under half the interval. Beyond that a late firing
  is indistinguishable from an early next one, so the reader MUST flag any boundary-sourced
  drift exceeding a quarter of the interval as unreliable rather than reporting it as fact.

## Requirements *(mandatory)*

### Functional Requirements

**Recording — heartbeat**

- **FR-001**: The heartbeat script MUST, by default, record exactly one beat per invocation
  and then exit, so that the schedule's cadence comes from the scheduler rather than from
  the script.
- **FR-002**: Each beat MUST record both the moment its run started and the moment it
  finished, the host, the operating-system process identity, a per-invocation session
  identity, a monotonically increasing sequence position, and the outcome of the
  invocation. Recording both endpoints is what makes overlap a decidable question rather
  than an inference from start times alone.
- **FR-003**: Each beat MUST record its expected moment, the difference between its actual
  and expected moments, and **which of three sources the expected moment came from**:
  a scheduler-supplied environment value; a snap to the nearest boundary of the
  caller-supplied interval; or none available. A drift figure MUST never be presented
  without its source, so that an inferred figure is not mistaken for a measured one.
- **FR-003a**: The three sources MUST be consulted in that order of precedence, so that a
  future release supplying the scheduled moment directly is consumed without changing these
  scripts.
- **FR-004**: The heartbeat script MUST offer an opt-in continuous mode bounded by a maximum
  beat count, a maximum duration, or both. **Whichever bound is reached first ends the loop.**
  When the caller supplies neither, a default maximum duration MUST still apply — continuous
  mode has no unbounded form, because an unbounded loop launched under a scheduler is a
  resource incident rather than a test.
- **FR-004a**: The duration bound MUST be evaluated between beats, not during one. A single
  deliberately-slow run may therefore overrun the bound by up to that run's length; this is
  the intended behavior, because interrupting a run mid-write would corrupt the very record
  the bound exists to protect. The documentation MUST state this.
- **FR-005**: The heartbeat script MUST offer a knob that deliberately extends a run's
  duration, so that overlap policies can be exercised without a contrived external workload.
- **FR-006**: The heartbeat script MUST offer a knob that deliberately reports failure with
  a caller-chosen code, and MUST still record the beat when it does. The knob MUST reject
  the codes reserved by FR-019 for success and for usage failure, so that an induced failure
  can never be mistaken for either a success or an unmet prerequisite.

**Recording — system information**

- **FR-007**: The system-information script MUST record, per invocation, one snapshot
  carrying: the moment in both epoch and human-readable form, the local-versus-universal
  offset, the host name, the logged-in user, the count of running processes, the host
  uptime, and the platform.
- **FR-008**: The snapshot MUST carry the host's IPv4 and IPv6 network addresses and its
  listening ports as separately queryable detail records, not as embedded text.
- **FR-009**: A probe that cannot run on the current host MUST record its datum as absent
  and emit a warning, and MUST NOT abort the snapshot.
- **FR-010**: The snapshot MUST record a caller-supplied label identifying what invoked it,
  so snapshots from different scheduled tasks are distinguishable.

**Reading**

- **FR-011**: The reader script MUST accept either a well-known database name or an explicit
  path, and MUST select among a set of named canned queries.
- **FR-011a**: A well-known name MUST resolve against a **user-writable per-user test
  directory**, never the daemon's system-wide data directory — writing there would require
  elevation and would mix disposable test artifacts with production scheduler state. The
  location MUST be overridable by environment variable, and every script MUST report the
  resolved path it used.
- **FR-012**: The reader MUST be able to list its available queries.
- **FR-013**: The reader MUST provide, at minimum, queries answering: overall summary; most
  recent records; observed cadence distribution; drift distribution; missed or delayed
  firings; overlapping executions; failed runs; restart and session boundaries; host
  breakdown; listening-port change between snapshots; and the stored structure.
- **FR-013a**: Any query that excludes records from its result MUST report how many it
  excluded and why. A distribution computed over an unstated subset invites a confident
  conclusion drawn from a fraction of the evidence.
- **FR-013b**: A result set spanning records with different expected-moment sources MUST
  break its figures down by source rather than pooling them, since a measured value and a
  derived one are not the same kind of number and averaging them produces neither.
- **FR-014**: The reader MUST offer both a human-readable and a machine-readable output
  form.

**Shared behavior across all scripts**

- **FR-015**: Every script MUST be delivered as a matched pair — one for PowerShell and one
  for POSIX shell — with equivalent behavior and one-to-one corresponding options.
- **FR-016**: Every script MUST locate the external database tool by checking, in strict
  order of precedence: an explicitly supplied path, a repository-local location, and then
  the system search path. The first candidate that exists **and is usable** wins.
- **FR-016a**: A candidate that exists but cannot be executed, or that reports a version
  below the documented minimum, MUST be treated as not found and the search MUST continue to
  the next candidate — a stale tool earlier in the order must not shadow a good one later in
  it. A minimum version MUST be documented, chosen as the lowest version providing every
  output form FR-014 requires.
- **FR-017**: When the database tool cannot be located, every script MUST exit with the
  reserved prerequisite-failure code and MUST name both the opt-in installer option and the
  platform-appropriate package-manager command.
- **FR-018**: Every script MUST offer an opt-in installer that retrieves the official
  database tool for the current platform, verifies it against a checksum pinned in the
  repository, and refuses to install on mismatch. No script may access the network unless
  this option is explicitly given.
- **FR-018a**: The installer MUST cover only the platform and architecture combinations for
  which the upstream project publishes prebuilt tools, and MUST NOT attempt to build from
  source. On any other combination it MUST exit with the prerequisite-failure code and name
  the platform's package manager.
- **FR-018b**: An artifact that fails verification MUST be discarded, and MUST NOT be left on
  disk where a subsequent run could find it. An unverified binary parked next to a verified
  one is worse than no binary at all.
- **FR-018c**: The pinned checksums MUST live in the repository beside the code that checks
  them, and MUST record the upstream version they correspond to. Updating them is a
  deliberate, reviewable act — never a convenience performed to make a failing download
  succeed.
- **FR-018d**: Each host-inspection probe's fallback ordering MUST be documented and fixed,
  not left to implementation discretion. An ordering that varies by machine produces data
  whose provenance differs between hosts while looking identical in the database.
- **FR-019**: Every script MUST use a distinct exit code for each of: success, runtime
  failure, and usage or unmet-prerequisite failure.
- **FR-020**: Every script MUST write results to standard output and diagnostics to standard
  error, consistent with the project's existing command-line behavior.
- **FR-021**: Every script MUST create its database and structure on demand, and MUST
  tolerate concurrent writers by waiting and retrying rather than failing. The wait and the
  retry count MUST both be bounded and documented; the bound MUST be generous enough that
  ordinary overlap-policy testing never reaches it.
- **FR-021a**: Exhausting the contention bound MUST be reported as a runtime failure naming
  contention as the cause. It MUST NOT be silently swallowed, and it MUST NOT be reported in
  a way that resembles a scheduler defect — a maintainer debugging a lost record needs to
  know the loss happened in the test harness, not in the product.
- **FR-021b**: A verification artifact that cannot record its observation MUST fail loudly.
  Silence from a test payload is indistinguishable from a scheduler that never fired, which
  is precisely the failure this feature exists to detect.
- **FR-021c**: A beat MUST be written **once, at the end of the run**, carrying a start
  moment captured in memory when the run began. Two writes per beat would double contention
  for no gain. The consequence is deliberate and MUST be documented: a run interrupted
  mid-flight records nothing and therefore appears as a missed firing — which is the correct
  signal, because from the maintainer's position an interrupted run and an absent one have
  the same meaning.
- **FR-021d**: Within each language, the prerequisite resolution, installation, and database
  access logic MUST have exactly one implementation shared by all three scripts of that
  twin. Three copies of the resolution order would be three chances for them to disagree,
  and a disagreement here presents as an intermittent, platform-specific test failure.
- **FR-021e**: Both twins MUST record timestamps in one documented format and precision, so
  that records produced by different twins are directly comparable in the same query.
- **FR-022**: Every script MUST provide built-in usage help describing every option.
- **FR-022a**: The databases MUST NOT be pruned, rotated, or size-capped automatically.
  They are disposable test artifacts and deleting one is the documented reset; automatic
  retention would silently destroy the very history a maintainer is trying to inspect. The
  documentation MUST state the approximate growth rate so the choice is informed.

**Documentation**

- **FR-023**: All documentation for these scripts MUST be consolidated into a single
  document, covering prerequisites, a quickstart, per-script option reference for both
  twins, the recorded structures, the query catalog, the exit-code contract, and
  troubleshooting.
- **FR-024**: The document MUST include worked end-to-end recipes for each verification the
  feature exists to support: on-time firing, restart survival, downtime catch-up, each
  overlap policy, and failure reporting.
- **FR-025**: The document MUST state the authoring conventions the POSIX-shell twins follow,
  since no separate governing standard exists for them.

**Repository configuration**

- **FR-026**: The repository MUST track the agent skills directory while continuing to
  exclude all other agent configuration, so that credential-bearing files are not tracked by
  default.
- **FR-027**: The repository MUST NOT track the databases the scripts produce, nor any
  locally installed database tool.
- **FR-026a**: The exclusion rule MUST be expressed as exclude-everything-then-narrowly-admit,
  not as a list of things to exclude. A denylist admits every file nobody thought of; an
  allowlist admits only what was named. Since the excluded material is credential-bearing by
  assumption, the failure directions are not symmetric.
- **FR-026b**: Tracking the skills subtree MUST NOT be assumed to make that subtree safe. The
  change MUST be accompanied by a check that no credential-bearing file is being introduced
  by it, performed before the work is committed.
- **FR-028**: The tracked skills MUST include the project's spec-kit command skills, the
  house standards governing the artifacts this feature produces, and a project-native skill
  carrying the verification procedure.
- **FR-029**: The repository-configuration change MUST be recorded as a dated decision,
  because the file it changes is one of the project's pinned process artifacts and the
  governing procedure requires such a record. This is a process obligation the other
  requirements do not imply, and it is the one most easily lost between planning and commit.

### Key Entities

- **Beat**: One recorded execution of the heartbeat script. Carries its start and finish
  moments, its expected moment together with the source that expected moment came from and
  the resulting difference, the host and process that produced it, its session and sequence
  position, and its outcome.
- **Session**: A set of beats sharing one continuous run of the recording process, used to
  identify restart and catch-up boundaries.
- **Snapshot**: One recorded execution of the system-information script, describing the host
  at a moment.
- **Address**: One network address belonging to a snapshot, with its family and interface.
- **Listening port**: One port the host was accepting connections on at a snapshot, with its
  protocol and owning process where discoverable.
- **Canned query**: A named, pre-written inspection the reader exposes, with a stated
  question it answers.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A maintainer with a freshly installed release can go from a clean machine to a
  quantified on-time-firing verdict in under fifteen minutes of elapsed time and under five
  minutes of hands-on work, following the documentation alone.
- **SC-002**: Over a ten-firing observation window on an otherwise idle machine, the reader
  reports zero missed firings and a 99th-percentile drift figure directly comparable against
  the project's documented dispatch budget, labelled with the source that figure was derived
  from so it is never mistaken for a precision it does not have.
- **SC-003**: All three verification scenarios that today have no maintainer-facing check —
  restart survival, downtime catch-up, and each of the three overlap policies — are
  answerable by running one documented recipe each.
- **SC-004**: Every script behaves equivalently across the supported platforms: the same
  option produces the same recorded result, and the matched pairs are interchangeable.
- **SC-005**: A maintainer with no database tool installed reaches a working setup via a
  single documented option or a single package-manager command, and the failure message they
  first encounter names both.
- **SC-006**: A fresh clone of the repository arrives with the spec-kit commands and house
  standards already available, with no post-clone installation step, and with no
  credential-bearing configuration file tracked.
- **SC-007**: Concurrent writers under the most permissive overlap policy lose no records.

## Assumptions

- **The scripts are test payloads, not product surface.** They are maintainer tooling that
  ships in the repository; they are not installed by the release packages and carry no
  compatibility promise to end users.
- **The external database tool is a dependency, not a vendored artifact.** The repository
  does not commit third-party executables; it detects, and optionally installs on request.
- **The scheduler supplies no context through the environment today.** Verified against the
  executor: a spawned task receives the inherited process environment plus the task's own
  configured variables, and nothing scheduler-generated. Every recorded field must therefore
  be derivable by the script itself, and the expected-moment precedence in FR-003a exists so
  that a future release which does inject that context is picked up for free. Nothing in
  this feature depends on that release happening.
- **PowerShell means the cross-platform edition.** Scripts target PowerShell 7 or later,
  which is available on all three supported platforms; Windows PowerShell 5.1 compatibility
  is not a goal.
- **The POSIX twins target Linux and macOS.** Windows maintainers use the PowerShell twins,
  though the POSIX twins remain usable under a POSIX-compatible environment on Windows.
- **No product code changes.** This feature adds maintainer tooling, documentation, and
  repository configuration only. The daemon, the command-line client, the graphical client,
  and the stored data are untouched, so the shipped binaries are byte-for-byte unaffected.
- **Continuous-integration configuration is untouched.** The feature's automated tests run
  within the existing test invocation, so no workflow change is required.
- **Skills are vendored by copying, not by reference.** Plugin-managed skills are
  deliberately excluded, because copying them would register them twice and let the copies
  drift from their source. Each vendored copy records where it came from, so that drift is
  at least *detectable*; refreshing a copy is a deliberate act, and no automatic
  synchronization is claimed or implied.
- **No shell authoring standard exists yet.** The PowerShell twins are governed by a
  vendored house standard; the POSIX twins are governed by a conventions section in the
  documentation plus an automated shell linter, with a dedicated standard deferred.
