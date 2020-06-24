package http

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/code"

	"github.com/oinume/todomvc/backend/model"
	"github.com/oinume/todomvc/backend/modeltest"
	"github.com/oinume/todomvc/backend/testings"
	"github.com/oinume/todomvc/proto-gen/go/proto/todomvc"
)

func Test_Server_CreateTodo(t *testing.T) {
	type response struct {
		statusCode int
		todo       *todomvc.Todo
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
				todo: &todomvc.Todo{
					Title:     "NewServer task",
					Completed: false,
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

			req := newHTTPRequest(t, m, "POST", "/todos", tt.request)
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
			testings.RequireEqual(t, tt.request.Title, got.Title, "unexpected Title")
		})
	}
}

func Test_Server_CreateTodo_Error(t *testing.T) {
	type response struct {
		statusCode int
		error      *todomvc.Error
	}
	tests := map[string]struct {
		request      *todomvc.CreateTodoRequest
		wantResponse response
	}{
		"BadRequest_TitleIsEmpty": {
			request: &todomvc.CreateTodoRequest{
				Title: "",
			},
			wantResponse: response{
				statusCode: http.StatusBadRequest,
				error: &todomvc.Error{
					Code:    code.Code_INVALID_ARGUMENT,
					Message: "Validation error",
				},
			},
		},
		"BadRequest_TitleIsTooLong": {
			request: &todomvc.CreateTodoRequest{
				Title: "012345678901234567890123456789012345678901234567890", // 51 char
			},
			wantResponse: response{
				statusCode: http.StatusBadRequest,
				error: &todomvc.Error{
					Code:    code.Code_INVALID_ARGUMENT,
					Message: "Validation error",
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

			req := newHTTPRequest(t, m, "POST", "/todos", tt.request)
			rr := httptest.NewRecorder()
			defer func() { _ = rr.Result().Body.Close() }()

			s.CreateTodo(rr, req)

			result := rr.Result()
			if got, want := result.StatusCode, tt.wantResponse.statusCode; got != want {
				body, _ := ioutil.ReadAll(result.Body)
				t.Fatalf("unexpected status code: got=%v, want=%v: body=%v", got, want, string(body))
			}

			gotResponse := &todomvc.Error{}
			if err := u.Unmarshal(result.Body, gotResponse); err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}
			testings.RequireEqual(t, tt.wantResponse.error, gotResponse, "CreateTodo unexpected response")
		})
	}
}

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

			req := newHTTPRequest(t, m, "PATCH", "/todos", tt.request)
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
			testings.RequireEqual(t, tt.wantResponse.todo, got, "UpdateTodo unexected response")
		})
	}
}

func Test_server_UpdateTodo_Error(t *testing.T) {
	const title = "New frontend task"
	type response struct {
		statusCode int
		error      *todomvc.Error
	}
	tests := map[string]struct {
		request      *todomvc.UpdateTodoRequest
		wantResponse response
	}{
		"NotFound": {
			request: &todomvc.UpdateTodoRequest{
				Todo: &todomvc.Todo{
					Id:        "not_found",
					Title:     title,
					Completed: true,
				},
			},
			wantResponse: response{
				statusCode: http.StatusNotFound,
				error:      &todomvc.Error{},
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

			req := newHTTPRequest(t, m, "PATCH", "/todos", tt.request)
			rr := httptest.NewRecorder()
			defer func() { _ = rr.Result().Body.Close() }()

			s.UpdateTodo(rr, req)

			result := rr.Result()
			if got, want := result.StatusCode, tt.wantResponse.statusCode; got != want {
				body, _ := ioutil.ReadAll(result.Body)
				t.Fatalf("unexpected status code: got=%v, want=%v: body=%v", got, want, string(body))
			}

			gotError := &todomvc.Error{}
			if err := u.Unmarshal(result.Body, gotError); err != nil {
				testings.RequireEqual(t, tt.wantResponse.error, gotError, "UpdateTodo unexpected response")
			}
		})
	}
}

func newHTTPRequest(t *testing.T, m *jsonpb.Marshaler, method, path string, request proto.Message) *http.Request {
	t.Helper()
	var reqBody bytes.Buffer
	if err := m.Marshal(&reqBody, request); err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	req, err := http.NewRequest(method, path, &reqBody)
	if err != nil {
		t.Fatalf("http.NewRequest failed: %v", err)
	}
	return req
}
