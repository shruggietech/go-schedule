# Connection Contract

## Trust Boundary

The production Wails bridge exposes sanitized values only. Endpoint paths, pipe names, socket names, OS error text, raw API errors, credentials, HTTP objects, and backend implementations remain in Go.

## Bound Methods

### `ConnectionSnapshot() ConnectionSnapshot`

Returns the latest immutable snapshot immediately. Before startup negotiation completes it returns generation zero in `connecting` state with local target identity.

### `RetryConnection() ActionResult`

Requests one fresh generation. It cancels an active attempt, stream, or retry wait, coalesces concurrent retry requests, and returns a safe accepted or unavailable result. It never blocks for connection completion.

### `Quit() ActionResult`

Requests orderly application shutdown through the native runtime. Browser adapters return a safe unavailable result without closing the browser.

## Event Channel

One Wails runtime channel named `desktop:event` carries `DesktopEvent` values. `connection.changed` instructs the frontend store to replace its snapshot only when its generation is newer or its revision is at least the current revision within the same generation. Sanitized daemon-domain events notify later workflows to refresh through their connection-facing application methods.

## Compatibility

Development versions and daemon major version 1 are compatible with the S061 local adapter. Another explicit major version maps to `incompatible`, exposes `update_daemon`, and does not open the event stream.

## Safe Error Mapping

| Source | State | Title | Action | Automatic retry |
| --- | --- | --- | --- | --- |
| endpoint absent or refused | unavailable | Daemon unavailable | open_service_help | yes |
| OS permission rejection | access_denied | Access denied | refresh_login | no |
| request deadline | timed_out | Connection timed out | retry | yes |
| compatible health, event stream failure | degraded | Live updates delayed | retry | yes |
| incompatible explicit major | incompatible | Update required | update_daemon | no |
| other bounded transport failure | unavailable | Connection interrupted | retry | yes |

No mapping includes the wrapped cause.

## Lifecycle

Startup creates one root context and one manager loop. Every request, delay, and event stream is a child of the active generation context. Manual retry cancels that child and signals the loop. Shutdown cancels the root, prevents new generations, waits for the loop, and completes within two seconds.
