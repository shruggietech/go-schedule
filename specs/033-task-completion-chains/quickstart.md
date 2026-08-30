# Quickstart: Validate Task-Completion Chains

## Primary Flow

Create two deterministic tasks, then:

```sh
gosched chain create --source <source-id> --target <target-id> --on success
gosched task run <source-id>
gosched runs --task <target-id>
```

Expected: one target run identifies trigger `completion`, source task, and source run. Both tasks keep their normal schedules.

## Failure Selection

Create failure and any chains from a failing source. One failure should match both while a success-only relationship does not.

## Validation Safety

Attempt a self-link. Then create A to B and B to C and attempt C to A. Every invalid command exits as validation failure and chain state remains unchanged.

## Restart Recovery

Automated integration persists pending, claimed, and completed deliveries, reopens the store/engine, and demonstrates that pending dispatches, claimed becomes replayable, completed does not repeat, and source correlation remains stable.

## Desktop Flow

Open **Chains**, select source, target, and outcome, create the relationship, edit the outcome, then delete it. Task and chain changes refresh live. Empty and invalid states remain actionable.

## CI-Parity Verification

```sh
sh scripts/verify.sh all
```

All eight gates must pass. Targeted store, engine, and fan-out evidence is recorded in `verification.md`.
