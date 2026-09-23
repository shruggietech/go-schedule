# S102 quickstart

1. Build the Linux desktop archive and install its binaries and `share/` assets using `docs/INSTALL-linux.md`. Install the supplied autostart entry and log in to a session with a compatible status host.
2. With the installed service running and local IPC healthy, inspect the item and Connections page: both say Running. Close the GUI; the item remains.
3. Use Stop with confirmation; the item and GUI say Stopped and explain that local scheduled tasks are not running. Reopen the GUI and confirm it does not auto-start the service.
4. Use Start and Restart, then exercise denied/cancelled authorization. Only observed success is reported; failed actions retain truthful status.
5. Launch the indicator and GUI repeatedly; only one item and one GUI appear. Remove the autostart entry and log out; the item does not return at next login.
6. In a session without a status host or bus, verify the daemon and GUI still work and the Connections page offers local service controls.

Automated parity: `sh scripts/verify.sh all`; Linux release packaging and hosted CI additionally build the indicator binary and inspect the archive contract.
