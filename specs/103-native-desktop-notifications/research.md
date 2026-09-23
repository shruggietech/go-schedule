# Research: Native desktop notifications

## Platform delivery

Wails v2.15 already exposes native notification initialization, availability and authorization checks, presentation, response callbacks, and cleanup for Windows, macOS, and Linux. Use this runtime rather than another OS library or a child process. Initialize after DOM readiness, request macOS authorization only when a user enables popups, and treat unavailable or denied OS facilities as an explicit local status. Windows toast activation depends on packaged application identity; Linux click activation depends on the desktop notification server, so neither should be promised where unsupported. Sources: [Wails notification runtime](https://wails.io/docs/reference/runtime/notification/), [Microsoft desktop toast identity](https://learn.microsoft.com/en-us/windows/win32/shell/enable-desktop-toast-with-appusermodelid), [Apple notification authorization](https://developer.apple.com/documentation/usernotifications/asking-permission-to-use-notifications), [Desktop Notifications Specification](https://specifications.freedesktop.org/notification/latest-single/).

## Event source and identity

The selected-connection manager observes only one target and is deliberately allowed to switch. It cannot safely drive all-daemon popups. Registered-target observation must use independent read-only clients with pinned daemon identity, following the All Systems target model. The event stream provides resource IDs, while an exact read supplies terminal run outcome or alert severity. The server's live-only subscription does not replay stored history or need a startup snapshot. Deduplicate by daemon identity and source record ID, and bound the cache and stream count. The existing observation projection in the API removes sensitive event payloads while preserving run and alert IDs.

## Navigation and local state

The existing Activity drilldown can open an exact run or alert after the frontend's identity-checked connection selection. Persist only desktop-local popup preferences in the existing atomic preferences file. Daemon webhook policies remain separate; this slice does not implement the unified notification-channel model represented by #19.

## Release boundary

Prepare cumulative v1.5.0 notes only for delivered features. SMTP #176 and parent #19 remain open. A PR merge is not a public release: tag and publication require a later release authorization and must not be simulated by issue closure or release copy.
