package releasegate

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

var (
	tagPattern         = regexp.MustCompile(`^v([0-9]+)\.([0-9]+)\.([0-9]+)$`)
	commitPattern      = regexp.MustCompile(`^[0-9a-f]{40}$`)
	digestPattern      = regexp.MustCompile(`^[0-9a-f]{64}$`)
	productCodePattern = regexp.MustCompile(`^\{[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12}\}$`)
	idPattern          = regexp.MustCompile(`^[a-z0-9]+(?:[.-][a-z0-9]+)*$`)
)

var requiredScenarios = []string{
	"access.intended-user", "access.unrelated-user-denied", "access.fresh-path-resolution", "access.path-removed",
	"window.clean-standard", "window.clean-high-or-mixed", "window.retained-profile", "window.state-transitions", "window.subsequent-launch",
	"error.daemon-unavailable", "error.access-denied", "error.timeout", "error.stream-disconnect", "error.repeated-refresh", "error.manual-retry", "error.recovery",
	"task.manual-success", "task.scheduled-success", "task.nonzero-exit", "task.start-failure",
	"setup.shortcut-defaults", "setup.shortcut-matrix", "setup.completion-matrix", "setup.finish-launch-integrity", "setup.cancel", "setup.maintenance", "setup.upgrade", "setup.invalid-input", "setup.rollback",
	"remove.preserve", "remove.wipe", "remove.cancel", "remove.multiple-profiles", "remove.locked-partial", "remove.reinstall-after-preserve", "remove.reinstall-after-wipe",
	"desktop.appearance-standard", "desktop.appearance-scaled", "desktop.interaction-states", "desktop.interaction-states-scaled", "desktop.navigation-options", "desktop.navigation-options-scaled", "desktop.scroll-input", "desktop.tasks-table", "desktop.tasks-table-scaled", "desktop.schedule-activity-tables", "desktop.schedule-activity-tables-scaled",
}

var validStatuses = map[string]bool{
	"pass": true, "fail": true, "unavailable": true, "skipped": true, "timed-out": true, "partial": true,
}

// RequiredScenarioIDs returns a copy of the canonical scenario identifiers.
func RequiredScenarioIDs() []string {
	result := make([]string, len(requiredScenarios))
	copy(result, requiredScenarios)
	return result
}

// Validate returns every independently discoverable release-readiness defect.
func Validate(e Evidence, root, artifactPath string, expected ExpectedIdentity) []string {
	return validateWithClass(e, root, artifactPath, expected, "attended-windows")
}

// ValidateFixture exercises all semantic rules for explicitly synthetic test
// data. It is deliberately separate from the promotion validation path.
func ValidateFixture(e Evidence, root, artifactPath string, expected ExpectedIdentity) []string {
	return validateWithClass(e, root, artifactPath, expected, "automated-fixture")
}

func validateWithClass(e Evidence, root, artifactPath string, expected ExpectedIdentity, requiredClass string) []string {
	v := validator{root: root, evidence: &e}
	v.validateRoot(requiredClass)
	v.validateCandidate(artifactPath, expected)
	v.validateEnvironments()
	v.validateAttachments()
	v.validateObservations()
	sort.Strings(v.failures)
	return v.failures
}

// ValidateCandidate verifies a staged candidate manifest independently of
// attended evidence.
func ValidateCandidate(candidate Candidate, artifactPath string, expected ExpectedIdentity) []string {
	evidence := Evidence{Candidate: candidate}
	v := validator{evidence: &evidence}
	v.validateCandidate(artifactPath, expected)
	sort.Strings(v.failures)
	return v.failures
}

// ValidateCandidateManifest compares independently staged identity with the
// identity embedded in attended evidence.
func ValidateCandidateManifest(evidence, manifest Candidate) []string {
	var failures []string
	checks := []struct {
		name  string
		left  any
		right any
	}{
		{"repository", evidence.Repository, manifest.Repository},
		{"tag", evidence.Tag, manifest.Tag},
		{"commit", evidence.Commit, manifest.Commit},
		{"workflow", evidence.Workflow, manifest.Workflow},
		{"run_id", evidence.RunID, manifest.RunID},
		{"run_attempt", evidence.RunAttempt, manifest.RunAttempt},
		{"filename", evidence.Filename, manifest.Filename},
		{"bytes", evidence.Bytes, manifest.Bytes},
		{"sha256", evidence.SHA256, manifest.SHA256},
		{"product_version", evidence.ProductVersion, manifest.ProductVersion},
		{"product_code", evidence.ProductCode, manifest.ProductCode},
	}
	for _, check := range checks {
		if fmt.Sprint(check.left) != fmt.Sprint(check.right) {
			failures = append(failures, fmt.Sprintf(
				"candidate manifest %s does not match attended evidence",
				check.name,
			))
		}
	}
	return failures
}

type validator struct {
	root         string
	evidence     *Evidence
	failures     []string
	environments map[string]Environment
	attachments  map[string]Attachment
}

func (v *validator) add(format string, args ...any) {
	v.failures = append(v.failures, fmt.Sprintf(format, args...))
}

func (v *validator) validateRoot(requiredClass string) {
	e := v.evidence
	if e.SchemaVersion != 1 {
		v.add("schema_version is %d; expected 1", e.SchemaVersion)
	}
	if e.EvidenceClass != requiredClass {
		v.add("evidence_class %q must be %q", e.EvidenceClass, requiredClass)
	}
	if e.StartedAt.IsZero() || e.CompletedAt.IsZero() || e.CompletedAt.Before(e.StartedAt) {
		v.add("evidence timestamps must be present and ordered")
	}
	if strings.TrimSpace(e.Operator.Role) == "" {
		v.add("operator.role must not be blank")
	}
	if e.Operator.Statement != OperatorAttestation {
		v.add("operator.statement does not match the required attestation")
	}
	if e.Operator.AttestedAt.IsZero() || e.Operator.AttestedAt.Before(e.StartedAt) || e.Operator.AttestedAt.After(e.CompletedAt) {
		v.add("operator.attested_at must fall within the evidence interval")
	}
}

func (v *validator) validateCandidate(artifactPath string, expected ExpectedIdentity) {
	c := v.evidence.Candidate
	if c.Repository == "" || (expected.Repository != "" && c.Repository != expected.Repository) {
		v.add("candidate repository %q does not match expected %q", c.Repository, expected.Repository)
	}
	match := tagPattern.FindStringSubmatch(c.Tag)
	if match == nil {
		v.add("candidate tag %q is not vMAJOR.MINOR.PATCH", c.Tag)
	} else if c.ProductVersion != strings.TrimPrefix(c.Tag, "v") {
		v.add("candidate product_version %q does not match tag %q", c.ProductVersion, c.Tag)
	}
	if expected.Tag != "" && c.Tag != expected.Tag {
		v.add("candidate tag %q does not match expected %q", c.Tag, expected.Tag)
	}
	if !commitPattern.MatchString(c.Commit) {
		v.add("candidate commit %q must be 40 lowercase hexadecimal characters", c.Commit)
	}
	if expected.Commit != "" && c.Commit != expected.Commit {
		v.add("candidate commit %q does not match expected %q", c.Commit, expected.Commit)
	}
	if c.Workflow != "Release" {
		v.add("candidate workflow %q must be canonical \"Release\"", c.Workflow)
	}
	if c.RunID <= 0 || c.RunAttempt <= 0 {
		v.add("candidate workflow, run_id, and run_attempt must identify the staging run")
	}
	expectedFilename := "go-schedule_" + c.Tag + "_windows_amd64.msi"
	if c.Filename != expectedFilename {
		v.add("candidate filename %q does not match canonical %q", c.Filename, expectedFilename)
	}
	if !digestPattern.MatchString(c.SHA256) {
		v.add("candidate sha256 is not a lowercase SHA-256 digest")
	}
	if !productCodePattern.MatchString(c.ProductCode) {
		v.add("candidate product_code %q is not a canonical uppercase braced GUID", c.ProductCode)
	}
	info, digest, err := hashRegularFile(artifactPath)
	if err != nil {
		v.add("candidate artifact: %v", err)
		return
	}
	if filepath.Base(artifactPath) != c.Filename {
		v.add("candidate artifact basename %q does not match %q", filepath.Base(artifactPath), c.Filename)
	}
	if info.Size() != c.Bytes {
		v.add("candidate artifact bytes are %d; evidence requires %d", info.Size(), c.Bytes)
	}
	if digest != c.SHA256 {
		v.add("candidate artifact SHA-256 %s does not match evidence %s", digest, c.SHA256)
	}
}

func (v *validator) validateEnvironments() {
	v.environments = make(map[string]Environment, len(v.evidence.Environments))
	for i, environment := range v.evidence.Environments {
		prefix := fmt.Sprintf("environments[%d]", i)
		if !idPattern.MatchString(environment.ID) {
			v.add("%s.id %q is invalid", prefix, environment.ID)
		}
		if _, exists := v.environments[environment.ID]; exists {
			v.add("%s.id %q is duplicate", prefix, environment.ID)
		} else {
			v.environments[environment.ID] = environment
		}
		if !environment.CleanSnapshot {
			v.add("%s must identify a clean snapshot", prefix)
		}
		if !strings.Contains(environment.WindowsEdition, "Windows 11") || strings.Contains(environment.WindowsEdition, "Server") {
			v.add("%s.windows_edition must identify Windows 11 client", prefix)
		}
		if strings.TrimSpace(environment.Snapshot) == "" || strings.TrimSpace(environment.WindowsBuild) == "" {
			v.add("%s snapshot and Windows build must not be blank", prefix)
		}
		if !oneOf(environment.AccountRole, "intended-user", "unrelated-user", "administrator") {
			v.add("%s.account_role %q is invalid", prefix, environment.AccountRole)
		}
		if !oneOf(environment.Integrity, "medium", "high", "system") {
			v.add("%s.integrity %q is invalid", prefix, environment.Integrity)
		}
		if !strings.HasPrefix(environment.AccountSID, "S-1-") {
			v.add("%s.account_sid must be a Windows SID", prefix)
		}
		expectedRID := map[string]int{"medium": 8192, "high": 12288, "system": 16384}[environment.Integrity]
		if expectedRID != 0 && environment.IntegrityRID != expectedRID {
			v.add("%s.integrity_rid does not match %s integrity", prefix, environment.Integrity)
		}
		if !oneOf(environment.DisplayClass, "standard-dpi", "high-dpi", "mixed-dpi", "not-applicable") {
			v.add("%s.display_class %q is invalid", prefix, environment.DisplayClass)
		}
		if !oneOf(environment.ProfileState, "clean", "retained-v0.9.1", "not-applicable") {
			v.add("%s.profile_state %q is invalid", prefix, environment.ProfileState)
		}
		if environment.DisplayClass != "not-applicable" && environment.EffectiveDPI <= 0 {
			v.add("%s.effective_dpi must be positive for a display environment", prefix)
		}
	}
}

func (v *validator) validateAttachments() {
	v.attachments = make(map[string]Attachment, len(v.evidence.Attachments))
	for i, attachment := range v.evidence.Attachments {
		prefix := fmt.Sprintf("attachments[%d]", i)
		if reason := unsafeRelativePath(attachment.Path); reason != "" {
			v.add("%s.path %q is unsafe: %s", prefix, attachment.Path, reason)
			continue
		}
		if !strings.HasPrefix(attachment.Path, "attachments/") {
			v.add("%s.path %q must be beneath attachments/", prefix, attachment.Path)
		}
		if _, exists := v.attachments[attachment.Path]; exists {
			v.add("%s.path %q is duplicate", prefix, attachment.Path)
			continue
		}
		v.attachments[attachment.Path] = attachment
		if !digestPattern.MatchString(attachment.SHA256) {
			v.add("%s.sha256 is not a lowercase SHA-256 digest", prefix)
		}
		if strings.TrimSpace(attachment.MediaType) == "" || strings.TrimSpace(attachment.Purpose) == "" {
			v.add("%s media_type and purpose must not be blank", prefix)
		}
		full := filepath.Join(v.root, filepath.FromSlash(attachment.Path))
		if err := rejectLinkedPath(v.root, attachment.Path); err != nil {
			v.add("%s: %v", prefix, err)
			continue
		}
		info, digest, err := hashRegularFile(full)
		if err != nil {
			v.add("%s: %v", prefix, err)
			continue
		}
		if info.Size() != attachment.Bytes {
			v.add("%s bytes are %d; evidence requires %d", prefix, info.Size(), attachment.Bytes)
		}
		if digest != attachment.SHA256 {
			v.add("%s SHA-256 %s does not match evidence %s", prefix, digest, attachment.SHA256)
		}
	}
}

func rejectLinkedPath(root, relative string) error {
	current := root
	for _, component := range strings.Split(filepath.FromSlash(relative), string(filepath.Separator)) {
		current = filepath.Join(current, component)
		info, err := os.Lstat(current)
		if err != nil {
			return fmt.Errorf("inspect attachment path component %s: %w", current, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("attachment path component %s is a symbolic link", current)
		}
	}
	return nil
}

func (v *validator) validateObservations() {
	seen := make(map[string]int, len(v.evidence.Observations))
	observations := make(map[string]*Observation, len(v.evidence.Observations))
	required := make(map[string]bool, len(requiredScenarios))
	for _, id := range requiredScenarios {
		required[id] = true
	}
	for i := range v.evidence.Observations {
		observation := &v.evidence.Observations[i]
		prefix := fmt.Sprintf("observations[%d] %q", i, observation.ID)
		seen[observation.ID]++
		if seen[observation.ID] > 1 {
			v.add("%s is a duplicate observation", prefix)
		} else {
			observations[observation.ID] = observation
		}
		if !required[observation.ID] {
			v.add("%s is not a required scenario", prefix)
		}
		if !validStatuses[observation.Status] {
			v.add("%s status %q is invalid", prefix, observation.Status)
		} else if observation.Status != "pass" {
			v.add("%s has non-passing status %q", prefix, observation.Status)
		}
		if strings.TrimSpace(observation.Summary) == "" {
			v.add("%s summary must not be blank", prefix)
		}
		if observation.StartedAt.IsZero() || observation.CompletedAt.IsZero() || observation.CompletedAt.Before(observation.StartedAt) {
			v.add("%s timestamps must be present and ordered", prefix)
		} else if observation.StartedAt.Before(v.evidence.StartedAt) || observation.CompletedAt.After(v.evidence.CompletedAt) {
			v.add("%s timestamps must fall within the evidence interval", prefix)
		}
		environment, exists := v.environments[observation.EnvironmentID]
		if !exists {
			v.add("%s references unknown environment %q", prefix, observation.EnvironmentID)
		}
		for _, attachmentPath := range observation.AttachmentPaths {
			if _, exists := v.attachments[attachmentPath]; !exists {
				v.add("%s references unknown attachment %q", prefix, attachmentPath)
			}
		}
		if requiresAttachment(observation.ID) && len(observation.AttachmentPaths) == 0 {
			v.add("%s requires at least one attachment", prefix)
		}
		if strings.HasPrefix(observation.ID, "window.") {
			v.requireAttachmentPurpose(prefix, observation, "native window measurement")
			v.requireRasterImage(prefix, observation)
		}
		if strings.HasPrefix(observation.ID, "error.") ||
			strings.HasPrefix(observation.ID, "setup.") ||
			strings.HasPrefix(observation.ID, "remove.") ||
			strings.HasPrefix(observation.ID, "desktop.") {
			v.requireRasterImage(prefix, observation)
		}
		if strings.HasPrefix(observation.ID, "task.") {
			v.requireAttachmentPurpose(prefix, observation, "task run evidence")
		}
		v.validateScenario(prefix, observation, environment)
	}
	for _, id := range requiredScenarios {
		if seen[id] == 0 {
			v.add("missing required observation %q", id)
		}
	}
	v.validateRelationships(observations)
}

func (v *validator) requireAttachmentPurpose(prefix string, o *Observation, purpose string) {
	for _, name := range o.AttachmentPaths {
		if attachment, ok := v.attachments[name]; ok && attachment.Purpose == purpose {
			return
		}
	}
	v.add("%s requires an attachment with purpose %q", prefix, purpose)
}

func (v *validator) requireRasterImage(prefix string, o *Observation) {
	for _, name := range o.AttachmentPaths {
		if _, ok := v.attachments[name]; !ok {
			continue
		}
		full := filepath.Join(v.root, filepath.FromSlash(name))
		if rasterImageFile(full) {
			return
		}
	}
	v.add("%s requires an attachment containing a supported raster image", prefix)
}

func rasterImageFile(name string) bool {
	file, err := os.Open(name)
	if err != nil {
		return false
	}
	defer file.Close()

	header := make([]byte, 512)
	n, err := io.ReadFull(file, header)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return false
	}
	header = header[:n]
	switch http.DetectContentType(header) {
	case "image/png", "image/jpeg", "image/gif", "image/bmp":
		return true
	}
	if len(header) >= 12 && string(header[:4]) == "RIFF" && string(header[8:12]) == "WEBP" {
		return true
	}
	return validPlainPPM(header)
}

func validPlainPPM(data []byte) bool {
	reader := bytes.NewReader(data)
	var magic string
	var width, height, maximum int
	if _, err := fmt.Fscan(reader, &magic, &width, &height, &maximum); err != nil ||
		magic != "P3" || width <= 0 || height <= 0 || maximum <= 0 || maximum > 65535 ||
		width > 10000 || height > 10000 {
		return false
	}
	samples := width * height * 3
	for i := 0; i < samples; i++ {
		var sample int
		if _, err := fmt.Fscan(reader, &sample); err != nil || sample < 0 || sample > maximum {
			return false
		}
	}
	return true
}

func (v *validator) validateRelationships(observations map[string]*Observation) {
	runIDs := make(map[string]string)
	for _, id := range []string{"task.manual-success", "task.scheduled-success", "task.nonzero-exit", "task.start-failure"} {
		o := observations[id]
		if o == nil {
			continue
		}
		runID, _ := o.Metrics["run_id"].(string)
		if prior, exists := runIDs[runID]; runID != "" && exists {
			v.add("task observations %q and %q must use distinct run_id values", prior, id)
		} else if runID != "" {
			runIDs[runID] = id
		}
	}
	first := observations["window.clean-standard"]
	subsequent := observations["window.subsequent-launch"]
	if first != nil && subsequent != nil {
		firstPID, firstOK := numberMetric(first.Metrics, "pid")
		priorPID, priorOK := numberMetric(subsequent.Metrics, "prior_process_id")
		currentPID, currentOK := numberMetric(subsequent.Metrics, "pid")
		if !firstOK || !priorOK || priorPID != firstPID {
			v.add("window.subsequent-launch prior_process_id must match window.clean-standard pid")
		}
		if firstOK && currentOK && currentPID == firstPID {
			v.add("window.subsequent-launch must use a fresh process id")
		}
	}
}

func (v *validator) validateScenario(prefix string, o *Observation, env Environment) {
	switch {
	case strings.HasPrefix(o.ID, "window."):
		v.validateWindow(prefix, o, env)
	case strings.HasPrefix(o.ID, "error."):
		v.validateError(prefix, o, env)
	case strings.HasPrefix(o.ID, "task."):
		v.validateTask(prefix, o, env)
	case strings.HasPrefix(o.ID, "access."):
		v.validateAccess(prefix, o, env)
	case strings.HasPrefix(o.ID, "setup."):
		v.validateSetup(prefix, o, env)
	case strings.HasPrefix(o.ID, "remove."):
		v.validateRemoval(prefix, o)
	case strings.HasPrefix(o.ID, "desktop."):
		v.validateDesktop(prefix, o, env)
	}
}

func (v *validator) validateDesktop(prefix string, o *Observation, env Environment) {
	v.requireRoutine(prefix, env)
	v.validateDesktopDPI(prefix, o.ID, env)
	switch o.ID {
	case "desktop.appearance-standard", "desktop.appearance-scaled":
		v.validateDesktopAppearance(prefix, o, env)
	case "desktop.interaction-states", "desktop.interaction-states-scaled":
		v.requireExactSet(prefix, o.Metrics, "palettes", "dark", "light")
		v.requireExactSet(prefix, o.Metrics, "control_families", "navigation", "selector", "ordinary", "primary", "danger", "dialog", "table-row")
		v.requireExactSet(prefix, o.Metrics, "states", "rest", "hover", "focus", "pressed", "selected", "disabled")
		v.requireNumber(prefix, o.Metrics, "minimum_text_contrast", 4.5, math.MaxFloat64)
		v.requireNumber(prefix, o.Metrics, "minimum_non_text_contrast", 3, math.MaxFloat64)
		v.requireTrue(prefix, o.Metrics, "labels_readable", "glyphs_readable", "selection_identifiable", "focus_visible", "non_color_cues_present")
	case "desktop.navigation-options", "desktop.navigation-options-scaled":
		v.requireExactSet(prefix, o.Metrics, "palettes", "dark", "light")
		v.requireExactSet(prefix, o.Metrics, "content_sizes", "1280x800", "800x600")
		v.requireString(prefix, o.Metrics, "destination_order", "tasks,groups,chains,schedule,activity,options,info")
		v.requireTrue(prefix, o.Metrics, "rail_spacing_balanced", "labels_unclipped", "boundary_full_height", "boundary_subtle", "exit_bottom_right", "exit_never_selected", "exit_semantic_glyph", "storage_rows_compact", "unavailable_rows_muted", "copy_exact", "selector_current_omitted")
		v.requireFalse(prefix, o.Metrics, "horizontal_scrollbar_present")
	case "desktop.scroll-input":
		v.requireExactSet(prefix, o.Metrics, "sensitivities", "1x", "2x", "4x")
		v.requireExactSet(prefix, o.Metrics, "surfaces", "options", "info", "editor-command", "editor-schedule", "editor-help")
		v.requireTrue(prefix, o.Metrics, "wheel_detents_responsive", "immediate_apply", "persistence_verified", "nested_multiplier_absent", "keyboard_scroll_preserved")
		available, ok := boolMetric(o.Metrics, "touchpad_available")
		if !ok {
			v.add("%s metric touchpad_available must be a boolean", prefix)
		} else if available {
			v.requireTrue(prefix, o.Metrics, "touchpad_fine_deltas_preserved")
		} else {
			v.requireNonEmpty(prefix, o.Metrics, "touchpad_unavailable_reason")
		}
	case "desktop.tasks-table", "desktop.tasks-table-scaled":
		v.requireInteger(prefix, o.Metrics, "row_count", 100, math.MaxInt32)
		v.requireExactSet(prefix, o.Metrics, "palettes", "dark", "light")
		v.requireExactSet(prefix, o.Metrics, "content_sizes", "1280x800", "800x600")
		v.requireExactSet(prefix, o.Metrics, "headers", "task", "enabled", "lifecycle", "time-zone", "group")
		v.requireExactSet(prefix, o.Metrics, "row_states", "odd", "even", "hover", "focus", "selected")
		v.requireTrue(prefix, o.Metrics, "headers_frozen", "status_dimensions_distinct", "bracket_decoration_absent", "full_values_discoverable", "refresh_identity_stable", "removed_selection_clears", "toolbar_actions_work", "double_click_edits")
		v.requireFalse(prefix, o.Metrics, "horizontal_scrollbar_present")
	case "desktop.schedule-activity-tables", "desktop.schedule-activity-tables-scaled":
		v.requireInteger(prefix, o.Metrics, "schedule_row_count", 100, math.MaxInt32)
		v.requireInteger(prefix, o.Metrics, "activity_row_count", 100, math.MaxInt32)
		v.requireExactSet(prefix, o.Metrics, "palettes", "dark", "light")
		v.requireExactSet(prefix, o.Metrics, "content_sizes", "1280x800", "800x600")
		v.requireExactSet(prefix, o.Metrics, "schedule_headers", "when", "task", "event", "outcome")
		v.requireExactSet(prefix, o.Metrics, "activity_headers", "when", "severity", "source", "summary")
		v.requireExactSet(prefix, o.Metrics, "schedule_states", "scheduled", "success", "failure", "skipped", "caught-up", "queued", "missing", "unknown")
		v.requireExactSet(prefix, o.Metrics, "severities", "INFO", "WARNING", "ERROR")
		v.requireExactSet(prefix, o.Metrics, "row_states", "odd", "even", "hover", "focus", "selected")
		v.requireTrue(prefix, o.Metrics, "headers_frozen", "semantic_text_glyphs_match", "non_color_cues_present", "full_values_discoverable", "refresh_identity_stable", "removed_selection_clears", "detail_activation_accurate", "range_calendar_switching", "filter_clear_acknowledge")
		v.requireFalse(prefix, o.Metrics, "horizontal_scrollbar_present")
	}
}

func (v *validator) validateDesktopDPI(prefix, id string, env Environment) {
	scaled := id == "desktop.appearance-scaled" || strings.HasSuffix(id, "-scaled")
	if scaled && env.EffectiveDPI <= 96 {
		v.add("%s must use an environment greater than 96 DPI", prefix)
	}
	if !scaled && env.EffectiveDPI != 96 {
		v.add("%s must use an environment at exactly 96 DPI", prefix)
	}
}

func (v *validator) validateDesktopAppearance(prefix string, o *Observation, env Environment) {
	v.requireExactSet(prefix, o.Metrics, "palettes", "dark", "light")
	v.requireExactSet(prefix, o.Metrics, "fonts_exercised", "system", "geist", "inter", "ubuntu", "monospace")
	v.requireTrue(prefix, o.Metrics, "system_font_default", "system_font_restored", "font_persistence_verified", "info_text_sharp", "body_text_sharp", "labels_centered", "labels_unclipped", "resize_verified", "minimize_restore_verified", "reopen_verified")
	dpi, ok := numberMetric(o.Metrics, "effective_dpi")
	if !ok || dpi != math.Trunc(dpi) || int(dpi) != env.EffectiveDPI {
		v.add("%s metric effective_dpi must be an integer matching its environment", prefix)
	}
	if o.ID == "desktop.appearance-standard" && dpi != 96 {
		v.add("%s metric effective_dpi must be 96", prefix)
	}
	if o.ID == "desktop.appearance-scaled" && dpi <= 96 {
		v.add("%s metric effective_dpi must be greater than 96", prefix)
	}
}

func (v *validator) requireRoutine(prefix string, env Environment) {
	if env.AccountRole != "intended-user" || env.Integrity != "medium" {
		v.add("%s must use the intended user at medium integrity", prefix)
	}
	if env.ServiceIdentity != "LocalSystem" {
		v.add("%s must use the installed LocalSystem service", prefix)
	}
}

func (v *validator) validateAccess(prefix string, o *Observation, env Environment) {
	switch o.ID {
	case "access.intended-user":
		v.requireRoutine(prefix, env)
		v.requireTrue(prefix, o.Metrics, "health_ok", "gui_task_list_ok")
		v.requireFalse(prefix, o.Metrics, "routine_elevation_required")
		if env.ServiceIdentity != "LocalSystem" {
			v.add("%s service identity must be LocalSystem", prefix)
		}
	case "access.unrelated-user-denied":
		if env.AccountRole != "unrelated-user" || env.Integrity != "medium" {
			v.add("%s must use the unrelated user at medium integrity", prefix)
		}
		v.requireFalse(prefix, o.Metrics, "pipe_opened")
		v.requireString(prefix, o.Metrics, "error_kind", "access-denied")
	case "access.fresh-path-resolution":
		v.requireRoutine(prefix, env)
		v.requireTrue(prefix, o.Metrics, "fresh_process", "matches_installed_cli")
		v.requireNonEmpty(prefix, o.Metrics, "resolved_path")
	case "access.path-removed":
		v.requireRoutine(prefix, env)
		v.requireTrue(prefix, o.Metrics, "fresh_process")
		v.requireFalse(prefix, o.Metrics, "resolves")
		v.requireNumber(prefix, o.Metrics, "registry_cardinality", 0, 0)
	}
}

func (v *validator) validateWindow(prefix string, o *Observation, env Environment) {
	v.requireRoutine(prefix, env)
	for _, key := range []string{"outer_rect", "client_rect", "monitor_rect", "work_area_rect"} {
		if !validRect(o.Metrics[key]) {
			v.add("%s metric %s must be a non-empty rectangle", prefix, key)
		}
	}
	v.requireInteger(prefix, o.Metrics, "pid", 1, math.MaxInt32)
	v.requireDigest(prefix, o.Metrics, "executable_sha256")
	v.requireNonEmpty(prefix, o.Metrics, "hwnd", "monitor_id", "process_user_sid")
	v.requireInteger(prefix, o.Metrics, "process_session_id", 0, math.MaxInt32)
	v.requireInteger(prefix, o.Metrics, "process_integrity_rid", 1, math.MaxInt32)
	v.requireInteger(prefix, o.Metrics, "effective_dpi", 1, 1000)
	if observedDPI, ok := numberMetric(o.Metrics, "effective_dpi"); ok && int(observedDPI) != env.EffectiveDPI {
		v.add("%s effective_dpi does not match its environment", prefix)
	}
	v.requireNumber(prefix, o.Metrics, "fyne_scale", 0.01, 100)
	width, widthOK := numberMetric(o.Metrics, "fyne_content_width")
	height, heightOK := numberMetric(o.Metrics, "fyne_content_height")
	workWidth, workWidthOK := numberMetric(o.Metrics, "logical_work_area_width")
	workHeight, workHeightOK := numberMetric(o.Metrics, "logical_work_area_height")
	if !widthOK || !heightOK || !workWidthOK || !workHeightOK || width <= 0 || height <= 0 || workWidth <= 0 || workHeight <= 0 {
		v.add("%s Fyne content and logical work-area dimensions must be positive", prefix)
	} else if strings.HasPrefix(o.ID, "window.clean-") || o.ID == "window.retained-profile" {
		if workWidth >= 1280.0/0.9 && workHeight >= 800.0/0.9 {
			if width != 1280 || height != 800 {
				v.add("%s sufficiently large work area requires 1280 by 800 Fyne content", prefix)
			}
		} else if width > workWidth*0.9+0.01 || height > workHeight*0.9+0.01 {
			v.add("%s Fyne content exceeds 90 percent of the logical work area", prefix)
		}
	}
	v.requireTrue(prefix, o.Metrics, "restored", "margins_visible", "title_bar_reachable", "resize_borders_reachable", "taskbar_reachable")
	v.requireFalse(prefix, o.Metrics, "maximized", "minimized", "fullscreen")
	outer, outerOK := rectMetric(o.Metrics, "outer_rect")
	work, workOK := rectMetric(o.Metrics, "work_area_rect")
	if outerOK && workOK && !(outer.Left > work.Left && outer.Top > work.Top && outer.Right < work.Right && outer.Bottom < work.Bottom) {
		v.add("%s outer_rect must leave positive margins inside work_area_rect", prefix)
	}
	if o.ID == "window.clean-standard" && env.DisplayClass != "standard-dpi" {
		v.add("%s must use a standard-dpi environment", prefix)
	}
	if o.ID == "window.clean-high-or-mixed" && !oneOf(env.DisplayClass, "high-dpi", "mixed-dpi") {
		v.add("%s must use a high-dpi or mixed-dpi environment", prefix)
	}
	if strings.HasPrefix(o.ID, "window.clean-") && env.ProfileState != "clean" {
		v.add("%s clean launch must use clean profile state", prefix)
	}
	if o.ID == "window.retained-profile" && env.ProfileState != "retained-v0.9.1" {
		v.add("%s must use retained-v0.9.1 profile state", prefix)
	}
	if o.ID == "window.state-transitions" {
		v.requireTrue(prefix, o.Metrics, "maximize_worked", "restore_worked", "resize_worked", "minimize_worked", "close_worked")
	}
	if o.ID == "window.clean-standard" {
		v.requireInteger(prefix, o.Metrics, "launch_sequence", 1, 1)
	}
	if o.ID == "window.subsequent-launch" {
		v.requireTrue(prefix, o.Metrics, "fresh_process", "prior_process_closed")
		v.requireInteger(prefix, o.Metrics, "prior_process_id", 1, math.MaxInt32)
		v.requireInteger(prefix, o.Metrics, "launch_sequence", 2, math.MaxInt32)
	}
	v.validateNativeWindow(prefix, o, env)
}

type nativeRect struct {
	Left   int `json:"left"`
	Top    int `json:"top"`
	Right  int `json:"right"`
	Bottom int `json:"bottom"`
}

type fyneWindowEvidence struct {
	SchemaVersion int       `json:"schema_version"`
	ProcessID     int       `json:"process_id"`
	CapturedAt    time.Time `json:"captured_at"`
	ContentWidth  float64   `json:"content_width"`
	ContentHeight float64   `json:"content_height"`
	CanvasScale   float64   `json:"canvas_scale"`
}

type nativeWindowEvidence struct {
	SchemaVersion       int                `json:"schema_version"`
	Kind                string             `json:"kind"`
	ObservationID       string             `json:"observation_id"`
	CapturedAt          time.Time          `json:"captured_at"`
	ProcessID           int                `json:"process_id"`
	ProcessPath         string             `json:"process_path"`
	ProcessSHA256       string             `json:"process_sha256"`
	ProcessSessionID    int                `json:"process_session_id"`
	ProcessUserSID      string             `json:"process_user_sid"`
	ProcessIntegrityRID int                `json:"process_integrity_rid"`
	HWND                string             `json:"hwnd"`
	OuterRect           nativeRect         `json:"outer_rect"`
	ClientRect          nativeRect         `json:"client_rect"`
	MonitorRect         nativeRect         `json:"monitor_rect"`
	WorkAreaRect        nativeRect         `json:"work_area_rect"`
	MonitorID           string             `json:"monitor_id"`
	EffectiveDPI        int                `json:"effective_dpi"`
	ShowCommand         int                `json:"show_command"`
	Visible             bool               `json:"visible"`
	Maximized           bool               `json:"maximized"`
	Minimized           bool               `json:"minimized"`
	Fullscreen          bool               `json:"fullscreen"`
	Restored            bool               `json:"restored"`
	Fyne                fyneWindowEvidence `json:"fyne"`
}

func (v *validator) validateNativeWindow(prefix string, o *Observation, env Environment) {
	var attachment *Attachment
	for _, name := range o.AttachmentPaths {
		candidate, ok := v.attachments[name]
		if ok && candidate.Purpose == "native window measurement" {
			copy := candidate
			attachment = &copy
			break
		}
	}
	if attachment == nil {
		return
	}
	file, err := os.Open(filepath.Join(v.root, filepath.FromSlash(attachment.Path)))
	if err != nil {
		v.add("%s native window attachment cannot be opened: %v", prefix, err)
		return
	}
	defer file.Close()
	data, err := readStrictJSON(file, "native window attachment")
	if err != nil {
		v.add("%s native window attachment: %v", prefix, err)
		return
	}
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	var record nativeWindowEvidence
	if err := decoder.Decode(&record); err != nil {
		v.add("%s native window attachment: %v", prefix, err)
		return
	}
	if record.SchemaVersion != 1 || record.Kind != "native-window-v1" || record.ObservationID != o.ID {
		v.add("%s native window schema, kind, or observation identity is invalid", prefix)
	}
	if record.CapturedAt.Before(v.evidence.StartedAt) || record.CapturedAt.After(v.evidence.CompletedAt) {
		v.add("%s native window capture timestamp is outside the evidence interval", prefix)
	}
	processBase := path.Base(strings.ReplaceAll(record.ProcessPath, `\`, "/"))
	if strings.TrimSpace(record.ProcessPath) == "" || !strings.EqualFold(processBase, "gosched-gui.exe") ||
		!record.Visible || record.Fyne.SchemaVersion != 1 || record.Fyne.ProcessID != record.ProcessID {
		v.add("%s native window process/Fyne identity is invalid", prefix)
	}
	if !digestPattern.MatchString(record.ProcessSHA256) || strings.TrimSpace(record.ProcessUserSID) == "" ||
		strings.TrimSpace(record.HWND) == "" || strings.TrimSpace(record.MonitorID) == "" ||
		record.ProcessID <= 0 || record.ProcessSessionID < 0 || record.ProcessIntegrityRID <= 0 ||
		record.EffectiveDPI <= 0 || record.ShowCommand <= 0 || record.Fyne.ContentWidth <= 0 ||
		record.Fyne.ContentHeight <= 0 || record.Fyne.CanvasScale <= 0 {
		v.add("%s native window attachment contains invalid required measurements", prefix)
	}
	if !nativeRectValid(record.OuterRect) || !nativeRectValid(record.ClientRect) ||
		!nativeRectValid(record.MonitorRect) || !nativeRectValid(record.WorkAreaRect) {
		v.add("%s native window attachment contains an invalid rectangle", prefix)
	}
	if record.Fyne.CapturedAt.Before(v.evidence.StartedAt) || record.Fyne.CapturedAt.After(v.evidence.CompletedAt) {
		v.add("%s Fyne capture timestamp is outside the evidence interval", prefix)
	}
	compareIntegerMetric(v, prefix, o.Metrics, "pid", record.ProcessID)
	compareStringMetric(v, prefix, o.Metrics, "executable_sha256", record.ProcessSHA256)
	compareIntegerMetric(v, prefix, o.Metrics, "process_session_id", record.ProcessSessionID)
	compareStringMetric(v, prefix, o.Metrics, "process_user_sid", record.ProcessUserSID)
	compareIntegerMetric(v, prefix, o.Metrics, "process_integrity_rid", record.ProcessIntegrityRID)
	compareStringMetric(v, prefix, o.Metrics, "hwnd", record.HWND)
	compareStringMetric(v, prefix, o.Metrics, "monitor_id", record.MonitorID)
	compareIntegerMetric(v, prefix, o.Metrics, "effective_dpi", record.EffectiveDPI)
	compareBoolMetric(v, prefix, o.Metrics, "restored", record.Restored)
	compareBoolMetric(v, prefix, o.Metrics, "maximized", record.Maximized)
	compareBoolMetric(v, prefix, o.Metrics, "minimized", record.Minimized)
	compareBoolMetric(v, prefix, o.Metrics, "fullscreen", record.Fullscreen)
	compareRectMetric(v, prefix, o.Metrics, "outer_rect", record.OuterRect)
	compareRectMetric(v, prefix, o.Metrics, "client_rect", record.ClientRect)
	compareRectMetric(v, prefix, o.Metrics, "monitor_rect", record.MonitorRect)
	compareRectMetric(v, prefix, o.Metrics, "work_area_rect", record.WorkAreaRect)
	compareNumberMetric(v, prefix, o.Metrics, "fyne_content_width", record.Fyne.ContentWidth)
	compareNumberMetric(v, prefix, o.Metrics, "fyne_content_height", record.Fyne.ContentHeight)
	compareNumberMetric(v, prefix, o.Metrics, "fyne_scale", record.Fyne.CanvasScale)
	if record.ProcessUserSID != env.AccountSID || record.ProcessIntegrityRID != env.IntegrityRID {
		v.add("%s native window token does not match environment", prefix)
	}
	if record.EffectiveDPI > 0 {
		scale := float64(record.EffectiveDPI) / 96
		compareNumberMetric(v, prefix, o.Metrics, "logical_work_area_width", float64(record.WorkAreaRect.Right-record.WorkAreaRect.Left)/scale)
		compareNumberMetric(v, prefix, o.Metrics, "logical_work_area_height", float64(record.WorkAreaRect.Bottom-record.WorkAreaRect.Top)/scale)
	}
}

func nativeRectValid(value nativeRect) bool {
	return value.Right > value.Left && value.Bottom > value.Top
}

func compareStringMetric(v *validator, prefix string, metrics map[string]any, key, expected string) {
	value, ok := metrics[key].(string)
	if !ok || value != expected {
		v.add("%s metric %s does not match native window attachment", prefix, key)
	}
}

func compareIntegerMetric(v *validator, prefix string, metrics map[string]any, key string, expected int) {
	value, ok := numberMetric(metrics, key)
	if !ok || value != float64(expected) {
		v.add("%s metric %s does not match native window attachment", prefix, key)
	}
}

func compareNumberMetric(v *validator, prefix string, metrics map[string]any, key string, expected float64) {
	value, ok := numberMetric(metrics, key)
	if !ok || math.Abs(value-expected) > 0.02 {
		v.add("%s metric %s does not match native window attachment", prefix, key)
	}
}

func compareBoolMetric(v *validator, prefix string, metrics map[string]any, key string, expected bool) {
	value, ok := boolMetric(metrics, key)
	if !ok || value != expected {
		v.add("%s metric %s does not match native window attachment", prefix, key)
	}
}

func compareRectMetric(v *validator, prefix string, metrics map[string]any, key string, expected nativeRect) {
	value, ok := rectMetric(metrics, key)
	if !ok || value.Left != float64(expected.Left) || value.Top != float64(expected.Top) || value.Right != float64(expected.Right) || value.Bottom != float64(expected.Bottom) {
		v.add("%s metric %s does not match native window attachment", prefix, key)
	}
}

func (v *validator) validateError(prefix string, o *Observation, env Environment) {
	v.requireRoutine(prefix, env)
	v.requireString(prefix, o.Metrics, "category", strings.TrimPrefix(o.ID, "error."))
	v.requireNonEmpty(prefix, o.Metrics, "trigger")
	v.requireInteger(prefix, o.Metrics, "sample_count", 1, math.MaxInt32)
	v.requireInteger(prefix, o.Metrics, "max_in_frame_count", 1, 1)
	v.requireInteger(prefix, o.Metrics, "max_modal_count", 0, 0)
	v.requireInteger(prefix, o.Metrics, "max_additional_top_level_count", 0, 0)
	v.requireTrue(prefix, o.Metrics, "retry_reachable", "exit_reachable")
	v.requireFalse(prefix, o.Metrics, "guidance_recommends_elevation")
	if o.ID != "error.recovery" {
		v.requireInteger(prefix, o.Metrics, "duration_seconds", 120, math.MaxInt32)
		if o.CompletedAt.Sub(o.StartedAt) < 120*time.Second {
			v.add("%s timestamp interval must be at least 120 seconds", prefix)
		}
	}
	if o.ID == "error.manual-retry" {
		v.requireTrue(prefix, o.Metrics, "retry_invoked")
	}
	if o.ID == "error.recovery" {
		v.requireTrue(prefix, o.Metrics, "incident_cleared", "interface_restored")
		v.requireFalse(prefix, o.Metrics, "reinstall_required")
		v.requireString(prefix, o.Metrics, "recovery_method", "daemon-restored")
	}
}

func (v *validator) validateTask(prefix string, o *Observation, env Environment) {
	v.requireRoutine(prefix, env)
	v.requireNonEmpty(prefix, o.Metrics, "public_interface", "task_id", "run_id")
	v.requireTrue(prefix, o.Metrics, "production_runner", "history_recorded")
	v.requireDigest(prefix, o.Metrics, "task_definition_sha256", "output_sha256", "marker_sha256", "history_sha256")
	switch o.ID {
	case "task.manual-success":
		v.requireString(prefix, o.Metrics, "invocation_mode", "manual")
		v.requireInteger(prefix, o.Metrics, "expected_exit_code", 0, 0)
		v.requireInteger(prefix, o.Metrics, "actual_exit_code", 0, 0)
	case "task.scheduled-success":
		v.requireString(prefix, o.Metrics, "invocation_mode", "scheduled")
		v.requireInteger(prefix, o.Metrics, "expected_exit_code", 0, 0)
		v.requireInteger(prefix, o.Metrics, "actual_exit_code", 0, 0)
	case "task.nonzero-exit":
		v.requireString(prefix, o.Metrics, "invocation_mode", "manual-fault")
		v.requireNonEmpty(prefix, o.Metrics, "trigger")
		v.requireInteger(prefix, o.Metrics, "expected_exit_code", 1, math.MaxInt32)
		expected, eOK := numberMetric(o.Metrics, "expected_exit_code")
		actual, aOK := numberMetric(o.Metrics, "actual_exit_code")
		if !eOK || !aOK || expected != actual {
			v.add("%s actual_exit_code must match the deliberate nonzero expected_exit_code", prefix)
		}
		v.requireString(prefix, o.Metrics, "diagnostic_category", "nonzero-exit")
	case "task.start-failure":
		v.requireString(prefix, o.Metrics, "invocation_mode", "manual-fault")
		v.requireNonEmpty(prefix, o.Metrics, "trigger")
		v.requireInteger(prefix, o.Metrics, "expected_exit_code", -1, -1)
		v.requireString(prefix, o.Metrics, "diagnostic_category", "process-start-failure")
	}
	v.validateTaskAttachment(prefix, o)
}

type taskRunEvidence struct {
	ObservationID   string `json:"observation_id"`
	TaskID          string `json:"task_id"`
	RunID           string `json:"run_id"`
	InvocationMode  string `json:"invocation_mode"`
	PublicInterface string `json:"public_interface"`
	TaskDefinition  string `json:"task_definition"`
	Output          string `json:"output"`
	Marker          string `json:"marker"`
	History         string `json:"history"`
	ProductionRun   bool   `json:"production_runner"`
	ExpectedExit    int    `json:"expected_exit_code"`
	ActualExit      int    `json:"actual_exit_code"`
	Diagnostic      string `json:"diagnostic_category"`
}

type taskEvidenceBundle struct {
	SchemaVersion int               `json:"schema_version"`
	Kind          string            `json:"kind"`
	Runs          []taskRunEvidence `json:"runs"`
}

func (v *validator) validateTaskAttachment(prefix string, o *Observation) {
	var attachment *Attachment
	for _, name := range o.AttachmentPaths {
		candidate, ok := v.attachments[name]
		if ok && candidate.Purpose == "task run evidence" {
			copy := candidate
			attachment = &copy
			break
		}
	}
	if attachment == nil {
		return
	}
	file, err := os.Open(filepath.Join(v.root, filepath.FromSlash(attachment.Path)))
	if err != nil {
		v.add("%s task evidence cannot be opened: %v", prefix, err)
		return
	}
	defer file.Close()
	data, err := readStrictJSON(file, "task evidence attachment")
	if err != nil {
		v.add("%s task evidence attachment: %v", prefix, err)
		return
	}
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	var bundle taskEvidenceBundle
	if err := decoder.Decode(&bundle); err != nil {
		v.add("%s task evidence attachment: %v", prefix, err)
		return
	}
	if bundle.SchemaVersion != 1 || bundle.Kind != "task-run-evidence-v1" {
		v.add("%s task evidence schema or kind is invalid", prefix)
	}
	if len(bundle.Runs) != 4 {
		v.add("%s task evidence must contain exactly four canonical runs", prefix)
	}
	seen := make(map[string]bool, len(bundle.Runs))
	for _, candidate := range bundle.Runs {
		if seen[candidate.ObservationID] || !oneOf(candidate.ObservationID,
			"task.manual-success", "task.scheduled-success", "task.nonzero-exit", "task.start-failure") {
			v.add("%s task evidence contains a duplicate or unknown observation identity", prefix)
		}
		seen[candidate.ObservationID] = true
	}
	var matches []taskRunEvidence
	for _, run := range bundle.Runs {
		if run.ObservationID == o.ID {
			matches = append(matches, run)
		}
	}
	if len(matches) != 1 {
		v.add("%s task evidence must contain exactly one matching run", prefix)
		return
	}
	run := matches[0]
	compareStringMetric(v, prefix, o.Metrics, "task_id", run.TaskID)
	compareStringMetric(v, prefix, o.Metrics, "run_id", run.RunID)
	compareStringMetric(v, prefix, o.Metrics, "invocation_mode", run.InvocationMode)
	compareStringMetric(v, prefix, o.Metrics, "public_interface", run.PublicInterface)
	compareBoolMetric(v, prefix, o.Metrics, "production_runner", run.ProductionRun)
	compareIntegerMetric(v, prefix, o.Metrics, "expected_exit_code", run.ExpectedExit)
	if o.ID != "task.start-failure" {
		compareIntegerMetric(v, prefix, o.Metrics, "actual_exit_code", run.ActualExit)
	}
	if o.ID == "task.nonzero-exit" || o.ID == "task.start-failure" {
		compareStringMetric(v, prefix, o.Metrics, "diagnostic_category", run.Diagnostic)
	}
	compareStringMetric(v, prefix, o.Metrics, "task_definition_sha256", digestString(run.TaskDefinition))
	compareStringMetric(v, prefix, o.Metrics, "output_sha256", digestString(run.Output))
	compareStringMetric(v, prefix, o.Metrics, "marker_sha256", digestString(run.Marker))
	compareStringMetric(v, prefix, o.Metrics, "history_sha256", digestString(run.History))
	if strings.TrimSpace(run.TaskDefinition) == "" || strings.TrimSpace(run.Output) == "" ||
		strings.TrimSpace(run.Marker) == "" || strings.TrimSpace(run.History) == "" {
		v.add("%s task evidence must retain task definition, output, marker, and history", prefix)
	}
}

func digestString(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func (v *validator) validateSetup(prefix string, o *Observation, env Environment) {
	v.requireNonEmpty(prefix, o.Metrics, "install_session_id", "installer_process_owner_role", "installer_process_integrity")
	v.requireInteger(prefix, o.Metrics, "installer_process_id", 1, math.MaxInt32)
	v.requireInteger(prefix, o.Metrics, "installer_session_id", 0, math.MaxInt32)
	v.requireDigest(prefix, o.Metrics, "selected_options_sha256", "observed_targets_sha256", "before_fingerprint", "after_fingerprint")
	v.requireFalse(prefix, o.Metrics, "owned_data_cleanup_invoked")
	switch o.ID {
	case "setup.shortcut-defaults":
		v.requireTrue(prefix, o.Metrics, "start_menu_default")
		v.requireFalse(prefix, o.Metrics, "desktop_default")
		v.requireTrue(prefix, o.Metrics, "defaults_visible", "effects_verified")
	case "setup.shortcut-matrix":
		v.requireInteger(prefix, o.Metrics, "combinations_verified", 4, 4)
		v.requireTrue(prefix, o.Metrics, "targets_verified")
	case "setup.completion-matrix":
		v.requireInteger(prefix, o.Metrics, "combinations_verified", 4, 4)
		v.requireTrue(prefix, o.Metrics, "independent_choices", "default_handler_verified")
		v.requireTrue(prefix, o.Metrics, "launch_default")
		v.requireFalse(prefix, o.Metrics, "documentation_default")
	case "setup.finish-launch-integrity":
		v.requireRoutine(prefix, env)
		v.requireString(prefix, o.Metrics, "process_integrity", "medium")
		v.requireInteger(prefix, o.Metrics, "launch_count", 1, 1)
	case "setup.cancel":
		v.requireTrue(prefix, o.Metrics, "state_unchanged")
		v.requireFalse(prefix, o.Metrics, "owned_data_cleanup_invoked")
		v.requireEqualMetrics(prefix, o.Metrics, "before_fingerprint", "after_fingerprint")
	case "setup.maintenance":
		v.requireTrue(prefix, o.Metrics, "transitions_verified", "repair_verified", "completion_actions_absent")
		v.requireFalse(prefix, o.Metrics, "owned_data_cleanup_invoked")
	case "setup.upgrade":
		v.requireTrue(prefix, o.Metrics, "choices_preserved", "completion_actions_absent")
		v.requireFalse(prefix, o.Metrics, "owned_data_cleanup_invoked")
	case "setup.invalid-input":
		v.requireTrue(prefix, o.Metrics, "input_rejected", "state_unchanged")
		v.requireFalse(prefix, o.Metrics, "owned_data_cleanup_invoked")
		v.requireEqualMetrics(prefix, o.Metrics, "before_fingerprint", "after_fingerprint")
	case "setup.rollback":
		v.requireTrue(prefix, o.Metrics, "rollback_completed", "state_unchanged")
		v.requireEqualMetrics(prefix, o.Metrics, "before_fingerprint", "after_fingerprint")
	}
}

func (v *validator) validateRemoval(prefix string, o *Observation) {
	v.requireInteger(prefix, o.Metrics, "owned_roots_count", 1, math.MaxInt32)
	v.requireDigest(prefix, o.Metrics, "before_content_sha256", "after_content_sha256", "controls_before_sha256", "controls_after_sha256")
	v.requireNonEmpty(prefix, o.Metrics, "security_state_disposition", "reinstall_result")
	v.requireString(prefix, o.Metrics, "security_state_disposition", "preserved")
	beforeControl, beforeOK := o.Metrics["controls_before_sha256"].(string)
	afterControl, afterOK := o.Metrics["controls_after_sha256"].(string)
	if beforeOK && afterOK && beforeControl != afterControl {
		v.add("%s unaffected control fingerprints must match", prefix)
	}
	switch o.ID {
	case "remove.preserve":
		v.requireTrue(prefix, o.Metrics, "software_removed", "owned_bytes_preserved", "controls_unchanged")
		v.requireTrue(prefix, o.Metrics, "preserve_default_visible", "owned_inventory_reviewed")
		v.requireFalse(prefix, o.Metrics, "owned_data_cleanup_invoked")
		v.requireString(prefix, o.Metrics, "mode", "preserve")
		v.requireString(prefix, o.Metrics, "reinstall_result", "preserved-state-restored")
		v.requireEqualMetrics(prefix, o.Metrics, "before_content_sha256", "after_content_sha256")
	case "remove.wipe":
		v.requireTrue(prefix, o.Metrics, "software_removed", "owned_roots_removed", "controls_unchanged", "security_state_preserved")
		v.requireTrue(prefix, o.Metrics, "wipe_explicitly_selected", "wipe_confirmed", "owned_inventory_reviewed", "owned_data_cleanup_invoked")
		v.requireString(prefix, o.Metrics, "mode", "wipe")
		v.requireString(prefix, o.Metrics, "reinstall_result", "clean-state")
	case "remove.cancel":
		v.requireTrue(prefix, o.Metrics, "software_unchanged", "data_unchanged")
		v.requireFalse(prefix, o.Metrics, "owned_data_cleanup_invoked")
		v.requireString(prefix, o.Metrics, "reinstall_result", "not-run")
		v.requireEqualMetrics(prefix, o.Metrics, "before_content_sha256", "after_content_sha256")
	case "remove.multiple-profiles":
		v.requireInteger(prefix, o.Metrics, "profile_count", 2, math.MaxInt32)
		v.requireTrue(prefix, o.Metrics, "all_profiles_accounted")
	case "remove.locked-partial":
		v.requireString(prefix, o.Metrics, "cleanup_result", "partial")
		v.requireInteger(prefix, o.Metrics, "residual_count", 1, math.MaxInt32)
		v.requireTrue(prefix, o.Metrics, "truthfully_reported")
	case "remove.reinstall-after-preserve":
		v.requireTrue(prefix, o.Metrics, "prior_tasks_available", "prior_preferences_available", "prior_config_available", "prior_logs_available")
		v.requireString(prefix, o.Metrics, "reinstall_result", "preserved-state-restored")
	case "remove.reinstall-after-wipe":
		v.requireFalse(prefix, o.Metrics, "prior_tasks_available", "prior_preferences_available", "prior_config_available", "prior_logs_available")
		v.requireString(prefix, o.Metrics, "reinstall_result", "clean-state")
	}
}

func (v *validator) requireEqualMetrics(prefix string, metrics map[string]any, leftKey, rightKey string) {
	left, leftOK := metrics[leftKey].(string)
	right, rightOK := metrics[rightKey].(string)
	if !leftOK || !rightOK || left != right {
		v.add("%s metrics %s and %s must match", prefix, leftKey, rightKey)
	}
}

func (v *validator) requireTrue(prefix string, metrics map[string]any, keys ...string) {
	for _, key := range keys {
		value, ok := boolMetric(metrics, key)
		if !ok || !value {
			v.add("%s metric %s must be true", prefix, key)
		}
	}
}

func (v *validator) requireFalse(prefix string, metrics map[string]any, keys ...string) {
	for _, key := range keys {
		value, ok := boolMetric(metrics, key)
		if !ok || value {
			v.add("%s metric %s must be false", prefix, key)
		}
	}
}

func (v *validator) requireNonEmpty(prefix string, metrics map[string]any, keys ...string) {
	for _, key := range keys {
		value, ok := metrics[key].(string)
		if !ok || strings.TrimSpace(value) == "" {
			v.add("%s metric %s must be a non-empty string", prefix, key)
		}
	}
}

func (v *validator) requireString(prefix string, metrics map[string]any, key, expected string) {
	value, ok := metrics[key].(string)
	if !ok || value != expected {
		v.add("%s metric %s must be %q", prefix, key, expected)
	}
}

func (v *validator) requireExactSet(prefix string, metrics map[string]any, key string, expected ...string) {
	value, ok := metrics[key].(string)
	if !ok {
		v.add("%s metric %s must be a comma-separated string", prefix, key)
		return
	}
	seen := make(map[string]bool, len(expected))
	actual := make([]string, 0, len(expected))
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			v.add("%s metric %s must not contain blank values", prefix, key)
			return
		}
		if seen[item] {
			v.add("%s metric %s must not contain duplicate value %q", prefix, key, item)
			return
		}
		seen[item] = true
		actual = append(actual, item)
	}
	want := append([]string(nil), expected...)
	sort.Strings(actual)
	sort.Strings(want)
	if len(actual) != len(want) {
		v.add("%s metric %s must contain exactly %s", prefix, key, strings.Join(want, ", "))
		return
	}
	for index := range want {
		if actual[index] != want[index] {
			v.add("%s metric %s must contain exactly %s", prefix, key, strings.Join(want, ", "))
			return
		}
	}
}

func (v *validator) requireNumber(prefix string, metrics map[string]any, key string, minimum, maximum float64) {
	value, ok := numberMetric(metrics, key)
	if !ok || value < minimum || value > maximum {
		v.add("%s metric %s must be between %v and %v", prefix, key, minimum, maximum)
	}
}

func (v *validator) requireInteger(prefix string, metrics map[string]any, key string, minimum, maximum float64) {
	value, ok := numberMetric(metrics, key)
	if !ok || value != math.Trunc(value) || value < minimum || value > maximum {
		v.add("%s metric %s must be an integer between %v and %v", prefix, key, minimum, maximum)
	}
}

func (v *validator) requireDigest(prefix string, metrics map[string]any, keys ...string) {
	for _, key := range keys {
		value, ok := metrics[key].(string)
		if !ok || !digestPattern.MatchString(value) {
			v.add("%s metric %s must be a lowercase SHA-256 digest", prefix, key)
		}
	}
}

func boolMetric(metrics map[string]any, key string) (bool, bool) {
	value, ok := metrics[key].(bool)
	return value, ok
}

func numberMetric(metrics map[string]any, key string) (float64, bool) {
	value, exists := metrics[key]
	if !exists {
		return 0, false
	}
	switch number := value.(type) {
	case int:
		return float64(number), true
	case int64:
		return float64(number), true
	case float64:
		return number, !math.IsNaN(number) && !math.IsInf(number, 0)
	case json.Number:
		parsed, err := number.Float64()
		return parsed, err == nil && !math.IsNaN(parsed) && !math.IsInf(parsed, 0)
	default:
		return 0, false
	}
}

func validRect(value any) bool {
	rectangle, ok := value.(map[string]any)
	if !ok {
		return false
	}
	left, lOK := numberMetric(rectangle, "left")
	top, tOK := numberMetric(rectangle, "top")
	right, rOK := numberMetric(rectangle, "right")
	bottom, bOK := numberMetric(rectangle, "bottom")
	return lOK && tOK && rOK && bOK && right > left && bottom > top
}

type metricRect struct {
	Left, Top, Right, Bottom float64
}

func rectMetric(metrics map[string]any, key string) (metricRect, bool) {
	rectangle, ok := metrics[key].(map[string]any)
	if !ok {
		return metricRect{}, false
	}
	left, lOK := numberMetric(rectangle, "left")
	top, tOK := numberMetric(rectangle, "top")
	right, rOK := numberMetric(rectangle, "right")
	bottom, bOK := numberMetric(rectangle, "bottom")
	return metricRect{Left: left, Top: top, Right: right, Bottom: bottom}, lOK && tOK && rOK && bOK
}

func requiresAttachment(id string) bool {
	return strings.HasPrefix(id, "window.") || strings.HasPrefix(id, "error.") || strings.HasPrefix(id, "task.") || strings.HasPrefix(id, "setup.") || strings.HasPrefix(id, "remove.") || strings.HasPrefix(id, "desktop.")
}

func unsafeRelativePath(name string) string {
	if name == "" {
		return "path is empty"
	}
	if strings.Contains(name, "\\") {
		return "backslashes are not allowed"
	}
	if strings.HasPrefix(name, "/") || filepath.IsAbs(name) || filepath.VolumeName(name) != "" {
		return "absolute or drive-qualified paths are not allowed"
	}
	clean := path.Clean(name)
	if clean != name || clean == "." || clean == ".." || strings.HasPrefix(clean, "../") {
		return "path traversal or non-canonical segments are not allowed"
	}
	return ""
}

func hashRegularFile(name string) (os.FileInfo, string, error) {
	pathInfo, err := os.Lstat(name)
	if err != nil {
		return nil, "", fmt.Errorf("inspect %s: %w", name, err)
	}
	if pathInfo.Mode()&os.ModeSymlink != 0 {
		return nil, "", fmt.Errorf("%s is a symbolic link", name)
	}
	file, err := os.Open(name)
	if err != nil {
		return nil, "", fmt.Errorf("open %s: %w", name, err)
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return nil, "", fmt.Errorf("stat %s: %w", name, err)
	}
	if !info.Mode().IsRegular() {
		return nil, "", fmt.Errorf("%s is not a regular file", name)
	}
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, "", fmt.Errorf("hash %s: %w", name, err)
	}
	return info, hex.EncodeToString(hash.Sum(nil)), nil
}

func oneOf(value string, allowed ...string) bool {
	for _, candidate := range allowed {
		if value == candidate {
			return true
		}
	}
	return false
}
