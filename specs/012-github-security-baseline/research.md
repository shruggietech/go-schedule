# Research: GitHub Security Baseline

**Date**: 2026-08-26

## D1 - Use advanced, repository-owned CodeQL setup

**Decision**: Commit an advanced CodeQL workflow instead of enabling opaque default setup.

**Rationale**: GitHub recommends default setup for many repositories, but this slice specifically requires a reviewable pull-request/default-branch/weekly trigger contract, pinned action families, a stable named check, and offline drift tests. Advanced setup makes those properties repository-owned and lets the analysis build avoid the desktop graphics path.

**Alternatives considered**:

- Default setup: simpler operationally, but its configuration cannot satisfy the repository's offline workflow-contract requirements.
- CodeQL CLI outside Actions: adds installation and result-upload complexity without improving this repository's security boundary.

**Primary sources**:

- <https://docs.github.com/en/code-security/concepts/code-scanning/setup-types>
- <https://docs.github.com/en/code-security/how-tos/find-and-fix-code-vulnerabilities/configure-code-scanning/configuring-advanced-setup-for-code-scanning>

## D2 - Use CodeQL Action v4 and existing Node 24 action majors

**Decision**: Use `github/codeql-action/init@v4` and `github/codeql-action/analyze@v4`, alongside the already approved `actions/checkout@v7` and `actions/setup-go@v7`.

**Rationale**: The v4 CodeQL action line uses the supported Node 24 runtime and is the maintained major line. The repository's policy intentionally approves floating major selectors after research, so the CodeQL entries follow the same maintenance model established by Feature 011.

**Alternatives considered**:

- CodeQL v3: rejected because it is the previous action major and does not satisfy the repository's Node 24 runtime baseline.
- Commit-SHA pinning: stronger immutability, but changes the repository's established floating-major policy and belongs in a dedicated supply-chain hardening slice.

**Primary sources**:

- <https://github.com/github/codeql-action>
- <https://github.com/github/codeql-action/blob/main/CHANGELOG.md>

## D3 - Analyze Go with a manual headless build

**Decision**: Initialize CodeQL for `go` with `build-mode: manual`, select Go from `go.mod`, and run `CGO_ENABLED=0 go build ./...` before analysis.

**Rationale**: Go is compiled, so CodeQL needs an observable build. The full repository includes a Fyne/OpenGL GUI entry point whose native dependencies are not required to analyze the normal cgo-free source tree. The existing project already guarantees that `go build ./...` works with cgo disabled. A manual build makes the analyzed path deterministic and avoids autobuild guessing.

**Alternatives considered**:

- CodeQL autobuild: convenient, but it may select a cgo/desktop build path and fail on hosted graphics dependencies.
- Install OpenGL development packages: expands the workflow and attack surface for no analysis benefit because the GUI entry point is build-tagged.

**Primary sources**:

- <https://docs.github.com/en/code-security/concepts/code-scanning/codeql/codeql-for-compiled-languages>
- <https://docs.github.com/en/code-security/reference/code-scanning/workflow-configuration-options>

## D4 - Use push, pull-request, and weekly triggers with least privilege

**Decision**: Run on pushes to `main`, pull requests targeting `main`, and a weekly cron. Grant workflow-level `contents: read` and `security-events: write` only.

**Rationale**: Pull requests provide pre-integration feedback, pushes establish the default-branch baseline, and weekly runs pick up new queries even when code does not change. Reading source and publishing SARIF are the only required repository privileges. The workflow consumes no secrets and does not use `pull_request_target`.

**Alternatives considered**:

- Event-only analysis: misses new query findings in unchanged source.
- Broad `write-all` or repository token defaults: unnecessary and contrary to least privilege.

## D5 - Extend the existing offline automation checker

**Decision**: Keep one policy script. Add the two CodeQL action references to its explicit allowlist and validate one canonical CodeQL workflow for required events, branch filters, permissions, language, build mode, headless build, and analysis steps. Extend temporary fixtures with each required negative case.

**Rationale**: Feature 011 already established `scripts/automation-check.sh` as the repository's independent, offline automation policy. A second parser would duplicate discovery and failure reporting. Static shell checks are sufficient for the deliberately narrow YAML contract and introduce no network or parser dependency.

**Alternatives considered**:

- Add a YAML parser: more structurally precise, but adds a runtime dependency for a fixed workflow schema.
- Depend only on GitHub's hosted validation: cannot detect drift before push and violates the local CI-parity requirement.

## D6 - Activate and report each GitHub setting separately

**Decision**: After authorization, enable private vulnerability reporting and request enablement for secret scanning, push protection, non-provider patterns, and validity checks through the GitHub API. Read each result back separately and classify it as `enabled`, `unavailable`, `unverified`, or `failed`.

**Rationale**: GitHub plan and token capabilities can differ per control. A single green summary would hide partial support. No alert contents or fabricated advisories are needed to prove activation.

**Alternatives considered**:

- Treat an accepted PATCH as proof: misses policy or plan normalization after the request.
- Create a sample advisory or seed secret: creates misleading security data and is explicitly outside scope.

**Primary sources**:

- <https://docs.github.com/en/code-security/how-tos/report-and-fix-vulnerabilities/configure-vulnerability-reporting/configure-for-a-repository>
- <https://docs.github.com/en/code-security/how-tos/secure-your-secrets/detect-secret-leaks/enable-secret-scanning>
- <https://docs.github.com/en/enterprise-cloud@latest/code-security/secret-scanning/enabling-secret-scanning-features/enabling-push-protection-for-your-repository>
- <https://docs.github.com/en/enterprise-cloud@latest/repositories/managing-your-repositorys-settings-and-features/enabling-features-for-your-repository/managing-security-and-analysis-settings-for-your-repository>
