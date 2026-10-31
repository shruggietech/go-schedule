# S132 First-use Contract

The guide supplies platform scripts and an explicit local workspace before task creation. A successful snapshot copies source.txt to snapshot.txt and emits a completion message. Missing input emits an actionable explanation and exits nonzero without overwriting prior snapshot output. Restoring source.txt and running the same task yields a new successful record.

Run now acknowledgement is accepted work, not completion and not a run identifier. Use runs --task TASK_ID --details to find the particular new run, runs active for active executions and runs get RUN_ID for its exact state/output. In the desktop, select a run in Activity and inspect its own Run ID, State, Exit status and Combined stdout/stderr in Run details; Source run remains separate.

Local package and changed-source/desktop evidence remain distinct. No automatic destructive cleanup or access-control bypass is introduced. Current endpoints, schemas and output contracts are unchanged.
