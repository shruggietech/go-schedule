# S107 isolation contract

Production defaults retain native secrets, local configuration, services and launcher. Explicit Go dependencies replace them in a focused harness; no production environment-variable credential downgrade. Independent roots cannot change each other's client/output/callbacks.

Explicit --ipc-path is a valid absolute local endpoint, incompatible with remote flags. Missing/malformed selectors fail without fallback. Local persistence commands never contact a daemon.

Before every packaged fixture mutation, validate unique IPC, installation identity, exact PID and effective config/state/database/log/lock paths against a resolved owned UUID workspace. Reject ownership mismatch and symlink/junction escapes. These checks belong to fixtures, not ordinary CLI product use.

Recording service/launcher/secret adapters exercise actual handlers. Windows child execution is hidden, redirected and noninteractive, with job containment before execution; Unix processes use an owned group. Cleanup is bounded and exact, retains evidence and never changes installed services, real user stores, desktop/session or global WSL configuration.

The future #289 runner consumes this boundary; complete scenario orchestration, resource budgets and canonical reports remain #289/#290.

Native credential reads use LoadContext and an owned hidden helper; tokens move only through redirected in-memory pipes. Request cancellation joins the helper and any Unix native-tool descendants. Injected backends must support cancellation-aware reads; legacy blocking Get is never a fallback. Cleanup has its own bounded five-second drain deadline and errors are reported separately from request cancellation.

Unix Stop waits for the complete owned process group to disappear. Default macOS workspaces use a short private UUID parent; an explicit parent that exceeds socket capacity fails before owned state is created.
