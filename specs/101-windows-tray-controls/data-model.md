# S101 data model

## LocalServiceSnapshot

- `scmState`: NotInstalled, Stopped, StartPending, Running, StopPending, or Unknown
- `health`: Healthy, Unreachable, NotChecked, with observed timestamp
- `displayState`: Not installed, Stopped, Starting, Running, Stopping, or Unknown/unreachable
- `detail`: Actionable bounded error text, never a credential

Running is valid only for SCM Running and fresh Healthy local IPC. A pending state takes precedence over health. NotInstalled is distinct from failure to query SCM.

## ServiceAction

- `kind`: Start, Stop, Restart
- `phase`: Requested, Confirmed, Elevating, Executing, Observing, Succeeded, Cancelled, Failed
- `result`: Observed snapshot and optional actionable error

Stop and Restart cannot enter Elevating without confirmation. A command exit alone cannot enter Succeeded; the target state must be observed.

## SessionCompanion

- Session-scoped instance ownership
- Icon registration generation, recreated on Explorer restart
- GUI activation endpoint for the same signed-in session

One owner per session; shutdown unregisters icon. Icon and GUI lifetimes are independent.
