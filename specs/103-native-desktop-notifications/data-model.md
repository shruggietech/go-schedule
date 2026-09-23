# Data model: Native desktop notifications

## DesktopNotificationPreferences

Desktop-local versioned preference with `enabled` defaulting to false, enabled conditions, severity threshold or selection, and optional daemon profile IDs. Empty daemon selection means all registered daemons. Validation normalizes unknown conditions and missing profiles without broadening an explicitly narrowed selection. Restore defaults disables popups.

## NotificationCandidate

Transient record containing pinned daemon ID, profile ID or local target, record kind (`run` or `alert`), exact record ID, condition or severity, presentation title and body, and observation time. No credential, command output, webhook endpoint, or secret is copied into a popup. The dedupe key is daemon ID plus record kind plus record ID.

## ActivationIntent

Transient identity-safe routing object with target profile ID, expected daemon ID, record kind, and exact record ID. The frontend must verify the selected target's actual daemon identity before opening Activity. Missing or stale records lead to a non-destructive explanatory state.

## ReleaseBoundary

The v1.5.0 draft lists only merged features and known limitations. It expressly records SMTP as deferred, with #176 and parent #19 still open. It contains no publication timestamp, tag, download URL, or shipped claim before release.
