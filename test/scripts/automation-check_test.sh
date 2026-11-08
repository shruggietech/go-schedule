#!/bin/sh
# Contract regressions for the maintainer automation baseline.

set -eu

ROOT=$(CDPATH='' cd "$(dirname "$0")/../.." && pwd)
CHECK="$ROOT/scripts/automation-check.sh"
MODE=${1:-all}
CASE_COUNT=0

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
  CASE_COUNT=$((CASE_COUNT + 1))
  label=$1
  shift
  if ! output=$("$@" 2>&1); then
    fail "$label should pass; got: $output"
  fi
}

run_expect_fail() {
  CASE_COUNT=$((CASE_COUNT + 1))
  label=$1
  needle=$2
  shift 2
  if output=$("$@" 2>&1); then
    fail "$label should fail"
  fi
  assert_contains "$output" "$needle"
}

. "$ROOT/test/scripts/lib/automation-fixture.sh"

run_automation_cases() {
  [ -f "$CHECK" ] || fail "missing $CHECK"

  tmp=$(mktemp -d "$ROOT/.automation-check-test.XXXXXX")
  trap 'rm -rf "$tmp"' EXIT HUP INT TERM
  (cd "$ROOT" && "${GO:-go}" build -o "$tmp/maintainer-check.exe" ./scripts/maintainer-check)
  CHECK="$tmp/maintainer-check.exe"
  good="$tmp/good"
  make_fixture "$good" 'format vet lint race gui coverage docs automation'

  run_expect_pass approved "$CHECK" automation "$good"

  for surface in runtime makefile verify ci release mixed-ci; do
    mismatch="$tmp/wails-$surface"
    cp -R "$good" "$mismatch"
    case "$surface" in
      runtime) target=desktop/go.mod ;;
      makefile) target=Makefile ;;
      verify) target=scripts/verify.sh ;;
      ci) target=.github/workflows/ci.yml ;;
      release) target=.github/workflows/release.yml ;;
      mixed-ci) target=.github/workflows/ci.yml ;;
    esac
    if [ "$surface" = mixed-ci ]; then
      printf '\n      - run: go run github.com/wailsapp/wails/v2/cmd/wails@v2.16.0 build\n' >> "$mismatch/$target"
    else
      sed 's/v2.17.1/v2.16.0/g' "$good/$target" > "$mismatch/$target"
    fi
    run_expect_fail "wails-$surface" 'Wails' "$CHECK" automation "$mismatch"
  done

  for prefix in stale current empty; do
    repeated="$tmp/wails-repeated-$prefix"
    cp -R "$good" "$repeated"
    case "$prefix" in
      stale) pin=v2.16.0 ;;
      current) pin=v2.17.1 ;;
      empty) pin='' ;;
    esac
    printf 'WAILS_VERSION = %s\n' "$pin" > "$repeated/Makefile"
    cat "$good/Makefile" >> "$repeated/Makefile"
    run_expect_fail "wails-repeated-$prefix" 'Wails build version' "$CHECK" automation "$repeated"
  done

  for release_file in release.yml promote-release.yml promote-release-waived.yml; do
    for defect in missing bypass late conditional soft-failure; do
      comparison_fixture="$tmp/comparison-$release_file-$defect"
      cp -R "$good" "$comparison_fixture"
      target="$comparison_fixture/.github/workflows/$release_file"
      case "$defect" in
        missing) sed '/comparison-check/d' "$good/.github/workflows/$release_file" > "$target" ;;
        bypass) sed '/comparison-check/s/$/ || true/' "$good/.github/workflows/$release_file" > "$target" ;;
        conditional) sed '/comparison-check/a\        if: false' "$good/.github/workflows/$release_file" > "$target" ;;
        soft-failure) sed '/comparison-check/a\        continue-on-error: true' "$good/.github/workflows/$release_file" > "$target" ;;
        late)
          sed '/comparison-check/d' "$good/.github/workflows/$release_file" > "$target"
          grep 'go run ./scripts/comparison-check' "$good/.github/workflows/$release_file" >> "$target"
          ;;
      esac
      run_expect_fail "comparison-$release_file-$defect" 'comparison review' "$CHECK" automation "$comparison_fixture"
    done
  done


  for defect in missing bypass conditional soft-failure wrong-arguments late no-go; do
    source_fixture="$tmp/source-preflight-$defect"
    cp -R "$good" "$source_fixture"
    target="$source_fixture/.github/workflows/release.yml"
    case "$defect" in
      missing) sed '/release-source-check/d' "$good/.github/workflows/release.yml" > "$target" ;;
      bypass) sed '/release-source-check/s/$/ || true/' "$good/.github/workflows/release.yml" > "$target" ;;
      conditional) sed '/release-source-check/a\        if: false' "$good/.github/workflows/release.yml" > "$target" ;;
      soft-failure) sed '/release-source-check/a\        continue-on-error: true' "$good/.github/workflows/release.yml" > "$target" ;;
      wrong-arguments) sed '/release-source-check/s/-root ./-root other/' "$good/.github/workflows/release.yml" > "$target" ;;
      no-go) sed '/actions\/setup-go@/d' "$good/.github/workflows/release.yml" > "$target" ;;
      late)
        sed '/release-source-check/d' "$good/.github/workflows/release.yml" > "$target"
        grep 'run: go run ./scripts/release-source-check' "$good/.github/workflows/release.yml" >> "$target"
        ;;
    esac
    run_expect_fail "source-preflight-$defect" 'release source preflight' "$CHECK" automation "$source_fixture"
  done

  stale_tag_boundary="$tmp/stale-tag-boundary"
  cp -R "$good" "$stale_tag_boundary"
  sed 's/reviewed S049 merge commit/reviewed S048 merge commit/' \
    "$good/specs/048-v100-release-cut/contracts/publication.md" > \
    "$stale_tag_boundary/specs/048-v100-release-cut/contracts/publication.md"
  run_expect_fail stale-tag-boundary 'S049 tag-boundary identity' \
    "$CHECK" automation "$stale_tag_boundary"

  missing_packet_command="$tmp/missing-packet-command"
  cp -R "$good" "$missing_packet_command"
  sed 's/windows-release-gate render-dispositions/manual packet copying/' \
    "$good/test/windows/README.md" > \
    "$missing_packet_command/test/windows/README.md"
  run_expect_fail missing-packet-command 'operator disposition command' \
    "$CHECK" automation "$missing_packet_command"

  missing_non_authority="$tmp/missing-non-authority"
  cp -R "$good" "$missing_non_authority"
  sed 's/does not update or close/is complete and may close/' \
    "$good/test/windows/README.md" > \
    "$missing_non_authority/test/windows/README.md"
  run_expect_fail missing-non-authority 'non-authority boundary' \
    "$CHECK" automation "$missing_non_authority"

  generated_notes="$tmp/generated-notes"
  cp -R "$good" "$generated_notes"
  sed 's/generate_release_notes: false/generate_release_notes: true/' \
    "$good/.github/workflows/release.yml" > \
    "$generated_notes/.github/workflows/release.yml"
  run_expect_fail generated-notes 'generated release notes must remain disabled' \
    "$CHECK" automation "$generated_notes"

  fixed_body_path="$tmp/fixed-body-path"
  cp -R "$good" "$fixed_body_path"
  sed 's/${{ github.ref_name }}/v0.9.0/' \
    "$good/.github/workflows/release.yml" > \
    "$fixed_body_path/.github/workflows/release.yml"
  run_expect_fail fixed-body-path 'dynamic tag-specific release-note body path' \
    "$CHECK" automation "$fixed_body_path"

  direct_main_push="$tmp/direct-main-push"
  cp -R "$good" "$direct_main_push"
  printf '      - run: git push origin HEAD:main\n' >> \
    "$direct_main_push/.github/workflows/release.yml"
  run_expect_fail direct-main-push 'must not invoke git' \
    "$CHECK" automation "$direct_main_push"

  quoted_main_push="$tmp/quoted-main-push"
  cp -R "$good" "$quoted_main_push"
  printf "      - run: git push origin 'HEAD:main'\n" >> \
    "$quoted_main_push/.github/workflows/release.yml"
  run_expect_fail quoted-main-push 'must not invoke git' \
    "$CHECK" automation "$quoted_main_push"

  full_main_refspec="$tmp/full-main-refspec"
  cp -R "$good" "$full_main_refspec"
  printf '      - run: git push origin HEAD:refs/heads/main\n' >> \
    "$full_main_refspec/.github/workflows/release.yml"
  run_expect_fail full-main-refspec 'must not invoke git' \
    "$CHECK" automation "$full_main_refspec"

  global_option_push="$tmp/global-option-push"
  cp -R "$good" "$global_option_push"
  # Preserve the variable reference in the generated workflow.
  # shellcheck disable=SC2016
  printf '      - run: git -C "$GITHUB_WORKSPACE" push origin HEAD:main\n' >> \
    "$global_option_push/.github/workflows/release.yml"
  run_expect_fail global-option-push 'must not invoke git' \
    "$CHECK" automation "$global_option_push"

  continued_push="$tmp/continued-push"
  cp -R "$good" "$continued_push"
  # A literal backslash produces the continued-command negative fixture.
  # shellcheck disable=SC1003
  printf '%s\n' '      - run: |' '          git \' \
    '            push origin HEAD:main' >> \
    "$continued_push/.github/workflows/release.yml"
  run_expect_fail continued-push 'must not invoke git' \
    "$CHECK" automation "$continued_push"

  folded_push="$tmp/folded-push"
  cp -R "$good" "$folded_push"
  printf '%s\n' '      - run: >' '          git' \
    '          push origin HEAD:main' >> \
    "$folded_push/.github/workflows/release.yml"
  run_expect_fail folded-push 'must not invoke git' \
    "$CHECK" automation "$folded_push"

  main_checkout="$tmp/main-checkout"
  cp -R "$good" "$main_checkout"
  printf '      - uses: actions/checkout@v7\n        with:\n          ref: main\n' >> \
    "$main_checkout/.github/workflows/release.yml"
  run_expect_fail main-checkout 'must inspect the tagged checkout' \
    "$CHECK" automation "$main_checkout"

  quoted_main_checkout="$tmp/quoted-main-checkout"
  cp -R "$good" "$quoted_main_checkout"
  printf '      - uses: actions/checkout@v7\n        with:\n          ref: "main"\n' >> \
    "$quoted_main_checkout/.github/workflows/release.yml"
  run_expect_fail quoted-main-checkout 'must inspect the tagged checkout' \
    "$CHECK" automation "$quoted_main_checkout"

  qualified_main_checkout="$tmp/qualified-main-checkout"
  cp -R "$good" "$qualified_main_checkout"
  printf "      - uses: actions/checkout@v7\n        with:\n          ref: 'refs/heads/main'\n" >> \
    "$qualified_main_checkout/.github/workflows/release.yml"
  run_expect_fail qualified-main-checkout 'must inspect the tagged checkout' \
    "$CHECK" automation "$qualified_main_checkout"

  ungated_release="$tmp/ungated-release"
  cp -R "$good" "$ungated_release"
  sed 's/needs: \[readme-version, ci-success\]/needs: readme-version/' \
    "$good/.github/workflows/release.yml" > \
    "$ungated_release/.github/workflows/release.yml"
  printf '  gui: # desktop release\n    needs: readme-version\n' >> \
    "$ungated_release/.github/workflows/release.yml"
  run_expect_fail ungated-release 'README and exact-commit CI preflight dependencies' \
    "$CHECK" automation "$ungated_release"

  nested_dependency="$tmp/nested-dependency"
  cp -R "$good" "$nested_dependency"
  sed 's/^    needs: \[readme-version, ci-success\]$/    env:\
      needs: [readme-version, ci-success]/' \
    "$good/.github/workflows/release.yml" > \
    "$nested_dependency/.github/workflows/release.yml"
  run_expect_fail nested-dependency 'README and exact-commit CI preflight dependencies' \
    "$CHECK" automation "$nested_dependency"

  missing_webkit="$tmp/missing-webkit"
  cp -R "$good" "$missing_webkit"
  sed 's/libwebkit2gtk-4.1-dev//' "$good/.github/workflows/release.yml" > \
    "$missing_webkit/.github/workflows/release.yml"
  run_expect_fail missing-webkit 'Linux Wails WebKit development headers' \
    "$CHECK" automation "$missing_webkit"

  public_staging="$tmp/public-staging"
  cp -R "$good" "$public_staging"
  sed '/draft: true/d' "$good/.github/workflows/release.yml" > \
    "$public_staging/.github/workflows/release.yml"
  run_expect_fail public-staging 'explicitly use draft: true' \
    "$CHECK" automation "$public_staging"

  staging_checksums="$tmp/staging-checksums"
  cp -R "$good" "$staging_checksums"
  printf '      - run: sha256sum * > SHA256SUMS.txt\n' >> \
    "$staging_checksums/.github/workflows/release.yml"
  run_expect_fail staging-checksums 'must not publish or create final checksums' \
    "$CHECK" automation "$staging_checksums"

  missing_promotion="$tmp/missing-promotion"
  cp -R "$good" "$missing_promotion"
  rm -f "$missing_promotion/.github/workflows/promote-release.yml"
  run_expect_fail missing-promotion 'promotion workflow not found' \
    "$CHECK" automation "$missing_promotion"

  rebuild_promotion="$tmp/rebuild-promotion"
  cp -R "$good" "$rebuild_promotion"
  printf '          go build ./cmd/gosched-gui\n' >> \
    "$rebuild_promotion/.github/workflows/promote-release.yml"
  run_expect_fail rebuild-promotion 'must not rebuild release artifacts' \
    "$CHECK" automation "$rebuild_promotion"

  ungated_promotion="$tmp/ungated-promotion"
  cp -R "$good" "$ungated_promotion"
  sed '/windows-release-gate verify-bundle/d' \
    "$good/.github/workflows/promote-release.yml" > \
    "$ungated_promotion/.github/workflows/promote-release.yml"
  run_expect_fail ungated-promotion 'shared exact-candidate evidence gate' \
    "$CHECK" automation "$ungated_promotion"

  no_candidate_gate="$tmp/no-candidate-gate"
  cp -R "$good" "$no_candidate_gate"
  sed '/windows-release-gate verify-candidate/d' \
    "$good/.github/workflows/promote-release.yml" > \
    "$no_candidate_gate/.github/workflows/promote-release.yml"
  run_expect_fail no-candidate-gate 'independent candidate-manifest gate' \
    "$CHECK" automation "$no_candidate_gate"

  no_actions_read="$tmp/no-actions-read"
  cp -R "$good" "$no_actions_read"
  sed '/actions: read/d' \
    "$good/.github/workflows/promote-release.yml" > \
    "$no_actions_read/.github/workflows/promote-release.yml"
  run_expect_fail no-actions-read 'workflow-run read permission' \
    "$CHECK" automation "$no_actions_read"

  no_asset_allowlist="$tmp/no-asset-allowlist"
  cp -R "$good" "$no_asset_allowlist"
  sed '/expected-assets.txt/d' \
    "$good/.github/workflows/promote-release.yml" > \
    "$no_asset_allowlist/.github/workflows/promote-release.yml"
  run_expect_fail no-asset-allowlist 'exact draft asset allowlist' \
    "$CHECK" automation "$no_asset_allowlist"

  no_attempt_job="$tmp/no-attempt-job"
  cp -R "$good" "$no_attempt_job"
  sed '/attempts\/.*\/jobs?per_page=100/d' \
    "$good/.github/workflows/promote-release.yml" > \
    "$no_attempt_job/.github/workflows/promote-release.yml"
  run_expect_fail no-attempt-job \
    'candidate-producing workflow-attempt job validation' \
    "$CHECK" automation "$no_attempt_job"

  no_tag_recheck="$tmp/no-tag-recheck"
  cp -R "$good" "$no_tag_recheck"
  sed '/git ls-remote origin/d' \
    "$good/.github/workflows/promote-release.yml" > \
    "$no_tag_recheck/.github/workflows/promote-release.yml"
  run_expect_fail no-tag-recheck 'last-moment tag identity check' \
    "$CHECK" automation "$no_tag_recheck"

  early_promotion="$tmp/early-promotion"
  cp -R "$good" "$early_promotion"
  sed '/gh api -F draft=false/d; /go run .\/scripts\/windows-release-gate verify-candidate/i\
          gh api -F draft=false' \
    "$good/.github/workflows/promote-release.yml" > \
    "$early_promotion/.github/workflows/promote-release.yml"
  run_expect_fail early-promotion 'not ordered safely' \
    "$CHECK" automation "$early_promotion"

  missing_changelog="$tmp/missing-changelog"
  cp -R "$good" "$missing_changelog"
  sed '/Read the \[full changelog\]/d' \
    "$good/.github/release-notes/v0.9.0.md" > \
    "$missing_changelog/.github/release-notes/v0.9.0.md"
  run_expect_fail missing-changelog 'tagged full changelog link' \
    "$CHECK" automation "$missing_changelog"

  too_few_highlights="$tmp/too-few-highlights"
  cp -R "$good" "$too_few_highlights"
  sed '/consistent visual identity/d' \
    "$good/.github/release-notes/v0.9.0.md" > \
    "$too_few_highlights/.github/release-notes/v0.9.0.md"
  run_expect_fail too-few-highlights 'four to six highlight bullets' \
    "$CHECK" automation "$too_few_highlights"

  too_many_highlights="$tmp/too-many-highlights"
  cp -R "$good" "$too_many_highlights"
  sed '/Read the \[full changelog\]/i\
- Improve local daemon maintenance.\
- Refresh contributor automation.\
- Strengthen release packaging checks.' \
    "$good/.github/release-notes/v0.9.0.md" > \
    "$too_many_highlights/.github/release-notes/v0.9.0.md"
  run_expect_fail too-many-highlights 'four to six highlight bullets' \
    "$CHECK" automation "$too_many_highlights"

  exhaustive_copy="$tmp/exhaustive-copy"
  cp -R "$good" "$exhaustive_copy"
  printf '\n## Installation\n\nDownload the package for your platform.\n' >> \
    "$exhaustive_copy/.github/release-notes/v0.9.0.md"
  run_expect_fail exhaustive-copy 'highlights-only release copy' \
    "$CHECK" automation "$exhaustive_copy"

  paragraph_copy="$tmp/paragraph-copy"
  cp -R "$good" "$paragraph_copy"
  printf '\nInstall the package and register the service before first use.\n' >> \
    "$paragraph_copy/.github/release-notes/v0.9.0.md"
  run_expect_fail paragraph-copy 'highlights-only release copy' \
    "$CHECK" automation "$paragraph_copy"

  h3_copy="$tmp/h3-copy"
  cp -R "$good" "$h3_copy"
  printf '\n### Installation\n\nRegister the service after extracting the package.\n' >> \
    "$h3_copy/.github/release-notes/v0.9.0.md"
  run_expect_fail h3-copy 'highlights-only release copy' \
    "$CHECK" automation "$h3_copy"

  old="$tmp/old"
  cp -R "$good" "$old"
  sed 's#actions/checkout@v7#actions/checkout@v4#' \
    "$good/.github/workflows/ci.yml" > "$old/.github/workflows/ci.yml"
  run_expect_fail obsolete 'actions/checkout@v4' "$CHECK" automation "$old"

  old_node="$tmp/old-node"
  cp -R "$good" "$old_node"
  sed 's#actions/setup-node@v7#actions/setup-node@v6#' \
    "$good/.github/workflows/ci.yml" > "$old_node/.github/workflows/ci.yml"
  run_expect_fail obsolete-node 'actions/setup-node@v6' "$CHECK" automation "$old_node"

  for action in configure-pages deploy-pages upload-pages-artifact; do
    case "$action" in configure-pages) former=v5 ;; *) former=v4 ;; esac
    old_pages="$tmp/old-$action"
    cp -R "$good" "$old_pages"
    printf '      - uses: actions/%s@%s\n' "$action" "$former" >> "$old_pages/.github/workflows/ci.yml"
    run_expect_fail "obsolete-$action" "actions/$action@$former" "$CHECK" automation "$old_pages"
  done

  old_node_runtime="$tmp/old-node-runtime"
  cp -R "$good" "$old_node_runtime"
  awk '
    !changed && /node-version: 26/ { sub(/node-version: 26/, "node-version: 24"); changed = 1 }
    { print }
  ' "$good/.github/workflows/ci.yml" > "$old_node_runtime/.github/workflows/ci.yml"
  run_expect_fail old-node-runtime 'every CI Node runtime pin must use Node 26; found 24' \
    "$CHECK" automation "$old_node_runtime"

  old_release_node_runtime="$tmp/old-release-node-runtime"
  cp -R "$good" "$old_release_node_runtime"
  awk '
    !changed && /node-version: 26/ { sub(/node-version: 26/, "node-version: 24"); changed = 1 }
    { print }
  ' "$good/.github/workflows/release.yml" > \
    "$old_release_node_runtime/.github/workflows/release.yml"
  run_expect_fail old-release-node-runtime 'every release Node runtime pin must use Node 26; found 24' \
    "$CHECK" automation "$old_release_node_runtime"

  incomplete_wails_matrix="$tmp/incomplete-wails-matrix"
  cp -R "$good" "$incomplete_wails_matrix"
  sed 's/os: \[ubuntu-latest, macos-latest, windows-latest\]/os: [ubuntu-latest, windows-latest]/' \
    "$good/.github/workflows/ci.yml" > "$incomplete_wails_matrix/.github/workflows/ci.yml"
  run_expect_fail incomplete-wails-matrix 'three-platform production Wails matrix' \
    "$CHECK" automation "$incomplete_wails_matrix"

  missing_production_desktop="$tmp/missing-production-desktop"
  cp -R "$good" "$missing_production_desktop"
  sed '/^  wails-desktop:$/,/^  wails-desktop-browser-contract:$/ { /^  wails-desktop-browser-contract:$/!d; }' \
    "$good/.github/workflows/ci.yml" > "$missing_production_desktop/.github/workflows/ci.yml"
  run_expect_fail missing-production-desktop 'production Wails desktop job' \
    "$CHECK" automation "$missing_production_desktop"

  missing_v13_qualification="$tmp/missing-v13-qualification"
  cp -R "$good" "$missing_v13_qualification"
  sed '/^  v13-release-qualification:$/,$d' \
    "$good/.github/workflows/ci.yml" > "$missing_v13_qualification/.github/workflows/ci.yml"
  run_expect_fail missing-v13-qualification 'named v1.3 release qualification job' \
    "$CHECK" automation "$missing_v13_qualification"

  incomplete_v13_matrix="$tmp/incomplete-v13-matrix"
  cp -R "$good" "$incomplete_v13_matrix"
  awk '
    $0 == "  v13-release-qualification:" { inside = 1 }
    inside && $0 == "        os: [ubuntu-latest, macos-latest, windows-latest]" {
      print "        os: [ubuntu-latest]"
      next
    }
    { print }
  ' "$good/.github/workflows/ci.yml" > "$incomplete_v13_matrix/.github/workflows/ci.yml"
  run_expect_fail incomplete-v13-matrix 'v13-release-qualification job missing three-platform matrix' \
    "$CHECK" automation "$incomplete_v13_matrix"

  old_codeql="$tmp/old-codeql"
  cp -R "$good" "$old_codeql"
  sed 's#github/codeql-action/init@v4#github/codeql-action/init@v3#' \
    "$good/.github/workflows/codeql.yml" > \
    "$old_codeql/.github/workflows/codeql.yml"
  run_expect_fail obsolete-codeql 'github/codeql-action/init@v3' \
    "$CHECK" automation "$old_codeql"

  unknown="$tmp/unknown"
  cp -R "$good" "$unknown"
  printf '      - uses: example/unknown@v1\n' >> \
    "$unknown/.github/workflows/ci.yml"
  run_expect_fail unknown 'example/unknown@v1' "$CHECK" automation "$unknown"

  missing_trigger="$tmp/missing-trigger"
  cp -R "$good" "$missing_trigger"
  sed '/  schedule:/,/cron:/d' "$good/.github/workflows/codeql.yml" > \
    "$missing_trigger/.github/workflows/codeql.yml"
  run_expect_fail missing-trigger 'weekly schedule trigger' \
    "$CHECK" automation "$missing_trigger"

  insufficient_permission="$tmp/insufficient-permission"
  cp -R "$good" "$insufficient_permission"
  sed '/  security-events: write/d' "$good/.github/workflows/codeql.yml" > \
    "$insufficient_permission/.github/workflows/codeql.yml"
  run_expect_fail insufficient-permission 'security-events: write' \
    "$CHECK" automation "$insufficient_permission"

  job_permission="$tmp/job-permission"
  cp -R "$good" "$job_permission"
  sed '/    runs-on: ubuntu-latest/a\    permissions: write-all' \
    "$good/.github/workflows/codeql.yml" > \
    "$job_permission/.github/workflows/codeql.yml"
  run_expect_fail job-permission 'job-level permissions override' \
    "$CHECK" automation "$job_permission"

  invalid_schedule="$tmp/invalid-schedule"
  cp -R "$good" "$invalid_schedule"
  sed "s/17 4 \* \* 1/0 0 31 2 */" \
    "$good/.github/workflows/codeql.yml" > \
    "$invalid_schedule/.github/workflows/codeql.yml"
  run_expect_fail invalid-schedule 'weekly schedule trigger' \
    "$CHECK" automation "$invalid_schedule"

  missing_analysis="$tmp/missing-analysis"
  cp -R "$good" "$missing_analysis"
  sed '/github\/codeql-action\/analyze@v4/d' \
    "$good/.github/workflows/codeql.yml" > \
    "$missing_analysis/.github/workflows/codeql.yml"
  run_expect_fail missing-analysis 'CodeQL analyze step' \
    "$CHECK" automation "$missing_analysis"

  missing_gomod="$tmp/missing-gomod"
  cp -R "$good" "$missing_gomod"
  sed '/  - package-ecosystem: gomod/,/  - package-ecosystem: github-actions/{
    /  - package-ecosystem: github-actions/!d
  }' "$good/.github/dependabot.yml" > \
    "$missing_gomod/.github/dependabot.yml"
  run_expect_fail missing-gomod 'gomod entry' "$CHECK" automation "$missing_gomod"

  noisy_cadence="$tmp/noisy-cadence"
  cp -R "$good" "$noisy_cadence"
  sed '0,/interval: weekly/s//interval: daily/' \
    "$good/.github/dependabot.yml" > \
    "$noisy_cadence/.github/dependabot.yml"
  run_expect_fail noisy-cadence 'weekly cadence' "$CHECK" automation "$noisy_cadence"

  missing_label="$tmp/missing-dependency-label"
  cp -R "$good" "$missing_label"
  sed 's/      - dependencies/      - maintenance/' \
    "$good/.github/dependabot.yml" > \
    "$missing_label/.github/dependabot.yml"
  run_expect_fail missing-dependency-label 'dependencies label' \
    "$CHECK" automation "$missing_label"

  extra_ecosystem="$tmp/extra-ecosystem"
  cp -R "$good" "$extra_ecosystem"
  printf '%s\n' \
    '  - package-ecosystem: npm' \
    '    directory: /' \
    '    schedule:' \
    '      interval: weekly' >> "$extra_ecosystem/.github/dependabot.yml"
  run_expect_fail extra-ecosystem 'expected root Go, desktop Go, desktop npm, website npm, and GitHub Actions entries' \
    "$CHECK" automation "$extra_ecosystem"

  grouped_major="$tmp/grouped-major"
  cp -R "$good" "$grouped_major"
  sed '/          - minor/a\          - major' \
    "$good/.github/dependabot.yml" > \
    "$grouped_major/.github/dependabot.yml"
  run_expect_fail grouped-major 'must not group major updates' \
    "$CHECK" automation "$grouped_major"

  grouped_security="$tmp/grouped-security"
  cp -R "$good" "$grouped_security"
  sed '/      routine-minor-and-patch:/a\        applies-to: security-updates' \
    "$good/.github/dependabot.yml" > \
    "$grouped_security/.github/dependabot.yml"
  run_expect_fail grouped-security 'must not group security updates' \
    "$CHECK" automation "$grouped_security"

  bracket_secret="$tmp/bracket-secret"
  cp -R "$good" "$bracket_secret"
  sed "/    steps:/a\      - run: echo \"\${{ secrets['NAME'] }}\"" \
    "$good/.github/workflows/codeql.yml" > \
    "$bracket_secret/.github/workflows/codeql.yml"
  run_expect_fail bracket-secret 'must not consume' \
    "$CHECK" automation "$bracket_secret"

  stale_brand="$tmp/stale-brand"
  cp -R "$good" "$stale_brand"
  sed 's#brand/platform/macos/go-schedule.icns#gui/assets/icon.png#' \
    "$good/.github/workflows/release.yml" > \
    "$stale_brand/.github/workflows/release.yml"
  run_expect_fail stale-brand 'canonical macOS ICNS' \
    "$CHECK" automation "$stale_brand"

  missing_brand_check="$tmp/missing-brand-check"
  cp -R "$good" "$missing_brand_check"
  rm -rf "$missing_brand_check/scripts/brand-check"
  run_expect_fail missing-brand-check 'brand integrity command not found' \
    "$CHECK" automation "$missing_brand_check"

  missing="$tmp/missing"
  make_fixture "$missing" 'format vet lint race gui coverage docs'
  run_expect_fail missing 'gate manifest differs' "$CHECK" automation "$missing"

  duplicate="$tmp/duplicate"
  make_fixture "$duplicate" \
    'format vet lint race gui coverage docs automation automation'
  run_expect_fail duplicate 'gate manifest differs' "$CHECK" automation "$duplicate"

  extra="$tmp/extra"
  make_fixture "$extra" \
    'format vet lint race gui coverage docs automation build'
  run_expect_fail extra 'gate manifest differs' "$CHECK" automation "$extra"

  rm -rf "$tmp"
  trap - EXIT HUP INT TERM
}

run_verify_cases() {
  sh "$ROOT/test/scripts/verify-driver_test.sh"
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

printf 'automation-check-test: OK (%s), %s assertion case(s)\n' "$MODE" "$CASE_COUNT"
