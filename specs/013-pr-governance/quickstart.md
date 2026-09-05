# Quickstart: Lightweight Pull-Request Workflow

## Local validation

```sh
sh scripts/verify.sh all
git diff --check
git status --short
```

Confirm the diff contains only Feature 013 specifications and synchronized governance guidance. It must not contain product code, dependency, CI workflow, release, CodeQL, or repository-settings artifacts.

Observed 2026-08-26: all eight aggregate gates pass in the foreground. Core coverage is engine 88.1%, schedule 91.5%, timezone 88.9%, store 87.0%, catchup 87.5%, and logbus 91.1%. The docs check reports 11 clean pages. The diff contains only the files named by the plan and no hosted-settings artifact.

## Publication after authorization

1. Push `codex/013-pr-governance`.
2. Open a pull request targeting `main` with `Closes #23`.
3. Let existing CI and third-party AI reviewers report naturally.
4. Consider each review comment and reply with the change made or the reason it is not warranted.
5. Leave final review and merge to the maintainer.

No branch protection, required approvals, fixed check contexts, tag, release, or settings mutation is part of this feature.

## Review follow-up (2026-08-27)

- Clarified that PR publication authorization covers verified, in-scope review-fix pushes to the same open PR without another halt.
- Synchronized the root README contributor entry point with the PR-first model.
- Re-ran all eight aggregate gates successfully after both corrections.
