# Installing go-schedule on Windows

go-schedule installs as a **formal Windows application** via an `.msi` package. It puts the program
in *Program Files*, runs the scheduler as an auto-starting **Windows service**, and adds a
Start-Menu shortcut. There is no "extract a zip and run an exe from Downloads" step.

## Install

1. From the [latest release](https://github.com/shruggietech/go-schedule/releases/latest),
   download **`go-schedule_<ver>_windows_amd64.msi`**.

2. (Recommended) Verify the download against `SHA256SUMS.txt`:

   ```powershell
   Get-FileHash .\go-schedule_*_windows_amd64.msi -Algorithm SHA256
   ```

3. **Double-click the `.msi`** and complete the wizard. Windows will prompt for administrator
   approval (UAC) — this is required because the installer registers a system service. Approve it.

That's it. The installer:

- installs `gosched-gui.exe`, `goschedd.exe`, and `gosched.exe` to
  `C:\Program Files\go-schedule\`;
- registers **`goschedd`** as a Windows service set to **start automatically** (so your tasks run
  in the background and survive reboots, even with no one logged in);
- adds a **go-schedule** shortcut to the Start Menu.

Launch **go-schedule** from the Start Menu to open the desktop app. It connects to the already-
running service (it does not start a second copy) and never shows a console window.

> Data location: tasks and logs live under `C:\ProgramData\goschedule\` (the database
> `goschedule.db` and the `logs\` folder). This is created automatically on first run.

## Upgrading

Download the newer `.msi` and run it — it performs an in-place major upgrade (the old version is
removed and the new one installed). Your data under `C:\ProgramData\goschedule\` is preserved.

## Using the CLI (optional)

The CLI is installed alongside the app. Open PowerShell:

```powershell
& "C:\Program Files\go-schedule\gosched.exe" health
& "C:\Program Files\go-schedule\gosched.exe" service status   # expect: running

# A recurring task
& "C:\Program Files\go-schedule\gosched.exe" task add backup `
  --command "C:\Windows\System32\cmd.exe" --arg "/c" --arg "echo backup" `
  --schedule "every weekday at 09:00"

& "C:\Program Files\go-schedule\gosched.exe" task list
& "C:\Program Files\go-schedule\gosched.exe" logs --severity error   # recent error logs
```

(Add `C:\Program Files\go-schedule\` to your `PATH` to drop the full path.)

## Uninstalling

Use **Settings → Apps → Installed apps → go-schedule → Uninstall** (or *Apps & features*). This
stops and removes the service, deletes the program files, and removes the Start-Menu shortcut.

Your data under `C:\ProgramData\goschedule\` is **left in place** so a later reinstall keeps your
tasks. To remove it completely, delete that folder manually after uninstalling:

```powershell
Remove-Item -Recurse -Force "C:\ProgramData\goschedule"
```

## Troubleshooting

- **UAC / "Do you want to allow this app to make changes?"** → expected; the service install needs
  elevation. If you decline, the install cancels cleanly (nothing is partially installed).
- **SmartScreen / antivirus warning** → the installer is currently unsigned. Verify the SHA-256
  hash against `SHA256SUMS.txt` and choose *More info → Run anyway* if it matches.
- **The GUI opens but says "daemon unreachable"** → check the service:
  `& "C:\Program Files\go-schedule\gosched.exe" service status`. Start it with
  `service start` if needed (run PowerShell as Administrator).
- **Where are the logs?** → `C:\ProgramData\goschedule\logs\goschedule.log` (and rotated
  siblings), or open the **Logs** view in the app.
