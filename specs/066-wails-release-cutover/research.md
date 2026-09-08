# Research: Wails Release Cutover

## Decision 1: Preserve `gosched-gui` as the external application name

**Decision**: Configure the production Wails build to emit `gosched-gui` on every platform.

**Rationale**: The CLI launcher, Windows installer upgrade identity, shortcuts, macOS bundle executable, Linux desktop entry, documentation, and existing user habits already depend on this name. Keeping it avoids a second migration that adds no user value.

**Alternatives considered**: Rename all surfaces to `go-schedule` (unnecessary compatibility break); build `go-schedule` and copy it to `gosched-gui` (creates two identities and weakens artifact inspection).

## Decision 2: Build the desktop natively in the release matrix

**Decision**: Use Wails on native Windows, macOS, and Linux runners, then assemble each existing distribution format around the native output and root-module daemon and CLI binaries.

**Rationale**: Wails depends on platform webviews and native packaging. Native runners already prove the production desktop builds, and the release workflow already has the correct platform matrix.

**Alternatives considered**: Cross-compile Wails from Linux (unsupported and cannot honestly validate native webview integration); publish a browser-only frontend (does not satisfy the desktop contract).

## Decision 3: Keep the desktop in its own Go module

**Decision**: Retain `desktop/go.mod` and remove Fyne only from the root module.

**Rationale**: The separation keeps daemon and CLI builds cgo-free and prevents native Wails dependencies from expanding the root server dependency surface. Release jobs can build both modules explicitly.

**Alternatives considered**: Fold Wails into the root module (unnecessarily couples server and desktop dependencies); move daemon and CLI into the desktop module (reverses dependency direction and duplicates core code).

## Decision 4: Retire both Fyne and the S060 executable proof

**Decision**: Remove `gui/`, `cmd/gosched-gui/`, and `experiments/wails-foundation/`, while preserving their design and verification history under `specs/` and `CHANGELOG.md`.

**Rationale**: Maintaining three desktop implementations after cutover would invite drift, duplicate vulnerability updates, and ambiguous contributor guidance. The production Wails module now covers the proof's purpose.

**Alternatives considered**: Keep Fyne indefinitely behind no release path (still burdens tests and dependencies); keep the Wails proof in CI (duplicates production evidence); archive source in-tree (Git history already provides an immutable archive).

## Decision 5: Treat candidate qualification as failure-closed automation

**Decision**: Extend repository automation checks to require the production desktop jobs and release workflow fragments, reject Fyne build/dependency residue, and tie hosted results to the PR head revision.

**Rationale**: Deleting the old implementation is safe only if a missing platform, missing payload, or stale evidence is mechanically visible. Existing CI and PR checks naturally bind results to one commit.

**Alternatives considered**: Rely on a prose checklist alone (easy to drift); use previous slice evidence (not the exact cutover candidate); trigger a draft release during PR validation (publication-adjacent and unnecessary).

## Decision 6: Separate package readiness from public release

**Decision**: Update and validate release automation without creating a tag, draft release, public release, milestone closure, or public availability claim in S066.

**Rationale**: The operator authorized a pull request, not a release. The merge ritual is the appropriate boundary for accepting the cutover before release operations begin.

**Alternatives considered**: Publish a release candidate automatically (outside authorization); leave release automation unchanged until after merge (would make cutover unreviewable).
