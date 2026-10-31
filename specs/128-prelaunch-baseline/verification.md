# S128 verification

The complete foreground `sh scripts/verify.sh all` passed on 2026-10-31T07:30:41.585Z in 1250.997 seconds. Source was the S128 working tree based on accepted main `e745b82a321d4efcc29e8d2ff1ef0b7b1f7e2ecf`. No tag or release acceptance is asserted.

## Complete aggregate

| Gate | Result |
| --- | --- |
| format | PASS, Go formatting and actual GitHub prose formatter |
| vet | PASS |
| lint | PASS, zero issues |
| race | PASS, complete root package/integration suite |
| gui | PASS, desktop race suite, native Wails application build, 167 component tests and frontend build |
| coverage | PASS, canonical cross-package measurement |
| docs | PASS, 24 canonical documents, policy/architecture checks and production export |
| automation | PASS, 22 job timeout alignment, workflow/lifecycle audits, coverage-profile and negative-case fixtures |

Core coverage: engine 83.6%, schedule 89.1%, timezone 91.3%, store 80.9%, catchup 88.9%, logbus 91.1%. The retained accepted S127 reading had store 80.6%; the other five readings are unchanged. The existing 80% thresholds remain.

## Focused evidence

- The accepted pre-S128 application database copy reopens with all 55 authored schema objects (22 tables, 33 indexes), current data and 24 historical marker rows retained. Fresh state writes one current marker 24 directly. Initialization fault injection rolls back partial DDL; unsupported state retains original bytes.
- Metrics format 4 is complete for both writers. Real payload/reader twins, constraints, concurrent writes, literal bindings, bounded sessions, atomicity and refusal checks pass. The deterministic 13-second exclusive-lock regression passes in both twins through the existing 5000ms/three-attempt retry. Non-SQLite refusal retains its bytes, original diagnostic and reset guidance.
- Shared preference, CLI/evaluator and desktop settings race suites pass; absent/current/malformed/unknown-field/no-import/no-transition behavior is covered. Current locking, scoped writes and atomic replacement remain.
- Clean desktop/site npm restores pass. Desktop browser: 35/35; site tests: 7/7; site browser: 10/10. New six site captures are isolated under the private S128 directory.
- Current qualification packages and Windows/V14 source contracts pass. Evidence revision 2 has 45 current observations; session manifest revision 1 retains three exact inputs. The independent current synthetic packet is not native evidence or a conversion of the archived revision-1 packet.
- ShellCheck passes for all test-script shell files/libraries. Both touched helper PowerShell files and all four touched Windows harnesses pass parsing/compliance checks. The existing library Write-Log advisory remains unchanged.
- The small initial three-iteration fresh-open sample measured 9,974,533 ns/op before simplification and 1,703,467 ns/op after it. This precedes the final rollback-error cleanup and is not an end-to-end or final-package performance claim.

## Preserved first failures and corrections

Private receipts under `.test-temp/s128/` retain baseline TDD failures and subsequent failures separately. The first aggregate stopped at unchecked rollback cleanup (fixed with explicit error handling). The second stopped at the date-sensitive desktop assertion (#409, fixed-clock future/equal/past and repeated-hour first/last cases). The third stopped at a lost concurrent metrics record (#410, bounded read-only retry). A separate malformed-file regression retained the missing-reset diagnostic failure before its correction. The first site restore exposed a relative specification link escaping the docs export; the absolute source link passes. No failed result was relabeled successful. Original #355/#356 causes remain unconfirmed.

## Source and evidence integrity

The successful aggregate retained 2,535 path records, including 19 deliberate retired-file deletions, and 2,516 present source files. Before/after drift: zero. Verified present-source manifest digest: `6f5d2ee7662486be2e90a1a985b6aef352bdf0bbe18d5b84f1c249af99737026`. All 209 frozen historical inputs retain their recorded sizes and SHA-256 values. Historical schema-1 evidence, release notes, assets and comparison bindings are unchanged.

After this aggregate, five factual specification/inventory records are finalized: this file, delivery.md, spec.md, tasks.md and specs/README.md. The staged whitespace check also removed one trailing blank line from current_schema.sql; the DDL text is identical after trimming trailing whitespace, and the full storage race suite is rechecked. Publication formatting, lifecycle and source/frozen-input guards verify that limited finalization before commit. It does not reassign prior QA to later bytes.

## Unverified boundaries

No software was installed on the active workstation. Native MSI execution, attended Windows observations, physical input/scaling and tagged-package acceptance are untested here. Build, automated source checks and headless browser results do not supply that evidence. S128 is unreleased source; published v1.7.0 is unchanged. External CI/review results are recorded on the eventual exact PR head.
