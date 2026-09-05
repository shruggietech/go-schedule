# Data Model: GitHub Security Baseline

This feature adds no application persistence. Its entities are repository and hosted-security contracts.

## Reporting Route

| Field | Type | Rule |
| --- | --- | --- |
| primary route | URL | Must target this repository's private advisory form. |
| private reporting | state | Must be enabled before the feature is active. |
| fallback | email | Must be `info@shruggie.tech`. |
| triage owner | actor group | Repository administrators. |
| acknowledgement target | duration | One week. |
| public disclosure warning | boolean | Must remain true. |

## Security Control

| Field | Type | Rule |
| --- | --- | --- |
| name | enum | Private reporting, secret scanning, push protection, non-provider patterns, validity checks, CodeQL. |
| requested state | enum | `enabled`. |
| observed state | enum | `enabled`, `unavailable`, `unverified`, or `failed`. |
| evidence source | string | API field, workflow run/check, or route inspection. |
| limitation | string/null | Exact plan, permission, or API reason when not enabled. |
| compensating control | string/null | Required when unavailable. |

### State transitions

`unverified -> enabled | unavailable | failed`

An `unavailable` result is not green. A request accepted without a successful read-back remains `unverified`.

## Security Analysis Contract

| Field | Required value |
| --- | --- |
| language | `go` |
| push branch | `main` |
| pull-request target | `main` |
| schedule | At least weekly |
| source permission | `contents: read` |
| result permission | `security-events: write` |
| build mode | `manual` |
| build command | `CGO_ENABLED=0 go build ./...` |
| action runtime line | CodeQL v4; checkout/setup-go v7 |
| secrets | None |

## Security Evidence Record

| Field | Type | Sensitivity |
| --- | --- | --- |
| collected at | timestamp | Non-sensitive |
| control | string | Non-sensitive |
| status | observed-state enum | Non-sensitive |
| source | API/check/route | Non-sensitive |
| alert count | integer/null | Non-sensitive aggregate only |
| limitation | string/null | Must not include alert or vulnerability contents |

## Pinned-Artifact Decision

| Field | Rule |
| --- | --- |
| date | `2026-08-26` for this slice |
| artifact | Exact workflow or policy path |
| decision | Security control added or policy extended |
| rationale | Why advanced setup/manual build/offline validation was chosen |
| compatibility | Node 24 majors, Go module version, cgo-free build, unchanged CI |
