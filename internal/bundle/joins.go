package bundle

import (
	"strings"
)

// Join is reusable portable parallel intent; native work and instances are excluded.
type Join struct {
	PortableID   string       `json:"portable_id"`
	Name         string       `json:"name"`
	EntryTaskID  string       `json:"entry_task_portable_id"`
	Branches     []JoinBranch `json:"branches"`
	TargetTaskID string       `json:"target_task_portable_id"`
}

// JoinBranch keeps a distinct named slot even when several slots reference one task.
type JoinBranch struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	TaskID string `json:"task_portable_id"`
}

func validateJoins(doc Document, tasks map[string]bool) []Issue {
	var issues []Issue
	if len(doc.Joins) > 0 && doc.Schema != SchemaV3 {
		issues = append(issues, Issue{Kind: "schema", Message: "parallel joins require bundle v3"})
	}
	for _, join := range doc.Joins {
		add := func(message string) {
			issues = append(issues, Issue{Kind: "join", Identity: join.PortableID, Message: message})
		}
		if !tasks[join.EntryTaskID] || !tasks[join.TargetTaskID] {
			add("entry or target task is absent from bundle")
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
		cyclic := false
		for _, branch := range join.Branches {
			graph[join.EntryTaskID] = append(graph[join.EntryTaskID], branch.TaskID)
			cyclic = cyclic || reaches(graph, branch.TaskID, join.EntryTaskID, map[string]bool{})
			graph[branch.TaskID] = append(graph[branch.TaskID], join.TargetTaskID)
			cyclic = cyclic || reaches(graph, join.TargetTaskID, branch.TaskID, map[string]bool{})
		}
		if cyclic {
			issues = append(issues, Issue{Kind: "join", Identity: join.PortableID, Message: "parallel relationship would create a cycle with completion chains or joins"})
		}
	}
	return issues
}
