# Linux desktop service contract

The existing Wails `LocalServiceSnapshot()` response retains `state`, `scmState`, `detail`, and `observedAt`. `ControlLocalService(action, confirmed)` continues to accept only `start`, `stop`, and `restart`, and returns `action`, `outcome`, `message`, and `snapshot`. On Linux these refer solely to the installed `goschedd.service` on This computer, regardless of selected remote daemon.

The indicator presents the same state names and commands. Stop and Restart require impact confirmation. Open activates one GUI. Quit removes the indicator but leaves the installed daemon and GUI unchanged.

Supported indicator environment: a Linux graphical session with an accessible D-Bus session bus and StatusNotifier watcher/host. Missing infrastructure produces no indicator but does not prevent the daemon or GUI from running.
