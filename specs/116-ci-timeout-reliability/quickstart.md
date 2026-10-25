# S116 validation guide

1. Inspect contracts/timeouts.md and run `go run ./scripts/ci-timeout-check` after implementation.
2. Run focused frontend identity and configured-timeout regressions, including one controlled delayed experiment beyond the old 5s whole-test cap.
3. Run native six-writer diagnostics and focused real CLI/owned IPC metadata cases with the selected supported budget; compare exact receipts without changing the S115 frozen report.
4. Run `sh scripts/verify.sh all` through the verified hidden Windows launcher and await all eight gates. The portable automation portion may use the existing private WSL bridge; identify that platform split.
5. Run source/comparison checks for pending v1.6.0 and preserve public v1.5.1. No tag, staging or promotion is performed.
6. After commit/push, audit exact-head hosted checks and up to two Codex review rounds before owner merge.
