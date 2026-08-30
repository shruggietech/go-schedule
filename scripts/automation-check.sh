#!/bin/sh
# Offline contract check for workflow action majors and local verification gates.

set -eu

ROOT=${1:-.}
GO=${GO:-go}

if [ ! -d "$ROOT/.github/workflows" ]; then
  printf 'automation-check: workflow directory not found: %s\n' \
    "$ROOT/.github/workflows" >&2
  exit 1
fi

FAILURES=$(mktemp)
trap 'rm -f "$FAILURES"' EXIT HUP INT TERM

report() {
  printf '%s\n' "$*" >> "$FAILURES"
}

approved_for() {
  case "$1" in
    actions/checkout) printf '%s' 'actions/checkout@v7' ;;
    actions/setup-go) printf '%s' 'actions/setup-go@v7' ;;
    actions/upload-artifact) printf '%s' 'actions/upload-artifact@v7' ;;
    github/codeql-action/init) printf '%s' 'github/codeql-action/init@v4' ;;
    github/codeql-action/analyze) printf '%s' 'github/codeql-action/analyze@v4' ;;
    softprops/action-gh-release) printf '%s' 'softprops/action-gh-release@v3' ;;
    *) return 1 ;;
  esac
}

check_action() {
  file=$1
  line=$2
  ref=$3
  family=${ref%@*}

  if approved=$(approved_for "$family"); then
    if [ "$ref" != "$approved" ]; then
      report "$file:$line: action reference $ref is not approved; expected $approved"
    fi
  else
    report "$file:$line: unaudited action reference $ref"
  fi
}

found_workflow=0
for file in "$ROOT"/.github/workflows/*.yml "$ROOT"/.github/workflows/*.yaml; do
  [ -f "$file" ] || continue
  found_workflow=1
  awk '
    match($0, /uses:[[:space:]]*[^[:space:]#]+/) {
      ref = substr($0, RSTART, RLENGTH)
      sub(/^uses:[[:space:]]*/, "", ref)
      gsub(/["'"'"']/, "", ref)
      print FNR "|" ref
    }
  ' "$file" | while IFS='|' read -r line ref; do
    check_action "$file" "$line" "$ref"
  done
done

if [ "$found_workflow" -eq 0 ]; then
  report "$ROOT/.github/workflows: no .yml or .yaml workflow files found"
fi

CODEQL="$ROOT/.github/workflows/codeql.yml"

require_codeql_pattern() {
  pattern=$1
  description=$2
  if ! grep -Eq "$pattern" "$CODEQL"; then
    report "$CODEQL: missing $description"
  fi
}

trigger_targets_main() {
  trigger=$1
  awk -v trigger="$trigger" '
    $0 == "  " trigger ":" { inside = 1; next }
    inside && /^  [[:alnum:]_-]+:/ { exit }
    inside && /^    branches: \[main\]$/ { found = 1 }
    END { exit !found }
  ' "$CODEQL"
}

if [ ! -f "$CODEQL" ]; then
  report "$CODEQL: canonical CodeQL workflow not found"
else
  require_codeql_pattern '^name: CodeQL$' 'stable CodeQL workflow name'

  for trigger in push pull_request; do
    if ! trigger_targets_main "$trigger"; then
      report "$CODEQL: missing $trigger trigger targeting main"
    fi
  done

  require_codeql_pattern '^  schedule:$' 'weekly schedule trigger'
  require_codeql_pattern \
    "^    - cron: '17 4 \* \* 1'$" \
    "weekly schedule trigger (expected '17 4 * * 1')"

  if grep -Eq '^[[:space:]]+permissions:' "$CODEQL"; then
    report "$CODEQL: job-level permissions override is not allowed"
  fi

  permissions=$(awk '
    /^permissions:$/ { inside = 1; next }
    inside && /^[^[:space:]]/ { exit }
    inside && NF { print }
  ' "$CODEQL")
  expected_permissions='  contents: read
  security-events: write'
  if [ "$permissions" != "$expected_permissions" ]; then
    report "$CODEQL: permissions must be exactly contents: read and security-events: write"
  fi

  require_codeql_pattern '^    runs-on: ubuntu-latest$' \
    'Ubuntu analysis runner'
  require_codeql_pattern '^          go-version-file: go\.mod$' \
    'module-selected Go version'
  require_codeql_pattern '^          languages: go$' 'CodeQL Go language'
  require_codeql_pattern '^          build-mode: manual$' \
    'CodeQL manual build mode'
  require_codeql_pattern '^[[:space:]]+(-[[:space:]]+)?uses: github/codeql-action/init@v4$' \
    'CodeQL init step'
  require_codeql_pattern \
    '^        run: CGO_ENABLED=0 go build \./\.\.\.$' \
    'cgo-free headless build step'
  require_codeql_pattern '^[[:space:]]+(-[[:space:]]+)?uses: github/codeql-action/analyze@v4$' \
    'CodeQL analyze step'

  if grep -Eq '\$\{\{[[:space:]]*secrets([[:space:]]|\.|\[)' "$CODEQL"; then
    report "$CODEQL: CodeQL workflow must not consume repository or organization secrets"
  fi
fi

EXPECTED_GATES='format
vet
lint
race
gui
coverage
docs
automation'

DEPENDABOT="$ROOT/.github/dependabot.yml"

dependabot_block() {
  ecosystem=$1
  awk -v ecosystem="$ecosystem" '
    /^  - package-ecosystem:/ {
      if (inside) exit
      if ($0 == "  - package-ecosystem: " ecosystem) inside = 1
    }
    inside { print }
  ' "$DEPENDABOT"
}

require_dependabot_pattern() {
  block=$1
  pattern=$2
  description=$3
  if ! printf '%s\n' "$block" | grep -Eq "$pattern"; then
    report "$DEPENDABOT: missing $description"
  fi
}

if [ ! -f "$DEPENDABOT" ]; then
  report "$DEPENDABOT: Dependabot configuration not found"
else
  ecosystem_count=$(grep -Ec '^  - package-ecosystem:' "$DEPENDABOT" || true)
  if [ "$ecosystem_count" -ne 2 ]; then
    report "$DEPENDABOT: expected exactly two package ecosystems"
  fi

  invalid_ecosystems=$(awk '
    /^  - package-ecosystem:/ {
      ecosystem = $0
      sub(/^  - package-ecosystem:[[:space:]]*/, "", ecosystem)
      if (ecosystem != "gomod" && ecosystem != "github-actions") print ecosystem
    }
  ' "$DEPENDABOT")
  if [ -n "$invalid_ecosystems" ]; then
    report "$DEPENDABOT: unapproved package ecosystem(s): $invalid_ecosystems"
  fi

  for ecosystem in gomod github-actions; do
    count=$(grep -Ec "^  - package-ecosystem: $ecosystem$" "$DEPENDABOT" || true)
    if [ "$count" -ne 1 ]; then
      report "$DEPENDABOT: expected exactly one $ecosystem entry"
      continue
    fi

    block=$(dependabot_block "$ecosystem")
    require_dependabot_pattern "$block" '^    directory: /$' \
      "$ecosystem root directory"
    require_dependabot_pattern "$block" '^    open-pull-requests-limit: 5$' \
      "$ecosystem five-PR limit"
    require_dependabot_pattern "$block" '^      routine-minor-and-patch:$' \
      "$ecosystem routine update group"
    require_dependabot_pattern "$block" '^          - minor$' \
      "$ecosystem grouped minor updates"
    require_dependabot_pattern "$block" '^          - patch$' \
      "$ecosystem grouped patch updates"
    require_dependabot_pattern "$block" '^      - dependencies$' \
      "$ecosystem dependencies label"

    case "$ecosystem" in
      gomod) cadence=weekly ;;
      github-actions) cadence=monthly ;;
    esac
    require_dependabot_pattern "$block" "^      interval: $cadence$" \
      "$ecosystem $cadence cadence"

    if printf '%s\n' "$block" | grep -Eq '^          - major$'; then
      report "$DEPENDABOT: $ecosystem routine group must not group major updates"
    fi
    if printf '%s\n' "$block" | grep -Eq \
      '^        applies-to: security-updates$'; then
      report "$DEPENDABOT: $ecosystem routine group must not group security updates"
    fi

    if [ "$ecosystem" = gomod ]; then
      require_dependabot_pattern "$block" '^      gui-minor-and-patch:$' \
        'gomod isolated GUI update group'
      require_dependabot_pattern "$block" '^          - fyne\.io/\*$' \
        'gomod GUI dependency pattern'
      require_dependabot_pattern "$block" '^      storage-minor-and-patch:$' \
        'gomod isolated storage update group'
      require_dependabot_pattern "$block" '^          - modernc\.org/sqlite$' \
        'gomod storage dependency pattern'
      require_dependabot_pattern "$block" '^      platform-minor-and-patch:$' \
        'gomod isolated platform update group'
      require_dependabot_pattern "$block" \
        '^          - github\.com/kardianos/service$' \
        'gomod service dependency pattern'
      require_dependabot_pattern "$block" '^          - golang\.org/x/sys$' \
        'gomod system dependency pattern'
    fi
  done
fi

LIFECYCLE="$ROOT/scripts/spec-lifecycle-check.sh"
if [ ! -f "$LIFECYCLE" ]; then
  report "$LIFECYCLE: specification lifecycle checker not found"
elif ! lifecycle_output=$(sh "$LIFECYCLE" "$ROOT" 2>&1); then
  report "$LIFECYCLE: lifecycle contract failed"
  report "$lifecycle_output"
fi

RELEASE="$ROOT/.github/workflows/release.yml"
RELEASE_NOTES_DIR="$ROOT/.github/release-notes"
require_release_text() {
  text=$1
  description=$2
  if ! grep -Fq -- "$text" "$RELEASE"; then
    report "$RELEASE: missing $description"
  fi
}

if [ ! -f "$RELEASE" ]; then
  report "$RELEASE: release workflow not found"
else
  require_release_text '-icon=brand/platform/windows/go-schedule.ico' \
    'canonical Windows ICO'
  require_release_text \
    "cp brand/platform/macos/go-schedule.icns \"\$app/Contents/Resources/icon.icns\"" \
    'canonical macOS ICNS'
  require_release_text \
    "cp brand/platform/linux/go-schedule.desktop \"\$stage/share/applications/\"" \
    'canonical Linux desktop entry'
  require_release_text \
    "cp -R brand/platform/linux/hicolor \"\$stage/share/icons/\"" \
    'canonical Linux hicolor tree'
  if grep -Fq -- 'generate_release_notes: true' "$RELEASE"; then
    report "$RELEASE: generated release notes must remain disabled"
  fi
  if grep -Eq '(^|[[:space:]])git[[:space:]]+push([[:space:]]|$)' "$RELEASE"; then
    report "$RELEASE: release workflow must not run git push"
  fi
  if grep -Eq "^[[:space:]]*ref:[[:space:]]*(main|'main'|\"main\")[[:space:]]*(#.*)?$" \
    "$RELEASE"; then
    report "$RELEASE: release validation must inspect the tagged checkout"
  fi
  if ! awk '
    /^  binaries:[[:space:]]*$/ { in_binaries = 1; next }
    in_binaries && /^  [[:alnum:]_-]+:[[:space:]]*$/ { exit }
    in_binaries && /^    needs:[[:space:]]*readme-version[[:space:]]*$/ {
      found = 1
    }
    END { exit !found }
  ' "$RELEASE"; then
    report "$RELEASE: binaries job missing README version preflight dependency"
  fi
  require_release_text 'generate_release_notes: false' \
    'disabled generated release notes contract'
  require_release_text \
    "body_path: .github/release-notes/\${{ github.ref_name }}.md" \
    'dynamic tag-specific release-note body path'
  require_release_text 'libwayland-dev' \
    'Linux desktop Wayland development headers'
  require_release_text 'wayland-protocols' \
    'Linux desktop Wayland protocols'
fi

if [ ! -d "$RELEASE_NOTES_DIR" ]; then
  report "$RELEASE_NOTES_DIR: release-note directory not found"
else
  found_release_notes=false
  for notes in "$RELEASE_NOTES_DIR"/v[0-9]*.[0-9]*.[0-9]*.md; do
    [ -f "$notes" ] || continue
    found_release_notes=true
    tag=$(basename "$notes" .md)

    if [ "$(grep -c '^## Highlights$' "$notes" || true)" -ne 1 ]; then
      report "$notes: expected exactly one Highlights heading"
    fi

    highlight_count=$(grep -c '^- ' "$notes" || true)
    if [ "$highlight_count" -lt 4 ] || [ "$highlight_count" -gt 6 ]; then
      report "$notes: expected four to six highlight bullets"
    fi

    line_count=$(awk 'END { print NR }' "$notes")
    expected_line_count=$((highlight_count + 4))
    last_bullet_line=$((highlight_count + 2))
    separator_line=$((highlight_count + 3))
    link_line=$((highlight_count + 4))
    invalid_body=false

    if [ "$line_count" -ne "$expected_line_count" ] ||
       [ -n "$(sed -n '2p' "$notes")" ] ||
       [ -n "$(sed -n "${separator_line}p" "$notes")" ] ||
       sed -n "3,${last_bullet_line}p" "$notes" | grep -vq '^- ' ||
       ! sed -n "${link_line}p" "$notes" |
         grep -Eq '^Read the \[full changelog\]\(.+\) for every change\.$'; then
      invalid_body=true
    fi

    if [ "$invalid_body" = true ] ||
       [ "$(grep -c '^## ' "$notes" || true)" -ne 1 ] ||
       grep -E '^#{1,6} ' "$notes" | grep -qv '^## Highlights$'; then
      report "$notes: expected highlights-only release copy"
    fi

    changelog_prefix="https://github.com/shruggietech/go-schedule/blob/$tag/CHANGELOG.md#"
    if [ "$(grep -Fc '[full changelog](' "$notes" || true)" -ne 1 ] ||
       [ "$(grep -Fc "$changelog_prefix" "$notes" || true)" -ne 1 ]; then
      report "$notes: expected exactly one tagged full changelog link"
    fi
  done

  if [ "$found_release_notes" = false ]; then
    report "$RELEASE_NOTES_DIR: no tag-specific release notes found"
  fi
fi

BRAND_CHECK="$ROOT/scripts/brand-check"
if [ ! -d "$BRAND_CHECK" ]; then
  report "$BRAND_CHECK: brand integrity command not found"
elif ! brand_output=$(cd "$ROOT" && "$GO" run ./scripts/brand-check 2>&1); then
  report "$BRAND_CHECK: brand integrity contract failed"
  report "$brand_output"
fi

VERIFY="$ROOT/scripts/verify.sh"
if [ ! -f "$VERIFY" ]; then
  report "$VERIFY: verification driver not found"
elif ! observed=$(sh "$VERIFY" list 2>&1); then
  report "$VERIFY: could not read gate manifest: $observed"
elif [ "$observed" != "$EXPECTED_GATES" ]; then
  report "gate manifest differs"
  report "expected:"
  report "$EXPECTED_GATES"
  report "observed:"
  report "$observed"
fi

if [ -s "$FAILURES" ]; then
  cat "$FAILURES" >&2
  printf 'automation-check: FAILED with %s issue(s)\n' \
    "$(grep -c '^[^[:space:]]' "$FAILURES")" >&2
  exit 1
fi

printf '%s\n' \
  'automation-check: OK - actions, CodeQL, Dependabot, release notes, brand, lifecycle, and 8 gates'
