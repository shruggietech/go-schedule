# S140 release checkpoint contract

## Source checkpoint

Inputs are accepted S139 source, the previous comparison, official vendor documentation and existing release notes. Outputs are a complete new v2.0.2 snapshot/review, the bounded #453 guide correction and a release continuation. Historical files, public selection and protected evidence are immutable.

The existing comparison-check rejects missing cells, malformed metadata, placeholder decisions and stale dataset/changelog/notes hashes. release-source-check validates prepared source metadata; it does not establish a tag, package or public release. Existing site tests run against both candidate-selected output and restored default-public output using isolated headless browser artifacts.

## Later same-S140 continuation

1. Owner reviews and merges the pre-tag official PR. Select its exact accepted main and record its SHA/tree and applicable CI.
2. Present the concrete tag target and draft-stage action under the existing exact-source approval boundary. After approval, create v2.0.2 at that source and use the existing release workflow to stage draft packages.
3. Inspect the actual complete asset inventory, names, sizes, SHA-256 values, artifact provenance and package identities. Preserve any failure. Present the inspected release ID/inventory for the existing separate promotion decision.
4. After promotion approval, publish that inspected release as public/latest, independently read back every asset and metadata identity, and preserve the standing native-testing disclosure.
5. Update current discovery/install/first-task availability to the actually promoted v2.0.2 in the second official source PR. Monitor CI and resolve findings within the same two-total-Codex-review budget, then stop for owner final review/merge.
6. Verify deployed guide/readback after owner merge and reconcile #444 only when its functional public-delivery outcome exists. Do not turn verification activity into an issue closure gate; file any revealed product defect separately.

No step launches the shared native desktop, changes host networking, requests an attended walkthrough or claims actual native/cross-OS execution that was not observed. The current run ends at step 1 readiness.
