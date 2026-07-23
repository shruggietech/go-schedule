# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- **Pinned artifact — `.github/workflows/release.yml` (2026-07-23).** The
  `readme-badge` job now rewrites *two* README version strings on each `v*.*.*`
  tag, not one. Alongside the static release-badge line it also rewrites the
  quick-start's illustrative `gosched health` output line — `daemon ok (version
  X.Y.Z)` — to the tag's version (leading `v` stripped, matching
  `buildinfo.Version`'s form) via a second `sed` anchored on `daemon ok
  (version `, committed and pushed in the same commit as the badge bump. That
  line was previously untouched by the release automation and so drifted one
  release behind after every tag (it was hand-fixed for v0.7.0 and again for
  v0.8.0). A no-match is a no-op — `sed` still exits `0`, so the drift-fix can
  never fail a release, preserving the job's standalone-and-non-blocking design.
  The job and its step were renamed from "badge" to "version lines" to reflect
  the widened scope.

### Decisions

- **2026-07-23** — **The release automation syncs the two version strings the
  README bakes in, and only those two.** README carries three kinds of version
  reference: the release badge, the `gosched health` sample output, and the
  `<ver>` placeholders in the download table. The first two are concrete and must
  track the current release, so both are now rewritten on tag. The third is a
  deliberately generic placeholder that teaches the *shape* of an asset filename
  (`go-schedule_<ver>_<os>_<arch>`); substituting a real version there would
  imply a single download rather than a family and is intentionally left alone.
  The `daemon ok` line strips the leading `v` because that is the form the daemon
  actually prints (`buildinfo.Version`) and the form the README already used.

## [0.8.0] - 2026-07-23

### Added

- **The p99 dispatch-latency budget is now measured and enforced, and the engine
  benchmarks run in CI (closes #14).** The constitution budgets dispatch latency
  at p99 < 100 ms, but nothing measured it: `internal/engine/engine_bench_test.go`
  had `BenchmarkDispatch`/`BenchmarkNextRun` that no workflow ran, and
  `testing.B` reports a mean, not the p99 the budget is stated in. Two changes
  close the loop. First, `TestDispatchLatencyP99`
  (`internal/engine/latency_test.go`) dispatches 2000 runs serially through the
  worker pool, measures each run's scheduled-time→execution-start latency
  (command execution excluded), and asserts the p99 against a new
  `engine.DispatchLatencyBudget = 100 * time.Millisecond` constant that lives next
  to the dispatch code it governs. It runs in the standard suite (locally and in
  the race job) and is cgo-free; observed p99 is microscopic against the ceiling,
  so it is stable on loaded CI hardware. Second, a `bench` CI job runs the engine
  benchmarks on every push/PR and publishes their output as a build artifact.
  The goroutine-leak test (`test/integration/leak_test.go`) already runs under
  `-race` in CI — confirmed, no change needed.

- **Documentation is now published as a site, and `docs/` is the single source
  of truth (closes #11).** The `docs/` set — the install guides, the CLI and GUI
  references, cron interoperability, and the test-scripts and build-autopilot
  guides — is published as a searchable, navigable GitHub Pages site served
  branch-based from the `docs/` folder on `main` using the just-the-docs remote
  theme, so the Markdown in the repository is both the reviewable source and the
  served page. Every page gained `title`/`nav_order` front matter, the three
  install guides are grouped under an Installation section, and the eleven links
  that pointed out of `docs/` were rewritten to absolute repository URLs so
  nothing 404s on a `docs/`-rooted site. A new `scripts/docs-check.sh` gate
  (front matter, on-disk link integrity, no links escaping `docs/`,
  pointer-README validity; no network) runs locally and as a `docs` CI job. The
  README and the issue-form contact links now point at the site, and the README
  quick-start version was corrected from 0.6.0 to 0.7.0. Going live is one
  repository setting: Pages → Deploy from a branch → `main` / `docs`.

### Changed

- **Pinned artifact — `.github/workflows/ci.yml` (2026-07-23).** Added a `bench`
  job (`ubuntu-latest`, `CGO_ENABLED=0`) that runs
  `go test -run '^$' -bench . -benchmem ./internal/engine/...` and uploads the
  output via `actions/upload-artifact@v4`. It is informational — the enforced
  dispatch-latency gate is `TestDispatchLatencyP99` in the `test` job, not this
  job — so a benchmark run never fails the build.

- **Pinned artifact — `.github/workflows/ci.yml` (2026-07-23).** Added a `docs`
  job (`ubuntu-latest`, no Go toolchain) that runs `sh scripts/docs-check.sh` on
  every push/PR to `main`, guarding the documentation-site sources — front
  matter, on-disk link integrity, no links escaping `docs/`, and pointer-README
  validity. It runs the exact script contributors run locally, so the two cannot
  drift.

- **Pinned artifact — `docs/INSTALL-windows.md` (2026-07-23).** Added a
  `title`/`parent`/`nav_order` YAML front-matter block, matching every other
  `docs/` page, so the Windows install guide appears under the site's
  Installation section. The prose is unchanged; GitHub renders front matter
  invisibly, so reading the file in the repository is unaffected.

### Decisions

- **2026-07-23** — **The dispatch-latency regression gate asserts the absolute
  p99 budget, not a relative benchmark delta.** The constitution requires CI to
  enforce "benchmark regression checks" and, separately, that a benchmark not
  regress by more than 10 % "without explicit, recorded justification." A
  benchstat percentage-delta gate against a stored `bench.txt` baseline was
  considered and rejected: on shared CI runners a 10 % delta fires on scheduler
  noise, and a gate that fires on noise gets disabled — worse than no gate. The
  absolute p99 assertion is stable and is the exact property the constitution
  budgets (a change that pushes p99 over 100 ms fails the build), so it satisfies
  the regression-check obligation; this entry is the recorded justification the
  10 % clause calls for. The benchmarks still run in CI and their output is
  published as an artifact, preserving the raw trend for spotting a within-budget
  slowdown by eye.

- **2026-07-23** — **The documentation site is served branch-based from `docs/`
  with the just-the-docs remote theme, not built by Hugo or MkDocs via GitHub
  Actions.** Branch-based serving keeps the `docs/` Markdown as both the
  reviewable source and the served content, with no deploy workflow and a single
  operator settings change to go live; Hugo and MkDocs with an Actions deploy
  (issue #11's alternatives) each add a build pipeline and a second toolchain for
  no gain at this size. The theme is pinned to `just-the-docs@v0.4.2`, the last
  release that builds under GitHub Pages' bundled Jekyll 3.9 (libsass); adopting
  a newer theme would require a Jekyll 4 Actions build, a deliberately deferred
  future change.

- **2026-07-23** — **`docs-check.sh` validates link targets, not anchors.**
  Fragments (`#section`) are stripped and not resolved: validating them would
  mean replicating Jekyll/kramdown heading-slug rules for little value, whereas
  file existence, the no-escape rule, and front-matter presence catch the drift
  that actually breaks the published site.

### Fixed

- **README release badge no longer breaks on shields.io token starvation
  (2026-07-23).** The badge used shields.io's dynamic
  `img.shields.io/github/v/release` endpoint, which calls the GitHub API
  server-side from shields.io's shared token pool; when that pool is exhausted
  the badge renders the literal error `Unable to select next GitHub token from
  pool` instead of the version. It is now a static `badge/release-vX.Y.Z-blue`
  URL (no API call, so the error is structurally impossible), matching the
  License badge.

### Changed

- **Pinned artifact — `.github/workflows/release.yml` (2026-07-23).** Added a
  standalone `readme-badge` job that, on each `v*.*.*` tag push, rewrites the
  static README release badge to the tag version and commits it back to `main`
  as `github-actions[bot]`. The job carries no `needs`, so a badge-bump failure
  (e.g. branch protection) cannot fail the release artifacts.

## [0.7.0] - 2026-07-23

### Added

- **Cron interoperability — `gosched cron import`, `explain`, and `export`
  (closes #12).** Everyone who would adopt this project already has a crontab,
  and until now the only way across was to read each line, hold it in your head,
  and retype it. `cron import --file <path>` reads a crontab and creates a task
  per line; `--dry-run` prints the identical report and creates nothing, which
  is both the migration preview and the answer to "did it understand my
  crontab?". `cron explain "0 9 * * 1-5"` translates one expression with no side
  effects. `cron export` gives the task set back as crontab lines.

  The conversion never approximates. Everything it will not carry is refused by
  name — `@reboot`, six-field Quartz expressions, `L`/`W`/`#`, a step that does
  not divide its range (`*/7` restarts at :00, which a fixed interval cannot
  reproduce), and an expression restricting both day-of-month and day-of-week
  (cron means "either"; the recurrence model means "both"). `MAILTO` and shell
  variable assignments are reported as warnings rather than dropped. The import
  summary states the fidelity facts outright: cron carries no timezone, no
  catch-up, no overlap policy and no restart recovery, so it names which zone was
  applied and which defaults the imported tasks received.

  Cron remains an interchange format and never becomes an authoring syntax:
  `--schedule "0 9 * * 1-5"` is still an error and no GUI field accepts an
  expression. `docs/cron.md` carries the full fidelity table in both directions.

- **Schedules addressed by calendar date and by year.** `on the 15th of every
  month`, `the 31st monthly at 09:00`, `every year on february 29`, `annually on
  4 july`, `every 12 months`. Without these, ordinary cron lines like
  `0 9 1 * *` had no target representation at all, so cron import could not have
  been complete.

- **A per-task missing-date policy (closes the calendar half of #8).** A rule on
  the 31st, on 29 February, or on the fifth Friday meets periods that have no
  such date. Until now the behavior was an implicit rule the task owner could not
  see, state, or change — and the stored summary lied about it: "The 5th Friday
  of every month" for a rule that fires four times a year. Each task now states
  its intent: `skip` (the default, and exactly what every existing task already
  did), `last_valid` (Feb 29 → Feb 28, the 31st → the 30th, a missing fifth
  Friday → the last Friday), or `next_valid` (roll into the next period without
  displacing that period's own run). Settable from `gosched task add|edit
  --missing-date`, shown by `task show`, and present in the GUI editor's Advanced
  Settings.

  Schedule descriptions now name the policy instead of asserting "every month"
  for a rule that skips months. The description is rendered when a task is read
  rather than stored, because the policy can change without the phrase changing —
  a stored sentence would go stale the moment an operator switched it.

  Deliberately still open on #8: DST anchoring (wall-clock versus elapsed-time
  versus UTC) and per-task skipped-hour/repeated-hour resolution.

- **The ShruggieTech attribution in the README footer is now a link** (closes
  #9). It was the only proper noun in a document that links every other one.

### Changed

- **Store migration v5** adds `tasks.missing_date_policy`, additive with a total
  default of `skip`. Forward-only and non-destructive: no existing column, row,
  or value is read or rewritten, and `skip` is the behavior every pre-v5 task
  already had, so no installed task's run times move. Pinned by
  `internal/store/migration_v5_test.go`, which asserts a v4-era database upgrades
  with every task row otherwise byte-identical and the schedules table untouched.

- **`schedule.NextRun` and `schedule.UpcomingRuns` take the missing-date
  policy.** Six call sites across the engine, catch-up, and API packages pass it
  through; all already held the task.


- **`TODO.md` removed; the roadmap is now the GitHub issue tracker.** The file
  had become a second, worse issue tracker: eight open items written as prose
  bullets that could not be labelled, discussed, assigned, or closed, in a
  document a reader had to be pointed at. Its "Delivered" section duplicated
  `CHANGELOG.md`, and its "Open" section duplicated nothing — that was the
  problem, since the work it described was invisible to anyone browsing issues.

  Each remaining item was filed with the context needed to act on it rather than
  transcribed: **#13** IPC access control (records that `config.AdminGroup`
  already exists and is inert, and that the Windows `AU` ACE is load-bearing for
  the non-elevated case), **#14** benchmarks and the p99 budget, **#15** signing
  and notarization, **#16** end-to-end verification of the `PATH` fix, and the
  four deferred post-v1 items as **#17**–**#20**, each marked deferred rather
  than rejected.

  One item was corrected in the move. `TODO.md` claimed a goroutine-leak test
  and a dispatch benchmark did not exist; both do —
  `test/integration/leak_test.go` and `internal/engine/engine_bench_test.go`.
  What is actually missing is that nothing runs the benchmarks (no CI job
  invokes `-bench`) and that `testing.B` reports a mean while the constitution
  budgets a p99. #14 states the real gap.

### Decisions

- **2026-07-23** — **The missing-date policy lives on the task, not the
  schedule.** Storing it on the schedule row looked cheaper: `NextRun` already
  receives the schedule, so no signature would have changed. Reading
  `internal/api/server/update.go` settled it against that. An edit supplying a
  new schedule phrase *creates a new schedule row* and repoints the task at it,
  so a policy stored there would be silently reset to the default by any phrase
  edit unless a carry-over were remembered at that one site. That is a silent
  change to run times — the class of defect issue #4 already produced in the task
  editor — and correctness that depends on remembering to copy a field is not
  correctness. The cost is a parameter added to two functions at six
  compile-checked call sites, which is the cheaper half of the trade.

- **2026-07-23** — **The cron parser is written in-tree rather than taken as a
  dependency.** The constitution's Engineering Constraints prefer the standard
  library where it suffices and require every dependency to be justified. A cron
  library (`robfig/cron`, `adhocore/gronx`) parses an expression into a compiled
  schedule, which is the opposite of what this feature needs: the work is
  inspecting *field structure* to decide what cannot be represented. Such a
  library accepts `*/7` and hands back a working schedule, when the honest answer
  is a refusal — so we would parse twice, and the second parse would be the one
  that mattered. The grammar is about 120 lines and fully covered by table tests.
  No new dependency was added.

- **2026-07-23** — **A cron expression becomes a schedule only by way of the
  human phrase.** The converter renders the phrase a user would have typed and
  hands it to the existing grammar; an expression with no phrase is refused. This
  is what makes the import preview trustworthy rather than advisory — the string
  shown is literally the string parsed and stored — and it keeps "cron is not an
  authoring syntax" structural rather than a matter of discipline, since the
  converter has no route into the engine that an operator does not also have. The
  cost is real and accepted: cron can express schedules the phrase grammar
  cannot (arbitrary by-minute lists), and those are refused rather than given a
  privileged back door into a task nobody could subsequently edit.


- **2026-07-23** — **The roadmap moves to the issue tracker rather than being
  reorganized in place.** The alternative was to keep `TODO.md` as a curated
  index pointing at issues, which reads tidy and reintroduces the failure it was
  meant to fix: a second list to keep in sync, silently wrong the first time an
  issue closes without someone remembering the file. `gh issue list` is
  generated from the state itself and cannot drift. References in `CLAUDE.md`,
  `docs/build-autopilot.md`, and the constitution now name the tracker;
  historical `CHANGELOG.md` and `specs/` entries mentioning `TODO.md` were left
  alone, because they are records of what happened and not instructions.

- **2026-07-23** — **Constitution amended to v2.0.1** (PATCH). Principle V's
  standing autopilot authorization named `TODO.md` as the second source of
  traceable scope; it now names the issue tracker. What autopilot may run
  without further authorization is unchanged — only where that scope is
  recorded — which is a clarification rather than a governance change.

## [0.6.0] - 2026-07-23

**The documented commands now work as written, and an ordinary user can ask
whether the scheduler is running.** Two Windows defects, both on the seam
between the shipped product and the person trying to use it, plus the
documentation set the project had been shipping without.

The minor bump is carried by the installer: it now writes to the machine `PATH`,
which is install-time behavior a user observes, and it changes two pinned
artifacts. Scheduling behavior is untouched.

### Fixed

- **The Windows `.msi` never added its install directory to `PATH`** (#5). Every
  command in the README and in `docs/test-scripts.md` is written as a bare
  `gosched ...`, and after a normal install none of them resolved — the first
  thing a new Windows user typed failed, and failed in a way that reads as a
  broken package rather than a missing `PATH` entry.

  It survived several releases for a reason worth recording: every machine where
  this project is developed or tested already has that directory on `PATH`, put
  there by hand. The defect was invisible from inside the project and unmissable
  from outside it.

- **`gosched service status` demanded elevation it did not need** (#6). The
  installed service's ACL grants Interactive Users `SERVICE_QUERY_STATUS`, so a
  read-only status query is permitted by policy — yet it failed with
  `Access is denied` for any non-elevated user.

  The cause was the access mask, not the ACL. The status path opened the service
  handle with `SERVICE_QUERY_CONFIG|SERVICE_QUERY_STATUS|SERVICE_START|SERVICE_STOP`,
  and `OpenService` evaluates the whole requested mask at once, so the call was
  refused over rights the query never used. What made this worth fixing rather
  than documenting is that the message was actively wrong: it reported that
  permission was withheld when in fact it was granted, sending the reader to look
  in exactly the wrong place.

  `status` now opens the service control manager with `SC_MANAGER_CONNECT` and
  the service with `SERVICE_QUERY_STATUS`, and nothing more. `start`, `stop`,
  `restart`, `install`, and `uninstall` still require elevation — the ACL
  withholds those rights deliberately, and that has not been relaxed. Output
  wording is unchanged on every platform, and the Linux and macOS paths are
  untouched.

### Added

- **Issue and pull-request templates** (#7). `.github/ISSUE_TEMPLATE/` now holds
  YAML forms rather than an empty box, with version, component, install method,
  OS, and elevation state all required. Each of those has already decided a
  diagnosis on this repository: #5 reproduces only via the MSI path, #6 turns
  entirely on whether the reporter is an administrator, and #3's version had to
  be reconstructed from its title. Blank issues are disabled.
- **`docs/cli.md`** — a user-facing reference for every command the binary
  exposes, with flags, examples, exit codes, and which service subcommands need
  elevation. Written from `internal/cli/`, not from the spec contract, which
  remains a contract.
- **`docs/INSTALL-linux.md` and `docs/INSTALL-macos.md`** — Windows was
  previously the only platform with a real guide. The macOS guide states plainly
  that the desktop bundle's auto-started daemon does **not** survive a reboot
  unless the service is registered, which was previously a README blockquote.
- **`docs/README.md`** — an index separating user-facing guides from maintainer
  material.
- **`CONTRIBUTING.md`, `SECURITY.md`, `CODE_OF_CONDUCT.md`.** `SECURITY.md`
  states the threat model the project actually holds rather than a generic
  posture: tasks run with the daemon's privileges, the Windows pipe admits
  Authenticated Users, the Unix socket admits any local user who can traverse
  the data directory, and release artifacts are unsigned.

### Changed

- **`README.md` rewritten** to the house Markdown style. A newcomer now reaches
  a running task without opening a spec artifact; the feature list became prose;
  the architecture section gained a diagram; "Project layout (target)" became the
  actual layout; and command-level questions route to `docs/cli.md` rather than
  into `specs/`.
- **`TODO.md` rewritten** to reflect delivered state. Every checkbox was
  unticked though the roadmap had shipped, and it still listed event triggers —
  removed outright in 0.4.0 along with their store migration — as pending work,
  which advertised a feature that does not exist.
- **`build/windows/verify_wxs.ps1`** asserts the `PATH` element and each of its
  attributes separately, so a partial edit (a per-user entry, or one that
  survives uninstall) is reported for what it is rather than passing.

### Decisions

- **2026-07-23** — **The `PATH` entry is declared on the CLI's own installer
  component** (`build/**`, pinned). MSI reference-counts by component, so binding
  the entry to `gosched.exe` gives correct install, upgrade, and uninstall
  behavior for free: written on install, replaced in place on a major upgrade,
  removed when the CLI is removed. A custom action editing `PATH` by hand would
  have to implement all three itself, and hand-rolled `PATH` editing is the
  classic source of duplicated and truncated `PATH` values.

  `System="yes"` matches the `perMachine` package scope; a per-user entry would
  be written for whoever ran the installer and stay invisible to everyone else on
  a machine that hosts a system-wide service. `Part="last"` appends, so an
  existing same-named tool keeps winning — the conservative choice for an
  installer running elevated. `Permanent="no"` is what removes it on uninstall.

  The rejected alternative was to document the full-path invocation everywhere.
  That is what the Windows guide did, and it is the reason the README and
  `docs/test-scripts.md` disagreed with reality.

- **2026-07-23** — **`docs/INSTALL-windows.md` now writes bare `gosched`
  commands** (pinned). This is correct only *because* of the decision above, and
  the two must move together. The guide keeps the full-path form, but demoted to
  what it actually is: the fallback for a shell that was already open when the
  installer ran, since the environment broadcast does not reach those. The
  troubleshooting section names that case first, because it is the one report
  this change will generate.

- **2026-07-23** — **The status fix lives in this repository rather than
  upstream.** The library helper responsible is shared with paths that
  legitimately need start and stop rights, so a correct upstream fix is a larger
  change than this defect needs — and a fork would sit on the critical path of
  every future dependency upgrade. Only the status path is reimplemented; every
  other action still goes through the library, and non-Windows platforms fall
  through to it unchanged rather than being reimplemented.

- **2026-07-23** — **The issue and PR templates are unpinned, and this was
  checked rather than assumed.** The pinned list in `CLAUDE.md` names
  `.github/workflows/**`, not `.github/**`, so `.github/ISSUE_TEMPLATE/**` and
  `.github/PULL_REQUEST_TEMPLATE.md` need no decision entry. Recorded here
  because the issue asked for the confirmation, and because "it lives under
  `.github/`" is exactly the inference someone would otherwise make.

### Known limitation

The `PATH` fix **cannot be verified from a development machine**, which already
has the install directory on `PATH` — precisely what hid the defect. Pre-release
evidence is the WiX sanity check and a read of the generated element; the
end-to-end check needs a clean machine and the released `.msi`. Issue #5 stays
open until that check passes.

## [0.5.3] - 2026-07-23

**`listeners` could present "the probe could not run" as "nothing is
listening".** Tooling only; program binaries remain identical to 0.5.0.

### Fixed

- **`listeners` read the newest snapshot unconditionally.** When that snapshot
  came from a host or twin where no port tool existed — or was taken with
  `-SkipPorts` — the query returned an empty table, which reads as a finding
  about the machine rather than an absence of data. Found while walking the
  quickstart: the POSIX twin under Git Bash has neither `ip` nor `ifconfig`, its
  snapshot landed with zero rows, and `listeners` went silent.

  This was the same class of error as the drift defect fixed in 0.5.1 —
  presenting a number, or an absence, without the provenance needed to read it.
  The schema now records **why** each list looks the way it does, and `listeners`
  reads the most recent snapshot whose port probe actually ran, naming the
  snapshot it chose and warning when the newest one was passed over.

- **The `listeners` diff compared against the previous snapshot regardless of
  whether that snapshot had port data.** A skipped snapshot in the middle of a
  series made every port on the next one read as `NEW`. The baseline is now the
  most recent *comparable* snapshot, and `no-comparable-snapshot` is reported
  when there isn't one.

### Added

- **`addresses_probe` and `ports_probe` on every snapshot**, recording `ok`
  (the probe ran; zero rows is a real answer), `unavailable` (no tool on this
  host could answer; zero rows means nothing), or `skipped`.

### Changed

- System schema version 3, via a **forward-only, non-destructive migration**:
  the two columns are added to existing databases, pre-existing rows are kept,
  and their probe status stays `NULL` rather than being back-filled with a
  guess. A snapshot recorded before this release has genuinely unknown
  provenance, and `listeners` flags it as such rather than assuming it was fine.

### Decisions

- **2026-07-23** — **The probe columns carry no `CHECK` constraint.** SQLite
  cannot add a constrained column to an existing table, so enforcing the
  vocabulary in the schema would leave fresh databases stricter than migrated
  ones — a divergence nobody would think to look for, and one that would only
  surface as a confusing write failure on some machines and not others. The
  writer enforces it instead, and the reader treats anything unexpected the same
  way it treats `NULL`: as unknown provenance.
- **2026-07-23** — **"No usable port data" exits 0 with a sentence, not 1.** The
  query ran and the honest answer is that nothing has been collected yet; that
  is a state of the data, not a failure of the run. What made the original bug
  harmful was silence, so the fix is words rather than an exit code.

## [0.5.2] - 2026-07-23

**Fixes a macOS-only defect introduced in 0.5.1.** Tooling only; program binaries
are identical to 0.5.0.

> **0.5.1 is broken on macOS.** Its POSIX twins cannot parse an `--anchor-iso`
> timestamp there and exit 2. The PowerShell twins are unaffected on every
> platform, and the POSIX twins are unaffected on Linux. macOS users should use
> 0.5.2.

### Fixed

- **`--anchor-iso` was unparseable on macOS.** The POSIX twins parsed timestamps
  with `date -d`, which is a GNU extension. macOS ships BSD `date`, where `-d` is
  not a parse flag at all, so every anchored query exited 2 with "not a parseable
  timestamp".

  Timestamp parsing now goes through one `parse_iso_epoch` helper that tries the
  GNU form first and falls back to BSD's explicit-format form, normalising the
  `+05:00` offset spelling to the `+0500` that BSD requires and dropping
  fractional seconds it will not accept.

  This could not reproduce locally: the development host is Windows with Git
  Bash, which ships GNU `date`, so the GNU path always succeeded. Only a real
  macOS runner exercises the other branch — which is exactly what caught it, on
  the CI run for the 0.5.1 tag. Regression coverage now exercises both the `Z`
  and numeric-offset spellings on every platform, because they take different
  branches of the fallback.

### Decisions

- **2026-07-23** — **0.5.1 was tagged and published before its CI run finished.**
  The local gates were green, but the local machine has no macOS and no C
  compiler, so two of the seven CI jobs had no local equivalent. Publishing on
  the strength of a partial signal is what put a known-broken artifact on the
  releases page. Tags are immutable once public, so the fix ships as 0.5.2 rather
  than a re-cut 0.5.1, and 0.5.1's release notes now say plainly that macOS users
  should skip it. For future releases: wait for CI to go green on the tag before
  treating a release as done, particularly when the change touches shell code that
  the development platform cannot fully exercise.

## [0.5.1] - 2026-07-23

**The drift measurement in 0.5.0 was wrong, and this fixes it.** Tooling only —
the daemon, CLI, GUI, and stored schema are untouched, so 0.5.0 and 0.5.1 ship
identical program binaries.

### Fixed

- **Dispatch drift reported a schedule's phase offset as though it were lateness.**
  0.5.0 accepted `-IntervalSeconds` alone and snapped each run's start to the nearest
  multiple of that interval *counted from the Unix epoch*. That is correct only when a
  schedule happens to sit on the epoch grid — and this scheduler anchors an interval
  schedule to the **task's creation time**, so a task created at `:06` fires at `:06`
  forever.

  Measured on a live daemon: drift of 6505 / 6262 / 6254 ms, apparently 64x over the
  project's 100 ms dispatch budget, while the same run's `cadence` query showed intervals
  of 59757–60006 ms. The scheduler was on time to within a quarter second; the 6.4 s was
  entirely the `:06` anchor. The figure was not merely imprecise, it was measuring a
  different quantity, and nothing in its presentation said so.

  Epoch snapping is removed. Drift now comes from a caller-supplied **anchor** — one real
  firing time from `gosched task show` — which reconstructs the whole `anchor + k x interval`
  grid. With no anchor, **no drift is recorded at all**: reporting nothing is better than
  reporting a confident wrong number, because nothing about a wrong number's presentation
  tells you which one you got. Verified after the fix on the same daemon: 259–312 ms, mean
  277 ms, against an independent `cadence` of 59949–59998 ms. The two agree.

### Added

- **`-AnchorIso` / `--anchor-iso` on `Test-ReadTestDB`**, which is now the primary path.
  The anchor cannot be known before the task exists — the scheduler derives an interval
  schedule's phase from the task's creation moment, so supplying it to the recorder is a
  chicken-and-egg problem. Drift is a derived quantity, so it is derived at read time from
  the raw start timestamps. This works on beats **already recorded**, and a wrong anchor is
  fixed by re-running the query rather than re-running the experiment.
- **`-AnchorIso` / `--anchor-iso` on `Test-Heartbeat`**, for the case where the firing grid
  genuinely is known in advance (a fixed-time schedule). Records `expected_source = 'anchor'`.
- **A `jitter` query**, for when no anchor is available. It derives the schedule's phase from
  the data and reports variation around it. The reader states on every run that jitter
  **cannot detect uniform lateness** — a scheduler consistently late by a fixed amount has
  zero jitter — because that limitation is the whole reason an anchor exists.

### Changed

- Heartbeat schema version 2. `expected_source` admits `anchor`; `boundary` remains
  readable so pre-0.5.1 databases still open, but is never written. The `drift` query flags
  any legacy `boundary` rows as phase offset rather than latency. Forward-only and
  non-destructive, per the constitution.

### Decisions

- **2026-07-23** — **Drift is derived at read time, not write time.** Three options were
  considered: keep epoch snapping and document the caveat (rejected — a caveat does not stop
  a wrong number being read as a right one); take the anchor at record time (rejected as the
  primary path — the anchor is unknowable until the task exists, and a wrong one is only
  fixable by discarding the data and starting over); derive at read time from raw
  timestamps. The third was chosen because the recorder already stores everything needed,
  the anchor is knowable by then, and the computation is re-runnable. The record-time option
  is retained as a secondary path for genuinely known grids.
- **2026-07-23** — **This defect was found by walking the quickstart end to end against a
  live daemon**, which was the one verification task left outstanding at the 0.5.0 halt. No
  unit test would have caught it: every unit test agreed with the implementation, because
  both shared the same wrong assumption about how schedules are anchored. The lesson is
  recorded in the spec's Clarifications section as a superseded decision rather than edited
  away, so the reasoning that produced the error stays visible next to its correction.

## [0.5.0] - 2026-07-23

Maintainer tooling and repository configuration only. The daemon, CLI, GUI, and
stored schema are untouched -- 0.4.1 and 0.5.0 ship identical program binaries.
The minor bump reflects new tracked tooling and two pinned-artifact changes, not
a behavior change in the scheduler.

### Added

- **Maintainer test scripts** (`test/scripts/`, documented in
  [`docs/test-scripts.md`](docs/test-scripts.md)). Three cross-platform script pairs — a
  PowerShell `.ps1` and a POSIX `.sh` twin each — that let a maintainer prove an installed
  `goschedd` actually fires on time, survives restarts, catches up after downtime, and honors
  its overlap policies. `Test-Heartbeat` records one beat per invocation into `heartbeat.db`
  with a measured dispatch drift; `Test-GetSystemInfo` records host snapshots into `system.db`;
  `Test-ReadTestDB` reads either back through eleven canned queries. `gosched runs` could say
  a task ran, but not how late it was, nor that a firing you expected never happened — those
  are the two questions this answers.
- **`.claude/skills/` is now tracked**, so a fresh clone arrives with the `/speckit-*` commands
  and the house-standard skills already present. `docs/build-autopilot.md` had named the
  missing-commands-on-a-fresh-clone problem as a setup failure; this closes it. Vendored:
  `shruggie-powershell`, `shruggie-markdown`, `shruggie-speckit`, `gh-fix-ci`, and a new
  project-native `go-schedule-verify` carrying the CI-parity procedure, its coverage-gate
  semantics, and both local-environment traps.

### Decisions

- **2026-07-23** — **Pinned artifact changed**: `.gitignore` moves from ignoring all of
  `.claude/` to `.claude/*` plus `!.claude/skills/`, and adds `test/scripts/.bin/`. Expressed as
  exclude-everything-then-narrowly-admit rather than a denylist, because the excluded material is
  credential-bearing by assumption and the two failure directions are not symmetric: a denylist
  admits every agent file nobody thought of, an allowlist admits only what was named. Verified
  before commit with `git status --porcelain -uall .claude`, which listed skills and nothing else.
- **2026-07-23** — **Pinned artifact changed**: `.gitattributes` gains an LF exemption for
  `test/scripts/**/*.ps1` and `.claude/skills/**/*.ps1`. The existing `*.ps1 text eol=crlf` rule
  is justified in-file as
  "Windows-only scripts keep CRLF", but these particular `.ps1` files are cross-platform by
  design — they run under `pwsh` 7 on Linux and macOS — so that rationale does not reach them,
  and the ShruggieTech compliance checker they are authored against requires LF. Scoped as
  narrowly as possible rather than flipping the global rule. The skills path is included for a
  second-order reason found while staging: the vendored `shruggie-powershell` skill ships the very
  checker that enforces LF, so storing its own scripts and examples as CRLF would have made them
  fail their own compliance check on a fresh clone.
- **2026-07-23** — **Dispatch drift is derived, not reported, and every figure carries its
  source.** Inspecting `internal/executor/executor.go` established that a spawned task receives
  the inherited environment plus its own configured variables and nothing scheduler-generated —
  no scheduled time, no run ID. Three options: infer drift from the observed cadence; change the
  executor to inject the scheduled moment; or snap the run's start to the nearest boundary of a
  caller-declared interval. Cadence inference was rejected because it measures *jitter* — a
  scheduler uniformly five seconds late scores perfectly, and that is the defect class this most
  needs to catch. Modifying the executor was rejected because it changes a safety-critical
  product surface for maintainer tooling's benefit and would forfeit this release's provable
  "the shipped binaries did not change" property. Boundary snapping yields genuine absolute
  latency for a wall-clock-aligned schedule, with an `env` tier kept ahead of it so a future
  release that does export the scheduled moment is consumed with no change. Every drift value
  records which of the three sources produced it, and the reader refuses to pool them.
- **2026-07-23** — **The scripts bind SQL parameters via `sqlite3`'s `.param set`**, which sets
  the 3.33.0 minimum version (with `.mode json`). The values written include hostnames,
  usernames, and interface names: string-interpolated SQL there is both an injection vector on a
  machine someone else administers and an ordinary bug for any user named `O'Brien`.
- **2026-07-23** — **No product code, CI workflow, or retention policy changed.** The daemon,
  CLI, GUI, and stored schema are untouched, so 0.4.1 and this release ship identical binaries.
  The new tests run inside the existing `go test ./...` invocation, so no workflow edit was
  needed. The test databases are never pruned or rotated: deleting the file is the documented
  reset, and automatic retention would silently destroy the history a maintainer is inspecting.

## [0.4.1] - 2026-07-23

Release-packaging fixes only. No change to the scheduler, the GUI, the CLI, or
the stored data — 0.4.0 and 0.4.1 are the same program.

### Fixed

- **`SHA256SUMS.txt` now covers every published asset.** It was generated in the
  job that builds the daemon and CLI tarballs, which cannot see the artifacts built
  by the GUI job, so the Windows `.msi` and the desktop bundles — the files most
  people actually download — were never checksummed. A final job now runs after all
  the others, downloads every attached asset, and publishes one complete checksum
  file.
- **The Windows `.wixpdb` is no longer published.** `wix build` writes a debug-symbol
  file next to the `.msi`, and the release step attached everything in `dist/` with a
  bare glob. The publish patterns are now explicit. (Present in 0.3.0 and 0.4.0;
  harmless, but not something anyone should download.)

### Decisions

- **2026-07-23** — **Pinned artifact changed**: `.github/workflows/release.yml` gains a third
  job, `checksums`, and both publish steps now name their file patterns explicitly instead of
  globbing `dist/*`. Pinned artifacts change only with a dated decision, hence this entry.
  Checksums move to a job gated on `needs: [binaries, gui]` because the completeness problem is
  structural, not a missing filename: the job that wrote `SHA256SUMS.txt` runs before the GUI
  artifacts exist and on a different runner, so no edit to it could ever cover them. The
  alternative — one checksum file per job — was rejected as it pushes the reassembly onto whoever
  is verifying a download. The new job is idempotent on re-run (it discards any prior checksum
  file before recomputing) and writes to a temp path so a failed run cannot leave a truncated
  file over a good one.

## [0.4.0] - 2026-07-23

**Groups work from the GUI, and the task editor tells the truth about a task's
schedule.** The two defects reported against 0.3.0 are fixed
([#3](https://github.com/shruggietech/go-schedule/issues/3),
[#4](https://github.com/shruggietech/go-schedule/issues/4)), and group
assignment is reachable without the command line for the first time.

Upgrading is a normal install; the store migrates forward automatically. Note
that a pre-rebrand `goscheduler` data directory is no longer picked up — see
**Removed**.

### Fixed

- **Task editor showed the wrong schedule** ([#4](https://github.com/shruggietech/go-schedule/issues/4)):
  opening a task for editing always displayed Mode as *Recurring* with the Schedule and one-off
  date/time fields blank, regardless of how the task was actually scheduled. The dialog now fetches
  the task's schedule and shows its real mode, its schedule phrase, or its one-off date and time in
  the task's own timezone. Saving an untouched dialog leaves the schedule byte-identical.
  Switching Mode now requires the new mode's timing, closing a hole where an empty date/time
  silently kept a recurring schedule on a task the user believed was one-off. Changing only a
  task's timezone now re-interprets its recurrence in the new zone.
- **Groups were unusable from the GUI** ([#3](https://github.com/shruggietech/go-schedule/issues/3)):
  there was no way to put a task into a group without the CLI, and no way at all — from any client —
  to take one back out, because an empty group value meant "leave unchanged". The task editor now
  has a Group field (including `(none)`), the Groups tab shows each group's member tasks plus an
  always-present **Ungrouped** area and a **Move to group…** action, and the task list shows each
  task's group. `gosched task edit --group ""` un-groups a task; omitting `--group` still leaves
  membership unchanged.

### Added

- **Build-Phase Autopilot Protocol** (`docs/build-autopilot.md`): the operating procedure for
  running a spec-kit feature end to end on one verbal kickoff, with the routine decisions made
  and recorded by the agent and exactly one halt before anything is pushed. Constitution
  principle V (**v1.1.0**) is the governing law; `CLAUDE.md` carries the standing authorization,
  the CI-parity verification commands, and the non-negotiable safety-critical test surfaces.

### Changed

- **Development is trunk-based; the pull-request requirement is gone**
  (**constitution v2.0.0**). Work is committed directly onto `main` — no feature branches, no
  pull requests. The old requirement ("every change lands via pull request; no direct pushes to
  the default branch") never described how this project actually works: it has one-to-two
  developers, has never used pull requests for review, and a PR with no reviewer adds latency
  without adding scrutiny. Nothing is relaxed. The single pre-push halt is retained and becomes
  the sole human review point; deviations from a principle are recorded in the commit message
  rather than a PR description; and the local CI-parity requirement is *strengthened*, because
  CI now reports after a push to `main` instead of blocking a merge — a red local run is a halt,
  not something to push and sort out afterwards. `.github/workflows/ci.yml` needed no change: it
  already triggers on push to `main`. Mirrored in `CLAUDE.md` and `docs/build-autopilot.md`.

### Removed

- **The pre-rebrand data-directory migration** (`config.MigrateLegacyPaths`, added in 0.3.0):
  the daemon no longer moves a `goscheduler` data directory onto the `goschedule` name at
  startup. Nothing on disk is deleted — an existing `goscheduler` directory is simply left
  alone and ignored, and the daemon creates a fresh `goschedule` beside it.

### Fixed (CI)

- **The coverage gate could fail for code that no longer exists.** `.github/workflows/ci.yml`
  measured core-package coverage with `go test -coverpkg=<six packages> ./...` and no
  `-count=1`. Under `-coverpkg` every test binary is instrumented for all six target packages,
  so a cached test result replays a coverage profile enumerating the file set as it stood when
  that result was cached. Packages whose own sources are unchanged are served from the cache
  (`actions/setup-go` restores it via `cache: true`) and drag stale blocks — including blocks
  belonging to deleted files — into the merged profile. Deleting a well-covered file therefore
  left its statements in the denominator with nothing covering them. Observed on the first push
  after `internal/schedule/render.go` was removed: `schedule` reported 50.5% against an 80%
  gate, exactly `168/333` where 333 is the current 191 statements plus the deleted file's 142.
  Adding `-count=1` to that step fixes it.

### Decisions

- **2026-07-22** — Store migration **v4** adds `schedules.expression`, retaining the human-readable
  phrase a recurring schedule was parsed from. Forward-only and non-destructive: one column with a
  total default, no existing value read or rewritten, so no stored timing moves. The phrase is
  inert with respect to execution — `RRULE` remains the only input the engine evaluates — and
  exists solely so a client can show the user their own wording again. Pinned by an explicit
  upgrade test asserting a v3 database migrates with every schedule row otherwise unchanged and
  re-opens as a no-op.
- **2026-07-22** — **Pinned artifact changed**: the coverage gate moves out of
  `.github/workflows/ci.yml` into `scripts/coverage-gate.sh`, and CI now invokes that script.
  Previously the gate existed only as inline Python in the workflow, so there was no way to run
  it locally without transcribing it — which is exactly how a push went out that CI then
  rejected: the local check used `go test -cover` (per-package) while the gate measures
  cross-package coverage with `-coverpkg`, two different numbers. One implementation removes the
  drift and makes the gate a first-class CI-parity command in `CLAUDE.md`. Written in POSIX `sh`
  + `awk` rather than Python so it runs unchanged in Git Bash on Windows, in WSL, and on the
  runner; the previous inline version required `python3`, which is absent on a stock Windows
  workstation. Threshold, package list, and aggregation semantics are unchanged, and the awk
  aggregation was verified to reproduce the Python output exactly. Both the pass path (exit 0)
  and the fail path (exit 1 at a raised threshold) were exercised on Windows and Linux.
- **2026-07-22** — **Pinned artifact changed**: `.github/workflows/ci.yml` gains `-count=1` on the
  coverage-gate command. Pinned artifacts change only with a dated decision, hence this entry. The
  gate was measuring a denominator that included deleted files, because `-coverpkg` plus Go's test
  cache replays stale coverage profiles from packages whose own sources did not change. This is a
  correctness fix to the measurement, not a relaxation: the 80% threshold, the six core packages,
  and the aggregation script are all unchanged, and the gate now measures the tree as it actually
  is. Verified by reproducing the gate locally on both Windows and Linux, which agree at
  `schedule` 88.0% / `store` 86.8%.
- **2026-07-22** — The pre-rebrand path migration is removed for the same reason as the schedule
  renderer below: it carries data forward from an installed base that does not exist. Unlike the
  renderer it was not merely inert. Inspecting the one machine where it would still fire found
  `C:\ProgramData\goscheduler` holding a `schema_version = 2` database — one *disabled* task, 24
  runs of which **all 24 failed**, 24 `run_failed` alerts, and no groups, spanning 45 minutes on
  2026-06-20. Keeping the migration would rename that directory onto the new name and run store
  migrations v3 and v4 over it, importing a broken database into an otherwise clean install.
  Removing it is non-destructive: the legacy directory is left untouched for manual recovery or
  deletion, and the daemon starts fresh.
- **2026-07-22** — Nothing reconstructs schedule phrases for rows stored before the `expression`
  column existed. An earlier revision of this work added `schedule.Render`, an RRULE→phrase
  inverse applied at read time, so already-installed databases would also show their schedule on
  edit. That was built on a wrong premise — the defects were filed against v0.3.0 and the design
  inferred an installed base to protect. There is none: the software has no working deployments
  and the only databases in existence are the maintainers' own, none of them functional. The
  renderer and its round-trip test suite served exclusively that phantom population and were
  removed. `schedule.Parse` is the only producer of recurring schedules, so every schedule created
  from here on retains its phrase; a database predating the column shows a blank schedule field on
  edit, which means "keep unchanged" and is harmless. Migration v4 is kept — it is what creates the
  column, and folding it into the v1 `CREATE TABLE` would leave existing databases at
  `schema_version = 3` with the column silently absent, failing every schedule query.
- **2026-07-22** — `TaskUpdateRequest.GroupID` becomes `*string` so group membership can carry
  three intents: nil leaves it unchanged, `""` removes the task from its group, and an id assigns
  it. Previously `""` meant "unchanged" and un-grouping was unreachable from every client. This
  reuses the convention already set by `GroupUpdateRequest.Parent` rather than introducing a
  sentinel value that could collide with a real group id. Wire-compatible: omission still means
  unchanged, and the CLI preserves that by only sending the field when `--group` is passed.
- **2026-07-22** — ~~Autopilot halts before the *branch push and pull request*, not before a push
  to `main`. The constitution forbids direct pushes to the default branch, so the halt is placed
  at the last point before work leaves the machine. This diverges deliberately from the
  trunk-based variant of the protocol used in other projects.~~ **Superseded the same day** by
  the constitution v2.0.0 amendment below: the project is trunk-based and the halt precedes the
  push to `main`.
- **2026-07-22** — Autopilot's standing scope is features traceable to
  `specs/001-task-scheduler/spec.md` and the `TODO.md` roadmap. This project has no separate
  build-sequence document, so the master spec plus the roadmap serve that role. Any other work
  can still be placed under autopilot by explicit operator request, which is itself the renewal.
- **2026-07-22** — The safety-critical test surfaces that autopilot may never weaken are named
  explicitly for this project: clock injection, timezone/DST resolution, forward-only store
  migrations, restart and catch-up recovery, goroutine termination under the race detector, and
  local IPC access control. Autopilot grants autonomy of execution only and relaxes no quality
  gate.
- **2026-07-22** — `.claude/` stays gitignored (the agent folder may hold credentials). The
  `/speckit-*` command skills the protocol drives are therefore per-clone local state, restored
  with `specify integration upgrade claude`; this is recorded as a precondition in the protocol
  rather than by tracking the folder.

## [0.3.0] - 2026-06-21

### Changed

- **Rebranded `go-scheduler` → `go-schedule`** (`specs/004-rebrand-gui-overhaul/`): module path,
  build/release config, user-facing strings, and on-disk identity (data dir `goschedule`, DB
  `goschedule.db`, logs under `goschedule/logs/`). The daemon performs a best-effort one-time move
  of a pre-rebrand `goscheduler` data directory on startup (non-fatal; never deletes data).
- **Windows is now distributed as a formal `.msi`** built with WiX v5
  (`build/windows/goschedule.wxs`): installs to Program Files, registers `goschedd` as an
  auto-start Windows service, adds a Start-Menu shortcut, and uninstalls cleanly (user data under
  `C:\ProgramData\goschedule` is preserved). The portable Windows zip and "run the exe" flow are
  removed; the Windows install guide was rewritten.
- **GUI "Alerts" replaced by a unified "Logs" view**: a new `internal/logbus` slog handler tees
  every daemon log record to a rotating on-disk JSONL file (`logs/goschedule.log`), a bounded
  in-memory ring (served by `GET /v1/logs`), and the live event stream. The view merges daemon
  logs and scheduler alerts, with severity filters, click-through detail, and "Dismiss All". A new
  `gosched logs` CLI command mirrors it (`alerts` is deprecated).
- **GUI updates in real time across all views**: the event broker now also publishes task/group
  change events from the API mutation handlers, the view-model folds them, and the GUI
  re-synchronizes on stream reconnect. All manual **Refresh** controls were removed.

### Added

- **Calendar view under Schedule**: a toggleable month-grid view over the existing calendar API,
  alongside the agenda list; the selected window is preserved across toggles and it updates live.

### Removed

- **Event Triggers feature removed entirely** (GUI tab, CLI commands, API routes/client, engine
  dispatcher, store tables, and domain types). Store **migration v3** drops the `triggers` and
  `dedup_ledger` tables (a no-op on databases that never had them).

### Added (earlier)

- Spec-driven development scaffolding via Spec Kit:
  - Project constitution (v1.0.0) — code quality, testing standards, UX consistency, performance.
  - Feature specification for the cross-platform task scheduler (`specs/001-task-scheduler/`),
    including clarifications and a one-off (non-recurring) scheduling mode.
  - Implementation plan, research, data model, CLI & local-API contracts, and quickstart.
  - Dependency-ordered task breakdown (78 tasks across 8 phases).
- Repository basics: Apache 2.0 license, README, changelog, and TODO.
- **Foundational implementation (Phases 1–2, tasks T001–T019):**
  - Go module, `golangci-lint` config, `Makefile`, and `.gitattributes`.
  - `internal/platform` — build-tagged data dirs and windowless process-spawn helper.
  - `internal/clock` — injectable `Clock` with real and deterministic fake implementations.
  - `internal/config` — single config schema, fail-fast validation, structured `slog` logging.
  - `internal/domain` + `internal/store` — core entities and durable SQLite persistence
    (pure-Go, cgo-free) with migrations and CRUD.
  - `internal/ipc` — local transport (Unix socket / Windows named pipe).
  - `internal/api` — local HTTP/JSON API server (health, error envelope) and shared client.
  - `cmd/goschedd` (daemon) and `cmd/gosched` (CLI): the daemon serves health over IPC and the
    CLI reaches it — end-to-end architecture verified.
- **User Story 1 — MVP (Phase 3, tasks T020–T037, T074–T078):**
  - `internal/timezone` — IANA resolution and DST rules (next-valid spring-forward,
    first-occurrence fall-back), verified against 2026 US transitions.
  - `internal/schedule` — RFC 5545 RRULE recurrence (rrule-go), one-off, and a human-readable
    parser with plain-language summaries (no cron syntax); cron-parity suite.
  - `internal/engine` — timer-driven scheduling loop over an injected clock, bounded worker
    pool, one-off completion, failure alerts; overlap policies (queue_one / skip /
    allow_concurrent) with warning + alert.
  - `internal/executor` — windowless command execution with bounded output capture; build-tagged
    `run_as` (Unix credential impersonation; rejected on Windows for now).
  - Local API: task CRUD + edit (PATCH), `schedules/preview`, `run-now`, enable/disable, and
    run/alert queries. Full cobra CLI: `task`, `runs`, `alerts`, `service`, `gui`, with `--json`
    and contract-compliant exit codes.
  - `internal/service` — cross-platform system-service control (install/start/stop/status) via
    kardianos; the daemon runs under the OS service manager (start on boot).
  - Verified end-to-end: create recurring + one-off tasks via CLI, run them, inspect history and
    failure alerts; DST handled correctly across the year.
- **User Story 3 — Nested task groups (Phase 5, tasks T049–T054):**
  - `internal/task` — pure, testable group-tree logic: cascading enabled-state resolution,
    descendant enumeration, cycle detection, forest building.
  - `internal/store` — group chain-enabled queries, parent validation, reparent with cycle
    rejection, rename, and tree retrieval.
  - Engine respects the group chain: disabling an ancestor group stops its tasks from being
    scheduled (without mutating each task's own enabled flag); re-enabling restores them.
  - Local API: group CRUD, tree view, reparent (PATCH), enable/disable. CLI: `group add/list
    [--tree]/enable/disable/rm`.
  - Verified end-to-end: 3-level hierarchy, cascade disable, cycle rejection.
  - Note: the GUI group tree (T055) is deferred until the US2 GUI exists.
- **User Story 4 — Event triggers (Phase 6, tasks T056–T061):**
  - `internal/trigger` — completion-event dispatcher: matches a source task's
    success/failure/any outcome to triggers and fires target tasks, with durable
    de-duplication (window + key) and at-least-once recovery across restarts.
  - `internal/store` — triggers and a dedup ledger (claim/mark-executed/pending),
    schema migration v2.
  - Engine wiring: a completion hook fires triggers after each run; a startup hook
    recovers unexecuted events. New `FireEvent` dispatches targets as event runs.
  - Local API: trigger CRUD; CLI: `trigger add/list/rm`.
  - Verified end-to-end: source completion fires the target once (recorded as an
    `event` run); duplicates within the window are de-duplicated.
  - Note: the GUI trigger editor field (T062) is deferred until the US2 GUI exists.
- **User Story 5 — Downtime catch-up (Phase 7, tasks T063–T066):**
  - `internal/catchup` — pure detection: given a task's schedule, last run, and
    policy, decide whether a scheduled run was missed during downtime.
  - Engine startup performs at most one catch-up run per eligible task (recorded
    as a `catchup` trigger at startup time, so a restart never re-triggers it),
    raises a `missed_run` alert, then resumes normal scheduling. Honors the
    per-task catch-up policy (`one` / `none`) and the overlap policy via dispatch.
  - Verified end-to-end: a short-interval task left across real downtime performs
    exactly one catch-up run and then resumes.
- **Polish & hardening (Phase 8, tasks T067–T071; T072/T073 partial):**
  - `internal/lock` — cross-platform single-instance guard (flock / LockFileEx); a
    second daemon now fails fast instead of double-executing every task (T070).
  - Goroutine-leak test (no leak after 500 executions) and a dispatch benchmark
    (~36µs per run — far under the 100ms budget) (T068, T069).
  - Test coverage raised to ≥80% statements on all core packages — engine, schedule,
    timezone, store, trigger, catchup (T071).
  - README updated to reflect functional CLI/daemon; daemon + CLI cross-compile
    cleanly for linux/macOS/windows on amd64 + arm64 (T067, T072 partial).
  - Deferred (need the US2 GUI): windowless-GUI verification (T072) and the GUI
    success criterion SC-008 (T073). Other success criteria verified via live CLI
    tests.
- **User Story 2 — Material Design desktop GUI (Phase 4, tasks T038–T048, T055, T062):**
  - `gui/` — Fyne desktop app with tabs for Tasks, Schedule (calendar/timeline),
    Groups (tree), Triggers, and Alerts, using Fyne's Material-style theme. The
    guided task editor shows a live plain-language schedule preview (FR-006); the
    alerts panel updates live and carries an unacknowledged badge.
  - `internal/events` — in-process pub/sub broker; API `GET /v1/events` streams
    run/alert events over SSE and `GET /v1/calendar` materializes occurrences.
  - `gui/viewmodel` — pure, unit-tested GUI state; the Fyne widget layer is
    cgo-free and unit-tested headlessly. Only `cmd/gosched-gui` (the OpenGL
    window) needs cgo; a cgo-free stub keeps `go build ./...` working everywhere.
  - CI builds the GUI with cgo + OpenGL and runs the headless GUI tests; releases
    publish `gosched-gui` for Linux, macOS, and Windows (windowless on Windows).
- **Zero-config desktop experience:**
  - `internal/autostart` — the GUI now starts the background daemon automatically
    (detached, windowless) if none is reachable, and reuses an already-running one
    (e.g. the installed service); the daemon's single-instance lock prevents
    duplicates.
  - Releases now publish a self-contained `go-scheduler-desktop_<os>_<arch>`
    archive bundling the GUI + daemon + CLI, so desktop users download one file and
    just run the GUI.

[Unreleased]: https://github.com/shruggietech/go-schedule/compare/v0.8.0...HEAD
[0.8.0]: https://github.com/shruggietech/go-schedule/compare/v0.7.0...v0.8.0
[0.7.0]: https://github.com/shruggietech/go-schedule/compare/v0.6.0...v0.7.0
[0.6.0]: https://github.com/shruggietech/go-schedule/compare/v0.5.3...v0.6.0
[0.5.3]: https://github.com/shruggietech/go-schedule/compare/v0.5.2...v0.5.3
[0.5.2]: https://github.com/shruggietech/go-schedule/compare/v0.5.1...v0.5.2
[0.5.1]: https://github.com/shruggietech/go-schedule/compare/v0.5.0...v0.5.1
[0.5.0]: https://github.com/shruggietech/go-schedule/compare/v0.4.1...v0.5.0
[0.4.1]: https://github.com/shruggietech/go-schedule/compare/v0.4.0...v0.4.1
[0.4.0]: https://github.com/shruggietech/go-schedule/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/shruggietech/go-schedule/releases/tag/v0.3.0
