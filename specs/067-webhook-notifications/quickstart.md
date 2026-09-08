# Quickstart: Dependable Webhook Notifications

## Prerequisites

- Build `goschedd` and `gosched` from the S067 review branch.
- Start the daemon with an isolated data directory and restricted local IPC configuration.
- Start a loopback HTTP receiver that records headers and JSON bodies and can return controlled status codes.

## Validate a Successful Task Notification

1. Create a webhook channel pointing at the loopback receiver with a recognizable authorization test value.
2. Create a task that exits successfully.
3. Replace the task's notification assignments with the channel configured for success only.
4. Run the task manually and wait for one terminal delivery record.
5. Validate the received body against `contracts/webhook-v1.schema.json` and correlate its delivery, task, and run IDs.
6. Search API output, CLI output, daemon logs, events, and delivery history for the submitted endpoint query and authorization values; neither may occur.

## Validate Group Precedence

1. Create a parent group, child group, and task in the child group.
2. Assign different channels at both groups and query the task's effective policy; only the child group is selected.
3. Add a task assignment and query again; only the task scope is selected.
4. Replace the task assignment list with an empty list; the child group becomes effective again.

## Validate Failure, Retry, and Restart

1. Configure a receiver to reject requests and run a matching task.
2. Observe no more than three attempts using one `X-Go-Schedule-Delivery` value and a terminal failed record.
3. Configure a receiver to hold a request, stop the daemon after the delivery is claimed, and restart it.
4. Observe recovery with the same delivery ID and preserved attempt count, never exceeding three total requests.
5. Confirm the source task run retains its original outcome throughout.

## Validate Lifecycle and Removal

1. Test, disable, enable, rename, change endpoint, rotate authorization, and list the channel through both API and CLI.
2. Complete at least one test delivery, then remove the channel.
3. Confirm assignments and unfinished work are gone, the completed history remains with safe snapshots, and no secret is returned.

## Verification

Run focused notification, store, engine, API, client, CLI, and installed-daemon tests, then run `sh scripts/verify.sh all` in the foreground. Record all eight canonical gates in `verification.md`.
