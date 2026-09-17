# Security Checklist: MCP Operate Authority

- [x] Observe sessions cannot discover or invoke mutation tools.
- [x] Session secrets are random, returned once, stored only as digests, and never logged.
- [x] Unknown, expired, revoked, or malformed sessions fail closed.
- [x] Tool calls execute under an MCP actor rather than the built-in local actor.
- [x] Every daemon-dispatched attempt records actor, daemon, operation, target, and result; local schema rejection occurs before dispatch and returns a bounded structured result.
- [x] Wrong-daemon requests stop before task mutation.
- [x] Request identifier reuse cannot create duplicate work or change targets.
- [x] Operate cannot create, edit, delete, reveal secrets, or administer permissions.
- [x] Tool outputs contain no command, arguments, environment, stdin, credentials, or protected paths.
- [x] Hostile resource content cannot alter tool discovery or authority.
