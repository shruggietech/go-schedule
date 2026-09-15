# Session data model

Version-1 input manifests bind repository, tag, commit, staging run ID/attempt, candidate MSI, public baseline MSI, portable PowerShell ZIP, and offline WebView2 EXE. Inputs have absolute path, positive bytes, lowercase SHA-256, and source URL. Helpers are hashed when packaged.

Generated scenario records name fresh or upgrade and fixed packaged filenames. Phase results retain UTC start/end, elapsed seconds, exit code when available, completed/failed/timed-out status, and log location. No generated record is an attended pass.
