// Command github-format checks or repairs repository-authored GitHub text.
package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

var (
	listItemPattern = regexp.MustCompile(`^(?:[-+*]|[0-9]+[.)])\s+`)
	linkDefinition  = regexp.MustCompile(`^\[[^]]+\]:\s*`)
	metadataLine    = regexp.MustCompile(`^\*\*[^*]+\*\*:\s*`)
	emDash          = string(rune(0x2014))
)

func main() {
	fix := flag.Bool("fix", false, "rewrite violations in place")
	stdin := flag.Bool("stdin", false, "format Markdown from standard input")
	flag.Parse()
	if *stdin {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "github-format: read stdin: %v\n", err)
			os.Exit(2)
		}
		formatted := formatMarkdown(replaceEmDashes(string(data)))
		fmt.Print(formatted)
		return
	}
	root := "."
	if flag.NArg() == 1 {
		root = flag.Arg(0)
	} else if flag.NArg() > 1 {
		fmt.Fprintln(os.Stderr, "usage: github-format [-fix] [repository-root]")
		os.Exit(2)
	}
	problems, err := checkRepository(root, *fix)
	if err != nil {
		fmt.Fprintf(os.Stderr, "github-format: %v\n", err)
		os.Exit(2)
	}
	if len(problems) != 0 {
		for _, problem := range problems {
			fmt.Fprintln(os.Stderr, problem)
		}
		fmt.Fprintf(os.Stderr, "github-format: FAILED with %d file(s)\n", len(problems))
		os.Exit(1)
	}
	fmt.Println("github-format: OK - no em dashes or hard-wrapped Markdown prose")
}

func checkRepository(root string, fix bool) ([]string, error) {
	var problems []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if path != root && skipDirectory(entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if !textFile(path) {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if !utf8.Valid(data) {
			problems = append(problems, filepath.ToSlash(relative)+": invalid UTF-8")
			return nil
		}
		original := string(data)
		content := original
		if strings.HasPrefix(content, "\ufeff") {
			if !fix {
				problems = append(problems, filepath.ToSlash(relative)+": UTF-8 BOM")
				return nil
			}
			content = strings.TrimPrefix(content, "\ufeff")
		}
		formatted := replaceEmDashes(content)
		if strings.EqualFold(filepath.Ext(path), ".md") {
			formatted = formatMarkdown(formatted)
		}
		if formatted == original {
			return nil
		}
		if !fix {
			problems = append(problems, filepath.ToSlash(relative))
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		return os.WriteFile(path, []byte(formatted), info.Mode().Perm())
	})
	sort.Strings(problems)
	return problems, err
}

func skipDirectory(name string) bool {
	switch name {
	case ".git", ".idea", ".vscode", "dist", "node_modules":
		return true
	default:
		return false
	}
}

func textFile(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".css", ".desktop", ".go", ".html", ".js", ".json", ".jsx", ".md", ".ps1", ".py", ".scss", ".sh", ".svg", ".toml", ".ts", ".tsx", ".txt", ".wxs", ".xml", ".yaml", ".yml":
		return true
	}
	switch filepath.Base(path) {
	case ".gitattributes", ".gitignore", "CODEOWNERS", "Dockerfile", "LICENSE", "Makefile":
		return true
	default:
		return false
	}
}

func replaceEmDashes(input string) string {
	lines := strings.Split(input, "\n")
	for index, line := range lines {
		if !strings.Contains(line, emDash) {
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "|") {
			lines[index] = strings.ReplaceAll(line, emDash, "N/A")
			continue
		}
		line = strings.ReplaceAll(line, "// "+emDash+" ", "// ")
		line = strings.ReplaceAll(line, "# "+emDash+" ", "# ")
		line = strings.ReplaceAll(line, " "+emDash+" ", ", ")
		line = strings.ReplaceAll(line, emDash+" ", "")
		line = strings.ReplaceAll(line, " "+emDash, ",")
		lines[index] = strings.ReplaceAll(line, emDash, "-")
	}
	return strings.Join(lines, "\n")
}

func formatMarkdown(input string) string {
	newline := "\n"
	if strings.Contains(input, "\r\n") {
		newline = "\r\n"
	}
	normalized := strings.ReplaceAll(input, "\r\n", "\n")
	hasFinalNewline := strings.HasSuffix(normalized, "\n")
	lines := strings.Split(strings.TrimSuffix(normalized, "\n"), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return input
	}
	out := make([]string, 0, len(lines))
	inFence := false
	inFrontMatter := len(lines) > 0 && strings.TrimSpace(lines[0]) == "---"
	inHTMLComment := false
	for index := 0; index < len(lines); {
		line := lines[index]
		trimmed := strings.TrimSpace(line)
		if index == 0 && inFrontMatter {
			out = append(out, line)
			index++
			continue
		}
		if inFrontMatter {
			out = append(out, line)
			index++
			if trimmed == "---" {
				inFrontMatter = false
			}
			continue
		}
		if inHTMLComment {
			out = append(out, line)
			index++
			if strings.Contains(line, "-->") {
				inHTMLComment = false
			}
			continue
		}
		if strings.HasPrefix(trimmed, "<!--") {
			out = append(out, line)
			index++
			inHTMLComment = !strings.Contains(line, "-->")
			continue
		}
		if isFence(trimmed) {
			inFence = !inFence
			out = append(out, line)
			index++
			continue
		}
		if inFence || isLiteralMarkdownLine(line) {
			out = append(out, line)
			index++
			continue
		}
		if content, ok := blockquoteProse(line); ok {
			joined := content
			index++
			for index < len(lines) {
				content, ok = blockquoteProse(lines[index])
				if !ok {
					break
				}
				joined += " " + content
				index++
			}
			out = append(out, "> "+joined)
			continue
		}
		if listItemPattern.MatchString(trimmed) {
			joined := strings.TrimRight(line, " \t")
			marker := listItemPattern.FindString(trimmed)
			contentIndent := leadingIndent(line) + len(marker)
			index++
			for index < len(lines) && isProseContinuation(lines[index], contentIndent) {
				joined += " " + strings.TrimSpace(lines[index])
				index++
			}
			out = append(out, joined)
			continue
		}
		joined := strings.TrimRight(line, " \t")
		index++
		for index < len(lines) && isPlainProse(lines[index]) {
			joined += " " + strings.TrimSpace(lines[index])
			index++
		}
		out = append(out, joined)
	}
	result := strings.Join(out, newline)
	if hasFinalNewline {
		result += newline
	}
	return result
}

func isFence(trimmed string) bool {
	return strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~")
}

func isProseContinuation(line string, contentIndent int) bool {
	if strings.TrimSpace(line) == "" {
		return false
	}
	trimmed := strings.TrimSpace(line)
	if listItemPattern.MatchString(trimmed) || markdownBlockStart(trimmed) {
		return false
	}
	indent := leadingIndent(line)
	return indent < contentIndent+4
}

func leadingIndent(line string) int {
	indent := 0
	for _, character := range line {
		switch character {
		case ' ':
			indent++
		case '\t':
			indent += 4
		default:
			return indent
		}
	}
	return indent
}

func blockquoteProse(line string) (string, bool) {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "> ") {
		return "", false
	}
	content := strings.TrimSpace(strings.TrimPrefix(trimmed, ">"))
	if content == "" || listItemPattern.MatchString(content) || markdownBlockStart(content) {
		return "", false
	}
	return content, true
}

func markdownBlockStart(trimmed string) bool {
	return strings.HasPrefix(trimmed, ">") || isFence(trimmed) || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "<") || strings.HasPrefix(trimmed, "|") || strings.HasPrefix(trimmed, "![") || strings.HasPrefix(trimmed, ":::") || strings.HasPrefix(trimmed, "$$")
}

func isPlainProse(line string) bool {
	if strings.TrimSpace(line) == "" || isLiteralMarkdownLine(line) {
		return false
	}
	trimmed := strings.TrimSpace(line)
	return !listItemPattern.MatchString(trimmed) && !strings.HasPrefix(trimmed, ">")
}

func isLiteralMarkdownLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || isFence(trimmed) || strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t") || strings.HasSuffix(line, "\\") || strings.HasSuffix(line, "  ") {
		return true
	}
	if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "<") || strings.HasPrefix(trimmed, "|") || strings.HasPrefix(trimmed, "![") || strings.HasPrefix(trimmed, ":::") || strings.HasPrefix(trimmed, "$$") {
		return true
	}
	if trimmed == "---" || trimmed == "***" || trimmed == "___" || linkDefinition.MatchString(trimmed) || metadataLine.MatchString(trimmed) {
		return true
	}
	return false
}
