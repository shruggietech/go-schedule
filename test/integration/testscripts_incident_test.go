package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/incident"
)

func TestScriptIncidentRetainsOrdinaryNativeFailure(t *testing.T) {
	parent := canonicalScriptIncidentParent(t)
	t.Setenv("GOSCHEDULE_INCIDENT_DIR", parent)
	args := []string{os.Args[0], "-test.run=^TestScriptObservationHelper$", "--", "overflow", "37", "unused", "unused"}
	got := observeScript(args, os.Args[0], os.Environ())
	if got.code != 37 || got.collectionErr == nil {
		t.Fatalf("original outcome changed: %+v", got)
	}
	var record incident.Record
	if err := json.Unmarshal([]byte(got.incident), &record); err != nil {
		t.Fatalf("missing incident: %v %s", err, got.incident)
	}
	if record.Original.NativeExit != 37 || record.Target.PID <= 0 || len(record.Target.SHA256) != 64 || record.Collection.Status != "absent" {
		t.Fatalf("incorrect incident: %+v", record)
	}
	if strings.Contains(got.incident, "SCRIPT_BIND_VALUE") {
		t.Fatal("incident retained environment values")
	}
}

func TestScriptIncidentUnavailableDoesNotReplaceNativeFailure(t *testing.T) {
	t.Setenv("GOSCHEDULE_INCIDENT_DIR", filepath.Join(canonicalScriptIncidentParent(t), "not-created"))
	got := observeScript([]string{os.Args[0], "-test.run=^TestScriptObservationHelper$", "--", "overflow", "37", "unused", "unused"}, os.Args[0], os.Environ())
	if got.code != 37 || !strings.Contains(got.incident, "unavailable") {
		t.Fatalf("original or independent collection lost: %+v", got)
	}
}

func TestScriptIncidentMissingInterpreterIsExplicitlyUnavailable(t *testing.T) {
	t.Setenv("GOSCHEDULE_INCIDENT_DIR", canonicalScriptIncidentParent(t))
	got := observeScript([]string{filepath.Join(canonicalScriptIncidentParent(t), "absent-interpreter")}, "", os.Environ())
	if got.code != -1 || got.err == nil || !strings.Contains(got.incident, "unavailable") {
		t.Fatalf("launch failure lost incident disposition: %+v", got)
	}
}

func TestPowerShellIncidentControlledNativeCrash(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("native Windows CoreCLR crash collection")
	}
	tw := allTwins()[0]
	for _, candidate := range allTwins() {
		if candidate.name == "powershell" {
			tw = candidate
		}
	}
	requireTwin(t, tw)
	parent := os.Getenv("GOSCHEDULE_INCIDENT_DIR")
	if parent == "" {
		parent = canonicalScriptIncidentParent(t)
	}
	t.Setenv("GOSCHEDULE_INCIDENT_DIR", parent)
	harness := filepath.Join(canonicalScriptIncidentParent(t), "owned-crash.ps1")
	if err := os.WriteFile(harness, []byte(powerShellIncidentPrelude()+"\n[Environment]::FailFast('owned S121 crash fixture')\n"), 0600); err != nil {
		t.Fatal(err)
	}
	got := observeScript(tw.command(harness), harness, os.Environ())
	if got.code == 0 || got.code == -1 {
		t.Fatalf("expected native crashing exit, got %d: %s", got.code, got.output)
	}
	var record incident.Record
	if err := json.Unmarshal([]byte(got.incident), &record); err != nil {
		t.Fatal(err)
	}
	if record.Original.NativeExit != got.code || record.Target.PID <= 0 {
		t.Fatalf("crash identity lost: %+v", record)
	}
	identityBytes, err := os.ReadFile(filepath.Join(parent, "incident-"+record.ID, "runtime.json"))
	if err != nil {
		t.Fatal(err)
	}
	var identity struct {
		PID           int    `json:"pid"`
		Observation   string `json:"observation"`
		Framework     string `json:"framework"`
		RuntimeSHA256 string `json:"runtime_sha256"`
	}
	if err := json.Unmarshal(identityBytes, &identity); err != nil || identity.PID != record.Target.PID || identity.Observation != "pre-payload" || identity.Framework == "" || len(identity.RuntimeSHA256) != 64 {
		t.Fatalf("not same-process runtime identity: %s %v", identityBytes, err)
	}
	// A missing runtime collector remains an explicit unavailable observation.
	if record.Collection.Status != "collected" && record.Collection.Status != "absent" && record.Collection.Status != "oversized" {
		t.Fatalf("unexplained collection: %+v", record)
	}
	if record.Collection.Status == "collected" {
		found := false
		for _, artifact := range record.Artifacts {
			if strings.HasSuffix(artifact.Name, ".dmp") && artifact.Status == "collected" && artifact.Bytes > 0 {
				found = true
			}
		}
		if !found {
			t.Fatal("claimed collection without native dump")
		}
	}
	t.Logf("controlled native exit=0x%08x incident=%s", uint32(got.code), got.incident)
}

func canonicalScriptIncidentParent(t *testing.T) string {
	t.Helper()
	parent, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return parent
}
