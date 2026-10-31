# S127 Research and Decisions

**Date**: 2026-10-31

## Direction

- Decision: Record explicit owner-directed deferral.
- Rationale: The S127 clarification answer was "Explicitly defer clustered execution." It resolves the choice S125 could not make from missing real-use evidence.
- Alternatives considered: immediate architecture without demand findings would bypass existing gates; permanent no-go would contradict the selected answer; repeating S125's same evidence search would not deliver the new decision.

## Issue and capability boundaries

- Decision: Preserve #186-#189, #20 and #146 open and distinguish deferred planning from shipped capability.
- Rationale: Deferral supplies direction, not actual workload findings, semantics, prototype results or final go/no-go.
- Alternatives considered: closing as completed or not planned would overstate the owner's choice.
- Sources: existing native issue requirements; [S125 decision](../125-reliability-delivery-decision/decision.md); [current page](../../docs/distributed-scheduling.md).

## Reactivation

- Decision: A later owner revisit or explicit gate revision is necessary. Reuse the existing intake and dependency sequence.
- Rationale: New evidence may justify a revisit but is not authorization to resume a deferred effort.
- Alternatives considered: date-based or automatic resumption adds unrequested work; permanent product limits narrow the product without authority.

## Current release attribution

- Decision: Refer to public v1.7.0 and its immutable comparison source without modifying release notes, assets, reviewed comparison data or historical source claims.
- Rationale: PR #406 merged at 6483fc81bbaa75ec0bddeb0fdf17fc4cc66be426; release 428427523 is public from df54fc7b2d49d3b71125a049c42bfd4225503b1d.
- No new external product recommendation or dependency selection is needed.
