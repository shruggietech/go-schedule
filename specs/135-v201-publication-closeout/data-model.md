# S135 Data Model

## Release candidate

Version v2.0.1, source commit 25f42d6eb35ff0cec2753f0f2934396ba50df954, immutable tag, successful Release run ID and attempt. Lifecycle: accepted source, approved staging, inspected draft, approved promotion, public release. A changed source or asset selection invalidates the applicable downstream selection.

## Asset selection

Release ID plus asset ID, exact name, size, SHA256 and source/run/attempt association. Eight original content assets precede the checksum inventory; nine final public assets must agree with the inspected selection and complete checksum list. Original content bytes are immutable through promotion.

## Discovery selection

Current published fallback tag selects an existing reviewed comparison snapshot. Live public/latest selection and the later source-owned fallback/guide revision are separately recorded. Historical comparison/review records remain unchanged.

## Evidence and workflow record

Observation timestamp, exact source or asset identity, outcome and limitation. Original failures remain diagnostic. Review ledger counts automatic request and any manual request with a maximum of two. Source integration and Pages deployment remain distinct from public package delivery.
