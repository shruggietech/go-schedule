# S108 publication and bounded review record

## 2026-10-21: Authorized publication

Owner authorization covers pushing S108, PR #310, the 12 evaluated product defect issues, CI monitoring and up to two Codex review rounds, resolving findings autonomously and stopping before merge. Initial published source is 8522df934ac46e3b45067d78d98f66c82372ce17. No authentication was required.

Product findings F001 through F012 map in order to #311 through #322. Every issue has milestone 15, existing labels, parent #277, a bounded outcome, functional acceptance criteria, reproduction and source-pinned evidence. The full open/closed title inventory and targeted searches found no equivalent defects. Parent #277 maintains all child links and completion states.

## 2026-10-21: CI and Codex round 1

Round 1 reviewed 8522df934ac46e3b45067d78d98f66c82372ce17 and completed with four findings: collected WSL resume (#325), outcome exits (#326), mutable cleanup claims (#327), and dirty-source identity (#328). CI additionally identified Linux-only lint declarations/diagnostics (#323), unchecked redactor allocation arithmetic (#324), and a policy-harness timeout below the shared configured case budget (#329). These seven independently bounded corrections are implemented in PR #310, pending final verification and hosted rechecks.

Frozen capture metadata now binds cleanup and immutable case observations. Source identity includes exact tracked/nonignored untracked bytes, paths, modes and deletions. Resume accepts only explicit not-found/inactive or established stopped owned-unit properties; ambiguous state remains rejected. Outcome codes are 0 for successful processing/passing evaluation or complete pending-review capture, 1 for processing/integrity errors, 2 for failed evaluation and 3 for blocked evaluation. The PowerShell entry point preserves outcome codes and artifact locations.

Focused regressions pass for the new contracts. Current-source canonical eight-gate verification, native Linux checks, fresh actual capture-review-finalize and exact-head hosted results are being reconciled. A second Codex round has not yet been requested.
