package bundle

import (
	"sort"
	"strings"
)

// Join is reusable portable parallel intent; native work and instances are excluded.
type Join struct {
	PortableID     string       `json:"portable_id"`
	Name           string       `json:"name"`
	EntryTaskID    string       `json:"entry_task_portable_id"`
	Branches       []JoinBranch `json:"branches"`
	TargetTaskID   string       `json:"target_task_portable_id"`
	RecoveryTaskID string       `json:"recovery_task_portable_id,omitempty"`
	RecoveryOn     *[]string    `json:"recovery_on,omitempty"`
	CleanupTaskID  string       `json:"cleanup_task_portable_id,omitempty"`
}

// JoinBranch keeps a distinct named slot even when several slots reference one task.
type JoinBranch struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	TaskID string `json:"task_portable_id"`
}

func validateJoins(doc Document, tasks map[string]bool) []Issue {
	var issues []Issue
	if len(doc.Joins) > 0 && doc.Schema != SchemaV3 && doc.Schema != SchemaV4 {
		issues = append(issues, Issue{Kind: "schema", Message: "parallel joins require bundle v3 or v4"})
	}
	for _, join := range doc.Joins {
		add := func(message string) {
			issues = append(issues, Issue{Kind: "join", Identity: join.PortableID, Message: message})
		}
		if !tasks[join.EntryTaskID] || !tasks[join.TargetTaskID] {
			add("entry or target task is absent from bundle")
		}
		if joinUsesRecovery(join) && doc.Schema != SchemaV4 {
			add("recovery and cleanup require bundle v4")
		}
		if join.RecoveryTaskID != "" && !tasks[join.RecoveryTaskID] {
			add("recovery task is absent from bundle")
		}
		if join.CleanupTaskID != "" && !tasks[join.CleanupTaskID] {
			add("cleanup task is absent from bundle")
		}
		if join.RecoveryOn != nil {
			if join.RecoveryTaskID == "" && len(*join.RecoveryOn) > 0 {
				add("recovery conditions require a recovery task")
			}
			for _, cause := range *join.RecoveryOn {
				switch cause {
				case "failure", "canceled", "skipped", "not_started":
				default:
					add("unknown recovery condition " + cause)
				}
			}
		}
		if len(join.Branches) == 0 {
			add("at least one named required branch is needed")
		}
		ids, names := map[string]bool{}, map[string]bool{}
		for _, branch := range join.Branches {
			if branch.ID == "" || ids[branch.ID] {
				add("branch IDs must be nonempty and unique")
			}
			ids[branch.ID] = true
			if strings.TrimSpace(branch.Name) == "" || names[branch.Name] {
				add("branch names must be nonempty and unique")
			}
			names[branch.Name] = true
			if !tasks[branch.TaskID] {
				add("branch task is absent from bundle")
			}
		}
	}
	return issues
}

func validateAutomationGraph(chains []Chain, joins []Join) []Issue {
	issues := validateChainGraph(chains)
	graph := map[string][]string{}
	for _, chain := range chains {
		graph[chain.SourceTaskID] = append(graph[chain.SourceTaskID], chain.TargetTaskID)
	}
	for _, join := range joins {
		destinations := []string{join.TargetTaskID, join.RecoveryTaskID, join.CleanupTaskID}
		for _, branch := range join.Branches {
			destinations = append(destinations, branch.TaskID)
		}
		for _, destination := range destinations {
			if destination != "" {
				graph[join.EntryTaskID] = append(graph[join.EntryTaskID], destination)
			}
		}
	}
	for _, join := range joins {
		for _, destination := range graph[join.EntryTaskID] {
			if reaches(graph, destination, join.EntryTaskID, map[string]bool{}) {
				issues = append(issues, Issue{Kind: "join", Identity: join.PortableID, Message: "parallel relationship would create a cycle with completion chains or joins"})
				break
			}
		}
	}
	return issues
}

func joinUsesRecovery(join Join) bool {
	return join.RecoveryTaskID != "" || join.CleanupTaskID != "" || join.RecoveryOn != nil
}

func normalizeJoinRecovery(join *Join) {
	if join.RecoveryTaskID != "" && join.RecoveryOn == nil {
		values := []string{"failure"}
		join.RecoveryOn = &values
	}
	if join.RecoveryOn != nil {
		values := append([]string{}, (*join.RecoveryOn)...)
		sort.Strings(values)
		unique := values[:0]
		for _, value := range values {
			if len(unique) == 0 || unique[len(unique)-1] != value {
				unique = append(unique, value)
			}
		}
		join.RecoveryOn = &unique
	}
}
