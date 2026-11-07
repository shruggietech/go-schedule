# S135 Inspected Promotion Inventory

**Decision status**: Owner approved this exact inventory. Promotion run 42335391283 succeeded, publishing release 437271055 as public/latest v2.0.1 on 2026-11-07T03:45:47Z (November 6 Eastern). The original pending-decision packet is preserved privately.

**Release**: [v2.0.1, release 437271055](https://github.com/shruggietech/go-schedule/releases/tag/v2.0.1).

**Immutable source**: [25f42d6eb35ff0cec2753f0f2934396ba50df954](https://github.com/shruggietech/go-schedule/commit/25f42d6eb35ff0cec2753f0f2934396ba50df954), tree 81b4f2fa133bfca59b8ac7343f9780357782979d, annotated tag object f44358627c5e84b9e47681a1c3ecea4b1b241e1d.

**Staging**: [Release run 42329729594](https://github.com/shruggietech/go-schedule/actions/runs/42329729594), attempt 1, all seven jobs completed successfully. Each checkout log identifies the approved source. Fresh authenticated downloads bind the draft release ID and all asset IDs. GitHub's declared sizes/digests match every downloaded byte. Final authenticated readback found no metadata or asset-selection change.

## Eight inspected original assets

| Asset | ID | Bytes | SHA-256 |
| --- | --- | --- | --- |
| `go-schedule-desktop_v2.0.1_darwin_arm64.tar.gz` | 719364832 | 17894740 | `bf5e8ea7acb7d1e5cfefd67839bc57da4eb961440533568579f063d56e2c8668` |
| `go-schedule-desktop_v2.0.1_linux_amd64.tar.gz` | 719362236 | 24981889 | `9641c70dbe647b522306049ed76d147c60f3a3c6f458d6e1edd64fd1fa3bea02` |
| `go-schedule_v2.0.1_darwin_amd64.tar.gz` | 719357969 | 12291158 | `9cdce7dc6d90c2ba0ff29d6b261ae3fbe311847ffcbb2860a8aae8804c15a17d` |
| `go-schedule_v2.0.1_darwin_arm64.tar.gz` | 719357966 | 11580595 | `b86f9d8b7c30c0e183a3dc6d79d7ffdba17cebc7e6a4fcb93fd7e729e692a991` |
| `go-schedule_v2.0.1_linux_amd64.tar.gz` | 719357965 | 12574505 | `72c7ec3b85cc26017b6c101a240e8abf221e50ce11ab31c9fc292484cecb878d` |
| `go-schedule_v2.0.1_linux_arm64.tar.gz` | 719357967 | 11594930 | `0cbd8e2facfeb6884532602ef185b5660e01411f9c940c58b61a4f76ca49a1a7` |
| `go-schedule_v2.0.1_windows_amd64.msi` | 719363688 | 24064000 | `29fba3a0fd05185a011221748a779a6118c92f3c63cd04e9a43171ada72cc329` |
| `windows-candidate-manifest.json` | 719363684 | 449 | `b9b14595c3f7069022f93294d4079324731a6dd105e9a43f39da924db30ca397` |

## Checksum addition during promotion

The existing waived promotion workflow generates and uploads SHA256SUMS.txt as the ninth asset. The independently generated expected file has eight entries, 844 bytes and SHA-256 `e6fc144048ed7a93244134e3f6ab0021fccca1f89169a1e0f296f8dbeafa6e69`. It is retained privately under .test-temp/s135/expected-SHA256SUMS.txt and matches published checksum asset 719508417 byte-for-byte. Approval covers publishing the eight original assets unchanged, adding that checksum inventory and selecting v2.0.1 as latest through the existing workflow; no rebuild or original-asset replacement is proposed.

## Inspection results and limits

All six archives contain their expected executables and source-identical README, LICENSE and CHANGELOG (18 document comparisons). Fifteen Go binaries carry the expected v2.0.1 build flag and target platform. Thirteen contain the accepted Git revision; both Wails GUI binaries omit embedded VCS metadata, so their source binding comes from the immutable tag, exact successful staging checkout/run and downloaded asset digest. Eight binaries report vcs.modified=true after generated staging output exists; these flags are retained and are not represented as clean-source metadata. No conflicting embedded revision was found. The macOS bundle identifies tech.shruggie.goschedule, version 2.0.1.

The existing Windows candidate guard passed. Fresh read-only compiled-MSI table inspection passed and reproduced the downloaded candidate manifest byte-for-byte, including product version 2.0.1, product code {C1743E66-7314-4EC5-B294-200AF58D4967}, source, run/attempt, size and digest. No installed-product, service-lifecycle or attended GUI result is claimed.

The first private archive probe stopped on its overstrict assumption that every executable embeds a Git revision. That original log/output remains diagnostic. The corrected fresh probe rejects conflicting revisions, validates version/target and records absence explicitly; it does not substitute an inferred embedded revision. All original 423 retained inputs remain unchanged at the pre-inspection audit. At inspection, prior v2.0.0 public release 430039348 retained its original nine asset identities and remained latest; approved promotion subsequently selected v2.0.1.

## Publication continuation

The owner approved promotion and existing promote-release-waived.yml completed successfully from reviewed main. All nine assets were downloaded independently without credentials, all eight checksum entries passed and the approved originals retained their IDs and bytes. Current discovery and the already-authorized source PR follow; integration and guide deployment remain separately reported. Native installation and attended GUI remain untested under the standing waiver. #419 remains open for actual public package/discovery delivery. At publication readback before PR creation, zero S135 Codex requests had been used. The official PR records subsequent review evidence; the maximum remains two including automatic PR-open review. Stop before final owner PR merge.
