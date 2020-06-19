package http

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"github.com/golang/protobuf/jsonpb"

	"github.com/oinume/todomvc/proto-gen/go/proto/todomvc"
)

func Test_Server_CreateTodo(t *testing.T) {
	m := &jsonpb.Marshaler{OrigName: true}
	u := &jsonpb.Unmarshaler{}
	s := NewServer(todoRepo, zap.NewNop())

	type response struct {
		statusCode int
		todoItem   *todomvc.TodoItem
	}
	tests := map[string]struct {
		request      *todomvc.CreateTodoRequest
		wantResponse response
	}{
		"OK_Created": {
			request: &todomvc.CreateTodoRequest{
				Title: "NewServer task",
			},
			wantResponse: response{
				statusCode: http.StatusCreated,
				todoItem: &todomvc.TodoItem{
					Title:     "NewServer task",
					Completed: false,
				},
			},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var reqBody bytes.Buffer
			if err := m.Marshal(&reqBody, test.request); err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest("POST", "/todos", &reqBody)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			defer func() { _ = rr.Result().Body.Close() }()

			s.CreateTodo(rr, req)

			result := rr.Result()
			if result.StatusCode != http.StatusCreated {
				body, _ := ioutil.ReadAll(result.Body)
				t.Fatalf("unexpected status code: got=%v, want=%v: body=%v", result.StatusCode, http.StatusCreated, string(body))
			}
			got := &todomvc.TodoItem{}
			if err := u.Unmarshal(result.Body, got); err != nil {
				t.Fatal(err)
			}
			if got.Id == "" {
				t.Fatal("got.ID is empty")
			}
			if got, want := got.Title, test.request.Title; got != want {
				// TODO: Use go-cmp
				t.Fatalf("unexpected Title: got=%v, want=%v", got, want)
			}
		})
	}
}

// TODO: Add error cases
