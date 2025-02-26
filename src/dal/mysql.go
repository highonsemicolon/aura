package dal

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/highonsemicolon/aura/config"
	"github.com/highonsemicolon/aura/src/utils"
)

type MySQLDAL struct {
	db *sql.DB
}

func NewMySQLDAL(config config.MySQL) *MySQLDAL {
	if config.CAPath != "" {
		tlsConfig, err := utils.InitTLS(config.CAPath)
		if err != nil {
			panic(fmt.Errorf("TLS setup failed for CA Cert '%s': %w", config.CAPath, err))
		}

		mysql.RegisterTLSConfig("custom", tlsConfig)
	}

	db, err := sql.Open("mysql", config.DSN)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	if err := db.Ping(); err != nil {
		log.Fatal("Database is unreachable:", err)
	}

	return &MySQLDAL{db: db}
}

func (m *MySQLDAL) Exec(query string, args ...interface{}) (sql.Result, error) {
	return m.db.Exec(query, args...)
}

func (m *MySQLDAL) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return m.db.Query(query, args...)
}

func (m *MySQLDAL) QueryRow(query string, args ...interface{}) *sql.Row {
	return m.db.QueryRow(query, args...)
}

func (m *MySQLDAL) Ping() error {
	return m.db.Ping()
}

func (m *MySQLDAL) Close() error {
	return m.db.Close()
}

func (m *MySQLDAL) withTransaction(fn func(tx *sql.Tx) error) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
