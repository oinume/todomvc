package repository

import (
	"context"

	"github.com/oinume/todomvc/backend/model"
)

// TodoRepository is an interface to define CRUD methods.
type TodoRepository interface {
	// Create creates a new todo.
	Create(ctx context.Context, e Executor, todo *model.Todo) error

	// Update updates database with given `todo`.
	Update(ctx context.Context, e Executor, todo *model.Todo) error

	// Delete deletes todo from database.
	Delete(ctx context.Context, e Executor, todo *model.Todo) (int64, error)

	// DeleteByID deletes todo from database with `id`
	DeleteByID(ctx context.Context, e Executor, id string) (int64, error)

	// Find returns slice of `*model.Todo`
	Find(ctx context.Context, e Executor) ([]*model.Todo, error)

	// FindOne returns single `*model.Todo`
	FindOne(ctx context.Context, e Executor, id string) (*model.Todo, error)
}
