# Contract: Remote API v1 and Enrollment

## Transport

- HTTPS only, TLS 1.3 minimum.
- URI major `/api/v1`.
- Protected requests use exactly one `Authorization: Bearer <opaque>` header.
- Requests with `Origin`, credential query parameters, or credential cookies are rejected.
- JSON request bodies require `application/json` and are limited to 1 MiB unless an operation declares a smaller limit.

## Public discovery and enrollment

| Method | Path | Authentication | Retry |
| --- | --- | --- | --- |
| GET | `/api/v1/health` | None | Safe read |
| GET | `/api/v1/manifest` | None | Safe read |
| POST | `/api/v1/enroll` | One active pairing phrase plus expected daemon identity | No replay |

Health and manifest contain no runtime paths or secrets. Enrollment returns a raw bearer value once and never returns its phrase, salt, verifier, or digest.

## Protected initial allowlist

The runtime allowlist exposes the safe task, group, run, alert, calendar, event, actor, and audit operations explicitly enumerated in `api/openapi/remote-v1.yaml`. Each route maps to the corresponding `/v1` application operation and S076 operation ID. Secret reveal or rotation routes, runtime information, MCP administration, trigger administration, filesystem watcher administration, notification-channel administration, and undeclared routes are absent.

## Stable errors

| Status | Code | Meaning |
| --- | --- | --- |
| 400 | `validation_failed` | Bounded syntax or field failure |
| 401 | `unauthorized` | Generic credential or enrollment failure |
| 403 | `forbidden` | Authenticated actor lacks authority |
| 404 | `not_found` | Resource or undeclared path is absent |
| 409 | `conflict` | Domain or compatibility conflict |
| 413 | `request_too_large` | Declared body bound exceeded |
| 415 | `unsupported_media_type` | JSON operation received another media type |
| 429 | `rate_limited` | Source or actor bucket is exhausted |
| 503 | `audit_unavailable` | Required durable evidence could not be written |

Authentication failures include `WWW-Authenticate: Bearer` and never distinguish missing, malformed, unknown, expired, revoked, or rotated credentials.

## Local administration

Local Enroll authority adds pairing create, list, and cancel operations plus credential list, rotate, and revoke operations under `/v1/access`. Raw pairing phrases appear only in create responses; raw credential values appear only in enrollment and rotation responses.

## Completeness

Automated checks compare the OpenAPI operation IDs to the runtime remote table, assert unique route and operation identities, and prove representative denied local-only paths return 404 before local application dispatch.
