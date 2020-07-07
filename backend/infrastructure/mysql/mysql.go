package mysql

import (
	"database/sql"
)

// NewDB connects to database with given `dsn` then returns `*sql.DB` instance.
func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	// TODO: db.SetMaxIdleConns, etc...
	if err != nil {
		return nil, err
	}
	return db, nil
}
