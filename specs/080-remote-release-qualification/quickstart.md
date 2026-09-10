# Quickstart: Verify S080

Run focused service and package lifecycle tests:

```sh
go test -race ./internal/cli ./test/integration -run 'TestServiceInstall|TestV14RemoteRelease' -count=1
```

Run the detailed remote qualification packages:

```sh
go test -race ./internal/config ./internal/remote ./internal/enrollment ./internal/authorization ./internal/api/client ./internal/cli ./internal/clientprofile
```

```sh
cd desktop && go test -race ./connection ./connections
```

Validate documentation and publication formatting:

```sh
sh scripts/docs-check.sh
go run ./scripts/github-format
```

Run the canonical foreground suite:

```sh
sh scripts/verify.sh all
```

Review the three-platform `v1.4 remote qualification` pull-request jobs before merge. Do not tag or publish a release from this slice.
