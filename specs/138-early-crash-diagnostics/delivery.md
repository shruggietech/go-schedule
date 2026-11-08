# Delivery: Early-crash incident evidence

## 2026-11-08 local completion

S138 implements [#441](https://github.com/shruggietech/go-schedule/issues/441) from accepted S137 main 8009da0809d77c8083ca4c5213f15ae71a2e2e96 on review branch codex/138-early-crash-diagnostics. Existing opt-in failed Windows PowerShell incident capture now retains a separate bounded runtime-events.json receipt, including failures before script identity executes. Exact process-handle times, PID, provider, event ID and executable identity constrain attribution. Primary native exit, cancellation, stream results and dump status remain unchanged.

The same PR also corrects [#443](https://github.com/shruggietech/go-schedule/issues/443), the macOS receipt fixture setup found by initial hosted CI, by reusing the existing canonical temporary-root helper. Production link rejection remains unchanged.

All nine implementation tasks are complete; [verification](verification.md) records the final full foreground qualification, fresh native proof, independent artifact readback and evidence audit. [Investigation](investigation.md) records the recovered historical event and its limits. #355 and #356 remain open because their original causal outcomes are unresolved. Issue closure is tied to each issue's functional outcome; verification activity is not a closure gate.

The owner explicitly approved this diagnostic scope and authorized automatic push, official PR creation, CI monitoring and in-scope findings resolution with at most two Codex review requests including automatic review. The default pre-push halt is therefore already authorized. Hosted results, exact candidate SHA and review disposition are maintained on the official PR. Final owner review and merge remain pending. No merge, tag, release or installed-service/runtime change was performed.
