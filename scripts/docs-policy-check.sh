#!/bin/sh
# Enforce the current dual-syntax product posture over a deliberately bounded
# inventory. Historical feature artifacts and CHANGELOG.md are excluded.
set -eu

ROOT=${1:-.}
FAILURES=$(mktemp)
trap 'rm -f "$FAILURES"' EXIT

require_text() {
  file=$1
  pattern=$2
  description=$3
  if ! grep -Eiq "$pattern" "$ROOT/$file"; then
    printf '%s: missing policy: %s\n' "$file" "$description" >> "$FAILURES"
  fi
}

CURRENT_FILES="
README.md
docs/README.md
docs/brand.md
docs/cli.md
docs/cron.md
docs/gui-fields.md
internal/cli/cli.go
internal/cli/cron_test.go
specs/001-task-scheduler/spec.md
specs/001-task-scheduler/data-model.md
specs/001-task-scheduler/quickstart.md
specs/001-task-scheduler/contracts/cli.md
specs/001-task-scheduler/contracts/local-api.md
"

for file in $CURRENT_FILES; do
  if [ ! -f "$ROOT/$file" ]; then
    printf '%s: missing policy input\n' "$file" >> "$FAILURES"
    continue
  fi
  if grep -Fiq "cron's power, without its syntax" "$ROOT/$file" ||
     grep -Fiq "desktop editor remains plain-language only" "$ROOT/$file" ||
     grep -Fiq "cron must never be an authoring syntax" "$ROOT/$file" ||
     grep -Fiq "any schedule that standard cron can express" "$ROOT/$file" ||
     grep -Fiq "body = human-readable schedule spec" "$ROOT/$file"; then
    printf '%s: obsolete categorical cron posture\n' "$file" >> "$FAILURES"
  fi
done

require_text README.md 'readable schedules.*supported cron|readable phrase' \
  'human-first dual-syntax positioning'
require_text README.md 'weekdays at 09:00' 'equivalent human example'
require_text README.md '0 9 \* \* 1-5' 'equivalent cron example'
require_text docs/cron.md 'minute hour day-of-month month day-of-week' \
  'five-field order'
require_text docs/cron.md 'Expression versus crontab file' \
  'expression/file distinction'
require_text docs/cron.md 'task timezone' 'task-owned timezone semantics'
require_text docs/cron.md 'lists, ranges, and field-local steps' \
  'standard composite cron breadth'
require_text docs/cron.md 'Both day-of-month and day-of-week restricted' \
  'day-field OR fidelity refusal'
require_text docs/gui-fields.md 'exact retained cron expression' \
  'exact cron edit identity'
require_text specs/001-task-scheduler/spec.md 'cron knowledge is optional' \
  'accessible master-product promise'
require_text specs/001-task-scheduler/contracts/local-api.md 'schedule_syntax' \
  'API syntax discriminator contract'
require_text docs/brand.md 'go-schedule-mark-color\.svg' \
  'canonical full-mark download'
require_text docs/brand.md 'go-schedule-mark-reduced\.svg' \
  'reduced small-size mark download'
require_text docs/brand.md 'go-schedule-horizontal-(color|black)\.svg' \
  'horizontal lockup download'
require_text docs/brand.md 'go-schedule-mark-(white|black)\.svg' \
  'monochrome mark download'
require_text docs/brand.md 'go-schedule-social-preview-1280x640\.png' \
  'social-preview download'
require_text docs/brand.md 'brand-guide\.pdf' 'long-form brand guide download'
require_text docs/brand.md 'brand/manifest\.json' 'complete kit inventory link'
require_text docs/brand.md '#071014.*#62D9B7.*#58A6FF' \
  'core palette values'
require_text docs/brand.md 'Space Grotesk.*Geist.*Geist Mono' \
  'brand typography roles'
require_text docs/brand.md 'A ShruggieTech project' \
  'parent-brand endorsement rule'
require_text docs/brand.md 'Do not' 'logo misuse guidance'

if [ -s "$FAILURES" ]; then
  cat "$FAILURES" >&2
  printf 'docs-policy-check: FAILED with %s issue(s)\n' \
    "$(wc -l < "$FAILURES" | tr -d ' ')" >&2
  exit 1
fi

printf '%s\n' 'docs-policy-check: OK (current product and brand surfaces aligned)'
