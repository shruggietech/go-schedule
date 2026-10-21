# S109 data model

## Invocation outcome

Existing errUsage and API validation code identify exit 2. Parser/operand wrappers preserve causes through errors.Is/errors.As. Runtime, transport and file I/O errors remain exit 1. No persistent state is added.

## Endpoint ownership

The invocation owns remote target fields. Global --daemon-endpoint and legacy --endpoint select the same remote endpoint with mutually exclusive explicit intent; the notification channel's local --endpoint belongs to its destination. Existing profile, daemon identity, certificate, credential and IPC fields remain intact.

## Binding identity

Each supplied portable watcher ID must occur exactly once in the selected document and at most once in CLI bindings. The shared preview service refuses unmatched IDs before saving a target-bound single-use plan.

## Issued credential

Retain the existing actor, daemon_id, credential and token shape. Real actor/daemon context is supplied from the authoritative enrollment/store state. No synthetic identity or zero timestamp is emitted as a successful association. Token remains one-time response data and never a log field.

## Watcher observation

Retain saved definition, health and readiness. Add observation/reconciliation metadata that binds an observation to the version of the definition actually processed, separately from health transition time. Immediate update replies can show pending while retaining the last observation; successful reconciliation of unchanged health must advance observation metadata. Runtime-unavailable and disabled states remain explicit. Existing field meanings and consumers remain compatible.
