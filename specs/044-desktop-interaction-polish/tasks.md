# Tasks: Desktop Interaction and Appearance Polish

**Input**: Design documents from `/specs/044-desktop-interaction-polish/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Regression tests are mandatory and precede each behavioral correction.

**Organization**: Tasks are grouped by independently testable user story. The six linked issues share appearance preferences, semantic theme roles, Options, navigation, and native Windows evidence, so they form one coherent slice.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel because it changes a different file after shared expectations are fixed
- **[Story]**: Maps a task to one user story from spec.md

## Phase 1: Setup and Baseline

**Purpose**: Fix traceability, authoritative dependency behavior, and the S043 baseline before product changes.

- [x] T001 Verify synchronized `main`, create `codex/044-desktop-interaction-polish`, and record issues #101, #104, #105, #106, #109, and #111 in `specs/044-desktop-interaction-polish/spec.md`
- [x] T002 Record Fyne 2.8.1 button blending, Select, ScrollEvent, scroll-driver, theme, and custom-widget behavior in `specs/044-desktop-interaction-polish/research.md`
- [x] T003 Record Inter 4.1 and Ubuntu Sans upstream source, license, and unmodified asset provenance in `specs/044-desktop-interaction-polish/research.md` and `gui/assets/fonts/`
- [x] T004 Record the S043 System-font A/B observation and the remaining exact-candidate native evidence boundary in `specs/044-desktop-interaction-polish/spec.md`, `plan.md`, and `contracts/desktop-interaction.md`

---

## Phase 2: Foundational Preference and State Contracts

**Purpose**: Establish safe defaults, selector alternatives, contrast composition, and scroll normalization before changing views.

**CRITICAL**: Tests must fail against the S043 implementation before production behavior is added.

- [x] T005 [P] Add failing System-default, legacy-choice, new-font, scroll normalization, persistence, and reset tests in `gui/appearance_test.go`
- [x] T006 [P] Add failing semantic overlay, importance/state composite, contrast, and new-font resource tests in `gui/theme_test.go`
- [x] T007 [P] Add failing current-choice exclusion and ordered-alternative tests in `gui/options_test.go`
- [x] T008 [P] Add failing discrete-versus-precision delta and sensitive vertical-scroll delegation tests in `gui/scroll_test.go`
- [x] T009 Run the foundational focused tests and record expected failures in `specs/044-desktop-interaction-polish/verification.md`
- [x] T010 Implement System-default font migration, curated font identifiers, scroll sensitivity, validation, persistence, and reset in `gui/appearance.go`
- [x] T011 Implement shared selector-current/alternative synchronization in `gui/options.go`
- [x] T012 Implement the bounded sensitive vertical-scroll wrapper in `gui/scroll.go`
- [x] T013 Re-run foundational focused and race tests and confirm the shared state layer is green

**Checkpoint**: Preference and event calculations are deterministic, bounded, and independently testable.

---

## Phase 3: User Story 1 - Read Every Interactive State (Priority: P1) MVP

**Goal**: Keep labels, glyphs, selected identity, focus, and disabled meaning readable across every shared control state and palette.

**Independent Test**: Exercise the semantic state matrix and representative active navigation and primary dialog buttons in dark, light, and follow-system modes.

### Tests for User Story 1

- [x] T014 [P] [US1] Add failing active-plus-hover, focus, pressed, disabled, danger, success, warning, and light/dark composite contrast tests in `gui/theme_test.go`
- [x] T015 [P] [US1] Add failing representative navigation and primary-button state persistence tests in `gui/widgets_test.go` and `gui/navigation_test.go`

### Implementation for User Story 1

- [x] T016 [US1] Replace opaque interaction-role colors with readable translucent semantic overlays in `gui/theme.go`
- [x] T017 [US1] Preserve selected and focused non-color identity through existing button importance and destination-content state in `gui/widgets.go` and `gui/navigation.go`
- [x] T018 [US1] Run the complete state matrix and representative control tests in dark, light, and follow-system modes

**Checkpoint**: The S043 active-tab and Save/Okay hover regressions are prevented by one shared semantic contract.

---

## Phase 4: User Story 2 - Start With Crisp, Comfortable Typography (Priority: P1)

**Goal**: Default to System while offering and persisting System, Geist (brand), Inter, Ubuntu, and Monospace consistently.

**Independent Test**: Start from clean, invalid, and every explicit supported preference; switch every family and inspect current and subsequently constructed Info/dialog content.

### Tests for User Story 2

- [x] T019 [P] [US2] Add failing font catalog ordering, label mapping, regular/bold/monospace/symbol resource, and default/reset tests in `gui/appearance_test.go` and `gui/theme_test.go`
- [x] T020 [P] [US2] Extend Info invariants and appearance-combination tests for all curated fonts in `gui/info_test.go` and `gui/options_test.go`

### Implementation for User Story 2

- [x] T021 [US2] Add unmodified licensed Inter regular/bold and Ubuntu Sans regular/bold resources in `gui/assets/fonts/`
- [x] T022 [US2] Embed and map the curated regular/bold/monospace/symbol faces in `gui/theme.go`
- [x] T023 [US2] Present System, Geist (brand), Inter, Ubuntu, and Monospace in canonical order and restore Dark/System defaults in `gui/options.go`
- [x] T024 [US2] Run all font/default/Info tests and confirm every explicit valid preference survives restart semantics

**Checkpoint**: New profiles use the sharp native-tested choice while explicit older preferences remain intact.

---

## Phase 5: User Story 3 - Scan and Copy Storage Information Efficiently (Priority: P2)

**Goal**: Replace storage cards with compact aligned rows, keep exact data visible, and mute unavailable rows completely.

**Independent Test**: Render mixed available/unavailable and long-path entries at default and small widths, then inspect columns, wrapping, selection, Copy, muted treatment, and absence of horizontal scrolling.

### Tests for User Story 3

- [x] T025 [P] [US3] Add failing storage header, aligned-row, flexible-path, fixed-action, wrapping, and no-horizontal-scroll tests in `gui/options_test.go`
- [x] T026 [P] [US3] Add failing whole-row muted, unselectable, disabled-Copy, and still-visible lifecycle-detail tests in `gui/options_test.go`

### Implementation for User Story 3

- [x] T027 [US3] Replace card-per-location construction with header and compact row layout in `gui/options.go`
- [x] T028 [US3] Implement semantic available/unavailable row text and surface treatment while preserving exact Copy and lifecycle wording in `gui/options.go`
- [x] T029 [US3] Run storage resolver and Options interaction/layout tests and confirm no ownership, traversal, or wipe semantics changed

**Checkpoint**: Storage is dense, readable, truthful, keyboard-accessible, and vertically responsive.

---

## Phase 6: User Story 4 - Navigate and Scroll Comfortably (Priority: P2)

**Goal**: Finish the rail boundary and Exit glyph, then apply immediate bounded discrete-wheel sensitivity to every owned vertical view.

**Independent Test**: Exercise rail geometry and Exit at 1280 by 800 and 800 by 600, then drive each owned vertical scroll with precision and discrete deltas at 1x, 2x, and 4x.

### Tests for User Story 4

- [x] T030 [P] [US4] Add failing full-height boundary, balanced inset, danger Exit glyph, placement, and selection-independence tests in `gui/navigation_test.go`
- [x] T031 [P] [US4] Add failing Options slider label/range/immediate/persist/reset and every-owned-scroll coverage tests in `gui/options_test.go` and `gui/scroll_test.go`
- [x] T032 [P] [US4] Retain repeated title-bar/navigation close exactly-once coverage in `gui/app_test.go`

### Implementation for User Story 4

- [x] T033 [US4] Add the rail/content separator and semantic danger Exit treatment without changing destination order or shutdown routing in `gui/navigation.go`
- [x] T034 [US4] Add the Scroll sensitivity slider and live multiplier label to Appearance in `gui/options.go`
- [x] T035 [US4] Replace direct application-owned vertical scroll construction in `gui/options.go`, `gui/info.go`, and `gui/editor.go` with the shared wrapper
- [x] T036 [US4] Run navigation, scroll, Options, editor, and shutdown tests at supported dimensions and sensitivity bounds

**Checkpoint**: Navigation is clearly bounded and long views respond to conventional wheel input without multiplying precision input.

---

## Phase 7: Polish and Cross-Cutting Verification

**Purpose**: Complete documentation, lifecycle state, repository-wide evidence, provenance, and integrity.

- [x] T037 [P] Document appearance defaults, font choices/licenses, storage rows, selector behavior, and scroll sensitivity in `docs/options.md`
- [x] T038 [P] Extend native Windows contrast, DPI, font, storage, wheel, touchpad, and navigation checks in `test/windows/README.md`
- [x] T039 [P] Add S044 behavior and architecture decisions to `CHANGELOG.md`
- [x] T040 [P] Add S044 chronologically to `specs/README.md`
- [x] T041 Advance `specs/044-desktop-interaction-polish/spec.md` through In Progress to Implemented and record red/green, focused, provenance, native-boundary, and gate evidence in `verification.md`
- [x] T042 Run `C:\Program Files\Git\bin\bash.exe scripts/verify.sh all` in the foreground and record all eight canonical gate results in `specs/044-desktop-interaction-polish/verification.md`
- [x] T043 Perform a systematic requirements, code, renderer-state, concurrency, history, license, safety, and evidence review; fix every high-confidence finding
- [x] T044 Validate strict UTF-8 without BOM, scan changed text for mojibake, run `git diff --check`, and mark every resolved task in `specs/044-desktop-interaction-polish/tasks.md`

---

## Dependencies and Execution Order

### Phase Dependencies

- **Setup (Phase 1)** establishes scope and upstream behavior.
- **Foundational contracts (Phase 2)** block all four user stories.
- **Readable states (Phase 3)** is the MVP and establishes shared contrast rules.
- **Typography (Phase 4)** depends on preference/font catalog foundations but is independent of storage layout.
- **Storage presentation (Phase 5)** depends on selector helpers and semantic disabled treatment.
- **Navigation and scrolling (Phase 6)** depends on the semantic state and scroll foundations.
- **Polish (Phase 7)** depends on every user story.

### User Story Dependencies

- **US1** establishes the semantic interaction-state treatment consumed by US3 and US4.
- **US2** shares appearance state but can be validated independently after Phase 2.
- **US3** preserves the existing storage data contract and consumes US1's disabled colors.
- **US4** consumes US1's danger/interaction roles and Phase 2's scroll preference/wrapper.

### Parallel Opportunities

- T005 through T008 touch independent test surfaces after expectations are fixed.
- T014 and T015 separate theme math from representative widgets.
- T019 and T020 separate font resource contracts from view invariants.
- T025 and T026 separate geometry from unavailable-state semantics.
- T030 through T032 separate navigation, scroll settings, and shutdown.
- T037 through T040 touch independent documentation surfaces.
- No multi-agent delegation is used because current session policy reserves it for explicit operator requests.

## Implementation Strategy

1. Pin the failures against S043 with preference, contrast, selector, and scroll tests.
2. Implement the smallest shared state layer and turn the foundation green.
3. Fix systemic theme blending before adding visual consumers.
4. Add licensed fonts and default migration, then compact Options storage.
5. Finish navigation and apply scroll sensitivity to every owned vertical view.
6. Run focused, full GUI, race, integrity, and canonical verification.
7. Publish the already-authorized review branch and PR, resolve automatic review and CI, then request and resolve no more than one second Codex review round.

## Notes

- The operator explicitly authorized S044 push, PR creation, and verified in-scope review-fix pushes in the kickoff.
- At most one manual `@Codex review` request may follow the automatic first review round.
- Exact-candidate native visual and physical-device evidence remains a release gate; headless tests do not impersonate it.
- Merge, tag, release, and material scope expansion remain unauthorized.
