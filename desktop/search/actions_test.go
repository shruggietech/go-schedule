package search

import (
	"context"
	"errors"
	"testing"

	"github.com/shruggietech/go-schedule/desktop/connection"
	"github.com/shruggietech/go-schedule/desktop/systems"
	"github.com/shruggietech/go-schedule/internal/api/client"
	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestExecuteRevalidatesIdentityAndObjectBeforeMutation(t *testing.T) {
	called := false
	daemon := &daemonFake{task: server.TaskResponse{Task: domain.Task{ID: "task", Enabled: true}}, action: func(domain.SearchAction, string) error { called = true; return nil }}
	serviceTarget := target{registration: systems.Registration{Key: "local", Label: "This computer"}, backend: backendFake{health: connection.Health{ID: "changed", Permissions: []string{"operate"}}}, client: daemon}
	service := &Service{resolveTarget: func(string) (target, error) { return serviceTarget, nil }}
	results := service.executeTarget(context.Background(), domain.SearchActionDisable, "local", []Selection{{RegistrationKey: "local", ExpectedDaemonID: "expected", Kind: domain.SearchKindTask, ObjectID: "task", DisplayName: "Task"}})
	if len(results) != 1 || results[0].Outcome != "rejected" || called {
		t.Fatalf("results=%+v called=%v", results, called)
	}
}

func TestExecuteRejectsObserveOnlyAuthorityBeforeMutation(t *testing.T) {
	called := false
	daemon := &daemonFake{task: server.TaskResponse{Task: domain.Task{ID: "task", Enabled: true}}, action: func(domain.SearchAction, string) error { called = true; return nil }}
	value := target{registration: systems.Registration{Key: "remote", Label: "Remote"}, backend: backendFake{health: connection.Health{ID: "daemon", Permissions: []string{"read"}}}, client: daemon}
	service := &Service{resolveTarget: func(string) (target, error) { return value, nil }}
	result := service.Execute(context.Background(), ActionIntent{Action: domain.SearchActionDisable, Selections: []Selection{{RegistrationKey: "remote", ExpectedDaemonID: "daemon", Kind: domain.SearchKindTask, ObjectID: "task"}}})
	if result.Outcome != "rejected" || called {
		t.Fatalf("result=%+v called=%v", result, called)
	}
}

func TestExecutePreservesUncertainMutationEvidence(t *testing.T) {
	daemon := &daemonFake{task: server.TaskResponse{Task: domain.Task{ID: "task", Enabled: true}}, action: func(domain.SearchAction, string) error {
		return &client.MutationUncertainError{Operation: "disable task", Cause: errors.New("transport")}
	}}
	value := target{registration: systems.Registration{Key: "remote", Label: "Remote"}, backend: backendFake{health: connection.Health{ID: "daemon", Permissions: []string{"operate"}}}, client: daemon}
	service := &Service{resolveTarget: func(string) (target, error) { return value, nil }}
	result := service.Execute(context.Background(), ActionIntent{Action: domain.SearchActionDisable, Selections: []Selection{{RegistrationKey: "remote", ExpectedDaemonID: "daemon", Kind: domain.SearchKindTask, ObjectID: "task"}}})
	if result.Outcome != "partial" || len(result.Outcomes) != 1 || result.Outcomes[0].Outcome != "uncertain" {
		t.Fatalf("result=%+v", result)
	}
}

func TestRevalidateAndDispatchRejectsStaleTaskState(t *testing.T) {
	daemon := &daemonFake{task: server.TaskResponse{Task: domain.Task{ID: "task", Enabled: true}}}
	err := revalidateAndDispatch(context.Background(), daemon, domain.SearchActionEnable, Selection{Kind: domain.SearchKindTask, ObjectID: "task"})
	if err == nil {
		t.Fatal("expected already-enabled task to be rejected")
	}
}

func TestRevalidateAndDispatchAcknowledgesOnlyCurrentAlert(t *testing.T) {
	called := false
	daemon := &daemonFake{alerts: []domain.Alert{{ID: "alert"}}, action: func(action domain.SearchAction, id string) error {
		called = action == domain.SearchActionAcknowledge && id == "alert"
		return nil
	}}
	if err := revalidateAndDispatch(context.Background(), daemon, domain.SearchActionAcknowledge, Selection{Kind: domain.SearchKindAlert, ObjectID: "alert"}); err != nil || !called {
		t.Fatalf("err=%v called=%v", err, called)
	}
	if err := revalidateAndDispatch(context.Background(), daemon, domain.SearchActionAcknowledge, Selection{Kind: domain.SearchKindAlert, ObjectID: "missing"}); err == nil {
		t.Fatal("missing alert accepted")
	}
}

func TestRevalidateAndDispatchFindsSelectedAlertBeyondFirstPage(t *testing.T) {
	alerts := make([]domain.Alert, 201)
	for index := range alerts {
		alerts[index] = domain.Alert{ID: "older"}
	}
	alerts[200] = domain.Alert{ID: "selected"}
	called := false
	daemon := &daemonFake{alerts: alerts, action: func(action domain.SearchAction, id string) error {
		called = action == domain.SearchActionAcknowledge && id == "selected"
		return nil
	}}
	if err := revalidateAndDispatch(context.Background(), daemon, domain.SearchActionAcknowledge, Selection{Kind: domain.SearchKindAlert, ObjectID: "selected"}); err != nil || !called {
		t.Fatalf("err=%v called=%v", err, called)
	}
}
