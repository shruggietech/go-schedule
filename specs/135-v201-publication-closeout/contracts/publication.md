# S135 Publication Contract

The existing tag-triggered Release workflow consumes accepted source at refs/tags/v2.0.1 and stages four daemon/CLI archives, two desktop archives, one Windows MSI and windows-candidate-manifest.json in a draft release. The tag resolves to 25f42d6eb35ff0cec2753f0f2934396ba50df954. Staging requires the owner's explicit exact-source decision.

An inspected selection binds the draft release, eight original asset IDs/names/sizes/digests, successful staging run/attempt and Windows candidate source. The owner decides promotion against that actual selection. Existing promote-release-waived validates immutable source, complete allowlist, successful run/attempt, manifest and checksums, then publishes unchanged content with its standing native-testing disclosure. It adds the ninth SHA256SUMS.txt asset covering all eight original assets.

Public/latest and all nine downloaded assets are read back independently. Current source discovery changes follow publication through a normal reviewed PR. Pages can select the newly public comparison after promotion before that PR merges; current guide deployment waits for owner integration. Failed or incomplete phases are reported separately without replacing original observations or claiming a completed #419 outcome.
