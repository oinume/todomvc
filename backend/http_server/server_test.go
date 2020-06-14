package http_server_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/protobuf/jsonpb"
	"github.com/oinume/todomvc/backend/http_server"

	"github.com/oinume/todomvc-example/proto-gen/go/proto/todomvc"
)

func Test_Server_CreateTodo(t *testing.T) {
	m := &jsonpb.Marshaler{OrigName: true}
	var reqBody bytes.Buffer
	if err := m.Marshal(&reqBody, &todomvc.CreateTodoRequest{
		Title: "New task",
	}); err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/todos", &reqBody)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	defer func() { _ = rr.Result().Body.Close() }()

	s := http_server.New()
	s.CreateTodo(rr, req)
	result := rr.Result()
	if result.StatusCode != 201 {
		t.Errorf("unexpected status code: got=%v, want=%v", result.StatusCode, 201)
	}
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = result.Body.Close() }()

	if len(body) == 0 {
		t.Error("empty response reqBody")
	}
	// TODO: Unmarshal to TodoItem
}
