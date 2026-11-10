package mcpmanage

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

// JoinInput changes one parallel block through ordinary Manage authority.
type JoinInput struct {
	Envelope
	Action   string      `json:"action" jsonschema:"create, update, or delete"`
	ObjectID string      `json:"object_id,omitempty"`
	Create   *JoinCreate `json:"create,omitempty"`
	Update   *JoinUpdate `json:"update,omitempty"`
}

func (e *Executor) joinHandler() mcp.ToolHandlerFor[JoinInput, Result] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input JoinInput) (*mcp.CallToolResult, Result, error) {
		op := "joins." + input.Action
		if msg := validateAction(input.Action, input.ObjectID, input.Create != nil, input.Update != nil); msg != "" {
			return toolResult(rejected(input.Envelope, op, "join", input.ObjectID, msg))
		}
		invoke := func(ctx context.Context) (string, error) {
			switch input.Action {
			case "create":
				v, err := e.client.CreateJoin(ctx, input.Create.request())
				return v.ID, err
			case "update":
				v, err := e.client.UpdateJoin(ctx, input.ObjectID, input.Update.request())
				return v.ID, err
			default:
				return input.ObjectID, e.client.DeleteJoin(ctx, input.ObjectID)
			}
		}
		return toolResult(e.execute(ctx, input.Envelope, op, "join", input.ObjectID, fingerprint(input), "Parallel join mutation accepted.", invoke))
	}
}

// JoinBranch names one writable slot; task execution inputs remain absent.
type JoinBranch struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name"`
	TaskID string `json:"task_id"`
}

// JoinCreate carries only writable parallel block fields.
type JoinCreate struct {
	Name         string       `json:"name"`
	EntryTaskID  string       `json:"entry_task_id"`
	TargetTaskID string       `json:"target_task_id"`
	Branches     []JoinBranch `json:"branches"`
}

// JoinUpdate preserves omission and revision preconditions for one block.
type JoinUpdate struct {
	Name              *string       `json:"name,omitempty"`
	EntryTaskID       *string       `json:"entry_task_id,omitempty"`
	TargetTaskID      *string       `json:"target_task_id,omitempty"`
	Branches          *[]JoinBranch `json:"branches,omitempty"`
	ExpectedUpdatedAt *time.Time    `json:"expected_updated_at,omitempty"`
}

func joinBranches(in []JoinBranch) []domain.ParallelJoinBranch {
	out := make([]domain.ParallelJoinBranch, len(in))
	for i, v := range in {
		out[i] = domain.ParallelJoinBranch{ID: v.ID, Name: v.Name, TaskID: v.TaskID}
	}
	return out
}
func (v JoinCreate) request() server.JoinCreateRequest {
	return server.JoinCreateRequest{Name: v.Name, EntryTaskID: v.EntryTaskID, TargetTaskID: v.TargetTaskID, Branches: joinBranches(v.Branches)}
}
func (v JoinUpdate) request() server.JoinUpdateRequest {
	r := server.JoinUpdateRequest{Name: v.Name, EntryTaskID: v.EntryTaskID, TargetTaskID: v.TargetTaskID, ExpectedUpdatedAt: v.ExpectedUpdatedAt}
	if v.Branches != nil {
		b := joinBranches(*v.Branches)
		r.Branches = &b
	}
	return r
}
