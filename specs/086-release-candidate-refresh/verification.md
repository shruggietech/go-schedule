# S086 verification

**Date**: 2026-09-15

**Review branch**: `codex/086-release-candidate-refresh`

## Local results

`sh scripts/verify.sh all` completed in the foreground with exit code 0 using the installed `C:\Program Files\Git\bin\sh.exe` through a no-window launcher. Format, vet, lint (zero issues), root race, GUI, coverage, docs, and automation passed. GUI includes Wails module race tests, native Windows production build, 24 frontend files/128 tests, and the production frontend bundle. Six core coverage results: engine 82.9%, schedule 89.1%, timezone 91.3%, store 80.1%, catchup 88.9%, logbus 91.1%. Automation policy and rejection fixtures passed.

PowerShell compliance passed: UTF-8 without BOM, LF, no trailing whitespace, one trailing newline, no pictographs, four ordered dividers, and help before CmdletBinding. Focused preparation tests passed; Windows regression tests execute hidden nondestructive children and verify success (0), failure (7), timeout (null exit), reboot-required (3010), unavailable attended status, and host invocation refusal. No fixture installs software.

Preparation tests were first run before implementation and failed to compile because the required command types/functions did not exist. Process-control fixtures were added alongside the runner and subsequently asserted by Windows tests; they were not independently demonstrated red before runner implementation. No claim of native MSI success follows from these fixtures.

## Findings addressed

Source review found and corrected an overbroad msiexec-process guard, incorrect treatment of MSI 3010, and inaccessible collector handoff. Idle service processes no longer block upgrades; native conflicting transactions remain nonpassing. Reboot requirements stop subsequent work and are distinct from installation failures. Collector and portable runtime are packaged, unavailable templates initialize automatically, and the final operator acknowledgment exports the workspace. Finalization remains through the original host repository collector and shared gate. Canonical MSI filenames are preserved.

An initial canonical run passed seven gates but failed automation because the specification used an unsupported Approved state and lacked an inventory row. Both were corrected before the successful complete rerun. The agent-context extension's generated width-wrapped paragraph was also corrected before format passed. Git's shell was installed but absent from PATH; the verified absolute executable was used without installing or changing host tools.

## PR review follow-up

Initial hosted macOS race testing exposed that the test temporary-directory prefix uses a system symlink. Tests now resolve their temporary root before constructing inputs and outputs; production rejection of linked inputs and destinations remains unchanged. Codex's initial review identified an overbroad reverse-containment check that rejected safe sibling output directories. A new sibling-output test reproduced the rejection before the check was changed to compare against the input file itself. Fixture helpers are now present for negative preparation tests, avoiding unrelated missing-helper failures. The corrected full foreground `sh scripts/verify.sh all` run passed all eight gates with exit code 0 on 2026-09-15. Follow-up hosted verification and final external review remain required before merge readiness.

## Unverified delivery boundaries

No native candidate MSI or Sandbox bootstrap was executed in this implementation round. The fifteen-minute delay remains an unverified reported symptom, not a reproduced or resolved cause. Generated packages and diagnostic logs do not satisfy attended observations. Suitable clean Windows 11, normal-user, multiple-profile, standard/high/mixed-DPI, and complete installer/removal evidence remain required under the existing gate.

No branch push, PR creation, tag replacement, draft asset mutation, release promotion, issue closure, or milestone closure is included in this local verification. #226 and #228 remain open. After push/PR authority, hosted CI and external reviews are separate evidence. After merge, obtain explicit tag/draft refresh authority, stage the exact corrected commit, qualify exact hosted bytes, and obtain promotion authority. Do not use closing keywords for #226/#228 in the tooling PR.
