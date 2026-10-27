package clievaluate

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/cli"
	"github.com/shruggietech/go-schedule/internal/clifixture"
	"github.com/shruggietech/go-schedule/internal/domain"
)

type execution struct {
	Code                         int
	Args                         []string
	Input, Out, Err, Logs, State []byte
	Checks                       []bool
	Blocked                      string
	Detail                       string
	Cleanup                      error
}

func (s *session) execute(ctx context.Context, definition step) (observed execution) {
	diagnostics := &executionDiagnostics{}
	defer func() { observed = diagnostics.finish(observed) }()
	finishIncident := s.incidentWatchdog(ctx, diagnostics, definition.Spec.ExpectedExit)
	defer func() { observed = finishIncident(observed) }()
	ctx = diagnostics.phase(ctx, "guard")
	if definition.Human {
		definition.Spec.ID = strings.TrimSuffix(definition.Spec.ID, ".human")
	}
	if definition.Hook == "readonly-output" {
		if err := s.readonlySetup(); err != nil {
			return execution{Blocked: "fixture_setup", Detail: err.Error()}
		}
		defer func() { _ = os.Chmod(s.vars["readonly_file"], 0600) }()
	}
	if err := s.w.Guard(ctx, s.api); err != nil {
		return execution{Blocked: "ownership_mismatch", Detail: err.Error()}
	}
	if handled, result := s.special(diagnostics.phase(ctx, "special"), definition); handled {
		diagnostics.CommandRan = result.Blocked == ""
		return result
	}
	ctx = diagnostics.phase(ctx, "setup")
	if definition.Spec.ID == "alerts.ack" {
		if err := s.before(ctx, definition); err != nil {
			return execution{Blocked: "fixture_setup", Detail: err.Error()}
		}
	}
	args, err := s.arguments(definition.Spec.Argv)
	if err != nil {
		return execution{Blocked: "fixture_prerequisite", Detail: err.Error()}
	}
	input, err := s.expand(definition.Input)
	if err != nil {
		return execution{Blocked: "fixture_prerequisite", Detail: err.Error()}
	}
	if err := s.before(ctx, definition); err != nil {
		return execution{Blocked: "fixture_setup", Detail: err.Error()}
	}
	if err := s.w.Guard(ctx, s.api); err != nil {
		return execution{Blocked: "ownership_mismatch", Detail: err.Error()}
	}
	metadata := slices.Contains([]string{"help", "alias", "completion", "documentation", "version"}, definition.Spec.ScenarioClass)
	var before map[string]any
	if !metadata {
		ctx = diagnostics.phase(ctx, "before")
		before, err = s.snapshot(ctx)
		diagnostics.BeforeComplete = err == nil
		if err != nil {
			return execution{Blocked: "observation_unavailable", Detail: err.Error(), State: mustObservation(map[string]any{"before": before, "before_complete": false})}
		}
	}
	ctx = diagnostics.phase(ctx, "command")
	diagnostics.CommandRan = true
	result := execution{Args: args, Input: []byte(input)}
	if definition.Spec.ExecutionKind == "command_harness" {
		var output, diagnostic bytes.Buffer
		s.recorder.Failure = nil
		if definition.Hook == "service-failure" || definition.Hook == "gui-failure" {
			s.recorder.Failure = errors.New("fixture adapter rejected the request")
		}
		result.Code = cli.ExecuteWithDependencies(ctx, args, s.dependencies(strings.NewReader(input), &output, &diagnostic))
		result.Out = output.Bytes()
		result.Err = diagnostic.Bytes()
		s.recorder.Failure = nil
	} else {
		capture := captureCommand(ctx, clifixture.Command{Executable: s.binaries["cli"], Arguments: args, Directory: s.w.Runtime, Environment: quietEnvironment()}, []byte(input), s.budget.ChannelBytes)
		result.Code = capture.Code
		result.Out = capture.Out
		result.Err = capture.Err
		result.Cleanup = capture.Cleanup
		if capture.Exceeded {
			result.Detail = "output_budget_exceeded"
		}
		if capture.Failure != nil {
			result.Detail = capture.Failure.Error()
		}
	}
	registerCapturedSecrets(s.redactor, definition, result.Out)
	registerCapturedSecrets(s.redactor, definition, result.Err)
	outputOK := true
	if definition.OutputContains != "" {
		outputOK = strings.Contains(string(result.Out)+string(result.Err), definition.OutputContains)
	} else if definition.Spec.ExpectedExit == 0 && definition.JSON {
		outputOK = json.Valid(result.Out)
	} else if definition.Spec.ExpectedExit != 0 {
		outputOK = len(result.Err)+len(result.Out) > 0
	}
	if result.Detail != "" || result.Cleanup != nil {
		outputOK = false
	}
	if metadata {
		result.Checks = []bool{outputOK}
		for len(result.Checks) < len(definition.Spec.Checks) {
			result.Checks = append(result.Checks, true)
		}
		result.State = mustObservation(map[string]any{"process_cleanup": result.Cleanup == nil, "metadata_only": true})
		return result
	}
	var body any
	if json.Valid(result.Out) {
		if err := json.Unmarshal(result.Out, &body); err != nil {
			outputOK = false
		}
	}
	ctx = diagnostics.phase(ctx, "output_inspection")
	if definition.Human && result.Code == 0 {
		observed, err := s.humanResponse(ctx, definition, result.Out, before)
		if err != nil {
			outputOK = false
			result.Detail = err.Error()
		} else {
			body = observed
		}
	}
	if strings.HasPrefix(definition.Spec.CommandPath, "trigger ") {
		registerTriggerKeys(s.redactor, body)
	}
	for _, save := range definition.Save {
		value, ok := fieldAt(body, save.Value)
		text, valid := value.(string)
		if !ok || !valid || text == "" {
			outputOK = false
			continue
		}
		s.vars[save.Name] = text
	}
	ctx = diagnostics.phase(ctx, "postcondition")
	postOK, postData, postErr := s.after(ctx, definition, result, body)
	if postErr != nil {
		postOK = false
		postData = map[string]any{"partial": postData, "error": postErr.Error()}
	}
	if definition.Hook == "mcp-http" {
		result.Logs = mustObservation(postData)
	}
	ctx = diagnostics.phase(ctx, "after")
	after, observationErr := s.snapshot(ctx)
	diagnostics.AfterComplete = observationErr == nil
	stateOK := observationErr == nil
	if stateOK {
		for _, assertion := range definition.Assert {
			path, err := s.expand(assertion.Path)
			if err != nil {
				stateOK = false
				continue
			}
			want, err := s.expand(assertion.Want)
			if err != nil {
				stateOK = false
				continue
			}
			value, ok := fieldAt(after, path)
			actual := "<absent>"
			if ok {
				actual = printable(value)
			}
			stateOK = stateOK && (actual == want || want == "<nonempty>" && ok && actual != "")
		}
		if definition.Spec.ExpectedExit != 0 || definition.Spec.ScenarioClass == "read" || definition.Hook == "unchanged" {
			stateOK = stateOK && equalReadObservation(definition, before, after)
		} else {
			stateOK = stateOK && s.persistedIntent(definition, args, before, after)
		}
	}
	stateEvidence := map[string]any{"before": observationIdentity(before), "after": observationIdentity(after), "before_complete": diagnostics.BeforeComplete, "after_complete": diagnostics.AfterComplete, "changes": observationChanges(before, after), "postcondition": postData, "observation_error": errorText(observationErr), "execution_diagnostic": result.Detail, "process_cleanup": result.Cleanup == nil}
	if observationErr != nil {
		stateEvidence["partial_after"] = after
	}
	result.State = mustObservation(stateEvidence)
	result.Checks = []bool{outputOK, stateOK, postOK}
	if len(definition.Spec.Checks) == 2 {
		result.Checks = result.Checks[:2]
		result.Checks[1] = stateOK && postOK
	}
	return result
}

func notificationFilter(channel string) domain.NotificationDeliveryFilter {
	return domain.NotificationDeliveryFilter{ChannelID: channel, Limit: 20}
}

func (s *session) observeRun(ctx context.Context, hook string) (bool, any, error) {
	variable := "base_task"
	switch hook {
	case "run-failure":
		variable = "failure_task"
	case "run-missing":
		variable = "missing_task"
	case "run-cap":
		variable = "output_task"
	case "run-edited":
		variable = "task"
	}
	var observed domain.Run
	err := waitUntil(ctx, func() bool {
		runs, err := s.api.ListRuns(ctx, s.vars[variable], 20)
		if err != nil {
			return false
		}
		for _, run := range runs {
			if run.EndedAt != nil && !s.priorRuns[run.ID] {
				observed = run
				return true
			}
		}
		return false
	})
	if err != nil {
		return false, observed, err
	}
	good := observed.ExitCode != nil && *observed.ExitCode == 0 && !observed.StartFailed
	switch hook {
	case "run-edited":
		good = good && strings.Contains(observed.Output, "space quote") && strings.Contains(observed.Output, "fixture stdin") && strings.Contains(observed.Output, "QA_VALUE=\"owned\"")
	case "run-success":
		s.vars["run"] = observed.ID
		good = good && strings.Contains(observed.Output, "fixture")
	case "run-failure":
		s.vars["failure_run"] = observed.ID
		good = observed.ExitCode != nil && *observed.ExitCode == 7 && !observed.StartFailed
	case "run-missing":
		s.vars["missing_run"] = observed.ID
		good = observed.StartFailed && observed.Output != ""
	case "run-cap":
		good = observed.OutputTruncated && len(observed.Output) <= int(s.w.Config.OutputCapBytes)
	case "chain-propagation":
		var child domain.Run
		err = waitUntil(ctx, func() bool {
			runs, e := s.api.ListRuns(ctx, s.vars["second_task"], 20)
			if e != nil {
				return false
			}
			for _, run := range runs {
				if run.SourceTaskID == s.vars["base_task"] && run.SourceRunID == observed.ID && run.EndedAt != nil {
					child = run
					return true
				}
			}
			return false
		})
		return good && err == nil, child, err
	}
	return good, observed, nil
}

func mustObservation(value any) []byte {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return []byte(`{"observation_error":"cannot encode product observation"}`)
	}
	return append(data, '\n')
}
func errorText(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
func equalObservation(a, b any) bool {
	first, err := json.Marshal(a)
	if err != nil {
		return false
	}
	second, err := json.Marshal(b)
	return err == nil && bytes.Equal(first, second)
}

func (s *session) before(ctx context.Context, definition step) error {
	if strings.HasPrefix(definition.Spec.ID, "smtp.") {
		if err := s.smtpBefore(ctx, definition); err != nil {
			return err
		}
	}
	if strings.HasPrefix(definition.Hook, "run-") || definition.Hook == "chain-propagation" {
		variable := "base_task"
		switch definition.Hook {
		case "run-failure":
			variable = "failure_task"
		case "run-missing":
			variable = "missing_task"
		case "run-cap":
			variable = "output_task"
		case "run-edited":
			variable = "task"
		}
		runs, err := s.api.ListRuns(ctx, s.vars[variable], 20)
		if err != nil {
			return err
		}
		s.priorRuns = map[string]bool{}
		for _, run := range runs {
			s.priorRuns[run.ID] = true
		}
	}
	if definition.Spec.ID == "alerts.ack" {
		alerts, err := s.api.ListAlerts(ctx, true)
		if err != nil {
			return err
		}
		if len(alerts) == 0 {
			return errors.New("owned failure scenario produced no acknowledgeable alert")
		}
		s.vars["alert"] = alerts[0].ID
	}
	if definition.Spec.ID == "profile.pair" {
		if _, ok := s.vars["phrase"]; !ok {
			return errors.New("pairing phrase was not captured from the owned creation case")
		}
	}
	if definition.Hook == "repair-preferences" {
		return s.writeOwned(filepath.Join(s.w.Runtime, "client", "preferences.json"), []byte("{malformed"))
	}
	if definition.Spec.ID == "runs.active" {
		return s.api.RunNow(ctx, s.vars["active_task"])
	}
	if definition.Hook == "watcher-event" {
		return s.writeOwned(filepath.Join(s.vars["watch_root"], "event.txt"), []byte("owned filesystem event\n"))
	}
	return nil
}

func (s *session) persistedIntent(definition step, args []string, before, after map[string]any) bool {
	path := definition.Spec.CommandPath
	root := strings.Fields(path)
	if len(root) == 0 {
		return true
	}
	if definition.Spec.ScenarioClass == "create" {
		if definition.Spec.ID == "profile.pair" {
			return !equalObservation(before["profiles"], after["profiles"])
		}
		if len(definition.Save) > 0 {
			for _, save := range definition.Save {
				if s.vars[save.Name] == "" {
					return false
				}
			}
		}
		if strings.HasPrefix(path, "cron ") {
			return !equalObservation(before["tasks"], after["tasks"])
		}
		return !equalObservation(before, after)
	}
	if definition.Spec.ScenarioClass == "delete" {
		return !equalObservation(before, after)
	}
	if definition.Spec.ScenarioClass == "host-adapter" {
		return !equalObservation(before["host_intents"], after["host_intents"])
	}
	if definition.Spec.ScenarioClass == "local-persistence" {
		return true
	}
	entity := map[string]string{"task": "tasks", "group": "groups", "chain": "chains", "trigger": "triggers", "watcher": "watchers", "actor": "actors", "credential": "credentials"}[root[0]]
	if root[0] == "trigger" && len(root) > 1 && root[1] == "set" {
		entity = "sets"
	}
	if entity == "" {
		return true
	}
	id := ""
	for i, arg := range args {
		if i > 0 && slices.Contains([]string{"edit", "update", "enable", "disable", "retarget", "revoke"}, args[i-1]) {
			id = arg
			break
		}
	}
	if id == "" {
		return true
	}
	if _, exists := fieldAt(after, entity+"."+id); !exists {
		return false
	}
	fields := map[string]string{"name": "name", "task": "target_task_id", "source": "source_task_id", "target": "target_task_id", "capability": "capability", "state": "state", "pattern": "pattern", "kind": "kind", "path": "path", "recursive": "recursive", "on": "on_outcome"}
	for i, arg := range args {
		if !strings.HasPrefix(arg, "--") {
			continue
		}
		name, value, hasValue := strings.Cut(strings.TrimPrefix(arg, "--"), "=")
		if !hasValue && i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
			value = args[i+1]
			hasValue = true
		}
		field, known := fields[name]
		if !known || !hasValue {
			continue
		}
		if root[0] == "actor" && field == "name" {
			field = "display_name"
		}
		actual, ok := fieldAt(after, entity+"."+id+"."+field)
		if !ok || printable(actual) != value {
			return false
		}
	}
	if strings.HasSuffix(path, " enable") || strings.HasSuffix(path, " disable") {
		if entity == "sets" {
			count, ok := fieldAt(after, entity+"."+id+".enabled_count")
			members, valid := fieldAt(after, entity+"."+id+".member_count")
			if strings.HasSuffix(path, " enable") {
				return ok && valid && count == members
			}
			return ok && count == float64(0)
		}
		actual, ok := fieldAt(after, entity+"."+id+".enabled")
		return ok && actual == strings.HasSuffix(path, " enable")
	}
	return true
}

func (s *session) after(ctx context.Context, definition step, result execution, body any) (bool, any, error) {
	if strings.HasPrefix(definition.Spec.ID, "smtp.") {
		return s.smtpPost(ctx, definition, result)
	}
	if definition.Spec.ID == "task.edit-inputs" && result.Code == 0 {
		detail, err := s.api.GetTask(ctx, s.vars["task"])
		task := detail.Task
		return len(task.Args) == 0 && len(task.Env) == 0 && task.Stdin == "" && task.WorkingDir == "" && task.RunAs == "" && task.ScheduleID == "" && task.GroupID == "", detail, err
	}
	if definition.Spec.ID == "actor.clear-expiration" && result.Code == 0 {
		actors, err := s.api.ListActors(ctx)
		for _, actor := range actors {
			if actor.ID == s.vars["actor"] {
				return actor.ExpiresAt == nil && actor.State == domain.ActorStateRevoked, actor, err
			}
		}
		return false, actors, err
	}
	if definition.Hook == "inherited-policy" {
		policy, err := s.api.EffectiveTaskNotificationPolicy(ctx, s.vars["base_task"])
		return policy.SourceScopeType == domain.NotificationScopeGroup && policy.SourceScopeID == s.vars["base_group"] && len(policy.Assignments) == 2, policy, err
	}
	if definition.Hook == "readonly-output" {
		data, err := os.ReadFile(s.vars["readonly_file"])
		return string(data) == "owned prior contents\n", map[string]any{"path": s.vars["readonly_file"], "preserved_sha256": Digest(data)}, err
	}
	if definition.Hook == "local-state" || definition.Hook == "repair-preferences" {
		return s.localPost(definition)
	}
	if definition.Hook == "mcp-http" && result.Code == 0 {
		return s.httpPost(ctx, definition, body)
	}
	if result.Code == 0 && (definition.Spec.CommandPath == "notification task set" || definition.Spec.CommandPath == "notification group set") {
		return s.assignmentPost(ctx, definition)
	}
	if definition.Spec.ID == "notification.channel-create" && result.Code != 0 {
		channel, err := s.api.CreateNotificationChannel(ctx, server.NotificationChannelCreateRequest{Name: "QA recovery destination", Endpoint: s.vars["webhook"]})
		if err == nil {
			s.vars["channel"] = channel.ID
		}
		return false, map[string]any{"failed_command_retained": true, "subsequent_case_fixture": channel}, err
	}
	if definition.Hook == "bundle-file" || definition.Hook == "plan-file" || definition.Hook == "audit-file" {
		path := s.vars["bundle_file"]
		if definition.Hook == "plan-file" {
			path = s.vars["plan_file"]
		}
		if definition.Hook == "audit-file" {
			path = s.vars["audit_file"]
		}
		data, err := os.ReadFile(path)
		if definition.Hook == "audit-file" {
			valid := true
			for _, line := range bytes.Split(bytes.TrimSpace(data), []byte("\n")) {
				if len(line) > 0 && !json.Valid(line) {
					valid = false
				}
			}
			return valid, map[string]any{"path": path, "bytes": len(data), "sha256": Digest(data)}, err
		}
		return json.Valid(data), map[string]any{"path": path, "bytes": len(data), "sha256": Digest(data)}, err
	}
	if definition.Spec.CommandPath == "pairing create" {
		if phrase, ok := fieldAt(body, "phrase"); ok {
			s.vars["phrase"], _ = phrase.(string)
			s.redactor.Register(s.vars["phrase"])
		}
	}
	if definition.Hook == "profile-pair" {
		if credential, ok := fieldAt(body, "credential_id"); ok {
			s.vars["credential"], _ = credential.(string)
		}
	}
	if definition.Hook == "identity-reset" {
		manifest, err := s.api.Manifest(ctx)
		return err == nil && manifest.InstallationID != "" && manifest.InstallationID != s.w.DaemonID, manifest, err
	}
	if definition.Hook == "run-success" || definition.Hook == "run-edited" || definition.Hook == "run-failure" || definition.Hook == "run-missing" || definition.Hook == "run-cap" || definition.Hook == "chain-propagation" {
		return s.observeRun(ctx, definition.Hook)
	}
	if definition.Hook == "watcher-event" {
		var observed any
		err := waitUntil(ctx, func() bool {
			runs, err := s.api.ListRuns(ctx, s.vars["base_task"], 20)
			if err != nil {
				return false
			}
			for _, run := range runs {
				if run.SourceWatcherID == s.vars["watcher"] {
					observed = run
					return true
				}
			}
			return false
		})
		return err == nil, observed, err
	}
	if definition.Hook == "webhook-delivery" {
		var deliveries any
		err := waitUntil(ctx, func() bool {
			items, err := s.api.ListNotificationDeliveries(ctx, notificationFilter(s.vars["channel"]))
			if err != nil {
				return false
			}
			deliveries = items
			for _, item := range items {
				if item.State == "succeeded" {
					return true
				}
			}
			return false
		})
		return err == nil, deliveries, err
	}
	if definition.Hook == "service-intent" || definition.Hook == "gui-intent" || definition.Hook == "service-failure" || definition.Hook == "gui-failure" {
		intents := s.recorder.Intents()
		operation := "gui.launch"
		if strings.HasPrefix(definition.Hook, "service") {
			parts := strings.Fields(definition.Spec.CommandPath)
			operation = "service." + parts[len(parts)-1]
		}
		found := false
		for _, intent := range intents {
			found = found || intent.Operation == operation
		}
		return found, intents, nil
	}
	if definition.Hook == "mcp-stdio" {
		return false, nil, errors.New("stdio protocol must use the owned interactive dispatcher")
	}
	if result.Code == definition.Spec.ExpectedExit {
		logs, err := s.api.ListLogs(ctx, "", 5)
		return err == nil, logs, err
	}
	return true, map[string]string{"result": "command exit mismatch is independently recorded"}, nil
}

func observationIdentity(value map[string]any) map[string]any {
	result := map[string]any{}
	for key, item := range value {
		data := mustObservation(item)
		result[key] = map[string]any{"sha256": Digest(data), "bytes": len(data)}
	}
	return result
}
func observationChanges(before, after map[string]any) map[string]any {
	result := map[string]any{}
	for key, value := range after {
		if equalObservation(before[key], value) {
			continue
		}
		old, oldOK := before[key].(map[string]any)
		next, nextOK := value.(map[string]any)
		if oldOK && nextOK {
			changed := map[string]any{}
			for id, item := range next {
				if !equalObservation(old[id], item) {
					changed[id] = item
				}
			}
			for id := range old {
				if _, ok := next[id]; !ok {
					changed[id] = nil
				}
			}
			result[key] = changed
		} else {
			result[key] = value
		}
	}
	return result
}

func equalReadObservation(definition step, before, after map[string]any) bool {
	normalize := func(snapshot map[string]any) map[string]any {
		data, _ := json.Marshal(snapshot)
		var result map[string]any
		_ = json.Unmarshal(data, &result)
		if definition.Spec.ScenarioClass == "host-adapter" {
			delete(result, "host_intents")
		}
		if credentials, ok := result["credentials"].(map[string]any); ok {
			for _, item := range credentials {
				if record, ok := item.(map[string]any); ok {
					delete(record, "last_used_at")
				}
			}
		}
		if collection, ok := result["profiles"].(map[string]any); ok {
			if profiles, ok := collection["profiles"].([]any); ok {
				for _, item := range profiles {
					if profile, ok := item.(map[string]any); ok {
						delete(profile, "last_successful_at")
						delete(profile, "updated_at")
					}
				}
			}
		}
		return result
	}
	return equalObservation(normalize(before), normalize(after))
}
