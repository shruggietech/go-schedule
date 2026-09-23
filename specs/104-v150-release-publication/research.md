# Research: v1.5.0 release publication

## Current boundary

S103 merged as `8d04509f41c3ecd3dd12ff31c9cc2f993fa7d781` with successful exact-main CI, but the tagged-source workflow would reject it: README still shows a v1.4.0 health example, the v1.5.0 changelog heading lacks a publication date, and the tag-specific changelog link uses the undated anchor. A reviewed metadata correction is necessary before tagging.

## Staging and promotion

`.github/workflows/release.yml` builds a draft only after exact-commit main CI and source metadata checks. It stages four headless archives, two desktop archives, a Windows MSI, and a Windows candidate manifest. `.github/workflows/promote-release.yml` additionally requires a genuinely passing attended Windows evidence archive. Neither a staging success nor a waiver is equivalent to that archive.

## Historical precedent and decision

The public v1.4.0 release used an explicit maintainer-authorized native-testing waiver documented on #226. Its untested checks were disclosed, not converted to passes, and publication used a manual waiver path after exact-source CI, successful staging, candidate verification, asset auditing, and checksums. The maintainer clarified during S104 that this native-testing waiver is standing policy across releases, not version-specific. Keep the full promotion workflow intact and use the new reviewed standing-waiver promotion workflow for this and later waived releases.

## Deferred scope

SMTP #176, notification coordinator #19, and distributed scheduling research #20 remain open outside this release. Public copy must not imply their delivery. The optional native popup exists only while the desktop GUI is running.
