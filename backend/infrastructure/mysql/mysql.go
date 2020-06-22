package mysql

import (
	"database/sql"
)

func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	// TODO: db.SetMaxIdleConns, etc...
	if err != nil {
		return nil, err
	}
	return db, nil
}
