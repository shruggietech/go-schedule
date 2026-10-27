# S123 research and decisions

## Original source recovery (#355)

Read-only Spec-Kit research found qualified-source-index.json, qualified-source.patch and all 61 sparse overlay paths under the retained S112 artifact root. The index records base 2b827e0, dirty=true and ff465894a3ab9e7a90720f5b7388a5f44d7a40b25bf3bddaf2c8637f61f991af. The patch has 45 tracked changes. Recovery must verify the full source digest, not just apply a plausible later commit.

Original binding test blob d3a1afe381e4b7af5040354575076dd374fc28f9 and SQLite adapter blobs 6ab7db0ce4fb813c94ce60bac2d4e093d2301d56 / 0e3038ad877e48422064348d9c0481e4540897a3 are untouched by that overlay. The failed test generated wrappers that source adapters. Heartbeat entrypoints and hidden Windows launcher are unchanged, whereas current wrappers add preflight/phase/incident instrumentation and bounded observation. Decision: run original and current tests separately and identify their controller differences.

The original PowerShell stack and native 0xc0000005 establish a crash, not an upstream cause. No matching original dump exists. Bash refusal is a separate result. Alternatives rejected: asserting PowerShell updates fixed it, comparing only unchanged heartbeat entrypoints, or calling current instrumentation an original invocation.

## Historical IPC compatibility (#356)

The five target definitions plus group prerequisites remain unchanged. Existing replay verifies all six prepared binary hashes. It currently fixes case allowance at 120 seconds and enables collection, so historical selection needs explicit duration and unavailable-collector attribution.

Current service.restart uses an in-process recorder command harness. Other targets invoke the selected CLI against its selected packaged WAL daemon. Component memory/WAL modes host current code and cannot claim historical server identity. Current MCP HTTP postconditions use 15-second requests; historical controller used three seconds. Decision: initially compare packaged candidates with the current controller, clearly labeled, and avoid claiming exact original controller reproduction.

Alternatives rejected: modifying original binaries, treating a missing historical private channel as a recurrence, silently reverting supported QA defaults, or expanding into mixed component/dependency matrices without an observed divergence.

## Delivery and governance

Reuse #355/#356 and their milestone/labels. A new demonstrated defect gets its own issue with functional acceptance criteria and explicit parent/dependencies; passing follow-ups cannot close the original causal issues. Reconcile S122 owner merge e8cefc1 and 23 successful main checks separately from 22 PR checks and both clean reviews on 8b0dda2 before the later CI-only correction.

No unresolved technical decisions remain. Protocol prerequisites, actual command contracts and retained source files are authoritative for local experimentation.

## Reconstruction preparation result

Base archive plus all 61 retained dirty paths reconstructs 2257 source paths, but its native Windows digest is bdcfca0f4b1505df92451986bf4a118f60076fe5c4529cf510143a34043a7950 rather than recorded ff465894a3ab9e7a90720f5b7388a5f44d7a40b25bf3bddaf2c8637f61f991af. The mismatch is retained; no exact whole-source recovery claim is made. Original binding wrapper/test and adapter blobs are independently identified and untouched by the dirty overlay. The compiled historical-target test is a reconstruction with that explicit whole-tree limitation, separate from the six original receipt-matched application binaries. No original is modified.
