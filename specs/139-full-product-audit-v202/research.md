# Research: Full-product audit

## Decision: Reuse complete CLI evidence contracts

The existing evaluator inventories gosched only, executes ordinary cases as binaries and labels harness/deterministic cases. Six-binary preparation, per-case progress, frozen source/corpus identities and semantic reviews already exist. Reuse them and add a supplemental entrypoint register rather than introducing another report framework. Research source: internal/clievaluate/inventory.go, build.go, manifest.go, scenarios_catalog.go and docs/cli-qa.md.

## Decision: Treat native UI isolation as a prerequisite

Windows same-logon desktop isolation is insufficient because desktop/main.go calls SignalGUI on a duplicate and activation_windows.go uses fixed Local session objects. Existing WSL has GTK3/WebKit4.1 but no Xvfb/Xorg/Xwayland or window manager. Do not launch native UI in the owner's desktop or default WSLg. Capture headless frontend fixtures separately and report unavailable native evidence literally. Microsoft's kernel-object namespace documentation confirms Local session scope: https://learn.microsoft.com/en-us/windows/win32/termserv/kernel-object-namespaces.

## Decision: Observe real cross-platform feasibility without host reconfiguration

The selected existing WSL environment has native Linux execution but network mode none and no nonloopback route/address. Actual endpoint attempts will distinguish unavailable transport from product refusal. Existing remote authority fixtures are same-host and cannot establish cross-OS behavior. Native Secret Service calls may show unlock prompts; use existing isolated adapters for catalogue work and report native credential limits explicitly.

## Decision: Investigate source hypotheses with failing observations

Daemon flag parsing appears not to reject positional operands. Several task/group/chain/service/GUI acknowledgements print human text despite the documented global JSON contract, while catalogue JSON validation excludes them. Confirm exact affected paths with source-bound regressions, file deduplicated defects and make proportional corrections. Do not infer every JSON exclusion is broken.

## Alternatives rejected

Another general QA framework duplicates strict existing evidence contracts. A default native UI launch, shared-logon Windows desktop or inherited Linux keyring could interrupt the owner. Reconfiguring WSL networking, installing display prerequisites or replacing the installed service exceeds the approved existing-environment audit boundaries. Same-host mocks cannot substitute for actual cross-platform observations. Raw public logs/screenshots are unacceptable without private-data inspection.
