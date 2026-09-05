# Research: Documentation Dark-Theme Quality

## Decision 1: Safe token fallback

**Decision**: Reset the foreground of every classified descendant to inherit the readable code-surface base color, then apply intentional role colors after that fallback.

**Rationale**: The current allowlist misses token classes and fails open to the light theme's near-black ink. A safe default remains correct when Rouge adds a class or a new lexer emits an existing unlisted class.

**Alternatives considered**: Enumerating every current Rouge class remains fragile at theme or lexer changes. Upgrading just-the-docs is incompatible with the branch-based GitHub Pages/libsass constraint.

## Decision 2: Contrast palette

**Decision**: Retain Panel `#0d171c` with base `#f3f7f8`, Muted `#9baeb6`, Anchor Blue `#58a6ff`, Interval Mint `#62d9b7`, Hold Amber `#f2b84b`, and Stop Red `#e05f5f`. Use Anchor Blue selection with Panel foreground and a dark `#17262d` highlighted-line surface.

**Rationale**: Existing issue measurements place the role colors between 5.2:1 and 16.8:1 against Panel, all above the 4.5:1 floor. Panel foreground on Anchor Blue also exceeds the floor. The new highlighted-line surface stays dark and retains strong base-text contrast.

**Alternatives considered**: White-on-white inherited selection and the light `#f9f9f9` highlighted-line plate are unreadable. Introducing new colors is unnecessary because the established brand palette already passes.

## Decision 3: Fence vocabulary

**Decision**: Allow `sh`, `bash`, `powershell`, and `text`; require a language on every published fence.

**Rationale**: These four categories describe all existing examples. `sh` denotes portable commands, `bash` denotes Bash-specific scripts, `powershell` denotes Windows shell examples, and `text` denotes output or grammar. Lexer coverage may be sparse, but base ink makes unclassified text intentionally legible.

**Alternatives considered**: Converting command examples to `console` would add prompts that readers might copy and does not solve palette correctness. Treating all commands as `text` discards useful option/string highlighting.

## Decision 4: Responsive endorsement inset

**Decision**: Use just-the-docs v0.4.2's `$gutter-spacing-sm` below `md` and `$gutter-spacing` at and above `md`, with `$sp-3` vertical spacing.

**Rationale**: The upstream navigation rules use those exact variables for link and category insets. Reusing them produces alignment without arbitrary pixels or a parallel responsive system.

**Alternatives considered**: A fixed inset aligns at only one breakpoint; hard-coded pixel values duplicate upstream layout knowledge.
