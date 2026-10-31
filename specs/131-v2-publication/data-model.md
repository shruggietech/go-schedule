# S131 Release Identity Model

An accepted source has commit/tree and terminal workflow receipts. A candidate has release ID, immutable tag/source, staging run ID/attempt and eight assets, each with asset ID/name/size/digest. The Windows manifest and MSI identify the same product, version, source and staging run. A publication preserves all original assets and adds SHA256SUMS.txt with eight entries, public/latest state and Pages/discovery receipts. The final source PR has its own commit, canonical validation and at most two review requests; earlier source/staged/attended evidence is not interchangeable.

State progression: accepted source -> explicitly approved staging -> complete inspected draft -> separately approved promotion -> actual public/latest -> discovery -> source reconciliation ready -> owner merge. Failure/pending states retain exact receipts and do not advance automatically. No application data model changes.
