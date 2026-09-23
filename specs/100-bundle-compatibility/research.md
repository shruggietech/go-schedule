# Research: Target-Aware Portable Bundle Compatibility

## Current boundary

The v1/v2 bundle validator checks syntax and several target-domain constraints. The manifest advertises feature families but not the bundle API itself. The shared client currently calls bundle endpoints directly and its switchable facade can resolve a different selected daemon on consecutive requests. The target daemon already interprets watcher path bindings using its local `filepath` rules and owns the single-use plan identity and fingerprint.

## Selected approach

Add `bundles` to the current daemon manifest. Derive required capability names from nonempty bundle families, with a fixed deterministic order. The shared client snapshots `target()` once, reads that target's manifest, evaluates requirements, and only then forwards the operation to that same target. Validation returns local compatibility findings when the endpoint is unavailable. Existing server-side schema and plan validation remain authoritative after preflight.

## Alternatives considered

- A bundle-embedded requirement list is rejectable as untrusted and would change v1/v2 canonical bytes and digests.
- Calling the bundle endpoint and interpreting 404 is less actionable and cannot list missing source-family capabilities before preview.
- Using the caller's operating system to validate watcher paths would reject valid remote paths and accept invalid target paths.
- Caching the manifest across operations risks stale capabilities and target changes; one discovery call is acceptable on this cold path.

## Compatibility details

`bundles` is always required. Groups, chains, triggers, watchers, and notifications are required only when the document contains their corresponding objects; tasks require both `schedule` and `tasks`. A watcher also requires a recognized target OS. Unknown capabilities do not stand in for known names. No bundle schema migration or store migration is necessary; v1 documents remain readable and v2 remains the default export.
