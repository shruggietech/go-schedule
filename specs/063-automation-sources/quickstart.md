# Quickstart: Connected Automation Sources

## Focused backend verification

```bash
go test ./desktop/automation ./desktop
go test -race ./desktop/automation ./desktop
```

## Focused frontend verification

```bash
cd desktop/frontend
npm test -- --run
npm run build
```

## Canonical repository gate

```bash
sh scripts/verify.sh all
```

## Manual smoke test

1. Start the local scheduler daemon and the Wails development application.
2. Open Automation Sources and confirm completion chains, external triggers, Trigger Sets, and filesystem watchers appear in one workspace.
3. Search and filter mixed source types, then navigate every action with the keyboard.
4. Create and edit one chain and one watcher, including a watcher path that produces a daemon health warning.
5. Create one trigger and one Trigger Set. Confirm keys appear only in the explicit secret dialog and are gone after closing it.
6. Reveal, copy, rotate, fire, retarget, toggle, and delete the applicable trigger sources.
7. Change an entity outside an open editor and confirm save offers reload or explicit overwrite.
8. Stop the daemon and confirm the last workspace remains visible, mutation controls are disabled, and retry restores a complete refreshed snapshot.
