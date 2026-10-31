# S132 Verification Guide

1. Review spec.md, plan.md and contracts/first-use.md, then the complete canonical docs/first-task.md after implementation.
2. Run the ActivityPage component regression before adding the Run ID field; retain its expected failure. Apply the presentation correction and rerun under the existing frontend test command.
3. Verify immutable downloaded MSI/Linux archive hashes, extract only to owned private paths and launch explicit isolated runtimes with existing IPC controls. Execute the guide success/failure/recovery operations and inspect exact records; stop owned runtimes and retain outputs. No installed-service or native GUI acceptance claim.
4. Run docs export and existing site browser cases, including the first-task entry/route test, with fresh screenshot paths. Verify current labels and release claims from production components.
5. Run sh scripts/verify.sh all in the foreground with source/input pre/post hashes, then final format/lifecycle/immutable audits. Bind any later factual completion records separately.
6. Commit, automatically push/open official PR, wait for exact-head CI and at most two review requests, handle feedback and hand off unmerged.
