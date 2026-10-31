# S129 Research and Decisions

Date: 2026-10-31. Source baseline: 0c1ab32129a3bd3fcb4694915aaccef4166a502c. Installed Spec-Kit Phase 0 dispatched three read-only research agents for version/module identity, comparison binding and delivery. Their original reports remain private under .test-temp/s129.

## Candidate identity

**Decision**: Prepare v2.0.0, with root Go module github.com/shruggietech/go-schedule/v2 and nested desktop module github.com/shruggietech/go-schedule/v2/desktop. Update active Go imports, local module dependency/replace and both existing release linker targets together. The nested desktop module is built locally and is not independently tagged. Keep buildinfo's development fallback unchanged.

**Rationale**: S128 removed the documented preferences.transition JSON/human CLI output, Fyne import and historical-state conversion. The changelog declares SemVer. Having no users removes migration obligations but does not make an incompatible interface change backward compatible. Go's module-aware v2 tag requires the /v2 root namespace; the desktop importer must share that parent for internal-package access. HTTP API routes and current persisted formats do not change.

**Alternatives considered**: v1.7.1 would misclassify incompatible removal as a fix; v1.8.0 would label it a compatible increment. A v2 tag with an unsuffixed root module would violate Go module version rules. New HTTP API routes, a release framework and historical converters are unnecessary.

**Sources**: [Semantic Versioning](https://semver.org/spec/v2.0.0.html), [Go major-version module guidance](https://go.dev/doc/modules/major-version), [Go module version suffix rules](https://go.dev/ref/mod#major-version-suffixes), current CHANGELOG.md, docs/cli.md, internal/userprefs/model.go and root/desktop go.mod.

## Source and public identity

**Decision**: README's exact health example identifies explicitly unpublished prepared v2.0.0 source, satisfying the existing release-source-check. Its public package availability, latest badge/link and site/data/published.json continue to identify actual v1.7.0. Current install guidance identifies prepared v2.0.0 behavior and links this same S129 continuation. Preserve historical version notes and comparison records.

**Rationale**: The existing source guard checks one candidate-specific health example and exact dated changelog/notes identity. Explicitly labeling that example avoids a new validator mode while preserving truthful download claims.

**Alternatives considered**: Switching public discovery early implies unavailable packages. Weakening the source guard or creating a second release mechanism adds complexity without a requirement.

## Comparison

**Decision**: Add v2.0.0 snapshot and content-bound review, covering all ten capability rows and eighty qualified cells. Preserve all statuses; refresh only claims whose official source was actually reviewed, retain source-specific historical dates otherwise, and update the clustered-execution opportunity to the owner's explicit deferral. Keep historical snapshots/reviews and public selection unchanged.

**Rationale**: Current-only initialization/preferences/helper/qualification simplify maintenance without adding a matrix capability. All seven named competitor contexts and specialized MCP/HA/access/workflow sources were reviewed from official documentation. Documentation review does not establish hands-on qualification or resolve conflicting vendor availability claims. ScriptRunner MCP remains unknown; Control-M and Logic Apps MCP retain preview/partial qualifiers; vendor editions/add-ons remain qualified.

**Alternatives considered**: Relabeling unknown as unavailable, editing immutable historical evidence or advertising cluster ownership/failover would overclaim. No new comparison engine is needed.

Primary sources include [JAMS](https://jamsscheduler.com/), [ActiveBatch](https://www.advsyscon.com/activebatch/), [ScriptRunner](https://www.scriptrunner.com/platform), [Control-M](https://documents.bmc.com/supportu/9.0.22/en-US/Documentation/Introduction_to.htm), [Azure Logic Apps](https://learn.microsoft.com/en-us/azure/logic-apps/logic-apps-overview), [Resolve Actions](https://docs.resolve.io/pro/Getting%20Started/architecture-overview/) and [OpCon](https://help.smatechnologies.com/opcon/core/v26.0/Files/UI/Enterprise-Manager/Using-History-Management). Exact cell sources remain in the dataset; the private research receipt retains the specialized documents and unresolved-claim rationale.

## Delivery and evidence

**Decision**: Reuse Release and Promote Release (standing native-testing waiver). This authorization ends at reviewed source PR readiness. Later owner merge, exact tag/draft approval and concrete-asset promotion approval continue within S129 and #412. Expected staging is eight assets; promotion adds only SHA256SUMS.txt and publishes the accepted original bytes.

**Rationale**: Existing workflows already enforce exact source, successful CI/staging attempt, manifest/MSI identity, frozen release/asset IDs, digests and idempotency. Attended native observations stay untested. Source tests, CI, artifact readback and public delivery are separate evidence classes.

**Alternatives considered**: Publishing now exceeds authorization; rebuilding accepted assets invalidates acceptance; another preparation slice or framework adds no functional outcome. #355/#356 causes, SQLite #400/#404 and all cluster outcomes remain unresolved/deferred.

## Research completion

All technical unknowns are resolved. No new dependency, schema conversion, HTTP interface or scheduling change is proposed. Implementation remains blocked until tasks and read-only cross-artifact analysis pass.

## Source health identity correction (#413)

**Decision**: Require the full selected tag in the README health sample and validator. Existing release.yml injects v2.0.0 unchanged; API health returns buildinfo.Version unchanged and CLI prints it unchanged. The earlier validator stripped v, enforcing inaccurate source documentation. The two-case regression fails before the fix (prefixed rejected, unprefixed accepted) and passes after. Existing fixed-version synthetic negative fixtures remain, updated only for the actual tag-output contract. Immutable old tagged documents are not edited.

**Rationale and alternatives**: Correcting the current source contract resolves FR-002 without changing daemon/API behavior or existing workflow injection. Changing runtime output would expand a documentation/tooling correction into an interface change. Leaving the enforced mismatch for another slice would ship knowingly inaccurate release metadata. Native #413 is a child of #412 and closes on its own source outcome, while actual release delivery remains open.
