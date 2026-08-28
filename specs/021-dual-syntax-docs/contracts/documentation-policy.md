# Contract: Dual-Syntax Documentation Policy

## Current-surface contract

Current product surfaces MUST:

1. describe recurring authoring as readable phrases or go-schedule's supported
   five-field cron subset;
2. preserve the message that cron knowledge is optional;
3. avoid promising arbitrary or standard cron parity;
4. identify named fidelity refusals and retained source identity where relevant;
5. distinguish cron expressions from crontab files and commands; and
6. state that the task timezone controls recurrence.

## Automated inventory

The policy helper scans only named current-facing files. It MUST not scan the
changelog or historical specifications globally. The canonical docs check runs
the helper and fixture harness, proving aligned copy passes and stale claims fail.

## Historical contract

S008 content remains unchanged except for prominent supersession notices on its
spec, tasks, quickstart, and fidelity checklist. Notices point to S019-S021.

## Publication contract

The eventual pull request uses `Closes #52`, `Closes #50`, and `Refs #22`.
