# Implementation Plan: Corrected candidate qualification

**Branch**: `codex/086-release-candidate-refresh` | **Date**: 2026-09-15 | **Spec**: [spec.md](spec.md)

## Summary

Use a standard-library Go preparation command and guest-only Windows bootstrap. Bind candidate, public v1.1.1 baseline, portable PowerShell, offline WebView2, and helpers to explicit provenance and hashes. Generate separate Sandbox configurations and export directories. Stage software locally before intentional attended installation, expose flushed logs and diagnostic deadlines, and never manufacture observations.

## Technical Context

**Language/Version**: Existing Go toolchain; Windows PowerShell 5.1 bootstrap and supplied PowerShell 7 collector runtime.

**Primary Dependencies**: Standard library, Windows Sandbox/Installer, supplied offline prerequisites.

**Storage**: Create-only JSON, guest-local staging, dedicated host exports.

**Testing**: Go preparation unit tests, guest nondestructive process fixtures, integration source contracts.

**Target Platform**: Disposable Windows 11; portable preparation tests.

**Project Type**: Release tooling, no product behavior redesign.

**Performance Goals**: Progress within 30 seconds, explicit deadlines, no latency-fix claim without native evidence.

**Constraints**: No host install, hidden consoles, no safety-guard bypass or release mutation.

**Scale/Scope**: #226 and #228 with #229-#233 native regressions.

## Constitution Check

Spec Kit artifacts precede implementation. Full foreground `sh scripts/verify.sh all` precedes commit. Missing prerequisites are unrun, not green. Local completion is distinct from attended qualification. Halt once before branch push/PR. Tag replacement, staging, and promotion require separate authority. No issue closure or new dependency. Pre-research and post-design checks pass without exceptions.

## Project Structure

```text
specs/086-release-candidate-refresh/
  spec.md, plan.md, research.md, data-model.md, quickstart.md
  contracts/session.md
  checklists/requirements.md, checklists/release-safety.md
  tasks.md
scripts/windows-qualification-session/
  main.go, main_test.go
test/windows/
  Start-QualificationGuest.ps1
  README.md
```

**Structure Decision**: Small generation command plus guest-only runner. Reuse the attended helper and gate, not another evidence framework.
