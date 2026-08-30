# Data Model: Repository Brand System Integration

This slice adds no runtime domain data. Its model is a version-controlled artifact contract.

## CanonicalArtifact

| Field | Meaning | Validation |
| --- | --- | --- |
| `path` | Path relative to `brand/` | Normalized, unique, cannot escape `brand/` |
| `bytes` | Approved file length | Non-negative integer matching disk |
| `sha256` | Approved content digest | Lowercase 64-character hexadecimal digest matching disk |
| `kind` | Inferred artifact family | Text, SVG, raster, font, PDF, platform, or reference |

The standalone `brand/manifest.json` contains the approved canonical artifacts. `manifest.json` and `VERIFY.md` are self-describing control files and are validated separately rather than hashing themselves.

## ConsumerMapping

| Field | Meaning | Validation |
| --- | --- | --- |
| `source` | Canonical path relative to `brand/` | Must exist and remain inside `brand/` |
| `targets` | One or more repository-relative consumer paths | Each must exist, remain inside the repository, and be byte-identical to `source` |
| `purpose` | Human-readable reason for the copy | Required and non-empty |

### Relationships

- One `CanonicalArtifact` can feed multiple `ConsumerMapping` targets.
- Every target belongs to exactly one mapping.
- A canonical artifact need not have a consumer copy when tooling can consume it directly.

## PortableSVG

A canonical or mapped SVG with these invariants:

- valid UTF-8 without BOM;
- one accessible `<title>` element;
- no `<text>` element;
- no `font-family` element attribute or style declaration.

## BrandIntegrityResult

| Field | Meaning |
| --- | --- |
| `checked_artifacts` | Number of manifest records examined |
| `checked_svgs` | Number of SVG portability checks |
| `checked_consumers` | Number of target comparisons |
| `failures` | Actionable messages naming the violated source or target |

The process exits zero only when `failures` is empty. Findings accumulate so one run reports all known drift.
