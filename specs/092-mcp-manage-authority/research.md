# Research: MCP Manage Authority

## Decision 1: Reuse capability ordering

Manage sessions use the existing `manage` capability. Because capability ordering already allows lower operations, a Manage server exposes Observe resources, the three Operate tools, and six additional Manage tools. Separate permission flags or parallel actor models would duplicate established authorization logic.

## Decision 2: Use one tool per definition family

Tasks, groups, completion chains, triggers, and filesystem watchers each receive one tool with a create, update, or delete action. Notification assignments receive one atomic replacement tool. This keeps discovery at six Manage tools and avoids a schema-free universal endpoint.

## Decision 3: Exclude sensitive task fields

Task environment variables and stdin are absent from MCP create and update schemas. Stored commands, arguments, paths, and credentials are never returned. Trigger create discards the API secret response and returns only the trigger identity. Human clients retain their existing full API surface.

## Decision 4: Make confirmation connection-scoped

Stdio and localhost HTTP accept an optional `require-confirmation` policy. When enabled, every Manage input must set `confirmed=true`; when disabled, a deliberately granted Manage session runs unattended. The policy is runtime-only and cannot be changed through a Manage tool.

## Decision 5: Reject bulk mutation

Each call performs one lifecycle operation or one atomic assignment replacement. This eliminates partial multi-object success. Transport ambiguity is returned as uncertain, never silently retried, and request UUID deduplication prevents accidental replay within one session.
