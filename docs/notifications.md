---
title: Notifications
nav_order: 8
---

# Notifications

go-schedule can send a versioned JSON webhook for selected task outcomes, persistent problems, recoveries, excessive duration, process start failures, and daemon healthy-presence heartbeats. Webhook delivery is durable and asynchronous: a receiver outage cannot change a task result or occupy a scheduler task worker.

Webhook is the durable notification transport for unattended monitoring. The running desktop application also offers optional native operating-system popups for new task outcomes and alerts. SMTP email is planned separately and is not available.

## Desktop popups

Open Notifications in the desktop app and use Desktop popups to enable or silence local popups. They are off by default and remember your choice on this computer. Choose failed or successful task outcomes, scheduler alerts and their severity, and which registered daemons may interrupt you. Leaving all scheduler checkboxes empty includes all registered schedulers. These choices do not change webhook assignments, task outcomes, or another desktop's preferences.

Windows, macOS, and Linux graphical sessions with a notification service use the operating system's notification facility. macOS asks for permission when popups are enabled. Denied permission, missing facilities, and a headless session leave scheduling and webhooks unaffected. The app must be running for these popups; closing it does not create a resident popup service or an unattended delivery guarantee. For unattended signals, configure a webhook receiver.

Each new matching terminal run or alert is attributed to its source daemon and suppressed if the same record is observed again. Old activity is not replayed when the app starts or reconnects. On platforms that report activation, selecting a popup opens its exact Activity record after the app confirms the daemon identity. If activation is unsupported, the record has aged out, or the registration changed, open Activity manually; the app will not substitute another daemon's record.

## Create and assign a channel

Create one reusable channel, then assign it to a task or group. HTTPS is required except for loopback HTTP receivers used in local development.

```sh
gosched notification channel add "Operations" --endpoint https://receiver.example/hooks/go-schedule --authorization "Bearer replace-me"
gosched notification channel list
gosched notification task set <task-id> --channel <channel-id> --on failure
gosched notification task effective <task-id>
gosched notification channel test <channel-id>
gosched notification deliveries --channel <channel-id> --limit 20
```

Use `--on success`, `--on failure`, `--on success,failure`, or `--on none` with another condition. Repeat `--channel` to assign multiple destinations with the same conditions. Calling `task set` or `group set` without a channel clears that scope and resumes inheritance.

## Conditions and suppression

```sh
gosched notification task set <task-id> --channel <channel-id> --on failure --failure-threshold 3 --failure-to-start --duration-threshold 30m --recovery --reminder 2h
gosched notification group set <group-id> --channel <channel-id> --on success --quiet-period 24h
gosched notification channel update <channel-id> --health-interval 5m
```

Failure thresholds count consecutive failed terminal runs and reset on a successful terminal run or a policy change. A threshold of one reports the first failure. Failure-to-start means the operating system could not start the configured process, not merely that the process returned a failing exit code. Duration is evaluated after completion and does not stop the task. When one run matches several problem conditions, one notification is created using this precedence: failure to start, consecutive failure, then duration exceeded.

The first transition into a problem creates one notification. Repeated runs matching the same active problem are suppressed unless a reminder interval has elapsed. A different problem creates a new notification. Recovery is optional and fires on the first run that clears the active problem. Success notifications are optional and a quiet period suppresses repeated routine successes. Policy changes reset prior counters and active-problem state so obsolete settings cannot produce a delayed notification. State is stored in SQLite and therefore survives daemon restart; runs that never reached durable history are not reconstructed as missed evaluations.

Daemon health uses opt-in healthy-presence heartbeats because a stopped process cannot send its own outage message. Set `--health-interval` on a channel. Each heartbeat includes `daemon.status=healthy` and `daemon.next_expected_at`; the receiver reports an outage when that deadline passes without another heartbeat. Zero disables heartbeats. Disabled channels do not create heartbeats.

## Policy precedence

One scope supplies the complete effective policy. Direct task assignments replace all inherited assignments. Without direct task assignments, the nearest ancestor group containing assignments supplies the policy and replaces more distant group policies. A task with no configured scope in its ancestry sends no notifications. This replacement rule prevents duplicate alerts and makes `gosched notification task effective` explain the source unambiguously.

## Receiver contract

Every request is `POST` with `Content-Type: application/json`, `User-Agent: go-schedule/<version>`, `X-Go-Schedule-Delivery: <delivery-id>`, and `X-Go-Schedule-Event: run.completed`, `daemon.health`, or `test`. Run notifications include a safe `condition` object with the reason, explanation, threshold context, and reminder state. Heartbeats omit task and run data. The additive JSON contract is documented in [`specs/095-notification-conditions/contracts/notification-conditions.md`](https://github.com/shruggietech/go-schedule/blob/main/specs/095-notification-conditions/contracts/notification-conditions.md).

```text
{
  "schema": "go-schedule.webhook.v1",
  "event": "run.completed",
  "delivery": {
    "id": "8dc4da91-8094-42c7-898b-5bfeb7627712",
    "created_at": "2026-09-07T16:30:00Z"
  },
  "daemon": {
    "version": "1.3.0"
  },
  "task": {
    "id": "task-id",
    "name": "Nightly backup",
    "group_id": "group-id",
    "group_name": "Operations"
  },
  "run": {
    "id": "run-id",
    "outcome": "success",
    "trigger": "schedule",
    "scheduled_for": "2026-09-07T16:30:00Z",
    "started_at": "2026-09-07T16:30:00.012Z",
    "ended_at": "2026-09-07T16:30:02.115Z",
    "duration_ms": 2103,
    "exit_code": 0
  }
}
```

A minimal Go receiver needs no go-schedule package:

```text
http.HandleFunc("POST /go-schedule", func(w http.ResponseWriter, r *http.Request) {
    var event struct {
        Schema string `json:"schema"`
        Event string `json:"event"`
        Delivery struct{ ID string `json:"id"` } `json:"delivery"`
    }
    if r.Header.Get("X-Go-Schedule-Delivery") == "" || json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10)).Decode(&event) != nil || event.Schema != "go-schedule.webhook.v1" {
        http.Error(w, "invalid webhook", http.StatusBadRequest)
        return
    }
    log.Printf("received %s delivery %s", event.Event, event.Delivery.ID)
    w.WriteHeader(http.StatusNoContent)
})
log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
```

Generic automation systems, chat relays, incident routers, and private HTTP services can consume the same contract. Put vendor transformation and authentication details at the receiver rather than in task commands or environment fields.

## Retries and duplicate handling

Each request has a five-second timeout. A 2xx response succeeds; transport failures and all other HTTP statuses retry after one second and then two seconds, for no more than three total attempts. Redirects are refused. A daemon restart returns interrupted work to pending without resetting its attempt count.

Delivery is at least once. A receiver can process a request and lose its response, causing a retry. Use `X-Go-Schedule-Delivery` or `delivery.id` as an idempotency key when duplicate side effects matter; the value remains stable across every retry and restart.

## Security boundary

The full endpoint and optional Authorization value are write-only. On Windows, they are encrypted with DPAPI under the daemon service identity before entering SQLite. On Linux and macOS, the daemon data directory is mode 0700 and the database is mode 0600. Management remains behind restricted local IPC. go-schedule does not claim one universal cross-platform encryption-at-rest guarantee; use operating-system disk encryption where that threat matters.

Ordinary channel, policy, delivery, event, log, export, and evidence representations exclude authorization, URL user information, URL paths and queries, task command, arguments, working directory, environment, standard input, and captured output. Only scheme and host identify the destination in evidence. Redirect refusal prevents authorization from moving to another host. Terminal delivery transitions erase their endpoint and authorization snapshots.

## Lifecycle and evidence

```sh
gosched notification channel get <channel-id>
gosched notification channel update <channel-id> --name "Primary operations" --endpoint https://new.example/hooks/go-schedule
gosched notification channel rotate <channel-id> --authorization "Bearer replacement"
gosched notification channel disable <channel-id>
gosched notification channel enable <channel-id>
gosched notification channel rm <channel-id>
```

Disabling a channel prevents new delivery creation. Removing it atomically removes assignments and unfinished work and erases its stored credentials. Completed delivery evidence remains with safe channel and destination snapshots. The daemon retains the newest 1,000 succeeded or failed delivery records; pending and claimed work is never pruned by that limit.

Filter evidence with `--channel`, `--task`, `--run`, `--state`, and `--limit`. Delivery state is `pending`, `claimed`, `succeeded`, or `failed`. The last HTTP status and a bounded safe diagnostic explain terminal failure without echoing an endpoint or credential.
