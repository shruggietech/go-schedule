---
title: Remote access architecture
nav_order: 8.6
---

# Remote access architecture

**Current status:** Implemented behind explicit daemon configuration. Remote access remains disabled by default.

**Primary remote transport:** HTTPS with versioned HTTP/JSON.

This page defines the v1.4 remote-access boundary and its operator contract. The daemon always retains its protected local Unix socket or Windows named pipe and remains fully usable offline. It opens the separate network listener only when an operator supplies complete remote configuration and explicitly enables it.

## Design posture

The daemon continues to own scheduling, persistence, execution, authorization decisions, and audit creation. The local IPC mux and remote operation allowlist are separate adapters over shared daemon operations. A route being safe and authenticated over local OS-controlled IPC does not make it safe for remote exposure.

The remote adapter is deny-by-default. Each admitted operation must have a stable OpenAPI operation identifier, minimum capability, input bound, audit class, retry class, secret exclusions, local-equivalence mapping, and acceptance tests. Local-only routes are absent from the allowlist, including runtime filesystem paths, localhost MCP lifecycle and credential issuance, external-trigger secret reveal and rotation, and every future route without a completed remote record.

The design adds no custom encryption, signed token, wire protocol, user directory, or general policy engine. It uses operating-system trust and credential stores, Go's standard cryptographic and HTTP packages, and established narrowly scoped libraries.

## Ordered request boundary

Every remote connection or request crosses these controls in order:

1. The operator has explicitly enabled remote access and supplied an exact bind address.
2. Go's TLS server completes a TLS 1.3 handshake with an operator-trusted certificate.
3. A coarse source limiter bounds unauthenticated connection and credential-guessing work.
4. The server accepts only a supported `/api/v1` path, bounded headers, and an expected method.
5. The request supplies a durable bearer credential, except for the isolated enrollment exchange.
6. Authentication loads current actor state and rejects unknown, expired, or revoked authority.
7. A second limiter bounds the authenticated actor.
8. The remote allowlist resolves the operation and checks its required capability.
9. Content type, OpenAPI schema, body size, collection size, and operation-specific constraints are validated.
10. The request receives a correlation identity and audit classification.
11. The adapter invokes the shared daemon application operation.
12. The server emits a bounded response and completes any required secret-free audit event.

Authentication and capability authorization run before daemon application behavior. A failure at any step prevents every later step. Remote authentication responses do not reveal whether a credential or enrollment phrase was unknown, malformed, expired, consumed, or revoked.

Browser origins are not a v1.4 client surface. Requests carrying an `Origin` header are rejected by default, and no cross-origin resource sharing response headers are emitted. A later browser-client proposal would require an explicit origin allowlist and its own review.

## Remote operation record

Every remotely exposed operation must declare all of these properties in or beside the OpenAPI source:

| Property | Required contract |
| --- | --- |
| Operation identity | Stable and unique within one API major |
| Method and path | Under `/api/v1`; enrollment is isolated within that major |
| Local equivalent | Shared daemon operation, or a reasoned discovery or enrollment exception |
| Capability | Exactly Observe, Operate, Manage, or Enroll |
| Input bound | Header and body bytes plus collection, pagination, and duration limits |
| Audit class | None, protected read, mutation, or enrollment |
| Retry class | Safe read, refreshable stream, or no automatic replay |
| Secret contract | Protected request fields and excluded response, log, event, and audit fields |
| Verification | Authentication, authorization, schema, bounds, equivalence, and failure tests |

The implemented initial allowlist is defined by `api/openapi/remote-v1.yaml`. No prose statement on this page exposes an endpoint by itself.

## Capability model

Observe, Operate, Manage, and Enroll are the complete capability vocabulary.

| Capability | Permitted intent | Explicit exclusions |
| --- | --- | --- |
| Observe | Read safe scheduler state, history, alerts, and live activity | Executable configuration, secret reveal, task execution, mutation |
| Operate | Observe plus deliberate task execution and non-configuration operational acknowledgement | Task, group, automation, credential, or notification configuration |
| Manage | Operate plus scheduler object and notification-policy administration | Actor authority, enrollment, credential issuance, and server exposure configuration |
| Enroll | Manage plus pairing phrase creation, actor capability assignment, credential rotation, and revocation | Federation, arbitrary policy, certificate-authority operations |

The classes are monotonic for product comprehension, but every operation still declares one minimum class. The daemon loads current server-owned authority for every request. Capability data is never trusted from a bearer value, a forwarded header, a display name, or client-supplied claims.

Local IPC authorization remains independent. A local OS caller does not receive an implicit bearer credential, and enabling HTTPS cannot widen Unix-socket ownership or Windows named-pipe access.

## Credential lifecycle

An enrollment phrase is short-lived and single-use. It exists only to bootstrap one named client relationship, is visibly tied to the intended daemon identity, has a strict expiry and attempt budget, is retained only through a salted Argon2id verifier, and cannot be used as an ordinary API credential. Successful exchange atomically consumes the phrase and issues unrelated durable material. Cancellation, expiry, exhaustion, or success permanently ends that phrase.

A durable credential is an opaque 256-bit bearer value retained only as a digest. Go's `crypto/rand` creates the raw value, base64url is its opaque transport encoding, RFC 6750 `Authorization: Bearer` syntax carries it, SHA-256 stores a verifier, and constant-time comparison checks it. The token contains no claims, authority, identity, expiry, or signature. Server-owned actor and credential records provide those facts so revocation takes effect immediately.

Each client installation receives an independent credential and safe fingerprint. Issuance returns the raw value once. The desktop and named CLI profiles store it through supported operating-system credential storage, never in application JSON, process arguments, shell-history examples, logs, exports, diagnostics, or screenshots. Rotation replaces the verifier atomically and invalidates the prior value. A lost credential is revoked and replaced, never recovered or re-derived from an enrollment phrase.

The daemon revalidates actor and credential state periodically during a live stream. Revocation, expiry, or actor disablement terminates authority for new requests immediately and closes existing streams within the documented revalidation bound defined by #167 and #168.

## Client connection profiles

The desktop and CLI share one versioned user-scoped `profiles.json` document under the go-schedule desktop configuration directory. Each record contains a random profile ID, local label, canonical HTTPS origin, pinned daemon installation ID, native credential reference, trusted certificate PEM and SHA-256 fingerprint, client kind, granted capability, safe daemon display and platform facts, and timestamps. The document contains no bearer value or pairing phrase and is replaced atomically under an exclusive cross-process lock.

The desktop retains one active profile ID and restores only that exact target after restart. This computer always remains available through protected local IPC. Selection loads the bearer from native operating-system credential storage, constructs one immutable TLS client, verifies the manifest installation ID, then switches the request router and event generation together. A failed lookup, credential read, certificate check, or identity comparison never falls back to a different remote target.

The Connections workspace lists same-named profiles with their endpoint and shortened daemon ID, supports local rename, accepts fresh pairing material for same-identity repair, and requires confirmation before removal. Repair replaces the profile only after the new identity and credential are valid, then deletes the superseded native credential. Removing the active profile switches to This computer first; any credential deletion failure keeps the profile metadata visible.

The CLI deliberately has no persistent active remote target. `--profile` selects one saved profile by exact ID or unambiguous label for one invocation. Alternatively, `--endpoint`, `--daemon-id`, `--credential-id`, and `--certificate-file` must all be supplied together. The bearer and one-time phrase are never command arguments. Human mode identifies the selected endpoint on stderr, while JSON stdout retains its existing machine-readable contract.

Private-network HTTPS is the preferred direct mode. For an SSH tunnel, keep application TLS enabled, forward a local port to the configured daemon listener, and create a profile whose HTTPS origin and certificate hostname match the tunneled endpoint. Do not add an insecure verification flag. If the certificate or daemon identity changes, stop and repair the profile from newly authenticated operator material.

For a daemon bound to loopback with a certificate containing the IP SAN `127.0.0.1`, a clean client can establish the tunnel and pair a CLI profile as follows:

```sh
ssh -N -L 127.0.0.1:8443:127.0.0.1:8443 operator@daemon-host
gosched profile pair production-tunnel --address https://127.0.0.1:8443 --expected-daemon-id DAEMON_ID --pairing-id PAIRING_ID --trusted-certificate daemon.pem --client-name "Tunnel CLI" --capability observe
gosched --profile production-tunnel health
```

## API compatibility and live updates

OpenAPI 3.1 is the source of truth for the remote contract. The repository pins `github.com/oapi-codegen/oapi-codegen/v2` as a Go tool and generates strict standard-library HTTP server and client boundaries. Authentication and authorization stay explicit middleware because generated routing does not implement product security policy. Generated files are committed, and a clean regeneration check prevents source, server, and client drift.

Local IPC retains `/v1`. Remote paths use `/api/v1` so network compatibility can evolve without silently changing local transport assumptions. Additive optional response fields, new operations, and new manifest capabilities are compatible within a major. Removing or renaming fields, changing their meaning, tightening previously accepted values, or changing success semantics requires a new path major.

The capability manifest reports stable daemon identity, product version, supported API majors, features, and safe compatibility facts before clients present actions. A new major overlaps the prior supported major for a documented migration interval defined by its delivery issue. Deprecation appears in OpenAPI, client diagnostics, documentation, and release notes before removal. This architecture invents no calendar deadline.

Server-Sent Events provide authenticated one-way live activity over HTTPS. Desktop clients treat each event as an invalidation hint, not a durable event log. After a disconnect they revalidate TLS trust, daemon identity, compatibility, and current credential authority through `/api/v1/access/current`, refresh mounted workspaces from authoritative reads, and then open a new stream. Mutations have no automatic replay when a disconnect leaves their outcome uncertain; the client immediately enters reconciliation, preserves user input, keeps mutation controls blocked until the authoritative refresh completes, and then requires a deliberate resubmission.

## Failure behavior

| Condition | Remote result | Required evidence |
| --- | --- | --- |
| Missing, malformed, unknown, or expired bearer credential | `401 unauthorized` and `WWW-Authenticate: Bearer` | Bounded secret-free security evidence, aggregated where necessary |
| Revoked bearer credential | `401 credential_revoked` and `WWW-Authenticate: Bearer` | Bounded secret-free security evidence that lets a holder repair the selected profile |
| Authenticated actor lacks capability | `403` and stable denial code | Actor-attributed audit denial |
| Unsupported API path or negotiated incompatibility | `404` for an unknown path or `409` for known incompatibility | Safe compatibility diagnostic |
| Invalid content type or schema | `400` or `415` with bounded field detail | Actor-attributed failed request without protected input |
| Header or body exceeds its limit | `413` or connection-level header rejection | Bounded security diagnostic |
| Source or actor rate is exceeded | `429` with bounded `Retry-After` | Aggregate limiter evidence without credential disclosure |
| Daemon conflict or domain validation fails | Existing stable JSON error envelope | Actor, operation, result, and target identity |
| Client loses a mutation response | No automatic replay | Authoritative refresh before explicit retry |
| Event stream disconnects | Revalidate, refresh authoritative state, and open a new stream | No state depends on complete event replay |

HTTP server configuration includes finite read-header, read, idle, and operation-aware write behavior, bounded maximum headers, bounded request bodies, graceful shutdown, and connection-lifecycle ownership. Streaming endpoints use heartbeat and maximum-idle rules instead of an ordinary fixed response timeout. Shutdown stops accepting connections, allows bounded in-flight operations to finish, terminates streams, and leaves local IPC lifecycle independently available until daemon shutdown.

## Deployment modes

All network modes require explicit enablement, an exact bind address, a certificate and key readable only by the daemon identity, and explicit exposure acknowledgement for wildcard or public addresses. Invalid configuration fails before bind. HTTPS is retained inside VPN, tunnel, and proxy deployments so the application boundary has one security contract.

| Mode | Listener and TLS | Product ownership | Operator ownership | Support posture |
| --- | --- | --- | --- | --- |
| Local IPC | Unix socket or Windows named pipe; no TCP | IPC permissions, local API, daemon lifecycle | OS account and group administration | Default and unchanged |
| Private-network HTTPS | Explicit private address with direct go-schedule TLS | HTTPS boundary, authentication, authorization, limits, audit | VPN or private routing, firewall, DNS, trusted certificate | Recommended remote mode |
| SSH-tunneled HTTPS | Explicit loopback or private address with direct go-schedule TLS inside SSH | HTTPS boundary and all application controls | SSH identity, tunnel lifecycle, port forwarding, trusted certificate | Recommended for headless administration |
| Reverse-proxied HTTPS | Explicit private address with direct go-schedule TLS; proxy adds public TLS | Application TLS configuration validation, authentication, authorization, trusted-proxy allowlist, limits, audit | Proxy lifecycle; public and backend certificate issuance and renewal; backend hostname trust and key permissions; routing, forwarding policy, firewall | Supported with explicit immediate-peer trust |
| Direct HTTPS | Explicit network address with direct go-schedule TLS | HTTPS boundary, authentication, authorization, limits, audit, fail-closed startup | Certificate issuance and renewal, key permissions, DNS, firewall, exposure monitoring | Supported for experienced operators, never automatic |

Forwarded headers are ignored unless trusted-proxy configuration explicitly matches the immediate peer. Trusting a proxy changes only attribution of client address and scheme metadata; it never delegates bearer authentication or capability authorization. A proxy cannot inject actor identity or bypass application limits.

Disabling remote access stops acceptance, drains bounded in-flight operations, closes streams, clears ephemeral limiter state, and preserves local IPC. Installation, upgrade, restore, and ordinary daemon restart never enable a listener implicitly.

## Operator runbook

### 1. Prepare configuration and TLS

Remote access is supported only with a complete JSON configuration and TLS certificate. A daemon started without `--config` checks the platform data directory for `config.json`: `C:\ProgramData\goschedule\config.json` on Windows, `/var/lib/goschedule/config.json` on Linux, or `/Library/Application Support/goschedule/config.json` on macOS. If the file is absent, built-in defaults apply and remote access stays disabled. Unix operators may instead retain another absolute path with `gosched service install --config <file>`; Windows MSI operators should use the default file because Windows Installer owns the packaged service definition. Relative data, IPC, log, certificate, and private-key paths inside the JSON file resolve from the file's directory.

The following configuration is an illustrative private-network example. Replace every address and path with values owned by the deployment. `remote.bind_address` must contain a numeric IP and port. The certificate must be valid for the hostname or IP used by clients, the daemon identity must be able to read both files, and the private key must be restricted to administrators and the daemon identity.

```text
{
  "remote": {
    "enabled": true,
    "bind_address": "10.0.0.20:8443",
    "certificate_file": "/etc/goschedule/tls/server.crt",
    "private_key_file": "/etc/goschedule/tls/server.key",
    "acknowledge_public_exposure": false
  }
}
```

Use a certificate issued by an organization-trusted or publicly trusted authority when possible. A private self-signed certificate is suitable only when the operator distributes and verifies that exact certificate through an authenticated channel. Certificate issuance, renewal, DNS, routing, firewall policy, key backup, and monitoring are operator-owned infrastructure; go-schedule does not automate them.

For a new Unix service installation, use the install command's pre-registration validation and then start the service:

```sh
sudo gosched service install --config /etc/goschedule/daemon.json
sudo gosched service start
gosched health
```

`service install --config` parses and validates the complete configuration before it changes the service definition. Do not use `goschedd --config` as a check-only command: it starts the scheduler, local IPC, and any configured remote listener. On Windows, place the reviewed file at `C:\ProgramData\goschedule\config.json`, then use an elevated PowerShell session to run `gosched service restart`; startup validation fails closed if that file is invalid.

### 2. Choose a deployment mode

Private-network HTTPS is the recommended direct deployment. Bind the daemon to its private address, restrict the port to approved client networks, and use a certificate valid for the client-facing name or address.

SSH-tunneled HTTPS is the recommended path for occasional headless administration. Bind go-schedule to loopback or a private address, retain application TLS, and forward a client port through an independently authenticated SSH connection. The certificate must still match the HTTPS endpoint selected by the client. There is no insecure verification flag.

Reverse-proxied HTTPS is an advanced supported deployment. The proxy must connect to the go-schedule backend over TLS, validate the backend certificate, preserve request bounds, and be the only immediate peer allowed by firewall policy. The operator owns both public and backend certificate lifecycles, proxy hardening, routing, forwarding policy, and monitoring. Forwarded headers do not confer actor identity or bypass bearer authorization.

Direct public HTTPS is an advanced supported deployment, never a zero-configuration recommendation. Before setting `acknowledge_public_exposure` to `true`, the operator must provide a publicly trusted certificate with automated renewal outside go-schedule, stable DNS, least-access firewall rules, exposure and certificate-expiry monitoring, protected local administrative access, and a tested credential-revocation response. Do not expose the listener merely to avoid configuring a VPN, SSH tunnel, or reverse proxy.

### 3. Pair one client

From a protected local shell on the daemon host, create the least-capable relationship the client needs:

```sh
gosched pairing create "Operations laptop" --kind desktop --capability operate
```

The output contains a daemon ID, pairing ID, one-time phrase, and expiration. Move those values and the trusted certificate through authenticated channels. Do not capture them in screenshots, tickets, shell scripts, process arguments, or committed files.

For the desktop, open **Settings**, then **Connections**, choose **Pair connection**, enter the HTTPS address and expected daemon ID, select the trusted certificate file, enter the pairing ID and phrase, and confirm the requested capability. The desktop verifies certificate trust and daemon identity before saving secret-free profile metadata and placing the bearer in native operating-system credential storage.

For the CLI, read the phrase through the command's protected prompt:

```sh
gosched profile pair production --address https://scheduler.example.internal:8443 --expected-daemon-id DAEMON_ID --pairing-id PAIRING_ID --trusted-certificate daemon.pem --client-name "Operations CLI" --capability operate
gosched --profile production health
gosched --profile production task list
```

Direct JSON clients perform `POST /api/v1/enroll` once, store the returned bearer in an operating-system or application secret store, verify `GET /api/v1/manifest`, and send `Authorization: Bearer <opaque-value>` only in request headers. The [OpenAPI document](https://github.com/shruggietech/go-schedule/blob/main/api/openapi/remote-v1.yaml) is authoritative. Browser origins, credential query parameters, credential cookies, redirects, and plaintext HTTP are unsupported.

Choose Observe for read-only state and history, Operate for deliberate runs and acknowledgements, Manage for scheduler-object configuration, and Enroll only for administrators who must manage actors and credentials. Create separate relationships for separate client installations.

### 4. Rotate, revoke, disable, and upgrade

List safe credential metadata locally, then rotate or revoke by ID:

```sh
gosched credential list
gosched credential rotate CREDENTIAL_ID
gosched credential revoke CREDENTIAL_ID
```

Rotation prints a new bearer once and immediately invalidates the old bearer. Update the intended client's native secret store without placing the value in a command argument. Revocation disables the credential and its actor relationship immediately; live streams are closed at their bounded revalidation interval.

To disable network access, set `remote.enabled` to `false` or remove the `remote` object, validate the file, and restart the service. Confirm `gosched health` succeeds locally and use an operating-system socket inventory to confirm the configured TCP port is absent. Do not treat a client connection failure alone as proof that the listener is gone.

For an upgrade, stop the service, replace the daemon and CLI through the platform install guide, and start it again. The database, daemon identity, remote configuration, actors, credentials, and audit records remain in their existing locations. Confirm the local manifest identity, local health, expected listener state, one least-privilege remote read, and one revocation check after the upgrade. Installation and upgrade do not create or enable a remote listener.

### 5. Back up and recover

Back up the daemon database, configuration, certificate, and private key under the platform's protected administrative procedure. Backing up user profile metadata does not back up native credential-store values. A database restore preserves daemon identity and client relationships; a deliberate clone must reset one copy's daemon identity and re-pair its clients before remote use.

If a bearer may be exposed, revoke its credential locally and create a new pairing. If a private key may be exposed, disable or firewall the listener, replace the certificate and key, restart, and deliberately repair every client profile after independently verifying the new certificate and daemon identity. If the address changes, update the operator-owned DNS, routing, tunnel, or profile through an authenticated workflow; clients never rewrite endpoints automatically. If the daemon identity changes unexpectedly, stop and investigate the data restore or clone rather than accepting the new identity.

Network loss and daemon restart use bounded automatic recovery. Credential rejection, capability reduction, certificate change, daemon identity mismatch, and incompatible API versions stop automatic recovery and require the targeted repair shown by the client. A mutation that loses its response is never replayed automatically; refresh authoritative state before deciding whether to submit it again.

## Threats and verification

The threat boundary covers individually administered daemons and their approved clients. It does not claim defense against an administrator who controls the daemon process, database, executable, or host operating system.

| ID | Risk | Required control | Acceptance-test class |
| --- | --- | --- | --- |
| T01 | Installation or upgrade accidentally exposes a port | Default disabled, explicit address, pre-bind validation | Clean install, upgrade, restart, and socket-inventory tests |
| T02 | Passive capture or transport downgrade reveals data or bearer values | HTTPS everywhere, TLS 1.3 minimum, trusted certificate | Protocol-version, certificate, hostname, and plaintext-refusal tests |
| T03 | An unauthenticated source guesses credentials or phrases | Coarse source limiter, phrase attempt budget, generic failure | Burst, refill, expiry, exhaustion, and enumeration tests |
| T04 | A captured phrase or old credential is replayed | Single-use phrase, independent credential, atomic rotation | Concurrent exchange, replay, rotation, and cancellation tests |
| T05 | A valid client invokes excessive authority | Server-owned capability mapping before handler dispatch | Per-operation allow, deny, stale-authority, and no-side-effect tests |
| T06 | A client spoofs proxy metadata or actor identity | Immediate-peer proxy allowlist, ignored forwarding headers, bearer-owned actor | Trusted and untrusted proxy integration tests |
| T07 | Oversized or slow requests exhaust daemon resources | Header, body, collection, rate, timeout, and connection bounds | Slow-header, large-body, concurrency, cancellation, and leak tests |
| T08 | Secrets appear in responses, logs, events, audit, or diagnostics | Explicit secret inventory, allowlisted response types, structural exclusion | Canary credential, phrase, trigger, environment, and webhook leakage tests |
| T09 | Revoked or expired authority survives in caches or streams | Current-state lookup and periodic stream revalidation | Revocation race, expiry boundary, open-stream, and restart tests |
| T10 | Same-named daemons cause a wrong-target action | Stable daemon identity, manifest negotiation, target-visible mutations | Clone, rename, profile mismatch, confirmation, and stale-cache tests |
| T11 | Version skew changes request meaning or hides capability | URI major, OpenAPI generation, manifest negotiation, overlap policy | Older-client, older-daemon, unknown-major, and deprecation tests |
| T12 | A hostile web origin drives a credentialed browser request | Reject `Origin` by default and emit no permissive CORS policy | Cross-origin preflight, simple-request, and trusted-native-client tests |

Security tests are bounded to these risks and supported modes. A review may add a threat only when it identifies an exposed asset, attacker capability, required control, and executable acceptance class.

## Dependency ownership

Go standard library owns HTTP, TLS, randomness, digests, and constant-time comparison. `golang.org/x/time/rate` owns rate limiting. OpenAPI 3.1 and pinned `oapi-codegen/v2` own API description and generated boundaries. `golang.org/x/crypto/argon2` owns enrollment phrase verification. Pinned `zalando/go-keyring` owns native client credential storage.

| Component | Introduction owner | Maintenance rule | Replacement trigger |
| --- | --- | --- | --- |
| `net/http` and `crypto/*` | #168 and #169 through the supported Go toolchain | Follow Go security releases and retain focused boundary tests | A proven requirement cannot be met by the supported standard library |
| `golang.org/x/time/rate` | #168 | Pin directly, include Dependabot coverage, review releases, and test with explicit timestamps | Maintenance or security posture fails, or required bounds cannot be expressed |
| OpenAPI 3.1 and `oapi-codegen/v2` | #168 | Pin the tool, commit source and output, verify clean regeneration, review license and advisories | Contract drift, maintenance failure, or inability to represent a required stable operation |
| `golang.org/x/crypto/argon2` | #169 | Pin directly, use Argon2id parameters benchmarked and recorded on supported platforms, review Go security releases | Supported guidance changes or resource bounds cannot be met |
| `zalando/go-keyring` | #169 | Permit native Keychain, Credential Manager, and Secret Service only; run native tests; review releases, license, and advisories | Maintenance failure, interactive or plaintext fallback, or supported-platform contract failure |

The owning implementation issues pin each reviewed version, prove clean restoration and supported-platform behavior, and record any deviation. Unsupported native credential storage fails closed; there is no application-file fallback.

## Non-goals

No JWT, user-account system, SSO, custom certificate authority, automatic public exposure, general policy language, offline mutation queue, or remote MCP mutation authority is included.

The boundary also excludes federation, teams, multi-tenancy, OAuth authorization-server behavior, automatic ACME, NAT traversal, service discovery, per-resource ACLs, deny expressions, background mutation replay, multi-master synchronization, conflict resolution, credential export, and credential recovery. Lost credentials are replaced through independently authorized administration.

## Required implementation order

1. #166 establishes daemon identity and discovery.
2. #167 establishes actors, capabilities, and audit.
3. #168 implements HTTPS and OpenAPI.
4. #169 implements enrollment and credential storage.
5. #170 and #171 implement desktop and CLI clients.
6. #172 implements resilience.
7. #173 qualifies the release.

No network implementation begins until S074 and [issue #165](https://github.com/shruggietech/go-schedule/issues/165) are reviewed and merged. Each downstream issue remains open until its own acceptance criteria and verification are complete.

## Implemented HTTPS and enrollment operations

The daemon now implements the S074 transport and enrollment core. Remote access remains disabled unless `remote.enabled` is true and an exact numeric `remote.bind_address`, `remote.certificate_file`, and `remote.private_key_file` are configured. A wildcard or public address additionally requires `remote.acknowledge_public_exposure: true`. Certificate loading succeeds before the TCP bind, and the listener accepts TLS 1.3 only.

```text
{
  "remote": {
    "enabled": true,
    "bind_address": "10.0.0.20:8443",
    "certificate_file": "/etc/go-schedule/server.crt",
    "private_key_file": "/etc/go-schedule/server.key",
    "acknowledge_public_exposure": false
  }
}
```

Use a private address or an SSH tunnel for the normal deployment. Direct public and reverse-proxy deployments retain application TLS, require explicit exposure acknowledgement where applicable, and leave certificate issuance, renewal, DNS, routing, and firewall ownership with the operator. Disabling the setting removes the TCP listener on the next daemon start without changing local IPC or stored client relationships.

Create a desktop phrase from a protected local shell, then move the displayed pairing ID, daemon ID, phrase, HTTPS address, and trusted certificate into Connections in the desktop application. The phrase is displayed once, expires after ten minutes, permits five failed attempts, and cannot be used as an API credential.

```text
gosched pairing create "Admin laptop" --kind desktop --capability manage
gosched pairing list
gosched pairing cancel <pairing-id>
gosched credential list
gosched credential rotate <credential-id>
gosched credential revoke <credential-id>
```

Rotation prints the new bearer value once and invalidates the former value immediately. Revocation also revokes its actor. Treat terminal output containing a phrase or newly rotated credential as sensitive and follow protected shell-output practices. There is no credential export or recovery path.

Remote clients use `Authorization: Bearer <credential>` against the documented `/api/v1` operations. Missing, malformed, unknown, and expired values receive `401 unauthorized`; a credential that matches a revoked record receives `401 credential_revoked` so an authenticated profile holder can distinguish repair from a transient outage. Browser-origin requests are rejected, rate limits return `429` with `Retry-After`, and unsupported local-only routes return `404`. See `api/openapi/remote-v1.yaml` for the machine-readable contract.
