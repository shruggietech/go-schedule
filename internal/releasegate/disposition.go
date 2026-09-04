package releasegate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// IssueDispositionMapping binds one release-readiness issue to the formal
// candidate observations that support its individual acceptance review.
type IssueDispositionMapping struct {
	Issue          int
	ObservationIDs []string
	RelatedIssues  []int
}

// DispositionFile is one deterministic file in a rendered issue packet.
type DispositionFile struct {
	Name string
	Data []byte
}

type dispositionPacketIndex struct {
	SchemaVersion int                    `json:"schema_version"`
	Candidate     Candidate              `json:"candidate"`
	Issues        []dispositionIssueFile `json:"issues"`
}

type dispositionIssueFile struct {
	Issue        int      `json:"issue"`
	File         string   `json:"file"`
	Observations []string `json:"observations"`
}

type dispositionWriteFile func(string, []byte, fs.FileMode) error

// IssueDispositionMappings returns a deep copy of the reviewed v1.0.0
// observation-to-issue contract in deterministic issue order.
func IssueDispositionMappings() []IssueDispositionMapping {
	required := RequiredScenarioIDs()
	mappings := []IssueDispositionMapping{
		{Issue: 96, ObservationIDs: required[:36], RelatedIssues: []int{97, 98, 94, 89, 90}},
		{Issue: 98, ObservationIDs: required[20:36]},
		{Issue: 101, ObservationIDs: []string{
			"desktop.appearance-standard", "desktop.appearance-scaled",
		}},
		{Issue: 104, ObservationIDs: []string{
			"desktop.navigation-options", "desktop.navigation-options-scaled",
			"desktop.interaction-states", "desktop.interaction-states-scaled",
		}},
		{Issue: 105, ObservationIDs: []string{
			"desktop.navigation-options", "desktop.navigation-options-scaled",
			"desktop.interaction-states", "desktop.interaction-states-scaled",
		}},
		{Issue: 106, ObservationIDs: []string{
			"desktop.appearance-standard", "desktop.appearance-scaled",
			"desktop.navigation-options", "desktop.navigation-options-scaled",
			"desktop.scroll-input",
		}},
		{Issue: 109, ObservationIDs: []string{
			"desktop.interaction-states", "desktop.interaction-states-scaled",
			"desktop.tasks-table", "desktop.tasks-table-scaled",
			"desktop.schedule-activity-tables", "desktop.schedule-activity-tables-scaled",
		}},
		{Issue: 111, ObservationIDs: []string{"desktop.scroll-input"}},
		{Issue: 112, ObservationIDs: []string{
			"desktop.tasks-table", "desktop.tasks-table-scaled",
			"desktop.interaction-states", "desktop.interaction-states-scaled",
		}},
		{Issue: 113, ObservationIDs: []string{
			"desktop.schedule-activity-tables", "desktop.schedule-activity-tables-scaled",
			"desktop.interaction-states", "desktop.interaction-states-scaled",
		}},
	}
	result := make([]IssueDispositionMapping, len(mappings))
	for i, mapping := range mappings {
		result[i] = IssueDispositionMapping{
			Issue:          mapping.Issue,
			ObservationIDs: append([]string(nil), mapping.ObservationIDs...),
			RelatedIssues:  append([]int(nil), mapping.RelatedIssues...),
		}
	}
	return result
}

// RenderDispositionPacket renders the complete deterministic packet in memory.
// Callers must validate the evidence and candidate inputs before invoking it.
func RenderDispositionPacket(evidence Evidence) ([]DispositionFile, error) {
	if !validRepository(evidence.Candidate.Repository) {
		return nil, fmt.Errorf("candidate repository %q is not a GitHub owner/repository", evidence.Candidate.Repository)
	}
	if !tagPattern.MatchString(evidence.Candidate.Tag) {
		return nil, fmt.Errorf("candidate tag %q is not vMAJOR.MINOR.PATCH", evidence.Candidate.Tag)
	}
	if !commitPattern.MatchString(evidence.Candidate.Commit) {
		return nil, fmt.Errorf("candidate commit %q is not 40 lowercase hexadecimal characters", evidence.Candidate.Commit)
	}
	if evidence.Candidate.RunID <= 0 || evidence.Candidate.RunAttempt <= 0 {
		return nil, fmt.Errorf("candidate workflow run identity must be positive")
	}

	observations := make(map[string]Observation, len(evidence.Observations))
	for _, observation := range evidence.Observations {
		if _, exists := observations[observation.ID]; exists {
			return nil, fmt.Errorf("observation %q is duplicate", observation.ID)
		}
		observations[observation.ID] = observation
	}
	environments := make(map[string]Environment, len(evidence.Environments))
	for _, environment := range evidence.Environments {
		if _, exists := environments[environment.ID]; exists {
			return nil, fmt.Errorf("environment %q is duplicate", environment.ID)
		}
		environments[environment.ID] = environment
	}

	mappings := IssueDispositionMappings()
	files := make([]DispositionFile, 0, len(mappings)+1)
	index := dispositionPacketIndex{SchemaVersion: 1, Candidate: evidence.Candidate}
	for _, mapping := range mappings {
		name := fmt.Sprintf("issue-%03d.md", mapping.Issue)
		data, err := renderIssueDisposition(evidence.Candidate, mapping, observations, environments)
		if err != nil {
			return nil, fmt.Errorf("render issue #%d: %w", mapping.Issue, err)
		}
		files = append(files, DispositionFile{Name: name, Data: data})
		index.Issues = append(index.Issues, dispositionIssueFile{
			Issue:        mapping.Issue,
			File:         name,
			Observations: append([]string(nil), mapping.ObservationIDs...),
		})
	}
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode packet index: %w", err)
	}
	files = append(files, DispositionFile{Name: "packet.json", Data: append(data, '\n')})
	return files, nil
}

// WriteDispositionPacket commits a complete packet to an absent directory. It
// never merges with or replaces an existing path.
func WriteDispositionPacket(outputDir string, evidence Evidence) error {
	return writeDispositionPacket(outputDir, evidence, os.WriteFile)
}

func writeDispositionPacket(outputDir string, evidence Evidence, writeFile dispositionWriteFile) (err error) {
	if strings.TrimSpace(outputDir) == "" {
		return fmt.Errorf("output directory must not be blank")
	}
	files, err := RenderDispositionPacket(evidence)
	if err != nil {
		return err
	}
	target, err := filepath.Abs(outputDir)
	if err != nil {
		return fmt.Errorf("resolve output directory: %w", err)
	}
	if _, err := os.Lstat(target); err == nil {
		return fmt.Errorf("output directory %q already exists", target)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect output directory %q: %w", target, err)
	}
	parent := filepath.Dir(target)
	parentInfo, err := os.Lstat(parent)
	if err != nil {
		return fmt.Errorf("inspect output parent %q: %w", parent, err)
	}
	if parentInfo.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("output parent %q is a symbolic link", parent)
	}
	if !parentInfo.IsDir() {
		return fmt.Errorf("output parent %q is not a directory", parent)
	}
	hasLink, err := pathContainsSymlink(parent)
	if err != nil {
		return fmt.Errorf("inspect output parent path %q: %w", parent, err)
	}
	if hasLink {
		return fmt.Errorf("output parent %q contains a symbolic link", parent)
	}

	staging, err := os.MkdirTemp(parent, ".go-schedule-dispositions-")
	if err != nil {
		return fmt.Errorf("create disposition staging directory: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			if cleanupErr := os.RemoveAll(staging); cleanupErr != nil && err == nil {
				err = fmt.Errorf("clean disposition staging directory: %w", cleanupErr)
			}
		}
	}()
	for _, file := range files {
		if filepath.Base(file.Name) != file.Name {
			return fmt.Errorf("packet filename %q is not a basename", file.Name)
		}
		if err := writeFile(filepath.Join(staging, file.Name), file.Data, 0o600); err != nil {
			return fmt.Errorf("write packet file %q: %w", file.Name, err)
		}
	}
	if err := os.Rename(staging, target); err != nil {
		return fmt.Errorf("commit disposition packet %q: %w", target, err)
	}
	committed = true
	return nil
}

func renderIssueDisposition(candidate Candidate, mapping IssueDispositionMapping, observations map[string]Observation, environments map[string]Environment) ([]byte, error) {
	var builder strings.Builder
	fmt.Fprintf(&builder, "# Formal %s candidate evidence for #%d\n\n", candidate.Tag, mapping.Issue)
	builder.WriteString("Generated from the exact attended Windows candidate after the ")
	builder.WriteString("production candidate, archive, attachment, and manifest validation passed.\n\n")
	builder.WriteString("## Candidate identity\n\n")
	builder.WriteString("| Field | Value |\n| --- | --- |\n")
	writeCandidateRow(&builder, "Repository", candidate.Repository)
	writeCandidateRow(&builder, "Tag", candidate.Tag)
	writeCandidateRow(&builder, "Commit", candidate.Commit)
	writeCandidateRow(&builder, "Workflow", candidate.Workflow)
	writeCandidateRow(&builder, "Run ID", strconv.FormatInt(candidate.RunID, 10))
	writeCandidateRow(&builder, "Run attempt", strconv.Itoa(candidate.RunAttempt))
	writeCandidateRow(&builder, "MSI", candidate.Filename)
	writeCandidateRow(&builder, "Bytes", strconv.FormatInt(candidate.Bytes, 10))
	writeCandidateRow(&builder, "SHA-256", candidate.SHA256)
	writeCandidateRow(&builder, "ProductVersion", candidate.ProductVersion)
	writeCandidateRow(&builder, "ProductCode", candidate.ProductCode)
	runURL := fmt.Sprintf("https://github.com/%s/actions/runs/%d/attempts/%d", candidate.Repository, candidate.RunID, candidate.RunAttempt)
	archiveName := fmt.Sprintf("go-schedule_%s_windows-attended-evidence.zip", candidate.Tag)
	archiveURL := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", candidate.Repository, candidate.Tag, archiveName)
	fmt.Fprintf(&builder, "\n- Workflow run: [%d attempt %d](%s)\n", candidate.RunID, candidate.RunAttempt, runURL)
	fmt.Fprintf(&builder, "- Evidence archive: [`%s`](%s)\n", archiveName, archiveURL)
	if len(mapping.RelatedIssues) != 0 {
		builder.WriteString("- Coordinator references: ")
		for i, issue := range mapping.RelatedIssues {
			if i != 0 {
				builder.WriteString(", ")
			}
			fmt.Fprintf(&builder, "#%d", issue)
		}
		builder.WriteString("\n")
	}

	environmentIDs := make(map[string]bool)
	issueObservations := make([]Observation, 0, len(mapping.ObservationIDs))
	for _, id := range mapping.ObservationIDs {
		observation, ok := observations[id]
		if !ok {
			return nil, fmt.Errorf("required observation %q is missing", id)
		}
		if _, ok := environments[observation.EnvironmentID]; !ok {
			return nil, fmt.Errorf("observation %q references unknown environment %q", id, observation.EnvironmentID)
		}
		for _, attachment := range observation.AttachmentPaths {
			if reason := unsafeRelativePath(attachment); reason != "" {
				return nil, fmt.Errorf("observation %q attachment %q is unsafe: %s", id, attachment, reason)
			}
		}
		environmentIDs[observation.EnvironmentID] = true
		issueObservations = append(issueObservations, observation)
	}

	builder.WriteString("\n## Referenced environments\n\n")
	builder.WriteString("| ID | Snapshot | Windows | Build | Account role | Integrity | Display | DPI | Profile |\n")
	builder.WriteString("| --- | --- | --- | --- | --- | --- | --- | ---: | --- |\n")
	orderedEnvironments := make([]string, 0, len(environmentIDs))
	for id := range environmentIDs {
		orderedEnvironments = append(orderedEnvironments, id)
	}
	sort.Strings(orderedEnvironments)
	for _, id := range orderedEnvironments {
		environment := environments[id]
		values := []string{
			environment.ID, environment.Snapshot, environment.WindowsEdition,
			environment.WindowsBuild, environment.AccountRole, environment.Integrity,
			environment.DisplayClass, strconv.Itoa(environment.EffectiveDPI), environment.ProfileState,
		}
		builder.WriteString("|")
		for _, value := range values {
			fmt.Fprintf(&builder, " %s |", markdownCell(value))
		}
		builder.WriteString("\n")
	}

	builder.WriteString("\n## Mapped formal observations\n")
	for _, observation := range issueObservations {
		fmt.Fprintf(&builder, "\n### `%s`\n\n", observation.ID)
		builder.WriteString("| Field | Value |\n| --- | --- |\n")
		writeCandidateRow(&builder, "Status", observation.Status)
		writeCandidateRow(&builder, "Environment", observation.EnvironmentID)
		writeCandidateRow(&builder, "Started", formatEvidenceTime(observation.StartedAt))
		writeCandidateRow(&builder, "Completed", formatEvidenceTime(observation.CompletedAt))
		fmt.Fprintf(&builder, "| Summary | %s |\n", markdownCell(observation.Summary))
		metrics, err := json.MarshalIndent(observation.Metrics, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("observation %q metrics: %w", observation.ID, err)
		}
		builder.WriteString("\nMetrics:\n\n")
		for _, line := range strings.Split(string(metrics), "\n") {
			builder.WriteString("    ")
			builder.WriteString(line)
			builder.WriteString("\n")
		}
		builder.WriteString("\nAttachments in the evidence archive:\n")
		if len(observation.AttachmentPaths) == 0 {
			builder.WriteString("\n- None declared for this observation.\n")
		} else {
			for _, attachment := range observation.AttachmentPaths {
				fmt.Fprintf(&builder, "\n- `%s`", attachment)
			}
			builder.WriteString("\n")
		}
	}
	builder.WriteString("\n## Disposition boundary\n\n")
	builder.WriteString("This record supports individual acceptance review. It does not itself comment on or close the issue, authorize promotion, or replace review of the issue's complete acceptance criteria.\n")
	return []byte(builder.String()), nil
}

func writeCandidateRow(builder *strings.Builder, field, value string) {
	fmt.Fprintf(builder, "| %s | `%s` |\n", field, markdownCell(value))
}

func markdownCell(value string) string {
	value = strings.ReplaceAll(value, "&", "&amp;")
	value = strings.ReplaceAll(value, "<", "&lt;")
	value = strings.ReplaceAll(value, ">", "&gt;")
	value = strings.ReplaceAll(value, "@", "&#64;")
	value = strings.ReplaceAll(value, "`", "&#96;")
	value = strings.ReplaceAll(value, "|", `\|`)
	value = strings.ReplaceAll(value, "\r\n", "<br>")
	value = strings.ReplaceAll(value, "\r", "<br>")
	return strings.ReplaceAll(value, "\n", "<br>")
}

func formatEvidenceTime(value time.Time) string {
	return value.UTC().Format(time.RFC3339Nano)
}

func validRepository(repository string) bool {
	parts := strings.Split(repository, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return false
	}
	for _, part := range parts {
		for _, char := range part {
			if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
				(char >= '0' && char <= '9') || char == '-' || char == '_' || char == '.' {
				continue
			}
			return false
		}
	}
	return true
}

func pathContainsSymlink(name string) (bool, error) {
	absolute, err := filepath.Abs(name)
	if err != nil {
		return false, err
	}
	volume := filepath.VolumeName(absolute)
	remainder := strings.TrimPrefix(absolute, volume)
	current := volume
	if strings.HasPrefix(remainder, string(filepath.Separator)) {
		current += string(filepath.Separator)
		remainder = strings.TrimLeft(remainder, string(filepath.Separator))
	}
	for _, component := range strings.Split(remainder, string(filepath.Separator)) {
		if component == "" {
			continue
		}
		current = filepath.Join(current, component)
		info, err := os.Lstat(current)
		if err != nil {
			return false, err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return true, nil
		}
	}
	return false, nil
}
