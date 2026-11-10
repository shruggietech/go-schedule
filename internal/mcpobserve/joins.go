package mcpobserve

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

const joinDefinitionsURI = "goschedule://joins/definitions"
const joinInstancesURI = "goschedule://joins/instances"

type joinReader interface {
	ListJoinsPage(context.Context, int, int, int, int, int) ([]domain.ParallelJoin, bool, error)
	GetJoinPage(context.Context, string, int, int, int) (domain.ParallelJoin, error)
	ListJoinInstancesPage(context.Context, domain.ParallelJoinInstanceFilter) ([]domain.ParallelJoinInstance, bool, error)
	GetJoinInstancePage(context.Context, string, domain.ParallelJoinInstanceFilter) (domain.ParallelJoinInstance, error)
}

func (a *adapter) joinEnvelope(ctx context.Context, kind, base string, offset int) (Envelope, error) {
	r, ok := a.client.(joinReader)
	if !ok {
		return Envelope{}, errors.New("unavailable: parallel join observations are unavailable")
	}
	e := Envelope{SchemaVersion: SchemaVersion, Permission: Permission, GeneratedAt: timestamp(a.now()), Trust: TrustNotice, UntrustedFields: []string{"data[].name", "data[].definition_name", "data[].branches[].name", "data[].reason", "data[].branches[].reason"}}
	items := []map[string]any{}
	more := false
	if kind == "join_definitions" {
		page, next, err := r.ListJoinsPage(ctx, offset, PageLimit, 0, PageLimit, TextLimit)
		if err != nil {
			return e, err
		}
		more = next
		for _, v := range page {
			items = append(items, joinDefinitionData(v))
		}
	} else {
		page, next, err := r.ListJoinInstancesPage(ctx, domain.ParallelJoinInstanceFilter{Offset: offset, Limit: PageLimit, BranchLimit: PageLimit, TextLimit: TextLimit})
		if err != nil {
			return e, err
		}
		more = next
		for _, v := range page {
			items = append(items, joinInstanceData(v))
		}
	}
	e.Data = items
	e.Page = &Page{Count: len(items), Limit: PageLimit}
	if more {
		e.Page.NextURI = base + "/page/" + encodeCursor(offset+len(items))
	}
	return e, nil
}
func joinDefinitionData(v domain.ParallelJoin) map[string]any {
	branches := []map[string]any{}
	for _, b := range v.Branches {
		branches = append(branches, map[string]any{"id": b.ID, "name": b.Name, "task_id": b.TaskID, "task_name": b.TaskName, "truncated_fields": b.TruncatedFields})
	}
	out := map[string]any{"id": v.ID, "name": v.Name, "entry_task_id": v.EntryTaskID, "entry_task_name": v.EntryTaskName, "target_task_id": v.TargetTaskID, "target_task_name": v.TargetTaskName, "branches": branches, "branch_count": v.BranchCount, "branch_offset": v.BranchOffset, "has_more_branches": v.HasMoreBranches, "created_at": timestamp(v.CreatedAt), "updated_at": timestamp(v.UpdatedAt), "truncated_fields": v.TruncatedFields}
	if v.HasMoreBranches {
		out["branches_next_uri"] = joinDefinitionsURI + "/" + url.PathEscape(v.ID) + "/branches/" + encodeCursor(v.BranchOffset+len(v.Branches))
	}
	return out
}
func joinInstanceData(v domain.ParallelJoinInstance) map[string]any {
	branches := []map[string]any{}
	for _, b := range v.Branches {
		branches = append(branches, map[string]any{"id": b.ID, "branch_id": b.BranchID, "name": b.Name, "task_id": b.TaskID, "task_name": b.TaskName, "delivery_id": b.DeliveryID, "delivery_state": b.DeliveryState, "run_id": b.RunID, "outcome": b.Outcome, "reason": b.Reason, "started_at": optionalTimestamp(b.StartedAt), "ended_at": optionalTimestamp(b.EndedAt), "cancellation_status": b.CancellationStatus, "cancellation_error": b.CancellationError, "truncated_fields": b.TruncatedFields})
	}
	out := map[string]any{"id": v.ID, "definition_id": v.DefinitionID, "definition_name": v.DefinitionName, "entry_task_id": v.EntryTaskID, "entry_task_name": v.EntryTaskName, "entry_run_id": v.EntryRunID, "target_task_id": v.TargetTaskID, "target_task_name": v.TargetTaskName, "state": v.State, "reason": v.Reason, "branches": branches, "branch_count": v.BranchCount, "branch_offset": v.BranchOffset, "has_more_branches": v.HasMoreBranches, "target_delivery_id": v.TargetDeliveryID, "target_delivery_state": v.TargetDeliveryState, "target_run_id": v.TargetRunID, "target_outcome": v.TargetOutcome, "target_started_at": optionalTimestamp(v.TargetStartedAt), "target_ended_at": optionalTimestamp(v.TargetEndedAt), "target_cancellation_status": v.TargetCancellationStatus, "target_cancellation_error": v.TargetCancellationError, "created_at": timestamp(v.CreatedAt), "updated_at": timestamp(v.UpdatedAt), "completed_at": optionalTimestamp(v.CompletedAt), "truncated_fields": v.TruncatedFields}
	if v.HasMoreBranches {
		out["branches_next_uri"] = joinInstancesURI + "/" + url.PathEscape(v.ID) + "/branches/" + encodeCursor(v.BranchOffset+len(v.Branches))
	}
	return out
}
func (a *adapter) readJoinExact(base string) mcp.ResourceHandler {
	return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		r, ok := a.client.(joinReader)
		if !ok {
			return nil, errors.New("unavailable: parallel join observations are unavailable")
		}
		suffix, ok := strings.CutPrefix(req.Params.URI, base+"/")
		if !ok {
			return nil, errors.New("invalid_resource: resource URI is not approved")
		}
		parts := strings.Split(suffix, "/")
		if len(parts) != 1 && (len(parts) != 3 || parts[1] != "branches") {
			return nil, errors.New("invalid_resource: resource URI is not approved")
		}
		id, err := url.PathUnescape(parts[0])
		if err != nil || id == "" || id == "page" {
			return nil, errors.New("invalid_resource: exact join ID is required")
		}
		offset := 0
		if len(parts) == 3 {
			offset, err = decodeCursor(parts[2])
			if err != nil || offset <= 0 {
				return nil, errInvalidCursor
			}
		}
		ctx, cancel := context.WithTimeout(ctx, a.timeout)
		defer cancel()
		e := Envelope{SchemaVersion: SchemaVersion, Permission: Permission, GeneratedAt: timestamp(a.now()), Trust: TrustNotice, UntrustedFields: []string{"data.name", "data.definition_name", "data.reason", "data.branches[].name", "data.branches[].reason"}}
		if base == joinDefinitionsURI {
			v, err := r.GetJoinPage(ctx, id, offset, PageLimit, TextLimit)
			if err != nil {
				return nil, safeError(err)
			}
			e.Data = joinDefinitionData(v)
		} else {
			v, err := r.GetJoinInstancePage(ctx, id, domain.ParallelJoinInstanceFilter{BranchOffset: offset, BranchLimit: PageLimit, TextLimit: TextLimit})
			if err != nil {
				return nil, safeError(err)
			}
			e.Data = joinInstanceData(v)
		}
		body, err := json.Marshal(e)
		if err != nil {
			return nil, errors.New("internal_error: resource could not be encoded")
		}
		return &mcp.ReadResourceResult{Contents: []*mcp.ResourceContents{{URI: req.Params.URI, MIMEType: "application/json", Text: string(body)}}}, nil
	}
}
func addJoinTemplates(s *mcp.Server, a *adapter) {
	for _, base := range []string{joinDefinitionsURI, joinInstancesURI} {
		name := "join_definition"
		if base == joinInstancesURI {
			name = "join_instance"
		}
		handler := a.readJoinExact(base)
		s.AddResourceTemplate(&mcp.ResourceTemplate{URITemplate: base + "/{id}", Name: name, Title: "Exact parallel join observation", Description: "One exact definition or instance with bounded branches and a continuation URI.", MIMEType: "application/json"}, handler)
		s.AddResourceTemplate(&mcp.ResourceTemplate{URITemplate: base + "/{id}/branches/{cursor}", Name: name + "_branches", Title: "Parallel join branch page", Description: "Continue exact branch evidence without inferring omitted contributions.", MIMEType: "application/json"}, handler)
	}
}
