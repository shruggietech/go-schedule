# S117 candidate evidence contract

Reuse `internal/clievaluate/schemas/v1` without adding a parallel schema. The review document is schema_version string 1 and kind go-schedule.cli-qa-reviews, bound to the exact run ID and corpus SHA-256. Each distinct required message has exactly one review with exact ordered cases, text and channel, all seven rubric reasons, and a consistent overall verdict. Failed messages and execution failures have actionable findings with retained path/hash/inclusive line spans and functional acceptance criteria.

The generated report retains seven sections: run identity/outcome; environment/budget; coverage; functional results; help/error/log review; defects; cleanup/artifact index. Never edit generated report Markdown. Build/capture/finalization retain the same source and binaries. Later fixes use independently named acceptance evidence.

The source delivery record names #354, #355 and #356 individually, actual candidate identity/outcome, original causes still unresolved, any demonstrated corrections, final source CI/review evidence and later owner actions. It does not close the release issue or advertise unpublished behavior.
