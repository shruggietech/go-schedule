#!/bin/sh
set -eu

ROOT=$(CDPATH='' cd -- "$(dirname "$0")/../.." && pwd)
CHECK="$ROOT/scripts/remote-architecture-check.sh"
TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT HUP INT TERM

fail() {
  printf '%s\n' "remote-architecture-check fixture: $*" >&2
  exit 1
}

write_good() {
  dest=$1
  mkdir -p "$dest/docs"
  cat > "$dest/docs/remote-access.md" <<'EOF'
# Remote access architecture

**Current status:** Architecture contract only. No remote listener is implemented.

**Primary remote transport:** HTTPS with versioned HTTP/JSON.

The local IPC mux and remote operation allowlist are separate adapters over shared daemon operations.

## Ordered request boundary

Authentication and capability authorization run before daemon application behavior.

## Capability model

Observe, Operate, Manage, and Enroll are the complete capability vocabulary.

## Credential lifecycle

An enrollment phrase is short-lived and single-use. A durable credential is an opaque 256-bit bearer value retained only as a digest.

## API compatibility and live updates

OpenAPI 3.1 is the source of truth. Remote paths use `/api/v1`. Server-Sent Events use `Last-Event-ID`. Mutations have no automatic replay.

## Deployment modes

| Mode | Product ownership | Operator ownership | Support posture |
| --- | --- | --- | --- |
| Local IPC | Protected local API | OS accounts and groups | Default |
| Private-network HTTPS | Application controls | Private routing | Recommended |
| SSH-tunneled HTTPS | Application controls | SSH tunnel | Recommended |
| Reverse-proxied HTTPS | Application controls | Proxy and public TLS | Supported |
| Direct HTTPS | Application controls | Certificate and firewall | Advanced |

## Threats and verification

| ID | Risk | Control | Acceptance-test class |
| --- | --- | --- | --- |
| T01 | Accidental exposure | Opt-in bind | Default exposure test |
| T02 | Downgrade | TLS 1.3 | TLS test |
| T03 | Guessing | Source limiter | Guessing test |
| T04 | Replay | Independent credential | Replay test |
| T05 | Privilege escalation | Capability check | Authorization test |
| T06 | Proxy spoofing | Trusted peer | Proxy test |
| T07 | Request exhaustion | Bounded input | Bounds test |
| T08 | Secret leakage | Redaction | Leakage test |
| T09 | Stale authority | Stream revalidation | Revocation test |
| T10 | Wrong target | Stable daemon identity | Identity test |
| T11 | Version skew | Capability negotiation | Compatibility test |
| T12 | Cross-origin misuse | Reject Origin by default | Origin test |

## Dependency ownership

Go standard library owns HTTP, TLS, randomness, digests, and constant-time comparison. `golang.org/x/time/rate` owns rate limiting. OpenAPI 3.1 and pinned `oapi-codegen/v2` own API description and generated boundaries. `golang.org/x/crypto/argon2` owns enrollment phrase verification. Pinned `zalando/go-keyring` owns native client credential storage.

## Non-goals

No JWT, user-account system, SSO, custom certificate authority, automatic public exposure, general policy language, offline mutation queue, or remote MCP mutation authority is included.

## Required implementation order

1. #166 establishes daemon identity and discovery.
2. #167 establishes actors, capabilities, and audit.
3. #168 implements HTTPS and OpenAPI.
4. #169 implements enrollment and credential storage.
5. #170 and #171 implement desktop and CLI clients.
6. #172 implements resilience.
7. #173 qualifies the release.
EOF
}

run_expect_fail() {
  name=$1
  expected=$2
  shift 2
  if output=$(sh "$CHECK" "$@" 2>&1); then
    fail "$name passed unexpectedly"
  fi
  case "$output" in
    *"$expected"*) : ;;
    *) fail "$name did not report '$expected': $output" ;;
  esac
}

[ -f "$CHECK" ] || fail "missing checker: $CHECK"

GOOD="$TMP/good"
write_good "$GOOD"
sh "$CHECK" "$GOOD"

missing_transport="$TMP/missing-transport"
cp -R "$GOOD" "$missing_transport"
sed 's/HTTPS with versioned HTTP\/JSON/unspecified transport/' "$GOOD/docs/remote-access.md" > "$missing_transport/docs/remote-access.md"
run_expect_fail missing-transport 'primary HTTPS transport' "$missing_transport"

merged_mux="$TMP/merged-mux"
cp -R "$GOOD" "$merged_mux"
sed 's/local IPC mux and remote operation allowlist are separate adapters/local IPC mux is exported directly/' "$GOOD/docs/remote-access.md" > "$merged_mux/docs/remote-access.md"
run_expect_fail merged-mux 'separate local and remote adapters' "$merged_mux"

for mode in 'Local IPC' 'Private-network HTTPS' 'SSH-tunneled HTTPS' 'Reverse-proxied HTTPS' 'Direct HTTPS'; do
  slug=$(printf '%s' "$mode" | tr '[:upper:] ' '[:lower:]-')
  dest="$TMP/missing-$slug"
  cp -R "$GOOD" "$dest"
  grep -Fv "| $mode |" "$GOOD/docs/remote-access.md" > "$dest/docs/remote-access.md"
  run_expect_fail "missing-$slug" "deployment mode: $mode" "$dest"
done

missing_owner="$TMP/missing-owner"
cp -R "$GOOD" "$missing_owner"
sed 's/Operator ownership/External responsibility/' "$GOOD/docs/remote-access.md" > "$missing_owner/docs/remote-access.md"
run_expect_fail missing-owner 'operator ownership' "$missing_owner"

missing_threat="$TMP/missing-threat"
cp -R "$GOOD" "$missing_threat"
grep -Fv '| T12 |' "$GOOD/docs/remote-access.md" > "$missing_threat/docs/remote-access.md"
run_expect_fail missing-threat '12 threat-to-test rows' "$missing_threat"

missing_non_goal="$TMP/missing-non-goal"
cp -R "$GOOD" "$missing_non_goal"
sed 's/offline mutation queue/background synchronization/' "$GOOD/docs/remote-access.md" > "$missing_non_goal/docs/remote-access.md"
run_expect_fail missing-non-goal 'offline mutation queue non-goal' "$missing_non_goal"

wrong_order="$TMP/wrong-order"
cp -R "$GOOD" "$wrong_order"
awk '
  /#167 establishes actors, capabilities, and audit/ { print "2. #168 implements HTTPS and OpenAPI."; next }
  /#168 implements HTTPS and OpenAPI/ { print "3. #167 establishes actors, capabilities, and audit."; next }
  { print }
' "$GOOD/docs/remote-access.md" > "$wrong_order/docs/remote-access.md"
run_expect_fail wrong-order 'issue #167 must appear before #168' "$wrong_order"

printf '%s\n' 'remote-architecture-check fixtures: OK (required boundary mutations rejected)'
