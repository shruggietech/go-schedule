# Contract: Exact-Candidate Issue Disposition Packet

## Invocation

```text
windows-release-gate render-dispositions
  --bundle <formal-evidence.zip>
  --candidate-manifest <windows-candidate-manifest.json>
  --artifact <exact-candidate.msi>
  --repository <owner/repository>
  --tag <vMAJOR.MINOR.PATCH>
  --commit <40-lowercase-hex>
  --output-dir <absent-directory>
```

All options are mandatory. Positional arguments are rejected. The operation is
offline and performs no GitHub mutation.

## Exit and stream contract

| Exit | Meaning | Output |
| --- | --- | --- |
| 0 | Complete packet atomically committed | One concise success line on stdout |
| 1 | Candidate or evidence failed production validation | Every independently discoverable validation failure on stderr |
| 2 | Invocation, input loading, rendering, or destination failure | One contextual diagnostic on stderr |

No failure may leave a requested destination created by this operation or a
sibling staging directory. A destination created concurrently by another
process is preserved and causes the operation to fail.

## Fixed output inventory

Exactly these files are committed:

```text
packet.json
issue-096.md
issue-098.md
issue-101.md
issue-104.md
issue-105.md
issue-106.md
issue-109.md
issue-111.md
issue-112.md
issue-113.md
```

`packet.json` uses UTF-8 without BOM, a trailing newline, and this logical
shape:

```json
{
  "schema_version": 1,
  "candidate": {
    "repository": "shruggietech/go-schedule",
    "tag": "v1.0.0",
    "commit": "<40-lowercase-hex>",
    "workflow": "Release",
    "run_id": 123,
    "run_attempt": 1,
    "filename": "go-schedule_v1.0.0_windows_amd64.msi",
    "bytes": 123456,
    "sha256": "<64-lowercase-hex>",
    "product_version": "1.0.0",
    "product_code": "{UPPERCASE-GUID}"
  },
  "issues": [
    {
      "issue": 96,
      "file": "issue-096.md",
      "observations": ["access.intended-user"]
    }
  ]
}
```

The example observation list is abbreviated; production #96 contains all 36
pre-desktop observations. The complete evidence archive contains all 47.

## Canonical mappings

The exact mappings are the table in `spec.md`. #96 contains the first 36
pre-desktop values from `releasegate.RequiredScenarioIDs()` in canonical order;
#98 contains the 16 setup/removal values from that same canonical list. Packet
files are ordered by issue number, and observations retain canonical order.

## Markdown contract

Each record includes:

1. a fixed title identifying the issue;
2. a candidate-identity table;
3. exact tag-specific workflow-run and evidence-archive links;
4. an explicit statement that production candidate, archive, attachment, and
   manifest validation passed;
5. a sorted table of environments referenced by the issue's observations;
6. one ordered section per mapped observation with status, timing, summary,
   deterministic metrics JSON, and archive-relative attachment paths; and
7. a fixed boundary statement explaining that the record supports, but does not
   perform, issue closure.

Evidence-controlled table text escapes Markdown delimiters, HTML delimiters,
newlines, code delimiters, and `@` mentions. Attachment paths are included only
after existing safe-relative-path and archive-integrity validation.

## Validation contract

Before rendering, the command must run the same production chain used by
promotion:

1. strict evidence and candidate-manifest decoding;
2. physical MSI existence and exact identity validation;
3. formal evidence root, environment, attachment, metric, status, and all-47
   observation validation;
4. exact archive membership validation; and
5. independent candidate-manifest equality validation.

The mandatory expected repository, tag, and commit options prevent a
self-consistent but wrong candidate from producing records.

## Atomicity contract

- The requested target must not exist.
- The parent must be a real directory and not a symbolic link.
- All content is rendered before filesystem mutation.
- Files are written with restrictive permissions in a private sibling staging
  directory.
- The staging directory is renamed onto the target only after every write and
  close succeeds.
- Cleanup removes the staging directory on every unsuccessful path.
- Existing targets are never deleted, merged, or overwritten.

## Non-authority contract

Generation proves only that the supplied exact-candidate evidence passes the
production gate and has been rendered without omission. It does not prove that
an issue's prose acceptance criteria were reviewed, and it never comments,
closes, tags, uploads, dispatches, promotes, or changes a milestone.
