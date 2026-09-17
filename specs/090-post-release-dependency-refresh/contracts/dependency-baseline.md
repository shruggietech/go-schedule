# Dependency Baseline Contract

## Included Sources

The replacement pull request accounts for #240, #241, and #242 and closes #243 only when merged.

## Go Contract

- Root and desktop modules declare Go 1.26.0.
- The root selects go-sdk 1.8.0, x/crypto 0.57.0, x/sys 0.48.0, and x/time 0.16.0.
- Both graphs accept native transitive resolution and pass tidy plus verification without residual changes.
- CI and release workflows continue deriving their Go version from the appropriate module file.

## Frontend Contract

- React and React DOM are 19.3.0 with matching 19.3.0 type packages.
- Node types are 26.5.1 and Vite is 8.3.0.
- Node 26 remains the declared and automated runtime.
- `npm ci` restores the lockfile without force, legacy-peer, or engine bypass.

## Compatibility Contract

- Public daemon, CLI, desktop, IPC, MCP, remote access, notification, scheduling, and packaging behavior remains unchanged.
- Existing security and accessibility assertions remain enabled.
- A dependency-caused compatibility repair must be narrow, tested, and recorded.

## Evidence Contract

- Focused affected-surface checks and all eight canonical gates pass.
- Hosted checks and review dispositions refer to the final PR head.
- #240, #241, and #242 receive explicit supersession comments linked to the replacement PR.

