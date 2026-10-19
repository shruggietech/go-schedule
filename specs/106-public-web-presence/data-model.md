# S106 data model

- Document: filename, route, legacy aliases, title, navigation order, headings, body and assets. Reject route collisions and links outside declared content. Generated copies are not authored.
- Product: ID, display name, edition/deployment and research date. Eight records include go-schedule and the seven named alternatives.
- Capability: ID, category, name, product cells and opportunity links. Cells contain supported/partial/planned/unavailable/unknown status, explanation, official URL and date. Snapshots identify the released go-schedule baseline.
- Review: exact tag, `refs/tags/<tag>`, date, dataset/release-content SHA-256 and complete unique per-capability change/no-change decisions. Record the resolved runtime commit rather than a self-referential committed SHA.
- Domain handoff: intended records, DNS-only state, Pages settings, challenge TXT and redacted before-state. Credentials are never model fields.

A snapshot becomes public only when deployment resolves its matching public release. A missing matching snapshot stops that build instead of advertising main-only capability.
