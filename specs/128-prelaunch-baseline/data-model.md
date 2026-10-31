# S128 Current Data Contracts

## Application database

Current schema marker 24;22 tables and 33 authored indexes. Empty databases transition atomically to complete current state with one marker row. Complete existing marker 24 state reopens without conversion or marker rewriting. Unsupported/foreign/partial state is rejected with an explicit reset instruction, preserving its bytes/application state.

Current tasks, runs, identity, credentials, notifications, audit, triggers, watchers, chains and portable identifiers retain existing relationships and constraints. No ordered historical schema collection remains.

## Desktop preferences

Version1 contains appearance and complete popup preferences (enabled, conditions, severities, daemon IDs). Defaults follow System with popup delivery disabled. Existing current validation, atomic writes and lock scope remain. Transition, retired keys and legacy import paths are absent from the current model.

Unknown JSON fields are not interpreted. Missing required popup filters or invalid current values return an error; only explicit reset repairs malformed current documents.

## Helper metrics

One current complete format 4 contains meta, beat, snapshot, address and port tables with existing current constraints/indexes. Both writers may use an explicit common file. Fresh initialization is atomic; current format 4 reopens/appends, obsolete/foreign/partial layouts receive explicit fresh-path/reset guidance without conversion.

## Qualification

Current full evidence contract revision 2 describes current candidate fresh/current-state observations. Historical baseline role, upgrade session, retained-profile/setup.upgrade scenarios and old-package identities are removed. Required current identity, hash, environment, token, attachment, access, repair, rollback and uninstall fields stay authoritative. Historical revision 1 evidence is preserved as historical, never converted or relabeled.
