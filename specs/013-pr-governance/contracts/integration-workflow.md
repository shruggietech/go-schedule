# Contract: Lightweight PR Integration

1. Synchronize local `main` and create a review branch.
2. Complete Spec-Kit work when the slice is feature-shaped.
3. Run `sh scripts/verify.sh all` and commit locally.
4. Halt before pushing the branch or opening the pull request.
5. After authorization, push and open a PR targeting `main`.
6. Use `Closes #N` only when the PR completes the issue; otherwise use `Refs #N`.
7. Consider each third-party review comment and respond with either a warranted change or a concise rationale.
8. Push verified, in-scope review fixes to the same branch under the PR's publication authorization. Material scope expansion or another PR requires separate authorization.
9. Leave final review and merge to the maintainer.
10. After merge, delete the merged branch and synchronize local `main`.

This contract does not require branch protection, approvals, resolved conversations, a fixed hosted-check set, or repository settings changes.
