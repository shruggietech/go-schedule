# Research: Production Wails Shell and Local Connection

## Decision 1: Graduate into a separate production module

**Decision**: Create `desktop/` as the production-intent Wails module and retain `experiments/wails-foundation/` unchanged as S060 evidence.

**Rationale**: A separate module isolates Wails and Node dependency ownership while the existing Fyne application still ships. It prevents accidental installer or root-module changes and gives later migration slices a stable destination.

**Alternatives considered**: Rename the experiment, which destroys its evidence boundary; add Wails to the root module, which expands current binary dependency operations; replace `cmd/gosched-gui` now, which violates #157 sequencing.

## Decision 2: Use one actor-style connection owner

**Decision**: One manager loop owns generation changes, attempt contexts, retry delays, event-stream cancellation, and shutdown.

**Rationale**: Ownership is auditable, stale results have one comparison point, and manual retry cannot overlap an automatic attempt or event stream.

**Alternatives considered**: Independent request and event goroutines coordinated by shared flags, which multiplies races; frontend-owned retries, which expose transport policy and keep work alive beyond Wails shutdown.

## Decision 3: Inject retry waiting, not time itself

**Decision**: The manager accepts a small scheduler interface that returns a wait channel for a requested duration; production uses timers and tests use controlled channels.

**Rationale**: The contract proves the exact 250 millisecond, one second, and five second progression without real sleeps and lets shutdown and manual retry interruption be asserted deterministically.

**Alternatives considered**: A clock-only abstraction cannot advance a blocking wait; direct timers create slow flaky tests; a broad executor abstraction adds unnecessary machinery.

## Decision 4: Sanitize at the Go-to-frontend boundary

**Decision**: Map client connection categories, compatibility outcomes, and event kinds into closed enums and approved guidance before emission.

**Rationale**: Local endpoint paths, named-pipe identifiers, OS errors, and future credentials never enter the web runtime, logs, screenshots, or browser test fixtures.

**Alternatives considered**: Passing wrapped errors and formatting in React leaks sensitive details; parsing message strings in the frontend is unstable and not transport-neutral.

## Decision 5: Negotiate compatibility without a daemon API change

**Decision**: Use the existing health version and a desktop-owned capability manifest for the current v1 daemon contract. Accept development builds and major version 1, reject other explicit majors, and expose only capabilities implemented by the connection adapter.

**Rationale**: #152 needs stable feature-facing capability data, but changing the daemon protocol would expand S061 into public API versioning. The adapter can truthfully describe the current method surface and later replace this derivation when #18 introduces capability negotiation.

**Alternatives considered**: Add a daemon capability endpoint now, which broadens scope; assume every version is compatible, which makes incompatibility untestable; exact-version equality, which would reject compatible patch updates.

## Decision 6: Derive components from contracts, not prototype markup

**Decision**: Recreate tokens and shared primitives under documented component contracts, then assemble the shell and placeholder routes from those primitives.

**Rationale**: S060 explicitly marks fixture routes and page markup disposable. Contract-driven components prevent page-specific styles and give #153 through #156 stable composition points.

**Alternatives considered**: Copy `App.tsx` wholesale, which carries fixture tasks and hard-coded screens; adopt a third-party component library, which adds weight and visual constraints without need.

## Decision 7: Keep appearance session-local in S061

**Decision**: Support system, light, and dark appearance immediately but do not persist the selection yet.

**Rationale**: #156 owns durable Fyne preference migration and needs to decide preservation versus retirement. Local persistence now would create an unapproved second preference source.

**Alternatives considered**: Browser local storage, which creates a migration contract prematurely; no appearance switch, which fails #151.

## Decision 8: Extend existing CI with a production foundation matrix

**Decision**: Add dedicated production desktop build and browser-contract jobs while retaining the S060 proof jobs until cutover.

**Rationale**: The proof remains historical evidence, while `desktop/` needs independent fail-closed evidence on all supported platforms. Parallel jobs keep failures attributable.

**Alternatives considered**: Replace proof jobs, which erases comparison evidence; build only Windows locally, which cannot satisfy #151 and #152 platform gates.
