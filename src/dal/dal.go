package dal

import "database/sql"

type Database interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Ping() error
}

type DAL[T any] interface {
	GetByID(id string) (*T, error)
	GetAll() ([]T, error)
	Create(entity *T) error
	Update(entity *T) error
	Delete(id string) error
}
