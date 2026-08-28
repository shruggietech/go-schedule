# Documentation Theme Contract

## Published code blocks

- Every opening triple-backtick fence under `docs/*.md` declares exactly one of
  `sh`, `bash`, `powershell`, or `text`.
- Closing fences are plain triple backticks.
- Fence failures identify the source file, line, and unsupported or missing
  category.

## Syntax palette

- The code surface is Panel (`#0d171c`) with base ink (`#f3f7f8`).
- Every classified descendant inherits base ink before role mappings apply.
- Role mappings use only the documented base, muted, blue, mint, amber, and red
  colors.
- Highlighted lines use a dark background; selected code uses Anchor Blue with
  Panel foreground.
- The offline gate rejects absence or drift of these contract elements.

## Sidebar endorsement

- Horizontal padding equals the theme's small navigation gutter below `md`.
- Horizontal padding equals the theme's standard navigation gutter from `md`
  upward.
- Top and bottom padding use the same established spacing step.
