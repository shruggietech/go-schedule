package clievaluate

import (
	"bytes"
	"strings"
	"testing"
)

func TestSchemasMatchConcreteStrictContracts(t *testing.T) {
	docs, err := SchemaDocuments()
	if err != nil {
		t.Fatal(err)
	}
	for name, want := range docs {
		got, err := schemaFiles.ReadFile("schemas/v1/" + name)
		if err != nil || !bytes.Equal(got, want) {
			t.Fatalf("schema drift %s: regenerate from SchemaDocuments: %v", name, err)
		}
	}
	for name := range contractTypes {
		if err := ValidateSchema(name, []byte(`{}`)); err == nil {
			t.Fatalf("missing fields accepted by %s", name)
		}
	}
	corpus := Corpus{SchemaVersion: Version, Kind: "go-schedule.cli-qa-corpus", RunID: "run-1"}
	data, err := Encode(corpus)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateSchema("corpus", data); err != nil {
		t.Fatal(err)
	}
}

func TestStrictDecodeRejectsAmbiguousAndUnboundedJSON(t *testing.T) {
	for _, input := range []string{`{"schema_version":"1","schema_version":"2"}`, `{"unknown":true}`, `{} {}`, "\ufeff{}", `{"schema_version":`} {
		var value ReviewSet
		if err := Decode(strings.NewReader(input), &value); err == nil {
			t.Fatalf("invalid JSON accepted: %q", input)
		}
	}
}

func TestPortableArtifactPathsRejectCrossPlatformEscapes(t *testing.T) {
	for _, path := range []string{"", ".", "../out", "/absolute", `C:/file`, `\\server\share`, `a\b`, "a/../b", "CON", "a:stream", "x/COM1.txt", "a/", "a//b"} {
		if err := ValidateArtifactPath(path); err == nil {
			t.Fatalf("unsafe artifact path accepted: %q", path)
		}
	}
	if err := ValidateArtifactPath("windows/case-001/stdout.txt"); err != nil {
		t.Fatal(err)
	}
}

func TestLegacyReportContractAcceptsHashBoundDiagnosticState(t *testing.T) {
	store, r, m, c := reportFixture(t)
	diagnostic, err := store.Write("env/state.json", "application/json", []byte(`{"diagnostics":{"phases":[{"name":"after","duration_ms":2}],"snapshot_complete":false,"failed_request":"GET /v1/access/credentials"}}`))
	if err != nil {
		t.Fatal(err)
	}
	r.Artifacts = append(r.Artifacts, diagnostic)
	r.Cases[0].Observed.Checks[0].Evidence = append(r.Cases[0].Observed.Checks[0].Evidence, diagnostic.Path)
	data, err := Encode(r)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateSchema("report", data); err != nil {
		t.Fatal(err)
	}
	var decoded Report
	if err := Decode(bytes.NewReader(data), &decoded); err != nil {
		t.Fatal(err)
	}
	if err := ValidateReport(store, decoded, m, c); err != nil {
		t.Fatal(err)
	}
	roundTrip, err := Encode(decoded)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, roundTrip) {
		t.Fatal("legacy report serialization changed")
	}
}
