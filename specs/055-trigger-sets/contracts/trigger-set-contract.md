# Contract: Trigger Sets

## Local API

| Method | Path | Result |
| --- | --- | --- |
| POST | `/v1/trigger-sets` | Atomically create a set and return ordered secrets |
| GET | `/v1/trigger-sets` | List redacted sets with ordered redacted members |
| GET | `/v1/trigger-sets/{id}` | Return one redacted set with ordered redacted members |
| PATCH | `/v1/trigger-sets/{id}` | Atomically retarget every member |
| DELETE | `/v1/trigger-sets/{id}` | Atomically delete the set and members |
| POST | `/v1/trigger-sets/{id}/enable` | Atomically enable every member |
| POST | `/v1/trigger-sets/{id}/disable` | Atomically disable every member |
| POST | `/v1/trigger-sets/{id}/rotate` | Atomically rotate all keys and return ordered replacement secrets |
| POST | `/v1/trigger-sets/{id}/reveal` | Return ordered current secrets without mutation |

Create accepts `name`, `target_task_id`, `count`, and optional `enabled`. Retarget accepts `target_task_id`. Count outside 1 through 99, missing names or targets, individual member target changes, and missing resources use the standard API error envelope with stable codes and actionable field names.

Ordinary set responses contain set identity, name, target identity and label, member counts, timestamps, and ordered ordinary trigger responses. They never contain raw keys. Secret responses contain the ordinary set response plus ordered members containing `position`, `trigger_id`, `key`, and `command`.

## CLI

```text
gosched trigger set create --name <name> --task <task-id> --count <1..99> [--disabled]
gosched trigger set list
gosched trigger set show <set-id>
gosched trigger set retarget <set-id> --task <task-id>
gosched trigger set enable <set-id>
gosched trigger set disable <set-id>
gosched trigger set reveal <set-id>
gosched trigger set rotate <set-id>
gosched trigger set rm <set-id>
```

Every command honors global `--json`. Human create, reveal, and rotate output writes exactly one complete invocation command per nonblank line in ascending position with one final newline. JSON secret output preserves the API structure.

## Desktop

The Triggers table gains Set and Position columns. Standalone triggers display `Not assigned`; members display the set name and permanent position. A New Set action accepts name, task, count, and initial enabled state. Selecting any member enables Copy set commands, Retarget set, Enable or disable set, Rotate set, and Delete set actions. Broad and destructive actions confirm the set name and current member count. Secret dialogs expose ordered commands and visible copy confirmation without placing raw keys in ordinary widgets or notifications.

## Events and Redaction

One `trigger_set` lifecycle event is published after each successful broad mutation. It contains set identity, lifecycle verb, target identity, and member count, with no key or command. Individual member operations retain ordinary `trigger` events. Failed transactions publish no lifecycle event.
