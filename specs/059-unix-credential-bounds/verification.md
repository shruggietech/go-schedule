# Verification: Unix Credential Bounds

## Spec-Kit readiness

- Requirements checklist: 16/16 complete.
- Security requirements checklist: 16/16 complete.
- Clarification: no unresolved material ambiguity after three autopilot decisions were recorded from issue #145 and the existing executor contract.
- Initial analysis: 14/14 functional requirements and 7/7 measurable outcomes map to 20 valid, dependency-ordered tasks. No ambiguity, duplication, unmapped task, critical/high finding, or constitution conflict was found.
- Checklist procedure: the installed prerequisite required `plan.md` despite the mandated checklist-before-plan order. The established direct-feature-path workaround preserved the required order without changing Spec-Kit tooling.

## Test-first evidence

- PASS: the Linux-targeted executor test build failed before implementation with undefined `credentialForUser` and `applyResolvedRunAs`, demonstrating that the new boundary and atomicity cases exercised behavior absent from the baseline.

## Implementation evidence

- `credentialForUser` parses UID and GID with base 10 and a 32-bit unsigned limit, returns field/account context on failure, and returns no partial pair.
- `applyResolvedRunAs` receives the complete pair before allocating `SysProcAttr`, replacing `Credential`, or changing `Env`.
- Synthetic tests cover zero, maximum, leading zeroes, empty, negative, whitespace-padded, fractional, hexadecimal-prefixed, alphabetic, and overflowing values for UID and GID.
- Command-state tests preserve caller-owned process attributes, credential identity, and environment after UID failure and after valid-UID/invalid-GID failure.
- Compatibility tests retain named-account credentials, numeric-account lookup, explicit home, inherited home, duplicate identity cleanup, and empty `run_as` behavior.
- The current Windows package test passed and the complete focused test suite cross-compiled successfully for Linux and macOS. This Windows host has no WSL, Docker, or Unix runtime, so execution of the build-tagged cases remains pending the hosted Linux and macOS race jobs.
- Audit scope found one production construction and one assignment of `syscall.Credential`; both receive values returned only after `ParseUint(..., 10, 32)` succeeds. No `Atoi` or signed/architecture-sized intermediate remains in the affected executor path.

## Canonical verification

- PASS: format, including `github-format: OK - no em dashes or hard-wrapped Markdown prose`.
- PASS: `go vet ./...`.
- PASS: golangci-lint v2.12.0 with zero issues.
- PASS: race across the canonical non-GUI package set, including `internal/executor` on the local Windows target and the 52.912-second integration suite.
- PASS: headless GUI and view-model tests.
- PASS: coverage, engine 81.9%, schedule 89.2%, timezone 91.3%, store 80.1%, catchup 88.9%, and logbus 91.1%.
- PASS: documentation policy, fixtures, 15-page link/front-matter/fence/theme validation, and product-policy checks.
- PASS: automation policy and fixture suite, including actions, CodeQL, Dependabot, release operations, lifecycle, and the eight-gate manifest.
- PASS: changed-file UTF-8, BOM, mojibake, whitespace, and GitHub publication-format audit.
- The initial bare `sh scripts/verify.sh all` invocation did not start because `sh` was absent from the PowerShell PATH. The unchanged canonical script was then run in the foreground through the installed `C:\Program Files\Git\bin\bash.exe` and completed successfully; no gate was substituted or skipped.
- Final Spec-Kit analysis retains 100 percent requirement/outcome task coverage, zero unmapped tasks, zero ambiguity or duplication findings, and zero constitution conflicts. Both 16-item checklists remain complete.

## Issue disposition

- Issue #145 is locally implementation-complete. Hosted Linux/macOS execution, CodeQL, and review evidence remain required before final maintainer review.
