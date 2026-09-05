# Verification Update: GUI Editor Refinements

## 2026-08-31 Intentional Reversal

Specification 037 supersedes specification 003 FR-001, User Story 5, SC-001, and the former maximize assumption. First-run Windows observation showed that the old full-work-area calculation was visually indistinguishable from a maximized window and could enlarge a constrained display beyond its work area.

The replacement contract requests a restored 1280x800 logical window, capped independently to 90 percent of the selected monitor's logical work area. The work area remains a hard upper bound, while explicit user maximize and restore behavior is unchanged.

Evidence for the correction is maintained in `specs/037-windows-release-polish/verification.md`. The original delivery evidence for the remaining specification 003 behaviors stays commit `a21c22e`.
