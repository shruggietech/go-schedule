# Feature Specification: GUI Navigation and Information

**Feature Branch**: `codex/017-gui-info-navigation`

**Created**: 2026-08-28

**Status**: Implemented

**Delivery**: [PR #49](https://github.com/shruggietech/go-schedule/pull/49)

**Input**: User description: "Combine issues #29 and #32 into one focused GUI navigation-and-information slice: group related tabs and add an Info tab with the application identity, version, attribution, repository, and documentation."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Follow a coherent workflow through the sidebar (Priority: P1)

As a desktop user, I see related management views grouped together and the remaining views ordered by the normal flow from configuration to outcomes and application information.

**Why this priority**: Navigation order affects every GUI session. Grouping the views makes the application easier to scan without changing any task behavior.

**Independent Test**: Open the application and inspect the first four leading sidebar views. They appear as Tasks, Groups, Schedule, Activity, and retain their current content and label behavior. This independently delivers the navigation grouping before the separate Info view is added.

**Acceptance Scenarios**:

1. **Given** the desktop application opens, **When** the user reads the sidebar from top to bottom, **Then** the views appear as Tasks, Groups, Schedule, Activity, and Info.
2. **Given** Activity has unacknowledged alerts, **When** its count changes, **Then** its existing bounded badge remains attached to the Activity item and Info remains the final item.
3. **Given** an existing navigation item is selected, **When** the user switches among Tasks, Groups, Schedule, and Activity, **Then** its existing content and behavior are unchanged by the reorder.

---

### User Story 2 - Identify the installed application (Priority: P1)

As a desktop user, I can open one local information view to confirm the running version, recognize the product, identify its maintainer, and find its official project resources.

**Why this priority**: Version and support identity are basic diagnostic facts. Users should not need to guess an address or leave the application to determine what build they are running.

**Independent Test**: Open Info without a daemon or network dependency and observe the official mark, the exact running build version, ShruggieTech attribution, and labeled links to the company, repository, and documentation.

**Acceptance Scenarios**:

1. **Given** any application build, **When** the user opens Info, **Then** the official application mark and the exact build version are plainly visible.
2. **Given** Info is open, **When** the user reviews attribution, **Then** ShruggieTech is identified as the builder and maintainer and is linked to `https://shruggie.tech`.
3. **Given** Info is open, **When** the user reviews project resources, **Then** labeled links identify `https://github.com/shruggietech/go-schedule` as the source repository and `https://shruggietech.github.io/go-schedule/` as the documentation site.
4. **Given** the daemon is unavailable or the machine is offline, **When** Info opens, **Then** all local identity and version content still appears; opening an external link remains the operating system's responsibility.

### Edge Cases

- Development builds show their actual development version rather than substituting a release-looking value.
- A long future version string remains readable and does not obscure the links.
- The full application mark scales without distortion and is accompanied by text, so product identity never depends on the image alone.
- External resources are presented as descriptive text links rather than raw or ambiguous addresses.
- The Activity badge may change after the tab collection is built; its update must not change tab order or the Info item.
- The historical issue #29 wording says Logs. The current product term is Activity, established by the completed Activity clarity slice, and this feature must not revert it.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The leading navigation MUST contain exactly five items ordered as Tasks, Groups, Schedule, Activity, and Info.
- **FR-002**: Existing Tasks, Groups, Schedule, and Activity content and interaction behavior MUST remain unchanged apart from their positions.
- **FR-003**: The existing Activity count label MUST continue to update in place without moving Activity or Info.
- **FR-004**: Info MUST render entirely from resources already available in the desktop process and MUST NOT require daemon access or a network request.
- **FR-005**: Info MUST show the existing official full application mark without introducing or duplicating a brand asset.
- **FR-006**: Info MUST show the exact version identifier embedded in the running build, including development identifiers.
- **FR-007**: Info MUST identify ShruggieTech as the application builder and maintainer and provide a descriptive link to `https://shruggie.tech`.
- **FR-008**: Info MUST provide descriptive links to the official source repository at `https://github.com/shruggietech/go-schedule` and documentation at `https://shruggietech.github.io/go-schedule/`.
- **FR-009**: The information hierarchy MUST remain readable at the application's supported window sizes, preserve the mark's aspect ratio, and identify every destination with visible text.
- **FR-010**: Automated tests MUST protect the exact navigation order, version presentation, canonical link labels and destinations, brand-resource reuse, and Activity badge position.
- **FR-011**: This slice MUST NOT add backend endpoints, configuration, stored state, new dependencies, new branding assets, update checks, telemetry, or changes to external-link behavior supplied by the desktop toolkit and operating system.

### Key Entities

- **Navigation item**: A stable sidebar entry with a visible label, ordered position, and existing content surface.
- **Application identity**: The local official mark, exact build version, and maintainer attribution shown in Info.
- **Official resource link**: A descriptive label paired with one canonical HTTPS destination for the company, source repository, or documentation.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A user can locate each of the five views in one top-to-bottom scan, with Tasks adjacent to Groups and Info in the final position.
- **SC-002**: Info presents one official mark, one exact version identifier, one maintainer attribution, and all three canonical external destinations without contacting the daemon or network.
- **SC-003**: One automated navigation assertion detects any changed, missing, duplicated, or reordered item among the five required labels.
- **SC-004**: Automated checks cover 100 percent of the required identity fields and link destinations named in FR-005 through FR-008.
- **SC-005**: All eight repository verification gates pass with no regression to existing Activity badge behavior, core coverage, or GUI build support.
- **SC-006**: Merging the slice fully satisfies issues #29 and #32 without partially implementing either issue or expanding into unrelated GUI work.

## Assumptions

- The Activity label is authoritative even though issue #29 predates its rename and refers to Logs.
- The full mark already embedded for application-level identity is the correct Info image; creating a banner or alternate brand treatment belongs to the broader branding backlog.
- The existing build-version value is the source of truth for both development and release binaries.
- The desktop toolkit's standard hyperlink control provides the expected keyboard, focus, and operating-system browser behavior.
- The current application window sizing and leading-tab layout remain supported; this slice does not redesign responsive navigation.
