# S107 data model

## Invocation

One command tree owns JSON mode, target flags, switchable API client, parent cancellation and injected input/output/error, clock, profile/preference stores, secrets, service operations and GUI launcher. Local persistence commands skip daemon resolution. No mutable client/output/MCP callback is process-global.

## Edit intent

Omitted retains state, replace supplies a value, explicit clear removes task overrides/input. ClearEnv does not erase the inherited process environment. Set/clear conflicts fail before mutation. Group parent nil means unchanged, empty means root, ID means reparent. Disabled create supplies enabled=false in one request.

## Notification request

An assignments array holds channel IDs and independent supported conditions. Empty restores inheritance. Missing/null array, unknown fields, trailing data and conflicting shorthand fail before replacement.

## Preferences and startup connection

Preference version 1, appearance, transition status/retired keys and popup enabled/conditions/severities/daemonIds preserve existing names/defaults. Scoped updates lock, read/migrate, validate and atomically replace. Invalid requested values do not trigger migration writes. Startup choice remains ActiveDesktopProfileID in the separate profile store.

## Observation

Registration key/label/endpoint, expected/observed daemon IDs, observation time and data or classified error identify each result. Four workers maximum, per-target deadlines and caller cancellation bound work. No credential values are returned. Calendar predictions and recorded executions are distinct; exact detail preserves available output/truncation/provenance.

## Owned fixture

UUID root, installation ID, explicit IPC, daemon PID and effective config/database/log/lock paths identify ownership. Validate canonical containment before mutations. Cleanup owns exact jobs/groups/processes, never executable names, and retains artifacts separately. Evidence identifies packaged_cli, command_harness or deterministic_engine; full report schema remains #290.
