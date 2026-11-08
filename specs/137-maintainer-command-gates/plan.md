# Implementation Plan: Maintainer command gates

**Branch**: `codex/137-maintainer-command-gates` | **Date**: 2026-11-08 | **Spec**: [spec.md](spec.md)

## Summary

Deliver #439 by composing existing lifecycle checks, real public-wrapper probes and an isolated driver guard into the canonical automation gate and existing native test matrix. Keep the current automation fixture corpus and eight-gate interface.

## Technical Context

**Language/Version**: Existing Go toolchain and POSIX shell.

**Primary Dependencies**: Existing YAML node parser, standard library and internal/platform; no added dependency.

**Storage**: Owned temporary fixtures and ignored evidence only.

**Testing**: Native policy regressions, existing automation/lifecycle corpus, bounded driver mutation, wrapper root matrix, full foreground canonical aggregate and exact-head hosted checks.

**Target Platform**: Windows, Linux and macOS.

**Project Type**: Internal maintainer tooling and CI integration.

**Performance Goals**: One compiled checker for the policy corpus; a fixed small wrapper matrix rather than per-policy compilation. No timing guarantee or timeout change.

**Constraints**: Eight gates unchanged, no real-suite replay or recursion inside the driver guard, hidden noninteractive Windows children, preserved historical inputs.

**Scale/Scope**: One issue, two public wrappers, five root classes, three native hosts.

## Constitution Check

Pre-design and post-design: PASS. I: reuse existing modules and explicit diagnostics. II: failing-first integration cases plus existing corpus and full canonical checks. III: unchanged public commands and gate interface. IV: bounded shared fixtures with no per-policy compile overhead. V: installed sequential Spec-Kit workflow, review branch and explicit publication authority, stop before owner merge. No waiver or architecture expansion is needed.

## Research and decisions

- The skill-directed read-only research agent identified real-suite recursion in the old failure-propagation mode. Copy the exact driver into an owned fixture and substitute all expensive child roles with inert stubs; never add a production bypass.
- Extract the existing automation fixture writer into a source-only helper, retaining its policy data. Public probes execute the real wrappers from source against target fixtures lacking checker code.
- Compose lifecycle, public adapter and driver guards in one focused shell entry point used by both canonical automation and jobs.test. Keep the policy corpus separately compiled once.
- Enforce native coverage structurally in the named test job, including its three-platform matrix and an unconditional Bash corpus step. Require the canonical automation branch invocation without constructing a general shell interpreter.
- Optional after-specify/after-plan context hooks are omitted because the core plan context update records the current plan. No blocking clarification or Codex review request occurred during planning.

## Project Structure

- test/scripts/lib/: source-only existing fixture writer.
- test/scripts/maintainer-commands_test.sh: focused composition.
- test/scripts/public-entrypoints_test.sh: real wrapper/root corpus.
- test/scripts/verify-driver_test.sh: isolated driver guard and mutation proof.
- scripts/verify.sh and .github/workflows/ci.yml: shared corpus integration.
- internal/maintainercheck/: scoped integration policy and regression cases.
- specs/137-maintainer-command-gates/: specification, design, tasks, analysis and evidence summary.

## Verification and delivery

Capture failing-first integration policy checks, run focused corpora, then all eight canonical gates in the foreground. Freeze implementation inputs and audit all prior retained evidence. Record exact candidate CI and review state, publish the authorized official PR, resolve findings within two total Codex requests and leave auto-merge disabled. Stop for owner review/merge. No release or installed-service operation belongs to this slice.
