# Portable Sources Contract

`go-schedule.bundle/v2` adds `external_triggers`, `trigger_sets`, `watchers`, and `notification_policies`. The daemon continues to accept `go-schedule.bundle/v1` with its original fields. Export always writes v2. Preview and compare accept `watcher_paths` keyed by watcher portable identity; this map is not a document field and never appears in canonical JSON or the bundle digest. Apply uses the server-retained preview and does not accept new bindings.

Preview reports a conflict for a new watcher without a target-local path or a notification assignment whose channel name is absent or ambiguous. Fresh trigger credentials never appear in a bundle or apply response.
