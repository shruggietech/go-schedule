# Windows local service UI contract

The tray tooltip and menu show a single local service display state. The icon uses approved reduced branding legible against light, dark, and high-contrast taskbars. The context menu offers Open, status/details, available Start/Stop/Restart actions, and Quit companion. Stop/Restart require confirmation; operation progress is distinct from final success.

Open activates an existing GUI in the current user session or starts one. A remote-selected GUI must continue to show the local Windows service card independently of the remote daemon connection.

The GUI control and tray control use identical state names and operation semantics. A service started by standalone GUI autospawn is distinct from an installed Windows service; an installed service is never auto-restarted merely by opening the GUI.

An elevation cancellation, SCM error, health timeout, or operation timeout presents an actionable failure/cancellation state. No success message may be emitted before the observed target state.
