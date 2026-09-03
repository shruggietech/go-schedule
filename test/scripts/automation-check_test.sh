#!/bin/sh
# Contract regressions for the maintainer automation baseline.

set -eu

ROOT=$(CDPATH='' cd "$(dirname "$0")/../.." && pwd)
CHECK="$ROOT/scripts/automation-check.sh"
VERIFY="$ROOT/scripts/verify.sh"
MODE=${1:-all}

fail() {
  printf 'automation-check-test: FAIL: %s\n' "$*" >&2
  exit 1
}

assert_contains() {
  haystack=$1
  needle=$2
  case "$haystack" in
    *"$needle"*) : ;;
    *) fail "expected output to contain '$needle'; got: $haystack" ;;
  esac
}

run_expect_pass() {
  label=$1
  shift
  if ! output=$("$@" 2>&1); then
    fail "$label should pass; got: $output"
  fi
}

run_expect_fail() {
  label=$1
  needle=$2
  shift 2
  if output=$("$@" 2>&1); then
    fail "$label should fail"
  fi
  assert_contains "$output" "$needle"
}

make_fixture() {
  fixture=$1
  manifest=$2
  mkdir -p "$fixture/.github/workflows" \
    "$fixture/.github/release-notes" \
    "$fixture/scripts/brand-check" \
    "$fixture/specs"
  cat > "$fixture/go.mod" <<'EOF'
module fixture

go 1.25.0
EOF
  cat > "$fixture/scripts/brand-check/main.go" <<'EOF'
package main
func main() {}
EOF
  cat > "$fixture/.github/workflows/ci.yml" <<'EOF'
name: fixture
jobs:
  test:
    steps:
      - uses: actions/checkout@v7
      - uses: actions/setup-go@v7
      - uses: actions/upload-artifact@v7
      - uses: softprops/action-gh-release@v3
EOF
  cat > "$fixture/.github/dependabot.yml" <<'EOF'
version: 2
updates:
  - package-ecosystem: gomod
    directory: /
    schedule:
      interval: weekly
    open-pull-requests-limit: 5
    groups:
      gui-minor-and-patch:
        patterns:
          - fyne.io/*
        update-types:
          - minor
          - patch
      storage-minor-and-patch:
        patterns:
          - modernc.org/sqlite
        update-types:
          - minor
          - patch
      platform-minor-and-patch:
        patterns:
          - github.com/kardianos/service
          - golang.org/x/sys
        update-types:
          - minor
          - patch
      routine-minor-and-patch:
        update-types:
          - minor
          - patch
    labels:
      - dependencies
  - package-ecosystem: github-actions
    directory: /
    schedule:
      interval: monthly
    open-pull-requests-limit: 5
    groups:
      routine-minor-and-patch:
        update-types:
          - minor
          - patch
    labels:
      - dependencies
EOF
  cat > "$fixture/scripts/spec-lifecycle-check.sh" <<'EOF'
#!/bin/sh
printf 'fixture lifecycle OK\n'
EOF
  cat > "$fixture/.github/workflows/codeql.yml" <<'EOF'
name: CodeQL
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
  schedule:
    - cron: '17 4 * * 1'
permissions:
  contents: read
  security-events: write
jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v7
      - uses: actions/setup-go@v7
        with:
          go-version-file: go.mod
      - uses: github/codeql-action/init@v4
        with:
          languages: go
          build-mode: manual
      - name: Build
        run: CGO_ENABLED=0 go build ./...
      - uses: github/codeql-action/analyze@v4
EOF
  cat > "$fixture/.github/workflows/release.yml" <<'EOF'
name: Release
jobs:
  readme-version:
    steps:
      - run: |
          BADGE_LINE_COUNT=$(sed 's/^[[:space:]]*//; s/[[:space:]]*$//' README.md | grep -Fxc -- "$BADGE" || true)
          test "$BADGE_LINE_COUNT" -eq 1
  release-state:
    needs: readme-version
    steps:
      - name: Refuse to mutate an already-public release
        run: test "$(printf '%s' "$BODY" | jq -r '.draft')" = true
  binaries:
    needs: release-state
    steps:
      - run: go run github.com/josephspurrier/goversioninfo/cmd/goversioninfo@v1.4.1 -icon=brand/platform/windows/go-schedule.ico
      - run: cp brand/platform/macos/go-schedule.icns "$app/Contents/Resources/icon.icns"
      - run: cp brand/platform/linux/go-schedule.desktop "$stage/share/applications/"
      - run: cp -R brand/platform/linux/hicolor "$stage/share/icons/"
      - run: sudo apt-get install -y libwayland-dev wayland-protocols
      - uses: softprops/action-gh-release@v3
        with:
          draft: true
          generate_release_notes: false
          body_path: .github/release-notes/${{ github.ref_name }}.md
          files: windows-candidate-manifest.json
EOF
  cat > "$fixture/.github/workflows/promote-release.yml" <<'EOF'
name: Promote Release
on:
  workflow_dispatch:
permissions:
  actions: read
jobs:
  promote:
    steps:
      - run: |
          if ! printf '%s' "$RELEASE_TAG" | grep -Eq '^v[0-9]+\.[0-9]+\.[0-9]+$'; then exit 2; fi
          test "$(printf '%s' "$RELEASE_JSON" | jq -r '.draft')" = true
           MANIFEST=assets/windows-candidate-manifest.json
           EVIDENCE=assets/go-schedule_${RELEASE_TAG}_windows-attended-evidence.zip
           expected-assets.txt
           go run ./scripts/windows-release-gate verify-candidate
           go run ./scripts/windows-release-gate verify-bundle \
             --candidate-manifest "$MANIFEST"
          LC_ALL=C sha256sum * > SHA256SUMS.txt
          sha256sum -c SHA256SUMS.txt
           gh release upload "$RELEASE_TAG" SHA256SUMS.txt
           Re-download and verify the final draft asset set
           git ls-remote origin "refs/tags/${RELEASE_TAG}"
           gh api -F draft=false
EOF
  cat > "$fixture/.github/release-notes/v0.9.0.md" <<'EOF'
## Highlights

- Complete scheduled work with durable follow-on task chains.
- Translate and validate broader cron expressions in either supported syntax.
- Preserve scheduling intent through daylight-saving transitions.
- Use a consistent visual identity across the desktop application and packages.

Read the [full changelog](https://github.com/shruggietech/go-schedule/blob/v0.9.0/CHANGELOG.md#090---2026-08-30) for every change.
EOF
  cat > "$fixture/scripts/verify.sh" <<EOF
#!/bin/sh
if [ "\${1:-}" = list ]; then
  printf '%s\\n' $manifest
  exit 0
fi
exit 2
EOF
}

run_automation_cases() {
  [ -f "$CHECK" ] || fail "missing $CHECK"

  tmp=$(mktemp -d)
  trap 'rm -rf "$tmp"' EXIT HUP INT TERM
  good="$tmp/good"
  make_fixture "$good" 'format vet lint race gui coverage docs automation'

  run_expect_pass approved sh "$CHECK" "$good"

  generated_notes="$tmp/generated-notes"
  cp -R "$good" "$generated_notes"
  sed 's/generate_release_notes: false/generate_release_notes: true/' \
    "$good/.github/workflows/release.yml" > \
    "$generated_notes/.github/workflows/release.yml"
  run_expect_fail generated-notes 'generated release notes must remain disabled' \
    sh "$CHECK" "$generated_notes"

  fixed_body_path="$tmp/fixed-body-path"
  cp -R "$good" "$fixed_body_path"
  sed 's/${{ github.ref_name }}/v0.9.0/' \
    "$good/.github/workflows/release.yml" > \
    "$fixed_body_path/.github/workflows/release.yml"
  run_expect_fail fixed-body-path 'dynamic tag-specific release-note body path' \
    sh "$CHECK" "$fixed_body_path"

  direct_main_push="$tmp/direct-main-push"
  cp -R "$good" "$direct_main_push"
  printf '      - run: git push origin HEAD:main\n' >> \
    "$direct_main_push/.github/workflows/release.yml"
  run_expect_fail direct-main-push 'must not invoke git' \
    sh "$CHECK" "$direct_main_push"

  quoted_main_push="$tmp/quoted-main-push"
  cp -R "$good" "$quoted_main_push"
  printf "      - run: git push origin 'HEAD:main'\n" >> \
    "$quoted_main_push/.github/workflows/release.yml"
  run_expect_fail quoted-main-push 'must not invoke git' \
    sh "$CHECK" "$quoted_main_push"

  full_main_refspec="$tmp/full-main-refspec"
  cp -R "$good" "$full_main_refspec"
  printf '      - run: git push origin HEAD:refs/heads/main\n' >> \
    "$full_main_refspec/.github/workflows/release.yml"
  run_expect_fail full-main-refspec 'must not invoke git' \
    sh "$CHECK" "$full_main_refspec"

  global_option_push="$tmp/global-option-push"
  cp -R "$good" "$global_option_push"
  printf '      - run: git -C "$GITHUB_WORKSPACE" push origin HEAD:main\n' >> \
    "$global_option_push/.github/workflows/release.yml"
  run_expect_fail global-option-push 'must not invoke git' \
    sh "$CHECK" "$global_option_push"

  continued_push="$tmp/continued-push"
  cp -R "$good" "$continued_push"
  printf '%s\n' '      - run: |' '          git \' \
    '            push origin HEAD:main' >> \
    "$continued_push/.github/workflows/release.yml"
  run_expect_fail continued-push 'must not invoke git' \
    sh "$CHECK" "$continued_push"

  folded_push="$tmp/folded-push"
  cp -R "$good" "$folded_push"
  printf '%s\n' '      - run: >' '          git' \
    '          push origin HEAD:main' >> \
    "$folded_push/.github/workflows/release.yml"
  run_expect_fail folded-push 'must not invoke git' \
    sh "$CHECK" "$folded_push"

  main_checkout="$tmp/main-checkout"
  cp -R "$good" "$main_checkout"
  printf '      - uses: actions/checkout@v7\n        with:\n          ref: main\n' >> \
    "$main_checkout/.github/workflows/release.yml"
  run_expect_fail main-checkout 'must inspect the tagged checkout' \
    sh "$CHECK" "$main_checkout"

  quoted_main_checkout="$tmp/quoted-main-checkout"
  cp -R "$good" "$quoted_main_checkout"
  printf '      - uses: actions/checkout@v7\n        with:\n          ref: "main"\n' >> \
    "$quoted_main_checkout/.github/workflows/release.yml"
  run_expect_fail quoted-main-checkout 'must inspect the tagged checkout' \
    sh "$CHECK" "$quoted_main_checkout"

  qualified_main_checkout="$tmp/qualified-main-checkout"
  cp -R "$good" "$qualified_main_checkout"
  printf "      - uses: actions/checkout@v7\n        with:\n          ref: 'refs/heads/main'\n" >> \
    "$qualified_main_checkout/.github/workflows/release.yml"
  run_expect_fail qualified-main-checkout 'must inspect the tagged checkout' \
    sh "$CHECK" "$qualified_main_checkout"

  ungated_release="$tmp/ungated-release"
  cp -R "$good" "$ungated_release"
  sed 's/needs: readme-version/needs: binaries/' \
    "$good/.github/workflows/release.yml" > \
    "$ungated_release/.github/workflows/release.yml"
  printf '  gui: # desktop release\n    needs: readme-version\n' >> \
    "$ungated_release/.github/workflows/release.yml"
  run_expect_fail ungated-release 'README version preflight dependency' \
    sh "$CHECK" "$ungated_release"

  nested_dependency="$tmp/nested-dependency"
  cp -R "$good" "$nested_dependency"
  sed 's/^    needs: readme-version$/    env:\
      needs: readme-version/' \
    "$good/.github/workflows/release.yml" > \
    "$nested_dependency/.github/workflows/release.yml"
  run_expect_fail nested-dependency 'README version preflight dependency' \
    sh "$CHECK" "$nested_dependency"

  sigpipe_preflight="$tmp/sigpipe-preflight"
  cp -R "$good" "$sigpipe_preflight"
  sed '/BADGE_LINE_COUNT=/c\          sed '\''s/^[[:space:]]*//; s/[[:space:]]*$//'\'' README.md | grep -Fqx -- "$BADGE"' \
    "$good/.github/workflows/release.yml" > \
    "$sigpipe_preflight/.github/workflows/release.yml"
  run_expect_fail sigpipe-preflight 'SIGPIPE-safe README badge count' \
    sh "$CHECK" "$sigpipe_preflight"

  missing_wayland="$tmp/missing-wayland"
  cp -R "$good" "$missing_wayland"
  sed 's/libwayland-dev //' "$good/.github/workflows/release.yml" > \
    "$missing_wayland/.github/workflows/release.yml"
  run_expect_fail missing-wayland 'Wayland development headers' \
    sh "$CHECK" "$missing_wayland"

  public_staging="$tmp/public-staging"
  cp -R "$good" "$public_staging"
  sed '/draft: true/d' "$good/.github/workflows/release.yml" > \
    "$public_staging/.github/workflows/release.yml"
  run_expect_fail public-staging 'explicitly use draft: true' \
    sh "$CHECK" "$public_staging"

  staging_checksums="$tmp/staging-checksums"
  cp -R "$good" "$staging_checksums"
  printf '      - run: sha256sum * > SHA256SUMS.txt\n' >> \
    "$staging_checksums/.github/workflows/release.yml"
  run_expect_fail staging-checksums 'must not publish or create final checksums' \
    sh "$CHECK" "$staging_checksums"

  missing_promotion="$tmp/missing-promotion"
  cp -R "$good" "$missing_promotion"
  rm -f "$missing_promotion/.github/workflows/promote-release.yml"
  run_expect_fail missing-promotion 'promotion workflow not found' \
    sh "$CHECK" "$missing_promotion"

  rebuild_promotion="$tmp/rebuild-promotion"
  cp -R "$good" "$rebuild_promotion"
  printf '          go build ./cmd/gosched-gui\n' >> \
    "$rebuild_promotion/.github/workflows/promote-release.yml"
  run_expect_fail rebuild-promotion 'must not rebuild release artifacts' \
    sh "$CHECK" "$rebuild_promotion"

  ungated_promotion="$tmp/ungated-promotion"
  cp -R "$good" "$ungated_promotion"
  sed '/windows-release-gate verify-bundle/d' \
    "$good/.github/workflows/promote-release.yml" > \
    "$ungated_promotion/.github/workflows/promote-release.yml"
  run_expect_fail ungated-promotion 'shared exact-candidate evidence gate' \
    sh "$CHECK" "$ungated_promotion"

  no_candidate_gate="$tmp/no-candidate-gate"
  cp -R "$good" "$no_candidate_gate"
  sed '/windows-release-gate verify-candidate/d' \
    "$good/.github/workflows/promote-release.yml" > \
    "$no_candidate_gate/.github/workflows/promote-release.yml"
  run_expect_fail no-candidate-gate 'independent candidate-manifest gate' \
    sh "$CHECK" "$no_candidate_gate"

  no_actions_read="$tmp/no-actions-read"
  cp -R "$good" "$no_actions_read"
  sed '/actions: read/d' \
    "$good/.github/workflows/promote-release.yml" > \
    "$no_actions_read/.github/workflows/promote-release.yml"
  run_expect_fail no-actions-read 'workflow-run read permission' \
    sh "$CHECK" "$no_actions_read"

  no_asset_allowlist="$tmp/no-asset-allowlist"
  cp -R "$good" "$no_asset_allowlist"
  sed '/expected-assets.txt/d' \
    "$good/.github/workflows/promote-release.yml" > \
    "$no_asset_allowlist/.github/workflows/promote-release.yml"
  run_expect_fail no-asset-allowlist 'exact draft asset allowlist' \
    sh "$CHECK" "$no_asset_allowlist"

  no_tag_recheck="$tmp/no-tag-recheck"
  cp -R "$good" "$no_tag_recheck"
  sed '/git ls-remote origin/d' \
    "$good/.github/workflows/promote-release.yml" > \
    "$no_tag_recheck/.github/workflows/promote-release.yml"
  run_expect_fail no-tag-recheck 'last-moment tag identity check' \
    sh "$CHECK" "$no_tag_recheck"

  early_promotion="$tmp/early-promotion"
  cp -R "$good" "$early_promotion"
  sed '/gh api -F draft=false/d; /go run .\/scripts\/windows-release-gate verify-candidate/i\
          gh api -F draft=false' \
    "$good/.github/workflows/promote-release.yml" > \
    "$early_promotion/.github/workflows/promote-release.yml"
  run_expect_fail early-promotion 'not ordered safely' \
    sh "$CHECK" "$early_promotion"

  missing_changelog="$tmp/missing-changelog"
  cp -R "$good" "$missing_changelog"
  sed '/Read the \[full changelog\]/d' \
    "$good/.github/release-notes/v0.9.0.md" > \
    "$missing_changelog/.github/release-notes/v0.9.0.md"
  run_expect_fail missing-changelog 'tagged full changelog link' \
    sh "$CHECK" "$missing_changelog"

  too_few_highlights="$tmp/too-few-highlights"
  cp -R "$good" "$too_few_highlights"
  sed '/consistent visual identity/d' \
    "$good/.github/release-notes/v0.9.0.md" > \
    "$too_few_highlights/.github/release-notes/v0.9.0.md"
  run_expect_fail too-few-highlights 'four to six highlight bullets' \
    sh "$CHECK" "$too_few_highlights"

  too_many_highlights="$tmp/too-many-highlights"
  cp -R "$good" "$too_many_highlights"
  sed '/Read the \[full changelog\]/i\
- Improve local daemon maintenance.\
- Refresh contributor automation.\
- Strengthen release packaging checks.' \
    "$good/.github/release-notes/v0.9.0.md" > \
    "$too_many_highlights/.github/release-notes/v0.9.0.md"
  run_expect_fail too-many-highlights 'four to six highlight bullets' \
    sh "$CHECK" "$too_many_highlights"

  exhaustive_copy="$tmp/exhaustive-copy"
  cp -R "$good" "$exhaustive_copy"
  printf '\n## Installation\n\nDownload the package for your platform.\n' >> \
    "$exhaustive_copy/.github/release-notes/v0.9.0.md"
  run_expect_fail exhaustive-copy 'highlights-only release copy' \
    sh "$CHECK" "$exhaustive_copy"

  paragraph_copy="$tmp/paragraph-copy"
  cp -R "$good" "$paragraph_copy"
  printf '\nInstall the package and register the service before first use.\n' >> \
    "$paragraph_copy/.github/release-notes/v0.9.0.md"
  run_expect_fail paragraph-copy 'highlights-only release copy' \
    sh "$CHECK" "$paragraph_copy"

  h3_copy="$tmp/h3-copy"
  cp -R "$good" "$h3_copy"
  printf '\n### Installation\n\nRegister the service after extracting the package.\n' >> \
    "$h3_copy/.github/release-notes/v0.9.0.md"
  run_expect_fail h3-copy 'highlights-only release copy' \
    sh "$CHECK" "$h3_copy"

  old="$tmp/old"
  cp -R "$good" "$old"
  sed 's#actions/checkout@v7#actions/checkout@v4#' \
    "$good/.github/workflows/ci.yml" > "$old/.github/workflows/ci.yml"
  run_expect_fail obsolete 'actions/checkout@v4' sh "$CHECK" "$old"

  old_codeql="$tmp/old-codeql"
  cp -R "$good" "$old_codeql"
  sed 's#github/codeql-action/init@v4#github/codeql-action/init@v3#' \
    "$good/.github/workflows/codeql.yml" > \
    "$old_codeql/.github/workflows/codeql.yml"
  run_expect_fail obsolete-codeql 'github/codeql-action/init@v3' \
    sh "$CHECK" "$old_codeql"

  unknown="$tmp/unknown"
  cp -R "$good" "$unknown"
  printf '      - uses: example/unknown@v1\n' >> \
    "$unknown/.github/workflows/ci.yml"
  run_expect_fail unknown 'example/unknown@v1' sh "$CHECK" "$unknown"

  missing_trigger="$tmp/missing-trigger"
  cp -R "$good" "$missing_trigger"
  sed '/  schedule:/,/cron:/d' "$good/.github/workflows/codeql.yml" > \
    "$missing_trigger/.github/workflows/codeql.yml"
  run_expect_fail missing-trigger 'weekly schedule trigger' \
    sh "$CHECK" "$missing_trigger"

  insufficient_permission="$tmp/insufficient-permission"
  cp -R "$good" "$insufficient_permission"
  sed '/  security-events: write/d' "$good/.github/workflows/codeql.yml" > \
    "$insufficient_permission/.github/workflows/codeql.yml"
  run_expect_fail insufficient-permission 'security-events: write' \
    sh "$CHECK" "$insufficient_permission"

  job_permission="$tmp/job-permission"
  cp -R "$good" "$job_permission"
  sed '/    runs-on: ubuntu-latest/a\    permissions: write-all' \
    "$good/.github/workflows/codeql.yml" > \
    "$job_permission/.github/workflows/codeql.yml"
  run_expect_fail job-permission 'job-level permissions override' \
    sh "$CHECK" "$job_permission"

  invalid_schedule="$tmp/invalid-schedule"
  cp -R "$good" "$invalid_schedule"
  sed "s/17 4 \* \* 1/0 0 31 2 */" \
    "$good/.github/workflows/codeql.yml" > \
    "$invalid_schedule/.github/workflows/codeql.yml"
  run_expect_fail invalid-schedule 'weekly schedule trigger' \
    sh "$CHECK" "$invalid_schedule"

  missing_analysis="$tmp/missing-analysis"
  cp -R "$good" "$missing_analysis"
  sed '/github\/codeql-action\/analyze@v4/d' \
    "$good/.github/workflows/codeql.yml" > \
    "$missing_analysis/.github/workflows/codeql.yml"
  run_expect_fail missing-analysis 'CodeQL analyze step' \
    sh "$CHECK" "$missing_analysis"

  missing_gomod="$tmp/missing-gomod"
  cp -R "$good" "$missing_gomod"
  sed '/  - package-ecosystem: gomod/,/  - package-ecosystem: github-actions/{
    /  - package-ecosystem: github-actions/!d
  }' "$good/.github/dependabot.yml" > \
    "$missing_gomod/.github/dependabot.yml"
  run_expect_fail missing-gomod 'gomod entry' sh "$CHECK" "$missing_gomod"

  noisy_cadence="$tmp/noisy-cadence"
  cp -R "$good" "$noisy_cadence"
  sed '0,/interval: weekly/s//interval: daily/' \
    "$good/.github/dependabot.yml" > \
    "$noisy_cadence/.github/dependabot.yml"
  run_expect_fail noisy-cadence 'weekly cadence' sh "$CHECK" "$noisy_cadence"

  missing_label="$tmp/missing-dependency-label"
  cp -R "$good" "$missing_label"
  sed 's/      - dependencies/      - maintenance/' \
    "$good/.github/dependabot.yml" > \
    "$missing_label/.github/dependabot.yml"
  run_expect_fail missing-dependency-label 'dependencies label' \
    sh "$CHECK" "$missing_label"

  extra_ecosystem="$tmp/extra-ecosystem"
  cp -R "$good" "$extra_ecosystem"
  printf '%s\n' \
    '  - package-ecosystem: npm' \
    '    directory: /' \
    '    schedule:' \
    '      interval: weekly' >> "$extra_ecosystem/.github/dependabot.yml"
  run_expect_fail extra-ecosystem 'unapproved package ecosystem' \
    sh "$CHECK" "$extra_ecosystem"

  grouped_major="$tmp/grouped-major"
  cp -R "$good" "$grouped_major"
  sed '/          - minor/a\          - major' \
    "$good/.github/dependabot.yml" > \
    "$grouped_major/.github/dependabot.yml"
  run_expect_fail grouped-major 'must not group major updates' \
    sh "$CHECK" "$grouped_major"

  grouped_security="$tmp/grouped-security"
  cp -R "$good" "$grouped_security"
  sed '/      routine-minor-and-patch:/a\        applies-to: security-updates' \
    "$good/.github/dependabot.yml" > \
    "$grouped_security/.github/dependabot.yml"
  run_expect_fail grouped-security 'must not group security updates' \
    sh "$CHECK" "$grouped_security"

  bracket_secret="$tmp/bracket-secret"
  cp -R "$good" "$bracket_secret"
  sed "/    steps:/a\      - run: echo \"\${{ secrets['NAME'] }}\"" \
    "$good/.github/workflows/codeql.yml" > \
    "$bracket_secret/.github/workflows/codeql.yml"
  run_expect_fail bracket-secret 'must not consume' \
    sh "$CHECK" "$bracket_secret"

  stale_brand="$tmp/stale-brand"
  cp -R "$good" "$stale_brand"
  sed 's#brand/platform/macos/go-schedule.icns#gui/assets/icon.png#' \
    "$good/.github/workflows/release.yml" > \
    "$stale_brand/.github/workflows/release.yml"
  run_expect_fail stale-brand 'canonical macOS ICNS' \
    sh "$CHECK" "$stale_brand"

  missing_brand_check="$tmp/missing-brand-check"
  cp -R "$good" "$missing_brand_check"
  rm -rf "$missing_brand_check/scripts/brand-check"
  run_expect_fail missing-brand-check 'brand integrity command not found' \
    sh "$CHECK" "$missing_brand_check"

  missing="$tmp/missing"
  make_fixture "$missing" 'format vet lint race gui coverage docs'
  run_expect_fail missing 'gate manifest differs' sh "$CHECK" "$missing"

  duplicate="$tmp/duplicate"
  make_fixture "$duplicate" \
    'format vet lint race gui coverage docs automation automation'
  run_expect_fail duplicate 'gate manifest differs' sh "$CHECK" "$duplicate"

  extra="$tmp/extra"
  make_fixture "$extra" \
    'format vet lint race gui coverage docs automation build'
  run_expect_fail extra 'gate manifest differs' sh "$CHECK" "$extra"

  rm -rf "$tmp"
  trap - EXIT HUP INT TERM
}

run_verify_cases() {
  [ -f "$VERIFY" ] || fail "missing $VERIFY"

  run_expect_fail unknown-mode 'usage:' sh "$VERIFY" unknown

  if output=$(GOFMT=false sh "$VERIFY" all 2>&1); then
    fail 'controlled aggregate child failure should fail'
  fi
  assert_contains "$output" '[format]'
  case "$output" in
    *'[vet]'*) fail 'aggregate continued after the controlled format failure' ;;
    *) : ;;
  esac
}

case "$MODE" in
  all)
    run_automation_cases
    run_verify_cases
    ;;
  automation)
    run_automation_cases
    ;;
  failure-propagation)
    run_verify_cases
    ;;
  *)
    fail "unknown test mode: $MODE"
    ;;
esac

printf 'automation-check-test: OK (%s)\n' "$MODE"
