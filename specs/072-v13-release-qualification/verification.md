# Verification: v1.3 Notifications and Local Agent Access Qualification

**Branch**: `codex/072-v13-release-qualification`

**Date**: 2026-09-09

**Issue**: [#190](https://github.com/shruggietech/go-schedule/issues/190)

**Release boundary**: Qualification evidence only. S072 does not create a tag, GitHub release, or public artifact.

## Spec Kit analysis

The specification, plan, research, data model, contract, checklists, and 18 tasks were audited before implementation. All 13 functional requirements and six buildable success criteria map to concrete tasks, all three user stories retain independent tests, and no constitution conflict or unresolved clarification remains. One medium coverage gap was found: the initial package-shaped task proved daemon health but did not exercise scheduling access required by FR-004. T008 and the integration contract were tightened to create, list, and recover a draft task through protected local IPC. A low wording mismatch that implied an already-correct product behavior must fail before implementation was also removed. Analysis then reported complete coverage with no critical or high findings.

## Package-shaped default evidence

`test/integration/v13_release_qualification_test.go` builds the current daemon into a distribution-like path containing spaces and non-ASCII text. On each supported platform it starts that candidate with isolated configuration, data, and IPC paths, confirms health, creates or recovers one local draft task, observes zero notification channels and deliveries, and observes a disabled MCP HTTP status with no endpoint, allowed origin, credential fingerprint, enable time, client name, last-access time, or request count. The daemon is stopped and the same executable and state are started again, covering both fresh and retained behavior. Windows subprocesses use the hidden-console helper.

## Notification evidence map

- `internal/notification` covers endpoint validation, safe summaries, redirect refusal, stable delivery headers, asynchronous dispatch, bounded retries, worker isolation, and shutdown.
- `internal/store` covers disabled channels, nonterminal runs, retry and restart recovery, terminal completion, bounded history, assignment precedence, migration preservation, and secret erasure.
- `internal/api/server` and `internal/api/client` cover write-only authorization, redacted channel and delivery responses, validation, lifecycle operations, and compatible serialization.
- `internal/cli` and `desktop/notifications` cover explicit configuration and redacted administrative projections.
- `test/integration/notifications_test.go` proves a successful webhook delivery does not change the source run's failure outcome or captured output.

## MCP evidence map

- `internal/mcpobserve/conformance_test.go` covers both supported protocol revisions, exact five-resource and four-template discovery, zero tools, schema stability, bounded safe errors, hostile-content isolation, and invariant server authority.
- `internal/mcpobserve/resources_test.go` covers pagination, cursor rejection, deadline propagation, bounded text and output, redaction, and secret-free projections.
- `test/integration/mcp_packaged_test.go` builds the package-shaped CLI, connects through the official SDK over stdio, verifies exact discovery, opens no listener, and closes the process transport with hidden Windows launch behavior.
- `internal/mcphttp/manager_test.go` covers official-SDK HTTP discovery for both revisions, exact Host and Origin policy, invalid and rotated credential rejection, runtime-only restart defaults, revocation, request bounds, concurrent lifecycle changes, access evidence, cancellation, and listener shutdown.
- `internal/api/server/mcp_http_test.go`, `internal/api/client/mcp_http_test.go`, and `internal/cli/mcp_test.go` cover protected local control, non-secret status, explicit enablement, rotation, disablement, and compatibility.

## Evidence status

The package-shaped notification and MCP integration command passed under the race detector on Windows in 3.699 seconds. The notification, MCP Observe, MCP HTTP, protected API server and client, and CLI packages passed under the race detector. The complete store package, including migration and notification lifecycle tests, passed under the race detector. The repository automation contract first failed with the S072 job absent, then passed after the named matrix and its fixture were added. The documentation audit confirms webhook is the sole shipped notification transport, SMTP and native desktop delivery remain future work, stdio and localhost HTTP are distinct local Observe transports, and remote MCP, Operate, Manage, and mutation tools remain unavailable.

## Canonical verification

`scripts/verify.sh all` passed in the foreground on the completed review-branch tree:

- `format`: repository-authored Markdown and GitHub text contain no Unicode em dash or hard-wrapped prose defects.
- `vet`: passed.
- `lint`: passed with zero issues.
- `race`: passed across daemon, CLI, core, scripts, and integration packages, including the package-shaped first-start and retained-state journey.
- `gui`: desktop Go suites, native Wails Windows build, all 82 frontend tests, and the production frontend bundle passed.
- `coverage`: engine 82.9%, schedule 89.1%, timezone 91.3%, store 80.7%, catchup 88.9%, and logbus 91.1%.
- `docs`: documentation policy, fixtures, links, front matter, fences, theme, and product policy passed across 17 pages.
- `automation`: workflow, CodeQL, Dependabot, release operations, release notes, brand, lifecycle, eight-gate, and negative fixture audits passed.

## Hosted evidence

The first hosted run passed the S072 matrix on Windows and Linux but failed on macOS because its long per-test temporary directory exceeded the Unix-domain socket path limit. The standard macOS race job reproduced the same new-test failure. The harness now keeps candidate data isolated while assigning Unix IPC a unique short `/tmp` endpoint with cleanup. The pull request's named `v1.3 release qualification` matrix supplies the final Windows, macOS, and Linux results for the exact reviewed commit. Those hosted checks must be green before maintainer review. S072 itself creates no release tag or public artifact.

First-round Codex review identified that the automation fragments were searched across the complete workflow rather than within the S072 job, allowing another matrix or job to mask qualification erosion. The checker now extracts only the `v13-release-qualification` block and requires its visible name, fail-fast setting, three-platform matrix, integration command, detailed notification and MCP command, and store command. A regression fixture narrows only the S072 matrix while retaining the Wails matrix and must fail.

Second-round Codex review identified that an empty `admin_group` made the package-shaped test exercise broad compatibility IPC despite the qualification's protected-local-access claim. The harness now resolves the current primary group on Unix and uses the built-in Users group on hosted Windows, creates a short daemon-managed Unix socket directory, and requires each daemon launch to emit `access_mode=restricted`. Existing installer qualification separately retains the exact `goschedadmin` creation and membership contract. No third review was requested.

## Publication audit

All 19 changed files passed strict UTF-8 decoding without BOM or mojibake markers, `git diff --check`, the GitHub publication formatter, and specification lifecycle validation. ShellCheck at warning severity passed every repository test script; its three informational findings are unchanged from `origin/main`. The branch contains only issue #190 qualification, documentation, workflow, automation-contract, and Spec Kit artifacts. Every S072 task is complete, and no tag, release workflow, issue closure, or public artifact was created.
