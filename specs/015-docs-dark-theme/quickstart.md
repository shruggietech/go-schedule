# Quickstart: Documentation Dark-Theme Quality

## Red regression step

After adding the theme and fence assertions but before changing SCSS or Markdown, run:

```sh
sh scripts/docs-check.sh
```

Expected: failure identifying the missing safe token fallback and selection/highlight contract, plus the untagged fence in `docs/gui-fields.md`.

## Focused green step

After implementing the theme, fence, and spacing changes, rerun:

```sh
sh scripts/docs-check.sh
```

Expected: all 11 published pages pass front matter, links, fence vocabulary, and dark-theme contract checks.

## Full verification

Run the canonical aggregate in the foreground:

```sh
sh scripts/verify.sh all
```

All eight gates must pass before the publication halt.

## Review points

- Unknown syntax classes inherit readable base ink.
- PowerShell identifiers no longer inherit near-black light-theme colors.
- Highlighted lines and selected code remain legible.
- Shell examples retain honest lexer behavior without artificial markup.
- The endorsement follows the sidebar gutter at small and desktop widths.
