# S109 research and decisions

## 2026-10-21: Invocation and targeting

Decision: use typed usage causes at parser/operand/document boundaries and retain operational errors. Rationale: matching arbitrary error text or marking every command error as usage would conceal genuine runtime failures. Alternatives rejected: broad exit-code rewriting and runtime-message substring classification.

Decision: add global --daemon-endpoint while retaining legacy global --endpoint and channel-local destination --endpoint. Inspect global selector ownership, and reject simultaneous remote alias flags. Rationale: a same-name child flag is currently confused with a remote selector in invocation.resolve. A distinct selector permits both explicit remote daemon and destination in one command. Alternatives rejected: breaking existing webhook scripts by removing --endpoint, or fixing only IPC while leaving explicit remote authoring inexpressible.

Decision: enforce supplied watcher binding identities at the shared service before plan creation and duplicates before CLI map conversion. Rationale: CLI-only unknown checks leave API consumers with the same defect; a map cannot retain duplicate evidence. Do not create a plan and reject it afterward.

## 2026-10-21: Reports and diagnostics

Decision: preserve the complete cron report, then return usage failure for invalid lines, preserving task-creation and I/O runtime failures. Unsupported lines retain their existing classification. Alternatives rejected: returning after the first invalid line or treating a mixed result as success.

Decision: enrich typed API errors at the CLI boundary with object kind/identity and field/correction, and add explicit safe context to local parser/certificate failures. Rationale: shared API consumers keep structured error codes and CLI users get actionable context. Never use credential or authorization payloads as diagnostic context.

## 2026-10-21: Association and observation

Decision: populate real actor and daemon context for credential rotation while preserving IssuedCredential's shape and one-time token. Rotation currently constructs only credential/token. Gather fallible context before mutation; do not authenticate the newly rotated token merely to populate the result because authentication imposes expiry/revocation semantics and updates last-used state.

Decision: add observed-definition context to watcher health and compare it with the saved definition. Existing reload requests are asynchronous/coalesced. Health ChangedAt is a transition timestamp, not proof that an unchanged state reflects the new definition. Reconciliation must be correlated to the actual definition rather than guessed from health text or a recent timestamp. Alternatives rejected: blocking until reload, claiming immediate readiness, or equating unchanged health with stale state forever.

## Sources and unresolved work

Primary sources are the merged repository, issue bodies #311-#322 and source-pinned S108 findings. The plan workflow's read-only research delegates inspect existing tests and desktop/API consumers before implementation. All technical unknowns have supported local decisions; implementation may refine field names while preserving the written contracts.
