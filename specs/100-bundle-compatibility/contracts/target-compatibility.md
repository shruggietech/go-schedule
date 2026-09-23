# Target Compatibility Contract

1. A current daemon's `GET /v1/manifest` response advertises `bundles` alongside existing capabilities.
2. A client operation snapshots the selected target once, obtains that target's manifest, and derives requirements from the supplied bundle's actual nonempty collections.
3. `bundle validate` returns a normal validation result with `valid: false` and deterministic findings when the selected target lacks a required capability or watcher platform support. It does not call `/v1/bundles/validate` on that target in this case.
4. Export, compare, preview, and apply reject incompatible targets before forwarding their respective bundle requests. A failed manifest discovery also blocks forwarding.
5. Compatible operations continue to the existing daemon endpoint, where schema validation, target-bound plan checks, and target-local watcher path rules remain authoritative.
6. No compatibility metadata is added to bundle JSON, and no target-only record is removed by comparison or apply.
