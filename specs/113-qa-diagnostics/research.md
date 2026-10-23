# S113 research and decisions

## #357: budget and terminal capture

Decision: share the Go-validated timing policy with the PowerShell entry point and separate execution cancellation from bounded worker shutdown/report finalization. Skip unchanged sanitized artifact rewrites only after exact regular-file verification.

Evidence: the wrapper fixes 1900000 ms regardless of supported selected totals up to 7200000 ms. The coordinator's total-plus-45-second child is under the already-expiring total parent, so it supplies no independent reserve. Capture process cancellation allows only seven seconds before owned-tree termination. freezeChannels/freezeProgress perform repeated fsync/rename of unchanged files; fallback repeats the work. Original root remained capturing/running while 644 per-case files retained more progress than the periodic 635-pass/five-failure checkpoint.

Rationale: correct deadline ownership and unnecessary I/O rather than increase evaluation limits. The exact historical finalization stall remains unproven. Alternatives rejected: fixed larger wrapper cap, silently selecting a different budget, discarding late progress, or presenting process absence as successful runtime cleanup.

## #356: partial observations and causal attribution

Decision: preserve typed completeness/failed-request facts and partial data, distinguish execution/protocol/observation errors, and collect bounded context-scoped phase/request timings. Instrument and reproduce before choosing a transport fix.

Evidence: normal nonmetadata cases share 15 seconds across four ownership-guard queries, twenty snapshot queries and command/protocol work. snapshot returns nil on its first failure; postcondition errors overwrite partial details. service.restart and mcp.http.enable returned successful command results before inspection failed; mcp.http.rotate failed independent group cleanup; two other commands were canceled. None proves a named-pipe, dependency or host cause.

Rationale: preserve diagnostic truth and fail-closed complete assertions. Alternatives rejected: infer a go-winio culprit, blind mutation retries, or raise case deadlines to hide stalls.

## #355: runtime and prerequisite failures

Decision: numeric candidate validation, bounded candidate stderr/exit/version diagnostics and interpreter outcome metadata, while preserving six concurrent writers and exact record assertions. Keep the CLR cause unconfirmed unless reproduced and explained.

Evidence: retained log SHA256 4eded1dd882f92b7f09aad62265065c46631212d79d97ce7b12922a4b2d3cacf contains a PowerShell 0xc0000005 / CLR0x80131506 reflection/type-metadata crash and a separate Bash prerequisite refusal. Both twins discard failed candidate stderr. Bash accepts malformed nonnumeric version tokens using sort comparison alone. The fixture prerequisite checks only executable lookup despite its stronger comment.

Official related reports have different stacks/architectures and do not establish the local cause: [PowerShell #26665](https://github.com/PowerShell/PowerShell/issues/26665), [runtime #131636](https://github.com/dotnet/runtime/issues/131636). [Microsoft crash collection guidance](https://learn.microsoft.com/en-us/dotnet/core/diagnostics/collect-dumps-crash) is future controlled-reproduction context, not permission for machine-wide dump configuration or runtime replacement.

Alternatives rejected: retry a crashed writer into success, serialize writers, lower expected counts, or claim an unrelated upstream issue fixes this signature.

## #358: successful SMTP diagnostics

Decision: additive diagnostic plus migration24; copy original SMTP last_error into diagnostic and clear old last_error only when saved state is succeeded. Preserve authored error details for partial rejection/retry/uncertain outcomes and map both fields consistently across clients.

Evidence: smtp.go authors neutral SMTPAttemptResult.Diagnostic, then FinishSMTPAttempt saves it unconditionally into last_error. Success already derives from every recipient accepted, so no heuristic status/prose parsing is necessary. Domain and OpenAPI already permit additive fields; explicit desktop adapters require mapping updates.

Rationale: preserve exact useful history and explicit actual-error semantics through a forward-only migration. Alternatives rejected: remove the summary, relabel state, parse SMTP prose, or change recipient/retry behavior.

## PR #360

Decision: incorporate core/UI16.17.7 and MDX15.4.8 with the existing two-file lockfile patch. Official [core/UI release](https://github.com/fuma-nama/fumadocs/releases/tag/fumadocs%4016.17.7) and [MDX release](https://github.com/fuma-nama/fumadocs/releases/tag/fumadocs-mdx%4015.4.8) describe Windows cross-drive/UNC and dot-folder import corrections. These are patch proposals, not compatibility evidence; real builds and policy checks remain required.

## Authorization

Push and official PR publication are explicitly authorized in this S113 kickoff. No merge, auto-merge, tag, release, global runtime change or DNS mutation is authorized. Up to two Codex invocations, with every actual comment handled individually. Issue closure remains functional rather than dependent on a verification packet.
