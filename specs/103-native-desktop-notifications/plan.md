# Implementation Plan: Optional native desktop notifications and release preparation

**Branch**: `codex/103-native-desktop-notifications` | **Date**: 2026-09-23 | **Spec**: [spec.md](spec.md)
**Input**: S103 and #177, with release preparation for #185. SMTP #176 remains open.

## Summary

Add opt-in, filterable native popups to the running Wails desktop using its existing cross-platform notification runtime. Observe registered daemons through independent identity-pinned, read-only clients rather than switching the user's selected target. Persist popup settings in desktop-local preferences, suppress startup/backlog replay and duplicates, and route activations through the existing identity-checked activity drilldown. Prepare honest v1.5.0 release copy but do not tag or publish before the post-merge release authorization.

## Technical Context

**Language/Version**: Go 1.26 and React/TypeScript in Wails v2.15.0.
**Primary Dependencies**: Existing Wails notification runtime, existing local IPC/remote JSON clients, existing protected profile credentials. No new OS notification library or child process.
**Storage**: Existing atomic desktop preference JSON; no daemon migration or new listener. In-memory bounded per-daemon dedupe state only.
**Testing**: Deterministic Go service tests, frontend interaction tests, Windows/Linux build checks, macOS cross-platform CI and packaged behavior checks where available, canonical `sh scripts/verify.sh all`.
**Target Platform**: Windows desktop, packaged macOS desktop, and Linux graphical sessions with a notification service. Headless Linux gracefully unsupported.
**Project Type**: Wails desktop client over independent schedulers.
**Performance Goals**: Notifications arrive within seconds after observed events without unbounded goroutines, polling, or retained event memory.
**Constraints**: No new resident service, no new network listener, no implicit daemon mutation, no popup startup storm, no cold-start click-through promise without a platform contract.
**Scale/Scope**: One local daemon plus registered remote profiles, bounded concurrency and event IDs; v1.5.0 release preparation only.

## Constitution Check

- **Code quality**: A bounded monitor owns every stream context and terminates with application shutdown; platform presentation is injected for deterministic tests.
- **Testing**: Cover mute, filters, identity-safe dedupe/activation, reconnect, denial, and cancellation. Run race detector and canonical gates.
- **UX consistency**: Desktop-local settings clearly separate popups from daemon-owned webhook policies. State limitations and closed-app behavior in UI and docs.
- **Performance**: Bound registered streams, event fetches, and dedupe memory; do not query unbounded histories or spawn OS notification helper processes.
- **Autopilot**: Spec-kit artifacts and analyze gate precede implementation. Review branch and PR are explicitly authorized; release tag/publication is not authorized by this PR request.

**Pre-design gate**: Pass. No constitution exception required.

## Design Decisions

1. Use Wails v2.15 native notifications. Its Windows, macOS, and Linux backends already provide delivery and response callbacks. An additional library or per-platform subprocess would duplicate the existing dependency and violate the no-visible-console rule.
2. Start observation after Wails DOM readiness and stop on shutdown. Request macOS permission only when the user enables popups; denied or absent platform facilities become a quiet status.
3. Observe every opted-in registered daemon with a bounded independent client, never by changing the selected connection. Pin saved daemon identity and credential as existing All Systems does. Consume secret-free event IDs, fetch an exact terminal run or alert, and dedupe by daemon identity plus record ID. The server's live-only subscription is the startup boundary; it does not replay historical records or need a preliminary snapshot.
4. Keep popup preferences desktop-local. They do not create daemon notification channels or modify task/group assignments; #19 remains open for the future unified-channel model.
5. On activation, emit an immutable source-and-record intent to React, where existing `selectAndWaitForTarget` verifies current identity before opening Activity. If an OS does not report activation, users can still open Activity manually and the UI does not promise deep linking there.
6. Treat v1.5.0 as a cumulative feature release. Release documentation will describe only delivered functionality. SMTP is explicitly deferred without closing #176 or #19, and release issue #185 remains open until actual publication.

## Project Structure

```text
desktop/notifications/                popup monitor, source adapters, candidates, tests
desktop/settings/                     persistent desktop-only popup preferences
desktop/app.go, desktop/main.go        lifecycle, native runtime, bridge
desktop/frontend/src/notifications/  popup preferences and explanatory UI
desktop/frontend/src/App.tsx          identity-safe activity activation
internal/api/client/events.go         shared secret-free observation stream
docs/notifications.md                 user behavior and platform limitations
.github/release-notes/                cumulative v1.5.0 preparation
specs/103-native-desktop-notifications/ spec-kit artifacts
```

**Structure Decision**: Reuse the existing desktop notification, settings, connection, and operations seams. The monitor does not change daemon scheduling or delivery persistence.

## Post-design Constitution Check

Pass if the implementation preserves bounded lifecycle, opt-in persistence, exact identity checks, and honest release claims. The spec-kit analyze gate must map every FR to a concrete task before coding.

## Complexity Tracking

No constitution violation. A per-registration read-only observer is necessary because the existing connection manager intentionally watches only one selected daemon and switching it for background popups would violate target-safety expectations.
