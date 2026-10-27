# Source Handoff Contract: S124

The source proposal identifies v1.7.0 and refs #394. It is not a tagged release candidate manifest.

- Record accepted base and exact proposed head separately from the assembled-source digest and actual 0.0.1-dev QA binaries.
- Record source-only release/comparison guard results without claiming that future tag URLs exist.
- Record current generated QA run/corpus/report identities, semantic review, cleanup and frozen historical inputs separately from canonical CI and review receipts.
- Preserve public v1.6.0 selection, v1.6.0 historical artifacts, development version fallback and existing package upgrade identity.
- Preserve the standing native-testing disclosure: attended native observations remain untested under the standing native-testing waiver; automated CI and artifact checks remain required. The existing authorization category is "not tested, release authorized by maintainer". No new release authorization is implied.
- Later staging requires owner-authorized exact tag/commit, successful Release run/attempt, actual package hashes and Windows manifest identities. Eight staged assets are expected; promotion adds SHA256SUMS.txt as the ninth public asset. These identities remain pending during S124.
- Final source handoff lists open #355/#356 causes without a resolution claim, leaves #394 open for delivery, and stops before merge/tag/staging/promotion.
