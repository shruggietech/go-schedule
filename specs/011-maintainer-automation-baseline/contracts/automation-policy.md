# Contract: Offline Automation Policy

## Invocation

```sh
sh scripts/automation-check.sh [repository-root]
```

`repository-root` defaults to the current repository. Tests supply a temporary
fixture root so negative cases never edit the real workflows or driver.

## Approved action majors

```text
actions/checkout@v7
actions/setup-go@v7
actions/upload-artifact@v7
softprops/action-gh-release@v3
```

Every workflow `uses:` reference MUST match one of these entries. A new action
family or major is rejected until its runtime and used contract are audited and
the policy is deliberately updated.

## Required verification manifest

```text
format
vet
lint
race
gui
coverage
docs
automation
```

The policy compares this independent list with `sh scripts/verify.sh list`.
Missing, extra, duplicated, or reordered gates fail. Ordering is contractual so
aggregate diagnostics and the autopilot breakdown remain stable.

## Failure diagnostics

Failures exit non-zero and name:

- the workflow and unaudited action reference;
- the expected approved major when a known family uses an old major; or
- the expected and observed gate manifests when they differ.

The check performs no network requests and no repository mutation.

## Required fixture scenarios

1. Current approved workflows and complete manifest pass.
2. `actions/checkout@v4` fails as an obsolete known major.
3. An unknown external action fails as unaudited.
4. A missing required gate fails with both manifests in the diagnostic.
5. A duplicated or extra gate fails.
