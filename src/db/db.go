package db

import (
	"database/sql"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type DB interface {
	AssignRole(userID, role, resourceID string) error
	RemoveRole(userID, resourceID string) error
	GetRole(userID, resourceID string) (string, error)
	Close() error
}

type sqlDB struct {
	conn      *sql.DB
	validator *validator.Validate
}

func NewSqlDB(conn *sql.DB) *sqlDB {
	if err := conn.Ping(); err != nil {
		panic(err)
	}
	return &sqlDB{conn: conn, validator: validator.New()}
}

func (db *sqlDB) Close() error {
	return db.conn.Close()
}

func (db *sqlDB) validateInput(input interface{}) error {
	if err := db.validator.Struct(input); err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}
	return nil
}
