# Quickstart: S104 release operator sequence

1. Merge the reviewed source metadata PR and confirm exact-main CI and CodeQL success.
2. Confirm `v1.5.0` has no pre-existing remote tag or public release, then create one annotated tag on reviewed main.
3. Wait for the `Release` workflow to finish successfully and inspect the draft's eight staged assets.
4. Download the exact MSI and manifest. Verify candidate identity against repository, tag, and commit, then audit all asset sizes and digests.
5. Record the maintainer's v1.5.0-specific waiver on #185. Use the separately documented manual publication path and disclose untested checks. Never invent evidence, dispatch the full-evidence `Promote Release` workflow without genuine evidence, or hide an automated failure.
6. Generate the all-asset checksum inventory, publish the draft, download public assets afresh, and verify checksums and latest-release identity.
7. Comment on and close #185 for completed publication; leave SMTP and distributed-scheduling issues open. Reconcile milestone and project status.
