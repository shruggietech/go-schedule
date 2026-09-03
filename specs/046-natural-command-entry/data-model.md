# Data Model: Natural Command Entry

S046 changes the authoring projection only. The stored task schema remains unchanged.

## Command Draft

Transient text owned by the task editor.

| Field | Meaning | Validation |
| --- | --- | --- |
| `text` | User-authored portable command-line text | May be incomplete while editing; must be valid UTF-8; NUL is invalid |
| `state` | empty, invalid, or valid | Derived on each edit refresh |
| `error` | Plain-language syntax failure | Present only when invalid; identifies position/line |

State transitions:

```text
empty -> invalid or valid as text is entered
invalid -> empty or valid after correction
valid -> empty or invalid or another valid invocation after editing
```

No invalid or empty draft can be submitted.

## Direct Invocation

The authoritative value already represented by `domain.Task` and API task requests.

| Field | Meaning | Validation |
| --- | --- | --- |
| `program` | Exact executable name or path | Non-empty valid UTF-8; no NUL |
| `arguments` | Ordered argument values | Zero or more valid UTF-8 values; empty values allowed; no NUL |

Identity is the exact program string plus the length, order, and exact text of every argument. Cosmetic quotation spelling is not part of identity.

## Canonical Command Line

A transient, lossless representation generated from a Direct Invocation when an existing task opens.

Invariants:

1. Parsing the canonical text returns the identical Direct Invocation.
2. Formatting that result is stable.
3. The same text parses identically on every supported platform.
4. No canonical text is persisted separately.

## Launch Preview

A display projection derived only from a valid Direct Invocation.

| Field | Meaning |
| --- | --- |
| `program_display` | Quoted escaped exact program value |
| `argument_displays` | Numbered quoted escaped values in original order |
| `argument_count` | Explicit count, including empty values |
| `guidance` | Empty-state or syntax-error guidance when no valid invocation exists |

The preview is never an execution input.

## Existing Task Compatibility

```text
stored command + args
        |
        v
canonical formatter -> editor text -> portable parser
        |                                  |
        +---------- identity check <-------+
                                           |
                                           v
                                existing API/store/executor
```

No table, column, JSON field, or migration is added.
