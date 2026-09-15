# Qualification contract

Reuse the existing `scripts/windows-qualification-session` manifest and sessions, `test/windows/Invoke-ReleaseCandidateAttended.ps1` collector, and `scripts/windows-release-gate` validator without schema or gate changes.

Candidate sources must identify the successful refresh at reviewed commit `951d864d9683ec3bdcb1e37535ecddd67738bf13`. Upgrade baseline is the independent public v1.1.1 MSI with its expected identity from `test/windows/README.md`. All native scenarios and attachment bytes must validate. Logs alone remain unavailable as attended evidence. S087 does not dispatch promotion.
