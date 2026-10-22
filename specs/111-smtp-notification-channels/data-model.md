# S111 data model

## Channel

NotificationChannel retains ID/name/kind/enablement/health interval/timestamps, protected endpoint/authorization and presence metadata. Add optional SMTPConfig with host, port, TLS mode, sender, recipients, timeout seconds and HTML selection. SMTPAuth username/password are a write-only input, serialized into the existing protected authorization slot. Kind defaults to webhook for old requests and is immutable after creation.

Validation: host is a bounded DNS name/IP, port 1-65535, mode starttls/tls/plain, sender/recipient addresses are single ASCII-envelope mailboxes, at most 100 unique recipients, no CR/LF/NUL, timeout 1-60 seconds (default 10), bounded field lengths and message. Unicode event names are MIME encoded. Clear credentials explicitly; omission retains them.

## Delivery

Add kind and frozen typed settings to each queued delivery; existing protected authorization snapshot carries SMTP auth and is never serialized. Add bounded recipient outcomes: address, state, SMTP numeric code and authored diagnostic. States: pending, submitting, accepted, temporary_failure, permanent_failure, uncertain. Recipient list is frozen with the event. Only accepted after final DATA success; accepted/permanent recipients are omitted from future attempts.

The existing pending/claimed/succeeded/failed aggregate lifecycle and three attempts remain. Submission checkpoint records the eligible DATA subset as submitting. Final SMTP result persists recipient states with aggregate retry/completion in one transaction. Recovery retains accepted/permanent states and changes interrupted submitting states to uncertain, preserving duplicate-risk evidence. Completed/exhausted rows erase credential snapshots.

## Portable SMTP definition

A unique logical channel name identifies nonsecret typed settings plus enablement/health interval and whether destination credentials are required. Provision exactly one target-local SMTP channel with that name and matching credential presence before preview. Stored credentials require the same host, port and encryption mode. Name-based binding is resolved during preview and checked again during apply with ID, UpdatedAt and target fingerprint. Export, comparison and fingerprints contain no credential values. Missing binding blocks apply before writes. Existing webhook-only v2 documents retain their shape and behavior.

## Migration

Forward v22-to-v23 rebuild uses the existing foreignKeysOff facility and post-migration foreign-key check. Copy exact prior IDs/protected strings and restore the enabled index; do not rename the old parent first. Add delivery columns with webhook-compatible defaults. Retain task/group assignments, condition state, daemon-health state and queued/history records.
