package mysql

import (
	"database/sql"

	"github.com/oinume/todomvc/backend/repository"
)

type mysqlTodoRepository struct {
	db *sql.DB
	repository.TodoRepository
}

func NewTodoRepository(db *sql.DB) repository.TodoRepository {
	return mysqlTodoRepository{
		db: db,
	}
}
