# S133 Validation Guide

Use hidden, noninteractive launchers and fresh owned output paths. No installed service or visible GUI operation is required.

1. Run the existing automation fixture suite before and after the #424 correction, preserving the failing-first result.
2. Resolve root and desktop graphs with -mod=readonly and verify checksums. Run frontend/site installation/audit via their lockfiles.
3. Run go run ./scripts/release-source-check -root . -repository shruggietech/go-schedule -tag v2.0.1 and the existing comparison checker for the candidate source identity.
4. Run sh scripts/verify.sh all in the foreground, then relevant existing desktop/site browser contracts and current-candidate CLI QA.
5. Compare protected historical hashes and quiescent source manifests, record exact results in verification.md, and finalize factual delivery records.

An official PR with complete exact-head checks and resolved findings is this execution's handoff. Public availability and attended native behavior remain separately attributable.
