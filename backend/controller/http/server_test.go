package http

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/oinume/todomvc/backend/infrastructure/mysql"
	"github.com/oinume/todomvc/backend/repository"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/jsonpb"
	"go.uber.org/zap"

	"github.com/oinume/todomvc/backend/config"
	"github.com/oinume/todomvc/proto-gen/go/proto/todomvc"
)

var (
	db       *sql.DB
	todoRepo repository.TodoRepository
)

func TestMain(m *testing.M) {
	_ = os.Setenv("MYSQL_DATABASE", "todomvc_test")
	config.MustProcessDefault()
	ldb, err := mysql.NewDB(config.DefaultVars.DBURL())
	if err != nil {
		panic("Failed to mysql.NewDB: " + err.Error())
	}
	db = ldb
	todoRepo = mysql.NewTodoRepository(db)
	os.Exit(m.Run())
}

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
				t.Fatalf("unexpected status code: got=%v, want=%v", result.StatusCode, http.StatusCreated)
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
