# Data Model: v1.3 Release Qualification Evidence

S072 adds no persisted production entity or migration. The following verification-only concepts define the evidence model.

## Qualification Matrix

Represents the complete release gate.

| Field | Meaning | Validation |
|---|---|---|
| operating_system | Hosted execution platform | Exactly Windows, macOS, or Linux |
| candidate_commit | Pull-request commit under test | Non-empty immutable commit identifier |
| default_state | Fresh and retained state result | Pass or fail, never silently skipped |
| notification_suite | Webhook, migration, redaction, recovery, and isolation result | Pass or fail |
| mcp_suite | SDK, protocol, transport, authorization, redaction, and lifecycle result | Pass or fail |
| overall_state | Platform qualification result | Pass only when every required component passes |

## Package-Shaped Daemon State

Represents one isolated candidate execution.

| Field | Meaning | Validation |
|---|---|---|
| executable | Built `goschedd` path | Exists below the isolated candidate directory |
| config_path | Full daemon configuration | UTF-8 JSON, isolated paths, no optional surface enabled |
| data_directory | Retained SQLite and runtime state | Unique to the test and reused for the second start |
| ipc_endpoint | Protected local client endpoint | Unique per test; named pipe on Windows, Unix socket elsewhere |
| start_number | Fresh or retained launch | `1` or `2` |
| channel_count | Configured notification channels | `0` before explicit configuration |
| delivery_count | Notification delivery records | `0` before explicit configuration |
| mcp_http_enabled | Localhost HTTP listener state | `false` before explicit enablement |
| exposed_mcp_metadata | Endpoint, fingerprint, client identity, and access evidence | Empty while disabled |

## Evidence Relationships

- One qualification matrix contains exactly three platform results.
- Each platform result contains two launches of one package-shaped daemon state.
- The second launch reuses the first launch's configuration and data directory.
- Detailed notification and MCP suite results are linked to the same candidate commit.
- The release boundary is qualified only when all platform and canonical verification evidence passes.

## State Transitions

```text
candidate built -> fresh start healthy -> defaults inspected -> daemon stopped -> retained start healthy -> defaults inspected -> platform qualified
```

Any build, startup, inspection, shutdown, or focused-suite failure moves that platform result directly to failed. There is no skipped-success transition.
