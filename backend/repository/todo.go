package repository

import (
	"context"

	"github.com/oinume/todomvc/backend/model"
)

type TodoRepository interface {
	// TODO: documentation
	Create(ctx context.Context, executor ContextExecutor, todo *model.Todo) error
	Update(ctx context.Context, todo *model.Todo) error
	Delete(ctx context.Context, todo *model.Todo) (int64, error)
	DeleteByID(ctx context.Context, id string) (int64, error)
	Find(ctx context.Context) ([]*model.Todo, error)
	FindOne(ctx context.Context, id string) (*model.Todo, error)
}
