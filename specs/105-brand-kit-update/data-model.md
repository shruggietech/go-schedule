# Data model: brand release import

## Pinned release

- `package`: exact archive filename and package identity.
- `release_url`: formal release asset URL.
- `sha256`: lowercase 64-character archive digest.
- `brand_version` and `brandbuilder_version`: independent version identities.
- State: proposed pin -> archive and checksum verified -> installed canonical kit.

## Canonical kit

- `manifest.json` inside the archive lists every distributable file by relative path, byte count, and SHA-256.
- `enforcement/bundle.json` identifies the package, version, and formal publication tag.
- The complete ZIP is stored unchanged in `brand/`; mapped presentation files are repository consumers rather than an extracted second copy of the kit.

## Consumer mapping

- Each mapping has one canonical source path, one or more repository-relative target paths, and a purpose.
- Source and target paths are normalized, unique, and constrained to the checkout.
- A target may be a repository-specific `brand/platform/` packaging path; its bytes must still equal its mapped official source unless documented as a local derivative.
