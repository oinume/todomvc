package http

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/oinume/todomvc/backend/model"
	"github.com/oinume/todomvc/backend/proto"
	"github.com/oinume/todomvc/proto-gen/go/proto/todomvc"
)

func (s *server) CreateTodo(w http.ResponseWriter, r *http.Request) {
	req := &todomvc.CreateTodoRequest{}
	if err := s.unmarshaler.Unmarshal(r.Body, req); err != nil {
		internalServerError(s.logger, w, err)
		return
	}

	id := uuid.New().String()
	todo := &model.Todo{
		ID:    id,
		Title: req.Title,
	}
	if err := s.todoRepo.Create(r.Context(), todo); err != nil {
		internalServerError(s.logger, w, err)
		return
	}

	writeJSON(w, http.StatusCreated, proto.NewTodoConverter().Convert(todo))
}

func (s *server) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	req := &todomvc.UpdateTodoRequest{}
	if err := s.unmarshaler.Unmarshal(r.Body, req); err != nil {
		internalServerError(s.logger, w, err)
		return
	}

	id := uuid.New().String()
	todo := &model.Todo{
		ID:    id,
		Title: req.Title,
	}
	if err := s.todoRepo.Create(r.Context(), todo); err != nil {
		internalServerError(s.logger, w, err)
		return
	}

	writeJSON(w, http.StatusCreated, proto.NewTodoConverter().Convert(todo))
}
