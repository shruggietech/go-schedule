package releasegate

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestValidateAcceptsCompleteEvidence(t *testing.T) {
	t.Parallel()

	root, artifact, evidence := passingEvidence(t)
	if failures := Validate(evidence, root, artifact, ExpectedIdentity{
		Repository: "shruggietech/go-schedule",
		Tag:        "v1.0.0",
		Commit:     strings.Repeat("a", 40),
	}); len(failures) != 0 {
		t.Fatalf("Validate() failures = %v", failures)
	}
}

func TestRequiredScenarioIDsIncludeDesktopQualification(t *testing.T) {
	t.Parallel()

	want := []string{
		"desktop.appearance-standard",
		"desktop.appearance-scaled",
		"desktop.interaction-states",
		"desktop.interaction-states-scaled",
		"desktop.navigation-options",
		"desktop.navigation-options-scaled",
		"desktop.scroll-input",
		"desktop.tasks-table",
		"desktop.tasks-table-scaled",
		"desktop.schedule-activity-tables",
		"desktop.schedule-activity-tables-scaled",
	}
	ids := RequiredScenarioIDs()
	if len(ids) != 47 {
		t.Fatalf("RequiredScenarioIDs() count = %d, want 47", len(ids))
	}
	for _, id := range want {
		found := false
		for _, actual := range ids {
			if actual == id {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("RequiredScenarioIDs() missing %q", id)
		}
	}
}

func TestValidateRejectsDesktopQualificationMutations(t *testing.T) {
	t.Parallel()

	if len(RequiredScenarioIDs()) != 47 {
		t.Skip("desktop qualification scenarios are not implemented yet")
	}
	tests := []struct {
		name   string
		mutate func(*Evidence)
		want   string
	}{
		{"desktop needs routine user", func(e *Evidence) {
			findObservation(e, "desktop.interaction-states").EnvironmentID = "admin"
		}, "medium integrity"},
		{"desktop needs image", func(e *Evidence) {
			findObservation(e, "desktop.interaction-states").AttachmentPaths = []string{"attachments/fixture.txt"}
		}, "supported raster image"},
		{"appearance palettes", func(e *Evidence) {
			findObservation(e, "desktop.appearance-standard").Metrics["palettes"] = "dark"
		}, "palettes"},
		{"duplicate exact set", func(e *Evidence) {
			findObservation(e, "desktop.appearance-standard").Metrics["palettes"] = "dark,light,dark"
		}, "duplicate"},
		{"standard appearance dpi", func(e *Evidence) {
			findObservation(e, "desktop.appearance-standard").Metrics["effective_dpi"] = 120
		}, "effective_dpi"},
		{"scaled appearance dpi", func(e *Evidence) {
			findObservation(e, "desktop.appearance-scaled").Metrics["effective_dpi"] = 96
		}, "greater than 96"},
		{"standard visual scenario dpi", func(e *Evidence) {
			findObservation(e, "desktop.interaction-states").EnvironmentID = "high"
		}, "exactly 96 DPI"},
		{"scaled visual scenario dpi", func(e *Evidence) {
			findObservation(e, "desktop.interaction-states-scaled").EnvironmentID = "standard"
		}, "greater than 96 DPI"},
		{"appearance sharpness", func(e *Evidence) {
			findObservation(e, "desktop.appearance-standard").Metrics["body_text_sharp"] = false
		}, "body_text_sharp"},
		{"interaction states", func(e *Evidence) {
			findObservation(e, "desktop.interaction-states").Metrics["states"] = "rest,hover,focus"
		}, "states"},
		{"text contrast", func(e *Evidence) {
			findObservation(e, "desktop.interaction-states").Metrics["minimum_text_contrast"] = 4.49
		}, "minimum_text_contrast"},
		{"non-text contrast", func(e *Evidence) {
			findObservation(e, "desktop.interaction-states").Metrics["minimum_non_text_contrast"] = 2.99
		}, "minimum_non_text_contrast"},
		{"navigation order", func(e *Evidence) {
			findObservation(e, "desktop.navigation-options").Metrics["destination_order"] = "tasks,options,info"
		}, "destination_order"},
		{"navigation horizontal scrollbar", func(e *Evidence) {
			findObservation(e, "desktop.navigation-options").Metrics["horizontal_scrollbar_present"] = true
		}, "horizontal_scrollbar_present"},
		{"scroll sensitivities", func(e *Evidence) {
			findObservation(e, "desktop.scroll-input").Metrics["sensitivities"] = "1x,2x"
		}, "sensitivities"},
		{"touchpad unavailable reason", func(e *Evidence) {
			findObservation(e, "desktop.scroll-input").Metrics["touchpad_unavailable_reason"] = ""
		}, "touchpad_unavailable_reason"},
		{"touchpad fine deltas", func(e *Evidence) {
			metrics := findObservation(e, "desktop.scroll-input").Metrics
			metrics["touchpad_available"] = true
			metrics["touchpad_fine_deltas_preserved"] = false
		}, "touchpad_fine_deltas_preserved"},
		{"tasks population", func(e *Evidence) {
			findObservation(e, "desktop.tasks-table").Metrics["row_count"] = 99
		}, "row_count"},
		{"tasks headers", func(e *Evidence) {
			findObservation(e, "desktop.tasks-table").Metrics["headers"] = "task,enabled,lifecycle,time-zone"
		}, "headers"},
		{"schedule population", func(e *Evidence) {
			findObservation(e, "desktop.schedule-activity-tables").Metrics["schedule_row_count"] = 99
		}, "schedule_row_count"},
		{"activity severity casing", func(e *Evidence) {
			findObservation(e, "desktop.schedule-activity-tables").Metrics["severities"] = "info,warning,error"
		}, "severities"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, artifact, evidence := passingEvidence(t)
			tt.mutate(&evidence)
			if failures := Validate(evidence, root, artifact, ExpectedIdentity{}); !containsFailure(failures, tt.want) {
				t.Fatalf("Validate() failures = %v, want substring %q", failures, tt.want)
			}
		})
	}
}

func TestValidateRejectsDeclaredImageWithNonRasterBytes(t *testing.T) {
	t.Parallel()

	root, artifact, evidence := passingEvidence(t)
	data := []byte("ordinary text deliberately mislabeled as image/png\n")
	name := "attachments/fixture.svg"
	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(name)), data, 0o600); err != nil {
		t.Fatal(err)
	}
	for i := range evidence.Attachments {
		if evidence.Attachments[i].Path == name {
			evidence.Attachments[i].Bytes = int64(len(data))
			evidence.Attachments[i].SHA256 = digest(data)
			evidence.Attachments[i].MediaType = "image/png"
		}
	}

	failures := Validate(evidence, root, artifact, ExpectedIdentity{})
	if !containsFailure(failures, "supported raster image") {
		t.Fatalf("Validate() failures = %v", failures)
	}
}

func TestValidateRequiresStandardAndScaledDPIForEveryVisualFamily(t *testing.T) {
	t.Parallel()

	standardIDs := []string{
		"desktop.appearance-standard",
		"desktop.interaction-states",
		"desktop.navigation-options",
		"desktop.tasks-table",
		"desktop.schedule-activity-tables",
	}
	for _, id := range standardIDs {
		id := id
		t.Run(id, func(t *testing.T) {
			root, artifact, evidence := passingEvidence(t)
			findObservation(&evidence, id).EnvironmentID = "high"
			if failures := Validate(evidence, root, artifact, ExpectedIdentity{}); !containsFailure(failures, "exactly 96 DPI") {
				t.Fatalf("Validate() failures = %v", failures)
			}
		})
	}

	scaledIDs := []string{
		"desktop.appearance-scaled",
		"desktop.interaction-states-scaled",
		"desktop.navigation-options-scaled",
		"desktop.tasks-table-scaled",
		"desktop.schedule-activity-tables-scaled",
	}
	for _, id := range scaledIDs {
		id := id
		t.Run(id, func(t *testing.T) {
			root, artifact, evidence := passingEvidence(t)
			findObservation(&evidence, id).EnvironmentID = "standard"
			if failures := Validate(evidence, root, artifact, ExpectedIdentity{}); !containsFailure(failures, "greater than 96 DPI") {
				t.Fatalf("Validate() failures = %v", failures)
			}
		})
	}
}

func TestValidateRejectsCriticalMutations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(*Evidence)
		want   string
	}{
		{"schema", func(e *Evidence) { e.SchemaVersion = 2 }, "schema_version"},
		{"evidence class", func(e *Evidence) { e.EvidenceClass = "automated-fixture" }, "evidence_class"},
		{"repository", func(e *Evidence) { e.Candidate.Repository = "fork/repo" }, "repository"},
		{"commit", func(e *Evidence) { e.Candidate.Commit = strings.Repeat("b", 40) }, "commit"},
		{"workflow", func(e *Evidence) { e.Candidate.Workflow = "Other" }, "workflow"},
		{"outcome", func(e *Evidence) { e.Observations[0].Status = "unavailable" }, "unavailable"},
		{"ordered timestamps", func(e *Evidence) { e.Observations[0].CompletedAt = e.Observations[0].StartedAt.Add(-time.Second) }, "timestamps"},
		{"duplicate", func(e *Evidence) { e.Observations[1].ID = e.Observations[0].ID }, "duplicate"},
		{"duplicate environment", func(e *Evidence) { e.Environments[1].ID = e.Environments[0].ID }, "duplicate"},
		{"unknown environment", func(e *Evidence) { e.Observations[0].EnvironmentID = "missing" }, "unknown environment"},
		{"unknown attachment", func(e *Evidence) {
			findObservation(e, "error.timeout").AttachmentPaths = []string{"attachments/missing.svg"}
		}, "unknown attachment"},
		{"missing", func(e *Evidence) { e.Observations = e.Observations[1:] }, "missing required observation"},
		{"elevated routine user", func(e *Evidence) { e.Environments[0].Integrity = "high" }, "medium integrity"},
		{"token mismatch", func(e *Evidence) { e.Environments[0].IntegrityRID = 12288 }, "integrity_rid"},
		{"wrong service identity", func(e *Evidence) { e.Environments[0].ServiceIdentity = "LocalService" }, "LocalSystem"},
		{"server", func(e *Evidence) { e.Environments[0].WindowsEdition = "Windows Server 2025" }, "Windows 11"},
		{"short error interval", func(e *Evidence) { findObservation(e, "error.timeout").Metrics["duration_seconds"] = 119 }, "duration_seconds"},
		{"modal spam", func(e *Evidence) { findObservation(e, "error.daemon-unavailable").Metrics["max_modal_count"] = 1 }, "max_modal_count"},
		{"fractional sample count", func(e *Evidence) { findObservation(e, "error.timeout").Metrics["sample_count"] = 120.5 }, "integer"},
		{"retry exit unreachable", func(e *Evidence) { findObservation(e, "error.manual-retry").Metrics["exit_reachable"] = false }, "exit_reachable"},
		{"manual retry interval", func(e *Evidence) {
			findObservation(e, "error.manual-retry").CompletedAt = findObservation(e, "error.manual-retry").StartedAt.Add(119 * time.Second)
		}, "timestamp interval"},
		{"wrong error category", func(e *Evidence) { findObservation(e, "error.timeout").Metrics["category"] = "daemon-unavailable" }, "category"},
		{"oversize", func(e *Evidence) {
			metrics := findObservation(e, "window.clean-standard").Metrics
			metrics["logical_work_area_width"] = 1800
			metrics["logical_work_area_height"] = 800
			metrics["fyne_content_width"] = 1700
		}, "90 percent"},
		{"unsafe attachment", func(e *Evidence) { e.Attachments[0].Path = "../outside.txt" }, "unsafe"},
		{"attachment digest", func(e *Evidence) { e.Attachments[0].SHA256 = strings.Repeat("e", 64) }, "attachments[0] SHA-256"},
		{"missing visual evidence", func(e *Evidence) {
			findObservation(e, "error.timeout").AttachmentPaths = []string{"attachments/fixture.txt"}
		}, "supported raster image"},
		{"missing native window evidence", func(e *Evidence) {
			findObservation(e, "window.clean-standard").AttachmentPaths = []string{"attachments/fixture.svg"}
		}, "native window measurement"},
		{"native metric disagreement", func(e *Evidence) { findObservation(e, "window.clean-standard").Metrics["hwnd"] = "0x00000002" }, "does not match native window attachment"},
		{"minimized window", func(e *Evidence) { findObservation(e, "window.clean-standard").Metrics["minimized"] = true }, "minimized"},
		{"missing outer margin", func(e *Evidence) {
			findObservation(e, "window.clean-standard").Metrics["outer_rect"] = rect(0, 0, 2560, 1400)
		}, "positive margins"},
		{"missing close proof", func(e *Evidence) { findObservation(e, "window.state-transitions").Metrics["close_worked"] = false }, "close_worked"},
		{"reused subsequent process", func(e *Evidence) { findObservation(e, "window.subsequent-launch").Metrics["pid"] = 100 }, "fresh process id"},
		{"reused task run", func(e *Evidence) {
			findObservation(e, "task.scheduled-success").Metrics["run_id"] = findObservation(e, "task.manual-success").Metrics["run_id"]
		}, "distinct run_id"},
		{"missing task evidence", func(e *Evidence) { findObservation(e, "task.manual-success").AttachmentPaths = nil }, "task run evidence"},
		{"task output mismatch", func(e *Evidence) {
			findObservation(e, "task.manual-success").Metrics["output_sha256"] = strings.Repeat("9", 64)
		}, "does not match native window attachment"},
		{"start failure exit", func(e *Evidence) { findObservation(e, "task.start-failure").Metrics["expected_exit_code"] = 0 }, "expected_exit_code"},
		{"shortcut default", func(e *Evidence) { findObservation(e, "setup.shortcut-defaults").Metrics["start_menu_default"] = false }, "start_menu_default"},
		{"completion matrix", func(e *Evidence) { findObservation(e, "setup.completion-matrix").Metrics["combinations_verified"] = 3 }, "combinations_verified"},
		{"finish integrity", func(e *Evidence) {
			findObservation(e, "setup.finish-launch-integrity").Metrics["process_integrity"] = "high"
		}, "process_integrity"},
		{"missing setup fingerprint", func(e *Evidence) { delete(findObservation(e, "setup.upgrade").Metrics, "before_fingerprint") }, "before_fingerprint"},
		{"preserve bytes", func(e *Evidence) { findObservation(e, "remove.preserve").Metrics["owned_bytes_preserved"] = false }, "owned_bytes_preserved"},
		{"wipe controls", func(e *Evidence) { findObservation(e, "remove.wipe").Metrics["controls_unchanged"] = false }, "controls_unchanged"},
		{"changed control fingerprint", func(e *Evidence) {
			findObservation(e, "remove.wipe").Metrics["controls_after_sha256"] = strings.Repeat("8", 64)
		}, "control fingerprints"},
		{"one profile", func(e *Evidence) { findObservation(e, "remove.multiple-profiles").Metrics["profile_count"] = 1 }, "profile_count"},
		{"locked result", func(e *Evidence) { findObservation(e, "remove.locked-partial").Metrics["cleanup_result"] = "complete" }, "cleanup_result"},
		{"wipe reinstall retained", func(e *Evidence) {
			findObservation(e, "remove.reinstall-after-wipe").Metrics["prior_tasks_available"] = true
		}, "prior_tasks_available"},
		{"preserve config", func(e *Evidence) {
			findObservation(e, "remove.reinstall-after-preserve").Metrics["prior_config_available"] = false
		}, "prior_config_available"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, artifact, evidence := passingEvidence(t)
			tt.mutate(&evidence)
			failures := Validate(evidence, root, artifact, ExpectedIdentity{
				Repository: "shruggietech/go-schedule",
				Tag:        "v1.0.0",
				Commit:     strings.Repeat("a", 40),
			})
			if !containsFailure(failures, tt.want) {
				t.Fatalf("Validate() failures = %v, want substring %q", failures, tt.want)
			}
		})
	}
}

func TestValidateRejectsEveryNonPassingOutcome(t *testing.T) {
	t.Parallel()

	for _, status := range []string{"fail", "unavailable", "skipped", "timed-out", "partial"} {
		t.Run(status, func(t *testing.T) {
			root, artifact, evidence := passingEvidence(t)
			evidence.Observations[0].Status = status
			failures := Validate(evidence, root, artifact, ExpectedIdentity{})
			if !containsFailure(failures, status) {
				t.Fatalf("Validate() failures = %v, want %q", failures, status)
			}
		})
	}
}

func TestValidateRejectsSymlinkAttachment(t *testing.T) {
	t.Parallel()

	root, artifact, evidence := passingEvidence(t)
	target := filepath.Join(root, "attachments", "fixture.svg")
	link := filepath.Join(root, "attachments", "linked.svg")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}
	info, err := os.Stat(target)
	if err != nil {
		t.Fatal(err)
	}
	bytes, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	evidence.Attachments = append(evidence.Attachments, Attachment{
		Path: "attachments/linked.svg", Bytes: info.Size(), SHA256: digest(bytes),
		MediaType: "image/svg+xml", Purpose: "must be rejected",
	})
	findObservation(&evidence, "error.timeout").AttachmentPaths = []string{"attachments/linked.svg"}
	failures := Validate(evidence, root, artifact, ExpectedIdentity{})
	if !containsFailure(failures, "symbolic link") {
		t.Fatalf("Validate() failures = %v", failures)
	}
}

func TestValidateRejectsSymlinkAttachmentParent(t *testing.T) {
	t.Parallel()

	root, artifact, evidence := passingEvidence(t)
	external := t.TempDir()
	data := []byte("outside evidence")
	if err := os.WriteFile(filepath.Join(external, "outside.svg"), data, 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "attachments", "linked")
	if err := os.Symlink(external, link); err != nil {
		t.Skipf("directory symlink creation unavailable: %v", err)
	}
	evidence.Attachments = append(evidence.Attachments, Attachment{
		Path: "attachments/linked/outside.svg", Bytes: int64(len(data)), SHA256: digest(data),
		MediaType: "image/svg+xml", Purpose: "must be rejected",
	})
	findObservation(&evidence, "error.timeout").AttachmentPaths = []string{"attachments/linked/outside.svg"}
	failures := Validate(evidence, root, artifact, ExpectedIdentity{})
	if !containsFailure(failures, "symbolic link") {
		t.Fatalf("Validate() failures = %v", failures)
	}
}

func TestValidateRejectsChangedCandidateBytes(t *testing.T) {
	t.Parallel()

	root, artifact, evidence := passingEvidence(t)
	file, err := os.OpenFile(artifact, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.WriteString("changed"); err != nil {
		file.Close()
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	failures := Validate(evidence, root, artifact, ExpectedIdentity{})
	if !containsFailure(failures, "artifact SHA-256") {
		t.Fatalf("Validate() failures = %v", failures)
	}
}

func TestDecodeEvidenceRejectsUnknownFieldsAndTrailingJSON(t *testing.T) {
	t.Parallel()

	for _, input := range []string{
		`{"schema_version":1,"unknown":true}`,
		`{"schema_version":1}{"schema_version":1}`,
		`{"schema_version":1,"observations":[{"status":"fail","status":"pass"}]}`,
		"\xef\xbb\xbf{}",
		"{\"schema_version\":1,\"operator\":{\"role\":\"\xff\"}}",
	} {
		if _, err := DecodeEvidence(strings.NewReader(input)); err == nil {
			t.Fatalf("DecodeEvidence(%q) unexpectedly succeeded", input)
		}
	}
}

func TestDecodeCandidateRejectsDuplicateKeysAndInvalidUTF8(t *testing.T) {
	t.Parallel()

	for _, input := range []string{
		`{"repository":"first","repository":"second"}`,
		"{\"repository\":\"\xff\"}",
	} {
		if _, err := DecodeCandidate(strings.NewReader(input)); err == nil {
			t.Fatalf("DecodeCandidate(%q) unexpectedly succeeded", input)
		}
	}
}

func TestCheckedInFixtureIsCompleteButExplicitlyNonNative(t *testing.T) {
	t.Parallel()

	root := filepath.Join("..", "..", "test", "fixtures", "windows-release-gate", "passing")
	file, err := os.Open(filepath.Join(root, "evidence.json"))
	if err != nil {
		t.Fatal(err)
	}
	evidence, err := DecodeEvidence(file)
	file.Close()
	if err != nil {
		t.Fatal(err)
	}
	if failures := ValidateFixture(evidence, root, filepath.Join(root, evidence.Candidate.Filename), ExpectedIdentity{
		Repository: "shruggietech/go-schedule",
		Tag:        "v1.0.0",
		Commit:     strings.Repeat("a", 40),
	}); len(failures) != 0 {
		t.Fatalf("checked-in fixture failures = %v", failures)
	}
	if failures := Validate(evidence, root, filepath.Join(root, evidence.Candidate.Filename), ExpectedIdentity{}); !containsFailure(failures, "evidence_class") {
		t.Fatalf("production validation unexpectedly accepted fixture evidence: %v", failures)
	}
	for _, observation := range evidence.Observations {
		if !strings.Contains(observation.Summary, "fixture") || strings.Contains(observation.Summary, "native proof") {
			t.Fatalf("observation %s is not clearly labeled fixture data", observation.ID)
		}
	}
}

func TestValidateCandidateManifestRejectsIndependentIdentityDrift(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(*Candidate)
	}{
		{"repository", func(c *Candidate) { c.Repository = "other/repo" }},
		{"tag", func(c *Candidate) { c.Tag = "v9.9.9" }},
		{"commit", func(c *Candidate) { c.Commit = strings.Repeat("b", 40) }},
		{"workflow", func(c *Candidate) { c.Workflow = "Other" }},
		{"run_id", func(c *Candidate) { c.RunID++ }},
		{"run_attempt", func(c *Candidate) { c.RunAttempt++ }},
		{"filename", func(c *Candidate) { c.Filename = "other.msi" }},
		{"bytes", func(c *Candidate) { c.Bytes++ }},
		{"sha256", func(c *Candidate) { c.SHA256 = strings.Repeat("f", 64) }},
		{"product_version", func(c *Candidate) { c.ProductVersion = "9.9.9" }},
		{"product_code", func(c *Candidate) { c.ProductCode = "{AAAAAAAA-BBBB-CCCC-DDDD-EEEEEEEEEEEE}" }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, evidence := passingEvidence(t)
			manifest := evidence.Candidate
			tt.mutate(&manifest)
			failures := ValidateCandidateManifest(evidence.Candidate, manifest)
			if !containsFailure(failures, tt.name) {
				t.Fatalf("ValidateCandidateManifest() failures = %v", failures)
			}
		})
	}
}

func TestValidateCandidateIsIndependentOfAttendedEvidence(t *testing.T) {
	t.Parallel()

	_, artifact, evidence := passingEvidence(t)
	if failures := ValidateCandidate(evidence.Candidate, artifact, ExpectedIdentity{
		Repository: "shruggietech/go-schedule",
		Tag:        "v1.0.0",
		Commit:     strings.Repeat("a", 40),
	}); len(failures) != 0 {
		t.Fatalf("ValidateCandidate() failures = %v", failures)
	}
	evidence.Candidate.SHA256 = strings.Repeat("f", 64)
	if failures := ValidateCandidate(evidence.Candidate, artifact, ExpectedIdentity{}); !containsFailure(failures, "artifact SHA-256") {
		t.Fatalf("ValidateCandidate() failures = %v", failures)
	}
}

func passingEvidence(t *testing.T) (string, string, Evidence) {
	t.Helper()

	root := t.TempDir()
	artifact := filepath.Join(root, "go-schedule_v1.0.0_windows_amd64.msi")
	artifactBytes := []byte("inert fixture; not a native MSI")
	if err := os.WriteFile(artifact, artifactBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	attachmentPath := filepath.Join(root, "attachments", "fixture.txt")
	if err := os.MkdirAll(filepath.Dir(attachmentPath), 0o700); err != nil {
		t.Fatal(err)
	}
	attachmentBytes := []byte("fixture evidence; not a native observation")
	if err := os.WriteFile(attachmentPath, attachmentBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	nativePath := filepath.Join(root, "attachments", "fixture.json")
	nativeBytes := []byte("{\"fixture\":true}\n")
	if err := os.WriteFile(nativePath, nativeBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	visualPath := filepath.Join(root, "attachments", "fixture.svg")
	visualBytes := []byte("P3\n1 1\n255\n36 90 114\n")
	if err := os.WriteFile(visualPath, visualBytes, 0o600); err != nil {
		t.Fatal(err)
	}

	start := time.Date(2026, 9, 2, 12, 0, 0, 0, time.UTC)
	end := start.Add(4 * time.Hour)
	evidence := Evidence{
		SchemaVersion: 1,
		EvidenceClass: "attended-windows",
		Candidate: Candidate{
			Repository:     "shruggietech/go-schedule",
			Tag:            "v1.0.0",
			Commit:         strings.Repeat("a", 40),
			Workflow:       "Release",
			RunID:          1234,
			RunAttempt:     1,
			Filename:       filepath.Base(artifact),
			Bytes:          int64(len(artifactBytes)),
			SHA256:         digest(artifactBytes),
			ProductVersion: "1.0.0",
			ProductCode:    "{11111111-2222-3333-4444-555555555555}",
		},
		Operator: Operator{
			Role:       "release maintainer",
			AttestedAt: end,
			Statement:  OperatorAttestation,
		},
		StartedAt:   start,
		CompletedAt: end,
		Environments: []Environment{
			environment("standard", "intended-user", "medium", "standard-dpi", "clean", 96),
			environment("high", "intended-user", "medium", "high-dpi", "clean", 144),
			environment("retained", "intended-user", "medium", "standard-dpi", "retained-v0.9.1", 96),
			environment("unrelated", "unrelated-user", "medium", "not-applicable", "not-applicable", 0),
			environment("admin", "administrator", "high", "not-applicable", "not-applicable", 0),
		},
		Attachments: []Attachment{
			{Path: "attachments/fixture.txt", Bytes: int64(len(attachmentBytes)), SHA256: digest(attachmentBytes), MediaType: "text/plain", Purpose: "automated fixture only"},
			{Path: "attachments/fixture.json", Bytes: int64(len(nativeBytes)), SHA256: digest(nativeBytes), MediaType: "application/json", Purpose: "automated native-metric fixture only"},
			{Path: "attachments/fixture.svg", Bytes: int64(len(visualBytes)), SHA256: digest(visualBytes), MediaType: "image/x-portable-pixmap", Purpose: "automated visual fixture only"},
		},
	}

	for i, id := range RequiredScenarioIDs() {
		environmentID := "standard"
		switch id {
		case "window.clean-high-or-mixed":
			environmentID = "high"
		case "desktop.appearance-scaled", "desktop.interaction-states-scaled", "desktop.navigation-options-scaled", "desktop.tasks-table-scaled", "desktop.schedule-activity-tables-scaled":
			environmentID = "high"
		case "window.retained-profile":
			environmentID = "retained"
		case "access.unrelated-user-denied":
			environmentID = "unrelated"
		case "remove.wipe", "remove.locked-partial":
			environmentID = "admin"
		}
		observation := Observation{
			ID:              id,
			EnvironmentID:   environmentID,
			Status:          "pass",
			StartedAt:       start.Add(time.Duration(i) * time.Minute),
			CompletedAt:     start.Add(time.Duration(i+1) * time.Minute),
			Summary:         "automated fixture only; no native claim",
			Metrics:         passingMetrics(id),
			AttachmentPaths: nil,
		}
		if strings.HasPrefix(id, "error.") && id != "error.recovery" {
			observation.CompletedAt = observation.StartedAt.Add(120 * time.Second)
		}
		if strings.HasPrefix(id, "window.") {
			observation.AttachmentPaths = []string{"attachments/fixture.json", "attachments/fixture.svg"}
		} else if strings.HasPrefix(id, "error.") || strings.HasPrefix(id, "setup.") || strings.HasPrefix(id, "remove.") || strings.HasPrefix(id, "desktop.") {
			observation.AttachmentPaths = []string{"attachments/fixture.svg"}
		}
		evidence.Observations = append(evidence.Observations, observation)
	}
	addNativeWindowFixtures(t, root, &evidence)
	addTaskEvidenceFixture(t, root, &evidence)

	return root, artifact, evidence
}

func addNativeWindowFixtures(t *testing.T, root string, evidence *Evidence) {
	t.Helper()
	for i := range evidence.Observations {
		o := &evidence.Observations[i]
		if !strings.HasPrefix(o.ID, "window.") {
			continue
		}
		dpi, _ := numberMetric(o.Metrics, "effective_dpi")
		pid, _ := numberMetric(o.Metrics, "pid")
		session, _ := numberMetric(o.Metrics, "process_session_id")
		rid, _ := numberMetric(o.Metrics, "process_integrity_rid")
		record := nativeWindowEvidence{
			SchemaVersion: 1, Kind: "native-window-v1", ObservationID: o.ID,
			CapturedAt: evidence.StartedAt.Add(time.Minute), ProcessID: int(pid),
			ProcessPath:   `C:\Program Files\go-schedule\gosched-gui.exe`,
			ProcessSHA256: strings.Repeat("b", 64), ProcessSessionID: int(session),
			ProcessUserSID: "S-1-5-21-1000", ProcessIntegrityRID: int(rid),
			HWND: "0x00000001", OuterRect: nativeRect{100, 100, 1400, 940},
			ClientRect: nativeRect{110, 130, 1390, 930}, MonitorRect: nativeRect{0, 0, 2560, 1440},
			WorkAreaRect: nativeRect{0, 0, 2560, 1400}, MonitorID: "fixture-monitor", EffectiveDPI: int(dpi),
			ShowCommand: 1, Visible: true, Restored: true,
			Fyne: fyneWindowEvidence{SchemaVersion: 1, ProcessID: int(pid), CapturedAt: evidence.StartedAt.Add(time.Minute), ContentWidth: 1280, ContentHeight: 800, CanvasScale: dpi / 96},
		}
		data, err := json.MarshalIndent(record, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		data = append(data, '\n')
		name := "attachments/native-" + strings.ReplaceAll(o.ID, ".", "-") + ".json"
		full := filepath.Join(root, filepath.FromSlash(name))
		if err := os.WriteFile(full, data, 0o600); err != nil {
			t.Fatal(err)
		}
		evidence.Attachments = append(evidence.Attachments, Attachment{
			Path: name, Bytes: int64(len(data)), SHA256: digest(data), MediaType: "application/json", Purpose: "native window measurement",
		})
		o.AttachmentPaths = []string{name, "attachments/fixture.svg"}
	}
}

func addTaskEvidenceFixture(t *testing.T, root string, evidence *Evidence) {
	t.Helper()
	bundle := taskEvidenceBundle{SchemaVersion: 1, Kind: "task-run-evidence-v1"}
	for i := range evidence.Observations {
		o := &evidence.Observations[i]
		if !strings.HasPrefix(o.ID, "task.") {
			continue
		}
		expected, _ := numberMetric(o.Metrics, "expected_exit_code")
		actual, _ := numberMetric(o.Metrics, "actual_exit_code")
		bundle.Runs = append(bundle.Runs, taskRunEvidence{
			ObservationID: o.ID, TaskID: "fixture-task", RunID: "fixture-" + o.ID,
			InvocationMode: o.Metrics["invocation_mode"].(string), PublicInterface: "gosched CLI",
			TaskDefinition: "fixture definition " + o.ID, Output: "fixture output " + o.ID,
			Marker: "fixture marker " + o.ID, History: "fixture history " + o.ID,
			ProductionRun: true, ExpectedExit: int(expected), ActualExit: int(actual),
			Diagnostic: stringMetric(o.Metrics, "diagnostic_category"),
		})
	}
	data, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, '\n')
	name := "attachments/task-runs.json"
	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(name)), data, 0o600); err != nil {
		t.Fatal(err)
	}
	evidence.Attachments = append(evidence.Attachments, Attachment{
		Path: name, Bytes: int64(len(data)), SHA256: digest(data), MediaType: "application/json", Purpose: "task run evidence",
	})
	for i := range evidence.Observations {
		if strings.HasPrefix(evidence.Observations[i].ID, "task.") {
			evidence.Observations[i].AttachmentPaths = []string{name}
		}
	}
}

func stringMetric(metrics map[string]any, key string) string {
	value, _ := metrics[key].(string)
	return value
}

func environment(id, role, integrity, display, profile string, dpi int) Environment {
	return Environment{
		ID:              id,
		Snapshot:        "fixture-clean-snapshot",
		CleanSnapshot:   true,
		WindowsEdition:  "Windows 11 Pro",
		WindowsBuild:    "fixture-build",
		AccountRole:     role,
		AccountSID:      "S-1-5-21-1000",
		Integrity:       integrity,
		IntegrityRID:    map[string]int{"medium": 8192, "high": 12288, "system": 16384}[integrity],
		ServiceIdentity: "LocalSystem",
		DisplayClass:    display,
		EffectiveDPI:    dpi,
		ProfileState:    profile,
	}
}

func passingMetrics(id string) map[string]any {
	m := map[string]any{"verified": true}
	switch id {
	case "desktop.appearance-standard", "desktop.appearance-scaled":
		dpi := 96
		if id == "desktop.appearance-scaled" {
			dpi = 144
		}
		m = map[string]any{
			"palettes": "dark,light", "effective_dpi": dpi,
			"system_font_default": true, "system_font_restored": true,
			"font_persistence_verified": true, "info_text_sharp": true,
			"body_text_sharp": true, "labels_centered": true,
			"labels_unclipped": true, "resize_verified": true,
			"minimize_restore_verified": true, "reopen_verified": true,
			"fonts_exercised": "system,geist,inter,ubuntu,monospace",
		}
	case "desktop.interaction-states", "desktop.interaction-states-scaled":
		m = map[string]any{
			"palettes":              "dark,light",
			"control_families":      "navigation,selector,ordinary,primary,danger,dialog,table-row",
			"states":                "rest,hover,focus,pressed,selected,disabled",
			"minimum_text_contrast": 4.5, "minimum_non_text_contrast": 3.0,
			"labels_readable": true, "glyphs_readable": true,
			"selection_identifiable": true, "focus_visible": true,
			"non_color_cues_present": true,
		}
	case "desktop.navigation-options", "desktop.navigation-options-scaled":
		m = map[string]any{
			"palettes": "dark,light", "content_sizes": "1280x800,800x600",
			"destination_order":     "tasks,groups,chains,schedule,activity,options,info",
			"rail_spacing_balanced": true, "labels_unclipped": true,
			"boundary_full_height": true, "boundary_subtle": true,
			"exit_bottom_right": true, "exit_never_selected": true,
			"exit_semantic_glyph": true, "storage_rows_compact": true,
			"unavailable_rows_muted": true, "copy_exact": true,
			"selector_current_omitted": true, "horizontal_scrollbar_present": false,
		}
	case "desktop.scroll-input":
		m = map[string]any{
			"sensitivities":            "1x,2x,4x",
			"surfaces":                 "options,info,editor-command,editor-schedule,editor-help",
			"wheel_detents_responsive": true, "immediate_apply": true,
			"persistence_verified": true, "nested_multiplier_absent": true,
			"keyboard_scroll_preserved": true, "touchpad_available": false,
			"touchpad_fine_deltas_preserved": false,
			"touchpad_unavailable_reason":    "fixture has no physical input hardware",
		}
	case "desktop.tasks-table", "desktop.tasks-table-scaled":
		m = map[string]any{
			"row_count": 100, "palettes": "dark,light",
			"content_sizes":  "1280x800,800x600",
			"headers":        "task,enabled,lifecycle,time-zone,group",
			"row_states":     "odd,even,hover,focus,selected",
			"headers_frozen": true, "status_dimensions_distinct": true,
			"bracket_decoration_absent": true, "full_values_discoverable": true,
			"horizontal_scrollbar_present": false, "refresh_identity_stable": true,
			"removed_selection_clears": true, "toolbar_actions_work": true,
			"double_click_edits": true,
		}
	case "desktop.schedule-activity-tables", "desktop.schedule-activity-tables-scaled":
		m = map[string]any{
			"schedule_row_count": 100, "activity_row_count": 100,
			"palettes": "dark,light", "content_sizes": "1280x800,800x600",
			"schedule_headers": "when,task,event,outcome",
			"activity_headers": "when,severity,source,summary",
			"schedule_states":  "scheduled,success,failure,skipped,caught-up,queued,missing,unknown",
			"severities":       "INFO,WARNING,ERROR",
			"row_states":       "odd,even,hover,focus,selected",
			"headers_frozen":   true, "semantic_text_glyphs_match": true,
			"non_color_cues_present": true, "full_values_discoverable": true,
			"horizontal_scrollbar_present": false, "refresh_identity_stable": true,
			"removed_selection_clears": true, "detail_activation_accurate": true,
			"range_calendar_switching": true, "filter_clear_acknowledge": true,
		}
	case "access.intended-user":
		m = map[string]any{"health_ok": true, "gui_task_list_ok": true, "routine_elevation_required": false}
	case "access.unrelated-user-denied":
		m = map[string]any{"pipe_opened": false, "error_kind": "access-denied"}
	case "access.fresh-path-resolution":
		m = map[string]any{"fresh_process": true, "matches_installed_cli": true, "resolved_path": `C:\Program Files\go-schedule\gosched.exe`}
	case "access.path-removed":
		m = map[string]any{"fresh_process": true, "resolves": false, "registry_cardinality": 0}
	case "setup.shortcut-defaults":
		m = map[string]any{"start_menu_default": true, "desktop_default": false, "defaults_visible": true, "effects_verified": true}
	case "setup.shortcut-matrix":
		m = map[string]any{"combinations_verified": 4, "targets_verified": true}
	case "setup.completion-matrix":
		m = map[string]any{"combinations_verified": 4, "independent_choices": true, "default_handler_verified": true, "launch_default": true, "documentation_default": false}
	case "setup.finish-launch-integrity":
		m = map[string]any{"process_integrity": "medium", "launch_count": 1}
	case "setup.cancel":
		m = map[string]any{"state_unchanged": true, "owned_data_cleanup_invoked": false}
	case "setup.maintenance":
		m = map[string]any{"transitions_verified": true, "repair_verified": true, "completion_actions_absent": true, "owned_data_cleanup_invoked": false}
	case "setup.upgrade":
		m = map[string]any{"choices_preserved": true, "completion_actions_absent": true, "owned_data_cleanup_invoked": false}
	case "setup.invalid-input":
		m = map[string]any{"input_rejected": true, "state_unchanged": true, "owned_data_cleanup_invoked": false}
	case "setup.rollback":
		m = map[string]any{"rollback_completed": true, "state_unchanged": true, "owned_data_cleanup_invoked": false}
	case "remove.preserve":
		m = map[string]any{"mode": "preserve", "software_removed": true, "owned_bytes_preserved": true, "controls_unchanged": true, "preserve_default_visible": true, "owned_inventory_reviewed": true, "owned_data_cleanup_invoked": false}
	case "remove.wipe":
		m = map[string]any{"mode": "wipe", "software_removed": true, "owned_roots_removed": true, "controls_unchanged": true, "security_state_preserved": true, "wipe_explicitly_selected": true, "wipe_confirmed": true, "owned_inventory_reviewed": true, "owned_data_cleanup_invoked": true}
	case "remove.cancel":
		m = map[string]any{"software_unchanged": true, "data_unchanged": true, "owned_data_cleanup_invoked": false}
	case "remove.multiple-profiles":
		m = map[string]any{"profile_count": 2, "all_profiles_accounted": true}
	case "remove.locked-partial":
		m = map[string]any{"cleanup_result": "partial", "residual_count": 1, "truthfully_reported": true}
	case "remove.reinstall-after-preserve":
		m = map[string]any{"prior_tasks_available": true, "prior_preferences_available": true, "prior_config_available": true, "prior_logs_available": true}
	case "remove.reinstall-after-wipe":
		m = map[string]any{"prior_tasks_available": false, "prior_preferences_available": false, "prior_config_available": false, "prior_logs_available": false}
	}
	if strings.HasPrefix(id, "window.") {
		m = map[string]any{
			"pid": 100, "executable_sha256": strings.Repeat("b", 64), "hwnd": "0x00000001",
			"process_session_id": 1, "process_user_sid": "S-1-5-21-1000", "process_integrity_rid": 8192,
			"outer_rect": rect(100, 100, 1400, 940), "client_rect": rect(110, 130, 1390, 930),
			"monitor_rect": rect(0, 0, 2560, 1440), "work_area_rect": rect(0, 0, 2560, 1400),
			"fyne_content_width": 1280, "fyne_content_height": 800, "fyne_scale": 1,
			"logical_work_area_width": 2560, "logical_work_area_height": 1400,
			"effective_dpi": 96, "monitor_id": "fixture-monitor", "restored": true,
			"maximized": false, "minimized": false, "fullscreen": false,
			"margins_visible": true, "title_bar_reachable": true,
			"resize_borders_reachable": true, "taskbar_reachable": true,
		}
	}
	if id == "window.clean-standard" {
		m["launch_sequence"] = 1
	}
	if id == "window.clean-high-or-mixed" {
		m["effective_dpi"] = 144
		m["fyne_scale"] = 1.5
		m["logical_work_area_width"] = 2560.0 / 1.5
		m["logical_work_area_height"] = 1400.0 / 1.5
	}
	if id == "window.state-transitions" {
		m["maximize_worked"] = true
		m["restore_worked"] = true
		m["resize_worked"] = true
		m["minimize_worked"] = true
		m["close_worked"] = true
	}
	if id == "window.subsequent-launch" {
		m["pid"] = 101
		m["prior_process_id"] = 100
		m["launch_sequence"] = 2
		m["fresh_process"] = true
		m["prior_process_closed"] = true
	}
	if strings.HasPrefix(id, "error.") {
		m = map[string]any{
			"category": strings.TrimPrefix(id, "error."), "duration_seconds": 120,
			"trigger":      "fixture-" + strings.TrimPrefix(id, "error."),
			"sample_count": 120, "max_in_frame_count": 1, "max_modal_count": 0,
			"max_additional_top_level_count": 0, "retry_reachable": true,
			"exit_reachable": true, "guidance_recommends_elevation": false,
		}
		if id == "error.manual-retry" {
			m["retry_invoked"] = true
		}
		if id == "error.recovery" {
			m["incident_cleared"] = true
			m["interface_restored"] = true
			m["reinstall_required"] = false
			m["recovery_method"] = "daemon-restored"
		}
	}
	if strings.HasPrefix(id, "task.") {
		m = map[string]any{
			"public_interface": "gosched CLI", "production_runner": true,
			"task_id": "fixture-task", "run_id": "fixture-" + id, "history_recorded": true,
		}
		switch id {
		case "task.manual-success", "task.scheduled-success":
			if id == "task.manual-success" {
				m["invocation_mode"] = "manual"
			} else {
				m["invocation_mode"] = "scheduled"
			}
			m["expected_exit_code"] = 0
			m["actual_exit_code"] = 0
			m["output_sha256"] = strings.Repeat("c", 64)
			m["marker_sha256"] = strings.Repeat("d", 64)
		case "task.nonzero-exit":
			m["invocation_mode"] = "manual-fault"
			m["trigger"] = "fixture-nonzero"
			m["expected_exit_code"] = 7
			m["actual_exit_code"] = 7
			m["diagnostic_category"] = "nonzero-exit"
		case "task.start-failure":
			m["invocation_mode"] = "manual-fault"
			m["trigger"] = "fixture-missing-executable"
			m["diagnostic_category"] = "process-start-failure"
			m["expected_exit_code"] = -1
		}
		m["task_definition_sha256"] = digest([]byte("fixture definition " + id))
		m["output_sha256"] = digest([]byte("fixture output " + id))
		m["marker_sha256"] = digest([]byte("fixture marker " + id))
		m["history_sha256"] = digest([]byte("fixture history " + id))
	}
	if strings.HasPrefix(id, "setup.") {
		m["owned_data_cleanup_invoked"] = false
		m["install_session_id"] = "fixture-" + id
		m["installer_process_id"] = 4242
		m["installer_session_id"] = 1
		m["installer_process_owner_role"] = "administrator"
		m["installer_process_integrity"] = "high"
		m["selected_options_sha256"] = strings.Repeat("1", 64)
		m["observed_targets_sha256"] = strings.Repeat("2", 64)
		m["before_fingerprint"] = strings.Repeat("3", 64)
		m["after_fingerprint"] = strings.Repeat("4", 64)
		if id == "setup.cancel" || id == "setup.invalid-input" || id == "setup.rollback" {
			m["after_fingerprint"] = m["before_fingerprint"]
		}
	}
	if strings.HasPrefix(id, "remove.") {
		m["owned_roots_count"] = 2
		m["before_content_sha256"] = strings.Repeat("5", 64)
		m["after_content_sha256"] = strings.Repeat("6", 64)
		m["controls_before_sha256"] = strings.Repeat("7", 64)
		m["controls_after_sha256"] = strings.Repeat("7", 64)
		m["security_state_disposition"] = "preserved"
		m["reinstall_result"] = "verified"
		if id == "remove.preserve" || id == "remove.cancel" {
			m["after_content_sha256"] = m["before_content_sha256"]
		}
		switch id {
		case "remove.preserve", "remove.reinstall-after-preserve":
			m["reinstall_result"] = "preserved-state-restored"
		case "remove.wipe", "remove.reinstall-after-wipe":
			m["reinstall_result"] = "clean-state"
		case "remove.cancel":
			m["reinstall_result"] = "not-run"
		}
	}
	return m
}

func rect(left, top, right, bottom int) map[string]any {
	return map[string]any{"left": left, "top": top, "right": right, "bottom": bottom}
}

func findObservation(e *Evidence, id string) *Observation {
	for i := range e.Observations {
		if e.Observations[i].ID == id {
			return &e.Observations[i]
		}
	}
	panic("missing fixture observation " + id)
}

func containsFailure(failures []string, want string) bool {
	for _, failure := range failures {
		if strings.Contains(strings.ToLower(failure), strings.ToLower(want)) {
			return true
		}
	}
	return false
}

func digest(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
