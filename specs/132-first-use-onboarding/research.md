# S132 Research and Decisions

## Existing first-use contracts

Decision: reuse native shell scripts, the local task API and existing run history. Rationale: the product already supports file work, explicit working directories, inactive tasks, immediate runs and exact run inspection. Alternatives: new wizard or test framework add no demonstrated value.

## Demonstrated gaps

Linux and macOS First task examples point to /usr/local/bin/backup.sh without supplying it. docs/gui-fields.md still labels cron and SMTP unreleased and describes obsolete editor/Activity behavior. Current Create task, Preview, Save inactive task, Timing mode, One time, Run at, Schedule syntax and Run details/Combined stdout/stderr labels are authoritative. docs/cli.md also has an Unreleased S111 heading.

## Exact run correlation

Decision: display own run.id in Run details, separate from sourceRunId. Rationale: the CLI supports runs get ID and run details already include source identity, but its own selected stable ID is absent. Alternatives: infer from task/time or add a new API were rejected because the existing ID solves the demonstrated visibility gap.

## Package evidence and native boundary

Decision: use exact published v2.0.0 bytes in isolated hidden local/WSL runtimes, with source/UI evidence separately attributed. MSI extraction is read-only and does not establish attended setup or SYSTEM execution acceptance. No GUI or UAC is launched on the active workstation.
