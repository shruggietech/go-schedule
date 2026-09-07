# Research: Desktop Settings, Information, and Recovery

## Decision 1: Migrate durable intent, not framework mechanics

**Decision**: Migrate a valid legacy `appearance.mode` once. Retire `appearance.font` and `appearance.scroll_sensitivity`.

**Rationale**: Light, dark, and system are user intent shared by both desktops. The old font resource selection and toolkit scroll multiplier are Fyne implementation controls that do not map truthfully to browser typography or native web scrolling.

**Alternatives considered**: Migrating every key was rejected because two controls would become inert or misleading. Resetting everything was rejected because it discards a valid durable appearance choice.

## Decision 2: Use a small independent versioned JSON store

**Decision**: Store desktop preferences below the user's configuration directory in a go-schedule-owned path, with an explicit schema version, appearance, and migration disposition. Replace the file through a same-directory temporary file and rename, with restrictive permissions where supported.

**Rationale**: The new desktop must not retain a runtime dependency on Fyne. A small standard-library document is inspectable, portable, easy to migrate, and requires no new dependency.

**Alternatives considered**: Continuing to use Fyne preferences was rejected because it keeps the retired framework in the new desktop lifecycle. Registry or platform-specific settings were rejected because they multiply implementations without delivery value. A daemon-owned setting was rejected because appearance must remain usable offline and is not scheduler policy.

## Decision 3: Default to system and never block startup on legacy damage

**Decision**: New installations, missing values, unknown values, malformed JSON, and unreadable legacy data use system appearance. The transition result remains visible but non-blocking.

**Rationale**: System follows platform preference and matches the current Wails direction. A damaged file should not strand the control center before connection recovery or storage help can be reached.

**Alternatives considered**: Preserving the Fyne dark default was rejected because it is toolkit history rather than current product intent. Treating malformed legacy data as fatal was rejected because preferences are optional presentation state.

## Decision 4: Consolidate Options and Info in Settings

**Decision**: One Settings route contains Appearance, Preference transition, Storage, and About sections. Connections remains a dedicated operational route.

**Rationale**: The old separate windows contained small related sets of local application facts. Consolidation reduces navigation while keeping daemon recovery distinct and visible in the existing shell order.

**Alternatives considered**: Adding a separate Info route was rejected because it would add navigation for a small static section. Embedding connection recovery in Settings was rejected because connection state affects the entire application.

## Decision 5: Preserve storage truth through source-specific resolution

**Decision**: Resolve desktop-owned paths locally, adapt daemon-owned effective paths only from `RuntimeInfo`, and preserve external ownership and removal semantics from the legacy inventory model.

**Rationale**: Installed paths and daemon configuration can differ from defaults. Source-specific resolution prevents support guidance from promising that absent, external, or relocated data will be removed.

**Alternatives considered**: Recomputing daemon defaults in the desktop was rejected because configured paths can differ. Omitting unavailable rows was rejected because it hides which information requires a connection.

## Decision 6: Authorize native actions by identifiers

**Decision**: React sends a storage record identifier for copy and a product-link key for open. Go resolves the current path or fixed HTTPS destination before calling Wails native integration.

**Rationale**: Native clipboard and browser methods are privileged application boundaries. Stable keys prevent compromised or stale frontend state from turning these methods into arbitrary text-copy or URL-launch primitives.

**Alternatives considered**: Passing raw paths and URLs from React was rejected as unnecessarily broad. Browser-only clipboard access was rejected because desktop permission and secure-context behavior is less predictable.

## Decision 7: Reuse connection state and make recovery inline

**Decision**: Connections consumes the existing connection manager snapshot and retry operation. The header dialog remains user-invoked detail, while failures never auto-open or repeatedly interrupt the user.

**Rationale**: A single connection authority avoids competing retries and state disagreement. A stable route supports deliberate troubleshooting and preserves focus during status changes.

**Alternatives considered**: A second settings-specific health client was rejected because it duplicates retry policy. Automatic modal display was rejected because recurring daemon failures would repeatedly steal focus.

## Decision 8: Verify GitHub branch cleanup with the S065 merge

**Decision**: The native `delete_branch_on_merge` repository setting is already enabled. Keep #195 open and use the S065 same-repository squash merge as the required behavioral proof.

**Rationale**: The API setting alone proves configuration, not the full merge lifecycle. S065 provides a normal repository topic branch without adding disposable code or custom automation.

**Alternatives considered**: Creating a custom workflow was rejected by #195 and would introduce broader credentials and race conditions. Closing #195 before a real merge was rejected because one acceptance criterion would remain unverified.
