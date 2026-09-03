# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Desktop interaction and appearance are now consistently readable and
  responsive (Refs #101, #104, #105, #106, #109, #111).** S044 makes System the
  clean-profile font default, adds packaged Inter and Ubuntu choices alongside
  Geist and Monospace, and ensures appearance menus offer only unapplied
  alternatives. Application storage is presented as compact aligned rows,
  unavailable rows are muted as a whole, the navigation rail has a full-height
  boundary and semantic Exit treatment, and a persisted 1x-to-4x wheel setting
  controls every application-owned vertical view while preserving precision
  deltas. Shared translucent interaction overlays keep text and glyphs readable
  across ordinary, selected, primary, and semantic button states in both
  palettes.

- **Windows release readiness now has a pre-publication demo qualification step
  (Refs #94, #98, #96, #101, #104, #105, #106).** S043 produces a distinctly
  labeled local MSI, exhausts non-attended verification, and binds operator
  testing to its exact hash before the review branch is published. Passing demo
  evidence reduces review risk but does not replace the later byte-identical
  draft-release candidate gate. Compiled inspection now labels these artifacts
  as `local-demo` instead of borrowing formal candidate provenance. The attended
  walkthrough accepted setup, shortcuts, Finish actions, task double-click,
  Exit, guided uninstall, and wipe behavior; its UI findings remain tracked in
  #101, #104, #105, #106, #109, #110, #111, #112, and #113 instead of
  expanding the qualification slice.

- **The desktop GUI now has an Options view and an owned navigation rail (Refs
  #104, #105, #106).** Users can switch among Dark, Light, and Follow system
  modes plus Brand, System, and Monospace fonts, restore the Dark and Brand
  defaults, and inspect or copy a read-only inventory of machine, per-user,
  runtime, installation, and Windows maintenance paths. Daemon paths reflect
  the connected daemon's effective configuration, while lifecycle text avoids
  claiming Windows-only wipe behavior on other platforms. The wider rail keeps
  Options above Info and a separate one-shot Exit command at the bottom-right.

- **Task rows now open Edit on double-click (Refs #103).** Single-click and
  keyboard selection remain available, while stable task identities survive
  list reorder, stale rows fail closed, detail-fetch fallback is retained, and
  one editor guard prevents stacked dialogs.

- **Windows releases now have an exact-candidate attended gate (Refs #94,
  #98, #96).** Version tags stage draft assets and a manifest instead of
  publishing immediately. A shared validator binds native Windows 11 evidence
  to the MSI hash, ProductCode, source commit, staging run, required display
  and account contexts, all core behavior scenarios, and hashed attachments.
  A separate promotion workflow publishes only the same tested MSI after the
  complete evidence archive passes.

- **Windows setup now exposes user-controlled shell integration and completion
  actions (Refs #97).** Start Menu and desktop shortcuts are independent MSI
  features with safe defaults, maintenance support, and upgrade migration.
  Fresh attended setup offers independent options to launch the unelevated GUI
  and open project documentation; unattended and non-fresh flows launch
  neither.

- **Windows uninstall now offers preserve-by-default or explicitly confirmed
  application-data cleanup (Refs #98).** Managed removal uses one exact wipe
  property. The installer-private helper derives only declared ProgramData and
  registered local-profile preference roots, refuses unsafe or redirected
  targets before deletion, avoids reparse traversal, and retains protected
  evidence for incomplete cleanup. Shared `goschedadmin` state is preserved.

### Fixed

- **Info version and attribution labels now use the unwrapped body-text layout
  path (Refs #101).** Their text, centered alignment, and dynamic build version
  remain unchanged. Headless tests pin the construction invariant; #94 retains
  exact-candidate Windows proof at standard and scaled DPI.

- **Windows Settings can no longer bypass the guided preserve-or-wipe
  uninstall flow (Refs #98).** The MSI registers maintenance as its supported
  attended application-management action and suppresses the reduced-interface
  direct Remove entry. Users select Modify, then Remove in the full wizard;
  direct silent preserve and explicit-wipe commands remain supported.

- **Installed Windows access now works for direct `goschedadmin` users whose
  standard token omits the new alias SID (Closes #90).** The restricted pipe
  retains SYSTEM, Built-in Administrators, and configured-group ACEs while
  adding validated, deterministic direct-user SID ACEs at service startup.
  Authenticated Users and unrelated accounts remain excluded, and a native
  standard-token regression proves the formerly denied connection.

- **Installed Windows task execution now has a real service-boundary contract
  (Closes #93).** The lifecycle probe creates tasks through the installed CLI
  and requires manual and scheduled LocalSystem execution, exit code 0,
  expected output, and marker effects. Nonzero child exits remain distinct from
  process-start failures, whose diagnostics name only the executable and OS
  error without disclosing arguments, stdin, or environment values.

- **The Windows desktop now recovers from named-pipe access denial (Refs
  #90).** Shared client errors distinguish daemon absence, access denial,
  timeout, other transport failures, and API responses. Concurrent model,
  calendar, and event-stream failures coalesce into one reachable in-frame
  state with Retry and Exit, while one cancelable reconnect loop uses bounded
  backoff and clears the incident after connectivity returns.

- **The desktop GUI now opens in a bounded restored window (Closes #89).**
  Fresh launches prefer 1280x800 logical units and cap each dimension to 90
  percent of the selected monitor's work area after display scaling. The
  work area is a hard bound on small displays, top-level actions remain
  reachable at effective 800x600, and user-controlled window states remain
  unchanged.

- **The Windows MSI now exposes concise Explorer Subject copy (Closes #88).**
  Summary Information PID 3 is `go-schedule: cross-platform task scheduler`.
  Source verification prevents drift, and candidate-MSI inspection records the
  artifact hash, exact Summary value, unchanged identity fields, and the native
  Shell `System.Subject` consumed by Explorer.

### Decisions

- **2026-09-03 - Separate pre-PR demo testing from formal candidate evidence.**
  The maintainer requires attended testing before a pull request, while #94
  requires a Release-workflow MSI staged from reviewed and merged source. S043
  therefore uses a clearly marked local demo for early testing and preserves the
  full exact-candidate matrix for the later authorized tag. Relabeling or
  promoting the local build would weaken provenance and is prohibited. The
  compiled inspector accepts `local-demo` for this narrow purpose while keeping
  candidate-manifest generation and published-origin validation unchanged.

- **2026-09-03 - Own desktop navigation and bound appearance choices.** Fyne's
  leading AppTabs cannot provide both controlled rail spacing and a bottom
  command without treating Exit as selectable content. S042 replaces that shell
  with seven ordinary destinations and one separate command. Appearance remains
  a small validated per-user preference over bundled or framework fonts; no
  arbitrary font parser, filesystem scan, path mutation, or daemon setting is
  introduced. Native text sharpness and scaled-DPI placement remain release
  candidate observations under #94.

- **2026-09-02 - Route attended Windows removal through MSI maintenance.** The
  system's direct MSI Remove entry can use reduced UI and skip package-authored
  data choices, as the maintainer observed against the unreleased S039 MSI.
  S041 therefore uses Windows Installer's supported `ARPNOREMOVE` model, an
  MSI-owned current-ProductCode `/I` registry value, and a package-owned
  maintenance page. Hosted evidence proved that `ARPNOREMOVE` also omits the
  generated `ModifyPath`, while WiX's stock page disables its Remove control.
  A shadow ARP entry, bootstrapper, or execute-sequence popup would duplicate
  product ownership or contaminate silent administration.

- **2026-09-02 - Build Windows release assets once, then promote the tested
  draft.** Rebuilding after attended testing cannot guarantee byte identity,
  while immediate tag publication leaves no interval for exact-MSI evidence.
  S040 therefore makes the tag workflow draft-only, moves final all-asset
  checksums into a fail-closed manual promotion, and retains the evidence ZIP
  beside the published assets. Tooling completion does not close #94, #98, or
  #96 until a real candidate passes every attended requirement.

- **2026-09-02 - Distinguish native window geometry from literal Fyne canvas
  size.** The attended collector records HWND, client, monitor, work-area, DPI,
  and state measurements through Win32. A narrow opt-in GUI evidence file
  records Fyne canvas size and scale from the installed process because an
  external HWND probe cannot truthfully supply that toolkit-level value.

- **2026-09-02 - Model Windows setup choices as native MSI features and one
  package-owned dialog flow.** S039 replaces the minimal dialog set because it
  cannot expose optional features or two independent completion actions. Stable
  child features give install, maintenance, managed deployment, and later
  upgrades one consistent shortcut state model. Completion actions exist only
  as guarded full-UI Finish events and use the interactive user's shell.

- **2026-09-02 - Delay explicit data wipe until uninstall commit and preserve
  uncertain security state.** A windowless helper embedded in the MSI derives
  fixed product-owned roots, preflights every root before deleting any, and
  records residual outcomes in a protected ledger. Its ignored commit return
  prevents Windows Installer from pretending irreversible partial deletion was
  rolled back. `goschedadmin`, exports, redirects, and unregistered or detached
  profile storage remain outside automatic cleanup.

- **2026-09-02 - Pre-authorize S039 review-branch publication only.** The
  operator explicitly directed autopilot to push this verified branch and open
  its pull request without the usual intermediate halt. This narrows the
  process pause for S039 but does not authorize merge, tag, release, scope
  expansion, or reduced verification.

- **2026-09-01 - Expand only verified direct users at Windows pipe startup.**
  S038 intentionally supersedes S036's ACL exclusion because the presentation
  repair left ordinary installed access broken. Netapi32 enumerates direct
  local-group members once, the daemon validates and de-duplicates user SIDs,
  and the configured-group ACE remains for fresh and nested membership tokens.
  Enumeration failures fail closed; access is never broadened to Authenticated
  Users or Everyone.

- **2026-09-01 - Require real service-hosted execution evidence for the core
  Windows promise.** In-process executor tests and the scheduler's recording
  runner cannot prove LocalSystem process creation. The pinned Windows
  lifecycle probe now uses an absolute inbox executable, manual and scheduled
  markers, run history, and controlled failure cases. The pinned CI workflow
  also runs a disposable binary-level LocalSystem probe and uploads its JSON
  evidence. The broader candidate-MSI release gate remains tracked separately
  by #94.

- **2026-08-31 - Preserve restricted IPC authorization while diagnosing the
  current Windows session.** S036 reads service, local-group, account-member,
  and process-token evidence only to select accurate recovery guidance. It does
  not mutate membership, elevate the GUI, or change named-pipe ACL identities.
  The pinned Windows installation guide now records the diagnostic and
  sign-out/sign-in walkthrough required by issue #90.

- **2026-08-31 - Replace specification 003's full-work-area startup with a
  DPI-aware restored-window contract.** Before a native window has stable
  placement, Windows monitor discovery uses the monitor nearest the launch
  pointer. The preferred 1280x800 logical size is capped independently to 90
  percent of that monitor's work area, without the former 640x480 enlargement
  that could exceed a constrained display. Saved placement and task-editor
  dialog sizing remain out of scope.

- **2026-08-31 - Pin the Windows installer Subject at both source and compiled
  artifact boundaries.** The pinned WiX source and verifier under `build/`
  approve one colon-separated Subject, while the pinned release workflow runs
  the Windows artifact inspector after WiX compilation and before attachment.
  This changes no installer product identity, UpgradeCode, filename, signing,
  lifecycle behavior, or other platform output.

- **2026-09-01 - Separate native token diagnosis from deterministic recovery
  verification.** Issue #90 no longer requires an operator to terminate an
  active Windows session solely to refresh a token during development. The
  installed host proves real service, group, account, token, restricted-pipe,
  and guidance behavior; an authorization-transition regression proves Retry
  restores the complete GUI without weakening the IPC boundary.

## [0.9.1] - 2026-08-31

### Fixed

- **Windows MSI group provisioning now uses WiX's elevated local-group path
  (Refs #83).** The installer no longer qualifies `goschedadmin` with the
  computer name, which WiX 6 misclassified as an impersonated domain operation
  and failed with `0x80070005`. Source and compiled-MSI contracts reject any
  non-empty group domain, and the Windows lifecycle verifier now records repair,
  reinstall, v0.9.0 upgrade, uninstall, group/member identity, service ordering,
  exit codes, verbose logs, and post-refresh non-elevated client access.

- **Release-time README synchronization now respects reviewed changes.** The
  v0.9.0 tag exposed that the old release job tried to push badge and health
  example edits directly to `main`. Release preparation now owns those edits,
  while the tag workflow verifies both lines without mutating the repository.
  Every publication job is gated on that tagged-tree preflight, and offline
  automation rejects both an ungated release and any future direct-main push.
  The Linux desktop release job also installs the same Wayland development
  headers and protocols as GUI CI, preventing a tag-only build failure after
  the GLFW dependency gained Wayland support.

### Decisions

- **2026-08-31 — Ship the Windows installer recovery as a v0.9.1 patch
  release.** The patch release contains the corrected WiX local-group
  authoring, stronger installer contracts and lifecycle tooling, and the
  release-pipeline corrections made after v0.9.0. Issue #83 closes only after
  the published artifact set and checksums are audited; the unavailable native
  Windows 11 lifecycle run remains disclosed rather than being inferred from
  source or package inspection.

- **2026-08-31 — Treat an omitted WiX group domain as a pinned Windows
  installer security contract.** S030 encoded `[ComputerName]` as local-group
  intent, but WiX 6 selects its unelevated domain action for every non-empty
  value. S035 deliberately removes that attribute while preserving the
  installing-user domain, group membership, service ordering, and uninstall
  state. Static source or MSI-table inspection remains necessary but cannot
  substitute for disposable Windows 11 lifecycle evidence.

## [0.9.0] - 2026-08-30

### Added

- **Task-completion chains now run end to end (Closes #72, #73, #74, #75,
  #76, and #77).** Success, failure, and any-outcome relationships supplement
  normal schedules and cascade through an acyclic graph. SQLite-backed delivery
  survives restart with an explicit at-least-once replay contract and correlated
  run history. Full lifecycle management is available through the local API,
  `gosched chain`, and a live desktop Chains view.

- **The complete go-schedule brand system is now repository-owned (Closes #10
  and #34).** The full BrandBuilder-standard kit now lives under `brand/` with
  its guide, portable masters, generated outputs, tokens, licensed fonts,
  platform assets, build source, verification report, and checksum inventory.
  A public Brand system page explains selection, accessibility, attribution,
  misuse, and downloads. An offline Go verifier checks all 108 inventoried
  artifacts and 59 declared consumer copies, and routine automation rejects
  checksum, encoding, SVG-portability, or downstream-copy drift.
- **`@reboot` is now a first-class scheduler-startup trigger (Closes #64).**
  `@reboot` and `at scheduler startup` round-trip through conversion, preview,
  task authoring, desktop editing, crontab import, and faithful export. Eligible
  tasks run exactly once per daemon process start with a distinct `startup`
  history origin, never on reload or mutation, and use normal overlap handling.
  Disabled task/group behavior, two independent starts, no catch-up, empty
  next-run previews, prior-database persistence, and retained crontab execution
  context are covered without requiring a manual VM lifecycle session.
- **Routine dependency discovery now opens reviewable proposals (Refs #40).**
  Dependabot checks Go modules weekly and GitHub Actions monthly, applies the
  existing `dependencies` label, caps each ecosystem at five open proposals,
  and groups compatible minor and patch version updates. Major and security
  updates remain individually visible, and existing pull-request verification
  applies without new branch or approval policy.
- **Cron parity now has a complete, faithful boundary (Closes #22).** Crontab
  import preserves ordered `CRON_TZ`, process environment, effective `SHELL -c`
  command text, cron percent stdin, and optional system-crontab run-as users.
  An explicit Unix/Quartz file dialect prevents command-field guessing. The
  task model and additive schema v8 persist stdin through API updates, restart,
  and execution. Supported six-field schedules compile seconds and Quartz `?`
  directly into the existing recurrence engine, preserve dialect-specific
  weekday numbering, and export six fields only when seconds are required.
  Five-field output remains stable. Export now refuses operational task context
  it cannot serialize instead of dropping it. The fidelity matrix records a
  decision for every issue #22 row and public copy promises the documented
  subset rather than universal cron parity.
- **Tasks now state their DST scheduling intent end to end (Closes #8).** Each
  task can follow local wall-clock readings, fixed elapsed intervals, or UTC
  readings, with per-task spring-gap `next_valid`/`skip` and fall-overlap
  `first`/`both`/`last` choices. Compatibility defaults preserve existing
  behavior. Schema v7, API, CLI, desktop Advanced Settings, live preview,
  detail, calendar, catch-up, restart, dispatch, and real IANA transition tests
  share one policy-aware evaluator. Elapsed mode refuses calendar shapes that
  do not have a fixed duration instead of approximating them. Elapsed schedules
  persist an absolute phase that presentation-timezone edits cannot move, dense
  fall-overlap rules jump directly into the repeated interval, and the desktop
  blocks incompatible elapsed submissions before Save.
- **Standard five-field cron combinations now run faithfully end to end (Refs
  #22).** Lists, ranges, wildcard and range steps, names, arbitrary weekday
  sets, and safe cross-field restrictions now share exact parsing, preview,
  execution, persistence, catch-up, editing, and canonical export behavior.
  This includes uneven field-local steps and policy-aware date sets across leap
  years and DST. Cron day-of-month/day-of-week OR semantics, Quartz layouts,
  boot triggers, and modifier composites remain named refusals, so #22 remains
  open for the deliberately excluded and operational crontab work.
- **Monthly calendar selectors now work end to end in both cron directions
  (Refs #22).** Day-of-month `L`, `1W` through `31W`, and `LW` can be
  explained, previewed, imported, retained, authored in plain language, and
  exported without approximation. A forward-only schema migration persists the
  typed nearest-weekday adjustment that RFC 5545 cannot express, with exact
  weekend, short-month, missing-date, leap-year, DST, restart, CLI, API, and
  desktop coverage. The desktop preview now also refreshes and sends the chosen
  missing-date policy. Composite modifiers and Quartz field layouts remain
  explicit refusals.
- **Last weekdays now convert faithfully in both cron directions (Refs #22).**
  One five-field day-of-week `weekdayL` term can be explained, previewed,
  imported, and retained as cron source, while native last-weekday schedules
  export canonically with Sunday as `0L`. Round-trip coverage spans DST and
  month boundaries; all missing-date policies are equivalent for this schedule,
  while broader `L` combinations remain named refusals.
- **Ordinal weekdays now convert faithfully in both cron directions (Refs
  #22).** One five-field `weekday#1..5` term can be explained, previewed,
  imported, and retained as cron source, while native first-through-fifth
  monthly weekday schedules export canonically with Sunday as `0`. Round-trip
  coverage spans DST, month boundaries, and absent fifth occurrences;
  incompatible missing-date policies and broader `#` forms remain named
  refusals.
- **Current product documentation now defines human-or-cron task authoring
  end to end (closes #52 and #50; refs #22).** README, CLI, GUI, API, and the
  authoritative S001 contract share equivalent weekday examples and name
  go-schedule's exact five-field subset, conversion workflows, source retention,
  field-local steps, expression-versus-crontab boundaries, task timezones, DST
  behavior, and fidelity refusals. S008 remains intact as historical evidence
  with explicit supersession notices.
- **The desktop Schedule field now accepts human phrases or supported five-field
  cron (Refs #50).** Local validation, live preview, create, and update use the
  shared syntax boundary. Cron tasks reopen with their retained expression,
  edits can switch syntax from the current text, and invalid or faithfully
  unsupported cron keeps Save disabled with a named reason. Human Start at,
  one-off, and expressionless legacy behavior remain unchanged. Broad
  cross-product documentation remains tracked by #52.
- **Task input now accepts human schedules or supported five-field cron (Refs
  #50).** Preview, create, update, and CLI task commands share one syntax
  boundary with optional explicit selection and no parser fallback. Task reads
  report the retained source syntax, and crontab import now keeps its normalized
  expression while continuing to show the plain-language explanation. Invalid
  or unsupported cron is refused before mutation. The desktop editor remains
  plain-language only in this slice.
- **The CLI now converts one schedule string locally in either direction
  (closes #51; groundwork for #50).** `gosched cron convert` produces a single
  cron or human string without the daemon or task mutation, supports explicit
  destination selection and stable JSON, and refuses lossy translations with
  exit code 2. Existing sub-daily cron translation now retains its `00:00`
  phase explicitly, and human-to-cron conversion refuses implicit or misaligned
  phases instead of manufacturing different run times.
- **The desktop sidebar now ends with a local Info view (closes #29 and #32).**
  Navigation follows Tasks, Groups, Schedule, Activity, and Info. The new view
  reuses the application mark, reports the exact running version, identifies
  ShruggieTech, and links to the company, source repository, and documentation
  without adding a daemon or network dependency.
- **The go-schedule brand system is now applied across the app, README, and docs
  site.** The desktop GUI ships a custom dark-first Fyne theme (`gui/theme.go`)
  built from the brand palette — Anchor Blue for selection/focus/links, Interval
  Mint for success, Hold Amber and Stop Red kept scarce — with the brand fonts
  (Geist, Geist Mono, Space Grotesk Bold) embedded and brand-radius geometry, and
  the window/app/exe icons replaced with the brand mark (a terminal prompt over
  five cron cells). The README gains the horizontal logo, brand-hued badges, the
  "A ShruggieTech project" endorsement, and drops the off-brand "Material Design"
  framing. The just-the-docs site gains a brand dark color scheme, the header
  logo, brand web fonts, favicons, a social-preview image, and the endorsement
  footer. Fonts ship under the SIL Open Font License (texts included beside the
  embedded/served fonts). No pinned artifact was modified.
- **The docs site now has a Changelog page.** A new `docs/changelog.md` appears
  in the site navigation and points to the canonical `CHANGELOG.md` in the
  repository root, and the docs home index links to it. The changelog stays a
  single source of truth at the repo root (edited as part of each release); the
  site page is a themed pointer rather than a duplicated copy, matching how the
  docs already link out to Contributing, the constitution, and the master spec.

### Changed

- **Privileged IPC now uses a dedicated administrative group (Closes #13).**
  The default `goschedadmin` group is resolved before the daemon listens.
  Windows uses its SID in a protected named-pipe descriptor; Linux and macOS
  apply and read back exact group ownership plus `0770` directory and `0660`
  socket modes. Missing or invalid non-empty groups fail closed. Explicitly
  empty `admin_group` retains broad local compatibility access with a warning.
  The Windows MSI creates or reuses the group, adds the installing account, and
  preserves group state on uninstall; Unix install guides provide the one-time
  setup commands.
- **Specification lifecycle history is now accurate and enforced (Closes
  #42).** All 29 Spec-Kit slices use one documented state vocabulary and cite
  objective delivery evidence. Stale historical task markers are reconciled
  with merge, release, or explicit waiver evidence. A fixture-backed offline
  checker runs inside the existing automation gate and rejects invalid states,
  missing evidence, contradictory task completion, and inventory drift.
- **The remaining useful secret-scanning controls are active (Closes #39).**
  Push protection, non-provider patterns, and validity checks now report
  enabled alongside the existing secret scanning and Dependabot security
  updates. Unrelated AI detection and delegated-dismissal controls remain
  unchanged.
- **The unverified Windows title/taskbar icon report is deferred from v0.9.**
  Issue #33 remains open with `needs: verification`, but now lives in `Post-v1`
  at P3 because its pre-brand-rollout premise has not been reproduced and the
  maintainer explicitly deferred the required manual install observation.
- **Pinned maintenance policy - lifecycle, dependency automation, and
  `scripts/automation-check.sh` (2026-08-30).** The existing offline automation
  gate now includes the Spec-Kit lifecycle contract and the exact two-ecosystem
  Dependabot policy. Spec templates and autopilot guidance separate local
  implementation completion from later PR publication bookkeeping.
- **Pinned documentation policy - `scripts/docs-check.sh` (2026-08-28).** The
  bounded current-surface check now requires the delivered standard composite
  cron breadth and retains the day-field OR refusal instead of requiring the
  superseded calendar-step refusal. Its generated fixtures cover the new claim.
- **Activity diagnostics now lead operators to the complete daemon log (closes
  #27 and #31).** Activity explicitly describes its records as a limited recent
  view and displays the daemon's exact configured rotating-log path, including
  a truthful unavailable state before metadata arrives. The daemon now emits
  one `daemon startup complete` record with structured `endpoint`, `db`, and
  `log_path` fields, while CLI log output remains compatible.
- **Pinned documentation policy - `scripts/docs-check.sh` (2026-08-28).** The
  offline docs gate now runs a bounded current-surface policy check plus generated
  acceptance and rejection fixtures. It rejects targeted obsolete human-only
  claims without scanning historical specifications or the changelog.
- **Calendar-field cron wildcard steps now retain field-local behavior.** The
  earlier safety refusal is superseded by direct compilation of the selected
  calendar values, so forms such as `0 9 */2 * *` no longer simplify or fail.
- **The documentation site's dark theme is complete and accessible (closes #35,
  #36, and #37).** Every syntax token now fails safely to readable base ink,
  deliberate brand colors cover named roles, highlighted lines and selected
  code remain dark-theme legible, and all published fences use a consistent
  `sh`, `bash`, `powershell`, or `text` category. The sidebar endorsement now
  follows the theme's responsive navigation gutters with balanced vertical
  spacing.
- **Pinned documentation policy - `scripts/docs-check.sh` (2026-08-27).** The
  existing network-free gate now rejects drift in the safe syntax-palette
  contract, selected/highlighted code treatment, responsive endorsement
  spacing, and published fence vocabulary. Front-matter and link behavior are
  unchanged.
- **The GUI's mixed log-and-alert view is now Activity (closes #26, #28, and
  #30).** Its unacknowledged-alert badge shows exact counts through 99 and caps
  higher counts at `99+`. The former Dismiss All/trash control is now Clear View
  with a clear-content icon and visible copy explaining that it hides current
  activity, acknowledges visible alerts, and does not delete records. Existing
  filtering, view cutoff, acknowledgement, and persistence behavior is
  unchanged.
- **Pinned artifact - `docs/INSTALL-windows.md` (2026-08-27).** Updated the
  stale GUI tab reference from Logs to Activity. Installation paths, commands,
  permissions, and behavior are unchanged.
- **All development now integrates through pull requests (closes #23).**
  Constitution v3.0.0 replaces the stale maintainer-direct-to-`main` rule with
  one review-branch workflow for maintainers, automation agents, and outside
  contributors. Its purpose is deliberately narrow: give third-party AI
  reviewers a durable place to comment before the maintainer decides whether to
  merge. Local parity and the autopilot halt precede publication. The README,
  contributor guidance, agent context, the autopilot protocol, and the pull
  request template now describe the same workflow and use `Closes` when a pull
  request fully completes an issue. No branch protection or approval machinery
  is added.
- **The repository now carries a reviewable GitHub security baseline (addresses
  #38 and #39).** `SECURITY.md` retains the private advisory route and now names
  `info@shruggie.tech` as its fallback plus repository administrators as the
  initial triage owners. A repository-owned advanced CodeQL workflow analyzes
  Go on pull requests, `main`, and weekly with a manual cgo-free build. The
  offline automation gate rejects obsolete or unknown security actions and
  drift in the workflow's triggers, permissions, and required analysis steps.
  Private reporting and supported secret protections are activated only after
  the autopilot publication authorization, with each hosted result reported
  separately.
- **Pinned artifact - `.github/workflows/codeql.yml` (2026-08-26).** Added
  advanced Go analysis using checkout/setup-go v7 and CodeQL v4, the
  module-selected Go version, a manual `CGO_ENABLED=0 go build ./...`,
  least-privilege `contents: read` and `security-events: write` permissions,
  and push-to-main, pull-request-to-main, and weekly triggers. The existing CI
  workflow and product build remain unchanged.
- **Pinned automation policy - `scripts/automation-check.sh` (2026-08-26).**
  Approved the researched CodeQL v4 init/analyze action families and added an
  offline, fixture-backed contract for the canonical security workflow. The
  policy remains network-free and preserves the existing eight-gate manifest.
- **Maintainer automation now has a supported hosted runtime and one local
  definition of green (closes #21 and #41).** CI and release actions use their
  researched Node 24 majors. `sh scripts/verify.sh all` runs the exact eight-gate
  pre-push contract, and CI, Make, contributor guidance, and autopilot delegate
  to it. An offline allowlist and fixture suite reject obsolete or unknown
  action references and verification-manifest drift.
- **Pinned artifact - `.github/workflows/ci.yml` (2026-08-26).** Updated
  checkout, Go setup, and artifact upload to their Node 24 majors, routed the
  seven existing verification commands through `scripts/verify.sh`, and added
  the offline automation-contract gate. Triggers, permissions, job boundaries,
  matrices, build commands, artifact names, and package selections are
  preserved.
- **Pinned artifact - `.github/workflows/release.yml` (2026-08-26).** Updated
  checkout, Go setup, and GitHub release publication to their Node 24 majors.
  Tag triggers, permissions, release inputs, output consumption, names, and file
  globs are unchanged; no tag or release was created during validation.
- **Pinned artifact - `Makefile` (2026-08-26).** Verification targets are thin
  delegates to `scripts/verify.sh`; `fmt` remains available and is explicitly
  identified as mutating.
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

### Fixed

- **The canonical lint gate again supports the repository's Go 1.25 target.**
  The pinned golangci-lint runner moved from v2.1.6, which was built with Go
  1.24 and refused to analyze the module, to v2.12.0, the latest upstream
  release whose declared build baseline remains Go 1.25. The verifier and
  current contributor instructions now use the same compatible pin.
- **The highlights-only release contract validates the entire document.**
  Tag-specific notes must contain only their heading, four to six one-line
  highlights, and the tagged full-changelog link. Extra paragraphs and nested
  headings are rejected instead of slipping into the public release body.
- **S029 maintenance checks now reject permissive near-misses.** Lifecycle
  validation requires recognizable commit, release, pull-request, or review-
  branch evidence instead of accepting the template placeholder. Dependabot
  validation rejects extra ecosystems plus major or security updates added to
  routine groups.
- **Windows installs now carry the go-schedule mark on the Start Menu shortcut
  and installed-apps entry (closes #24 and #25).** The MSI declares one
  canonical icon row sourced from the same multi-resolution icon embedded in
  the GUI executable, and both installer-created surfaces reference it
  explicitly. Portable mutation tests, the release-time WiX sanity check, and
  Windows Installer table inspection guard the relationships. Evidence tooling
  records candidate or published provenance separately, retains prompted native
  observations, and records cleanup plus final machine state after lifecycle
  failures. Clean published-MSI PATH verification (#16) and native title/taskbar
  observation (#33) remain open until a clean Windows desktop supplies their
  required evidence.
- **The docs site homepage no longer 404s at its root.** The Pages site was
  live and every content page served, but `https://shruggietech.github.io/go-schedule/`
  — the URL the README release badge and the issue-template links point at —
  returned 404. `docs/README.md` is the intended home page (front matter
  `title: Home`, `nav_order: 1`), but Jekyll builds a source file named
  `README.md` to `README.html`, not `index.html`, and nothing remapped it, so
  the site root had no index document. Adding `permalink: /` to that page's
  front matter emits it as `index.html` at the `baseurl` root, so the root now
  serves the docs home. `docs/README.md` is not a pinned artifact, so no dated
  decision was required.

### Decisions

- **2026-08-30 - Keep public release notes concise, reviewed, and tied to the
  exact tag.** Each release reads four to six highlights from
  `.github/release-notes/<tag>.md` and links to the tagged changelog for the
  complete record. GitHub-generated notes and long installation copy are
  disabled because they duplicate the changelog and obscure the release's main
  outcomes. The tag-derived path also prevents a later release from silently
  reusing stale copy.

- **2026-08-30 - Use durable source-run identities and an honest at-least-once
  replay boundary for completion chains.** A run record and its outgoing
  deliveries commit atomically, with uniqueness on chain plus immutable source
  run. Completed work never replays, while an interrupted claimed delivery is
  made pending at startup. This replaces the older configurable time-window
  dedup concept, which could either suppress legitimate later work or admit a
  crash-window duplicate and therefore could not provide the promised
  exactly-once external effect.

- **2026-08-30 - Consume committed platform artwork directly and verify brand
  copies without expanding the gate set.** The release workflow and WiX source
  now consume the canonical Windows ICO, macOS ICNS, Linux desktop entry, and
  hicolor icons rather than resampling artwork during release assembly.
  `scripts/automation-check.sh` runs the repository brand verifier inside the
  existing automation gate, preserving the eight-gate contract and adding no
  application runtime dependency. `.gitattributes` preserves canonical and
  declared consumer asset bytes so cross-platform line-ending normalization
  cannot invalidate approved hashes or create false drift. Issue #33 remains
  deferred because this slice does not claim manual Windows install/uninstall
  observation.
- **2026-08-30 - Reuse the retained event schedule shape without restoring the
  removed trigger subsystem.** Startup uses `kind=event` with
  `trigger_id=scheduler_startup` and the existing storage column.
  `Engine.Start` owns one startup snapshot outside `recompute`, so reload cannot
  fire it. Run history adds the migration-free text origin `startup`; legacy
  generic `event` values remain readable. No schema migration is warranted
  because every required persisted field already exists.
- **2026-08-30 - Fail closed for configured IPC groups and leave custom Unix
  directories operator-owned.** A non-empty `admin_group` never falls back to
  broad access. The daemon may secure its default data directory or create a
  missing custom parent, but it only verifies an existing custom parent. This
  avoids silently changing `/tmp` or another shared path. WiX and both required
  extensions move together from 5.0.2 to 6.0.2 because v6 is the first stable
  line with declarative local-group creation; the installer preserves the group
  and membership on uninstall to avoid deleting shared administrator state.
- **2026-08-30 - Close v0.9 maintenance gaps without expanding solo-project
  ceremony.** S029 combines lifecycle repair, low-noise dependency proposals,
  and supported secret-scanning activation into one review slice. It adds no
  branch protection, required approvals, auto-merge, dependency upgrade, or
  application change. The evidence-blocked Windows icon report moves to
  Post-v1 rather than interrupting the maintainer for a VM or manual install
  session.
- **2026-08-30 - Treat publication as evidence, not an implementation task.**
  Feature state now describes requirement, implementation, and verification
  maturity. Push, PR, review, merge, and branch cleanup are recorded through
  delivery evidence, preventing completed specs from retaining stale open
  bookkeeping tasks after publication.
- **2026-08-29 - Treat five/seven-hour transition gaps as wall-clock behavior,
  not an interval defect.** A six-hour local cycle remains six clock hours apart
  when the UTC offset changes, so its elapsed gap is legitimately five or seven
  hours. S027 deliberately corrects issue #8's defect classification without
  changing stored tasks: `wall_clock` remains the default, while explicit
  `elapsed` supplies exact-duration intent. UTC remains a distinct calendar
  basis, and elapsed rejects variable calendar periods rather than inventing a
  duration.
- **2026-08-28 - Keep one GUI Schedule field with shared syntax selection.** A
  syntax selector or cached response identity could disagree with edited text,
  so validation, preview, and save classify the current field through the S019
  boundary and carry its normalized expression and identity together. The GUI
  retains cron exactly, keeps Start at specific to human sub-daily intervals,
  and limits documentation changes to GUI adoption while #52 owns the broader
  rewrite.
- **2026-08-28 - Add dual-syntax task input without a storage migration.** The
  existing `schedules.expression` column retains either accepted source, and
  response-only syntax identity is derived from it. Cron is converted into the
  same RRULE and anchor model used by human input, so raw cron never becomes an
  engine input. The GUI sends an explicit human hint and defers cron editing to
  later work on #50.
- **2026-08-27 - Use one canonical WiX icon row and layer Windows installer
  evidence.** `build/windows/goschedule.wxs` and
  `build/windows/verify_wxs.ps1` are pinned artifacts. The package-level Icon,
  `ARPPRODUCTICON`, and Start Menu shortcut now share
  `cmd/gosched-gui/icon.ico`; the existing 64-bit executable-resource step is
  preserved because it is already correctly declared. Source tests, compiled
  MSI table inspection, published-artifact lifecycle checks, and native desktop
  observations report independently as proven, failed, or unavailable. The
  current installed development host is deliberately not modified to imitate a
  clean machine, and unavailable evidence does not close #16 or #33.
- **2026-08-26 - Constitution v3.0.0 reinstates PR-first integration.** Recent
  work has repeatedly used third-party PR review, so the direct-push policy no
  longer described the useful operating model. This is a major governance
  amendment because it reintroduces a mandatory integration gate for every
  author. The once-per-feature autopilot halt now precedes review-branch
  publication and PR creation. That publication authorization covers verified,
  in-scope review-fix pushes to the same open PR, so responding to CI or review
  does not create contradictory extra halts.
- **2026-08-26 - Keep PR governance proportional to the project.** The PR exists
  to obtain AI review and preserve discussion for a one-developer application
  with no users. Branch protection, approval counts, mandatory conversation
  resolution, fixed required-check contexts, and a governance policy checker
  were considered and rejected as maintenance-heavy ceremony. Existing local
  verification and hosted CI remain evidence, AI comments remain advisory, and
  the maintainer retains final merge judgment.
- **2026-08-26 - Use repository-owned advanced CodeQL with a manual headless Go
  build.** Reviewable triggers and least-privilege permissions are part of this
  repository's security contract. `CGO_ENABLED=0 go build ./...` traces compiled
  Go without pulling the Fyne/OpenGL entry point into the hosted analysis build.
- **2026-08-26 - Validate the security workflow with the existing offline
  automation policy.** CodeQL v4 joins the researched Node 24 action allowlist,
  and temporary fixtures reject five drift classes before publication:
  obsolete CodeQL, unknown actions, missing triggers, insufficient permissions,
  and missing analysis.
- **2026-08-26 - Report hosted security capabilities individually.** Private
  reporting and each secret-scanning control must read back as enabled or carry
  an exact unavailable, unverified, or failed result. Plan limitations and token
  scope are not treated as green, and validation creates no fake advisory or
  test secret.
- **2026-08-26 - Use the latest researched Node 24 action majors behind an
  offline allowlist.** CI uses `actions/checkout@v7`, `actions/setup-go@v7`, and
  `actions/upload-artifact@v7`; release automation also uses checkout/setup-go
  v7 and `softprops/action-gh-release@v3`. Major selectors retain upstream
  security fixes while the repository policy makes runtime regression visible.
- **2026-08-26 - Make `scripts/verify.sh` the command source for CI, Make, and
  maintainer guidance.** A POSIX shell driver owns gate order, package
  selection, the pinned linter version, coverage thresholds, and failure
  propagation. The independent automation check validates its exact manifest,
  preventing consumers from silently defining different meanings of green.
- **2026-08-26 - Preserve workflow behavior while changing pinned automation.**
  Runtime upgrades and verification delegation do not alter triggers,
  permissions, matrices, artifacts, release globs, or publication behavior.
  Release validation is static so this maintenance feature has no external
  release-side effect.
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
  was then a refusal. S026 later added faithful field-set compilation while
  retaining this in-tree parser and adding no dependency.

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

- **2026-08-28** — **Cron and human input compile independently into the same
  durable schedule contract.** This deliberately supersedes the 2026-07-23
  phrase-only conversion rule. Re-parsing a generated English phrase cannot
  preserve standard cron lists, ranges, uneven field-local steps, or safe
  cross-field combinations. The original normalized cron remains the editable
  source, while its complete readable description is display metadata and the
  compiled recurrence is authoritative for preview, execution, persistence,
  catch-up, and export. The shared durable model keeps the two authoring paths
  aligned without pretending the natural-language grammar expresses every cron
  set.


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

[Unreleased]: https://github.com/shruggietech/go-schedule/compare/v0.9.1...HEAD
[0.9.1]: https://github.com/shruggietech/go-schedule/compare/v0.9.0...v0.9.1
[0.9.0]: https://github.com/shruggietech/go-schedule/compare/v0.8.0...v0.9.0
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
