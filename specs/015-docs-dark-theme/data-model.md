# Data Model: Documentation Dark-Theme Quality

This feature introduces no persisted application data. Its validation model has
two static concepts.

## Token role

- **Safe base**: default foreground for every classified or unclassified code
  token.
- **Muted**: comments and prompts.
- **Blue**: keywords, functions, types, names, decorators, and headings.
- **Mint**: strings, tags, and inserted diff content.
- **Amber**: numbers, literals, variables, and attributes.
- **Base**: operators, punctuation, whitespace, and lexer errors.
- **Red**: deleted diff content.
- **Highlighted line**: dark secondary code surface.
- **Selection**: Anchor Blue background with Panel foreground.

Every token first receives Safe base and may then receive exactly one named
foreground role. Styling-only roles such as bold and italic do not change the
foreground contract.

## Fence category

| Category | Meaning |
| --- | --- |
| `sh` | Portable POSIX shell command or snippet |
| `bash` | Bash-specific script or behavior |
| `powershell` | PowerShell command or script |
| `text` | Output, grammar, or non-executable example |

An opening fence without one of these values is invalid. Closing fences carry
no category.
