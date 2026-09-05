# Research: Filesystem Watchers

## Decision 1: Use fsnotify v1.10.0 as a direct dependency behind an internal observer interface

**Decision**: Promote `github.com/fsnotify/fsnotify` from an indirect GUI dependency to a direct v1.10.0 dependency and adapt it behind a narrow internal interface.

**Rationale**: The upstream package is the established cross-platform native-notification abstraction for Go. Its documentation requires consuming both Events and Errors, states that paths must exist before registration, and documents non-recursive registration. Version 1.10.0 includes Windows watch-list race corrections and kqueue descriptor-leak fixes that materially support this long-running feature. The internal boundary permits deterministic fake-observer tests and keeps native quirks out of domain logic. Sources: [fsnotify package documentation](https://pkg.go.dev/github.com/fsnotify/fsnotify) and [fsnotify releases](https://github.com/fsnotify/fsnotify/releases/tag/v1.10.0).

**Alternatives considered**:

- Hand-written inotify, ReadDirectoryChangesW, and kqueue backends would multiply platform risk and testing cost without product value.
- Polling every watched path would work on more network filesystems but violates the resource-efficiency goal and loses native rename signals.
- Remaining on v1.9.0 would retain known lifecycle and race defects in precisely the code path S056 depends upon.

## Decision 2: Watch parent directories for file-kind watchers

**Decision**: Register a file watcher's parent directory and filter its exact cleaned path instead of asking the observer to watch the file object.

**Rationale**: Upstream explicitly warns that watching individual files loses observation during atomic-save replacement. Parent observation sees a final-name create after rename into place and also works while the target file itself is missing. This aligns file identity with the configured path rather than an inode or platform handle. Source: [fsnotify FAQ, Watching files](https://github.com/fsnotify/fsnotify#watching-files).

**Alternatives considered**:

- Watching the file directly fails common atomic-save workflows and cannot register a missing file.
- Polling the file after deletion narrows but does not eliminate the notification gap and adds a second observation model.

## Decision 3: Implement recursion explicitly and exclude link traversal

**Decision**: Walk existing real subdirectories during registration and add newly created real subdirectories when recursion is enabled. Do not follow symbolic-link directories or Windows reparse-point directories.

**Rationale**: Upstream registration is non-recursive and requires each directory to be added. Explicit traversal gives consistent filtering and lets the product enforce a hard scope boundary against loops and path escape. Source: [fsnotify FAQ, Are subdirectories watched?](https://github.com/fsnotify/fsnotify#are-subdirectories-watched).

**Alternatives considered**:

- Relying on undocumented native recursion would differ by backend and is not available through the stable dependency contract.
- Following links makes scope surprising and can create cycles or observe outside the configured root.

## Decision 4: Treat notifications as hints and settle through an injected-clock state machine

**Decision**: Accept create and write hints, then debounce and compare two regular-file snapshots separated by the stability duration before dispatch. Rename into place is represented by the final path's create signal where supported.

**Rationale**: Upstream documents that a single logical write can emit many events and that low-level sequences differ by platform. The product requirement is a stable resulting file, not a particular kernel opcode. Separating quiet-time coalescing from stability verification avoids early task invocation and is deterministic under a fake clock. Sources: [fsnotify Event documentation](https://pkg.go.dev/github.com/fsnotify/fsnotify#Event) and [fsnotify source event contract](https://github.com/fsnotify/fsnotify/blob/main/fsnotify.go).

**Alternatives considered**:

- Dispatching directly on Write is vulnerable to partial files and write storms.
- Platform-specific close-write operations are not portable and would produce different feature semantics.
- Content hashing reads potentially large files and exceeds the scheduler's responsibility.

## Decision 5: Rebuild the observer on configuration changes and observer errors

**Decision**: Use generation-scoped full rebuilds rather than incremental watch mutation. A single recovery deadline rebuilds all enabled definitions after errors or degraded registrations.

**Rationale**: At the supported scale of 100 definitions, rebuild cost is small. It eliminates stale registrations, provides a clean generation boundary for pending events, and handles overflow whose exact affected watch may be unknowable. It also avoids concurrent Add, Remove, WatchList, and Close complexity in upstream native backends.

**Alternatives considered**:

- Incremental diffing reduces registration calls but introduces subtle stale-watch and event-generation races.
- One observer and goroutine per configured watcher simplifies attribution but creates unnecessary descriptors and goroutines and makes shared recursive roots wasteful.

## Decision 6: Keep health ephemeral and publish transitions only

**Decision**: Persist intended configuration only. The running daemon owns current health and suppresses identical state and reason repetitions.

**Rationale**: `active` is true only for a specific daemon process and observer generation. Persisting it would make stopped or crashed state misleading. Transition-only reporting gives operators actionable recovery visibility without flooding logs during a missing-root retry loop.

**Alternatives considered**:

- Persisted health becomes stale outside daemon lifecycle and complicates crash semantics.
- Logging every retry makes a predictable missing path look like an application failure storm.

## Decision 7: Use at-most-once observation with no downtime replay

**Decision**: Start observation prospectively after registration and do not scan trees to infer missed events.

**Rationale**: Native watchers are not durable event logs. A scan cannot determine which files changed, whether a prior version was already dispatched, or the order of changes. Claiming recovery would be false. Operators needing durable delivery should invoke opaque external triggers from their producer or use a domain-specific queue.

**Alternatives considered**:

- Startup scans create duplicate or fabricated runs without a durable per-file checkpoint model.
- Persisting every candidate expands the feature into a durable ingestion queue, which is outside #135.

## Decision 8: Keep matched paths out of durable run provenance

**Decision**: Add `source_watcher_id` and `filesystem_watcher` to run history, but do not record the matched path.

**Rationale**: Watcher identity is sufficient to explain the configured cause and remains valid if a root moves. Full local paths often reveal usernames or sensitive directory structure and are not necessary for scheduler history. Authorized watcher configuration output still shows the configured root.

**Alternatives considered**:

- Reusing `source_trigger_id` would blur opaque external trigger and filesystem source contracts.
- Recording each matched path improves per-file diagnostics but creates unbounded sensitive history and is not requested.
