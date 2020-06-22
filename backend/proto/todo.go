package proto

import (
	"github.com/oinume/todomvc/backend/model"
	"github.com/oinume/todomvc/proto-gen/go/proto/todomvc"
)

type TodoConverter struct{}

func NewTodoConverter() *TodoConverter {
	return &TodoConverter{}
}

func (c *TodoConverter) Convert(m *model.Todo) *todomvc.Todo {
	completed := false
	if m.Completed == 1 {
		completed = true
	}
	return &todomvc.Todo{
		Id:        m.ID,
		Title:     m.Title,
		Completed: completed,
	}
}
