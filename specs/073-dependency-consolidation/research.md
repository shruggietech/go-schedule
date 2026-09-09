# Research: Dependency Consolidation

## Apply Updates to Current Main

**Decision**: Recreate the version changes on the S072-merged `main` tree instead of cherry-picking Dependabot commits.

**Rationale**: The eight automated branches were created before later v1.3 work merged. Comparing their heads directly with current `main` shows unrelated deletions and reversions. Native package-manager updates retain the current product tree and isolate the intended dependency changes.

**Alternatives considered**: Sequentially merge each Dependabot branch, which creates overlapping lockfile conflicts and risks reverting newer work; cherry-pick each bot commit, which preserves intended manifest patches but still serializes stale lockfiles and obscures the final combined graph.

## Version Selection

**Decision**: Use Wails 2.15.0, modernc.org/sqlite 1.58.0, fsnotify 1.10.1, React and React DOM 19.2.8, React types 19.2.18, React DOM types 19.2.7, Testing Library jest-dom 7.0.1, jsdom 30.0.1, Vite React plugin 6.1.1, and Node types 26.5.0.

**Rationale**: These are the exact proposed versions except Node types, where 26.5.0 is a newer compatible patch than #208's 26.4.1. Current module and registry metadata confirms no newer stable Wails, sqlite, fsnotify, React, jest-dom, jsdom, or React plugin version is available through the repository's package sources at planning time.

**Alternatives considered**: Freeze Node types at 26.4.1, which would create immediate patch staleness without compatibility benefit; broaden to unrelated outdated packages such as TypeScript 7, which is not represented by #201 through #208 and is excluded.

## Node Runtime Alignment

**Decision**: Move the frontend package and hosted build baseline from Node 24 to Node 26.

**Rationale**: Requested Node 26 type definitions should match the supported runtime major. More importantly, jsdom 30.0.1 declares support for current Node 26 while its Node 24 range begins above the workstation's and potentially pinned runner's available 24.x baseline. Node 26 satisfies jest-dom 7, jsdom 30, Vite 8, and Vitest 5 without ignoring engines. Aligning runtime and types is the smallest coherent baseline.

**Alternatives considered**: Keep Node 24 with Node 26 types, which creates a false compile-time API surface; install jsdom with engine warnings, which violates clean restoration; omit jsdom 30, which fails issue #215 without an upstream defect that prevents a supported alignment.

## Vite Peer Compatibility

**Decision**: Update Vite from 7.3.6 to 8.2.2 alongside Vite React plugin 6.1.1.

**Rationale**: Plugin 6.1.1 declares Vite `^8.0.0` as its peer range. The failed #207 jobs report `ERESOLVE` against Vite 7.3.6. Vite 8.2.2 is the current compatible stable peer and supports Node 26.

**Alternatives considered**: Use `--force` or `--legacy-peer-deps`, both of which hide an invalid graph; retain plugin 5, which omits #207; choose an older plugin 6 release, which still requires Vite 8 and provides no benefit.

## Go Module Reconciliation

**Decision**: Update the root direct dependencies first, run root tidy and verification, then update Wails in the desktop module and run desktop tidy and verification against the local root replacement.

**Rationale**: The failed #202 and #203 desktop production jobs report that `go.mod` needs updates because each bot updated only the root integrity state while the desktop module imports the replaced root module. Reconciling both module graphs together is required for a clean native build.

**Alternatives considered**: Commit only the bot-generated root sums, which reproduces the failures; use read-only module flags to suppress updates, which leaves the graph inconsistent.

## Security and Provenance

**Decision**: Preserve package-manager integrity files, run clean installation and audit evidence, retain CodeQL and security-bot review, and make no dependency install-script bypass part of the committed workflow.

**Rationale**: The package sources and lockfiles provide provenance, while final-head hosted checks and security review cover combined-update risk. No new production network service, credential, or permission surface is introduced.

**Alternatives considered**: Rely only on individually green bot checks, which do not test the combined graph; add a new policy or scanning service, which exceeds the issue.
