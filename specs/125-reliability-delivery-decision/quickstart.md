# S125 Quickstart

Use the repository's verified hidden launcher on Windows; no visible console or headed browser. Keep local browser workers at one.

1. Run the deterministic adverse-order browser regression against the original/fixed fixture and preserve both outcomes.
2. Run `node node_modules/@playwright/test/cli.js test --workers=1` from desktop/frontend through a verified hidden runner.
3. Run `sh scripts/verify.sh all` in the foreground through the existing hidden launcher and watch it complete.
4. Review `docs/distributed-scheduling.md` against the cited source evidence; unsupported actor/scale/availability fields must remain unknown, and #187/#188 must remain gated unless their prerequisite decision is supported.
5. Follow `handoff.md` only up to the current authorization boundary. No tag, draft creation or promotion runs during S125 source work.
6. Check GitHub formatting, UTF-8/no-BOM, lifecycle inventory, source/evidence integrity and final exact-head CI/reviews before owner handoff.
