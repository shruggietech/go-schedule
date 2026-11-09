# S140 Inspected Promotion Inventory

**Decision status**: The owner approved this exact inventory with #458 open. Promotion run 42574277858 succeeded, publishing release 438780666 as public/latest v2.0.2 at 2026-11-09T06:39:01Z (November 9 Eastern). The original pending-approval packet is preserved privately.

**Release**: [v2.0.2, release 438780666](https://github.com/shruggietech/go-schedule/releases/tag/v2.0.2), now public/non-prerelease.

**Immutable source**: [fba288addcfce1769443b31ab25bd18b83a5e1e4](https://github.com/shruggietech/go-schedule/commit/fba288addcfce1769443b31ab25bd18b83a5e1e4), tree b7f56cc92bea72399ba5b64dd8919133c1b21cde, annotated tag object e038ab476c91528703c8d0d0470074f4507089d3.

**Staging**: [Release run 42568819780](https://github.com/shruggietech/go-schedule/actions/runs/42568819780), attempt 1, all seven jobs successful. All five source checkouts identify the approved commit. Authenticated downloads bind this draft ID and all eight asset IDs; every declared size and SHA-256 digest matches the downloaded bytes. Final metadata and asset readback found no drift from the repaired frozen selection.

## Eight inspected original assets

| Asset | ID | Bytes | SHA-256 |
| --- | --- | --- | --- |
| go-schedule-desktop_v2.0.2_darwin_arm64.tar.gz | 726396021 | 17911299 | 2d5d29f86b2ac37caa75fc698ededf89cc765af162a4324276f9b9c3e9597fb2 |
| go-schedule-desktop_v2.0.2_linux_amd64.tar.gz | 726385117 | 24989722 | 6e56dbf8a2fa96e7361013b7c01e0994d03397b6bdf596976c21a7616e90d42d |
| go-schedule_v2.0.2_darwin_amd64.tar.gz | 726379836 | 12303027 | 97d4620e178bb76451fc45d60f762e08bb90cf7e32b40dd90199a110af929e40 |
| go-schedule_v2.0.2_darwin_arm64.tar.gz | 726379835 | 11593679 | f60ee25d8e54d444b309ede488107c4e063879a40cbb46aea5e807106963f720 |
| go-schedule_v2.0.2_linux_amd64.tar.gz | 726379832 | 12585750 | 6c408b6d85adac7e9be2fb1a8841744d7660e0c162f98c3ca6fbd13aa3316a7c |
| go-schedule_v2.0.2_linux_arm64.tar.gz | 726379833 | 11606981 | d2c921ccea0c2bee61149fb918c49fd033042efa6bb65f27176889712e040b50 |
| go-schedule_v2.0.2_windows_amd64.msi | 726386417 | 24072192 | 7bdb7f1bf32e817ac0bfd8a8b4f10b01616b8fe192c1c2edc204d37d9361b410 |
| windows-candidate-manifest.json | 726386416 | 449 | 9a33c4a132bb599a052b88faba671a8a1f8beb202b591358d93adb2e6f74869f |

## Published checksum inventory

The approved promotion workflow added SHA256SUMS.txt as asset 726489981 after complete asset validation. It has eight entries, 844 bytes and SHA-256 b22163597713d67d37f4ca729f046bafc84de9f2e8ef67db6b493a3f8f0be786. Independent GNU measurements match all eight original hashes; its bytes match the expected Linux-workflow two-space inventory, and all entries pass checksum verification. The initial Windows binary-mode-marker comparison and its corrected interpretation remain retained diagnostically. All eight original assets were published unchanged, with no rebuild or replacement.

## Inspection results

Six archives contain the expected executable layouts and bounded member paths, with 18 README/LICENSE/CHANGELOG comparisons against frozen approved-source bytes. All 15 archived Go binaries carry the expected v2.0.2 build flag and target platform. Thirteen contain the approved embedded Git revision; both Wails GUI binaries omit that field and are bound through exact-source staging and immutable asset digests. Eight report vcs.modified=true after generated staging output exists; these flags remain literal observations. No conflicting embedded revision exists. The macOS bundle identifies tech.shruggie.goschedule, version 2.0.2.

The Windows candidate guard passed. Fresh read-only MSI inspection passed and reproduced the downloaded candidate manifest byte-for-byte: ProductVersion 2.0.2, ProductCode {8A138562-7D16-41F5-A53A-DA4DBF02F726}, x64;1033 template, exact source/run/attempt, filename, bytes and digest. The File table contains the four expected installed executable payloads. MSI identity, icon, PATH, administrative-group, lifecycle-dialog and cleanup/control relationships were inspected without installation. These are compiled-artifact observations, not installed service or attended GUI observations.

## Remaining version-resource finding

[#458](https://github.com/shruggietech/go-schedule/issues/458) remains open: the MSI File row for gosched-gui.exe reports Version 1.0.0.0 while MSI ProductVersion is 2.0.2. The Wails resource template uses Info.ProductVersion, which is not release-specific in the current project configuration. No Explorer UI, runtime-version error or upgrade failure is claimed. Correcting packaged resource bytes requires a future source/build increment; this approved immutable tag and the inspected assets remain unchanged. This metadata deviation is disclosed for the owner's promotion decision and is not represented as fixed.

## Resolved draft metadata findings

[#456](https://github.com/shruggietech/go-schedule/issues/456): the draft body now includes the existing explicit standing-waiver sentence, preserving all six maintenance highlights and the exact changelog link. Tagged notes/review remain byte-identical, with notes SHA-256 be06cfc8d4c21ec2657677634522c4fedcb4bb0f6be6c04e453d6e46c34b946c. The original helper refusal remains retained. Final draft-body SHA-256 is ebf7e74c6aa1e49a6be27a2dc176e2564749a27f2efe16d36727423461011afb.

[#457](https://github.com/shruggietech/go-schedule/issues/457): fresh GET after the body-only edit exposed a temporary untagged release association while the remote Git tag stayed unchanged. The same draft ID was restored to v2.0.2 with the approved exact target and explicit draft/non-prerelease state. All original asset IDs, sizes and digests stayed unchanged. The cause of that metadata transition is not established; original response/refusal receipts remain retained. Fresh selection and authenticated downloads passed after repair.

Frozen release-content identity is 89aa6a8c2d3ebebd5e1d37a64bc1df80051f3bf0da8251ddcbaeb1692f8ee589, using the existing helper's length-prefixed tag/name/body/target/prerelease contract. No public release metadata was edited.

## Current delivery and actual limits

The owner approved promotion with #458 open. The existing standing-waiver workflow published this same release ID and original assets, added the exact checksum inventory and completed independent public readback. Parent #444 remains open for current-discovery source integration/deployment, and #458 remains an independent open functional defect. Native installation, attended GUI/keyring and actual host/WSL coverage gaps remain explicit. No packaged program was installed or executed locally; no operator desktop focus or non-admin interaction occurred. Public/latest is v2.0.2. S140 review usage is one of two before the final PR; its automatic review is the second and final request. The original pending-decision packet remains immutable private evidence.

## Approved promotion and public readback (2026-11-09)

The owner approved promotion with #458 open. Existing standing-waiver promotion run 42574277858, attempt 1, passed. The eight original assets retain their IDs and bytes. Checksum asset 726489981 is the expected 844-byte/eight-entry inventory, SHA-256 b22163597713d67d37f4ca729f046bafc84de9f2e8ef67db6b493a3f8f0be786. All nine public URLs were independently downloaded without credentials, all declared sizes/digests and all checksum entries passed, and the candidate manifest remains byte-identical. Public/latest selects release 438780666. Current discovery source integration/deployment remains pending the final owner PR merge.
