# MCP Manage Tool Contract

Manage adds exactly six tools to the existing Observe and Operate server:

| Tool | Actions | Object identity |
| --- | --- | --- |
| `tasks_manage` | create, update, delete | created or supplied task ID |
| `groups_manage` | create, update, delete | created or supplied group ID |
| `chains_manage` | create, update, delete | created or supplied chain ID |
| `triggers_manage` | create, update, delete | created or supplied trigger ID |
| `watchers_manage` | create, update, delete | created or supplied watcher ID |
| `notification_assignments_replace` | replace | exact task or group ID |

Every input includes `daemon_id`, `request_id`, and `confirmed`. Update and delete include `object_id`. Family definition inputs map to the existing API request types after sensitive fields are removed.

Every result is the same redacted envelope: `schema_version`, `permission`, `operation`, `daemon_id`, `request_id`, `object_kind`, `object_id`, `outcome`, and `message`.

Bulk arrays, actor or permission fields, enrollment values, trigger keys, task environment, task stdin, and raw storage paths are not part of any schema.
