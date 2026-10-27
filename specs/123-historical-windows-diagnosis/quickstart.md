# S123 validation guide

1. Verify the original source index, dirty overlay and six prepared binaries without changing originals. Recover source in ignored scratch and prove its recorded digest before compiling its original integration test.
2. Prepare current source-bound binaries and race-enabled integration/IPC test executables. Use each test package's normal working directory. Record preparation timing independently.
3. Run replay selection regressions before changing the existing replay controls; retain their first failures.
4. Freeze the declared matrix and a single UTC deadline thirty minutes after execution starts. Invoke compiled tests sequentially with failfast, selected repetitions and remaining timeout; preserve six concurrent writers.
5. For packaged replay, select GOSCHEDULE_REPLAY_BIN, positive whole-millisecond GOSCHEDULE_REPLAY_CASE_TIMEOUT and explicit GOSCHEDULE_REPLAY_INCIDENTS. Historical private collection is unavailable, not a successful current collector observation.
6. Preserve all original result channels, exact expected exits, state/credential assertions and owned cleanup. Stop unexpected conditions for analysis under the same deadline.
7. Publish the six-section investigation using the existing maintainer recipe. Correct demonstrated application defects with regressions and then run existing integrated QA. Run canonical sh scripts/verify.sh all for the final source proposal.
8. Preserve original artifacts, report unresolved causes truthfully, push the review branch and handle at most two authorized Codex rounds. Stop before merge.
