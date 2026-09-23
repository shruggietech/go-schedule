# Feature Specification: Target-Aware Portable Bundle Compatibility

**Feature Branch**: `codex/100-bundle-compatibility`

**Created**: 2026-09-22

**Status**: Implemented

**Delivery**: S100 review branch implementation with all eight local verification gates passed on 2026-09-22; pull request pending.

**Input**: S100, finishing the functional acceptance criteria of [#184](https://github.com/shruggietech/go-schedule/issues/184) after S098 and S099.

## User Scenarios & Testing

### User Story 1 - Know whether the selected daemon supports the bundle (Priority: P1)

An operator validates portable definitions against a selected local or remote daemon and sees specific missing capabilities or unsupported platform requirements before attempting to apply them.

**Independent Test**: Validate a source-bearing bundle against a daemon advertising only tasks, then confirm that validation names every missing feature and does not call a bundle mutation endpoint.

**Acceptance Scenarios**:

1. **Given** a selected daemon without bundle support, **when** the operator validates or previews a bundle, **then** the result explains that bundle operations are unavailable and no apply plan is issued.
2. **Given** a bundle containing watchers and notifications and a target missing those capabilities, **when** the operator validates it, **then** each missing capability is named separately.
3. **Given** a watcher bundle and an unsupported or unknown target platform, **when** the operator validates it, **then** the watcher platform requirement is reported before apply.

---

### User Story 2 - Review and apply to exactly the checked target (Priority: P1)

An operator previews or applies a bundle while changing the selected daemon elsewhere in the control center. The in-flight operation remains attached to the daemon whose compatibility was checked.

**Independent Test**: Switch the selected daemon between manifest discovery and preview and confirm that both requests use the same target, with no preview on the newly selected daemon.

**Acceptance Scenarios**:

1. **Given** target A is selected, **when** a bundle operation begins and the UI selects target B before the request finishes, **then** that operation remains bound to A.
2. **Given** a target that cannot support the bundle, **when** preview or apply is attempted, **then** no bundle mutation request is sent to that target.

---

### User Story 3 - Keep cross-platform drift review deterministic (Priority: P2)

An operator moves v1 or v2 portable intent between supported platforms, provides watcher paths on the target, and sees stable, read-only drift rather than changes caused by the origin machine's path syntax.

**Independent Test**: Exercise supported and unsupported target-platform profiles, v1 and v2 documents, and read-only comparison; confirm stable issue ordering and no inferred removals.

**Acceptance Scenarios**:

1. **Given** a v1 or v2 bundle and a compatible Windows, Linux, or macOS target, **when** it is validated, **then** the same portable intent is accepted independently of the source machine's path rules.
2. **Given** target-only items, **when** a bundle is compared, **then** those items are reported as drift and remain untouched.
3. **Given** a target-local watcher path, **when** a plan is previewed, **then** the target daemon, not the caller's platform, decides whether that path is absolute and usable.

### Edge Cases

- A manifest missing a capability list or platform fails closed for operations requiring those facts.
- A target may advertise bundle support but lack one source family; report that family rather than a generic transport failure.
- Unknown future capability names do not imply support for known bundle families.
- A failed manifest request does not fall through to bundle apply.
- Untrusted bundle content determines requirements; it cannot claim fewer requirements to bypass preflight.
- Existing v1 and v2 bundle bytes and digest rules do not change merely to add target compatibility checks.

## Requirements

### Functional Requirements

- **FR-001**: The selected daemon MUST advertise bundle-operation support distinctly from task, schedule, group, chain, trigger, watcher, and notification capabilities.
- **FR-002**: Validation MUST derive required capabilities from actual bundle content and return one actionable finding per missing capability before a plan is created.
- **FR-003**: A watcher bundle MUST identify unknown or unsupported target operating systems before apply; supported targets remain Windows, Linux, and macOS.
- **FR-004**: Export, validation, comparison, preview, and apply MUST use one immutable selected target for compatibility discovery and the corresponding operation.
- **FR-005**: Incompatible preview or apply MUST fail before sending that operation, and validation MUST expose findings in ordinary machine-readable results.
- **FR-006**: Target-specific watcher paths MUST remain outside the portable document and be interpreted by the selected daemon's operating system.
- **FR-007**: Compatibility preflight MUST preserve v1 and v2 document meaning, canonical digests, read-only drift, target-only preservation, and explicit no-removal behavior.
- **FR-008**: Missing or stale target discovery MUST fail closed without treating an unknown platform or capability as supported.

### Key Entities

- **Target manifest**: Daemon-owned identity, platform, and advertised capabilities for the selected target.
- **Derived requirements**: Capability and platform needs inferred from bundle object families, not provided by the bundle author.
- **Compatibility finding**: A stable, item-specific explanation of a missing target requirement.
- **Target-bound operation**: A bundle request and its preceding compatibility discovery performed against one immutable daemon selection.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Every missing capability required by a bundle appears exactly once in validation findings, in a deterministic order.
- **SC-002**: No preview or apply request reaches a target whose discovered capabilities cannot support it.
- **SC-003**: Switching selected daemons during a bundle operation cannot cause its discovery and operation requests to reach different daemons.
- **SC-004**: Existing v1 and v2 canonical documents retain their digest and round-trip behavior; comparison never mutates target-only state.

## Assumptions

- The current daemon's capability manifest is the authoritative preflight contract; a target without the new bundle capability is treated as an older unsupported daemon.
- Bundle requirements are derived from object families already present in schema v1 or v2, so no bundle schema revision is needed.
- Platform-specific watcher paths continue to be supplied at preview to the selected target; source paths are never exported.
- Service capability discovery is advisory to the client, while the target daemon remains authoritative for its own path and domain validation.
- S100 does not add deletion, background synchronization, new notification transports, or desktop tray behavior.
