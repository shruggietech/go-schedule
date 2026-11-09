# Validation guide

1. Run go test ./scripts/desktop-build, including malformed versions, development defaults, restoration after failure and concurrent-build refusal.
2. Build a selected Windows desktop through go run ./scripts/desktop-build -version v2.0.3 -- -clean -platform windows/amd64. This is a private metadata fixture, not a release selection. Inspect with build/windows/verify_version.ps1 without launching it.
3. Compile an MSI from the same selected version and inspect matching ProductVersion and GUI File.Version. Deliberately request a different version and require refusal.
4. Run scripts/verify.sh all in the existing hidden foreground driver and watch all eight gates. Preserve historical evidence.
