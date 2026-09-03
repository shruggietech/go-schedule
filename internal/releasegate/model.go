// Package releasegate validates exact Windows release candidate evidence.
package releasegate

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
	"unicode/utf8"
)

const maxJSONBytes = 16 << 20

// OperatorAttestation is the exact statement required in release evidence.
const OperatorAttestation = "I attest that these observations were made against the identified candidate in the identified environments."

// Evidence is the canonical versioned release-candidate evidence record.
type Evidence struct {
	SchemaVersion int           `json:"schema_version"`
	EvidenceClass string        `json:"evidence_class"`
	Candidate     Candidate     `json:"candidate"`
	Operator      Operator      `json:"operator"`
	StartedAt     time.Time     `json:"started_at"`
	CompletedAt   time.Time     `json:"completed_at"`
	Environments  []Environment `json:"environments"`
	Observations  []Observation `json:"observations"`
	Attachments   []Attachment  `json:"attachments"`
}

// Candidate binds evidence to one staged MSI and its GitHub origin.
type Candidate struct {
	Repository     string `json:"repository"`
	Tag            string `json:"tag"`
	Commit         string `json:"commit"`
	Workflow       string `json:"workflow"`
	RunID          int64  `json:"run_id"`
	RunAttempt     int    `json:"run_attempt"`
	Filename       string `json:"filename"`
	Bytes          int64  `json:"bytes"`
	SHA256         string `json:"sha256"`
	ProductVersion string `json:"product_version"`
	ProductCode    string `json:"product_code"`
}

// Operator records a role-oriented attestation without requiring a name.
type Operator struct {
	Role       string    `json:"role"`
	AttestedAt time.Time `json:"attested_at"`
	Statement  string    `json:"statement"`
}

// Environment describes one real Windows observation context.
type Environment struct {
	ID              string `json:"id"`
	Snapshot        string `json:"snapshot"`
	CleanSnapshot   bool   `json:"clean_snapshot"`
	WindowsEdition  string `json:"windows_edition"`
	WindowsBuild    string `json:"windows_build"`
	AccountRole     string `json:"account_role"`
	AccountSID      string `json:"account_sid"`
	Integrity       string `json:"integrity"`
	IntegrityRID    int    `json:"integrity_rid"`
	ServiceIdentity string `json:"service_identity"`
	DisplayClass    string `json:"display_class"`
	EffectiveDPI    int    `json:"effective_dpi"`
	ProfileState    string `json:"profile_state"`
}

// Observation records one fixed release scenario.
type Observation struct {
	ID              string         `json:"id"`
	EnvironmentID   string         `json:"environment_id"`
	Status          string         `json:"status"`
	StartedAt       time.Time      `json:"started_at"`
	CompletedAt     time.Time      `json:"completed_at"`
	Summary         string         `json:"summary"`
	Metrics         map[string]any `json:"metrics"`
	AttachmentPaths []string       `json:"attachment_paths"`
}

// Attachment integrity-protects one supporting file in the evidence bundle.
type Attachment struct {
	Path      string `json:"path"`
	Bytes     int64  `json:"bytes"`
	SHA256    string `json:"sha256"`
	MediaType string `json:"media_type"`
	Purpose   string `json:"purpose"`
}

// ExpectedIdentity is supplied by promotion from authoritative GitHub state.
type ExpectedIdentity struct {
	Repository string
	Tag        string
	Commit     string
}

// DecodeEvidence performs strict JSON decoding and rejects trailing values.
func DecodeEvidence(reader io.Reader) (Evidence, error) {
	data, err := readStrictJSON(reader, "evidence")
	if err != nil {
		return Evidence{}, err
	}
	buffered := bufio.NewReader(bytes.NewReader(data))
	prefix, _ := buffered.Peek(3)
	if bytes.Equal(prefix, []byte{0xef, 0xbb, 0xbf}) {
		return Evidence{}, fmt.Errorf("evidence JSON must be UTF-8 without BOM")
	}

	decoder := json.NewDecoder(buffered)
	decoder.DisallowUnknownFields()
	decoder.UseNumber()
	var evidence Evidence
	if err := decoder.Decode(&evidence); err != nil {
		return Evidence{}, fmt.Errorf("decode evidence: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return Evidence{}, fmt.Errorf("decode evidence: multiple JSON values")
		}
		return Evidence{}, fmt.Errorf("decode evidence trailing content: %w", err)
	}
	return evidence, nil
}

// DecodeCandidate performs strict decoding of a staged candidate manifest.
func DecodeCandidate(reader io.Reader) (Candidate, error) {
	data, err := readStrictJSON(reader, "candidate manifest")
	if err != nil {
		return Candidate{}, err
	}
	buffered := bufio.NewReader(bytes.NewReader(data))
	prefix, _ := buffered.Peek(3)
	if bytes.Equal(prefix, []byte{0xef, 0xbb, 0xbf}) {
		return Candidate{}, fmt.Errorf("candidate manifest must be UTF-8 without BOM")
	}
	decoder := json.NewDecoder(buffered)
	decoder.DisallowUnknownFields()
	var candidate Candidate
	if err := decoder.Decode(&candidate); err != nil {
		return Candidate{}, fmt.Errorf("decode candidate manifest: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return Candidate{}, fmt.Errorf("decode candidate manifest: multiple JSON values")
		}
		return Candidate{}, fmt.Errorf(
			"decode candidate manifest trailing content: %w",
			err,
		)
	}
	return candidate, nil
}

func readStrictJSON(reader io.Reader, label string) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(reader, maxJSONBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read %s JSON: %w", label, err)
	}
	if len(data) > maxJSONBytes {
		return nil, fmt.Errorf("%s JSON exceeds %d bytes", label, maxJSONBytes)
	}
	if !utf8.Valid(data) {
		return nil, fmt.Errorf("%s JSON is not valid UTF-8", label)
	}
	if err := rejectDuplicateJSONKeys(data); err != nil {
		return nil, fmt.Errorf("decode %s: %w", label, err)
	}
	return data, nil
}

func rejectDuplicateJSONKeys(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := scanJSONValue(decoder); err != nil {
		return err
	}
	if _, err := decoder.Token(); err != io.EOF {
		if err == nil {
			return fmt.Errorf("multiple JSON values")
		}
		return err
	}
	return nil
}

func scanJSONValue(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delimiter, ok := token.(json.Delim)
	if !ok {
		return nil
	}
	switch delimiter {
	case '{':
		seen := make(map[string]bool)
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return err
			}
			key, ok := keyToken.(string)
			if !ok {
				return fmt.Errorf("object member name is not a string")
			}
			if seen[key] {
				return fmt.Errorf("duplicate object key %q", key)
			}
			seen[key] = true
			if err := scanJSONValue(decoder); err != nil {
				return err
			}
		}
		closing, err := decoder.Token()
		if err != nil {
			return err
		}
		if closing != json.Delim('}') {
			return fmt.Errorf("object is not terminated")
		}
	case '[':
		for decoder.More() {
			if err := scanJSONValue(decoder); err != nil {
				return err
			}
		}
		closing, err := decoder.Token()
		if err != nil {
			return err
		}
		if closing != json.Delim(']') {
			return fmt.Errorf("array is not terminated")
		}
	default:
		return fmt.Errorf("unexpected delimiter %q", delimiter)
	}
	return nil
}
