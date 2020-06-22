package proto

import (
	"github.com/oinume/todomvc/backend/model"
	"github.com/oinume/todomvc/proto-gen/go/proto/todomvc"
)

type TodoConverter struct{}

func NewTodoConverter() *TodoConverter {
	return &TodoConverter{}
}

func (c *TodoConverter) ToProto(todo *model.Todo) *todomvc.Todo {
	completed := false
	if todo.Completed == 1 {
		completed = true
	}
	return &todomvc.Todo{
		Id:        todo.ID,
		Title:     todo.Title,
		Completed: completed,
	}
}

func (s *TodoConverter) ToModel(todo *todomvc.Todo) *model.Todo {
	completed := 0
	if todo.Completed {
		completed = 1
	}
	return &model.Todo{
		ID:        todo.Id,
		Title:     todo.Title,
		Completed: uint8(completed),
	}
}
