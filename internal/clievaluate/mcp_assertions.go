package clievaluate

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func toolObject(result *mcp.CallToolResult) map[string]any {
	data, _ := json.Marshal(result.StructuredContent)
	var object map[string]any
	_ = json.Unmarshal(data, &object)
	return object
}

func (s *session) manageProtocol(ctx context.Context, connection *mcp.ClientSession) (passed bool, observation map[string]any, resultErr error) {
	observations := map[string]any{}
	input := map[string]any{"daemon_id": s.w.DaemonID, "request_id": uuid.NewString(), "action": "create", "create": map[string]any{"name": "QA confirmed MCP group"}}
	before, err := s.api.ListGroups(ctx)
	if err != nil {
		return false, observations, err
	}
	first, err := connection.CallTool(ctx, &mcp.CallToolParams{Name: "groups_manage", Arguments: input})
	if err != nil {
		return false, observations, err
	}
	observations["unconfirmed"] = first
	after, err := s.api.ListGroups(ctx)
	if err != nil {
		return false, observations, err
	}
	good := equalObservation(before, after) && first.IsError && toolObject(first)["outcome"] == "rejected"
	input["confirmed"] = true
	input["request_id"] = uuid.NewString()
	accepted, err := connection.CallTool(ctx, &mcp.CallToolParams{Name: "groups_manage", Arguments: input})
	if err != nil {
		return false, observations, err
	}
	observations["confirmed"] = accepted
	id, _ := toolObject(accepted)["object_id"].(string)
	good = good && id != "" && !accepted.IsError && toolObject(accepted)["outcome"] == "accepted"
	if id != "" {
		defer func() {
			cleanupCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			err := s.w.Guard(cleanupCtx, s.api)
			if err == nil {
				err = s.api.DeleteGroup(cleanupCtx, id)
			}
			resultErr = errors.Join(resultErr, err)
			if err != nil {
				passed = false
				observation["cleanup_error"] = err.Error()
			}
		}()
	}
	groups, err := s.api.ListGroups(ctx)
	if err != nil {
		return false, observations, err
	}
	found := false
	for _, group := range groups {
		found = found || group.ID == id
	}
	good = good && found && len(groups) == len(before)+1
	retry, err := connection.CallTool(ctx, &mcp.CallToolParams{Name: "groups_manage", Arguments: input})
	if err != nil {
		return false, observations, err
	}
	observations["identical_retry"] = retry
	repeated, err := s.api.ListGroups(ctx)
	if err != nil {
		return false, observations, err
	}
	good = good && toolObject(retry)["object_id"] == id && equalObservation(groups, repeated)
	input["create"] = map[string]any{"name": "QA changed retry must not create"}
	changed, err := connection.CallTool(ctx, &mcp.CallToolParams{Name: "groups_manage", Arguments: input})
	if err != nil {
		return false, observations, err
	}
	observations["changed_retry"] = changed
	unchanged, err := s.api.ListGroups(ctx)
	good = good && changed.IsError && equalObservation(groups, unchanged)
	return good, observations, err
}

func (s *session) operateProtocol(ctx context.Context, connection *mcp.ClientSession) (bool, any, error) {
	request := map[string]any{"daemon_id": s.w.DaemonID, "request_id": uuid.NewString(), "task_id": s.vars["base_task"]}
	response, err := connection.CallTool(ctx, &mcp.CallToolParams{Name: "tasks_run_now", Arguments: request})
	if err != nil {
		return false, nil, err
	}
	wrong := map[string]any{"daemon_id": uuid.NewString(), "request_id": uuid.NewString(), "task_id": s.vars["base_task"]}
	rejected, e := connection.CallTool(ctx, &mcp.CallToolParams{Name: "tasks_run_now", Arguments: wrong})
	err = errors.Join(err, e)
	return !response.IsError && toolObject(response)["outcome"] == "accepted" && rejected != nil && rejected.IsError, map[string]any{"authorized": response, "wrong_daemon": rejected}, err
}
