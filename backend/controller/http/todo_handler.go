package http

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/oinume/todomvc/backend/model"
	"github.com/oinume/todomvc/backend/proto"
	"github.com/oinume/todomvc/backend/repository"
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
	if err := s.todoRepo.Create(ctx, s.db, todo); err != nil {
		internalServerError(s.logger, w, err)
		return
	}

	writeJSON(w, http.StatusCreated, proto.NewTodoConverter().ToProto(todo))
}

func (s *server) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		writeJSON(w, http.StatusBadRequest, errors.New("no id"))
		return
	}

	req := &todomvc.UpdateTodoRequest{}
	if err := s.unmarshaler.Unmarshal(r.Body, req); err != nil {
		internalServerError(s.logger, w, err)
		return
	}

	if id != req.Todo.Id {
		writeJSON(w, http.StatusBadRequest, errors.New("wrong id"))
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
	if err := repository.Transaction(ctx, s.db, func(e repository.Executor) error {
		if _, err := s.todoRepo.FindOne(ctx, e, todo.ID); err != nil {
			return err
		}
		return s.todoRepo.Update(ctx, e, todo)
	}); err != nil {
		if err == sql.ErrNoRows {
			// TODO: not found error body
			writeJSON(w, http.StatusNotFound, struct{}{})
			return
		}
		internalServerError(s.logger, w, err)
		return
	}

	writeJSON(w, http.StatusOK, converter.ToProto(todo))
}

func (s *server) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		writeJSON(w, http.StatusBadRequest, errors.New("no id"))
		return
	}

	ctx := r.Context()
	deleted, err := s.todoRepo.DeleteByID(ctx, s.db, id)
	if err != nil {
		internalServerError(s.logger, w, err)
	}
	if deleted == 0 {
		writeJSON(w, http.StatusNotFound, struct{}{})
		return
	}

	writeNoContent(w)
}
