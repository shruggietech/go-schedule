# S117 execution quickstart

Use the verified hidden noninteractive launcher for non-Git tools. Inspect existing PowerShell/Go/SQLite and nonroot WSL identity and delegation first. Do not install prerequisites or alter workstation services.

From an unchanged candidate snapshot run:

~~~powershell
pwsh -NoProfile -NonInteractive -File scripts/cli-qa.ps1 -Action build -BinaryDirectory A:\Code\go-schedule\.test-temp\s117\qa-inputs-01
pwsh -NoProfile -NonInteractive -File scripts/cli-qa.ps1 -Action capture -BinaryDirectory A:\Code\go-schedule\.test-temp\s117\qa-inputs-01 -OutputDirectory A:\Code\go-schedule\.test-temp\s117\qa-run-01 -WSLDistribution Ubuntu-24.04 -WSLUser h8rt3rmin8r
~~~

Leave Budget omitted. Keep source frozen while actually reading the fresh corpus and occurrence artifacts and producing reviews-01.json. Run finalize with that review file and the same input/output directories, then validate and render. Preserve native exit 0 pass, 2 valid fail, 3 valid blocked; 1 is processing/integrity error. Do not edit generated report.md.

After report finalization, inspect each actual failure, reproduce narrowly if needed, write a failing regression before its correction, and retain original versus corrected acceptance separately. Then run `sh scripts/verify.sh all` in the foreground through the verified hidden environment, source/comparison guards and GitHub formatter. Commit only after all canonical gates pass. Push/open the authorized PR, monitor exact-head CI, address individual feedback within two Codex rounds and stop before merge/tag/staging/promotion.
