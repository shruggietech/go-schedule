# Implementation Plan: Windows version identity

**Branch**: codex/141-windows-version-identity | **Date**: 2026-11-09 | **Spec**: [spec.md](spec.md)

## Summary

Deliver #458 through a maintained desktop builder and existing local, CI and release packaging paths.

## Technical Context

Go 1.26.0, Wails 2.17.1, PowerShell 7 and WiX 6.0.2. No new dependency or application storage. Windows executable and MSI inspection is read-only. Tests use injected command execution and fixture configuration, followed by actual build metadata readback. No scheduling hot path changes or performance expansion.

## Constitution Check

Before and after design: PASS. Regression-first behavior checks, all eight foreground verification gates, hidden noninteractive child processes, unchanged access/state semantics, normal review branch, no release mutation. Research delegation follows the installed planning command. Optional agent-context hooks are skipped; the required plan reference is maintained in CLAUDE.md.

## Project Structure

- scripts/desktop-build/: selected version parsing, temporary configuration, hidden Wails invocation and tests.
- build/windows/verify_version.ps1: compiled PE/MSI identity refusal.
- Makefile, scripts/verify.sh, .github/workflows/ci.yml and release.yml: maintained invocation points.
- specs/141-windows-version-identity/: specification, research, model, validation and delivery records.

## Complexity Tracking

No constitutional violations. Configuration is restored exactly, with exclusive build ownership to reject simultaneous writers.
