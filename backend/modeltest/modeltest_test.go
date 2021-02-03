package modeltest

import (
	"reflect"
	"testing"
	"time"

	"github.com/oinume/todomvc/backend/model"
)

func TestNewTodo(t *testing.T) {
	createdAt := time.Date(2021, 1, 1, 0, 1, 30, 0, time.UTC)
	tests := map[string]struct {
		setter func(*model.Todo)
		want   *model.Todo
	}{
		"normal": {
			setter: func(todo *model.Todo) {
				todo.ID = "aaa"
				todo.Title = "aaa"
				todo.Completed = 1
				todo.CreatedAt = createdAt
				todo.UpdatedAt = createdAt
			},
			want: &model.Todo{
				ID:        "aaa",
				Completed: 1,
				Title:     "aaa",
				CreatedAt: createdAt,
				UpdatedAt: createdAt,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := NewTodo(tt.setter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTodo() = %v, want %v", got, tt.want)
			}
		})
	}
}
