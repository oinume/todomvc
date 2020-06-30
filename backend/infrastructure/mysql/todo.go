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

func (r *todoRepository) Create(ctx context.Context, e repository.Executor, todo *model.Todo) error {
	return todo.Insert(ctx, e, boil.Infer())
}

func (r *todoRepository) Update(ctx context.Context, e repository.Executor, todo *model.Todo) error {
	_, err := todo.Update(ctx, e, boil.Infer())
	return err
}

func (r *todoRepository) FindOne(ctx context.Context, e repository.Executor, id string) (*model.Todo, error) {
	return model.FindTodo(ctx, e, id)
}

func (r *todoRepository) Find(ctx context.Context, e repository.Executor) ([]*model.Todo, error) {
	return model.Todos().All(ctx, e)
}

func (r *todoRepository) Delete(ctx context.Context, e repository.Executor, todo *model.Todo) (int64, error) {
	return todo.Delete(ctx, e)
}

func (r *todoRepository) DeleteByID(ctx context.Context, e repository.Executor, id string) (int64, error) {
	return r.Delete(ctx, e, &model.Todo{ID: id})
}
