# Research: Windows Setup Lifecycle Control

## Decision: Model shortcuts as MSI features

**Rationale**: WiX's FeatureTree maintenance flow exposes native installed feature state. `StartMenuShortcut` is level 1 and selected by default; `DesktopShortcut` is level 2 and absent at the package's default install level. Standard `ADDLOCAL` and `REMOVE` support managed deployment without parallel custom properties. Stable IDs allow `MigrateFeatureStates` to preserve S039-and-later choices.

**Alternatives considered**: Component conditions, transient checkbox properties, and an always-created shortcut removed by a custom action. All obscure maintenance state or make repair/upgrade behavior fragile.

**References**: [WiX Feature schema](https://docs.firegiant.com/wix/schema/wxs/feature/), [WiX MajorUpgrade](https://docs.firegiant.com/wix/schema/wxs/majorupgrade/)

## Decision: Own a FeatureTree-derived UI sequence

**Rationale**: WiX's built-in ExitDialog exposes one optional checkbox. The feature needs two independent completion selections plus custom removal inventory and wipe confirmation. WiX recommends copying and customizing a dialog-set definition when inserting or removing dialogs. One package-owned sequence avoids duplicate success-exit rows and retains native maintenance Modify behavior.

**Alternatives considered**: `WixUI_Minimal` cannot display optional features and disables modification. Layering a second success dialog over stock FeatureTree risks duplicate sequence rows. Mutually exclusive completion options misrepresent two compatible actions.

**References**: [WiX UI dialog library](https://docs.firegiant.com/wix/tools/wixext/wixui/), [Microsoft checkbox guidance](https://learn.microsoft.com/en-us/windows/apps/develop/ui/controls/checkbox)

## Decision: Invoke completion targets only from ordered UI events

**Rationale**: `WixUnelevatedShellExec` hands an application or URL to the interactive user's shell. The utility action uses one shared target property, so each Finish event sets that target immediately before its corresponding action. Distinct actions and repeated fresh-install/full-UI conditions keep selections independent and structurally prevent unattended launch.

**Alternatives considered**: Execute-sequence launch actions can run in silent, elevated, repair, or upgrade contexts. A single combined launcher couples failures and obscures independent selection.

**Reference**: [WiX unelevated shell execution](https://docs.firegiant.com/wix/tools/wixext/util/)

## Decision: Preserve unless one exact secure property opts into wipe

**Rationale**: Destructive removal must never be inferred. `GOSCHEDULE_REMOVE_DATA=1` is the only unattended opt-in; absence and `0` preserve. Full UI requires a separate confirmation before the same property becomes `1`. The property is secure across the client/server installer boundary but contains only a mode, never a path.

**Alternatives considered**: Default wipe is unsafe. Persisting the choice can turn a later generic uninstall destructive. Multiple boolean aliases create ambiguous managed behavior.

**References**: [SecureCustomProperties](https://learn.microsoft.com/en-us/windows/win32/msi/securecustomproperties), [Windows setup guidance](https://learn.microsoft.com/en-us/windows/win32/uxguide/exper-setup)

## Decision: Embed a windowless native cleanup executable

**Rationale**: The helper must exist at commit even after installed files are removed. Embedding one Windows executable in the Binary table removes installed-file ordering and optional-runtime dependencies. It accepts only a fixed operation verb; all roots are internally derived from trusted Windows locations.

**Alternatives considered**: Installed helpers may be removed before commit. Managed custom actions can fail when their runtime is absent during uninstall. PowerShell custom actions inherit script-policy and process-launch complexity. A C/C++ MSI DLL could report directly through the MSI handle but would introduce another language, build system, and much larger unsafe-code review surface. A Go shared library would load the Go runtime inside `msiexec`.

**References**: [Windows Installer custom-action type 2](https://learn.microsoft.com/en-us/windows/win32/msi/custom-action-type-2), [Deferred execution context](https://learn.microsoft.com/en-us/windows/win32/msi/obtaining-context-information-for-deferred-execution-custom-actions)

## Decision: Delete only as a commit custom action

**Rationale**: Commit custom actions run after the installation script completes successfully. The helper preflights every root before deleting any, then records progress and outcome atomically. It uses ignored-return semantics because a checked commit failure can trigger software rollback after user data has already been partly deleted, but that rollback cannot restore the bytes. A retained result ledger and completion guidance surface incomplete cleanup without pretending the software transaction failed safely. Post-commit helper results cannot be read reliably by the already-authored MSI UI, so the UI makes no unconditional cleanup-success claim.

**Alternatives considered**: Immediate or deferred recursive deletion can erase data before uninstall success. Pre-commit quarantine still mutates data during a transaction that may cancel or roll back. `RemoveFolderEx` precomputes RemoveFile rows but does not prove all-profile ownership or safe traversal.

**References**: [Commit custom actions](https://learn.microsoft.com/en-us/windows/win32/msi/commit-custom-actions), [Changing system state with custom actions](https://learn.microsoft.com/en-us/windows/win32/msi/changing-the-system-state-using-a-custom-action), [WiX RemoveFolderEx](https://docs.firegiant.com/wix/schema/util/removefolderex/)

## Decision: Bound roots to ProgramData and registered local profiles

**Rationale**: Machine state lives under `C:\ProgramData\goschedule`. Fyne stores this application's preferences under each profile's `AppData\Roaming\fyne\tech.shruggie.goschedule`. Profile discovery uses the local ProfileList registration and accepts only accessible fixed-local-volume profiles. Root and ancestor reparse points, redirection, device/UNC paths, empty/root paths, and containment mismatches are refused. Descendant reparse entries are unlinked without traversal.

**Alternatives considered**: Scanning `C:\Users` invents ownership and misses relocated profiles. Reading configurable data paths could delete exports or administrator-managed locations. Claiming remote or detached profile cleanup cannot be proven by a local MSI.

**References**: [Windows user profiles](https://learn.microsoft.com/en-us/windows/win32/shell/userenv), [Reparse point operations](https://learn.microsoft.com/en-us/windows/win32/fileio/reparse-point-operations)

## Decision: Preserve `goschedadmin` and memberships

**Rationale**: Existing authoring deliberately preserves the group and membership. No provenance distinguishes installer-created state from pre-existing or shared state. Both uninstall modes therefore identify this security state separately and retain it.

**Alternatives considered**: Removing it during wipe can disrupt unrelated administrators or later installations. Guessing ownership from current membership is insufficient.

## Decision: Separate implementation proof from #94 release proof

**Rationale**: Source, unit, compiled MSI, and disposable silent lifecycle verification can run automatically and block regressions. The current workstation is a non-elevated retained-state environment and project policy forbids automation that flashes attended setup, GUI, or browser windows. The clean interactive candidate matrix therefore remains #94's downstream release gate, and #97/#98 remain open until that evidence satisfies their attended acceptance items.

**Alternatives considered**: Treating XML assertions or a session-zero CI runner as attended desktop proof would be false evidence. Running destructive UI tests on the maintainer workstation would be unsafe.
