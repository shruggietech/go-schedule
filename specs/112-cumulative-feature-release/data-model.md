# S112 data model

No database schema changes. The v23 notification migration is an existing S111 prerequisite.

## Probe evidence

The evidence record has status, primary failure, owned task/run observations, health/service observations and cleanup outcomes. Missing diagnostics are explicit errors. Only bounded fixture data is retained; no credential or environment enumeration is allowed.

## Release state

Source preparation, reviewed source merge, tag/draft staging and public promotion remain distinct states. The pending v1.6.0 comparison cannot replace the published v1.5.1 comparison until actual public promotion.
