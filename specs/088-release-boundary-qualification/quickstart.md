# S088 Validation

Run from the repository root using the verified hidden launcher on Windows.

1. Audit CHANGELOG.md and .github/release-notes/v1.4.0.md for reviewed task, controls, navigation, feedback, administration, and Notifications repairs. Retain the heading/anchor and distinguish unpublished status from qualification.
2. Run `go run ./scripts/github-format`, then `sh scripts/verify.sh all` in the foreground. Require all eight gates; existing automation fixtures validate release metadata and workflow contracts.
3. Publish with `gh pr create --body-file` and read back stored Markdown. No closing keywords for incomplete criteria.
4. Dispose every review finding within two rounds and require green latest-head checks before requesting maintainer merge.
5. After merge and exact-main checks, follow test/windows/README.md for authorized staging and Sandbox preparation. Verify exact assets before fresh/upgrade UI observations.
6. Record supported results and untested checks with the maintainer waiver. Do not finalize a fully qualified archive unless actual evidence validates.
