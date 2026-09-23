# Data Model: v1.5.0 publication identities

## Release source

- `tag`: exactly `v1.5.0`.
- `commit`: the immutable reviewed main revision with successful exact-commit CI.
- `metadata`: README health example, dated changelog section, and tag-specific notes matching the tag.

## Staged candidate

- `workflow run and attempt`: successful `Release` run for the tagged commit.
- `assets`: four headless archives, two desktop archives, Windows MSI, and Windows candidate manifest.
- `Windows identity`: manifest repository, tag, commit, run, attempt, MSI filename, byte size, digest, ProductVersion, and ProductCode.

## Publication disposition

- `mode`: full attended qualification, or the standing maintainer-authorized native-check waiver applied to this release.
- `observations`: genuine pass/fail/untested results, never inferred or backfilled from another candidate.
- `public record`: release notes plus issue comment explaining the chosen mode and known limitations.

## Public release

- `assets`: unchanged staged assets plus checksum inventory, and genuine evidence archive only in full-qualification mode.
- `state`: draft until all selected-path conditions pass, then public and latest.
- `traceability`: tag, workflow run, release URL, and issue #185.
