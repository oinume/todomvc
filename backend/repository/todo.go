package repository

import (
	"context"

	"github.com/oinume/todomvc/backend/model"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *model.Todo) error
	Update(ctx context.Context, todo *model.Todo) error
	Delete(ctx context.Context, id string) (int64, error)
	Find(ctx context.Context) ([]*model.Todo, error)
	FindOne(ctx context.Context, id string) (*model.Todo, error)
}
