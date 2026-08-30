# Feature Specification: Repository Brand System Integration

**Feature Branch**: `codex/032-brand-system-integration`

**Created**: 2026-08-30

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Review branch `codex/032-brand-system-integration`; local verification completed 2026-08-30

**Input**: User description: "Finish the full go-schedule brand kit, bring as much of it into the repository as practical, spec the integration as one slice, and run it end-to-end under autopilot."

## Clarifications

### Session 2026-08-30

- Q: Which artifact is authoritative when current repository assets and the completed external kit differ? → A: The completed BrandBuilder-standard kit is authoritative; repository consumers receive explicit derived copies only where their tooling requires them.
- Q: How much of the completed kit belongs in version control? → A: Import the full usable kit and reproducible source, excluding transport archives, temporary environments, caches, and audit previews.
- Q: Should this slice absorb the environment-dependent Windows icon observation in #33? → A: No. It integrates and verifies canonical assets without requiring a VM or manual install-uninstall observation; #33 remains deferred.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Obtain the Complete Canonical Kit (Priority: P1)

As a maintainer or contributor, I can find one complete, self-explaining brand system in the repository, including portable masters, generated outputs, tokens, fonts and licenses, usage guidance, platform assets, and its reproducible build and verification material.

**Why this priority**: Without a canonical in-repository source, every product and documentation surface can drift and future contributors cannot reliably identify the approved asset.

**Independent Test**: Starting from a fresh checkout, inspect the brand-system entry point and inventory, then verify that every listed file exists, every recorded checksum matches, all SVG masters are portable, and the guide can be opened without relying on the original external export location.

**Acceptance Scenarios**:

1. **Given** a fresh repository checkout, **When** a contributor follows the brand-system entry point, **Then** the complete approved kit and its usage rules are available under one canonical repository directory.
2. **Given** a contributor needs a logo, color, font, icon, platform asset, or UI reference, **When** they consult the kit inventory, **Then** they can select an approved artifact without recreating or guessing it.
3. **Given** the kit is validated in a clean checkout, **When** its integrity checks run, **Then** missing files, checksum drift, live SVG text, encoding damage, or incorrect downstream copies produce a failing result.

---

### User Story 2 - Discover and Download the Brand (Priority: P1)

As a person viewing the documentation site, I can understand the go-schedule identity and download the appropriate approved assets from a dedicated brand page.

**Why this priority**: The current site uses brand assets internally but gives integrators, contributors, and press no canonical public page explaining or distributing them.

**Independent Test**: Validate the documentation source and links, then inspect the brand page at desktop and narrow widths to confirm that the mark, lockups, palette, typography, usage rules, attribution, and downloads are present and readable.

**Acceptance Scenarios**:

1. **Given** a visitor needs go-schedule artwork, **When** they open the documentation navigation, **Then** a dedicated Brand system page presents the identity and approved downloads.
2. **Given** a visitor chooses a light, dark, monochrome, compact, or social context, **When** they consult the page, **Then** it identifies the suitable variant and prohibited misuse without duplicating the complete long-form guide.
3. **Given** a referenced asset is renamed or omitted, **When** documentation integrity validation runs, **Then** the broken brand-page reference is reported.

---

### User Story 3 - Keep Product Surfaces Synchronized (Priority: P2)

As a maintainer, I can use canonical outputs across the README, documentation site, desktop application, Windows packaging, macOS packaging, favicons, and social previews, with automated detection when a required embedded copy drifts.

**Why this priority**: Some build systems require assets beside their consumers, but those copies must not become competing masters.

**Independent Test**: Compare each declared consumer copy with its canonical kit counterpart, run the desktop asset tests and packaging checks, and confirm that modifying any checked copy causes the brand-integrity gate to fail.

**Acceptance Scenarios**:

1. **Given** product tooling requires a local copy, **When** the approved kit is integrated, **Then** the copy is byte-identical to or reproducibly derived from a declared canonical artifact.
2. **Given** an application or documentation asset drifts, **When** the canonical verification suite runs, **Then** it fails with the exact source and consumer paths.
3. **Given** a normal release build, **When** Windows and macOS identities are assembled, **Then** they consume approved current brand artwork without adding a new runtime dependency.

### Edge Cases

- The full-size mark and reduced small-size mark have different intended consumers; verification must reject accidental substitution where the distinction matters.
- White, black, interval-color, light-surface, and dark-surface variants can have similar geometry but must remain semantically distinct and correctly named.
- SVG assets must remain usable on systems without the bundled fonts installed; live text and font-family dependencies are therefore invalid in distributed SVGs.
- Documentation links must work under the configured `/go-schedule` base path rather than only from a local filesystem root.
- The transport ZIP, pre-rebuild backup, rendering previews, dependency environments, and Python caches are not product assets and must not enter the repository.
- Existing historical Spec-Kit artifacts and changelog records can describe superseded artwork; integrity checks must target current consumer files rather than rewriting history.
- The repository must remain buildable when optional brand-generation dependencies are absent; ordinary CI verifies committed outputs without regenerating them.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The repository MUST contain one clearly named canonical brand-system directory with a single documented entry point.
- **FR-002**: The canonical directory MUST include the approved brand guide, identity and usage guidance, vector masters, raster outputs, tokens, fonts and licenses, icons, UI references, platform assets, build source, verification report, and checksum inventory.
- **FR-003**: Transport archives, historical backup archives, build environments, caches, and visual-audit scratch files MUST remain outside the repository.
- **FR-004**: Every distributed SVG MUST be UTF-8 without a byte-order mark, contain an accessible title, and avoid live text and font-family dependencies.
- **FR-005**: The repository MUST retain the distinct full mark, reduced small-size mark, horizontal lockups, stacked lockups, wordmarks, monochrome variants, header artwork, and social-preview artwork defined by the approved kit.
- **FR-006**: Brand colors, typography, spacing, radii, motion, semantic roles, and measured contrast claims MUST be available in human-readable guidance and machine-readable tokens.
- **FR-007**: Font files MUST be accompanied by their applicable license texts and documented roles.
- **FR-008**: The documentation site MUST contain a navigable Brand system page covering identity, logo selection, color, typography, usage constraints, parent attribution, accessibility, and direct downloads.
- **FR-009**: The brand page MUST use repository-hosted assets and remain readable at desktop and narrow layouts without introducing a new documentation build system.
- **FR-010**: Current README and documentation identity assets MUST use approved kit outputs or declared repository derivatives.
- **FR-011**: Desktop full-size and window-size embedded icons MUST use the approved full and reduced marks respectively.
- **FR-012**: Windows executable and installer identity MUST use the approved multi-resolution Windows icon.
- **FR-013**: macOS release packaging MUST consume the approved icon source or approved platform asset while preserving the existing release shape.
- **FR-014**: Favicons, web manifest icons, repository header, and social-preview imagery MUST be synchronized with the approved kit.
- **FR-015**: An offline integrity check MUST validate required inventory, canonical checksums, portable SVG constraints, UTF-8 integrity, and every declared consumer copy in a fresh checkout without optional brand-generation dependencies or additional setup.
- **FR-016**: The canonical automation gate MUST run the brand integrity check and fail with actionable file paths when a contract is violated.
- **FR-017**: Existing application, documentation, packaging, and release verification MUST remain green with no new runtime dependency.
- **FR-018**: Maintainer documentation MUST explain how to select assets, when to regenerate the kit, how to validate committed outputs, and which files are canonical versus consumer copies.
- **FR-019**: The delivered slice MUST fully satisfy and close GitHub issues #10 and #34 while leaving issue #33 open with no claim of environment-based verification.
- **FR-020**: The full brand-kit import MUST add no more than 3 MiB to the repository checkout so the comprehensive source remains proportional to this project's scale.

### Key Entities

- **Canonical brand artifact**: An approved master, token, guide, font, generated output, platform asset, or reference file stored in the repository brand-system directory.
- **Consumer copy**: A required copy or deterministic derivative located beside documentation, embedded application resources, or packaging inputs.
- **Consumer mapping**: A declared relationship between one canonical artifact and one consumer path, including whether byte identity or a documented derivation is required.
- **Brand integrity report**: Verification evidence covering inventory, hashes, SVG portability, encoding, consumer synchronization, and integration checks.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A fresh checkout contains 100% of the files listed in the canonical brand inventory, and every listed checksum validates.
- **SC-002**: All distributed SVG assets pass portability and accessible-title checks with zero live-text or font-family findings.
- **SC-003**: The dedicated documentation page provides working links to at least one approved asset for each of these needs: full mark, reduced mark, horizontal lockup, monochrome use, social preview, brand guide, and complete inventory.
- **SC-004**: Every declared README, documentation, desktop, favicon, Windows, macOS, and social consumer mapping passes the synchronization check; changing any one mapped copy makes the check fail.
- **SC-005**: The complete eight-gate local verification command passes with the brand integrity check included in the automation gate.
- **SC-006**: The committed brand-system directory is no larger than 3 MiB and requires zero new application runtime dependencies.
- **SC-007**: Issues #10 and #34 have complete closing evidence in the eventual pull request, while #33 receives no closure claim.

## Assumptions

- The completed kit at the operator-provided GPT-app output location is approved as version 1.0.0 and is the import source for this slice.
- Generated PNG, ICO, ICNS, PDF, and font binaries are legitimate source-distribution artifacts because consumers cannot reconstruct them with the repository's ordinary Go-only toolchain.
- Optional Python and browser dependencies are acceptable for deliberate brand regeneration but are not installed or run by routine CI.
- Existing repository asset paths may remain as consumer locations when changing them would add build complexity; their canonical relationship must instead be explicit and verified.
- GitHub repository social-preview configuration is hosted state and cannot be changed by a repository commit, so this slice supplies the approved image and documents its location rather than mutating that setting.
- Issue #33 remains a future environment-observation task as previously directed by the maintainer.
