# Data Model: Windows Release Polish

This slice introduces no persistent product data. Its test and evidence contracts use the following value objects.

## Work Area

- **Physical width and height**: Positive pixel dimensions for one selected monitor's reserved work area.
- **Scale**: Positive device scale used to derive logical dimensions.
- **Logical cap**: 90 percent of each logical work-area dimension.
- **Requested size**: Each preferred dimension independently bounded by its logical cap.
- **Fallback state**: Unknown work area or invalid scale returns the preferred 1280x800 request.

## MSI Summary Evidence

- **Artifact path**: Resolved candidate or published MSI.
- **Artifact origin**: Provenance string or release URL.
- **SHA-256**: Lowercase content hash.
- **Summary Subject**: PID 3 string, exactly `go-schedule: cross-platform task scheduler`.
- **Evidence status**: Proven only when every existing compiled contract and the Subject contract pass.

## State Transitions

- Work area unknown -> preferred fallback.
- Work area known -> logical conversion -> independent 90 percent caps -> restored size request.
- MSI opened -> hash and Summary Subject read -> pass or failure evidence emitted.
