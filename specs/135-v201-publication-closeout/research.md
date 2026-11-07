# S135 Research

## Accepted source and release mechanism

**Decision**: Stage v2.0.1 from owner-accepted main `25f42d6eb35ff0cec2753f0f2934396ba50df954`, using the existing Release workflow only after explicit exact-source authorization. Use the existing standing-waiver promotion workflow after a separate decision on actual inspected assets.

**Rationale**: S134's accepted tree is identical to its reviewed candidate. Fresh exact-main CI has 20 successful jobs, CodeQL and Pages succeeded, and source/comparison checks pass. At the planning baseline, public latest was v2.0.0 and no v2.0.1 tag or release existed. Source PR authority does not authorize release mutation.

**Alternatives considered**: Another preparation-only slice would not deliver #419. Tagging the unaccepted S135 documentation branch would replace the approved release candidate with a different source. A new publication framework duplicates existing identity/byte guards.

## Current discovery surface

**Decision**: After actual publication, update README.md, site/data/published.json, docs/first-task.md, docs/gui-fields.md, docs/install.md, docs/INSTALL-windows.md, docs/website.md and test/windows/README.md. Update the existing site/e2e/site.spec.ts availability assertion.

**Rationale**: Read-only planning research identified each present-tense v2.0.0 or source-only notice. The README badge already follows latest; the fallback selects the existing reviewed v2.0.1 comparison without regenerating it. Historical v2.0.0 state-policy and release records remain correct.

**Alternatives considered**: Changing only the initial four entry points leaves contradictory current installation and desktop guidance. Global version replacement corrupts historical facts. Regenerating comparison/review data invalidates accepted content-bound records without functional need.

## Package readback and limits

**Decision**: Download the actual selected draft assets into a fresh S135 directory, bind IDs/sizes/digests to successful Release run/attempt and accepted source, inspect archive/build and MSI metadata, then independently download the public nine-asset set after approved promotion and verify all eight checksum entries.

**Rationale**: Source QA, a hosted build, published metadata and downloaded bytes prove different facts. Existing MSI inspection reads tables without installation. Attended native installation/GUI remains untested under the standing waiver.

**Alternatives considered**: Reusing historical packages cannot establish new byte identity. Replaying unchanged broad CLI QA adds no coverage for this documentation/discovery change; S134's qualified source remains separately attributable. Native installation is outside the standing waiver.

## Verification and preservation

**Decision**: Run all eight canonical gates in the foreground through new hidden S135 wrappers, plus the existing production site browser suite. Preserve S134's diagnostic and QA records. Treat site/data/published.json as an intentionally mutable live selector after publication, retaining its original baseline bytes/hash separately.

**Rationale**: Reusing S134 runner paths would overwrite receipts. Its protected inventory includes the live selector, whose exact intended transition must be disclosed rather than silently excluded or called unchanged. No scheduler/runtime changes require new tenancy or performance architecture.

**Alternatives considered**: Background test polling violates project instructions. Claiming every protected source path unchanged after changing the live fallback is false. Rewriting prior inventories loses historical evidence.

## Research completion

Installed speckit-plan dispatched read-only publication research. Findings were consolidated above; no unresolved technical choice remains. At planning, exact-source and actual-asset operation approvals were pending decisions, not unspecified product behavior. Subsequent approvals and execution are recorded in delivery.md.
