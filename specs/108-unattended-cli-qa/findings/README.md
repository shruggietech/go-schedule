# S108 product findings and published issues

These 12 independently actionable product findings come from the completed Windows/WSL capture and agent communication review. Their implementation is outside S108. The canonical evaluation intentionally reports failure, with 590 execution failures and 76 additional cases affected by communication findings; every applicable case ran. The runner/report capability for #289/#290 is delivered.

The initial implementation stopped at the requested pre-publication boundary. The owner subsequently authorized PR and defect-issue publication on 2026-10-21. All 12 findings are now published as #311 through #322 under milestone 15 and parent #277; they remain open as independently closeable product outcomes.

| Finding | Published issue | Affected case/environment combinations |
| --- | --- | --- |
| F001 | [#311: Use the documented usage exit code for parser and local validation failures](https://github.com/shruggietech/go-schedule/issues/311) | 550 |
| F002 | [#312: Reject unexpected operands before executing CLI commands](https://github.com/shruggietech/go-schedule/issues/312) | 26 |
| F003 | [#313: Keep alerts deprecation notices out of JSON stdout](https://github.com/shruggietech/go-schedule/issues/313) | 2 |
| F004 | [#314: Separate webhook endpoint flags from remote target selection](https://github.com/shruggietech/go-schedule/issues/314) | 8 |
| F005 | [#315: Reject watcher path bindings for unknown bundle identities](https://github.com/shruggietech/go-schedule/issues/315) | 2 |
| F006 | [#316: Return failure when a cron dry-run report contains invalid input lines](https://github.com/shruggietech/go-schedule/issues/316) | 2 |
| F007 | [#317: Replace missing and placeholder command help with useful descriptions](https://github.com/shruggietech/go-schedule/issues/317) | 28 |
| F008 | [#318: Identify affected objects and invalid fields in CLI diagnostics](https://github.com/shruggietech/go-schedule/issues/318) | 24 |
| F009 | [#319: Avoid empty actor identity scaffolding in credential rotation output](https://github.com/shruggietech/go-schedule/issues/319) | 4 |
| F010 | [#320: Label columns in human-readable operational lists](https://github.com/shruggietech/go-schedule/issues/320) | 12 |
| F011 | [#321: Explain pending watcher health after an enable or disable operation](https://github.com/shruggietech/go-schedule/issues/321) | 4 |
| F012 | [#322: Explain the IPC compatibility-mode warning and its configuration context](https://github.com/shruggietech/go-schedule/issues/322) | 8 |

The [publication manifest](issues.json) contains exact titles, existing labels, milestone 15, parent #277, published issue identities and body filenames. The open/closed title inventory and targeted searches found no equivalent defects before publication. Parent #277 maintains the linked child checklist. The local records preserve functional acceptance criteria, dependencies, source references and representative bounded evidence. These outcomes do not require CI, QA packets or human walkthroughs for closure. The frozen canonical report retains its original null issue URLs; this manifest supplies the subsequent publication mapping without rewriting captured evidence.

The complete local artifact package is A:/Code/go-schedule/.test-temp/s108-run-17. Its report.json and deterministic report.md include all case references, the 1,040 complete communication reviews, exact six-binary identities, four successful cleanup actions and no remaining owned resources. The original inputs are A:/Code/go-schedule/.test-temp/s108-inputs-17.
