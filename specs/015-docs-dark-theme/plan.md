# Implementation Plan: Documentation Dark-Theme Quality

**Branch**: `codex/015-docs-dark-theme` | **Date**: 2026-08-27 |
**Spec**: `specs/015-docs-dark-theme/spec.md`

## Summary

Close #36, #37, and #35 as one documentation-site quality slice. Make syntax highlighting safe by default, retain deliberate brand-role accents, correct line-highlight and selection colors, enforce a four-category fence vocabulary, and align the sidebar endorsement with just-the-docs' responsive navigation gutters. Extend the existing offline docs gate before changing the sources so the known defect produces a red regression signal first.

## Technical Context

**Language/Version**: libsass-compatible SCSS, POSIX shell, Markdown

**Primary Dependencies**: GitHub Pages' Jekyll 3.9 environment and `just-the-docs@v0.4.2` (unchanged)

**Storage**: None

**Testing**: `sh scripts/docs-check.sh`, then `sh scripts/verify.sh all`

**Target Platform**: GitHub Pages documentation site at desktop and mobile viewport widths

**Project Type**: Branch-published static documentation site within the Go application repository

**Performance Goals**: No additional runtime scripts, requests, fonts, or generated assets

**Constraints**: Dark-only site; libsass syntax; no theme upgrade; offline validation; `scripts/docs-check.sh` is a pinned process artifact

**Scale/Scope**: 11 published Markdown pages, one SCSS customization file, one offline documentation gate, and one fence normalization

## Constitution Check

| Gate | Result | Evidence |
| --- | --- | --- |
| Code quality | PASS | The safe-default rule replaces an incomplete class allowlist; the gate documents failures with locations. |
| Testing | PASS | Regression assertions land before the style and fence fixes and must fail against the current defect. |
| UX consistency | PASS | One accessible palette, one fence vocabulary, and responsive gutter variables apply across every page. |
| Performance | PASS | Static CSS and offline checks add no page requests or runtime execution. |
| Autopilot | PASS | The slice is traceable to #36, #37, and #35 and halts once before publication. |
| Integration | PASS | Work remains on a review branch and the PR will close all three issues. |
| Pinned artifacts | PASS | The docs gate change is required by FR-007 and receives a dated CHANGELOG decision. |

No constitution exception is required.

## Project Structure

### Documentation (this feature)

```text
specs/015-docs-dark-theme/
├── checklists/
│   ├── accessibility.md
│   └── requirements.md
├── contracts/docs-theme.md
├── data-model.md
├── plan.md
├── quickstart.md
├── research.md
├── spec.md
└── tasks.md
```

### Repository files

```text
docs/
├── _sass/custom/custom.scss
└── gui-fields.md

scripts/docs-check.sh
CHANGELOG.md
CLAUDE.md
```

**Structure Decision**: Keep the complete theme in the existing custom SCSS entry point and extend the existing docs gate. A new partial or test framework would add indirection without improving this small surface.

## Design

1. Extend `scripts/docs-check.sh` with a static theme contract and a Markdown fence scanner. The gate requires the safe token fallback, named role groups, dark highlighted-line background, legible selection colors, and the approved `sh|bash|powershell|text` vocabulary.
2. Run the docs gate against the existing repository and retain the expected red evidence for missing safe fallback, selection/highlight treatment, and the untagged `gui-fields.md` block.
3. In `custom.scss`, put `[class] { color: inherit; }` before intentional token mappings. This neutralizes any inherited light-theme foreground, including future Rouge classes, while later equal-specificity role rules supply brand accents.
4. Explicitly map the previously leaked literal/name/decorator/subheading classes, use a dark `.hll` background, and override code `::selection` with Anchor Blue plus Panel ink. All selected pairs exceed 4.5:1.
5. Change the one untagged published block to `text`; retain existing `sh`, `bash`, `powershell`, and `text` fences because they already match their content. Do not add markup to force lexer token density.
6. Set endorsement spacing from theme variables: `$gutter-spacing-sm` on small viewports, `$gutter-spacing` from the `md` breakpoint upward, with `$sp-3` top and bottom. These are the exact navigation insets used by v0.4.2.

## Post-Design Constitution Re-check

PASS. The design replaces fragile enumeration with a safe default, adds a test-first offline contract, uses the upstream theme's existing spacing tokens, and preserves all product and hosting behavior.

## Pinned-Artifact Decision

`scripts/docs-check.sh` is intentionally expanded because static style and fence regressions cannot be caught by its current link/front-matter checks. The change remains POSIX-only and network-free and will be recorded in `CHANGELOG.md` dated 2026-08-27.
