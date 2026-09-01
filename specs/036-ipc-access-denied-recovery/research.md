# Research: IPC Access-Denied Recovery

## Typed error boundary

**Decision**: Add a client `ConnectionError` with unavailable, access-denied,
timeout, and other-transport kinds; keep `StatusError` for API responses.

**Rationale**: `errors.Is` can classify wrapped OS errors without parsing user
copy, while callers retain `Unwrap` access to the cause.

**Alternatives considered**: GUI string matching and returning raw errors were
rejected because they caused fragile and misleading presentation.

## Incident presentation

**Decision**: Render one persistent panel above tabs with status, guidance,
Retry, and Exit.

**Rationale**: It preserves the application frame and provides a stable target
for updates without modal allocation.

**Alternatives considered**: A single modal still hides the app; transient
toasts cannot carry durable recovery actions.

## Retry ownership

**Decision**: One stream-owned reconnect loop with cancelable waits, 2-second
initial delay, doubling to a 30-second cap, and an immediate retry signal.

**Rationale**: Single ownership prevents duplicate loops and gives the goroutine
a clear shutdown path.

**Alternatives considered**: Per-request retries reproduce amplification;
unbounded backoff harms recovery; fixed polling creates avoidable churn.

## Windows evidence

**Decision**: Read service, local-group, account membership, and process-token
membership only for guidance; report unknown values explicitly.

**Rationale**: Issue #90 requires cause verification, while the security scope
forbids changing membership or ACLs.

**Alternatives considered**: Assuming every denial is stale membership is
inaccurate; automatically adding membership or elevating is out of scope.
