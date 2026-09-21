# S098 Implementation Tasks: Portable Automation Bundles and Drift Reporting

**Feature**: S098, GitHub issue #184

## Completed

- [x] T001 Create the versioned `internal/bundle` document, canonical serialization, digesting, structural validation, drift planner, and deterministic ordering.
- [x] T002 Persist stable portable identities independently from daemon-local record IDs with migration v22.
- [x] T003 Add protected API export, validate, preview, compare, and apply routes plus Observe and Manage authorization catalog entries.
- [x] T004 Bind apply to a one-time, short-lived, target-fingerprinted preview and reject stale, altered, or replayed plans.
- [x] T005 Apply only groups, task scheduling intent, and completion chains. Create imported tasks as disabled command-free drafts and preserve execution inputs on updates.
- [x] T006 Report trigger keys, watcher paths, and task execution inputs as explicit exclusions without exposing their values.
- [x] T007 Add shared local and remote API client methods and the CLI export, validate, preview, compare, and apply workflow.
- [x] T008 Add the desktop bundle page with selected-target context, validation findings, plan review, explicit confirmation, and outcome display.
- [x] T009 Cover canonical output, omission safety, secret redaction, target drift, preview replay, migration, authorization route coverage, and desktop build behavior.

## Finalization

- [x] T010 Run repository formatting, focused and full test suites, desktop build, spec-kit analysis, and publication formatting checks.
- [ ] T011 Update GitHub issue #184 and project status, publish the S098 pull request, and complete the authorized review protocol.
