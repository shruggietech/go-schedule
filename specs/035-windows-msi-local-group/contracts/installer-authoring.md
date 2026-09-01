# Installer Authoring Contract

- `GoScheduleAdminGroup` names `goschedadmin` and has an empty effective domain;
  source should omit the attribute.
- Any non-empty group domain, including `[ComputerName]`, fails the contract.
- Creation, reuse, update, vital failure, feature linkage, and uninstall
  preservation remain enabled.
- `InstallingUser` remains `[LogonUser]` in `[%USERDOMAIN]`, references the
  group, and preserves membership.
- The compiled `Wix4Group` row records `goschedadmin` with an empty domain.
- Runtime evidence shows provisioning completes before `goschedd` starts.
