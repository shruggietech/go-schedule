# Verification: Remote MCP Authorization

## Local evidence

- `go test -race ./internal/remotemcp ./internal/remote ./internal/config ./internal/enrollment ./internal/api/client`
- `go run ./scripts/github-format .`
- `scripts/spec-lifecycle-check.sh .`
- `scripts/verify.sh all`

The focused matrix covers independent enablement, canonical HTTPS resource validation, RFC 9728 and authorization-server metadata, exact client and resource binding, MCP-only source credentials, scope ceilings, header-only bearer handling, 0/3/9 authority discovery, official SDK Streamable HTTP negotiation, persistent actor delegation, token expiry, credential rotation and revocation, and daemon identity reset. The canonical aggregate passed all eight gates: format, vet, lint, race, GUI, coverage, docs, and automation.

## Hosted evidence

Hosted CI and external Codex and security review evidence are recorded on the S093 pull request before the maintainer merge ritual.
