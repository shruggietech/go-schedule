# S130 Data Model

No application or helper schema changes. Application schema 24, preferences 1, metrics 4 and evidence 2 remain current.

| Entity | Identity and validation | Transition |
| --- | --- | --- |
| Dependency selection | Published SQLite1.61.1 plus exact libc1.79.0, complete root/desktop checksums | Held pair to coherent selected pair |
| Current store | Existing schema and durable records; unsupported state refused without mutation | Empty to current; valid current to reopened current |
| Prepared statement | Same owned connection and query, alternating bound values with stable results | Prepare, repeated bind/query, close |
| Cancelled operation | Supplied context, owned lock and transaction, no partial durable change | Active to lock wait to context error and rollback |
| Release evidence | Existing v2.0.0 source identity, later immutable assets and public release | Source PR, owner merge, exact tag/draft approval, concrete promotion approval, public discovery |

Source acceptance does not assign tagged-package or attended-native identity. Existing immutable historical records stay unchanged.
