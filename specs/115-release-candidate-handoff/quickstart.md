# S115 maintainer quickstart

1. Run `go run ./scripts/release-source-check -root . -repository shruggietech/go-schedule -tag v1.6.0` on the reviewed source. Run existing comparison checks separately using their actual source identities.
2. Run focused release preflight and IPC regressions under race, then the original six-writer fixture once per declared Windows/Linux experiment with opt-in diagnostics and unchanged assertions. Keep original failure output if anything fails.
3. Run `sh scripts/verify.sh all` in the foreground, preserving every canonical gate. On this workstation the previously verified private policy bridge may delegate only shell-policy checks to native WSL; report that platform split explicitly.
4. Freeze the candidate source after canonical acceptance. Use `scripts/cli-qa.ps1` build and capture on fresh ignored input/run directories with an explicitly selected existing WSL2 distribution and unelevated user, the supported two-hour total budget and original per-case/resource fields. Follow `docs/maintainers/cli-qa-run-v1.md` and review every actual message before finalize/validate/render. Preserve command exits separately from observations and cleanup.
5. Shared summaries name source, receipt and executable hashes, platform/runtime identities, current outcomes and unresolved historical causes. Later documentation is a distinct source change and never relabels the captured subject.

After owner merge and separate tag authorization, the existing Release workflow may stage the exact immutable tag. That future operation and promotion are outside S115. Do not create a tag or release to follow this quickstart.
