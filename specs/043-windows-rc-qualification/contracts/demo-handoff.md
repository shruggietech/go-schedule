# Contract: S043 Demo Installer Handoff

The installer handoff is complete only when all fields below are supplied.

## Identity

- Slice: `S043`
- Artifact class: `local demo, not an official release candidate`
- Absolute local MSI link
- Source commit (full SHA)
- Embedded binary version
- MSI ProductVersion and ProductCode
- Byte size and SHA-256
- Build and inspection timestamp

## Automated disposition

- Eight canonical gates: each named with result
- Focused Go/race suites: exact commands and result
- WiX source validation: result
- Compiled MSI inspection: result and linked report
- Destructive installed lifecycle: either a result from an approved disposable environment or explicit `unavailable on this non-disposable host`

## Operator request

The handoff asks the operator to test only the exact linked hash and report:

1. setup shortcut choices and Finish actions;
2. startup window size and absence of error spam;
3. Options, Info text, navigation spacing, Exit, storage rows, and task editing;
4. one manual and one scheduled task;
5. Windows Settings maintenance entry, cancellation, preserve uninstall, and wipe uninstall as distinct journeys.

## Claim boundary

The handoff must say that success authorizes preparing the S043 pull request but does not itself close #94/#98/#96. Those issues require the later formal staged candidate and full attended evidence matrix.
