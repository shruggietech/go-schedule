# Pull request

Every change to `main` uses this pull-request workflow, whether it comes from a
maintainer, an automation agent, or an outside contributor.

If you are about to propose something substantial, please open an issue first.
Every feature on this project is specified through
[Spec Kit](https://github.com/github/spec-kit) before it is built, and a large
change that arrives without a spec is likely to need reworking rather than
merging.

## What this changes

<!-- What the change does and why. Remove the next line if no issue is closed. -->

Closes #

## Verification

Run the canonical CI-parity aggregate in the **foreground**, watch it finish,
and paste its result. The gates are documented in
[CONTRIBUTING.md](../CONTRIBUTING.md).

```sh
sh scripts/verify.sh all
```

```text

```

- [ ] `sh scripts/verify.sh all` passed locally, or any unavailable prerequisite
      is disclosed above.

## Checklist

- [ ] No safety-critical test surface was weakened or skipped: clock
      injection, timezone and DST resolution, store migrations, restart and
      catch-up recovery, goroutine termination, IPC access control.
- [ ] No pinned artifact changed, or: it did, and `CHANGELOG.md` carries a
      dated decision entry explaining why.
- [ ] Documentation updated where behavior changed.
