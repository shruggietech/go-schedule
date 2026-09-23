# Contract: v1.5.0 publication

1. The tag points to one reviewed main commit with successful CI. Once staged, it is not moved or silently rebuilt.
2. The draft contains exactly the documented platform artifact set and Windows candidate manifest before any qualification-specific asset is added.
3. The candidate verifier binds the MSI and manifest to repository, tag, commit, workflow run, byte count, digest, and product identity.
4. Full qualification requires real passing attended evidence. The standing waiver path requires disclosure of this release's untested checks and the same exact-source and artifact-integrity controls; it never feeds synthetic evidence to the full validator.
5. The final checksum inventory covers every public binary, manifest, and genuine evidence asset. Freshly downloaded public bytes must match it.
6. The release is not marked latest, and #185 is not closed, until all selected-path checks and public asset publication complete.
