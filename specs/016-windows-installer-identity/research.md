# Research: Windows Installer Identity and Verification

## Decision 1: Use one package-level icon for both installed surfaces

**Decision**: Define one WiX `Icon` sourced from `cmd/gosched-gui/icon.ico`; point both `ARPPRODUCTICON` and the Start Menu shortcut at its identifier.

**Rationale**: Windows Installer defines `ARPPRODUCTICON` as a foreign key into the Icon table, and WiX exposes an Icon identifier on Shortcut. One row gives both consumers an explicit supported relationship and one asset to maintain.

**Alternatives considered**:

- Rely on the shortcut target executable. Rejected because the published v0.8.0 shortcut has an empty Icon column and the reported surface is generic.
- Add separate icon files for shortcut and installed-apps. Rejected because it creates needless drift.
- Wrap the MSI in an executable to change the `.msi` file icon. Rejected as an unrelated distribution expansion; Windows' installed-apps icon is the actual controllable surface.

**Sources**:

- <https://docs.firegiant.com/wix/schema/wxs/icon/>
- <https://docs.firegiant.com/wix/schema/wxs/shortcut/>
- <https://learn.microsoft.com/en-us/windows/win32/msi/arpproducticon>

## Decision 2: Keep the executable icon pipeline, add a regression contract

**Decision**: Do not rewrite `.github/workflows/release.yml`. Assert its current 64-bit `goversioninfo` source/output contract in portable tests.

**Rationale**: The workflow already generates `cmd/gosched-gui/resource_windows_amd64.syso` from the canonical `.ico` before the GUI build. Issue #33 is therefore an observation gap, not a demonstrated missing declaration. Changing a working pinned release step without evidence would replicate speculative logic, contrary to the project directives.

**Alternatives considered**:

- Add a second resource embedding tool. Rejected as duplicate logic and a new dependency.
- Add Win32 icon calls inside the GUI. Rejected until native observation proves the existing Fyne plus executable-resource path is insufficient.

## Decision 3: Test XML relationships structurally

**Decision**: Use Go's standard XML decoder in a table-driven regression test; retain the release-time PowerShell sanity check and add specific icon messages.

**Rationale**: The behavioral artifact is XML. Structural tests are portable, run under the normal race suite, and tolerate harmless formatting changes. PowerShell remains useful in the Windows release step for staged-file checks.

**Alternatives considered**:

- Add more regex-only checks. Rejected for the regression suite because attribute order and formatting are not part of the contract.
- Add an XML library or Pester dependency. Rejected because the standard library is sufficient.

## Decision 4: Separate four evidence classes

**Decision**: Record source definition, candidate artifact, published artifact, and native desktop observation independently as `proven`, `failed`, or `unavailable`.

**Rationale**: The current host can inspect a candidate MSI but cannot provide a clean baseline: v0.8.0 is installed, its machine PATH entry exists, and Windows Sandbox is unavailable without elevation. Candidate evidence cannot satisfy issue #16's published-artifact requirement, and static data cannot replace issue #33's native visual observation.

**Observed baseline (2026-08-27)**:

- Latest release: v0.8.0, published 2026-07-23.
- Published MSI has no Icon table, no `ARPPRODUCTICON`, and an empty Shortcut Icon column.
- Published MSI has one machine PATH row bound to the CLI component.
- Installed v0.8.0 registry entry has an empty DisplayIcon.
- Current machine PATH contains exactly one go-schedule install-directory row.
- Windows Sandbox executable is absent; feature-state inspection requires elevation.

## Decision 5: Make lifecycle tooling refuse contaminated hosts

**Decision**: The state-changing lifecycle script stops before installation if the CLI resolves, the install directory exists, or the machine PATH already contains go-schedule.

**Rationale**: A contaminated development machine was the original reason the PATH defect could be masked. Refusal is safer than trying to clean or reinterpret the operator's real installation, and it avoids destructive changes here.

**Alternatives considered**:

- Uninstall the operator's current application. Rejected as destructive and still not equivalent to a never-installed machine.
- Enable Windows Sandbox or switch Docker to Windows containers. Rejected as a host-wide virtualization change outside this slice.
