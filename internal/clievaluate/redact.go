package clievaluate

import (
	"bytes"
	"encoding/json"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var sensitiveField = regexp.MustCompile(`(?i)"(?:token|secret|credential|password|authorization|pairing_code|hmac_key|bearer_token)"\s*:\s*"((?:\\.|[^"\\])*)"`)
var authorizationValue = regexp.MustCompile(`(?i)((?:authorization\s*[:=]\s*bearer\s+|(?:--key|--token|--secret|--password)(?:=|\s+)|[?&](?:token|key|secret)=)[A-Za-z0-9_./+%=-]{8,}|\bbearer\s+[A-Za-z0-9_./+%=-]{20,})`)

// Redactor owns fixture-only secret values in memory; the registry is never serialized.
type Redactor struct {
	mu       sync.Mutex
	values   []string
	seen     map[string]bool
	replacer *strings.Replacer
}

// NewRedactor creates the shared redaction sink for one run.
func NewRedactor(values ...string) *Redactor { r := &Redactor{}; r.Register(values...); return r }

// Register adds dynamic fixture secrets before captured output is persisted.
func (r *Redactor) Register(values ...string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.seen == nil {
		r.seen = map[string]bool{}
	}
	changed := false
	for _, value := range values {
		if len(value) >= 8 && value != "[REDACTED]" {
			variants := []string{value, url.QueryEscape(value)}
			encoded := strconv.Quote(value)
			if len(encoded) > 2 {
				variants = append(variants, encoded[1:len(encoded)-1])
			}
			for _, variant := range variants {
				if !r.seen[variant] {
					r.values = append(r.values, variant)
					r.seen[variant] = true
					changed = true
				}
			}
		}
	}
	if changed {
		sort.Slice(r.values, func(i, j int) bool { return len(r.values[i]) > len(r.values[j]) })
		r.replacer = nil
	}
}

// Sanitize operates on a complete bounded capture, including secrets split across writes.
func (r *Redactor) Sanitize(data []byte) []byte {
	for _, match := range sensitiveField.FindAllSubmatch(data, -1) {
		var value string
		if json.Unmarshal(append(append([]byte{'"'}, match[1]...), '"'), &value) == nil {
			r.Register(value)
		}
	}
	r.mu.Lock()
	if r.replacer == nil && len(r.values) > 0 {
		pairs := make([]string, 0, len(r.values)*2)
		for _, value := range r.values {
			pairs = append(pairs, value, "[REDACTED]")
		}
		r.replacer = strings.NewReplacer(pairs...)
	}
	replacer := r.replacer
	r.mu.Unlock()
	text := string(data)
	if replacer != nil {
		text = replacer.Replace(text)
	}
	clean := sensitiveField.ReplaceAllFunc([]byte(text), func(field []byte) []byte {
		colon := bytes.IndexByte(field, ':')
		return append(append([]byte{}, field[:colon+1]...), []byte(` "[REDACTED]"`)...)
	})
	return authorizationValue.ReplaceAll(clean, []byte("[REDACTED]"))
}
