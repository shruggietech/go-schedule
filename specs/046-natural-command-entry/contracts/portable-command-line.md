# Contract: Portable Direct Command Line

## Purpose

Define the user-facing text grammar that maps one task-editor field to the existing direct invocation of one program plus ordered arguments. This is not a shell language.

## Lexical contract

1. Unicode whitespace separates values outside quotes.
2. Single (`'`) and double (`"`) quotation marks group content and are removed from the resulting value.
3. Empty quoted content creates an intentional empty value.
4. Adjacent quoted and unquoted segments without separating whitespace form one value.
5. Inside single quotes, every character except the closing single quote is literal.
6. Inside double quotes, backslash escapes a following double quote. Other backslashes are literal.
7. Outside quotes, backslash escapes following whitespace or a quotation mark. Before any other character, including another backslash or end of input, it is literal.
8. A newline outside quotes is whitespace. A newline inside quotes is part of the value.
9. Input and stored values must be valid UTF-8 text. Invalid UTF-8 is rejected instead of being silently replaced.
10. NUL is invalid because supported process APIs cannot launch it as program or argument content.
10. Unmatched quotes are invalid and report their opening position and line.

## Non-shell contract

The parser does not interpret variables, `%...%`, `$...`, wildcards, comments, semicolons, ampersands, pipes, redirects, command substitutions, or shell built-ins. Those characters remain ordinary value content.

A user can request shell behavior only by naming a shell as the program, for example:

```text
cmd /c "echo hello > output.txt"
sh -c 'echo hello > output.txt'
```

In those cases go-schedule still launches one direct program with ordered arguments. The named shell owns all subsequent parsing, expansion, and security behavior.

## Examples

| Editor text | Program | Arguments in order |
| --- | --- | --- |
| `python -m http.server` | `python` | `-m`, `http.server` |
| `"C:\Program Files\Tool\tool.exe" --name "Ada Lovelace"` | `C:\Program Files\Tool\tool.exe` | `--name`, `Ada Lovelace` |
| `/usr/bin/printf '%s\n' 'hello world'` | `/usr/bin/printf` | `%s\n`, `hello world` |
| `program --tag one --tag two ''` | `program` | `--tag`, `one`, `--tag`, `two`, empty string |
| `program "first` + literal newline + `second"` | `program` | one argument containing a literal newline |
| `program a\ b` | `program` | `a b` |
| `program '$HOME' '|' '*.txt'` | `program` | `$HOME`, `|`, `*.txt` |

## Canonical formatting contract

- Empty values are quoted.
- Values without whitespace or quotation ambiguity are emitted bare.
- Values requiring grouping use a reversible quote style that preserves backslashes literally.
- Values containing both quote styles use adjacent segments and an escaped quote.
- For every valid invocation `I`, `Parse(Format(I))` equals `I` exactly.
- Formatting is stable: `Format(Parse(Format(I)))` equals `Format(I)`.

## Preview contract

- Label the program as **Program**.
- Label the ordered values as **Arguments in order** and include the exact count.
- Number arguments starting at 1 for users.
- Display every exact value in a quoted escaped notation so empty strings and control characters are visible.
- On invalid text, show no stale program or arguments. Show only the actionable syntax error.

## Compatibility contract

- The API continues receiving `command` and `args` separately.
- SQLite continues storing the program and JSON argument array separately.
- CLI `--command` and repeated `--arg` behavior is unchanged.
- The executor continues direct process invocation and hidden-console configuration.
- Existing tasks are represented with canonical formatting and require no migration.
