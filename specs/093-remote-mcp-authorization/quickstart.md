# Quickstart: Remote MCP Authorization

## Configure the daemon

Enable the existing remote HTTPS listener and independently enable its MCP resource with the exact public URL clients will use:

```json
{
  "remote": {
    "enabled": true,
    "bind_address": "10.0.0.2:8443",
    "certificate_file": "/etc/go-schedule/tls/server.crt",
    "private_key_file": "/etc/go-schedule/tls/server.key",
    "mcp": {
      "enabled": true,
      "resource_url": "https://scheduler.example.internal:8443/mcp"
    }
  }
}
```

Create a pairing for kind `mcp` and the intended capability through protected local administration, exchange it once, and retain the returned credential ID as OAuth `client_id` and returned credential token as OAuth `client_secret`. The pairing phrase itself is never an MCP bearer or OAuth client secret.

## Discover authorization

Request the protected resource metadata URL advertised by an unauthenticated `/mcp` challenge. Follow its authorization server entry to the authorization-server metadata document.

## Obtain an access token

Use HTTP Basic client authentication and an exact resource indicator:

```text
grant_type=client_credentials
resource=https://scheduler.example.internal:8443/mcp
scope=mcp:observe
```

Use the returned short-lived bearer only with the exact `/mcp` resource. Obtain a new token after expiry, restart, credential rotation, or revocation.

## Verify

Connect with an official Streamable HTTP client, negotiate the current protocol, list resources and tools, and confirm the discovered mutation surface matches the issued scope.
