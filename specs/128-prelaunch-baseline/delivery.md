# S128 delivery

Status: source implementation and local verification complete; publication and owner integration pending.

## Delivered source outcomes

Direct application schema 24 replaces historical replay while valid current data/identity reopen. Current preferences have no Fyne import or transition surface. Both metrics writers use complete format 4 with bounded non-mutating preflight. Historical MSI/profile qualification is retired in favor of current evidence revision 2, while repair, maintenance, removal and security/identity assertions remain. Constitution 4.0.0 and active guidance implement the owner-approved prelaunch policy.

| Issue | Source outcome | Integration state |
| --- | --- | --- |
| #408 | Migration-free current baseline and policy | Ready in this review branch, owner merge pending |
| #409 | Deterministic one-off preview clock/DST expectations | Ready in this review branch, owner merge pending |
| #410 | Current metrics read-only contention and refusal diagnostics | Ready in this review branch, owner merge pending |

[Verification](verification.md) records all eight foreground gates, focused checks, preserved first failures and zero drift across 2,516 source files and 209 historical inputs. Five factual records were finalized afterward, plus removal of one trailing SQL blank line with identical trimmed DDL and a focused storage race recheck.

## Authorization and remaining boundary

The owner explicitly authorizes push, one official PR and verified in-scope CI/review corrections. At most two Codex requests include PR-open automation. CI and every received finding must be reconciled before owner handoff. Merge, auto-merge, tag and release are not authorized. No scheduled monitor is created.

SQLite #400/#404, original #355/#356 causes and clustered execution remain independent/deferred. The previous conditional S128 dependency allocation on #404 is superseded by #408. Published v1.7.0 and immutable historical records are unchanged.
