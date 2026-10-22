# S112 release requirement checklist

- [x] CHK001 Are queued, skipped, completed and missing observations distinguished explicitly? [FR-001]
- [x] CHK002 Are primary failure, diagnostic collection and cleanup outcomes independent? [FR-002]
- [x] CHK003 Are setup and runtime lifetimes independent with unchanged authentication assertions? [FR-003]
- [x] CHK004 Is every dependency proposal accounted for, including refreshed and overlapping changes? [FR-004]
- [x] CHK005 Are permissions, artifact exclusions, environment and domain preserved for Pages upgrades? [FR-005]
- [x] CHK006 Are candidate and published comparison states distinguishable? [FR-006, FR-007]
- [x] CHK007 Does #354 remain open across the source/publication boundary? [FR-008]
- [x] CHK008 Are local service mutation and foreground windows excluded? [FR-009]
- [x] CHK009 Are exact-head CI and every review thread required before handoff, with two rounds maximum? [FR-010]
