# Requirements Quality Checklist: Windows Installed Core Recovery

**Purpose**: Validate that S038 requirements are complete, unambiguous, secure, and independently verifiable before planning **Created**: 2026-09-01 **Feature**: [spec.md](../spec.md)

## Scope and Traceability

- [x] CHK001 The specification names both authorizing GitHub issues and gives each issue an independently testable user story.
- [x] CHK002 The specification explicitly distinguishes S038 outcomes from the broader S039 release-smoke gate in #94.
- [x] CHK003 S036's superseded authorization exclusion and preserved recovery presentation are stated without rewriting historical delivery claims.
- [x] CHK004 Out-of-scope items prevent release publication, merge, general shell parsing, Windows `run_as`, and unrelated authorization expansion from entering the slice.

## Security and Authorization

- [x] CHK005 Authorized principals are bounded to SYSTEM, Built-in Administrators, the configured group, and verified direct user members.
- [x] CHK006 Requirements explicitly prohibit Authenticated Users, Everyone, UAC disablement, permanent elevation, and interactive-user service identity.
- [x] CHK007 Lookup/enumeration failure behavior is fail-closed and distinguishes an unusable member record from failure to inspect the group.
- [x] CHK008 SDDL validation, deterministic ordering, de-duplication, and identity-safe logging are measurable requirements.
- [x] CHK009 An unrelated-user denial scenario independently verifies that the access boundary remains restricted.

## Execution Correctness

- [x] CHK010 Manual and scheduled execution are separate acceptance scenarios using the production service-hosted executor.
- [x] CHK011 Success requires run outcome, exit code, output, and side-effect evidence rather than accepting a fake or recording runner.
- [x] CHK012 A nonzero child exit and a process-start failure have distinct diagnostic expectations.
- [x] CHK013 Service isolation, hidden child processes, noninteractive behavior, and output bounds remain explicit invariants.
- [x] CHK014 The task command, arguments, working directory, service account, trigger, and redacted environment evidence are enumerated.

## Verification Quality

- [x] CHK015 Automated descriptor tests cover omitted alias SID, empty group, duplicate/unstable members, malformed SIDs, and failure paths.
- [x] CHK016 Native verification requires real named-pipe and service-hosted execution evidence and forbids a fake-backend substitution.
- [x] CHK017 Unavailable disposable-host evidence must be reported honestly rather than inferred from source inspection.
- [x] CHK018 Success criteria are measurable and include all eight canonical repository gates.

## Notes

- All checklist items passed during the 2026-09-01 autopilot specification review.
- Implementation tests remain separate from this requirements-quality checklist.
