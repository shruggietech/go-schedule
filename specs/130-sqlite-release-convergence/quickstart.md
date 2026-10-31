# S130 Validation Guide

Use a headless, noninteractive runner on Windows and watch foreground commands to terminal completion. Required prerequisites are existing Go1.26, C/race toolchain, POSIX shell and Wails platform/frontend tools.

1. Run the focused storage driver tests and repeated prepared benchmark on the held pair before updating; retain the expected lock-cancellation failure and raw measurements privately.
2. Select SQLite1.61.1/libc1.79.0 in root, retain the expected stale-desktop read-only failure, regenerate desktop metadata, then run the same focused tests on the coherent pair.
3. In both root and desktop, run go list -mod=readonly -m all and go mod verify; run desktop connection tests with -mod=readonly and compare metadata hashes before/after.
4. Compare repeated same-host benchmark medians using identical source/workloads. Investigate changes exceeding the existing 10 percent policy rather than hiding them.
5. Run sh scripts/verify.sh all in the foreground and require all eight gates. Existing root and desktop race/security/current-state suites remain unchanged.
6. Run existing release-source-check for v2.0.0 and comparison-check with the actual eventual commit. Check candidate note/changelog digests, unchanged comparison/history/public selection and all 221 frozen inputs.
7. After authorized source push, observe actual final-head CI and every review/comment/thread/body reaction. Consume at most two total Codex requests including the automatic first request. Stop for owner final review/merge; no scheduled checks.

Detailed actual-source, timing, failure and terminal receipts belong in verification.md and private .test-temp/s130, separate from any later staged-package acceptance. Continue actual release operations through ../129-prelaunch-release/handoff.md after its owner decisions.
