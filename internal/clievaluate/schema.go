package clievaluate

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
)

//go:embed schemas/v1/*.schema.json
var schemaFiles embed.FS

var contractTypes = map[string]reflect.Type{
	"manifest": reflect.TypeFor[Manifest](), "report": reflect.TypeFor[Report](), "corpus": reflect.TypeFor[Corpus](), "reviews": reflect.TypeFor[ReviewSet](), "event": reflect.TypeFor[Event](), "evidence-manifest": reflect.TypeFor[EvidenceManifest](),
	"build-receipt": reflect.TypeFor[BuildReceipt](),
}

// SchemaDocuments deterministically derives strict schemas for the concrete contracts.
// Checked-in copies are required to match by the contract regression tests.
func SchemaDocuments() (map[string][]byte, error) {
	result := map[string][]byte{}
	for name, t := range contractTypes {
		defs := map[string]any{}
		root := schemaType(t, defs)
		doc := map[string]any{"$schema": "https://json-schema.org/draft/2020-12/schema", "$id": "https://go-schedule.io/schemas/cli-qa/v1/" + name + ".schema.json", "$ref": root["$ref"], "$defs": defs}
		data, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			return nil, err
		}
		result[name+".schema.json"] = append(data, '\n')
	}
	return result, nil
}

func schemaType(t reflect.Type, defs map[string]any) map[string]any {
	if t.Kind() == reflect.Pointer {
		return map[string]any{"anyOf": []any{schemaType(t.Elem(), defs), map[string]any{"type": "null"}}}
	}
	if t.Kind() == reflect.Struct {
		name := t.Name()
		ref := map[string]any{"$ref": "#/$defs/" + name}
		if _, ok := defs[name]; ok {
			return ref
		}
		defs[name] = nil
		properties := map[string]any{}
		required := []string{}
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			key := field.Tag.Get("json")
			if key == "" || strings.Contains(key, ",") {
				panic("QA contracts require explicit non-omitted JSON fields")
			}
			properties[key] = schemaField(name, key, field.Type, defs)
			required = append(required, key)
		}
		defs[name] = map[string]any{"type": "object", "additionalProperties": false, "required": required, "properties": properties}
		return ref
	}
	switch t.Kind() {
	case reflect.Slice:
		return map[string]any{"type": "array", "maxItems": 20000, "items": schemaType(t.Elem(), defs)}
	case reflect.Bool:
		return map[string]any{"type": "boolean"}
	case reflect.Int, reflect.Int64:
		return map[string]any{"type": "integer", "minimum": 0, "maximum": int64(1<<53 - 1)}
	case reflect.String:
		return map[string]any{"type": "string", "maxLength": 262144}
	default:
		panic("unsupported QA contract field: " + t.String())
	}
}

func schemaField(owner, key string, t reflect.Type, defs map[string]any) map[string]any {
	s := schemaType(t, defs)
	if key == "exit_code" || key == "expected_exit" {
		number := map[string]any{"type": "integer", "minimum": -1, "maximum": int64(4294967295)}
		if t.Kind() == reflect.Pointer {
			return map[string]any{"anyOf": []any{number, map[string]any{"type": "null"}}}
		}
		return number
	}
	if t.Kind() != reflect.String {
		return s
	}
	if key == "schema_version" {
		s["const"] = Version
	}
	if key == "kind" {
		names := map[string]string{"Manifest": "manifest", "Report": "report", "Corpus": "corpus", "ReviewSet": "reviews", "Event": "event", "EvidenceManifest": "evidence-manifest", "BuildReceipt": "build-receipt"}
		s["const"] = "go-schedule.cli-qa-" + names[owner]
	}
	if key == "started_at" || key == "finished_at" || key == "at" {
		s["format"] = "date-time"
		s["pattern"] = "Z$"
	}
	if key == "execution_kind" {
		s["enum"] = []string{"packaged_cli", "command_harness", "deterministic_engine"}
	}
	if key == "status" || key == "execution_status" || key == "outcome" {
		if owner == "Environment" {
			s["enum"] = []string{"ready", "blocked"}
		} else {
			s["enum"] = []string{"pass", "fail", "blocked", "not_applicable"}
		}
	}
	if key == "verdict" {
		s["enum"] = []string{"pass", "fail", "not_applicable"}
		if owner == "Review" {
			s["enum"] = []string{"pass", "fail"}
		}
	}
	if key == "category" {
		s["enum"] = []string{"functional", "cli-parity", "help", "error-message", "logging", "performance", "harness"}
	}
	if key == "target" {
		s["enum"] = []string{"product", "harness"}
	}
	if key == "severity" {
		s["enum"] = []string{"P1", "P2", "P3"}
	}
	if key == "phase" {
		s["enum"] = []string{"capturing", "pending_review", "finalized"}
	}
	if key == "completion_state" {
		s["enum"] = []string{"running", "completed", "interrupted"}
	}
	if key == "channel" {
		s["enum"] = []string{"help", "description", "stdout", "stderr", "log"}
	}
	if owner == "Event" && key == "type" {
		s["enum"] = eventTypes
	}
	return s
}

// Encode produces stable JSON with required empty collections represented as arrays.
func Encode(value any) ([]byte, error) {
	if value == nil {
		return nil, errors.New("QA contract cannot be null")
	}
	v := reflect.New(reflect.TypeOf(value))
	v.Elem().Set(reflect.ValueOf(value))
	normalizeArrays(v)
	data, err := json.MarshalIndent(v.Elem().Interface(), "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

func normalizeArrays(v reflect.Value) {
	if v.Kind() == reflect.Pointer {
		if !v.IsNil() {
			normalizeArrays(v.Elem())
		}
		return
	}
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			normalizeArrays(v.Field(i))
		}
	}
	if v.Kind() == reflect.Slice {
		if v.IsNil() && v.CanSet() {
			v.Set(reflect.MakeSlice(v.Type(), 0, 0))
		}
		for i := 0; i < v.Len(); i++ {
			normalizeArrays(v.Index(i))
		}
	}
}

// ValidateSchema uses only the embedded schema, with required date-time assertions.
func ValidateSchema(name string, data []byte) error {
	if _, ok := contractTypes[name]; !ok {
		return errors.New("unknown QA schema")
	}
	source, err := schemaFiles.ReadFile("schemas/v1/" + name + ".schema.json")
	if err != nil {
		return err
	}
	var doc any
	if err := json.Unmarshal(source, &doc); err != nil {
		return err
	}
	compiler := jsonschema.NewCompiler()
	compiler.DefaultDraft(jsonschema.Draft2020)
	compiler.AssertFormat()
	uri := "https://go-schedule.io/schemas/cli-qa/v1/" + name + ".schema.json"
	if err := compiler.AddResource(uri, doc); err != nil {
		return err
	}
	schema, err := compiler.Compile(uri)
	if err != nil {
		return err
	}
	var value any
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&value); err != nil {
		return err
	}
	if err := schema.Validate(value); err != nil {
		return fmt.Errorf("QA %s schema: %w", name, err)
	}
	return nil
}
