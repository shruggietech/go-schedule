# Release-source preflight contract

Command: `go run ./scripts/release-source-check -root . -repository shruggietech/go-schedule -tag v1.6.0`.

Read local files only. Exit 0 after all source metadata agrees; nonzero with actionable stderr on invalid arguments, missing/unreadable files or inconsistent metadata. Success is source validation, not tag, package, CI or publication verification.

Require one Release badge occurrence and one canonical latest-release badge line (surrounding whitespace permitted), one health example matching the plain version, one exact version heading with calendar-valid YYYY-MM-DD, tag-specific notes and exactly one whole-line canonical changelog link derived from that heading. CRLF behavior follows ordinary line parsing; extra prose around the required link is rejected. Reject malformed tags/repositories before path construction.

Release invokes the shared command after setup-go and before comparison; failures propagate and cannot be conditional or continue-on-error. Existing exact-main CI, draft-only staging, candidate identity, checksum, allowlist and promotion contracts remain unchanged.
