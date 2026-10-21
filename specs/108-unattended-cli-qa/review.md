# S108 publication and bounded review record

## 2026-10-21: Authorized publication

Owner authorization covers pushing S108, PR #310, the 12 evaluated product defect issues, CI monitoring and up to two Codex review rounds, resolving findings autonomously and stopping before merge. Initial published source is 8522df934ac46e3b45067d78d98f66c82372ce17. No authentication was required.

Product findings F001 through F012 map in order to #311 through #322. Every issue has milestone 15, existing labels, parent #277, a bounded outcome, functional acceptance criteria, reproduction and source-pinned evidence. The full open/closed title inventory and targeted searches found no equivalent defects. Parent #277 maintains all child links and completion states.

## 2026-10-21: CI and Codex round 1

Round 1 reviewed 8522df934ac46e3b45067d78d98f66c82372ce17 and completed with four findings: collected WSL resume (#325), outcome exits (#326), mutable cleanup claims (#327), and dirty-source identity (#328). CI additionally identified Linux-only lint declarations/diagnostics (#323), unchecked redactor allocation arithmetic (#324), and a policy-harness timeout below the shared configured case budget (#329). These seven independently bounded corrections are implemented in PR #310, pending final verification and hosted rechecks.

Frozen capture metadata now binds cleanup and immutable case observations. Source identity includes exact tracked/nonignored untracked bytes, paths, modes and deletions. Resume accepts only explicit not-found/inactive or established stopped owned-unit properties; ambiguous state remains rejected. Outcome codes are 0 for successful processing/passing evaluation or complete pending-review capture, 1 for processing/integrity errors, 2 for failed evaluation and 3 for blocked evaluation. The PowerShell entry point preserves outcome codes and artifact locations.

Focused regressions pass for the new contracts. Current-source canonical eight-gate verification, native Linux checks, fresh actual capture-review-finalize and exact-head hosted results are being reconciled. A second Codex round has not yet been requested.

## 2026-10-21: Codex round 2

Round 2 reviewed f9c3b82200f019f18c327c89b0a8e4084a3f85d0 and completed with one finding, #330: source identity also needs a build-time association with prepared executable bytes. A seventh strict version 1 schema defines build-receipt.json, which binds the source digest, commit, dirty state and all six executable hashes. Preparation checks source stability across the build, and capture/resume reject mismatched or missing receipts. Dirty-source changes, dirty-to-clean pairing, executable substitution and incomplete receipts have focused regressions.

Native Linux-target lint also identified one remaining capitalized Unix diagnostic under #323; its wording is corrected without changing the control. Exactly two Codex review rounds were requested. No third round will be started; final corrections remain subject to focused regressions and exact-head hosted checks.

## 2026-10-21: Receipt-bound runtime checks

The pinned Linux-target linter reports zero issues, and current affected native Linux race checks pass. A public-entry-point mismatch probe added one exclusively created nonignored source file after preparation while keeping commit, dirty state and all six executables fixed. Capture blocked all 1,675 applicable cases before execution, retained the complete matrix and actionable subject_error, and left no owned resources remaining. The exact probe was removed. Public validate returned blocked outcome exit 3. The subsequent normal capture accepted source_sha256 c83829e61d742fbf91d59e908af258a7813cd68f1b2ee151b7cda2303c1ff6f2, matching build-receipt.json, and is running under .test-temp/s108-run-19 with .test-temp/s108-inputs-19.

The complete first review-fix sh scripts/verify.sh all aggregate exited zero across all eight gates. The final receipt-aware aggregate is .test-temp/verify-s108-review2.log and remains in progress while hosted checks run independently.
