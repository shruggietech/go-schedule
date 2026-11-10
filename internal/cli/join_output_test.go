package cli

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
)

func TestJoinHumanOutputPropagatesWriterFailure(t *testing.T) {
	for _, args := range [][]string{
		{"join", "create", "--name", "X", "--entry", "entry", "--target", "target", "--branches-file", "-"},
		{"join", "update", "join", "--name", "Renamed"},
		{"join", "list"}, {"join", "show", "join"}, {"join", "instances"}, {"join", "instance", "instance"},
	} {
		t.Run(strings.Join(args[:2], "-"), func(t *testing.T) {
			calls := 0
			c := client.NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				w.Header().Set("Content-Type", "application/json")
				switch r.URL.Path {
				case "/v1/joins":
					if r.Method == http.MethodGet {
						io.WriteString(w, `{"joins":[{"id":"join","branches":[]}]}`)
					} else {
						w.WriteHeader(http.StatusCreated)
						io.WriteString(w, `{"id":"join","branches":[]}`)
					}
				case "/v1/join-instances":
					io.WriteString(w, `{"instances":[{"id":"instance","branches":[]}]}`)
				default:
					io.WriteString(w, `{"id":"join","branches":[]}`)
				}
			}))
			var diagnostic bytes.Buffer
			root := newRootWithDependencies(Dependencies{Client: c, Out: failingOutput{}, Err: &diagnostic, In: strings.NewReader(`[{"id":"a","name":"A","task_id":"branch"}]`)})
			root.SetArgs(args)
			err := root.Execute()
			if !errors.Is(err, io.ErrClosedPipe) || calls != 1 {
				t.Fatalf("writer failure swallowed or operation repeated: err=%v calls=%d", err, calls)
			}
		})
	}
}
