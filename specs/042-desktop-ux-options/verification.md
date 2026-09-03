# Verification: Desktop UX and Options

## 2026-09-03 - Foundational red phase

Command:

```powershell
go test ./gui -run "Appearance|Theme|Storage|EditorOwnership|ShutdownCoordinator" -count=1
```

Result: expected failure. The test package reported undefined appearance,
storage-location, editor-ownership, and shutdown-coordinator symbols. This proves
the S041 baseline does not implement the new shared contracts before production
code is added.

## Green evidence

Focused headless command:

```powershell
go test ./gui -run "Appearance|Theme|Options|Storage|Navigation|TaskRow|TaskList|TaskEditor|Shutdown|Info" -count=1
```

Result: PASS. This covers all nine appearance combinations, invalid-value
fallback, immediate application and persistence, storage resolution and copy,
Info label invariants, navigation order and geometry, Activity badge updates,
single/double task activation across reorder/removal, degraded detail lookup,
editor ownership, and one-shot shutdown.

Focused race command:

```powershell
go test -race ./gui -run "AppearancePreferences|ResolveStorage|StorageExistence|EditorOwnership|Shutdown" -count=1
```

Result: PASS. The non-widget preference, path, editor-ownership, backend-context,
and shutdown coordination state passed under the race detector.

Full headless GUI command:

```powershell
go test ./gui -count=1
```

Result: PASS in 25.968 seconds after the final behavior corrections.

## Canonical gate

The first full attempt passed format and vet, then static analysis rejected the
deprecated window-scoped clipboard call. After switching to the application
clipboard, the next attempt passed format, vet, lint, race, GUI, coverage, and
documentation. Automation then correctly rejected a Draft/In Progress lifecycle
mismatch. The specification and inventory were aligned without changing product
behavior. The authoritative rerun then passed all eight gates:

```text
format      PASS (no files reported)
vet         PASS
lint        PASS (0 issues)
race        PASS
gui         PASS (25.839 seconds)
coverage    PASS (engine 86.4%, schedule 89.2%, timezone 91.3%,
                  store 84.1%, catchup 88.9%, logbus 91.1%)
docs        PASS (15 pages; links, front matter, fences, theme, and policy clean)
automation  PASS
```

## Local review

The requirements, complete diff, recent integration history, concurrent state,
filesystem boundaries, Fyne gesture dispatch, and native-evidence claims were
reviewed locally. Corrections made before the final gate:

- forwarded a row's single tap back to stable-ID list selection because Fyne
  sends gestures to the deepest tap-capable child;
- added distinct configuration and per-user application-root storage rows;
- prevented external configured logs and development executables from being
  represented as wipe-owned;
- centralized the uninstall cleanup-evidence path instead of copying its GUID;
- derived backend call contexts from the application's cancelable run context;
- replaced the deprecated clipboard call and aligned Spec Kit lifecycle state;
- removed publication bookkeeping from the implementation task checklist.

No unresolved high-confidence correctness, security, data-loss, concurrency,
or requirements finding remains. The expected native visual risk is retained
explicitly below.

## Integrity

All 42 changed or added files decoded as strict UTF-8, none had a UTF-8 BOM, the
mojibake scan found no suspect sequence, and `git diff --check` passed.

## Native evidence boundary

Headless evidence in S042 will not be described as native Windows text or DPI
proof. Exact-candidate Windows 11 observations remain required by #94 before
the visual acceptance criteria in #101, #104, and #105 close.
