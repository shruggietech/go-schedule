# S118 research

Date: 2026-10-25. All unknowns resolved from repository protocols, retained receipts and current GitHub records.

## Completed implementation versus external workflow

Decision: Mark S117 Implemented with actual PR #377 receipts. Move T016/T017 publication-only entries to a historical receipt table and explicitly attribute missed T018 source reconciliation to S118.

Rationale: docs/build-autopilot.md excludes publication from implementation checkboxes; specs/README.md separates maturity from publication. S117 source is already merged at 9ad13821c825d01efcc038b6b1634b187a854758.

Alternatives considered: Leaving stale tasks open obscures delivered work; marking every old task complete without attribution rewrites history. Neither is accurate.

## Existing pending release content

Decision: Preserve the selected v1.6.0 changelog section, release notes, comparison dataset/review and public v1.5.1 selection. Update current procedure links to S118 with a dated pinned-file decision.

Rationale: Existing three reviewed content hashes match actual bytes and already describe #375/#376 plus unresolved #355/#356. S114 procedure links are stale as current pointers, although their historical procedure remains useful.

Alternatives considered: Rewriting correct release prose or refreshing unchanged hashes adds no functional value. Reusing an old candidate as a future release identity is invalid.

## Issue-level outcomes and candidate identity

Decision: Keep #354 open for actual public release/discovery and #355/#356 open for their independently unconfirmed original causes. Preserve original failed reports and distinguish later corrective source acceptance, reviewed head and squash merge.

Rationale: Passing repetitions and verification are evidence, not causal fixes or issue closure gates. Accepted merged-main CI/deployment is independently recorded and does not establish a tag, draft or package.

Alternatives considered: Another unchanged full sweep or speculative runtime/IPC fixes would not resolve the demonstrated record gap.

## Delivery boundary

Decision: Use the existing exact-source tag, draft staging and separately approved promotion sequence only as future handoff guidance. Publish one source PR and handle at most two Codex rounds; stop unmerged.

Rationale: Latest owner authorization permits source publication and review, while withholding merge/tag/promotion. No staging identity exists to record.
