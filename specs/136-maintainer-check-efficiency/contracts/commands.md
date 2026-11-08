# Maintainer command contract

Existing entry points: `sh scripts/automation-check.sh [root]` and `sh scripts/spec-lifecycle-check.sh [root]`. Default root is the caller working directory; explicit relative roots resolve from that directory. Wrappers load checker code from their own source checkout, not an isolated target module.

The shared native CLI accepts `automation` or `lifecycle` and an optional root. Successful commands preserve the existing OK summary, failing commands return nonzero and report file/rule context to stderr. Invalid command arguments fail with usage. No per-spec or per-assertion subprocess occurs. The automation command retains only bounded brand and gate-manifest helpers, with hidden noninteractive Windows launches and original errors. Workflow command strings are inspected, never executed.

Existing eight-gate order remains format, vet, lint, race, gui, coverage, docs, automation. No release, installation, state-format or scheduling behavior changes.
