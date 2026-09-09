#!/bin/sh
# Enforce the reviewed remote-access architecture before network implementation.
set -eu

ROOT=${1:-.}
DOC="$ROOT/docs/remote-access.md"
FAILURES=$(mktemp)
ORDER=$(mktemp)
trap 'rm -f "$FAILURES" "$ORDER"' EXIT HUP INT TERM

report() {
  printf '%s: %s\n' "$DOC" "$1" >> "$FAILURES"
}

require_fixed() {
  text=$1
  description=$2
  if ! grep -Fq -- "$text" "$DOC"; then
    report "missing architecture contract: $description"
  fi
}

if [ ! -f "$DOC" ]; then
  printf '%s: architecture page not found\n' "$DOC" >&2
  exit 1
fi

require_fixed 'Architecture contract only. No remote listener is implemented.' 'current-state boundary'
require_fixed 'HTTPS with versioned HTTP/JSON' 'primary HTTPS transport'
require_fixed 'local IPC mux and remote operation allowlist are separate adapters' 'separate local and remote adapters'
require_fixed 'Authentication and capability authorization run before daemon application behavior.' 'authorization-before-operation order'
require_fixed 'Observe, Operate, Manage, and Enroll' 'bounded capability vocabulary'
require_fixed 'enrollment phrase is short-lived and single-use' 'ephemeral enrollment phrase'
require_fixed 'opaque 256-bit bearer value retained only as a digest' 'digest-only durable credential'
require_fixed 'OpenAPI 3.1' 'machine-readable API source'
require_fixed 'oapi-codegen/v2' 'pinned Go contract generation owner'
require_fixed 'golang.org/x/time/rate' 'rate-limit library owner'
require_fixed 'golang.org/x/crypto/argon2' 'enrollment verifier owner'
require_fixed 'zalando/go-keyring' 'native credential-store owner'
require_fixed 'Server-Sent Events' 'one-way live transport'
require_fixed 'no automatic replay' 'uncertain mutation replay prohibition'
require_fixed 'TLS 1.3' 'minimum TLS policy'
require_fixed 'Origin' 'cross-origin default-deny control'
require_fixed 'Product ownership' 'product ownership'
require_fixed 'Operator ownership' 'operator ownership'
require_fixed 'offline mutation queue' 'offline mutation queue non-goal'
require_fixed 'remote MCP mutation authority' 'remote MCP mutation non-goal'

for mode in \
  'Local IPC' \
  'Private-network HTTPS' \
  'SSH-tunneled HTTPS' \
  'Reverse-proxied HTTPS' \
  'Direct HTTPS'; do
  if ! grep -Fq -- "| $mode |" "$DOC"; then
    report "missing deployment mode: $mode"
  fi
done

threat_count=$(grep -Ec '^\| T[0-9][0-9] \|' "$DOC" || true)
if [ "$threat_count" -lt 12 ]; then
  report "missing threat coverage: expected at least 12 threat-to-test rows, found $threat_count"
fi

awk '
  /^## Required implementation order$/ { inside = 1; next }
  inside && /^## / { exit }
  inside { print }
' "$DOC" > "$ORDER"

previous=0
for issue in 166 167 168 169 170 171 172 173; do
  line=$(grep -n -m 1 -F -- "#$issue " "$ORDER" | cut -d: -f1 || true)
  if [ -z "$line" ]; then
    report "missing downstream issue #$issue"
    continue
  fi
  if [ "$line" -lt "$previous" ]; then
    report "invalid downstream sequence: issue #$issue appears before its prerequisite"
  fi
  previous=$line
done

line167=$(grep -n -m 1 -F -- '#167 ' "$ORDER" | cut -d: -f1 || true)
line168=$(grep -n -m 1 -F -- '#168 ' "$ORDER" | cut -d: -f1 || true)
if [ -n "$line167" ] && [ -n "$line168" ] && [ "$line167" -ge "$line168" ]; then
  report 'invalid downstream sequence: issue #167 must appear before #168'
fi

if [ -s "$FAILURES" ]; then
  cat "$FAILURES" >&2
  printf 'remote-architecture-check: FAILED with %s issue(s)\n' "$(wc -l < "$FAILURES" | tr -d ' ')" >&2
  exit 1
fi

printf '%s\n' 'remote-architecture-check: OK (transport, trust, deployment, threats, dependencies, and issue order preserved)'
