# Quickstart: Validate Persisted Adjustable Columns

## Automated validation

```bash
go test ./gui
go test -race ./gui -run "Test(ColumnProfile|ProfileWidths)"
sh scripts/verify.sh all
```

Expected: focused GUI tests prove the [column-layout contract](contracts/column-layout-contract.md),
the pure profile/allocation race run is clean, and all eight canonical gates
pass. The repository intentionally excludes the full Fyne widget package from
its race gate because of documented upstream font-cache behavior.

## Headless interaction scenarios

1. Build Schedule and Activity with distinct in-memory preferences.
2. Drag each boundary and operate it with arrow keys.
3. Assert adjacent-only change, conservation, minimums, and row alignment.
4. Recreate views at wide, default, and narrow widths.
5. Load wrong-version, wrong-schema, malformed, non-finite, and impossible data.
6. Customize both views, reset each independently, and recreate them.

## Native visual check

Inspect both views at default and narrow sizes, in light and dark modes and each
font. Separators and focus remain visible without excessive contrast, **When**
has a practical default share, rows align, and no horizontal scrollbar appears.
