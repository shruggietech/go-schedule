# Repository brand integration

The complete official kit is the immutable archive named in [`source.json`](source.json). Its `manifest.json` lists every official file and digest. The selected files outside the archive are installed according to [`repository-consumers.json`](repository-consumers.json); edit the map when a consumer path changes instead of editing a copied asset.

## Update from the official release

1. Find the current go-schedule archive on [brand.shruggie.tech](https://brand.shruggie.tech) and its formal `shruggie-brand` GitHub release. Download the exact ZIP and that release's `SHA256SUMS` file. Confirm the archive name, brand version, BrandBuilder version, and checksum, then update `source.json`.
2. Review the archive's migration notes and manifest. Update `repository-consumers.json` with each approved source path and repository target. Keep repository-owned Linux `.desktop` metadata separate; use the official size-matched raster exports for Linux icons when the release provides no Linux icon suite.
3. From the repository root, run `go run ./scripts/brand-import <archive.zip> <SHA256SUMS> <new-consumer-map.json>`. The importer validates the pin, release checksum, ZIP entries, manifest, bundle identity, source paths, and destination paths before writing. It then stores the unchanged archive, synchronizes declared consumers, and removes retired declared assets.
4. Update documentation, packaging references, and tests for any renamed assets. Run `go run ./scripts/brand-check`, `go run ./scripts/github-format`, and `sh scripts/verify.sh all` before publishing the branch.

The importer accepts only the version 2 consumer map. A mapping with `transform: "docs-site-manifest"` rewrites absolute favicon URLs to relative paths for the documentation site's `/go-schedule` base path. Every other mapping is byte-identical. The offline `brand-check` command verifies the stored archive against the pin, rechecks its internal manifest, validates all copied or transformed consumers, and reports undeclared files in the selected brand tree.

Use the current SVG logos in `logos/svg/`, social image `logos/png/go-schedule-social-preview-1280.png`, Windows icon `platform/windows/go-schedule.ico`, macOS icon `platform/macos/go-schedule.icns`, and Linux hicolor files in `platform/linux/hicolor/`. The official archive retains the full design guide and source material without introducing upstream prose formatting into this repository.

Source and packaging checks prove asset identity and installer input wiring. They do not establish how an icon appears on every installed operating system or theme; report any such defect separately against the affected release.
