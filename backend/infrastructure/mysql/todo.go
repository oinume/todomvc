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

func (r *todoRepository) Create(ctx context.Context, executor repository.ContextExecutor, todo *model.Todo) error {
	return todo.Insert(ctx, executor, boil.Infer())
}

func (r *todoRepository) Update(ctx context.Context, todo *model.Todo) error {
	_, err := todo.Update(ctx, r.db, boil.Infer())
	return err
}

func (r *todoRepository) FindOne(ctx context.Context, id string) (*model.Todo, error) {
	return model.FindTodo(ctx, r.db, id)
}

func (r *todoRepository) Delete(ctx context.Context, todo *model.Todo) (int64, error) {
	return todo.Delete(ctx, r.db)
}

func (r *todoRepository) DeleteByID(ctx context.Context, id string) (int64, error) {
	return r.Delete(ctx, &model.Todo{ID: id})
}
