# Research: Windows Release Candidate Gate

## Decision 1: Stage a draft release, then promote it

**Decision**: Keep the tag-triggered workflow, but make every asset upload explicitly target a draft release. Add a manually dispatched promotion workflow that validates the attended evidence and exact draft assets before changing the release to public.

**Rationale**: The current workflow publishes Linux and macOS archives before the Windows MSI exists. Attended proof cannot precede the MSI produced by that workflow, and a later rebuild cannot be assumed byte-identical. GitHub draft releases provide a durable, non-public staging surface whose assets can be tested and later promoted without rebuilding.

**Alternatives considered**:

- Build in CI, commit an approval record, and rebuild on the release tag. Rejected because WiX and embedded build output do not provide a repository contract of reproducible byte identity.
- Require a self-hosted interactive Windows 11 Actions runner. Rejected because the repository has no such runner and unattended Actions execution still cannot provide credible attended visual observations.
- Rely on a GitHub Environment approval. Rejected because approval alone does not bind the approved MSI bytes or prove scenario completeness.

## Decision 2: Validate evidence with one cross-platform Go command

**Decision**: Add `scripts/windows-release-gate`, backed by a focused `internal/releasegate` package, to validate candidate identity, scenario coverage, outcomes, measurements, attachment integrity, and the exact MSI digest. The same command runs locally, in canonical automation, and in promotion.

**Rationale**: A single standard-library implementation avoids semantic drift between PowerShell collection and Linux promotion. It is testable through table-driven and mutation tests and follows the existing repository-tool pattern established by `scripts/brand-check`.

**Alternatives considered**:

- Validate entirely in PowerShell. Rejected because the promotion runner is Linux and maintaining a second validation path would weaken the single-contract requirement.
- Validate with shell and `jq`. Rejected because deeply nested schema and cross-observation rules become brittle and difficult to unit test.
- Introduce JSON Schema tooling. Rejected because semantic rules such as display-class diversity, time windows, candidate identity, and scenario-specific measurements still need custom logic and would add a dependency.

## Decision 3: Use a versioned manifest with fixed scenario identifiers

**Decision**: One `evidence.json` contains schema version 1, exact candidate identity, environment records, and exactly one result for every fixed scenario identifier. Every result is one of `pass`, `fail`, `unavailable`, `skipped`, `timed-out`, or `partial`; only `pass` satisfies readiness.

**Rationale**: Fixed identifiers make omissions, duplicates, and substitutions detectable. Explicit non-pass states preserve the difference between failure and unavailable infrastructure without allowing either to pass.

**Alternatives considered**:

- Free-form checklist names. Rejected because spelling changes can conceal missing evidence.
- Boolean pass fields. Rejected because they cannot distinguish failure, cancellation, unavailable infrastructure, timeout, and partial cleanup.
- Multiple independent fragments as the promotion input. Rejected because aggregation rules and duplicate precedence would be ambiguous. Collector fragments may exist during the run, but final validation consumes one canonical manifest.

## Decision 4: Combine native measurements with attended visual evidence

**Decision**: Add an attended PowerShell collector that captures the exact GUI PID, HWND, client and outer rectangles, monitor and work-area rectangles, effective DPI, visibility, and placement state. Required visual observations also retain timestamped screenshots and operator surface counts as hashed attachments.

**Rationale**: Native API measurements are necessary for sizing and state. Fyne modal dialogs are canvas overlays, not separate top-level HWNDs, so native enumeration alone cannot count all visible error surfaces. Structured operator counts plus screenshots cover that gap without mislabeling headless widget state as native proof.

**Alternatives considered**:

- Screenshots alone. Rejected because they cannot prove logical content size, DPI, monitor work area, process identity, or state flags.
- HWND enumeration alone. Rejected because Fyne overlays share the main HWND.
- Add shipping telemetry solely for the gate. Rejected unless native evidence later proves derived client sizing is insufficient; S040 avoids changing production behavior speculatively.

## Decision 5: Keep product corrections outside the gate unless reproduced

**Decision**: Do not change GUI sizing, monitor choice, or connection-error behavior in S040 unless the exact candidate reproduces a failure and the required baseline evidence is retained. Any recurring-error correction also requires passing native proof against the uncommitted worktree before commit.

**Rationale**: Existing unit and headless tests support the intended behavior, but repository analysis found an unproven mixed-monitor risk. A speculative correction would violate #94's proof-before-commit requirement and could introduce a new display regression.

## Decision 6: Retain evidence as a release asset

**Decision**: Upload a canonical evidence archive to the draft release. Promotion validates the archive and includes it in final `SHA256SUMS.txt` before publication.

**Rationale**: Release assets keep candidate, proof, and final checksums together and avoid depending on an undeclared external store. The evidence remains reviewable and publicly auditable after promotion.

**Alternatives considered**:

- Commit screenshots and logs to Git. Rejected because binary evidence would bloat history and the evidence is specific to one staged artifact.
- Use transient Actions artifacts. Rejected because expiration can break delayed promotion and public release auditability.
- Store only an approval digest in Git. Rejected because reviewers and users would lose the underlying evidence bundle.

## Repository Findings

- `test/windows/Invoke-InstallerLifecycle.ps1` already verifies Windows 11 client identity, MSI hash/version/ProductCode, LocalSystem service state, normal-user installed access, and real manual/scheduled task execution. S040 should preserve and reference those results rather than replace them.
- S039's hosted Windows job proves compiled MSI relationships and silent shortcut, maintenance, preserve, wipe, locked-cleanup, repair, and upgrade behavior. It remains supporting evidence, not attended proof.
- Current GUI startup requests 1280 by 800 logical content and caps each dimension at 90 percent of a selected work area. A mixed-monitor risk exists because sizing selects the pointer monitor while pre-show centering may select the primary monitor; native evidence must determine whether it is real.
- Connection incidents are deduplicated in one in-frame card in application state, but native proof must observe Fyne overlays and top-level windows for two minutes per induced condition.
- The repository has no self-hosted runner or release environment and main is not branch-protected. Workflow controls govern the supported release path but cannot prevent a repository administrator from bypassing GitHub manually.
