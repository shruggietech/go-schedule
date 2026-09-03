# Implementation Plan: Windows Setup Lifecycle Control

**Branch**: `codex/039-windows-setup-lifecycle-control` | **Date**: 2026-09-02 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/039-windows-setup-lifecycle-control/spec.md`

## Summary

Replace the stock minimal Windows installer UI with a package-owned feature-selection and maintenance flow. Model Start Menu and desktop shortcuts as stable optional MSI features, add two independent success-page actions that execute only from a fresh attended install, and add a preserve-by-default uninstall inventory with an explicitly confirmed wipe mode. Implement destructive cleanup through an embedded, windowless helper that derives and validates only trusted product roots, deletes them only as a commit custom action after successful software removal, and retains protected evidence when cleanup is incomplete. Extend source, helper, compiled-MSI, silent lifecycle, documentation, and CI verification while reserving clean interactive desktop proof for #94.

## Technical Context

**Language/Version**: Go 1.25, WiX Toolset 6.0.2 XML, PowerShell 7, Markdown

**Primary Dependencies**: Go standard library, existing `golang.org/x/sys/windows` dependency, WiX UI and Util extensions 6.0.2, Windows Installer

**Storage**: Existing `C:\ProgramData\goschedule` runtime root; Fyne preference leaf `AppData\Roaming\fyne\tech.shruggie.goschedule`; protected uninstall-result evidence retained only for refused or incomplete cleanup

**Testing**: Go unit and integration tests, PowerShell source verifier and parser checks, compiled-MSI database inspection, disposable hosted Windows silent lifecycle matrix, canonical eight-gate repository suite

**Target Platform**: Per-machine x64 Windows MSI on Windows 11; cross-platform source and regression gates remain mandatory

**Project Type**: Go desktop application, service, CLI, Windows installer, and installer-only cleanup helper

**Performance Goals**: Installer UI remains immediate; cleanup performs one bounded preflight and traversal of declared roots and adds no runtime or scheduler overhead

**Constraints**: Preserve by default; validate all targets before deleting any; no arbitrary input paths, reparse traversal, permanent elevation, PowerShell custom action, visible helper console, shared-security deletion, completion action outside fresh full UI, or false claim of attended release proof

**Scale/Scope**: Two shortcut features, two completion actions, one machine data root, every safely registered accessible local profile preference leaf, and preserve/wipe lifecycle modes

## Constitution Check

| Principle | Gate | Design evidence |
| --- | --- | --- |
| I. Code Quality | PASS | Shortcut/UI authoring is structurally validated; cleanup discovery, validation, preflight, deletion, and reporting are separated behind testable operations with contextual errors. |
| II. Testing Standards | PASS | Contract and cleanup tests are written first; helper and installer behavior receive unit, source, compiled-package, silent native, and canonical race coverage. |
| III. UX Consistency | PASS | Defaults are explicit, completion choices are independent, destructive cleanup requires confirmation, and every lifecycle mode uses matching language and documentation. |
| IV. Performance | PASS | No scheduler path changes; cleanup is bounded to declared roots and runs only on explicit complete uninstall. |
| V. Autonomous Execution | PASS WITH RECORDED OPERATOR OVERRIDE | Full Spec Kit and analyze gates remain mandatory. The operator explicitly pre-authorized this branch push and PR, so the normal pre-publication pause is satisfied in advance for S039 only; merge, tag, and release remain outside authorization. |

### Post-design re-check

All engineering gates remain PASS. The operator's explicit automatic-push direction is recorded as a narrow process override rather than a relaxation of implementation or verification requirements. Changes to `build/windows/**`, `.github/workflows/**`, and `docs/INSTALL-windows.md` are required by #97/#98 and receive a dated changelog decision. Destructive cleanup never accepts caller-supplied roots and does not delete `goschedadmin` or membership state.

## Architecture and Decision Log

### Use stable child MSI features for shortcut ownership

Keep core application components inside a non-optional `Main` feature. Move the existing Start Menu component into a stable `StartMenuShortcut` child feature at the default install level, and add `DesktopShortcut` as a stable child feature above the default install level. A FeatureTree-style selection surface then supports fresh install and maintenance modification using Windows Installer's native feature state. Stable feature IDs allow later major upgrades to migrate choices. The first upgrade from pre-S039 releases necessarily maps the old mandatory Start Menu state to the new defaults because those releases contain no matching feature IDs.

Property-conditioned components were rejected because component conditions are sticky across repair and maintenance and do not expose native feature state. Bespoke shortcut booleans were rejected because standard `ADDLOCAL` and `REMOVE` already provide deterministic managed deployment semantics.

### Own one coherent installer UI sequence

Author one package-owned flow based on WiX FeatureTree rather than combining two stock dialog sets. It includes welcome/license, shortcut feature selection, install verification/progress, maintenance selection, a removal inventory, a separate wipe confirmation, and a two-checkbox success dialog. This avoids duplicate success-exit sequence rows and keeps removal decisions inside the same maintenance flow. Direct full-UI removal is intercepted before progress; reduced, silent, external-UI, or policy-driven removals preserve by default unless the exact managed wipe property is supplied.

The completion checkboxes are independent. Finish-button events set the shared WiX unelevated shell target immediately before each selected action, invoke two distinct actions, and then close the dialog. Every event repeats the fresh-install/full-UI guard. The actions never appear in the execute sequence, so silent, repair, upgrade, removal, rollback, cancellation, and failure cannot launch a visible process.

### Use one exact secure wipe opt-in

`GOSCHEDULE_REMOVE_DATA` is a secure public property with default `0`; only exact values `0` and `1` are valid. A full-UI removal begins in preserve mode. Selecting wipe routes to a separate confirmation, and only confirmation changes the property to `1`. Back or cancel returns it to `0`. Silent removal preserves unless `GOSCHEDULE_REMOVE_DATA=1` is supplied explicitly. The value is not persisted.

Every destructive action uses the same complete-removal condition: the product is installed, `REMOVE` includes all features, the operation is not an upgrade or reinstall, and the opt-in equals `1`. Invalid explicit values fail before the execute sequence.

### Delete only from a post-transaction commit action

Build a narrow `gosched-cleanup.exe` only for embedding in the MSI Binary table. It has no public product command and is built with the Windows GUI subsystem so execution cannot create a console. The helper derives ProgramData and registered profile roots from Windows, appends fixed product-owned leaf names, and accepts no path from MSI properties, configuration, environment overrides, or command-line input.

The helper first discovers and validates every existing target. Any unsafe target refuses the whole wipe before deletion starts. After preflight succeeds, it removes each declared root without following reparse entries and records progress atomically in a protected installer-result ledger. The MSI invokes the helper as a non-impersonated commit custom action, so it runs only after the installation script has completed successfully. `Return=ignore` prevents an incomplete irreversible cleanup from triggering a fictitious software rollback that cannot restore deleted bytes. A refused or residual cleanup leaves a bounded report and registry summary; success removes stale result evidence.

Pre-commit quarantine was rejected because even an atomic rename violates the requirement that canceled, failed, or rolled-back uninstall cannot mutate application data. Direct commit deletion with checked return was rejected because Windows Installer can start rollback after partial irreversible deletion. `RemoveFolderEx` was rejected because it cannot enumerate registered profiles, establish canonical trusted bounds, or report all partial cleanup conditions. A managed, PowerShell, C/C++, or Go shared-library custom action was rejected because it adds an optional runtime or a second unsafe build/test surface without solving commit-result propagation reliably.

### Preserve security state and bound “all user data” honestly

Both removal modes preserve the local `goschedadmin` group and memberships because the installer records no reliable pre-existing ownership provenance and other software or administrators may depend on them. Wipe covers the machine root and the application-specific Fyne leaf of every registered, accessible local profile on a fixed local volume. Disconnected roaming copies, detached profile containers, unregistered/orphaned profiles, redirected roots, and encrypted or inaccessible state are reported as outside or incomplete rather than claimed erased.

### Layer proof without absorbing #94

Source and Go tests prove authoring semantics and cleanup safety. A Windows CI job builds the exact branch MSI, inspects its database, and exercises silent preserve/wipe/install/uninstall behavior on the disposable runner without launching GUI or browser. S039 records interactive attended evidence as unavailable on the current standard-user workstation. Issues #97 and #98 remain open with `Refs` until #94 supplies their required attended Windows 11 evidence. Issue #94 remains the clean candidate gate for visible defaults, user interaction, unelevated launches, window sizing, recurring-error observation, and complete release readiness.

## Project Structure

```text
build/windows/goschedule.wxs
build/windows/verify_wxs.ps1
cmd/gosched-cleanup/
internal/winuninstall/
test/integration/windows_installer_contract_test.go
test/windows/inspect-installer.ps1
test/windows/Invoke-InstallerContractCI.ps1
test/windows/README.md
.github/workflows/ci.yml
.github/workflows/release.yml
docs/INSTALL-windows.md
specs/039-windows-setup-lifecycle-control/
CHANGELOG.md
CLAUDE.md
specs/README.md
```

**Structure Decision**: Extend the existing installer, verification, and lifecycle boundaries. The only new production code is an installer-private cleanup command plus a focused internal package; no application runtime, database schema, public API, or external dependency changes.

## Complexity Tracking

| Added complexity | Why needed | Simpler alternative rejected because |
| --- | --- | --- |
| Package-owned installer dialogs | #97 requires two independent completion choices while #98 requires removal inventory and confirmation in the same maintenance flow. | Stock WiX dialog sets provide at most one finish checkbox and cannot safely compose two separate exit flows. |
| Commit-only cleanup helper | Multi-profile cleanup must be path-bounded, independent of installed files, and delayed until successful software removal. | Broad recursive directives cannot establish ownership, enumerate profiles, avoid reparse traversal, or retain accurate partial-cleanup evidence. |
| Windows built-MSI/silent lifecycle CI job | Destructive sequencing and compiled database relationships require proof beyond XML inspection. | Source-only tests can pass while the linker emits incorrect sequence, feature, or dialog tables. |
