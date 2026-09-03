# Quickstart: Validate Natural Command Entry

## Prerequisites

- Go toolchain selected from `go.mod`
- C compiler available where the race detector requires it
- Repository checkout on `codex/046-natural-command-entry`
- Windows, macOS, or Linux for the corresponding native process-boundary test

## 1. Parser and formatter verification

```bash
go test ./internal/commandline -count=1
go test -race ./internal/commandline -count=1
```

Expected: the contract in [contracts/portable-command-line.md](contracts/portable-command-line.md) passes for spaces, both quote styles, empty values, Unicode, Windows/POSIX paths, repeated flags, backslashes, tabs, CR/LF, shell punctuation, invalid quotes, canonical stability, and seeded fuzz round trips.

## 2. Headless editor verification

```bash
go test ./gui -run 'TestEditor_(Command|ExistingCommand|CommandHelp|CommandSize)' -count=1
```

Expected: one Command line field exposes at least six lines, grows vertically, updates an exact program/argument preview, clears stale previews on invalid input, gates Save, submits exact create/update values, and retains existing tasks without mutation.

## 3. Compatibility and native launch verification

```bash
go test ./internal/api/server ./internal/store ./internal/executor ./internal/commandline -count=1
```

Expected: the existing API and SQLite shapes retain unusual values and order, the executor invokes one direct program, and the host-specific helper process receives the expected argument vector. Windows execution remains console-hidden.

## 4. Full repository verification

```bash
go test ./... -count=1
sh scripts/verify.sh all
```

Expected: the full suite and all eight canonical gates pass in order: format, vet, lint, race, gui, coverage, docs, automation.

## 5. Headless/manual editor matrix

Create and edit tasks from these representative inputs:

```text
python -m http.server --bind "127.0.0.1"
"C:\Program Files\Tool\tool.exe" --name "Ada Lovelace" --empty ''
/usr/bin/printf '%s\n' 'héllo 世界'
program --tag one --tag two '$HOME' '|' '*.txt'
```

Also enter one double-quoted value containing a literal newline. For every case:

1. Read the separately labeled program and numbered arguments.
2. Confirm empty and multiline values are visibly escaped.
3. Save, reopen, and compare the preview.
4. Enter an unmatched quote and confirm Save disables and the old preview disappears.
5. Confirm the editor shows at least six lines at default size and gains height when the dialog grows.

## 6. Shell boundary

Enter shell punctuation without a shell and confirm it remains literal. Then explicitly name `cmd /c`, PowerShell, or `sh -c` only in controlled tests and confirm the help accurately assigns interpretation to that shell. S046 must introduce no implicit shell process.
