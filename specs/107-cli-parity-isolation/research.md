# S107 research and decisions

**Date**: 2026-10-19. Baseline: `bae303071d1dd1af050cffb201267052ae58b2c7`.

## Invocation and isolation

Decision: Command-context invocation dependencies and a switchable client replace mutable package globals, direct process I/O and MCP callbacks. Rationale: independent roots must not affect each other; receiver conversion would churn every command. Captured-client constructors need the scoped client explicitly. Source: internal/cli/cli.go, target.go, mcp.go, access.go and daemon.go.

Decision: Add --ipc-path only. Temporary DataDir does not change the fixed Windows pipe; service install already owns --config. Explicit local selection never reads native remote credentials or falls back. Fixture ownership uses Manifest.InstallationID plus RuntimeInfo paths and new authoritative IPC/PID fields.

Decision: Reuse qualification daemon/config patterns but add a small owned supervisor. Existing exact-PID Kill does not bound descendants. Windows uses a kill-on-close Job Object assigned before execution; Unix uses owned process groups. Hidden redirected noninteractive children, canonical path containment including symlinks/junctions, retained file evidence and identity checks precede fixture mutation. Ownership gates apply to QA, not ordinary product eligibility.

## Authoring and notifications

Decision: Reuse current task/group/preview request types and add ClearEnv, consistent with existing explicit clear fields. Current Env omitempty drops a non-nil empty map. Fix desktop and CLI together, preserving base environment. Bundles omit execution inputs; create-then-disable has an activation window.

Decision: Strict bounded NotificationAssignmentsRequest file/stdin input plus canonical export enables heterogeneous round trips. Existing APIs validate complete replacement and inheritance. Preserve shorthand.

## Views

Decision: Reuse GetCalendar/ListActiveRuns/GetRun/RuntimeInfo/ListLogs/SystemSummary/Search. Preserve legacy default formats with opt-in modes. Validate malformed RFC3339 before HTTP because legacy calendar parsing defaults invalid values; maximum 30 days matches desktop.

Decision: Keep runtime local-only, as current remote route tests explicitly exclude machine paths. Share only transport-neutral bounded observation plumbing; preserve desktop presentation/cache behavior and all target authorization/identity semantics.

## Preferences

Decision: Extract schema-v1 models, paths, Fyne migration, popup validation and atomic updates into internal/userprefs, retaining desktop aliases. Add cross-process locking around read-modify-write; existing per-Service mutex is insufficient across binaries. Preserve restrictive permissions, unsupported-version rejection and explicit malformed-file restore.

Decision: Reuse profile Store.SetActive/Find and ActiveDesktopProfileID. Local administration does not connect to a daemon, read credentials or change invocation selection. No schema merge/change.

## Alternatives and scope

Retaining globals fails isolation; importing desktop couples the CLI to Wails; duplicate scheduler/preference implementations drift; a plaintext credential environment switch changes production behavior. All rejected.

Installed /speckit-plan dispatched focused research agents for isolation, views and preferences; no source edits occurred during research. No dependency, owner decision or irreversible ambiguity remains. #289/#290 stay the complete evaluator/report follow-on.

## Publication review, 2026-10-19

Codex round 1 on PR #305 found unbounded credential setup and incomplete Unix descendant drainage (#306/#307). A context-only wrapper around the existing synchronous native read cannot interrupt it; a goroutine wrapper would strand work and violate the termination contract. Native reads therefore run in a hidden owned CLI subprocess with private request/response pipes. Injected adapters expose LoadContext explicitly. The normal keyring backend and persistent credential format remain unchanged.

The cancellation design follows [Go exec CommandContext/Cancel/WaitDelay](https://pkg.go.dev/os/exec). Unix drainage follows the [POSIX kill process-group and null-signal contract](https://pubs.opengroup.org/onlinepubs/009604499/functions/kill.html): signal only the owned group and wait for absence, with bounded errors. The same group lifecycle helper serves fixtures and native reads, including macOS security-tool descendants.

Hosted macOS race failures showed that the default per-user temporary prefix produced an unusable socket (#308). A short private UUID root beneath canonical /private/tmp preserves ownership while fitting the supported socket address. Explicit overlong parents fail before runtime creation.

## Second review correction, 2026-10-19

The ordinary private-flag entry still allowed a caller to supply profile IDs on stdin and obtain a native token on stdout (#309). An internal name and a same-user ACL alone do not establish parent authority. The corrected design uses mutually verified OS peer process identities, the actual parent relationship and executable-image equality before credential access; the parent accepts only its exact child. All token traffic uses a dedicated operation-owned local channel with bounded framing, and standard streams are discarded.

Windows peer identity follows [GetNamedPipeServerProcessId](https://learn.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-getnamedpipeserverprocessid) and its client counterpart. macOS executable authentication follows [Apple libproc proc_pidpath](https://github.com/apple-oss-distributions/xnu/blob/main/libsyscall/wrappers/libproc/libproc.c) and its kernel executable path, rather than mutable process arguments. No new dependency, native credential format, end-user eligibility rule or authentication prompt is introduced.
