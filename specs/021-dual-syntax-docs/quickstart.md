# Quickstart: Validate S021

```sh
go test ./internal/cron ./internal/scheduleinput ./internal/cli
sh test/scripts/docs-policy-check_test.sh
sh scripts/docs-check.sh
go run ./cmd/go-schedule task add --help
go run ./cmd/go-schedule task edit --help
sh scripts/verify.sh all
```
