package modeltest

import (
	mrand "math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/oinume/todomvc/backend/model"
)

var (
	letters = []rune(`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`)
	//commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	random = mrand.New(mrand.NewSource(time.Now().UnixNano()))
)

func NewTodo(setters ...func(*model.Todo)) *model.Todo {
	todo := &model.Todo{}
	for _, setter := range setters {
		setter(todo)
	}
	if todo.ID == "" {
		todo.ID = uuid.New().String()
	}
	if todo.Title == "" {
		todo.Title = RandomString(10)
	}
	return todo
}

func RandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[random.Intn(len(letters))]
	}
	return string(b)
}
