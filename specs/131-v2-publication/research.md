# S131 Release Research

## Decision: reuse accepted delivery machinery

Use accepted main 75eae1ae11ad40db38341dc7d63f32876cba6483 and existing Release/standing-waiver workflows. The source already has successful main CI and corrected SQLite graphs. Rebuilding accepted assets after inspection or creating another release-preparation framework would add inconsistent identities without a missing capability.

## Rationale and alternatives

release.yml requires successful exact-main CI, valid source/comparison metadata and absent/draft state before staging. promote-release-waived.yml checks immutable tag/source, complete eight-asset allowlist, Windows candidate manifest and successful staging attempt, checksum generation/readback and frozen release selection. Direct mutation and a new framework are rejected.

## Existing authorization

Source autopilot/push/official PR and in-scope feedback corrections are approved. Exact v2.0.0 tag/draft and separate actual-asset promotion remain pending decisions. Earlier v1.7.0 approvals do not apply. Current-only state and owner-deferred clustering remain unchanged.
