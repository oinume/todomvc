package http

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/oinume/todomvc/backend/model"
	"github.com/oinume/todomvc/backend/proto"
	"github.com/oinume/todomvc/proto-gen/go/proto/todomvc"
)

type CreateTodoRequestValidation struct {
	Title string `validate:"required,min=1,max=50"`
}

type UpdateTodoRequestValidation struct {
	Id    string `validate:"required"`
	Title string `validate:"required,min=1,max=50"`
}

func (s *server) CreateTodo(w http.ResponseWriter, r *http.Request) {
	req := &todomvc.CreateTodoRequest{}
	if err := s.unmarshaler.Unmarshal(r.Body, req); err != nil {
		internalServerError(s.logger, w, err)
		return
	}

	ctx := r.Context()
	validation := CreateTodoRequestValidation{
		Title: req.Title,
	}
	if err := s.validator.StructCtx(ctx, validation); err != nil {
		validationError(w, err)
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

	writeJSON(w, http.StatusCreated, proto.NewTodoConverter().ToProto(todo))
}

func (s *server) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	req := &todomvc.UpdateTodoRequest{}
	if err := s.unmarshaler.Unmarshal(r.Body, req); err != nil {
		internalServerError(s.logger, w, err)
		return
	}

	ctx := r.Context()
	validation := UpdateTodoRequestValidation{
		Id:    req.Todo.Id,
		Title: req.Todo.Title,
	}
	if err := s.validator.StructCtx(ctx, validation); err != nil {
		validationError(w, err)
		return
	}

	converter := proto.NewTodoConverter()
	todo := converter.ToModel(req.Todo)
	if _, err := s.todoRepo.FindOne(ctx, todo.ID); err != nil {
		// TODO: not found error body
		writeJSON(w, http.StatusNotFound, struct{}{})
		return
	}
	if err := s.todoRepo.Update(ctx, todo); err != nil {
		internalServerError(s.logger, w, err)
		return
	}

	writeJSON(w, http.StatusOK, converter.ToProto(todo))
}
