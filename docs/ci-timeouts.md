---
title: Verification completion timeouts
nav_order: 8.1
---

# Verification completion timeouts

Owner authorization dated 2026-10-25 supersedes earlier advice against increasing incidental verification deadlines. These are completion allowances, not performance targets or guarantees that an earlier failure has a proven cause.

| Boundary | Previous | New | Classification |
| --- | --- | --- | --- |
| Workflow jobs (22 definitions) | inherited 360 min | explicit 360 min | hosting ceiling, no effective increase beyond maximum |
| Go package tests / coverage / benchmarks | 10 min default | 60 min | canonical and all direct hosted commands |
| Frontend test / hook / teardown | 5 s / 10 s / 10 s | 60 s each | full asynchronous scenarios |
| DOM async lookup | 1 s default; selected case 5 s | 10 s shared | inner lookup below whole-test budget |
| Browser tests / assertions | 30 s / 5 s | 120 s / 30 s | both browser suites |
| Desktop browser server / site readiness | 120 s / 5 s | 300 s / 60 s | incidental startup |
| Static page generation | 60 s | 300 s | Next build |
| Packaged build preparation | 2 min | 10 min | immutable preparation only |
| Packaged runtime / SQLite test version probe | 20 s / 3 s | 120 s / 30 s | normal fixture work |
| QA case / execution total | 15 s / 30 min | 120 s / 2 h | within current supported maxima; explicit/resume preserved |
| QA setup / worker evidence freeze | 30 s / 30 s | 120 s / 120 s | preparation and evidence completion |
| QA shutdown / reporting / wrapper drainage | 45 s / 60 s / 15 s | 180 s / 240 s / 60 s | outer selected total plus 480 s |
| QA bootstrap version / build / timing query | 5 s / 300 s / 10 s | 30 s / 900 s / 60 s | fallback 8,000,000 ms |
| Owned QA protocol / child exit / webhook close / ownership probe | 3 s / 3 s / 2 s / 10 s | 15 s / 30 s / 20 s / 60 s | not production transport defaults |
| Service CLI / health / task completion | 10 s / 30 s / 45 s | 120 s / 180 s / 180 s | disposable CI only |
| Installer task completion | 45 s | 180 s | MSI wait remains outer-job bounded |
| Pages action deployment | 10 min | 60 min | milliseconds action input |
| Release exact-main waiter | 30 min | 300 min | 15 s poll cadence, 1200 attempts |
| Promotion request / command | 5 min / 20 min | 15 min / 120 min | immutable identity/refusal preserved |
| Qualification guest default / checks / init / WebView / drainage | 1200 / 30 / 120 / 600 / 3 s | 3600 / 120 / 600 / 1800 / 30 s | no attended guest run on owner host |

## Deliberate contracts retained

- Intentional injected QA wrapper 150 ms expiry and explicit short budgets
- SQLite adapter 10-second candidate-probe timeout and inherited-stream assertions
- Notification transport 50 ms cancellation and negative group observation 150ms
- Qualification self-test 3s timeout against 6s child
- Production reconnect 6500ms, RPC/dial/TLS/SMTP bounds and scheduler dispatch latency
- Unlimited lint, Node site whole-test and MSI Start-Process waits remain outer-job bounded

## Evidence and references

Prior failure receipts remain immutable and source-bound. The new default budget does not replace explicit or resumed budgets. Resource, ownership, row counts, identity and exit/output assertions remain active. Original CLR/IPC causes remain unconfirmed unless independently demonstrated.

- [GitHub Actions hosted limits](https://docs.github.com/en/actions/reference/limits)
- [Go timeout default](https://go.dev/src/cmd/go/testdata/script/test_timeout.txt)
- [Vitest options](https://main.vitest.dev/guide/cli.html)
- [Testing Library configuration](https://testing-library.com/docs/dom-testing-library/api-configuration/#asyncutiltimeout)
- [Playwright nested timeouts](https://playwright.dev/docs/test-timeouts)
- [Pages deployment inputs](https://github.com/actions/deploy-pages/blob/v5/action.yml)
- [Next static generation timeout](https://nextjs.org/docs/messages/static-page-generation-timeout)
