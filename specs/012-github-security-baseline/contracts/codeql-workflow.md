# Contract: Advanced CodeQL Workflow

The canonical workflow is `.github/workflows/codeql.yml`.

## Required structure

- A stable human-readable workflow name containing `CodeQL`.
- `push.branches` includes only the required default-branch baseline (`main`).
- `pull_request.branches` includes `main`.
- `schedule` contains the reviewed weekly expression `17 4 * * 1`; changing
  cadence requires updating the offline contract and its fixtures together.
- Workflow permissions are exactly sufficient for source read and result
  publication: `contents: read` and `security-events: write`.
- The analysis job runs on a GitHub-hosted Ubuntu runner.
- Checkout uses `actions/checkout@v7`.
- Go setup uses `actions/setup-go@v7` and `go-version-file: go.mod`.
- Initialization uses `github/codeql-action/init@v4` with `languages: go` and
  `build-mode: manual`.
- The build step runs `CGO_ENABLED=0 go build ./...`.
- Result publication uses `github/codeql-action/analyze@v4`.
- No step consumes a repository or organization secret.
- No job-level `permissions` key can override the exact workflow-level grant.
- `.github/workflows/ci.yml` is unchanged.

## Offline rejection contract

`scripts/automation-check.sh` must fail with an actionable path/message when:

1. either CodeQL action uses an obsolete major such as v3;
2. any workflow introduces an unknown action family;
3. push, pull request, `main`, or schedule coverage is absent;
4. `contents: read` or `security-events: write` is absent or weakened; or
5. initialization, manual headless build, or analysis publication is absent.

The checker also rejects impossible or unreviewed schedule expressions,
job-level permission overrides, and both dotted and bracket-form secret access.

The checker intentionally validates the repository's fixed contract rather than
attempting to be a general YAML validator.
