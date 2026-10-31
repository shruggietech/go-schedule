---
title: First task and recovery
nav_order: 2.5
---

# Your first useful task

Create a local text-file snapshot, inspect its exact run, then diagnose and repair a missing-input failure. This example supplies its own script and input. The task needs only the installed package and native shell/file tools, with no network account or additional interpreter.

## 1. Connect to the intended local scheduler

Install using the [Windows](INSTALL-windows.md), [Linux](INSTALL-linux.md) or [macOS](INSTALL-macos.md) guide. On Windows, open a new terminal after installation to inherit the updated PATH.

~~~sh
gosched health
gosched daemon manifest
~~~

Health must identify a reachable daemon; the manifest identifies its platform and installation. These instructions assume the default local target, without remote profile/endpoint options. For a custom local socket or named pipe, add the matching global --ipc-path option to every gosched command; see the [CLI target reference](cli.md). In the desktop, select **This computer** and confirm it is connected.

If health fails, use gosched service status and the platform guide's troubleshooting section. A stopped service needs an explicit Start action; access denied needs the documented account/group access. Do not erase data, start a second default daemon or disable access checks to solve a connection problem.

All task paths belong to the daemon's computer. Windows's installed service runs as SYSTEM; system services on Linux/macOS usually run as root, while a manually started daemon uses its process identity. Both you and that identity need read/write access to the example directory and traversal access through its parents. This local example contains no credentials.

## 2. Supply the script and input

Use a new workspace. Setup refuses an existing directory so it cannot overwrite your files. A successful run intentionally replaces only snapshot.txt in that workspace; missing input is checked before writing. This is a small file-copy example, not a complete backup or versioning system.

### Windows (PowerShell)

The default ProgramData location is visible to the installed service. Keep the quoted working-directory value together even though it contains spaces.

~~~powershell
$Workspace = Join-Path $env:ProgramData 'go-schedule first task'
if (Test-Path -LiteralPath $Workspace) { throw 'Workspace already exists. Choose a new directory before continuing.' }
New-Item -ItemType Directory -Path $Workspace | Out-Null
$Utf8 = New-Object System.Text.UTF8Encoding($false)
[IO.File]::WriteAllText((Join-Path $Workspace 'source.txt'), ('A useful local note.' + [Environment]::NewLine), $Utf8)
$Script = @'
@echo off
setlocal
if not exist "source.txt" (
  echo Missing source.txt. Restore it and run this task again. 1>&2
  exit /b 2
)
copy /y "source.txt" "snapshot.txt" >nul
if errorlevel 1 exit /b 1
echo Snapshot saved to snapshot.txt.
'@
[IO.File]::WriteAllText((Join-Path $Workspace 'snapshot.cmd'), $Script, $Utf8)
~~~

Name the shell explicitly: go-schedule does not interpret shell built-ins or redirects by itself.

~~~powershell
gosched task add "First snapshot" --command "$env:SystemRoot\System32\cmd.exe" --arg /d --arg /c --arg snapshot.cmd --cwd "$Workspace" --schedule "every day at 09:00" --tz UTC
~~~

The response includes the new task ID, readiness and upcoming runs. This complete CLI add creates an active schedule at **09:00 UTC**, not 09:00 in your local zone. Copy the returned task ID below and disable automatic runs during the exercise. Run now remains available for this command-ready task.

~~~powershell
$TaskId = 'paste-the-created-task-id'
gosched task disable $TaskId
~~~

### Linux or macOS (POSIX shell)

Choose a local directory the daemon identity can access. If a restricted service identity cannot traverse your home directory, choose another shared local workspace with the required ownership/permissions; do not make it world-writable.

~~~sh
WORKSPACE="$HOME/go-schedule first task"
mkdir "$WORKSPACE" && (
  cd "$WORKSPACE" || exit 1
  printf '%s\n' 'A useful local note.' > source.txt
  cat > snapshot.sh <<'SCRIPT'
#!/bin/sh
set -eu
if [ ! -f source.txt ]; then
  printf '%s\n' 'Missing source.txt. Restore it and run this task again.' >&2
  exit 2
fi
cp source.txt snapshot.txt
printf '%s\n' 'Snapshot saved to snapshot.txt.'
SCRIPT
)
~~~

Continue only if setup succeeded. mkdir refuses an existing workspace; do not silently reuse it. /bin/sh reads the script, so it needs no executable permission bit.

~~~sh
gosched task add "First snapshot" --command /bin/sh --arg snapshot.sh --cwd "$WORKSPACE" --schedule "every day at 09:00" --tz UTC
~~~

Copy the returned task ID below. The new schedule is active at **09:00 UTC**; disable automatic runs so the deliberate failure does not race scheduled work.

~~~sh
TASK_ID='paste-the-created-task-id'
gosched task disable "$TASK_ID"
~~~

## 3. Run and inspect the result

Windows:

~~~powershell
gosched task run-now $TaskId
gosched runs --task $TaskId --details
~~~

Linux/macOS:

~~~sh
gosched task run-now "$TASK_ID"
gosched runs --task "$TASK_ID" --details
~~~

The message triggered run for task acknowledges submission. It does **not** say execution finished and contains no run ID. Find the new record for this task and its scheduled/start time. Refresh history if it has not appeared yet; gosched runs active separately shows executions still running. Queued/running work is not successful completion, and an older successful record is not evidence for the new request.

A completed successful record has its own **Run ID**, a populated **Ended** timestamp, **State: success**, **Exit code: 0**, and captured output Snapshot saved to snapshot.txt. Substitute that record's Run ID for RUN_ID below, without angle brackets:

~~~sh
gosched runs get RUN_ID
~~~

The --json flag is available for automation. Plain gosched runs is a compact table; --details includes IDs and output. On Windows, read the snapshot with Get-Content -LiteralPath (Join-Path $Workspace 'snapshot.txt'); on Linux/macOS, use cat "$WORKSPACE/snapshot.txt". Its contents should match source.txt. File existence alone is not a run result: a failed attempt leaves the prior snapshot in place.

## 4. Explain a failure and recover

Rename only the supplied input, then request another run of the **same task**.

Windows:

~~~powershell
Move-Item -LiteralPath (Join-Path $Workspace 'source.txt') -Destination (Join-Path $Workspace 'source.saved.txt')
gosched task run-now $TaskId
gosched runs --task $TaskId --details
~~~

Linux/macOS:

~~~sh
mv "$WORKSPACE/source.txt" "$WORKSPACE/source.saved.txt"
gosched task run-now "$TASK_ID"
gosched runs --task "$TASK_ID" --details
~~~

Wait for the new terminal record, then use its own Run ID with gosched runs get. Expect **State: failure**, **Exit code: 2**, and Missing source.txt. Restore it and run this task again. in captured output. The old snapshot should still exist; it does not establish this attempt's success. If you see Start failed instead, inspect the executable/working-directory path and permissions: the payload did not reach its missing-input check.

Restore the original input and request a third run.

Windows:

~~~powershell
Move-Item -LiteralPath (Join-Path $Workspace 'source.saved.txt') -Destination (Join-Path $Workspace 'source.txt')
gosched task run-now $TaskId
gosched runs --task $TaskId --details
~~~

Linux/macOS:

~~~sh
mv "$WORKSPACE/source.saved.txt" "$WORKSPACE/source.txt"
gosched task run-now "$TASK_ID"
gosched runs --task "$TASK_ID" --details
~~~

The third terminal record should succeed with exit code 0. One task retains success, explained failure and successful recovery as separate run IDs. You did not delete the task, reset data or erase the failure.

## 5. Use the desktop for the same workflow

Open the installed desktop shortcut or gosched gui. In **Tasks**, select **Create task**. Use the supplied script and workspace:

| Daemon platform | Command line | Advanced settings: Working directory |
| --- | --- | --- |
| Windows | cmd.exe /d /c snapshot.cmd | Complete value of $Workspace |
| Linux/macOS | /bin/sh snapshot.sh | Complete value of $WORKSPACE |

Select **Timing mode: Recurring**, enter every day at 09:00 in **Schedule**, and set **Timezone: UTC**. Select **Preview** for command boundaries and upcoming times, then **Save inactive task**. Unlike the complete CLI add above, a new desktop task is saved inactive; leave it that way during the exercise. If you already created the CLI task, select it instead of creating a duplicate.

Choose **Run now**, then **Confirm run**. Acceptance opens **Activity**, but work may still be in progress. Set **Record type: Runs**, use **Search** with the task ID, select the particular run's time button, and inspect **State**, **Exit status** and **Combined stdout/stderr** in its **Run details** side pane. Use **Refresh** for a new snapshot. **Clear View** hides displayed records and acknowledges visible alerts; it does not delete history or repair a task.

Builds containing S132 also display the selected run's own **Run ID**, separately from **Source run** (an upstream completion-chain run). Published v2.0.0 does not yet show the own-ID field; use CLI --details history to find it until a package containing that source change is published. See the [desktop reference](gui-fields.md) for all fields.

## 6. Keep or retire the example

The schedule remains disabled after the exercise. When ready, gosched task enable followed by the task ID enables its 09:00 UTC schedule; the desktop's **Enable** action does the same for a command-ready task. Review the preview and timezone first.

To leave it inactive, keep it disabled. If you no longer want it, save any desired run evidence and remove only this task with gosched task rm followed by its ID. Then remove only the example workspace you deliberately created, after confirming its exact path and that no task uses it. This guide does not delete daemon data or any other directory. Continue with the [CLI reference](cli.md), [desktop fields](gui-fields.md) or [notifications](notifications.md).
