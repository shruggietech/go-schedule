# S143 verification

## Candidate and boundary

Accepted main baseline: 3d772b99c7684062c80cebb410d4dcf3d0d489f2. Review branch: codex/143-active-run-cancellation. Public v2.0.3 remains unchanged. No release or native desktop interaction is part of this source increment.

## Focused evidence

- Executor: failing-first explicit cancellation/tree and Windows breakaway regressions retained. Final Windows race suite passed twice (5.064s), actual Linux/WSL race suite passed twice (4.175s), Windows vet and macOS arm64 test compilation passed. Slow startup, one original cleanup deadline, failure feedback, unreaped ownership, completion races, input/environment/working-directory/run-as and intentional breakaway are covered.
- Engine/API: exact identities, repeated requests, natural completion/finishing persistence, delayed cleanup, overlap skip/queue reservations, canceled chain suppression, local/remote authority, identity, revocation, mutation audit and RetryNever checks passed under race. The original #466 and #467 baseline failures are retained.
- Store: cancellation preserves failure/recovery streaks, ignores duration policies, suppresses any/success/failure completion deliveries, retains native result/output and reopens valid current schema-24 data.
- Interfaces: client acknowledgement validation, human/JSON CLI, MCP deduplication/revocation and bounded active discovery passed; full desktop module race and TypeScript passed. Active MCP pages retain exact IDs and authored cleanup guidance, exclude raw output/execution inputs, and fail explicitly on unavailable active reads.
- Desktop: component regressions cover capability/authority/current-target guards, synchronous duplicate submission, target changes, exact Refresh (#468), bounded reconciliation, timeout recovery and uncertainty retained across a failed Refresh. All 175 frontend tests and headless interaction evidence passed. Browser screenshots use synthetic fixtures only.
- Full desktop browser coverage: 35 of 36 cases passed on private port 4183; the local-assets assertion expected canonical fixture port 4173 and misclassified local requests. The sole unresolved case passed on canonical port 4173. No production or test assertion was weakened. Both receipts remain retained.
- Actual packaged CLI/daemon: native Windows and Linux/WSL independently stopped one owned parent/child tree while its same-task sibling remained active, retained native exit/output, handled repeated cancellation without duplicate history, canceled the sibling through human output and retained healthy daemon service. The first WSL attempt failed before product execution because host mirrored networking could not initialize and a dependency was uncached. The successful retry used the existing trusted workstation module cache as an offline proxy without changing host or distro networking.
- CLI catalogue: new command metadata was regenerated from initialized Cobra metadata after a retained drift failure. Finished-run human/JSON and missing-ID behavioral cases accompany the real active packaged integration.

Private failure and green receipts are retained in .test-temp/s143/. Raw local paths and process identifiers are not published. Historical retained inputs remain protected by the baseline hash manifest.

## Paired performance

Same host/toolchain, GOMAXPROCS=2, 50 iterations and five samples per accepted-main/candidate benchmark. Median dispatch: 91,102 ns to 92,462 ns (+1.49%). Median completion fan-out 100: 23,662,560 ns to 22,882,368 ns (-3.30%). Both remain within the existing ten-percent budget; the existing p99 dispatch contract remains mandatory.

## Canonical aggregate

The watched foreground command `sh scripts/verify.sh all` passed all eight gates on 2026-11-10 in 901.596 seconds: format, vet, lint, race, GUI, coverage, docs and automation. All 2,771 frozen source paths and 8,708 protected historical inputs had zero drift. The root integration package passed in 264.855 seconds, the desktop race suite/build and all 175 frontend tests passed, and automation passed 97 policy assertions plus 30 public-adapter cases. Coverage: engine 85.4%, schedule 89.1%, timezone 91.3%, store 81.0%, catchup 88.9%, logbus 91.1%. Earlier import-grouping/formatting failures remain retained, were corrected without weakening checks, and preceded the successful complete run.

Only factual lifecycle/delivery records and availability/pagination documentation clarifications followed the frozen implementation run. Their final format, documentation, lifecycle, encoding, privacy and historical-integrity checks are recorded at the publication checkpoint; no Go, TypeScript, OpenAPI or command inventory source changed after the successful freeze. Exact final-head hosted CI remains a separate PR receipt.

## Limits

Headless React/browser evidence does not establish native Wails window behavior. No operator focus, native GUI launch, service installation or host network mutation was used. Actual macOS execution remains for hosted CI; local macOS evidence is cross-compilation only. Explicitly escaped process work remains outside the run-owned job/group target. Cleanup expiry remains unconfirmed and owned until genuine reaping.
