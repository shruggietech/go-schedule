# Quickstart: updating the official brand kit

1. Inspect the formal `shruggie-brand` release, hosted go-schedule migration guidance, archive name, and release `SHA256SUMS`.
2. Update `brand/source.json` to that exact release and review `brand/repository-consumers.json` against the new archive inventory. Record any required migration changes in the slice plan.
3. Copy the proposed consumer map to a separate JSON file, then run `go run ./scripts/brand-import <archive.zip> <SHA256SUMS> <proposed-map.json>` from repository root. The importer stores the complete unmodified ZIP and synchronizes only declared assets. Do not regenerate official files locally.
4. Review the archive pin, consumer map, removed stale assets, documentation, and package inputs. Run `go run ./scripts/brand-check` to detect archive or consumer drift.
5. Run `go run ./scripts/github-format` and `sh scripts/verify.sh all` before pushing. Record which native appearances were not observed.

For this slice, the pinned package is `go-schedule-brand-2.0.0-bb2.4.0.zip` and its digest is `cee5372c5554e5f7d1c87eb5fe1b8bab824dbc9c57a414c4b5d41af703094266`.
