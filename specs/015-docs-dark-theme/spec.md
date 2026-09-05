# Feature Specification: Documentation Dark-Theme Quality

**Feature Branch**: `codex/015-docs-dark-theme`

**Created**: 2026-08-27

**Status**: Implemented

**Delivery**: [PR #47](https://github.com/shruggietech/go-schedule/pull/47)

**Traceability**: Closes GitHub issues [#36](https://github.com/shruggietech/go-schedule/issues/36), [#37](https://github.com/shruggietech/go-schedule/issues/37), and [#35](https://github.com/shruggietech/go-schedule/issues/35).

## Overview

The documentation site uses a dark branded code surface, but some syntax tokens still inherit a light-theme palette and become nearly unreadable. Code examples also use inconsistent fence conventions, and the sidebar endorsement sits against the container edge. This slice makes the documentation theme complete, accessible, consistent, and resilient to new syntax token classes.

### Scope in

- Ensure every syntax token has a safe dark-theme fallback.
- Preserve a small, deliberate brand palette for meaningful token roles.
- Make highlighted text selection legible.
- Establish and enforce one documentation fence vocabulary.
- Add regression checks for the theme and fence contract.
- Align the sidebar endorsement with surrounding navigation spacing.

### Scope out

- Replacing or upgrading the documentation-site theme.
- Adding a light-theme toggle.
- Artificially coloring command words that the selected lexer does not classify.
- Rewriting example content or changing documented commands.
- Changing the application, daemon, API, CLI, packaging, or release behavior.

## User Story 1 - Read every code token (Priority: P1)

As a documentation reader, I want every part of a code example to remain legible on the dark code surface so identifiers and commands do not disappear.

**Why this priority**: Current identifiers can render at 1.58:1 contrast, making substantive PowerShell examples unreadable and failing accessibility guidance.

**Independent Test**: Review representative PowerShell, shell, text, and future-token examples against the dark code surface and confirm every text role meets the contrast floor, including selected text and unrecognized token types.

### Acceptance Scenarios

1. **Given** any syntax token class, **when** it has no special brand role, **then** it uses the readable base code color rather than a light-theme inherited color.
2. **Given** a named token role such as comment, keyword, string, number, variable, inserted text, or deleted text, **when** it is highlighted, **then** it uses the documented dark-theme brand mapping.
3. **Given** highlighted code text, **when** a reader selects it for copying, **then** foreground and selection-background colors remain legible.
4. **Given** a highlighted-line token, **when** it appears, **then** it uses a dark background compatible with the surrounding code surface.

---

## User Story 2 - Encounter consistent code examples (Priority: P2)

As a documentation reader, I want equivalent example types identified consistently so highlighting differences communicate content type rather than accidental authoring drift.

**Why this priority**: A total palette fixes legibility first; a documented fence convention then keeps presentation stable across existing and new pages.

**Independent Test**: Scan all published documentation pages and confirm every fenced block declares one approved content type and equivalent examples use the same category.

### Acceptance Scenarios

1. **Given** a POSIX command example, **when** it is authored, **then** it uses the POSIX shell category.
2. **Given** a Bash-specific script, PowerShell example, or non-executable output/grammar example, **when** it is authored, **then** it uses the matching approved category.
3. **Given** an untagged or unsupported fenced block, **when** documentation validation runs, **then** validation fails with the file and fence category.
4. **Given** a shell lexer leaves bare command text unclassified, **when** the page is displayed, **then** that text remains intentionally legible in the base code color without invented tokenization.

---

## User Story 3 - Read a comfortably spaced endorsement (Priority: P2)

As a documentation reader, I want the sidebar endorsement aligned with nearby navigation content so it does not look clipped or unfinished.

**Why this priority**: It is a small, visible defect in the same documentation stylesheet and can be corrected without expanding the slice.

**Independent Test**: Inspect the sidebar at desktop and mobile widths and confirm the endorsement has non-zero horizontal and bottom spacing aligned with the navigation inset.

### Acceptance Scenarios

1. **Given** the desktop sidebar, **when** the endorsement appears, **then** it aligns with the horizontal inset of nearby navigation content and has bottom breathing room.
2. **Given** a mobile-width sidebar, **when** the endorsement appears, **then** the same alignment and bottom spacing remain.

### Edge Cases

- A newly introduced syntax token class receives readable base ink without a stylesheet update.
- Text selection remains legible inside nested highlighted spans.
- A highlighted-line background never falls back to a light plate.
- Untagged fences and plausible but unsupported aliases are rejected rather than silently accepted.
- Historical specifications outside the published documentation are not subject to the fence vocabulary.

## Requirements

### Functional Requirements

- **FR-001**: Every code token MUST have a dark-theme fallback with at least 4.5:1 contrast against the code surface.
- **FR-002**: Comments, keywords/functions, strings/tags, numbers/variables, operators/punctuation, prompts, inserted text, deleted text, and lexer error tokens MUST have one documented, consistent role mapping.
- **FR-003**: Highlighted-line backgrounds and selected code text MUST remain dark-theme legible at a minimum 4.5:1 foreground/background contrast.
- **FR-004**: Published documentation MUST use only these fence categories: `sh` for portable POSIX commands, `bash` for Bash-specific scripts, `powershell` for PowerShell, and `text` for output or non-executable examples.
- **FR-005**: Every published fenced block MUST declare a category; untagged and unsupported categories MUST fail offline documentation validation with an actionable location.
- **FR-006**: Bare command text not classified by a lexer MUST remain in the readable base code color; the feature MUST NOT add artificial markup solely to increase token density.
- **FR-007**: Offline validation MUST guard the safe token fallback, highlighted-line background, selection colors, named role mappings, and fence vocabulary.
- **FR-008**: The sidebar endorsement MUST have non-zero horizontal and bottom spacing aligned with the existing sidebar navigation inset at desktop and mobile widths.
- **FR-009**: The feature MUST preserve the current dark-only site, theme pin, documented commands, page content, and application behavior.

## Success Criteria

### Measurable Outcomes

- **SC-001**: All syntax foreground/background combinations defined by the documentation theme meet or exceed 4.5:1 contrast.
- **SC-002**: A previously unknown token class displays with the readable base code color without requiring a new class-specific rule.
- **SC-003**: All fenced blocks in the 11 published pages declare exactly one of the four approved categories.
- **SC-004**: Offline validation fails for a missing safe fallback, missing selection/highlight treatment, missing role mapping, untagged fence, or unsupported fence category.
- **SC-005**: The endorsement has horizontal and bottom spacing on both desktop and mobile layouts, with no viewport-specific exception needed.
- **SC-006**: The full local verification aggregate passes with no application, dependency, theme-version, or hosted-workflow changes.

## Assumptions

- WCAG AA's 4.5:1 normal-text threshold is the appropriate floor because code examples are not guaranteed to render as large text.
- Consistency means a complete palette and a stable content vocabulary, not equal token density across lexers.
- The current four fence categories cover every existing published example.
- The existing sidebar navigation inset is the spacing authority; this feature does not introduce a new responsive layout system.
