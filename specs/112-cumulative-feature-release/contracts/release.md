# S112 release and diagnostic contracts

- Terminal selection: match task/trigger and a completion timestamp or explicit skipped outcome before selecting the newest available matching record; an unsuccessful terminal outcome remains visible to the caller.
- Failure retention: capture before destructive cleanup, record cleanup independently, rethrow the original failure even if a secondary diagnostic fails.
- Credential fixture: bounded setup and fresh runtime contexts; original authorization ordering and two-second channel read deadline remain intact.
- Publication: source PR closes only delivered #352/#353 outcomes. #354 remains open until the separately authorized cumulative public release exists. No tags, release workflows or installed service changes occur in this source-preparation phase.
- Dependency disposition: every pending PR remains traceable; duplicates are superseded only after the replacement is delivered.
