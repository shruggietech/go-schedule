# Data Model: Target-Aware Bundle Compatibility

## Target profile

The daemon manifest supplies installation identity, operating system, and a set of advertised capability names. The client snapshots the concrete selected daemon before reading this profile, and the profile applies only to that operation.

## Derived requirements

Requirements are a deterministic set derived from bundle content: `bundles` always; `groups`, `chains`, `triggers`, `watchers`, and `notifications` for nonempty corresponding collections; and both `schedule` and `tasks` for task definitions. Watcher content additionally requires a supported target operating system. No requirement is stored in or trusted from the bundle document.

## Compatibility finding

One finding names the missing capability or unsupported watcher platform, identifies its affected bundle family, and tells the operator to select or upgrade a capable daemon. Findings are stable and unique. They are returned by validation and transformed into an actionable error for operations that cannot proceed.

## Existing entities

The v1/v2 bundle document, canonical digest, target-bound plan, and item outcomes retain their existing structure. Watcher path bindings remain out-of-document request data. No persistent schema changes occur.
