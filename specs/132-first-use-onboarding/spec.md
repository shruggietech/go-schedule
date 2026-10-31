# Feature Specification: S132 First-use Onboarding and Real-workload Polish

**Feature Branch**: `codex/132-first-use-onboarding`

**Created**: 2026-10-31

**Status**: Implemented

**Delivery**: Review branch codex/132-first-use-onboarding; self-contained snapshot/failure/recovery guidance and selected Run ID delivered for #417; package, component/browser and canonical evidence in [verification.md](verification.md). Full foreground eight-gate run passed in 1108.337 seconds, with zero drift across 2570 evaluated source paths and 255 retained inputs.

**Input**: Owner authorizes S132 first-use onboarding and real-workload polish through Spec-Kit/autopilot, automatic push and one official PR, in-scope CI/review corrections and at most two review rounds. Final merge remains owner-controlled.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Complete a useful first task (Priority: P1)

A new operator installs the public package, confirms the local daemon connection, creates a self-contained task that writes a useful local file snapshot, runs it immediately and sees the completed result without inventing an executable or waiting for tomorrow.

**Why this priority**: The current Linux first-task example requires an uncreated backup script. Installation should lead to an achievable first result.

**Independent Test**: Follow the published first-task instructions in an isolated workspace using actual v2.0.0 package binaries, retaining the resulting snapshot and task/run identifiers separately from source-build results.

**Acceptance Scenarios**:

1. **Given** a fresh supported installation, **When** the operator follows the linked first-task guide, **Then** every input file and command needed by the example is supplied and a local snapshot is produced.
2. **Given** the task has been submitted, **When** the operator inspects its run, **Then** the guide distinguishes accepted work from a terminal outcome and explains how to find captured output.
3. **Given** a path containing spaces, **When** the operator follows the platform-specific instructions, **Then** command arguments and the execution identity's required path access are explicit.

### User Story 2 - Understand and repair a failed run (Priority: P1)

The operator can cause an intentional missing-input failure in the owned example workspace, find the failed run's output, restore the input and run the same task successfully while the original failure remains visible.

**Why this priority**: A successful hello message alone does not teach an operator how to recover when useful work fails.

**Independent Test**: Use the same task for success, missing-input failure and recovery, inspecting distinct terminal records and output and preserving the failed record.

**Acceptance Scenarios**:

1. **Given** a completed snapshot task, **When** its input is deliberately removed and the task runs again, **Then** the operator can identify the new failed run and the missing-input explanation.
2. **Given** that failure, **When** the operator restores the input and runs again, **Then** the new run succeeds without erasing history or recreating the task.
3. **Given** a stopped or unreachable daemon, **When** the operator starts the guide, **Then** the guide points to platform connection/service recovery and does not advise disabling access controls or replacing existing data.

### User Story 3 - Move between guides and interfaces (Priority: P2)

The operator can reach the same first-task workflow from the repository, public installation guidance and desktop reference, using current product terminology and accurate release claims.

**Why this priority**: Contradictory platform and desktop guidance creates avoidable setup friction.

**Independent Test**: Follow those entry links in the exported public site and compare desktop field names/actions against the existing production interface without opening a local GUI.

**Acceptance Scenarios**:

1. **Given** a supported platform guide or README, **When** the operator seeks the first useful workflow, **Then** a direct link identifies the self-contained guide and platform-specific command form.
2. **Given** the desktop reference, **When** the operator reads schedule and notification instructions, **Then** shipped cron and SMTP are described as available, and task/run inspection actions match the current interface.

### Edge Cases

- An accepted run may still be queued or running, and a task may have multiple historical runs; identify the selected run rather than assume the newest command completed synchronously.
- A system daemon may run under SYSTEM/root rather than the operator; workspace location, permissions and explicit shell behavior must be stated.
- Existing tasks/data, input/output files outside the named example workspace and the installed service must not be silently replaced or deleted.
- Unsupported native desktop environments remain disclosed; headless browser evidence is separate from attended native observations.
- Published package behavior and changed source behavior retain distinct version/commit attribution.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Supply a self-contained local snapshot workflow with explicit setup, daemon readiness, task creation, immediate execution, terminal status and captured-output inspection.
- **FR-002**: Supply accurate Windows and POSIX command forms, including executable, ordered arguments, working paths, execution identity and timezone/schedule choices.
- **FR-003**: Explain accepted/queued/running versus completed outcomes and identification of the particular run before interpreting its output. Desktop Run details must expose its own stable Run ID separately from any source-run identity so the operator can correlate it with the exact CLI record.
- **FR-004**: Supply a missing-input failure and same-task repair path that preserves earlier task/run history and provides an actionable explanation.
- **FR-005**: Link the workflow from README and relevant platform/desktop documentation and remove demonstrated contradictory current cron/SMTP status claims.
- **FR-006**: Preserve existing task authorization, IPC boundaries, current state/refusal behavior, and explicit ownership of any example-file cleanup; do not install or change services on the active workstation.
- **FR-007**: Correct any demonstrated first-use defect within this local workflow with issue-level traceability; when no application defect is demonstrated, preserve existing application behavior and report the documentation outcome accurately.
- **FR-008**: Preserve release assets, tagged notes/comparison records and historical diagnostic evidence; keep migrations, release publication, cluster work and unchanged historical replay campaigns outside S132.

### Key Entities

- **Example workspace**: One explicitly named local directory containing example input, snapshot output and a platform script, owned by the operator and accessible to the daemon identity.
- **Snapshot task**: One task with a stable ID, explicit execution command, local paths and a readable schedule; immediate runs reuse that ID.
- **Run**: A stable record of accepted work and its terminal status/output; successful recovery does not overwrite the prior failure.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A supported-platform operator has one complete route from installation to a local snapshot with zero undisclosed executable or input-file prerequisites.
- **SC-002**: One task supports three distinguishable outcomes in order: success, explained missing-input failure and success after repair.
- **SC-003**: README, all three platform installation guides and the desktop reference provide an unbroken route to the workflow and consistent shipped-feature status.
- **SC-004**: No existing daemon data, release asset or unrelated operator file is changed by following the owned-workspace setup and cleanup instructions.

## Assumptions

- Product scope is local first use with the actual published v2.0.0 packages; the existing install and permission models remain authoritative.
- The owner confirms no users; no migrations or old-version upgrade qualification are required.
- The Windows MSI may be inspected/extracted read-only and its binaries run in an owned isolated runtime. Attended installation/UI checks remain untested under the standing waiver; no local GUI focus or system-service mutation is authorized by this exercise.
- Windows and Linux package workflows can be exercised headlessly. macOS command shape may be assessed through existing CI and portable shell checks without claiming a local macOS native run.
- Routine decisions are resolved under autopilot. Automatic PR-open review consumes round one; at most one manual second request is permitted, and the PR remains unmerged.

## Clarifications

### Session 2026-10-31

- Q: Which workload should make the first-use path concrete? A: A local text-file snapshot using native shell/file tools, with no network service, interpreter dependency or credentials.
- Q: Should S132 add an onboarding subsystem or duplicate verification framework? A: Reuse existing product, docs export, package probes and browser/runtime checks; application changes require a demonstrated defect.
- Q: Does automatic publication require another permission halt? A: The owner explicitly authorizes push and official PR creation for S132; no redundant halt applies. Tag/release and merge remain outside that authorization.
