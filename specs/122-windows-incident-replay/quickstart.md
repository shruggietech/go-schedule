# S122 validation guide

Use the verified hidden noninteractive launcher. Prepare exact-source binaries through existing gosched-qa build before campaign execution. Keep original S112/S115/S117/S118 and S121 artifacts unchanged.

Runtime entry points: TestScriptsBindingsKeepConcurrentSessionIdentities and TestScriptsRuntimeAndSelectionContrasts, selecting powershell with current and retained child PATH separately, then posix. Use race, JSON, failfast and at most count=10. Run deliberate prerequisite fault matrices once. Set GOSCHEDULE_PROBE_DIAGNOSTICS=1 and a fresh canonical GOSCHEDULE_INCIDENT_DIR for each group.

IPC entry point: TestWindowsIncidentReplay, opt-in GOSCHEDULE_INCIDENT_REPLAY=1 and GOSCHEDULE_REPLAY_BIN pointing to the source-bound prepared input directory. Exercise memory/WAL component fixtures and packaged-daemon WAL separately with at most ten repeats. Preserve unchanged catalog assertions and fresh cleanup/incident ownership. The maintainer guide supplies exact runnable commands and fixed report fields.

Before claiming a correction, retain a failing existing-unlaunchable-executable or MCP lifecycle cancellation regression and rerun it after the smallest demonstrated fix. Then run the full foreground sh scripts/verify.sh all and, after application changes, the existing full Windows/WSL evaluation/finalization process. Two Codex rounds maximum; stop before merge/tag/promotion.
