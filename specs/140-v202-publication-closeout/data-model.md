# S140 data model

The existing comparison schema contains tag, assessment date, eight product contexts, ten stable capability rows and 80 cells. Each cell has status, bounded note, official source URL and its actual evidence date. Preserve opportunity links and all row/product identities. No schema migration or source_commit self-reference is added.

The existing review identifies refs/tags/v2.0.2 and binds SHA-256 identities of the exact dataset, selected changelog section and tag-specific notes. Its ten decisions are changed or no-change; review-required is a deliberate rejected intermediate state. The final commit is supplied to validation externally, avoiding a self-referential committed hash.

Publication state remains separate: prepared source, owner-merged main, approved exact-source draft staging, inspected package inventory, approved public promotion and discovery merge. Current published.json remains v2.0.1 at the first source checkpoint. Parent #444 closes only when actual public delivery is complete; #452 and #453 close with their individual source outcomes.
