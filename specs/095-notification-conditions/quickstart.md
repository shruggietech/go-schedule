# Quickstart: Actionable Notification Conditions

## Consecutive failure and recovery

```sh
gosched notification task set <task-id> --channel <channel-id> --on failure --failure-threshold 3 --recovery --reminder 30m --quiet-period 10m
gosched notification task show <task-id>
gosched notification task effective <task-id>
```

The third consecutive failure opens the problem. Equivalent failures remain quiet until the 30-minute reminder interval. The first run that no longer matches a configured problem emits one recovery.

## Failure to start and duration

```sh
gosched notification group set <group-id> --channel <channel-id> --on none --failure-to-start --duration-threshold 15m --recovery
```

`--on none` disables routine success and ordinary failure outcomes while the advanced problem conditions remain active.

## Daemon-health heartbeat

```sh
gosched notification channel update <channel-id> --health-interval 5m
```

The destination receives a healthy-presence event every five minutes while the daemon and notification pipeline are operating. Configure the receiver to alert after at least two missing intervals. A resumed heartbeat is recovery evidence; go-schedule does not claim it can emit while stopped.

## Inspect why a delivery fired

```sh
gosched notification deliveries --channel <channel-id> --limit 20
```

Condition deliveries include a safe condition kind and explanation. Open Notifications in the desktop for the same explanation, current effective policy, inheritance source, and advanced condition controls.
