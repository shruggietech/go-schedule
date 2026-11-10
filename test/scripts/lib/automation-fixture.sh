# Source-only shared automation/lifecycle fixture writer.

make_fixture() (
  fixture=$1
  manifest=$2
  mkdir -p "$fixture/.github/workflows" \
    "$fixture/desktop" \
    "$fixture/.github/release-notes" \
    "$fixture/scripts/brand-check" \
    "$fixture/specs/048-v100-release-cut/contracts" \
    "$fixture/test/windows"
  mkdir -p "$fixture/scripts/desktop-build"
  printf 'package main\n// github.com/wailsapp/wails/v2/cmd/wails@v2.17.1\n' > "$fixture/scripts/desktop-build/main.go"
  cat > "$fixture/go.mod" <<'EOF'
module fixture

go 1.25.0
EOF
  cat > "$fixture/desktop/go.mod" <<'EOF'
module fixture/desktop

go 1.25.0

require github.com/wailsapp/wails/v2 v2.17.1
EOF
  cat > "$fixture/Makefile" <<'EOF'
WAILS_VERSION = v2.17.1
EOF
  cat > "$fixture/scripts/brand-check/main.go" <<'EOF'
package main
func main() {}
EOF
  cat > "$fixture/.github/workflows/ci.yml" <<'EOF'
name: fixture
jobs:
  test:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
    steps:
      - name: Maintainer command contracts
        shell: bash
        run: sh test/scripts/maintainer-commands_test.sh
      - uses: actions/checkout@v7
      - uses: actions/setup-go@v7
      - uses: actions/setup-node@v7
        with:
          node-version: 26
      - uses: actions/upload-artifact@v7
      - uses: softprops/action-gh-release@v3
  wails-desktop:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
    defaults:
      run:
        working-directory: desktop
    steps:
      - run: go test -timeout=60m -race ./...
      - run: go run ./scripts/desktop-build
      - run: go run github.com/wailsapp/wails/v2/cmd/wails@v2.17.1 build
      - run: npm audit --audit-level=high && npm test && npm run build
      - run: 'echo cache-dependency-path: desktop/go.sum'
      - run: 'echo cache-dependency-path: desktop/frontend/package-lock.json'
      - run: echo desktop/build/bin/gosched-gui.exe
      - run: echo "gosched-tray.exe') ./cmd/gosched-tray"
      - run: go build -o desktop/build/bin/gosched-indicator ./cmd/gosched-indicator
      - run: echo 'Copy-Item README.md, LICENSE, CHANGELOG.md -Destination $stage'
      - name: Inspect stable desktop identity
        run: test -f build/bin/go-schedule.app/Contents/MacOS/gosched-gui
  wails-desktop-browser-contract:
    steps:
      - run: npm run test:e2e
  v13-release-qualification:
    name: v1.3 release qualification (${{ matrix.os }})
    strategy:
      fail-fast: false
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
    steps:
      - run: go test -timeout=60m -race ./test/integration -run '^(TestV13PackageDefaultsRemainOptIn|TestWebhookNotificationEndToEndPreservesRunOutcome|TestPackagedMCPCommandDiscovery)$' -count=1
      - run: go test -timeout=60m -race ./internal/notification ./internal/mcpobserve ./internal/mcphttp ./internal/api/server ./internal/api/client ./internal/cli
      - run: go test -timeout=60m -race ./internal/store
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
  - package-ecosystem: gomod
    directory: /desktop
    schedule:
      interval: weekly
    open-pull-requests-limit: 5
    groups:
      desktop-minor-and-patch:
        patterns:
          - github.com/wailsapp/wails/v2
        update-types:
          - minor
          - patch
      routine-minor-and-patch:
        update-types:
          - minor
          - patch
    labels:
      - dependencies
  - package-ecosystem: npm
    directory: /desktop/frontend
    schedule:
      interval: weekly
    open-pull-requests-limit: 5
    groups:
      routine-minor-and-patch:
        update-types:
          - minor
          - patch
    labels:
      - dependencies
  - package-ecosystem: npm
    directory: /site
    schedule:
      interval: weekly
    open-pull-requests-limit: 5
    groups:
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
# Real lifecycle data below is evaluated in-process.
EOF
  mkdir -p "$fixture/specs/001-fixture"
  printf "**Status**: Implemented\n**Delivery**: commit abc1234\n" > "$fixture/specs/001-fixture/spec.md"
  printf -- "- [x] T001 Fixture task\n" > "$fixture/specs/001-fixture/tasks.md"
  printf "| 001-fixture | Implemented | commit abc1234 |\n" > "$fixture/specs/README.md"
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
  ci-success:
    name: Require successful CI for tagged commit
    permissions:
      actions: read
      contents: read
    steps:
      - run: |
          RESPONSE="https://api.github.com/repos/${GITHUB_REPOSITORY}/actions/workflows/ci.yml/runs?head_sha=${GITHUB_SHA}&event=push"
          RUN='.head_branch == "main" and .event == "push"'
          if [ "$STATUS" = completed ]; then
            exit 1
          fi
          echo "Timed out waiting for successful main CI"
  readme-version:
    steps:
      - uses: actions/setup-go@v7
      - name: Verify release source metadata
        run: go run ./scripts/release-source-check -root . -repository "$GITHUB_REPOSITORY" -tag "$GITHUB_REF_NAME"
      - run: go run ./scripts/comparison-check -tag "$GITHUB_REF_NAME" -commit "$GITHUB_SHA"
  release-state:
    needs: [readme-version, ci-success]
    steps:
      - name: Refuse to mutate an already-public release
        run: test "$(printf '%s' "$BODY" | jq -r '.draft')" = true
  binaries:
    needs: release-state
    steps:
      - uses: actions/setup-node@v7
        with:
          node-version: 26
      - run: 'echo go-version-file: desktop/go.mod'
      - run: 'echo cache-dependency-path: desktop/frontend/package-lock.json'
      - run: go run ./scripts/desktop-build
      - run: go run github.com/wailsapp/wails/v2/cmd/wails@v2.17.1 build
      - run: cp desktop/build/bin/gosched-gui.exe "$stage/gosched-gui.exe"
      - run: go build -o "$stage/gosched-tray.exe" ./cmd/gosched-tray
      - run: go build -o "$stage/gosched-indicator" ./cmd/gosched-indicator
      - run: cp -R desktop/build/bin/go-schedule.app "$app"
      - run: app="$stage/gosched-gui.app"
      - run: cp brand/platform/macos/go-schedule.icns "$app/Contents/Resources/icon.icns"
      - run: cp brand/platform/linux/go-schedule.desktop "$stage/share/applications/"
      - run: cp brand/platform/linux/go-schedule-indicator.desktop "$stage/share/autostart/"
      - run: cp -R brand/platform/linux/hicolor "$stage/share/icons/"
      - run: sudo apt-get install -y libwebkit2gtk-4.1-dev
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
          gh api "repos/x/actions/runs/${RUN_ID}/attempts/${RUN_ATTEMPT}/jobs?per_page=100"
          Build & stage GUI (windows-latest)
          go run ./scripts/comparison-check -tag "$RELEASE_TAG" -commit "$TAG_COMMIT"
          LC_ALL=C sha256sum * > SHA256SUMS.txt
          sha256sum -c SHA256SUMS.txt
           gh release upload "$RELEASE_TAG" SHA256SUMS.txt
           Re-download and verify the final draft asset set
           git ls-remote origin "refs/tags/${RELEASE_TAG}"
           gh api -F draft=false
EOF
  # Preserve literal workflow variables in this mutation fixture.
  # shellcheck disable=SC2016
  sed 's|gh release upload.*|"$RUNNER_TEMP/release-promotion" ensure-checksum --state "$RUNNER_TEMP/selection.json" --file SHA256SUMS.txt|; s|gh api -F draft=false|"$RUNNER_TEMP/release-promotion" publish --state "$RUNNER_TEMP/selection.json"|' \
    "$fixture/.github/workflows/promote-release.yml" > "$fixture/.github/workflows/promote-release-waived.yml"
  cat > "$fixture/.github/release-notes/v0.9.0.md" <<'EOF'
## Highlights

- Complete scheduled work with durable follow-on task chains.
- Translate and validate broader cron expressions in either supported syntax.
- Preserve scheduling intent through daylight-saving transitions.
- Use a consistent visual identity across the desktop application and packages.

Read the [full changelog](https://github.com/shruggietech/go-schedule/blob/v0.9.0/CHANGELOG.md#090---2026-08-30) for every change.
EOF
  cat > "$fixture/specs/048-v100-release-cut/contracts/publication.md" <<'EOF'
Require HEAD, origin/main, and the reviewed S049 merge commit to match.
Run windows-release-gate render-dispositions after formal qualification.
EOF
  cat > "$fixture/test/windows/README.md" <<'EOF'
Run windows-release-gate render-dispositions against the exact candidate.
The packet does not update or close GitHub issues.
EOF
  cat > "$fixture/scripts/verify.sh" <<EOF
#!/bin/sh
if [ "\${1:-}" = list ]; then
  printf '%s\\n' $manifest
  exit 0
fi
if [ "\${1:-}" = gui ]; then
  go run github.com/wailsapp/wails/v2/cmd/wails@v2.17.1 build
fi
exit 2
EOF
 )
