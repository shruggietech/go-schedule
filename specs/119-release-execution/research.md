# S119 research

Date: 2026-10-26. Read-only plan-command research completed for existing staging and post-publication discovery. No demonstrated staging source blocker.

## Existing immutable staging

Decision: Use the existing tag-only Release workflow on owner-approved 0f8834ba3b88b69c8f55aaa5d72add70e5b6252b. Retain four daemon/CLI archives, two desktop archives, Windows MSI and candidate manifest, plus actual release/asset/run/attempt identities.

Rationale: Existing exact-main CI, source/comparison and absent-or-draft gates already express the required contract. SHA256SUMS.txt is added during promotion, not expected among eight draft assets.

Alternatives considered: Moving the tag, local packaging or substitute workflows would discard approved provenance. No new framework is warranted.

## Separate promotion authority

Decision: Request approval only after all actual staged identities and successful Windows manifest attempt are available. Reuse promote-release-waived.yml from reviewed main, with authenticated exact-release selection, original assets, complete checksums and immutable-tag checks.

Rationale: Owner approved draft staging, not public promotion. Existing standing waiver already covers native observations and requires truthful disclosure; no host installation/session change is needed.

Alternatives considered: Premature approval of nonexistent assets, full-attended gate under a waiver and guard bypass are inconsistent with the established route.

## Source discovery reconciliation

Decision: After actual promotion, update the existing published selector and current canonical status paragraphs, preserving immutable tagged notes/dataset/original review bytes. Correct the current main v1.6.0 introductory status and refresh only its reviewed changelog digest/date after inspecting changed prose, while retaining the dated heading/anchor and recording both original and reconciled hashes.

Rationale: Pages already follows successful promotion and selects actual releases/latest; generated site content comes from canonical docs and must not be hand-edited. Current pinned Windows installation status needs a dated changelog decision. Existing rollback and #355/#356 disclosures remain applicable.

Alternatives considered: Rewriting immutable tag/notes/dataset would misrepresent published provenance. Leaving the main introduction claiming public v1.5.1 after publication could mislead direct-anchor readers. Existing source/comparison validators permit a reviewed main-only status/digest refresh; broad website/workflow redesign has no demonstrated need.

## Release-aware browser assertions

Decision: Exercise the existing responsive browser case against the prepared v1.6.0 snapshot before changing its hardcoded v1.5.1 expectation. If reproduced, track its maintainer-QA defect separately and bind the assertion to the actual prepared comparison metadata while retaining exact visibility/column assertions.

Rationale: Current site/e2e/site.spec.ts hardcodes the previous public version although prepare.mjs supports a published selector and SITE_RELEASE_TAG override. This is directly compatible with publication/discovery reconciliation.

Alternatives considered: Hardcoding v1.6.0 only repeats the next-release defect; removing the assertion weakens the contract. No new test framework is needed.
