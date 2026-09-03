# Contract: Windows Release Candidate Evidence and Promotion

## Supported Inputs

The release-gate command accepts:

```text
windows-release-gate verify-candidate --candidate-manifest <windows-candidate-manifest.json> \
  --artifact <candidate.msi> --repository shruggietech/go-schedule \
  --tag vX.Y.Z --commit <40-hex>
windows-release-gate validate --evidence <bundle/evidence.json> --artifact <candidate.msi>
windows-release-gate validate --evidence <bundle/evidence.json> --artifact <candidate.msi> \
  --candidate-manifest <windows-candidate-manifest.json> \
  --repository shruggietech/go-schedule --tag vX.Y.Z --commit <40-hex>
windows-release-gate verify-bundle --bundle <evidence.zip> --artifact <candidate.msi> \
  --candidate-manifest <windows-candidate-manifest.json> \
  --repository shruggietech/go-schedule --tag vX.Y.Z --commit <40-hex>
```

The manifest-only form proves candidate identity independently. The local evidence form supports operator iteration. Promotion supplies all expected identity flags and uses both independent candidate validation and the archive form.

## Exit Contract

- `0`: complete passing evidence for the exact candidate and expected identity.
- `1`: well-formed invocation whose candidate or evidence is not release-ready. Every discovered defect is printed.
- `2`: invalid command usage or an internal validation failure that prevents a reliable decision.

The command writes human-readable diagnostics to standard error and a one-line success summary to standard output. It never repairs or rewrites evidence during validation.

## Candidate Contract

- The artifact must be a regular `.msi` file with a canonical filename for the tag.
- Its byte length and SHA-256 must match `candidate`.
- `product_version` must equal the tag without `v`.
- Promotion-supplied repository, tag, and commit must match the record exactly.
- ProductCode is obtained from compiled MSI inspection during staging and preserved in both the candidate manifest and attended evidence.
- Run ID and attempt identify the staging workflow that created the draft asset.

## Evidence Contract

- JSON decoding rejects unknown fields, multiple JSON values, invalid UTF-8, and a UTF-8 BOM.
- Promotion requires `evidence_class: attended-windows`; inert checked-in data is marked `automated-fixture` and cannot pass the production validator.
- Every fixed scenario identifier occurs exactly once.
- Every environment and attachment reference resolves exactly once.
- Every observation is explicitly `pass`; all other states block while retaining their distinct diagnostic.
- Scenario-specific metrics satisfy the rules in `data-model.md`.
- Every attachment path remains beneath the evidence root, exists as a regular file, and matches its recorded size and SHA-256.
- Timestamps are RFC 3339, ordered, and bounded by the evidence-run interval.
- Routine-client observations use a medium-integrity intended-user environment. Administrative preparation cannot substitute for them.
- Window scenarios collectively cover standard DPI, high or mixed DPI, clean profile, and retained v0.9.1 profile. They require native JSON plus image evidence.
- Task scenarios require one strict `task-run-evidence-v1` attachment with exactly four canonical records. Retained task definition, output, marker, history, invocation, exit, and diagnostic values must match the observation metrics.
- Error, setup, and removal scenarios require image evidence; lifecycle records also require process/session identity, option/target digests, populated-state fingerprints, unaffected-control fingerprints, security-state disposition, and reinstall result.

## Attended Collector Contract

`test/windows/Invoke-ReleaseCandidateAttended.ps1` is a resumable operator tool with explicit actions:

- `Initialize`: validate the candidate identity and create a new evidence workspace without overwriting one.
- `CaptureWindow`: sample the exact GUI PID and HWND into a UTF-8 JSON observation fragment, including process, window, monitor, DPI, work-area, and state data.
- `RecordObservation`: import one operator-reviewed observation fragment only when its ID is not already present.
- `Finalize`: assemble canonical `evidence.json`, hash attachments, invoke the shared gate, and create the evidence ZIP only when validation passes.

All launched console applications use hidden, redirected, noninteractive child processes. Destructive or machine-mutating actions honor PowerShell `ShouldProcess`. Cancellation and unavailable infrastructure are recorded explicitly and never converted to pass.

## Workflow Contract

### Staging

- A semantic-version tag triggers release staging.
- Every `softprops/action-gh-release` upload explicitly sets `draft: true`.
- All platform artifacts are attached to the same draft release.
- The Windows job writes and uploads the candidate manifest beside the MSI.
- No staging job generates the final checksum file or publishes the draft.

### Promotion

- Manual dispatch accepts exactly one semantic-version tag.
- Promotion requires `contents: write` and `actions: read`; validation jobs otherwise use least privilege.
- It retrieves the existing draft release and rejects missing, already-public, wrong-target, or mismatched-tag state.
- It rejects missing or unrecognized assets, then downloads the exact allowlisted draft asset set.
- It validates the tag target, independently checks the candidate manifest, validates evidence, ProductCode, version, size, SHA-256, repository, workflow run, and attempt, then rechecks the remote tag immediately before publication.
- It generates `SHA256SUMS.txt` after all final assets exist, verifies the file, uploads it, and only then changes the release from draft to public.
- A failure at any point leaves the release draft. The workflow does not rebuild or substitute the MSI.
- A second invocation after publication fails safely instead of mutating the public release.

## Automation Contract

Repository automation rejects:

- a release-stage upload without explicit `draft: true`;
- final checksum or publication logic in the staging workflow;
- a promotion workflow without manual dispatch, draft-state guard, exact gate invocation, all-assets checksum generation, and final promotion after validation;
- an unapproved GitHub Action family or version;
- any path that can publish before the release-gate dependency succeeds.
