# Verification: Desktop Interaction and Appearance Polish

## 2026-09-03 - Foundational red phase

Command:

```powershell
go test ./gui -run "Appearance|Theme|Options|Storage|Navigation|Scroll|Info|Shutdown" -count=1
```

Result: expected failure. The S043 production code does not define the scroll
sensitivity preference, Inter or Ubuntu choices, sensitive vertical scroll,
selector-alternative state, compact storage-row fields, or navigation boundary.
The compiler stopped first at those missing preference and font symbols. This
preserves objective evidence that the new tests preceded the implementation.

## Green evidence

Commands:

```powershell
go test ./gui -run "Appearance|Theme|Options|Storage|Navigation|Scroll|Info|Shutdown|CursorButtonKeeps" -count=1
go test -race ./gui -run "Appearance|Theme|Options|Storage|Navigation|Scroll|Info|Shutdown|CursorButtonKeeps" -count=1
go test ./... -count=1
```

Result: passed. The focused suite proves preference defaults and migration,
font resources, current-choice exclusion, compact and unavailable storage rows,
small-viewport wrapping, semantic composite contrast, navigation geometry,
live discrete-wheel scaling, precision-delta preservation, Info construction,
representative control identity, and one-shot shutdown. The same focused set
passed with the race detector, and the complete Go test suite passed.

## Canonical gate

Command:

```powershell
& 'C:\Program Files\Git\bin\bash.exe' scripts/verify.sh all
```

Result: passed on 2026-09-03. All eight canonical gates completed successfully:
format, vet, lint (zero issues), race, GUI, coverage, docs, and automation.

## Systematic audit

The final local review covered the six authoritative GitHub issue bodies,
requirement-to-test traceability, Fyne renderer state composition, preference
and callback concurrency, repository history and scope, embedded-font licensing
and hashes, shutdown behavior, storage ownership wording, and evidence claims.
It corrected the storage text widget to use Fyne's selectable Label contract,
paired bright standalone Stop Red with dark text on filled danger controls so
both indicator and text thresholds survive every state, kept scroll-only changes
from reinstalling the theme, added a long-path small-window layout regression,
reflows measured storage columns after font changes, and removed a false
system-native font claim. No unresolved high-confidence code or evidence
finding remains.

## Integrity

All 33 changed text files decode as strict UTF-8, contain no UTF-8 BOM or common
mojibake markers, and `git diff --check` reports no whitespace errors. The six
font/license files match the SHA-256 values pinned in `research.md`.

## Native evidence boundary

The operator's S043 attended A/B observation established that choosing System
made the fuzzy normal body text sharp on the QHD ultrawide Windows display. It
supports System as S044's new default. It does not prove the after-change
candidate. Native contrast, font sharpness, DPI, conventional-wheel, and
precision-touchpad evidence remains assigned to exact-candidate qualification;
headless results will not be represented as native proof.
