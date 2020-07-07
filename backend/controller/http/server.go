package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"go.opencensus.io/plugin/ochttp"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/code"

	"github.com/oinume/todomvc/backend/repository"
	"github.com/oinume/todomvc/proto-gen/go/proto/todomvc"
)

type server struct {
	httpServer  *http.Server
	db          *sql.DB
	todoRepo    repository.TodoRepository
	logger      *zap.Logger
	unmarshaler *jsonpb.Unmarshaler
	validator   *validator.Validate
}

// NewServer creates HTTP server for backend
func NewServer(addr string, db *sql.DB, todoRepo repository.TodoRepository, logger *zap.Logger) *server {
	s := &server{
		db:       db,
		todoRepo: todoRepo,
		logger:   logger,
		unmarshaler: &jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		},
		validator: validator.New(),
	}
	router := s.newRouter()
	ochttpHandler := &ochttp.Handler{
		Handler: router,
	}
	httpServer := &http.Server{
		Addr:    addr,
		Handler: ochttpHandler,
	}
	s.httpServer = httpServer
	return s
}

// ListenAndServe is a wrapper of `http.ListenAndServe`.
func (s *server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown this server gracefully.
func (s *server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *server) newRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(accessLogMiddleware(s.logger))
	r.Handle("/todos", ochttp.WithRouteTag(http.HandlerFunc(s.ListTodos), "/todos")).Methods("GET")
	r.Handle("/todos/{id}", ochttp.WithRouteTag(http.HandlerFunc(s.GetTodo), "/todos")).Methods("GET")
	r.Handle("/todos", ochttp.WithRouteTag(http.HandlerFunc(s.CreateTodo), "/todos")).Methods("POST")
	r.Handle("/todos/{id}", ochttp.WithRouteTag(http.HandlerFunc(s.UpdateTodo), "/todos/{id}")).Methods("PATCH")
	r.Handle("/todos/{id}", ochttp.WithRouteTag(http.HandlerFunc(s.DeleteTodo), "/todos/{id}")).Methods("DELETE")
	return r
}

func validationError(w http.ResponseWriter, err error) {
	_, ok := err.(validator.ValidationErrors)
	if !ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	r := &todomvc.Error{
		Code:    code.Code_INVALID_ARGUMENT,
		Message: "Validation error",
	}
	// TODO: Map errors to ErrorResponse
	//for _, e := range errors {
	//	e.Field()
	//}

	writeJSONProto(w, http.StatusBadRequest, r)
}

func internalServerError(logger *zap.Logger, w http.ResponseWriter, err error) {
	logger.Error("caught error", zap.Error(err))
	http.Error(w, fmt.Sprintf("Internal server Error\n\n%v", err), http.StatusInternalServerError)
}

func writeJSON(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, `{ "status": "Failed to Encode as writeJSON" }`, http.StatusInternalServerError)
	}
}

func writeJSONProto(w http.ResponseWriter, code int, message proto.Message) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	m := &jsonpb.Marshaler{
		EmitDefaults: false,
		OrigName:     true,
	}
	if err := m.Marshal(w, message); err != nil {
		http.Error(w, `{ "status": "Failed to Encode as writeJSON" }`, http.StatusInternalServerError)
	}
}

func writeNoContent(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)
}
