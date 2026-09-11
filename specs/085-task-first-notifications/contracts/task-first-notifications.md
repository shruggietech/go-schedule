# Task-First Notifications Contract

## Workspace projection

- The desktop workspace retains channels, tasks, groups, deliveries, and loaded time.
- The workspace adds `coverage`, an ordered list of secret-free configured task and group summaries, plus `coverageComplete`.
- Task summaries use effective assignments so inherited group configuration is visible.
- Group summaries use direct assignments because group inheritance applies to descendant tasks rather than parent groups.
- A summary with no assignments is omitted.
- Destination counts distinguish configured channels from enabled configured channels and preserve enabled success and failure coverage independently.
- A failed per-scope lookup does not discard channels or delivery history; `coverageComplete` becomes false.
- Serialized workspace data must not contain full endpoints, authorization values, credentials, or notification payloads.

## Initial overview

- The heading describes notification status and recent results without leading with webhook or policy terminology.
- Status includes active destination count, configured task and group count, recent task-outcome count, and a text state.
- State guidance identifies a relevant advanced section or states that no action is needed.
- At most five newest results appear before all advanced summaries.
- A recent row names task context or Test notification, destination, state, time, and guidance.

## Advanced disclosures

- Destinations, Assignment rules, and Delivery diagnostics are three independent details and summary regions.
- All three are collapsed on page entry and expose their expanded state to accessibility APIs.
- Destinations retains channel listing, create and edit, secret replacement, enable and disable, test, and remove confirmation.
- Assignment rules retains scope selection, direct assignments, effective policy, inheritance, warnings, save, and override clearing.
- Delivery diagnostics retains filters, all 200 bounded records, selection, and redacted detail.

## Responsive and accessibility boundary

- Initial and expanded states avoid document-level horizontal overflow at 800 by 600 and 200 percent zoom.
- Summary controls, form controls, result selection, and dialog actions remain keyboard reachable with visible focus.
- State meaning never depends on color alone.
- Specialist explanations are visible inline or associated through accessible help text.
