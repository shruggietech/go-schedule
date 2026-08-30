# Implementation Plan: Repository Brand System Integration

**Branch**: `codex/032-brand-system-integration` | **Date**: 2026-08-30 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/032-brand-system-integration/spec.md`

## Summary

Import the approved BrandBuilder-standard kit into `brand/` as the canonical repository source, publish a concise `docs/brand.md` entry point, synchronize required docs/application/packaging copies through an explicit consumer map, and enforce inventory and copy integrity with a standard-library Go verifier included in the existing automation gate. Release packaging will consume the committed Windows, macOS, and Linux platform assets instead of rebuilding identity ad hoc.

## Technical Context

**Language/Version**: Go 1.25 for the repository integrity command and tests; Python 3.12, Node.js, SVG/CSS/HTML, and bundled font files remain optional brand-generation inputs

**Primary Dependencies**: Go standard library for routine verification; existing Fyne embedding, Jekyll 3.9/just-the-docs 0.4.2, WiX 6.0.2, GitHub Actions, and the kit's documented optional graphics toolchain for deliberate regeneration

**Storage**: Version-controlled files and SHA-256 inventory; no runtime persistence or schema change

**Testing**: Table-driven Go tests for inventory, SVG, encoding, consumer mappings, and negative drift; existing docs, GUI, packaging, release-contract, lifecycle, and eight CI-parity gates

**Target Platform**: Repository consumers on Linux, macOS, Windows, GitHub Pages, GitHub repository surfaces, and release archives

**Project Type**: Cross-platform Go desktop/CLI project with static documentation and packaging automation

**Performance Goals**: Verify the complete kit and consumer map in under 5 seconds on a normal checkout; no effect on scheduler runtime or dispatch latency

**Constraints**: Full canonical import at or below 3 MiB; no new application runtime dependency; routine CI must not require graphics/font-generation packages; assets must remain UTF-8 without BOM where textual; no live text in distributed SVGs; review branch and pre-publication halt required

**Scale/Scope**: Approximately 110 canonical kit files, 28 SVGs, 13 raster logo outputs, platform icon families, one public documentation page, and about 60 declared consumer copies

## Constitution Check

*GATE: Passed before research and re-checked after design.*

| Principle | Gate | Result |
| --- | --- | --- |
| I. Code Quality | The verifier has one responsibility, actionable errors, idiomatic Go, documented public contracts, and no ignored errors. | PASS |
| II. Testing Standards | Behavioral integrity work is test-first, including negative missing/hash/live-text/BOM/consumer-drift cases; full race and coverage gates remain mandatory. | PASS |
| III. User Experience Consistency | The canonical tokens and public page unify identity across docs, GUI, packaging, and downloads; accessibility rules are explicit. | PASS |
| IV. Performance Requirements | The verifier is bounded by the manifest and consumer map, uses streaming hashes, has a 5-second target, and does not touch scheduler hot paths. | PASS |
| V. Autonomous Execution | S032 follows specify through local verified commit on a review branch and halts before publication. | PASS |

Pinned artifacts in scope are `.github/workflows/release.yml` and `build/windows/goschedule.wxs`; the automation contract also changes through `scripts/automation-check.sh`. The changelog records the dated decision because identity generation and verification become repository architecture rather than ad hoc release behavior. No constitutional deviation is required.

## Project Structure

### Documentation (this feature)

```text
specs/032-brand-system-integration/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── verification.md
├── checklists/
│   ├── requirements.md
│   └── brand-system.md
├── contracts/
│   └── brand-integrity.md
└── tasks.md
```

### Source Code (repository root)

```text
brand/                              # complete canonical kit
├── README.md                       # identity and selection guide
├── REPOSITORY.md                   # canonical/copy/regeneration workflow
├── brand-guide.pdf
├── manifest.json                   # upstream kit checksums
├── repository-consumers.json       # exact canonical-to-consumer mappings
├── build/                          # reproducible optional generation pipeline
├── components/ guidelines/ icons/ logos/ specimens/
├── favicons/ fonts/ platform/ tokens/ ui_kits/
└── VERIFY.md

scripts/brand-check/
├── main.go                         # standard-library verifier command
├── manifest_test.go                # inventory and path contract tests
├── text_test.go                    # encoding and portable-SVG contract tests
└── consumers_test.go               # exact-copy contract tests

docs/
├── brand.md                        # public kit page
├── assets/brand/                   # mapped web-consumer copies
├── assets/favicons/                # mapped web-consumer copies
└── assets/fonts/                   # mapped web-consumer copies

gui/assets/                         # mapped embedded full/reduced marks and fonts
cmd/gosched-gui/icon.ico            # mapped Windows executable consumer
.github/workflows/release.yml       # canonical Windows/macOS/Linux packaging
build/windows/goschedule.wxs        # canonical installer identity
scripts/automation-check.sh         # brand integrity in existing automation gate
```

**Structure Decision**: `brand/` is the source of truth because it keeps the complete identity system discoverable without coupling it to a single consumer. Consumer-local copies remain only where Jekyll publication or Go embedding requires them, and `brand/repository-consumers.json` makes each duplication explicit and mechanically checked. The verifier lives under `scripts/brand-check/` because it is a repository-maintenance command, not application behavior.

## Planned Delivery Phases

1. Import the complete validated kit and add repository-specific canonical/copy guidance.
2. Write failing integrity tests, implement the standard-library verifier, and declare consumer mappings.
3. Publish the dedicated documentation page and synchronize all documentation assets.
4. Synchronize desktop and Windows identity, then simplify macOS and enrich Linux release packaging around canonical platform assets.
5. Wire brand integrity into automation, update maintainership guidance and changelog decisions, and run focused plus full verification.

## Post-Design Constitution Re-check

The design remains compliant. The only duplication is forced by existing publishing and embedding boundaries, is enumerated rather than implicit, and is protected by byte-level comparison. Optional generation dependencies stay isolated from routine CI. Importing the full 2.1 MiB kit is simpler and safer than inventing a partial second kit, while remaining below the specification's proportionality cap.

## Complexity Tracking

No constitution violations require justification.
