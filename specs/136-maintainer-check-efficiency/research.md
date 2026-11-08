# Research: Maintainer check efficiency

## 2026-11-08 baseline and contract research

Unchanged accepted main 5e4f4f4ca373ca31b607f7d456fb4d3d3339a638 completed lifecycle in 31.541 seconds and automation in 44.701 seconds on this Windows host. Both exited zero. Automation includes lifecycle, so timings overlap. S135 aggregate 4355.058 seconds is historical context, not an isolated attribution. Private original commands/results are retained.

Decision: use existing Go/YAML tooling with an evaluation-local read cache and in-process lifecycle traversal. Alternatives: retain shell utility loops (measured startup cost), embed a shell interpreter (new dependency and copied fragile logic), or skip policy work (violates scope). Existing internal/platform supports hidden Windows children.

Read-only research identified comment/run-string action false matches (#434), invalid colon-space plain scalars in synthetic workflows, and a desktop-job removal fixture that also deleted the next header. Correct the latter fixture syntax while preserving the intended cases. Publication order must be evaluated within decoded run blocks and across steps; duplicate mappings must fail. No unresolved research question remains.

The original publication checker also returned success when promotion comparison moved to a separate job without a needs edge. #435 records that demonstrated defect. The native evaluator follows actual needs edges and command order. [GitHub workflow reuse documentation](https://docs.github.com/en/actions/reference/workflows-and-actions/reusing-workflow-configurations) confirms support for anchors/aliases; anchored fields are evaluated and cyclic alias graphs fail explicitly.

Expanded pending-source probes demonstrated #436: substring checks accepted near-match dependency values. Exact decoded membership and unique Makefile assignments retain the prior policy. Native focused regressions and three additional Makefile cases cover the correction. The final full aggregate passes with all 96 cases.
