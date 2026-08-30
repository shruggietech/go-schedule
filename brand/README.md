# go-schedule Brand System

**Status:** Approved identity, version 1.0.0  
**Parent:** ShruggieTech  
**Repository:** `shruggietech/go-schedule`  
**Audience:** Product, engineering, documentation, packaging, and communications

go-schedule is a cross-platform task scheduler written in Go. Its identity combines terminal familiarity, cron-field structure, and visible execution policy.

## Contents

- [Brand foundation](#brand-foundation)
- [Logo system](#logo-system)
- [Color](#color)
- [Typography](#typography)
- [Visual language](#visual-language)
- [Voice and writing](#voice-and-writing)
- [Parent brand relationship](#parent-brand-relationship)
- [Digital implementation](#digital-implementation)
- [Asset inventory](#asset-inventory)

## Brand foundation

### Governing principle: familiar scheduling, explicit behavior

People should recognize a scheduler immediately and still understand what the application will do. go-schedule keeps authored expressions visible, translates them into plain language, and surfaces the policies that change execution: timezone, catch-up, overlap, enabled state, and daemon lifecycle.

The identity follows the same rule. Familiar terminal and cron cues establish the category. Clear rails, anchored points, and restrained status colors make execution state legible.

### Positioning

**Category:** Cross-platform task scheduling.  
**Role:** A readable scheduler for service-managed automation.  
**Audience:** Developers, operators, maintainers, and technical users who want exact scheduling behavior across operating systems.  
**Functional descriptor:** Cross-platform task scheduling with a CLI, desktop application, and background daemon.

### Brand idea

**Know the next run.**

The product turns schedule intent into observable state. A task has an authored expression, an explanation, a next run when one exists, and an execution history. Event schedules state their lifecycle boundary instead of inventing a clock time.

### Personality

| Trait | Expression | Avoid |
| --- | --- | --- |
| Precise | Exact schedule expressions, timestamps, timezones, and policy names | Vague automation claims |
| Practical | Commands, examples, and observable outcomes | Aspirational productivity language |
| Calm | Measured status and sparse emphasis | Alarm-clock urgency |
| Transparent | Limitations and lifecycle boundaries stated plainly | Implied guarantees |
| Familiar | Terminal, cron, service, task, and run vocabulary | Decorative calendar metaphors |

### Brand promises

- Authored schedule intent remains visible.
- Plain-language explanations agree with execution behavior.
- Timezone, catch-up, overlap, and lifecycle rules are stated when relevant.
- Run state always includes a label or symbol in addition to color.
- Unknown or unavailable next runs are presented honestly.

## Logo system

### Mark construction

The mark contains a terminal prompt, a command line, and five schedule cells. The center cell uses Anchor Blue to identify an exact run point. The remaining cells use Interval Mint to represent recurrence. Neutral rails keep the lower structure readable at compact sizes.

All distributable logo SVGs use vector paths. Font-dependent text is outlined during the build, including the wordmark and schedule-cell asterisks. A recipient can open the assets without installing the kit fonts.

### Approved lockups

Use the horizontal lockup for repository headers, documentation, release pages, and wide surfaces. Use the stacked lockup for square compositions and title cards. Use the mark alone for application icons, avatars, and compact controls. Use the wordmark alone when nearby context already establishes the product.

The cover of the printable guide intentionally uses the transparent mark by itself. The page title supplies the product name, so a second logo image does not repeat it. The mark shares the exact page surface and has no mismatched square tile.

### Clear space

Maintain at least 32 units of clear space around the 512-unit mark. For lockups, extend that same proportional space around the complete artwork. Keep text, borders, icons, and crops outside this boundary.

### Minimum size

| Asset | Minimum digital width |
| --- | ---: |
| Full mark | 36 px |
| Reduced mark | 16 px |
| Horizontal lockup | 180 px |
| Stacked lockup | 104 px |
| Wordmark | 120 px |

At 32 px and below, use `go-schedule-mark-reduced.svg` or the supplied favicon exports. The reduced mark replaces the small asterisks and rails with five clear schedule cells while preserving the prompt and command line.

### Backgrounds

The primary presentation uses Interval Mint and Anchor Blue on Night. The light lockup uses the deeper accessible accents on Paper. White and black variants support single-ink reproduction and arbitrary compatible backgrounds.

Transparent lockups contain no background rectangle. Assets with `-dark` or `-light` in the name include their declared surface. This naming rule prevents a tile from being mistaken for transparent artwork.

### Prohibited treatments

- Do not rotate, skew, stretch, outline, bevel, or add glow.
- Do not recolor individual schedule cells outside the approved variants.
- Do not replace the outlined wordmark with live text.
- Do not place the mark on a nearly matching square tile inside another surface.
- Do not add clock hands, notification bells, speed lines, or decorative calendar pages.
- Do not combine the go-schedule and ShruggieTech marks into one logo.
- Do not crop inside the clear-space boundary.

## Color

The palette is dark-first and restrained. Interval Mint represents recurrence and ready state. Anchor Blue identifies exact points, links, and focus. Hold Amber covers policies or states that require attention. Stop Red is reserved for failure and destructive actions.

| Token | Hex | Role | Contrast |
| --- | --- | --- | ---: |
| Interval Mint | `#62D9B7` | Recurrence, ready state, identity | 11.09:1 on Night |
| Anchor Blue | `#58A6FF` | Exact run points, links, focus | 7.60:1 on Night |
| Hold Amber | `#F2B84B` | Catch-up, overlap, pending policy, warning | 10.73:1 on Night |
| Stop Red | `#E05F5F` | Failure and destructive actions | 5.46:1 on Night |
| Night | `#071014` | Primary background | - |
| Panel | `#0D171C` | Cards and code surfaces | - |
| Raised | `#13232A` | Elevated and selected surfaces | - |
| Line | `#28414B` | Rails, borders, inactive structure | - |
| Text | `#F3F7F8` | Primary dark-surface text | 17.81:1 on Night |
| Muted | `#9BAEB6` | Secondary dark-surface text | 8.35:1 on Night |
| Paper | `#F6F8F7` | Light reading background | - |
| Ink | `#132027` | Primary light-surface text | 15.58:1 on Paper |
| Light Muted | `#4A616B` | Secondary light-surface text | 6.13:1 on Paper |
| Light Interval | `#087A62` | Accessible recurrence accent | 4.96:1 on Paper |
| Light Anchor | `#0067C5` | Accessible link and focus accent | 5.26:1 on Paper |
| Light Hold | `#805500` | Accessible policy warning | 6.12:1 on Paper |
| Light Stop | `#B4232A` | Accessible failure state | 6.12:1 on Paper |

Every contrast value is re-derived from the token file during verification.

### Light-surface rule

Interval Mint measures 1.62:1 on Paper and cannot carry text there. Use Light Interval. Apply the corresponding light variants for Anchor Blue, Hold Amber, and Stop Red when those colors carry text or essential interface meaning.

### Color ratio and semantic use

Aim for roughly 80 percent neutral surfaces, 15 percent mint or blue, and no more than 5 percent amber or red. Keep large fills neutral. Use accents for relationships, selection, status, and focus.

Color never acts alone. Pair ready, pending, failed, paused, skipped, and running states with labels or distinct symbols.

## Typography

| Function | Typeface | Weights |
| --- | --- | --- |
| Display and headings | Space Grotesk | 500, 700 |
| Body and interface | Geist | 400, 500 |
| Schedules, paths, logs, metadata | Geist Mono | 400 |
| Product wordmark | Outlined Space Grotesk artwork | Fixed vector asset |

Geist ships regular and medium faces in this kit. Geist Mono ships regular only. Use color, spacing, or labels for emphasis in technical text instead of requesting an unavailable bold face.

Display headings use tracking near `-0.025em`. Body copy uses line height near `1.65`. Compact labels use Geist Mono with positive tracking. Schedule expressions, command lines, timestamps, paths, and log fragments disable ligatures.

Use the outlined specimen in `specimens/` to evaluate schedule and timestamp readability without relying on local font installation.

## Visual language

### Scheduling structure

Visual references come from cron fields, terminal prompts, timelines, task queues, service status, and run history. Use aligned rails and points to explain sequence or policy. Decorative clocks and calendar illustrations do not add meaning.

### Iconography

Use line icons on a 24-unit grid with strokes between 1.5 and 2 units. Prefer direct symbols: calendar, clock, terminal, repeat, play, history, task, service, pause, and warning. Six starter icons ship in `icons/` and inherit `currentColor`.

### Surfaces and geometry

Panels use 1 px Line borders and radii from 4 to 8 px. Keep shadows minimal. Selected states may use a mint border or a low-opacity mint fill. Amber and red backgrounds should remain small and paired with text.

### Motion

Use motion to show state transitions, refreshed next-run data, or task activity. Keep transitions between 120 and 240 ms on `cubic-bezier(.2,.6,.2,1)`. Respect `prefers-reduced-motion` everywhere.

## Voice and writing

Write like a maintainer explaining a runbook. Use active voice, exact nouns, and observable outcomes. State prerequisites before commands. Name the timezone and lifecycle boundary whenever they change the result.

### Preferred language

| Prefer | Avoid |
| --- | --- |
| task, schedule, expression, run | automation magic, workflow wizard |
| daemon, service, startup, reload | always-on intelligence |
| timezone, catch-up, overlap, policy | set it and forget it |
| next run unavailable | never miss a beat |
| supported, experimental, deferred | flawless, effortless, revolutionary |

### Examples

**Product statement:** Cross-platform task scheduling with a CLI, desktop application, and background daemon.

**Schedule explanation:** At 09:30 every weekday in America/New_York.

**Lifecycle explanation:** Runs once when the scheduler daemon starts. Reloading the task list does not run it again.

**Empty state:** No tasks match the current filters.

**Error:** The task could not start because its working directory is unavailable.

### Casing and terminology

The product name is always lowercase: **go-schedule**. Command names and file formats preserve their canonical casing. Use `daemon start` for the process lifecycle boundary and `host boot` only when discussing operating-system startup behavior.

## Parent brand relationship

go-schedule is an independent ShruggieTech product identity. It shares Space Grotesk, Geist, Geist Mono, dark-first discipline, and direct technical writing. It retains its own mint-and-blue scheduling palette and terminal-field mark.

The approved endorsement is **A ShruggieTech project**. Set it in Geist Mono with uppercase letters and positive tracking. Keep it subordinate and outside the logo clear space. Suitable placements include footers, About views, title-page colophons, repository metadata, and social previews.

## Digital implementation

Load the complete foundation and optional component layer:

```html
<link rel="stylesheet" href="styles.css">
<link rel="stylesheet" href="components/components.css">
```

The default surface is dark. Add `class="gs-light"` to a container for the accessible light token mapping.

| File | Provides |
| --- | --- |
| `tokens/colors.css` | Raw and semantic colors plus `.gs-light` |
| `tokens/typography.css` | Families, weights, tracking, leading, and scale |
| `tokens/spacing.css` | Spacing, radii, stroke, and motion |
| `tokens/base.css` | Resets, typography defaults, focus, technical text, reduced motion |
| `components/components.css` | Cards, buttons, fields, badges, and schedule rows |

### Favicons and platform assets

Favicon entries at 32 px and below use the reduced mark. Larger browser, Android, Apple, Windows, macOS, and Linux assets use the full mark on Night. The background protects the terminal structure in both light and dark system chrome.

Windows receives a genuine seven-entry ICO. macOS receives an ICNS bundle. Linux receives a hicolor tree and a desktop-entry template. The browser set includes SVG, ICO, PNG, Apple, Android, and manifest files.

## Asset inventory

| Directory | Contents |
| --- | --- |
| `logos/svg/` | Portable vector masters, approved lockups, marks, wordmarks, header, and social preview |
| `logos/png/` | High-resolution raster exports |
| `favicons/` | Browser, Apple, Android, SVG, ICO, and web manifest assets |
| `platform/` | Windows ICO, macOS ICNS, Linux hicolor icons, and desktop-entry template |
| `fonts/` | WOFF2, TTF, CSS declarations, and OFL licenses |
| `tokens/` | CSS and JSON design tokens with measured contrast claims |
| `icons/` | Six starter scheduling icons |
| `components/` | Framework-neutral CSS and thin React wrappers |
| `guidelines/index.html` | Live visual reference for the system |
| `ui_kits/go-schedule-web/` | Demonstration interface using shipped tokens and components |
| `specimens/` | Fully outlined type and schedule-readability specimen |
| `styles.css` | Single CSS foundation entry point |
| `SKILL.md` | Compact instructions for agents and contributors |
| `build/` | Reproducible generators and verification |
| `brand-guide.pdf` | Eleven-page printable reference manual |
| `VERIFY.md` | Measured dimensions, contrast, PDF checks, and SHA-256 inventory |

Use SVG in product interfaces and documentation whenever supported. Use PNG for social platforms, raster-only systems, and external listings. Treat `build/geometry.py`, the bundled fonts, and the generation scripts as the reproducible source of truth.
