# S138 investigation and approved disposition

## 2026-11-08 recovered historical event

Read-only Application log inspection recovered .NET Runtime event 1023, record 229042, at 2026-10-22T17:43:35.8766937Z. It reports application pwsh.exe, PID 479476, CoreCLR 10.0.1226.41902 and .NET 10.0.12, with native fault offset 0x1296dc relative to the recorded module base. This is near the original #355 writer failure timestamp. The original writer log contains no PID, so exact writer attribution is unproven. This event does not recover a module hash or an original dump.

Exact-window PowerShell Operational and Security process-creation queries found no matching retained records. User crash-dump and Windows WER archive/queue inventories contained no original PowerShell artifact. Raw event XML and query receipts remain private; machine/account data is not published. All historical evidence remains unchanged.

The nearest upstream custom-attribute report and the newer collectible-assembly lifetime correction differ in evidence and execution context. No supported upstream resolution or shared #355/#356 cause is established.

## Demonstrated diagnostic gap

The original stack is in metadata/type resolution before authored script execution. The existing identity prelude therefore cannot cover every early crash, while a dump can also be unavailable. The new controlled pre-identity probe failed first because runtime-events.json was missing, with the native failure preserved. This is the independently tracked diagnostic gap #441, not an original CLR recurrence.

The corrected probe produced controlled exit 0x80131623, no runtime.json, and an event-1025 receipt for exact child PID 61120. Its native creation/exit interval was 2026-11-08T22:00:40.0732568Z through 22:00:40.556016Z; the matching event timestamp was 22:00:40.2644034Z. Event readback took 57 milliseconds and reported CoreCLR 10.0.1426.52326/.NET 10.0.14. This initial controlled proof is separate from final frozen-source verification. Native dump status and the original operation remain independent of event evidence.

## Remaining outcomes

The owner approved this bounded diagnostic improvement after reviewing the preflight limitation. #441 can be delivered by the new capture behavior and guidance. #355 and #356 remain open with needs: reproduction because their original causal outcomes are still unconfirmed. No original failure was retried into success, no broad unchanged replay was performed, and no installed runtime/service, global event/dump policy, tag or release was changed.
