# Data Model

- `Document` v2: groups, tasks, chains, external triggers, trigger sets, watchers, notification policies, and exclusions.
- `ExternalTrigger`: portable identity, name, target task identity, enabled intent. No key.
- `TriggerSet`: portable identity, name, target task identity, member count, enabled intent. No member keys.
- `Watcher`: portable identity, name, kind, pattern, recursion, timing, target task identity. No path.
- `NotificationPolicy`: scope kind and portable identity, channel name and outcome conditions. No endpoint or authorization.
- `BundleRequest`: document plus optional target-local watcher path bindings for preview.
- `storedBundlePlan`: reviewed document, plan, and path bindings. Apply receives only the single-use plan token.
