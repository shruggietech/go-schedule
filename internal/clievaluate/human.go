package clievaluate

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
)

func jsonObject(value any) (map[string]any, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	err = json.Unmarshal(data, &result)
	return result, err
}

func newIdentity(before, after map[string]any, entity string) (string, error) {
	old, _ := before[entity].(map[string]any)
	next, _ := after[entity].(map[string]any)
	found := ""
	for id := range next {
		if _, exists := old[id]; !exists {
			if found != "" {
				return "", errors.New("human-output identity observation is ambiguous")
			}
			found = id
		}
	}
	if found == "" {
		return "", errors.New("human-output command created no independently observable identity")
	}
	return found, nil
}

func humanLine(output, label string) string {
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(line)), strings.ToLower(label)) {
			return strings.TrimSpace(strings.TrimSpace(line)[len(label):])
		}
	}
	return ""
}

func (s *session) humanResponse(ctx context.Context, definition step, output []byte, before map[string]any) (any, error) {
	body := map[string]any{}
	if json.Valid(output) {
		if err := json.Unmarshal(output, &body); err != nil {
			body = map[string]any{}
		}
	}
	for _, save := range definition.Save {
		if save.Value != "id" && !strings.HasSuffix(save.Value, ".id") {
			continue
		}
		after, err := s.snapshot(ctx)
		if err != nil {
			return nil, err
		}
		parts := strings.Fields(definition.Spec.CommandPath)
		entity := map[string]string{"task": "tasks", "group": "groups", "chain": "chains", "actor": "actors", "pairing": "pairings", "trigger": "triggers", "watcher": "watchers"}[parts[0]]
		if len(parts) > 1 && parts[0] == "trigger" && parts[1] == "set" {
			entity = "sets"
		}
		id := ""
		if parts[0] == "profile" {
			old, _ := before["profiles"].(map[string]any)
			oldItems, _ := old["profiles"].([]any)
			known := map[string]bool{}
			for _, item := range oldItems {
				if profile, ok := item.(map[string]any); ok {
					value, _ := profile["id"].(string)
					known[value] = true
				}
			}
			collection, err := s.profiles.Load()
			if err != nil {
				return nil, err
			}
			for _, profile := range collection.Profiles {
				if !known[profile.ID] {
					id = profile.ID
					body["credential_id"] = profile.CredentialID
				}
			}
		} else {
			id, err = newIdentity(before, after, entity)
			if err != nil {
				return nil, err
			}
		}
		if id == "" {
			return nil, errors.New("human-output profile identity is unavailable")
		}
		if (parts[0] != "trigger" || len(parts) <= 1 || parts[1] != "set") && !strings.Contains(string(output), id) {
			return nil, errors.New("human command output does not identify the independently created object")
		}
		s.vars[save.Name] = id
		body["id"] = id
		for _, key := range []string{"task", "trigger", "trigger_set", "watcher"} {
			body[key] = map[string]any{"id": id}
		}
	}
	path := definition.Spec.CommandPath
	if path == "pairing create" {
		phrase := humanLine(string(output), "Phrase (shown once):")
		if phrase == "" {
			return nil, errors.New("human pairing output omitted its one-time phrase")
		}
		s.redactor.Register(phrase)
		body["phrase"] = phrase
	}
	if path == "trigger create" || path == "trigger rotate" || definition.Hook == "trigger-secret" {
		secret, err := s.api.RevealTrigger(ctx, s.vars["trigger"])
		if err != nil {
			return nil, err
		}
		s.redactor.Register(secret.Key)
		if !strings.Contains(string(output), secret.Key) {
			return nil, errors.New("human trigger output does not contain the daemon's current key")
		}
		return jsonObject(secret)
	}
	if path == "trigger set create" || path == "trigger set rotate" || path == "trigger set reveal" {
		id := s.vars["trigger_set"]
		for _, save := range definition.Save {
			if save.Value == "trigger_set.id" {
				id = s.vars[save.Name]
			}
		}
		secret, err := s.api.RevealTriggerSet(ctx, id)
		if err != nil {
			return nil, err
		}
		for _, member := range secret.Members {
			s.redactor.Register(member.Key)
		}
		for _, member := range secret.Members {
			if !strings.Contains(string(output), member.Key) {
				return nil, errors.New("human Trigger Set output omitted a current member key")
			}
		}
		return jsonObject(secret)
	}
	if path == "mcp http enable" || path == "mcp http rotate" {
		token := humanLine(string(output), "credential (shown once):")
		if token == "" {
			return nil, errors.New("human localhost MCP output omitted its one-time credential")
		}
		s.redactor.Register(token)
		body["credential"] = token
	}
	return body, nil
}
