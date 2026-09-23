# Contract: Desktop popups

## Preferences

The desktop bridge exposes current popup preferences and OS capability/status, and accepts validated preference updates. Disabled is the default. Saving a disabled preference immediately stops presentation and observation; restoring desktop defaults disables it. Permission denial must not be interpreted as delivery success.

## Presentation

With the app running, a newly observed terminal run or alert matching the user's condition, severity, and daemon selection can produce at most one native popup per daemon and exact record. Startup history, duplicate event delivery, non-terminal run updates, and events outside the selected filters produce none. The popup contains a short safe summary and source name, not command output or secrets.

## Activation

Where the OS reports a notification response, a click opens the exact Activity record only after the current registered target resolves to the expected daemon ID. When that identity check fails, the app does not display a different daemon's record. Platforms without click callbacks still deliver the popup where supported and state that Activity can be opened manually.

## Lifecycle

Observation and native notification resources start after desktop readiness and stop at app shutdown. Reconnect does not replay historical records. A closed desktop app sends no native popup. This contract does not create a daemon notification channel or a persistent delivery queue.
