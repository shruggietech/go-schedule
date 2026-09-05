# Quickstart: Validate GUI Navigation and Information

## Focused headless regression

From the repository root, run:

```powershell
go test ./gui -run 'TestUI_(BuildsAllTabs|ActivityBadgeReflectsUnacked|Info)' -count=1
```

Expected:

- the tab sequence is Tasks, Groups, Schedule, Activity, Info;
- Activity badge updates leave Info last;
- Info uses the embedded full mark and exact build version;
- all three descriptive labels map to their canonical HTTPS destinations.

## Manual desktop observation

Run the GUI against the normal local daemon and inspect the leading navigation. Select Info and confirm the hierarchy remains readable at the default window size and after making the window smaller. The mark should retain its proportions, and the version, maintainer attribution, repository, and documentation should be readable without contacting the daemon.

External-link activation may be observed with the machine's default browser, but the automated contract validates destinations without launching a browser.

## Repository verification

```sh
sh scripts/verify.sh all
```

All eight gates must pass before the local feature commit. The PR description will use `Closes #29` and `Closes #32` because this slice completes both issues.
