# S108 product findings and publication drafts

These 12 independently actionable product findings come from the completed Windows/WSL capture and agent communication review. Their implementation is outside S108. The canonical evaluation intentionally reports failure, with 590 execution failures and 76 additional cases affected by communication findings; every applicable case ran. The runner/report capability for #289/#290 is delivered.

All files below are local publication drafts. The owner selected the normal autopilot pre-publication boundary, and CLAUDE.md requires halting before publication. No new GitHub issue has been filed and no issue has been closed by this work.

| Finding | Proposed issue | Affected case/environment combinations |
| --- | --- | --- |
| F001 | [Use the documented usage exit code for parser and local validation failures](F001.md) | 550 |
| F002 | [Reject unexpected operands before executing CLI commands](F002.md) | 26 |
| F003 | [Keep alerts deprecation notices out of JSON stdout](F003.md) | 2 |
| F004 | [Separate webhook endpoint flags from remote target selection](F004.md) | 8 |
| F005 | [Reject watcher path bindings for unknown bundle identities](F005.md) | 2 |
| F006 | [Return failure when a cron dry-run report contains invalid input lines](F006.md) | 2 |
| F007 | [Replace missing and placeholder command help with useful descriptions](F007.md) | 28 |
| F008 | [Identify affected objects and invalid fields in CLI diagnostics](F008.md) | 24 |
| F009 | [Avoid empty actor identity scaffolding in credential rotation output](F009.md) | 4 |
| F010 | [Label columns in human-readable operational lists](F010.md) | 12 |
| F011 | [Explain pending watcher health after an enable or disable operation](F011.md) | 4 |
| F012 | [Explain the IPC compatibility-mode warning and its configuration context](F012.md) | 8 |

The [publication manifest](issues.json) contains exact canonical titles, existing labels, milestone 15, parent #277, body filenames and null existing-issue links. Match existing open/closed records before each create, publish idempotently only after authorization, update #277 with actual child links, and reconcile finding IDs with the resulting issue URLs. Drafts have functional acceptance criteria, dependencies, source references and representative bounded evidence. They do not require CI, QA packets or human walkthroughs for closure.

The complete local artifact package is A:/Code/go-schedule/.test-temp/s108-run-17. Its report.json and deterministic report.md include all case references, the 1,040 complete communication reviews, exact six-binary identities, four successful cleanup actions and no remaining owned resources. The original inputs are A:/Code/go-schedule/.test-temp/s108-inputs-17.
