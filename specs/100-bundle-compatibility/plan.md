# Implementation Plan: Target-Aware Portable Bundle Compatibility

**Branch**: `codex/100-bundle-compatibility` | **Date**: 2026-09-22 | **Spec**: [spec.md](spec.md)

**Input**: S100 specification and issue #184, after S098 and S099.

## Summary

Finish the target-compatibility gap in portable bundles. Advertise bundle support in the daemon manifest, derive required target capabilities from each bundle's actual object families, and make the shared client check one frozen selected target before forwarding every bundle operation. Return ordinary validation findings for incompatibility, leave the daemon authoritative for watcher path syntax, and preserve the existing v1/v2 wire format and read-only drift behavior.

## Technical Context

**Language/Version**: Go 1.26 and the existing Wails/TypeScript desktop frontend.

**Primary Dependencies**: Existing Go standard library, daemon manifest endpoint, shared API client, bundle model.

**Storage**: No schema or data migration.

**Testing**: Go unit and client transport tests, bundle API integration tests, full `scripts/verify.sh all` parity.

**Target Platform**: Windows, Linux, and macOS for watcher-bound bundles; all existing daemon targets for platform-neutral portable objects when capabilities are advertised.

**Project Type**: Daemon API plus shared desktop and CLI client.

**Performance Goals**: One manifest request per bundle operation and deterministic linear capability checking over fixed feature families, including schedule and task support for task definitions.

**Constraints**: No bundle schema change, secret export, inferred deletion, background sync, or new external dependency.

**Scale/Scope**: One selected daemon and one bundle per operation; existing plan lifetime and item limits remain unchanged.

## Constitution Check

- **I. Code Quality**: Keep target compatibility in a small pure bundle helper and shared client adapter; format, vet, and lint remain mandatory.
- **II. Testing**: Add missing-capability, platform, target-switch, and no-forwarding regressions before code; run race and full CI parity.
- **III. UX Consistency**: CLI and desktop already share the client, so both receive identical actionable findings and errors.
- **IV. Performance**: This is an operator-initiated cold path; one manifest request and fixed-size capability set need no benchmark or cache.
- **V. Autopilot**: Spec-kit analyze blocks implementation; use the review branch and PR. The user's explicit push/PR authorization satisfies the publication halt for this slice.
- **Security**: Derive requirements from untrusted document contents rather than a claim in the bundle; freeze the selected client before discovery and operation; never include watcher paths in bundle digests.

## Decisions

1. **No bundle v3**: We considered adding `required_capabilities` to bundle JSON, deriving requirements from content, or requiring operators to supply a compatibility profile. Derived requirements are authoritative and preserve v1/v2 digests, so choose that approach.
2. **Advertise `bundles` in the manifest**: We considered probing the bundle endpoint or adding a separate negotiation endpoint. Existing daemon manifest capabilities are the established target contract, so add one capability to that list and treat its absence as a legacy target.
3. **Freeze selection per operation**: We considered relying on server plan identity alone or calling the switchable client twice. The latter can select different targets between discovery and request. Resolve the selected concrete client once and send both calls through it. Server target-bound plan checks remain the final mutation boundary.
4. **Platform scope**: Only watchers impose a platform-specific requirement because their paths are target-local. Supported Windows, Linux, and macOS targets are accepted; unknown OS values produce a watcher compatibility finding. Platform-neutral bundles are not rejected solely due to an unfamiliar OS when they advertise all required capabilities.
5. **Validation shape**: Return compatibility findings through the existing validation response without invoking a missing endpoint. Preview, comparison, export, and apply return a concise actionable compatibility error before forwarding. No new wire route is needed.

## Project Structure

```text
internal/bundle/                 Derived requirement and target-compatibility model and tests
internal/api/server/             Manifest capability advertisement and tests
internal/api/client/             Frozen-target preflight across bundle operations and transport tests
specs/100-bundle-compatibility/  Specification, decision record, contract, tasks, and delivery evidence
```

**Structure Decision**: Extend the existing portable bundle, manifest, and shared client modules. CLI and desktop use the shared client without separate compatibility logic.

## Complexity Tracking

No constitution violation or new dependency is required.
