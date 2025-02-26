package dal

import (
	"context"
	"database/sql"
)

type Database interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Ping() error
}

type DAL[T any] interface {
	GetByID(ctx context.Context, id string) (*T, error)
	GetAll(ctx context.Context) ([]T, error)
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id string) error
}
