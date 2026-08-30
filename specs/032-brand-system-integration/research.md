# Research: Repository Brand System Integration

## R1. Canonical repository location

**Decision**: Import the complete approved kit at repository root as `brand/`.

**Rationale**: The kit serves documentation, GUI, packaging, release, and external consumers equally. A root directory makes that cross-cutting ownership visible and preserves the kit's existing internal structure.

**Alternatives considered**:

- `docs/assets/brand/` as the master would incorrectly make Jekyll's publish tree own platform files, build scripts, and UI references.
- `gui/assets/` as the master would repeat the earlier error of treating one consumer as the whole identity system.
- Importing only a few assets would leave the guide, tokens, platform files, licenses, and reproducible generation outside version control.

## R2. Consumer duplication policy

**Decision**: Keep exact consumer-local copies only where Jekyll or `go:embed` requires them, declare every copy in `brand/repository-consumers.json`, and require byte identity.

**Rationale**: Jekyll publishes only `docs/`, and Go embedding is simplest and most reviewable beside `gui/icon.go`. An explicit mapping converts unavoidable duplication into a verifiable contract.

**Alternatives considered**:

- Symlinks are unreliable on Windows checkouts and GitHub Pages.
- Rewriting consumers to reach outside their build contexts adds packaging fragility.
- Undeclared copies recreate the drift that this slice exists to remove.

## R3. Routine integrity tooling

**Decision**: Implement the repository verifier in Go using only the standard library and invoke it from the existing automation gate.

**Rationale**: Go is already guaranteed by local and hosted verification. It provides portable hashing, JSON, UTF-8, filesystem, and XML/token scanning without installing the optional graphics stack.

**Alternatives considered**:

- Running the kit's Python verifier would require CairoSVG, HarfBuzz, Pillow, pikepdf, and browser dependencies in routine CI.
- POSIX shell lacks a single portable SHA-256 command across Linux, macOS, and Git Bash.
- A separate ninth CI gate would expand process surface without improving feedback over the existing automation contract.

## R4. Canonical manifest and repository mapping

**Decision**: Preserve the kit's `manifest.json` as the record of the imported build and add repository-only mapping and guidance beside it without rewriting historical kit hashes.

**Rationale**: The completed kit has already passed standalone validation. The repository layer needs additional consumer relationships, not a mutation that makes the imported package differ from its approved archive.

**Alternatives considered**:

- Rebuilding the manifest after repository-specific additions would blur standalone kit evidence and repository integration evidence.
- Maintaining two complete manifests would invite disagreement about authority.

## R5. Platform release assets

**Decision**: Use the canonical Windows ICO and macOS ICNS directly during release builds; include the canonical `.desktop` file and hicolor tree in Linux desktop archives.

**Rationale**: The platform outputs are already generated, verified, and reviewed. Direct use removes release-time resampling, while Linux finally receives the desktop-integration files anticipated by issue #10.

**Alternatives considered**:

- Continue rebuilding `.icns` from the GUI PNG repeats generation logic in CI and can diverge from the approved platform output.
- Keeping Linux assets only in `brand/` would make them available but not useful to release consumers.
- Adding a system installer is outside this slice and disproportionate.

## R6. Public documentation depth

**Decision**: Add a concise Jekyll page that explains selection and usage, previews approved assets, and links to downloads and the full guide.

**Rationale**: Visitors need a practical web entry point, while the eleven-page guide remains the authoritative long-form reference. Duplicating the entire guide into Markdown would create two prose masters.

**Alternatives considered**:

- Publishing only the PDF is poorly discoverable and does not participate in site navigation/search.
- Reproducing every guide page in Markdown adds maintenance without additional user value.

## R7. Issue and scope boundary

**Decision**: S032 closes #10 and #34, and explicitly leaves #33 open.

**Rationale**: The canonical kit, generation pipeline, platform assets, wiring, and public page satisfy the two brand issues. Issue #33 requires environment observation the maintainer has already deferred, and asset integration is not evidence that the observed Windows behavior is fixed.

**Alternatives considered**:

- Closing #33 by inference would misrepresent unperformed Windows verification.
- Splitting #10 and #34 into separate thin slices would add more review ceremony than implementation value.
