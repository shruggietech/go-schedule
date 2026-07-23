# TODO

**Audience:** maintainers\
**Authoritative task list:** [`specs/001-task-scheduler/tasks.md`](specs/001-task-scheduler/tasks.md)\
**Release history:** [`CHANGELOG.md`](CHANGELOG.md)

A high-level roadmap. Everything the master specification scoped for v1 is
delivered; what remains below is genuinely open.

## Delivered

The v1 scope is complete and released. Each line corresponds to a phase of
`specs/001-task-scheduler/tasks.md`, and the features that followed it have
their own specs under `specs/`.

- [x] **Setup and foundations** — Go module, lint and CI, injected clock,
      config, structured logging, SQLite store with forward-only migrations,
      local IPC, API, daemon.
- [x] **Scheduling (MVP)** — human-readable recurrence and one-off runs, per-task
      IANA timezones with DST resolution, windowless executor, task CLI, system
      service and start-on-boot, run history, `run_as`, cron-parity suite.
- [x] **Desktop GUI** — calendar and schedule views, guided task editor with live
      schedule preview, real-time updates driven by broker events, unified Logs
      view with filters and on-disk JSONL.
- [x] **Nested task groups** — groups within groups, cascading enable and
      disable, assignment from both the CLI and the GUI.
- [x] **Downtime catch-up** — one catch-up run per task after missed runs, then
      normal resumption.
- [x] **Packaging** — Windows `.msi` with an auto-starting service and a `PATH`
      entry, macOS desktop bundle, Linux and macOS archives, checksummed
      releases.
- [x] **Maintainer test scripts** — evidence that an installed daemon fires on
      time, survives restarts, catches up, and honors overlap policies
      (`docs/test-scripts.md`).
- [x] **Documentation set** — per-platform install guides, CLI reference,
      contribution, security, and conduct documents
      (`specs/007-issue-cleanup-docs/`).

Event triggers, once planned as Phase 6, were **removed** from the product in
`specs/004-rebrand-gui-overhaul/` along with their store migration. They are not
pending work.

## Open

- [ ] **Harden local IPC access control.** Today the Windows named pipe grants
      Authenticated Users read and write, and the Unix socket is reachable by
      any local user who can traverse the data directory. Both are deliberate,
      both suit a single-user machine, and both are too coarse for a multi-user
      one. Narrowing to a dedicated administrative group is the intended fix;
      `config.AdminGroup` exists for it. Documented in
      [`SECURITY.md`](SECURITY.md).
- [ ] **Goroutine-leak test and a dispatch-latency benchmark.** The p99 budget
      of 100 ms is documented and believed, but not measured by a committed
      benchmark.
- [ ] **Sign the release artifacts.** The `.msi` is not Authenticode-signed and
      the macOS builds are neither signed nor notarized, so every install path
      shows a warning that users are told to click past. A checksum file is a
      weaker guarantee than a signature and asks more of the reader.
- [ ] **Verify the `PATH` fix end to end.** The fix for issue #5 cannot be
      confirmed from a development machine, which already has the directory on
      `PATH` — precisely what hid the defect. It needs a clean machine and a
      released `.msi`.

## Later — out of scope for v1

- [ ] External trigger sources (CLI or API-delivered events) and file or folder
      watching.
- [ ] Remote or multi-user GUI access.
- [ ] External notification channels — email, push, webhooks.
- [ ] Distributed, multi-machine scheduling.
