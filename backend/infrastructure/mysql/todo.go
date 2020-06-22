package mysql

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/boil"

	"github.com/oinume/todomvc/backend/model"
	"github.com/oinume/todomvc/backend/repository"
)

type todoRepository struct {
	db *sql.DB
	repository.TodoRepository
}

func NewTodoRepository(db *sql.DB) repository.TodoRepository {
	return &todoRepository{
		db: db,
	}
}

func (r *todoRepository) Create(ctx context.Context, todo *model.Todo) error {
	return todo.Insert(ctx, r.db, boil.Infer())
}

func (r *todoRepository) Update(ctx context.Context, todo *model.Todo) error {
	_, err := todo.Update(ctx, r.db, boil.Infer())
	return err
}
