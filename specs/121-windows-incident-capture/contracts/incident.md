# Incident contract v1

Incident collection is opt-in. capture/resume --incidents records the selection; resume cannot change it. Script experiments use GOSCHEDULE_INCIDENT_DIR as an existing private absolute parent. Every child gets a unique subdirectory and fixed expected artifact names. No existing incident is overwritten and no linked/escaped path is read or written.

A collector never changes original native exit, partial-state completeness or cleanup results. A collection failure remains separate. Runtime selection and same-process pre-payload identity are distinct from crash-time module evidence. Collection status is collected, absent, unavailable, failed, truncated, oversized timed_out or canceled. Retained artifacts include role, byte count and SHA-256; unknown or oversize files do not enter report publication.

Private daemon requests are <=1KiB and bind nonce/ID/PID. Responses use exclusive/atomic finalization and stay <=128KiB. At most eight per daemon are served; active-handler entries cap 64 with omitted counts; stack buffer 64 KiB. Client collection allowance is 2 seconds and does not replace the caller budget. Paths, raw stacks and native dumps stay private; structured report references go through existing redaction/artifact controls.
