# Report and Review Contract

JSON Schema 2020-12 v1 documents are strict at every object, require all contract fields, preserve explicit nulls, bound strings/collections and use UTC RFC3339 timestamps and integer nonnegative counts/durations. Shared concrete Go types and embedded schemas must agree.

The report follows #290's required top-level shape: schema_version/kind/run_id, started_at/finished_at, subject, runner, environments, summary, cases, communication_reviews, findings, artifacts and cleanup. Add explicit phase/completion_state/corpus_sha256 and execution_status to distinguish capture from final review. Every manifest/environment case occurs exactly once.

Pass requires actual exit/check evidence, all applicable manifest outcomes, complete contextual message reviews and successful cleanup. An intentional rejection passes only against its specified exit/output/state contract. Not_applicable requires a manifest-defined reason. Missing prerequisites or reviews are blocked; product/harness failures are fail with separate reproducible findings. Counts and overall status are derived and independently revalidated.

Reviews bind run_id, schema_version and corpus_sha256. Require exact immutable message/case/channel references and excerpt, explicit rubric verdicts for plain language, action/outcome, correction, severity, terminology, identifiers/paths/units and secret absence, plus finding links. Missing/duplicate/stale/extra reviews or fabricated no-message text are rejected.

Evidence uses confined portable paths, regular files, bounded bytes, hashes and redaction metadata. Capture sanitizes complete bounded streams in memory, registering fixture-only dynamic secrets before persistence. Evidence manifest covers frozen case artifacts/events/corpus; derived reports do not recursively hash themselves. Semantic validation checks uniqueness/references, expected checks, completeness, timestamps, totals, spans and cleanup truthfulness.

Human report outline: 1 identity/outcome, 2 environment/budget, 3 coverage, 4 functional results, 5 help/error/log review, 6 defects in fixed field order, 7 cleanup/artifact index. All content derives from canonical JSON with escaped Markdown and stable sorting.
