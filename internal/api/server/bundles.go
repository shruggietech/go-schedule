package server

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/shruggietech/go-schedule/internal/bundle"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/scheduleinput"
	"github.com/shruggietech/go-schedule/internal/store"
	"github.com/shruggietech/go-schedule/internal/timezone"
)

const (
	bundlePlanLifetime = 15 * time.Minute
	maxBundlePlans     = 8
)

// BundleRequest carries an untrusted portable document for validation or review.
type BundleRequest struct {
	Bundle bundle.Document `json:"bundle"`
}

// BundleApplyRequest must echo the exact preview supplied by the daemon. A caller cannot turn a validation result into an apply authorization.
type BundleApplyRequest struct {
	PlanID            string `json:"plan_id"`
	BundleDigest      string `json:"bundle_digest"`
	TargetDaemonID    string `json:"target_daemon_id"`
	TargetFingerprint string `json:"target_fingerprint"`
}

type BundleApplyResponse struct {
	Plan     bundle.Plan   `json:"plan"`
	Outcomes []bundle.Item `json:"outcomes"`
}

type storedBundlePlan struct {
	Document bundle.Document
	Plan     bundle.Plan
	Created  time.Time
}

func (s *Server) handleBundleExport(w http.ResponseWriter, _ *http.Request) {
	doc, err := s.bundleDocument()
	if err != nil {
		s.internal(w, err)
		return
	}
	canonical, _, _, _ := bundle.Canonicalize(doc)
	writeJSON(w, http.StatusOK, canonical)
}

func (s *Server) handleBundleValidate(w http.ResponseWriter, r *http.Request) {
	req, err := decodeBundleRequest(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON")
		return
	}
	canonical, _, digest, issues := bundle.Canonicalize(req.Bundle)
	issues = append(issues, bundleCompatibilityIssues(canonical)...)
	writeJSON(w, http.StatusOK, map[string]any{"bundle": canonical, "digest": digest, "issues": issues, "valid": len(issues) == 0})
}

func (s *Server) handleBundlePreview(w http.ResponseWriter, r *http.Request) {
	s.handleBundlePlan(w, r, true)
}
func (s *Server) handleBundleCompare(w http.ResponseWriter, r *http.Request) {
	s.handleBundlePlan(w, r, false)
}

func (s *Server) handleBundlePlan(w http.ResponseWriter, r *http.Request, retain bool) {
	req, err := decodeBundleRequest(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON")
		return
	}
	canonical, _, _, issues := bundle.Canonicalize(req.Bundle)
	issues = append(issues, bundleCompatibilityIssues(canonical)...)
	if len(issues) > 0 {
		writeError(w, http.StatusBadRequest, CodeValidation, "bundle", "bundle is invalid")
		return
	}
	target, err := s.bundleDocument()
	if err != nil {
		s.internal(w, err)
		return
	}
	issues = append(issues, bundle.ValidateProjectedChains(canonical, target)...)
	if len(issues) > 0 {
		writeError(w, http.StatusBadRequest, CodeValidation, "bundle", "bundle is incompatible with the target")
		return
	}
	identity, err := s.store.DaemonIdentity()
	if err != nil {
		s.internal(w, err)
		return
	}
	plan := bundle.Preview(bundlePlanID(), identity.InstallationID, canonical, target)
	if retain {
		s.rememberBundlePlan(canonical, plan)
	}
	writeJSON(w, http.StatusOK, plan)
}

func decodeBundleRequest(body io.Reader) (BundleRequest, error) {
	var req BundleRequest
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		return BundleRequest{}, err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return BundleRequest{}, errors.New("body contains more than one JSON value")
		}
		return BundleRequest{}, err
	}
	return req, nil
}

func bundleCompatibilityIssues(doc bundle.Document) []bundle.Issue {
	issues := []bundle.Issue{}
	for _, task := range doc.Tasks {
		if _, err := timezone.Resolve(task.Timezone); err != nil {
			issues = append(issues, bundle.Issue{Kind: "task", Identity: task.PortableID, Message: "timezone is unsupported on target"})
		}
		if task.Schedule != "" {
			if _, err := scheduleinput.Parse(task.Schedule, scheduleinput.Syntax(task.ScheduleSyntax), task.Timezone, time.Now().UTC()); err != nil {
				issues = append(issues, bundle.Issue{Kind: "task", Identity: task.PortableID, Message: "schedule is unsupported on target"})
			}
		}
		if !validOverlap(domain.OverlapPolicy(task.OverlapPolicy)) || (domain.CatchupPolicy(task.CatchupPolicy) != domain.CatchupOne && domain.CatchupPolicy(task.CatchupPolicy) != domain.CatchupNone) || !validMissingDate(domain.MissingDatePolicy(task.MissingDatePolicy)) || !validTimeBasis(domain.TimeBasis(task.TimeBasis)) || !validDSTGap(domain.DSTGapPolicy(task.DSTGapPolicy)) || !validDSTOverlap(domain.DSTOverlapPolicy(task.DSTOverlapPolicy)) {
			issues = append(issues, bundle.Issue{Kind: "task", Identity: task.PortableID, Message: "task policy is unsupported on target"})
		}
	}
	return issues
}

func (s *Server) handleBundleApply(w http.ResponseWriter, r *http.Request) {
	var req BundleApplyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, CodeValidation, "body", "invalid JSON")
		return
	}
	stored, ok := s.takeBundlePlan(req)
	if !ok {
		writeError(w, http.StatusConflict, CodeConflict, "plan", "preview is missing or expired; preview again before apply")
		return
	}
	target, err := s.bundleDocument()
	if err != nil {
		s.internal(w, err)
		return
	}
	identity, err := s.store.DaemonIdentity()
	if err != nil {
		s.internal(w, err)
		return
	}
	if stored.Plan.TargetDaemonID != identity.InstallationID || stored.Plan.TargetFingerprint != bundle.Fingerprint(target) {
		writeError(w, http.StatusConflict, CodeConflict, "plan", "target changed after preview; preview again before apply")
		return
	}
	outcomes := s.applyBundle(stored.Document, stored.Plan)
	s.reload()
	writeJSON(w, http.StatusOK, BundleApplyResponse{Plan: stored.Plan, Outcomes: outcomes})
}

func (s *Server) rememberBundlePlan(doc bundle.Document, plan bundle.Plan) {
	s.bundlePlanMu.Lock()
	defer s.bundlePlanMu.Unlock()
	now := time.Now().UTC()
	for id, stored := range s.bundlePlans {
		if now.Sub(stored.Created) > bundlePlanLifetime {
			delete(s.bundlePlans, id)
		}
	}
	for len(s.bundlePlans) >= maxBundlePlans {
		var oldestID string
		var oldest time.Time
		for id, stored := range s.bundlePlans {
			if oldestID == "" || stored.Created.Before(oldest) {
				oldestID, oldest = id, stored.Created
			}
		}
		delete(s.bundlePlans, oldestID)
	}
	s.bundlePlans[plan.ID] = storedBundlePlan{Document: doc, Plan: plan, Created: now}
}

func (s *Server) takeBundlePlan(req BundleApplyRequest) (storedBundlePlan, bool) {
	s.bundlePlanMu.Lock()
	defer s.bundlePlanMu.Unlock()
	stored, ok := s.bundlePlans[req.PlanID]
	if !ok || time.Since(stored.Created) > bundlePlanLifetime || stored.Plan.BundleDigest != req.BundleDigest || stored.Plan.TargetDaemonID != req.TargetDaemonID || stored.Plan.TargetFingerprint != req.TargetFingerprint {
		return storedBundlePlan{}, false
	}
	delete(s.bundlePlans, req.PlanID)
	return stored, true
}

// applyBundle creates or updates only portable intent. Imported tasks remain disabled drafts because commands and all other execution inputs are excluded.
func (s *Server) applyBundle(doc bundle.Document, plan bundle.Plan) []bundle.Item {
	outcomes := make([]bundle.Item, 0, len(plan.Items))
	blocked := s.detachBundleGroupParents(doc, plan, &outcomes)
	groups := append([]bundle.Group(nil), doc.Groups...)
	for len(groups) > 0 {
		remaining := make([]bundle.Group, 0, len(groups))
		progressed := false
		for _, group := range groups {
			if !bundlePlanChanges(plan, "group", group.PortableID) {
				continue
			}
			if blocked[group.PortableID] {
				continue
			}
			if group.ParentPortableID != "" {
				if _, err := s.store.ObjectIDForPortableID("group", group.ParentPortableID); errors.Is(err, store.ErrNotFound) {
					remaining = append(remaining, group)
					continue
				}
			}
			outcomes = append(outcomes, s.applyBundleGroup(group))
			progressed = true
		}
		if !progressed {
			for _, group := range remaining {
				outcomes = append(outcomes, bundle.Item{Kind: "group", PortableID: group.PortableID, Name: group.Name, Action: bundle.ActionFailed, Message: "parent group could not be resolved"})
			}
			break
		}
		groups = remaining
	}
	for _, task := range doc.Tasks {
		if bundlePlanChanges(plan, "task", task.PortableID) {
			outcomes = append(outcomes, s.applyBundleTask(task))
		}
	}
	outcomes = append(outcomes, s.applyBundleChains(doc.Chains, plan)...)
	for _, item := range plan.Items {
		if item.Action == bundle.ActionTargetOnly || item.Action == bundle.ActionUnchanged || item.Action == bundle.ActionConflict {
			outcomes = append(outcomes, item)
		}
	}
	return outcomes
}

func (s *Server) applyBundleChains(sources []bundle.Chain, plan bundle.Plan) []bundle.Item {
	outcomes := []bundle.Item{}
	updates := []domain.CompletionChain{}
	updateItems := []bundle.Item{}
	creates := []bundle.Chain{}
	for _, source := range sources {
		if !bundlePlanChanges(plan, "chain", source.PortableID) {
			continue
		}
		item := bundle.Item{Kind: "chain", PortableID: source.PortableID, Name: source.PortableID}
		sourceID, sourceErr := s.store.ObjectIDForPortableID("task", source.SourceTaskID)
		targetID, targetErr := s.store.ObjectIDForPortableID("task", source.TargetTaskID)
		if sourceErr != nil || targetErr != nil {
			item.Action, item.Message = bundle.ActionFailed, "chain task was not applied"
			outcomes = append(outcomes, item)
			continue
		}
		id, err := s.store.ObjectIDForPortableID("chain", source.PortableID)
		if errors.Is(err, store.ErrNotFound) {
			creates = append(creates, source)
			continue
		}
		if err != nil {
			item.Action, item.Message = bundle.ActionFailed, err.Error()
			outcomes = append(outcomes, item)
			continue
		}
		updates = append(updates, domain.CompletionChain{ID: id, SourceTaskID: sourceID, TargetTaskID: targetID, OnOutcome: domain.CompletionOutcome(source.OnOutcome)})
		updateItems = append(updateItems, item)
	}
	if err := s.store.ReplaceCompletionChains(updates); err != nil {
		for _, item := range updateItems {
			item.Action, item.Message = bundle.ActionFailed, err.Error()
			outcomes = append(outcomes, item)
		}
	} else {
		for _, item := range updateItems {
			item.Action = bundle.ActionApplied
			outcomes = append(outcomes, item)
		}
	}
	for _, source := range creates {
		outcomes = append(outcomes, s.applyBundleChain(source))
	}
	return outcomes
}

// detachBundleGroupParents removes old parent edges before installing the
// reviewed hierarchy. This prevents an otherwise-valid hierarchy reversal from
// being rejected because an intermediate state would temporarily be cyclic.
func (s *Server) detachBundleGroupParents(doc bundle.Document, plan bundle.Plan, outcomes *[]bundle.Item) map[string]bool {
	blocked := map[string]bool{}
	for _, source := range doc.Groups {
		if !bundlePlanChanges(plan, "group", source.PortableID) {
			continue
		}
		id, err := s.store.ObjectIDForPortableID("group", source.PortableID)
		if errors.Is(err, store.ErrNotFound) {
			continue
		}
		if err != nil {
			blocked[source.PortableID] = true
			*outcomes = append(*outcomes, bundle.Item{Kind: "group", PortableID: source.PortableID, Name: source.Name, Action: bundle.ActionFailed, Message: err.Error()})
			continue
		}
		group, err := s.store.GetGroup(id)
		if err != nil || group.ParentID == "" {
			continue
		}
		desiredParent, err := s.groupObjectID(source.ParentPortableID)
		if err != nil || group.ParentID == desiredParent {
			continue
		}
		if err := s.store.SetGroupParent(id, ""); err != nil {
			blocked[source.PortableID] = true
			*outcomes = append(*outcomes, bundle.Item{Kind: "group", PortableID: source.PortableID, Name: source.Name, Action: bundle.ActionFailed, Message: err.Error()})
		}
	}
	return blocked
}

func bundlePlanChanges(plan bundle.Plan, kind, portableID string) bool {
	for _, item := range plan.Items {
		if item.Kind == kind && item.PortableID == portableID {
			return item.Action == bundle.ActionCreate || item.Action == bundle.ActionUpdate
		}
	}
	return false
}

func (s *Server) applyBundleGroup(source bundle.Group) bundle.Item {
	item := bundle.Item{Kind: "group", PortableID: source.PortableID, Name: source.Name}
	id, err := s.store.ObjectIDForPortableID("group", source.PortableID)
	parentID, parentErr := s.groupObjectID(source.ParentPortableID)
	if parentErr != nil {
		item.Action, item.Message = bundle.ActionFailed, "parent group was not applied"
		return item
	}
	if errors.Is(err, store.ErrNotFound) {
		group := &domain.Group{Name: source.Name, ParentID: parentID, Enabled: source.Enabled}
		if err := s.store.CreateGroup(group); err != nil {
			item.Action, item.Message = bundle.ActionFailed, err.Error()
			return item
		}
		if err := s.store.BindPortableID("group", group.ID, source.PortableID); err != nil {
			item.Action, item.Message = bundle.ActionUncertain, err.Error()
			return item
		}
		item.Action = bundle.ActionApplied
		return item
	}
	if err != nil {
		item.Action, item.Message = bundle.ActionFailed, err.Error()
		return item
	}
	group, err := s.store.GetGroup(id)
	if err != nil {
		item.Action, item.Message = bundle.ActionFailed, err.Error()
		return item
	}
	if group.Name != source.Name {
		if err := s.store.RenameGroup(id, source.Name); err != nil {
			item.Action, item.Message = bundle.ActionFailed, err.Error()
			return item
		}
	}
	if group.ParentID != parentID {
		if err := s.store.SetGroupParent(id, parentID); err != nil {
			item.Action, item.Message = bundle.ActionFailed, err.Error()
			return item
		}
	}
	if group.Enabled != source.Enabled {
		if err := s.store.SetGroupEnabled(id, source.Enabled); err != nil {
			item.Action, item.Message = bundle.ActionFailed, err.Error()
			return item
		}
	}
	item.Action = bundle.ActionApplied
	return item
}

func (s *Server) groupObjectID(portableID string) (string, error) {
	if portableID == "" {
		return "", nil
	}
	return s.store.ObjectIDForPortableID("group", portableID)
}

func (s *Server) applyBundleTask(source bundle.Task) bundle.Item {
	item := bundle.Item{Kind: "task", PortableID: source.PortableID, Name: source.Name}
	groupID, err := s.groupObjectID(source.GroupPortableID)
	if err != nil {
		item.Action, item.Message = bundle.ActionFailed, "task group was not applied"
		return item
	}
	scheduleID, err := s.createBundleSchedule(source)
	if err != nil {
		item.Action, item.Message = bundle.ActionFailed, err.Error()
		return item
	}
	id, err := s.store.ObjectIDForPortableID("task", source.PortableID)
	if errors.Is(err, store.ErrNotFound) {
		task := bundleTask(source, groupID, scheduleID)
		if err := s.store.CreateTask(&task); err != nil {
			item.Action, item.Message = bundle.ActionFailed, err.Error()
			return item
		}
		if err := s.store.BindPortableID("task", task.ID, source.PortableID); err != nil {
			item.Action, item.Message = bundle.ActionUncertain, err.Error()
			return item
		}
		item.Action, item.Message = bundle.ActionApplied, "created as a disabled draft; configure a local command before enabling"
		return item
	}
	if err != nil {
		item.Action, item.Message = bundle.ActionFailed, err.Error()
		return item
	}
	task, err := s.store.GetTask(id)
	if err != nil {
		item.Action, item.Message = bundle.ActionFailed, err.Error()
		return item
	}
	updated := bundleTask(source, groupID, scheduleID)
	updated.ID, updated.Command, updated.Args, updated.WorkingDir, updated.Env, updated.Stdin, updated.RunAs, updated.State, updated.Enabled = task.ID, task.Command, task.Args, task.WorkingDir, task.Env, task.Stdin, task.RunAs, task.State, source.Enabled
	if err := s.store.UpdateTask(&updated); err != nil {
		item.Action, item.Message = bundle.ActionFailed, err.Error()
		return item
	}
	item.Action, item.Message = bundle.ActionApplied, "updated portable intent; local execution inputs were preserved"
	return item
}

func bundleTask(source bundle.Task, groupID, scheduleID string) domain.Task {
	return domain.Task{Name: source.Name, GroupID: groupID, Timezone: source.Timezone, ScheduleID: scheduleID, OverlapPolicy: domain.OverlapPolicy(source.OverlapPolicy), CatchupPolicy: domain.CatchupPolicy(source.CatchupPolicy), MissingDatePolicy: domain.MissingDatePolicy(source.MissingDatePolicy), TimeBasis: domain.TimeBasis(source.TimeBasis), DSTGapPolicy: domain.DSTGapPolicy(source.DSTGapPolicy), DSTOverlapPolicy: domain.DSTOverlapPolicy(source.DSTOverlapPolicy), State: domain.TaskActive}
}

func (s *Server) createBundleSchedule(source bundle.Task) (string, error) {
	if source.Schedule == "" {
		return "", nil
	}
	input, err := scheduleinput.Parse(source.Schedule, scheduleinput.Syntax(source.ScheduleSyntax), source.Timezone, time.Now().UTC())
	if err != nil {
		return "", err
	}
	if err := s.store.CreateSchedule(&input.Schedule); err != nil {
		return "", err
	}
	return input.Schedule.ID, nil
}

func (s *Server) applyBundleChain(source bundle.Chain) bundle.Item {
	item := bundle.Item{Kind: "chain", PortableID: source.PortableID, Name: source.PortableID}
	sourceID, err := s.store.ObjectIDForPortableID("task", source.SourceTaskID)
	if err != nil {
		item.Action, item.Message = bundle.ActionFailed, "source task was not applied"
		return item
	}
	targetID, err := s.store.ObjectIDForPortableID("task", source.TargetTaskID)
	if err != nil {
		item.Action, item.Message = bundle.ActionFailed, "target task was not applied"
		return item
	}
	id, err := s.store.ObjectIDForPortableID("chain", source.PortableID)
	chain := domain.CompletionChain{ID: id, SourceTaskID: sourceID, TargetTaskID: targetID, OnOutcome: domain.CompletionOutcome(source.OnOutcome)}
	if errors.Is(err, store.ErrNotFound) {
		chain.ID = ""
		if err := s.store.CreateCompletionChain(&chain); err != nil {
			item.Action, item.Message = bundle.ActionFailed, err.Error()
			return item
		}
		if err := s.store.BindPortableID("chain", chain.ID, source.PortableID); err != nil {
			item.Action, item.Message = bundle.ActionUncertain, err.Error()
			return item
		}
		item.Action = bundle.ActionApplied
		return item
	}
	if err != nil {
		item.Action, item.Message = bundle.ActionFailed, err.Error()
		return item
	}
	if err := s.store.UpdateCompletionChain(&chain); err != nil {
		item.Action, item.Message = bundle.ActionFailed, err.Error()
		return item
	}
	item.Action = bundle.ActionApplied
	return item
}

func (s *Server) bundleDocument() (bundle.Document, error) {
	doc := bundle.Document{Schema: bundle.SchemaV1, Groups: []bundle.Group{}, Tasks: []bundle.Task{}, Chains: []bundle.Chain{}, Exclusions: []bundle.Issue{}}
	groups, err := s.store.ListGroups()
	if err != nil {
		return doc, err
	}
	groupIDs := map[string]string{}
	for _, group := range groups {
		id, err := s.store.PortableID("group", group.ID)
		if err != nil {
			return doc, err
		}
		groupIDs[group.ID] = id
	}
	for _, group := range groups {
		doc.Groups = append(doc.Groups, bundle.Group{PortableID: groupIDs[group.ID], Name: group.Name, ParentPortableID: groupIDs[group.ParentID], Enabled: group.Enabled})
	}
	tasks, err := s.store.ListTasks("", "")
	if err != nil {
		return doc, err
	}
	taskIDs := map[string]string{}
	for _, task := range tasks {
		id, err := s.store.PortableID("task", task.ID)
		if err != nil {
			return doc, err
		}
		taskIDs[task.ID] = id
	}
	for _, task := range tasks {
		value := bundle.Task{PortableID: taskIDs[task.ID], Name: task.Name, GroupPortableID: groupIDs[task.GroupID], Enabled: task.Enabled, Timezone: task.Timezone, OverlapPolicy: string(task.OverlapPolicy), CatchupPolicy: string(task.CatchupPolicy), MissingDatePolicy: string(task.MissingDatePolicy), TimeBasis: string(task.TimeBasis), DSTGapPolicy: string(task.DSTGapPolicy), DSTOverlapPolicy: string(task.DSTOverlapPolicy)}
		if task.ScheduleID != "" {
			schedule, err := s.store.GetSchedule(task.ScheduleID)
			if err != nil {
				return doc, err
			}
			value.Schedule, value.ScheduleSyntax = schedule.Expression, string(scheduleinput.SourceSyntax(schedule))
		}
		doc.Tasks = append(doc.Tasks, value)
		if task.Command != "" || task.WorkingDir != "" || len(task.Env) > 0 || task.RunAs != "" || task.Stdin != "" {
			doc.Exclusions = append(doc.Exclusions, bundle.Issue{Kind: "task_execution_input", Identity: value.PortableID, Message: "command and machine-specific execution inputs were excluded"})
		}
	}
	chains, err := s.store.ListCompletionChains()
	if err != nil {
		return doc, err
	}
	for _, chain := range chains {
		id, err := s.store.PortableID("chain", chain.ID)
		if err != nil {
			return doc, err
		}
		doc.Chains = append(doc.Chains, bundle.Chain{PortableID: id, SourceTaskID: taskIDs[chain.SourceTaskID], TargetTaskID: taskIDs[chain.TargetTaskID], OnOutcome: string(chain.OnOutcome)})
	}
	triggers, err := s.store.ListExternalTriggers()
	if err != nil {
		return doc, err
	}
	if len(triggers) > 0 {
		doc.Exclusions = append(doc.Exclusions, bundle.Issue{Kind: "external_trigger", Message: "external trigger keys were excluded"})
	}
	watchers, err := s.store.ListFilesystemWatchers()
	if err != nil {
		return doc, err
	}
	if len(watchers) > 0 {
		doc.Exclusions = append(doc.Exclusions, bundle.Issue{Kind: "filesystem_watcher", Message: "machine-specific watcher paths were excluded"})
	}
	return doc, nil
}

func bundlePlanID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return time.Now().UTC().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(bytes)
}
