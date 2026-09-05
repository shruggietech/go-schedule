# Quickstart: GUI Activity Clarity

## Focused red-green loop

From the repository root, add the Activity label, badge-boundary, and control presentation tests first, then run:

```sh
go test ./gui/...
```

The new expectations must fail before production behavior changes. Implement the smallest change, rerun the focused GUI tests, and keep existing merge, filter, cutoff, and acknowledgement tests green.

## Full verification

Run the canonical aggregate in the foreground:

```sh
sh scripts/verify.sh all
```

Expected gates, in order: format, vet, lint, race, gui, coverage, docs, and automation. Do not publish until every gate passes and the autopilot halt is explicitly authorized.

## Manual review points

- The fourth navigation tab reads `Activity`.
- Counts display exactly through 99 and cap at `99+`.
- The action reads `Clear View` and uses a clear-content icon.
- Visible help states that current activity is hidden, visible alerts are acknowledged, and records are not deleted.
- New activity after clearing still appears.
