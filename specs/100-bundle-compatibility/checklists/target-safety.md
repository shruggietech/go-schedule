# Target Safety Checklist

- [x] Requirements come from bundle content, not user-supplied capability claims.
- [x] Unsupported or unknown watcher platforms fail before apply.
- [x] Selection is frozen across manifest discovery and bundle operation.
- [x] Manifest failure does not forward preview or apply.
- [x] Existing daemon plan identity and fingerprint checks remain intact.
- [x] Target-local watcher paths stay out of canonical bundle bytes.
- [x] Read-only drift and no-inferred-removal semantics remain intact.
