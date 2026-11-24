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

This project opts into the `workflowy` skill and Workflowy Desktop MCP for substantive work. Load the skill when available. A running, signed-in desktop app with Local MCP enabled and a connected agent session is required.

- Read **AI Agent Workspace → ABOUT** and locate the project under **PROJECTS**. Reuse its branch or insert it alphabetically. Use **PROJECTS → _general** for work without a named project. Confirm the destination before writing.
- Keep named projects in `info`, `notes`, and `spec-kit` when the project uses Spec Kit. Keep `info` current and link its claims to evidence. Use `notes` for candid drafts, decisions, questions, and sessions in chronological order. Start each major session and substantive update title with one local ISO 8601 timestamp and its offset, followed by a short fragment. Keep updates and replies with their session or discussion.
- Put each formal slice directly under `spec-kit` in code order. Its title contains only the formal code and short title, like `S094 Check source order`. Add a `Date` child and an `Epic` child when applicable. Put scope, outcome, and links in children. Use the native Workflowy date and project hashtag together on the timeline child for dated records.
- Find the existing record by meaning and source before writing. Link to the exact finding or authoritative Git, tracker, or project artifact. Add new evidence or a delta in place and preserve the original node ID. Follow ABOUT for shared catalog fields, alphabetical order, and link roles.
- Keep private work with the operator and their authorized agents. The operator decides what is published. Read back MCP edits and check titles, parents, order, links, dates, and duplicates. If MCP is unavailable, continue the work and report what still needs to be saved.

Repository-local skill: [.agents/skills/workflowy/SKILL.md](.agents/skills/workflowy/SKILL.md). Use the bundled guides and templates. Keep private notebook content and connection credentials out of repository files.
