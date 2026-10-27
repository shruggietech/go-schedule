# S121 verification

## Source and evidence identity

Accepted source is S120 squash merge 6e4834eea519dd0b5626f7919a101c4e9d4b4cf7. The implementation branch is codex/121-windows-incident-capture. The retained 26 original S112/S115/S117/S118 artifacts remain byte-identical. These original failed receipts are not replaced by later passing runs.

## Targeted contracts

Native Windows race checks passed for exclusive and bounded incident artifacts, immutable capture destinations, same-process runtime identity, ordinary exit 37 with stream overflow, unavailable/missing interpreters, native controlled FailFast, independent collection failure, fast unexpected failures, disabled mode and resume identity. The standard PowerShell wrapper forwards Incidents without altering its explicit 150 ms timeout/output regression. Native Linux race checks passed on existing Ubuntu-24.04 as h8rt3rmin8r, including real Unix-socket pool waits and handler-only contrasts. Linux PowerShell was unavailable; Windows supplied that interpreter coverage.

All four original IPC request shapes (tasks, runtime-info, credentials and group deletion) captured live single-connection SQL waits in memory and WAL stores before cancellation. Cancellation still completes the server before the connection is released. The handler-only contrast reports handler stage with an idle pool instead.

The retained native controlled crash produced a real Windows minidump: PID 26464, 7953714 bytes, SHA-256 d6934348510332c636419e52b22b70459d7a5e629b9899c921ddc9df2cec8411, native exit 0x80131623. Same-process pre-payload runtime/module identity matches that PID. The artifact stays private. Current PowerShell 7.6.9/.NET 10.0.14 and retained portable 7.6.6/.NET 10.0.12 remain separately identified; the retained runtime passed three repetitions of all PowerShell writer and SQLite-selection contrasts. No installed runtime change or historical crash explanation is inferred.

## Full source-bound CLI evaluation

Run 9ebb6317-af9a-47a8-afe2-f2dae97e83c2 used six frozen binaries built with existing Go 1.26.0, base 6e4834eea519dd0b5626f7919a101c4e9d4b4cf7, dirty=true, source SHA-256 8e9390af0e6d8b40ac9b98d319897d8718918964d3abff9bc55944ba328f0151. Capture selected Incidents=true, existing nonroot Ubuntu-24.04 and the default 120-second case/two-hour total policy, with CPU, memory, process, output and namespace controls unchanged. Capture ran from 2026-10-27T15:20:02.9342273Z through 2026-10-27T15:25:25.8443228Z.

| Environment | Passed | Failed | Blocked | Not applicable |
| --- | ---: | ---: | ---: | ---: |
| Native Windows 11 | 881 | 0 | 0 | 2 |
| Existing WSL2 Ubuntu | 882 | 0 | 0 | 1 |
| Total | 1763 | 0 | 0 | 3 |

All 1109 distinct messages received semantic review across all seven rubrics, including full help, expected invalid inputs, native errors, generated shell completions, SMTP outcomes, redacted credential output, daemon logs and MCP resource/tool descriptions and mutation refusals. No new CLI communication finding was identified. Finalize, validate and render all returned exit 0 with outcome pass. The canonical seven-section human report and JSON remain private under .test-temp/s121/qa-run-01/. All four disposable runtimes were removed with no remaining owned resources.

| Retained file | SHA-256 |
| --- | --- |
| report.json | 79fcc34107a7c652ece98227796eaf1f621dec5a6488dabdd300baa8018d9934 |
| report.md | e72f1fb18c0a454494d62d6d796ca44982c42f1ff26dc649e249612c5be5ff1c |
| corpus.json | 835c32573d3d6481f0665538cc6a295cc15f8e2d8df31144ef5b30d2a5e8d92d |
| evidence-manifest.json | a10cef4956cc000aca04251d0473fba61c7f17fa12feb2cd9e7d4481cea31beb |

- windows cli: 4af071ea68e6d475ae436b770c2bb0a5e12844bb784098b02836a92bdeef0587
- windows daemon: 4ca2bd766000e1622d4956d7c98d40a53e27d0e52b110dcaa8e1296a84fd9ded
- windows runner: a14b75958d6c0424ba026295ae8444341e1fb7806a3972d7d99ff35447ddf58d
- linux cli: 297b7b180226d361cfa47fb85a60f012b5be7e61b097aec5a3b689a6d1bd84ce
- linux daemon: bb5fa3dc6a302d35a551f035115174f627c9a3f81f290f1eaa8ab030068b0215
- linux runner: ab84efb8ad2dbd77cb622cb0d68757c17c1658cbc08b03114816bf797c836a11

The later #385 change touches only desktop/frontend/src/App.test.tsx. Subsequent factual Spec-Kit/verification finalization changes do not alter the evaluated CLI/daemon implementation. The frozen catalog is attributable to the digest above, not a claim that later documentation bytes existed during capture. Original input binaries, build receipt and corpus remain unchanged.

## Canonical aggregate

The first complete aggregate passed all eight gates in 1273.007 seconds before the final incident edge corrections. A subsequent source-stable aggregate passed format, vet, lint and race but failed the independent desktop search fixture (#385), with native exit 1 after 420.346 seconds. Its DOM showed an empty-query refusal, not an identity bypass. The failed log is retained with SHA-256 f61156e22ff2ee69e846ba37a88753cfd3a46f870b3861ee1b6a695270467cc8. The correction waits for authored route focus and asserts the exact intended query, preserving every original mismatch and no-Activity assertion. All 17 focused shell tests then passed.

The pre-hosted-CI corrected foreground aggregate passed all eight gates in 952.911 seconds (2026-10-27T15:32:19.943Z through 2026-10-27T15:48:12.854Z), native exit 0, 2401 frozen inputs, no drift and all 26 original artifacts unchanged. Core coverage: engine 83.6%, schedule 89.1%, timezone 91.3%, store 80.6%, catchup 88.9%, logbus 91.1%. All 167 frontend tests, production bundle/native desktop build, 23 canonical documentation pages, website static export and automation fixtures passed. The receipt and before/after input manifests remain private under .test-temp/s121/.

Final factual delivery/verification records are authored after that frozen aggregate. The repository GitHub formatter and automation/lifecycle checker both passed on the finalized records before committing. No application source changes follow the aggregate.

## Causal and publication boundaries

#355 original CLR/prerequisite causes and #356 original interrupted IPC/inspection cause remain unconfirmed and open. The source delivers incident-capture capabilities and separately corrects #385; passing verification is not original causal resolution or an issue closure gate. No installed service, account session, public DNS, tag, release asset or public promotion changes occurred. Source push, official PR and at most two Codex rounds are authorized; merge, tagging and promotion remain excluded.

## Hosted macOS fixture correction (#387)

The first hosted CI run (40535147345, macOS job 122298271500) failed because test-owned t.TempDir parents retained the normal macOS temporary-root alias. Production clifixture.NewWorkspace already canonicalizes its system temporary parent. Five test files now resolve only their test-owned temporary parent before invoking the unchanged strict collector. A POSIX alias regression reproduced the same rejection before the correction and passed afterward; caller-created parent links and artifact links remain rejected. The retained failed hosted log has SHA-256 4561953f7a85eef8adc3290ca88089df7e6d106f434c85a475f50fe48a26eaf9. This correction changes test setup only, with no CLI/daemon implementation or frozen catalog input change.

The automatic Code Review and owner-authorized second Codex request both completed on 1fa11413217ddfa6a896242df4bab12152989a13 with no findings. The second request returned a standard review receipt, not a dedicated security verdict. CodeQL passed independently on that head. The two-round allowance is exhausted; subsequent test-only CI corrections receive exact-head CI verification without a third request.

## Corrected-source acceptance

The post-macOS-correction foreground aggregate passed all eight gates in 1162.916 seconds (2026-10-27T16:04:35.649Z through 2026-10-27T16:23:58.565Z), native exit 0, 2403 frozen inputs, no drift and all 26 original artifacts unchanged. Native Windows package/race checks, all 167 frontend tests, production/native desktop build, all six core coverage thresholds, 23 canonical documentation pages, website static export and automation fixtures passed. The affected nonroot Linux race suite and POSIX alias regression passed separately. Final factual records follow the frozen aggregate; only documentation changes follow this verification.

## Merged-main verification

Owner merge 6f1de579bed866eb2bb328a25a0ad165caaa278a independently passed all 23 checks, including completed/success CI 40546444215, CodeQL 40546444065 and Public website 40546444134. Exact-source Pages deployment 7364985902 succeeded, and live landing/docs/comparison/incident pages and the legacy alias returned HTTP 200 over validated HTTPS. This is merged-main acceptance, separate from prior PR-head and frozen catalog receipts. See delivery.md and the public post-merge receipt for exact links.
