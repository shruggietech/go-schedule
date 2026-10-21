# S109 validation quickstart

Use only scoped fixtures and recording adapters from S108. Run from the repository root using a verified hidden launcher on Windows.

1. Run the new focused regressions in internal/cli, internal/bundle, internal/enrollment, internal/watcher and internal/api/server before and after each fix. Expected failures before the fix demonstrate the specific defect, not an unrelated compile or fixture error.
2. Regenerate initialized command metadata through the existing inventory command/test mechanism. Preserve case IDs and extend negative cases for newly expressible endpoint flags and binding conflicts.
3. Run sh scripts/verify.sh all in the foreground and watch all eight gates complete.
4. Use scripts/cli-qa.ps1 to prepare fresh receipt-bound inputs, capture Windows and the existing selected WSL2 environment, review every frozen contextual message with the maintainer prompt, and finalize/validate/render the canonical report. Existing historical S108 reports remain immutable.
5. Record exact source/binary identities, issue-by-issue outcomes, report totals and any separate new defects in verification.md. Do not relabel failures as passes or use a prior receipt after source changes.

The local done gate is followed by a local feat(109) commit and the requested pre-publication review. No branch push, PR, merge or release is part of this guide's authorization.
