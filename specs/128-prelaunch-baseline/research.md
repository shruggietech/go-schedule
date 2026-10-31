# S128 Research Decisions

**Date**: 2026-10-31

Three read-only research investigations mapped application/metrics storage, preference/QA/UI contracts and policy/qualification dependencies. Reports are retained in .test-temp/s128; no source implementation was changed by research.

## Direct application baseline

Decision: Embed the accepted final22-table/33-index DDL, retain marker 24 and initialize atomically once. Existing current data/identity stays unchanged; unsupported state is refused before persistent changes. Alternatives: retaining 24 replay steps contradicts owner scope; renumbering the unchanged current application schema adds disruption without value; automatic reset violates explicit developer control.

Retain current behavior hidden in migration-named tests, particularly local actor/audit, grant expiry, conditions, bounded query indexes, portable identity, elapsed epoch, chains and current notification state. Retire only seeded historical conversion expectations.

## Current preferences

Decision: Remove Transition, legacy file discovery/import, missing-popup conversion and migration UI from shared store, desktop, CLI and evaluator. Keep current atomic-write/lock/patch/reset/validation contracts. Ordinary unknown JSON keys remain ignored rather than interpreted. Missing current popup filters fail existing validation; absent files initialize valid defaults.

## Current metrics and qualification

Decision: Both helper metrics writers/readers use one complete current layout4, supporting shared explicit database paths without old-format conversion. Preserve fresh/current probe vocabulary, skipped/unavailable/no-usable-snapshot reporting, constraints and concurrent writers.

Decision: Fresh-only qualification removes historical baseline roles, upgrade.wsb, retained-profile/setup.upgrade requirements and exact v1.1.1 identities. Full current evidence uses revision 2; archived revision 1 stays immutable. Keep current package identity, trusted environment, hashes, one-time binding, access, repair, rollback and bounded removal assertions.

## Owner policy and exclusions

Decision: Formally amend constitution 3.1.0 to 4.0.0 and synchronize active guidance. Historical constraints no longer govern this no-user prelaunch baseline. Preserve all current security, data correctness, scheduling and contribution gates.

Keep current driver versions; #355/#356/#400/#404 and deferred cluster work are not implemented by this slice. No tag, release or attended foreground test is authorized by source autopilot.
