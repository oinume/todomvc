package repository

import (
	"context"

	"github.com/oinume/todomvc/backend/model"
)

type TodoRepository interface {
	// TODO: documentation
	Create(ctx context.Context, e Executor, todo *model.Todo) error
	Update(ctx context.Context, e Executor, todo *model.Todo) error
	Delete(ctx context.Context, e Executor, todo *model.Todo) (int64, error)
	DeleteByID(ctx context.Context, e Executor, id string) (int64, error)
	Find(ctx context.Context, e Executor) ([]*model.Todo, error)
	FindOne(ctx context.Context, e Executor, id string) (*model.Todo, error)
}
