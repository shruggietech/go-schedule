# Research: Guided Windows Uninstall Entry

## Decision 1: Use the supported maintenance entry

**Decision**: Set `ARPNOREMOVE=1`, leave `ARPNOMODIFY` unset, write an MSI-owned `MsiExec.exe /I[ProductCode]` value as `ModifyPath`, and direct attended Windows Settings users through Modify, then Remove inside a package-owned maintenance wizard.

**Rationale**: Microsoft documents that `ARPNOREMOVE` suppresses direct Remove and `UninstallString`. Hosted CI proved that Windows Installer also omitted `ModifyPath`, and WiX 6.0.2's stock maintenance dialog disables its internal Remove control under the same property. An owned Registry-table value and maintenance page preserve the safe `/I` entry and guided Remove route without changing deletion code.

**Alternatives considered**:

- Force authored dialogs into Windows Settings' reduced native removal interface: rejected because reduced/basic MSI UI does not promise package-authored wizard pages.
- Hide the native MSI entry and author a second registry-owned application entry: rejected because it duplicates Windows Installer product identity and complicates repair, upgrade, and rollback.
- Add a bootstrapper or uninstall executable: rejected because it creates a new executable distribution/signing surface for one registration defect.
- Show a custom-action dialog from the execute sequence: rejected because it can affect silent administration and crosses elevated execution/UI concerns.

## Decision 2: Preserve direct administrative uninstall

**Decision**: Retain direct `msiexec /x` support and the existing `GOSCHEDULE_REMOVE_DATA` contract.

**Rationale**: `ARPNOREMOVE` controls only operating-system UI registration. Administrators and package managers still require deterministic, noninteractive removal. Preserve is the safe default, and wipe remains one exact explicit opt-in.

**Alternatives considered**:

- Require all removals to use full maintenance UI: rejected because it would break managed deployment and CI.
- Default direct removal to wipe because no prompt can appear: rejected because destructive behavior must never be implicit.

## Decision 3: Require installed-state proof

**Decision**: Validate the authored source, compiled MSI tables, and installed registry result, then defer the actual Windows 11 Settings interaction to #94.

**Rationale**: Source-only checks missed the S039 defect and the first S041 implementation's one-property assumption. Compiled tables prove property, registry, dialog, control-condition, and routing rows. A disposable-host install proves the expected `NoRemove`, owned `ModifyPath`, and absence-of-`UninstallString` contract. Only attended Windows 11 observation can prove the final Settings wording and navigation.

**Alternatives considered**:

- Treat WiX XML inspection as sufficient: rejected by the observed S039 mismatch.
- Automate Windows Settings UI in hosted CI: rejected because hosted Windows Server is not the clean attended Windows 11 acceptance environment defined by #94.

## Decision 4: Do not broaden cleanup scope

**Decision**: Leave cleanup helper behavior, profile discovery, reparse-point defenses, and security-state preservation unchanged.

**Rationale**: The reported failure is entry routing. Changing destructive cleanup in the same repair would increase risk without evidence of a cleanup defect.

**Alternatives considered**:

- Rework the entire #98 implementation: rejected as unnecessary and less reviewable.
