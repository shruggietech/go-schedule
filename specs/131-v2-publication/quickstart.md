# S131 Execution and Readback

1. Read spec.md, plan.md, tasks.md and the existing S129 handoff. Inspect gh release/list, immutable tag state and exact-main CI without mutation.
2. Run go run ./scripts/release-source-check -root . -repository shruggietech/go-schedule -tag v2.0.0 and go run ./scripts/comparison-check -tag v2.0.0 -commit 75eae1ae11ad40db38341dc7d63f32876cba6483 through the verified hidden runner. Expected: valid source/comparison, explicitly no tag/package acceptance claim.
3. Only after exact approval, create/push the immutable tag; watch Release to completion. Retain actual run/attempt/release/asset inventory and inspect downloaded archive/MSI/manifest identities.
4. Only after separate concrete-inventory approval, dispatch gh workflow run promote-release-waived.yml --ref main -f tag=v2.0.0. Retain actual result and independently verify public/latest nine files/eight checksums and original asset IDs/digests.
5. Read Pages deployment and actual HTTPS /, /compare/ and /docs/install/windows/; reconcile site/data/published.json and current wording. Run sh scripts/verify.sh all in the foreground on the final source tree.
6. Commit and automatically push/open one official PR, monitor terminal CI and every review comment/reaction with a two-request cap, then hand off unmerged.
