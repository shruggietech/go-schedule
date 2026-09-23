# Standing native-testing waiver applied to S104

The maintainer clarified that the prior testing waiver is universal, not version-specific: attended native Windows testing is not part of this project's release practice. S104 applies that standing direction to v1.5.0. No attended test result is claimed.

Automated CI, reviewed-source identity, successful platform staging, exact Windows MSI and candidate-manifest identity, complete asset inventory, and downloadable checksum verification remain mandatory. A failed, timed-out, or partial automated check is not covered by this waiver. No untested observation may be recorded as a pass, no historical candidate evidence may be reused as new evidence, and no synthetic evidence archive may be submitted to the full-attended-evidence promotion workflow.

The public release notes disclose the limitation. If actual v1.5.0 use reveals a defect, track it as a separate bug issue against the active release rather than revising the historical disposition or keeping the completed publication issue open.

At this preparation checkpoint, no v1.5.0 candidate has been staged or published. Final source, workflow run, MSI identity, checksum, asset, and public-release details belong in the #185 publication record after those facts exist. Future releases may apply the standing waiver without a new permission request, but each public release must disclose its own untested native scope accurately.
