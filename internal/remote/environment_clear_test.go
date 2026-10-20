package remote

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/enrollment"
	"github.com/shruggietech/go-schedule/internal/store"
)

func TestRemoteEnvironmentClearRetainsManageAuthority(t *testing.T) {
	for _, capability := range []domain.Capability{domain.CapabilityObserve, domain.CapabilityOperate, domain.CapabilityManage} {
		t.Run(string(capability), func(t *testing.T) {
			st, err := store.Open(":memory:")
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if err := st.Close(); err != nil {
					t.Error(err)
				}
			})
			local, err := st.LocalActor()
			if err != nil {
				t.Fatal(err)
			}
			task := domain.Task{Name: "fixture", Command: "echo", Env: map[string]string{"EXAMPLE": "old"}, State: domain.TaskActive}
			if err := st.CreateTask(&task); err != nil {
				t.Fatal(err)
			}
			enroll := enrollment.New(st)
			pairing, err := enroll.Create("S107 fixture", domain.ActorKindJSON, capability)
			if err != nil {
				t.Fatal(err)
			}
			credential, err := enroll.Exchange(pairing.ID, pairing.Phrase, pairing.DaemonID, pairing.DisplayName, pairing.Kind, pairing.Capability)
			if err != nil {
				t.Fatal(err)
			}
			api := server.New(st, nil, nil, nil, "", slog.New(slog.NewTextHandler(io.Discard, nil)))
			api.SetActorResolver(ActorID)
			handler := NewHandler(api.Handler(), enroll, local.ID)
			request := httptest.NewRequest(http.MethodPatch, "/api/v1/tasks/"+task.ID, bytes.NewBufferString(`{"clear_env":true}`))
			request.Header.Set("Authorization", "Bearer "+credential.Token)
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			want := http.StatusForbidden
			if capability == domain.CapabilityManage {
				want = http.StatusOK
			}
			if response.Code != want {
				t.Fatalf("status=%d want=%d body=%s", response.Code, want, response.Body.String())
			}
			got, err := st.GetTask(task.ID)
			if err != nil {
				t.Fatal(err)
			}
			if capability == domain.CapabilityManage && len(got.Env) != 0 || capability != domain.CapabilityManage && got.Env["EXAMPLE"] != "old" {
				t.Fatalf("authority=%s environment=%v", capability, got.Env)
			}
		})
	}
}
