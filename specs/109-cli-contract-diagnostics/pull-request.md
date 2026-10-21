# CLI contract and diagnostic repairs

CLI parser failures could return runtime status, extra operands could reach handlers, webhook destinations could shadow daemon targeting, and some reports and identity/state responses gave incomplete or misleading results. This change delivers all 12 S109 corrections in #311-#322: shared invocation contracts, independent --daemon-endpoint selection, strict watcher bindings, complete unsuccessful cron reports, clean JSON diagnostics, meaningful help and labeled lists, contextual errors, accurate atomic rotation responses, explicit watcher observation/reconciliation context across API/CLI/desktop, and an actionable IPC warning.

The existing unattended QA contracts are retained. The initialized inventory now has 152 nodes and 844 cases per environment, with historical case IDs preserved. Real capture also exposed an omitted channel-to-snapshot identity mapping in human-output QA; its correction retains stdout identity, unchanged-state and ambiguity assertions, and is tracked independently as #331, with source-pinned evidence in specs/109-cli-contract-diagnostics/findings/human-channel-identity.md.

Validation: all eight foreground local gates passed; the final Windows/WSL canonical evaluation passed 1,685 of 1,688 planned combinations with three manifest-defined exclusions, zero failures/blocked cases, 1,049 complete contextual reviews and no remaining owned resources. Focused native Linux race checks and Linux-target lint passed. Finalize/validate/render passed and repeated Markdown rendering was byte-identical. Detailed source identities, artifact hashes, earlier corrective runs, functional issue dispositions and platform observation boundaries are recorded in specs/109-cli-contract-diagnostics/verification.md. The captured subject is the receipt-bound dirty source digest recorded there, before subsequent delivery metadata; hosted CI is separate.

Decisions: reuse Cobra and typed error causes, retain legacy --endpoint compatibility, add observation fields without persisted schema changes, and assemble rotation identities before transaction commit without new eligibility rules. No new dependency or intentionally deferred functional outcome.

- Closes #311
- Closes #312
- Closes #313
- Closes #314
- Closes #315
- Closes #316
- Closes #317
- Closes #318
- Closes #319
- Closes #320
- Closes #321
- Closes #322
- Closes #331

Refs #277. Milestone: CLI completeness and unattended application QA. Spec-Kit slice: S109 (specs/109-cli-contract-diagnostics).
