# S118 validation guide

Prerequisites: Existing native Windows toolchain and verified hidden launcher; existing private WSL portable-policy bridge. No provisioning or installed-system changes.

1. Compare S117 spec, task dispositions and inventory with PR #377. Confirm reviewed head and squash merge remain distinct and original failed reports retain their hashes/outcomes.
2. Read [handoff.md](handoff.md). Confirm all three issue outcomes, current procedure links, pending/public version distinction and later owner decisions agree.
3. Run existing lifecycle and GitHub-format checks. Run `go run ./scripts/release-source-check -repository shruggietech/go-schedule -tag v1.6.0` and `go run ./scripts/comparison-check -tag v1.6.0 -commit 9ad13821c825d01efcc038b6b1634b187a854758` as read-only content guards. These do not establish an actual tag or staged artifact.
4. Run all eight foreground gates via `scripts/verify.sh all` with the existing Windows/WSL bridge; preserve source-input hashes across the run and record actual outcomes in verification.md.
5. Audit UTF-8/no-BOM, mojibake, current links and retained report/content hashes after factual lifecycle finalization. Publish the owner-authorized PR and record exact-head CI/review receipts externally, stopping before merge/tag/promotion.
