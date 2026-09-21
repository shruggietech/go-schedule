# Portable Bundle Contract v1

## Read operations

| Method | Path | Authority | Result |
|---|---|---|---|
| `GET` | `/v1/bundles/export` | Observe | Canonical portable document and exclusions |
| `POST` | `/v1/bundles/validate` | Observe | Schema and reference findings |
| `POST` | `/v1/bundles/preview` | Observe | Target-bound plan and drift |
| `POST` | `/v1/bundles/compare` | Observe | Read-only plan and drift |

## Apply operation

| Method | Path | Authority | Result |
|---|---|---|---|
| `POST` | `/v1/bundles/apply` | Manage | Per-item terminal outcomes |

## Apply request

```json
{
  "plan_id": "preview identity",
  "bundle_digest": "sha256 hex",
  "target_daemon_id": "selected daemon identity",
  "target_fingerprint": "preview fingerprint"
}
```

The daemon resolves the server-held, short-lived preview by `plan_id`, then rejects a mismatched daemon ID, digest, fingerprint, expired preview, or changed target before mutating the first item. Responses never contain protected source material.

## Error semantics

- Invalid document or unsupported schema: request rejected, no mutation.
- Conflict or incompatibility: represented per item in a preview, excluded from apply.
- Remote transport ambiguity: item outcome `uncertain`; no automatic retry.
- Omitted target state: `target_only` drift only, never deletion.
