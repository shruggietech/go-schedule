# S112 research

## CI baseline

- Decision: correct queued-record masking in the probe, not core overlap policy. Unchanged source selects the first matching trigger before ended_at; queued and executed observations can coexist. The hosted timeout cause remains unconfirmed until disposable CI evidence is available.
- Alternative: increasing waits or changing scheduling overlap. Rejected because neither addresses the demonstrated selection defect.
- Decision: separate bounded credential build setup from runtime. The shared 45-second context killed compilation before assertions; local warm baseline passed and is separate evidence.
- Alternative: skipping packaged helper checks or reusing an arbitrary binary. Rejected because actual executable identity is part of the contract.

## Dependencies and release

Pending upstream research is recorded below before corresponding changes. Use current PR patches and official action/Wails release notes. Preserve v1.5.1 public discovery; v1.6.0 is a prepared target until separately authorized publication. No schema or scheduler behavior expansion is proposed.

## Reviewed upstream migrations

- Root #351 includes #273 x/sys scope; preserve SMTP additions and reconcile desktop root replacement before tidying.
- Wails runtime #336 targets 2.17.0 while active CLI build paths still pin 2.15.0. Align runtime, CI, release workflow, verification and policy fixtures to 2.17.0 together. [Official release](https://github.com/wailsapp/wails/releases/tag/v2.17.0). Preserve dev-server authorization and hidden process behavior; do not restore the removed PowerShell toast fallback.
- Pages configure6/deploy5 use Node24; upload-pages-artifact5 uses upload-artifact7. Existing hosted runners support these. Preserve current permissions, environment, artifact name/path and deployment output. Hidden-file exclusion already exists in v4 and is not broadened. [configure6](https://github.com/actions/configure-pages/releases/tag/v6.0.0), [deploy5](https://github.com/actions/deploy-pages/releases/tag/v5.0.0), [upload5](https://github.com/actions/upload-pages-artifact/releases/tag/v5.0.0).
- Keep exact npm manifest/lock pairs from #337/#295; React and react-dom must align, as must Tailwind/PostCSS adapter. Browser tooling requires matching managed binaries. [Playwright browser contract](https://playwright.dev/docs/browsers#install-browsers).
- Release tag workflow requires main CI, version-bound health example, changelog/notes and comparison guard. Prepare v1.6.0 material while retaining site/data/published.json at v1.5.1. No legitimate immutable staged-candidate manifest can be made from a pre-tag local tree.

## Comparison review, 2026-10-22

Reviewed all 70 vendor cells across ten capabilities with official sources for JAMS, ActiveBatch, ScriptRunner, Control-M, Logic Apps, Resolve and OpCon. Updated precise access/workflow/MCP/HA/deployment claims with their edition and preview limits, replaced inaccessible Control-M XWiki and stale OpCon CLI evidence, and retained unknown where sources do not establish availability. MCP remains distinct from the Unisys MCP platform. The pending go-schedule snapshot records delivered SMTP and CLI coverage and removes completed #176/#277 opportunities; conditional #20 remains. Historical v1.5.1 content and the published pointer are untouched.

## Independent payload observation

Local canonical repetition exposed a PowerShell internal CLR crash and a separate SQLite prerequisite refusal in unchanged concurrent payload writers. The earlier suite and five focused race repetitions passed with the same six-writer assertions. #355 retains the original failure, runtime identities and evidence hash without assigning a shared cause or claiming a fix. Local verification uses two Go build workers to bound host load; writer concurrency and assertions are unchanged.

## Final lifecycle review

Reused existing clifixture contained process trees for packaged fixture compilation so cancellation drains compiler descendants on Windows Jobs and Unix groups. Removed Recorder.SecretStore composition from the generic fixture package and constructed the unchanged native adapter at its seven caller composition roots, avoiding an import cycle and duplicated launcher logic. Probe review added explicit skipped-run handling and cleanup-only failed status, each with executable regressions. No scheduler policy or credential authority changes result.

## QA reporting disposition

The original default-budget capture exceeded its catalog budget and the fixed PowerShell wrapper cap prevented aggregate finalization (#357). A fresh direct-coordinator run selected a bounded two-hour total catalog budget while preserving the 15-second case deadline and all CPU/memory/process/output limits. This is a recorded invocation decision, not a source fix. The original five Windows deadline observations (#356) did not recur; cause remains unconfirmed. Full execution passed on both platforms, while contextual review identified neutral SMTP acceptance text serialized as last_error (#358). The valid failed report and all original evidence remain retained rather than relabelled as passing.
