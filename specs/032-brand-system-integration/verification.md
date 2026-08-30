# S032 Verification

## Pre-implementation evidence

- Review branch: `codex/032-brand-system-integration` from synchronized `origin/main` at `c7f6e20`.
- Approved import source: GPT-app output `go-schedule-brand-kit`, completed 2026-08-30 under the BrandBuilder-standard workflow.
- Source status: standalone `VERIFY.md` reports PASS for 108 inventoried artifacts plus manifest/report, 28 portable SVGs, an 11-page/11-bookmark tagged PDF, measured contrast, platform ICO sizes, UTF-8, and page fit.
- Import scope: 110 files and 2,140,943 bytes (2.042 MiB) after excluding the two `build/__pycache__` files. ZIP transports, the pre-BrandBuilder backup, temporary dependency environments, and rendered audit previews are outside the source directory and excluded.
- Proportionality gate: PASS, 2.042 MiB is below the 3 MiB specification limit before repository-only mapping and guidance.

## Spec-Kit analysis

Blocking pre-implementation analysis completed after tasks generation:

- Requirements: 20 functional requirements plus 7 measurable outcomes.
- Tasks: 36 dependency-ordered tasks across setup, canonical import, three user stories, and delivery.
- Coverage: 100%; every functional requirement and buildable success outcome maps to at least one task.
- Ambiguities: 0.
- Duplications: 0.
- Constitution conflicts: 0.
- Critical/high findings: 0 after correcting the verifier requirement to remain implementation-neutral and assigning parallel test tasks to separate files.
- Unmapped tasks: 0; setup and delivery tasks cover lifecycle, analysis, verification, and traceability obligations.

Implementation is authorized to proceed.

## Focused implementation evidence

- `go test ./scripts/brand-check`: PASS, including missing/hash/path manifest failures, UTF-8/BOM/mojibake and SVG portability failures, and missing/different/duplicate consumer failures.
- `go run ./scripts/brand-check`: PASS (`108 artifacts, 28 SVGs, 59 consumers`). The changed-consumer negative test proves a byte-different mapped target reports both canonical and consumer paths.
- `sh scripts/docs-check.sh`: PASS (`12 pages`), including brand-page policy fixtures for full, reduced, horizontal, monochrome, social, guide, inventory, palette, typography, attribution, and misuse coverage.
- Brand-page presentation: PASS. The page uses responsive Markdown images and tables with no fixed-width wrapper, the approved transparent lockup rather than a mismatched square tile, concise selection guidance, and working local downloads. The canonical 1280 x 640 social composition was visually inspected at repository resolution; the imported eleven-page guide retained the completed BrandBuilder visual audit.
- `go test ./gui/...`: PASS, including decoded embedded dimensions of 1024 x 1024 for the full application mark and 256 x 256 for the reduced window mark.
- `go test ./test/integration -run 'TestWindowsInstallerIdentityContract|TestWindowsInstallerGUIResourceContract' -count=1`: PASS.
- `pwsh -NoProfile -File build/windows/verify_wxs.ps1`: PASS.
- `sh test/scripts/automation-check_test.sh automation`: PASS, including missing brand-check and stale release-artifact rejection fixtures.

## CI-parity evidence

`sh scripts/verify.sh all` completed in the foreground on 2026-08-30 with all eight gates passing:

- format: PASS, no files reported;
- vet: PASS;
- lint: PASS, `0 issues`;
- race: PASS, including `scripts/brand-check` and the 57.851-second integration package;
- GUI: PASS;
- coverage: PASS (engine 89.5%, schedule 89.2%, timezone 91.3%, store 86.9%, catchup 88.9%, logbus 91.1%);
- docs: PASS, 12 pages and brand policy fixtures clean;
- automation: PASS, Actions, CodeQL, Dependabot, brand, lifecycle, and the eight-gate manifest clean.

## Final audits and issue disposition

- Final Spec-Kit analysis: PASS. All 20 functional requirements and 7 measurable outcomes remain covered by 36 dependency-ordered tasks; ambiguity, duplication, constitution-conflict, critical/high-finding, and unmapped-task counts are all zero. Both the 15-item requirements checklist and the 20-item brand-system checklist remain fully checked against the delivered result.
- Live issue audit on 2026-08-30: #10 and #34 are open and fully satisfied by the canonical kit, public page, reproducible source, derived asset family, and wired consumers, so the eventual pull request is eligible to use `Closes #10` and `Closes #34`. Issue #33 remains open with `needs: verification`; this slice makes no install/uninstall observation claim and must not close it.
- Imported-size audit: PASS, 112 files and 2,152,025 bytes (2.052 MiB) under `brand/`, including the 108 inventoried artifacts and four repository/manifest controls. This remains below the 3 MiB limit and adds no application runtime dependency.
- Encoding audit: PASS, 129 changed or added text files decode as strict UTF-8 without BOM. A direct mojibake-marker scan is clean; the verifier's corruption fixtures use Unicode escapes so intentional negative-test values cannot be mistaken for damaged source.
- Diff audit: PASS. `git diff --check` reports no whitespace errors. Canonical and declared consumer asset directories are marked `-text` in `.gitattributes`, preventing Git line-ending normalization from invalidating approved hashes or producing checkout-dependent copy drift.
- Historical inventory correction: the stale local-review evidence for already merged S030 and S031 was updated to PR #68 and PR #69. This small deviation is required to keep the lifecycle index truthful and does not expand S032 product scope.
