# Contract: S050 Release Audit Record

## Purpose

`verification.md` is the PR-reviewed, human-readable audit for the external
v1.0.0 release operation. It records enough immutable evidence to reproduce and
review the decision without copying the binary evidence archive into Git.

## Required sections

1. Authorization and immutable source boundary.
2. Tag object and peeled-commit identity.
3. Accepted Release workflow run, attempt, job conclusions, and draft release.
4. Exact eight-asset staging inventory.
5. Candidate manifest and MSI identity, size, SHA-256, ProductVersion, and
   ProductCode.
6. Formal workspace environment summary and 47-observation result table.
7. Evidence ZIP identity, independent validation output, and upload URL.
8. Disposition packet inventory and per-issue evidence/comment/closure results.
9. Promotion workflow run and its ordered gate conclusions.
10. Exact ten-asset public inventory and independent checksum results.
11. Latest release, release notes, README, changelog, tag, and binary identity.
12. Final issue, coordinator, milestone, branch, and repository audit.
13. Deviations, failures, recovery, and residual risks.

## Allowed committed data

- Public GitHub URLs and numeric identifiers.
- Cryptographic hashes, byte counts, filenames, product identifiers, counts,
  statuses, and timestamps.
- Non-personal environment roles, Windows version/build, architecture, display
  dimensions, DPI, and token-integrity classes.
- Concise summaries of observations and validation output.

## Prohibited committed data

- Screenshots or archives containing unnecessary local account names, profile
  paths, desktop contents, secrets, or unrelated applications.
- Authentication tokens, cookies, credentials, private keys, or raw environment
  variables.
- Invented pass results or copied local-demo results.
- Mutable local paths presented as durable release evidence.

## Completion rule

The record may state `PASS` only after the public release and all project state
have been independently audited. Until then it must identify the latest
completed boundary and every remaining stop condition.
