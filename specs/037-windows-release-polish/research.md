# Research: Windows Release Polish

## R1 - DPI-aware startup sizing

- **Decision**: Treat work-area values as physical pixels and divide by the positive Fyne canvas scale before applying the 90 percent cap.
- **Rationale**: Fyne window resize units are device-independent; applying the cap before conversion would request twice the intended logical size at 200 percent scaling.
- **Alternatives considered**: Fixed pixel sizes and fixed title-bar subtraction, both of which vary incorrectly with scaling.

## R2 - Launch monitor discovery

- **Decision**: Use `GetCursorPos`, `MonitorFromPoint(MONITOR_DEFAULTTONEAREST)`, and `GetMonitorInfoW.rcWork`.
- **Rationale**: This uses the Windows monitor work area, respects taskbar reservation, and supports a smaller secondary display before the first window is shown.
- **Alternatives considered**: `SystemParametersInfoW(SPI_GETWORKAREA)`, which returns the primary work area; native window handles, which are unavailable at this construction point.

## R3 - Top-level reachability

- **Decision**: Keep existing tab-specific scroll regions and add a headless minimum-size contract for the root content at the 720x540 capped viewport produced by an 800x600 work area.
- **Rationale**: The current top-level tabs already use bounded toolbars plus lists or scroll containers. A root-size regression test protects the primary actions without redesigning dialogs.
- **Alternatives considered**: Wrapping the entire application in another scroll layer, which would introduce nested scrolling and unnecessary layout change.

## R4 - MSI Summary Information

- **Decision**: Author the approved copy through WiX `SummaryInformation Description`, statically assert it, and inspect compiled PID 3 through Windows Installer COM.
- **Rationale**: PID 3 is the Explorer-visible Subject. The compiled artifact, not source text, is the authoritative release evidence.
- **Alternatives considered**: Querying only MSI tables, which do not expose the Summary Information stream, or trusting WiX source compilation without inspection.

## R5 - Evidence honesty

- **Decision**: Record candidate hash and Subject when a compiled MSI is available, and otherwise mark compiled and Explorer checks unavailable.
- **Rationale**: The repository prohibits treating an unavailable prerequisite as green.
- **Alternatives considered**: Inferring native behavior from source, rejected because it would overstate verification.
