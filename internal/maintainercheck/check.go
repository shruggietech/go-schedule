// Package maintainercheck evaluates repository policy without per-rule utility processes.
package maintainercheck

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"go.yaml.in/yaml/v3"

	"github.com/shruggietech/go-schedule/v2/internal/platform"
)

type checker struct {
	root     string
	files    map[string]string
	findings []string
}

func newChecker(root string) *checker { return &checker{root: root, files: map[string]string{}} }
func (c *checker) report(file, message string) {
	c.findings = append(c.findings, filepath.Join(c.root, file)+": "+message)
}
func (c *checker) read(file string) string {
	if data, ok := c.files[file]; ok {
		return data
	}
	data, err := os.ReadFile(filepath.Join(c.root, file))
	if err != nil {
		c.report(file, fmt.Sprintf("read required input: %v", err))
	}
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	c.files[file] = text
	return text
}
func (c *checker) need(file, text, description string) {
	if !strings.Contains(c.read(file), text) {
		c.report(file, "missing "+description)
	}
}
func (c *checker) yaml(file string) *yaml.Node {
	data := c.read(file)
	var doc yaml.Node
	decoder := yaml.NewDecoder(strings.NewReader(data))
	if err := decoder.Decode(&doc); err != nil {
		c.report(file, fmt.Sprintf("invalid YAML: %v; restore unambiguous workflow policy", err))
		return nil
	}
	var extra yaml.Node
	if err := decoder.Decode(&extra); err != io.EOF {
		c.report(file, "invalid YAML: expected one document")
		return nil
	}
	if len(doc.Content) != 1 || doc.Content[0].Kind != yaml.MappingNode {
		c.report(file, "invalid YAML: expected mapping")
		return nil
	}
	if err := validateNode(doc.Content[0]); err != nil {
		c.report(file, "invalid YAML: "+err.Error())
		return nil
	}
	return doc.Content[0]
}
func validateNode(n *yaml.Node) error {
	active, done := map[*yaml.Node]bool{}, map[*yaml.Node]bool{}
	var walk func(*yaml.Node) error
	walk = func(node *yaml.Node) error {
		if node == nil {
			return fmt.Errorf("nil YAML alias target")
		}
		if active[node] {
			return fmt.Errorf("line %d: cyclic YAML alias", node.Line)
		}
		if done[node] {
			return nil
		}
		active[node] = true
		defer delete(active, node)
		if node.Kind == yaml.MappingNode {
			seen := map[string]bool{}
			for i := 0; i < len(node.Content); i += 2 {
				key := node.Content[i]
				if key.Kind != yaml.ScalarNode {
					return fmt.Errorf("line %d: non-scalar mapping key", key.Line)
				}
				if seen[key.Value] {
					return fmt.Errorf("line %d: duplicate key %q", key.Line, key.Value)
				}
				seen[key.Value] = true
			}
		}
		if node.Kind == yaml.AliasNode {
			if err := walk(node.Alias); err != nil {
				return err
			}
		}
		for _, child := range node.Content {
			if err := walk(child); err != nil {
				return err
			}
		}
		done[node] = true
		return nil
	}
	return walk(n)
}
func unalias(n *yaml.Node) *yaml.Node {
	for n != nil && n.Kind == yaml.AliasNode {
		n = n.Alias
	}
	return n
}

func mappingField(n *yaml.Node, key string) *yaml.Node {
	n = unalias(n)
	if n == nil || n.Kind != yaml.MappingNode {
		return nil
	}
	var merge *yaml.Node
	for i := 0; i < len(n.Content); i += 2 {
		if n.Content[i].Value == key {
			return n.Content[i+1]
		}
		if n.Content[i].Value == "<<" {
			merge = n.Content[i+1]
		}
	}
	merge = unalias(merge)
	if merge != nil && merge.Kind == yaml.SequenceNode {
		for _, item := range merge.Content {
			if found := mappingField(item, key); found != nil {
				return found
			}
		}
		return nil
	}
	if merge != nil {
		return mappingField(merge, key)
	}
	return nil
}
func field(n *yaml.Node, keys ...string) *yaml.Node {
	for _, key := range keys {
		n = mappingField(n, key)
	}
	return n
}

func value(n *yaml.Node) string {
	n = unalias(n)
	if n == nil {
		return ""
	}
	return n.Value
}
func items(n *yaml.Node) []*yaml.Node {
	n = unalias(n)
	if n == nil {
		return nil
	}
	if n.Kind == yaml.SequenceNode {
		return n.Content
	}
	return nil
}
func has(n *yaml.Node, key string) bool { return field(n, key) != nil }
func joined(n *yaml.Node) string {
	n = unalias(n)
	if n == nil {
		return ""
	}
	if n.Kind != yaml.SequenceNode {
		return value(n)
	}
	var values []string
	for _, v := range n.Content {
		values = append(values, value(v))
	}
	return strings.Join(values, ", ")
}
func steps(n *yaml.Node) []*yaml.Node {
	var result []*yaml.Node
	jobs := unalias(field(n, "jobs"))
	if jobs == nil {
		return nil
	}
	for i := 1; i < len(jobs.Content); i += 2 {
		result = append(result, items(field(jobs.Content[i], "steps"))...)
	}
	return result
}

// semanticText retains decoded executable values and keys, excluding YAML comments.
func semanticText(n *yaml.Node) string {
	if n == nil {
		return ""
	}
	var out strings.Builder
	var visit func(*yaml.Node)
	seen := map[*yaml.Node]bool{}
	visit = func(n *yaml.Node) {
		n = unalias(n)
		if n == nil || seen[n] {
			return
		}
		seen[n] = true
		switch n.Kind {
		case yaml.MappingNode:
			for i := 0; i < len(n.Content); i += 2 {
				v := n.Content[i+1]
				out.WriteString(n.Content[i].Value + ":")
				switch v.Kind {
				case yaml.ScalarNode:
					out.WriteString(" " + v.Value)
				case yaml.SequenceNode:
					var scalar = true
					for _, x := range v.Content {
						scalar = scalar && x.Kind == yaml.ScalarNode
					}
					if scalar {
						out.WriteString(" [" + joined(v) + "]")
					}
				}
				out.WriteByte('\n')
				if v.Kind != yaml.ScalarNode {
					visit(v)
				}
			}
		case yaml.SequenceNode:
			for _, v := range n.Content {
				visit(v)
			}
		}
	}
	visit(n)
	return out.String()
}
func (c *checker) requireText(file, text, description string, n *yaml.Node) {
	if !strings.Contains(semanticText(n), strings.TrimSpace(text)) {
		c.report(file, "missing "+description)
	}
}
func (c *checker) require(file string, ok bool, message string) {
	if !ok {
		c.report(file, message)
	}
}
func (c *checker) child(file, exe string, args ...string) string {
	cmd := exec.Command(exe, args...)
	cmd.Dir = c.root
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	platform.HideConsole(cmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.report(file, fmt.Sprintf("child failed: %v: %s", err, strings.TrimSpace(string(output))))
	}
	return strings.TrimSpace(string(output))
}

func containsScalar(n *yaml.Node, want string) bool {
	for _, v := range items(n) {
		if value(v) == want {
			return true
		}
	}
	return false
}
