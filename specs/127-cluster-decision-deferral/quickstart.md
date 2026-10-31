# S127 Validation Guide

1. Read [spec.md](spec.md), [decision.md](decision.md) and the [current page](../../docs/distributed-scheduling.md). The direction must be explicit owner deferral, not a permanent rejection or claimed completed research.
2. Follow the four-stage sequence and locate the owner revisit and maintainer semantics approval boundaries.
3. Compare all six issue dispositions with [contracts/decision.md](contracts/decision.md); retained issues must remain open.
4. Run the actual publication formatter, then sh scripts/verify.sh all in the foreground using the verified hidden Windows launcher. Record original output and failures in .test-temp/s127.
5. Run the existing site browser contract for the modified documentation route and lifecycle check; retain source/evidence drift receipts. No attended native runtime evidence is claimed.
6. Before PR publication, inspect the diff for source/runtime, release or historical evidence mutations. After publication, verify exact-head CI and all comments/threads/body reactions. Stop before owner merge.
