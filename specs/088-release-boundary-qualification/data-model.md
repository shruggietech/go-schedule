# S088 Data Model

No runtime schema changes.

- Release boundary: version, retained heading/anchor, cumulative entries, tag-specific highlights, reviewed commit, draft/public state. Summaries must include the merged UI repairs without claiming qualification or publication.
- Candidate identity: repository, tag, commit, workflow run/attempt, eight asset names, lengths/digests, and MSI identity. Sequence: review, merge, exact-main CI, authorized staging, identity verification, supported observations, release disposition.
- Observation: scenario, candidate identity, environment, actual result, attachments, unavailable capability, and waiver reference. Waived untested results cannot become passed without actual evidence; historical observations retain their original identity.
- Review handoff: PR, current head, local gates, findings/dispositions, round count, latest-head CI, and remaining issue criteria. Maximum rounds: two. Preparation merge does not implicitly close #231, #228, or #226.
