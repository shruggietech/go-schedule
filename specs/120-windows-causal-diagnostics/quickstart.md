# S120 focused validation guide

Prerequisites: existing native Windows Go/PowerShell/Bash/SQLite and nonroot WSL setup, verified hidden launchers, retained original receipts. Do not replace runtimes or operate the installed service. Keep experiment outputs under the ignored .test-temp/s120 directory.

Run targeted script process/binding/probe contracts and store/API cancellation contracts through the verified launcher, preserving the first before result separately from corrected results. The process fixture must distinguish original exit 0/nonzero from held output collection and clean its owned descendants. The pool fixture must demonstrate a real queued request, cancel it and observe server completion before releasing the pool. Preserve native exits and all six-writer/partial-state/audit assertions.

Use GOSCHEDULE_PROBE_DIAGNOSTICS=1 for native metadata/writer/prerequisite follow-ups, and record executable/harness/source identity. Existing supported QA defaults remain authoritative. Consult [research.md](research.md) for primary references and [contracts/causal-diagnostics.md](contracts/causal-diagnostics.md) for interpretation. Record unavailable collectors and unconfirmed original causes explicitly.

Freeze completed source inputs and run the canonical foreground aggregate:

```sh
sh scripts/verify.sh all
```

Only after all eight gates and final metadata audits pass, commit/push and open the owner-authorized PR. Hosted checks and at most two Codex rounds follow; stop before merge, tagging or promotion.
