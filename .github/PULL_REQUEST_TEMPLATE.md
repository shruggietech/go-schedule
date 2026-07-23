# Pull request

**A note on how this repository works.** go-schedule is developed trunk-based:
maintainer work lands directly on `main`, with a single pre-push review as the
human gate. There is no internal pull-request queue, so this template exists for
outside contributions — which are welcome, and which will be reviewed here.

If you are about to propose something substantial, please open an issue first.
Every feature on this project is specified through
[Spec Kit](https://github.com/github/spec-kit) before it is built, and a large
change that arrives without a spec is likely to need reworking rather than
merging.

## What this changes

<!-- What the change does and why. Link the issue it closes, if any. -->

## Verification

Paste the result of the CI-parity gates. They are documented in
[CONTRIBUTING.md](../CONTRIBUTING.md); run them in the **foreground** and watch
them finish, because a backgrounded `go test` cannot be told apart from a dead
one.

```text

```

- [ ] `gofmt -l internal cmd test` printed nothing
- [ ] `go vet ./...` clean
- [ ] `golangci-lint run ./...` clean
- [ ] `go test -race` green (or: the race gate needs a C toolchain I do not
      have, and I am saying so rather than reporting the suite as passing)
- [ ] `go test ./gui/...` green
- [ ] `sh scripts/coverage-gate.sh` green

## Checklist

- [ ] No safety-critical test surface was weakened or skipped — clock
      injection, timezone and DST resolution, store migrations, restart and
      catch-up recovery, goroutine termination, IPC access control.
- [ ] No pinned artifact changed, or: it did, and `CHANGELOG.md` carries a
      dated decision entry explaining why.
- [ ] Documentation updated where behavior changed.
