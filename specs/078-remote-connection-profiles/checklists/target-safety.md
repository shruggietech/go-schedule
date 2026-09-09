# Target Safety Requirements Checklist: Remote Connection Profiles and Target-Safe Clients

**Purpose**: Review the written requirements for wrong-target prevention, credential safety, lifecycle completeness, accessibility, and scripting stability
**Created**: 2026-09-09
**Feature**: [spec.md](../spec.md)

## Identity and Selection

- [x] CHK001 Are local labels, profile IDs, daemon IDs, endpoints, and credential IDs assigned distinct identity roles? [Clarity, Spec FR-002]
- [x] CHK002 Is the behavior for same-named daemons specified with objective disambiguation fields? [Coverage, Spec US1]
- [x] CHK003 Are desktop-persistent and CLI-invocation selection rules explicitly separated? [Consistency, Spec FR-004, FR-005]
- [x] CHK004 Are invalid, incomplete, and conflicting selections required to fail before daemon contact? [Exception Flow, Spec FR-006]

## Credential and Trust Boundaries

- [x] CHK005 Are profile-file secret exclusions exhaustive enough to prevent bearer or phrase persistence? [Completeness, Spec FR-003]
- [x] CHK006 Are transport trust, redirect, TLS, native credential, and identity-pinning requirements specified together? [Consistency, Spec FR-007]
- [x] CHK007 Are pair, repair, and removal partial-failure outcomes defined without plaintext fallback? [Recovery, Spec FR-008, FR-012, FR-013]
- [x] CHK008 Is protected CLI input distinguished from ordinary process arguments and history? [Security, Spec FR-009]

## User Experience and Accessibility

- [x] CHK009 Is target context required globally and again at every mutation decision point? [Coverage, Spec FR-014]
- [x] CHK010 Are unsupported capability and permission states required to explain disabled actions without local fallback? [Clarity, Spec FR-015]
- [x] CHK011 Are keyboard and assistive-technology requirements measurable for destructive confirmations? [Accessibility, Spec SC-007]
- [x] CHK012 Is stale generation rejection specified during target switches? [Concurrency, Spec FR-016]

## CLI and JSON Stability

- [x] CHK013 Is local-by-default CLI behavior measurable without weakening existing stdout contracts? [Compatibility, Spec FR-005, SC-003]
- [x] CHK014 Are remote human diagnostics and JSON stdout separated explicitly? [Consistency, Spec FR-017]
- [x] CHK015 Are direct JSON enrollment, authentication, identity, error, timeout, and mutation replay requirements complete? [Completeness, Spec FR-019]
- [x] CHK016 Are non-goals explicit enough to prevent resilience and remote-route scope expansion? [Scope, Spec FR-022]

## Dependencies and Completion

- [x] CHK017 Are completed foundations and downstream blocked issues identified? [Dependency, Spec Dependencies]
- [x] CHK018 Can every success criterion be verified through deterministic tests or documented clean-environment scenarios? [Measurability, Spec Success Criteria]
