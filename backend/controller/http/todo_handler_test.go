package http

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/protobuf/jsonpb"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.uber.org/zap"

	"github.com/oinume/todomvc/backend/model"
	"github.com/oinume/todomvc/backend/modeltest"
	"github.com/oinume/todomvc/proto-gen/go/proto/todomvc"
)

func Test_Server_CreateTodo(t *testing.T) {
	m := &jsonpb.Marshaler{OrigName: true}
	u := &jsonpb.Unmarshaler{}
	s := NewServer("", todoRepo, zap.NewNop())

	type response struct {
		statusCode int
		todoItem   *todomvc.Todo
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
				todoItem: &todomvc.Todo{
					Title:     "NewServer task",
					Completed: false,
				},
			},
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var reqBody bytes.Buffer
			if err := m.Marshal(&reqBody, tt.request); err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest("POST", "/todos", &reqBody)
			rr := httptest.NewRecorder()
			defer func() { _ = rr.Result().Body.Close() }()

			s.CreateTodo(rr, req)

			result := rr.Result()
			if got, want := result.StatusCode, tt.wantResponse.statusCode; got != want {
				body, _ := ioutil.ReadAll(result.Body)
				t.Fatalf("unexpected status code: got=%v, want=%v: body=%v", got, want, string(body))
			}
			got := &todomvc.Todo{}
			if err := u.Unmarshal(result.Body, got); err != nil {
				t.Fatal(err)
			}
			if got.Id == "" {
				t.Fatal("got.Id is empty")
			}
			if got, want := got.Title, tt.request.Title; got != want {
				t.Fatalf("unexpected Title: got=%v, want=%v", got, want)
			}
		})
	}
}

// TODO: Add error cases

func Test_server_UpdateTodo(t *testing.T) {
	const title = "New frontend task"
	type response struct {
		statusCode int
		todo       *todomvc.Todo
	}
	tests := map[string]struct {
		request      *todomvc.UpdateTodoRequest
		wantResponse response
	}{
		"OK_Updated": {
			request: &todomvc.UpdateTodoRequest{
				Todo: &todomvc.Todo{
					Id:        "", // Set later
					Title:     title,
					Completed: true,
				},
			},
			wantResponse: response{
				statusCode: http.StatusOK,
				todo: &todomvc.Todo{
					Id:        "", // Set later
					Title:     title,
					Completed: true,
				},
			},
		},
	}

	m := &jsonpb.Marshaler{OrigName: true}
	u := &jsonpb.Unmarshaler{}
	s := NewServer("", todoRepo, zap.NewNop())
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			todo := modeltest.NewTodo(func(todo *model.Todo) {
				todo.Title = "New frontend task"
			})
			if err := todoRepo.Create(ctx, todo); err != nil {
				t.Fatal(err)
			}
			tt.request.Todo.Id = todo.ID
			tt.wantResponse.todo.Id = todo.ID

			var reqBody bytes.Buffer
			if err := m.Marshal(&reqBody, tt.request); err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest("PATCH", "/todos", &reqBody)
			rr := httptest.NewRecorder()
			defer func() { _ = rr.Result().Body.Close() }()

			s.UpdateTodo(rr, req)

			result := rr.Result()
			if got, want := result.StatusCode, tt.wantResponse.statusCode; got != want {
				body, _ := ioutil.ReadAll(result.Body)
				t.Fatalf("unexpected status code: got=%v, want=%v: body=%v", got, want, string(body))
			}

			got := &todomvc.Todo{}
			if err := u.Unmarshal(result.Body, got); err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.wantResponse.todo, got, cmpopts.IgnoreUnexported(todomvc.Todo{})); diff != "" {
				t.Fatalf("UpdateTodo response mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
