# Research: v1.3 Notifications and Local Agent Access Qualification

## Decision 1: Qualify the merge candidate without publishing a release

**Decision**: Treat the reviewed commit and its hosted pull-request checks as the v1.3 qualification boundary. Do not create a tag, GitHub release, or public artifact in S072.

**Rationale**: Issue #190 asks for release qualification, while repository governance reserves release publication for a separate explicit ritual. Pull-request evidence can prove the implemented contract without making an irreversible release claim.

**Alternatives considered**: Publishing a draft release was rejected because it expands scope and creates a public artifact. Tagging a prerelease was rejected because tags require separate authorization and are unnecessary for behavioral qualification.

## Decision 2: Exercise a built daemon as a package-shaped candidate

**Decision**: Build `goschedd` into an isolated distribution-like directory, launch it with an isolated configuration and IPC endpoint, inspect state through the supported client, stop it, and launch the same candidate again against retained storage.

**Rationale**: This verifies process wiring, configuration defaults, migrations, local IPC, and restart behavior together. It is stronger than calling internal constructors and remains portable across all hosted operating systems.

**Alternatives considered**: Unit-only verification was rejected because it cannot prove executable wiring. Installing the native package on every runner was rejected because the repository currently has only a Windows installer contract and package installation would mix platform packaging work into a behavioral release gate.

## Decision 3: Compose existing detailed suites into the release gate

**Decision**: The named matrix will run the new default-state journey and existing focused tests for webhook delivery, migration preservation, MCP official-SDK interoperability, HTTP authorization, redaction, hostile content, and lifecycle behavior.

**Rationale**: S067 through S071 already contain detailed tests close to their owning packages. Reusing those tests preserves a single behavioral authority while making their combined release significance visible.

**Alternatives considered**: Reimplementing every case in one large end-to-end test was rejected because it would duplicate assertions, slow diagnosis, and increase flake risk. Running only the repository-wide race gate was rejected because it does not expose a stable, named v1.3 qualification signal.

## Decision 4: Make every supported platform result explicit

**Decision**: Use a required three-entry GitHub Actions matrix with fail-fast disabled and no platform conditionals or silent skips. Missing prerequisites fail the relevant job.

**Rationale**: A discrete job per operating system gives reviewers an explicit pass or failure for every supported platform and satisfies the issue's visibility requirement.

**Alternatives considered**: Cross-compilation was rejected because it cannot exercise runtime behavior. A single Linux job was rejected because it cannot qualify Windows named pipes or macOS execution. Optional jobs were rejected because skipped evidence could be mistaken for success.

## Decision 5: Protect qualification wiring as a repository contract

**Decision**: Extend `scripts/automation-check.sh` to require the named job, all three matrix entries, and the qualification test command.

**Rationale**: The repository already uses this script to prevent accidental erosion of pinned workflow behavior. Adding the new fragments makes loss of a platform or test command fail locally and in CI.

**Alternatives considered**: Relying on manual review was rejected because workflow drift is easy to miss. Adding a new workflow parser dependency was rejected because exact fragment checks are sufficient for this bounded contract.

## Decision 6: Record an explicit non-release boundary

**Decision**: Update the slice changelog and verification record to state what was qualified and that no release tag or public artifact was produced.

**Rationale**: Release readiness and release publication are different states. Durable wording prevents a green qualification PR from being misread as a shipped release.

**Alternatives considered**: Adding v1.3 release notes was rejected because release notes would imply a publication step outside this slice.
