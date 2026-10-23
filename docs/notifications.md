---
title: Notifications
nav_order: 8
---

# Notifications

go-schedule can send SMTP email or a versioned JSON webhook for selected task outcomes, persistent problems, recoveries, excessive duration, process start failures, and daemon healthy-presence heartbeats. Delivery is durable and asynchronous: a receiver outage cannot change a task result or occupy a scheduler task worker.

Webhooks and SMTP email provide durable unattended monitoring. The running desktop application also offers optional native operating-system popups for new task outcomes and alerts. SMTP is introduced in the v1.6.0 candidate, which is pending publication. Consult the latest published release before expecting it in an installed package.

## Desktop popups

Open Notifications in the desktop app and use Desktop popups to enable or silence local popups. They are off by default and remember your choice on this computer. Choose failed or successful task outcomes, scheduler alerts and their severity, and which registered daemons may interrupt you. Leaving all scheduler checkboxes empty includes all registered schedulers. These choices do not change webhook assignments, task outcomes, or another desktop's preferences.

Windows, macOS, and Linux graphical sessions with a notification service use the operating system's notification facility. macOS asks for permission when popups are enabled. Denied permission, missing facilities, and a headless session leave scheduling and webhooks unaffected. The app must be running for these popups; closing it does not create a resident popup service or an unattended delivery guarantee. For unattended signals, configure a webhook receiver or SMTP channel.

Each new matching terminal run or alert is attributed to its source daemon and suppressed if the same record is observed again. Old activity is not replayed when the app starts or reconnects. On platforms that report activation, selecting a popup opens its exact Activity record after the app confirms the daemon identity. If activation is unsupported, the record has aged out, or the registration changed, open Activity manually; the app will not substitute another daemon's record.

## Create and assign a channel

Create one reusable webhook channel, then assign it to a task or group. HTTPS is required except for loopback HTTP receivers used in local development.

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

Each webhook request has a five-second timeout. A 2xx response succeeds; transport failures and all other HTTP statuses retry after one second and then two seconds, for no more than three total attempts. Redirects are refused. A daemon restart returns interrupted work to pending without resetting its attempt count.

Delivery is at least once. A receiver can process a request and lose its response, causing a retry. Use `X-Go-Schedule-Delivery` or `delivery.id` as an idempotency key when duplicate side effects matter; the value remains stable across every retry and restart.

## Security boundary

The full webhook endpoint and transport credentials are write-only. On Windows, they are encrypted with DPAPI under the daemon service identity before entering SQLite. On Linux and macOS, the daemon data directory is mode 0700 and the database is mode 0600. Management uses restricted local IPC or the authenticated remote API with operation-specific capabilities. go-schedule does not claim one universal cross-platform encryption-at-rest guarantee; use operating-system disk encryption where that threat matters.

Ordinary channel, policy, delivery, event, log, export, and evidence representations exclude authorization, URL user information, URL paths and queries, task command, arguments, working directory, environment, standard input, and captured output. Webhook evidence identifies the destination by scheme and host. SMTP host, port, sender, recipients and encryption mode are public channel settings; username and password remain write-only. Redirect refusal prevents authorization from moving to another host. Terminal delivery transitions erase their endpoint and authorization snapshots.

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

Filter evidence with `--channel`, `--task`, `--run`, `--state`, and `--limit`. Delivery state is `pending`, `claimed`, `succeeded`, or `failed`. The last HTTP or SMTP status and bounded authored detail describe the transport outcome without echoing a credential or arbitrary server response. SMTP `diagnostic` carries neutral delivery detail for successful server acceptance as well as incomplete attempts. `last_error` is empty for succeeded SMTP deliveries and retains genuine failure detail for retrying or failed deliveries. CLI human output and desktop delivery details distinguish these fields. Migration v24 preserves historical SMTP text in `diagnostic`, clearing only succeeded SMTP errors; webhook history keeps its existing semantics.

## SMTP email channels

Introduced in the v1.6.0 candidate; publication is pending.

SMTP channels use the same task/group policies, health heartbeats, durable queue and delivery history as webhooks. Create a channel in Notifications in the desktop application or use the CLI:

```sh
gosched notification channel add "Operations email" --kind smtp --smtp-host smtp.example.com --smtp-port 587 --smtp-tls starttls --smtp-from scheduler@example.com --smtp-to ops@example.com,backup@example.com --smtp-timeout 10s --smtp-html --smtp-auth-file ./smtp-auth.json
gosched notification channel test <channel-id>
gosched notification task set <task-id> --channel <channel-id> --on failure --recovery
gosched notification deliveries --channel <channel-id>
```

The private UTF-8 credential file contains `{"username":"SMTP username","password":"SMTP password"}`. Restrict its file permissions and remove it when no longer needed. Use `--smtp-auth-file -` to read this object from standard input. Credentials are optional for relays that accept anonymous submission. A configured username and password must both be present; failed or unsupported authentication never falls back to anonymous submission. Avoid putting credentials in host names, sender fields, task commands or environment fields.

| Setting | Meaning |
| --- | --- |
| `--smtp-host` | DNS name or IP address, without a URL scheme or port |
| `--smtp-port` | Port 1-65535; omitted/zero selects the mode's default |
| `--smtp-tls starttls` | Required STARTTLS, default port 587; failure does not downgrade |
| `--smtp-tls tls` | TLS from connection start, default port 465 |
| `--smtp-tls plain` | Explicit unencrypted relay, default port 25 |
| `--smtp-from` | One sender mailbox |
| `--smtp-to` | Repeatable or comma-separated list of 1-100 unique recipient mailboxes |
| `--smtp-timeout` | One whole-second budget for the entire attempt, 1s-60s, default 10s |
| `--smtp-html` | Escaped HTML alternative alongside the always-present plain text body |

TLS modes verify the certificate chain and hostname using the daemon's trusted roots, with TLS 1.2 as the minimum. Plain mode sends both the message and any configured authentication without encryption, so select it only for the intended unencrypted relay. Sender and recipient envelope addresses must be ASCII, including correctly quoted local parts when needed. For example, pass a recipient list as `--smtp-to '"ops,night"@example.com,ops@example.com'` so the comma inside the quoted mailbox stays part of that address; Unicode task names and event text are supported in MIME subjects and bodies. Header controls and malformed addresses are rejected before a delivery is queued. SMTP PLAIN and LOGIN authentication are supported when advertised by the server.

Email includes the task/run or daemon-health context and stable delivery identity. It excludes task commands, arguments, environment, input and captured output. Messages are bounded to 256 KiB. The subject and optional HTML are generated from the same event as the webhook, with HTML text escaped.

```sh
gosched notification channel update <channel-id> --smtp-to new-ops@example.com --smtp-html=false
gosched notification channel rotate <channel-id> --smtp-auth-file ./replacement-auth.json
gosched notification channel rotate <channel-id> --smtp-auth-file -
```

For the final command, supply `{}` on standard input to clear authentication explicitly. Updating settings retains credentials; rotation changes credentials without returning the previous values. The desktop form offers equivalent replace/remove controls and never prefills stored credentials. Transport kind is immutable; create another channel to change between webhook and SMTP.

For [Amazon SES SMTP](https://docs.aws.amazon.com/ses/latest/dg/send-email-smtp.html), use its regional SMTP endpoint, verified sender and region-specific SMTP credentials. SES SMTP credentials differ from AWS access keys and console sign-in credentials. For Google accounts that support SMTP with an [app password](https://support.google.com/accounts/answer/185833), use the account's SMTP username and generated app password; availability depends on the account and its two-step verification settings. Provider-specific OAuth setup is outside this username/password transport.

### SMTP results and retries

Each attempt has one shared timeout across connection, greeting, TLS, authentication, recipients and message submission. There are at most three attempts, with the existing one-second and two-second retry delays. Recipient-level history distinguishes `pending`, `submitting`, `accepted`, `temporary_failure`, `permanent_failure` and `uncertain`. A 4xx rejection is temporary; a 5xx rejection is permanent. Valid recipients can receive the message while other recipients fail.

A recipient becomes accepted only after the server confirms completion of DATA. Confirmed recipients and permanent rejections are excluded from later attempts. A failed QUIT after that confirmation does not undo acceptance. An interrupted submission without a saved confirmation becomes uncertain and may be retried within the same budget; history explicitly warns that this can duplicate email. The stable Message-ID helps correlate attempts but does not force a mail server to deduplicate. Server acceptance records handoff to that server, not delivery to a recipient's inbox.

Edits and credential rotation affect newly queued work. Existing work retains its original settings and credential snapshot. Disabling a channel stops new work; already queued work keeps its existing lifecycle. Removing a channel removes unfinished work and assignments. Terminal deliveries erase their credential snapshots and retain bounded, authored diagnostics and recipient results without copying server response text.

### Remote access and portable configuration

Local IPC and authenticated remote administration expose the same notification operations. Remote reads require Observe, channel/policy changes require Manage, a test requires Operate, and credential rotation requires Enroll. Retargeting a channel that already has stored credentials to another webhook endpoint or SMTP host/port/encryption mode also requires Enroll. This prevents ordinary settings management from redirecting credentials entrusted to another destination. The desktop keeps history and assignment inspection available to read-only connections and disables actions that need more authority. Select the intended daemon in the desktop application or CLI profile.

Portable v2 exports include typed nonsecret `smtp_channels` definitions, enablement, health interval and a credential-presence requirement. They never contain usernames or passwords. Before applying to another daemon, provision exactly one SMTP channel with the same logical name and the required local credentials. If credentials are present, its host, port and encryption mode must already match. Preview freezes the selected channel ID and revision; missing, ambiguous, replaced or changed bindings block apply before task or policy writes. A successful apply changes nonsecret settings and retains the target daemon's credentials. Duplicate source names produce an explicit export exclusion rather than guessing which identity to transfer.
