# Research: official brand 2.0.0 adoption

## Formal source

Decision: Pin `go-schedule-brand-2.0.0-bb2.4.0.zip` from the `shruggie-brand` `v2.4.0` release. The hosted portfolio and download page present this filename, and the downloaded 6,024,943-byte archive hashes to `cee5372c5554e5f7d1c87eb5fe1b8bab824dbc9c57a414c4b5d41af703094266`, matching that release's `SHA256SUMS` entry. The bundle identifies brand 2.0.0, BrandBuilder 2.4.0, canon 1.6.0, and source revision `d33cb8c658eac2f66b58e715274e0c791bad9ba3`.

Alternative considered: Download loose files from the hosted site. Rejected because it loses one release identity and complete checksum inventory.

## Migration scope

Decision: Store the complete distributable archive unchanged and extract only mapped product assets into repository presentation paths. The official migration contract marks identity, platform assets, documentation, and recovery as required. It makes the existing reduced prompt and cursor geometry the current primary mark; palette, typography, and Web/React and egui adapter APIs remain unaffected.

Alternative considered: Regenerate the kit locally. Rejected because the formal release already ships verified assets and a compiler-pinned recovery path; regeneration would create a new unreviewed artifact identity.

## Repository consumers

Decision: Keep a single explicit mapping of source paths inside the pinned archive to documentation, frontend, and packaging destinations. Derive Linux hicolor PNG copies from the official size-matched web icon files, preserve repository-authored `.desktop` metadata, and use the official Windows ICO for current packaging. Windows and macOS application icons come from the official platform suites. Stale copies from the previous mapping are removed during an update.

Alternative considered: Keep legacy file names and copies without mapping. Rejected because old artwork could continue to be distributed silently.

## Integrity and format boundary

Decision: Verify archive SHA-256 against both a repository pin and the formal release checksum inventory, validate archive paths and the bundled manifest before mutation, and run an offline repository check in CI. The exact upstream kit remains in its ZIP, so its width-wrapped prose and one em dash do not become published repository text. Repository-authored control files and GitHub text still pass `github-format`.

Alternative considered: Extract and reformat the entire upstream kit. Rejected because any byte edit breaks its manifest, while exact extraction publishes an em dash and width-wrapped Markdown contrary to repository rules.

## Limitations

Source and build checks can establish exact asset inputs. Native taskbar, title-bar, tray, and installer appearance remain untested until an attended environment observes them; this does not delay functional issue disposition under the current constitution.
