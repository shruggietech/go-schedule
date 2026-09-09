# Dependency Baseline Contract

## Included Sources

The replacement pull request accounts for #201, #202, #203, #204, #205, #206, #207, and #208. It closes #215 only when merged.

## Manifest Contract

- `go.mod` and `go.sum` contain the root sqlite and fsnotify selections plus their resolved transitive requirements.
- `desktop/go.mod` and `desktop/go.sum` contain Wails 2.15.0 and a graph tidy against the current local root module.
- `desktop/frontend/package.json` contains the selected frontend dependencies, Node 26 engine boundary, and Vite 8 companion peer.
- `desktop/frontend/package-lock.json` resolves exactly from the manifest without force, legacy-peer, or engine bypass.

## Runtime Contract

- Local development, continuous integration, and release packaging use Node major 26.
- Runtime declarations, Node type definitions, and workflow setup agree on Node major 26.
- The Go language and toolchain lines remain unchanged.

## Compatibility Contract

- Existing public daemon, CLI, desktop, local IPC, notification, and MCP behavior remains unchanged.
- Existing database files migrate and operate through the unchanged schema path.
- Existing watcher behavior remains covered by native and race tests.
- Existing frontend behavior remains covered by unit, production bundle, browser accessibility, responsive, and native package checks.
- No verification assertion is weakened to accommodate an update.

## Evidence Contract

- Both Go graphs pass tidy and module verification without residual changes.
- The frontend passes clean install, audit, unit test, type check, and production bundle.
- Canonical verification passes all eight gates locally.
- Hosted checks and all review dispositions refer to the final reviewed commit.
- Pull requests #201 through #208 remain open until the replacement is merged, then are closed as superseded with a link to the authoritative pull request.
