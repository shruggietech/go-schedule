---
title: Unattended CLI evaluation
nav_order: 7.6
---

# Unattended CLI evaluation

This maintainer workflow evaluates the CLI and shared application behavior against disposable state on Windows 11 and an explicitly selected existing WSL2 distribution. S108 implements issues [#289](https://github.com/shruggietech/go-schedule/issues/289) and [#290](https://github.com/shruggietech/go-schedule/issues/290). It is a source-tree maintainer tool; its presence does not announce a new product release.

Start with the [version 1 execution prompt](https://github.com/shruggietech/go-schedule/blob/main/docs/maintainers/cli-qa-run-v1.md). Give the resulting report to a separate agent with the [version 1 issue-triage prompt](https://github.com/shruggietech/go-schedule/blob/main/docs/maintainers/cli-qa-triage-v1.md).

## Run the coordinator

Use an existing unelevated PowerShell 7 session or a verified hidden process launcher. Select an existing WSL2 distribution and its existing nonroot user before starting. The script launches every console child with CreateNoWindow, separate redirected streams and noninteractive input.

```powershell
pwsh -NoProfile -NonInteractive -File scripts/cli-qa.ps1 -WSLDistribution Ubuntu-24.04 -WSLUser EXISTING_USER
```

The script prints the exact input and artifact directories. It bootstraps the Go coordinator with two Go workers and build parallelism two, then prepares the Windows and Linux CLI, daemon and coordinator binaries under bounded Windows Jobs. Dependencies must already exist in the module cache. Preparation uses GOTOOLCHAIN=local and GOPROXY=off and does not provision tools, distributions, users or host configuration.

To evaluate already prepared inputs, supply BinaryDirectory. Capture creates a new OutputDirectory and refuses to overwrite an existing run. Keep the six input files unchanged throughout capture and review.

```powershell
pwsh -NoProfile -NonInteractive -File scripts/cli-qa.ps1 -Action capture -BinaryDirectory A:\qa\inputs -OutputDirectory A:\qa\run-001 -WSLDistribution Ubuntu-24.04 -WSLUser EXISTING_USER
pwsh -NoProfile -NonInteractive -File scripts/cli-qa.ps1 -Action finalize -BinaryDirectory A:\qa\inputs -OutputDirectory A:\qa\run-001 -Reviews A:\qa\reviews.json
pwsh -NoProfile -NonInteractive -File scripts/cli-qa.ps1 -Action validate -BinaryDirectory A:\qa\inputs -OutputDirectory A:\qa\run-001
pwsh -NoProfile -NonInteractive -File scripts/cli-qa.ps1 -Action render -BinaryDirectory A:\qa\inputs -OutputDirectory A:\qa\run-001
```

The underlying command is cmd/gosched-qa. Its build, inventory, capture, resume, finalize, validate and render commands expose the same workflow. Validate reads the frozen artifacts without rewriting them. Render derives report.md from validated report.json. Capture never supplies semantic review verdicts on an agent's behalf.

## Execution boundaries and coverage

JSON and human-readable command output run against separate fresh fixtures, so rendering checks do not reuse IDs or mutation state across modes. The manifest comes from the initialized Cobra command tree, including aliases, generated commands, inherited flags and local shadows. A checked-in inventory detects metadata drift. Help and invalid-flag cases complement behavioral cases; they do not stand in for implementing a command's actual scenario. Generated IDs derive from command paths, so adding an unrelated command preserves existing case references. Stable IDs and the complete environment cross-product remain in every report.

Three execution kinds identify the observation boundary: packaged_cli runs an actual executable; command_harness invokes actual command handlers with scoped adapters or exercises an explicitly recorded application boundary; deterministic_engine invokes shared product logic with controlled time and recording execution. The manifest and each result preserve that distinction.

An owned daemon provides unique IPC, data, configuration, logs, TLS and loopback destinations. Real benign task children exercise success, failure, launch diagnostics, quoting, Unicode, stdin, environment, working directory and capture limits. Other cases cover definitions, schedules, chains, triggers and sets, watchers, notifications, history, audit, access, profiles, preferences, bundles, summaries, search and MCP. Service and GUI handlers record their requested operations. Credential cases use a disposable in-memory secret backend; the private native-helper boundary is exercised without supplying a credential request.

The negative matrix includes arity and typed-flag errors, malformed policies and documents, missing identities, unavailable IPC, occupied and unwritable output paths, invalid daemon configuration, conflicting daemon ownership, rejected authority, cancellation and partial observations. Expected refusals pass only when their documented exit, output and independent state checks agree. Functional defects remain failures and become independently actionable findings.

## Default budgets

| Control | Default |
| --- | --- |
| Fixture application concurrency | One daemon at a time, fresh JSON and human-output fixtures, sequential environments, worker pool 1, GOMAXPROCS 2 |
| Windows CPU | Hard Job rate 25 percent of aggregate system CPU, below-normal priority |
| Linux CPU | Hard cgroup quota 25 percent of one logical CPU, nice 10 |
| Memory | 1 GiB committed Job memory or cgroup memory.max; Linux swap allowance 0 |
| Processes | Windows Job maximum 32 processes; Linux pids.max 128 including threads |
| Per-case deadline | 15 seconds |
| Total coordinator deadline | 30 minutes, with bounded cleanup |
| Per captured channel | 256 KiB |
| Entire artifact directory | 64 MiB, including checkpoints and derived reports |

An optional Budget file supplies every field in the [Budget model](https://github.com/shruggietech/go-schedule/blob/main/internal/clievaluate/model.go). Unknown or missing fields, unbounded values and excessive budgets are rejected. The effective budget is in report.json. Linux workers verify their actual cgroup controls and run identity independently of the Windows WSL launcher. An existing user systemd manager with delegated cpu, memory and pids controllers is required. Missing controls produce blocked results.

## Capture, review and finalization

Capture persists a complete initial JSON and Markdown ledger before starting fixtures. It appends bounded events, persists each completed case separately and refreshes the full worker ledger every 32 cases and at lifecycle boundaries. Recovery replays verified incremental records. It freezes sanitized channels, the corpus and evidence inventory; capture-identity.json binds source, budget, environment and completion metadata to that inventory. events.jsonl records typed, ordered transitions. Runtime state stays in an owned UUID directory. Cleanup stops exact owned process trees, retains ownership evidence and records any remaining resource explicitly.

The executing agent reads every distinct frozen message in its command and scenario context, including generated help, descriptions, success output, diagnostics, MCP responses and daemon logs. It writes reviews.json under the [version 1 reviews schema](https://github.com/shruggietech/go-schedule/blob/main/internal/clievaluate/schemas/v1/reviews.schema.json). Each review binds the complete excerpt, message ID, occurrence IDs and channel to the frozen corpus hash. The seven rubric entries are plain_language, action_outcome, actionable_correction, severity, terminology, identifiers_paths_units and secret_absence. Each entry requires a verdict and a specific reason. Cases with no emitted text have explicit no-message records.

Missing or stale reviews prevent finalization and cannot yield overall pass. Finalized execution failures require actionable findings. Findings identify product or harness defects, exact source and binary identity, affected cases and platforms, reproduction, expected and actual behavior, user impact, hashed evidence spans, functional acceptance criteria and an existing issue URL or null. The coordinator does not edit source or GitHub records.

## Recovery and integrity

Resume uses the original inputs and OutputDirectory. It checks source, executable hashes, manifest, budget, selected distribution/user, previous process creation identity, owned WSL unit state, and retained evidence. An active prior worker blocks resumption until its bounded shutdown completes. Changed inputs require a new capture. A resumed attempt uses a fresh run identity and fresh fixture state, replays stateful scenarios, and retains prior attempt artifacts and cleanup obligations; it does not carry old successes forward as new observations. Cancellation first closes the owned worker input and permits up to seven seconds for ledger persistence and cleanup before exact-tree termination. Any remaining resource stays visible in subsequent attempts.

The validator checks Draft 2020-12 schemas, duplicate JSON keys, required fields, IDs, the full case matrix, expected versus observed checks, path containment, regular files, byte counts, hashes, event order, every nonempty captured channel, review completeness, finding references, summary totals and cleanup. Interrupted captures remain blocked overall. Incompatible future artifact contracts receive a new schema version; version 1 readers reject unknown versions instead of interpreting them as compatible.

## Fixed report outline

1. Run identity and outcome.
2. Environment and budget.
3. Coverage table.
4. Functional results.
5. Help, error and log review.
6. Defects.
7. Cleanup and artifact index.

Canonical JSON is the source for these sections. The deterministic Markdown renderer does not consult a clock, network service, model or current repository state. Keep the complete artifact directory with its relative paths intact when handing the report to another maintainer.

The human report summarizes execution boundaries and communication failures before the defect blocks. Expand its complete case ledger, large affected-case lists and artifact index as needed; canonical JSON retains every successful and failed rubric record.
