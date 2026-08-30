# Repository integration

The complete standalone brand kit lives in this directory. `README.md` explains the identity; `brand-guide.pdf` is the long-form visual guide; `manifest.json` and `VERIFY.md` preserve the approved build's integrity evidence.

## Source of truth

Files under `brand/` are canonical. A file under `docs/assets/`, `gui/assets/`, or `cmd/gosched-gui/` is a consumer copy only when `repository-consumers.json` declares the relationship. Do not edit a consumer copy by hand.

Run the ordinary offline check from repository root:

```text
go run ./scripts/brand-check
```

This validates the approved manifest, UTF-8 integrity, portable SVG rules, and every declared byte-identical consumer. It requires only the Go toolchain already used by the project and is part of the automation gate.

## Choosing an asset

- Use `logos/svg/` for scalable artwork and new design surfaces.
- Use the full mark at 36 px and above. Use `go-schedule-mark-reduced.svg` or the favicon family at 32 px and below.
- Use dark/light lockups on their named surfaces. Use white or black variants when a single-color mark is required.
- Use `logos/png/go-schedule-social-preview-1280x640.png` for repository and link-preview settings.
- Use the ready-made files under `platform/` for Windows, macOS, and Linux packaging.
- Use `tokens/` and `components/` when implementing a product surface. Do not sample colors from raster images.

## Updating the kit

Brand regeneration is a deliberate maintainer workflow, not a routine build step. The optional toolchain is documented in `build/README.md` and includes Python graphics/font packages plus Node.js, Playwright, and Chromium.

1. Change the canonical geometry, tokens, guide source, or components under `brand/`.
2. Run the complete generation sequence in `build/README.md`.
3. Run `python build/verify.py` and require PASS.
4. Copy changed outputs to every target declared in `repository-consumers.json`.
5. Run `go test ./scripts/brand-check` and `go run ./scripts/brand-check`.
6. Run `sh scripts/verify.sh all` before commit.

Do not commit `.venv`, `__pycache__`, rendered audit pages, ZIP transports, or backup archives. If an approved asset is renamed, update documentation and the consumer map in the same change.

## Evidence boundaries

The standalone `VERIFY.md` records visual-kit validation. The repository verifier records import and consumer integrity. Neither is evidence that an installed Windows title-bar or taskbar icon was manually observed; GitHub issue #33 remains intentionally separate.
