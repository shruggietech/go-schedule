package clievaluate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

const documentLimit = 32 << 20

// Decode rejects duplicate keys, invalid UTF-8/BOM, unknown fields and trailing values.
func Decode(input io.Reader, target any) error {
	data, err := io.ReadAll(io.LimitReader(input, documentLimit+1))
	if err != nil {
		return fmt.Errorf("read QA document: %w", err)
	}
	if len(data) > documentLimit || !utf8.Valid(data) || bytes.HasPrefix(data, []byte{239, 187, 191}) {
		return errors.New("QA document must be bounded UTF-8 without BOM")
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	if err := uniqueJSON(d, 0); err != nil {
		return fmt.Errorf("invalid QA JSON: %w", err)
	}
	if _, err := d.Token(); err != io.EOF {
		return errors.New("QA document must contain exactly one value")
	}
	d = json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(target); err != nil {
		return fmt.Errorf("decode QA contract: %w", err)
	}
	return nil
}

func uniqueJSON(d *json.Decoder, depth int) error {
	if depth > 64 {
		return errors.New("JSON nesting exceeds 64 levels")
	}
	t, err := d.Token()
	if err != nil {
		return err
	}
	delim, ok := t.(json.Delim)
	if !ok {
		return nil
	}
	switch delim {
	case '{':
		seen := map[string]bool{}
		for d.More() {
			key, err := d.Token()
			if err != nil {
				return err
			}
			name, ok := key.(string)
			if !ok || seen[name] {
				return errors.New("duplicate or invalid object key")
			}
			seen[name] = true
			if err := uniqueJSON(d, depth+1); err != nil {
				return err
			}
		}
	case '[':
		for d.More() {
			if err := uniqueJSON(d, depth+1); err != nil {
				return err
			}
		}
	default:
		return errors.New("unexpected closing delimiter")
	}
	_, err = d.Token()
	return err
}
