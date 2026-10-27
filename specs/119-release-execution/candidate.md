# S119 staged v1.6.0 candidate

Draft release ID: 423908120. [Authenticated draft](https://github.com/shruggietech/go-schedule/releases/tag/untagged-a8820cde2ce9b86f8105). Approved source: 0f8834ba3b88b69c8f55aaa5d72add70e5b6252b. [Release run 40278046423](https://github.com/shruggietech/go-schedule/actions/runs/40278046423), successful attempt 1, all seven jobs successful. This is draft staging, not public promotion.

All eight original assets were downloaded through the authenticated exact-release helper and their actual SHA-256 values match GitHub asset digests. Existing Windows candidate validation passed against the approved source and actual MSI. ProductVersion: 1.6.0; ProductCode: {2C4FA5F6-B1CE-4F5D-B338-573F21DB761F}.

| Original asset | Asset ID | Bytes | Downloaded SHA-256 |
| --- | --- | --- | --- |
| go-schedule-desktop_v1.6.0_darwin_arm64.tar.gz | 677205773 | 17835958 | 59631c34cdea239cbaa3efbc362c9a7fcfb467484625a0377d492b231ad51fa4 |
| go-schedule-desktop_v1.6.0_linux_amd64.tar.gz | 677206044 | 24854886 | b6144343a888bd0255888d7c07973d7cb0771ce18a8e7b318150379560443244 |
| go-schedule_v1.6.0_darwin_amd64.tar.gz | 677201007 | 12248285 | ae91466e773091257110d3e0d5dbdb2d145f168ce4dd61b3d3ed5a39ce01f6b9 |
| go-schedule_v1.6.0_darwin_arm64.tar.gz | 677201005 | 11536337 | e063949656b4004636cd461a207a1d216fddec00beb0da201e6b6c93e33f3e9c |
| go-schedule_v1.6.0_linux_amd64.tar.gz | 677201004 | 12501547 | 0080d17e0ce2b117b9ad7124f65a17e866fed61d4f60d98a7136446425e7ac66 |
| go-schedule_v1.6.0_linux_arm64.tar.gz | 677201008 | 11524711 | 1c56af78b674aa198a8ea9741de46ce3a8c61b25c10fffbe6b505af29f41c2ce |
| go-schedule_v1.6.0_windows_amd64.msi | 677206524 | 23982080 | 6dcb5551380202568fb2ce091debd4c01a2e0afef8907314908bd07e4c0475d4 |
| windows-candidate-manifest.json | 677206526 | 449 | 48421f19edd9044d6f64245dc105ba9daf02108da6f906714fb986aa6640abcb |

The manifest binds repository/tag/source, workflow Release, actual run/attempt and MSI filename/size/hash. The successful Windows staging job is 121425577076 in this exact attempt. SHA256SUMS.txt is intentionally added by the existing promotion workflow; it is not an absent draft asset.

The prepared release body retains original #355/#356 unknown-cause disclosures and the standing disclosure that attended native observations remain untested. No original failed evidence is relabelled. Promotion approval would permit the existing standing-waiver workflow from reviewed main to validate these identities again, add/verify checksums, publish this same release and make it latest. Source discovery changes will still use the owner-reviewed, unmerged S119 PR.

## Authorized promotion outcome (2026-10-27)

The owner approved this exact candidate after the inventory above was presented. Existing guarded promotion run 40450478494, attempt 1, completed successfully and published the same release ID as latest at 2026-10-27T04:35:10Z. Original asset identities/bytes and release notes are unchanged. Added SHA256SUMS.txt is asset 680449833, 844 bytes, SHA-256 42db590df8671d04baa1bcd4119ce51feed36c282daf1106ccf4023750e8950c. Its eight entries independently match all original downloaded bytes.
