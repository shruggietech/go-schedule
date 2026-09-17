# Implementation Plan: Post-release dependency refresh

**Branch**: `codex/090-post-release-dependency-refresh` | **Date**: 2026-09-17 | **Spec**: [spec.md](spec.md)

## Summary

Replace Dependabot pull requests #240, #241, and #242 with one current-main refresh. Upgrade the Go MCP, security, timing, and platform modules, adopt their required Go 1.26 baseline across both modules and operational references, update the React/Vite frontend line, regenerate native integrity files, and prove released behavior through focused and canonical verification.

## Technical Context

**Language/Version**: Go 1.26.0; Node.js 26; TypeScript 5.6.3  
**Primary Dependencies**: go-sdk 1.8.0, x/crypto 0.57.0, x/sys 0.48.0, x/time 0.16.0, React 19.3.0, Vite 8.3.0  
**Testing**: Go race suites, Vitest, TypeScript, Vite production build, Playwright, Wails native build, packaging contracts, canonical eight-gate verification  
**Target Platforms**: Linux, macOS, Windows, Wails desktop, Node 26 browser toolchain

## Constitution Check

The slice preserves pull-request integration, issue-level traceability, native dependency tooling, exact final-head evidence, and all existing safety-critical tests. The Go 1.26 baseline change is required by the selected upstream modules and receives a dated changelog decision. No architecture, public API, release artifact, or verification assertion is weakened.

## Execution Order

1. Establish #243, the source-proposal inventory, selected versions, runtime requirements, and spec-kit analysis.
2. Update the Go baseline contract and operational guidance before resolving both Go graphs.
3. Update the frontend manifest and regenerate its lockfile with Node 26 and npm.
4. Run clean-restoration checks and focused MCP, security, platform, desktop, frontend, browser, and packaging verification.
5. Run canonical verification, audit publication formatting and repository integrity, then publish the official replacement PR.
6. Close #240, #241, and #242 as superseded, process no more than two Codex review rounds, and wait for final-head CI to become green.

## Project Structure

Dependency changes remain in `go.mod`, `go.sum`, `desktop/go.mod`, `desktop/go.sum`, and `desktop/frontend/`. Baseline guidance remains in `CONTRIBUTING.md`, `CLAUDE.md`, `desktop/README.md`, and `CHANGELOG.md`. Automation continues to derive Go versions from module files. S090 planning and evidence live in this directory.

## Risk Controls

- Treat Go 1.26 as an intentional pinned baseline, not an incidental tidy side effect.
- Reconcile both Go graphs independently after the root update.
- Accept only native resolver output and reject force or legacy-peer restoration.
- Keep product behavior changes out of scope unless directly required to restore compatibility.
- Tie review and CI conclusions to the final pushed commit.

