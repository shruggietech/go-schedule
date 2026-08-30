# Quickstart: Verify S029

1. Run `sh test/scripts/spec-lifecycle-check_test.sh` and confirm the valid fixture passes while every contradiction fixture is rejected.
2. Run `sh scripts/spec-lifecycle-check.sh .` and inspect `specs/README.md` for all 29 feature records.
3. Run `sh test/scripts/automation-check_test.sh automation` and `sh scripts/automation-check.sh .`.
4. Inspect `.github/dependabot.yml` for the two bounded ecosystem entries and no dependency version edits.
5. Read repository security settings and confirm the three requested controls are enabled while unrelated controls retain their prior state.
6. Inspect issue #33 and confirm it is open, assigned to `Post-v1`, labeled P3 and `needs: verification`, with an honest deferral comment.
7. Run `sh scripts/verify.sh all` in the foreground through all eight gates.
8. Audit UTF-8 without BOM, mojibake, trailing whitespace, and the complete diff.
