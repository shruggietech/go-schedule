# Quickstart: S104 release operator sequence

1. Merge the reviewed source metadata PR and confirm exact-main CI and CodeQL success.
2. Confirm `v1.5.0` has no pre-existing remote tag or public release, then create one annotated tag on reviewed main.
3. Wait for the `Release` workflow to finish successfully and inspect the draft's eight staged assets.
4. Download the exact MSI and manifest. Verify candidate identity against repository, tag, and commit, then audit all asset sizes and digests.
5. Record the standing native-testing waiver's application on #185 and confirm the draft release notes explicitly disclose untested checks. Dispatch the reviewed `Promote Release (standing native-testing waiver)` workflow on `main` with the exact staged tag. It checks the eight-asset allowlist, successful staging run and Windows job attempt, MSI candidate identity, checksums, final draft download, immutable tag, and latest-release transition. Never invent evidence, dispatch the full-evidence `Promote Release` workflow without genuine evidence, or hide an automated failure.
6. Wait for waived promotion to pass, then download public assets afresh and verify `SHA256SUMS.txt` and latest-release identity.
7. Comment on and close #185 for completed publication; leave SMTP and distributed-scheduling issues open. Reconcile milestone and project status.
