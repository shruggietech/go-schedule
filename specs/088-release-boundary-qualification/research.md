# S088 Research

## Metadata

Decision: move release-bound S083-S085 product improvements and the S087 selector repair into the existing v1.4.0 section, retaining its heading and tagged anchor. Preserve incomplete candidate-refresh history as unfinished work.

Rationale: the next reviewed candidate must describe its actual desktop source. README health already identifies 1.4.0; its latest-release badge does not claim the draft is public.

Alternatives considered: stale release copy would omit the repairs; inventing a publication date would falsely imply publication.

## Sandbox and waiver

Decision: reuse Windows Sandbox without host changes or restart. Execute supported checks after reviewed-source staging and retain unavailable checks as untested under the explicit maintainer waiver recorded in #226.

Rationale: no VM management tools, Hyper-V namespace, configured remote session or active Sandbox was found. The maintainer rejected provisioning another platform and authorized release despite missing tests.

Alternatives considered: another VM platform contradicts the maintainer's direction; counting obsolete screenshots or hosted fixtures as current native evidence would be false.

## Sequencing

Decision: review and merge the source preparation before staging. Keep public promotion and any future waiver-aware gate implementation explicit, without bypassing the existing full-evidence validator.

Rationale: draft source 951d864d9683ec3bdcb1e37535ecddd67738bf13 predates the selector fix in reviewed main c48ee096251b33187deda61200eaa25edd0e36ab. The release workflow requires successful exact tagged main CI. Promotion currently requires a fully validated archive; a waiver is not a passing archive.

Alternatives considered: tagging the unreviewed branch or silently disabling validators violates source and evidence integrity; a replacement framework is disproportionate.
