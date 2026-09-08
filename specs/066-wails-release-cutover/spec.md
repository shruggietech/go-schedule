# Feature Specification: Wails Release Cutover

**Feature Branch**: `codex/066-wails-cutover`

**Created**: 2026-09-07

**Status**: Implemented

**Delivery**: Production packaging, retirement, compiled Windows MSI inspection, and canonical eight-gate verification completed 2026-09-07 on review branch `codex/066-wails-cutover`; see `verification.md`

**Input**: User description: "Complete the v1.2 desktop release boundary by packaging and qualifying Wails as the sole maintained desktop application, retiring Fyne only after parity and release gates pass, and preserving installer, daemon, CLI, upgrade, accessibility, and platform behavior. Tracks #157 and completes parent #147 when all criteria pass."

## User Scenarios & Testing

### User Story 1 - Install the complete Wails desktop (Priority: P1)

As a desktop user, I can install or unpack a supported go-schedule desktop package and receive the Wails control center together with the daemon and CLI, with the expected native application identity and no retired Fyne application.

**Why this priority**: The replacement is not shipped until production packages deliver the new application on every supported desktop platform.

**Independent Test**: Build each supported desktop artifact from one candidate revision, inspect its contents and native identity, and prove that the packaged desktop starts through the established executable or application entry point.

**Acceptance Scenarios**:

1. **Given** a supported Windows, macOS, or Linux release runner, **When** the desktop package is built, **Then** it contains the Wails application, daemon, CLI, required attribution, and platform integration assets.
2. **Given** a packaged desktop artifact, **When** its contents are inspected, **Then** no Fyne application binary, dependency, resource, or alternate desktop entry point is present.
3. **Given** a Windows installation, **When** the Start Menu or optional desktop shortcut is used, **Then** the Wails application starts without a console window and uses the established go-schedule identity.

---

### User Story 2 - Upgrade without losing operation or data (Priority: P1)

As an existing user, I can upgrade from the shipping Fyne-era package to the Wails package without losing scheduled work, daemon operation, CLI access, preferences eligible for migration, or the preserve-versus-wipe uninstall choice.

**Why this priority**: A visually successful replacement that breaks upgrades or data ownership is not releasable.

**Independent Test**: Exercise the supported upgrade, repair, uninstall-preserve, and explicit-uninstall-wipe paths and verify service, shortcut, executable, preference, and application-data outcomes.

**Acceptance Scenarios**:

1. **Given** a supported earlier installation with tasks and desktop preferences, **When** the candidate is upgraded in place, **Then** the daemon and CLI remain usable, task data is preserved, and the Wails desktop can migrate supported appearance intent.
2. **Given** the candidate installation, **When** it is uninstalled with preservation selected, **Then** binaries and integration entries are removed while application data remains.
3. **Given** the candidate installation, **When** explicit wipe is selected, **Then** only documented application-owned data is removed and the durable cleanup result remains available.

---

### User Story 3 - Qualify one auditable release candidate (Priority: P1)

As a maintainer, I can determine from one candidate revision whether functional parity, accessibility, native behavior, packaging, and supported-platform checks passed, with failures blocking cutover rather than being hidden by historical evidence.

**Why this priority**: The removal of the old desktop must be justified by current, reviewable evidence from the exact candidate.

**Independent Test**: Run the canonical verifier and hosted supported-platform matrix against one revision, then audit the parity and intentional-deviation records with no unexplained omission.

**Acceptance Scenarios**:

1. **Given** the candidate revision, **When** canonical and hosted checks complete, **Then** frontend, Go bridge, daemon integration, accessibility, responsive behavior, native builds, package contracts, and release automation are green.
2. **Given** the completed feature slices #149 through #156, **When** the final parity inventory is audited, **Then** every supported Fyne workflow maps to Wails behavior or a documented intentional retirement.
3. **Given** a failed required candidate gate, **When** release readiness is evaluated, **Then** the cutover remains incomplete until the same candidate or a newer reviewed candidate passes.

---

### User Story 4 - Follow current desktop guidance (Priority: P2)

As a user or contributor, I can rely on documentation, help surfaces, issue templates, and build guidance that describe the Wails desktop and its actual packaging without presenting Fyne as maintained or shipped.

**Why this priority**: Stale guidance would make the technical cutover confusing and invite reports against deleted behavior.

**Independent Test**: Search maintained documentation and repository automation for Fyne-only claims, validate links and examples, and confirm any remaining Fyne references are explicitly historical or migration-related.

**Acceptance Scenarios**:

1. **Given** maintained user and contributor documentation, **When** desktop instructions are followed, **Then** they name the current executable, supported packages, controls, and build path accurately.
2. **Given** an intentional historical or migration reference to Fyne, **When** it is read in context, **Then** it cannot be mistaken for a current implementation or support promise.

### Edge Cases

- A desktop build succeeds but its daemon or CLI companion is missing from the final archive or installer.
- A Wails build output name differs from the stable launcher and installer contract.
- A prior Windows GUI process is running during upgrade and holds the executable open.
- A previous installation has valid, malformed, missing, or unreadable Fyne preferences.
- A user requests preserved data, explicit wipe, repair, downgrade, or same-version reinstall.
- A supported-platform build passes while another matrix member fails or produces an incomplete artifact.
- Historical release notes, changelog entries, and migration documentation legitimately mention Fyne.
- The old proof application or old GUI code remains reachable from CI, release automation, dependencies, or contributor instructions.

## Requirements

### Functional Requirements

- **FR-001**: Production desktop release automation MUST build the Wails application natively on Windows amd64, macOS arm64, and Linux amd64 from the same candidate revision.
- **FR-002**: Each desktop distribution MUST include the daemon, CLI, license, README, changelog, and platform integration assets required by its documented format.
- **FR-003**: The stable desktop launch contract MUST remain `gosched-gui` for executable discovery, shortcuts, installer upgrades, scripts, and existing user habits even though the implementation changes.
- **FR-004**: Windows packaging MUST install the Wails desktop as a windowless executable and preserve service, PATH, shortcut, upgrade, repair, uninstall-preserve, and explicit-wipe behavior.
- **FR-005**: macOS packaging MUST provide a native application bundle with the established bundle identifier, icon, version, daemon, and CLI placement.
- **FR-006**: Linux packaging MUST provide the desktop executable and established freedesktop application and icon assets beside the portable bundle.
- **FR-007**: Release and continuous-integration automation MUST validate the production Wails frontend, Go boundary, native builds, accessibility contract, responsive contract, and package contents.
- **FR-008**: The exact candidate MUST pass the repository's canonical format, vet, lint, race, desktop, coverage, documentation, automation, and lifecycle gates.
- **FR-009**: A parity inventory MUST map every workflow named by #147 and every materially supported legacy desktop capability to the replacement, an intentional retirement, or a separately tracked blocker.
- **FR-010**: Intentional behavior and appearance deviations MUST state the user impact, rationale, and verification evidence.
- **FR-011**: The retired Fyne application source, entry point, test suite, dependencies, build gates, release steps, and maintained-only documentation MUST be removed after the parity inventory and replacement gates are established.
- **FR-012**: Remaining Fyne references MUST be limited to immutable history, release history, migration behavior, or explicit transition rationale and MUST be understandable as non-current.
- **FR-013**: The S060 proof application MUST no longer be built or tested as a parallel maintained desktop after the production application replaces it; durable design evidence MUST remain in the S060 specification.
- **FR-014**: Existing daemon and CLI server-only archives MUST remain cgo-free and continue covering their current Linux and macOS architectures.
- **FR-015**: Desktop package construction MUST fail when an expected binary, bundle, platform asset, or attribution file is absent.
- **FR-016**: Candidate qualification MUST reject stale evidence from a different revision and MUST report every required supported-platform result.
- **FR-017**: User documentation, contributor guidance, issue forms, release automation comments, and package documentation MUST describe the Wails desktop as current.
- **FR-018**: Existing Wails appearance migration MUST remain one-time, bounded to supported values, and non-destructive to daemon-owned data.
- **FR-019**: No public release, version tag, or final milestone closure MUST occur as part of this pull request; those actions remain part of the maintainer's post-merge release ritual.
- **FR-020**: The completed slice MUST preserve issue-level traceability for #157 and parent #147, and MUST report either fully satisfied acceptance criteria or explicit remaining blockers without closing them implicitly.

### Key Entities

- **Desktop candidate**: One immutable revision and version context from which all supported desktop artifacts and qualification evidence are derived.
- **Desktop distribution**: A platform-specific archive, application bundle, or installer containing the Wails desktop and required companions.
- **Parity entry**: A legacy capability mapped to its Wails destination, intentional retirement, or unresolved blocker with evidence.
- **Qualification result**: A required gate tied to a candidate revision, supported platform, outcome, and evidence location.
- **Intentional deviation**: A reviewed difference from the former desktop with user impact, rationale, and verification.

## Success Criteria

### Measurable Outcomes

- **SC-001**: All three supported desktop package targets build from one candidate revision and contain 100 percent of their required payload inventory.
- **SC-002**: Every workflow and capability in the final parity inventory has exactly one explained disposition, with zero unexplained omissions.
- **SC-003**: One hundred percent of required canonical and hosted candidate gates pass, with no skipped gate represented as successful.
- **SC-004**: Windows upgrade, repair, uninstall-preserve, and explicit-wipe contracts retain their documented outcomes while launching the replacement desktop.
- **SC-005**: Maintained source, dependencies, CI, release automation, and current documentation contain zero Fyne implementation dependencies or current-product claims.
- **SC-006**: The desktop remains keyboard-operable and produces zero serious or critical automated accessibility findings at the supported viewport and 80, 100, 150, and 200 percent zoom settings.
- **SC-007**: Existing daemon-owned tasks and run history survive the desktop replacement in every supported upgrade scenario.
- **SC-008**: A contributor can identify and run the production desktop verification path from current documentation without consulting the retired implementation.

## Assumptions

- Issue #157 is the only remaining implementation child of #147, and #147 is completed only after #157 merges and its acceptance criteria are satisfied.
- The stable `gosched-gui` filename remains a compatibility surface; changing the implementation does not require changing that external name.
- Windows amd64, macOS arm64, and Linux amd64 remain the supported desktop package targets, while server-only daemon and CLI archives retain their broader existing matrix.
- The current protected local IPC, daemon storage, and CLI contracts remain authoritative and do not change in this slice.
- Hosted native builds and package-contract jobs provide the cross-platform candidate evidence available during pull-request review; a public release and any separately attended signing or distribution ritual happen only after merge authorization.
- Historical specifications, changelog entries, and release notes may retain accurate Fyne references, while maintained operational guidance must clearly describe Wails.
