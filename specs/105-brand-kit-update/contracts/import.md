# Brand import contract

The maintainer command accepts an archive and the formal release `SHA256SUMS` file. It reads the explicit repository pin and fails if the exact archive filename, digest, bundle identity, or manifest inventory differs. It refuses absolute paths, parent traversal, duplicate names, and unexpected files before writing to the canonical kit. Failed preflight leaves the installed kit unchanged.

On success, it stores the complete ZIP unchanged, synchronizes the declared consumer paths from archive entries, and removes targets retired from the previous mapping and unpacked legacy kit. It preserves repository-owned control files and documented packaging metadata. The command then runs the same offline checks used in CI.

`go run ./scripts/brand-check` verifies the installed archive manifest, release pin and bundle identity, repository-specific allowances, UTF-8 and SVG integrity, and byte equality of every mapped consumer. Any mismatch exits nonzero with the affected path.
