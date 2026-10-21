# Coordinator Contract

The standalone maintainer command exposes inventory, capture, resume, finalize, validate and render. Capture accepts exact binaries/source identity, output root, budgets and selected existing WSL distro/user; missing prerequisites are reported as blocked environments. All diagnostics are bounded, redacted and actionable. Successful capture may still require semantic review and does not imply overall pass.

An environment worker verifies applied OS resource controls and its ownership before starting the complete manifest. It launches only owned binaries/tasks/receivers and uses scoped profile/preference/credential/service/GUI adapters. Interruption checkpoints completed evidence and blocks uncompleted planned cases. Total watchdog and OS lifetime independently bound execution.

Capture writes report.json, report.md, events.jsonl, manifest.json, corpus.json and evidence-manifest.json. Finalize consumes externally authored reviews.json bound to the frozen run/corpus hash, validates all artifacts and writes canonical final reports. Render takes validated canonical input only and is byte-deterministic. Exit codes are 0 for successful processing/passing evaluation or complete pending-review capture, 1 for usage/processing/integrity errors, 2 for a failed evaluation, and 3 for a blocked evaluation or incomplete capture. The PowerShell entry point preserves outcome codes; pending-review capture is not a passing QA result.

Resume requires exact subject/binary/manifest/runner/budget/environment identity, valid retained hashes and no active prior ownership; in-progress cases rerun in newly owned runtime state. Never reuse a previous success against changed inputs. The runner does not modify GitHub or source.
