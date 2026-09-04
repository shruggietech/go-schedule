# Contract: v1.0.0 Staging, Qualification, and Promotion

## Preparation contract

The S048 pull request is eligible for merge only when:

- `CHANGELOG.md` has an empty Unreleased section and a dated v1.0.0 section
  containing all 33 pre-cut entries;
- the README has exactly one badge derived from GitHub's latest published
  release and exactly one v1.0.0 health example;
- `.github/release-notes/v1.0.0.md` satisfies the established highlights-only
  policy and links to the tagged changelog;
- the Release/Promote workflow and S047 evidence contracts remain intact;
- all local, hosted, and review gates pass; and
- the PR uses `Refs #122`, leaving release and verification issues open.

## Tag authorization contract

After merge, do not create the tag until the maintainer explicitly authorizes
v1.0.0 tag staging. Immediately before mutation:

1. fetch and prune origin;
2. switch to `main` and fast-forward only;
3. require an empty working tree;
4. require `HEAD`, `refs/remotes/origin/main`, and the reviewed S049 merge commit
   to be identical;
5. require both local and remote `v1.0.0` tag absence; and
6. require no GitHub release for v1.0.0.

Create one annotated `v1.0.0` tag at that commit and push only that tag. Never
move or recreate the tag to absorb later changes.

S049 deliberately supersedes the S048 merge commit as the tag boundary because
the required release-operations pull request advances `main`. S049 is limited
to release tooling, tests, specifications, and operator documentation; it does
not alter packaged runtime behavior.

## Draft staging contract

Wait for the tag-triggered Release workflow. Accept exactly one successful run
whose event is `push`, head branch is `v1.0.0`, and head SHA is the tagged commit.
The release must remain draft. Download and inventory the seven build packages
plus `windows-candidate-manifest.json`. The manifest must identify that
run and its successful `Build & stage GUI (windows-latest)` attempt.

Verify the MSI and manifest with `windows-release-gate verify-candidate` before
installation. Any missing/extra asset, failed job, identity mismatch, public
release, or tag drift stops the ritual.

## Formal qualification contract

Initialize the attended workspace with the exact downloaded MSI, tagged source
workspace, tag, commit, run ID, and run attempt. Do not import S043/S047
local-demo pass results.

Complete every generated fragment and required attachment. The final evidence
must contain exactly 47 passing observations. Required screenshots must contain
supported raster bytes. Finalize once, then independently run
`windows-release-gate verify-bundle` against the MSI and candidate manifest.

Upload the resulting archive as:

`go-schedule_v1.0.0_windows-attended-evidence.zip`

Do not overwrite an unexplained existing archive. Requalification after a
candidate byte change requires a fresh tag/version rather than relabeling.

## Issue reconciliation contract

After formal qualification, run `windows-release-gate render-dispositions`
against the exact evidence archive, candidate manifest, MSI, repository, tag,
and commit. Review the resulting `packet.json` and ten Markdown records before
posting the relevant record to each of #98, #101, #104, #105, #106, #109, #111,
#112, and #113. Close only issues whose complete criteria pass. Update #96's
child index and checklist from actual issue state, then close it only when every
coordinator criterion passes.

The packet is offline review input, not closure authority. Generation never
comments on, closes, labels, or otherwise mutates a GitHub issue.

Issue #122 remains open until promotion and final audit. Post-v1 issues remain
unchanged.

## Promotion contract

Dispatch `Promote Release` with tag `v1.0.0` only after:

- the draft contains the exact nine pre-checksum payloads;
- candidate and evidence verification pass;
- all ten v1 readiness issues have passing dispositions; and
- the remote tag still resolves to the reviewed commit.

The workflow must re-download assets, generate and verify one checksum per
payload, re-download the final set, verify cardinality, and then change the
existing draft to public. It must never rebuild or substitute an artifact.

## Final audit contract

After promotion, verify:

- the release is public, latest, and tagged v1.0.0 at the reviewed S049 commit;
- all ten payload files exist, are non-empty, and every payload digest passes;
- release notes contain the reviewed five highlights and tagged changelog link;
- README, changelog, release URL, tag, manifest, and binaries identify v1.0.0;
- the ten readiness issues and issue #122 have accurate final evidence/state;
- the v1.0.0 milestone has no open issues before it is closed; and
- `main` remains clean and synchronized.

Only then may v1.0.0, S048, and S049 be reported complete.
