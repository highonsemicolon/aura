package dal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/highonsemicolon/aura/config"
	"github.com/highonsemicolon/aura/src/utils"
)

type MySQLDAL struct {
	conn *sql.DB
}

func NewMySQLDAL(config config.MySQL) *MySQLDAL {
	if config.CAPath != "" {
		tlsConfig, err := utils.InitTLS(config.CAPath)
		if err != nil {
			panic(fmt.Errorf("TLS setup failed for CA Cert '%s': %w", config.CAPath, err))
		}

		mysql.RegisterTLSConfig("custom", tlsConfig)
	}

	Conn, err := sql.Open("mysql", config.DSN)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	Conn.SetMaxOpenConns(config.MaxOpenConns)
	Conn.SetMaxIdleConns(config.MaxIdleConns)
	Conn.SetConnMaxLifetime(config.ConnMaxLifetime)

	if err := Conn.Ping(); err != nil {
		log.Fatal("Database is unreachable:", err)
	}

	return &MySQLDAL{conn: Conn}
}

func (m *MySQLDAL) Exec(query string, args ...any) (sql.Result, error) {
	return m.conn.Exec(query, args...)
}

func (m *MySQLDAL) Query(query string, args ...any) (*sql.Rows, error) {
	return m.conn.Query(query, args...)
}

func (m *MySQLDAL) QueryRow(query string, args ...any) *sql.Row {
	return m.conn.QueryRow(query, args...)
}

func (m *MySQLDAL) PingContext(ctx context.Context) error {
	if m.conn == nil {
		return errors.New("database connection is nil")
	}
	return m.conn.PingContext(ctx)
}

func (m *MySQLDAL) Close() error {
	return m.conn.Close()
}

func (m *MySQLDAL) withTransaction(fn func(tx *sql.Tx) error) error {
	tx, err := m.conn.Begin()
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
