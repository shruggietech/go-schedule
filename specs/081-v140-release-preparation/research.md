# Research: Cumulative v1.4.0 Release Preparation

## Public version boundary

**Decision**: Publish v1.4.0 as the next cumulative release after v1.1.1.

**Rationale**: Current main already contains the completed v1.2.0 desktop, v1.3.0 notification and local MCP, and v1.4.0 remote-access increments. A tag created now cannot truthfully isolate either earlier milestone.

**Alternatives considered**: Retroactively tag v1.2.0 and v1.3.0 from current main, which would falsely include later work; reconstruct historical commits, which would create unsupported artifacts without the current release corrections; skip directly to v1.5.0, which would overstate the shipped multi-daemon roadmap.

## Milestone history

**Decision**: Preserve v1.2.0 and v1.3.0 as closed delivery milestones and reopen v1.4.0 only for publication issue #226.

**Rationale**: Milestones record completed planning increments, while tags and GitHub Releases record public immutable distributions. Keeping those roles distinct preserves the actual project history.

**Alternatives considered**: Reopen all three milestones, which would blur completed implementation with the remaining publication operation; erase or rename milestones, which would rewrite planning history.

## Durable documentation wording

**Decision**: Prepare documentation that identifies v1.4.0 content but directs readers to the dynamic latest-release link for current publication status.

**Rationale**: The reviewed commit becomes the tagged artifact unchanged. Text that says v1.1.1 is latest would be true before tagging but stale inside the eventual v1.4.0 package. Publication-neutral wording remains accurate both before and after promotion.

**Alternatives considered**: A post-release documentation commit, which would leave tagged files stale; claiming v1.4.0 is already public, which would be false during review and staging.

## Release metadata preflight

**Decision**: Extend the existing README release-metadata gate to require one matching dated changelog section and a tag-specific release-note file before release state or artifacts can be mutated.

**Rationale**: The workflow already blocks a mismatched README health example. Changelog and note identity are equally source-owned prerequisites, and failing them before draft mutation makes recovery simpler.

**Alternatives considered**: Rely on the release upload action to fail for a missing note, which occurs later and does not validate the changelog; add a new workflow, which would duplicate tag and CI authority.

## Release-copy shape

**Decision**: Use exactly four concise one-line highlights and one final tagged changelog link.

**Rationale**: This is the established repository contract and can summarize the Wails desktop, automation sources, notification and local-agent access, and remote daemon access without copying the complete changelog.

**Alternatives considered**: Exhaustive notes, which duplicate a very large changelog; fewer than four bullets, which cannot represent the four principal outcome groups and fails the current validator.

## Windows qualification

**Decision**: Require separate clean-install and public-v1.1.1-upgrade snapshots using the exact staged v1.4.0 MSI, with the established 47 attended observations bound to the candidate.

**Rationale**: v1.4.0 replaces the desktop implementation and introduces new persisted state and opt-in listeners. Source and silent package checks cannot prove native interaction, retained state, or default-off behavior across the real upgrade boundary.

**Alternatives considered**: Fresh install only, which misses migration and service replacement; reuse older attended evidence, which is bound to different bytes and a Fyne-era candidate; infer observations from hosted runners, which lack a credible attended Windows desktop.

## Native window evidence evolution

**Decision**: Capture version-two generic desktop measurements directly from the installed Wails window and retain read compatibility for version-one Fyne attachments.

**Rationale**: The current collector requires a Fyne metrics file that the production Wails application cannot emit. Win32 already supplies the client rectangle and effective DPI needed for logical content dimensions and scale, so a toolkit-neutral attachment removes the impossible input without weakening native evidence. Retaining the version-one decoder preserves historical auditability.

**Alternatives considered**: Store Wails values under Fyne field names, which would corrupt meaning; add release-only instrumentation to the product, which is unnecessary because the operating system already exposes the measurements; reject all historical evidence, which would break prior audit records.

## Authorization boundary

**Decision**: S081 ends with a green reviewed pull request. Tagging, draft staging, attended installation, promotion, and final administrative closure remain sequential post-merge operations under issue #226.

**Rationale**: The user authorized branch publication and the pull request, but repository governance reserves tags and releases for explicit authorization. Physical evidence also cannot exist before the exact reviewed candidate is staged.

**Alternatives considered**: Tag from the PR head, which bypasses the reviewed main boundary; manufacture a local candidate, which cannot satisfy immutable release provenance.
