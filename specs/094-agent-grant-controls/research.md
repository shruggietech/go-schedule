# Research: Agent Grant Controls

## Decision 1: Treat the actor as the authoritative grant

The actor already owns display name, capability, state, creation, update, and optional expiry. Every local mutation session, durable remote credential, and remote MCP token revalidates the actor on request. The desktop therefore projects actors rather than adding a second grant table.

Alternatives rejected: a parallel grant table would duplicate capability and state, while treating credentials as grants would omit runtime stdio and localhost actors and would make revocation depend on transport-specific records.

## Decision 2: Persist the intended grant deadline on the pairing

Finite access must survive the ten-minute pairing window. A nullable grant deadline travels with the pairing request and becomes the actor expiry during atomic exchange. Non-expiring access uses a null deadline.

Alternatives rejected: calculating the deadline at exchange would silently lengthen access when pairing is delayed, while encoding duration only in desktop state would lose the administrator's approved boundary.

## Decision 3: Keep one-time enrollment material behind the native bridge

The Go desktop service creates the pairing, builds the enrollment bundle, copies it natively, clears its local phrase value, and returns a refreshed secret-free workspace. React never receives the phrase.

Alternatives rejected: rendering a phrase or placing it in React state creates screenshot, accessibility-tree, test-snapshot, and diagnostic exposure paths. Requiring a shell would fail the simple desktop-control outcome.

## Decision 4: Infer transport from authoritative metadata

An MCP actor with a durable client credential is Remote HTTPS. The actor currently attached to the active localhost listener is Localhost HTTP. A remaining runtime MCP actor is Stdio. Observe-only localhost remains represented by its listener status because it has no mutation actor.

Alternatives rejected: adding a transport field to actors would duplicate runtime and credential ownership. Guessing from the client name would be unreliable.

## Decision 5: Make desktop edits monotonic

The edit surface offers only lower capabilities, a new future expiry for non-expiring access, an earlier future expiry for finite access, and revocation. Expansion requires a fresh pairing so high-impact authority is deliberate without adding confirmation prompts to ordinary actions.

Alternatives rejected: unrestricted actor editing makes accidental elevation easy. Confirmation on every agent call would violate the requested low-ceremony workflow and belongs in the separate Manage session policy.

## Decision 6: Use shared audit with a small per-actor window

The workspace loads the newest 25 audit events for a selected grant. It displays only structured actor, daemon, operation, target, result, and timestamp fields already guaranteed by the shared audit model.

Alternatives rejected: a new agent activity log would split evidence and risk inconsistent redaction. Loading full retention would be unnecessary and unbounded.
