# Quickstart: Review the Remote Access Architecture

## 1. Confirm current behavior is unchanged

```sh
git diff --stat origin/main...HEAD
git diff origin/main...HEAD -- internal cmd desktop go.mod go.sum desktop/go.mod desktop/go.sum
```

The second command must show no runtime implementation or dependency change. S074 defines the reviewed boundary only.

## 2. Validate the architecture contract

```sh
sh scripts/remote-architecture-check.sh .
sh test/scripts/remote-architecture-check_test.sh
```

The positive check must accept the maintained architecture page. The fixture suite must prove that removing the primary transport, any deployment mode, a control owner, a threat test, a non-goal, or the issue sequence fails.

## 3. Review issue traceability

```sh
gh issue view 165 --repo shruggietech/go-schedule
for issue in 166 167 168 169 170 171 172 173; do gh issue view "$issue" --repo shruggietech/go-schedule --json number,state,title; done
```

Issue #165 and every downstream issue must remain open until their own acceptance criteria and merge evidence are complete.

## 4. Run documentation and lifecycle gates

```sh
sh scripts/docs-check.sh
sh scripts/spec-lifecycle-check.sh .
go run ./scripts/github-format
```

## 5. Run canonical verification

```sh
sh scripts/verify.sh all
```

All eight gates must pass in the foreground. S074 adds no listener, no runtime dependency, and no public artifact.
