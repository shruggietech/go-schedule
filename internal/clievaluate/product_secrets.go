package clievaluate

import (
	"encoding/json"
	"regexp"
	"strings"
)

var capturedTriggerField = regexp.MustCompile(`"key"\s*:\s*"((?:\\.|[^"\\])*)"`)
var capturedTriggerCommand = regexp.MustCompile(`gosched(?:\.exe)?\s+trigger\s+fire\s+([A-Za-z0-9_-]+)`)

func registerCapturedSecrets(redactor *Redactor, definition step, data []byte) {
	path := definition.Spec.CommandPath
	if strings.HasPrefix(path, "trigger ") {
		for _, match := range capturedTriggerField.FindAllSubmatch(data, -1) {
			var key string
			if json.Unmarshal(append(append([]byte{'"'}, match[1]...), '"'), &key) == nil {
				redactor.Register(key)
			}
		}
		for _, match := range capturedTriggerCommand.FindAllSubmatch(data, -1) {
			redactor.Register(string(match[1]))
		}
		redactor.Register(humanLine(string(data), "key:"))
	}
	if path == "pairing create" {
		redactor.Register(humanLine(string(data), "Phrase (shown once):"))
	}
	if path == "mcp http enable" || path == "mcp http rotate" {
		redactor.Register(humanLine(string(data), "credential (shown once):"))
	}
}

// Trigger responses define key as a secret. System observations use key as identity.
func registerTriggerKeys(redactor *Redactor, value any) {
	switch item := value.(type) {
	case map[string]any:
		for key, child := range item {
			if key == "key" {
				if secret, ok := child.(string); ok {
					redactor.Register(secret)
				}
			}
			registerTriggerKeys(redactor, child)
		}
	case []any:
		for _, child := range item {
			registerTriggerKeys(redactor, child)
		}
	}
}
