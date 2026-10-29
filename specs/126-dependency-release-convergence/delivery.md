# S126 Delivery

## Integrated source and published release

The owner merged [PR #405](https://github.com/shruggietech/go-schedule/pull/405) at df54fc7b2d49d3b71125a049c42bfd4225503b1d. Exact-main CI, CodeQL and website checks passed. Separately approved staging and concrete-asset promotion published [v1.7.0](https://github.com/shruggietech/go-schedule/releases/tag/v1.7.0), release 428427523, from that immutable source. Publication-triggered Pages now selects v1.7.0 over HTTPS.

| Record | Functional outcome | Disposition |
| --- | --- | --- |
| PR #401 | API runtime 1.8.0 and YAML 3.0.5, complete intended Go graphs | Integrated through #405; redundant proposal closed by Dependabot |
| PR #402 | Vite 8.3.5 and selected lock graph | Integrated through #405; redundant proposal closed by Dependabot |
| PR #403 | Coupled Next/Fumadocs update and complete site lock graph | Integrated through #405; redundant proposal closed by Dependabot |
| PR #400 / #404 | Suitable published SQLite update and complete secondary graph | Deferred, no suitable corrected driver published at selection; remain open |
| #394 | v1.7.0 reliability packages and actual public discovery | Release is public/latest; this follow-through reconciles the remaining source fallback and current prose before issue closure |
| #355 / #356 | Original Windows historical causes | Unresolved; no causal closure inferred from current success |
| #146 / #20 / #186-#189 | Remaining distributed-scheduling decision | Separate conditional work; no clustered capability or final disposition claimed |

## Actual assets and publication

Release run 40960084352, attempt 1, passed all seven staging jobs. The owner approved the [concrete eight-asset inventory](https://github.com/shruggietech/go-schedule/issues/394#issuecomment-6509164486). Promotion run 40961590136, attempt 1, revalidated the same source, staging run and asset identities, added SHA256SUMS.txt and published the same release. All nine public downloads match their GitHub digests; all eight original asset IDs and hashes are unchanged. The checksum file has exactly eight matching entries. Independent MSI inspection confirms ProductVersion 1.7.0 and ProductCode {46A1517E-8ECD-4F7F-8D5E-F77D2AC2BAD5}, matching its manifest.

[handoff.md](handoff.md) records the immutable source, tag object, workflow/deployment IDs, recovery contract and exact artifact evidence. Earlier application QA remains distinct from staged package inspection. Attended native observations remain **not tested, release authorized by maintainer** under the standing waiver.

## Review and source reconciliation boundary

At owner merge, both authorized Codex requests lacked a terminal completion receipt. The summary still reported the final manual request running; no two-clean-round claim is made. The owner deleted the scheduled monitor. No third review was requested.

The published release is complete. This same-slice source follow-through updates the repository fallback, current README/installation wording and factual delivery records without altering the immutable release, its notes or its content-bound comparison. It uses a new review branch and stops before pushing/opening another PR pending publication authorization. [verification.md](verification.md) separates source validation for this follow-through from earlier QA and actual release acceptance.
