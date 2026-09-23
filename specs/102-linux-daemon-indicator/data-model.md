# S102 data model

## Local service observation

Fields: `state` (running, stopped, starting, stopping, not_installed, unreachable, unknown), `scmState` (raw service-manager classification), `detail` (plain-language explanation), and `observedAt` (UTC timestamp). Running is valid only when the installed unit is active and a fresh local IPC health request succeeds.

## Service action

Fields: requested verb (start, stop, restart), outcome (accepted, rejected, unavailable, cancelled, failed, indeterminate), user-facing message, and latest observation. A mutation begins only after any required confirmation and ends with an observed state, never an optimistic status.

## Session indicator

One item per Linux user session, owned by a lock and D-Bus registration. Its menu mirrors the observation and exposes Open, available service actions, and Quit indicator. Quit closes only the session indicator, not the daemon.

## GUI activation

One local, same-user activation endpoint per session. A second GUI launch signals the existing instance and exits. The endpoint lifecycle follows the GUI instance, not the daemon.
