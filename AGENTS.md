# AI Agent Project Directives

## Mermaid diagrams

- Lay out Mermaid diagrams from top to bottom. For flowcharts, use `flowchart TB`; for nested subgraphs with an explicit direction, use `direction TB`.
- Do not use left-to-right, right-to-left, or bottom-to-top flow directions (`LR`, `RL`, or `BT`) unless the user explicitly requests one for a specific diagram.
- Keep the primary reading order vertically stacked. When a diagram would become too wide, split it into multiple focused diagrams or vertically arranged subgraphs instead of changing it to a horizontal layout.
- Treat top-to-bottom layout as a required project convention for every new or modified Mermaid diagram, not as an optional formatting preference.
- Apply this convention only to project-authored diagrams. Do not reject, rewrite, or warn about a different valid direction contained in a user-owned file being processed by the product.

## GitHub publishing format

- Never use the Unicode em dash in repository content, code comments, documentation, issues, pull requests, reviews, release notes, or any other text published to GitHub. Use commas, parentheses, colons, semicolons, or standard hyphens instead.
- Never hard-wrap Markdown prose to a source-width limit. Each paragraph, heading, and list item must occupy one physical source line, with GitHub responsible for visual wrapping.
- Use multiple physical lines only where syntax or literal content requires them, including fenced blocks, tables, front matter, separate list items, and deliberate semantic line breaks.
- Run `go run ./scripts/github-format` before publishing repository or GitHub content.
- Treat an em dash or width-based prose wrap as a publication-blocking formatting defect.

## Implementation slices

- Bundle as many compatible GitHub Issues into one implementation slice as can be completed, validated, and reviewed coherently.
- Do not default to one GitHub Issue per implementation slice. Split work only when dependencies, risk, reviewability, platform boundaries, conflicting validation needs, or independent delivery value make separate slices materially clearer or safer.
- Preserve issue-level traceability inside a bundled slice by naming every included issue, satisfying each issue's functional acceptance criteria, and reporting any issue whose functional outcome remains incomplete instead of closing it implicitly.
- Prefer slices that deliver a meaningful end-to-end capability or release increment over artificially narrow issue-by-issue churn.
- Do not close an issue merely because it was included in a slice. Close it when its individual functional acceptance criteria are delivered. Tests, CI, attended walkthroughs, qualification packets, evidence archives, and other verification activity MUST NOT be GitHub issue closure gates; any defect they reveal MUST be filed or updated as a separate bug issue against the affected release or pending source revision.

## GitHub project management

- Treat native GitHub Issues, labels, and milestones as the authoritative project-management surface unless the repository explicitly designates another system.
- When approved specifications, plans, requirements, review findings, or implementation discoveries create actionable work, proactively create or update the corresponding GitHub planning records. Do not leave durable work tracked only in chat, temporary notes, or an implementation plan.
- Before creating or modifying GitHub records, verify that the local Git remote identifies the intended repository and inspect the repository's existing issues, milestones, labels, title prefixes, sibling repositories, and organizational conventions.
- Search open and closed issues before creating anything. Reuse or update an existing issue when it already represents the work, and never create a duplicate merely because its wording differs.
- Organize substantial work under outcome-oriented or release-oriented milestones. Give every planned issue the appropriate milestone unless the repository convention explicitly permits an unmilestoned backlog.
- Use coordinator or epic issues for work spanning several child issues or milestones. Each child issue must identify its parent, and each coordinator issue must maintain a linked child-issue index and a completion checklist.
- Record dependencies with explicit GitHub issue references. State what an issue depends on, what it blocks when useful, and the order in which dependent work can safely proceed.
- Give every issue a specific outcome, bounded scope, functional acceptance criteria, dependency section, and links to the specifications, decisions, incidents, or review findings that authorize it. Acceptance criteria MUST describe product or planning outcomes and MUST NOT require verification activity for closure.
- Apply the repository's established label taxonomy consistently. Prefer structured label families for work type, functional area, priority, and estimated effort when those families exist.
- Treat “tags” as GitHub labels. Do not create synonymous, differently punctuated, or differently capitalized labels when an existing label already represents the same category.
- Create or revise labels when the roadmap introduces a durable new category. Every label must have a concise description, an intentional color, and a scope that does not overlap confusingly with existing labels.
- Follow the repository's established title-prefix convention, such as lifecycle, stability, release-stage, or work-type prefixes. Inspect existing usage before selecting a prefix, apply it consistently, and do not invent a new prefix for a single issue.
- Do not invent milestone deadlines, release dates, priorities, assignees, or completion claims without supporting project authority.
- Keep planned capabilities visibly distinct from shipped capabilities. Closing an issue or milestone must not cause documentation, associations, release notes, or public claims to advertise behavior that has not been implemented or published as stated.
- When completed historical work must be represented, create a clearly identified retrospective issue or milestone, attach the available implementation references and known defect links, and close it as completed without rewriting project history.
- After bulk creation or reorganization, audit the complete planning surface for duplicate titles, missing milestones, missing required labels, broken issue references, missing parent links, invalid dependency ordering, and incorrect open or closed states.
- Prefer idempotent automation for bulk GitHub management. Match existing records by stable identity or exact canonical title so an interrupted run can resume without creating duplicates.
- Never create, modify, close, or reorganize issues in a repository that does not match the verified Git remote or the repository explicitly named by the user.

## Workflowy

This project opts into the `workflowy` skill and its private agent workspace.
Load the skill when available. Use the operator's project link or find this
repo or folder under **AI Agent Workspace > PROJECTS**. Reuse its branch or
insert it alphabetically. Search **STANDARDS** by title or topic tag. Read the
applicable rules, defaults, and scopes before governed work.

Use normal bullets for candid discussion, drafts, decisions, and evidence.
Never store work in Workflowy comments or comment drafts. Give `info` exactly
three direct children: `tags`, `links`, and `details`. Put purpose, current state,
next actions, blockers, and supporting evidence under `details`.
Use lowercase letters and underscores for reusable field names. Use inline
code for exact technical literals. Link to the exact canonical finding or evidence.
Git and the established tracker govern official progress where used.
Otherwise use the project artifacts. The operator decides what is published.

In `notes`, keep major sessions and substantive updates as a flat, oldest-first
list. Use only the ISO 8601 date or timestamp as the title. Retain the finest
supported precision and offset. Put the description, scope tags, and prose
under `title`, `tags`, and `details`. Use label-only field bullets with each
value or link nested in its own bullet. Require the Workflowy bullet format
for fields and values. Only the contents of `details` may use other formats.
Keep tags out of names and other fields. Use one hashtag per `tags` value.
Put any project-level tags in `info > tags`. Link related updates. The operator
owns timeline construction. Agents record dates without creating timeline
metadata, views, saved searches, or exports.

For Spec Kit, use one direct `spec-kit` bullet per formal slice in code order.
Use only its code and short title. Nest the verified date under `date` and
the documented epic under `epic` when applicable. Use the skill for shared
reference records, migration, and precise field rules.

Desktop MCP requires a running, signed-in Workflowy app with Local MCP
enabled and a connected client. A cloud or other nonlocal agent uses a
separately configured, operator-authorized connector when available. Confirm
the destination and read back writes. If no connector is available, continue
under the available instructions and provide an appropriately private,
ready-to-paste nested outline with its exact destination. Mark it unsynced.
