# Quickstart: Scheduler Startup Trigger

## Local string behavior

Run `gosched cron explain "@reboot"`, convert `@reboot`, and convert `at scheduler startup`. Explain should report `At scheduler startup` without upcoming times. Each conversion should return the opposite canonical syntax.

## Author a startup task

Run `gosched task add startup-check --command <program> --schedule "@reboot"`, then show the task. Detail should report `At scheduler startup`, no next-runs section, and active state. The desktop Schedule field accepts either canonical syntax and previews the same summary.

## Import fidelity

Use a crontab fixture with environment, `SHELL`, and one `@reboot` job. Dry-run must report one creatable job and mutate nothing. Real import must create one startup task, duplicate re-import must follow existing policy, and faithful export must begin with `@reboot`.

## Controlled lifecycle validation

Persist enabled, disabled, and disabled-group startup tasks. Start an engine with injected time, wait for its run callback, request reload, and confirm no additional run. Cancel and drain it, start a fresh engine on the same store, and confirm the eligible task has exactly two startup runs total while ineligible tasks have none.

## Verification

Run `sh scripts/verify.sh all`. Format, vet, lint, race, GUI, coverage, docs, and automation must all pass.

